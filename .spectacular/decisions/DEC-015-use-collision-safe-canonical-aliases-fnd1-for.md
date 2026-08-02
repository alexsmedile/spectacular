---
id: DEC-015
type: decision
status: superseded
origin: user-derived
derived_from: "User confirmations in request-workflow interview, 2026-08-02"
supersedes: ""
evidence: "user direction"
tags: [ids, aliases, slug, cli]
updated: 2026-08-02
superseded_by: DEC-017
---

# DEC-015 — Use collision-safe canonical aliases: fnd1 for FND, f1 for fixes, bug1 for BUG, and b1 for roadmap builds

**Context:**
The proposed compact aliases collided with existing verified-fix F<N> and roadmap build b<N> identities.

**Evidence:**
User direction recorded by explicit command.

**Alternatives considered:**
Recorded in the cited user context or decision discussion.

**Decision:**
Use collision-safe canonical aliases: fnd1 for FND, f1 for fixes, bug1 for BUG, and b1 for roadmap builds.

**Consequences:**
Canonical IDs remain unambiguous; context-free f1 resolves to a fix and b1 to a build; findings and bugs use fnd1 and bug1. #ids #aliases #slug #cli
