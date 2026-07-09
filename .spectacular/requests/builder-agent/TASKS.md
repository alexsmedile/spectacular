---
status: review
updated: 2026-07-10
related:
  - PLAN.md
---

# Tasks — builder-agent

<!--
  Executable checklist for one request.
  Lives at: .spectacular/requests/<slug>/TASKS.md

  Rules:
  - Group tasks by milestone using `### M<N> — <name>` headings.
  - Use `- [ ]` for open, `- [x]` for done. No other bullet syntax.
  - `status:` in frontmatter should match parent PLAN.md.
  - Tasks are owned by the user. Engine never adds/removes/reorders tasks.
-->

## v1

### M1 — Builder agent def shipped
- [x] Write `.claude/agents/spec-builder.md` matching the debug-fixer contract shape (frontmatter, protocol, bounce rail, output block) — shipped as `spec-builder` (build-from-spec), not `debug-builder`
- [x] Adapt the 5-slot brief to milestone granularity (Goal / Constraints / Approach / Expected output / Success criteria)
- [x] Encode the bounce rule (under-specified milestone, vague validation line, undecided design → stop and report)
- [x] Verify: def parses as valid agent; a closed brief from a real request milestone runs end-to-end (built commit-discipline M1); a vague brief bounces (verified live 2026-07-06)

### M2 — build-workflow orchestrator arc shipped
- [x] Write `skills/spectacular/references/build-workflow.md` mirroring `bug-workflow.md`'s arc
- [x] Document the context-assembly chain (task row → milestone block → PLAN §2/§3/§6/§7 → brief)
- [x] Document the worth-it / fan-out decision table (dispatch vs build inline; ≥N independent milestones; same-file serialize rule)
- [x] Document ledger discipline (Builder never ticks checkboxes / moves lifecycle — main-thread only)
- [x] Wire the SKILL.md trigger row + `doc-index.md` row
- [~] Verify: frontmatter parity with bug-workflow.md ✓; trigger + index rows resolve ✓; worth-it gate walkthrough on a real multi-milestone *fan-out* still pending (the one review→verified item)

### M3 — (gated) trace + CLI signal
> Do not start until M1/M2 prove a durable trace or CLI emit earns its keep (P11).
- [ ] Decide whether Builder jobs need their own trace folder (analog of `debug/<slug>/`)
- [ ] Decide whether a `--delegable`/`--brief` CLI emit is worth the surface

## v2 (deferred)

- [ ] Test Agent (idea candidate 2) — verify Builder output independently
- [ ] Review fleet (idea candidates 1, 5, 6) — correctness / inefficiency / dead-code
