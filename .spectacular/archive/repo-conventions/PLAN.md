---
status: superseded
priority: medium
owner: alex
updated: 2026-05-23
superseded_by:
  - ../convention-pack-schema/PLAN.md
  - ../convention-pack-fabricator/PLAN.md
  - ../convention-pack-application/PLAN.md
summary: "Bake the user's repo-structure conventions into Spectacular so init/new-request scaffolds the right shape by default — SUPERSEDED 2026-05-23. Original plan baked one opinionated shape directly into the skill. Replaced by an extensible convention-pack system (3 new requests). The 10 sections of conventions from this plan become the schema the fabricator grills against, plus an opinionated alex-default pack in the repo's app-store folder for opt-in download."
related:
  - ../../../skills/spectacular/SKILL.md
  - ../../STACK.md
references:
  - ~/code/NAMING_RULES.md
  - ~/code/README.md
  - ~/.claude/CLAUDE.md
---

# Plan — Repo Conventions

## Goal

Give Spectacular an **opinionated default** for host-repo structure (the surrounding repo, not just `.spectacular/`). On `spectacular init` and `spectacular new <slug>`, the skill should scaffold and place files using the user's established conventions instead of guessing.

## Why

Spectacular currently only opinionates `.spectacular/`. The host repo around it gets no help — folder layout, naming, README contract, and file placement decisions all fall back to LLM defaults, which drift from the user's actual style.

The user has documented conventions (`~/code/NAMING_RULES.md`, `~/code/README.md`, `~/.claude/CLAUDE.md`) and follows them in practice across `~/code/`, `~/vault/projects/`, `~/vault/data/skills_db/`. Spectacular should encode those conventions as the default scaffolding behavior.

## Scope

**In scope (v1)**
- `references/repo-layout.md` — single reference doc, loaded on demand, encoding all conventions below
- Update `references/init-workflow.md` — call repo-layout during init
- Update `references/new-request.md` — call repo-layout when scaffolding artifacts (research → `_research/`, screenshots → `requests/<slug>/artifacts/screenshots/`, etc.)
- `templates/repo/` — minimal scaffolds for the project types listed below
- README contract template (Type/Stack/Run header)
- `.gitignore` defaults template
- AGENTS.md + CLAUDE.md pairing convention

**Out of scope (v2)**
- A `spectacular scaffold <type>` retrofit command for existing repos
- Auto-detection of project type from existing files
- Convention linter (`spectacular check repo`)
- Per-language deeper conventions (Python packaging specifics, Node monorepo specifics)
- Migration command (`spectacular migrate` — rename `_archived/` → `_archive/`, etc.)

## Approach

Markdown-only. One reference doc + templates folder. No new CLI commands in v1. The skill reads `repo-layout.md` during init or new-request workflows and applies the relevant conventions.

### Conventions encoded (source: user's existing docs + observed practice)

#### 1. Naming

- `kebab-case` only at repo and folder level
- Format: `{anchor}-{descriptor}[-{role}][-{qualifier}]`
- 2-3 words max
- Role suffixes (fixed set): `ctrl`, `manager`, `svc`, `orch`, `worker`, `dash`, `viz`, `exp`
- Internal Python packages: `snake_case` (for imports)
- Code identifiers: `camelCase` (JS) / `snake_case` (Python) per language
- No `app-`/`svc-` prefixes (anchor first, suffix communicates role)
- No `v2`, `new`, `project` in names
- Dates only for `sandbox/` and `archive/` entries: `-YYYYMM` or `-YYYY-MM`

#### 2. Top-level taxonomy (for `~/code`-style mono-collections)

```
apps/       # deployable UIs & complete products
libs/       # reusable libraries/packages       (on demand)
tools/      # CLIs, scripts, automations, utilities
design/     # design-code hybrid projects
dash/       # dashboard UIs (admin, monitoring)
sandbox/    # experiments & spikes (dated)
templates/  # starters/boilerplates
archive/    # retired stuff (dated)
infra/      # infrastructure                    (on demand)
```

> Only relevant for *mono-collection* roots like `~/code`. Individual projects use the per-project layout below.

#### 3. Per-project standard folders

