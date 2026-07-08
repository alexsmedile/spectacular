---
status: verified
updated: 2026-05-29
related:
  - PLAN.md
---

# Tasks — spec-refactor

## M1 — Density audit

- [x] Read SPEC.md bullet-by-bullet; record line count for each capability
- [x] Note cross-references: which bullets are referenced from other bullets, references/, or kit docs
- [x] Note hypothetical or real agent-load friction (would the agent have benefited from a focused spec?)
- [x] Produce a ranked candidate table — top 3-5 promotion candidates
- [x] Decision: which 1-2 to promote (per "Don't promote everything" constraint) → **doc-engine** + **roadmap**

## M2 — Justification per pick

- [x] For each pick, write a one-paragraph rationale: *why this one, why now* → recorded in `specs/doc-engine.md` § Decisions (doc-engine = densest, highest agent-load friction; roadmap = most multi-axis)
- [x] Record the rationales — single source is the spec's Decisions log (not duplicated here)
- [x] Sanity check: would promotion *actually* improve agent ergonomics, or is it cosmetic? → real: both bullets were 5-6 lines, multi-axis

## M3 — Promotion

- [x] `spectacular snapshot .spectacular/SPEC.md` (before edit) → `snapshots/SPEC/@v2.md`, v1.4→1.5
- [x] For each picked capability:
  - [x] Create `.spectacular/specs/<capability>/SPEC.md` with valid frontmatter
  - [x] Lift the current SPEC.md bullet text + expand to standalone narrative
  - [x] Add `related:` links to the spec file
  - [x] Compress the SPEC.md bullet to one-line + link to `specs/<capability>/SPEC.md`
- [x] Verify SPEC.md still reads cleanly as an index (no half-promoted state)

## M4 — Doctor green

- [x] `bash cli/spectacular doctor specs` exits 0
- [x] New spec files appear in the doctor report
- [x] `bash cli/spectacular doctor frontmatter` validates the new files
- [x] No new warnings introduced anywhere

## M5 — Ship

- [x] CHANGELOG entry naming the promoted capabilities (target: **v1.10.0** — retargeted from the stale v1.9.0 placeholder)
- [x] Verify: each compressed SPEC.md bullet is exactly 1 line + link (no multi-line survivors)
- [x] Verify: promoted `specs/<capability>/SPEC.md` files each have ≥3 sections beyond frontmatter (both have 7)
- [ ] Standard release flow: bump-manifests → tag → push → GitHub Release → marketplace (user-triggered) — *in progress via /wrap-up v1.10.0*

## Scope note

This request grew beyond the original 2-capability promotion: the doc-engine spec drove a 9-decision design tree, which surfaced (and fixed) staleness in SKILL.md + grill.md plus a broader skill-file audit (self-describing refs + catalog.sh + legacy cleanup) — all shipped in the same v1.10.0 line. See CHANGELOG [Unreleased].
