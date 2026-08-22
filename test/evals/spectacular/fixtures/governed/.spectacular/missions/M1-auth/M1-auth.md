---
type: Mission
id: 019fe381-5d61-7223-b362-03a5f99a7b20
ref: M1
title: Repair authentication behavior
status: active
owner: Alex
outcome: Authentication behavior is bounded, tested, and recoverable.
completion:
  - claim: auth-check
    pass_boundary: The auth fixture rejects an empty token and accepts a non-empty token.
    proof_requirement: Running `sh tests/check.sh` exits zero.
objectives:
  - ref: O1
    outcome: Repair the bounded auth fixture.
    claims: [auth-check]
run:
  ref: R1
  current_objective: O1
  status: active
---
# Mission

## Recovery

Last safe point: inspect `src/auth.txt`; no implementation change has been accepted.
