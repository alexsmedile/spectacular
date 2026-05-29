---
status: planned
updated: 2026-05-29
related:
  - PLAN.md
---

# Tasks — cli-debt-removal

## v1

### M1 — Audit + ADR
- [ ] Enumerate every deprecated surface (verbs, --global, docs-* refs, banner machinery)
- [ ] Write DECISIONS ADR: MINOR classification + exact removal list

### M2 — Remove verbs + banner
- [ ] Delete `docs init|export|new|review|status` verbs
- [ ] Remove `deprecation_notice()` machinery
- [ ] Preserve pageworks install hint where docs verbs used to point

### M3 — Remove refs + alias
- [ ] Delete `docs-contract` / `docs-rules` / `docs-renderer-adapters` reference docs
- [ ] Remove legacy back-compat PRD references
- [ ] Remove `--global` alias (keep `--skill-scope global`)

### M4 — Surface + tests
- [ ] Update `--help` + usage text for trimmed surface
- [ ] Update test suite for removed surface
- [ ] Confirm `doctor docs` (discovery-only) still passes

### M5 — Ship
- [ ] CHANGELOG [1.13.0] Removed entry; plugin bump to v1.13.0

## v2 (deferred)

- [ ] Verb/flag renames (only if any awkward surface genuinely needs it — fold into v2.0.0)
