# Changelog

All notable changes are documented here.
Format: [Keep a Changelog](https://keepachangelog.com/en/1.0.0/)

---

## [1.0.1] — 2026-05-23

### Added
- `spectacular --version` / `-v` / `version` prints the CLI version.
- `spectacular` / `spectacular help` / `spectacular --help` / `-h` (no subcommand) prints a top-level usage with the full verb list. `spectacular init --help` still prints the init-specific flag reference.

### Fixed
- `.codex-plugin/plugin.json` version drifted to `0.6.0` during the v1.0.0 cut; realigned to match the Claude plugin manifest.

---

## [1.0.0] — 2026-05-23 — first stable release

**Spectacular reaches v1.** The surface developed across v0.6.0 → v0.7.5 is now frozen as the stable contract. Future changes follow strict semver: breaking changes require a major bump; new capabilities require a minor; fixes ship as patch releases.

This release is a **tag, not new code** — it captures the dev arc's endpoint. All capabilities below are inherited from the v0.x lineage; per-version detail lives in the entries that follow.

### v1 surface (what's stable)

**Workspace foundation**
- `.spectacular/` directory contract: PRD, SPEC, config, AGENTS, requests/, specs/ as always-set
- Frontmatter as the signal layer — skill reads frontmatter, never full file bodies, for briefings
- 5-state lifecycle (planned → active → review → verified → archived) stored in PLAN.md frontmatter
- Convention packs (4-tier scope: project-local → user → app-store → bundled; 3 modes: suggest / scaffold / enforce; 6 rule categories)
- Snapshot-before-edit on canonical docs (`<FILE>@vN.md`)
- Memory (`spectacular remember this`) writes to `.spectacular/memory/`, git-committed
- Workspace schema versioning (`workspace_schema:` field) + migration registry + `spectacular migrate` verb

**Doc-writing engine**
- One generic grill / refine / review engine + registry-driven doc types
- 13 registered docs: prd, spec, plan, tasks, principles, architecture, roadmap, stack, agents, decisions, convention-pack, docs-manifest, docs-page
- Kit system: PRD has 5 kits (blank, coding, content, product, research) with extension contract
- Structured ROADMAP with 3 precision tiers (full / themed / vision), 9-phase chain, Icebox section
- 2-of-6 verification rule (when VERIFY.md is required vs. folded into PLAN/TASKS)

**CLI surface** (`spectacular`)
- `init` (kit + flags + interactive mode + --update for skill refresh)
- `doctor` (10 areas: skill, workspace, frontmatter, snapshots, links, lifecycle, kits, conventions, specs, docs; --fix for mechanical repairs; precondition for destructive verbs)
- `pack` (list, install, show, remove, with 4-tier scope resolution)
- `docs` (init, doctor passes)
- `migrate` (--dry-run, --to, --from, --list; registry-driven)
- **Mutator verbs**: `new`, `promote`, `snapshot`, `archive`, `touch` — CLI mutates atomically; skill orchestrates
- `status --against-latest` (workspace-schema verdict)
- `status --since <date>` (activity report; YYYY-MM-DD / 7d / 2w / 1m / yesterday)

**Skill surface** (`/spectacular`)
- Lean orchestrator (~120 lines) + on-demand reference loading
- Routing for status, new, archive, doctor, migrate, doc verbs, memory, snapshot, lifecycle transitions
- Skill walks judgment-requiring flows; CLI handles mechanical mutations

**Distribution**
- Claude Code plugin marketplace (`alexsmedile/spectacular`)
- Codex plugin marketplace
- CLI via curl one-liner installer

### Testing

- 7 test files, 216 asserts
- `init.test.sh`: 66 asserts
- `doctor.test.sh`: 48 asserts
- `migrate.test.sh`: 39 asserts
- `mutator.test.sh`: 59 asserts
- `pack.test.sh`: 48 asserts
- `conventions.test.sh`: 18 asserts
- `specs.test.sh`: 29 asserts

### Semver discipline starts here

- Patch (1.0.x): bug fixes, doc clarifications, internal refactors that preserve all public contracts
- Minor (1.x.0): new CLI verbs, new doctor areas, new registered doc types, new kits, new convention-pack rule categories
- Major (2.0.0): any change that breaks an existing workspace's ability to load, validate, or migrate

The `workspace_schema:` field tracks structural shape independently of CLI version. A v2.0.0 CLI with `workspace_schema: "1.0"` workspaces will run the migration chain on first invocation.

### Migration from v0.x

Nothing to do. Anyone on any v0.6.0+ CLI keeps working — v1.0.0 is the same surface with a fresh version label. If you're on a pre-v0.6.0 workspace, run `spectacular migrate` to bring your workspace schema forward.

### Pre-1.0 development arc

The history that produced v1.0.0 is preserved below as per-version entries. Each entry documents the capability that shipped in that version and is the place to look for "when did X land?" questions. Future releases (v1.0.1+, v1.1.0+) will continue this entry pattern.

---

## [0.7.5] — 2026-05-23

### Added — `spectacular status --since <date>` activity report

Closes task #57. Lists requests + canonical doc changes since a cutoff date.

- **Frontmatter-only scan** — reads `updated:` fields from `requests/*/PLAN.md`, `archive/*/PLAN.md`, and `.spectacular/*.md` (skipping snapshot files). Lexicographic YYYY-MM-DD compare against cutoff. No git involvement (matches Spectacular's "frontmatter is the signal layer" principle).
- **Date input formats**:
  - `--since 2026-05-20` — absolute YYYY-MM-DD
  - `--since 7d` / `--since 2w` / `--since 1m` — relative (days / weeks / months-approx-30d)
  - `--since yesterday` — keyword
  - `--since=<date>` — equals form also accepted
- **BSD + GNU date compatibility** — uses macOS `date -v-Nd` when available, falls back to GNU `date -d "N days ago"`.
- **Output groups requests by status** — archived / verified / review / active / planned, in lifecycle order. Each group only renders if non-empty. Empty bucket renders `(none)`.
- **Canonical docs section** lists `.spectacular/*.md` files (excluding snapshot @v* files) with their `updated:` date.

Exit codes: `0` success (including empty buckets); `1` outside a workspace; `2` missing or unparseable `--since` argument.

### Implementation note

Intercepts `--since` in the same `status` dispatch path as `--against-latest` (introduced v0.6.1). Both are CLI-mechanical flags on the otherwise-skill-owned `status` verb. Plain `spectacular status` (no flags) still routes to the skill stub for AI-driven briefing.

### Testing

- 7 test files, all green
- `init.test.sh` +1 scenario (status --since: 7 asserts covering absolute/relative/equals/empty/missing-arg/bad-format/outside-workspace) — 66 asserts (was 59)
- Plugin bumped 0.7.4 → 0.7.5

### Closes

- Task #57 — `spectacular status --since <date>`

---

## [0.7.4] — 2026-05-23

### Added — Doctor precondition for destructive verbs + VERIFY-as-tests convention

Two safety/process improvements that close out long-standing pending tasks (#55, #56).

**Doctor precondition on `archive` + `migrate` + `promote --archive`:**
- Before running destructive ops, CLI runs a scoped doctor check. Errors block (exit 1); warnings/info pass through silently.
- `archive` runs `workspace + specs + links` — full structural validation. Refuses if any errors found.
- `migrate` runs `links` only — workspace/specs would self-trigger on the v0.4-shape drift migrate is designed to repair. Dry-run skips the check entirely.
- `promote --archive` chain passes `--skip-doctor` through to the chained `cmd_archive`.
- Bypass: `--skip-doctor` flag on each verb. Caller responsibility to know what they're skipping.
- Shared helper `_doctor_precondition` for consistent behavior; thin wrappers `require_doctor_clean` (3-area) and `require_doctor_clean_links` (links-only).

**VERIFY-as-tests convention (documentation-only):**
- New `skills/spectacular/references/verify-tests.md` documents the `tests/verify/<slug>.test.sh` pattern: when to author one (multi-verb workflows not covered by `tests/cli/` suites), template, naming convention.
- No backfill: existing 7 `tests/cli/` suites already cover most of what archived VERIFY.md files would script. Pattern reserved for new requests that ship behavior not already in an area-level test suite.
- New VERIFY-tagging convention: `[x] mechanically verified` (auto-checkable) vs `[x] manually verified` (interactive UX, judgment calls). Lets future-agents grep to know which scenarios have a safety net.
- Wired into SKILL.md routing.

### Testing

- 7 test files, all green
- `mutator.test.sh` +2 scenarios (doctor precondition + clean passes through) — 59 asserts (was 49)
- Plugin bumped 0.7.3 → 0.7.4

### Closes

- Task #55 — Promote VERIFY scenarios → tests/verify/<slug>.sh (convention-doc only per scope cut)
- Task #56 — Doctor-as-precondition for destructive ops (archive + migrate + promote --archive)

---

## [0.7.3] — 2026-05-23

### Fixed — Unknown `convention_pack.mode` values fall back to `suggest` with info note

Previously, unrecognized `mode:` values in `config.yaml` (e.g. `strict`) silently passed through and were treated as scaffold-equivalent (warnings, not errors). Behavior diverged from the spec, which says: validate against `{suggest, scaffold, enforce}`, fall back to `suggest` if invalid, emit info note.

- New helper `config_pack_mode_raw()` returns the raw config value verbatim
- `config_pack_mode()` now validates and falls back to `suggest` for unknown values
- `check_conventions` emits info line when the raw value doesn't match the allow-list: "unknown convention_pack.mode 'strict' (valid: suggest / scaffold / enforce) — falling back to 'suggest'"
- Valid modes (`suggest`, `scaffold`, `enforce`) and absent `mode:` field both stay silent (latter defaults to `suggest`)

Found during S18 verification in `convention-pack-application` (task #62).

### Testing

- 7 test files, all green
- `pack.test.sh` +1 scenario (unknown mode fallback) — 48 asserts (was 44)
- Plugin bumped 0.7.2 → 0.7.3

---

## [0.7.2] — 2026-05-23

### Added — Roadmap-research findings applied (Outcome slot, Icebox, gate guardrails)

Tightens the v0.7.1 structured ROADMAP based on convergent advice from Pichler (GO Product Roadmap), Torres (Opportunity Solution Trees), Cagan (SVPG), Gilad (GIST), Productside, and beginner-tool onboarding patterns (GitHub Projects, Trello, Aha!). Closes the gap between v0.7.1 and what dominant roadmap frameworks recommend without ballooning the slot set.

**Structural changes:**
- **`Outcome:` slot added** between Phase and Scope-in. Required for `full` + `themed` tiers; absent for `vision` (Direction covers). One paragraph: "what business or product outcome does this version move?" Pichler/Torres/Cagan/Gilad convergent: the #1 missing slot in feature-list-shaped roadmaps. Forces goal-before-features discipline.
- **"Bucket list" → "Icebox"** rename across template, overrides doc, and live ROADMAP. Convergent dev-tool idiom (GitHub Projects, Pivotal Tracker, Linear; GIST's "Idea Bank"). Distinguishes "unbound idea" from "planned but vague" (which is what `vision`-tier blocks are for).

**Review gate additions** (all warnings/info — preserves recommend-not-enforce stance):
- **Date guards extended to themed/vision blocks** (gate check 12) — Cagan's "#1 sin." Scans entire block for `YYYY-MM-DD`, `Q[1-4] YYYY`, `MMM YYYY`. Warning, not error.
- **Outcome required by tier** (gate check 16) — full + themed must have Outcome; vision must not (Direction covers).
- **Full-tier row count** (gate check 17) — tiered: silent ≤7 (sweet spot per Cagan), info 8-10 ("consider demoting older versions to themed"), warning 11+ ("roadmap-as-backlog anti-pattern").
- **Scope-out push** (gate check 18) — when Scope-in has ≥4 items and Scope-out is empty, warning ("every item you add implies others you're not building" — Productside).

**Phase taxonomy extension:**
- **Meta-phase aliases** — `Phase:` field accepts both individual values (`mvp`, `release-prep`) AND coarser meta-phase values (`discover`, `build`, `release`). Coexist. Document the "start coarse, refine as work crystallizes" rule. Maps Cagan's discovery-vs-delivery split onto our 9-phase chain.

**Documentation-only additions** (no automation):
- **Beginner pattern** in `roadmap-overrides.md` — start at vision tier (one paragraph), graduate to themed when 2nd version exists, unlock full when first request links via `target_version:`. Mirrors GitHub Projects/Trello/Notion progressive-disclosure onboarding.
- **Icebox-promotion ritual** — explicit 4-step walk (pick item → choose version → choose tier → fill slots → delete from Icebox). Skill executes on `/spectacular roadmap` invocation. No new CLI verb; manual ritual is the point.

**Doctor extension:**
- `check_workspace` flags pre-v0.7.2 "Bucket list" heading in ROADMAP, suggests Icebox rename. Mechanical fix tag (sed substitution). Silent once renamed.

**Dogfood:** live `.spectacular/ROADMAP.md` snapshotted (`ROADMAP@v2.md`) then updated with Outcome paragraphs on v0.7.1 + v0.7.x + v0.11.x + v1.0.0, renamed Bucket list → Icebox, and added the 4 deferred research items to the Icebox.

### Out of scope (deferred per interview)

- Confidence rating per row (GIST/ProductBoard) — overlaps tier
- Audience field (Pichler internal-vs-external) — over-engineered for solo/small-team
- Opportunity-Solution-Tree as separate doc type (Torres) — heavyweight
- ICE/RICE scoring for icebox items (GIST signature) — convention-pack territory
- CLI verb for icebox promotion — manual via skill is enough
- `--beginner` flag — doc-only enough

### Testing

- 7 test files, all green
- doctor.test.sh: +1 scenario (Bucket-list → Icebox info) — 48 asserts (was 47)
- Plugin bumped 0.7.1 → 0.7.2

---

## [0.7.1] — 2026-05-23

### Added — Structured ROADMAP with precision tiers + 9-phase chain

ROADMAP graduates from `mode: freeform` to `mode: structured`. Each version block has a precision tier (`full | themed | vision`) that controls which slots are required, capturing the natural precision gradient — active work is detailed; long-term direction is intentionally fuzzy.

- **9-phase chain** for version progression: `intent → discover → prototype → spec-refine → mvp → iterate → test → release-prep → release`. Skill recommends the next phase; user can skip with reason. Skips recorded explicitly in `Phase:` field (e.g. `Phase: spec-refine (skipped: discover, prototype)`).
- **Prototype phase broadened** — any artifact produced to validate a decision against real tooling or against the user counts: data/schema drafts run through parsers, fake datasets tested against downstream scripts, mock API responses, ASCII wireframes, interactive mocks, sample CLI output. The artifact isn't the deliverable; the decision it informs is.
- **Three precision tiers**:
  - `full` — Status, Phase, Scope (in), Scope (out), Exit criteria, Linked requests. Use for active + near-term planned versions.
  - `themed` — Status, Phase, Themes (list), Exit criteria (directional). Use for mid-term (2-3 versions out).
  - `vision` — Status, Direction (free-text paragraph). Use for long-term + speculative.
- **Bucket list section** at the end of `ROADMAP.md` for ideas not yet tied to any version. Promoting an item: pick a version label + `Tier: vision` minimum.
- **`spectacular new --target-version <ver>`** flag — adds `target_version:` to PLAN.md frontmatter. Used by `spectacular roadmap refine` to autopopulate Linked requests in matching version blocks.
- **Doctor extension** — `check_workspace` flags pre-v0.7.1 freeform ROADMAP shape with info line (never error). Skill walk to migrate available via `spectacular roadmap grill`.

### Reference + template

- `skills/spectacular/references/roadmap-overrides.md` — tier-aware slot prompts (full/themed/vision variants), mini-refine patterns, vibe→spec rewrite tables, 15-check review gate
- `skills/spectacular/templates/roadmap/base.md` — structured template showing all 3 tiers + bucket list
- `skills/spectacular/references/doc-registry.md` — ROADMAP entry switched `mode: freeform` → `mode: structured`; mode reference table extended with `structured` and `reference` modes

### Dogfood

`.spectacular/ROADMAP.md` rewritten against new shape (snapshotted as `ROADMAP@v1.md`). Current versions tiered as: v0.7.1 full, v0.7.x/v0.11.x/v1.0.0 themed, v2.x/v3+ vision. Bucket list populated with previously-roadmapped ideas that no longer fit a specific version (hook automation, multi-agent orchestration, burndown viz, etc.).

### Testing

- 7 test files, all green
- mutator.test.sh: +1 scenario (target_version field) — 49 asserts (was 48)
- doctor.test.sh: +1 scenario (ROADMAP shape detection: old triggers info, new is silent) — 47 asserts (was 45)
- Plugin bumped 0.7.0 → 0.7.1

---

## [0.7.0] — 2026-05-23

### Added — CLI mutator verbs (skill orchestrates, CLI mutates)

Five CLI verbs replace skill-side manual file edits for the most common lifecycle mutations. Establishes the **mutation principle**: lifecycle changes go through CLI verbs; the skill orchestrates (reads, decides, communicates); the CLI mutates (atomically, deterministically, tested). Manual file edits remain available for edge cases but become the exception.

- **`spectacular new <slug> [--summary ...] [--status planned|active|review] [--priority low|medium|high]`** — scaffolds `.spectacular/requests/<slug>/PLAN.md` + `TASKS.md` from templates with frontmatter prefilled. Validates slug (kebab-case, max 64 chars). Refuses duplicates in `requests/` or `archive/`.
- **`spectacular promote <slug> [--to <state>] [--force] [--archive]`** — advances request through the lifecycle (planned → active → review → verified). Refuses backward transitions without `--force`. Mutates PLAN.md + TASKS.md `status:` + `updated:` atomically. `--archive` chains into archive after promoting to verified.
- **`spectacular snapshot <file> [--major]`** — snapshots a canonical doc to `<base>@v<N>.md`, bumps `version:` field, sets `updated:`. Refuses non-canonical files. Idempotent: compares body (frontmatter-stripped) against latest snapshot; exits cleanly if unchanged.
- **`spectacular archive <slug> [--force]`** — moves request to `.spectacular/archive/<slug>/`, rewrites every inbound `related:` link in other request files (`../<slug>/...` → `../../archive/<slug>/...`), sets PLAN frontmatter `status: archived` + `archived: <today>`. Refuses unless status is `verified` or `review`.
- **`spectacular touch <file>`** — sets frontmatter `updated:` to today. Idempotent. Refuses files without frontmatter.

### Shared infrastructure

- New frontmatter helpers (`fm_get`, `fm_set`, `fm_touch`, `fm_add_to_list`) — single shared implementation used by all 5 verbs. Single-source-of-truth for YAML rewriting.
- `is_canonical_doc` helper — gates `snapshot` to registered canonical files only.

### Skill instruction sync

To resist drift, both surfaces updated:

- `SKILL.md` routing table: each verb points to its CLI verb (not a reference doc — the CLI is the runtime). Top-level **mutation principle** stated explicitly.
- Reference docs rewritten to instruct CLI verb usage:
  - `references/new-request.md` → `spectacular new <slug>`
  - `references/archive.md` → `spectacular archive <slug>`
  - `references/lifecycle.md` → `spectacular promote <slug>` (state transitions)
  - `references/versioning.md` → `spectacular snapshot <file>`

### Absorbs

- Task #54 (archive --check / auto-rewrite related: paths) — same behavior shipped as the standard `archive` verb behavior.

### Testing

- 7 test files, all green; mutator.test.sh adds 48 asserts across 7 scenarios covering all 5 verbs + the promote-archive combo + --help flags
- init.test.sh scenario_10 updated — only status + remember remain as skill stubs (archive/new/snapshot/promote no longer stubs)

### Why this matters

Before v0.7.0, lifecycle mutations happened via the skill writing free-form file edits. Two agents would do the same thing differently; frontmatter parsers drifted; archive link-rewriting was easy to forget. The CLI verbs collapse this: one implementation, deterministic output, tested. Skill instruction sync ensures the skill actually calls the verbs instead of falling back to manual edits.

Remaining as skill flows (intentionally): `status` (briefing requires AI judgment), `remember` (memory entries require user-context distillation).

---

## [0.6.2] — 2026-05-23

### Added — Workspace migrations: registry pattern + judgment skill walk

Stage 2 of workspace-migrations (Stage 1 shipped in v0.6.1). Replaces the v0.6.1 hardcoded migration list with a proper registry + adds judgment-migration support via the skill.

- **Migration registry** — migrations now live as .md files under `skills/spectacular/references/migrations/v<from>-to-v<to>.md`. Frontmatter declares `(id, from, to, mechanical, reversible, apply-fn, affects)`; body documents detection rule, steps, rollback, validation. Full schema: `skills/spectacular/references/migrations-contract.md`.
- **CLI loader** — `cmd_migrate` now scans the registry, sorts by `from` semver, resolves `apply-fn` to bash functions in `cli/spectacular`. Adding a new migration = one .md file + one bash function (mechanical) OR one skill-walk section (judgment).
- **`spectacular migrate --to <ver>`** — migrate up to a specific schema version (default: latest).
- **`spectacular migrate --from <ver>`** — re-run starting from a specific schema (for repair).
- **`spectacular migrate --list`** — show all registered migrations with descriptions.
- **`/spectacular migrate` skill walk** — `skills/spectacular/references/migrate.md` defines the flow for judgment migrations: snapshot-before-edit on affected canonical docs, y/n/q per step, validation phase, audit-trail memory write. No judgment migrations ship in v0.6.2; the skill flow is scaffolding for future migrations that need it.
- **Chain validation in doctor** — `check_kits` validates the migration registry: no gaps in the chain, no duplicate `(from, to)` edges, every `apply-fn` resolves to a defined bash function, reversible migrations have reverse functions. Reports as part of the `kits` area; no new doctor area needed.
- **Downgrade refused** — `migrate --to <older-version>` exits non-zero. Bidirectional migration is explicitly out of scope.

### Why this matters

In v0.6.1 the migration logic was hardcoded directly in `cmd_migrate` with two if-branches. Easy to ship but doesn't scale: every new migration would have meant editing `cmd_migrate` plus the bash function. The registry pattern collapses this — maintainers add a .md spec + a bash function; the loader handles dispatch + ordering + chain validation automatically.

### Migration

For consumers, no action needed. v0.6.2 reads the same `workspace_schema:` field that v0.6.1 wrote. The behavior of `spectacular migrate` (no flags) is identical.

### Testing

- 6 test files, all green
- migrate.test.sh: 39 asserts (+12 for `--list`, `--to`, downgrade refusal)
- doctor.test.sh: 45 asserts (+3 for chain validation pass + gap detection)
- specs.test.sh: 29 asserts (unchanged)
- init.test.sh: 63 asserts (unchanged)
- conventions.test.sh: 18 asserts (unchanged)
- pack.test.sh: 44 asserts (unchanged)

### What's still in v0.6.1 territory (no changes here)

`workspace_schema:` field, `spectacular status --against-latest`, doctor flat-contract-docs support, v0.6+ scaffold suggestion — all from v0.6.1, unchanged. Stage 2 only touches the registry/loader/skill-walk surfaces.

---

## [0.6.1] — 2026-05-23

### Added — Workspace migrations + scaffold discoverability

Bridge the gap between Spectacular versions for already-scaffolded workspaces. Triggered by a real audit on the [Octopus](https://github.com/alexsmedile/octopus) repo where a v0.1.x-shape workspace had no discoverability path to v0.6 conventions.

- **`spectacular migrate`** — applies pending workspace-schema migrations to bring `.spectacular/` to the shape this CLI expects. Idempotent.
  - `--dry-run` lists planned migrations without writing
  - Two backfilled migrations: **v0.4 → v0.5** (rename `current/` → `specs/`, preserve contents) + **v0.5 → v0.6** (ensure `specs/` exists as always-set)
- **`workspace_schema:` field** — new top-level key in `config.yaml`. Records the structural version. `init` writes `"0.6"` on fresh workspaces; absent value treated as `"0.4"`.
- **`spectacular status --against-latest`** — one-line discoverability check. CLI-mechanical (no skill); the rest of `status` stays in the skill.
- **Flat contract docs in `specs/`** — top-level `.md` files in `specs/` are valid alongside per-capability subfolders. Pattern for projects whose primary truth is on-disk contracts (e.g. `SCHEMA-TASK.md`, `AXIS-MODEL.md`).
- **v0.6+ scaffold suggestion** — `doctor workspace` surfaces missing PRINCIPLES/ARCH/ROADMAP as one info line with the exact `init --with` command. Silent when all present.

Stage 2 (v0.6.2) will move the migration list from hardcoded into a `references/migrations/` registry, add `--to`/`--from` flags, and bring `/spectacular migrate` skill walk for judgment migrations.

---

## [0.6.0] — 2026-05-23

### Added — Public docs as a first-class surface

The `docs/` tree is now a first-class Spectacular surface, sibling to `.spectacular/` (workspace) and `specs/` (system truth). Flat, opinionated, renderer-agnostic — single `docs.yaml` nav manifest, page-level frontmatter, CLI scaffold + doctor validation + skill verbs for interactive authoring.

**Audience clarification:** the spec/doc boundary lives at the *folder* level, not the page level. `docs/` is for users + agents consuming the product; `specs/` is for devs + coding agents building it. Per-page `audience` would be ceremony — no such field.

- **`spectacular docs init [--minimal]`** — CLI subcommand. Scaffolds `docs/docs.yaml` + `index.md` + 3 default sections (`getting-started`, `guides`, `reference`) with placeholder pages. `--minimal` skips the sections and ships just docs.yaml + index. Idempotent; re-running fills empty stubs without overwriting content.
- **`spectacular doctor docs`** — substrate validation: docs.yaml parseable, declared pages exist, no orphan files, required frontmatter present (`title`, `description`, `section`, `status`, `updated`). Supports both sectioned trees (`docs/<section>/<page>.md`) and flat-tree extras (`docs/<slug>.md` registered via `docs.yaml extras:`).
- **`doctor docs --fix`** — mechanical: injects frontmatter stubs into pages missing them (delimiter or individual fields). Title defaults to slug, status to `draft`, updated to today.
- **Skill verbs** (registry-driven via the existing engine):
  - `spectacular docs new <page>` — scaffolds a page, prompts for section if omitted, updates docs.yaml
  - `spectacular docs new --section <name>` — declares a new section + scaffolds the dir
  - `spectacular docs review` — quality gate (same checks as doctor)
  - `spectacular docs status` — briefing scoped to docs/
- **`references/docs-contract.md`** — schema spec: folder shape, docs.yaml manifest, page frontmatter contract, validation rules, anti-patterns. Documents the spec-vs-doc boundary at folder level.
- **`references/docs-overrides.md`** — engine rules: `docs new` flow with section-prompt UX, `docs review` gate checks, `docs status` briefing format. Vibe→spec patterns for future `docs refine` (deferred to v2).
- **Doc registry** — `docs-manifest` and `docs-page` entries registered. Doc IDs registered count rises from 11 to 13.
- **`templates/docs/`** — `docs.yaml.tmpl`, `index.md.tmpl`, `page.md.tmpl` for the engine + CLI.
- **`tests/cli/docs.test.sh`** — 12 scenarios, 38 asserts covering init (default + minimal + idempotent), doctor (skip / clean / missing-declared / orphan / missing-frontmatter / --fix injection / extras), skill-verb refusal by CLI, help text.

### Dogfood

- **This repo's own `docs/`** migrated to v0.6.0 shape: `docs.yaml` authored with three sections + 5 existing pages registered as `extras:` (flat-tree preservation — moving files would break README links). Each of `workflow.md`, `commands.md`, `configuration.md`, `scaffold.md`, `troubleshooting.md` got proper frontmatter (title, description, section, status, since, updated). Doctor `docs` clean — 0 errors / 0 warnings.

### Out of scope (deferred to `public-docs-advanced` v0.6.2+)

- Renderer adapters (Mintlify / Docusaurus / Fumadocs / MkDocs export)
- Versioned docs snapshots (`docs/versioned/v<x.y.z>/`)
- `docs sync-from-spec` (spec ↔ doc sync flow)
- Convention-pack `docs-layout` rule category

These ship only when real-world demand surfaces (per the same activation-trigger pattern as `convention-pack-modules`).

### Quality

- 191 asserts pass across 5 test files (init, doctor, pack, specs, docs)
- Pre-commit version-consistency check green across 7 sources
- Doctor on this repo: 0 errors / 0 warnings on all 10 areas

---

## [0.5.0] — 2026-05-23

### Breaking — `current/` folder renamed to `specs/` + new `SPEC.md` index

The legacy `.spectacular/current/` folder convention is replaced by a two-part surface: `.spectacular/SPEC.md` (always-on, present-tense index of what's built) plus an optional `.spectacular/specs/` folder for per-capability detail when a SPEC.md bullet outgrows one line.

- **Why:** "current" is a temporal word, not a content word — agents kept mis-routing it as recency state. "spec" is the industry term and ties cleanly to what the layer actually holds. The TODO had this flagged since v0.3.0; v0.5.0 ships the migration.
- **`SPEC.md` is always-on**, scaffolded by every init. Per-capability `specs/<capability>/SPEC.md` files are **optional** — only break out when the bullet in SPEC.md outgrows one line. Small projects ship with one file.
- **Mechanical migration via doctor**: `spectacular doctor specs --fix` renames any legacy `current/` → `specs/`, preserving contents. Conflict case (both dirs present) raises an error and refuses auto-fix.

### Added

- **`SPEC.md` doc type** registered in `doc-registry.md` — uses the generic engine, mode `freeform`, snapshot-on-edit, scope `project-wide`.
- **`templates/spec/base.md`** — index-style template with "What this system is" + "Capabilities" sections.
- **`doc_spec()` writer** in `cli/spectacular` — scaffolds SPEC.md as part of the always-set.
- **Doctor `specs` area** — validates SPEC.md presence/parseability, specs/ dir presence, per-capability SPEC.md frontmatter, legacy current/ migration detection, conflict detection.
- **`spectacular spec` skill verb** (via existing registry-driven engine — `grill`, `refine`, `review`).
- **`references/spec-sync.md`** — renamed from `current-sync.md`, updated to drive SPEC.md bullet edits + specs/ creation during archive flow.
- **`tests/cli/specs.test.sh`** — 8 scenarios, 25 asserts covering fresh init, kit init, doctor on clean v0.5.0 workspace, legacy detection, mechanical migration, conflict refusal, per-capability validation, re-init non-destructive.

### Changed

- **Always-set bumped from 5 → 6 files**: PRD.md, **SPEC.md**, config.yaml, `<agents-file>`, requests/, **specs/** (replacing `current/`).
- **`doc_agents()` template rewrite** — new build's `.spectacular/AGENTS.md` now leads with the four-layer model (Intent/Truth/Work/Memory), documents two-layer task tracking, and consistently uses `SPEC.md` + `specs/<capability>/SPEC.md` references throughout. Mirrors landed in `templates/agents/base.md`.
- **Doctor re-run dispatcher fix** — `conventions` area was missing from the `--fix` re-run loop, causing pack-driven gitignore fixes to not refresh detection state. Both `conventions` and the new `specs` area are now in both dispatcher loops.
- All `current/<capability>` references swept through skill references, `.spectacular/` live workspace, `docs/`, `README.md`, `CLAUDE.md` — replaced with `specs/<capability>/SPEC.md` per the new convention. Migration callouts retained where users upgrading from v0.4.x need them.
- Root `AGENTS.md` updated to point at `.spectacular/SPEC.md` for system-truth queries and clarify per-request loading discipline.
- **codex-plugin `longDescription`** — replaced "current truth" framing with "system spec" framing to match the new convention.

### Migration from v0.4.x

```bash
spectacular doctor specs        # detect legacy current/
spectacular doctor specs --fix  # rename current/ → specs/, preserve contents
spectacular init                # fills in SPEC.md if missing
```

Workspaces with both `current/` and `specs/` present require manual merge — doctor refuses auto-fix.

---

## [0.4.0] — 2026-05-23

### Added — Convention Pack system

A new opt-in layer for declaring repo-shape opinions: naming rules, folder taxonomy, required root files, gitignore defaults, file-placement rules, project-type scaffolds. Packs are mini-skills (folder + `pack.md` + `templates/` + `references/`), distributable via four scope locations (project-local → user → app-store → bundled, in precedence order).

- **`packs-contract.md`** — full schema spec covering 6 rule categories (`naming`, `taxonomy`, `root-files`, `gitignore`, `file-placement`, `project-types`), pack folder shape, 4-tier scope precedence, single-pack-only v1 (multi-pack composition in [convention-pack-modules](.spectacular/requests/convention-pack-modules/)).
- **`pack-overrides.md`** — pack-specific grill rules: 7 slot prompts, source-ingestion (`--from`), 8 mini-refine patterns, vibe→spec rewrite tables, review gate checks, reserved pack-id enforcement.
- **`spectacular pack` CLI subcommand** — `list` (shows all 4 scope locations), `install <name> [--from <path>]` (copies pack to `~/.spectacular/packs/<name>/`), `remove <name> [--force]` (refuses bundled/app-store/project-local without `--force`), `show <name>` (prints scope + frontmatter).
- **`config.yaml` `convention_pack:` block** — declares active pack per repo with `source`, `mode` (suggest|scaffold|enforce), and reserved `overrides` field.
- **Init wiring** — when a pack is declared with `mode: scaffold` or `enforce`, init appends the pack's `gitignore.always-add` entries (deduplicated). Always-set wins on conflicts; pack scaffold is purely additive.
- **Doctor `conventions` area** — validates pack source resolves; in `scaffold` mode flags gitignore drift as warnings (exit 1); in `enforce` mode escalates to errors (exit 2); `suggest` mode skips drift checks. `--fix` mechanically appends missing pack-declared gitignore entries.
- **Bundled `minimal` pack** — ships at `skills/spectacular/templates/packs/minimal/`. Enforces only README contract + `.gitignore` baseline. The implicit default when no other pack is declared.
- **App-store `alex-default` pack** — ships at `packs/alex-default/`. Fully-opinionated pack encoding kebab-case naming with role suffixes, mono-collection detection, 8 project-type scaffolds, full `.gitignore` baseline + language-specific blocks for Python/Node/Go.
- **`tests/cli/pack.test.sh`** — 12 scenarios, 44 asserts covering list/install/remove/show, init wiring, doctor across all 3 modes, mechanical fix repair, scope precedence, `--from` install, error paths.

### Added — Workflow conventions

- **Two-layer task tracking convention** documented in `.spectacular/AGENTS.md` and skill's `SKILL.md`: harness `TaskCreate`/`TaskUpdate` = ephemeral session micro-tracker (drives CLI live progress UI); on-disk `requests/<slug>/TASKS.md` = persistent milestone blocks. Anti-pattern: one-for-one duplication. Resolves the recurring "task tools haven't been used" warning by giving it a real role.

### Changed

- `references/init-workflow.md` — § "Convention packs (v0.4.0+)" added with 3-mode behavior table + 4-tier precedence table.
- `references/doctor.md` — `conventions` check area added; severity-per-mode table.
- `references/new-request.md` — `artifacts/` directory consults active pack's `file-placement.request-artifacts:` rule.
- `references/doc-registry.md` — `convention-pack` entry registered with new `scope: user` value (packs live under `$HOME`, not per-project).
- `cli/spectacular` — new `SCRIPT_DIR` constant; `pack` subcommand sibling of `init` + `doctor`; `check_conventions()` doctor area; `pack_apply_scaffold()` init hook; `config_pack_source/mode()` awk parsers.
- `cli/spectacular@v0.3.1` snapshot captured before v0.4.0 work.

### Convention pack chain (3 requests)

- `convention-pack-schema` (verified) — locks the schema + ships bundled `minimal`
- `convention-pack-fabricator` (review) — pack-overrides + alex-default dogfood; live grill walkthroughs remain for full signoff
- `convention-pack-application` (review) — CLI + init + doctor wiring; live three-mode + cross-machine scenarios remain

### Planned

- **`convention-pack-modules`** (planned, v0.5.0+) — split monolithic packs into composable rule-category modules. Stays planned until composition pain surfaces from v1 use.

---

## [0.3.1] — 2026-05-23

### Added

- **`spectacular doctor`** — environment/infrastructure self-check. CLI detects substrate drift across 7 areas (`skill`, `workspace`, `frontmatter`, `snapshots`, `links`, `lifecycle`, `kits`); skill handles judgment-requiring repairs via `/spectacular doctor --fix`. Severity model: ✅ pass / ⚠️ warning / ❌ error / ℹ️ info. Exit codes: 0 clean, 1 warnings, 2 errors. `--format text|json` for human + machine consumption. `--fix` applies content-free mechanical repairs (`.gitignore` append, missing always-set dirs, dangling symlinks, re-stub missing canonical files).
- **Skill-invoked doctor subsets** — `status.md`, `grill.md`, `onboarding.md`, `lifecycle.md` auto-run scoped doctor checks when substrate failures block their operation; findings surface inline.
- **`references/doctor.md`** — full spec: check definitions per area, severity model, report format (text + JSON), CLI vs skill split, repair-flow walkthrough with worked examples, anti-patterns.
- **`tests/cli/doctor.test.sh`** — 11 scenarios, 33 asserts covering detect, mechanical fix, scoped areas, JSON output, regression.

### Changed

- Smart-init's diagnostic placeholder (`"run diagnostics via spectacular doctor once available"`) replaced with explicit area pointer: `"run \`spectacular doctor frontmatter\` for details"`.
- `references/SKILL.md` routing extended with doctor triggers + skill-invoked-subset note.

### Workspace cleanup (from doctor dogfood)

- Fixed broken `related:` paths in `prd-craft/PLAN.md` + `repo-conventions/PLAN.md` (file-relative paths instead of repo-root-relative).
- Parked `cli-bootstrap` request (`active → planned`) since smart-init shipped at v0.3.0.
- Documented `PRD@v1.1.md` snapshot gap in `DECISIONS.md` (version bump skipped during canonical-docs-rework; no git trace exists).

---

## [0.3.0] — 2026-05-22

### Changed (breaking)

- **`spectacular init` scaffolds only what the project needs**, not all 7 root docs. Default = 5-file always-set (`PRD.md`, `config.yaml`, `<agents-file>`, `requests/`, `current/`). Extra docs come from the selected kit's `triggers-docs.always` and `triggers-docs.suggested` lists, or via explicit `--with <doc1,doc2>` flag.

### Added

- **PRD 8-slot shape** — base PRD template + all 5 kits versioned to v1.1 (Vision / Problem / Target users / Deliverable / Goals & success criteria / Non-goals / Constraints / First milestone).
- **Doc-writer engine** — generic `grill.md` / `refine.md` / `review.md` references that consume `doc-registry.md` to handle any doc type. PRD-specific logic moved to `prd-overrides.md`. Per-doc overrides for PLAN + TASKS. New templates for plan/tasks/principles/architecture/roadmap/stack/agents/decisions.
- **Kits as diff-only plugins** — kits now declare `adds-slots`, `modifies-slots`, `triggers-docs.always`, `triggers-docs.suggested` via frontmatter. Single-kit-only in v1. Documented in `references/kits-contract.md`. All 5 bundled kits (`blank`, `coding`, `content`, `product`, `research`) refactored to diff format.
- **Verification convention** — `references/verification.md` formalizes when VERIFY.md is needed (2-of-6 rule) vs folded into PLAN § Validation or TASKS § Verification. "Opt-in" refers to file scaffolding only; verification itself is mandatory before any `verified` transition.
- **CLI flags** — `--kit <name>`, `--with <doc1,doc2>`, `--minimal` flags for `spectacular init`. Backwards-compatible with existing `--name`, `--summary`, `--agents-file`, `--global`, `--update`, `-i` flags.
- **Pre-flight non-overwrite** — init detects existing/empty/malformed files and skips/fills/diagnoses without ever overwriting. Generic "run diagnostics via `spectacular doctor`" message emitted for malformed cases (to be replaced when doctor ships).
- **Test harness** — `tests/run.sh` discovers and runs `tests/**/*.test.sh`. First test suite at `tests/cli/init.test.sh` covers 6 smart-init scenarios (41 asserts).

### Anti-patterns formalized

- Per-doc skills (e.g. one skill per doc type) — superseded by all-in-one `/spectacular` with registry-driven verbs.
- `--force` flag — explicitly rejected. Re-init never overwrites; to regenerate a stub, delete the file first.
- Project-type inference in init — bare init uses `blank` kit unconditionally. Auto-detection deferred to v2.
- Skipping verification because "VERIFY.md is optional" — opt-in refers to the file, not the practice. Every request reaches `verified` through some artifact (VERIFY.md > TASKS § Verification > PLAN § Validation).

### Upgrading from v0.2.x

`spectacular init` on a v0.2.x workspace is safe — pre-flight skips every existing file. To add docs the v0.3.0 init no longer scaffolds by default (PRINCIPLES, ARCHITECTURE, etc.), run `spectacular init --with <docs>` or `spectacular init --kit <kit>`.

CLI snapshot preserved at `cli/spectacular@v0.2.0` for reference.

---

## [0.2.0] — 2026-05-21

### Changed

- **Canonical docs split.** The original 896-line `.spectacular/PRD.md` was split into four focused root docs:
  - `PRD.md` — product intent (now 121 lines, 6-slot shape: problem / who / success / non-goals / constraints / milestone)
  - `PRINCIPLES.md` — 8 operating principles, each with a runtime enforcement hook
  - `ARCHITECTURE.md` — workspace structure, frontmatter conventions, lifecycle, versioning
  - `ROADMAP.md` — versioned future work (v1 / v2 / v3+)
- **AGENTS.md rewritten** as the in-folder onboarding doc for any agent landing in `.spectacular/`. Authoritative source for per-task context loading rules.
- **PLAN.md template upgraded** to the 7-slot decomposition: goal / why / constraints / milestones / tasks / dependencies / validation / deliverables.
- **CLI scaffolds the full 7-doc root layer** on every `spectacular init`. PRD stub uses the new 6-slot shape; AGENTS stub uses the onboarding shape; new PRINCIPLES / ARCHITECTURE / ROADMAP stubs included.
- **Skill references aligned** with the new doc set — `status.md`, `onboarding.md`, `init-workflow.md`, `scaffold-reference.md`, `new-request.md`, `versioning.md`, and SKILL.md state-awareness all updated.
- **Project docs aligned** — README, CLAUDE.md, `docs/scaffold.md`, `docs/configuration.md`, `docs/commands.md`, `docs/troubleshooting.md`, `docs/workflow.md` all reflect the new 7-doc canonical set.

### Added

- `prd / prd refine / prd review` skill triggers — interactive PRD building with 5 kits (coding / product / content / research / blank), vibe→spec refine patterns, and a pass/fail quality gate.
- `requests/prd-craft/` and `requests/canonical-docs-rework/` — tracking artifacts for the v0.2.0 work.
- Snapshot history preserved: `PRD@v1.3.md`, `AGENTS@v1.0.md`, request PLAN/TASKS `@v1.0.md` and `@v1.1.md`.

### Anti-patterns formalized

- **Never create `requests/<slug>/PRD.md`.** Product intent is project-wide and lives at `.spectacular/PRD.md`. Per-request folders use `PLAN.md` + `TASKS.md` only.

---

## [0.1.1] — 2026-05-11

### Fixed

- `cli/spectacular`: progress log in `download_and_install_skill` redirected to stderr — was polluting `skills.lock` with the "Fetching skill..." line
- `cli/spectacular`: skill ref resolution now tries releases API first, tags API second, `main` as final fallback — previously fell straight to `main` when no GitHub release existed

---

## [0.1.0] — 2026-05-11

### Added

- `cli/spectacular` — Bash CLI binary (`spectacular init`) with zero-prompt default, `-i` interactive mode, and flags: `--name`, `--summary`, `--agents-file`, `--global`, `--update`
- `cli/install.sh` — curl-installable installer; places binary at `~/.local/bin/spectacular`
- `skills/spectacular/` — `/spectacular` Claude Code slash command; lean SKILL.md orchestrator routing to `references/` subdocs
- `.spectacular/config.yaml` — `agents.file` key (override primary agents file) and `agents.tool_overrides` map (per-tool supplementary files, e.g. `claude: CLAUDE.md`)
- `skills.lock` — CLI-written lockfile tracking installed skill ref, SHA, and source URL
- Skill install targets both `.agents/skills/spectacular/` (source) and `.claude/skills/spectacular/` (symlink) for multi-tool compatibility
- `CLAUDE.md` — project guidance for Claude Code
