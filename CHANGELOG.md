# Changelog

All notable changes are documented here.
Format: [Keep a Changelog](https://keepachangelog.com/en/1.0.0/)

---

## [0.1.1] — 2026-05-11

### Fixed

- `cli/spectacular`: progress log in `download_and_install_skill` redirected to stderr — was polluting `skills.lock` with the "Fetching skill..." line
- `cli/spectacular`: skill ref resolution now tries releases API first, tags API second, `main` as final fallback — previously fell straight to `main` when no GitHub release existed

---

## [0.1.0] — 2026-05-11

### Added

- `cli/spectacular` — Bash CLI binary (`spectacular init`) with zero-prompt default, `-i` interactive mode, and flags: `--name`, `--summary`, `--agents-file`, `--global`, `--update`
- `cli/install.sh` — curl-installable installer; places binary at `~/.local/bin/spectacular`
- `skills/spectacular/` — `/spectacular` Claude Code slash command; lean SKILL.md orchestrator routing to `references/` subdocs
- `.spectacular/config.yaml` — `agents.file` key (override primary agents file) and `agents.tool_overrides` map (per-tool supplementary files, e.g. `claude: CLAUDE.md`)
- `skills.lock` — CLI-written lockfile tracking installed skill ref, SHA, and source URL
- Skill install targets both `.agents/skills/spectacular/` (source) and `.claude/skills/spectacular/` (symlink) for multi-tool compatibility
- `CLAUDE.md` — project guidance for Claude Code
