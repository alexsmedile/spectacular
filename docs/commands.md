# Commands

Spectacular has two command surfaces:

- the local CLI command: `spectacular init`
- the agent skill triggers: `/spectacular`, `spectacular new`, `spectacular archive`, and related workflow commands

The CLI bootstraps the workspace. The skill operates the workspace.

---

## CLI commands

The installed shell command currently supports one subcommand:

```bash
spectacular init [options]
```

Use it from a project root to create `.spectacular/` and install the skill.

### `spectacular init`

Creates the default workspace with no prompts.

```bash
spectacular init
```

Defaults:

- project name is inferred from the current folder
- project summary is empty
- primary agent file is `AGENTS.md`
- skill install scope is project-local

### `spectacular init -i`

Runs interactive setup.

```bash
spectacular init -i
```

Prompts for:

- project name
- project summary
- primary agent file
- install scope

### `spectacular init --name <slug>`

Sets the project name written to `.spectacular/config.yaml`.

```bash
spectacular init --name my-app
```

### `spectacular init --summary <text>`

Sets the project summary written to `.spectacular/config.yaml`.

```bash
spectacular init --summary "Internal dashboard for support workflows"
```

### `spectacular init --agents-file <file>`

Sets the primary agent context file created inside `.spectacular/`.

```bash
spectacular init --agents-file CLAUDE.md
```

Use `CLAUDE.md` for Claude-only projects. Use `AGENTS.md` for multi-tool projects.

### `spectacular init --global`

Installs the skill globally instead of project-locally.

```bash
spectacular init --global
```

Global install paths:

```text
~/.agents/skills/spectacular/
~/.claude/skills/spectacular/ -> ~/.agents/skills/spectacular/
```

Project-local install paths:

```text
<project>/.agents/skills/spectacular/
<project>/.claude/skills/spectacular/ -> <project>/.agents/skills/spectacular/
```

### `spectacular init --update`

Re-downloads the latest skill release and updates `.spectacular/skills.lock`.

```bash
spectacular init --update
```

This updates the installed skill. It does not rewrite existing workspace documents such as `PRD.md`, `PRINCIPLES.md`, `ARCHITECTURE.md`, `ROADMAP.md`, `STACK.md`, or `DECISIONS.md`.

---

## Skill triggers

The following commands are not shell CLI subcommands. Use them in an agent conversation where the Spectacular skill is installed.

### `/spectacular`

Reads the workspace and returns a project briefing.

```text
/spectacular
```

Use this at the start of a session.

### `spectacular status`

Same intent as `/spectacular`: read state and surface the next action.

```text
spectacular status
```

### `spectacular new <description>`

Creates a new request folder under `.spectacular/requests/`.

```text
spectacular new add team billing
```

Expected output:

```text
.spectacular/requests/add-team-billing/
├── PLAN.md
└── TASKS.md
```

The skill should derive a kebab-case slug, check for collisions, and let the human override the slug.

### `spectacular archive <slug>`

Archives a verified request.

```text
spectacular archive add-team-billing
```

The skill should propose `current/` updates and memory entries before moving the request to `.spectacular/archive/`.

### `spectacular remember this`

Writes an operational lesson to `.spectacular/memory/` after human confirmation.

```text
spectacular remember this
```

Use this for team-visible lessons, not personal notes or secrets.

### `spectacular snapshot <file>`

Creates a versioned snapshot before editing a canonical document.

```text
spectacular snapshot .spectacular/PRD.md
```

Canonical documents include:

- `PRD.md`
- `PRINCIPLES.md`
- `ARCHITECTURE.md`
- `ROADMAP.md`
- `STACK.md`
- `DECISIONS.md`
- `AGENTS.md`
- `config.yaml`
- `current/` capability specs

### `spectacular promote <idea>`

Promotes an idea file into a full request.

```text
spectacular promote pricing-model-notes
```

The skill should create a request and move the original idea into `.spectacular/archive/ideas/`.

### `spectacular prd` / `spectacular prd grill`

Walks the user through the 6-slot canonical PRD (problem / who / success / non-goals / constraints / milestone), one question at a time, with kit-aware prompts. Asks for a kit first (`coding` / `product` / `content` / `research` / `blank`), then drives the interview.

```text
spectacular prd
```

Use on fresh projects or when `PRD.md` is empty / placeholder-only. If a real PRD already exists, the skill asks whether to refine in place or start fresh (snapshots first).

### `spectacular prd refine`

Runs vibe→spec rewrite patterns on an existing PRD — flags vague adjectives, plural users, unbounded success metrics, inserts `[NEEDS CLARIFICATION: …]` markers where it can't resolve.

```text
spectacular prd refine
```

### `spectacular prd review`

Runs the PRD quality gate — pass/fail check on all 6 required slots, vague-word scan, intent-trace check.

```text
spectacular prd review
```

---

## Common confusion

Do not run this in your shell:

```bash
spectacular new add team billing
```

The shell CLI does not implement `new`, `archive`, `remember`, `snapshot`, `promote`, or `status` as local subcommands. Those are skill triggers for your coding agent.

The shell command is for bootstrapping:

```bash
spectacular init
```

The agent skill is for ongoing project operation:

```text
/spectacular
spectacular new add team billing
spectacular archive add-team-billing
```
