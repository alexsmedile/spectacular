---
doc-id: build-workflow-doctrine
kind: reference
summary: "The rationale behind build-workflow.md's rules — why each gate exists, the failure modes it prevents, and the relation to bug-workflow. Load only when a routing call feels uncertain or when editing the workflow itself; never needed for routine dispatch."
status: active
---

# Build workflow — doctrine

The *why* behind every rule in [[build-workflow]]. That doc is the runtime core (arc, gates,
tables, procedures); this one holds the reasoning so the core stays lean in the main context.
Read the section matching the step you're unsure about.

## Why the orchestrator bookends everything (the arc)

The orchestrator is the only role that holds the whole request, so it's the only role that can
plan (close an open brief) and mutate (tick, move status). The Builder is deliberately narrow —
closed-brief, build-only, honest-fallback: it bounces to planning rather than freelance a missing
spec, and never writes the ledger. Splitting the verbs this way means no agent ever has to guess
what another decided: **you plan → Builder builds → you confirm + record.**

## A1 — why map-before-plan is a dispatch, and why it's optional

Sometimes the chain won't close not because the *plan* is open but because **you** don't know the
subsystem — you can't name the seams, the sibling pattern, or the blast radius, so any Approach
you wrote would be a guess. That's a knowledge gap, not a planning gap. It earns a `repo-explorer`
dispatch only when the exploration is substantial enough that a fresh window beats reading inline
— for a familiar subsystem the dispatch is pure overhead. And like the Investigator on the fix
side, the Explorer maps but never plans: if it chose the milestone's approach it would be planning
without holding the request, which is the exact failure the role split prevents.

## A2 — why "closed brief or nothing"

A vague brief doesn't fail loudly — it comes back as a *confident wrong build*. The five-slot test
is the same one the debug Fixer uses, at milestone granularity; the difference is size (a
milestone ≈ a small PR, a fix ≈ a patch). The bounce rail exists as a backstop for a brief that's
*wrong* (the codebase moved, the premise broke); using it to catch briefs that are merely *vague*
outsources your quality gate to an agent that can't see the request. That's why the vague-brief
gate is explicitly the orchestrator's, before dispatch.

Name-over-number matching exists because `M<n>` is a convention `doctor lifecycle` warns on but
never blocks — numbering drifts in real requests, and an orchestrator that gets stuck on "M2 vs
M3" stalls on a non-problem. The name text is the stable key.

## A3 — the two axes, and why conflating them is *the* mistake

Decomposability (can you name distinct pieces?) and independence (do the pieces touch?) are
different questions. A change can be decomposable but coupled — pieces exist, but B needs A's
output or they edit the same file — and treating "I can name pieces" as "I can parallelize" is how
fan-outs produce colliding diffs. Monolithic work is built better by one mind holding the whole
thing; slicing it just creates hand-off seams that didn't need to exist. Hence the default is
monolithic-self, and delegation must clear the gates, not be assumed.

### The puzzle-piece principle (why the join is never delegable)

A builder sees only its own piece. It *cannot* invent the interface it shares with a sibling it
never saw — so if you leave the join unstated, each builder guesses, and the guesses don't match.
That's not a build failure; it's a design failure you delegated. Hand two agents "a puzzle piece
each" without the join and you haven't delegated work — you've delegated a design problem. Two
builders picking their own name for the same registry key is the classic mismatch. Designing the
join is cheap at routing time and expensive at reconciliation; that's the entire economics of
doing it in A3. And if the join *can't* be stated cleanly, that's evidence the pieces aren't
disjoint — the module-clarity gate biting early, telling you to serialize.

### Why plan-readiness and module-clarity make parallelism safe

A good plan (each piece closes) + clean pieces with a named join are exactly the preconditions
under which N windows can work blind to each other and still compose. Without them, parallel work
produces locally-green pieces that don't fit, and you pay in D1 the cost you tried to save in C.

## B1 — the economics of the ≥3 gate

Dispatch has a fixed cost: spawn + brief-write + collect + confirm. At 1–2 units that cost roughly
equals the saving; at 3+ concurrency clearly wins. It's deliberately the same threshold as
[[bug-workflow]]'s fan-out gate because the economics are identical.

**Why a single milestone builds inline, even a big one:** milestones within one request are
usually sequentially dependent — M2 needs M1's schema, M3 needs M2's verb. A lone milestone has no
sibling to run concurrently *with*, so the spawn buys nothing but overhead and a context hand-off.
Dispatch earns its keep across **independent** milestones, often across *different* requests — not
one request sliced into pieces.

