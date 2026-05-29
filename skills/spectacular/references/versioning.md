---
description: Snapshot-before-edit rules and the <FILE>@vN.md naming convention.
when_to_use: Snapshotting a canonical doc before a substantive edit.
---

# Versioning — Snapshot Before Edit

Canonical documents are **never overwritten in place**. Always snapshot first.

---

## What counts as canonical

- Root layer files: `PRD.md`, `PRINCIPLES.md`, `ARCHITECTURE.md`, `ROADMAP.md`, `AGENTS.md`, `STACK.md`, `DECISIONS.md`
- `SPEC.md` (always-on index) + `specs/<capability>/SPEC.md` (per-capability)
- `config.yaml`

Requests files (`PLAN.md`, `TASKS.md`, `SESSION.md`) are operational/temporary — no snapshot required.

---

## Snapshot location (v1.6.0+)

Snapshots live in a dedicated tree: `.spectacular/snapshots/<DOC>/@v<N>.md`.

- Folder name **matches the canonical filename stem** (uppercase preserved): `snapshots/PRD/`, `snapshots/ROADMAP/`.
- Filename keeps the `@` prefix: `@v1.0.md`, `@v1.2.md`, `@v2.md`. Mixed integer + dotted versions coexist (the CLI scans any numeric suffix to pick the next N).
- Sub-doc snapshots **mirror their path**: `specs/cli/SPEC.md` → `snapshots/specs/cli/SPEC/@v1.0.md`. Avoids slug collisions.

Examples:
- `PRD.md` → `snapshots/PRD/@v1.md`, then `snapshots/PRD/@v2.md`, ...
- `STACK.md` → `snapshots/STACK/@v1.md`
- `specs/auth/SPEC.md` → `snapshots/specs/auth/SPEC/@v1.md`

The unversioned filename at root (`PRD.md`) always points to the **latest** version.

### Migration from pre-v1.6 layout

Before v1.6.0, snapshots lived alongside the canonical file as `PRD@v1.2.md`. Those still work — the CLI reads both locations — but `spectacular doctor snapshots` warns on root-level `*@v*.md` files until you migrate. Run `spectacular doctor --fix snapshots` to `git mv` them into the new tree.

The warning demotes to info in v1.7.0.

---

## Snapshot sequence (v0.7.0+ via CLI verb)

Use **`spectacular snapshot <file>`** — never do this by hand. The CLI verb:

1. Validates `<file>` is a registered canonical doc; refuses otherwise
2. Scans for existing snapshots in **both** `snapshots/<DOC>/@v*.md` (v1.6+) and `<base>@v*.md` (legacy); picks next N from the union
3. Compares current file body to latest snapshot — if unchanged, exits cleanly (idempotent)
4. Auto-creates `.spectacular/snapshots/<DOC>/` if missing
5. Copies current state to `snapshots/<DOC>/@v<N>.md`
6. Bumps `version:` field in the unversioned file (minor by default: `X.Y` → `X.(Y+1)`; `--major` for `(X+1).0`)
7. Sets `updated:` to today

Manual snapshotting (cp + sed) is fragile and gets the version bump wrong. The verb has tests; ad-hoc shell doesn't.

---

## Version bump guidance

| Change type | Bump |
|---|---|
| Minor corrections, wording | patch (1.0 → 1.1) |
| New section, significant update | minor (1.0 → 1.1, 1.1 → 1.2) |
| Major restructure or rewrite | major (1.x → 2.0) |

This is a soft guideline — the human decides what constitutes a major change.
