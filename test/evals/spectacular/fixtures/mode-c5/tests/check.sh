#!/bin/sh
set -eu

RUNNER=""
if [ -f "src/task_cli" ] && [ -x "src/task_cli" ]; then
  RUNNER="./src/task_cli"
elif [ -f "src/main.py" ]; then
  RUNNER="python3 src/main.py"
elif [ -f "src/task.py" ]; then
  RUNNER="python3 src/task.py"
elif [ -f "src/main.js" ]; then
  RUNNER="node src/main.js"
elif [ -f "src/task.js" ]; then
  RUNNER="node src/task.js"
elif [ -f "src/main.go" ]; then
  RUNNER="go run ./src/..."
elif [ -f "src/task.sh" ]; then
  RUNNER="sh src/task.sh"
else
  for f in src/*; do
    if [ -x "$f" ] && [ ! -d "$f" ]; then
      RUNNER="$f"
      break
    fi
  done
fi

if [ -z "$RUNNER" ]; then
  echo "check.sh: no executable runner found in src/" >&2
  exit 1
fi

rm -f tasks.db
trap 'rm -f tasks.db' EXIT

# 1. Add tasks
$RUNNER add "Setup repo" --tag dev >/dev/null
$RUNNER add "Write docs" --tag docs >/dev/null
$RUNNER add "Ship release" >/dev/null

# 2. List all
LIST_OUTPUT=$($RUNNER list)
echo "$LIST_OUTPUT" | grep -qi "Setup repo"
echo "$LIST_OUTPUT" | grep -qi "Write docs"
echo "$LIST_OUTPUT" | grep -qi "Ship release"

# 3. Tag filtering
DOCS_ONLY=$($RUNNER list --tag docs)
echo "$DOCS_ONLY" | grep -qi "Write docs"
if echo "$DOCS_ONLY" | grep -qi "Setup repo"; then
  echo "check.sh: tag filter failed" >&2
  exit 1
fi

# 4. Mark done and assert status
$RUNNER done 1 >/dev/null
UPDATED_LIST=$($RUNNER list)
echo "$UPDATED_LIST" | grep -qi "done"

# 5. Invalid command returns non-zero
if $RUNNER invalid_subcommand_test >/dev/null 2>&1; then
  echo "check.sh: expected failure on invalid command" >&2
  exit 1
fi

echo "TASK_CLI_GENESIS_CHECK_PASS"
