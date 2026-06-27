---
doc-id: pack
doc-id-aliases: [convention-pack]
mode: grill
location: ~/.spectacular/packs/<name>/pack.md
scope: user
template: templates/packs/minimal/pack.md
slots: [Name & scope, Naming, Taxonomy, Root files & README, Gitignore, File placement, Project types]
snapshot-on-edit: false
summary: "Repo-shape convention pack — naming + taxonomy + gitignore + file-placement rules"
status: active
---

# Pack Rules — pack-specific rules consumed by the skill

Loaded by `grill.md` / `refine.md` / `review.md` when the active doc is `pack` (per doc-index). `convention-pack` is a recognized back-compat alias for the same doc-id.

This file declares everything pack-specific. The skill handles the rest.

## What this doc produces

A `pack.md` manifest at `~/.spectacular/packs/<pack-id>/pack.md` plus the surrounding folder (`templates/`, `references/`). Schema and field semantics live in [[packs-contract]]. This file only declares *how the grill walks the user through filling in* that schema.

## Pre-flight (before slot loop)

1. **Resolve target path** — `~/.spectacular/packs/<pack-id>/` (creates folder if missing). `--scope project` overrides to `<project>/.spectacular/packs/<pack-id>/`.
2. **Source-ingestion mode** — if `--from <path1>,<path2>,...` was passed, read those files and pre-populate confident slots (see § Source ingestion below).
3. **Scaffold base** — copy `templates/packs/minimal/pack.md` to target location. The minimal pack is the baseline; the grill **overrides** sections as the user answers.
4. **Open snapshot** — packs use `snapshot-on-edit: false` (per registry); creating a brand-new pack snapshots nothing.

## Slot prompts

The skill uses these in the slot loop. One question per slot, short, with concrete guidance. **Each slot maps to one rule category** from [[packs-contract]] (except slot 1 which is metadata).

**Slot 1 — Name & scope**
> What's this pack called and what kind of projects is it for?
>
> Provide: pack-id (kebab-case), one-sentence description, and the project types it applies to (`cli`, `library`, `skill`, `plugin`, `content`, `research`, `vault-project`, or `any`).
>
> *Example:* "pack-id: `alex-default`. Description: 'Opinionated defaults for solo-dev mono-collections — kebab-case naming, mono-collection detection, AGENTS.md pattern, project-type scaffolds.' Applies to: any."

**Slot 2 — Naming**
> What naming rules does this pack enforce on folders and files?
>
> Cover: folder-case (kebab/snake/Pascal/any), file-case (often same), forbidden words (`v2`, `new`, `project`), role suffixes if any (`ctrl`, `svc`, `manager`), forbidden prefixes (`app-`, `svc-`), date-in-name scopes (usually only `sandbox/` and `archive/`).
>
> Leave empty (`{}`) if this pack doesn't opinionate on naming.
>
> *Example:* "Folders: kebab-case. Files: kebab-case. Max 3 words. Role suffixes: ctrl, manager, svc, orch, worker, dash, viz, exp. Forbid: v2, new, project. Forbid prefixes: app-, svc-."

**Slot 3 — Taxonomy**
> What top-level folder structure does this pack require or suggest?
>
> Cover: `required:` folders (must exist), `opt-in:` folders (pack knows about, scaffolds on demand), and — if relevant — mono-collection detection rules + the mono-collection-specific folders (`apps/`, `libs/`, `tools/`, etc.).
>
> Leave empty if this pack is shape-agnostic.
>
> *Example:* "Required at project root: none (varies by project type). Opt-in: src/, scripts/, tests/, docs/, examples/, _research/. Mono-collection detect: when 2+ children have .git/ OR package.json OR pyproject.toml. Mono folders: apps/, libs/, tools/, design/, dash/, sandbox/, templates/, archive/, infra/."

