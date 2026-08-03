---
id: DEC-025
type: decision
status: verified
origin: user-derived
derived_from: "workspace-migration-readiness interview and user confirmation on 2026-08-03"
supersedes: ""
evidence: "user direction"
tags: [local, workspace, schema, simplicity, security, compatibility]
updated: 2026-08-03
---

# DEC-025 — Do not version private local state until incompatibility earns it

**Context:**
The readiness interview considered a root local schema and feature-owned versions before the local features had real incompatible formats. Both would add migration ceremony to sparse, lazy, mostly disposable state.

**Evidence:**
User direction recorded by explicit command.

**Alternatives considered:**
Recorded in the cited user context or decision discussion.

**Decision:**
Do not version private local state until incompatibility earns it

**Consequences:**
No `.spectacular.local/` root schema and no version field by default. Plain Markdown stays plain, caches are recreated, and small settings files keep a stable shape. Add a narrowly owned format marker only after an actual incompatible change cannot be safely detected or migrated. Today the load-bearing protections are gitignore, non-override of shared truth/authority, and fail-closed tracked-path detection.
