---
type: owner-decision
decision: v1-deprioritization
version: 1.0
status: settled-current
decider: owner
decided_at: 2026-08-09
supersedes_scope:
  - CLEAN-BREAK-CUTOVER-AND-RECOVERY-CONTRACT.md@1.0#Frozen-v1-boundary
  - IMPLEMENTATION-ARCHITECTURE-AND-MIGRATION-CONTRACT.md@1.0#Migration-and-cutover
  - EXECUTABLE-REFACTOR-PROGRAM-CONTRACT.md@1.1#P0-and-v1-critical-path
next_action: W0-shared-scaffold-design-sufficiency
---

# V1 deprioritization decision

## Owner direction

> “I don't care about v1. If there is something we can skip. Pls do it.”

## Decision

Spectacular stops spending product and implementation effort on v1. P0 is abandoned without
merging its implementation branch. Its bounded implementation and independent-review evidence are
retained as history, but its repair, Pageworks reconciliation, and acceptance gates are cancelled.

No new final v1 release is required. Existing Git history and already-published tags are recovery
references, not a continuing support or compatibility promise. The clean v2 core and its release
must not wait for v1 stabilization, a universal v1 mapping inventory, or a generic migration
capsule.

Migration of an actual project is a later project-specific Mission. It may use disposable isolated
scripts, helper references, or a narrowly earned capsule, but it must preserve the source snapshot,
report ambiguity, validate a separate candidate, and require owner cutover. None of that logic may
enter v2 core.

## Program consequence

```text
W0 shared-scaffold Design Sufficiency
  → M1 semantic records + canonical workspace substrate
  → M2 durable governed Mission loop
  → M3 guided Skill + registry CLI + retrieval/integrity
  → M4 clean-v2 release readiness
```

The prior M4 frozen-v1/inventory Mission and M5 generic-capsule Mission leave the critical path.
Their useful safety properties remain requirements of any later project-specific migration.
