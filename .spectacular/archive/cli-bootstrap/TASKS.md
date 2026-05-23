---
updated: 2026-05-23
related:
  - PLAN.md
  - ../../archive/canonical-docs-rework/PLAN.md
---

# Tasks — cli-bootstrap

## Canonical doc set alignment (new — 2026-05-21)

Triggered by `canonical-docs-rework`. The CLI must scaffold the post-rework 7-doc root layer.

- [x] Update `PRD.md` stub — 6-slot shape + `related:` frontmatter pointing to siblings
- [x] Add `PRINCIPLES.md` stub — 8 default principles with `How the skill enforces this:` placeholders
- [x] Add `ARCHITECTURE.md` stub — section headers for Layout / Root layer / Frontmatter / Lifecycle / Versioning
- [x] Add `ROADMAP.md` stub — v1 / v2 / v3+ section headers
- [x] Update `AGENTS.md` stub — onboarding shape (What this folder is / How to operate / Context loading / Don'ts)
- [x] Success output already lists created files inline via `write_if_missing` — auto-extends to all 7 docs
- [x] `config.yaml` `agents.default_context:` left as PRD/STACK/DECISIONS for back-compat; full context map now in AGENTS.md per onboarding shape

## Repo layout

- [x] Create `cli/spectacular` (the binary, no extension, executable)
- [x] Create `cli/install.sh` (curl-installable, copies binary to `~/.local/bin/spectacular`)

## Flags + argument parsing

- [x] `--help` flag with usage output
- [x] `--name <slug>` — project name override
- [x] `--summary <text>` — project summary override
- [x] `--agents-file <filename>` — agents file override (default: `AGENTS.md`)
- [x] `--global` — install skill to `~/.agents/skills/` + `~/.claude/skills/`
- [x] `--update` — re-download latest skill, update `skills.lock`
- [x] `-i` / `--interactive` — prompt for all foundational settings

## Default init (bare, zero prompts)

- [x] Infer project name from current directory slug (lowercase, hyphens)
- [x] Create `.spectacular/current/` and `.spectacular/requests/`
- [x] Write `.spectacular/config.yaml` with name, empty summary, `agents.file` value
- [x] Create stub root files with frontmatter: `PRD.md`, `STACK.md`, `DECISIONS.md`
- [x] Create stub agents file (`AGENTS.md` or value of `--agents-file`)
- [x] Append `.spectacular.local/` to `.gitignore` (create if missing)

## Interactive mode (`-i`)

- [x] Prompt: project name (pre-filled with folder slug, enter to accept)
- [x] Prompt: project summary (optional, enter to skip)
- [x] Prompt: agents file preference (`AGENTS.md` / `CLAUDE.md`, default `AGENTS.md`)
- [x] Prompt: install scope (project-local / global, default project-local)

## Skill installation

- [x] Detect latest release tag from GitHub API (`/repos/alexsmedile/spectacular/releases/latest`)
- [x] Fall back to `main` tarball if no release tag exists
- [x] Download tarball, verify SHA, extract `skills/spectacular/` from it
- [x] Project-local: install extracted skill to `.agents/skills/spectacular/`
- [x] Project-local: create `.claude/skills/` and symlink `.claude/skills/spectacular/` → `.agents/skills/spectacular/`
- [x] Global (`--global`): install to `~/.agents/skills/spectacular/`, symlink `~/.claude/skills/spectacular/` → there
- [x] Write `.spectacular/skills.lock` with `ref`, `sha`, `installed`, `source`

## `--update` flow

- [x] Check `skills.lock` for current ref
- [x] Fetch latest release tag
- [x] If already at latest: report "already up to date (v1.0.0)", exit 0
- [x] If newer: download, verify, overwrite skill dir, update `skills.lock`

## Idempotency

- [x] Skip existing scaffold files, report each skipped file
- [x] Skip skill install if `skills.lock` exists and skill dir is present
- [x] Report installed ref from `skills.lock` when skipping
- [x] Safe to re-run as repair for partial scaffolds

## Error handling

- [x] Network failure during skill download: print actionable error, exit non-zero
- [x] Partial scaffold is kept on network failure (file writes already done)
- [x] SHA mismatch: abort install, print error, exit non-zero
- [x] Non-zero exit on any unrecoverable failure

## Output

- [x] Structured success summary (Created / Skill installed / Gitignore sections)
- [x] Skipped files listed inline when re-running
- [x] `--update` reports old ref → new ref

## Testing

- [x] Bare init on blank directory — correct scaffold, zero prompts (tested 2026-05-21)
- [x] `-i` mode — all prompts fire, flags pre-fill defaults (tested 2026-05-21 — name "my-test-app" captured into config.yaml + PRD heading)
- [x] `--global` — skill lands in `~/.agents/` and `~/.claude/` (tested 2026-05-23 with sandboxed `HOME=$(mktemp -d)` — both targets populated, symlink resolves correctly)
- [x] Re-run idempotency — all files skipped, correct report (tested 2026-05-21 — all 10 files reported skipped)
- [x] `--update` — fetches new version, updates `skills.lock` (tested 2026-05-23 — downgrade lock to v1.0.0 → `--update` detects v1.0.0 → v1.0.1, fetches tarball, rewrites lock; re-run reports "already up to date (v1.0.1)")
- [x] Network failure — exits non-zero, scaffold files intact
- [x] `install.sh` — syntax-checked OK; full curl-install test **blocked** until repo published to GitHub
- [x] Bash syntax check on `spectacular` binary — clean
- [x] Stub files match canonical-docs-rework shape — PRD 6-slot, PRINCIPLES 8 numbered, ARCHITECTURE 5-section, ROADMAP 3-tier, AGENTS onboarding shape (verified 2026-05-21)
