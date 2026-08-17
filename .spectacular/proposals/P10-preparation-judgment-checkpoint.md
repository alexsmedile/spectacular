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

**2. Does a non-`sufficient` verdict refuse activation, or warn? (decided: split)**
Refusing outright means an owner who knowingly accepts a risk cannot proceed without
editing the plan to lie.

`needs-decision` refuses: a missing owner decision is a genuine stop, and the correction is
available — make the decision. `needs-evidence` reports a notice: gathering that evidence
is frequently the Mission's own first Objective, so refusing would block the work that
resolves it.

**3. Do the two specialist roles become real agents? (decided: later)**
Not in this Mission. The two checks are written as Skill guidance first. If the judgment
turns out to be hard enough to need its own context, extraction is cheap then; an unused
agent definition rots now. This repository has no `.claude/agents/` directory at all, so
nothing is being removed by deferring.

**4. Does this Proposal also fix the `independent` review trigger? (decided: yes — already
written)**
The same defect: three files mention the level, none said when to choose it.
`prepare.md:44` asked the owner to pick from `automatic | clustered | independent` with no
criteria; `close.md:30` described what an earned review costs; `close.md:44` said freshness
is not independence without saying what is.

Both are now written — the trigger list where the level is chosen, and what actually earns
independence where the review is recorded. Kept in this Proposal's scope rather than split
out, because splitting would leave a known gap open for no benefit.

**5. How should Skill growth be governed? (decided — recorded here for the Mission)**
M11 froze a pass boundary reading "The Skill does not grow", and the Skill grew 33 lines.
The claim was passed with the delta logged, which is the correct outcome for the wrong
reason: the boundary was a proxy for "is the Skill getting bloated", stated as an absolute
that the real work immediately contradicted.

The owner's decision is that a hard line-count stop is the wrong instrument. Growth is
sometimes correct — guidance that does not exist cannot be followed, and three hollow seams
are the direct cost of guidance that was never written.

Replace the stop with a **reported delta and a review judgment**. A Mission touching
`skills/` records lines added and removed per file, and the reviewer answers one question:
*was the growth worth what it bought?* A Mission may grow the Skill by 200 lines and pass
if the guidance is load-bearing, and may grow it by 10 and fail if it is restatement.

This is the `contract-drift` pattern applied to Skill size: report the fact, let judgment
handle it, refuse nothing. It also means the boundary can actually fail, which "does not
grow" could not — that one was violated in letter and upheld in spirit at the same time.

Still worth preferring displacement where it is genuinely smaller: `prepare.md` is 96 lines
and its `Plan` section already weighs seven slice criteria and already asks for the
sufficiency verdict. That is the right content in the wrong frame, so rewriting it is
probably cheaper than appending beside it — as a drafting preference, not a limit.

The general form of this defect — a `pass_boundary` that freezes a countable proxy instead
of the concern it stands for — is recorded in `FEEDBACKS.md` under "A boundary was written
as a hard stop when it was a proxy for a concern". It may want a second boundary kind, a
hardness marker, or preparation-time detection. That is a separate Proposal; this one only
decides the Skill-size case.

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
