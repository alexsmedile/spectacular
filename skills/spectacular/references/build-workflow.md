---
doc-id: build-workflow
kind: reference
summary: "How the skill implements planned work: assemble a closed milestone brief from the request chain, decide build-inline vs dispatch (and when it's worth it), fan out spec-builders for independent milestones, then confirm + tick the ledger. The build-direction analog of bug-workflow.md."
status: active
---

# Build workflow — implement a milestone, dispatch when it's worth it

Loaded when actively implementing a request that's in `active` status — turning a PLAN/TASKS
milestone into shipped code. Governs how the skill assembles a **closed brief** from the request
chain, decides whether to **build inline or dispatch a `spec-builder`**, fans out for independent
milestones, and confirms results. This is the **build direction**; its mirror is [[bug-workflow]]
(the fix direction). Same shape: closed-brief test, bounce rail, single-threaded ledger.

## The arc — read this first (orchestrator's whole job)

**You are the orchestrator: the only role that holds the whole request and the only hand that
mutates.** You assemble the brief, decide dispatch, spawn builders, confirm diffs, tick the
`TASKS.md` checkbox, and move lifecycle state. The `spec-builder` is narrow delegable labor
(closed-brief, build-only); the two things you never delegate are **planning** (closing an
open brief) and **mutation** (the checkbox tick + status move). Verbs, no overlap:
**you plan → Builder builds → you confirm + record.**

The steps, in order:

```
0   Chain closes?      walk task-row → M-block → PLAN §2/§3/§6/§7 → can the 5 slots fill? → else bounce-to-planning
1   Worth-it gate      milestone closed AND independent? → dispatch; else build inline
1b  Fan-out gate       ≥3 independent closed disjoint-file milestones → fan out; else self-serve
2a  Build inline       apply the milestone yourself, verify                          (no dispatch)
2b  Dispatch           spawn spec-builder(s) with closed brief(s), collect results
3   Confirm + record   confirm diff · run Success criteria · tick TASKS checkbox · decide lifecycle move
```

> **@Implementation policy gate.** First, run `spectacular policy @Implementation` and follow every
> active policy returned — chiefly `understand-before-change` (the request's Understanding slot must
> be filled before you build — Step 0 *is* that understanding for the milestone) and `build-order`
> (P11 — build the milestone the plan ordered, not a later one). All are `warn`. See [[policy-injection]].

---

## Step 0 — Does the chain close into a brief? (the context-assembly chain)

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

---

## Step 1 — Build inline, or dispatch? (the worth-it gate)

Once Step 0 confirms the milestone closes, one question: do *you* (the main thread) build it, or
dispatch a `spec-builder`? Dispatch has a fixed cost (spawn + brief-write + collect + confirm). It
pays back only when the milestone is **substantial enough that a fresh window builds it better than
crowding the main thread** — and when you have **more than one** to run.

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
same gate as [[bug-workflow]] Step 1b, because the economics are identical: fan-out's fixed cost
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

---

## Step 2a — Build inline

