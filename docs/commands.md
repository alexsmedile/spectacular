# Commands

Spectacular has two command surfaces:

- the local CLI command: `spectacular init`, `spectacular doctor`
- the agent skill triggers: `/spectacular`, `spectacular new`, `spectacular archive`, doc verbs, and related workflow commands

The CLI bootstraps the workspace and runs substrate self-checks. The skill operates the workspace and writes content.

---

## CLI commands

The installed shell command currently supports two subcommands:

```bash
spectacular init [options]
spectacular doctor [<area>] [options]
```

### `spectacular init`

Scaffolds `.spectacular/` and installs the skill. As of v0.3.0, init scaffolds the **5-file always-set** by default — `PRD.md`, `config.yaml`, `<agents-file>` (default `AGENTS.md`), `requests/`, `current/`. Extra docs come from the selected kit or explicit `--with` flag.

```bash
spectacular init                              # always-set + blank kit
spectacular init -i                           # interactive: kit menu + per-doc prompts
spectacular init --kit coding                 # always-set + coding kit's STACK + ARCHITECTURE
spectacular init --with principles,roadmap    # additive: those two on top of always-set
spectacular init --kit coding --minimal       # always-set only; kit identity preserved
spectacular init --name my-app
spectacular init --summary "Internal dashboard for support workflows"
spectacular init --agents-file CLAUDE.md      # for Claude-only teams
spectacular init --global                     # install skill to ~/.agents and ~/.claude
spectacular init --update                     # re-download latest skill release
```

Pre-flight is **always non-destructive** — re-running init on an existing workspace never overwrites existing files. Empty files are filled; non-empty files are skipped; malformed files trigger a `spectacular doctor frontmatter` hint.

### `spectacular init --update`

Re-downloads the latest skill release and updates `.spectacular/skills.lock`. Does not rewrite workspace documents.

```bash
spectacular init --update
```

### `spectacular doctor`

Environment / infrastructure self-check. Read-only by default; `--fix` applies content-free mechanical repairs (gitignore, missing dirs, dangling symlinks, missing always-set stubs). Judgment-requiring fixes route to the skill (`/spectacular doctor --fix`).

```bash
spectacular doctor                        # all areas
spectacular doctor frontmatter            # scoped: only frontmatter checks
spectacular doctor --fix                  # apply mechanical fixes interactively
spectacular doctor --format json          # JSON report for skill/tool consumption
```

Available areas: `skill`, `workspace`, `frontmatter`, `snapshots`, `links`, `lifecycle`, `kits`.

Exit codes:
- `0` clean (no warnings, no errors)
- `1` warnings only
- `2` one or more errors

Severity legend:
- `✅` pass
- `⚠️` warning (non-blocking drift)
- `❌` error (blocking — substrate is broken)
- `ℹ️` info (informational, e.g. snapshot counts)

---

## Skill triggers

The following commands are not shell CLI subcommands. Use them in an agent conversation where the Spectacular skill is installed.

### `/spectacular`

Reads the workspace and returns a project briefing. Use at the start of a session.

```text
/spectacular
```

If the skill can't parse the workspace state, it auto-runs `spectacular doctor workspace frontmatter kits` and surfaces findings inline.

### `spectacular status`

Same as `/spectacular`.

### `spectacular new <description>`

Creates a new request folder. The skill derives a kebab-case slug, checks for collisions, applies the [[verification]] 2-of-6 rule to decide whether to scaffold a `VERIFY.md`, and asks the user for confirmation before writing.

```text
spectacular new add team billing
```

Output:

```text
.spectacular/requests/add-team-billing/
├── PLAN.md
└── TASKS.md
(+ VERIFY.md if 2-of-6 rule triggers)
```

### `spectacular archive <slug>`

Archives a verified request. The skill proposes `current/` updates and memory entries before moving the request to `.spectacular/archive/`.

```text
spectacular archive add-team-billing
```

### `spectacular remember this`

Writes an operational lesson to `.spectacular/memory/` after human confirmation. Team-visible — not for personal notes or secrets.

### `spectacular snapshot <file>`

Creates a versioned snapshot (`<FILE>@vN.md`) before editing a canonical document.

```text
spectacular snapshot .spectacular/PRD.md
```

Canonical files: `PRD.md`, `PRINCIPLES.md`, `ARCHITECTURE.md`, `ROADMAP.md`, `STACK.md`, `DECISIONS.md`, `AGENTS.md`, `config.yaml`, `current/` capability specs.

### `spectacular promote <idea>`

Promotes an idea file into a full request and moves the original to `.spectacular/archive/ideas/`.

### `/spectacular doctor [<area>] [--fix]`

Skill side of the doctor command. Reads the CLI's JSON report and walks each judgment-requiring finding interactively (per-finding `y/n/q`). Snapshots canonical docs before any edit.

```text
/spectacular doctor             # walk findings conversationally
/spectacular doctor --fix       # interactive repair flow
```

### Doc-writing verbs

Spectacular's registry-driven engine treats every canonical doc the same way. For any registered doc type:

```text
spectacular <doc>               # grill if empty, review if filled
spectacular <doc> grill         # interactive slot-filling
spectacular <doc> refine        # vibe → spec rewrite
spectacular <doc> review        # quality gate (pass/fail)
```

Doc IDs in v0.3.x: `prd`, `plan`, `tasks`, `principles`, `architecture`, `roadmap`, `stack`, `agents`, `decisions`.

Legacy aliases (backwards-compatible from v0.2.x): `spectacular prd`, `spectacular prd grill`, `spectacular prd refine`, `spectacular prd review`.

### `spectacular prd` / `spectacular prd grill`

Walks the user through the **8-slot canonical PRD** (Vision / Problem / Target users / Deliverable / Goals & success criteria / Non-goals / Constraints / First milestone), one question at a time, with kit-aware prompts.

```text
spectacular prd
```

Asks for a kit first (`blank` / `coding` / `content` / `product` / `research`), then drives the interview. Each kit adds its own extra slots — `coding` adds Stack + Interfaces; `product` adds User stories + Metrics + Distribution; etc. See [kits-contract.md](../skills/spectacular/references/kits-contract.md).

### `spectacular prd refine`

Runs vibe→spec rewrite patterns on an existing PRD — flags vague adjectives, plural users, unbounded success criteria, vague deliverables, and inserts `[NEEDS CLARIFICATION: …]` markers where it can't resolve.

### `spectacular prd review`

Runs the PRD quality gate — 10 checks total (8 base + 2 kit-aware). Reports a punch list; user fixes.

---

## Common confusion

Do not run skill triggers in your shell:

```bash
spectacular new add team billing      # ✗ shell doesn't implement this
```

The shell CLI implements only `init` and `doctor`. Everything else is a skill trigger for your coding agent.

```bash
spectacular init                       # ✓ shell
spectacular doctor                     # ✓ shell

/spectacular                           # ✓ agent
spectacular new add team billing       # ✓ agent
spectacular archive add-team-billing   # ✓ agent
spectacular prd                        # ✓ agent
```
