---
status: review
updated: 2026-07-07
related:
  - PLAN.md
---

# Tasks — archive-closure-gate

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

### M1 — Skill-side flow (archive.md + spec-sync.md)
- [x] Rewrite `spec-sync.md` § Proposal format to structured deltas (ADDED / MODIFIED quotes current+replacement / REMOVED quotes bullet / NONE — <why>)
- [x] Add closure-gate step to `archive.md` § Archive sequence (before the CLI verb) + override etiquette
- [x] Dry-run the flow on a fixture request in a scratch workspace — fulfilled by the M2/M3 scratch smoke tests + tests/cli/archive-closure-gate.test.sh scenario matrix
- [x] → check: proposal output contains delta sections, not free-form prose

### M2 — CLI gate (cmd_archive)
- [x] TASKS closure check: every box `[x]` or `[~]` with a reason; unexplained `[ ]` blocks with the box named
- [x] VERIFY closure check: if VERIFY.md exists, VERIFY-LOG.md must contain ≥1 walk entry (✅ row)
- [x] Spec-delta presence check + `--override <check> --reason "<text>"` recording into `archive_overrides:` frontmatter
- [x] → check: blocked archive exits non-zero with actionable message; override run exits 0 and records the reason

### M3 — Doctor + tests
- [x] `doctor specs`: delta-integrity validation (MODIFIED/REMOVED quote an existing bullet; ADDED not a duplicate)
- [x] Tests: gate-blocks, override-recording, NONE delta, grandfathered legacy archive untouched
- [x] → check: run: full tests/ suite green; doctor exits non-zero on a bad fixture delta

## v2 (deferred)

- [ ] Auto-suggest delta content from the request's `related:` capabilities (skill-side convenience)
- [ ] `doctor lifecycle` non-vocab `status:` flag with `note:` remediation (companion to lifecycle.md rule shipped 2026-07-07)
