---
description: Judgment-fix repair flow — y/n/q walk, snapshot-before-edit, worked examples.
when_to_use: Applying a non-mechanical doctor fix that needs judgment.
---

# Doctor — judgment-fix repair flow

Loaded when the skill is invoked as `/spectacular doctor --fix`. Entry point: [[doctor]]; check-area reference: [[doctor-areas]]; auto-invocation triggers: [[doctor-substrate]].

## Repair flow

1. **Detect** — invoke `spectacular doctor --format json` and capture the report. Optionally re-use `.spectacular/.doctor-report.json` if the user just ran it.
2. **Filter** — keep only findings where `fix_type == "judgment"`. Mechanical findings are the CLI's job and have already been applied if the user ran `spectacular doctor --fix` first.
3. **Group by area** — present in order: `skill → workspace → frontmatter → snapshots → links → lifecycle → kits → conventions → specs → docs`. Within an area, errors before warnings before info.
4. **Per finding** — propose a context-aware fix. Show the user:
   - Which file
   - What's wrong (the `message`)
   - What I'd do to fix it (concrete, not abstract)
   - The `[y/n/q]` prompt
5. **On `y`** — snapshot first if the file is canonical (per [[versioning]]), then apply the edit. Confirm with `✓ fixed: <what changed>`.
6. **On `n`** — skip silently, move to next finding.
7. **On `q`** — print remaining-findings summary + exit. Already-applied fixes stand.
8. **Final** — re-run `spectacular doctor` to confirm clean state. Summarize what was fixed vs skipped.

## Snapshot-before-edit is mandatory

For every fix touching a **canonical doc** (per `doc-index.md`'s `snapshot-on-edit: true`):
1. Read current `version:` from frontmatter
2. Copy file to `<doc>@v<current>.md` if that snapshot doesn't already exist
3. Apply the edit
4. Bump the file's `version:` field if the edit was substantive

Per-request files (PLAN.md, TASKS.md, VERIFY.md) typically have `snapshot-on-edit: false` — no snapshot needed. Still read the file before editing and surface a diff to the user.

## Examples — context-aware proposals

### Frontmatter — missing required field

```
> .spectacular/STACK.md missing `version:` field (warning)

  The registry expects `version: 1.0` for canonical root docs. Existing
  frontmatter has `updated:` and `summary:` but no `version:`. I can add
  `version: 1.0` as the first frontmatter field.

  [y] apply this fix      [n] skip      [q] quit doctor
> y

  Snapshotting .spectacular/STACK.md → STACK@v0.9.md ... ✓
  Editing .spectacular/STACK.md ... ✓
  → fixed: added version: 1.0
```

### Snapshots — version gap

```
> .spectacular/PRD@v1.1.md missing — version-sequence gap (warning)

  Git log shows no commit ever referenced PRD@v1.1. Two interpretations:
  (a) v1.1 was never created — version bump skipped. Most likely.
  (b) v1.1 existed but was deleted before commit. Unrecoverable.

  For (a), the canonical fix is to document the skip in DECISIONS.md
  rather than fabricate content for a version that never existed.

  [y] append a DECISIONS entry explaining the skip
  [n] skip this finding (leave the gap unaddressed)
  [q] quit doctor
> y

  Appending to .spectacular/decisions/index.md ... ✓
  → fixed: gap acknowledged
```

### Links — broken `related:` target

```
> .spectacular/requests/foo/PLAN.md related: "../bar/PLAN.md" target missing (warning)

  The path resolves to .spectacular/requests/foo/../bar/PLAN.md
  which doesn't exist. Two options:
  (a) The path was meant to be a different request → fix the path
  (b) The referenced request was archived or never created → remove the link

  I checked .spectacular/requests/ — there's no `bar/` and no `archive/bar/`.
  Recommend (b): remove the stale link.

  [y] remove the link from related:    [n] skip    [q] quit
> y

  Editing .spectacular/requests/foo/PLAN.md ... ✓
  → fixed: removed "../bar/PLAN.md" from related:
```

### Lifecycle — stale verified request

```
> .spectacular/requests/old-feature/ verified 14 days ago — archive candidate (info)

  This request reached `verified` status 14 days ago without being archived.
  Per the lifecycle convention, verified requests should move to archive/
  unless they're being actively referenced.

  [y] run `spectacular archive old-feature` flow now
  [n] leave it; mark not-yet-ready-to-archive
  [q] quit doctor
> y

  → routing to references/archive.md for the archive sequence
```

## Anti-patterns

- **Bulk-applying without per-finding confirm** — even if the user typed `y y y y`, surface each finding individually
- **Inventing content** — for missing fields with no canonical default, ask the user or punt with `[NEEDS CLARIFICATION]`
- **Modifying snapshots** — `<doc>@vN.md` files are immutable. Never edit them.
- **Editing during read-only mode** — if user invoked `/spectacular doctor` without `--fix`, never propose writes; just walk the report conversationally

## Related

- [[doctor]] — entry point, severity model
- [[doctor-areas]] — what each check actually checks
- [[doctor-substrate]] — auto-invocation triggers from other skill flows
- [[versioning]] — snapshot rule every judgment fix follows
- [[archive]] — where stale-verified findings route on `y`
