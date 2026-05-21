# Configuration

Spectacular configuration lives in `.spectacular/config.yaml`. The file tells the skill how to name requests, which agent context file to treat as primary, and which root documents are default context.

---

## Default config

`spectacular init` creates:

```yaml
project:
  name: my-app
  summary: ""

naming:
  requests: kebab-case
  prefix: ""

required_files:
  requests:
    - PLAN.md
    - TASKS.md

agents:
  file: AGENTS.md
  tool_overrides:
    # claude: CLAUDE.md
  default_context:
    - PRD.md
    - STACK.md
    - DECISIONS.md

skills:
  symlink_on_init: []
```

---

## `project`

Project metadata used in briefings and generated files.

```yaml
project:
  name: my-app
  summary: "Internal dashboard for support workflows"
```

Fields:

| Field | Meaning |
|---|---|
| `name` | Short project slug or display name |
| `summary` | One-sentence project description |

The default `name` is inferred from the current directory when you run `spectacular init`.

---

## `naming`

Controls generated request names.

```yaml
naming:
  requests: kebab-case
  prefix: ""
```

Fields:

| Field | Meaning |
|---|---|
| `requests` | Slug style for request folders |
| `prefix` | Optional prefix for generated request slugs |

Recommended request slugs are kebab-case:

```text
add-team-billing
fix-session-timeout
replace-payment-provider
```

If `prefix` is set, the skill should apply it when proposing new request slugs.

---

## `required_files`

Defines the files every request must contain.

```yaml
required_files:
  requests:
    - PLAN.md
    - TASKS.md
```

`PLAN.md` owns request intent and lifecycle state. `TASKS.md` owns the executable checklist.

Avoid adding many required files. Use optional files such as `SESSION.md`, `RISKS.md`, and `VERIFY.md` only when the request needs them.

---

## `agents`

Controls which project context file the skill should read as the primary agent guide.

```yaml
agents:
  file: AGENTS.md
  tool_overrides:
    # claude: CLAUDE.md
  default_context:
    - PRD.md
    - STACK.md
    - DECISIONS.md
```

### `agents.file`

The primary agent instructions file inside `.spectacular/`.

Use this for multi-tool projects:

```yaml
agents:
  file: AGENTS.md
```

Use this for Claude-only projects:

```yaml
agents:
  file: CLAUDE.md
```

You can also create the Claude-only version during init:

```bash
spectacular init --agents-file CLAUDE.md
```

### `agents.tool_overrides`

Optional per-tool supplementary files.

For a multi-tool project that uses `AGENTS.md` as the shared base but also needs Claude-specific rules:

```yaml
agents:
  file: AGENTS.md
  tool_overrides:
    claude: CLAUDE.md
```

The shared file should contain rules that apply to every agent. Tool override files should contain only tool-specific differences.

### `agents.default_context`

Stable context the skill should treat as default project grounding:

```yaml
agents:
  default_context:
    - PRD.md
    - STACK.md
    - DECISIONS.md
```

Keep this list small. Spectacular works best when agents load targeted context instead of the full project history.

The **full per-task context map** lives in `.spectacular/AGENTS.md` (under "Context loading by task") — `default_context` is just the always-on baseline. AGENTS.md is the authoritative source; the skill reads it on every invocation.

---

## `skills`

Reserved for project-specific skill behavior.

```yaml
skills:
  symlink_on_init: []
```

In v1, Spectacular installs its own skill and leaves project-specific skill automation minimal. Keep this empty unless your project has a documented local convention.

---

## `.spectacular.local/`

`.spectacular.local/` is a personal override layer at the repository root.

Use it for:

- local development notes
- personal paths
- machine-specific settings
- private workflow preferences

Do not use it for:

- team decisions
- product context
- canonical system truth
- secrets that should be managed by a real secret store

`spectacular init` adds `.spectacular.local/` to `.gitignore`.

---

## Agent file strategies

### Multi-tool project

Use `AGENTS.md` as the primary file:

```bash
spectacular init --agents-file AGENTS.md
```

Config:

```yaml
agents:
  file: AGENTS.md
  tool_overrides:
    # claude: CLAUDE.md
```

This works well when Claude Code, Codex, Cursor, or other agents share the same workspace conventions.

### Claude-only project

Use `CLAUDE.md` as the primary file:

```bash
spectacular init --agents-file CLAUDE.md
```

Config:

```yaml
agents:
  file: CLAUDE.md
```

### Shared base with Claude-specific rules

Use `AGENTS.md` as the base and add a Claude override:

```yaml
agents:
  file: AGENTS.md
  tool_overrides:
    claude: CLAUDE.md
```

Put common context loading rules in `AGENTS.md`. Put Claude-only behavior in `CLAUDE.md`.

---

## Safe editing

`config.yaml` is canonical workspace configuration. Snapshot it before major edits:

```text
spectacular snapshot .spectacular/config.yaml
```

Then edit `.spectacular/config.yaml`.

Small corrections are usually fine, but changes that affect agent loading, request naming, or required files should be treated like project-level decisions.
