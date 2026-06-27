---
title: Commands
description: CLI subcommands and agent skill triggers, including the boundary between shell and skill.
section: ""
status: stable
since: 0.1.0
updated: 2026-05-23
---

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

Scaffolds `.spectacular/` and installs the skill. As of v0.5.0, init scaffolds the **6-file always-set** by default — `PRD.md`, `SPEC.md`, `config.yaml`, `<agents-file>` (default `AGENTS.md`), `requests/`, `specs/`. Extra docs come from the selected kit or explicit `--with` flag. (v0.3.0–v0.4.x scaffolded `current/` instead of `SPEC.md` + `specs/`; legacy workspaces are auto-migrated by `spectacular doctor specs --fix`.)

```bash
spectacular init                              # always-set + blank kit
spectacular init -i                           # interactive: kit menu + per-doc prompts
spectacular init --kit coding                 # always-set + coding kit's STACK + ARCHITECTURE
spectacular init --with principles,roadmap    # additive: those two on top of always-set
spectacular init --kit coding --minimal       # always-set only; kit identity preserved
spectacular init --name my-app
spectacular init --summary "Internal dashboard for support workflows"
spectacular init --agents-file CLAUDE.md      # for Claude-only teams
spectacular init --skill-scope global         # install skill to ~/.agents and ~/.claude
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

Available areas: `skill`, `workspace`, `frontmatter`, `snapshots`, `links`, `lifecycle`, `kits`, `conventions` *(v0.4.0+)*, `specs` *(v0.5.0+)*, `docs` *(v0.6.0+)*, `personas` *(v1.3.0+)*, `memory` / `sessions` *(v1.5.0+)*, `feedback` *(v1.6.0+)*, `ideas` *(v1.7.0+)*, `policies` *(v1.12.0+)*.

The `specs` area also flags **SPEC.md drift** *(v1.18.0+)* — a `⚠️` warning when `SPEC.md`'s `updated` date predates the newest archived request, signalling a likely missed spec-sync. It's a date heuristic, so it routes to the skill's spec-sync flow to confirm and reconcile content.

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

### `spectacular pack` *(v0.4.0+)*

Manage convention packs — opt-in repo-shape opinions (naming, taxonomy, gitignore, README contract, file-placement, project-type scaffolds).

```bash
spectacular pack list                         # show all installed packs across scopes
spectacular pack install <name>               # copy bundled/app-store pack to ~/.spectacular/packs/<name>/
spectacular pack install <name> --from <path> # install from arbitrary local folder
spectacular pack show <name>                  # print scope + path + pack.md frontmatter
spectacular pack remove <name>                # delete user-scope pack
spectacular pack remove <name> --force        # allow removing bundled/app-store/project-local
```

**Pack scopes** (precedence: project-local > user > app-store > bundled):

| Scope | Path | Notes |
|---|---|---|
| `bundled` | `skills/spectacular/templates/packs/<name>/` | Ships with the skill (`minimal` only) |
| `app-store` | `<spectacular-repo>/packs/<name>/` | Distributable via this repo |
| `user` | `~/.spectacular/packs/<name>/` | Installed via `pack install` |
| `project-local` | `<project>/.spectacular/packs/<name>/` | Per-project override |

**Activate a pack per-repo** by adding to `.spectacular/config.yaml`:

```yaml
convention_pack:
  source: <pack-name>
  mode: suggest | scaffold | enforce