**Slot 4 — Root files & README contract**
> What files must (or should) exist at project root, and what does this pack's README contract look like?
>
> Cover: `required:` files (README, .gitignore typically), `conditional:` (LICENSE when oss, CHANGELOG when post-v1, AGENTS.md when agentic), `optional:` files, and the README contract (header lines + section headings the README must contain).
>
> *Example required:* "README.md, .gitignore."
> *Example conditional:* "AGENTS.md when agentic, LICENSE when oss, CHANGELOG.md when post-v1."
> *Example README contract:* "Header must include `Type:`, `Stack:`, `Run:` in first 10 lines. Sections required: What it does, Setup, Usage."

**Slot 5 — Gitignore**
> What does this pack add to a project's `.gitignore` by default?
>
> Cover three lists: `always-add:` (appended unconditionally), `opt-in:` (offered to user interactively), `never-auto-add:` (pack will refuse to auto-ignore even on request — typically tool-generated hidden dirs the user must explicitly opt into). Optionally `language-specific:` blocks keyed by detected stack (`python:`, `node:`, `go:`).
>
> *Example always-add:* "_archive/, _archived/, _backup/, _backups/, _tmp/, scratch/, .env.local, .env.*.local, .spectacular.local/."
> *Example never-auto-add:* ".scrapekit/, .playwright-mcp/, .smart-env/, .obsidian/ — user must opt in."

**Slot 6 — File placement**
> Where do new files of each kind land when scaffolded?
>
> Cover the common kinds: helper script, architecture doc, skill reference, research artifact, backup, generated cache, sensitive data, temp work, request artifacts, large file threshold.
>
> Leave empty if this pack doesn't opinionate on file placement.
>
> *Example:* "helper-script: scripts/<name>.sh. architecture-doc: docs/<name>.md. research-artifact: _research/<topic>/. backup: _backups/<timestamp>/. generated-cache: .cache/<name>. sensitive-data: .env.local. temp-work: scratch/<name>. request-artifacts: .spectacular/requests/<slug>/artifacts/{kind}/. large-file-threshold: 5MB."

**Slot 7 — Project types**
> Which project types does this pack scaffold, and what does each add?
>
> For each type, declare: `adds:` (list of folders + root files this type creates) and `template-dir:` (where in this pack's `templates/repo/<type>/` folder the type's scaffold files live, if any).
>
> Leave empty (`{}`) if this pack doesn't ship project-type scaffolds.
>
> *Example:* "cli: adds [cli/, scripts/, install.sh, README.md, LICENSE]; template-dir: templates/repo/cli/. library: adds [src/, tests/, examples/, docs/]; template-dir: templates/repo/library/. skill: adds [SKILL.md, references/, templates/, scripts/]; template-dir: templates/repo/skill/."

## Source ingestion (`--from`)

If `--from <path1>,<path2>,...` is passed at grill start, the skill reads each file and uses heuristics to pre-populate slot answers. The user reviews and confirms (default `y`, edit on `n`).

| Source file pattern | Slots pre-populated |
|---|---|
| `.gitignore` (any path) | Slot 5: parse line-by-line; everything goes to `always-add` initially; user reclassifies during review |
| `*NAMING*` / `*naming*` (markdown) | Slot 2: extract bulleted rules; map `kebab-case` mentions, role-suffix lists, forbidden-word lists |
| `README.md` (any path) | Slot 4: detect header lines (`**Type**:`, `**Stack**:`, `**Run**:`), detect H2 sections |
| `CLAUDE.md` / `AGENTS.md` | Slot 4 (agentic flag) + Slot 6 (file placement preferences mentioned in prose) |
| Existing project folder tree | Slot 3 taxonomy: ls top-level folders, present as opt-in; let user mark required |

**Confidence rule:** only pre-populate a slot when the source unambiguously declares the answer. A `.gitignore` containing `_archive/` unambiguously sets `always-add: [_archive/]`. A README mentioning "we usually use kebab-case" is ambiguous → ask the user during the slot instead of pre-populating.

**Always show the user the source-derived answer before locking it in.** The grill prints `Pre-filled from <path>: <answer>. Keep? [Y/n]`. On `n`, falls back to the normal slot prompt.

## Mini-refine patterns

Applied inline by the grill after each slot answer.

