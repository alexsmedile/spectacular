---
type: Proposal
id: 01a00a93-4757-7547-b64e-e91d2c291ce4
title: Reduce repeated cost and brittle paths
status: draft
human_ref: P5
created_by: Alex
created: "2026-08-16T12:36:59Z"
updated: "2026-08-16T12:36:59Z"
scope:
    - v2
target_contract: Contract:01a00a20-63dd-7670-97f1-9eb8e12adc3a
---

# Reduce repeated cost and brittle paths

Exploration for a possible Mission 7. Nothing here is frozen. The problem framing is
firmer than the proposed shapes, and the shapes are firmer than the schema sketches.
Anything below may be dropped, split, or reversed at plan-freeze.

## Where this came from

A review of the v2 method against the efficiency dynamics of natural systems —
foraging trails, slime-mould transport networks, fractal branching, flocking. The
useful part of that comparison was not the metaphors. It was that biology and this
method are solving the same two problems: how to avoid paying the same cost twice,
and how to survive a chosen path turning out to be wrong.

The comparison also confirmed that several v2 properties are already at the right
answer, which is worth recording so a future Mission does not "improve" them:

- Inline-until-promoted Objectives and Runs. Structure is only paid for when flow
  through it justifies the wall. M6 keeps three inline Objectives; M2 promoted four
  because delegation earned them. That ratio is responding to the correct pressure.
- One file at Mission start, and no required Mission-local `index.md`. Minimum
  perimeter around the governed volume.
- Coordination through the workspace rather than between agents. Reading and writing
  `MISSION.md` instead of messaging is why cold recovery works at all, and it is why
  coordination cost does not grow quadratically with the number of sessions.
- The batched tree-bound gate, and reading detailed logs only on failure.

None of those need work. P3 and P4 already took the obvious remaining slack out of
context loading and control ceremony. What follows is what those two did not reach.

## Problem

**We pay full discovery cost every session.** SKILL.md asks every launch to derive
thirteen preflight dimensions — workspace, branch, cleanliness, upstream freshness,
release state, Mission, bindings, owner, activation fingerprint, validation mode,
current Objective and Run, dependencies, Gaps, stops. It derives all thirteen
uniformly, whether or not anything moved since the last session ended twenty minutes
ago. P3's progressive disclosure fixed how *deep* a session reads. It did not touch
how *often* a session re-derives the same unchanged fact. Those are different wastes
and only one of them has been addressed.

The workspace already hints that this was felt: `.spectacular/.last-mutation` persists
four fields across sessions. It records what changed, but not what was observed, so it
cannot be used to skip anything.

**Repair is a cliff.** The current loop takes the narrowest justified repair, consumes
budget, and refuses completion on exhausted repair. One path, monotonically depleting,
terminating in an escalation that hands the owner a stop rather than a choice. The
cost is compounded by something that happens earlier: planning compares genuinely
different, outcome-sized approaches, freezes one, and discards the rest. When repair
exhausts, those alternatives are gone, and the budget has been spent to learn exactly
one fact — that the chosen approach fails for a specific reason. That fact is often
precisely what would have selected a different approach at the start.

**Audits are aimed by memory.** Audit is correctly bounded to one named claim or
scope. But a human names the target, and humans name what they remember rather than
what has actually drifted. The signal needed to aim better is already in the bundle
and is currently unread: repair counts on the Run, evidence age and freshness,
claim verdict state, contract fingerprint age. M6's Run carries `repairs: 1` and
nothing consumes it.

**Serialization, not fan-out rules, is what blocks parallel work.** The fan-out rule
requires disjoint claim scopes, which is correct and should stay — it is what prevents
two runners returning contradictory evidence for the same claim. But on M6 disjointness
was satisfiable and parallelism still never became available, because the three
Objectives form a strict `after:` chain. The binding constraint is dependency topology.

M6's O3 shows the cost concretely. It bundles five separable proof activities —
atomicity, safe refusals, representation equivalence, legacy readability, distribution
behavior — into one Objective claiming one thing, gated behind all implementation.
Proof work is among the most parallelizable work available, and it has been serialized
to the end. The reason is that `after:` conflates two different dependencies: needing
the produced artifact, and needing only the interface shape. Golden fixtures and
table-driven negative tests need the second, which is frozen at plan time.

## Directions worth exploring

Four, listed cheapest-first. They are independent; adopting one does not require any
other.

> **Status as of 2026-08-18 (v2.4.0).** Three of the four shipped and one did not, which is
> why this Proposal stays live rather than being retired under `D11-proposal-retirement`.
>
> | Direction | State |
> |---|---|
> | A trace of the last preflight, with decay | **open** — nothing persists a trace |
> | Aim audits at drift instead of at memory | shipped — `mission check` reports per-claim drift and `audit` defaults to a claim |
> | Keep rejected approaches as thin fallbacks | shipped — `fallbacks:` with `invalidated_if:` validates |
> | Separate "needs the artifact" from "needs the interface" | shipped — `after_interface:` is a distinct edge kind with cycle detection |
>
> The open direction is the cheapest of the four and is unblocked. A Proposal is absorbed when
> the question it asked has been answered, not when most of it has.

### A trace of the last preflight, with decay

Persist what the last preflight observed, alongside hashes of its inputs. On the next
launch, recompute only the genuinely volatile inputs — git HEAD and dirty state are
cheap and should always be checked — and if those still match the trace, report the
traced outcome and cite the trace as the evidence line rather than re-deriving the
remaining dimensions.