The self-serve path (Step 1's default column). Read the milestone's chain, build the Approach,
run the Success criteria, verify. Then Step 3 (confirm + record). This is the common case — most
requests are a handful of sequentially-dependent milestones you build one at a time.

## Step 2b — Dispatch spec-builder(s)

When the table says dispatch:

1. **Write one closed brief per milestone** — the five slots, in prose, in each subagent's prompt:
   **Goal / Constraints / Approach / Expected output / Success criteria** (assembled from the chain in
   Step 0). The bar: **each slot concrete enough that a Builder who has never seen this request builds
   it without a decision.** If any slot leaves the Builder something to *decide* — which approach,
   which default, what "done" means — it isn't closed; close it first (Step 0). A slot you can't fill
   concretely means the milestone isn't dispatchable.
2. **Spawn one `spec-builder` per independent milestone, in parallel.** Each reads only its own
   target files, builds only its Approach, verifies against its Success criteria, and returns
   `built+diff` or `bounced+reason`. It never ticks a checkbox or moves lifecycle state.
3. **Collect all results before any ledger write.** Each Builder's returned block **is its Agent-tool
   result** (machine-read, not human prose). Gather the full set before recording, so the request's
   progress reflects one coherent batch.
4. **Route by verdict** (check the block is well-formed first — a `VERDICT` from the enum, a real
   `MILESTONE`, and for `built` a non-empty `DIFF` + a `VERIFY` that actually ran; a malformed/empty
   return is not a trustworthy `built` — treat it as a bounce):
   - **built** → Step 3: confirm the diff yourself, run the Success criteria, tick the milestone's
     `TASKS.md` checkbox. The checkbox tick + any lifecycle move stay single-threaded on you.
   - **bounced** → the brief was open (undecided design, vague check, or the milestone spans another).
     **Finish the planning the bounce named, then re-dispatch a closed brief — never auto-retry the
     same open brief.** A bounce is the boundary working; respect it.

**Dispatch anti-patterns:** dispatching a single request's sequential milestones as if independent
(they collide / block each other — build inline); dispatching a milestone that won't close and
leaning on the bounce as your planning gate; auto-retrying a bounced brief instead of closing it.
When in doubt, self-serve — the table's default column is the safe one.

---

## Step 3 — Confirm + record (the ledger stays yours)

After a milestone is **built and verified**, the orchestrator — and *only* the orchestrator —
records it:

1. **Confirm the diff.** Read what the Builder returned; it's your change now. Don't record a
   milestone you didn't confirm from the returned diff + verify.
2. **Run (or re-run) the Success criteria** if you didn't watch it run. Real check, real output.
3. **Tick the checkbox.** `- [ ]` → `- [x]` for the milestone's completed tasks in `TASKS.md`. This
   is a main-thread write, never the Builder's — the same invariant that keeps the debug fleet's
   `fixes/`/`audit/` writes single-threaded (`use-audit-fix-verbs` in spirit: the labor fans out, the
   state write doesn't).
4. **Decide the lifecycle move.** All of the request's milestones done → the request is ready for
   `active → review` (via `spectacular advance <slug>`; verification runs at `review → verified` per
   [[verify]]). A single milestone done mid-request → no status move, just the checkbox. The status
   move is a judgment call and a mutation — both yours.

**What never happens in this step by a subagent:** ticking a checkbox, moving `status:`, writing a
decision or memory. The Builder returns `LEDGER: not-written` precisely because this step is the
orchestrator's alone.

---

## The loop, in one line

**chain closes? (assemble brief) → worth-it gate (inline vs dispatch) → self-serve or fan out (≥3
independent closed disjoint-file milestones → N× spec-builder) → confirm diff + run check + tick
checkbox + decide lifecycle move.**

The orchestrator bookends the fleet: **you plan (close the brief) → Builder builds → you record
(tick + move).** The Builder is delegable labor with the honest-fallback invariant — it bounces to
planning rather than freelance a missing spec, and never writes the ledger.

---

## Relation to bug-workflow

Same substrate, opposite direction:

| | [[bug-workflow]] (fix) | build-workflow (build) |
|---|---|---|
| Trigger | a bug / quirk reported | a planned milestone to implement |
| Input | a symptom (unknown cause/site) | a closed milestone chain (known plan) |
| Discover role | `debug-investigator` (open → findings) | none — the plan *is* the findings |
| Apply role | `debug-fixer` (closed fix → applied) | `spec-builder` (closed milestone → built) |
| Brief size | 5 slots, single-fix | 5 slots, milestone (≈ small PR) |
| Fan-out gate | ≥3 independent closed disjoint-file fixes | ≥3 independent closed disjoint-file milestones |
| Bounce meaning | root cause wrong / cross-cutting | brief open / design undecided / spans another M |
| Ledger | `fixes/` + `audit/` (main-thread) | `TASKS.md` tick + lifecycle move (main-thread) |

The Builder is to a milestone what the Fixer is to a fix — the closed-brief, bounce-on-judgment,
never-writes-the-ledger contract, scaled up one granularity.

**Related:** [[spec-builder]] (the agent contract), [[bug-workflow]] (the mirror), [[new-request]]
(where the chain is authored), [[verify]] (the `review → verified` gate), [[lifecycle]], [[doc-index]].