**Why same-file work serializes:** two Builders editing the same file in separate windows produce
two diffs against the same base; the second clobbers or conflicts. The fix is not coordination
machinery — it's not having concurrent writers. You, building inline sequentially, *are* the
single serial writer; there is no race to solve. The re-verify between ordered steps exists
because an earlier step may have moved lines or changed what the next step builds against — a step
planned against the pre-change file can be wrong.

**Why no branches/worktrees:** branches solve *concurrent* writers; serializing inline means there
are none. `isolation: worktree` per Builder + merge-back is real merge machinery re-solving what
sequential-inline already dissolves — YAGNI until same-request fan-out is measurably too slow,
which it won't be.

## B2 — why decomposition is the visibility mechanism

A fat milestone built as one atomic pass runs for a long time as an opaque block. The platform
spawns subagents fire-and-forget — you can't stream from inside one — so visibility can't come
from watching a Builder mid-run. It comes from *shortening the units and checkpointing between
them*: decomposition turns one long opaque dispatch into several short confirmable ones. That's
the whole mechanism. The nested TASKS bullets are the existing acceptance-checklist convention
([[tasks-rules]]) put to a second use as decomposition checkpoints; they stay out of `x/total` so
they add structure without distorting the ledger. Most small milestones are single-phase — the
`decompose-large-milestone` policy is a `warn` precisely so a coherent one-pass change sails
through without manufactured ceremony.

## C3 — why dispatched work gets a distrust plan

A milestone you built inline, you watched take shape — you already know where it's sound. A
dispatched one, you have only the diff, and a long diff read end-to-end at equal attention is how
a subtle break slips past. Naming 1–2 sharp points *at dispatch time* (before you've seen the
diff) is what makes the return read targeted instead of uniform. A Builder can get its own files
right and still break the join or a caller it never saw — its public surface changed under code it
didn't read. That's why the join and the unseen caller top the list. A `code-reviewer` dispatch
complements the plan for a big/risky diff; it doesn't replace it — the plan is *your* targeted
first look, the reviewer a broad independent one.

## D1 — why reconciliation exists and when it's cheap

"The pieces are disjoint" and "the join holds" were *claims made at routing time*. Reconciliation
is where they're confirmed — checking a contract you already wrote, not discovering one now.
Builder A produced the interface, Builder B consumed it, neither saw the other; two builders
adding to the same registry from disjoint files can each be locally correct and jointly wrong.
A reconciliation failure is never a builder's fault (each met its brief) — the seam between briefs
was the orchestrator's join. This pass is cheap when A3's gates held and expensive when they
didn't, which is exactly why module-clarity is a *gate*: repeatedly hard reconciliation is the
signal the work wasn't disjoint, and the next routing should serialize.

## D2 — why delegated work is confirmed harder

Inline, you were the check. Delegated, the diff is all you have — so scrutiny scales up, the C3
points get read first, and the optional arms-length agents earn their dispatch on risk, not by
default: the reviewer when blast radius is medium+, the verifier when the builder graded its own
work (*the agent that built it shouldn't be the only one to grade it*). The ledger write stays
single-threaded on the main thread for the same reason the debug fleet's `fixes/`/`audit/` writes
do — the labor fans out; the state write never does. The `review → verified` gate is the
highest-leverage review point because it's the only place a *whole-request* diff can be read
coherently — per-milestone passes can't see cross-milestone effects.

## Relation to bug-workflow

Same substrate, opposite direction:

| | [[bug-workflow]] (fix) | build-workflow (build) |
|---|---|---|
| Trigger | a bug / quirk reported | a planned milestone to implement |
| Input | a symptom (unknown cause/site) | a closed milestone chain (known plan) |
| Discover role | `debug-investigator` (open bug → why/where findings) | `repo-explorer` (unfamiliar subsystem → seams/precedent map, *for planning*) — optional |
| Apply role | `debug-fixer` (closed fix → applied) | `spec-builder` (closed milestone → built) |
| Brief size | 5 slots, single-fix | 5 slots, milestone (≈ small PR) |
| Fan-out gate | ≥3 independent closed disjoint-file fixes | ≥3 independent closed disjoint-file milestones |
| The join | rarely — fixes are usually single-site | **the orchestrator's to design** when milestones fan out (A3 → C2 → D1) |
| Bounce meaning | root cause wrong / cross-cutting | brief open / design undecided / spans another M |
| Ledger | `fixes/` + `audit/` (main-thread) | `TASKS.md` tick + lifecycle move (main-thread) |

The Builder is to a milestone what the Fixer is to a fix — the closed-brief, bounce-on-judgment,
never-writes-the-ledger contract, scaled up one granularity.

**Related:** [[build-workflow]] (the runtime core this explains), [[bug-workflow]],
[[bug-workflow-doctrine]], [[spec-builder]], [[tasks-rules]], [[policy-injection]].
