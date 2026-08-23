---
type: Proposal
id: 01a02faa-ed22-736d-94a0-e1596184921f
ref: P12
title: Owner interaction design for decision presentation and anti-slop execution
status: draft
created_by: Alex
created: "2026-08-23T17:28:46Z"
updated: "2026-08-23T17:28:46Z"
scope:
    - v2
target_contract: Contract:019fe381-5d61-7223-b362-03a5f99a7b10
---

# Owner interaction design for decision presentation and anti-slop execution

Exploration for a possible Mission. Nothing here is frozen — this Proposal carries
no execution authority and binds only when a Mission plan freezes its claims.

## Provenance

These four directions were written during P11 exploration on 2026-08-22 and were
not carried into the P11 that shipped. P11's final form settled the *mechanics* of
bounded execution — perimeters, token budgets, the Git boundary, what refuses when.
It produced D12 through D16 and the M15–M22 campaign. The material below settles a
different question: how the owner is asked things, and what stops a worker from
producing bloat once it has been asked.

The original text survives at `b58185d` (unreachable; recover with
`git show b58185d:.spectacular/proposals/P11-context-sandwich-compilation-and-decision-steering-protocol.md`).
A search of `skills/`, `docs/`, `.spectacular/decisions/`, and `.spectacular/campaigns/`
found no implementation of any direction below except a single mention of YAGNI.

## The problem in one line

Spectacular governs what an agent may do, but says almost nothing about how an
agent should ask — so an owner is either grilled on table stakes or silently
overridden on architecture, and a worker with a valid charter can still return
correct, bloated, over-abstracted code.

## Direction 1: The 4-Class Decision Model

Decisions differ in what they cost the owner to answer, so they should differ in
how they are presented. Four classes:

| Class | Content | Presentation | Owner cost |
|---|---|---|---|
| 1. Assumption Baseline | table stakes, sane defaults, standard status codes | consolidated statement — "proceeding with X, Y, Z; any pushback?" | zero; silence locks |
| 2. Structural Fork | irreversible architecture, data model, isolation | 4-part card with recommended default | seconds, high value |
| 3. Taste & Ergonomics | layout, density, CLI verbosity, payload shape | 3 visual specimens to pick between | one glance |
| 4. Operational Invariant | deploy target, budget, SLO | constraint lock card before provisioning | one confirmation |

The 4-part card for Class 2 states: the concrete problem, the technical options
with their consequences, the concrete trade-off, and the recommended default with
its reason.

The rule that makes this worth having is Class 1: a baseline assumption is never
asked as an open question. The stated goal in the source was eliminating
"low-IQ grilling" — never asking whether users should have passwords.

**Why it matters here.** The eval suite already measures `interaction` by owner
question count and by whether settled questions are repeated. This model is the
rule that metric is implicitly testing for, and nothing currently implements it.

## Direction 2: The 3-Anchor Clarity Threshold

Kickoff question count scales with the quality of the owner's input rather than a
fixed quota. Three anchors are checked:

1. **North-Star Outcome** — one sentence of user-observable behavior, plus explicit
   non-goals.
2. **Mechanical Boundary** — input shape → transformation → output shape, plus the
   chosen stack.
3. **Pass Boundary** — one deterministic, failable verification command.

All three present → zero questions, straight to a flight plan. Outcome clear but
one architectural boundary open → one Class 2 card. Core boundaries unstated →
one Class 1 assumption baseline plus one or two Class 2 forks.

**Why it matters here.** `OR-01` asserts `maximum_owner_questions: 0` and `RT-01`
asserts `maximum_owner_questions: 1`. Those thresholds are currently asserted
without a stated rule that decides them.

## Direction 3: Anti-slop execution invariants

Five mechanical rules constraining what a worker may produce inside a valid
charter:

1. **Perimeter constriction** — 2–4 designated paths per slice; editing outside
   triggers `refusal: scope_escape` and rejects the Evidence receipt.
2. **Explicit non-goals injection** — every compiled charter carries a `NON-GOALS`
   block naming the abstractions and frameworks it must not introduce.
3. **Zero-dependency default** — standard library first; a new third-party
   dependency is a Class 2 Structural Fork requiring owner authorization.
4. **Earn your abstraction** — one usage inline, two duplicate, three abstract. No
   `Controller → Service → Manager → Repository → DAO` hierarchies.
5. **Static ceilings** — `SKILL.md` body under 90 lines, charters budgeted near
   1,200 tokens.

Rules 1 and 5 are partly settled by the shipped P11 (perimeters and the 1,200-token
budget with its 1,401–1,440 escalation band). Rules 2, 3, and 4 are not.

## Direction 4: The anti-abdication invariant

The framing that ties the other three together:

> The objective is not fewer decisions, but maximum decision density per second of
> owner attention.

Fast steering must never mean handing architectural control to an agent. Four
rules follow: zero silent assumptions on architecture, security, or data model;
dense option framing so a deep structural call takes seconds; the owner may always
modify rather than pick; every choice lands permanently in
`.spectacular/decisions/`.

This is the reason Direction 1 partitions by presentation rather than by
importance. Reducing owner questions is only safe when the questions removed are
the table-stakes ones.

## What was deliberately left out

The following were in the same exploration and are **not** proposed, because the
shipped P11 supersedes them or they have no engine behind them:

- The Context-Sandwich compilation protocol and its three-layer envelope — settled
  by P11 and D12.
- Worktree isolation, async dispatch, entry gates, and cleanup — settled by P11's
  Git boundary, which explicitly supersedes D15 and D16.
- The 4-layer verification architecture — describes what `test/verify.sh` and
  `test/evals/spectacular/` already do.
- Objective-bound evidence verification — shipped through M20.
- The `[A]/[B]/[C]/[M]/[G]/[F]` single-keystroke steering grammar and the 5-level
  fidelity ladder — interface speculation with no mechanism behind it.

## Open questions for review

- Is decision *presentation* a Spectacular concern at all, or does it belong to the
  Skill rather than the engine? Directions 1 and 2 are behavioral, not mechanical;
  they may be `skills/spectacular/references/` guidance rather than CLI surface.
- Direction 3 rules 2–4 imply a charter field and a refusal path. Does that grow
  the command surface, or is it charter content only?
- Direction 1 Class 3 assumes a visual preview channel that does not exist.
  Drop the class, or narrow it to text specimens?
- Should this Proposal be split? Directions 1, 2, and 4 are one coherent argument
  about asking; Direction 3 is about producing.

## Non-goals

- Do not replace owner authority over outcomes, architecture, or scope.
- Do not make decision records mandatory for trivial implementation details.
- Do not add runtime complexity or external service dependencies.
- Do not introduce a command without explicit owner authorization.
