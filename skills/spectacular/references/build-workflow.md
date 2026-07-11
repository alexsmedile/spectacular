---
doc-id: build-workflow
kind: reference
summary: "How the skill implements planned work: frame a closed milestone brief from the request chain and route its shape (monolith / serial / fan-out), decide inline-vs-dispatch and unit size, execute (build or dispatch spec-builders, planning what to distrust on return), then reconcile parallel outputs and record. Four phases: Frame → Decide → Execute → Reconcile+record. The build-direction analog of bug-workflow.md."
status: active
---

# Build workflow — implement a milestone, dispatch when it's worth it

Loaded when actively implementing a request that's in `active` status — turning a PLAN/TASKS
milestone into shipped code. Governs how the skill **frames** a closed brief from the request chain,
**routes** the work's shape, **decides** inline-vs-dispatch, **executes** (building or dispatching a
`spec-builder`), and **reconciles + records** the result. This is the **build direction**; its mirror
is [[bug-workflow]] (the fix direction). Same shape: closed-brief test, bounce rail, single-threaded
ledger.

## The arc — read this first (orchestrator's whole job)

**You are the orchestrator: the only role that holds the whole request and the only hand that
mutates.** You frame the brief, route the work, decide dispatch, spawn builders, confirm diffs, tick
the `TASKS.md` checkbox, and move lifecycle state. The `spec-builder` is narrow delegable labor
(closed-brief, build-only); the two things you never delegate are **planning** (closing an open brief)
and **mutation** (the checkbox tick + status move). Verbs, no overlap:
**you plan → Builder builds → you confirm + record.**

The four phases, in order — each phase's sub-steps are the mechanics:

```
A  FRAME THE WORK      — can it close, and what shape is it?
   A1  Map unfamiliar?    subsystem unknown → dispatch repo-explorer → then plan          (optional)
   A2  Chain closes?      walk task-row → M-block → PLAN §2/§3/§6/§7 → 5 slots fill? → else bounce-to-planning
   A3  Route the shape    monolith → self · coupled → serialize · disjoint → fan out — AND design the join now

B  DECIDE WHO & HOW BIG — inline vs dispatch, one pass vs phased
   B1  Worth-it gate      dispatch only if ≥3 independent closed disjoint-file milestones; else build inline
   B2  Size the unit      multi-phase milestone → nested sub-steps, one at a time, checkpoint between      (optional)

C  EXECUTE              — build it; plan what to distrust
   C1  Build inline       apply the milestone yourself, verify → straight to D
   C2  Dispatch           spawn spec-builder(s) with closed brief(s) — each brief carries its half of the join
   C3  Critical-check plan (dispatched only) name the 1–2 things you'll verify first on return

D  RECONCILE + RECORD   — do the pieces fit, then tick
   D1  Reconcile          (parallel disjoint fan-out only) does the join hold? seams fit, contract matches → else stitch
   D2  Confirm + record   confirm diff (C3 checks first; delegated ⇒ harder) · run Success criteria · [opt. reviewer/verifier] · tick · lifecycle move
```

> **@Implementation policy gate.** First, run `spectacular policy @Implementation` and follow every
> active policy returned — chiefly `understand-before-change` (the request's Understanding slot must
> be filled before you build — A2 *is* that understanding for the milestone), `build-order` (P11 —
> build the milestone the plan ordered, not a later one), and `decompose-large-milestone` (B2 — a
> multi-phase milestone is built as sequential sub-steps, not one opaque pass). All are `warn`. See
> [[policy-injection]].

