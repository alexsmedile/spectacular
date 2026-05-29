---
status: verified
priority: medium
owner: alex
updated: 2026-05-29
target_version: v1.10.0
summary: "Audit SPEC.md, promote 1-2 dense capability bullets to specs/<capability>/SPEC.md (target: v1.10.0 — retargeted from v1.9.0, which shipped the versioning doc instead)"
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
- M5 — Snapshot, CHANGELOG entry, release as v1.10.0.

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

## Candidate list (refreshed for v1.8.1 SPEC.md — to be confirmed in M1)

Candidates re-assessed against current SPEC.md v1.4 (35 bullets, as of v1.8.1):

| Capability | Approx. lines | Promotion rationale |
|---|---|---|
| `roadmap` | 6 | 9-phase chain + tiers + 18-check gate + outcome slot — longest bullet, multi-axis |
| `read-verbs` | 5 | 11 verbs + universal flags + skim/full semantics — dense enough to warrant own page |
| `feedback-loop` | 5 | 5-step loop + substrate + auto-promotion + 3-checkpoint rule + PRINCIPLES §9 |
| `doc-engine` | 4 | grill/refine/review × 18 doc types (was 13 at v1.2.1 — grown materially) |
| `migrations` | 5 | Contract + registry pattern + judgment-vs-mechanical distinction |
| `convention-packs` | 4 | 6 categories × 4 scopes × 3 modes |
| `doctor` | 3 | 15 areas + severity model + --fix taxonomy |

M1 re-scores these by actual line count in current SPEC.md and agent-load friction. Pick 1-2 in M2.

## Resolved design questions

These were open questions — resolved here so M3 has no ambiguity:

**Q: Snapshot history per promoted spec?**
A: Yes. Each new `specs/<capability>/SPEC.md` gets snapshotted before its first substantive edit, same as root-level canonical docs. Naming: `specs/<capability>/SPEC@v1.md` (or managed via `spectacular snapshot`).

**Q: Bullet style when compressed — what vs that-it-exists?**
A: Retain *what*. The compressed one-liner in SPEC.md should still convey the capability's core behavior (one clause), plus a link to the spec file. Example: `- **Read verbs (v1.8.0+)** — 11 read-only CLI verbs for cheap cold-start. See [[specs/read-verbs/SPEC]].`

**Q: ARCHITECTURE.md update?**
A: No change needed. ARCHITECTURE.md documents the `.spectacular/` directory structure and file contract — not individual capabilities. The new `specs/<capability>/` files already fit the documented pattern.

## Out of scope

- Reshaping SPEC.md's overall format
- Promoting everything (defeats the purpose of an index)
- Adding new capabilities — pure restructuring
- Changing the `specs/<capability>/SPEC.md` template (use what doctor already validates)

## Notes

This is the first request that explicitly addresses "is the index still the right shape?" Future requests of this type should follow the same audit-pick-promote pattern.
