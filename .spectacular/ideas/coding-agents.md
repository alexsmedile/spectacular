---
type: idea
status: parked
priority: medium
owner: alex
origin: (captured)
updated: 2026-07-07
promoted_to: null
related: []
---

# Idea — coding-agents

## Hypothesis

Spectacular should ship a small fleet of **coding agents** — reusable subagent definitions,
each with a concrete *protocol* (an explicit sequence of steps), not a vague role prompt. The
protocol is the product: a role prompt gives you an agent that improvises; a protocol gives you
one that follows a repeatable, auditable method. Three candidates below.

Each agent stays project-agnostic. Project-specific risk lists (e.g. "Swift 6 closure
isolation") are *injected* from the host repo's PRINCIPLES.md / DECISIONS.md / known-hard list,
not baked into the agent.

## Candidate agents

### 1. Code Review Agent

Catch correctness bugs (not style nits) with high signal, by targeting *where bugs live*
instead of reading uniformly.

1. **Hypothesize** — form failure hypotheses from the diff shape + the project's injected risk
   areas, before reading line-by-line.
2. **Target & fan out** — one sub-agent per hypothesis, each scoped to specific files/symbols,
   returning structured findings (file:line, claim, failure scenario, confidence).
3. **Adversarially verify** — a second agent tries to *refute* each finding (construct the
   input where it's actually fine). Kill what can't survive. This is the noise filter.
4. **Rank & report** — data loss > correctness > perf > cleanup; name the concrete failure
   input, not "this looks off".

### 2. Test Agent

Write and run *in-depth* tests that exercise real failure modes and drive behavior, not assert
on mocks.

1. **Map the surface** — public entry points, async/actor boundaries, state transitions. A test
   with no failure mode behind it is noise.
2. **Pick the level** — climb the ladder: fake/in-memory unit → integration through the store
   → end-to-end, only as high as the bug actually lives.
3. **Target the known-hard cases** — the project's injected list of historically-fragile seams.
4. **Run + observe** — actually execute, read output, report real pass/fail. Flaky → say so.
5. **Leave one runnable check** — smallest thing that fails if the logic breaks; no
   frameworks/fixtures unless the case needs them.

### 3. Researcher Agent

When a bug smells like a platform issue, find whether *others* already signaled it — forums,
SO, GitHub issues — before burning hours rediscovering a known bug. Emphasis: query design +
relevancy judging.

1. **Frame the symptom, not our code** — translate the internal failure into the vocabulary
   *other* users would use; strip project-specific names.
2. **Draft diverse queries** — exact error text / system symptom / API surface / platform
   version. One query finds one echo chamber; spread them.
3. **Search the right places** — weight by signal: official vendor forums, SO, issues on
   similar OSS projects, general web last. Escalate fetch tool on block/thin content.
4. **Judge relevancy hard** — same *symptom* ≠ same *cause*. Same version? same signature?
   confirmed or speculation? Drop plausible-but-unrelated; a wrong match misroutes the fix.
5. **Report with verdict + citations** — known platform bug (→ workaround) vs genuinely ours
   (→ fix). Link sources, note confidence, name what's unconfirmed. "No strong match" is a
   valid answer — cap rounds, don't force a weak one.

### 4. Fix-Applier Agent (debugging-delegation)

Apply a **closed** fix — one where the delegation prompt is complete and no mistake can be
made — via subagent, so the main thread fans out mechanical corrections instead of doing them
serially. Not a debugger: a debugger that already *knows the fix* and just executes it under a
contract.

**Load-bearing insight: the delegation boundary already exists.** `bug-workflow.md`'s Step-1
ceremony gate (audit-first vs just-fix) *is* the trust boundary for delegation. The "just-fix"
quadrant — reproducible, root cause known, single site, non-breaking — is definitionally the
set of fixes where a subagent can't go wrong. We don't invent a "can I trust an agent with
this?" heuristic; we read the one already shipped. Extend the 2-way gate to 3-way:

```
clear + single-site + mechanical  →  DELEGABLE   (spawn Fix-Applier)
clear + single-site + judgment    →  JUST-FIX    (main thread does it)
unclear OR multi-site             →  AUDIT-FIRST (main thread investigates)
```

The new split (mechanical vs judgment on an *already-clear* fix) mirrors the mechanical-vs-
judgment distinction the policy contract already draws for checks. "Delegable" = *closed*, not
merely *small*.

1. **Receive a closed brief, not a bug** — the prompt is the audit's five slots already filled:
   Problem / Intended / Root cause / **Proposed fix** / **Success criteria**. If any slot is
   empty, it's not delegable — it's audit-first. Filling the slots *is* the gate.
2. **Apply-only, no investigation** — make exactly the Proposed fix. Don't expand scope, don't
   refactor neighbours, don't re-diagnose.
3. **Verify against Success criteria** — run the check that already exists; report real
   pass/fail + the diff.
4. **Bounce on judgment** — if the Root cause turns out wrong or the fix needs a decision,
   **stop and report back** — never freelance. A delegated fix that discovers it needs judgment
   returns to the main thread; it doesn't improvise. This is the safety rail.
5. **Never writes the ledger** — returns diff + verdict only. The main thread confirms and runs
   `spectacular fix new` / `audit resolve --into-fix`. Fanning out the *labor* is fine; the
   *ledger write* stays single-threaded and CLI-gated (`use-audit-fix-verbs`). This is the one
   hard invariant — a subagent hand-writing `F<N>.md` violates the mutator contract.

Fan-out shape: one Fix-Applier per **independent** fix site. A multi-caller bug is *not* N
independent sites — it's one root fix (`fix-root-not-symptom`), so it stays single-threaded.
Delegation parallelizes genuinely-independent closed fixes, not one bug sliced into pieces.

### 5. Inefficiency-Finder Agent

Not correctness (that's Candidate 1) — this agent hunts *needless cost*: O(n²) where O(n)
exists, redundant re-computation/re-fetch, unbounded growth, blocking calls on a hot path.

1. **Scope the hot path** — profile data if present, else infer from call frequency (loops,
   request handlers, render paths). Cold/init-only code is out of scope — cost there is noise.
2. **Pattern-match known smells** — repeated work inside a loop that's loop-invariant, N+1
   queries, missing memoization/index, synchronous I/O blocking an async caller.
3. **Quantify, don't vibe** — state the actual complexity class or the repeated-call count, not
   "this looks slow." No complexity class/count → drop the finding.
4. **Weigh cost vs risk** — a correct-but-slow function used once a day isn't worth a risky
   rewrite. Rank by (frequency × cost) vs (diff size × blast radius).

### 6. Dead-Code Finder Agent

Find code with **zero live callers** — not "looks unused," a proven negative.

1. **Build the reachability set** — grep every export/symbol for call sites, including dynamic
   dispatch (string-based lookups, DI containers, reflection) before concluding "no callers."
2. **Exclude false positives by construction** — public API surface (library entry points),
   test-only helpers, framework-required overrides (lifecycle hooks) are not dead even with few
   callers.
3. **Confirm negative, not absence-of-search** — "grep found nothing" is a hypothesis; check
   re-exports, barrel files, and string-interpolated call sites before declaring dead.
4. **Report removal, not just detection** — dead code left in place accrues drift risk (source
   of truth confusion, next dev cargo-cults it, out-of-date next to what it once mirrored).
   Propose the deletion diff, not merely a list of names.

## Debugging taxonomy — jobs, roles, and the delegation boundary

The debugging fleet isn't "4 agents that seemed useful" — it's *derived*. Classify debug jobs
by **what's unknown**, map each to the **role** that resolves that unknown, and an agent exists
for a role **iff its prompt can be closed** (read-only, or a complete brief). That last rule is
the whole safety model: it draws the line the ceremony gate already draws, from the role side.

### The 8 debug jobs (by what's unknown)

| # | Job | Unknown | Ceremony | Primary role |
|---|---|---|---|---|
| 1 | Known-class | nothing (matches prior fix) | recognize | Recognizer |
| 2 | Mechanical | nothing (cause+fix clear, 1 site) | just-fix / delegate | Fixer |
| 3 | Localize | *where* | audit-first | Localizer |
| 4 | Diagnose | *why* | audit-first | Diagnostician |
| 5 | Cross-cutting | *how far* | audit-first | Diagnostician (find shared root) |
| 6 | External | *whose* (platform/dep?) | audit-first | Researcher |
| 7 | Heisenbug | *reproduction* | audit-first | Reproducer |
| 8 | Design-ambiguity | *is it a bug?* | audit-first | Diagnostician (adjudicate) |

Jobs 1–2 (nothing unknown) = the just-fix / delegable quadrant of `bug-workflow.md`'s ceremony
gate. Jobs 3–8 (a genuine unknown) = audit-first — and *an unknown is exactly what you can't
hand to a closed prompt.* The job taxonomy and the ceremony gate agree; that's the signal the
boundary is real.

### The 6 roles → 5 agents (Localizer + Diagnostician fuse into Investigator)

| Role | Agent? | Thread | Protocol (one line) |
|---|---|---|---|
| **Recognizer** | ✓ | delegable | `fix list` + signature-grep symptom → match(es) or "none, novel" |
| **Fixer** | ✓ built | delegable *when closed* | 5-slot brief → apply-only, verify vs criteria, bounce on judgment |
| **Investigator** | ✓ built | own window, *discover-only* | open bug → find where+why+who-else → report findings (root cause / ranked hypotheses / suspected sites). Never proposes the fix |
| **Researcher** | ✓ built | delegable | reframe symptom → diverse queries → verdict: platform vs ours, cited. Prefers scrapekit, falls back to harness web tools |
| **Reproducer** | ✓ | delegable | intermittent symptom → deterministic repro (or flakiness envelope). Never fixes |
| **Router** | ✗ **main** | main thread | triage: which job type? just-fix vs audit-first vs delegate (the ceremony gate) |
| **Validator** | ✗ **main** | main thread | collect Fixer/Investigator results, confirm diffs, route applied/bounced, own ledger writes |

**Investigator = Localizer (where) + Diagnostician (why) fused into one own-context-window agent.**
The earlier open question — "non-delegable = no agent, or = shared window?" — is now **resolved**:
they get an agent, but a *propose-never-mutate* one. Investigation wants a big window to reason in,
and it's safe to give it one **because it returns knowledge, never a mutation** — no code edits, no
ledger writes. That's a third safety mode alongside the other two.

**The invariant (revised):** an agent exists for a role iff its output can't do harm — one of:
*read-only* (Recognizer, Researcher), *closed brief* (Fixer applies a complete, verifiable spec),
or **propose-never-mutate** (Investigator returns a brief/audit; the main thread mutates). Router
and Validator stay on the main thread — they *are* the orchestrator: triage in, adjudication out.

### The MVP pipeline

```
Router (main) → Investigator (agent, own window) → orchestrator (main) → Fixer (agent ×N) → Validator (main)
 triage /        open bug → FINDINGS                 plans fixes           apply + verify     confirm, route,
 1-2 lines?      (root cause / hypotheses /          from findings         (or bounce)        write ledger
 fix inline      suspected sites) — never the fix    (closed briefs)
```

The orchestrator (main window) owns the judgment-heavy work: triage, **fix-planning**, and every
mutation. Two agents are the delegable labor — Investigator turns *open → findings* (discovers where
+ why + who-else), Fixer turns *closed → applied*. The key separation: **the Investigator discovers,
the orchestrator plans the fix, the Fixer applies it** — three verbs, no overlap. An obvious 1–2
line bug skips the whole pipeline; the orchestrator just fixes it. Both agents keep the
honest-fallback invariant: Investigator reports hypotheses-only rather than a fake root cause; Fixer
bounces rather than freelance. Neither writes the ledger.

### Build order

1. **Fixer** ✓ built + tested (closed → applied). The delegation MVP.
2. **Investigator** ✓ built (open → findings + plausible solutions, discover-only). Pairs with
   Fixer to cover the audit-first path without the main thread doing diagnosis inline.
3. **Researcher** ✓ built (external-smell → verdict + citations). Its own debugging protocol on top
   of scrapekit/harness web tools; routes the `needs-research` outcome from the Investigator.
4. Reproducer — deferred; only fires on Heisenbugs. Recognizer ≈ an inline grep, may never need its
   own agent.

Both agents ship as project `.claude/agents/` defs for testing; graduate to plugin agents once
proven.

### Review-fleet MVP proposal (candidates 5–6, not built yet)

Distinct from the debug fleet above — this is the code-review trio (correctness / inefficiency /
dead-code). Proposed scope, kept as an idea only:

- **Ship 2, not 3**: Code-Review Agent (candidate 1 — may already overlap with the existing
  `code-reviewer-lean` subagent, check before building new) + Dead-Code Finder (candidate 6 —
  reachability is mostly mechanical, verdict is a clean binary, cheapest to get right).
- **Defer Inefficiency-Finder** (candidate 5) — needs profiling data or frequency inference
  that's often unavailable, plus a subjective cost/risk weighing step. Ship after the harness
  shape is proven on the other two.
- **Reuse the debug fleet's shape**: read-only discover-only agents (findings only, no edits),
  one adversarial-verify pass per finding (skeptic tries to refute — the noise filter that made
  the debug fleet trustworthy), main thread ranks + decides whether to act.

### 7. Builder Agent (codes from spec, not from a bug report)

The opposite direction from review/debug: given planned work already captured in
`.spectacular/requests/<slug>/`, *implement* it — not critique existing code, not fix a bug,
build the thing the plan describes.

**A bare `TASKS.md` checkbox line is not a closed brief on its own** — it's a fragment. The
same "can this prompt be closed?" test the debug fleet uses (an agent exists iff its input can
be a complete, unambiguous brief) says: assemble the brief by walking up the existing chain
that already exists in every request folder, not by inventing a new context format.

**Context-assembly chain** (each link sourced from a file that already exists per request):

```
TASK ROW           →  MILESTONE BLOCK      →  PLAN SECTION        →  BRIEF                →  OUTPUT + TESTS
TASKS.md            TASKS.md's own          PLAN.md §3            synthesized from        implementation +
one `- [ ]` line     `## M<n> — ...`         milestone entry        all 4 links above:      validation run
                     header + siblings       ("M3 — CLI writers:    Goal / Constraints /
                                              implement...")        Expected output /
                                                                     Success criteria
```

1. **Task row** — the single checkbox line that triggered the work. Gives *what*, not *why* or
   *how much*.
2. **Milestone block** — the `## M<n>` section in `TASKS.md` containing that row plus its
   sibling `- [ ]` lines. Gives the full ordered step list and scope boundary (what's *in* this
   milestone vs the next one).
3. **PLAN section** — the matching milestone description in `PLAN.md` §3 "Milestones", plus
   §2 "Constraints" (applies to every milestone) and §6 "Validation" (per-milestone acceptance
   criteria — already written, not invented). This is where *why* and *done means* live.
4. **Brief** — the orchestrator synthesizes 1–3 into the same shape the Fixer already uses
   (Problem / Intended / Root cause → **here: Goal / Constraints / Approach** / Proposed
   change / **Success criteria** — reusing `PLAN.md §6`'s validation line verbatim where
   possible instead of re-deriving it).
5. **Expected output** — the deliverable named in `PLAN.md §7` "Deliverables" that this
   milestone contributes, so the agent knows the artifact shape (a new file? a CLI verb? a
   doc section?), not just "write some code."
6. **Tests and validation** — run exactly the check named in `PLAN.md §6` for this milestone
   (a CLI smoke test, a doctor area going green, a worked example in `discovery.md`) — never
   invent a new success bar the plan didn't specify.

**Bounce rule (mirrors Fixer):** if the milestone block turns out under-specified — the PLAN
section is missing, the validation line is vague ("works correctly"), or the task row implies a
design decision not yet made — the agent **stops and reports back**, same as Fixer bounces on
judgment. It never freelances a missing spec into existence.

**Why this isn't the same as Fixer:** Fixer's brief is 5 slots because a bug fix is bounded by
"make this one thing true again." A Builder's brief is bounded by a milestone, which can span
multiple files/verbs/docs — closer in size to a small PR than a patch. The closed-brief test
still applies, just at milestone granularity instead of single-fix granularity.

**Chain grep-ability — checked against real requests, mostly holds:** `## M<N>` headers,
`- [ ]`/`- [x]` checkboxes, and numbered `## <N>. <Section>` headers in PLAN.md are all
consistently grep/sed-able across every request sampled (`soft-db-substrate`, `spec-audit-mode`,
others). The one real gap: the link between a TASKS.md milestone and its PLAN.md §3/§6
counterpart is an **`M<N>` numbering convention, not an enforced ID** — nothing stopped them
drifting (confirmed: this repo's own `spec-audit-mode` TASKS.md was still the unfilled
boilerplate template, 3 milestones vs PLAN's 4). Decision: **the chain must stay walkable even
when M-numbers drift** — don't require alignment as a hard gate, an agent should tolerate
matching by name/prose if numbers disagree. **Shipped (2026-07-06):** `doctor lifecycle` now
flags M-label mismatches between TASKS.md ↔ PLAN §3 Milestones ↔ PLAN §6 Validation as an
advisory `judgment` warning (never blocks) — `check_lifecycle` in `cli/spectacular`, covered by
`scenario_19_milestone_label_drift` in `tests/cli/doctor.test.sh`. Also flags a non-standard
milestone ID prefix (e.g. `G1` instead of `M1`) as its own warning — reusing the project-wide
one-letter-per-entity convention now documented in `ARCHITECTURE.md` (`M`=milestones,
`D`=decisions, `F`=fixes, `b`=roadmap builds, `A`=debug findings) — **but only escalates to a
real "chain broken" warning when a name-based fallback match also fails**: it compares the
milestone's name text (the words after the em-dash) across files first, so a relettered- or
renumbered-but-same-named milestone is silently tolerated, not false-flagged. Covered by
`scenario_20_milestone_off_prefix_name_fallback`. This closes the gap for *detecting* drift; a
Builder agent consuming the chain should still lean on name-matching as its own primary
strategy, since the doctor check is advisory, not enforced.

**Open, not resolved:**
- Does the Builder write its own tests, or does a separate agent (existing Test Agent,
  candidate 2) verify after? Lean: Builder writes the check named in PLAN §6 (that's usually
  already a testable claim); a heavier suite is Test Agent's job, not duplicated here.
- Multi-milestone requests: one Builder per milestone (fan-out, since M-blocks are usually
  sequential/dependent) or a single Builder walking the whole `TASKS.md` in order? Milestones
  in the wild (see `soft-db-substrate/PLAN.md`) are often sequentially dependent (M2 needs M1's
  schema) — lean single-threaded per request, fan-out only across independent requests.
- Ledger discipline: does the Builder ever touch `TASKS.md` checkboxes itself (`- [ ]` → `- [x]`),
  or does that stay a main-thread write like the debug fleet's ledger rule? Lean: main thread
  only — mirrors "Fanning out the labor is fine; the ledger write stays single-threaded."

## Open questions

- ~~**Non-delegable = no agent, or = shared context window?**~~ **RESOLVED (2026-07-05):**
  own window, propose-only. Localizer + Diagnostician fused into the **Investigator** agent — it
  gets its own big context to reason in, and it's safe because it *proposes, never mutates* (returns
  a brief/audit; the main thread writes). Router + Validator stay main-thread. See the role map above.
- **Fix-Applier: does "DELEGABLE" need a CLI signal, or is it purely a skill judgment?** Today
  the ceremony gate is skill-side prose. A 3-way gate could stay prose, or `spectacular audit`
  could gain a `--delegable` disposition that emits the closed brief for a subagent. Lean prose
  first — don't build the CLI surface until the fan-out proves it earns its keep.
- **Where's the mechanical-vs-judgment line drawn?** String-sub / guard-insert / rename are
  clearly delegable; a fix touching control flow is clearly not. The grey middle needs a
  concrete rubric before this is safe to automate, or it silently delegates a judgment call.
- **Shared harness or three standalone agents?** All three are fan-out → verify → rank in
  shape. Is there one orchestration primitive underneath, or are they just siblings?
- **Where does the project-specific risk list come from?** A convention file the agent reads
  (`.spectacular/` known-hard list?), or injected per-invocation by the caller?
- **How do they compose?** Review → Test → (on platform-smell) Researcher is a natural loop.
  Does spectacular orchestrate that, or does the user chain them?
- **Test agent: writes, runs, or both?** Lean both — but report separately so a self-written
  green suite isn't mistaken for independent verification.
- **Packaging** — `.claude/agents/*.md` shipped in the plugin? A pack? A skill that spawns them?
- Overlap with the existing `code-reviewer-lean` / `scrapekit` subagents — reuse as the worker
  layer, agent adds the protocol brain on top?

## Promoted to

—
