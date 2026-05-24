<div align="center">

<img src="docs/assets/banner.svg" width="100%" alt="Spectacular — Other tools help you write code faster. Spectacular helps you stay coherent while you do." />

![License](https://img.shields.io/badge/license-MIT-blue)
![Claude Code](https://img.shields.io/badge/Claude%20Code-compatible-blueviolet)
![Platform](https://img.shields.io/badge/platform-Claude%20%7C%20Codex%20%7C%20Cursor-lightgrey)
![Version](https://img.shields.io/badge/version-1.3.0-green)

</div>

---

## The problem

AI coding agents are fast. They are not coherent.

Every new session starts cold. Decisions made last week are invisible. Context that took an hour to build evaporates the moment the conversation ends. Your agent doesn't know what was decided, what was tried, or what's next — so you re-explain, re-decide, and re-discover the same things, over and over.

The bottleneck isn't code generation. It's everything around it.

---

## What it is

Spectacular is an operational workspace for AI-assisted software projects. Drop a `.spectacular/` directory in any repo and it becomes the shared brain for you and every agent that touches the project.

Strategic context is split across seven focused canonical docs (PRD / PRINCIPLES / ARCHITECTURE / ROADMAP / STACK / DECISIONS / AGENTS) so agents load only what each task needs, not the entire repo.

It ships as three layers:

- **Convention** — a structured directory contract separating strategy, current truth, and active work
- **Skill** — a `/spectacular` Claude Code slash command that reads state, proposes actions, scaffolds files, and manages lifecycle transitions
- **CLI** — `spectacular init` bootstraps the workspace in any project in seconds

The skill is the primary interface. The CLI runs once.

---

## Quick start

```bash
# Install the CLI
curl -fsSL https://raw.githubusercontent.com/alexsmedile/spectacular/main/cli/install.sh | bash

# Bootstrap any project
cd your-project
spectacular init

# Open Claude Code or Codex and run
/spectacular
```

> [!TIP]
> `spectacular init` scaffolds the `.spectacular/` directory and installs the `/spectacular` skill into `.agents/skills/spectacular/` (source) and `.claude/skills/spectacular/` (symlink). After init, `/spectacular` is immediately available to Codex and Claude Code.

> [!NOTE]
> `spectacular init` is zero-prompt by default. It infers the project name from the folder slug and uses `AGENTS.md` as the primary agent context file. Pass `-i` for interactive setup, or use flags: `--name`, `--summary`, `--agents-file`, `--global`.

---

## How it works

### The three layers

<img src="docs/assets/layer-map.svg" width="100%" alt="Strategy, truth, and work mapped to the .spectacular/ directory" />

**Strategy** — changes rarely. Split across focused canonical docs:
- `PRD.md` — *what* the product is and *why*
- `PRINCIPLES.md` — operating principles + runtime enforcement hooks
- `ARCHITECTURE.md` — the workspace structure itself
- `ROADMAP.md` — versioned future work
- `STACK.md` — host project's tech choices
- `DECISIONS.md` — ADR-style decision log
- `AGENTS.md` — onboarding doc for any agent landing in `.spectacular/`

Agents load only what the current task needs (planning loads PRD + PRINCIPLES + DECISIONS; implementation loads STACK + PLAN + TASKS).

**Current truth** — reflects actual system behavior right now. Modular capability specs (auth, billing, editor). Never overwritten in place — the skill snapshots before proposing edits.

**Active work** — temporary. Each request gets a folder with `PLAN.md` (intent + 7-slot decomposition: goal, constraints, milestones, tasks, dependencies, validation, deliverables) and `TASKS.md` (execution checklist). Archived on completion, not deleted.

---

### The lifecycle

<img src="docs/assets/lifecycle-flow.svg" width="100%" alt="Request lifecycle: planned → active → review → verified → archived" />

State lives in `PLAN.md` frontmatter. The skill reads it on every invocation and surfaces the highest-priority next action. When all tasks are checked, it proposes moving to `review`. When the checklist passes, it proposes `archived` — and offers to update `current/` with anything that changed.

---

### The workspace

```
.spectacular/
│   ── always-set (created by every init) ─────────────────────────────
├── PRD.md              # product intent — what & why & for whom
├── SPEC.md             # system spec — index of what's built right now (present tense)
├── config.yaml         # naming, kit identity, agent file overrides
├── AGENTS.md           # onboarding doc for agents working in this folder
├── requests/           # active and planned work
└── specs/              # per-capability specs (optional; SPEC.md is the index)

│   ── opt-in (scaffolded by kit declaration or --with flag) ──────────
├── PRINCIPLES.md       # operating principles + enforcement hooks
├── ARCHITECTURE.md     # .spectacular/ structure, frontmatter, lifecycle, versioning
├── ROADMAP.md          # time-ordered "what's next"
├── STACK.md            # host project's tech choices
├── DECISIONS.md        # ADR-style decision log

│   ── created on demand ─────────────────────────────────────────────
├── memory/             # long-term operational learning (git-committed)
└── archive/            # completed requests (never deleted)
```

A typical coding project (`spectacular init --kit coding`) scaffolds the always-set + `STACK.md` + `ARCHITECTURE.md`. A doc-only or research project (`spectacular init --kit research` or `--kit blank`) gets only the always-set. Smart-init never overwrites existing files — re-running is always safe.

`.spectacular.local/` — personal overrides, always gitignored.

---

## Skill commands

| Command | What happens |
|---|---|
| `/spectacular` | Project briefing — active requests, draft capabilities, next action |
| `spectacular status` | Same as no-arg invocation |
| `spectacular new <slug>` | Scaffold a new request (PLAN.md + TASKS.md) |
| `spectacular promote <slug>` | Advance lifecycle: `planned → active → review → verified` |
| `spectacular snapshot <file>` | Snapshot a canonical document before editing |
| `spectacular touch <file>` | Bump `updated:` on a canonical doc |
| `spectacular archive <slug>` | Archive a completed request; propose `SPEC.md`/`specs/` sync + memory entries |
| `spectacular remember this` | Write an insight to `memory/` immediately |

> [!WARNING]
> **`spectacular docs *` verbs are deprecated as of v1.2.0** — public-facing documentation work moved to the standalone [pageworks](https://github.com/alexsmedile/pageworks) skill. The verbs still work and emit a deprecation banner pointing to the equivalent pageworks command; they will be removed in v2.0.0. Install pageworks via its own one-liner.

---

## CLI reference

```
spectacular init                              # always-set + blank kit (5 files only)
spectacular init -i                           # interactive — kit menu + per-doc prompts
spectacular init --kit coding                 # always-set + coding kit's STACK + ARCHITECTURE
spectacular init --with principles,roadmap   # additive — those two on top of always-set
spectacular init --kit coding --minimal       # always-set only; kit identity preserved
spectacular init --name my-app
spectacular init --agents-file CLAUDE.md      # use CLAUDE.md instead of AGENTS.md
spectacular init --global                     # install skill to ~/.agents/ and ~/.claude/
spectacular init --update                     # re-download latest skill release

spectacular doctor                            # substrate self-check (all areas)
spectacular doctor <area>                     # scoped: skill | workspace | frontmatter | snapshots | links | lifecycle | kits | conventions | specs | docs
spectacular doctor --fix                      # apply mechanical fixes (gitignore, missing dirs, dangling symlinks, pack drift, legacy current/ migration)
spectacular doctor --format json              # JSON report for the skill or other tools

spectacular migrate                           # apply pending workspace schema migrations
spectacular migrate --dry-run                 # preview the migration plan
spectacular migrate --list                    # show all available migrations + current schema version

spectacular pack list                         # show installed packs (bundled + app-store + user + project-local)
spectacular pack install <name>               # install pack to ~/.spectacular/packs/<name>/
spectacular pack install <name> --from <path> # install from arbitrary local folder
spectacular pack show <name>                  # print scope + pack.md frontmatter
spectacular pack remove <name>                # remove user-scope pack (--force for bundled/app-store/project-local)
```

### Convention packs (v0.4.0)

Packs encode opt-in repo-shape opinions — naming rules, folder taxonomy, gitignore baseline, README contract, file-placement rules, project-type scaffolds. Two ship out of the box:

- **`minimal`** (bundled) — README contract + safe gitignore baseline. The default.
- **`alex-default`** (app-store) — fully-opinionated: kebab-case naming with role suffixes, mono-collection detection, 8 project-type scaffolds, language-specific gitignore blocks.

Activate a pack per-repo by adding to `.spectacular/config.yaml`:

```yaml
convention_pack:
  source: alex-default
  mode: scaffold     # suggest | scaffold | enforce
```

| Mode | Init behavior | Doctor behavior |
|---|---|---|
| `suggest` | Pack read, not applied | Reports pack active, no drift checks |
| `scaffold` | Appends pack gitignore entries | Warnings on drift (exit 1) |
| `enforce` | Same as scaffold | Errors on drift (exit 2); `--fix` repairs |

Full schema in [`skills/spectacular/references/packs-contract.md`](skills/spectacular/references/packs-contract.md). App-store packs live in [`packs/`](packs/).

> [!TIP]
> Init scaffolds the **6-file always-set** by default (`PRD.md`, `SPEC.md`, `config.yaml`, `<agents-file>`, `requests/`, `specs/`). Kits add docs they need. Use `--with` for explicit extras. Use `--minimal` to ignore the kit's defaults. (v0.4.x scaffolded `current/` instead of `SPEC.md` + `specs/` — see [CHANGELOG](CHANGELOG.md) for the migration.)

> [!TIP]
> Claude-only team? Use `--agents-file CLAUDE.md`. Multi-tool team? Keep `AGENTS.md` as primary and add `tool_overrides.claude: CLAUDE.md` to `config.yaml` — the skill will surface both.

---

## Pairing with pageworks

Spectacular owns `.spectacular/` — strategy, current truth, active work. **Public-facing documentation (the `docs/` surface) is owned by [pageworks](https://github.com/alexsmedile/pageworks)**, a sibling skill extracted in v1.2.0.

The two compose: when a SPEC-touching request archives, spectacular asks whether `docs/` should be updated and hands off to pageworks if you confirm. `spectacular doctor docs` reports discovery only (folder + manifest presence; install hint if pageworks missing) — never validates schema. There is no automatic invocation across the boundary.

Install pageworks separately when you need a public docs surface:

```bash
curl -fsSL https://raw.githubusercontent.com/alexsmedile/pageworks/main/cli/install.sh | bash
```

---

## Who it's for

- Solo developers using Claude Code, Codex, or Cursor on projects that span weeks or months
- Small teams where AI agents need to share operational context
- Anyone who has re-explained the same architectural decision to an agent more than twice

## Who it's not for

- Projects that live and die in a single session — the structure has no value at that scale
- Teams that already have a working context management system
- Non-technical users — this is a developer-facing directory convention, not a GUI tool

---

## Install

**CLI** — curl one-liner, installs `spectacular` to `~/.local/bin/`:

```bash
curl -fsSL https://raw.githubusercontent.com/alexsmedile/spectacular/main/cli/install.sh | bash
```

**Claude Code plugin** — from Claude Code:

```text
/plugin marketplace add alexsmedile/spectacular
/plugin install spectacular@spectacular
/reload-plugins
```

You can also use the Claude Code CLI:

```bash
claude plugin marketplace add alexsmedile/spectacular
claude plugin install spectacular@spectacular
```

**Codex plugin** — from Codex, open `/plugins`, add this marketplace, then install or enable `spectacular`:

```text
alexsmedile/spectacular
```

Codex CLI builds that expose marketplace commands can add the marketplace first:

```bash
codex plugin marketplace add alexsmedile/spectacular
```

To update an already-added marketplace:

```bash
codex plugin marketplace upgrade spectacular
```

Then open Codex and use `/plugins` to install or enable `spectacular`.

**Skill only** (no CLI, no plugin marketplace):

```bash
# manual
cp -r skills/spectacular ~/.claude/skills/
mkdir -p ~/.agents/skills
cp -r skills/spectacular ~/.agents/skills/
```

**Skill install locations:**

| Scope | Source | Claude symlink |
|---|---|---|
| Project-local (default) | `.agents/skills/spectacular/` | `.claude/skills/spectacular/` → above |
| Global (`--global`) | `~/.agents/skills/spectacular/` | `~/.claude/skills/spectacular/` → above |

---

---

## Documentation

| Doc | What it covers |
|---|---|
| [docs/workflow.md](docs/workflow.md) | Practical end-to-end usage loop — init, briefing, requests, lifecycle, archive, current sync, memory |
| [docs/commands.md](docs/commands.md) | CLI command reference and agent skill triggers, including the boundary between shell commands and skill commands |
| [docs/troubleshooting.md](docs/troubleshooting.md) | Common setup, install, skill discovery, update, symlink, and workspace state issues |
| [docs/configuration.md](docs/configuration.md) | `config.yaml`, agent files, tool overrides, request naming, and `.spectacular.local/` |
| [docs/scaffold.md](docs/scaffold.md) | Complete `.spectacular/` directory spec — every file, frontmatter schema, creation rules, versioning |

---

<div align="center">

Built with [Claude Code](https://claude.ai/code)

</div>
