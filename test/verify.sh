#!/usr/bin/env bash
set -euo pipefail

script_dir="$(cd "$(dirname "$0")" && pwd)"
repo_root="$(cd "$script_dir/.." && pwd)"
mode="${1:-all}"
go_cache="${GOCACHE:-${TMPDIR:-/tmp}/spectacular-test-gocache}"
log_file="$(mktemp)"

cleanup() {
  rm -f "$log_file"
}
trap cleanup EXIT

mkdir -p "$go_cache"
export GOCACHE="$go_cache"
export GOPROXY=off
export GOFLAGS=-mod=readonly

check() {
  label="$1"
  shift
  if ! (cd "$repo_root" && "$@") >>"$log_file" 2>&1; then
    echo "Spectacular verification: FAIL check=$label" >&2
    cat "$log_file" >&2
    exit 1
  fi
}

tree_basis() {
  (
    cd "$repo_root"
    git ls-files -co --exclude-standard -- VERSION cmd internal skills install .spectacular test \
      | LC_ALL=C sort -u \
      | while IFS= read -r path; do
          [[ -f "$path" ]] || continue
          printf '%s %s\n' "$(git hash-object "$path")" "$path"
        done \
      | git hash-object --stdin
  )
}

manifest_checks() {
  local version
  version="$(sed -n '1p' "$repo_root/VERSION" | tr -d '[:space:]')"
  if [[ -z "$version" ]]; then
    echo "VERSION is empty" >&2
    exit 1
  fi
  if [[ -f "$repo_root/plugin.json" ]]; then
    if ! grep -q "\"version\": \"$version\"" "$repo_root/plugin.json"; then
      echo "plugin.json version drift: expected $version" >&2
      exit 1
    fi
  fi
  if [[ -f "$repo_root/.claude-plugin/plugin.json" ]]; then
    if ! grep -q "\"version\": \"$version\"" "$repo_root/.claude-plugin/plugin.json"; then
      echo ".claude-plugin/plugin.json version drift: expected $version" >&2
      exit 1
    fi
  fi
  if [[ -f "$repo_root/.codex-plugin/plugin.json" ]]; then
    if ! grep -q "\"version\": \"$version\"" "$repo_root/.codex-plugin/plugin.json"; then
      echo ".codex-plugin/plugin.json version drift: expected $version" >&2
      exit 1
    fi
  fi
  if [[ -f "$repo_root/skills/spectacular/SKILL.md" ]]; then
    if ! grep -q '^  version: "'"$version"'"$' "$repo_root/skills/spectacular/SKILL.md"; then
      echo "skills/spectacular/SKILL.md version drift: expected $version" >&2
      exit 1
    fi
  fi
}

security_checks() {
  if command -v gitleaks >/dev/null 2>&1; then
    check gitleaks gitleaks git --redact -v --no-banner
  fi
}

# --- Pre-Flight (Tier 0 + Tier 1) -------------------------------------------
# Read-only, sub-2s sanity gate. Emits a JSON receipt on stdout.
# Fails fast so heavy tiers (acceptance/release/all) are never spent on a
# workspace that is already syntactically or contractually broken.

preflight_failures=()

preflight_fail() {
  preflight_failures+=("$1")
}

preflight_json_escape() {
  printf '%s' "$1" | python3 -c 'import json,sys; print(json.dumps(sys.stdin.read()))'
}

# Tier 0: static syntax and tree sanity.
preflight_tier0() {
  local formatting
  formatting="$(cd "$repo_root" && gofmt -l cmd internal test/acceptance 2>&1)" || {
    preflight_fail "tier0/gofmt: gofmt failed to run"
    return
  }
  if [[ -n "$formatting" ]]; then
    preflight_fail "tier0/gofmt: unformatted files: $(echo "$formatting" | tr '\n' ' ')"
  fi

  local vet_out
  if ! vet_out="$(cd "$repo_root" && go vet ./... 2>&1)"; then
    preflight_fail "tier0/go-vet: $(echo "$vet_out" | head -5 | tr '\n' ' ')"
  fi

  local whitespace
  if ! whitespace="$(cd "$repo_root" && git diff --check 2>&1)"; then
    preflight_fail "tier0/git-diff-check: $(echo "$whitespace" | head -5 | tr '\n' ' ')"
  fi

  local conflicts
  conflicts="$(cd "$repo_root" && git ls-files -u -- go.sum go.mod 2>/dev/null | awk '{print $4}' | sort -u)"
  if [[ -n "$conflicts" ]]; then
    preflight_fail "tier0/lockfile-conflict: unmerged: $(echo "$conflicts" | tr '\n' ' ')"
  fi

  local bad_frontmatter=""
  while IFS= read -r path; do
    [[ -f "$repo_root/$path" ]] || continue
    if [[ "$(head -1 "$repo_root/$path")" != "---" ]]; then
      bad_frontmatter+="$path "
      continue
    fi
    if ! sed -n '2,200p' "$repo_root/$path" | grep -qx -- '---'; then
      bad_frontmatter+="$path "
    fi
  done < <(cd "$repo_root" && git ls-files -co --exclude-standard -- '.spectacular/**/*.md' 2>/dev/null | grep -v '/index\.md$' | grep -v '/README\.md$' | grep -v '\.amendments\.md$' | grep -vE '\.spectacular/[A-Z]+\.md$')
  if [[ -n "$bad_frontmatter" ]]; then
    preflight_fail "tier0/frontmatter: unterminated or missing frontmatter: $bad_frontmatter"
  fi
}

