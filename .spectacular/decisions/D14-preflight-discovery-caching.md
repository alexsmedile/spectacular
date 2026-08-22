---
type: Decision
id: 01a029c7-e306-7fd9-949a-b7c6160cae63
title: Enable fast-path session preflight caching bounded by Git tree SHA
created_by: Alex
created: "2026-08-22T16:02:40Z"
updated: "2026-08-22T16:02:40Z"
actor: Alex
actor_role: owner
ref: D14-preflight-discovery-caching
question: How should Spectacular minimize launch discovery latency across repeated sessions?
disposition: git-tree-bound-preflight-cache
rationale: >-
    Re-deriving all 13 preflight dimensions on every conversational turn inflates token consumption
    and adds multi-second latency. When the Git tree SHA and governance record fingerprints match
    the previous session's `.spectacular/.preflight-cache`, the agent reuses the validated preflight
    state, achieving sub-100ms session start while retaining zero-drift safety.
alternatives:
    - unconditional full 13-point re-derivation on every session
    - loose time-based caching without cryptographic Git tree binding
authority_basis: Owner explicitly approved Option A (Fast Git-bound Preflight Cache) in the design interview.
---

# Enable fast-path session preflight caching bounded by Git tree SHA

## Decision
- Persist a deterministic `.spectacular/.preflight-cache` recording verified governance and environment state tied to `git rev-parse HEAD`.
- If the tree SHA and governance records are unchanged, launch immediately on the fast path.
