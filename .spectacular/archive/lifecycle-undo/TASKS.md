---
status: verified
updated: 2026-06-28
related:
  - PLAN.md
---

# Tasks — lifecycle-undo

## v1

### M1 — Undo a status transition
- [x] Added shared breadcrumb helpers `_breadcrumb_write/_get/_clear` (writes `.spectacular/.last-mutation`)
- [x] `cmd_promote` records breadcrumb (prev_status) after the transition
- [x] `cmd_undo` reverses a status transition (PLAN + TASKS), `fm_touch`, clears breadcrumb
- [x] Wired `undo` into main dispatch + `--help` + verb lists
- [x] Added `.spectacular/.last-mutation` to `.gitignore`

### M2 — Undo an archive
- [x] `cmd_archive` records breadcrumb (captures pre-archive status before the move)
- [x] `cmd_undo` reverses dir move `archive/<slug> → requests/<slug>` (git-aware, plain-mv fallback)
- [x] Restores `status`, drops `archived:` field (new `fm_unset` helper)
- [x] Reverses inbound-link rewrites (`../../archive/<slug>/ → ../<slug>/`) across sibling requests

### M3 — Undo an idea promote
- [x] `cmd_idea_promote` records breadcrumb (request_slug + idea_prev_status)
- [x] Restores idea source from `archive/ideas/`, resets status, drops `promoted_to:`
- [x] Prompts: remove scaffolded request or keep (default: keep, per decision D9)

### M4 — Guardrails + clarity
- [x] No breadcrumb → "Nothing to undo", exit 0
- [x] `--dry-run` prints the reversal, mutates nothing (all three ops)
- [x] Staleness check: refuses if any affected file's mtime > breadcrumb ts (timestamp-vs-mtime, D9)
- [x] Confirm prompt before removing the scaffolded request (idea-promote undo)
- [x] `tests/cli/undo.test.sh` covering M1–M4 (30 assertions, all pass)
- [x] Skill hints in lifecycle.md + archive.md (↩ revert tier-reveal)
- [x] SPEC.md CLI-verb bullet synced (advance/undo/next); fixed dangling [[verification]] links from b13

## v2 (deferred)

- [ ] Multi-level undo (breadcrumb stack instead of single)
- [ ] Full git-state verification (refuse when working tree moved since the mutation)
