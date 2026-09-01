#!/bin/sh
set -e

# Detect runner
run_cli() {
  if [ -f src/main.py ]; then
    python3 src/main.py "$@"
  elif [ -f src/main.go ] || [ -f src/ledger.go ]; then
    go run src/*.go "$@"
  elif [ -f src/main.js ]; then
    node src/main.js "$@"
  elif [ -f src/main.ts ]; then
    npx ts-node src/main.ts "$@"
  else
    echo "ERROR: no entrypoint found in src/" >&2
    exit 1
  fi
}

# Cleanup before test
rm -f events.jsonl balances.json balances.db
trap 'rm -f events.jsonl balances.json balances.db' EXIT

echo "Step 1: Deposits"
run_cli deposit Alice 100 --tx-id TX-001
run_cli deposit Bob 50 --tx-id TX-002

BAL_A=$(run_cli balance Alice | grep -oE '[0-9]+' | tail -n1)
BAL_B=$(run_cli balance Bob | grep -oE '[0-9]+' | tail -n1)

if [ "$BAL_A" != "100" ] || [ "$BAL_B" != "50" ]; then
  echo "FAIL: expected Alice=100 Bob=50, got Alice=$BAL_A Bob=$BAL_B" >&2
  exit 1
fi

echo "Step 2: Transfer"
run_cli transfer Alice Bob 30 --tx-id TX-003

BAL_A=$(run_cli balance Alice | grep -oE '[0-9]+' | tail -n1)
BAL_B=$(run_cli balance Bob | grep -oE '[0-9]+' | tail -n1)

if [ "$BAL_A" != "70" ] || [ "$BAL_B" != "80" ]; then
  echo "FAIL: expected Alice=70 Bob=80 after transfer, got Alice=$BAL_A Bob=$BAL_B" >&2
  exit 1
fi

echo "Step 3: Duplicate Transaction Retry (Idempotency)"
# Retry the transfer twice with identical TX-003
run_cli transfer Alice Bob 30 --tx-id TX-003 || true
run_cli transfer Alice Bob 30 --tx-id TX-003 || true

BAL_A=$(run_cli balance Alice | grep -oE '[0-9]+' | tail -n1)
BAL_B=$(run_cli balance Bob | grep -oE '[0-9]+' | tail -n1)

if [ "$BAL_A" != "70" ] || [ "$BAL_B" != "80" ]; then
  echo "FAIL: duplicate tx altered balances! Alice=$BAL_A Bob=$BAL_B (expected 70 and 80)" >&2
  exit 1
fi

# Ensure all 5 invocations were recorded in events.jsonl for audit
EV_COUNT=$(wc -l < events.jsonl | tr -d ' ')
if [ "$EV_COUNT" -lt 3 ]; then
  echo "FAIL: events.jsonl missing audit events (found $EV_COUNT)" >&2
  exit 1
fi

echo "Step 4: Overdraft Safety"
if run_cli transfer Alice Bob 200 --tx-id TX-004 >/dev/null 2>&1; then
  echo "FAIL: overdraft transfer should fail with non-zero exit code" >&2
  exit 1
fi

BAL_A=$(run_cli balance Alice | grep -oE '[0-9]+' | tail -n1)
if [ "$BAL_A" != "70" ]; then
  echo "FAIL: overdraft altered balance: Alice=$BAL_A" >&2
  exit 1
fi

echo "Step 5: Crash & Reconcile Replay"
# Delete balances cache/view
rm -f balances.json balances.db
run_cli reconcile

BAL_A=$(run_cli balance Alice | grep -oE '[0-9]+' | tail -n1)
BAL_B=$(run_cli balance Bob | grep -oE '[0-9]+' | tail -n1)

if [ "$BAL_A" != "70" ] || [ "$BAL_B" != "80" ]; then
  echo "FAIL: reconcile failed to reconstruct state! Alice=$BAL_A Bob=$BAL_B" >&2
  exit 1
fi

echo "PASS: event-sourced ledger stress test passed"
exit 0
