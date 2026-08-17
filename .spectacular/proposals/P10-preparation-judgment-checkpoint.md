---
type: Proposal
id: 01a0107d-f32c-787f-962d-ed77d7338ecb
ref: P10
title: Preparation judgment checkpoint
status: draft
created_by: Alex
created: "2026-08-17T16:11:24Z"
updated: "2026-08-17T16:11:24Z"
scope:
    - v2
target_contract: Contract:01a00a20-63dd-7670-97f1-9eb8e12adc3a
---

# Preparation judgment checkpoint

Exploration for a possible Mission. Nothing here is frozen. The problem is demonstrated;
the mechanism is a direction. Anything below may be dropped, split, or reversed at
plan-freeze.

## The problem in one line

A Mission can be frozen and activated without anyone having asked whether the approach is
understood well enough to start, or whether the slice is the right size — and the two
verdicts that would answer those questions are already half-present in the workflow, one
of them unexplained and the other absent.

## Where this came from

An adversarial review, H16, was run on 2026-08-09 against the v1-era foundation at
`5b5a738`. It audited four SDLC scenarios — a fuzzy feature through maintenance, a
regulated Mission with a false assumption, concurrent Missions across several actors, and
a production incident through observed recovery — and found no missing foundation. Its
conclusion was that the accepted contracts already provided authority, stop, evidence, and
resume routes for all four without imposing fixed SDLC phases.

That part is history and needs nothing. The review's durable output was narrower: two
practical capabilities it argued were *already implied* by accepted commitments but not
implemented anywhere.

The review's return was never accepted. Its `central_disposition` reads `pending`, its
`next_action` names a central orchestration that no longer exists, and the branch it ran
on was superseded by the v2 clean break. It survived as an untracked file in an abandoned
worktree and was found during worktree cleanup.

**So this Proposal does not adopt H16.** It re-derives the two capabilities against the
v2 surface as it stands today, and drops everything else — the eleven v1 contracts, the
S10–S12 session plan, the specialist slate as products, the return schema.

## What v2 actually has today

Both verdicts H16 proposed already have a partial footprint. That is the finding that makes
this worth doing rather than an idea from an old document.

**Design sufficiency: the verdict exists, unexplained.**
`skills/spectacular/references/prepare.md:32` reads:

```text
Record one verdict: `sufficient | needs-evidence | needs-decision`.
```

Three states, no definition of any of them, no statement of what a non-`sufficient`
verdict blocks. Nothing consumes it, nothing stores it, and no record carries it. An agent
reading that line records a word.

**Slice quality: absent.**
The lines above it (`prepare.md:23-30`) ask the agent to compare approaches on observable
result, coherence, dependencies, reversibility, learning value, integration path, and
cancellation state — seven of the exact criteria H16 assigned to slice comparison. But it
frames them as *choosing between approaches*, never as *judging whether the chosen one is
the right size*. There is no `coherent | too-broad | fragmented | dependency-bound`
verdict, no requirement to generate more than one candidate, and no check that a Mission's
Objectives form one coherent outcome.

**Review level: a field with no trigger.**
`prepare.md:44` offers `automatic | clustered | independent`, "defaulted once when
shared". Nothing anywhere says when a Mission earns `independent`. The only other mention
in the Skill is `close.md:30`, which describes what an earned review costs, not what earns
it. H16 listed twelve triggers; v2 has zero.

So the gap is not conceptual. Three seams exist and two are hollow.

## Why this matters beyond tidiness

The evidence is in this repository's own governance history.

**M10 was activated and then archived without executing** because its design was wrong —
amendment had been coupled to Mission completion, which made the timing arbitrary, the
mechanism one-shot, and Contract typos uncorrectable. That was caught by the owner asking
"why should completion amend?", after the Mission was frozen and activated. A design
sufficiency check is exactly the question that was missing, and it cost a Decision (`D9`),
an archived Mission, and a re-plan.

**M11 nearly shipped a silent data loss.** `mission start` dropped `resolves_gaps:`
because the plan struct did not carry the field. Unit tests were green. It was caught only
by driving the real CLI with an invented Gap ref. That is an evidence-adequacy question — a
sufficiency reviewer's stated concern — and nothing asked it.

