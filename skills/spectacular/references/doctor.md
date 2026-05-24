# Doctor — environment/infrastructure self-check

Loaded when the user runs `spectacular doctor [<area>] [--fix]` or when another skill flow auto-invokes a doctor subset.

## Core principle

**Doctor checks the substrate, not the content.**

Content quality (PRD slot completeness, vague-word usage, success-criteria measurability) is what `<doc> review` does. Doctor verifies that the workspace can *do its job* — skill installed, registry parseable, frontmatter valid, snapshots continuous, links intact, lifecycle hygienic.

## CLI vs skill split

| Phase | Runs as | Output |
|---|---|---|
| Detect | CLI (`spectacular doctor`) | Structured findings: area, severity, file, message, proposed fix, fix type |
| Mechanical fix | CLI (`spectacular doctor --fix`) | Content-free repairs only |
| Judgment fix | Skill (`/spectacular doctor --fix`) | Per-finding context-aware proposals, `y/n/q` confirm per fix, snapshot-before-edit → see [[doctor-repair]] |
| Skill-invoked subset | Skill (auto) | When other flows hit substrate failures → see [[doctor-substrate]] |

The CLI never proposes fixes that require judgment. The skill never edits without explicit per-finding user confirmation.

## Severity model

| Symbol | Severity | Exit code contribution |
|---|---|---|
| ✅ | pass | none |
| ⚠️ | warning | sets exit to 1 (if no errors) |
| ❌ | error | sets exit to 2 |
| ℹ️ | info | none — informational only |

Exit codes: 0 (clean), 1 (warnings only), 2 (errors present).

## Check areas (v1)

Ten areas in current order: `skill`, `workspace`, `frontmatter`, `snapshots`, `links`, `lifecycle`, `kits`, `conventions`, `specs`, `docs`.

Per-area check tables, severity model per check, and mechanical-vs-judgment classification live in [[doctor-areas]]. Load that file when you need to explain what a finding means or implement a new check.

## Repair flows

| Need | Load |
|---|---|
| Walking judgment findings interactively (`/spectacular doctor --fix`) | [[doctor-repair]] |
| A skill flow hit a substrate failure and you need to auto-invoke a scoped check | [[doctor-substrate]] |
| Explain a specific finding or area check | [[doctor-areas]] |

## Report format

### Text (default)

```
Spectacular doctor

skill         ✅ spectacular@0.6.0 installed at .agents/skills/spectacular/
workspace     ✅ always-set present
frontmatter   ⚠️ .spectacular/STACK.md is missing required field: version
              → propose: add `version: 1.0` (requires agent — run /spectacular doctor --fix)
links         ⚠️ requests/foo/PLAN.md → "../bar/PLAN.md" target missing
              → propose: remove the stale link (requires agent)
kits          ✅ 5 bundled kits parse; no project-local overrides

0 errors, 2 warnings, 0 info

Run `spectacular doctor --fix` for mechanical repairs.
Run `/spectacular doctor --fix` for agent-driven repairs.
```

### JSON (`--format json`)

Shape: `{ version, ran_at, workspace, summary: { errors, warnings, info }, findings: [...] }`.

Each finding: `{ area, severity, file, message, proposed_fix, fix_type }`.

`fix_type` is the contract between CLI and skill:
- `mechanical` — CLI's `--fix` can apply without user prompt
- `judgment` — CLI reports; agent handles via `/spectacular doctor --fix` (see [[doctor-repair]])

Pass-entries are omitted from JSON for signal-to-noise.

## Mechanical fixes (CLI's `--fix`)

The CLI directly applies content-free repairs:

| Fix | When | Action |
|---|---|---|
| `.gitignore` append | Entry missing | Append `.spectacular.local/` to existing or create the file |
| Missing dir | `.spectacular/requests/` or `.spectacular/specs/` absent | `mkdir -p` |
| Dangling symlink | `.claude/skills/spectacular` broken but `.agents/...` exists | Recreate symlink |
| Missing always-set stub | `PRD.md` / `config.yaml` / `<agents-file>` entirely absent | Re-scaffold via init templates (only if file is *missing*) |
| Legacy `current/` → `specs/` | v0.4.x layout detected | `git mv` migration |
| docs page frontmatter | Page exists but lacks frontmatter | Inject minimal stub |
| Pack-driven gitignore drift | `enforce` mode + missing entry | Append to `.gitignore` |

Stdout per fix: `✓ fixed <area>: <action>`. After all mechanical fixes apply, the CLI re-runs detect; remaining findings need agent attention.

## What doctor does NOT do

- Run on every invocation (explicit trigger only — exception: skill-invoked subsets per [[doctor-substrate]])
- Modify content during detect (read-only)
- Apply judgment fixes from the CLI (route to the skill per [[doctor-repair]])
- Bulk-apply fixes without per-finding confirmation (`--yes-to-all` rejected for v1)
- Rewrite git history
- Check content quality — that's `<doc> review`
- Validate the host project's code/tests/build

## Related

- [[doctor-areas]] — per-area check tables (load when explaining/implementing a check)
- [[doctor-repair]] — judgment-fix repair flow (load on `--fix` invocations)
- [[doctor-substrate]] — auto-invocation spec (load from status/grill/onboarding/lifecycle)
- [[doc-index]] — source of truth for "what should each doc look like"
- [[kits-contract]] — schema doctor validates kit files against
- [[verification]] — convention doctor's lifecycle check enforces
- [[versioning]] — snapshot-before-edit rule every doctor judgment fix follows
- [[init-workflow]] — smart-init's diagnostic-deferral messages point here
