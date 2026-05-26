# Doctor — check areas

Per-area check tables. Load when explaining a specific finding, implementing a new check, or debugging area behavior. For overall orchestration see [[doctor]]; for the repair flow see [[doctor-repair]].

## `skill`

| Check | Severity | Detects |
|---|---|---|
| Install path exists | ❌ if missing | `.agents/skills/spectacular/` (or `~/.agents/...` if global) |
| Symlink valid | ❌ if broken | `.claude/skills/spectacular` resolves to install path |
| skills.lock present | ⚠️ if missing | `.spectacular/skills.lock` |
| skills.lock consistent | ⚠️ if drift | `ref:` matches installed version (best-effort) |
| SKILL.md parseable | ❌ if malformed | Top of install path; frontmatter loads |

## `workspace`

| Check | Severity | Detects |
|---|---|---|
| `.spectacular/` exists | ❌ if missing | Project's workspace directory |
| Always-set present | ❌ per missing file | `PRD.md`, `SPEC.md`, `config.yaml`, `<agents-file>`, `requests/`, `specs/` |
| `config.yaml` parses | ❌ if malformed | YAML well-formed (v1: `grep project:` only — see VERIFY.md carry-forward) |
| `<agents-file>` matches `config.yaml`'s `agents.file:` field | ⚠️ if mismatch | Frontmatter declaration vs reality |
| `.gitignore` includes `.spectacular.local/` | ⚠️ if missing | Mechanical fix available |

## `frontmatter`

| Check | Severity | Detects |
|---|---|---|
| Required fields per registry | ⚠️ per missing field | Each canonical doc has registry-declared required fields (`version`, `updated`, `summary`, etc.) |
| Frontmatter delimiter | ❌ if missing | First line of `*.md` canonical doc is `---` |
| Schema version current | ⚠️ if old | `version:` matches expected for the doc type (best-effort warning) |
| Date format | ⚠️ if non-ISO | `updated:` is `YYYY-MM-DD` |

## `snapshots`

| Check | Severity | Detects |
|---|---|---|
| Version sequence continuity | ⚠️ per gap | For each `<doc>@v<N>.<M>.md` family, flag missing intermediate versions |
| Snapshot frontmatter parses | ⚠️ if malformed | Same parse check as canonical docs |
| Orphan snapshots | ℹ️ | `<doc>@v<N>.md` exists but no `<doc>.md` next to it |

## `links`

| Check | Severity | Detects |
|---|---|---|
| `related:` targets exist | ⚠️ per missing | Each entry resolves to a file (relative paths from the file's location) |
| `[[wikilink]]` targets | ℹ️ | Wikilinks in body to known canonical refs (best-effort) |

## `lifecycle`

| Check | Severity | Detects |
|---|---|---|
| `active` requires SESSION.md | ⚠️ | Request with `status: active` should have a SESSION.md |
| `review` requires verification artifact | ❌ | One of: VERIFY.md, TASKS § Verification group, PLAN § Validation/Success/Acceptance |
| `verified` requires verification artifact | ❌ | Same as review; additionally VERIFY.md (if present) must have zero unchecked items |
| Stale `verified` | ℹ️ | `status: verified` more than 7 days old without archival — propose archive |
| Stale `active` | ℹ️ | `status: active` with no `updated:` change in 14+ days |

## `kits`

| Check | Severity | Detects |
|---|---|---|
| Bundled kits parse | ❌ per bad kit | Each `templates/prd/kits/<id>.md` in installed skill loads via the YAML parser |
| Project-local kits parse | ❌ per bad kit | Each `.spectacular/templates/prd/kits/<id>.md` (override) loads |
| `triggers-docs` references known doc IDs | ⚠️ per unknown | Each `always:` / `suggested:` entry is in the registry's `KNOWN_DOCS` |
| `extends:` is known | ❌ per unknown | Currently only `prd` is allowed |
| Kit declared in PRD frontmatter exists | ❌ if missing | If a PRD has `kit: <id>`, the kit file must be findable |

## `conventions` (v0.4.0+)

Active only when `.spectacular/config.yaml` declares `convention_pack:`. Otherwise reports `ℹ️ no convention_pack declared — skipped` and exits the area.

| Check | Severity | Detects |
|---|---|---|
| Pack source resolves | ❌ if missing | `convention_pack.source: <name>` must resolve via scope precedence (project-local → user → app-store → bundled) |
| Pack mode is valid | ⚠️ if unknown | `mode:` must be `suggest`, `scaffold`, or `enforce` (default `suggest`) |
| gitignore alignment (scaffold mode) | ⚠️ per missing entry | Each `pack.rules.gitignore.always-add` entry must appear in `.gitignore` |
| gitignore alignment (enforce mode) | ❌ per missing entry | Same check, severity escalated to error |

In `suggest` mode no drift checks run — the area only confirms the pack is active and reachable.

Mechanical fix: missing gitignore entries are appended by `spectacular doctor --fix` (extracts the entry from the finding message).

## `specs` (v0.5.0+)

| Check | Severity | Detects |
|---|---|---|
| `SPEC.md` present | ❌ if missing | `.spectacular/SPEC.md` exists |
| `specs/` dir present | ❌ if missing | `.spectacular/specs/` exists (may contain `.gitkeep` only) |
| Legacy `current/` migration | ⚠️ + mechanical | `current/` exists alongside `specs/` → propose `git mv current → specs` |
| Conflict (both `current/` and `specs/` with content) | ❌ | Manual resolution required |

## `docs` (v0.6.0+)

| Check | Severity | Detects |
|---|---|---|
| `docs.yaml` parses | ❌ if malformed | Top-level docs manifest loads |
| Declared pages exist on disk | ⚠️ per missing | Each `sections.*.pages[]` and `extras[]` entry resolves |
| Orphan pages | ⚠️ per orphan | Page file exists but not declared in docs.yaml |
| Required frontmatter (per page) | ❌ per missing | `title`, `description`, `section`, `status`, `updated` |
| Frontmatter delimiter | ❌ if missing | Mechanical fix injects stub |

## `feedback` (v1.6.0+)

Judgment-only area — no `--fix`. Scans `.spectacular/feedback/*.md` and `.spectacular/requests/*/feedback/*.md`.

| Check | Severity | Detects |
|---|---|---|
| Required frontmatter present | ⚠️ per entry | `type: feedback`, `status`, `opened`, `target` all set |
| Stale open entries | ⚠️ per entry | `status: open` with `opened` > 30 days ago |
| Orphan back-refs | ⚠️ per entry | `request:` field points to a missing request folder |

See [[feedback-rules]] and [[feedback-loop]] for the full mode spec.

## `ideas` (v1.7.0+)

Judgment-only area — no `--fix`. Scans `.spectacular/ideas/*.md`.

| Check | Severity | Detects |
|---|---|---|
| Required frontmatter present | ⚠️ per entry | `type: idea`, `status`, `updated` all set |
| Stale exploring | ⚠️ per entry | `status: exploring` with `updated` > 90 days ago |
| Orphan promoted | ⚠️ per entry | `status: promoted` but file still in `.spectacular/ideas/` (should be in `archive/ideas/`) |
| Unknown status | ⚠️ per entry | `status` not one of `parked\|exploring\|promoted` |

No `--fix` because every finding requires a human decision (promote? demote? delete? move?). See [[idea-rules]] for the full mode spec.

## Related

- [[doctor]] — entry point, severity model, report format
- [[doctor-repair]] — y/n/q judgment-fix flow
- [[doctor-substrate]] — when skill flows auto-invoke an area
