# Versioning — Snapshot Before Edit

Canonical documents are **never overwritten in place**. Always snapshot first.

---

## What counts as canonical

- Root layer files: `PRD.md`, `PRINCIPLES.md`, `ARCHITECTURE.md`, `ROADMAP.md`, `AGENTS.md`, `STACK.md`, `DECISIONS.md`
- `SPEC.md` (always-on index) + `specs/<capability>/SPEC.md` (per-capability)
- `config.yaml`

Requests files (`PLAN.md`, `TASKS.md`, `SESSION.md`) are operational/temporary — no snapshot required.

---

## Snapshot naming

`<filename>@v<N>.<ext>` where `<N>` is a monotonically increasing integer.

Examples:
- `PRD.md` → snapshot as `PRD@v1.md`, then `PRD@v2.md`, ...
- `STACK.md` → snapshot as `STACK@v1.md`
- `specs/auth/SPEC.md` → snapshot as `specs/auth/SPEC@v1.md`

Snapshots live **alongside** the current file (same directory).

The unversioned filename (`PRD.md`) always points to the **latest** version.

> **Note (v0.7.0+):** older snapshots in the repo use dotted version naming (`PRD@v1.0.md`) — that was the v0.7.x convention. The v0.7.0 CLI verb uses integer naming (`PRD@v1.md`). Both coexist; the CLI scans for any numeric suffix when picking the next N.

---

## Snapshot sequence (v0.7.0+ via CLI verb)

Use **`spectacular snapshot <file>`** — never do this by hand. The CLI verb:

1. Validates `<file>` is a registered canonical doc; refuses otherwise
2. Scans for existing `<base>@v*.md` snapshots; picks next N
3. Compares current file body to latest snapshot — if unchanged, exits cleanly (idempotent)
4. Copies current state to `<base>@v<N>.md`
5. Bumps `version:` field in the unversioned file (minor by default: `X.Y` → `X.(Y+1)`; `--major` for `(X+1).0`)
6. Sets `updated:` to today

Manual snapshotting (cp + sed) is fragile and gets the version bump wrong. The verb has tests; ad-hoc shell doesn't.

---

## Version bump guidance

| Change type | Bump |
|---|---|
| Minor corrections, wording | patch (1.0 → 1.1) |
| New section, significant update | minor (1.0 → 1.1, 1.1 → 1.2) |
| Major restructure or rewrite | major (1.x → 2.0) |

This is a soft guideline — the human decides what constitutes a major change.
