# Task

Build a concurrent in-memory job queue and worker pool engine in `src/` (Go, Python, or Node).

Requirements:
- Priority Scheduling: Jobs have priorities (`high`, `normal`, `low`). `high` priority jobs must execute before `normal` or `low`.
- Concurrency: Processes jobs concurrently across workers (default 4 workers).
- Dead-Letter Queue (DLQ): If a job fails processing 3 times (e.g. simulated failure), it must be routed to `dlq.json` with error details and attempt counts.
- Run `sh tests/check.sh` to verify your implementation before reporting.
