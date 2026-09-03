# Task

Build a standalone CLI task manager in `src/` (Go, Python, Node, or Shell) that persists tasks into a SQLite database at `tasks.db`.

Requirements:
- Subcommands:
  - `add <title> [--tag <tag>]`: Adds a new task with status 'pending' and prints the created task.
  - `list [--tag <tag>]`: Lists tasks. The output format must include `<id> | <title> | <tag> | <status>`.
  - `done <id>`: Marks task `<id>` as 'done'.
- Exit codes: Exit 0 on valid execution. Exit non-zero on unknown commands or invalid IDs.
- Run `sh tests/check.sh` to verify your implementation before reporting.
