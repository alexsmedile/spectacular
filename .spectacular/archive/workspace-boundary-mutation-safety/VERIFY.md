---
updated: 2026-08-06
related:
  - PLAN.md
  - TASKS.md
---

# Verify — workspace-boundary-mutation-safety

## Automated {run}

- [x] bash -n cli/spectacular tests/cli/init.test.sh tests/cli/migrate.test.sh
- [x] bash tests/run.sh

## Contract assertions {assert}

- [x] `CURRENT_SCHEMA` remains `2.0` and the migration registry has no `2.0 → 3.0` edge.
- [x] The schema write gate refuses a newer workspace before a representative mutator writes, while `status --against-latest` remains diagnostic and does not suggest migration.
- [x] Migration refuses a synthetic tracked `.spectacular.local/` pathname without printing its synthetic file content.
- [x] `init-workflow.md` describes local state as supplement-only and does not claim local precedence.

## Boundary judgment {judge}

- [x] The central command classifier preserves named read-only CLI views while conservatively gating unclassified write-capable invocations.
- [x] The delivered diff contains no schema bump, migration apply edge, private-content operation, GitHub configuration change, or root-artifact cleanup.
