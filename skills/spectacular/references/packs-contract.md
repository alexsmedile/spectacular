# Packs Contract — repo-shape conventions as installable mini-skills

Loaded when the orchestrator needs to read, validate, or apply a convention pack. A pack is the canonical place a project's *opinions* about folder structure, naming, file placement, README contract, and gitignore defaults live.

## Core principle

**A pack is a folder, not a config file.**

It bundles:
1. A `pack.md` manifest (frontmatter contract — see below)
2. `templates/` — files the pack scaffolds into a target project (`.gitignore`, `README.md`, etc.)
3. `references/` — narrative docs the skill may surface to the user (rationale, naming examples)
4. (optional) `scripts/` — pack-local helpers (rare; v1 doesn't ship any)

This mirrors the existing skill convention deliberately. A pack *is* a mini-skill — same shape, different purpose. Users who know skills already know pack folder layout.

## Why packs (vs kits or inline config)

| | Kits | Packs |
|---|---|---|
| Scope | One doc (PRD) | One project shape (whole repo) |
| Output | Slot additions to a single file | Folder structure + multiple template files |
| Bundling | Single `.md` (frontmatter only) | Folder (`pack.md` + `templates/` + `references/`) |
| Distribution | Bundled with skill | Bundled + user-installed + project-local + app-store |
| Composition (v1) | single-kit per PRD | single-pack per repo |

Kits decorate one document. Packs shape the whole repo. Both share the same registry pattern and override layering.

## Pack folder shape

```
<pack-name>/
├── pack.md                    # required — manifest + rule declarations
├── templates/                 # optional — files the pack scaffolds
│   ├── .gitignore             # gitignore stub
│   ├── README.md              # README contract stub
│   └── <other>                # any file the pack wants to scaffold by default
├── references/                # optional — explanatory docs the skill may load
│   └── why-<name>.md          # rationale / examples (loaded on demand)
└── scripts/                   # optional, rare — pack-local helpers
```

The only required file is `pack.md`. Everything else is opt-in. A schema-only pack with no templates is valid (it declares rules but ships no files).

## `pack.md` manifest

```yaml
---
pack: <pack-id>                # required, kebab-case, used in config.yaml + registry
version: 1.0                   # required, semver
description: |                 # required, one-paragraph; surfaced in pack-selection menu
  <what this pack opinionates about and who it's for>
extends: <parent-pack-id>      # optional, single-pack inheritance (v2 — declare-only in v1)
applies-to:                    # optional, project-type filter; "any" if omitted
  - cli
  - library
  - skill
rules:                         # required — the 6 rule categories below
  naming: { ... }
  taxonomy: { ... }
  root-files: { ... }
  gitignore: { ... }
  file-placement: { ... }
  project-types: { ... }
templates:                     # optional — list of files in templates/ this pack will scaffold
  - .gitignore
  - README.md
references:                    # optional — list of files in references/ surfaced to user
  - why-<name>.md
---

# <Pack Name>

<Optional longer-form body: rationale, examples, anti-patterns. Loaded only when the
skill explicitly surfaces this pack's docs to the user — not parsed.>
```

## Field semantics

**`pack`** (required) — kebab-case identifier. Matches the folder name. Used in `config.yaml`'s `convention_pack: { source: <pack-id> }` and in `spectacular pack install/remove <pack-id>`.

**`version`** (required) — semver. Bumped on every meaningful schema or rule change. Pack consumers (init / new-request / doctor) may warn on major-version mismatches between pack and skill.

**`description`** (required) — one paragraph, surfaced in the pack-selection menu (`spectacular pack list` and during interactive init). Lead with what the pack opinionates about, then who benefits.

**`extends`** (optional, v1: declare-only) — name of a parent pack this one inherits from. v1 parses the field but does not yet resolve inheritance. v2 will resolve via shallow merge (child overrides parent per rule category).

**`applies-to`** (optional, list) — project types this pack is designed for. `[any]` if omitted. Init may filter the pack-selection menu based on detected project type. Valid values match the `project-types` schema below.

**`rules`** (required) — the 6 rule categories. Each is its own block; see § "Rule categories" below. A pack may leave a category empty (`{}`) to inherit from parent or skip entirely.

**`templates`** (optional, list) — explicit list of files in `templates/` the pack will scaffold. The list is the contract; files in `templates/` not listed here are ignored. Lets a pack ship optional templates the user opts into individually (rare in v1).

**`references`** (optional, list) — explicit list of files in `references/` the skill may load. Same opt-in contract as `templates`.

## Rule categories (the 6)

Each rule category is a small structured block. The engine reads only the fields it understands; unknown fields are preserved but ignored. This keeps the schema forward-compatible.

### 1. `naming`

Folder/file naming conventions.

```yaml
naming:
  folder-case: kebab-case          # kebab-case | snake_case | PascalCase | any
  file-case: kebab-case            # same options; can differ from folders
  pattern: "{anchor}-{descriptor}[-{role}]"   # optional, descriptive only
  max-words: 3                     # optional, soft cap
  role-suffixes:                   # optional, fixed allow-list
    - ctrl
    - manager
    - svc
  forbidden-words:                 # optional, words to flag
    - v2
    - new
    - project
  forbidden-prefixes:              # optional
    - app-
    - svc-
  date-formats:                    # optional, restricted scopes for dates-in-names
    sandbox: "-YYYYMM"
    archive: "-YYYY-MM"
```

**Doctor enforcement (when mode=enforce):** scan folder + file names in repo root; warning per violation; error on forbidden words.

### 2. `taxonomy`

Required and opt-in folder structure at project root.

```yaml
taxonomy:
  required:                        # folders that must exist
    - src/
    - tests/
  opt-in:                          # folders the pack knows about; scaffold on demand
    - docs/
    - examples/
    - scripts/
  mono-collection-detect:          # heuristic for mono-collection root detection
    when-children-have:            # if N+ children contain these signals, treat parent as mono
      - .git/
      - package.json
      - pyproject.toml
    threshold: 2
  mono-collection-folders:         # different taxonomy applied to mono-collection roots
    - apps/
    - libs/
    - tools/
    - design/
    - sandbox/
    - archive/
```

**Doctor enforcement:** missing required folder → error; presence of unrecognized top-level folder → info (not warning — packs aren't exhaustive).

### 3. `root-files`

Required and optional files at project root, plus the README contract.

```yaml
root-files:
  required:
    - README.md
    - .gitignore
  conditional:                     # required-when conditions
    - file: AGENTS.md
      when: agentic
    - file: LICENSE
      when: oss
    - file: CHANGELOG.md
      when: post-v1
  optional:
    - CLAUDE.md
    - STACK.md
  readme-contract:                 # checked structure for README.md
    must-contain-header:           # first ~10 lines must include these labelled lines
      - Type
      - Stack
      - Run
    must-contain-sections:
      - What it does
      - Setup
      - Usage
```

**Doctor enforcement:** missing required → error; missing readme-contract header → warning.

### 4. `gitignore`

Defaults to include + tool-generated dirs to leave alone unless explicitly opted in.

```yaml
gitignore:
  always-add:                      # appended unconditionally
    - _archive/
    - _archived/
    - _backup/
    - _backups/
    - _tmp/
    - scratch/
    - .env.local
    - .env.*.local
    - .spectacular.local/
  opt-in:                          # offered but not auto-added
    - .scrapekit/
    - .playwright-mcp/
    - .smart-env/
    - .cache/
  never-auto-add:                  # the pack explicitly will not auto-gitignore these
    - .obsidian/                   # user must opt in (e.g. via vault projects)
  language-specific:               # added when stack detected
    python:
      - __pycache__/
      - "*.pyc"
      - .venv/
    node:
      - node_modules/
      - dist/
```

**Init behavior:** `always-add` lines appended to project's `.gitignore` during init (deduplicated). `opt-in` surfaced interactively. `language-specific` block keyed by detected stack.

**Doctor enforcement:** any `always-add` line missing from `.gitignore` → warning; mechanical fix available.

### 5. `file-placement`

Where new files of each kind land. The skill consults this when scaffolding artifacts.

```yaml
file-placement:
  helper-script: scripts/<name>.sh
  architecture-doc: docs/<name>.md
  skill-reference: references/<name>.md
  research-artifact: _research/<topic>/
  backup: _backups/<timestamp>/
  generated-cache: .cache/<name>
  sensitive-data: .env.local
  temp-work: scratch/<name>
  request-artifacts: .spectacular/requests/<slug>/artifacts/{kind}/
  large-file-threshold: 5MB        # files larger than this prompt user before commit
```

**Init / new-request behavior:** consulted whenever a new artifact is created and no explicit target is given.

**Doctor enforcement:** files matching a `kind` pattern found in the wrong location → info (not warning — placement is advisory, not strict).

### 6. `project-types`

Type-specific scaffold definitions. Each type names a folder under the pack's `templates/repo/<type>/` (if the pack ships repo scaffolds) and declares its add-on folders.

```yaml
project-types:
  cli:
    adds:
      - cli/
      - scripts/
      - install.sh
    template-dir: templates/repo/cli/
  library:
    adds:
      - src/
      - tests/
      - examples/
      - docs/
    template-dir: templates/repo/library/
  skill:
    adds:
      - SKILL.md
      - references/
      - templates/
      - scripts/
    template-dir: templates/repo/skill/
  plugin:
    adds:
      - .claude-plugin/
      - skills/
      - agents/
      - commands/
    template-dir: templates/repo/plugin/
```

**Init behavior:** `spectacular init --type <type>` consults active pack's `project-types.<type>` block and scaffolds the `adds:` folders + copies the `template-dir/` contents.

**Doctor enforcement:** none (project type is a one-shot init decision, not an ongoing invariant).

## Override / inheritance layering

A pack can be loaded from one of four scope locations. When the same pack-id exists in multiple scopes, **project-local wins, then user, then app-store, then bundled** — same precedence the engine uses for kit overrides.

| Scope | Path | Notes |
|---|---|---|
| bundled | `skills/spectacular/templates/packs/<name>/` | Ships with the skill. v1 includes only `minimal`. |
| app-store | `<repo-root>/packs/<name>/` | Distributable via this repo. Cloned/forked to install. |
| user | `~/.spectacular/packs/<name>/` | Installed via `spectacular pack install <name>` (v2 / request 3). |
| project-local | `<project>/.spectacular/packs/<name>/` | Per-project override; never auto-shared. |

Project-local lookup → user lookup → app-store lookup → bundled fallback. The first hit wins entirely; v1 does **not** merge across scopes (avoids the same conflict surface that drove single-kit-only).

`extends:` field in `pack.md` is reserved for v2 cross-pack inheritance (parsed-only in v1).

## What packs CAN do

- Declare rules in any subset of the 6 categories
- Ship template files (`.gitignore`, `README.md`, etc.) that init scaffolds
- Ship reference docs the skill surfaces to the user
- Be loaded from any of the four scope locations
- Be enforced (when `config.yaml` mode=enforce) by doctor

## What packs CANNOT do

- Run arbitrary code at install time (no `pack.json` with scripts; pack is data + markdown)
- Modify files outside the conventions they declare (no `post-install` hooks)
- Apply two packs simultaneously to one repo (v1: single-pack only — see Composition)
- Override the always-set (`PRD.md`, `requests/`, `current/`, `config.yaml`, `<agents-file>`) — packs only *add*

## Composition (v1: single-pack only)

Each repo declares exactly one pack in `config.yaml`. Multi-pack composition is out of scope for v1 for the same reasons that drove single-kit-only:

1. **Rule conflicts** — two packs could disagree on `naming.folder-case` (kebab vs snake)
2. **Taxonomy conflicts** — two packs could both require different `src/` layouts
3. **Gitignore unions** — easy to merge, but `never-auto-add` lists conflict
4. **Doctor enforcement noise** — two conflicting enforce-mode packs would double-flag

**v2 multi-pack sketch:** `convention_packs: [base, alex-default]` with declaration-order precedence (later wins for scalar rules, union for lists, explicit `never-auto-add` always wins).

## Schema coverage check

Walking the 10 conventions from `archive/repo-conventions/PLAN.md` to confirm the schema can encode them all:

| # | Archived convention | Encoded in | Status |
|---|---|---|---|
| 1 | Naming (kebab-case, role suffixes, forbidden words) | `naming.*` | ✓ full |
| 2 | Top-level taxonomy (mono-collection: `apps/`, `libs/`, etc.) | `taxonomy.mono-collection-folders` + `mono-collection-detect` | ✓ full |
| 3 | Per-project standard folders (`src/`, `scripts/`, `tests/`, `docs/`, `_research/`, etc.) | `taxonomy.required` + `taxonomy.opt-in` | ✓ full |
| 4 | Per-project root files (README, AGENTS, LICENSE, CHANGELOG, .gitignore, STACK) | `root-files.required` + `conditional` + `optional` | ✓ full |
| 5 | README contract (Type/Stack/Run header + What/Setup/Usage sections) | `root-files.readme-contract` | ✓ full |
| 6 | AGENTS.md pattern (most-specific wins) | `root-files.conditional[when: agentic]` + narrative in `references/` | ✓ structure only; "most-specific wins" lives in pack docs, not schema |
| 7 | .gitignore defaults (always / opt-in / never-auto-add / language-specific) | `gitignore.*` | ✓ full |
| 8 | File placement rules (helper script → `scripts/`, research → `_research/`, etc.) | `file-placement.*` | ✓ full |
| 9 | Project-type templates (cli, library, webapp, skill, plugin, content, research, vault-project) | `project-types.<type>.adds` + `template-dir` | ✓ full (template content ships in pack's `templates/repo/<type>/`) |
| 10 | Mono-collection vs project root detection | `taxonomy.mono-collection-detect` | ✓ full |

**Result: 10 / 10 expressible.** Conventions 1, 2, 3, 4, 5, 7, 8, 9, 10 fully schema-encoded. Convention 6's "most-specific AGENTS.md wins" rule is a runtime behavior, not a pack rule — it lives in `references/why-<pack>.md` narrative.

## Engine integration

### Pack loading (consumed by init / new-request / doctor — requests 2 & 3)

1. Read `config.yaml` for `convention_pack: { source: <pack-id>, mode: <mode> }`
2. Resolve `<pack-id>` via scope precedence (project-local → user → app-store → bundled)
3. Parse `pack.md` frontmatter; validate against this schema
4. Make rule blocks available to the calling surface

### Pack validation (consumed by `spectacular pack review` — request 2)

1. `pack.md` exists and has required frontmatter fields
2. `rules` block contains at least one of the 6 categories (empty pack is invalid)
3. Listed `templates:` files exist in `templates/`
4. Listed `references:` files exist in `references/`
5. `applies-to` values are recognized project types (warn on unknown)
6. `version` is semver-parseable

## Adding a new pack

1. Pick a pack-id (e.g. `alex-default`, `python-app`, `vault-project`)
2. Create a folder at one of the four scope locations
3. Write `pack.md` with required frontmatter
4. Ship any `templates/` and `references/` the pack needs
5. Validate with `spectacular pack review` (request 2)
6. Activate per repo via `config.yaml`'s `convention_pack` block (request 3)

No engine changes. No code edits to the skill. Pack discovery is path-based.

## Related

- [[doc-registry]] — declares the `convention-pack` doc type; pack registration entry
- [[kits-contract]] — sibling pattern for single-doc kit extensions
- [[pack-overrides]] — (request 2) pack-specific grill prompts the fabricator will consume
- [[init-workflow]] — (request 3) where pack consumption wires into init
- [[doctor]] — (request 3) where pack enforce-mode lands as a check area
