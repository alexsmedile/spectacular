---
status: verified
updated: 2026-06-30
related:
  - PLAN.md
---

# Tasks — snapshot-retention

## v1

### M1 — Allowlist: add DESIGN.md
- [ ] Add `DESIGN.md` to `is_canonical_doc` root-doc case (cli:1611)
- [ ] Test: `is_canonical_doc .spectacular/DESIGN.md` → 0; non-canonical → 1

### M2 — Version coupling (stop @vN ↔ version: drift)
- [ ] `cmd_snapshot`: when `version:` parseable, filename = `@v<current_version>.<ext>` (content's version, decision 1), then bump live doc; else integer counter
- [ ] Keep idempotence (body-only md5) + `--major` path working
- [ ] Tests: doc at 1.3 → `@v1.3.md` + live doc 1.4; DESIGN.md (no fm) → counter

### M3 — Tiered retention config + prune
- [ ] `config_snapshots_keep()` + `config_snapshots_period()` parsers; defaults 3 / `month`
- [ ] Tier union: origin (`@v1`) ∪ periodic (newest per `updated:`-date bucket, month/week/off) ∪ recent (newest `keep` by `@vN`) (decision 4)
- [ ] Bucket keys from snapshot frontmatter `updated:` (substring for month; `date` for week, guard BSD/GNU) — never mtime
- [ ] `cmd_snapshot_prune` / `doctor --fix snapshots`: move outside-union → `.trash/` (git-rm if tracked), dry-run default, live doc untouched
- [ ] Tests: `period: off` `@v1..@v6` keep 3 → keep `@v1 @v4 @v5 @v6`; multi-month → one survivor/month kept outside recent window

### M4 — Folder rename + configurable name + gitignore
- [ ] Collapse 15 hardcoded `snapshots/` refs → one `$snap_root` from `config_snapshots_folder()` (default `_snapshots`)
- [ ] `config_snapshots_folder()` + `config_snapshots_gitignore()` parsers
- [ ] `check_snapshots` detects old `snapshots/` vs configured `_snapshots/`; `--fix` renames (git-mv) losslessly
- [ ] `gitignore: true` → init + `doctor --fix` add `.spectacular/<folder>/` ignore line; default false → not ignored (decision 2)
- [ ] Ship step: migrate THIS repo's `snapshots/` → `_snapshots/` + commit (decision 3, dogfood)
- [ ] Tests: default lands in `_snapshots/`; migration lossless; gitignore toggle

### M5 — Doctor retention check + docs
- [ ] `check_snapshots`: flag docs exceeding `keep` (info/warning), relayed by `status`
- [ ] Docs: commands.md, configuration.md (`snapshots:` block — folder/keep/gitignore), ARCHITECTURE.md versioning + path mentions, SPEC.md + specs/cli/SPEC.md
- [ ] VERIFY-LOG

## v2 (deferred)
- [ ] Re-evaluate snapshots vs git history (do we keep the mechanism at all?) — separate TODO item
