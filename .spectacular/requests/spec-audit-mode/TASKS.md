---
status: planned
updated: 2026-07-06
related:
  - PLAN.md
---

# Tasks — spec-audit-mode

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

### M1 — Orphan capability bullets
- [ ] Extend `check_specs` to walk SPEC.md's Capabilities bullets against `specs/<slug>/SPEC.md` + archive slug/summary mentions
- [ ] Info/warn threshold by bullet length (short bullets are allowed to stay spec-file-less)
- [ ] Test: bullet with no spec file + no archive mention → warning; matching spec file → clean

### M2 — Orphan spec files
- [ ] Extend `check_specs` to flag `specs/<cap>/SPEC.md` files not referenced anywhere in SPEC.md's body
- [ ] Test: unreferenced `specs/ghost.md` → warning; referenced → clean

### M3 — Stale capability specs
- [ ] Extend `check_specs` to compare each `specs/<cap>/SPEC.md`'s `updated` against its newest related archive (via archive PLAN's `related:`)
- [ ] Test: cap spec older than a related archive → warning

### M4 — Docs + JSON summary
- [ ] `doctor specs --json` includes per-signal findings
- [ ] Update `doctor-areas.md` (specs table), `status.md` (signal table), `spec-sync.md` (standalone audit trigger)
- [ ] ROADMAP ledger row mapping build b11 → target version