```

See [configuration.md](configuration.md#convention-packs) for `convention_pack:` field semantics.

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

Creates a new request folder. The skill derives a kebab-case slug, checks for collisions, applies the [[verify]] 2-of-6 rule to decide whether to scaffold a `VERIFY.md`, and asks the user for confirmation before writing.

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

### `spectacular advance <slug> [--to <state>] [--force]` *(v1.19.0+)*

Advances a request one step through the lifecycle (`planned → active → review → verified`), or jumps with `--to`. Backward transitions require `--force`. The CLI is a dumb mutator — it edits `status:` in PLAN.md frontmatter only.

```text
spectacular advance add-team-billing            # one step forward
spectacular advance add-team-billing --to review
```

> **Renamed from `promote` in v1.19.0.** `spectacular promote <slug>` still works as a deprecated alias (prints a one-line notice). Not to be confused with `spectacular idea promote`, which promotes an *idea* into a request.

### `spectacular next` *(v1.19.0+)*

Prints the single highest-priority next action for the workspace. Read-only — mutates nothing. Order: an active request (keep going / advance) beats a review request (verify) beats a planned request (start it); an empty workspace is ushered into `spectacular new`.

```text
spectacular next
```

### `spectacular archive <slug>`

Archives a verified request. The skill proposes `current/` updates and memory entries before moving the request to `.spectacular/archive/`.

```text
spectacular archive add-team-billing
```

### `spectacular remember this`

Writes an operational lesson to `.spectacular/memory/` after human confirmation. Team-visible — not for personal notes or secrets. **Skill flow** — runs inside Claude Code/Codex.

### `spectacular remember "<text>" [--tag a,b]` (v1.5.0+)

CLI mutator. Writes one memory entry to `.spectacular/memory/<slug>.md` and regenerates `MEMORY.md` index. Auto-derives slug + summary. If a session is open, the entry frontmatter gets `session: <slug>` automatically.

```text
spectacular remember "haiku is fast enough for slug generation" --tag perf,cli
spectacular remember "..." --dry-run    # preview without writing
```

### `spectacular decide "<decision>" [--context "..."] [--consequences "..."]` (v1.5.0+)

CLI mutator. Appends one ADR-style entry (**Context / Decision / Consequences**) to `.spectacular/DECISIONS.md`. The positional argument fills `**Decision:**`; auto-derives a title slug from the first ~6 words of the decision. If a session is open, the entry includes a `Session:` link.

`--context` and `--consequences` (v1.8.4+) populate those sections at write time. Omitted sections are emitted as empty headers to fill in later — never invented from the decision text.

```text
spectacular decide "use bash for the CLI to keep install footprint zero"
spectacular decide "use bash for the CLI" \
  --context "want zero-install distribution across varied targets" \
  --consequences "ships everywhere with no runtime; harder to unit-test"
