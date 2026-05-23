# Doc Registry — the single source of truth for what Spectacular knows how to write

Loaded when the orchestrator needs to resolve any `spectacular <doc> <verb>` command. The registry maps each known document type to its template, slot list, mode, location, and override rules.

**Adding a new doc type to Spectacular = adding an entry to this file + creating a template.** No other code changes required.

## Core principle

**One registry, many docs. Generic engine, specific overrides.**

The engine (`grill.md` / `refine.md` / `review.md`) is doc-agnostic. It reads this registry to know:
- Which template to scaffold from
- Which slots to walk
- How to behave (`grill` interactive, `append` single entry, `freeform` no enforcement)
- Where the file lives
- Which override rules apply (per-doc gate logic, vague-word lists, etc.)

When no override exists, the engine runs base behavior.

## Schema

Each entry is YAML inside a fenced block. Required fields are marked **required**.

```yaml
<doc-id>:                              # required, kebab-case, used in triggers
  template: <path>                     # required, relative to skills/spectacular/
  mode: grill | append | freeform      # required
  location: <path>                     # required, where the file is created
  scope: project-wide | per-request | user  # required
  slots: [Slot1, Slot2, ...]           # required if mode=grill, omit otherwise
  snapshot-on-edit: true | false       # default false
  overrides: <path>                    # optional, path to per-doc override doc
  kit-support: true | false            # default false; whether kits can mutate this doc's slots
  description: <one-line>              # required, surfaced in `spectacular <doc> --help`
```

### Field semantics

**`template`** — markdown file the engine copies as the initial scaffold. Placeholders in the template (`<DATE>`, `<PROJECT NAME>`, slot markers) are filled by the engine using context.

**`mode`** — drives engine behavior:
- `grill` — walk slots in order, ask one question at a time, run inline mini-refine, exit on review gate pass
- `append` — append a single new entry to the existing file (e.g. ADR entries in DECISIONS.md); no slot loop
- `freeform` — scaffold the template and exit; no interactive slot-fill (used for docs where structure is loose: ROADMAP, AGENTS)

**`location`** — supports `<slug>` interpolation for per-request docs:
- `.spectacular/PRD.md` (project-wide)
- `.spectacular/requests/<slug>/PLAN.md` (per-request — slug supplied by `spectacular new`)

**`scope`** — `project-wide` files exist once at `.spectacular/` root; `per-request` files exist once per request folder; `user` files live under `$HOME` (e.g. `~/.spectacular/packs/<name>/`) and are shared across all projects on the host.

**`slots`** — ordered list of slot names. The engine walks these in order. Slot names appear as `## N. <Name>` section headings in the template.

**`snapshot-on-edit`** — when `true`, the engine snapshots the existing file (`<DOC>@vN.md`) before applying any edit. Always `true` for canonical root docs; usually `false` for per-request files.

**`overrides`** — path to a markdown file with per-doc rules the engine consumes:
- Vague-word lists specific to this doc
- Custom gate checks (e.g. "Goals slot must contain number + verb + date")
- Mini-refine patterns unique to this doc type
- See `prd-overrides.md` for the reference shape

**`kit-support`** — `true` only for PRD in v1. Kits (see [[kits-as-plugins]]) declare additions to this doc's slots and trigger other docs to scaffold.

## The registry (v1)

