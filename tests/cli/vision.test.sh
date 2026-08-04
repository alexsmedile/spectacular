#!/usr/bin/env bash
# Pre-request Vision lifecycle, human reactions, derivation gate, and legacy compatibility.

set -u

REPO_ROOT="$(cd "$(dirname "$0")/../.." && pwd)"
CLI="$REPO_ROOT/cli/spectacular"
WS="$(mktemp -d)"
pass_count=0
fail_count=0
trap 'rm -rf "$WS"' EXIT

pass() { echo "    ✓ $1"; pass_count=$((pass_count + 1)); }
fail() { echo "    ✗ $1"; fail_count=$((fail_count + 1)); }
assert_contains() { echo "$1" | grep -qF -- "$2" && pass "$3" || fail "$3 — missing: $2"; }
assert_file() { [[ -f "$1" ]] && pass "$2" || fail "$2 — missing file: $1"; }
assert_exit() { [[ "$1" -eq "$2" ]] && pass "$3" || fail "$3 — exit $1, expected $2"; }

mkdir -p "$WS/.spectacular/requests"
printf 'project:\n  name: vision-test\n  owner: test\n' > "$WS/.spectacular/config.yaml"

echo "── scaffold is pre-request and proportional ──"
out=$(cd "$WS" && bash "$CLI" imagine better-onboarding 2>&1)
assert_contains "$out" "status: draft · pre-request" "imagine scaffolds a pre-request Vision"
assert_file "$WS/.spectacular/visions/better-onboarding/VISION.md" "Vision spine exists without a request"
[[ ! -d "$WS/.spectacular/requests/better-onboarding" ]] && pass "imagine does not manufacture a request" || fail "imagine manufactured a request"

(cd "$WS" && bash "$CLI" vision add strategy guided --slug better-onboarding --caption "Guided entry" >/dev/null)
(cd "$WS" && bash "$CLI" vision add ui welcome --slug better-onboarding --caption "Welcome screen" >/dev/null)
assert_file "$WS/.spectacular/visions/better-onboarding/fragments/guided.md" "strategy fragment uses shared fragments folder"
assert_file "$WS/.spectacular/visions/better-onboarding/fragments/welcome.md" "UI fragment is optional and independently addressable"

code=0; out=$(cd "$WS" && bash "$CLI" vision propose better-onboarding 2>&1) || code=$?
assert_exit "$code" 1 "empty Vision cannot be proposed"
assert_contains "$out" "section 'Intent' is empty" "proposal gate names the first empty section"

SPINE="$WS/.spectacular/visions/better-onboarding/VISION.md"
tmp=$(mktemp)
awk '
  { print }
  /^## Intent$/ { print "\nMake first use immediately understandable." }
  /^## North star$/ { print "\nA newcomer confidently completes the first meaningful action." }
  /^### Current reality$/ { print "\nThe current entry point exposes concepts before an outcome." }
  /^### Users and needs$/ { print "\nNew users need a visible safe first action." }
  /^### Constraints$/ { print "\nKeep expert paths direct and preserve terminal readability." }
  /^### Material uncertainties$/ { print "\nWhether guidance should be inline or a dedicated screen." }
  /^## Experience signature$/ { print "\nCalm, obvious, and reversible." }
  /^## Strategies considered$/ { print "\nGuided entry versus contextual hints." }
  /^## Chosen direction$/ { print "\nUse a guided entry with an immediate escape hatch." }
  /^## Boundaries$/ { print "\nNo redesign of expert commands." }
' "$SPINE" > "$tmp" && mv "$tmp" "$SPINE"

(cd "$WS" && bash "$CLI" vision react better-onboarding guided --approve --note "Best north-star fit" >/dev/null)
(cd "$WS" && bash "$CLI" vision react better-onboarding welcome --reject --note "Too much ceremony" >/dev/null)
assert_contains "$(cat "$SPINE")" '**approved** (`fragments/guided.md`)' "manifest mirrors approved reaction"
assert_contains "$(cat "$SPINE")" '**rejected** (`fragments/welcome.md`)' "manifest mirrors rejected reaction"

