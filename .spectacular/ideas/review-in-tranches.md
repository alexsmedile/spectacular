---
type: idea
status: parked
priority: medium
owner: alex
origin: (captured from a Harbor review session)
updated: 2026-07-04
promoted_to: null
related: [coding-agents]
---

# Idea — review-in-tranches

## Hypothesis

A skill for reviewing/refactoring a whole codebase **in small tranches (2-3 files at a time)**,
ranked by architectural risk — not a diff review, a *whole-repo sweep* that a single context
can't hold. Multi-agent: inventory once, then fan out one reviewer per tranche, each writing
findings to a persistent ledger so the sweep survives context limits and resumes cleanly.

The differentiator from a normal code review: it's **structured for scale and resumability**.
Output lives on disk (an index + a feedback ledger), tranches are ordered by risk so the most
dangerous code is reviewed first, and it stops before touching code — review/suggest only.

## Protocol (from the captured session)

1. **Inventory first** — `ls`/`rg --files` all real source, *excluding* generated/build/cache
   (`.git`, `DerivedData`, `build`, project bundles). Get line counts + types so priority
   reflects complexity, not filenames. Report the total + per-dir breakdown.
2. **Rank by risk & coupling** — order review targets by architectural danger (state safety,
   sync correctness, teardown lifecycle, trust boundaries) — not alphabetical, not by size
   alone. Group into feature-aligned tranches of 2-3 files.
3. **Persist the map before reviewing** — write the exact ordered file list to
   `INDEX_CODEBASE.md` (root) with a one-line responsibility note per file. This doubles as a
   codebase map the user can read to understand each part.
4. **Fan out, tranche by tranche** — one reviewer agent per tranche (multi-agent). Each reviews
   its 2-3 files against project invariants and appends findings to `FEEDBACKS.md` — never
   edits production code.
5. **Preserve user work** — detect uncommitted worktree changes, treat them as user-owned, call
   them out, don't clobber.
6. **Resumable** — because index + ledger are on disk, the sweep can pause and continue across
   sessions; each tranche is an independent unit.

## Fit with spectacular

- The persistent ledger idea is native to spectacular (on-disk truth, resumable work). This is
  arguably a spectacular verb — `spectacular review-sweep` — more than a generic skill.
- Reuses the **Code Review Agent** protocol ([[coding-agents]]) as the per-tranche worker; this
  idea adds the inventory → rank → tranche → ledger orchestration on top.
- Output files (`INDEX_CODEBASE.md`, `FEEDBACKS.md`) — do they live at repo root (as in the
  session) or under `.spectacular/`? Root is more discoverable; `.spectacular/` is more
  consistent with the workspace model.

## Open questions

- Where do the output files live — root or `.spectacular/`?
- Does the ranking heuristic get injected (project declares its risk areas) or does the agent
  infer risk from structure each run?
- Tranche size fixed at 2-3, or scaled to file complexity?
- Multi-agent: parallel tranches (fast, but ledger write contention) or sequential (ordered,
  cheaper to reason about)? A pipeline fits — inventory once, then fan out tranches.
- Overlap with the review agent in [[coding-agents]] — is this just that agent + an
  orchestration wrapper, or a separate skill?

## Promoted to

—
