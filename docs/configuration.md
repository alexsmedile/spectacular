---
title: Configuration
description: config.yaml, agent files, tool overrides, request naming, and .spectacular.local/.
section: ""
status: stable
since: 0.1.0
updated: 2026-05-23
---

# Configuration

Spectacular configuration lives in `.spectacular/config.yaml`. The file tells the skill how to name requests, which agent context file to treat as primary, and which root documents are default context.

---

## Default config

`spectacular init` creates:

```yaml
project:
  name: my-app
  summary: ""

# Structural version of the .spectacular/ layout. Migrations bump this.
workspace_schema: "0.6"

# Provenance — which spectacular version scaffolded this workspace (write-once)
# and which last structurally touched it (bumped by migrate / doctor --fix).
created_with: "1.26.3"
last_touched_with: "1.26.3"

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
    - decisions/index.md

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

## `workspace_schema` + provenance *(v1.26.3+)*

Three top-level keys track *what version of Spectacular* built and touched this workspace — distinct from the project's own release version (which lives in `ROADMAP.md` / git tags).

```yaml
workspace_schema: "0.6"        # structural layout version — migrations bump this
created_with: "1.26.3"         # CLI/skill version at init (write-once, never changes)
last_touched_with: "1.26.3"    # version that last structurally touched it
```

| Field | Meaning | Written by |
|---|---|---|
| `workspace_schema` | Structural version of the `.spectacular/` layout. Advances through the migration chain (`0.4 → 0.5 → 0.6 → …`). Decoupled from the CLI version. | `init`, `migrate` |
| `created_with` | The CLI/skill version that scaffolded this workspace. **Write-once** — never overwritten, so it survives migrations as a permanent origin marker. On workspaces created before v1.26.3 it's backfilled as `"unknown"`. | `init` (once); backfilled by `doctor --fix` |
| `last_touched_with` | The CLI version that last structurally mutated the workspace. Bumped by `migrate` and `doctor --fix`. Answers "was this edited by a newer version than it was born with?" | `migrate`, `doctor --fix` |

**Checking alignment:**

- `spectacular status --against-latest` — one-line verdict: up-to-date or behind the CLI's expected schema, with the catch-up verb.
- `spectacular doctor workspace` — emits a warning when `workspace_schema` is behind (run `migrate`) *or* ahead (update the CLI).
- `.spectacular/migrations.log` — append-only history of every applied migration (`<date>  <from> → <to>  (spectacular <version>)`), so upgrade provenance is auditable, not just inferred from the current schema.

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
    - decisions/index.md
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
    - decisions/index.md
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

## `last_build` *(v1.17.0+)*

The monotonic counter behind roadmap **build ids**. Each `spectacular new` stamps the
next id (`build: bN`) on the request's `PLAN.md` and increments this field.

```yaml
last_build: 18
```

You normally never edit this by hand — `spectacular new` owns it. It exists so build
ids are globally unique and never reused, even across reslotted or abandoned requests
(gaps in the sequence are normal and fine). Build ids are the request's permanent
identity in the [roadmap ledger](versioning.md#the-roadmap-ledger--how-builds-map-to-versions);
the version a build targets lives in the ledger's `target-version` column, never here
and never on the request.

---

## `snapshots` *(v1.24.0+)*

Controls the snapshot store directory, retention, and whether snapshots are gitignored.
All fields are optional and default sanely when the block or any field is absent — a
workspace with no `snapshots:` block behaves exactly as the defaults below.

```yaml
snapshots:
  folder: _snapshots   # store dir name under .spectacular/   (default _snapshots)
  keep: 3              # recent-tier count per doc            (default 3)
  period: month        # periodic-tier bucket: month|week|off (default month)
  gitignore: false     # gitignore the store by default?      (default false)
```

### Fields

- **`folder`** — the store directory under `.spectacular/`. Default `_snapshots` (the
  `_` prefix marks it a non-content / scanner-skip layer, like `_archive/`). A workspace
  still on the old `snapshots/` dir is flagged by `doctor snapshots`; `doctor --fix`
  renames it (git-mv when tracked) losslessly.
- **`keep`** — how many of the most-recent snapshots (by version ordinal) the **recent**
  retention tier keeps. Default 3.
- **`period`** — the calendar bucket for the **periodic** tier: `month`, `week`, or `off`.
  The tier keeps the newest snapshot per bucket, keyed off each snapshot's `updated:`
  frontmatter date. `off` collapses retention to origin + recent. Default `month`.
- **`gitignore`** — when `true`, `init` and `doctor --fix` add `.spectacular/<folder>/`
  to `.gitignore`; when `false` (default), the store stays committed. Toggling the value
  and re-running `doctor --fix snapshots` adds or removes the ignore line.

Retention is **tiered and generational** — a snapshot kept by *any* of origin (`@v1`),
periodic, or recent survives. Apply it with [`spectacular snapshot prune`](commands.md#spectacular-snapshot-prune---apply)
(dry-run by default). See also [`spectacular snapshot`](commands.md#spectacular-snapshot-file).

---

## `convention_pack` *(v0.4.0+)*

Opt-in. Declares which convention pack the repo follows and how strictly it's enforced.

```yaml
convention_pack:
  source: alex-default        # pack name; resolved via 4-tier scope precedence
  mode: scaffold              # suggest | scaffold | enforce
  overrides: []               # reserved for v2 (modular packs) — unused in v0.4.0
