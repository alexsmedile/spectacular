---
id: DEC-024
type: decision
status: verified
origin: user-derived
derived_from: "workspace-migration-readiness report and user confirmation on 2026-08-03"
supersedes: "DEC-023"
evidence: "user direction"
tags: [migration, schema, compatibility, workspace, lifecycle]
updated: 2026-08-03
---

# DEC-024 — Keep workspace schema 2.0 until a real breaking contract exists

**Context:**
The migration-readiness audit found that the currently proposed local-boundary, compatibility, doctor, and safety improvements are additive. No concrete breaking layout or required-field delta presently earns schema 3.0.

**Evidence:**
User direction recorded by explicit command.

**Alternatives considered:**
Recorded in the cited user context or decision discussion.

**Decision:**
Keep workspace schema 2.0 until a real breaking contract exists

**Consequences:**
Ship compatible contract hardening under schema 2.0. Reserve schema 3.0 for a future explicitly approved breaking change. D23's staged safety, branch, compatibility, and recovery protocol remains applicable when that real migration exists, but its assumed immediate 2.0-to-3.0 transition is superseded.
