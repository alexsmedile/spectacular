# Start and resume

## Start

Start only an owner-confirmed Mission with sufficient design, no blocking Gap, exact Contract and
Git baseline, complete claim coverage, coherent Objectives, explicit authority/scope, review level,
budgets, and stops. Record owner and time plus a fingerprint over outcome, review, completion,
authority, scope, budgets, dependencies, Gaps, and stops. Do not include mutable status, Objective
progress, Run state, or repair count in that fingerprint.

A normal start creates only `MISSION.md` with inline Objectives and R1. A Proposal reference is
optional input and is never created by Mission start. A Decision is not activation authority.

Use supported typed tooling for atomic generation and validation. Under `manual-bootstrap`, create
the same canonical shape directly and verify the frozen fingerprint and structural invariants with
focused scripts; do not route the work through an incompatible legacy command sequence.

## Resume

Read the Mission card, current Objective section, and exact promoted/source pointers. Recheck the
Contract fingerprint, Git baseline, activation fingerprint, validation mode, authority, scope,
budgets, dependencies/Gaps, and stops. A material semantic change returns to the owner; reversible
implementation changes remain with the operator.

Plan outcome-sized clusters:
`[claims + dependencies] -> [work] -> [focused checks] -> [boundary integration] -> [local commit]`.
Run one full repository/release gate after integration rather than repeating it per small edit.
Inspect detailed logs only on failure. Keep Git/secret/distribution checks at the commit, push, PR,
or release boundary where they apply.
