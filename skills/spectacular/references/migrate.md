---
description: The /spectacular migrate skill walk — judgment migrations with snapshot + y/n/q per step.
when_to_use: spectacular migrate, or a workspace_schema upgrade.
---

# `/spectacular migrate` — Skill walk for judgment migrations

Loaded when the user invokes `/spectacular migrate` in their AI agent (Claude Code, Codex, etc.).

This file walks judgment migrations — the ones the CLI explicitly refused with "mechanical: false, run /spectacular migrate in your AI agent". For purely mechanical migrations, the CLI handles everything; the skill never gets involved.

## When to invoke

Three triggers:

1. **CLI told the user to**: `spectacular migrate` printed `⚠ <id>: this migration requires judgment` and exited non-zero
2. **User asked explicitly**: "walk me through any judgment migrations"
3. **Onboarding detected gap**: skill ran `spectacular status --against-latest` during onboarding and the chain includes a judgment migration

## Preflight (before walking any migration)

1. **Confirm location**: must be inside a Spectacular workspace (`.spectacular/` exists). If not, stop and instruct the user to `cd` to a workspace root.
2. **Read current state**: parse `workspace_schema:` from `config.yaml` (or assume `"0.4"` if absent).
3. **Read target**: the CLI's `CURRENT_SCHEMA` constant (run `spectacular migrate --list` and pick the highest `to:` value as the target).
4. **List pending migrations**: run `spectacular migrate --list` to see the full chain. Filter to those where `from >= current_schema AND to <= target_schema`. Split into mechanical vs judgment.
5. **Apply mechanical first**: if any mechanical migrations precede the judgment one, run `spectacular migrate --to <ver-just-before-judgment>` first. Don't try to walk judgment when the workspace isn't at the expected `from:` schema.
6. **Surface the plan**: tell the user how many judgment migrations are pending and what they do (read each migration's `description:` + body's `## Steps` section).

## Walk pattern (per judgment migration)

For each pending judgment migration, in order:

### Step 1 — Read the migration spec

Load `skills/spectacular/references/migrations/<id>.md`. Extract:
- Frontmatter: `id`, `from`, `to`, `description`, `reversible`, `affects`
- Body: `## Detection`, `## Steps`, `## Rollback`, `## Validation` sections

### Step 2 — Confirm prerequisites

Read the `## Detection` section. State to the user: "this migration applies because `<predicate>`". If the predicate doesn't hold, say so and skip (workspace doesn't need this migration).

### Step 3 — Snapshot affected canonical docs

For every path in `affects:` that's a canonical doc (PRD.md, SPEC.md, ARCHITECTURE.md, PRINCIPLES.md, ROADMAP.md, STACK.md, AGENTS.md, DECISIONS.md, or `specs/<cap>/SPEC.md`):

Run `spectacular snapshot <path>` BEFORE any edit. This creates `<file>@v<N>.md` so the user can compare before/after and revert if needed.

For non-canonical paths (directories, config.yaml), no snapshot — those are either mechanical (CLI handles) or unwritable (read-only validation paths).

### Step 4 — Walk steps with y/n/q

For each numbered step in `## Steps`:

1. Read the step aloud (paraphrased for context if helpful).
2. Show the planned change: file path + diff or replacement content.
3. Ask: `Apply this step? [y/n/q]`
   - `y` → apply, mark step done, advance
   - `n` → skip step, log to `.spectacular/memory/migration-skips-<date>.md`, advance (NOT recommended; surface the risk: "skipping may leave the workspace in an inconsistent state. Continue? [y/n]")
   - `q` → stop walking; leave workspace in whatever partial state; print final state + instructions to resume

### Step 5 — Bump workspace_schema

After all steps applied, write the new `workspace_schema: "<to>"` value to `.spectacular/config.yaml`.

(This is the same operation the CLI does via `bump_workspace_schema`. The skill can shell out: `python3 -c "..."` or run `sed` directly.)

### Step 6 — Run validation

Read the `## Validation` section. Walk each validation predicate; report PASS/FAIL. If any FAIL, instruct the user how to fix (point at `## Rollback` if available) and stop the chain.

### Step 7 — Advance to next pending migration

If more judgment migrations are pending, go to Step 1 with the next migration. If none, run a final `spectacular doctor` to confirm overall workspace health and report.

## Memory write

After a successful judgment walk, write to `.spectacular/memory/migrations-<date>.md`:

```yaml
---
type: migration-log
date: <YYYY-MM-DD>
migrations:
  - id: <id>
    from: <from>
    to: <to>
    walked-via: skill
    snapshots-created: [PRD@v2.md, ...]
    steps-skipped: []
---

# Migration log — <date>

<one paragraph: what was done, anything notable>
```

This gives the team an audit trail and survives the session.

## Stage 2 has no judgment migrations

As of v0.6.2, no migration ships with `mechanical: false`. The skill walk is shipped as scaffolding for future migrations that need it (e.g. content rewrites, conflict-aware merges). This file is the spec; when the first judgment migration is added, the maintainer follows the contract above and the skill is already wired.

## Failure modes

| Symptom | Cause | Response |
|---|---|---|
| `spectacular migrate --list` shows no migrations | Skill not installed in this workspace | Run `spectacular init --update` |
| Detection predicate fails on a "pending" migration | Workspace already past this migration | Bump `workspace_schema:` to its `to:` value; skip step 4-6 |
| Snapshot fails | File doesn't exist or perm error | Stop; instruct user to fix; don't proceed without snapshot |
| Validation fails after apply | Migration spec is wrong OR user edited mid-walk | Stop; load `## Rollback`; restore from snapshot; report |
| User runs `/spectacular migrate` with no pending migrations | Workspace is up to date | Reply: "workspace is up to date (schema X)"; suggest `doctor` if they wanted a check |

## Related

- [[migrations-contract]] — schema for migration .md files
- [[doctor]] — chain validation in `check_kits` area
- [[doctor-substrate]] — onboarding flow triggers migrate-check
- [[versioning]] — snapshot rules for canonical docs (used in Step 3)
- [[memory]] — write-back pattern (used in Memory write section)
