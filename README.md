<div align="center">

<img src="docs/assets/banner.svg" width="100%" alt="Spectacular — a spec-driven protocol for building with AI agents. Agents build. Humans decide." />

![License](https://img.shields.io/badge/license-MIT-blue)
![Claude Code](https://img.shields.io/badge/Claude%20Code-compatible-blueviolet)
![Platform](https://img.shields.io/badge/platform-Claude%20%7C%20Codex%20%7C%20Cursor-lightgrey)
![Version](https://img.shields.io/badge/version-1.26.1-green)

</div>

---

## No spec. No plan. No clue.

**Agents act before they understand.**

AI agents act on the first thing they find. Hand one a request and it starts coding immediately — without a spec, without a plan, without knowing *what* the project is or *why* the work matters. It can't see what was decided, what phase you're in, or what comes next.

The bottleneck isn't writing code. It's giving agents — and yourself — the context to write the *right* code. That context has to be planned, not assumed.

---

## A spec-driven protocol for agents

Spectacular is a spec-driven protocol for ideating, planning, scaffolding, and acting on projects with AI agents — keeping track of the changes, phases, tasks, and decisions so agents always know *what* they're doing and *why*. Drop a `.spectacular/` directory in any repo and it becomes the spec and shared context every agent works from.

Strategic context is split across seven focused canonical docs (PRD / PRINCIPLES / ARCHITECTURE / ROADMAP / STACK / DECISIONS / AGENTS) so agents load only what each task needs, not the entire repo.

It ships as three layers:

- **Convention** — a structured directory contract separating strategy, current truth, and active work
- **Skill** — a `/spectacular` Claude Code slash command that reads state, proposes actions, scaffolds files, and manages lifecycle transitions
- **CLI** — `spectacular init` bootstraps the workspace in any project in seconds

The skill is the primary interface. The CLI runs once.

> **Agents build. Humans decide.** Agents are made for building — humans, for deciding what's worth building.

In short, here's how spectacular helps you build your projects:

<img src="docs/assets/benefits.svg" width="100%" alt="Six benefits: Stay coherent · Plan, then build · Start in seconds · Roadmap the runway · Trust the state · Lose nothing" />

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
> `spectacular init` is zero-prompt by default. It infers the project name from the folder slug and uses `AGENTS.md` as the primary agent context file. Pass `-i` for interactive setup, or use flags: `--name`, `--summary`, `--agents-file`, `--skill-scope <project\|global\|none>`. If spectacular is already installed (globally, upstream, or as a plugin) init detects it, warns, and skips the redundant skill install.

---

## The shape of it

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

New here? Two things to understand. First, what lives at the top of `.spectacular/`:

<picture>
  <source media="(prefers-color-scheme: dark)" srcset="docs/assets/workspace-tree-dark.svg" />
  <img src="docs/assets/workspace-tree-light.svg" width="100%" alt=".spectacular/ holds three things: PRD.md (the intent), SPEC.md (the truth), and requests/ (the work)" />
</picture>

And second, what's inside any one request:

<picture>
  <source media="(prefers-color-scheme: dark)" srcset="docs/assets/request-tree-dark.svg" />
  <img src="docs/assets/request-tree-light.svg" width="100%" alt="A request folder holds PLAN.md (the approach), TASKS.md (the steps), VERIFY.md (the proof), and an optional specs/ folder" />
</picture>

That's the whole idea. The full layout, once a project fills in:

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
├── feedback/           # prototyping-mode feedback entries (v1.6.0+)
├── ideas/              # thinking scratchpad — not acted on automatically (v1.7.0+)
├── debug/              # live debug-job traces — one folder per bug (v1.26.0+)
├── audit/              # diagnostic examinations earned at resolution (v1.25.0+)
├── fixes/              # reusable verified-fix library, greppable (v1.25.0+)
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
| `spectacular decide "<decision>"` | Append an ADR entry; `--context`/`--consequences` fill the other sections (v1.8.4+) |
| `spectacular feedback-loop` | Prototyping-mode human-feedback loop — pick target, craft proposal, ask user, capture, decide (v1.6.0+) |
| `spectacular idea` | Thinking-scratchpad ideas — `new\|list\|promote`. Promotion scaffolds a request + archives the idea (v1.7.0+) |
| `spectacular summary` | One-page workspace overview: counts of requests/decisions/memories/sessions/ideas/feedback (v1.8.0+) |
| `spectacular requests` | List requests; filter with `--status`, `--active`, `--since`; `--json` for agents (v1.8.0+) |
| `spectacular request <slug>` | Skim view of one request (frontmatter + outline + milestone progress); `--full` for raw (v1.8.0+) |
| `spectacular decisions` / `decision <slug>` | List/inspect decisions (v1.8.0+) |
| `spectacular memories` / `memory <slug>` | List/inspect memories (v1.8.0+) |
| `spectacular sessions` / `sessions show <slug>` | List/inspect sessions read-only (v1.8.0+) |
| `spectacular show <doctype>` | Dump canonical doc (`prd\|spec\|principles\|...`); `--section <h2>` filters (v1.8.0+) |
| `spectacular progress <slug>` | Milestone tick rate from TASKS.md (v1.8.0+) |
| `spectacular paths` | JSON map of conventional workspace paths (v1.8.0+) |
| _(bug report to `/spectacular`)_ | Debug agent fleet — the skill checks prior `fixes/`, decides audit-first vs just-fix, delegates to read-only `debug-investigator` / apply-only `debug-fixer` / `debug-researcher`, then graduates a verified fix to the library. Traces land in `debug/`; reusable remedies in `fixes/` (v1.26.0+) |
| `spectacular fix new` / `audit new` | Log a verified fix / open a diagnostic audit (the debug-fleet library verbs); `fix new --debug-job <slug>` back-links the trace (v1.25.0+, `--debug-job` v1.26.0) |

> [!NOTE]
> **`spectacular docs *` verbs were removed in v1.17.0.** Public-facing documentation work lives in the standalone [pageworks](https://github.com/alexsmedile/pageworks) skill. `doctor docs` (discovery-only) remains. Install pageworks via its own one-liner.

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
spectacular init --skill-scope global         # install skill to ~/.agents/ and ~/.claude/
spectacular init --skill-scope none           # scaffold only; skip skill install (already installed elsewhere)
spectacular init --update                     # re-download latest skill release

spectacular doctor                            # substrate self-check (all areas)
spectacular doctor <area>                     # scoped: skill | workspace | frontmatter | snapshots | links | lifecycle | kits | conventions | specs | docs | personas | memory | sessions | feedback | ideas | debug | policies | vision | decisions | roadmap
spectacular doctor --fix                      # apply mechanical fixes (gitignore, missing dirs, dangling symlinks, pack drift, legacy current/ migration)
spectacular doctor --format json              # JSON report for the skill or other tools

spectacular feedback-loop new <target>        # scaffold a feedback entry (status: open); --request <slug> to scope
spectacular feedback-loop list                # list entries across .spectacular/feedback/ + per-request folders
spectacular feedback-loop resolve <slug> --next-action <a>   # close with a decision (required flag)
spectacular feedback-loop archive <slug>      # move to .spectacular/archive/feedback/<year>/

spectacular idea new <slug>                   # scaffold idea entry (status: parked); --hypothesis, --origin, --priority
spectacular idea list                         # list ideas; --status parked|exploring|promoted|all
spectacular idea promote <slug>               # scaffold request from idea + move source to archive/ideas/

spectacular summary                           # one-page workspace overview (v1.8.0+)
spectacular requests --active                 # list active requests (table + --json)
spectacular requests --status review --since 7d --limit 10
spectacular request <slug>                    # skim view (frontmatter + outline + progress); --full for raw
spectacular decisions [--tag t] [--since 30d] [--json]
spectacular decision <slug>                   # skim view of one decision
spectacular memories [--tag t] [--since 7d] [--json]
spectacular memory <slug>                     # skim view of one memory
spectacular sessions [--status open|closed|all] [--since 24h]
spectacular sessions show <slug>              # detail subverb (avoids `session start|end` collision)
spectacular show <doctype>                    # prd|spec|principles|architecture|roadmap|stack|agents|...
spectacular show <doctype> --section <name>   # filter to one H2 section
spectacular progress <slug>                   # M1: 8/8 ✓, M2: 3/5, ... (parses TASKS.md ticks)
spectacular paths                             # JSON map of workspace paths

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

## Works well with

Spectacular owns one thing — `.spectacular/`: strategy, current truth, active work. It composes with focused sibling tools rather than absorbing their jobs.

### Best with — pageworks

**Public-facing documentation (the `docs/` surface) is owned by [pageworks](https://github.com/alexsmedile/pageworks)**, a sibling skill extracted in v1.2.0. The two compose: when a SPEC-touching request archives, spectacular asks whether `docs/` should be updated and hands off to pageworks if you confirm. `spectacular doctor docs` reports discovery only (folder + manifest presence; install hint if pageworks missing) — never validates schema. There is no automatic invocation across the boundary.

```bash
curl -fsSL https://raw.githubusercontent.com/alexsmedile/pageworks/main/cli/install.sh | bash
```

### Also pairs with

- **Claude Code / Codex / Cursor** — the agents that read `.spectacular/` and act on it. Spectacular is the protocol; these are the runtimes.
- **Git** — `.spectacular/` is fully committed, so phases, decisions, and changelogs version alongside your code. `.spectacular.local/` stays gitignored.

> Building a tool that composes with the `.spectacular/` contract? [Open an issue](https://github.com/alexsmedile/spectacular/issues) — this list is meant to grow.

---

## Built for

- Solo developers using Claude Code, Codex, or Cursor on projects that span weeks or months
- Small teams where AI agents need to share operational context
- Anyone who has re-explained the same architectural decision to an agent more than twice

## Not built for

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
| Global (`--skill-scope global`) | `~/.agents/skills/spectacular/` | `~/.claude/skills/spectacular/` → above |

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
