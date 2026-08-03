---
id: DEC-026
type: decision
status: superseded
origin: user-derived
derived_from: "workspace-migration-readiness interview and user confirmation on 2026-08-03"
supersedes: ""
evidence: "user direction"
tags: [migration, retention, git, workspace, simplicity, docs]
updated: 2026-08-03
superseded_by: DEC-027
---

# DEC-026 — Remove the live migrations log and rely on existing provenance

**Context:**
The readiness audit found one tracked root migrations.log line duplicating the current workspace schema, migration registry, Spectacular provenance, and Git commit history. The minimal-root contract does not declare the log.

**Evidence:**
User direction recorded by explicit command.

**Alternatives considered:**
Recorded in the cited user context or decision discussion.

**Decision:**
Remove the live migrations log and rely on existing provenance

**Consequences:**
Stop appending migrations.log and remove it from the live workspace in the later cleanup. Current state remains in workspace_schema and last_touched_with; migration definitions remain in the registry; active evidence lives in the owning request/session; completed history remains in Git. Do not add a migration folder, index, or replacement provenance schema until real support cases earn it.
