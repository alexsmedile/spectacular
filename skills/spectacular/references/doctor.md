# Doctor — environment/infrastructure self-check

Loaded when the user runs `spectacular doctor [<area>] [--fix]` or when another skill flow auto-invokes a doctor subset (see § Skill-invoked checks).

## Core principle

**Doctor checks the substrate, not the content.**

Content quality (PRD slot completeness, vague-word usage, success-criteria measurability) is what `<doc> review` does. Doctor verifies that the workspace can *do its job* — skill installed, registry parseable, frontmatter valid, snapshots continuous, links intact, lifecycle hygienic.

## CLI vs skill split

| Phase | Runs as | Output |
|---|---|---|
| Detect | CLI (`spectacular doctor`) | Structured findings: area, severity, file, message, proposed fix, fix type |
| Mechanical fix | CLI (`spectacular doctor --fix`) | Content-free repairs only: gitignore append, missing dirs, dangling symlinks, missing always-set stubs |
| Judgment fix | Skill (`/spectacular doctor --fix`) | Per-finding context-aware proposals, `y/n/q` confirm per fix, snapshot-before-edit |
| Skill-invoked subset | Skill (auto) | When other flows hit substrate failures, run the relevant doctor area and surface inline |

The CLI never proposes fixes that require judgment. The skill never edits without explicit per-finding user confirmation.

## Severity model

Mirroring flutter doctor + claude doctor conventions:

| Symbol | Severity | Exit code contribution |
|---|---|---|
| ✅ | pass | none |
| ⚠️ | warning | sets exit to 1 (if no errors) |
| ❌ | error | sets exit to 2 |
| ℹ️ | info | none — informational only |

Exit codes: 0 (clean), 1 (warnings only), 2 (errors present).

## Check areas (v1)

### `skill`

| Check | Severity | Detects |
|---|---|---|
| Install path exists | ❌ if missing | `.agents/skills/spectacular/` (or `~/.agents/...` if global) |
| Symlink valid | ❌ if broken | `.claude/skills/spectacular` resolves to install path |
| skills.lock present | ⚠️ if missing | `.spectacular/skills.lock` |
| skills.lock consistent | ⚠️ if drift | `ref:` matches installed version (best-effort; only checks if lock readable) |
| SKILL.md parseable | ❌ if malformed | Top of install path; frontmatter loads |

### `workspace`

| Check | Severity | Detects |
|---|---|---|
| `.spectacular/` exists | ❌ if missing | Project's workspace directory |
| Always-set present | ❌ per missing file | `PRD.md`, `config.yaml`, `<agents-file>`, `requests/`, `current/` |
| `config.yaml` parses | ❌ if malformed | YAML well-formed |
| `<agents-file>` matches `config.yaml`'s `agents.file:` field | ⚠️ if mismatch | Frontmatter declaration vs reality |

### `frontmatter`

| Check | Severity | Detects |
|---|---|---|
| Required fields per registry | ⚠️ per missing field | Each canonical doc has its registry-declared required fields (`version`, `updated`, `summary`, etc.) |
| Frontmatter delimiter | ❌ if missing | First line of `*.md` canonical doc is `---` |
| Schema version current | ⚠️ if old | `version:` matches expected for the doc type (best-effort warning, not blocker) |
| Date format | ⚠️ if non-ISO | `updated:` is `YYYY-MM-DD` |

### `snapshots`

| Check | Severity | Detects |
|---|---|---|
| Version sequence continuity | ⚠️ per gap | For each `<doc>@v<N>.<M>.md` family, flag missing intermediate versions (e.g. v1.0 + v1.2 with no v1.1) |
| Snapshot frontmatter parses | ⚠️ if malformed | Same parse check as canonical docs |
| Orphan snapshots | ℹ️ | `<doc>@v<N>.md` exists but no `<doc>.md` next to it |

### `links`

