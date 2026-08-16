# Mission anatomy

The full field lists. `SKILL.md` carries the summary; this file carries the
inventory. Load it when you are writing or auditing a `MISSION.md`, not on every
session.

## Shape

A real `MISSION.md`, trimmed to one of each thing:

```yaml
---
type: Mission
id: 01a00af1-38c0-7268-9529-5856afc7b2f2   # UUIDv7, generated, never hand-written
ref: M7                                     # human navigation; sub-records are M7/O1, M7/R1
title: Render derived state and validate the Proposal record
status: active
activation:
    at: "2026-08-16T14:19:42Z"
    by: Alex
    fingerprint: sha256:ef7e695f...        # over the frozen envelope only
authority:
    operator: [inspect, edit-in-scope, run-checks, bounded-repair, commit-local]
    requires_owner: [activate-mission, change-outcome-or-completion, expand-scope, push]
baseline:
    branch: m7-derived-state
    commit: 127dac140467a462c3810c85c9ca325c18278a14
contract:
    ref: Contract:01a00aae-8921-7b27-96a9-1a4c175e7dc6
    fingerprint: sha256:1ffd39b4...
completion:
    - claim: drift-flags
      pass_boundary: ...      # what must be observably true
      proof_requirement: ...  # what would demonstrate it
objectives:
    - id: 01a00af1-38c0-72bf-a3b2-0c8f1595b50d
      claims: [drift-flags]   # every Objective serves at least one claim
      outcome: ...
dependencies:
    - M6 completed with independent review and owner acceptance.
gaps: []
---

# Origin, rationale, detailed plans, review instructions
```

Live examples are in `.spectacular/missions/*/MISSION.md`. Read one before
authoring a new Mission by hand.

## Frontmatter

Frozen at activation, covered by the activation fingerprint:

- `id` — UUIDv7, durable identity
- `ref` — human-readable navigation reference
- `owner` — who alone may change the frozen parts
- activation time and activation fingerprint
- exact Contract binding
- exact Git baseline binding
- `outcome` — what this Mission is for
- completion claims — the frozen list to prove
- review level, including whether review must be independent
- authority, and the forbidden-effect ceiling
- mechanical scope and semantic scope
- budgets, including the repair budget
- dependencies, Gaps, stops

Mutable, and deliberately **outside** the fingerprint:

- `status`
- inline Objectives and their progress
- current Run and its state
- repair count
- validation mode

## Markdown body

- origin — where this Mission came from, including any Proposal reference
- rationale — why this approach, and what was rejected
- detailed Objective plans that outgrew the frontmatter
- bootstrap conditions, when the Mission declares `manual-bootstrap`
- examples
- review instructions for the reviewer

## What the tooling owns

The plan never chooses these. The active schema registry does:

- schema and vocabulary validation
- UUIDv7 and ref allocation
- fingerprint computation
- baseline checks
- dependency integrity
- safe path handling
- atomic multi-file transitions
- concurrency and retry behavior
- compact projections
- exact refusal messages

## What the plan owns

The tooling never decides these:

- outcome
- completion criteria
- decomposition into Objectives
- semantic scope
- authority
- dependencies, Gaps, stops
- rationale and prose

## Growth

Start with one file: `<mission>/MISSION.md`.

- Add `objectives/` when an Objective earns its own detail, delegation, owner, or
  independent review.
- Add `runs/` when a Run has a distinct job, operator, baseline, or recovery
  boundary.

When you split, keep the same UUID and ref, and leave one pointer where the inline
detail was. `MISSION.md` stays the index. Do not add a Mission-local `index.md`.
