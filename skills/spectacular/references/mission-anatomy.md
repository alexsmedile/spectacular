# Mission anatomy

The full field lists. `SKILL.md` carries the summary; this file carries the
inventory. Load it when you are writing or auditing a `MISSION.md`, not on every
session.

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
