---
status: planned
updated: 2026-06-27
related:
  - PLAN.md
---

# Tasks — naming-coherence

## v1

### M1 — advance rename (lifecycle promote → advance)
- [ ] Add `advance` dispatch → `cmd_promote`
- [ ] Keep `promote <slug>` as alias with one-line deprecation notice
- [ ] Update `--help`, docs/commands.md, SKILL.md routing
- [ ] Tests: both `advance` and `promote` paths

### M2 — feedback rename (feedback-loop → feedback)
- [ ] Rename verb to `feedback`; `feedback-loop` → hidden alias
- [ ] Update skill refs + CLI + docs (8-file surface)
- [ ] Leave `feedback-rules.md` doc name unchanged

### M3 — drop convention-pack form
- [ ] Make `pack` canonical everywhere
- [ ] `convention-pack` → silent alias, flagged for removal next release

### M4 — spectacular next
- [ ] Add read-only `cmd_next` → prints single highest-priority next action
- [ ] Wire into dispatch + help
- [ ] Test: active workspace, empty workspace, mutates nothing

### M5 — tier-reveal suggestions (skill)
- [ ] new-request.md ends with grill suggestion
- [ ] active-request.md / lifecycle.md suggest next at checkpoints
- [ ] archive.md suggests spec-sync
- [ ] Rule: one suggestion max, never mid-flow

## v2 (deferred)
- [ ] Remove `convention-pack` alias entirely
- [ ] Tiered presentation in human docs/