| Pattern | Slots scope | Trigger | Proposed action |
|---|---|---|---|
| Reserved pack-id | 1 only | User picks `blank`, `minimal`, `default`, or names matching `pack`/`packs` | "That name is reserved or ambiguous. Pick a more specific id (e.g. `<role>-default`, `<lang>-app`)." |
| Conflicting naming rules | 2 only | `folder-case: snake_case` + `forbidden-words: [snake]` (or similar contradiction) | "These rules contradict — folder-case requires snake but forbidden-words rejects snake-style names. Pick one." |
| Required folder with no use case | 3 only | `required:` includes a folder type that doesn't appear in any `project-types.<type>.adds` | "Folder `<name>` is required but no project type uses it. Demote to `opt-in:` or document why it's universally required." |
| README contract too thin | 4 only | `readme-contract.must-contain-header` AND `must-contain-sections` both empty | "Without a contract, `root-files.required: [README.md]` only checks that the file exists. Add at least Type/Stack/Run header expectations." |
| Empty gitignore always-add | 5 only | `always-add` is empty AND `never-auto-add` is empty | "Empty gitignore block means this pack does nothing to a project's .gitignore. Either add safe defaults or remove the gitignore rule entirely." |
| Tool-generated in always-add | 5 only | `always-add` contains `.scrapekit/`, `.playwright-mcp/`, `.smart-env/`, or `.obsidian/` | "Tool-generated hidden dirs must live in `never-auto-add`, not `always-add` — per the global rule, the user opts in explicitly." |
| Unknown project type | 7 only | `project-types.<type>` uses a type name outside `{cli, library, webapp, skill, plugin, content, research, vault-project}` | "Unknown project type `<type>`. Stick to the documented set or add it to `applies-to:` first so users know it's supported." |
| Project-types without template-dir | 7 only | A type declares `adds:` but no `template-dir:` AND `adds:` includes files (not just folders) | "Type `<type>` adds files but has no `template-dir:` — where do the file contents come from? Add a template-dir or restrict `adds:` to folders only." |

## Mini-refine exemptions

**Slot 1 (Name & scope) generic vague-word scan is skipped** — pack descriptions are necessarily abstract ("opinionated defaults", "minimal baseline"); the reserved-pack-id check is the only slot-1 enforcement.

**Slots 2-7 each have a custom mini-refine pattern.** The PRD vague-word list does **not** apply to pack docs — packs encode rules, not feature descriptions.

## Vibe → spec rewrite tables (refine mode)

### Vague naming rules → concrete rules

| Vibe | Spec |
|---|---|
| "use kebab-case everywhere" | `folder-case: kebab-case`, `file-case: kebab-case`, `pattern: "{anchor}-{descriptor}[-{role}]"` |
| "avoid weird names" | `forbidden-words: [v2, new, project, temp, misc, util, utils, common, shared]` (extend per project) |
| "no random suffixes" | `role-suffixes: [<explicit allow-list>]` — if you don't have an allow-list, drop the rule |
| "short names" | `max-words: 3` |

### Vague taxonomy → concrete taxonomy

| Vibe | Spec |
|---|---|
| "standard project layout" | `required: [src/, tests/]`, `opt-in: [docs/, examples/, scripts/]` — pick concrete folders |
| "mono-repo support" | `mono-collection-detect: { when-children-have: [.git/, package.json, pyproject.toml], threshold: 2 }` + `mono-collection-folders: [<list>]` |
| "Python project layout" | `required: [src/, tests/]`, `opt-in: [docs/]` — and ship a `project-types.library` with src-layout |

### Vague root-files → concrete rules

| Vibe | Spec |
|---|---|
| "README is required" | `required: [README.md]` + readme-contract with at least Type/Stack/Run header |
| "OSS projects need a license" | `conditional: [{ file: LICENSE, when: oss }]` (and document the `oss` flag source) |
| "AGENTS pattern" | `conditional: [{ file: AGENTS.md, when: agentic }]` + narrative in `references/why-<pack>.md` explaining most-specific-wins |

### Vague gitignore → concrete blocks

