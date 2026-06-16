---
status: review
updated: 2026-06-16
related:
  - PLAN.md
---

# Tasks — cli-debt-removal

## v1

### M1 — Audit + ADR
- [x] Enumerate every deprecated surface (verbs, --global, docs-* refs, banner machinery)
- [x] Write DECISIONS ADR: MINOR classification + exact removal list (D6)

### M2 — Remove verbs + banner
- [x] Delete `docs init|export|new|review|status` verbs
- [x] Remove `deprecation_notice()` machinery
- [x] Preserve pageworks install hint where docs verbs used to point

### M3 — Remove refs + alias
- [x] Delete `docs-contract` / `docs-rules` / `docs-renderer-adapters` reference docs
- [x] Remove legacy back-compat PRD references
- [x] Remove `--global` alias (keep `--skill-scope global`)

### M4 — Surface + tests
- [x] Update `--help` + usage text for trimmed surface
- [x] Update test suite for removed surface (deleted docs*.test.sh; updated init + mutator + visual tests)
- [x] Confirm `doctor docs` (discovery-only) still passes

### M5 — Ship
- [x] CHANGELOG [1.17.0] Removed entry; plugin bump already at v1.17.0

## v2 (deferred)

- [ ] Verb/flag renames (only if any awkward surface genuinely needs it — fold into v2.0.0)
