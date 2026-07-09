# Repository Guidelines

## Project Structure & Module Organization

Spectacular is a Bash CLI plus an agent skill. `cli/spectacular` is the executable used by `spectacular init`; `cli/install.sh` installs it into `~/.local/bin`. The skill lives in `skills/spectacular/`: keep `SKILL.md` as the lean orchestrator and place detailed workflows in `skills/spectacular/references/`. Plugin metadata is under `.claude-plugin/` and `.codex-plugin/`. Hook definitions live in `hooks/`, with the implementation in `scripts/hooks/pre-commit`. Documentation and visual assets are in `docs/`. This repository also uses its own `.spectacular/` workspace for planning and project state — read `.spectacular/specs/index.md` for a one-page index of what's built right now (introduced in v0.5.0; replaces the legacy `.spectacular/current/` folder).

## Build, Test, and Development Commands

There is no package manager build step. Use these checks during development:

```bash
bash -n cli/spectacular cli/install.sh scripts/hooks/pre-commit
scripts/hooks/pre-commit --check
./cli/spectacular --help
```

`bash -n` catches shell syntax errors. The pre-commit check verifies version consistency across plugin metadata, README badge, changelog, tags, and skill frontmatter. Run `./cli/spectacular --help` after CLI edits to confirm argument parsing still loads.

## Coding Style & Naming Conventions

Shell scripts use Bash with `set -euo pipefail`, two-space indentation inside functions and conditionals, lowercase helper function names, and uppercase constants such as `GITHUB_REPO`. Prefer small helpers over repeated inline blocks. Request slugs, file names, and generated workspace paths should use kebab-case where applicable. Markdown files should use concise headings, frontmatter when already present, and relative paths in examples.

## Testing Guidelines

A shell test suite lives under `tests/` — run it with `bash tests/run.sh` (or a single area, e.g. `bash tests/cli/doctor.test.sh`). Coverage spans `tests/cli/` (CLI flags + scaffold output, asserting on findings + exit codes), `tests/pipeline/` (debug-fleet orchestration runbooks with real bug fixtures), and `tests/agents/` (single-agent judgment fixtures for the fixer/investigator). Treat shell syntax checks (`bash -n cli/spectacular`) and the version guard as the required baseline on top of that.

For a one-off manual check of a `cli/spectacular` command path, use a temporary directory:

```bash
tmpdir="$(mktemp -d)" && cd "$tmpdir"
/path/to/repo/cli/spectacular init --name demo
```

Do not commit scratch directories. Prefer shell-focused tests that exercise CLI flags and scaffold output; add new CLI-behavior scenarios to `tests/cli/*.test.sh` in the existing assert style.

## Commit & Pull Request Guidelines

The git history uses Conventional Commit prefixes such as `feat:`, `fix:`, and `chore:`. Keep subjects imperative and specific, for example `fix: preserve skill install ref fallback`. Pull requests should describe the behavior change, list verification commands, link related issues or Spectacular request folders, and include screenshots only for documentation asset changes. Versioned releases must update all guarded version sources together.

## Agent-Specific Instructions

When changing canonical Spectacular docs or skill behavior, consult `.spectacular/` for current project intent before editing — start with `.spectacular/AGENTS.md` (operating rules), `.spectacular/PRD.md` (intent), and `.spectacular/specs/index.md` (what's built). Keep `.spectacular.local/` personal and uncommitted. Do not overwrite versioned skill snapshots in `skills/spectacular/versions/`; add a new snapshot when intentionally releasing a new skill version. When working on a request inside `.spectacular/requests/<slug>/`, load that folder's `PLAN.md` and `TASKS.md`, plus any `specs/<capability>.md` it references — not the whole `specs/` tree.
