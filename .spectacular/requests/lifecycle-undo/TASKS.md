---
status: planned
updated: 2026-06-27
related:
  - PLAN.md
---

# Tasks — lifecycle-undo

## v1

### M1 — Undo a status transition
- [ ] Add shared `_record_mutation <verb> <slug> <prior-status>` helper (writes `.spectacular/.last-mutation`)
- [ ] Call it from `cmd_promote` before the `fm_set status`
- [ ] Add `cmd_undo`: read breadcrumb, reverse a status transition (PLAN + TASKS), `fm_touch`, clear breadcrumb
- [ ] Wire `undo` into main dispatch + `--help`
- [ ] Add `.spectacular/.last-mutation` to `.gitignore`

### M2 — Undo an archive
- [ ] Record breadcrumb in `cmd_archive` (capture prior status before the move)
- [ ] In `cmd_undo`: reverse dir move `archive/<slug> → requests/<slug>` (git-aware, plain-mv fallback)
- [ ] Restore `status`, drop `archived:` field
- [ ] Reverse inbound-link rewrites (`../../archive/<slug>/ → ../<slug>/`) across sibling requests

### M3 — Undo an idea promote
- [ ] Record breadcrumb in `cmd_idea` promote path
- [ ] Restore idea source from `archive/ideas/`, reset status
- [ ] Prompt: remove scaffolded request or keep (default: keep)

### M4 — Guardrails + clarity
- [ ] No breadcrumb → "nothing to undo", exit 0
- [ ] `--dry-run` prints the reversal, mutates nothing
- [ ] Staleness check: refuse if breadcrumb timestamp predates touched files (advisory)
- [ ] Confirm prompt before any directory move
- [ ] `tests/cli/undo.test.sh` covering M1–M4

## v2 (deferred)

- [ ] Multi-level undo (breadcrumb stack instead of single)
- [ ] Full git-state verification (refuse when working tree moved since the mutation)
