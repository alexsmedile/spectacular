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

`<filename>@v<version>.<ext>`

Examples:
- `PRD.md` → snapshot as `PRD@v1.0.md`
- `STACK.md` → snapshot as `STACK@v1.2.md`
- `current/auth/login.md` → snapshot as `current/auth/login@v1.0.md`

Snapshots live **alongside** the current file (same directory).

The unversioned filename (`PRD.md`) always points to the **current/latest** version.

---

## Snapshot sequence

1. Read the current file's frontmatter to get the current `version`
2. Copy file to `<name>@v<version>.<ext>` — do not modify the snapshot
3. Edit the new version of the unversioned file
4. Increment `version` in frontmatter (patch for small edits, minor for meaningful changes)
5. Update `updated` date in frontmatter

---

## Manual snapshot

User can run: `spectacular snapshot <file>`

Skill will:
1. Read current version from frontmatter
2. Create snapshot
3. Confirm to user: "Snapshotted `PRD.md` as `PRD@v1.1.md`. Ready to edit."

---

## Version bump guidance

| Change type | Bump |
|---|---|
| Minor corrections, wording | patch (1.0 → 1.1) |
| New section, significant update | minor (1.0 → 1.1, 1.1 → 1.2) |
| Major restructure or rewrite | major (1.x → 2.0) |

This is a soft guideline — the human decides what constitutes a major change.
