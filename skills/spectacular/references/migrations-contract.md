# Migrations Contract — workspace-schema migrations as a registry

Loaded when the CLI's `cmd_migrate` builds the migration chain, or when the skill walks a judgment migration (`/spectacular migrate`), or when doctor's `check_kits` validates chain integrity.

## Core principle

**A migration is one file that says "what changes between two workspace_schema versions and how to apply it."**

It bundles:
1. Frontmatter declaring the version edge (`from:` → `to:`), the kind (mechanical vs judgment), reversibility, and which paths it touches
2. A markdown body documenting: detection rule, steps, rollback (if reversible), validation
3. A reference to a bash function in `cli/spectacular` named `migration_apply_<id>` that implements mechanical steps. Judgment migrations have no apply-fn — they're handled by the skill walk.

The registry pattern means: adding a new migration = adding one .md file + one bash function (for mechanical migrations) OR a skill-walk section (for judgment). No special-casing in `cmd_migrate`.

## File layout

```
skills/spectacular/references/migrations/
├── v04-to-v05.md       # rename current/ → specs/, preserve layout
├── v05-to-v06.md       # ensure specs/ exists as always-set
└── <next>.md           # future migrations append here
```

Naming: `v<from>-to-v<to>.md`. From + to are dotted semver-lite (`0.4`, `0.5`, `1.0`). The file is the migration's source of truth; its filename is the human-readable ID.

## Frontmatter contract

```yaml
---
id: v04-to-v05                       # required, matches filename without .md
from: "0.4"                          # required, source workspace_schema
to: "0.5"                            # required, target workspace_schema
description: |                       # required, one-paragraph rationale
  Rename .spectacular/current/ → .spectacular/specs/. Preserves contents
  verbatim (flat .md files or per-capability subfolders). Idempotent.
mechanical: true                     # required, true | false
                                     #   true  → bash function migration_apply_<id> runs it
                                     #   false → judgment migration; skill walks it
reversible: false                    # required, true | false
                                     #   true  → reverse-fn must exist (migration_revert_<id>)
                                     #   false → CLI refuses downgrade; documented in body
apply-fn: migration_apply_v04_to_v05 # required when mechanical: true
                                     #   omitted when mechanical: false
affects:                             # required, list of paths this migration touches
  - .spectacular/current
  - .spectacular/specs
  - .spectacular/config.yaml         # workspace_schema bump
---
```

### Field semantics

**`id`** — kebab-case identifier. Matches the filename. Used in log output and skill verb routing. Convention: `v<from>-to-v<to>` (dashes around "to", dots in versions stripped or kept consistent — `v04-to-v05` not `v0.4-to-v0.5` to avoid dots in IDs).

**`from`** / **`to`** — `workspace_schema` versions. Quoted strings, dotted (e.g. `"0.4"`). The chain is built by matching every migration's `from` to a previous migration's `to` (or to the workspace's current schema as the starting `from`).

**`description`** — one paragraph. Surfaced by `spectacular migrate --dry-run` and `spectacular migrate --list`. Lead with what changes; explain why if non-obvious.

**`mechanical`** — boolean. `true` means the CLI runs the migration end-to-end via the named `apply-fn`. `false` means it requires judgment (e.g. content-aware rewrites, content moves with conflict resolution); CLI exits non-zero with a "run `/spectacular migrate` in your AI agent" message; skill walks it via `references/migrate.md`.

**`reversible`** — boolean. `true` means a `migration_revert_<id>` bash function exists for downgrade. `false` means the migration is one-way; CLI refuses downgrade.

**`apply-fn`** — name of the bash function in `cli/spectacular` that runs the mechanical steps. Required iff `mechanical: true`. Function signature: `apply-fn <dry_run>` where `<dry_run>` is `"true"` or `"false"`. Function prints progress to stdout, returns 0 on success.

**`affects`** — list of paths this migration touches. Used by doctor (to know what to re-validate after a migration runs) + by the skill walk (to know which canonical docs need snapshotting before edit). Convention: list directories (`.spectacular/current`) for renames; list files (`.spectacular/config.yaml`) for content edits.

## Body sections

Migration markdown body should contain at minimum:

```markdown
## Detection

How the CLI determines whether this migration is needed.
Usually a path check or frontmatter check. Stated as a predicate.

Example: ".spectacular/current/ exists AND .spectacular/specs/ does not"

## Steps

Numbered, idempotent.

1. <action>
2. <action>
3. Bump workspace_schema to "<to>"

For mechanical migrations, these mirror what `apply-fn` does. Documentation, not enforcement.

## Rollback

If `reversible: true`: numbered steps to reverse.
If `reversible: false`: explain why (e.g. "destructive content move; no reverse").

## Validation

How to verify the migration succeeded. Used by VERIFY scenarios.
```

## Loader behavior

The CLI loader (in `cmd_migrate`):

1. Scans `${CLAUDE_SKILL_DIR}/references/migrations/*.md` (or local dev path when run from the source tree)
2. Parses each file's frontmatter using the existing awk-based extractor pattern
3. Builds an in-memory list: `[(id, from, to, mechanical, apply-fn, description), ...]`
4. Sorts by `from` semver order
5. For each migration where `from >= current_schema AND to <= target_schema`:
   - If `mechanical: true` → look up `apply-fn` as a bash function name; call it with `<dry_run>` arg
   - If `mechanical: false` → print message routing to `/spectacular migrate`; exit non-zero
6. After successful apply, bumps `workspace_schema` to migration's `to`

## Chain validation (doctor)

`check_kits` validates the registry on every doctor run:

| Check | Failure mode |
|---|---|
| Every `from:` chains to a previous `to:` (or to baseline `"0.4"`) | Gap in chain — error |
| No two migrations share the same `(from, to)` pair | Duplicate migration — error |
| For `mechanical: true`, the `apply-fn` must reference a function actually defined in `cli/spectacular` | Broken reference — error |
| For `reversible: true`, a corresponding `migration_revert_<id>` must exist | Broken reverse — error |
| `id:` matches filename | Inconsistent naming — warning |

## Adding a new migration

1. Pick the `(from, to)` versions. Usually `(CURRENT_SCHEMA, CURRENT_SCHEMA + 1)`.
2. Create `references/migrations/v<from>-to-v<to>.md` with the frontmatter contract above.
3. Decide mechanical vs judgment:
   - **Mechanical**: write `migration_apply_<id>` bash function in `cli/spectacular`. Document steps in body.
   - **Judgment**: add a section to `references/migrate.md` walking the user through. No `apply-fn`.
4. Bump `CURRENT_SCHEMA` in `cli/spectacular` to the new `to:` value.
5. Update CHANGELOG + SPEC.md capability bullet (if structural enough to warrant it).
6. Add smoke test under `tests/cli/migrations/<id>.test.sh` (or extend `migrate.test.sh`).
7. Doctor `kits` will pick up the new entry automatically and validate it.

## Reserved migration IDs

- `v00-*` — reserved for placeholder / sentinel entries; should never ship
- IDs containing `manual`, `noop`, `test`, `skip` — reserved keywords; rejected by the loader

## Related

- [[doctor]] — chain validation lives in `check_kits` (per design lock 2026-05-23)
- [[doctor-areas]] — kits area row documents the chain check
- [[migrate]] — skill walk for judgment migrations
- [[doc-index]] — `migrations-contract` registered alongside other contract docs
- [[packs-contract]] — sibling registry pattern (registry + scope precedence + lifecycle)