spectacular decide "..." --dry-run    # preview without writing
```

### `spectacular session start|end` (v1.5.0+)

CLI mutator. Opens or closes a working session entry in `.spectacular/sessions/`.

```text
spectacular session start --tag substrate-work     # open
spectacular session end                            # close, recompute linked counts
```

At most **one** session can be open at a time. At `end`, the writer scans `DECISIONS.md` + `memory/*.md` for entries with matching `session: <slug>` and appends Linked-decisions / Linked-memories sections to the session body. `spectacular doctor sessions` warns on sessions open >4h.

### `spectacular snapshot <file>`

Creates a versioned snapshot (`<FILE>@vN.md`) before editing a canonical document.

```text
spectacular snapshot .spectacular/PRD.md
```

Canonical files: `PRD.md`, `PRINCIPLES.md`, `ARCHITECTURE.md`, `ROADMAP.md`, `STACK.md`, `DECISIONS.md`, `AGENTS.md`, `config.yaml`, `current/` capability specs.

### `spectacular idea promote <idea>`

Promotes an idea file into a full request and moves the original to `.spectacular/archive/ideas/`.

### `spectacular policy [@hook | <id> | --principle N | --json]` *(v1.12.0+)*

Reads the merged **policy contract** — `POLICY.md` (the always-set practice layer) with any `config.yaml` `policies:` overrides applied. The skill calls `spectacular policy @<hook>` on entering a work phase to retrieve only that phase's rules (progressive disclosure).

```text
spectacular policy                 # all policies, grouped by hook (skim)
spectacular policy @Verification   # one hook's policies + linked principle lines
spectacular policy understand-before-change   # one policy, full text + its principle
spectacular policy --principle 7   # reverse: which policies enforce principle 7
spectacular policy --json          # machine form (skill-consumed)
```

Hooks (the locked 8): `@Init`, `@Planning`, `@Implementation`, `@Verification`, `@Archive`, `@Remember`, `@Snapshot`, `@SessionEnd`. A policy blocks a transition only if it declares `severity: block`; `warn` and unset are surface-and-continue. `spectacular advance` prints an advisory at the spine transitions, and `spectacular doctor policies` reports structural + `## Understanding`-gate findings. See [policies-contract](../skills/spectacular/references/policies-contract.md) for the schema.

### `/spectacular doctor [<area>] [--fix]`

Skill side of the doctor command. Reads the CLI's JSON report and walks each judgment-requiring finding interactively (per-finding `y/n/q`). Snapshots canonical docs before any edit.

```text
/spectacular doctor             # walk findings conversationally
/spectacular doctor --fix       # interactive repair flow
```

### Doc-writing verbs

Spectacular treats every canonical doc the same way — each doc has a rules file (`skills/spectacular/references/<doc-id>-rules.md`) declaring its mode, slots, and template. For any known doc type:

```text
spectacular <doc>               # grill if empty, review if filled
spectacular <doc> grill         # interactive slot-filling
spectacular <doc> refine        # vibe → spec rewrite
spectacular <doc> review        # quality gate (pass/fail)
```

**Agentic vs mechanical verbs** *(v1.4.0+)*

| Verb | Runs in | Why |
|---|---|---|
| `grill`, `refine` | **Skill only** | Require an LLM to interview, mini-refine, vibe→spec rewrite |
| `review` | **Mixed** | Structural checks run in CLI (`doctor`); semantic checks need the skill |
| `new`, `archive`, `snapshot`, `init`, `doctor`, `pack`, `migrate` | **CLI primarily** | Mechanical scaffolding, file moves, integrity checks |

If you type `spectacular <doc> grill` at terminal, the CLI prints a friendly redirect telling you to run `/spectacular <doc> grill` inside Claude Code or Codex. This is by design — the agentic verbs need an LLM that the CLI doesn't have.

**Grill sub-modes** *(v1.4.0+)*

Each doc declares its grill style in `mode:`. User can override per session:

| Mode | Behavior | Example doc |
|---|---|---|
| `grill` / `grill-wide` | Walk all slots once | PRD, PLAN |
| `grill-each` | Per-block walk; "add another?" loop | ROADMAP, PERSONAS |
| `grill-loop` | Wide pass, then deep on vague slots | (opt-in via `--loop`) |

Flag override: `spectacular roadmap grill --wide` forces the wide style regardless of the doc's declared mode.

Doc IDs in v1.4.0: `prd`, `plan`, `tasks`, `principles`, `architecture`, `roadmap`, `stack`, `agents`, `decisions`, `personas`, `spec`, `convention-pack`. See [`doc-index.md`](../skills/spectacular/references/doc-index.md) for the full catalog.

Legacy aliases (backwards-compatible from v0.2.x): `spectacular prd`, `spectacular prd grill`, `spectacular prd refine`, `spectacular prd review`.

### `spectacular verify <slug>` *(v1.11.0+)*

Runs the **validation walk** — the skill-side ritual that moves a request from `review` → `verified`. Reads the request's `VERIFY.md` (or falls back to `PLAN § Validation`) and walks every check, one at a time, verifying each by its **kind**:

| Kind | Tag | Authority |
|---|---|---|
| executable | `` `run: <cmd>` `` | command exit code (deterministic) |
| assertable | `{assert}` | agent checks a binary property of files/state |
| judgable | `{judge}` | LLM reasons over named artifacts |
| observable | `{observable}` (default) | human looks & confirms (passive) |
| manual | `{manual}` | human performs an action, then confirms (active) |

Checks can be tagged **inline** per line, or grouped under a `## Title {kind}` section (section is absolute). Executable checks confirm before running (batch-allow at walk start).

The walk records to both `VERIFY.md` (ticks passed boxes) and an append-only `VERIFY-LOG.md` (timestamped, per-check evidence + the `[kind]` that confirmed it). On all-pass it proposes `verified` (configurable auto via `verify.auto_promote`); any blocker keeps the request at `review` with a punch list.

```text
spectacular verify add-team-billing
```

**Skill only** — needs an LLM to read each check and judge evidence. At the terminal the CLI prints a redirect to run `/spectacular verify <slug>` inside Claude Code or Codex. See [verify.md](../skills/spectacular/references/verify.md) — the single verification reference (Part 1 the walk · Part 2 the 2-of-6 rule · Part 3 promoting checks to scripts).

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
