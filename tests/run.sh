#!/usr/bin/env bash
# Test harness for Spectacular.
# Discovers tests/**/*.test.sh, runs each, reports pass/fail.
# Each test script must:
#   - Set up its own isolated workspace (typically under /tmp/)
#   - Clean up on exit (trap recommended)
#   - Exit 0 on pass, non-zero on fail
#   - Write progress to stdout

set -uo pipefail

REPO_ROOT="$(cd "$(dirname "$0")/.." && pwd)"
TESTS_DIR="$REPO_ROOT/tests"

# Allow filtering: ./tests/run.sh cli  → only runs tests under tests/cli/
FILTER="${1:-}"

# Discover tests (POSIX-portable, works on bash 3.2)
test_files=()
while IFS= read -r line; do
  test_files+=("$line")
done < <(find "$TESTS_DIR" -name "*.test.sh" -type f | sort)

if [[ ${#test_files[@]} -eq 0 ]]; then
  echo "No tests found under $TESTS_DIR"
  exit 1
fi

if [[ -n "$FILTER" ]]; then
  filtered=()
  for f in "${test_files[@]}"; do
    if [[ "$f" == *"/$FILTER/"* ]]; then
      filtered+=("$f")
    fi
  done
  if [[ ${#filtered[@]} -gt 0 ]]; then
    test_files=("${filtered[@]}")
  else
    test_files=()
  fi
fi

echo "Running ${#test_files[@]} test file(s)..."
echo ""

passed=0
failed=0
failures=()

for test_file in "${test_files[@]}"; do
  rel="${test_file#$REPO_ROOT/}"
  echo "── $rel ──────────────────────────────────────────────"
  if bash "$test_file"; then
    passed=$((passed + 1))
    echo "  ✓ $rel"
  else
    failed=$((failed + 1))
    failures+=("$rel")
    echo "  ✗ $rel"
  fi
  echo ""
done

echo "═══════════════════════════════════════════════════════"
echo "Results: $passed passed, $failed failed"
if [[ $failed -gt 0 ]]; then
  echo "Failed:"
  for f in "${failures[@]}"; do
    echo "  - $f"
  done
  exit 1
fi
exit 0
