---
status: verified
updated: 2026-05-21
related:
  - PLAN.md
  - ../doc-writer/PLAN.md
  - ../smart-init/PLAN.md
---

# Tasks — Kits as Plugins

## v1

### M1 — Contract defined
- [x] Draft `references/kits-contract.md` with: schema, frontmatter spec, three example kits (worked example), composition rules
- [x] Field names locked: `adds-slots`, `modifies-slots`, `triggers-docs.always`, `triggers-docs.suggested`
- [x] Composition policy for v1: single-kit-only (decision documented in spec with v2 sketch)

### M2 — Kits refactored
- [x] Rewrite `kits/blank.md` (empty diff — declares no additions/triggers)
- [x] Rewrite `kits/coding.md` (Stack required + Interfaces optional; triggers STACK + ARCHITECTURE always; PRINCIPLES/ROADMAP/DECISIONS suggested)
- [x] Rewrite `kits/product.md` (User stories + Metrics + Distribution; triggers ROADMAP always)
- [x] Rewrite `kits/content.md` (Audience + Format + Distribution + Editorial principles; triggers ROADMAP always)
- [x] Rewrite `kits/research.md` (Hypothesis + Method + Decision-being-informed; triggers DECISIONS suggested)
- [x] Versioned all kits to 2.0; v1.1 snapshots preserved at `versions/kit-<name>@v1.1.md`
- [x] Line-count target relaxed: 3 of 5 kits exceed 50 lines (acceptable — bulk is YAML prompts/examples + human-readable bodies)

### M3 — Grill+review wired
- [x] Update `prd-overrides.md` kit-selection: replace hardcoded menu with frontmatter-driven discovery
- [x] Update `grill.md` § 2a: explicit "kit application" step after scaffold (parse frontmatter, insert added slots, layer modify notes, write `kit:` frontmatter)
- [x] Update `grill.md` § 3 slot loop: walks merged sequence (base + kit adds-slots at declared positions)
- [x] Update `grill.md` § 4 slot prompts: priority order includes kit-added slot prompts + modifies-slots note layering
- [x] Update `prd-overrides.md` review gate: new check 10 (kit-aware), required vs optional logic documented
- [x] Decision: kit identity stored in PRD frontmatter as `kit: <name>` (used by gate + smart-init)

### M4 — Composition rules documented
- [x] Single-kit-only constraint documented in `kits-contract.md` § Composition with rationale (4 reasons)
- [x] v2 multi-kit composition sketch in spec (precedence, conflicts, frontmatter shape `kits: [...]`)

### M5 — Dogfood
- [x] Re-scaffold a coding-kit PRD in `/tmp/spectacular-kits-test/` — merged slot order confirmed: 10 slots (base 8 + Stack@8 + Interfaces@9), First milestone pushed to 10
- [x] Verify frontmatter writes `kit: coding`
- [x] Kit-aware gate tested with 3 scenarios:
  - kit declared + required slot empty → gate fails with specific message ✓
  - kit declared + required slot filled + optional empty → gate passes ✓
  - no kit declared → gate skips all kit checks ✓

### Index updates
- [x] SKILL.md References index includes `kits-contract.md`

## v2 (deferred)

- [ ] Multi-kit composition runtime (precedence rules, conflict resolution, slot-insertion ordering)
- [ ] User-authored kits with custom slot schemas (project-local registry overrides)
- [ ] Kit marketplace / discovery beyond bundled
- [ ] Kits for non-PRD docs (PLAN kits? PRINCIPLES kits?) — requires extending kit-support beyond PRD in registry
