---
status: review
updated: 2026-05-23
related:
  - PLAN.md
  - VERIFY.md
  - ../convention-pack-schema/PLAN.md
---

# Tasks — Convention Pack Fabricator

## v1

### M1 — Pack registry hooked up
- [x] `doc-registry.md` convention-pack entry references `references/pack-overrides.md` as `overrides:` (forward-declared in request 1; file now exists)
- [x] Add SKILL.md routing: `spectacular pack new <name>` / `pack grill <name>` / `pack refine <name>` / `pack review <name>`
- [x] Update SKILL.md References + Templates indexes (pack-overrides.md added)
- [x] Update SKILL.md frontmatter triggers list

### M2 — Slot prompts (pack-overrides.md)
- [x] Write `pack-overrides.md` mirroring `prd-overrides.md` structure
- [x] Slot 1: Name + scope + applies-to (description + reserved-id check)
- [x] Slot 2: Naming rules (folder/file case, forbidden words/prefixes, role suffixes, date-formats)
- [x] Slot 3: Top-level taxonomy (required + opt-in + mono-collection-detect + mono-collection-folders)
- [x] Slot 4: Root files + README contract (required + conditional + optional + readme-contract)
- [x] Slot 5: Gitignore (always-add + opt-in + never-auto-add + language-specific)
- [x] Slot 6: File placement rules (helper-script, architecture-doc, research-artifact, …)
- [x] Slot 7: Project types (adds + template-dir per supported type)
- [x] Mini-refine patterns for each slot (reserved-id, conflicting-rules, tool-generated-in-always-add, unknown-project-type, …)

### M3 — Grill flow tested cold
- [ ] **Interactive (VERIFY S1):** `spectacular pack new throwaway-test` walks all 7 slots → produces parseable pack.md → review gate passes
- [ ] Reserved pack-id rejection confirmed in a live session
- [ ] Mini-refine patterns confirmed firing in a live session
- [x] Mechanical routing confirmed — engine resolves `convention-pack` → grill + pack-overrides via registry; `pack` alias documented in SKILL.md

### M4 — Source-ingestion mode
- [x] `--from <path1>,<path2>` flag documented in pack-overrides.md
- [x] Source-file pattern → slot mapping table written (`.gitignore` → Slot 5; `*NAMING*` → Slot 2; `README.md` → Slot 4; existing folder tree → Slot 3)
- [x] Confidence rule documented (only pre-populate when unambiguous)
- [x] User-confirmation flow documented (`Pre-filled from <path>: <answer>. Keep? [Y/n]`)
- [ ] **Interactive (VERIFY S3):** live source-ingestion accuracy spot-check

### M5 — alex-default dogfood
- [x] Hand-author `packs/alex-default/pack.md` encoding all 10 archived conventions through the 6-rule-category schema
- [x] Ship `packs/alex-default/templates/.gitignore` and `templates/README.md`
- [x] Write `packs/alex-default/references/why-alex-default.md` — full rationale + lineage + when-not-to-use
- [x] Smoke-test frontmatter through awk parser AND Python YAML parser — 7 top-level keys, 6 rule categories, 8 project types, all parse clean
- [x] Removed placeholder `packs/alex-default/README.md` (superseded by pack.md + why-alex-default.md)
- [x] All 10 sections from archived repo-conventions PLAN expressible in the manifest
- [ ] **Interactive (VERIFY S5):** `spectacular pack review alex-default` passes the live skill gate

### M6 — Review gate
- [x] `pack-overrides.md` review gate spec — checks 4-12 (1 rule populated, frontmatter contract, declared templates exist, declared references exist, applies-to valid, version is semver, always-add ∩ never-auto-add empty, naming self-consistency, README contract minimum)
- [x] Universal base checks (placeholder, clarification, frontmatter) still run — explicit note in pack-overrides.md
- [ ] Doctor `packs` check area — deferred to convention-pack-application request (doctor wiring is request 3's surface; spec lives in pack-overrides.md review gate now and gets implemented in CLI later)

### M7 — Tests + VERIFY.md
- [x] VERIFY.md scaffolded with 6 scenario groups (cold grill, resume grill, source-ingestion, review-gate negatives, dogfood validation, schema parity)
- [x] Schema-parity scenario (S6) added — every grill prompt produces a documented schema field, no orphan surfaces
- [x] Dogfood scenarios (S5) ticked for the mechanical parts that completed during this session; live skill-engine items left open for VERIFY signoff
- [ ] `tests/cli/pack.test.sh` — **deferred** to request 3 (`pack new` produces files but the CLI surface for pack lifecycle lands in convention-pack-application; bash test infrastructure makes more sense there)

### Validation (folded)
- [x] All M tasks above completed mechanically OR queued in VERIFY.md for live verification
- [x] alex-default pack present at `packs/alex-default/`
- [x] Doctor + tests still green after this work (run before commit)
- [x] Schema (request 1) confirmed sufficient — every grill prompt maps cleanly to a schema field; no gaps surfaced

## v2 (deferred)

- [ ] `spectacular pack list` / `pack install` / `pack remove` CLI commands — owned by [[convention-pack-application]]
- [ ] Pack diff/merge
- [ ] Multi-pack composition
- [ ] Auto-fabricate from a fully-formed repo (analyze + propose)
- [ ] `--scope project` shortcut for project-local packs (declared in SKILL.md routing but no CLI surface yet)
- [ ] Doctor `packs` check area — implementation lands in request 3