(cd "$WS" && bash "$CLI" vision propose better-onboarding >/dev/null)
assert_contains "$(cat "$SPINE")" "status: proposed" "resolved Vision enters proposed state"
code=0; (cd "$WS" && bash "$CLI" vision add flow late --slug better-onboarding >/dev/null 2>&1) || code=$?
assert_exit "$code" 1 "proposed Vision cannot acquire new fragments"

code=0; (cd "$WS" && bash "$CLI" vision approve better-onboarding >/dev/null 2>&1) || code=$?
assert_exit "$code" 1 "whole-Vision approval requires a named human"
(cd "$WS" && bash "$CLI" vision approve better-onboarding --approved-by alex --note "Direction confirmed" >/dev/null)
assert_contains "$(cat "$SPINE")" "status: approved" "whole Vision records approved lifecycle state"
assert_contains "$(cat "$SPINE")" "approved_by: alex" "whole Vision records approval provenance"

out=$(cd "$WS" && bash "$CLI" vision list 2>&1)
assert_contains "$out" "better-onboarding" "vision list includes first-class workspaces"
out=$(cd "$WS" && bash "$CLI" vision show better-onboarding 2>&1)
assert_contains "$out" "Status: approved" "vision show reports whole-Vision state"

code=0; out=$(cd "$WS" && bash "$CLI" vision derive better-onboarding 2>&1) || code=$?
assert_exit "$code" 1 "derive redirects to the agentic flow"
assert_contains "$out" "Derive a draft SPC" "derive targets a draft SPC, not PLAN"

(cd "$WS" && bash "$CLI" imagine discarded-direction >/dev/null)
(cd "$WS" && bash "$CLI" vision reject discarded-direction --rejected-by alex --reason "Wrong north star" >/dev/null)
assert_contains "$(cat "$WS/.spectacular/visions/discarded-direction/VISION.md")" "status: rejected" "draft Vision can preserve an explicit rejection"

echo "── doctor and path discovery ──"
doctor=$(cd "$WS" && bash "$CLI" doctor vision 2>&1)
assert_contains "$doctor" "better-onboarding · approved" "doctor validates first-class Vision lifecycle"
paths=$(cd "$WS" && bash "$CLI" paths --text)
assert_contains "$paths" "visions_dir" "paths exposes the Vision collection"

FRAG="$WS/.spectacular/visions/better-onboarding/fragments/guided.md"
tmp=$(mktemp)
awk 'BEGIN{changed=0} !changed && /^reaction: approved$/ { print "reaction: rejected"; changed=1; next } { print }' "$FRAG" > "$tmp" && mv "$tmp" "$FRAG"
doctor=$(cd "$WS" && bash "$CLI" doctor vision 2>&1)
assert_contains "$doctor" "spine manifest out of date" "doctor detects manifest drift"
(cd "$WS" && bash "$CLI" doctor vision --fix >/dev/null)
doctor=$(cd "$WS" && bash "$CLI" doctor vision 2>&1)
if echo "$doctor" | grep -q "spine manifest out of date"; then fail "doctor fix did not repair manifest"; else pass "doctor fix repairs only manifest drift"; fi

echo "── legacy request Vision remains readable ──"
LEGACY="$WS/.spectacular/requests/legacy/vision"
mkdir -p "$LEGACY/ui"
printf '%s\n' '---' 'doc: vision' 'updated: 2026-01-01' '---' '# Vision — legacy' '' '## Manifest' '' '_No fragments yet. Add only what makes a material uncertainty discussable._' > "$LEGACY/VISION.md"
out=$(cd "$WS" && bash "$CLI" vision show legacy 2>&1)
assert_contains "$out" "Status: legacy" "legacy request Vision remains readable"
doctor=$(cd "$WS" && bash "$CLI" doctor vision 2>&1)
assert_contains "$doctor" "legacy · legacy" "doctor diagnoses legacy request Vision"

echo "── experiment is no longer a feedback alias ──"
code=0; out=$(cd "$WS" && bash "$CLI" experiment list 2>&1) || code=$?
assert_exit "$code" 1 "experiment no longer dispatches to feedback-loop"
assert_contains "$out" "Unknown command: 'experiment'" "ambiguous experiment asks for explicit routing"

echo ""
echo "Results: $pass_count passed, $fail_count failed"
[[ $fail_count -eq 0 ]]
