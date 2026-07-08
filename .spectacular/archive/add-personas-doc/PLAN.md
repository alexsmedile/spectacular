---
status: archived
priority: medium
owner: alex
updated: 2026-05-26
target_version: v1.3.0
summary: "Add PERSONAS.md as an opt-in canonical doc for proto-audience profiles + user stories; register in doc-registry so grill/refine/review work."
related:
  - ../../specs/index.md
  - ../../ARCHITECTURE.md
archived: 2026-05-26
---

# Plan — add-personas-doc

## 1. Goal

Add `PERSONAS.md` as an **opt-in** canonical doc capturing proto-audience profiles and the user stories tied to them — lightweight, narrative-driven, designed to ground prototype + build decisions in a perspective beyond "the goal." Register it in the doc-registry so the grill/refine/review engine works on it for free.

## 2. Constraints

- **Reasonable size.** Each persona ≈ 6-10 lines. Doc as a whole stays under ~120 lines for typical projects.
- **No JTBD framework apparatus.** A single "What they want to accomplish" line per persona — not a Jobs-to-be-Done methodology import (no job stories, no functional/emotional/social job decomposition).
- **Opt-in only.** Not in always-set. Scaffolded via `--with personas` or triggered by `product` and `content` kits.
- **One file, not a folder.** Personas + their stories live together. No `personas/<id>/SPEC.md` splitting in v1.
- **Stories belong to personas.** Live inside PERSONAS.md attached to a persona. Request-scoped stories still live in PLAN.md if needed.
- **Register, don't special-case.** New doc-registry entry + small `personas-overrides.md`. No special CLI code paths.
- **No breaking changes.** Existing projects continue working. Doctor reports presence as info, never warns on absence.

## 3. Milestones

- M1 — Template + scaffold reference: `skills/spectacular/templates/personas/base.md` (one filled example, one blank); document in `scaffold-reference.md`.
- M2 — Doc-registry entry: register `personas` in `doc-registry.md`; write `personas-overrides.md` (slot prompts, grill questions, anti-patterns).
- M3 — Kit wiring: `product` + `content` kits gain `triggers-docs: [personas]`; `--with personas` flag works; `init -i` offers personas.
- M4 — Doctor area: optional `personas` area validates frontmatter + structure when file exists; never errors on absence.
- M5 — Doc engine works: `spectacular personas grill|refine|review` produces useful output via generic engine + overrides.
- M6 — Snapshot, CHANGELOG, ship.

## 4. Tasks

See `TASKS.md`.

## 5. Dependencies

None. Uses existing kit-extension pattern (`triggers-docs`) and doc-registry.

## 6. Validation

- M1 — Template renders cleanly; example persona is concrete enough to be useful, generic enough to apply to any project type
- M2 — Registry lookup surfaces the new doc-type
- M3 — `spectacular init --with personas` scaffolds PERSONAS.md; `spectacular init --kit product` scaffolds it automatically
- M4 — `spectacular doctor personas` exits 0 on a valid file; gracefully skips when file absent
- M5 — Running grill on a near-empty PERSONAS.md surfaces the right questions (who, what they want, stories, anti-personas)
- M6 — v1.3.0 ships PERSONAS.md alongside spec-refactor

## 7. Deliverables

- `skills/spectacular/templates/personas/base.md`
- `skills/spectacular/references/personas-overrides.md`
- Updated `doc-registry.md` with `personas` entry
- Updated `scaffold-reference.md` documenting the file
- Updated `templates/prd/kits/product.md` + `content.md` with `triggers-docs: [personas]`
- CLI `--with` parser + `init -i` menu updates
- Doctor `personas` area (optional, info-only)
- CHANGELOG entry

## Template shape (provisional — finalize in M1)

```markdown
---
version: 1.0
updated: <date>
summary: "Audience profiles and the user stories that drive build decisions"
related:
  - ../../PRD.md
---

# Personas

> Opt-in. Lightweight. Stories live with the persona.

## <Persona name — short, evocative>

**Who** — One sentence. Role + context. ("Solo OSS maintainer juggling a day job.")

**Wants to** — One sentence. What they're trying to accomplish, outcome-focused.

**Pain** — 1-2 bullets. What's friction-y today.

**Stories** — Behaviors this product enables. "As X, I want Y, so Z."
- As <persona>, I want <action>, so that <outcome>
- ...

**Not for** — *(optional)* Who this is explicitly NOT serving. Helps scope.

---

## <Next persona>
...
```

**Target shape:** 2-5 personas per project, each 6-10 lines. Total doc 40-120 lines.

## Out of scope

- JTBD methodology apparatus (job stories, functional/emotional/social decomposition, outcome statements)
- Persona photos, demographics, age ranges, income brackets
- Per-industry persona templates
- Cross-linking stories from PLAN.md back to personas (v2)
- Auto-extracting personas from PRD § Target users (v2)
- Promoting personas to subfolder

## Decisions (provisional)

- **Name: PERSONAS.md.** Industry-standard. Reads correctly alongside PRD/DECISIONS/PRINCIPLES.
- **One "Wants to" line, not JTBD framework.** Captures the outcome without importing a methodology.
- **User stories live with their persona.** Avoids the STORIES.md trap of becoming a stale backlog.
- **Opt-in, not always-set.** Many projects (CLI tools, internal libs, infrastructure) don't benefit from formal personas.

## Risks

- **Becomes a museum piece** if scaffolded but never grilled. Mitigation: doctor reports presence info; briefings can soft-nudge stale files.
- **Forces premature formalization** in projects that don't need it. Mitigation: not in always-set; only `product` + `content` kits auto-trigger.
- **Drift from PRD § Target users.** Mitigation: position as a "deepening" of the PRD bullet, not a replacement. Both coexist; PRD stays terse, PERSONAS.md goes deeper.
