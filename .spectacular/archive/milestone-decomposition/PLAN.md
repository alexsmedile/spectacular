---
status: archived
priority: medium
owner: alex
updated: 2026-07-11
build: b27
summary: "Step 1.5 size-and-decompose gate in build-workflow + decompose-large-milestone @Implementation policy — break fat milestones into visible sequential sub-steps (decisions pre-locked in ideas/milestone-decomposition.md)"
related:
  - PRD.md
  - ../../ideas/milestone-decomposition.md
  - ../../../skills/spectacular/references/build-workflow.md
  - ../../POLICY.md
  - ../../AGENTS.md
archived: 2026-07-11
---

# Plan — milestone-decomposition

> **Design pre-locked** via grill (2026-07-11). Full rationale + the exact Step 1.5 text and
> policy block are in [[ideas/milestone-decomposition]]. This PLAN is the buildable extract.

## Goal

Give the build-workflow orchestrator a size-and-decompose gate: when a milestone spans multiple verify-points, it breaks the milestone into sequential sub-steps (nested TASKS bullets + harness tasks) and builds/dispatches them one at a time with a checkpoint between each — so no milestone runs for hours as one opaque block the orchestrator can't see into.

## Constraints

- **No new substrate.** Reuse what exists: nested `- [ ]` bullets in TASKS.md (durable checkpoints, already a documented convention) + harness `TaskCreate`/`TaskUpdate` (live progress signal, already the two-layer model in AGENTS.md). No trace folder, no doctor area, no CLI change.
- **Decomposition, not live streaming.** This fixes visibility via sub-step *boundaries*, not by watching a running agent. The live-trace approach (builder-agent's deferred M3) stays out of scope — it fights the platform's fire-and-forget subagent model.
- **Ordered sub-steps serialize.** Sub-steps within one milestone are the same-file/ordered case — build sequentially, one closed sub-brief at a time; never fan out sub-steps in parallel.
- **Sizing is judgment → warn, not block.** A single-phase milestone needs no sub-steps; the policy must not manufacture ceremony for a coherent one-pass change.
- **Authored legibly** per the v1.30.1 patch: the new warn policy gets an `**Override:**` clause and no `⛔` marker.

## Understanding

### How it works now
- `build-workflow.md` has a milestone-level arc: Step 0 (chain closes into a brief), Step 1 (worth-it: inline vs dispatch), Step 1b (fan-out ≥3), Step 2a/2b (build/dispatch), Step 3 (confirm + record). Every path treats a **whole milestone** as the atomic unit.
- The doc explicitly notes a milestone can be "≈ a small PR" — i.e. large — but offers no step to break one down. A fat milestone → one `spec-builder` → one atomic dispatch → hours of silence.
- Sub-step surfaces already exist but the build arc never invokes them: `tasks-rules.md` allows nested `- [ ]` bullets ("acceptance checklist, not counted"); `AGENTS.md` § Task tracking says harness tasks are "finer — decompose a single TASKS.md milestone into concrete edits/tests."
- `@Implementation` policies are injected when the skill enters that phase (build-workflow runs `spectacular policy @Implementation`).

### What changes
1. **New Step 1.5 — size-and-decompose gate** in `build-workflow.md`, between Step 1 (worth-it) and Step 2 (build/dispatch).
2. **New `@Implementation` warn policy** `decompose-large-milestone` in POLICY.md so the rule is injected at runtime.
3. **A mirror line in `bug-workflow.md` Step 3** — a fix can also be multi-phase (rarer), same discipline.
4. **A note in `AGENTS.md` § Task tracking** tying nested-bullet checkpoints to the two-layer model.
5. **A `tasks-rules.md` line** blessing nested bullets as decomposition checkpoints, not only acceptance criteria.

### What stays the same
- The milestone-level gates (worth-it, fan-out ≥3, same-file serialize) are unchanged — Step 1.5 sits *below* them.
- No change to `spec-builder`'s contract, `_policy_records`, or any CLI path.
- Single-phase milestones behave exactly as today (skip the gate).

## Decisions

- **Fix = decomposition, not live trace** — chose sub-step boundaries as the visibility mechanism over streaming from a running agent, because decomposition sidesteps the platform's fire-and-forget constraint and fixes both "runs for hours" and "no visibility" at once. Rejected: reviving M3's live trace (fights the platform; heavier).
- **Sub-steps live in ephemeral nested TASKS bullets** — chose the existing nested-bullet convention over a new `build/<m>/` trace folder (zero new substrate, no doctor area) and over a no-artifact habit (keeps a durable checkpoint record). This is the smaller-slice version of M3 and may retire M3 permanently.
- **Step 1.5 as its own step** — chose a distinct step over folding into Step 1, because sizing (one-pass vs multi-phase) is a different decision from dispatch-vs-inline.
- **Doc + policy, not doc alone** — chose to pair the doc step with a warn policy so the rule is injected at runtime (a doc step nobody loads is a poster — Principle 3 enforcement pattern).
- **Warn, not block** — sizing is judgment; a block would force ceremony on coherent single-phase milestones.

## Milestones

- M1 — Step 1.5 gate live in `build-workflow.md`: a new step between Step 1 and Step 2 that sizes the milestone, and on multi-phase decomposes into sequential nested-bullet sub-steps built/dispatched one at a time with a checkpoint between each. The arc diagram + one-line loop updated to include it.
- M2 — `decompose-large-milestone` @Implementation warn policy in POLICY.md: appears in `spectacular policy @Implementation`, parses cleanly, reads legibly (Override clause, no ⛔ marker).
- M3 — Mirror + convention notes: a `bug-workflow.md` Step 3 line for multi-phase fixes; an `AGENTS.md` § Task tracking note tying nested-bullet checkpoints to the two-layer model; a `tasks-rules.md` line blessing nested bullets as decomposition checkpoints.
- M4 — Docs synced: `specs/index.md` notes the gate; CHANGELOG; plugin bump.

## Tasks

See `TASKS.md`.

## Dependencies

- Builds on the v1.30.1 legibility patch (POLICY.md v1.5) — the new warn is authored in that style.
- Independent of [[stance-layer]] (b26). No ordering constraint between them.
- Relates to [[builder-trace-and-fanout]] — this request is the smaller-slice alternative to that idea's M3; if it lands, M3 likely stays permanently parked.

## Validation

- M1 — run: `grep -n "Step 1.5\|size-and-decompose\|sub-step" skills/spectacular/references/build-workflow.md` matches; observable: the arc diagram lists Step 1.5 between Step 1 and Step 2; observable: a walkthrough of a synthetic multi-phase milestone produces nested `- [ ]` checkpoints + a report between each, and a single-phase milestone skips the gate.
- M2 — run: `spectacular policy @Implementation` lists `decompose-large-milestone`; run: `spectacular doctor policies` exits 0 with the new block counted; observable: the block has an `**Override:**` line and NO `⛔` marker (it's a warn).
- M3 — grep: `bug-workflow.md` Step 3 references multi-phase decomposition; `AGENTS.md` § Task tracking references nested-bullet checkpoints; `tasks-rules.md` blesses nested bullets as decomposition checkpoints.
- M4 — grep: `specs/index.md` mentions the size-and-decompose gate; CHANGELOG entry present; `spectacular doctor` overall exits 0.

## Deliverables

- `skills/spectacular/references/build-workflow.md` — new Step 1.5 + updated arc diagram + one-line loop.
- `.spectacular/POLICY.md` — one new `### decompose-large-milestone` block under `## @Implementation`.
- `skills/spectacular/references/bug-workflow.md` — one mirror line in Step 3.
- `.spectacular/AGENTS.md` — a § Task tracking note.
- `skills/spectacular/references/tasks-rules.md` — a nested-bullet-as-checkpoint line.
- Docs — `specs/index.md`, CHANGELOG, plugin version bump.
