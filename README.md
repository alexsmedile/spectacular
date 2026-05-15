<div align="center">

<img src="docs/assets/banner.svg" width="100%" alt="Spectacular — Other tools help you write code faster. Spectacular helps you stay coherent while you do." />

![License](https://img.shields.io/badge/license-MIT-blue)
![Claude Code](https://img.shields.io/badge/Claude%20Code-compatible-blueviolet)
![Platform](https://img.shields.io/badge/platform-Claude%20%7C%20Codex%20%7C%20Cursor-lightgrey)
![Version](https://img.shields.io/badge/version-0.1.1-green)

</div>

---

## The problem

AI coding agents are fast. They are not coherent.

Every new session starts cold. Decisions made last week are invisible. Context that took an hour to build evaporates the moment the conversation ends. Your agent doesn't know what was decided, what was tried, or what's next — so you re-explain, re-decide, and re-discover the same things, over and over.

The bottleneck isn't code generation. It's everything around it.

---

## What it is

Spectacular is an operational workspace for AI-assisted software projects. Drop a `.spectacular/` directory in any repo and it becomes the shared brain for you and every agent that touches the project.

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

# Open Claude Code and run
/spectacular
```

> [!TIP]
> `spectacular init` scaffolds the `.spectacular/` directory and installs the `/spectacular` skill into `.agents/skills/spectacular/` (source) and `.claude/skills/spectacular/` (symlink). After init, `/spectacular` is immediately available in Claude Code.

> [!NOTE]
> `spectacular init` is zero-prompt by default. It infers the project name from the folder slug and uses `AGENTS.md` as the primary agent context file. Pass `-i` for interactive setup, or use flags: `--name`, `--summary`, `--agents-file`, `--global`.

---

## How it works

### The three layers

<img src="docs/assets/layer-map.svg" width="100%" alt="Strategy, truth, and work mapped to the .spectacular/ directory" />

**Strategy** — changes rarely. Documents why the project exists, what stack it uses, and why decisions were made. Agents load this once at the start of a session.

**Current truth** — reflects actual system behavior right now. Modular capability specs (auth, billing, editor). Never overwritten in place — the skill snapshots before proposing edits.

**Active work** — temporary. Each request gets a folder with `PLAN.md` (intent) and `TASKS.md` (execution checklist). Archived on completion, not deleted.

---

### The lifecycle

<img src="docs/assets/lifecycle-flow.svg" width="100%" alt="Request lifecycle: planned → active → review → verified → archived" />

State lives in `PLAN.md` frontmatter. The skill reads it on every invocation and surfaces the highest-priority next action. When all tasks are checked, it proposes moving to `review`. When the checklist passes, it proposes `archived` — and offers to update `current/` with anything that changed.

---

### The workspace

```
.spectacular/
├── PRD.md              # product intent
├── STACK.md            # tech + rules
├── DECISIONS.md        # why we chose what we chose
├── AGENTS.md           # context loading rules for this project
├── config.yaml         # naming, required files, agent file overrides
│
├── current/            # canonical system truth (capability specs)
│   ├── auth/
│   └── billing/
│
├── requests/           # active and planned work
│   └── add-team-billing/
│       ├── PLAN.md     # goal, approach, success criteria (owns lifecycle state)
│       ├── TASKS.md    # executable checklist
│       └── SESSION.md  # current execution state
│
├── memory/             # long-term operational learning (git-committed)
│   ├── lessons.md
│   └── failures.md
│
└── archive/            # completed requests (never deleted)
```

`.spectacular.local/` — personal overrides, always gitignored.

---

## Skill commands

| Command | What happens |
|---|---|
| `/spectacular` | Project briefing — active requests, draft capabilities, next action |
| `spectacular new <description>` | Scaffold a new request with PLAN.md + TASKS.md |
| `spectacular archive <slug>` | Archive a completed request, propose `current/` sync + memory entries |
| `spectacular remember this` | Write an insight to `memory/` immediately |
| `spectacular snapshot <file>` | Snapshot a canonical document before editing |
| `spectacular promote <idea>` | Promote an idea file to a full request |
| `spectacular status` | Same as no-arg invocation |

---

## CLI reference

```
spectacular init                          # zero prompts, all defaults
spectacular init -i                       # interactive — prompts for all settings
spectacular init --name my-app            # set project name
spectacular init --summary "..."          # set project summary
spectacular init --agents-file CLAUDE.md  # use CLAUDE.md instead of AGENTS.md
spectacular init --global                 # install to ~/.agents/ and ~/.claude/
spectacular init --update                 # re-download latest skill release
```

> [!TIP]
> Claude-only team? Use `--agents-file CLAUDE.md`. Multi-tool team? Keep `AGENTS.md` as primary and add `tool_overrides.claude: CLAUDE.md` to `config.yaml` — the skill will surface both.

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

**Skill only** (no CLI) — manual copy or via `apm`:

```bash
# via apm
apm --mode skills install spectacular

# manual
cp -r skills/spectacular ~/.claude/skills/
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