**The one idea that binds the phases: the join is yours.** When you split work across builders, *how
the pieces fit together is a design decision — and it's the orchestrator's, never a builder's.* A
builder sees only its own piece; it cannot invent the interface it shares with a sibling it never saw.
So the shared contract (the signature, the schema, the registry key, the shape of what's passed) is
**designed up front in A3, written into every brief in C2, and verified held in D1.** Hand two agents "a
puzzle piece each" without specifying the join, and you haven't delegated the work — you've delegated a
design problem, and it comes back as pieces that don't fit. This through-line is what makes fan-out
safe; watch for it flagged as **[join]** in the phases below.

---

# Phase A — Frame the work

*Can this milestone become a closed brief, and what shape is the work?* Phase A produces two things: a
**closed brief** (A2) and a **routing decision** (A3). Nothing dispatches until both exist.

## A1 — Map unfamiliar ground first (optional — dispatch `repo-explorer`)

Sometimes the chain won't close in A2 not because the *plan* is open, but because **you don't yet
understand the subsystem well enough to write the Approach**. You can't name the seams a milestone
touches, the sibling pattern it should mirror, or its blast radius — so any brief you'd write would be
a guess. That's not a planning gap in the request; it's a *knowledge* gap in you.

When that happens, and the subsystem is large enough that reading it inline would crowd the main
context window, **dispatch [[repo-explorer]]** — the build-side mirror of [[bug-workflow]]'s
Investigator. Give it a scoped question ("how does the CLI register a doctor area? what would a new
`spectacular <verb>` touch?"); it reads in its own window and returns a **structured map**: entry
points, the precedent to mirror, the integration seams (`file:line`), and blast radius. You then plan
A2 against that map.

This step is **optional and judgment-gated** — skip it when you already know the ground (most
milestones in a familiar request). It earns its dispatch only when the exploration is substantial
enough that a fresh window beats reading inline. Like the Investigator, `repo-explorer` **maps, it
doesn't plan**: it names the seams and the precedent; *you* decide what the milestone is and write the
Approach. The map informs the brief; it isn't the brief.

## A2 — Does the chain close into a brief? (the context-assembly chain)

**Before dispatching anything, confirm the milestone can become a *closed* brief.** A bare `- [ ]`
TASKS line is a fragment, not a brief. Walk the chain that already exists in every request folder —
you're not inventing a context format, you're reading the one the request lifecycle already wrote:

```
TASK ROW            →  MILESTONE BLOCK      →  PLAN SECTIONS         →  BRIEF                 →  OUTPUT + CHECK
TASKS.md             TASKS.md's own          PLAN.md §3 (this M's     synthesized from all:    implementation +
one `- [ ]` line      `### M<n> — ...`        line) + §2 Constraints   Goal / Constraints /     Success-criteria
                      header + siblings       + §6 Validation (this    Approach / Expected       run
                                              M's check) + §7 Deliv.    output / Success crit.
```

The five brief slots map straight onto the chain (see [[spec-builder]] for the slot contract):

| Brief slot | Sourced from |
|---|---|
| **Goal** | PLAN §3 — this milestone's line |
| **Constraints** | PLAN §2 (applies to every milestone) + inherited STACK/PRINCIPLES |
| **Approach** | the ordered `- [ ]` steps in this milestone's `### M<n>` TASKS block |
| **Expected output** | PLAN §7 — the deliverable this milestone contributes |
| **Success criteria** | PLAN §6 — this milestone's validation line (reuse verbatim; don't re-derive) |

**Match by name, not by number.** The TASKS↔PLAN milestone link is an `M<n>` *convention*, not an
enforced ID — `doctor lifecycle` flags drift but never blocks, so numbers can disagree while names
agree. When they disagree, **match on the milestone's name text** (the words after the em-dash), not
the number. A renumbered-but-same-named milestone is the same milestone; don't get stuck on `M2` vs
`M3`. (`spectacular doctor lifecycle` will have already surfaced any label drift as an advisory
warning — read it, then name-match through it.)

**If the chain won't close — bounce to planning (yourself).** The milestone isn't dispatchable when:
- PLAN §6's validation line for it is vague ("works correctly") — there's no runnable Success criteria.
- PLAN §3's line or the TASKS block implies a **design decision not yet made** (which format? which
  default?) — filling it is planning, not building.
- The `### M<n>` TASKS block is still the unfilled boilerplate template, or has no real step list.
- §3/§6/§7 have no matching entry for this milestone at all (the plan is incomplete).

This is the same "closed brief or nothing" test the Fixer uses, at milestone granularity. The
difference from a bug: a Fixer's brief is 5 slots because a fix is "make one thing true again"; a
milestone can span multiple files/verbs/docs — closer to a small PR than a patch. **When the chain
won't close, you finish the planning first** (grill/refine the PLAN, decide the open question, fill
§6's check) — *then* the milestone is dispatchable. Don't hand a Builder an open brief and lean on
its bounce to catch it: the bounce is the backstop for a *wrong* brief, not the quality gate for a
*vague* one — that gate is you, here.

## A3 — Route the shape: monolith · serial · fan-out (and design the join)

The brief closes — now **route the shape of the work** before you ask *who* builds it (Phase B). This
is the orchestrator's core maneuver; Phase B is the mechanics of the choice you make here.

**The default is monolithic-self.** A hard, coherent, tightly-coupled change is *built better by one
mind holding the whole thing* than by slicing it across windows that each see a fragment. Delegation is
not free and not the goal — it's the exception you reach for **only when the shape earns it.** When in
doubt, build it yourself. Decomposition and fan-out must clear a bar, not be assumed.

### Two axes decide the shape

Route on two independent questions — they are *not* the same axis, and conflating them is the mistake:

```
              DECOMPOSABLE?  ── one hard thing, or separable pieces?
                │
    ┌───────────┴───────────┐
  one hard,                many
  coupled thing            separable pieces
    │                        │
  MONOLITHIC-SELF          INDEPENDENT?  ── do the pieces touch each other?
  (build inline, C1)         │
                 ┌───────────┴───────────┐
              coupled /                disjoint,
              order-dependent          self-contained
                │                        │
              SERIALIZE-SELF           DECOMPOSE-AND-DELEGATE
              (inline, in order —        │  [join] design the shared contract now
               B1 same-file rule,      worth fanning out? (B1 economics)
               or B2 if one fat          │
               milestone, many        ┌──┴───────────────┐
               phases)               FAN OUT           PARALLEL + RECONCILE
                                     (N builders, C2)   (disjoint files → D1 verifies the join)
```

1. **Decomposability** — *one hard monolithic thing, or many separable pieces?* A change is monolithic
   when its parts share so much state/logic that splitting them just creates hand-off seams that didn't
   need to exist. It's decomposable when you can name **distinct pieces that each stand on their own** —
   a piece someone could build knowing only *its* brief plus the shared contract, not the whole.
   - Monolithic → **build it yourself** (C1). Don't manufacture pieces; that's ceremony (the
     `decompose-large-milestone` warn exists to let a genuine one-pass change sail through).
   - Decomposable → go to axis 2.

2. **Independence** — *do the pieces touch each other?* Given separable pieces, are they **coupled**
   (piece B needs piece A's output, or they edit the same file) or **disjoint** (each owns its own
   files and needs nothing from a sibling mid-build)?
   - Coupled / order-dependent → **serialize, yourself, in order** (B1's same-file rule, or B2 if it's
     one fat milestone of many phases). One serial writer, no race.
   - Disjoint & self-contained → **decompose-and-delegate** — the only shape that fans out or runs in
     parallel. Whether it's *worth* fanning out is B1's economics (the decision table owns the
     threshold). When disjoint pieces run in parallel, they need a **reconciliation pass** afterward
     (D1).

### [join] Design the join before you cut the pieces

The moment you choose decompose-and-delegate, **one design job is now yours and stays yours: how the
pieces fit together.** This is the puzzle-piece principle made concrete. A builder sees only its own
piece — it cannot invent the interface it shares with a sibling it never saw, so *if you leave the join
unstated, each builder guesses, and the guesses don't match.* That's not a build failure; it's a design
failure you delegated.

So before writing any brief, **name the shared contract explicitly** — whatever the pieces must agree
on to compose:
- **The signature / shape** — the function args, the return type, the data shape one piece produces and
  another consumes.
- **The name / key** — the exact identifier, registry key, config field, or route both sides use. Two
  builders picking their own name for the same thing is the classic mismatch.
- **The ordering / protocol** — who registers first, what state exists when, which piece owns the
  default.

Write it down once, as the orchestrator's decision. It becomes the **shared half of every brief** in C2
(each builder gets its own Goal *plus* the same contract), and the thing D1 verifies held. **If you
can't state the join cleanly, the pieces aren't as disjoint as you think — serialize them instead**
(that's the module-clarity gate in the table below biting early). Designing the join is cheap here and
expensive at D1; do it here.

### The gates on decompose-and-delegate

Even when the two axes point at delegation, three readiness gates must all hold — miss any and you
**build inline instead**:

| Gate | The question | If it fails |
|---|---|---|
| **Benefit of delegation** | Does a fresh window build this *better* than the main thread, and is there ≥1 sibling to run *with*? | A lone piece / a trivial piece → no concurrency to win → self-serve. |
| **Plan readiness** | Is each piece a **closed brief** (A2's five slots fill, no undecided design)? | Open brief → finish the planning first; delegating an open brief just bounces. |
| **Module clarity** | Are the boundaries **clean** — each piece owns its files, and **can you state the shared join** ([join] above)? | Fuzzy seams / an un-nameable contract → serialize instead; you'll spend more reconciling than you saved. |

**Plan readiness and module clarity are the two that make parallelism *safe*.** A good plan (readiness)
+ clear pieces with a named join (module clarity) are exactly what let disjoint work run in parallel and
reconcile cleanly. Without them, parallel work produces pieces that don't fit, and you pay the cost you
tried to save.

### Where each shape routes

| Shape | Route to |
|---|---|
| One hard monolithic thing | **C1** (build inline) |
| One fat milestone of several verify-point phases | **B2** (sequential sub-steps, checkpoint between) |
| Several coupled / ordered milestones | **B1** *same-file / ordered* rule (serialize inline, in order) |
| Disjoint pieces, too few to fan out | **C1** (self-serve — B1's economics say fan-out isn't worth it yet) |
| Disjoint pieces, enough to fan out | **B1 → C2** (fan out N builders, each brief carrying the join) |
| Disjoint pieces run in parallel | as above **+ D1** (reconcile the parallel outputs) |

Phase B owns the *mechanics* (the fan-out threshold, the same-file serialize rule, the phase-sizing) —
A3 is the doctrine that chooses among them. Read B next through the lens of the shape you just picked.

---

# Phase B — Decide who & how big

*Given the shape from A3: does the main thread build it or a dispatched Builder (B1), and is each unit
one pass or several phased sub-steps (B2)?*

## B1 — Build inline, or dispatch? (the worth-it gate)

One question: do *you* (the main thread) build it, or dispatch a `spec-builder`? Dispatch has a fixed
cost (spawn + brief-write + collect + confirm). It pays back only when the milestone is **substantial
enough that a fresh window builds it better than crowding the main thread** — and when you have **more
than one** to run.

### Decision table — self vs dispatch

| Situation in hand | Orchestrator builds it **itself** | **Dispatch** to spec-builder |
|---|---|---|
| A 1–2 line / trivial milestone | ✓ — a spawn costs more than the edit | |
| A single substantial milestone | ✓ — one milestone, one context; you already hold the plan | |
| **≥3 independent, closed, disjoint-file milestones** | | ✓ — concurrency wins; each Builder self-verifies |
| Milestones that share files / are order-dependent | ✓ — sequential, one context (see *Same-file* below) | |
| A milestone that won't close into a brief | ✓ — that's planning; finish it, don't dispatch | |
| A milestone spanning a design decision | ✓ — decide it first (or grill the PLAN); not delegable | |

**The rule under the table:** dispatch only when the milestones are **independent, closed, and touch
disjoint files** — and there are **3 or more**. Miss any → build inline. This is deliberately the
same gate as [[bug-workflow]]'s fan-out gate, because the economics are identical: fan-out's fixed cost
roughly equals the saving at 1–2 units and clearly loses to concurrency at 3+.

**Why a single milestone builds inline, even a big one:** unlike a bug fan-out (N independent fixes),
milestones within *one* request are usually **sequentially dependent** — M2 needs M1's schema, M3
needs M2's verb. A lone substantial milestone has no sibling to run concurrently *with*, so the spawn
buys nothing but overhead and a context hand-off. Dispatch earns its keep across **independent
milestones** (often across *different* requests), not one request sliced into pieces.

### Same-file / ordered milestones — serialize, don't parallelize

Fan-out's disjoint-file rule exists because **parallel writers to one file collide** — two Builders
editing `cli/spectacular` in separate windows produce two diffs against the same base; the second
clobbers or conflicts. When milestones share files or one blocks the next, you build them **inline,
sequentially, in one context** — you *are* the single serial writer, so there's no race. Three cases,
same as the debug fleet's:

- **Independent, disjoint files** — the only case that fans out. Each Builder owns its own files.
- **Ordered / blocking** (M2 only makes sense after M1) — build M1 → **re-verify** → *then* plan M2
  against the new codebase state (M1 may have moved lines or changed what M2 builds against) → build
  M2 → verify. The re-verify between steps is the point.
- **Genuinely entangled** (two "milestones" are one change) — it's one milestone; don't split what
  shares code.

**No git branches / worktrees for this.** Branches solve *concurrent* writers; serializing inline
means there are none. `isolation: worktree` per Builder + a merge-back step is real merge machinery
to re-solve what sequential-inline already dissolves — YAGNI until same-request fan-out is measurably
too slow, which it won't be. Parallelism is for genuinely independent milestones; a single request
gets one careful serial hand.

## B2 — Size the unit: one pass, or several phased sub-steps? (optional)

B1 decided *who* builds this milestone. B2 decides *how big the unit of work is* — and it's the answer
to the failure you'll otherwise hit: **a fat milestone built as one atomic pass runs for a long time as
an opaque block.** Whether you're building inline or dispatching, if the milestone is one coherent
change you proceed as-is; if it **spans several phases that each have their own verify point**, you
break it into sub-steps and do them one at a time, confirming between each. The sub-step boundaries
*are* the checkpoints — that's where you see progress, catch a wrong turn early, and keep any single
unit short.

### The sizing test

A milestone is **multi-phase** when it contains two or more phases that each end at their own
runnable check — e.g.:
- schema → CLI wiring → tests (three verify points)
- parser → renderer → doctor validation
- new data model → the code that reads it → the migration

A milestone is **single-phase** (skip this step) when it's one coherent change with one verify point
— a focused edit, a single function + its test, a doc section. **Most small milestones are
single-phase.** Don't manufacture sub-steps for a change that's already one pass — that's ceremony,
and the `decompose-large-milestone` policy is a `warn` precisely so a coherent one-pass milestone
sails through.

### If multi-phase — decompose into sequential sub-steps

1. **Write the sub-steps as nested `- [ ]` bullets** under this milestone's `### M<n>` block in
   `TASKS.md`. These are the existing "nested acceptance checklist" bullets ([[tasks-rules]]) put to
   a second use: **decomposition checkpoints.** They're not counted in progress (top-level milestones
   still own the `x/total`), so they add visible structure without distorting the ledger. Mirror them
   as harness `TaskCreate` items too — that's the *live* signal (drives the CLI progress UI); the
   nested bullets are the *durable* record. This is the two-layer task model ([[AGENTS]] § Task
   tracking) applied one level below the milestone.
2. **Do one sub-step at a time**, in order. Build (or dispatch) it, run its check, **report the
   result and tick its nested bullet before starting the next.** Each sub-step is short and visible;
   nothing runs unwatched for long. The checkpoint between sub-steps is where you confirm the phase
   landed and re-plan the next against the new state (an earlier sub-step may have moved lines).
3. **When dispatching a multi-phase milestone, send one closed sub-brief per sub-step in sequence** —
   confirm each returned diff before dispatching the next. Never hand a Builder the whole fat
   milestone as one brief that runs for hours; that's the exact opacity this step exists to prevent.
   Ordered sub-steps within one milestone are the **same-file / ordered case** from B1 —
   **serialize, never parallelize** (they share the milestone's files and build on each other).

### If single-phase — proceed to Phase C

One coherent pass, one check → go straight to C1 (inline) or C2 (dispatch) as B1 decided. No
sub-steps, no nested bullets. The step cost you one sizing judgment and nothing more.

> **Why this isn't a live-trace on a running agent.** The visibility comes from *shortening the
> units and checkpointing between them* — not from watching a Builder mid-run (the platform spawns
> subagents fire-and-forget; you can't stream from inside one). Decomposition turns one long opaque
> dispatch into several short confirmable ones. That's the whole mechanism.

---

# Phase C — Execute

*Build the milestone — yourself (C1) or via dispatched Builder(s) (C2) — and, for dispatched work,
plan what you'll distrust on return (C3).*

## C1 — Build inline

The self-serve path (B1's default column). Read the milestone's chain, build the Approach, run the
Success criteria, verify. Then go straight to **D2** (confirm + record) — **skip C3 and D1: the
critical-check plan is for dispatched work and reconciliation is for parallel fan-out, and here you
built it alone, so you *are* the check and there are no sibling pieces to reconcile.** This is the
common case — most requests are a handful of sequentially-dependent milestones you build one at a time.

## C2 — Dispatch spec-builder(s)

When B1 says dispatch:

1. **Write one closed brief per milestone** — the five slots, in prose, in each subagent's prompt:
   **Goal / Constraints / Approach / Expected output / Success criteria** (assembled from the chain in
   A2). The bar: **each slot concrete enough that a Builder who has never seen this request builds it
   without a decision.** If any slot leaves the Builder something to *decide* — which approach, which
   default, what "done" means — it isn't closed; close it first (A2). A slot you can't fill concretely
   means the milestone isn't dispatchable.

   **[join] Every brief carries its half of the shared contract.** For fanned-out pieces, the join you
   designed in A3 goes into *each* brief verbatim — the exact signature, name/key, or protocol the
   piece must produce or consume. This is the half of the brief a builder *cannot* infer from its own
   files, because the other side lives in a sibling's window. A brief that names the piece but not its
   join is the delegated-design-problem: the builder will pick something plausible, and it won't match.
   Put the contract in writing, same words in every brief that touches it.

   **Placeholder phrases are brief failures.** "TBD", "add appropriate error handling", "similar to
   M<N>", "as needed" in any slot means the brief isn't closed — the Builder must bounce it, and the
   orchestrator must not send it. A brief is written for a reader with zero request context; anything
   that assumes shared context is a placeholder wearing prose.
2. **Spawn one `spec-builder` per independent milestone, in parallel.** Each reads only its own
   target files, builds only its Approach, verifies against its Success criteria, and returns
   `built+diff` or `bounced+reason`. It never ticks a checkbox or moves lifecycle state.
3. **Collect all results before any ledger write.** Each Builder's returned block **is its Agent-tool
   result** (machine-read, not human prose). Gather the full set before recording, so the request's
   progress reflects one coherent batch.
4. **Route by verdict** (check the block is well-formed first — a `VERDICT` from the enum, a real
   `MILESTONE`, and for `built` a non-empty `DIFF` + a `VERIFY` that actually ran; a malformed/empty
   return is not a trustworthy `built` — treat it as a bounce):
   - **built** → **D2** (and **D1 first, only if you fanned out in parallel** — a single dispatched
     milestone has nothing to reconcile): confirm the diff yourself, run the Success criteria, tick the
     milestone's `TASKS.md` checkbox. The checkbox tick + any lifecycle move stay single-threaded on you.
   - **bounced** → the brief was open (undecided design, vague check, or the milestone spans another).
     **Finish the planning the bounce named, then re-dispatch a closed brief — never auto-retry the
     same open brief.** A bounce is the boundary working; respect it.

**Dispatch anti-patterns:** dispatching a single request's sequential milestones as if independent
(they collide / block each other — build inline); dispatching a milestone that won't close and
leaning on the bounce as your planning gate; dispatching pieces without their shared join (the
delegated-design-problem); auto-retrying a bounced brief instead of closing it. When in doubt,
self-serve — the table's default column is the safe one.

## C3 — Critical-check plan (dispatched work only)

**When you dispatch, you didn't watch the code happen.** A milestone you built inline (C1), you saw
take shape — you already know where it's sound. A milestone a Builder returns, you have **only the
diff**, and a diff can be long enough that reading it end-to-end at equal attention is how a subtle
break slips past. So before the Builder returns — *at dispatch time, while it runs* — decide **the 1–2
things you will check first and hardest on return.** Keep them in mind; they're what D2 confirms
against.

This step is **only for dispatched work.** If you built the milestone yourself, skip it — you were
the check. This is the concrete form of *delegated work is confirmed harder than self-built*: the
plan is the extra scrutiny delegation earns.

**What goes in the plan — the sharp points, not a re-review.** Name the parts where a *plausible-looking
diff could still be wrong* and it would matter most:
- **The join / shared interface** — the contract from A3 that this piece must uphold *and* any existing
  caller of what this piece exports. A Builder working in its own window can get its *own* files right and
  still break the join (mismatch the signature, pick a different key, register in the wrong order) or break
  a caller it never saw (its own public surface changed under code it didn't read). For fan-out this is the
  [join] check at the individual-piece level (D1 checks it across all pieces); for a lone dispatched
  milestone it's the same instinct — what depends on this that the Builder couldn't see?
- **The risky seam** — where this milestone meets code it didn't write (the integration point, the
  call into a shared module, the assumption about an upstream return value).
- **The security- or correctness-critical line** — the auth check, the boundary condition, the money
  math, the destructive operation. The place where "looks fine" isn't good enough.
- **The Success-criteria gap** — anything the milestone's check *doesn't* actually exercise. If the
  brief's Success criteria passes but wouldn't catch a specific failure, that failure is a C3 item.

Pick the **one or two** that matter for *this* milestone — a plan that lists ten things is a
re-review, not a critical-check plan. The test: *"if this diff is subtly wrong, where would the
damage be, and what's the fastest thing I can look at to find out?"* That's your plan.

**Fan-out note:** with N dispatched milestones, each gets its own short plan (they touch different
seams). Write them as you write each brief in C2 — the brief says what to build, the plan says what
you'll distrust on return. A [[code-reviewer]] dispatch (D2) *complements* the plan for a big/risky diff
— it doesn't replace it: the plan is *your* targeted first look, the reviewer is a broad independent
one.

---

# Phase D — Reconcile + record

*Do the pieces fit together (D1, parallel fan-out only), then confirm and record (D2). The ledger is
yours alone.*

## D1 — Reconcile the parallel outputs (only after disjoint fan-out)

**When you fanned out N builders on disjoint files (A3's parallel path → B1's fan-out → C2), each one
built in its own window, blind to its siblings.** Each got its *own* files right. But "the pieces are
disjoint" and "the join holds" were claims you made at routing time — before recording, you confirm they
held: **do the independently-built pieces actually compose into one coherent whole?** This is the pass
that makes "built in parallel, reconciled afterwards" safe, and it's where the [join] you designed in A3
gets its verdict. It runs **only** for parallel disjoint fan-out — a single milestone (C1) or a serial
sub-step chain (B2) has nothing to reconcile; you were the one serial writer.

**What reconciliation checks — the seams *between* the pieces, not the pieces themselves.** Each
builder's own Success criteria already graded its own files (that's C3 + D2). Reconciliation is the
layer they *couldn't* see:

- **[join] The shared contract holds across every piece.** The interface you designed in A3 and wrote
  into every brief — the signature, schema, or key the pieces agree on — must match on *both* sides.
  Builder A produced it; Builder B consumed it; neither saw the other. Verify the producer and consumer
  actually agree (same shape, same name, same assumptions). This is the direct payoff of designing the
  join up front: reconciliation is *checking a contract you already wrote*, not discovering one now.
- **The seams fit.** Where piece A's output feeds piece B, or two pieces register into the same list /
  config / dispatch table, confirm they don't collide, duplicate, or leave a gap. Two builders adding to
  the same registry from disjoint files can each be locally correct and jointly wrong (same id, clashing
  order, missing wiring).
- **The whole composes.** Run the check that only exists at the *assembled* level — the integration point,
  the end-to-end path that no single piece's Success criteria exercised. If nothing at that level runs,
  that gap is itself a reconciliation finding (mirror C3's Success-criteria-gap test, one level up).

**On a reconciliation failure, it's not a bounce — it's yours to resolve.** A builder didn't fail (each
met its brief); the *seam between briefs* is where the mismatch lives, and that seam was never any single
builder's job — **it was yours, the join you owned.** Plan the fix yourself — usually a small stitch
(align the contract, fix the registration, add the missing wire) — or, if the mismatch is real work, cut
a follow-up milestone. **Don't tick any of the parallel milestones' checkboxes until they reconcile.**
Locally-green pieces that don't compose are not done.

**This is cheap when A3's gates held and expensive when they didn't** — which is exactly why the
module-clarity gate (and its [join] test) is a gate. Clean boundaries + a named contract make
reconciliation a quick confirmation; fuzzy seams make it the debugging you routed *around*. If
reconciliation is repeatedly hard, that's the signal the work wasn't as disjoint as you routed it — next
time, serialize it (B1) instead of parallelizing.

## D2 — Confirm + record (the ledger stays yours)

After a milestone is **built and verified** (and, for fan-out, **reconciled** in D1), the orchestrator —
and *only* the orchestrator — records it:

1. **Confirm the diff — hit your critical-check plan first.** Read what the Builder returned; it's
   your change now. For dispatched work, start with the 1–2 points you named in C3 — verify *those*
   before a general read, because that's where a plausible-looking diff would hurt. **Delegated work is
   confirmed harder than self-built:** inline, you watched it happen; delegated, the diff is all you
   have, so the scrutiny is higher. Don't record a milestone you didn't confirm from the returned diff +
   verify. **For parallel disjoint fan-out, D1 comes first** — locally-green pieces that don't compose
   are not done, and no checkbox ticks until they fit.
2. **Run (or re-run) the Success criteria** if you didn't watch it run. Real check, real output.
3. **Optional — arms-length review + verify (judgment-gated).** Before you tick, consider dispatching
   the review/verify agents. Same worth-it economics as fan-out: skip for trivial changes; reach for
   them when the change earns it.
   - **[[code-reviewer]]** (read-only) — when the milestone diff is **substantial or medium+ blast
     radius** (touches a shared module, a schema, a widely-imported helper), dispatch it over the diff.
     It returns severity-ranked findings across five lenses (correctness/structure/security/perf/
     dead-code) with fix *direction*, never the fix. You triage: route real findings to a `debug-fixer`
     (single-site) or `spec-builder` (larger) with a closed brief, or accept them as known trade-offs.
   - **[[test-verifier]]** (apply-only, tests only) — when the Builder **self-reported the pass** or
     blast radius is medium+, dispatch it for *independent* pass/fail instead of trusting the self-check
     — *the agent that built it shouldn't be the only one to grade it*. It runs the named check (or
     writes a test to a closed spec) in its own window and returns honest pass/fail. On `fail`, treat
     the milestone as not-done: plan the fix, don't tick.
   Both are **optional** — the default path (you confirm the diff + run the check yourself) stands for
   ordinary work. They earn their dispatch on risk, not by default.
4. **Tick the checkbox.** `- [ ]` → `- [x]` for the milestone's completed tasks in `TASKS.md`. This
   is a main-thread write, never a subagent's — the same invariant that keeps the debug fleet's
   `fixes/`/`audit/` writes single-threaded (`use-audit-fix-verbs` in spirit: the labor fans out, the
   state write doesn't). A `code-reviewer` finding, a `test-verifier` `fail`, or an unresolved D1
   mismatch blocks the tick until resolved — reviewer/verifier/reconcile inform the record; they never
   write it.
5. **Decide the lifecycle move.** All of the request's milestones done → the request is ready for
   `active → review` (via `spectacular advance <slug>`; verification runs at `review → verified` per
   [[verify]]). A single milestone done mid-request → no status move, just the checkbox. The status
   move is a judgment call and a mutation — both yours.

   **At the `review → verified` gate**, a full-request-diff pass is worth considering even when you
   skipped per-milestone review: dispatch [[code-reviewer]] over the whole request's diff (a coherent
   review the per-milestone passes can't give), and/or [[test-verifier]] over the request's acceptance
   check — arms-length confirmation before the request is called done. Optional, but this is the
   highest-leverage place to spend a review dispatch.

**What never happens in this step by a subagent:** ticking a checkbox, moving `status:`, writing a
decision or memory. The Builder returns `LEDGER: not-written` precisely because this step is the
orchestrator's alone.

---

## The loop, in one line

**A Frame** (chain closes into a brief → route the shape: monolith-self / serialize / fan-out — and if
you split, *design the join*) **→ B Decide** (inline vs dispatch at ≥3 independent disjoint milestones;
size each unit — multi-phase → sequential sub-steps) **→ C Execute** (build inline, or dispatch builders
each carrying their half of the join; for dispatched work, name what you'll distrust on return) **→ D
Reconcile + record** (parallel fan-out → confirm the join held; then confirm diff + run check + tick +
lifecycle move).

The orchestrator bookends the fleet: **you plan (close the brief + own the join) → Builder builds → you
record (reconcile + tick + move).** The Builder is delegable labor with the honest-fallback invariant —
it bounces to planning rather than freelance a missing spec, and never writes the ledger.

---

## Relation to bug-workflow

Same substrate, opposite direction:

| | [[bug-workflow]] (fix) | build-workflow (build) |
|---|---|---|
| Trigger | a bug / quirk reported | a planned milestone to implement |
| Input | a symptom (unknown cause/site) | a closed milestone chain (known plan) |
| Discover role | `debug-investigator` (open bug → why/where findings) | `repo-explorer` (unfamiliar subsystem → seams/precedent map, *for planning*) — optional; when the ground is known, the plan *is* the findings |
| Apply role | `debug-fixer` (closed fix → applied) | `spec-builder` (closed milestone → built) |
| Brief size | 5 slots, single-fix | 5 slots, milestone (≈ small PR) |
| Fan-out gate | ≥3 independent closed disjoint-file fixes | ≥3 independent closed disjoint-file milestones |
| The join | rarely — fixes are usually single-site | **the orchestrator's to design** when milestones fan out (A3 [join] → C2 briefs → D1 verify) |
| Bounce meaning | root cause wrong / cross-cutting | brief open / design undecided / spans another M |
| Ledger | `fixes/` + `audit/` (main-thread) | `TASKS.md` tick + lifecycle move (main-thread) |

The Builder is to a milestone what the Fixer is to a fix — the closed-brief, bounce-on-judgment,
never-writes-the-ledger contract, scaled up one granularity.

**Related:** [[spec-builder]] (the agent contract), [[bug-workflow]] (the mirror), [[new-request]]
(where the chain is authored), [[verify]] (the `review → verified` gate), [[lifecycle]], [[doc-index]].