```

### Fields

**`source`** (required) — pack id (kebab-case). The CLI resolves this name via the precedence chain `project-local → user → app-store → bundled` and uses the first match. Run `spectacular pack list` to see what's available.

**`mode`** (required when `convention_pack:` is declared) — how the pack interacts with init + doctor:

| Mode | Init behavior | Doctor `conventions` area |
|---|---|---|
| `suggest` | Pack read, never applied automatically. Skill may surface pack opinions during interactive work. | Reports the pack is active; no drift checks. |
| `scaffold` | Init appends pack's `gitignore.always-add` entries to `.gitignore` (deduplicated). Always-set wins on conflicts. | Flags missing entries as warnings (exit 1). |
| `enforce` | Same as scaffold. | Flags missing entries as errors (exit 2). `spectacular doctor --fix` mechanically repairs. |

**`overrides`** (reserved) — declared but unused in v0.4.0. v2 ([convention-pack-modules](../.spectacular/requests/convention-pack-modules/)) will use this to skip specific rules from a pack while keeping the rest.

### Resolution and precedence

Same pack name in multiple scopes: most-specific wins. If `alex-default` exists in both `~/.spectacular/packs/` (user) and `<project>/.spectacular/packs/` (project-local), the project-local copy is used. v0.4.0 does **not** merge across scopes — first hit wins entirely.

### Errors and recovery

- **Pack source not found** — `doctor conventions` reports `❌ convention_pack source '<name>' not found in any scope`. Fix: `spectacular pack install <name>` or change the source field.
- **Pack file corrupt** — pack.md frontmatter unparseable → init logs a warning and skips pack scaffold. Doctor flags as error.
- **Unknown mode** — falls back to `suggest`; doctor adds an info note.

Full pack schema in [`skills/spectacular/references/packs-contract.md`](../skills/spectacular/references/packs-contract.md).

---

## `policies` *(v1.12.0+)*

An **override layer** over `POLICY.md` (the always-set practice layer). `POLICY.md` *defines* each policy — hook, check, severity, prose. This block *tunes* the contract for this project: enable/disable a shipped default, change a severity, or register a custom policy. They are layers, not competing copies — `spectacular policy` reads POLICY.md, applies these overrides, returns the merged result.

```yaml
policies:
  understand-before-change:
    enabled: false          # disable a shipped default
  scaffold-contract:
    severity: block         # override shipped severity (warn → block)
  no-secrets-in-memory:     # register a custom policy (id = key)
    hook: "@Remember"       # required for custom: one of the 9 hooks
    severity: warn
    check: "memory entry contains no API keys, tokens, or passwords"
```

### Fields (per policy id)

- **`enabled`** — `true` (default) or `false`. A disabled policy is still listed by `spectacular policy` but marked `[disabled]` and never enforced.
- **`severity`** — `block` (refuse to proceed) or `warn` (surface + continue). Overrides the value in POLICY.md. **Severity is opt-in to blocking**: a policy hard-stops only with an explicit `block`; absent/`warn`/unrecognized → non-blocking.
- **`hook`** *(custom only)* — one of `@Init @Planning @Implementation @Verification @Archive @Debugging @Remember @Snapshot @SessionEnd`. Required to register a policy that isn't in POLICY.md.
- **`check`** — the condition (required for blockers).
- **`principle`** — optional integer linking the principle this policy enforces; `spectacular policy` pulls that one line alongside it.

### Notes

- **Scope is config-only in v1.** A single `policies:` block in `.spectacular/config.yaml`. The 4-tier scope precedence used by convention packs is a v2 candidate, not built.
- `doctor policies` validates the merged contract's structure and the `## Understanding` gate on active requests.

Full policy schema in [`skills/spectacular/references/policies-contract.md`](../skills/spectacular/references/policies-contract.md).

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
