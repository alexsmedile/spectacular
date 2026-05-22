---
status: planned
priority: medium
owner: alex
updated: 2026-05-23
summary: "Skill flow that interviews the user and produces a convention pack — the alex-default dogfood lands here"
related:
  - ../../../skills/spectacular/references/grill.md
  - ../convention-pack-schema/PLAN.md
  - ../convention-pack-application/PLAN.md
---

# Plan — Convention Pack Fabricator

## Goal

Add a skill flow that grills the user (or reads an existing repo as input) and produces a valid convention pack. First dogfood: produce the opinionated `alex-default` pack and place it in the repo app-store folder for distribution.

## Why

The schema (request 1) defines what packs *look like* but not how they get *made*. Without a fabricator, users have to handwrite packs from scratch by reading the schema spec — high friction, error-prone. The fabricator turns pack creation into a 10-minute grill session.

This also stress-tests the schema: if the fabricator can't grill the full schema cleanly, the schema needs adjustment.

## Scope

**In scope (v1)**
- New skill verbs: `spectacular pack new <name>` / `spectacular pack grill <name>` / `spectacular pack refine <name>` / `spectacular pack review <name>`
- Routing in SKILL.md: pack verbs → grill engine with convention-pack registry entry
- `references/pack-overrides.md` — pack-specific slot prompts, mini-refine patterns, gate checks (mirrors prd-overrides.md structure)
- Slot prompts covering the 6 rule categories from the schema
- Mini-refine patterns for common ambiguities ("kebab-case or any-case?" → ask; vague names → flag; missing role suffix list → prompt)
- Review gate: pack passes when all required rule categories have content + schema validation passes
- Output location: `~/.spectacular/packs/<name>/` (user-scope, per registry)
- Source-ingestion mode (v1 stretch): if `--from <path>` flag is passed, read existing files (`~/code/NAMING_RULES.md`, an example `.gitignore`) to pre-fill answers; user confirms/edits

**In scope (v1) — dogfood**
- Produce `alex-default` pack by grilling against this user's existing conventions
- Place result in `<repo-root>/packs/alex-default/` (commit to repo as app-store distribution)

**Out of scope (v2)**
- `pack install <name>` / `pack list` / `pack remove` CLI commands (CLI lifecycle handled by `convention-pack-application`'s install verbs)
- Pack diff/merge (`spectacular pack diff <a> <b>`)
- Pack composition (multi-pack)
- Auto-fabricate from a fully-formed repo (would be doctor-territory + heavy inference)

**Explicit anti-patterns**
- Hardcoded pack templates in the fabricator — packs always come from grill output or user editing, never invented
- Writing to `<project>/.spectacular/packs/` by default — fabricator targets user scope (`~/.spectacular/packs/`); project-local lands via explicit `--scope project`
- Modifying an existing pack without snapshotting first — grill follows the same snapshot rules as PRD grill

## Verification routing

2-of-6 rule applied:
1. User-visible change — ✓ new skill verbs surface for users
2. Reversibility cost — ⚠️ partial (writing to ~/.spectacular/packs/ is user-state; reversible per-file)
3. Multi-surface verification — ✓ grill flow + refine + review + source-ingestion mode + dogfood produces alex-default
4. Risk surface — ✗ no canonical doc edits in the user's project workspace
5. External contract change — ✓ new skill verbs become part of user-facing surface
6. Rollback — ⚠️ partial (revert is rm -rf ~/.spectacular/packs/<name>/)

**Score: 3 of 6** → VERIFY.md scaffolded. Similar weight to smart-init.

## Milestones

1. **Pack registry hooked up** — `pack` doc-id registered; grill engine routes pack verbs correctly
2. **Slot prompts written** — `pack-overrides.md` covers all 6 rule categories with example/good/bad
3. **Grill flow tested cold** — fabricate a throwaway pack from scratch; confirm all slots write to disk
4. **Source-ingestion mode** — `--from <path>` pre-fills answers from existing files; user confirms
5. **alex-default dogfood** — grill against the user's conventions; produce `packs/alex-default/`
6. **Review gate** — pack passes the gate; doctor's `kits` area extended (or new `packs` area) for pack validity
7. **VERIFY.md** + tests

## Tasks

See `TASKS.md`.

## Dependencies

- **Hard dependency on [[convention-pack-schema]]** — needs the schema + registry entry to exist before this can build
- **Touches [[doc-writer]]** — extends the registry; grill/refine/review engines consume pack-overrides.md
- **Downstream:** [[convention-pack-application]] benefits — packs produced here are what application consumes

## Validation

- `spectacular pack new test-throwaway` walks the 6 categories and produces a parseable pack
- `spectacular pack review` correctly flags an incomplete pack
- `spectacular pack new alex-default --from ~/code/NAMING_RULES.md,~/.gitignore` pre-fills with > 50% of answers before grilling
- alex-default pack lives at `<repo>/packs/alex-default/` after dogfood; passes its own review gate
- The schema (from request 1) handles everything the grill produces — no field surprises

## Deliverables

- Updated `skills/spectacular/references/doc-registry.md` — `pack` entry refined with overrides path
- New `skills/spectacular/references/pack-overrides.md` — slot prompts, mini-refine patterns, gate rules
- Updated `skills/spectacular/SKILL.md` — routing for `pack new` / `pack grill` / `pack refine` / `pack review`
- `<repo>/packs/alex-default/` — first user-grade pack, ships as app-store entry
- `tests/cli/pack.test.sh` — automated coverage of pack grill / review / source-ingestion mode
- `VERIFY.md` — manual scenarios that automation can't cover (interactive grill, source-ingestion judgment calls)