| Vibe | Spec |
|---|---|
| "ignore the usual stuff" | `always-add: [_archive/, _backups/, _tmp/, scratch/, .env.local, .env.*.local]` |
| "ignore tool dirs" | move to `opt-in:` — never `always-add` (per global rule) |
| "Python venv" | `language-specific.python: [__pycache__/, "*.pyc", .venv/, .pytest_cache/]` |

### Vague placement → concrete rules

| Vibe | Spec |
|---|---|
| "scripts go in scripts/" | `helper-script: scripts/<name>.sh` |
| "no large files committed" | `large-file-threshold: 5MB` — prompt before commit |
| "research artifacts somewhere" | `research-artifact: _research/<topic>/` |

### Vague project types → concrete types

| Vibe | Spec |
|---|---|
| "CLI projects" | `cli: { adds: [cli/, scripts/, install.sh, README.md], template-dir: templates/repo/cli/ }` |
| "library projects" | `library: { adds: [src/, tests/, examples/, docs/, README.md, LICENSE], template-dir: templates/repo/library/ }` |
| "Claude skill" | `skill: { adds: [SKILL.md, references/, templates/, scripts/], template-dir: templates/repo/skill/ }` |
| "plugin" | `plugin: { adds: [.claude-plugin/, skills/, agents/, commands/], template-dir: templates/repo/plugin/ }` |

## Review gate checks (in addition to base)

| # | Check | How |
|---|---|---|
| 4 | At least one rule category populated | At least one of the 6 `rules.*` blocks is non-empty (`{}` and absent both fail). A pack with no opinions has no reason to exist. |
| 5 | Frontmatter contract satisfied | `pack`, `version`, `description` all present and non-empty. `applies-to:` is a list (default `[any]` if omitted is allowed). |
| 6 | Declared templates exist | Every entry in `templates:` corresponds to an actual file in `templates/`. |
| 7 | Declared references exist | Every entry in `references:` corresponds to an actual file in `references/`. |
| 8 | applies-to values are valid | All `applies-to` entries are in `{any, cli, library, webapp, skill, plugin, content, research, vault-project}`. Unknown values → warning, not error. |
| 9 | Version is semver | `version:` parses as semver (`MAJOR.MINOR[.PATCH]`). |
| 10 | Gitignore never-auto-add not overridden | No `always-add` entry appears in `never-auto-add` (would be self-contradictory). |
| 11 | Naming rule self-consistency | `forbidden-words` and `forbidden-prefixes` do not contradict `folder-case`/`file-case` choices. |
| 12 | README contract minimum | If `root-files.required` includes `README.md` AND `readme-contract:` is declared, contract must have at least one `must-contain-header` OR one `must-contain-section`. |

### Universal base checks still run

Placeholder check (no `<TODO>`, `<TBD>`, `???` in any slot), clarification check (no `[NEEDS CLARIFICATION: ...]` markers remaining), frontmatter integrity (frontmatter parseable as YAML). Apply regardless of pack rules.

## Tokenization rule

Same as PRD: vague-word matching and rule-content scans must preserve hyphenated compounds as single tokens (`[a-z]+(?:-[a-z]+)*`). Pack rule lists are full of hyphenated identifiers (`kebab-case`, `snake_case`, `.env.local`) — naive `\b\w+\b` tokenization breaks them.

## Reserved pack-ids

These names are reserved and the grill rejects them in slot 1:

- `blank` — would collide with the kit concept
- `minimal` — bundled default; user packs may not shadow
- `default` — ambiguous (which default?)
- `pack` / `packs` — too generic; collides with the doc-type name

User packs must pick a distinct id. Recommended: `<owner>-default` (e.g. `alex-default`), `<lang>-app` (e.g. `python-app`), `<role>-<scope>` (e.g. `mobile-team`).

## Related

- [[packs-contract]] — schema spec this grill produces against
- [[doc-index]] — registry entry pointing here as `rules:`
- [[grill]] — consumes the slot prompts + mini-refine patterns from this file
- [[refine]] — consumes the vibe→spec tables for full refine passes
- [[review]] — consumes the gate checks from this file
- [[prd-rules]] — sibling rules file (reference for shape)
