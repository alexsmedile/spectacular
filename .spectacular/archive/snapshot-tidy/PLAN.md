---
status: archived
priority: medium
owner: alex
updated: 2026-05-25
target_version: v1.5.0
summary: "Move canonical-doc snapshots (PRD@v1.2.md, ROADMAP@v4.md, etc.) out of .spectacular/ root into snapshots/<DOC>/v<x.y>.md to keep the workspace root clean. Add doctor area + --fix migration; update snapshot verb to write to new location."
related:
  - ../../roadmaps/index.md
  - ../../ARCHITECTURE.md
  - ../../specs/index.md
archived: 2026-05-25
---

# Plan — snapshot-tidy

## 1. Goal

Today `spectacular snapshot <file>` writes versioned snapshots next to the canonical doc — `.spectacular/PRD@v1.2.md` sits beside `.spectacular/PRD.md`. After a few iterations the workspace root becomes noisy: this repo already has `PRD@v1.2.md` and `ROADMAP@v4.md` mixed into the root listing, and the pattern will compound as `memory-protocols` and future work generate more snapshots.

This request introduces a dedicated `snapshots/` subtree:

```
.spectacular/
├── PRD.md                  # canonical, unchanged
├── ROADMAP.md              # canonical, unchanged
└── snapshots/
    ├── PRD/
    │   ├── v1.0.md
    │   └── v1.2.md
    └── ROADMAP/
        └── v4.md
```

One folder per canonical doc; snapshot filenames lose the `@v` and just become `v<x.y>.md`. The canonical file at root still points to "current" — that contract is unchanged.

## 2. Constraints

- **Backwards compatibility for one minor.** Doctor warns on root-level `*@v*.md` files; CLI auto-migrates on next `snapshot` call OR via `doctor --fix snapshots`. Root-level snapshots remain readable until cleaned up — never silently deleted.
- **No content changes.** Migration is filesystem-only: rename + move. Frontmatter, body, links untouched.
- **Composes with versioning ref.** `skills/spectacular/references/versioning.md` is the source of truth for the snapshot contract — update it as part of this request.
- **Doctor area is mechanical.** `doctor snapshots` (read-only) lists drift; `--fix` performs the migration. No skill-side judgment needed.

## 3. Scope

### In
- New layout: `.spectacular/snapshots/<DOC>/v<x.y>.md`
- `spectacular snapshot <file>` writes to new location
- `spectacular doctor snapshots` area (warns on root-level `*@v*.md`)
- `spectacular doctor --fix snapshots` migrates them
- Update `versioning.md` reference doc to describe new layout
- Update `ARCHITECTURE.md` directory tree
- Migrate this repo's own snapshots (dogfood) — `PRD@v1.2.md` → `snapshots/PRD/v1.2.md`, `ROADMAP@v4.md` → `snapshots/ROADMAP/v4.md`

### Out
- Bulk rename of historical snapshots in archived requests (those stay in their `archive/<slug>/` subtree untouched)
- Compression / pruning of old snapshots (separate request if needed)
- Cross-doc snapshot indexing (no `snapshots/INDEX.md`) — folder listing is enough

## 4. Locked decisions (grilled 2026-05-25)

1. **Folder casing:** **uppercase** — matches the canonical filename stem (`snapshots/PRD/`, `snapshots/ROADMAP/`). Visual parity with root file.
2. **Sub-doc snapshots:** **mirror path** — `snapshots/specs/cli/SPEC/@v1.0.md`. Preserves directory structure; no slug collisions.
3. **Filename:** **keep the `@`** — final filename is `@v<N>.md` (e.g. `@v1.2.md`). Preserves the existing `@`-prefix convention; grep-friendly across both old and new layouts.
4. **Doctor severity:** **warn** during migration window (v1.6 → v1.7); demote to info in v1.7.

### Final layout

```
.spectacular/
├── PRD.md
├── ROADMAP.md
└── snapshots/
    ├── PRD/
    │   ├── @v1.0.md
    │   ├── @v1.1.md
    │   ├── @v1.2.md
    │   └── @v1.3.md
    ├── ROADMAP/
    │   ├── @v1.md
    │   ├── @v2.md
    │   ├── @v3.md
    │   └── @v4.md
    └── specs/
        └── cli/
            └── SPEC/
                └── @v1.0.md
```

## 5. Milestones

### M1 — Migration spec + dogfood
- Lock layout decisions from §4
- Update `versioning.md` with new layout + migration notes
- Update `ARCHITECTURE.md` directory tree
- Migrate this repo's snapshots (`PRD@v1.2.md`, `ROADMAP@v4.md`) by hand as the reference case
- Update `CLAUDE.md` if it references snapshot layout

### M2 — CLI snapshot verb update
- `spectacular snapshot <file>` writes to `.spectacular/snapshots/<DOC>/v<x.y>.md`
- Auto-creates `snapshots/<DOC>/` if missing
- Existing snapshot reading paths (if any) — audit and update

### M3 — Doctor area
- New `doctor snapshots` area: scans `.spectacular/*@v*.md` and `.spectacular/specs/**/*@v*.md`
- Reports each as warn with suggested target path
- `doctor --fix snapshots` performs the moves (using `git mv` when in a repo)

### M4 — Release
- CHANGELOG entry
- Bump manifests to 1.6.0
- Update `scripts/hooks/pre-commit` only if it touches snapshot paths

## 6. Non-goals (locked)

- Deleting old snapshots
- Compression / archiving of snapshots
- Snapshot diff tooling (`spectacular diff PRD v1.0 v1.2`) — separate idea
- Adding snapshot-driven changelog generation

## 7. Dependencies

- None — pure substrate hygiene. Independent of `memory-protocols` (v1.6.0) and can ship in same or different minor.

## 8. References

- `skills/spectacular/references/versioning.md` — current snapshot contract
- `skills/spectacular/references/doctor.md` — doctor area pattern
- `.spectacular/ARCHITECTURE.md` § directory tree