**M9's Gap sweep was invented mid-Mission.** Gaps appear in the Skill eleven times, always
as blockers to read, never as things to resolve. The Mission needed a sweep, so it made one
up. A slice-quality check asking "is this one coherent outcome?" is where that surfaces.

Three Missions, three distinct preparation failures. None of them were implementation bugs.

## The direction

Two proportional judgment checks inside guided `define`, before the owner is asked to
activate:

```text
accepted direction
  → candidate outcome
  → candidate slices (two or three, when the work is not clearly small)
  → design sufficiency: sufficient | needs-evidence | needs-decision
  → slice quality:      coherent | too-broad | fragmented | dependency-bound
  → both permit continuation
  → owner activation
```

Held deliberately:

- **Not a new lifecycle phase.** No `DESIGN.md`, no Design record type, no fourth
  ceremony. This repository has five empty record directories already; it does not need a
  sixth type.
- **Not a new command.** The surface is twelve and `AGENTS.md` states that growth past
  twelve is a stop. If a mechanical check is wanted, it belongs inside `mission start`,
  which already validates the frozen envelope.
- **Proportional.** A clearly reversible local change compresses both checks to a
  sentence. Consequential work deepens them. The cost has to scale with the risk or it
  becomes ceremony people route around — the failure mode `P3` and `P4` were both about.
- **Skill-side judgment, CLI-side presence.** Whether a design is sufficient is semantic
  and belongs to the Skill. Whether a Mission that claimed a verdict actually carries one
  is mechanical.

## Open questions

These need owner decisions and are the reason this is a Proposal.

**1. Are the verdicts frozen in the Mission, or do they stay in chat?**
Freezing them makes them auditable and lets `mission check` report them — and the
`resolves_gaps:` precedent from M11 shows how a frozen semantic field behaves. It also
grows the fingerprinted envelope and adds two required fields to every plan.
*Recommendation: frozen, because a verdict nobody can check later is the state we already
have.*

**2. Does a non-`sufficient` verdict refuse activation, or warn?**
Refusing makes the check real but means an owner who knowingly accepts a
`needs-evidence` risk cannot proceed without editing the plan to lie. Warning preserves
owner authority and matches the `contract-drift` notice precedent.
*Recommendation: refuse on `needs-decision`, notice on `needs-evidence`. A missing owner
decision is a genuine stop; missing evidence is often the Mission's own first Objective.*

**3. Do the two specialist roles become real agents?**
H16 proposed a Mission Slice Advisor and a Design Sufficiency Reviewer. This repository has
no `.claude/agents/` directory at all, so this would be new machinery.
*Recommendation: no, not initially. Write the two checks as Skill guidance first and see
whether the judgment is actually hard enough to need a separate context. Extraction is
cheap later; an unused agent definition is not.*

**4. Does this Proposal also fix the `independent` review trigger, or is that separate?**
It is the same class of defect — a field with no guidance — and H16 supplied a trigger
list. But it is a distinct surface and could be its own slice.
*Recommendation: same Mission, separate Objective. It is the third hollow seam and
splitting it leaves a known gap open for no benefit — which is precisely the slice-quality
question this Proposal is about, applied to itself.*

**5. Is `prepare.md` the right home, given it must not grow?**
M11's `workflow-states-the-step` claim was not met because the Skill grew 33 lines against
a "does not grow" boundary. This work is larger than that. Either the boundary needs
restating, or this content displaces existing text rather than adding to it.
*Recommendation: name a net-line budget in the Mission plan and treat displacement as the
default. `prepare.md` is 96 lines and its `Plan` section already gestures at both checks;
rewriting that section is likely smaller than appending to it.*

## What this does not touch

- The four SDLC scenarios H16 audited. That conclusion was about the v1 foundation and is
  not being re-accepted.
- Any accepted Contract's stated behavior. If this changes what `mission start` requires,
  that is a `contract_version:` bump on `CC-missioncli`, not an amendment.
- The Handoff freeze in `P9`. Independent, and neither blocks the other.

## Provenance

Derived from `H16-sdlc-coherence-adversarial-review.md`, a handoff return found untracked
in an abandoned Codex worktree at `5b5a738` during worktree cleanup. Its
`central_disposition` was `pending` and remains so — this Proposal supersedes the need to
dispose of it. The source document is not a v2 record type and is not being imported;
`P9` would give future returns a real home.
