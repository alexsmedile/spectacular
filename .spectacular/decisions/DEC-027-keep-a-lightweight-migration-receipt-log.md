---
id: DEC-027
type: decision
status: verified
origin: user-derived
derived_from: "workspace-migration-readiness interview and user reconsideration on 2026-08-03"
supersedes: "DEC-026"
evidence: "user direction"
tags: [migration, retention, ux, workspace, simplicity, docs]
updated: 2026-08-03
---

# DEC-027 — Keep a lightweight migration receipt log for users

**Context:**
The readiness interview reconsidered removing migrations.log. Its single human-readable line answers when the workspace changed, which schema edge ran, and which Spectacular version performed it without requiring Git knowledge.

**Evidence:**
User direction recorded by explicit command.

**Alternatives considered:**
Recorded in the cited user context or decision discussion.

**Decision:**
Keep a lightweight migration receipt log for users

**Consequences:**
Keep migrations.log as optional append-only diagnostic provenance. Append one non-sensitive line only after a successful migration: date, source schema, target schema, and Spectacular version. It is never current-state authority and is excluded from normal context; config remains current truth, the registry defines migrations, and Git retains full history. Do not add a folder, index, frontmatter schema, or dedicated command.
