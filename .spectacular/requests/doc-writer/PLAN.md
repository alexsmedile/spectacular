---
status: verified
priority: high
owner: alex
updated: 2026-05-21
summary: "Generalize prd-grill/refine/review into a single doc-writing engine with a registry of doc types and shared verbs"
related:
  - ../prd-craft/PLAN.md
  - ../kits-as-plugins/PLAN.md
  - ../smart-init/PLAN.md
  - ../../ARCHITECTURE.md
---

# Plan — Doc Writer

## Goal

Extract the grill / refine / review pattern out of PRD-specific references into a **shared doc-writing engine** that handles any structured document type Spectacular knows about. One registry, one engine, many templates. Adding a new doc type = registry entry + template, not a new skill.

## Why

Today's `prd-grill.md`, `prd-refine.md`, `prd-review.md` are written as if PRD is the only doc that gets interactive treatment. But the same three verbs apply to every structured doc in `.spectacular/`:
- **PLAN.md** needs slot-filling (7-slot decomposition)
- **TASKS.md** needs format gating (checklist conventions)
- **PRINCIPLES.md**, **ARCHITECTURE.md**, **ROADMAP.md**, **STACK.md**, **AGENTS.md** — all benefit from guided creation when scaffolded
- **DECISIONS.md** needs an append-only entry flow

Without generalization, Spectacular either ends up with `plan-grill.md`, `tasks-grill.md`, `principles-grill.md` (duplication), or it stays PRD-only and every other doc gets neglected. Neither is acceptable.

The all-in-one skill decision (2026-05-21) makes this the right shape: one orchestrator skill, generic engine references, per-doc overrides only when needed.

## Scope

**In scope (v1)**
- `references/doc-registry.md` — single source of truth mapping doc-type → template + slots + mode + location + scope + snapshot-policy + overrides-path
- `references/grill.md` — generic interactive slot-filler (extracted from `prd-grill.md`)
- `references/refine.md` — generic vibe→spec rewriter (extracted from `prd-refine.md`)
- `references/review.md` — generic quality gate runner (extracted from `prd-review.md`)
- Per-doc override files only when a doc has unique rules:
  - `references/prd-overrides.md` — vague-word list, success-criteria regex, kit awareness
  - `references/plan-overrides.md` — milestone-before-tasks ordering, dependency-link checks
  - `references/tasks-overrides.md` — checklist format, frontmatter `depends_on` / `validates`
- New base templates for the other 8 doc types:
  - `templates/plan/base.md` (7-slot)
  - `templates/tasks/base.md`
  - `templates/principles/base.md`
  - `templates/architecture/base.md`
  - `templates/roadmap/base.md`
  - `templates/stack/base.md`
  - `templates/decisions/entry.md` (append template, not file template)
  - `templates/agents/base.md`
- Unified trigger surface in SKILL.md:
  - `spectacular <doc>` — shorthand: grill if empty, review if filled
  - `spectacular <doc> grill` — interactive slot-fill
  - `spectacular <doc> refine` — vibe→spec rewrite
  - `spectacular <doc> review` — quality gate
- Mode field in registry: `grill` (default) | `append` (DECISIONS) | `freeform` (no slot enforcement)
- Backwards compatibility: `spectacular prd` / `spectacular prd grill` / `spectacular prd refine` / `spectacular prd review` keep working as registry-driven aliases

**Out of scope (v2)**
- Cross-doc consistency checks (e.g. PRD goals referenced in PLAN.md)
- AI-suggested slot answers from existing project context (covered by [[smart-init]] v2 auto-detect)
- Custom slot schemas per-project (registry overrides via `.spectacular/doc-registry.yaml`)
- Doc-to-doc translation (e.g. expand PRD success criteria into PLAN milestones)

**Explicit anti-patterns**
- One skill per doc type (`/spectacular:prd`, `/spectacular:plan`) — superseded by all-in-one decision
- Forcing every doc through grill — DECISIONS uses `append` mode; ROADMAP/AGENTS may use `freeform`
- Duplicating grill/refine/review logic per doc — that's what the engine prevents

## Constraints

- All-in-one skill — no plugin namespacing, no sub-skills
- `prd-craft v1.1` must land first — needs the 8-slot PRD base as the first registry entry to model the schema against
- Existing `spectacular prd*` triggers must keep working unchanged for users — registry-driven aliasing handles this
- Markdown-only — registry is YAML inside `doc-registry.md`, parsed by the skill at runtime

## Milestones

1. **Registry defined** — `doc-registry.md` schema documented + first entry (PRD) written; mode/scope/snapshot fields finalized
2. **Engine extracted** — `grill.md` / `refine.md` / `review.md` written as doc-agnostic; consume registry to know what to do
3. **PRD migrated** — `prd-grill.md` / `prd-refine.md` / `prd-review.md` shrunk to override files (only PRD-unique rules)
4. **PLAN + TASKS added** — registry entries + templates + minimal override files; tested by grilling a throwaway PLAN
5. **Remaining docs added** — PRINCIPLES / ARCHITECTURE / ROADMAP / STACK / AGENTS / DECISIONS registered with appropriate modes
6. **SKILL.md routing unified** — single trigger handler routes `spectacular <doc> <verb>` to the engine; PRD aliases preserved
7. **Dogfood** — grill a fresh PLAN.md and a fresh PRINCIPLES.md from scratch using the engine; confirm no PRD-specific code paths involved

## Tasks

See `TASKS.md`.

## Dependencies

- **Hard dependency on [[prd-craft]] v1.1** — needs the 8-slot PRD base as the canonical first registry entry; engine schema is designed against it
- **Downstream: [[kits-as-plugins]]** — kits become registry-aware (they declare `triggers-docs` for which other registry entries to also scaffold); the kit contract is simpler once the engine is generic
- **Downstream: [[smart-init]]** — consumes registry to know what to scaffold for each kit + always-set; without doc-writer, smart-init has nothing to map kits onto

## Validation

- All 9 doc types have registry entries with required fields filled
- `grill.md` contains zero references to "PRD" — it's doc-agnostic
- Running `spectacular plan grill` on an empty PLAN.md produces a valid 7-slot PLAN without invoking PRD-specific code
- Running `spectacular decisions` appends a new entry (mode: append) without grilling
- `spectacular prd` (legacy trigger) produces identical behavior to pre-doc-writer
- A new doc type can be added in a single commit: registry entry + template + (optional) override file

## Deliverables

- `references/doc-registry.md` — the registry + schema spec
- `references/grill.md`, `references/refine.md`, `references/review.md` — engine
- `references/prd-overrides.md`, `references/plan-overrides.md`, `references/tasks-overrides.md` — per-doc deltas
- 8 new base templates (plan, tasks, principles, architecture, roadmap, stack, decisions, agents)
- Updated `SKILL.md` with unified routing
- Migration note in CHANGELOG for v0.3.0
