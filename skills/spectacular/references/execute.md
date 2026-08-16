# Start and resume

## Start

Start only an owner-confirmed Mission. Before activating, confirm all of it:

- the design is sufficient, with no blocking Gap
- exact Contract and Git baseline
- complete claim coverage
- coherent Objectives
- explicit authority and scope
- review level
- budgets and stops

Then record the owner and the activation time, plus a fingerprint. All three land
in the Mission's `activation:` block:

```yaml
activation:
    at: "2026-08-16T14:19:42Z"
    by: Alex
    fingerprint: sha256:ef7e695fc1f496fabb92c05a0369ca148fc6404e49b728521db3a2cdf120f389
```

**The fingerprint covers the frozen semantic envelope** — outcome, review,
completion, authority, scope, budgets, dependencies, Gaps, stops. It must not
cover mutable state: status, Objective progress, Run state, repair count.

The tooling computes it. Do not hand-roll a hash to make a check pass.

A normal start creates one file:
`.spectacular/missions/<slug>/MISSION.md`, with inline Objectives and R1.

Two things that are not activation authority:

- A Proposal is optional input. Mission start never creates one.
- A Decision is not activation authority. It records a choice; it does not
  authorize a start.

Use the typed tooling for atomic generation and validation:

```bash
spectacular mission start <plan.md|-> --json   # mutating; generates the canonical path
spectacular mission check <ref> --json         # read-only; verify after
```

`mission start` takes the approved plan as a Markdown file, or `-` for stdin. It
generates identities, bindings, and activation atomically — you do not author the
UUID or the fingerprint yourself.

Under `manual-bootstrap`, build the same canonical shape directly and verify it
with focused scripts — see [bootstrap.md](bootstrap.md). Never route the work
through an incompatible legacy command sequence to make it look validated.

## Resume

Read three things: the Mission card, the current Objective section, and the exact
pointers they name.

```bash
spectacular mission show <ref> --json    # card, current Objective, current Run
spectacular run show <ref>/<run> --json  # e.g. M7/R1, when the Run detail matters
```

Then recheck:

- Contract fingerprint
- Git baseline
- activation fingerprint
- validation mode
- authority and scope
- budgets
- dependencies and Gaps
- stops

A material semantic change goes back to the owner. Reversible implementation
changes stay with the operator.

## Work in outcome-sized clusters

```
[claims + dependencies] -> [work] -> [focused checks] -> [boundary integration] -> [local commit]
```

- Run one full repository and release gate after integration — `bash test/verify.sh all`.
  Do not repeat it per small edit.
- Read detailed logs only when something fails.
- Keep Git, secret, and distribution checks at the boundary where they apply:
  commit, push, PR, or release.
