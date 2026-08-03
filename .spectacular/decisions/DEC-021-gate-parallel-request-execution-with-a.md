---
id: DEC-021
type: decision
status: verified
origin: user-derived
derived_from: "GitHub Issue #3 and user confirmation on 2026-08-03"
supersedes: ""
evidence: "user direction"
tags: [github, git, branches, request, lifecycle, agents, traffic]
updated: 2026-08-03
---

# DEC-021 — Gate parallel request execution with a traffic preflight

**Context:**
Issue #3 exposed that isolated branches prevent direct overwrites but do not stop two requests or agents from making incompatible changes. The user confirmed four launch states: parallel, conditional, serialized, and unknown.

**Evidence:**
User direction recorded by explicit command.

**Alternatives considered:**
Recorded in the cited user context or decision discussion.

**Decision:**
Gate parallel request execution with a traffic preflight

**Consequences:**
Requests permanently retain blocked-by, blocks, and conflicts-with relationships. Spectacular recalculates traffic state at scaffold and activation against current requests, branches, and pull requests; unknown or changed conditions return to the orchestrator or user. Approved run scope avoids repeated confirmations for ordinary in-scope steps but never bypasses declared HITL gates.
