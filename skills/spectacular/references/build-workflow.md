---
doc-id: build-workflow
kind: reference
summary: "Runtime core for implementing planned work: frame a closed milestone brief from the request chain, route its shape (monolith / serial / fan-out), decide inline-vs-dispatch, execute, reconcile + record. Four phases: Frame → Decide → Execute → Reconcile+record. Rationale lives in build-workflow-doctrine.md; the fix-direction mirror is bug-workflow.md."
status: active
---

# Build workflow — implement a milestone, dispatch when it's worth it

Loaded when actively implementing a request in `active` status. **You are the orchestrator: the only
role that holds the whole request and the only hand that mutates.** The `spec-builder` is narrow
delegable labor (closed-brief, build-only); the two things you never delegate are **planning**
(closing an open brief) and **mutation** (the checkbox tick + status move). Verbs, no overlap:
**you plan → Builder builds → you confirm + record.**

> This is the runtime core: the arc, the gates, the tables, the procedures. The *why* behind each
> rule — rationale, failure modes, the relation to [[bug-workflow]] — is in
> [[build-workflow-doctrine]]. Load it only when a routing call feels uncertain or you're editing
> this workflow itself.

## The arc

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
> active policy returned — chiefly `understand-before-change` (A2 *is* the understanding),
> `build-order` (P11 — build the milestone the plan ordered), and `decompose-large-milestone` (B2).
> All are `warn`. See [[policy-injection]].

**The one binding idea: the join is yours.** When you split work across builders, the shared
contract (signature, schema, registry key, protocol) is a design decision only you can make — a
builder can't invent the interface it shares with a sibling it never saw. Design it in A3, write it
into every brief in C2, verify it held in D1. Flagged **[join]** below.

---

# Phase A — Frame the work

## A1 — Map unfamiliar ground first (optional — dispatch `repo-explorer`)

If the chain won't close because *you* don't understand the subsystem well enough to write the
Approach — and reading it inline would crowd the main window — dispatch [[repo-explorer]] with a
scoped question. It returns a structured map (entry points, precedent to mirror, seams `file:line`,
blast radius); you plan A2 against it. Skip when you know the ground (most milestones). It **maps,
it doesn't plan** — the map informs the brief; it isn't the brief.

## A2 — Does the chain close into a brief?

A bare `- [ ]` TASKS line is a fragment, not a brief. Walk the chain the request lifecycle already
wrote:

```
TASK ROW            →  MILESTONE BLOCK      →  PLAN SECTIONS         →  BRIEF                 →  OUTPUT + CHECK
TASKS.md             TASKS.md's own          PLAN.md §3 (this M's     synthesized from all:    implementation +
one `- [ ]` line      `### M<n> — ...`        line) + §2 Constraints   Goal / Constraints /     Success-criteria
                      header + siblings       + §6 Validation (this    Approach / Expected       run
                                              M's check) + §7 Deliv.    output / Success crit.
```

| Brief slot | Sourced from |
|---|---|
| **Goal** | PLAN §3 — this milestone's line |
| **Constraints** | PLAN §2 + inherited STACK/PRINCIPLES |
| **Approach** | the ordered `- [ ]` steps in this milestone's `### M<n>` TASKS block |
| **Expected output** | PLAN §7 — the deliverable this milestone contributes |
| **Success criteria** | PLAN §6 — this milestone's validation line (reuse verbatim) |

**Match milestones by name, not number** — `M<n>` is a convention, not an enforced ID; when
TASKS↔PLAN numbering drifts, match on the name text after the em-dash (`doctor lifecycle` flags
drift as advisory).

**The chain won't close — bounce to planning (yourself) — when:**
- PLAN §6's validation line is vague ("works correctly") — no runnable Success criteria.
- §3 / the TASKS block implies a **design decision not yet made** (which format? which default?).
- The `### M<n>` block is unfilled boilerplate or has no real step list.
- §3/§6/§7 have no matching entry for this milestone.

Finish the planning first (grill/refine the PLAN, decide the open question, fill §6), *then*
dispatch. Don't hand a Builder an open brief and lean on its bounce — the bounce catches a *wrong*
brief; the quality gate for a *vague* one is you, here.

## A3 — Route the shape (and design the join)

**Default: monolithic-self.** Delegation is the exception the shape must earn; when in doubt, build
it yourself. Two independent axes:

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

**[join] The moment you choose decompose-and-delegate, name the shared contract explicitly, in
writing, before any brief:** the signature/shape (args, return type, data shape), the name/key (the
exact identifier/registry key/route both sides use), the ordering/protocol (who registers first,
who owns the default). It becomes the shared half of every brief (C2) and the thing D1 verifies.
**If you can't state the join cleanly, the pieces aren't disjoint — serialize instead.**

**Three readiness gates on decompose-and-delegate — miss any → build inline:**

| Gate | The question | If it fails |
|---|---|---|
| **Benefit of delegation** | Does a fresh window build this *better*, with ≥1 sibling to run *with*? | Lone/trivial piece → self-serve. |
| **Plan readiness** | Is each piece a **closed brief** (A2's five slots fill)? | Open brief → finish planning first. |
| **Module clarity** | Clean boundaries — each piece owns its files, and the join is nameable? | Fuzzy seams → serialize instead. |

| Shape | Route to |
|---|---|
| One hard monolithic thing | **C1** (build inline) |
| One fat milestone of several verify-point phases | **B2** (sequential sub-steps) |
| Several coupled / ordered milestones | **B1** same-file/ordered rule (serialize inline) |
| Disjoint pieces, too few to fan out | **C1** (self-serve) |
| Disjoint pieces, enough to fan out | **B1 → C2** (fan out N builders, each brief carrying the join) |
| Disjoint pieces run in parallel | as above **+ D1** (reconcile) |

---

# Phase B — Decide who & how big

## B1 — Build inline, or dispatch? (the worth-it gate)

| Situation in hand | Orchestrator builds it **itself** | **Dispatch** to spec-builder |
|---|---|---|
| A 1–2 line / trivial milestone | ✓ — a spawn costs more than the edit | |
| A single substantial milestone | ✓ — one milestone, one context | |
| **≥3 independent, closed, disjoint-file milestones** | | ✓ — concurrency wins |
| Milestones that share files / are order-dependent | ✓ — sequential, one context | |
| A milestone that won't close into a brief | ✓ — that's planning; finish it | |
| A milestone spanning a design decision | ✓ — decide it first; not delegable | |

**The rule: dispatch only when the milestones are independent, closed, and touch disjoint files —
and there are 3 or more.** Miss any → build inline. (Fan-out's fixed cost ≈ the saving at 1–2
units; concurrency clearly wins at 3+. Same gate as [[bug-workflow]] Step 1b.)

**Same-file / ordered milestones — serialize, don't parallelize.** Parallel writers to one file
collide. Three cases:
- **Independent, disjoint files** — the only case that fans out.
- **Ordered / blocking** — build M1 → **re-verify** → *then* plan M2 against the new state → build →
  verify. The re-verify between steps is the point.
- **Genuinely entangled** — it's one milestone; don't split what shares code.

**No git branches / worktrees for this** — serializing inline means there are no concurrent
writers; branch machinery is YAGNI here.

## B2 — Size the unit: one pass, or phased sub-steps? (optional)

A milestone is **multi-phase** when it contains ≥2 phases that each end at their own runnable check
(schema → CLI wiring → tests; parser → renderer → doctor validation). **Single-phase** (one coherent
change, one verify point) → skip this step; don't manufacture sub-steps
(`decompose-large-milestone` is a `warn` so a one-pass milestone sails through).

If multi-phase:
1. **Write the sub-steps as nested `- [ ]` bullets** under the milestone's `### M<n>` block in
   `TASKS.md` (not counted in `x/total` progress). Mirror them as harness `TaskCreate` items — the
   live signal; the nested bullets are the durable record ([[AGENTS]] § Task tracking).
2. **Do one sub-step at a time, in order.** Build (or dispatch), run its check, report + tick its
   nested bullet before starting the next. The checkpoint is where you confirm the phase landed and
   re-plan the next against the new state.
3. **When dispatching a multi-phase milestone, send one closed sub-brief per sub-step in
   sequence** — confirm each return before dispatching the next. Never hand one opaque fat brief.
   Ordered sub-steps are the same-file/ordered case from B1 — **serialize, never parallelize**.

---

# Phase C — Execute

## C1 — Build inline

Read the milestone's chain, build the Approach, run the Success criteria, verify. Go straight to
**D2** — skip C3 (you were the check) and D1 (no sibling pieces). This is the common case.

## C2 — Dispatch spec-builder(s)

1. **Write one closed brief per milestone** — the five slots, in prose, in each subagent's prompt:
   **Goal / Constraints / Approach / Expected output / Success criteria** (from A2's chain). The
   bar: each slot concrete enough that a Builder with zero request context builds it without a
   decision. If a slot leaves the Builder something to *decide*, it isn't closed — close it first.
   - **[join] Every brief carries its half of the shared contract, verbatim, same words in every
     brief that touches it** — the one thing a builder cannot infer from its own files.
   - **A conventions capsule rides in every brief (Constraints slot), ~5 lines:** the test harness +
     how to run it, naming/error idiom, the precedent file to mirror. Siblings get the same capsule.
   - **Placeholder phrases are brief failures.** "TBD", "add appropriate error handling", "similar
     to M<N>", "as needed" in any slot → the brief isn't closed; don't send it.
2. **Spawn one `spec-builder` per independent milestone, in parallel.** Each builds only its
   Approach, verifies against its Success criteria, returns `built` or `bounced+reason`. It never
   ticks a checkbox or moves lifecycle state.
3. **Collect all results before any ledger write** — record one coherent batch.
4. **Route by verdict** (well-formed first: a `VERDICT` from the enum, a real `MILESTONE`, and for
   `built` a non-empty `CHANGED` list + a `VERIFY` that ran; malformed/empty = treat as bounce. The
   block carries no diff — the Builder's edits live in the working tree; pull with `git diff` in D2):
   - **built** → **D2** (D1 first, only if you fanned out in parallel).
   - **bounced** → finish the planning the bounce named, then re-dispatch a closed brief — **never
     auto-retry the same open brief.**

**Anti-patterns:** dispatching one request's sequential milestones as if independent; leaning on the
bounce as your planning gate; briefs without their shared join; auto-retrying a bounced brief. When
in doubt, self-serve.

## C3 — Critical-check plan (dispatched work only)

At dispatch time, while the Builder runs, decide **the 1–2 things you will check first and hardest
on return** — where a plausible-looking diff could be subtly wrong and it would matter most:
- **The join / shared interface** — the A3 contract this piece must uphold, and any existing caller
  of what it exports (the thing the Builder couldn't see).
- **The risky seam** — where this milestone meets code it didn't write.
- **The security- or correctness-critical line** — auth check, boundary condition, money math,
  destructive op.
- **The Success-criteria gap** — a failure the brief's check wouldn't catch.

Pick one or two — ten is a re-review, not a plan. With N dispatched milestones, each gets its own
short plan, written as you write each brief. Skip entirely for inline builds.

---

# Phase D — Reconcile + record

## D1 — Reconcile the parallel outputs (only after disjoint fan-out)

Each builder got its *own* files right; reconciliation checks the seams *between* pieces — the layer
none of them could see:
- **[join] The shared contract holds on both sides** — producer and consumer agree (same shape,
  same name, same assumptions).
- **The seams fit** — no collision, duplication, or gap where pieces feed each other or register
  into the same list/config/table.
- **The whole composes** — run the check that only exists at the assembled level; if nothing runs
  there, that gap is itself a finding.

**A reconciliation failure is yours, not a bounce** — each builder met its brief; the seam was your
join. Stitch it yourself (align the contract, fix the registration) or cut a follow-up milestone.
**No checkbox ticks until the pieces reconcile.** If reconciliation is repeatedly hard, the work
wasn't as disjoint as routed — next time serialize.

## D2 — Confirm + record (the ledger stays yours)

1. **Confirm the diff — C3 plan first.** Read the Builder's `CHANGED` list, then pull the real
   change from the tree: `git diff -- <those files>` (the return block carries no diff; read
   selectively, not wholesale). Start with the 1–2 C3 points. **Delegated work is confirmed harder
   than self-built.** For parallel fan-out, D1 comes first.
2. **Run (or re-run) the Success criteria** if you didn't watch it run.
3. **Optional — arms-length review + verify (judgment-gated, on risk not by default):**
   - **[[code-reviewer]]** — when the diff is substantial or medium+ blast radius. Returns ranked
     findings; you triage → `debug-fixer` (single-site) or `spec-builder` (larger) with a closed
     brief, or accept as trade-off.
   - **[[test-verifier]]** — when the Builder self-reported the pass or blast radius is medium+;
     independent pass/fail. On `fail`, the milestone is not-done: plan the fix, don't tick.
4. **Tick the checkbox** (`- [ ]` → `- [x]` in `TASKS.md`) — main-thread write, never a subagent's.
   A reviewer finding, a verifier `fail`, or an unresolved D1 mismatch blocks the tick.
5. **Decide the lifecycle move.** All milestones done → `active → review` via `spectacular advance
   <slug>` (verification at `review → verified` per [[verify]]). Mid-request → checkbox only. At
   the `review → verified` gate, consider a full-request-diff [[code-reviewer]] / [[test-verifier]]
   pass — the highest-leverage place to spend a review dispatch.

**Never done by a subagent:** ticking a checkbox, moving `status:`, writing a decision or memory.

---

## The loop, in one line

**A Frame** (chain closes → route the shape; if you split, design the join) **→ B Decide** (dispatch
at ≥3 independent disjoint milestones; multi-phase → sequential sub-steps) **→ C Execute** (inline,
or builders each carrying their half of the join + a critical-check plan) **→ D Reconcile + record**
(join held → diff confirmed → check run → tick + lifecycle move).

**Related:** [[build-workflow-doctrine]] (the why), [[spec-builder]] (the agent contract),
[[bug-workflow]] (the fix-direction mirror), [[new-request]], [[verify]], [[lifecycle]],
[[doc-index]].
