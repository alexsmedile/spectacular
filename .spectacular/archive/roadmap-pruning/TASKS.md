---
status: verified
updated: 2026-06-29
related:
  - PLAN.md
---

# Tasks — roadmap-pruning

> Depends on b17 (roadmap-contract-docs) — ledger must be specced first.

## v1

### M1 — Decide approach + spec it
- [x] Approach B (roadmap-index mode), keep newest 3 inline — resolved at design call (decisions 1-2)
- [x] specs/roadmap/SPEC.md § "Index mode" + ARCHITECTURE.md ledger rule updated (enforces "history → CHANGELOG"); spec snapshotted @v2, 1.2

### M2 — Detection (doctor)
- [x] New `doctor roadmap` area: orphan index lines, stale per-version files, flat/index prune nudge beyond keep-window. Registered in DOC_AREAS + dispatch + arg validation.

### M3 — Prune mechanism
- [x] `spectacular roadmap migrate [--dry-run] [--keep N]`: snapshot-safe, writes roadmap/v*.md before rewriting ROADMAP, ## Shipped index, dry-run default, idempotent (bash 3.2 safe)
- [x] Dogfood: migrated this repo's ROADMAP (528 → 410 lines; 7 oldest shipped → roadmap/v*.md; removed "Recently shipped" mirror + stale reconciliation notes; fixed roadmap-overrides→roadmap-rules refs)
- [x] tests/cli/roadmap-migrate.test.sh (22 assertions, all pass) + docs (commands.md, doctor area lists in commands.md/CLAUDE.md)
- [x] VERIFY-LOG
