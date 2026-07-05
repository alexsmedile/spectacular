---
type: idea
status: parked
priority: medium
owner: alex
origin: (captured)
updated: 2026-07-05
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
