# Why alex-default?

This pack distills a working set of conventions used across a personal monorepo (`~/code/`-style mono-collection) and an adjacent Obsidian vault. It's opinionated by design — if you'd rather opinionate yourself, install `minimal` instead.

## What this pack stands for

### 1. Kebab-case names everywhere (with language exceptions)

Folders and standalone files use `kebab-case`. The reason isn't aesthetic — it's that `kebab-case` is the one casing convention that survives every shell, every URL, every filesystem, every cross-tool reference without escaping or quoting.

Two principled exceptions:
- **Python package directories** stay `snake_case` because they get imported (`import my_package`)
- **Go directories** stay lowercase because that's the language norm

Inside code, follow the language norm (`camelCase` in JS, `snake_case` in Python). The naming rule is for **directory and file names visible to the shell and to humans browsing the tree**.

### 2. Role suffixes, not prefixes

Names follow `{anchor}-{descriptor}[-{role}]`. The role suffix is from a fixed allow-list: `ctrl`, `manager`, `svc`, `orch`, `worker`, `dash`, `viz`, `exp`.

So: `payments-dash`, not `dash-payments` or `app-payments`. The anchor comes first because *what it is* matters more for browsing than *what kind it is*. The role suffix sorts together (`grep -- '-dash'` finds all dashboards) and reads naturally in prose.

### 3. Mono-collection vs project root detection

`~/code/` is a *mono-collection* — a parent folder full of independent projects. It needs a different shape (`apps/`, `libs/`, `tools/`, `sandbox/`, `archive/`) than any individual project.

The pack detects mono-collection roots by checking whether 2+ children contain `.git/`, `package.json`, `pyproject.toml`, or `SKILL.md`. When true, init treats the parent as a mono-collection and scaffolds taxonomy at that level rather than per-project.

### 4. AGENTS.md, most-specific wins

`AGENTS.md` is the contract for agentic projects. When present, it governs agent behavior. Each subfolder may have its own `AGENTS.md`; the most-specific one wins (the same precedence rule as `.gitignore`, but for agent instructions).

The schema can only declare *that* AGENTS.md is required (when `agentic`). The *most-specific-wins* behavior is runtime, not schema — it lives here in narrative because that's where load-time semantics belong.

### 5. README contract

Every project must have a `README.md` with a Type/Stack/Run header in the first ~10 lines. Three lines, zero excuses. They mean any agent or human can triage a repo in 2 seconds: "is this a Python lib I can pip-install? a CLI? a content folder?"

Three required sections follow: **What it does**, **Setup**, **Usage**. Not exhaustive — just enough that someone can use the thing without reading the source.

### 6. Gitignore baseline — and what's NEVER auto-added

The `always-add` list covers safe defaults: `_archive/`, `_backups/`, `_tmp/`, `scratch/`, `.env.local`, `.spectacular.local/`. Almost every project benefits; almost no project regrets it.

The `never-auto-add` list is more important than the always-add list:
- `.scrapekit/`, `.playwright-mcp/`, `.smart-env/`, `.obsidian/`

These are tool-generated hidden dirs. The user might want them committed (an Obsidian vault often *does* commit `.obsidian/`); the pack must never make that decision unilaterally. Per the global rule: ask before adding.

### 7. File placement is advisory, not strict

Helper scripts go in `scripts/`. Research artifacts go in `_research/`. Backups go in `_backups/<timestamp>/`. These are defaults the skill uses when scaffolding new files — they're not enforced by doctor (info-level, not warning).

The reason placement is advisory: every project has exceptions. A one-file Python script doesn't need `scripts/`. A project that intentionally keeps research at the root is making a legitimate choice. Doctor flags placement drift as `info` so the user can see it without being nagged.

### 8. Project types

The pack ships scaffolds for 8 types: `cli`, `library`, `webapp`, `skill`, `plugin`, `content`, `research`, `vault-project`. Each type knows what folders it adds and what `template-dir:` to copy from.

Type selection happens at init time (`spectacular init --type cli`). The pack consults the active pack's `project-types.<type>` block and lays down the relevant scaffold.

> **Note on template-dir folders:** This pack declares `template-dir:` paths but ships only minimal stubs in `templates/` (`.gitignore` + `README.md`). The per-type `templates/repo/<type>/` folders are stub-only until the [convention-pack-application](../../.spectacular/requests/convention-pack-application/PLAN.md) request wires init to consume them. Users wanting full per-type scaffolds today should either (a) wait for v0.4.0 or (b) ship their own templates via project-local override.

## When NOT to use this pack

- **You're working on a single one-off script.** `minimal` is enough. The role-suffix and project-type machinery is tax, not benefit, on small surface area.
- **You disagree with kebab-case for files.** Don't fight the pack — fork it (`spectacular pack new my-pack`) and pick `snake_case` or whatever you prefer. Packs are designed to be forked.
- **You're on a team with existing conventions.** Use those. A pack is for projects where *you* are the one setting conventions.

## Lineage

This pack is the fabricator dogfood for [convention-pack-fabricator](../../.spectacular/requests/convention-pack-fabricator/PLAN.md). The conventions encoded here come from the archived [repo-conventions PLAN](../../.spectacular/archive/repo-conventions/PLAN.md) — the 10-section convention catalog that motivated the pack system in the first place.
