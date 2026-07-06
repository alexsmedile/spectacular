---
status: active
updated: 2026-07-06
related:
  - PLAN.md
  - TASKS.md
---

# Session — builder-agent

## Current state

Request scaffolded (b21) and moved `planned → active` (2026-07-06). PLAN's 7 slots + Understanding
are filled; the spec source is `.spectacular/ideas/coding-agents.md` §7, which had already worked
out the context-assembly chain and confirmed it grep-able against real requests. No design questions
gate implementation — the debug fleet's shape (`debug-fixer` + `bug-workflow.md`) is the template to
mirror, one granularity up.

## Active task

**M1 — Builder agent def shipped. ✅ DONE (2026-07-06).**
- `.claude/agents/spec-builder.md` — the Builder prompt. Mirrors `debug-fixer.md`'s contract shape
  (frontmatter, 5-slot closed brief, 7-step protocol, bounce-on-judgment rail, machine-read output
  block). Brief adapted to milestone granularity: **Goal / Constraints / Approach / Expected output /
  Success criteria**. Named `spec-builder` (build-from-spec) to signal the build direction, not
  `debug-*`. `tools: Read, Edit, Write, Bash, Grep, Glob` (adds Write over Fixer — milestones create
  files). Bounce fires on: open design decision, vague check, milestone spanning another, state moved.

**M2 — build-workflow orchestrator arc shipped. ✅ DONE (2026-07-06).**
- `skills/spectacular/references/build-workflow.md` — the orchestrator arc, mirror of
  `bug-workflow.md`. Steps 0–3: chain-closes (context-assembly chain + slot→source table +
  name-match rule) → worth-it gate (self-vs-dispatch table, why-single-milestone-builds-inline,
  same-file serialize rule, no-worktrees) → dispatch loop → confirm+record (ledger stays main-thread).
  Includes a bug-workflow↔build-workflow comparison table.
- `SKILL.md` — trigger row ("Implementing a milestone — decide build-inline vs dispatch").
- `doc-index.md` — `build-workflow` row in Skill-internal references.
- `ROADMAP.md` b21 row + `CLAUDE.md` Active Requests row.

**M3 — (gated) trace + CLI signal.** Not started; deferred per P11 until M1/M2 prove a durable trace
or a `--delegable` CLI emit earns its keep.

## Verification so far

- `spec-builder.md` frontmatter parses (name/description/tools/model present).
- `build-workflow.md` has frontmatter parity with `bug-workflow.md` (doc-id/kind/summary/status).
- `doctor lifecycle` shows no M-label drift (TASKS ↔ PLAN §3/§6 aligned); `doctor docs` clean.
- **End-to-end dispatch — DONE (2026-07-06).** Live-tested `spec-builder` against two real briefs:
  - **Built case:** dispatched the closed brief for `commit-discipline` M1 (assembled from its
    PLAN/TASKS chain per the build-workflow arc). Builder added the `commit-checkpoint` POLICY entry
    + extended `@SessionEnd` prose + `active-request.md` note, in house style, chose principle 11
    with justification, ran both Success-criteria checks. Orchestrator re-verified independently
    (`policy @Implementation` lists it; `doctor policies` green) and ticked `commit-discipline` M1's
    checkboxes. **This shipped commit-discipline M1 as a side effect of the test.**
  - **Bounce case:** dispatched a deliberately open brief (undecided design: commit info in status,
    no format/layer decided). Builder read the code, discovered `status` is skill-owned not a CLI
    command, named the exact CLI-vs-skill fork the plan left open, and bounced with a precise
    "orchestrator must decide X/Y/Z" report — refused to guess. Safety rail confirmed.
- **Still pending:** a worth-it-gate walkthrough on a real multi-milestone *dispatch* (≥3 independent
  milestones fanned out) — the fan-out path isn't exercised yet, only single dispatch + bounce.

## Next

The core contract is proven (built + bounce). Move `active → review` when ready; the fan-out
walkthrough is the remaining `review → verified` item.