# Tier 1: contract and schema drift for the live Mission(s).
preflight_tier1() {
  local refs=()
  if [[ -n "${PREFLIGHT_MISSION_REF:-}" ]]; then
    refs=("$PREFLIGHT_MISSION_REF")
  else
    # Default to the live Mission: the highest-numbered one in missions/.
    # PREFLIGHT_MISSION_REF pins a specific ref; PREFLIGHT_ALL_MISSIONS=1 sweeps all.
    local candidates
    candidates="$(find "$repo_root/.spectacular/missions" -mindepth 1 -maxdepth 1 -type d 2>/dev/null \
      | while IFS= read -r dir; do basename "$dir" | cut -d- -f1; done \
      | sed -n 's/^M\([0-9][0-9]*\)$/\1/p' | sort -n)"
    if [[ -z "$candidates" ]]; then
      return
    fi
    if [[ "${PREFLIGHT_ALL_MISSIONS:-0}" == "1" ]]; then
      while IFS= read -r n; do refs+=("M$n"); done <<<"$candidates"
    else
      refs=("M$(printf '%s' "$candidates" | tail -1)")
    fi
  fi

  if [[ ${#refs[@]} -eq 0 ]]; then
    return
  fi

  local ref out
  for ref in "${refs[@]}"; do
    if ! out="$(cd "$repo_root" && go run ./cmd/spectacular mission check "$ref" --json 2>&1)"; then
      preflight_fail "tier1/mission-check[$ref]: $(printf '%s' "$out" | head -c 400)"
      continue
    fi
    local verdict
    verdict="$(printf '%s' "$out" | python3 -c '
import json, sys
try:
    doc = json.load(sys.stdin)
except Exception as exc:
    print("unparseable mission check output: %s" % exc)
    sys.exit(0)
data = doc.get("data") or {}
problems = []
if not data.get("valid", False):
    problems.append("valid=false")
for entry in data.get("drift") or []:
    if entry.get("verdict") not in ("pass", None):
        problems.append("drift:%s=%s" % (entry.get("claim"), entry.get("verdict")))
print("; ".join(problems))
' 2>&1)"
    if [[ -n "$verdict" ]]; then
      preflight_fail "tier1/mission-check[$ref]: $verdict"
    fi
  done
}

preflight_checks() {
  local started_ns finished_ns elapsed_ms status
  started_ns="$(python3 -c 'import time; print(time.time_ns())')"

  preflight_tier0
  preflight_tier1

  finished_ns="$(python3 -c 'import time; print(time.time_ns())')"
  elapsed_ms=$(( (finished_ns - started_ns) / 1000000 ))

  if [[ ${#preflight_failures[@]} -eq 0 ]]; then
    status="pass"
  else
    status="fail"
  fi

  {
    printf '{\n'
    printf '  "schema_version": "spectacular.preflight-receipt.v1",\n'
    printf '  "status": "%s",\n' "$status"
    printf '  "failures": ['
    local first=1 failure
    for failure in ${preflight_failures[@]+"${preflight_failures[@]}"}; do
      if [[ $first -eq 1 ]]; then printf '\n    '; first=0; else printf ',\n    '; fi
      preflight_json_escape "$failure" | tr -d '\n'
    done
    if [[ $first -eq 0 ]]; then printf '\n  '; fi
    printf '],\n'
    printf '  "cost_units_used": %s,\n' "$elapsed_ms"
    printf '  "cost_units": "milliseconds",\n'
    printf '  "tiers": ["tier0-static-syntax", "tier1-contract-drift"],\n'
    printf '  "mutation": "none"\n'
    printf '}\n'
  }

  if [[ "$status" == "fail" ]]; then
    echo "Spectacular pre-flight: FAIL — repair before running acceptance/release/all" >&2
    exit 1
  fi
  echo "Spectacular pre-flight: PASS elapsed_ms=$elapsed_ms" >&2
}

static_checks() {
  formatting="$(cd "$repo_root" && gofmt -l cmd internal test/acceptance)"
  if [[ -n "$formatting" ]]; then
    echo "gofmt required:" >&2
    echo "$formatting" >&2
    exit 1
  fi
  manifest_checks
  security_checks
  check go-mod-verify go mod verify
  check go-vet go vet ./...
  check diff-check git diff --check
}

quick_checks() {
  preflight_checks >/dev/null
  static_checks
  check focused-go-test go test ./cmd/... ./install/... ./internal/...
}

acceptance_checks() {
  check acceptance go test -count=1 ./test/acceptance
}

release_checks() {
  check install-distribution bash install/test.sh
  check release-distribution bash test/release.sh
}

case "$mode" in
  preflight)
    preflight_checks
    exit 0
    ;;
  quick)
    quick_checks
    ;;
  acceptance)
    static_checks
    acceptance_checks
    ;;
  release)
    static_checks
    release_checks
    ;;
  all)
    static_checks
    acceptance_checks
    check race go test -race -count=1 ./...
    release_checks
    ;;
  *)
    echo "usage: bash test/verify.sh [preflight|quick|acceptance|release|all]" >&2
    exit 2
    ;;
esac

echo "Spectacular verification: PASS mode=$mode basis=$(tree_basis) logs=compact"