| Check | Severity | Detects |
|---|---|---|
| `related:` targets exist | ⚠️ per missing | Each entry in `related:` frontmatter resolves to a file (relative paths from the file's location) |
| `[[wikilink]]` targets | ℹ️ | Wikilinks in body to known canonical refs (best-effort; informational only) |

### `lifecycle`

| Check | Severity | Detects |
|---|---|---|
| `active` requires SESSION.md | ⚠️ | Request with `status: active` should have a SESSION.md |
| `review` requires verification artifact | ❌ | Request with `status: review` must have one of: VERIFY.md (any items), TASKS § Verification group, PLAN § Validation entries |
| `verified` requires verification artifact | ❌ | Same as review, but additionally — VERIFY.md (if present) must have zero unchecked items |
| Stale `verified` | ℹ️ | `status: verified` more than 7 days old without archival — propose archive |
| Stale `active` | ℹ️ | `status: active` with no `updated:` change in 14+ days |

### `kits`

| Check | Severity | Detects |
|---|---|---|
| Bundled kits parse | ❌ per bad kit | Each `templates/prd/kits/<id>.md` in installed skill loads via the YAML parser |
| Project-local kits parse | ❌ per bad kit | Each `.spectacular/templates/prd/kits/<id>.md` (override) loads |
| `triggers-docs` references known doc IDs | ⚠️ per unknown | Each `always:` / `suggested:` entry is in the registry's `KNOWN_DOCS` |
| `extends:` is known | ❌ per unknown | Currently only `prd` is allowed |
| Kit declared in PRD frontmatter exists | ❌ if missing | If a PRD has `kit: <id>`, the kit file must be findable |

## Report format

### Text (default)

```
Spectacular doctor

skill         ✅ spectacular@0.3.0 installed at .agents/skills/spectacular/
workspace     ✅ always-set present
frontmatter   ⚠️ .spectacular/STACK.md is missing required field: version
              → propose: add `version: 1.0` (requires agent — run /spectacular doctor --fix)
snapshots     ❌ PRD@v1.1.md missing (v1.0, v1.2, v1.3 present)
              → propose: backfill from git log OR document the skip (requires agent)
links         ⚠️ requests/foo/PLAN.md → "../bar/PLAN.md" target missing
              → propose: remove the stale link (requires agent)
lifecycle     ✅ all requests have verification artifacts
kits          ✅ 5 bundled kits parse; no project-local overrides

1 error, 2 warnings, 0 info

Run `spectacular doctor --fix` for mechanical repairs.
Run `/spectacular doctor --fix` for agent-driven repairs.
```

### JSON (`--format json`)

```json
{
  "version": "0.3.0",
  "ran_at": "2026-05-22T21:00:00Z",
  "workspace": "/Users/alex/projects/my-app",
  "summary": { "errors": 1, "warnings": 2, "info": 0 },
  "findings": [
    {
      "area": "frontmatter",
      "severity": "warning",
      "file": ".spectacular/STACK.md",
      "message": "missing required field: version",
      "proposed_fix": "add `version: 1.0` to frontmatter",
      "fix_type": "judgment"
    },
    {
      "area": "snapshots",
      "severity": "error",
      "file": ".spectacular/PRD@v1.1.md",
      "message": "version-sequence gap (v1.0 and v1.2 present)",
      "proposed_fix": "backfill from git log OR document the skip in DECISIONS.md",
      "fix_type": "judgment"
    },
    {
      "area": "links",
      "severity": "warning",
      "file": ".spectacular/requests/foo/PLAN.md",
      "message": "related: target '../bar/PLAN.md' does not exist",
      "proposed_fix": "remove the stale link",
      "fix_type": "judgment"
    }
  ]
}
```

`fix_type` is the contract between CLI and skill:
- `mechanical` — CLI's `--fix` can apply without user prompt (only on explicit `--fix`)
- `judgment` — CLI reports; agent handles via `/spectacular doctor --fix`

## Mechanical fixes (CLI's `--fix`)

The CLI can directly apply these without human judgment because they're content-free:

| Fix | When | Action |
|---|---|---|
| `.gitignore` append | Entry missing | Append `.spectacular.local/` to existing `.gitignore` (or create the file) |
| Missing dir | `.spectacular/requests/` or `.spectacular/current/` absent | `mkdir -p` |
| Dangling symlink | `.claude/skills/spectacular` broken but `.agents/skills/spectacular/` exists | Recreate symlink to existing target |
| Missing always-set stub | `PRD.md` / `config.yaml` / `<agents-file>` entirely absent | Re-scaffold using init's templates (only if file is *missing*, never overwrites existing content) |

Stdout per fix: `✓ fixed <area>: <action>`.

After all mechanical fixes apply, the CLI re-runs detect; remaining findings are reported as needing agent attention.

## Judgment fixes (skill's `/spectacular doctor --fix`)

The skill walks each `fix_type: judgment` finding and proposes context-aware actions:

```
> .spectacular/STACK.md missing `version:` field

  The registry expects `version: 1.0` for canonical root docs. Existing
  frontmatter has `updated:` and `summary:` but no `version:`. Add it?

  [y] add `version: 1.0` as the first frontmatter field
  [n] skip this finding
  [q] quit doctor

> y

  Snapshotting .spectacular/STACK.md to STACK@v0.9.md ... ✓
  Editing .spectacular/STACK.md ... ✓
  → fixed: added version field
```

**Per-finding behavior:**
- `y` — apply (snapshot first if file is canonical, per [[versioning]])
- `n` — skip, move to next
- `q` — print remaining findings summary, exit

**Snapshot-before-edit is mandatory** for every fix touching a canonical doc. The skill calls `spectacular snapshot <file>` semantics (or the equivalent inline operation).

## Skill-invoked checks

Doctor isn't only a CLI subcommand — the skill itself can invoke scoped subsets when its operations hit substrate problems. The skill **does not run a full doctor sweep**; it runs the relevant area and surfaces results inline.

| Invoking flow | Doctor subset | When |
|---|---|---|
| `references/status.md` briefing | `kits` + `frontmatter` (root docs only) | If `doc-registry.md` parse fails or root docs have malformed frontmatter |
| `references/grill.md` pre-flight | `kits` | If the active kit file fails to parse |
| `references/onboarding.md` first-invocation | `workspace` + `frontmatter` | First time the skill sees a workspace |
| `references/lifecycle.md` transition | `lifecycle` (scoped to that request) | When proposing `verified` — check the verification artifact exists per [[verification]] |

These invocations show the user what's broken without forcing them to learn about doctor first. Each surfaces "I noticed this issue with the workspace — run `spectacular doctor <area>` for details or `/spectacular doctor --fix` to repair."

## What doctor does NOT do

- Run on every invocation (explicit trigger only)
- Modify content during detect (read-only)
- Apply judgment fixes from the CLI (those route to the skill)
- Bulk-apply fixes without per-finding confirmation (`--yes-to-all` is rejected for v1)
- Rewrite git history (snapshot backfill proposes new commits only)
- Check content quality — that's `<doc> review`
- Replace `<doc> review` — the two have different scopes and both must exist
- Validate the host project's code/tests/build (out of scope — Spectacular doesn't run the host)

## Related

- [[doc-registry]] — source of truth for "what should each doc look like"
- [[kits-contract]] — schema doctor validates kit files against
- [[verification]] — convention doctor's lifecycle check enforces
- [[versioning]] — snapshot-before-edit rule every doctor judgment fix follows
- [[init-workflow]] — smart-init's diagnostic-deferral messages point here
