---
status: planned
updated: 2026-05-24
related:
  - PLAN.md
---

# Tasks — add-personas-doc

## M1 — Template + scaffold reference

- [ ] Draft `skills/spectacular/templates/personas/base.md` (frontmatter + structure + 1 filled example + 1 blank)
- [ ] Pick a concrete example persona — project-agnostic enough to apply across kits, specific enough to feel real
- [ ] Verify the template stays under ~40 lines including the filled example
- [ ] Update `skills/spectacular/references/scaffold-reference.md` with a PERSONAS.md stub entry
- [ ] Sanity check: a non-product project (coding/research) could ignore this entirely without confusion

## M2 — Doc-registry entry

- [ ] Add `personas` entry to `skills/spectacular/references/doc-registry.md`:
  - template: `templates/personas/base.md`
  - slots: `who`, `wants-to`, `pain`, `stories`, `not-for`
  - mode: `freeform` (per-persona blocks)
  - location: `.spectacular/PERSONAS.md`
  - overrides: `personas-overrides.md`
- [ ] Write `skills/spectacular/references/personas-overrides.md`:
  - Grill prompts (who is missing, are stories outcome-focused, is "Not for" named, etc.)
  - Slot definitions
  - Anti-patterns (no demographics, no photos, no JTBD framework, no 20+ personas)
  - Vague-word list ("users", "everyone", "people")
- [ ] Verify the generic grill/refine/review engine consumes the new entry without code changes

## M3 — Kit wiring

- [ ] Add `triggers-docs: [personas]` to `skills/spectacular/templates/prd/kits/product.md`
- [ ] Add `triggers-docs: [personas]` to `skills/spectacular/templates/prd/kits/content.md`
- [ ] Update CLI `--with` parser in `cli/spectacular` to accept `personas`
- [ ] Update `init -i` interactive menu to offer personas
- [ ] Confirm `spectacular init --kit product` scaffolds PERSONAS.md automatically
- [ ] Confirm `spectacular init --with personas` works on top of any kit
- [ ] Confirm other kits (blank/coding/research) do NOT scaffold it

## M4 — Doctor area

- [ ] Add `personas` area to `cli/spectacular` `cmd_doctor`:
  - File absent: skip with info note
  - File present: validate frontmatter + ≥1 `## <Persona>` block + ≥1 story per persona
- [ ] Add `personas` to doctor area list in usage + README
- [ ] Tests in `tests/cli/doctor.test.sh`: absent / present-valid / present-invalid

## M5 — Doc engine validation

- [ ] Run `spectacular personas grill` on a near-empty PERSONAS.md — confirm prompts are useful
- [ ] Run `spectacular personas refine` on a partial PERSONAS.md
- [ ] Run `spectacular personas review` on a complete PERSONAS.md

## M6 — Ship

- [ ] CHANGELOG entry under v1.3.0 (alongside spec-refactor if both ship together)
- [ ] Update `.spectacular/SPEC.md` with a Personas capability bullet (or extend doc-engine bullet)
- [ ] Manifest alignment + version bump
- [ ] Snapshot PLAN + TASKS, archive request
- [ ] Tag + push + GitHub Release
- [ ] `/plugin marketplace update spectacular`

## Open questions

- [ ] Should PERSONAS.md frontmatter include a `personas:` count field for briefings?
- [ ] Should the doctor warn on drift between PRD § Target users and PERSONAS.md persona names? (v2 — implies controlled vocabulary)
- [ ] `--from-prd` flag to pre-fill PERSONAS.md from PRD § Target users? (v2)
- [ ] Should `coding` kit also opt into personas? (Dev tools have users too — but probably leave opt-in for now)
