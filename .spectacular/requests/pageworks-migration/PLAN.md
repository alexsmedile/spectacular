---
status: planned
priority: high
owner: alex
updated: 2026-05-23
target_version: spectacular v1.2.0 (deprecation) + pageworks v0.1.0 (first ship)
summary: "Extract public-facing documentation work into a new standalone skill `pageworks`; deprecate docs verbs in spectacular without removing them"
related:
  - ../../archive/public-docs-advanced/PLAN.md
  - ../../archive/public-docs-foundation/PLAN.md
  - ../../../skills/spectacular/references/docs-contract.md
  - ../../../skills/spectacular/references/docs-renderer-adapters.md
  - ../../../skills/spectacular/references/docs-overrides.md
---

# Plan — pageworks migration

## Goal

Spin a new standalone skill `pageworks` at `skills_db/pageworks/` that owns the entire public-facing documentation surface: scaffold, schema, structure, authoring, review, maintenance, and renderer export. Deprecate (but don't yet remove) the equivalent surface in spectacular. Spectacular keeps only **discovery awareness** of docs/ — it knows the folder exists and that a manifest is present, nothing else.

## Why

Documentation authoring is a different cognitive job than operational scaffolding. Spectacular's strength is *internal* — PRDs, specs, plans, decisions, requests, lifecycle. Public docs are *external* — narrative, evergreen, multi-audience, prose-heavy. Bundling both into one skill has accumulated cost: 30+ references in spectacular, docs-specific overrides, renderer adapters, schema contracts. Splitting them lets each skill stay sharp.

Separately: a future `docs-writer` style agent (Tier 1 in pageworks's design) becomes natural once authoring lives in its own skill. Spectacular's role becomes orchestrator-of-skills, not omni-skill.

## Scope

### In scope (one request, multiple deliverables)

**pageworks the skill** (`skills_db/pageworks/`)

- New skill scaffolded with `SKILL.md`, `references/`, `templates/`, `scripts/`
- Lift-and-shift of working code/refs from spectacular:
  - `docs-contract.md` → `pageworks/references/contract.md` (schema becomes pageworks's contract)
  - `docs-overrides.md` → `pageworks/references/authoring.md` (the doc-engine verbs)
  - `docs-renderer-adapters.md` → `pageworks/references/renderers.md`
  - `docs_init` shell code → `pageworks` CLI binary
  - `docs_export` shell code → `pageworks export <renderer>` CLI
  - `check_docs` doctor area → `pageworks doctor` (or equivalent self-check)
  - Templates from `spectacular/templates/docs/` → `pageworks/templates/`
- Page-type templates added (Diátaxis: tutorial / how-to / reference / explanation)
- New `pageworks` CLI binary at `pageworks/cli/pageworks`
- Standalone scaffold: `pageworks init` works on any project with no spectacular present
- Plugin manifests for Claude Code + Codex (mirrors spectacular's pattern)

**Spectacular deprecation (kept working, marked deprecated)**

- `spectacular docs init|export|new|review|status` print a deprecation banner pointing at pageworks
- Deprecation banner is informational only — verbs still run for backward compatibility
- `doctor docs` slims to **discovery-only**: pass/info on whether docs/ exists, whether a manifest is present, whether pageworks is installed. No schema validation, no frontmatter checks, no renderer validation.
- `docs-overrides.md`, `docs-contract.md`, `docs-renderer-adapters.md` stay in spectacular's references for the deprecation cycle; each gains a `> Deprecated — see pageworks` banner at the top
- CHANGELOG entry under spectacular v1.2.0: deprecation notice + removal target (v2.0.0)

**Handoff signaling**

- When spectacular archives a request that touched SPEC.md, it prompts the user: *"This change may affect public docs/. Hand off to pageworks for updates?"* (User decides; no auto-invocation.)
- When `spectacular doctor docs` runs and discovers docs/ + no pageworks installed, it prints the install hint.
- Spectacular's `.spectacular/AGENTS.md` and CLAUDE.md updated with the new boundary.

### Out of scope

- **Subagents in pageworks v0.1.0.** No `docs-writer`, `docs-reviewer`, `docs-architect` agent definitions yet. Main Claude Code agent runs the skill directly. Subagent design deferred until skill is used enough to know what to extract — likely a follow-on `pageworks-agents` request.
- **Removing the deprecated surface from spectacular.** That's the v2.0.0 cut, separate request.
- **Public-docs dogfood** (authoring spectacular's own docs/ with pageworks). Separate request `public-docs-dogfood` after pageworks ships and is shaken out.
- **Bundle restructure of spectacular.** No longer needed — pageworks is standalone, not a sub-skill.
- **Cross-skill linting / contract version negotiation between pageworks and spectacular.** v1: pageworks is the source of truth for its own schema; spectacular's discovery check is loose ("is there a manifest?" not "does the manifest match a specific schema version?"). v2 may add a compatibility table.
- **Migrating existing pageworks-style docs that downstream users have generated with spectacular.** The schema is unchanged in v0.1.0 (lift-and-shift), so existing docs/ trees remain valid. No migration script needed.

## Decisions

- **Standalone, not sub-skill.** Pageworks lives at `skills_db/pageworks/`, installs independently. Spectacular references it by name in its handoff doc. Reverses the earlier "bundle" decision because pageworks-as-standalone-capable conflicts with sub-skill packaging.
- **One big migration request.** Solo work, locked architecture — splitting into 5 requests adds bookkeeping without changing the build order.
- **No subagents in v1.** Ship the skill; observe how it's used; extract subagents in a follow-on once the friction points are real, not theoretical.
- **Pageworks owns docs/ completely.** Scaffold, schema, renderer adapters, doctor — all move. Spectacular keeps discovery-only awareness: "is there a docs/?" and "is there a manifest?" — nothing about contents.
- **Deprecate, don't remove.** Spectacular v1.2.0 marks the docs surface deprecated; verbs keep working. Removal lands in spectacular v2.0.0, separate request.
- **Lift-and-shift, not rewrite.** Working code/refs copy verbatim into pageworks. Cleaner reorganization is a pageworks v0.2.x problem; ship working bits first.
- **Spectacular discovers pageworks, never installs it.** Per global "never install without explicit instruction" rule. If pageworks isn't present, spectacular prints the install command and gets out of the way.
- **Request lives in spectacular's `.spectacular/`.** Spectacular drives the migration since it owns the source code being moved. Once pageworks exists, ongoing pageworks work tracks in `skills_db/pageworks/.spectacular/`.

## Milestones

1. **M1 — Scaffold pageworks**
   - Create `skills_db/pageworks/` with `SKILL.md`, `README.md`, `CHANGELOG.md`
   - Plugin manifests (`.claude-plugin/plugin.json`, `.codex-plugin/plugin.json`)
   - Frontmatter declares standalone status + compatible-with-spectacular range
   - Bootstrap with `spectacular init` so pageworks has its own `.spectacular/` for tracking future work
2. **M2 — Migrate references (lift-and-shift)**
   - Copy spectacular's `docs-contract.md`, `docs-overrides.md`, `docs-renderer-adapters.md` into `pageworks/references/` with renames
   - Update internal cross-links to use pageworks naming
   - Add Diátaxis page-type templates: `pageworks/templates/pages/{tutorial,how-to,reference,explanation}.md.tmpl`
   - Add prose-pattern reference: `pageworks/references/prose-patterns.md` (callouts, code blocks, link conventions, voice/tone — new content, not migrated)
3. **M3 — pageworks CLI binary**
   - New binary at `pageworks/cli/pageworks` (Bash, mirrors spectacular's pattern)
   - `pageworks init` — scaffolds `docs/`, `docs.yaml`, `index.md`, default sections (lifted from spectacular's `docs_init`)
   - `pageworks export <renderer>` — renderer adapters (lifted from spectacular's `docs_export_*`)
   - `pageworks doctor` — schema + frontmatter + structure validation (lifted from `check_docs`)
   - `pageworks --version` / `--help` / top-level usage
   - Install script `pageworks/cli/install.sh`
4. **M4 — Spectacular deprecation banners**
   - Each of `docs init|export|new|review|status` prints a one-line deprecation notice pointing at `pageworks <verb>` and continues
   - `docs-contract.md`, `docs-overrides.md`, `docs-renderer-adapters.md` gain `> **Deprecated in v1.2.0** — see pageworks/references/` at the top
   - `doctor docs` slims to discovery (pass/info only — folder presence, manifest presence, pageworks install status)
   - `spectacular docs --help` updated to reflect deprecation status
5. **M5 — Handoff wiring**
   - New reference: `skills/spectacular/references/pageworks-handoff.md` — when to delegate, what the install command looks like, how spectacular detects "docs may need updating"
   - `spectacular archive` checks: if archived request touched SPEC.md or specs/, prompt user about updating docs
   - SKILL.md (spectacular) gets a small section on "Public docs work is owned by pageworks — when delegating, point users at `apm install pageworks` (or the user's preferred installer)"
   - `.spectacular/AGENTS.md` updated with the boundary
6. **M6 — Tests + release**
   - `pageworks/tests/cli/init.test.sh`, `export.test.sh`, `doctor.test.sh` — port equivalents from spectacular
   - `spectacular/tests/cli/docs-deprecation.test.sh` — new: verify each docs verb prints deprecation banner + still works
   - `spectacular/tests/cli/doctor.test.sh` — update docs scenarios to reflect discovery-only
   - Tag spectacular v1.2.0 (deprecation release)
   - Tag pageworks v0.1.0 (first release)
   - GitHub releases for both
   - CHANGELOG entries for both

## Validation

- **pageworks standalone**: in a fresh dir with no spectacular, `pageworks init` creates docs/ tree; `pageworks export mkdocs` works; `pageworks doctor` reports correctly.
- **pageworks alongside spectacular**: both installed, both work, no collision.
- **spectacular deprecation**: each docs verb prints deprecation banner once and still produces the same output. Existing downstream workflows that depend on `spectacular docs export mkdocs` keep working.
- **spectacular doctor docs**: when docs/ absent → silent. When docs/ present + pageworks not installed → info-level "consider installing pageworks." When docs/ present + pageworks installed → pass.
- **Handoff prompt**: after archiving a request that touched SPEC, spectacular prompts (doesn't auto-invoke).
- **No regression**: spectacular's full test suite green; pageworks's full test suite green; both releases pass `bash -n` syntax check.

## Risks

- **Double maintenance during deprecation cycle.** Bug fixes to docs surface need to land in both spectacular (deprecated) and pageworks (current). Mitigation: keep spectacular's docs code frozen unless critical; route enhancement requests to pageworks; remove from spectacular in v2.0.0 to end the window.
- **User confusion: "which one do I use?"** Mitigation: spectacular's deprecation banners are unambiguous ("this is deprecated; pageworks is the current home"). pageworks's README says "supersedes the docs surface from spectacular v1.x."
- **Schema drift.** Pageworks evolves its `docs.yaml` schema; spectacular's frozen copy diverges. Mitigation: spectacular's awareness is discovery-only, so it doesn't validate schema — drift can't break spectacular. Downstream users picking up new pageworks schemas don't break old spectacular installs because spectacular doesn't read those fields.
- **The handoff prompt becomes noise.** Spectacular prompting "update docs?" after every archive could annoy users. Mitigation: prompt only when SPEC.md or specs/ changed; allow `--no-docs-prompt` flag and `.spectacular/config.yaml` setting to suppress.
- **Pageworks ships immediately without subagents but needs them sooner than expected.** Mitigation: subagent design is explicitly deferred but tracked. If pain emerges within a release, scaffold `pageworks-agents` quickly.
- **Lift-and-shift introduces no immediate value beyond reorganization.** The skill's authoring features (Diátaxis templates, prose patterns) are the new value. Mitigation: M2 adds those as part of the migration, not after.

## Success criteria

- `skills_db/pageworks/` exists, installs as a standalone skill via apm/symlink/manual, works on a fresh project
- `pageworks --version` → `0.1.0`
- `pageworks init`, `pageworks export mkdocs`, `pageworks export docusaurus`, `pageworks doctor` all functional and tested
- pageworks `references/` contains the migrated contract + renderer-adapters + authoring docs, plus new prose-patterns reference and Diátaxis templates
- Spectacular v1.2.0 tagged: docs verbs marked deprecated, banners visible, behavior unchanged, doctor docs slimmed
- Both skills installed simultaneously do not conflict
- A user who runs `spectacular docs init` sees the deprecation banner and is pointed at pageworks
- A user who runs `pageworks init` on a project with no spectacular gets a working docs/ tree

## Follow-on requests (after this lands)

- `pageworks-agents` — design + ship `docs-writer`, `docs-reviewer`, optionally `docs-architect`
- `public-docs-dogfood` — author spectacular's own docs/ using pageworks
- `spectacular-v2-cleanup` — remove deprecated docs surface from spectacular (target spectacular v2.0.0)
