---
id: DEC-019
type: decision
status: verified
origin: user-derived
derived_from: "Explicit user direction, 2026-08-02"
supersedes: ""
evidence: "user direction"
tags: [github, pr, request, lifecycle]
updated: 2026-08-02
---

# DEC-019 — Every executed Spectacular request ends through a pull request before integration

**Context:**
Requests need a consistent collaborative review and integration boundary even when implementation begins locally or autonomously.

**Evidence:**
User direction recorded by explicit command.

**Alternatives considered:**
Recorded in the cited user context or decision discussion.

**Decision:**
Every executed Spectacular request ends through a pull request before integration

**Consequences:**
Request execution must produce or link a PR; merge remains human-gated; draft-versus-ready timing and exceptions are deferred to SPC-003 grilling.
