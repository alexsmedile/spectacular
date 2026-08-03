---
id: DEC-022
type: decision
status: verified
origin: user-derived
derived_from: "User migration interview confirmed on 2026-08-03"
supersedes: ""
evidence: "user direction"
tags: [local, workspace, security, privacy, git, lifecycle]
updated: 2026-08-03
---

# DEC-022 — Keep .spectacular.local subordinate, private, and lazily created

**Context:**
The GitHub-native design needs machine-specific accounts, remotes, private ideas, caches, and protected working material without allowing an ignored local layer to silently change shared project truth or authority.

**Evidence:**
User direction recorded by explicit command.

**Alternatives considered:**
Recorded in the cited user context or decision discussion.

**Decision:**
Keep .spectacular.local subordinate, private, and lazily created

**Consequences:**
Local state may customize machine/user operation but cannot override canonical repository identity, approved specs, decisions, request lifecycle, dependencies, authority mappings, permanent HITL/security rules, verification, or shared policy. Init guarantees gitignore protection; each feature creates only its needed path. Tracked local paths stop migration with filename-only reporting, and suspected leakage enters an explicit security response.
