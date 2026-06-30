---
description: Snapshot-before-edit rules and the <FILE>@vN.md naming convention.
when_to_use: Snapshotting a canonical doc before a substantive edit.
---

# Versioning — Snapshot Before Edit

Canonical documents are **never overwritten in place**. Always snapshot first.

> **@Snapshot policy gate.** Before overwriting a canonical doc, run `spectacular policy @Snapshot` and follow every active policy. The default blocker is `snapshot-before-overwrite`: a `<DOC>@v<N>.md` snapshot must exist first — or you stop. See [policy-injection.md](policy-injection.md).

---

## What counts as canonical

- Root layer files: `PRD.md`, `PRINCIPLES.md`, `ARCHITECTURE.md`, `ROADMAP.md`, `AGENTS.md`, `STACK.md`, `DECISIONS.md`
- `SPEC.md` (always-on index) + `specs/<capability>/SPEC.md` (per-capability)
- `DESIGN.md` — snapshot-able even though Spectacular doesn't otherwise manage it (b16). Usually has no `version:` field, so it uses the counter (see below).
- `config.yaml`

Requests files (`PLAN.md`, `TASKS.md`, `SESSION.md`) are operational/temporary — no snapshot required.

---

## Snapshot location + naming

Snapshots live in a dedicated tree: `.spectacular/<store>/<DOC>/@v<ver>.md`.

- **Store dir** is configurable — `config.yaml`'s `snapshots.folder`, default `_snapshots` (since v1.24.0; the pre-v1.24 default was `snapshots/`). The `_` prefix marks it a non-content layer, like `_archive/`.
- Folder name **matches the canonical filename stem** (uppercase preserved): `_snapshots/PRD/`, `_snapshots/ROADMAP/`.
- **Filename couples to the version the content IS** (v1.24.0+): a doc at `version: 1.3` snapshots to `@v1.3.md`, *then* the live doc bumps to `1.4`. The `@v` label and the `version:` field never drift apart.
- **Version-less docs use a plain `@v<N>` counter** and are not version-bumped: `DESIGN.md` (no `version:` frontmatter) → `@v1.md`, `@v2.md`, … The live doc's frontmatter is left untouched.
- Sub-doc snapshots **mirror their path**: `specs/cli/SPEC.md` → `_snapshots/specs/cli/SPEC/@v1.0.md`. Avoids slug collisions.

Examples:
- `PRD.md` at `1.3` → `_snapshots/PRD/@v1.3.md` (live doc → `1.4`)
- `DESIGN.md` (no version) → `_snapshots/DESIGN/@v1.md`, then `@v2.md`, …
- `specs/auth/SPEC.md` at `1.0` → `_snapshots/specs/auth/SPEC/@v1.0.md`

The unversioned filename at root (`PRD.md`) always points to the **latest** version.

### Migration from older layouts

- **Root-level legacy** (`PRD@v1.2.md`, pre-v1.6): still read; `doctor snapshots` warns until you `spectacular doctor --fix snapshots` to git-mv them into the tree.
- **Folder rename** (`snapshots/` → configured `_snapshots/`, b16): `doctor snapshots` flags it; `doctor --fix snapshots` renames the dir (git-mv when tracked) losslessly.

---

## Snapshot sequence (via CLI verb)

Use **`spectacular snapshot <file>`** — never do this by hand. The CLI verb:

1. Validates `<file>` is a registered canonical doc; refuses otherwise
2. Resolves the store dir from config; scans existing snapshots (new tree + legacy) for the counter fallback + idempotence
3. Compares current file body to latest snapshot — if unchanged, exits cleanly (idempotent)
4. Names the snapshot `@v<current-version>` when the doc has a parseable `version:`, else the next `@v<N>` counter
5. Copies current state into `<store>/<DOC>/@v<ver>.md` (creating the dir if missing)
6. Bumps `version:` in the live doc (minor by default; `--major` for `(X+1).0`) — **skipped for version-less docs**
7. Sets `updated:` to today

Manual snapshotting (cp + sed) is fragile and gets the version bump wrong. The verb has tests; ad-hoc shell doesn't.

---

## Retention + prune (v1.24.0+)

Snapshots are bounded by **tiered, generational retention** — a snapshot kept by *any* tier survives:

- **origin** — always keep the first snapshot (`@v1`).
- **periodic** — keep the newest snapshot per calendar bucket (`snapshots.period`: `month` default, `week`, or `off`). Bucketed by each snapshot's `updated:` frontmatter date — stable across clone, unlike mtime.
- **recent** — keep the newest `snapshots.keep` (default 3) by version ordinal.

Run **`spectacular snapshot prune`** (dry-run) → **`--apply`** to remove snapshots no tier claims. Tracked files are `git rm`'d (history holds them); untracked / non-git fall back to `.spectacular/.trash/`. The live canonical doc is never touched. `doctor snapshots` surfaces an info nudge when prunable snapshots accumulate.

This bounds a doc to roughly `1 + periods_alive + keep` snapshots instead of unbounded growth.

### gitignoring the store

`snapshots.gitignore` (default `false`) controls whether the store is committed. Set `true` and run `spectacular doctor --fix snapshots` (or re-`init`) to add `.spectacular/<store>/` to `.gitignore`; set back to `false` and `--fix` removes the line.

---

## Version bump guidance

| Change type | Bump |
|---|---|
| Minor corrections, wording | patch (1.0 → 1.1) |
| New section, significant update | minor (1.0 → 1.1, 1.1 → 1.2) |
| Major restructure or rewrite | major (1.x → 2.0) |

This is a soft guideline — the human decides what constitutes a major change.