The decay is the part that must not be dropped. A cache without expiry locks the
session onto a stale picture, because things change outside the hashed inputs: upstream
moves, a release ships, a dependency closes. The trace should be trusted only while its
hashes match *and* it is inside a freshness horizon. This mechanism already exists in
the workspace — `freshness_valid_until` on the Anchor — and simply has never been
pointed at the path that is walked most often.

Sketch, not a schema:

```yaml
# .spectacular/.preflight-trace
traced_at: "2026-08-16T13:54:02Z"
mission: M6
head: e0e2fdcc6874adb0750b9e07cae43d6f09febc6d
dirty: false
contract_fp: sha256:80336b15...
activation_fp: sha256:2bb4d9ff...
workspace_fp: d8b24fe7cfef...
validation: mission.v2
outcome: clean
```

Open: what the horizon should be, whether it is fixed or configured in
`workspace.yaml`, and whether the trace belongs beside `.last-mutation` or replaces
it by absorbing its fields.

### Aim audits at drift instead of at memory

Have `mission check` compute and report a per-claim drift signal from inputs the
bundle already carries — repairs consumed, evidence age and freshness, verdict state,
fingerprint age. Let `audit` default to the highest-signal claim when the owner names
none. The bounded-scope discipline is unchanged; only target selection improves.

This needs no schema change, since every input is already present and the score is
derived. It is the cheapest of the four and can ship alone.

Open: whether the signal is a single ranked score or a small set of named flags. A
score is easier to default on; flags are easier for an owner to argue with, which may
matter more.

### Keep rejected approaches as thin fallbacks

At plan-freeze, record each seriously-considered rejected approach in three parts:
what it was, why it lost, and — the load-bearing part — what observation would reverse
that rejection.

```yaml
fallbacks:
  - approach: Keep two package roots, one per representation
    rejected_because: Doubles the decode surface and the migration cost
    invalidated_if: A single decoder cannot preserve v2 readability without field loss
```

This costs three lines and no maintenance. It is not a second implementation kept
alive; it is a note about where a different path was. Its value appears exactly at
repair exhaustion: instead of escalating "budget gone, stuck," the owner gate can
report that the failure which consumed the budget matches a recorded `invalidated_if`,
and recommend re-freezing onto that approach. The owner receives a decision rather than
a dead end.

This also gives FROST's Operability dimension a structural answer to "the chosen
approach was wrong," which it currently lacks.

Note that this necessarily sits inside the frozen envelope and the activation
fingerprint. If fallbacks were mutable, an operator could invent one mid-Run to escape
a stop, which is the exact failure the freeze exists to prevent. Freezing them at plan
time is the point, not an implementation detail.

Open: whether an unmatched `invalidated_if` at exhaustion should still surface the
fallbacks to the owner as context, or stay silent to avoid implying false options.

### Separate "needs the artifact" from "needs the interface"

Split the single `after:` edge into two kinds. `after:` keeps its current meaning: this
Objective needs the produced artifact and is genuinely sequential. A second edge —
`after_interface:` is a working name, not a proposal about naming — means the Objective
needs only the contract shape, which exists the moment the plan freezes, and is
therefore startable at activation.

Under that split, M6's O3 decomposes into proof Objectives that begin at activation and
land as implementation lands, rather than queueing behind it. They remain disjoint by
claim, so the existing fan-out rule finally has something to apply to instead of being
moot.

This is the only one of the four that changes validator semantics. `validate.go`
currently validates one dependency graph and would need to validate two edge types,
including rejecting an interface-only dependency on something whose interface is not
in fact frozen. That is real work and a real source of new refusal modes. It may
deserve to be its own Mission rather than an Objective inside M7 — that call belongs
to plan-freeze, not to this document.

## Alternatives considered and set aside

**Cache the whole preflight without expiry.** Simpler and faster, and wrong for the
reason foraging trails evaporate: without decay the session locks onto an obsolete
observation and cannot detect that it has. Rejected on correctness, not cost.

**Drop the disjoint-claim requirement to widen fan-out.** Tempting, since flocking
achieves coordination with heavily overlapping awareness and no disjointness at all.
Rejected: overlapping *awareness* is not overlapping *authority*. Two runners owning
one claim can return contradictory evidence with no rule to arbitrate. The dependency
split gets the parallelism without touching the guarantee.

**Maintain runner-up approaches as live branches.** The strongest version of the
fallback idea, and the closest to how slime moulds actually keep redundant loops open.
Rejected as far too expensive here: real branches need rebasing, checks, and attention,
and the protoplasm analogy breaks because a stale branch decays silently while a
biological tube does not. The three-line note captures most of the value at
approximately none of the cost.

**Make Autopilot decide when to escalate.** Considered and dropped quickly. Nature has
no owner gate because nature has no accountability — a colony that optimizes into a
failure mode simply dies. The gates and frozen fingerprints are what make this system
correctable, and every direction above is deliberately inside the frozen envelope. Any
change that blurs who decides, or what counts as proof, is the wrong trade regardless
of how much waste it removes.

## Open questions for the owner

- Is M7 the right container for all four, or should the dependency split be its own
  Mission against the same Contract? The first three are additive and do not change
  proof or authority semantics; the fourth changes what the validator must prove.
- Should the first two land in the current worktree, given they touch no frozen
  semantics, or wait for a clean M7 activation?
- What freshness horizon is acceptable for a preflight trace before a full re-scan is
  forced?
- Does `fallbacks` belong in the activation fingerprint from the start, or should a
  first version be advisory and unfingerprinted to see whether it is used at all?

## Relationship to current work

M6 completed with independent review and `repairs: 1`, so it does not absorb any of this. Its frozen
schema and completion boundary stay as activated; anything here applies to a later
Mission. P3 and P4 are accepted and already reflected in the Skill — this explores what
remains after both, and does not revisit either.