| Folder | Purpose | When to create |
|---|---|---|
| `src/` | Source code | Most code projects (Node, Python with src-layout) |
| `scripts/` | Utility scripts | Prefer over loose root scripts |
| `tests/` or `test/` | Tests | Per language convention (Python `tests/`, Node `test/`, Go co-located) |
| `docs/` | Human-facing docs | When docs exceed README |
| `examples/` | Runnable examples | For libraries, CLIs, plugins |
| `assets/` | Static media | Images, fonts, icons |
| `bin/` | Executables | CLIs only |
| `_research/` | Research artifacts | When NotebookLM/scrape output exists |
| `_archive/` | Archived content | Gitignored. Prefer over `_archived/` |
| `_backups/` | Timestamped backups | Gitignored. Prefer over `_backup/` |
| `_tmp/` or `scratch/` | Temporary work | Gitignored |
| `cli/` | CLI implementation | For tools that ship a CLI plus other layers (see spectacular itself) |

#### 4. Per-project root files

| File | Required? | Notes |
|---|---|---|
| `README.md` | Required | Human-facing. Must follow README contract below. |
| `AGENTS.md` | Required for agentic projects | Root-level agent guidance. Authoritative over README for agents. |
| `CLAUDE.md` | Optional | Often a symlink to AGENTS.md, or a Claude-scoped variant. |
| `LICENSE` | Required for OSS | |
| `CHANGELOG.md` | Required at v1+ | Keep-a-Changelog format |
| `.gitignore` | Required | See defaults below |
| `STACK.md` | Optional | Captured by `.spectacular/STACK.md` when spectacular is in use |
| `PRD.md` / `PLAN.md` / `TASKS.md` | Optional at root | Migrate to `.spectacular/` when spectacular is initialized |

#### 5. README contract (every project)

```markdown
# project-name
One-line description.

**Type**: cli | script | automation | app | lib | skill | plugin | content
**Stack**: e.g. Python 3.11 | Node 20 | Bash
**Run**: one-liner to execute (e.g. `python main.py --help`)

## What it does
2-3 sentences. Input → transformation → output.

## Setup
Steps to install dependencies and configure env vars.

## Usage
Minimal working example.
```

**Rule:** the synopsis block (Type/Stack/Run) must be readable in the first ~10 lines so an agent or human can triage instantly.

#### 6. AGENTS.md pattern

When present, `AGENTS.md` governs agent behavior. Each subfolder may also have its own `AGENTS.md` (most specific wins). Pattern observed across `~/code`, `~/vault`, `~/vault/data/skills_db/`.

Spectacular respects this: if a project has root `AGENTS.md`, spectacular skill checks for it before scaffolding, and references it from `.spectacular/AGENTS.md` rather than duplicating.

#### 7. .gitignore defaults

Every new project's `.gitignore` includes:

```gitignore
# Spectacular
.spectacular.local/

# Archives & backups (kept locally, never committed)
_archive/
_archived/
_backup/
_backups/

# Temp work
_tmp/
scratch/

# Tool-generated (only if user opts in — ASK BEFORE ADDING)
# .scrapekit/
# .playwright-mcp/
# .smart-env/

# Sensitive
.env.local
.env.*.local

# Language-specific (added based on detected stack)
```

**Rule:** Tool-generated hidden dirs (`.scrapekit/`, `.playwright-mcp/`, `.smart-env/`) are **never auto-gitignored** — the skill asks the user first (per global CLAUDE.md).

#### 8. File placement rules (where new files go)

When the skill creates any file, it routes per these rules:

| File type | Default location |
|---|---|
| Helper script | `scripts/<name>.sh` (never at root unless single-file project) |
| Architecture/contributor doc | `docs/` |
| Skill reference doc | `references/` inside the skill folder |
| Research artifact (NotebookLM, scrape) | `_research/<topic>/` |
| Backup before edit | `_backups/<timestamp>/` |
| Generated/cached | `.cache/` (gitignored) |
| Sensitive data | `.env.local`, `.spectacular.local/` |
| Large files (>5MB) | Flag to user, don't commit silently |
| Temp work | `scratch/` or `_tmp/` (gitignored) |
| Per-request artifacts | `.spectacular/requests/<slug>/artifacts/{screenshots,benchmarks,user-feedback,research}/` |

