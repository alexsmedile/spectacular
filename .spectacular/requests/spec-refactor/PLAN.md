---
status: planned
priority: medium
owner: alex
updated: 2026-05-24
target_version: v1.3.0
summary: "Review SPEC.md, identify dense capabilities, and promote 1-2 to per-capability spec files under specs/<capability>/SPEC.md"
related:
  - ../../SPEC.md
  - ../../ARCHITECTURE.md
---

# Plan — spec-refactor

## 1. Goal

`.spectacular/SPEC.md` currently has 35 capability bullets, several of which are dense enough (3-6 lines) that they're approaching the point where the per-capability spec pattern (`specs/<capability>/SPEC.md`) becomes worthwhile. This request audits the index, picks the best candidates, and promotes them — without bloating the index, and without promoting prematurely.

## 2. Constraints

- **Don't promote everything.** SPEC.md works as an index *because* it's terse. Promotion is only justified by either (a) the bullet is too long to read in one glance, or (b) an agent has had to repeatedly load full context that should have been a focused spec.
- **Promotion is mechanically simple but semantically risky.** The bullet stays in SPEC.md but compresses to one line + link to `specs/<capability>/SPEC.md`. The promoted file needs frontmatter, related-links, and standalone narrative — not just the lifted bullet text.
- **Doctor rules must keep passing.** `spectacular doctor specs` validates each promoted spec has frontmatter + a non-empty body. Promotion must satisfy these rules.
- **No format changes to SPEC.md itself.** This isn't an opportunity to reshape the index format — just promote dense entries.

## 3. Milestones

- M1 — Density audit: read SPEC.md bullet-by-bullet, score each by line count + cross-reference frequency + agent load friction. Output: ranked candidate list.
- M2 — Pick 1-2 (per "Don't promote everything" constraint). Justify each pick against a real or hypothetical agent-load case.
- M3 — Promote: scaffold `specs/<capability>/SPEC.md` for each, lift + expand the bullet, compress the SPEC.md entry to one-line + link, snapshot SPEC.md before edit.
- M4 — Doctor green: `spectacular doctor specs` passes; the new spec files appear in the report.
- M5 — Snapshot, CHANGELOG entry, release as v1.3.0 (or fold into next release).

## 4. Tasks

See `TASKS.md`.

## 5. Dependencies

None. This is a pure cleanup pass; no upstream blockers.

## 6. Validation

- M1 — A scored audit table exists somewhere (PLAN, TASKS, or scratch file) showing density per capability
- M2 — For each promotion pick, a one-paragraph rationale in TASKS.md explaining *why this one*
- M3 — `specs/<capability>/SPEC.md` files exist with valid frontmatter and standalone body
- M4 — `spectacular doctor specs` exits 0 and reports the new files
- M5 — SPEC.md still readable as an index (no degradation), CHANGELOG mentions the promoted capabilities

## 7. Deliverables

- 1-2 new files at `.spectacular/specs/<capability>/SPEC.md`
- Updated `.spectacular/SPEC.md` with corresponding bullets compressed to one-line + link
- Snapshot of pre-promotion SPEC.md
- CHANGELOG entry under next release

## Candidate list (starting hypothesis — to be revisited in M1)

Strong promotion candidates surfaced during the v1.2.1 audit:

| Capability | Current bullet length | Promotion rationale |
|---|---|---|
| `doctor` | 3 lines | 10 areas + severity model + --fix taxonomy is real spec material |
| `migrations` | 5 lines | Contract + registry pattern + judgment-vs-mechanical distinction |
| `convention-packs` | 4 lines | 6 categories × 4 scopes × 3 modes is a matrix that wants its own page |
| `lifecycle` | 1 line in SPEC.md | Underspecified actually — could grow rather than promote |
| `cli-mutators` | 3 lines | 5 verbs + mutation principle + frontmatter helpers |
| `doc-engine` | 4 lines | grill/refine/review × 13 doc types |
| `roadmap` | 6 lines | 9-phase chain + tiers + 18-check gate — the longest current bullet |

Pick from these in M2 based on actual usage signal, not just density.

## Out of scope

- Reshaping SPEC.md's overall format
- Promoting everything (defeats the purpose of an index)
- Adding new capabilities — pure restructuring
- Changing the `specs/<capability>/SPEC.md` template (use what doctor already validates)

## Notes

This is the first request that explicitly addresses "is the index still the right shape?" Future requests of this type should follow the same audit-pick-promote pattern.
