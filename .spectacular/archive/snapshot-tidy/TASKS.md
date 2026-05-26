---
status: planned
updated: 2026-05-25
related:
  - PLAN.md
---

# Tasks — snapshot-tidy

## M1 — Migration spec + dogfood ✅

- [x] Lock answers to PLAN §4 open questions (uppercase folders, mirror sub-doc paths, keep `@` in filename, warn severity)
- [x] Update `skills/spectacular/references/versioning.md` with new layout + migration notes
- [x] Update `.spectacular/ARCHITECTURE.md` directory tree to show `snapshots/`
- [x] Migrate all 11 root snapshots into `snapshots/<DOC>/@v<N>.md` (PRD ×4, ROADMAP ×4, AGENTS, ARCHITECTURE, SPEC)
- [ ] Update `CLAUDE.md` if it references snapshot layout (none found — skipped)

## M2 — CLI snapshot verb update

- [ ] Update `cli/spectacular snapshot <file>` to write to `.spectacular/snapshots/<DOC>/v<x.y>.md`
- [ ] Auto-create `snapshots/<DOC>/` if missing
- [ ] Audit codebase for any read paths assuming root-level snapshots; update
- [ ] Add test for snapshot creation in new location

## M3 — Doctor area

- [ ] Add `snapshots` to `doctor` area list in `cli/spectacular`
- [ ] Implement scan: `.spectacular/*@v*.md` + `.spectacular/specs/**/*@v*.md`
- [ ] Report each as warn with suggested target path
- [ ] Implement `--fix snapshots`: `git mv` (or `mv` when not in repo) to target paths
- [ ] Update `skills/spectacular/references/doctor.md` with new area
- [ ] Add test for doctor + --fix behavior

## M4 — Release

- [ ] Bump manifests to 1.6.0 (or align with concurrent request)
- [ ] CHANGELOG entry under Added/Changed
- [ ] Tag + push via `/release`
- [ ] Archive this request post-ship

## Deferred

- [ ] Snapshot diff tool (`spectacular diff <doc> <v1> <v2>`)
- [ ] Snapshot pruning / compression
- [ ] Demote `doctor snapshots` severity from warn → info in v1.7
