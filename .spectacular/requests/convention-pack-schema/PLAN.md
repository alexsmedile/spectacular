---
status: verified
priority: high
owner: alex
updated: 2026-05-23
summary: "Define the convention-pack schema: folder shape, frontmatter contract, what a pack can declare. Ship a minimal bundled pack as the proof."
related:
  - ../../../skills/spectacular/references/kits-contract.md
  - ../../archive/repo-conventions/PLAN.md
  - ../../archive/convention-pack-fabricator/PLAN.md
  - ../../archive/convention-pack-application/PLAN.md
---

# Plan — Convention Pack Schema

## Goal

Lock the convention-pack schema before anything else gets built. A pack is a **mini-skill folder** that declares repo-structure conventions — naming rules, folder taxonomy, README contract, gitignore defaults, file placement rules, project-type templates. Self-contained, distributable, versionable.

Ship a bundled `minimal` pack as the proof the schema works.

## Why

Without a locked schema, the fabricator (request 2) has nothing to produce against and the application layer (request 3) has nothing to consume. Schema-first prevents the fabricator + application from diverging or having to refactor halfway through.

## Scope

**In scope (v1)**
- Pack folder shape (`templates/`, `references/`, `pack.md` frontmatter)
- `pack.md` frontmatter contract — `name`, `version`, `description`, `extends`, `applies-to`, `rules`, declared template list
- Schema for the 6 rule categories: naming / taxonomy / root-files / gitignore / file-placement / project-types
- `references/packs-contract.md` documenting the schema with examples
- Bundled `minimal` pack at `skills/spectacular/templates/packs/minimal/` — gitignore defaults + README contract stub, nothing else
- Repo "app store" folder at `<repo-root>/packs/` — empty in this request, but the location is established
- Registry entry for `convention-pack` doc type so the existing engine can grill packs (consumed by fabricator)

**Out of scope (v2 — separate requests handle these)**
- The fabricator skill (`spectacular pack new`) — request 2
- Init/new-request/doctor wiring — request 3
- The opinionated `alex-default` pack — produced by request 2's dogfood
- Pack composition (multi-pack per repo)
- Pack auto-detection from existing repo signals

**Explicit anti-patterns**
- Schema declared inline in CLI code — must live in `references/packs-contract.md` and be loadable
- Packs as single files — folders only (per user decision: packs need to bundle templates + references)
- Embedding pack content inside skill code — packs live in `templates/packs/` (bundled), `~/.spectacular/packs/` (user), `<project>/.spectacular/packs/` (project-local), or `<repo>/packs/` (distributable)

## Constraints

- Mirror the existing kit pattern (`kits-contract.md`) for schema shape — keep concepts consistent so users don't learn two grammars
- Schema must be parseable by Bash YAML utilities already in use (awk-based parser)
- v1 schema must be forward-compatible — additional rule categories must not break existing packs
- No new CLI commands in this request (fabricator + application bring those)

## Verification routing

2-of-6 rule applied:
1. User-visible change — ⚠️ partial (only the `minimal` pack is visible; full UX lands in requests 2+3)
2. Reversibility cost — ✗ low (markdown only, no migrations)
3. Multi-surface verification — ✗ schema doc + one bundled pack
4. Risk surface — ✗ no canonical doc edits
5. External contract change — ✓ defines a new contract other requests + users build against
6. Rollback — ✗ trivial

**Score: 1 of 6** → no VERIFY.md needed. Verification lives in PLAN § Validation.

## Milestones

1. **Schema locked** — `references/packs-contract.md` written; frontmatter fields documented; rule categories enumerated with examples
2. **Folder shape ratified** — pack folder convention (templates/ + references/ + pack.md) demonstrated by the minimal pack
3. **`minimal` pack shipped** — `skills/spectacular/templates/packs/minimal/` exists with: pack.md frontmatter, `templates/.gitignore`, `templates/README.md` (the README contract stub), `references/why-minimal.md`
4. **Registry entry added** — `doc-registry.md` gains a `convention-pack` entry (mode: grill, location: `~/.spectacular/packs/<name>/pack.md`, scope: user) so requests 2+3 have something to consume
5. **App-store folder established** — `<repo-root>/packs/` created with a README explaining the distribution model (even if empty for now)
6. **Dogfood** — manually walk through what fabricating a pack would look like *against the schema only* — confirm the schema can express the 10 conventions from the archived repo-conventions PLAN

## Validation

- `references/packs-contract.md` exists with full schema spec
- `skills/spectacular/templates/packs/minimal/pack.md` parses through the existing kit-frontmatter awk pattern
- `doc-registry.md` includes `convention-pack` entry
- A walk-through (mental or written) confirms all 10 conventions from the archived repo-conventions PLAN can be encoded in the schema (no fundamental gaps)
- `packs/` exists at repo root with a README

## Tasks

See `TASKS.md`.

## Dependencies

- **Hard dependency on [[doc-writer]]** ✓ verified — registry pattern + grill engine are the substrate the fabricator will consume in request 2
- **Hard dependency on [[kits-as-plugins]]** ✓ verified — schema pattern mirrors kit contract; awk parser reused
- **Downstream blockers:** [[convention-pack-fabricator]] and [[convention-pack-application]] both wait for this schema to land

## Deliverables

- `skills/spectacular/references/packs-contract.md` — schema spec
- `skills/spectacular/templates/packs/minimal/` — bundled minimal pack (pack.md + templates/ + references/)
- Updated `skills/spectacular/references/doc-registry.md` — new `convention-pack` entry
- `<repo-root>/packs/README.md` — explains the app-store distribution model
- Updated `SKILL.md` templates index — `templates/packs/minimal/`
