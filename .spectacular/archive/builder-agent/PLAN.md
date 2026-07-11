---
status: archived
priority: medium
owner: alex
updated: 2026-07-11
build: b21
summary: "Add a Builder agent + build-workflow orchestrator arc — implement planned request milestones via subagent, the build-direction analog of the debug fleet"
related:
  - ../../PRD.md
archived: 2026-07-11
---
<!-- Spec source: .spectacular/ideas/coding-agents.md §7 (cited in prose below; not a related: target — cross-folder paths don't resolve from the request dir). -->


# Plan — builder-agent

## Goal

Give Spectacular a **Builder** agent that implements a planned request milestone from
its already-written PLAN/TASKS chain, plus a `build-workflow.md` orchestrator arc that
decides *when* to dispatch it (and when not) — the build-direction analog of the
existing debug fleet.

## Constraints

- **Reuse the debug fleet's shape, don't invent a new one.** Same closed-brief test,
  same bounce-on-judgment rail, same ledger discipline (agent never writes state; the
  orchestrator mutates). A Builder is a Fixer scaled from single-fix to milestone
  granularity — not a new safety model.
- **Closed brief or bounce.** A bare `- [ ]` TASKS line is a fragment, not a brief. The
  orchestrator assembles the brief by walking the context-assembly chain (task row →
  milestone block → PLAN §2/§3/§6/§7) — see idea §7. If the chain won't close (vague
  validation line, undecided design), the Builder bounces; it never freelances a spec.
- **Ledger stays single-threaded.** The Builder never ticks `TASKS.md` checkboxes and
  never touches lifecycle state — that's a main-thread write, mirroring the debug
  fleet's `use-audit-fix-verbs` invariant.
- **Name-match, don't require ID alignment.** The TASKS↔PLAN milestone link is an
  `M<N>` convention, not an enforced ID (doctor flags drift but doesn't block). The
  Builder leans on milestone *name* matching, tolerating renumbered-but-same-named
  milestones — matching the advisory doctor check shipped 2026-07-06.
- **Small (P3).** One agent def + one orchestrator ref. No new CLI verb, no config
  surface. `--delegable`-style CLI signals stay prose-first (idea open question).
- **Distinct from Fixer.** Fixer's brief is 5 slots because a bug fix is "make one thing
  true again." A Builder's brief spans a milestone — closer to a small PR than a patch.
  The build-workflow arc must not blur the two.

## Understanding

### How it works now

Spectacular ships a **debug fleet** (`debug-fixer`, `debug-investigator`,
`debug-researcher`) and a fully-worked orchestrator arc in `bug-workflow.md`:
ceremony gate (Step 1) → fan-out gate (Step 1b) → open the job (Step 1c) →
investigate/plan (Step 2b) → resolve+log (Step 3). The trace substrate
(`debug/<slug>/`), the closed-brief contract, and the bounce rail are all built.

But this is the **fix direction only** — critique/repair existing code. There is
**no build-direction agent**: no way to hand a planned request milestone to a subagent
for implementation. The context to do so already exists in every request folder
(PLAN §2/§3/§6/§7 + TASKS.md M-blocks) but nothing assembles it into a closed brief.
The idea's §7 worked out the context-assembly chain and confirmed it's grep-able
against real requests; the M-label drift gap it flagged was closed by the
`doctor lifecycle` check shipped 2026-07-06.

### What changes

- New agent def `.claude/agents/spec-builder.md` (project scope for testing;
  graduates to a plugin agent once proven, like the debug fleet). *(Shipped as
  `spec-builder` — build-from-spec — not the working title `debug-builder`.)*
- New orchestrator reference `skills/spectacular/references/build-workflow.md` —
  mirrors `bug-workflow.md`'s arc for the build direction: context-assembly chain,
  worth-it/fan-out gate, bounce rule, ledger discipline.
- SKILL.md trigger row + doc-index row so the ref is discoverable.

### What stays the same

- No new CLI verb, no config field, no lifecycle transition changes. TASKS checkbox
  ticks and status moves stay main-thread writes.
- The debug fleet is untouched — Builder sits *alongside* it, sharing shape not code.
- The trace substrate (`debug/<slug>/`) is a debug concept; whether Builder jobs get
  their own trace folder is an M2 question, not assumed.

## Milestones

- M1 — **Builder agent def shipped.** ✅ `spec-builder.md` exists, matches the
  debug-fixer contract shape (frontmatter, protocol, bounce rail, output block), and
  reads as apply-from-spec not investigate. A worked closed brief runs end-to-end
  against a real milestone. *(Verified 2026-07-06: built commit-discipline M1; bounce case passed.)*
- M2 — **build-workflow orchestrator arc shipped.** ✅ `build-workflow.md` mirrors
  `bug-workflow.md`: the context-assembly chain, the worth-it/fan-out decision table
  (when to dispatch vs build inline), the bounce rule, ledger discipline. SKILL.md +
  doc-index wired so it's discoverable.
- M3 — **(gated) trace + CLI signal.** Only if M1/M2 prove the fan-out earns a durable
  trace or a `--delegable`-style CLI emit. Deferred per P11 — don't build speculatively.

## Tasks

See `TASKS.md`.

## Dependencies

- None blocking. Builds on shipped substrate: the debug fleet contract, `bug-workflow.md`,
  and the `doctor lifecycle` M-label drift check (2026-07-06). Reads the idea doc
  `.spectacular/ideas/coding-agents.md` §7 as the spec source.

## Validation

- M1 — the agent def parses as a valid `.claude/agents/*.md`; a hand-written closed
  brief (assembled from an existing request's milestone) dispatched to it produces a
  correct implementation + verification, and a deliberately-vague brief makes it bounce.
- M2 — `build-workflow.md` has frontmatter parity with `bug-workflow.md` (doc-id, kind,
  summary); the SKILL.md trigger row + doc-index row resolve; a walkthrough of the
  worth-it gate on a real multi-milestone request produces the right dispatch/inline call.
- M3 — n/a until gated open.

## Deliverables

- `.claude/agents/spec-builder.md` — the Builder agent prompt. ✅
- `skills/spectacular/references/build-workflow.md` — the orchestrator arc. ✅
- SKILL.md trigger row + `doc-index.md` row wiring the ref in. ✅
