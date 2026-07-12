---
version: 1.7
updated: 2026-07-12
summary: "Operating policies — the practice layer paired with PRINCIPLES.md"
---

# Spectacular — Operating Policies

<!--
  POLICY.md is the PRACTICE layer (PRINCIPLES.md is the THEORY layer).
  Policies are filed under named work-phase hooks (## @<hook>). The skill
  retrieves only the active hook's policies on entering a phase.

  Anatomy:  ### <verb>-<noun>
            - principle: N      (optional — the PRINCIPLES.md § it enforces)
            - severity: block   (block = refuse | warn = surface + continue)
            - check: <condition>
            - directive: <one imperative sentence — what the CLI injects at the
              hook gate; the body prose stays the full rationale>
            <prose: rationale + the instruction injected into context>

  Severity is OPT-IN to blocking: blocks ONLY if it explicitly says
  'severity: block'. Absent/warn/unrecognized → non-blocking.

  Valid hooks (9): @Init @Planning @Implementation @Debugging @Verification
  @Archive @Remember @Snapshot @SessionEnd.
-->

## @Init

### scaffold-contract
- principle: 4
- severity: warn
- directive: Verify README, committed `.spectacular/`, gitignored `.spectacular.local/`, and the always-set docs before finishing init.
The workspace must satisfy its scaffold contract — README present, `.spectacular/` committed, `.spectacular.local/` gitignored, always-set docs in place. Surface any gap on init; don't block.

## @Planning

### request-shape
- principle: 3
- severity: warn
- directive: Before work begins, make the PLAN carry a one-sentence Goal, explicit Constraints, and demoable Milestones.
A new request's PLAN must be well-shaped before work begins: a one-sentence Goal, explicit Constraints, and demoable Milestones — not a vague wish. Surface a thin plan; let the author proceed.

### scope-down
- principle: 10
- severity: warn
- directive: Cut the plan to the smallest slice that delivers the core value now; move everything else to ROADMAP `v2+`.
Before fixing milestones, name the smallest high-impact slice that delivers the core value now, and push the rest to ROADMAP as `v2+`. Prefer a finished MVP of the features actually needed today over a complete build of features that might be. Flag speculative generality and any feature without a current need. Surface the leaner cut; the human chooses the scope.

**Law: cut to the smallest slice that delivers the core value; defer the rest, don't build it.**
| Excuse | Reality |
|---|---|
| "While I'm in here anyway, I'll add…" | "While I'm here" is how a scoped request becomes a sprawl. The extra feature has no current need — it goes to ROADMAP `v2+`, not into this milestone. |
Red flag — stop: a milestone list that grows a feature nobody asked for *today*.

**Override:** legitimate to include more than the minimum when a later need is *already committed* and building it now is strictly cheaper than a second pass (a known migration, a schema both slices share). Name the committed need in the plan; "might need it" is not one.

### milestones-in-build-order
- principle: 11
- severity: warn
- check: milestones are ordered by dependency — M1 is demoable with nothing before it, and no milestone depends on a later one
- directive: Order milestones so each stands on the one before it — M1 demoable with nothing prior, nothing depending on a later slot.

A plan encodes its own build order. Milestones must be sequenced so each one stands on the one before it: M1 is demoable with no prerequisite, and nothing depends on a milestone that comes after it. An out-of-order milestone list is principle 11's failure written into the plan — reaching for the moon in slot 3 before the launchpad in slot 5. This promotes the refine-time ordering check (was human-only) to a surfaced policy. Surface any milestone whose prerequisite sits later in the list; the human reorders.

## @Implementation

### understand-before-change
> ⛔ **BLOCKING** — refuses `planned → active` until satisfied.
- principle: 7
- severity: block
- check: PLAN.md has a filled `## Understanding` section (How it works now / What changes / What stays the same), OR a `UNDERSTANDING.md` exists with the same three subheads
- directive: Write the three Understanding subheads — How it works now / What changes / What stays the same — before moving `planned → active`.

A request must not move `planned → active` until the agent has written down how the system works today, what this change touches, and what it leaves alone. Establish understanding before touching code. Satisfied by either the PLAN slot or a dedicated UNDERSTANDING.md.

**Law: no `planned → active` until the three understanding subheads are written down.**
| Excuse | Reality |
|---|---|
| "I've read enough to start." | Read-enough is unfalsifiable; a written `## Understanding` is the proof. If you can't fill *How it works now / What changes / What stays the same*, you haven't read enough — that's the gate working. |
Red flag — stop: moving to `active` with an empty or hand-waved Understanding section.

### build-order
- principle: 11
- severity: warn
- check: no task depends on a prior task that is still unbuilt or stubbed — each layer stands on a proven one below it
- directive: Build the lower layer first — never stack work on a stub, mock, or intention of your own code.

Build in prerequisite order: the thing you're building must stand on something that already works, not on a stub, a mock, or an intention. Don't write the abstraction before the first concrete case runs. Don't add the retry/cache/fallback before the happy path returns a real value. Don't wire integration before the unit it integrates is green. If a step needs a lower step that isn't real yet, build the lower step first — skipping it doesn't remove the work, it defers it to a worse moment. Surface an out-of-order jump; the human decides whether the foundation is actually there.

**Law: build each layer on a proven one below it — never on a stub, mock, or intention.**
| Excuse | Reality |
|---|---|
| "The stub's good enough to build on." | A stub is an intention wearing a return type. The layer you stack on it inherits its emptiness — and you'll rebuild both when the stub becomes real. Build the lower layer first; skipping it defers the work to a worse moment. |
Red flag — stop: writing the abstraction (or the retry/cache/integration) before its first concrete case runs green.

**Override:** legitimate to build against a stub when the lower layer is *external and contract-frozen* (a documented third-party API, a spec you can't run yet) — the stub encodes a real contract, not a guess. Note the frozen contract; a stub for your *own* unbuilt code is never this case.

### earn-the-verification
- principle: 11
- severity: warn
- check: the code path a verification exercises actually exists and runs — no green check on a stub or a mock standing in for the real thing
- directive: Before claiming verified, trace the check to the real code path — a green check on a stub or mock proves nothing.

A passing check on code that isn't built is worse than no check — it reports safety that isn't there. Before claiming a slice verified, confirm the check drives the real path, not a placeholder. An integrity gate on an empty thing has integrity to report about nothing.

**Law: a green check must exercise the real path — never a stub or mock standing in for it.**
| Excuse | Reality |
|---|---|
| "It's green, so it works." | Green proves the *check* ran, not that the *real path* did. A test hitting a mock is green and hollow. Trace the assertion to the real code before you call it verified — a false green is worse than a red. |
Red flag — stop: claiming "verified" when the check's target is a stub, mock, or not-yet-built path.

**Override:** a mock is legitimate when it stands in for something *genuinely out of scope to exercise* (a paid external call, a destructive side effect) AND a separate real-path check covers the seam. State what the mock replaces and where the real path is verified instead; "the real thing is hard to set up" is not this case.

### prefer-cli-mutator
- principle: 6
- severity: warn
- check: a structured mutation (lifecycle move, archive, snapshot, memory/decision/idea/audit/fix write) goes through its `spectacular` verb, not a free-form file edit — unless no verb covers it
- directive: Use the `spectacular` verb for any structured mutation; hand-edit only what no verb covers.

The CLI is the deterministic mutator; the skill orchestrates, reads, decides, communicates. Whenever a `spectacular` verb exists for a mutation, use it — hand-editing bypasses auto-numbering (`F<N>`/`A<N>`), signature and verified gates, index regeneration (`memories/index.md`), atomic link-rewriting on archive, and frontmatter bumps, and the drift lands as a `doctor` finding later. Manual edits remain the escape hatch for what no verb covers — that's the exception, not the default. If you catch yourself writing an entry file by hand where a verb exists, stop and run the verb. When a hand-edit already happened and left the substrate malformed, `spectacular doctor <area>` names the drift and `--fix` repairs the mechanical part.

### commit-checkpoint
- principle: 11
- severity: warn
- directive: When a milestone's tasks are all checked, suggest a local `git commit` before starting the next milestone.
Once a milestone's tasks are all checked off, the working tree holds an earned, working step — the natural moment to lock it in with a local `git commit` before starting the next milestone. An uncommitted milestone is a rung not yet earned to stand on: the next milestone builds on code that could still vanish with an errant `git checkout`/`reset`, and a debug session later has no checkpoint to bisect against. Spectacular never commits on your behalf and never blocks the next milestone on a clean tree — it surfaces the reminder at the boundary; committing (or explicitly deferring, with a reason) is the human's or agent's call. Distinct from `spectacular snapshot`, which versions canonical docs, not source.

### decompose-large-milestone
- principle: 10
- severity: warn
- check: a milestone that spans multiple verify-points is built/dispatched as sequential sub-steps (nested `- [ ]` TASKS bullets + harness tasks) with a checkpoint between each — not as one unbounded pass
- directive: Split any milestone spanning multiple verify-points into nested sub-step checkpoints and build them one at a time.

A milestone that spans several phases — each with its own verify point (schema → CLI wiring → tests; parser → renderer → doctor check) — is built one visible sub-step at a time, confirming at each boundary, not dispatched as a single opaque brief that runs for hours. Decompose it into nested `- [ ]` checkpoints under its `### M<n>` block (mirrored as harness tasks for the live signal), build or dispatch them sequentially, and report between each. The sub-step boundaries are the visibility — a fat milestone handed to one Builder as one brief disappears until it returns. See [[build-workflow]] B2. **Override:** a genuinely single-phase milestone — one coherent change, one check — needs no sub-steps; don't manufacture ceremony for a one-pass change.

## @Debugging

<!-- Entered when the user reports a bug, quirk, regression, or "why does X do Y".
     Loads [[bug-workflow]]. These policies gate the steps that doc prescribes. -->

### check-prior-fixes
- principle: 5
- severity: warn
- check: `.spectacular/fixes/` has been searched (fix list + signature grep) before diagnosis begins
- directive: Run `spectacular fix list` plus a signature grep of `fixes/` before forming any hypothesis.

Before diagnosing a bug, search the `fixes/` corpus — this is the payoff of operational memory. A solved bug rediscovered from scratch is the exact waste principle 5 exists to prevent. Grep the signatures for the symptom in hand; a match may resolve it immediately or name it as a known class. Surface what the search found (or that it was empty); don't block.

**Law: no diagnosis before `spectacular fix list` + a signature grep of `fixes/` has run.**
| Excuse | Reality |
|---|---|
| "This bug looks new" | Looking new is the signature of every rediscovered bug. The grep costs seconds. |
Red flag — stop: forming hypotheses before the signature grep has run.

### ceremony-matches-uncertainty
- principle: 11
- severity: warn
- check: an audit is opened only when root cause is unclear OR the fix spans multiple sites; a clear one-site fix is applied directly
- directive: Grep the callers first, then open an audit only if root cause is unclear or the fix spans multiple sites — otherwise just fix it.

Ceremony scales with uncertainty, not with every bug. Open an `audit/` only when the root cause is unknown, the symptom spans multiple callers/subsystems, it's not yet reproduced, or it might be intended behavior. If root cause is clear and the fix is one site, just fix it — an audit there is pure ceremony. Equally, don't skip the audit on a genuinely cross-cutting bug and patch one caller while siblings stay broken. Surface the mismatch; the human picks the ceremony level.

**Law: grep the callers of the shared path first, then choose — audit only if root cause is unclear or the fix spans multiple sites.**
| Excuse | Reality |
|---|---|
| "It's probably just this one site" | "Probably" without a caller grep is how siblings stay broken. Check, then choose. |
Red flag — stop: choosing just-fix or audit-first before grepping the callers.

### fix-root-not-symptom
- principle: 11
- severity: warn
- check: a multi-caller bug is fixed once at the shared root, not per-caller; the logged fix names the root cause, not the surface symptom
- directive: Grep every caller of the shared path and place the fix at the root, not at the reported site.

Fix the bug where all callers route through, not at the one site the report named. A guard in the shared function is a smaller, complete diff; a guard in one caller leaves every sibling still broken. This is the debugging face of principle 11 — the symptom stands on a root cause; fix the layer underneath, not the layer you can see.

**Law: grep every caller of the shared function before placing a guard at the reported site.**
| Excuse | Reality |
|---|---|
| "The ticket names this caller" | The ticket names a *symptom*. Grep every caller of the shared function before placing the guard. |
Red flag — stop: patching the named path without having listed its sibling callers.

### log-only-verified-reusable
- principle: 5
- severity: warn
- check: a `fixes/F<N>.md` entry is written only after the fix is verified AND carries reusable knowledge (non-obvious cause or recurring class), with a `--signature`
- directive: Log a `fixes/F<N>` entry only after verification and only when it carries reusable knowledge — always with a `--signature`.

A fix entry is earned, not automatic: log only once resolved *and* verified, and only when a future agent would benefit — a non-obvious cause, a recurring class, something that took real investigation. Skip typos, renames, one-off edits. The corpus is valuable because it's curated. The `--signature` is mandatory — it's what makes the entry findable next time (closes the loop back to `check-prior-fixes`).

**Law: write `fixes/F<N>` only via `spectacular fix new --signature …` after verification, and only when it carries reuse value.**
| Excuse | Reality |
|---|---|
| "Log it now so it isn't forgotten" | An unverified entry poisons the corpus for the next signature search. Verify, then log. |
Red flag — stop: an F&lt;N&gt; draft containing "should fix" phrasing.

### use-audit-fix-verbs
- principle: 6
- severity: warn
- check: audit and fix entries are created via `spectacular audit new` / `fix new`, not by hand-writing `A<N>.md` / `F<N>.md` files
- directive: Create audit and fix entries with `spectacular audit new` / `fix new`, never by hand-writing the files.

Write audit and fix entries with their verbs, never by hand. `spectacular audit new` / `fix new` auto-number (`A<N>`/`F<N>`), stamp frontmatter, enforce the `--signature` and `--verified-by` gates, and `audit resolve --into-fix` copies every slot forward. A hand-written entry skips all of that and lands as a `doctor fixes` finding. This is `prefer-cli-mutator` (@Implementation) applied to the debugging phase — stated here because a bug flow fires `@Debugging`, not `@Implementation`, and this is exactly the moment the temptation to hand-write an entry appears.

## @Verification

### verification-present
> ⛔ **BLOCKING** — refuses `review → verified` until every check is satisfied.
- principle: 7
- severity: block
- check: every check in VERIFY.md (or PLAN § Validation) is satisfied before `review → verified`
- directive: Satisfy every check in VERIFY.md (or PLAN § Validation) before moving `review → verified`.

A request must not reach `verified` while any verification check is unmet. Verification always happens; the only question is which artifact carries the checks. (Absorbs verify-walk's gate.)

## @Archive

### spec-sync
- principle: 2
- severity: warn
- directive: On archive, propose the `specs/index.md` and `specs/` updates the shipped work implies — the human confirms.
On archiving a request, propose the specs/index.md / `specs/` updates the shipped work implies. Intent and truth are different files — keep truth current. The human confirms.

### memory-propose
- principle: 5
- severity: warn
- directive: On archive, propose any operational lesson worth keeping as a memory — never write one unconfirmed.
On archiving, propose any operational lesson worth keeping as a memory. Operational memory compounds. Surface the candidate; never write memory without confirmation.

## @Remember

### confirm-before-write
> ⛔ **BLOCKING** — refuses the memory write until the user has confirmed the text.
- principle: 8
- severity: block
- check: the user has confirmed the memory text before it is written to `.spectacular/memories/`
- directive: Show the exact memory text and get the user's confirmation before writing it.

Memory is team-visible and git-committed. Humans decide, agents propose: never write a memory the user has not seen and confirmed.

## @Snapshot

### snapshot-before-overwrite
> ⛔ **BLOCKING** — refuses the overwrite until a snapshot of the current version exists.
- principle: 8
- severity: block
- check: a `<DOC>@v<N>.md` snapshot exists before a canonical doc is overwritten in place
- directive: Snapshot the current version (`<DOC>@v<N>.md`) before overwriting any canonical doc in place.

Canonical documents are never overwritten without a snapshot first. The unversioned filename always points to current; history is preserved.

## @SessionEnd

### summarize-before-handoff
- severity: warn
- directive: Before handing off, summarize what changed, what's left, and what's next — and suggest a commit of the session's work.
Before handing off, summarize what changed, what's left, and what's next, so the next session resumes without re-deriving context. Surface the summary; don't block. The same handoff moment is also the last natural checkpoint for outstanding source changes: suggest a local `git commit` of the session's work, or have the summary note explicitly why not (e.g. mid-edit, deliberately left uncommitted for review). Spectacular never runs the commit itself — it surfaces the suggestion alongside the summary; the human or agent decides.
