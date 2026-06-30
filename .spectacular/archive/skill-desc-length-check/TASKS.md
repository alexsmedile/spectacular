---
status: verified
updated: 2026-06-30
related:
  - PLAN.md
---

# Tasks — skill-desc-length-check

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

### M1 — Threshold confirmed
- [x] Confirm what Codex measures: `description` alone vs `description` + `when_to_use` concatenation — **alone** (D7)
- [x] Confirm the exact cap (1024) and decide warning band (>1000)
- [x] Record finding (D7 decision), citing v1.17.2 (1146 → 986)

### M2 — Doctor check
- [x] Add YAML-literal-block parser for `description` to measurement logic
- [x] Wire char-count sub-check into `check_skill()` in `cli/spectacular` (`check_skill_desc_len`)
- [x] Emit `error` >1024, `warning` >1000, `pass` otherwise — include actual count in the message
- [x] Fixture test: SKILL.md at >1024 (error), ~1010 (warning), ≤1000 (pass); live spectacular SKILL.md (986) passes — `tests/cli/doctor.test.sh` scenario 17

### M3 — Pre-commit guard
- [x] Factor measurement into a shared `scripts/check-skill-desc.sh` helper (doctor inlines the same awk; hook sources the script)
- [x] Add guard via `scripts/hooks/pre-commit-wrapper` + `.active/` symlink dir, so git-guard's `pre-commit` is untouched and survives regen
- [x] Test: staged SKILL.md edit >1024 rejected; trimmed-back edit commits clean (verified in throwaway repo)

## v2 (deferred)

- [ ] Generalize beyond this repo's single skill — scan all `**/SKILL.md` in a multi-skill repo
- [ ] Consider a doctor `--fix` that proposes a trimmed description (judgment fix, not mechanical)