```yaml
# ─── Project-wide canonical root docs ──────────────────────────────────────

prd:
  template: templates/prd/base.md
  mode: grill
  location: .spectacular/PRD.md
  scope: project-wide
  slots: [Vision, Problem, Target users, Deliverable, Goals & success criteria, Non-goals, Constraints, First milestone]
  snapshot-on-edit: true
  overrides: references/prd-overrides.md
  kit-support: true
  description: "Product Requirements Document — what & why & for whom (8 slots)"

principles:
  template: templates/principles/base.md
  mode: freeform
  location: .spectacular/PRINCIPLES.md
  scope: project-wide
  snapshot-on-edit: true
  description: "Operating principles + runtime enforcement hooks"

architecture:
  template: templates/architecture/base.md
  mode: freeform
  location: .spectacular/ARCHITECTURE.md
  scope: project-wide
  snapshot-on-edit: true
  description: ".spectacular/ structure, frontmatter, lifecycle, versioning"

spec:
  template: templates/spec/base.md
  mode: freeform
  location: .spectacular/SPEC.md
  scope: project-wide
  snapshot-on-edit: true
  description: "System spec — index of what the system actually is and how it behaves right now (capability list, link out to specs/<capability>/SPEC.md only when needed)"

roadmap:
  template: templates/roadmap/base.md
  mode: freeform
  location: .spectacular/ROADMAP.md
  scope: project-wide
  snapshot-on-edit: false
  description: "Time-ordered what's next"

stack:
  template: templates/stack/base.md
  mode: freeform
  location: .spectacular/STACK.md
  scope: project-wide
  snapshot-on-edit: true
  description: "Host project's tech choices"

agents:
  template: templates/agents/base.md
  mode: freeform
  location: .spectacular/AGENTS.md
  scope: project-wide
  snapshot-on-edit: true
  description: "Onboarding doc for agents working in .spectacular/"

decisions:
  template: templates/decisions/entry.md
  mode: append
  location: .spectacular/DECISIONS.md
  scope: project-wide
  snapshot-on-edit: false
  description: "ADR-style decision log — append one entry per decision"

# ─── Per-request docs ──────────────────────────────────────────────────────

plan:
  template: templates/plan/base.md
  mode: grill
  location: .spectacular/requests/<slug>/PLAN.md
  scope: per-request
  slots: [Goal, Constraints, Milestones, Tasks, Dependencies, Validation, Deliverables]
  snapshot-on-edit: false
  overrides: references/plan-overrides.md
  description: "Request-scoped plan — 7-slot decomposition (owns lifecycle state)"

tasks:
  template: templates/tasks/base.md
  mode: freeform
  location: .spectacular/requests/<slug>/TASKS.md
  scope: per-request
  snapshot-on-edit: false
  overrides: references/tasks-overrides.md
  description: "Executable checklist for one request"

# ─── Convention packs (folder-shape mini-skills) ──────────────────────────

convention-pack:
  template: templates/packs/minimal/pack.md
  mode: grill
  location: ~/.spectacular/packs/<name>/pack.md
  scope: user
  slots: [Name & scope, Naming, Taxonomy, Root files & README, Gitignore, File placement, Project types]
  snapshot-on-edit: false
  overrides: references/pack-overrides.md
  description: "Repo-shape convention pack — naming + taxonomy + gitignore + file-placement rules"

# ─── Public-facing docs (v0.6.0+) ─────────────────────────────────────────

docs-manifest:
  template: templates/docs/docs.yaml.tmpl
  mode: freeform
  location: docs/docs.yaml
  scope: project-wide
  snapshot-on-edit: false
  overrides: references/docs-overrides.md
  description: "Nav manifest for docs/ — sections + page order + site metadata"

docs-page:
  template: templates/docs/page.md.tmpl
  mode: freeform
  location: docs/<section>/<slug>.md
  scope: per-request
  snapshot-on-edit: false
  overrides: references/docs-overrides.md
  description: "Single user-facing docs page — driven by `spectacular docs new <page>`"
```

> **`docs-manifest`** is project-wide (one per repo). `docs-page` is conceptually per-page but uses `per-request` scope since the engine treats `<slug>` interpolation identically. The skill's `docs new` flow supplies both `<section>` and `<slug>`.

> **No `audience` in docs-page frontmatter** — see [[docs-contract]] § Core principle. Audience is folder-level (`docs/` vs `specs/`), never per-page.

> **`scope: user`** — packs live under `$HOME` by default, not per-project. Per-project override at `<project>/.spectacular/packs/<name>/` is allowed (precedence rules: see [[packs-contract]]). This differs from the rest of the registry, which uses `project-wide` / `per-request`.

> **`overrides: references/pack-overrides.md`** — file does not exist in v1 (lands in request 2: convention-pack-fabricator). The registry entry is forward-declared so request 2 has a target.

## How the engine uses this

For any `spectacular <doc> <verb>` invocation:

1. Look up `<doc>` in the registry. Missing → "unknown doc type" error.
2. If `<verb>` is omitted → infer: `grill` if location is empty, `review` if filled (only valid for `mode: grill`); `append` for `mode: append`; open in editor for `freeform`.
3. Load `template`, `slots` (if grill), `location`, `snapshot-on-edit`, `overrides`.
4. Run engine verb (`grill.md` / `refine.md` / `review.md`) against this config.
5. If `overrides` is set, the engine loads it and merges per-doc rules into base behavior.

## Override file shape

An override file is doc-specific. It contains:
- Mini-refine patterns unique to this doc (e.g. PRD's "plural-user → singular" rule)
- Custom gate checks (e.g. PLAN's "milestones before tasks" ordering)
- Vague-word lists scoped to specific slots
- Examples of good/bad answers per slot

The engine reads the override file *if present* and applies its rules in addition to base behavior. No override = base behavior only.

See `prd-overrides.md` as the reference example.

## Adding a new doc type

1. Pick a kebab-case ID (e.g. `verify`)
2. Create the template at `templates/<id>/base.md`
3. Add a registry entry here with required fields
4. (Optional) Create `references/<id>-overrides.md` if the doc has unique rules
5. (Optional) Add the triggers `spectacular <id>` / `spectacular <id> grill` / etc. to SKILL.md routing — usually unnecessary since the generic handler catches them

No engine code changes. No new skills. No plugin manifest.

## Project-local overrides (v2)

Out of scope for v1, planned for v2: `.spectacular/doc-registry.yaml` in the host project would override or extend the bundled registry. Use cases: custom slot schemas, project-specific doc types, alternate templates. Until then, kit overrides (see [[kits-as-plugins]]) cover the common case for PRD customization.

## Related

- [[grill]] — generic interactive slot-filler that consumes registry
- [[refine]] — generic vibe→spec rewriter
- [[review]] — generic quality gate runner
- [[prd-overrides]] — reference override file (PRD-specific rules)
- [[plan-overrides]] — milestone-before-tasks ordering check
- [[tasks-overrides]] — checklist format check
- [[scaffold-reference]] — file-template reference (separate concern: what the templates look like)
