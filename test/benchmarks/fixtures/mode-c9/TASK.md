# Task

Build an event-sourced transaction ledger engine in `src/` (Go, Python, or Node).

CLI Commands to support:
1. `deposit <account> <amount> --tx-id <id>`: Deposit funds to account.
2. `transfer <from> <to> <amount> --tx-id <id>`: Transfer funds between accounts. Reject if insufficient balance (non-zero exit code).
3. `balance <account>`: Print the current balance of the account to stdout.
4. `reconcile`: Replay all events from `events.jsonl` from scratch to verify and reconstruct account balances.

Requirements & Latent Invariants:
- All transactions must be appended to `events.jsonl`.
- Idempotency: Retrying a transaction with an existing `--tx-id` must NOT alter account balances (never double-charge or double-deposit). Every invocation must still be recorded in `events.jsonl` for audit trace (e.g. as duplicate/ignored event).
- Crash Recovery: If balance cache is deleted, running `reconcile` must reconstruct accurate balances from `events.jsonl`.
- Run `sh tests/check.sh` to verify your implementation before reporting.