#### 9. Project-type templates

`templates/repo/<type>/` contains a minimal scaffold per type:

| Type | Adds |
|---|---|
| `cli` | `cli/`, `scripts/`, `install.sh`, `README.md`, `LICENSE` |
| `library` | `src/`, `tests/`, `examples/`, `docs/`, `README.md`, `LICENSE` |
| `webapp` | `src/`, `public/`, `tests/`, `.env.example`, `README.md` |
| `skill` | `SKILL.md`, `references/`, `templates/`, `scripts/`, `README.md` |
| `plugin` | `.claude-plugin/plugin.json`, `skills/`, `agents/`, `commands/`, `README.md` |
| `content` | `articles/`, `_research/`, `assets/`, `drafts/`, `README.md` |
| `research` | `_research/`, `notebooks/`, `data/`, `reports/`, `README.md` |
| `vault-project` | Obsidian-style: `assets/`, `inbox/`, `brand/`, `business/`, `marketing/`, `offer/`, `projects/`, `tasks/`, `README.md` (mirrors `~/vault/projects/<name>/`) |

Each template ships only the scaffold + a README stub conforming to the contract.

#### 10. Mono-collection vs project root

Spectacular detects which one it's initializing into:

- **Mono-collection root** (`~/code/`, `~/vault/data/skills_db/`) → don't scaffold `apps/tools/libs/...`. Just install `.spectacular/` at the mono-root and let each child project be initialized separately.
- **Project root** (a single repo) → scaffold per project-type templates above.

Heuristic: if the parent has subfolders that are themselves git repos or have `package.json`/`pyproject.toml`, treat as mono-collection.

### Architecture

```
skills/spectacular/
├── references/
│   ├── repo-layout.md              # NEW — encodes all conventions above
│   ├── init-workflow.md            # UPDATED — calls repo-layout
│   ├── new-request.md              # UPDATED — calls repo-layout for artifact placement
│   └── ...
└── templates/
    ├── prd/                        # existing
    └── repo/                       # NEW
        ├── cli/
        ├── library/
        ├── webapp/
        ├── skill/
        ├── plugin/
        ├── content/
        ├── research/
        └── vault-project/
```

## Success criteria

- A user runs `spectacular init` in a fresh repo, picks a project type, and gets a scaffold that **matches** what they'd build by hand from `~/code/NAMING_RULES.md`
- `spectacular new <slug>` places artifacts (research, screenshots, benchmarks) in the right subfolders without prompting
- README contract is enforced for any new project README the skill writes
- `.gitignore` defaults are applied without auto-adding tool-generated dirs (user is asked)
- Mono-collection vs project-root detection picks the right behavior on `~/code`, `~/vault/projects/<single>`, and a fresh `~/tmp/test-repo/`

## Build steps

1. [ ] Write `references/repo-layout.md` encoding all 10 sections above
2. [ ] Write minimal scaffolds in `templates/repo/<type>/` (8 types)
3. [ ] Write `.gitignore` default template
4. [ ] Write README contract template
5. [ ] Update `references/init-workflow.md` to call repo-layout
6. [ ] Update `references/new-request.md` to call repo-layout for artifact placement
7. [ ] Add `templates/repo/` to SKILL.md templates index
8. [ ] Dogfood: run init on a fresh test repo, verify scaffold matches conventions
9. [ ] Dogfood: run `spectacular new test-slug` and verify artifact paths

## Open questions

- Should repo conventions live in `references/repo-layout.md` OR be split into `naming.md` + `layout.md` + `gitignore.md`? (Leaning single doc for v1, split if it gets long)
- Should `.spectacular/STACK.md` extend to capture repo-level conventions per project, or stay focused on tech stack? (Leaning: STACK = tech, add `CONVENTIONS.md` if needed in v2)
- Detection of project type — ask user, or infer? (v1: ask. v2: infer from `package.json`/`pyproject.toml`/file presence.)
- Should there be a `spectacular doctor` / `spectacular check` to audit existing repos against conventions? (v2)
- Mono-collection root behavior: scaffold or skip? (Decided: skip — let child projects initialize separately)
