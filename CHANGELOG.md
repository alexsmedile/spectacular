# Changelog

All notable changes are documented here.
Format: [Keep a Changelog](https://keepachangelog.com/en/1.0.0/)

---

## [0.2.0] — 2026-05-21

### Changed

- **Canonical docs split.** The original 896-line `.spectacular/PRD.md` was split into four focused root docs:
  - `PRD.md` — product intent (now 121 lines, 6-slot shape: problem / who / success / non-goals / constraints / milestone)
  - `PRINCIPLES.md` — 8 operating principles, each with a runtime enforcement hook
  - `ARCHITECTURE.md` — workspace structure, frontmatter conventions, lifecycle, versioning
  - `ROADMAP.md` — versioned future work (v1 / v2 / v3+)
- **AGENTS.md rewritten** as the in-folder onboarding doc for any agent landing in `.spectacular/`. Authoritative source for per-task context loading rules.
- **PLAN.md template upgraded** to the 7-slot decomposition: goal / why / constraints / milestones / tasks / dependencies / validation / deliverables.
- **CLI scaffolds the full 7-doc root layer** on every `spectacular init`. PRD stub uses the new 6-slot shape; AGENTS stub uses the onboarding shape; new PRINCIPLES / ARCHITECTURE / ROADMAP stubs included.
- **Skill references aligned** with the new doc set — `status.md`, `onboarding.md`, `init-workflow.md`, `scaffold-reference.md`, `new-request.md`, `versioning.md`, and SKILL.md state-awareness all updated.
- **Project docs aligned** — README, CLAUDE.md, `docs/scaffold.md`, `docs/configuration.md`, `docs/commands.md`, `docs/troubleshooting.md`, `docs/workflow.md` all reflect the new 7-doc canonical set.

### Added

- `prd / prd refine / prd review` skill triggers — interactive PRD building with 5 kits (coding / product / content / research / blank), vibe→spec refine patterns, and a pass/fail quality gate.
- `requests/prd-craft/` and `requests/canonical-docs-rework/` — tracking artifacts for the v0.2.0 work.
- Snapshot history preserved: `PRD@v1.3.md`, `AGENTS@v1.0.md`, request PLAN/TASKS `@v1.0.md` and `@v1.1.md`.

### Anti-patterns formalized

- **Never create `requests/<slug>/PRD.md`.** Product intent is project-wide and lives at `.spectacular/PRD.md`. Per-request folders use `PLAN.md` + `TASKS.md` only.

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
