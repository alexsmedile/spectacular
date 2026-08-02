---
id: DEC-017
type: decision
status: verified
origin: user-derived
derived_from: "User correction in conversation, 2026-08-02"
supersedes: "DEC-015"
evidence: "user direction"
tags: [ids, aliases, cli, fixes, findings]
updated: 2026-08-02
---

# DEC-017 — Prefer explicit fix1 and fnd1 aliases; refuse ambiguous f1 while preserving b1 for roadmap builds

**Context:**
The previously confirmed alias rule made f1 mean a legacy fix, but the user clarified that conversational aliases should name their entity explicitly.

**Evidence:**
User direction recorded by explicit command.

**Alternatives considered:**
Recorded in the cited user context or decision discussion.

**Decision:**
Prefer explicit fix1 and fnd1 aliases; refuse ambiguous f1 while preserving b1 for roadmap builds

**Consequences:**
Fix lookup remains backward-compatible through fix1 to legacy F1 without migrating stored IDs; findings use fnd1; f1 fails with corrective guidance; b1 remains unchanged.
