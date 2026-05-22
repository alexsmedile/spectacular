---
status: verified
priority: high
owner: alex
updated: 2026-05-21
summary: "Formalize PRD kits as diff-only extensions over a constant base, with adds-slots / modifies-slots / triggers-docs contract"
related:
  - ../prd-craft/PLAN.md
  - ../doc-writer/PLAN.md
  - ../smart-init/PLAN.md
  - ../doctor/PLAN.md
  - ../../ARCHITECTURE.md
---

# Plan — Kits as Plugins

## Goal

Refactor PRD kits from standalone templates into **diff-only extensions** over a constant base PRD template. Each kit declares what it adds, modifies, and which root docs it triggers — nothing more.

## Why

Today's kits (`coding.md`, `product.md`, `content.md`, `research.md`, `blank.md`) duplicate the full base template, then add a few extra sections. This:
- Makes the base hard to evolve — every change propagates 5 times
- Hides what each kit actually contributes
- Breaks composability — can't apply two kits to one project
- Couples PRD shape to kit shape (a kit can't say "I also need STACK.md")

Treating kits as plugins separates **what a kit is** (a typed delta) from **what the base is** (the canonical contract). This unlocks the [[smart-init]] flow, which needs kits to declare their doc requirements.

## Scope

**In scope (v1)**
- Define the kit extension contract: `adds-slots`, `modifies-slots`, `triggers-docs` (always + suggested)
- Refactor 5 existing kits to the new diff-only format
- Update `prd-grill.md` to apply base + active kit's diff at scaffold time
- Update `prd-review.md` so the gate runs against base slots always; kit slots only when that kit is declared active
- Document composition rules (can multiple kits apply? precedence?)
- Frontmatter on each kit file declaring which slots it owns and which docs it triggers

**Out of scope (v2)**
- Multi-kit composition runtime (define rules, defer enforcement)
- User-authored kits in `.spectacular/templates/prd/kits/` beyond override (covered by existing override mechanism)
- CLI flag plumbing for kit selection (lives in [[smart-init]])

**Explicit anti-patterns**
- Kits as full standalone templates — superseded by diff-only model
- Kits modifying the base PRD's required slot list — base is the contract, kits extend it

## Constraints

- Base PRD slot list (Vision / Problem / Target users / Deliverable / Goals & success criteria / Non-goals / Constraints / First milestone) is **fixed** after [[prd-craft]] v1.1 lands
- Existing PRDs scaffolded with the old kit format must remain valid — no breaking change to file shape
- Kit file shape itself can break (no shipped consumers yet)

## Milestones

1. **Contract defined** — kit-extension schema documented in `references/kits-contract.md` (new), with examples
2. **Kits refactored** — 5 existing kits rewritten as diff files; each <50 lines
3. **Grill+review wired** — `prd-grill.md` reads kit's `adds-slots`, `prd-review.md` gate respects kit declarations
4. **Composition rules documented** — single-kit-only in v1, multi-kit deferred with rationale
5. **Dogfood** — re-scaffold a test PRD with `coding` kit and confirm output matches base + coding diff

## Tasks

See `TASKS.md`.

## Dependencies

- **Hard dependency on [[prd-craft]] v1.1** — base slot list (Vision + Deliverable added) must land first
- **Hard dependency on [[doc-writer]]** — kits register their additions as deltas to the PRD registry entry; without the registry, "kit-aware gate" has nothing to look up. The kit contract is materially simpler once the doc-writer engine exists.
- **Downstream: [[smart-init]]** — consumes the `triggers-docs` field from kit frontmatter to decide which root docs to scaffold
- **Touches [[doctor]]** — doctor's gate checks may need to know which kit (if any) was used to scaffold a PRD; surface kit identity in PRD frontmatter

## Validation

- All 5 refactored kits load cleanly through grill
- Review gate passes on a base-only PRD with no kit
- Review gate flags missing kit-required slot when kit is active
- A kit's `triggers-docs.always` list is machine-readable (YAML) and consumable by [[smart-init]]

## Deliverables

- `references/kits-contract.md` (new) — the extension schema spec
- 5 rewritten kit files in `templates/prd/kits/` (diff-only)
- Updated `prd-grill.md` (apply diff at scaffold)
- Updated `prd-review.md` (kit-aware gate)
- Routing entry in SKILL.md if a new trigger emerges (likely not — kits remain selected inside grill)
