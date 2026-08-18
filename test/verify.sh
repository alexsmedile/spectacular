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

static_checks() {
  formatting="$(cd "$repo_root" && gofmt -l cmd internal test/acceptance)"
  if [[ -n "$formatting" ]]; then
    echo "gofmt required:" >&2
    echo "$formatting" >&2
    exit 1
  fi
  check go-mod-verify go mod verify
  check go-vet go vet ./...
  check diff-check git diff --check
}

quick_checks() {
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
    check race go test -race -count=1 ./...
    release_checks
    ;;
  *)
    echo "usage: bash test/verify.sh [quick|acceptance|release|all]" >&2
    exit 2
    ;;
esac

echo "Spectacular verification: PASS mode=$mode basis=$(tree_basis) logs=compact"
