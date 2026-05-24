---
status: planned
updated: 2026-05-24
related:
  - PLAN.md
---

# Tasks — spec-refactor

## M1 — Density audit

- [ ] Read SPEC.md bullet-by-bullet; record line count for each capability
- [ ] Note cross-references: which bullets are referenced from other bullets, references/, or kit docs
- [ ] Note hypothetical or real agent-load friction (would the agent have benefited from a focused spec?)
- [ ] Produce a ranked candidate table — top 3-5 promotion candidates
- [ ] Decision: which 1-2 to promote (per "Don't promote everything" constraint)

## M2 — Justification per pick

- [ ] For each pick, write a one-paragraph rationale: *why this one, why now*
- [ ] Record the rationales in this TASKS file under the M2 heading
- [ ] Sanity check: would promotion *actually* improve agent ergonomics, or is it cosmetic?

## M3 — Promotion

- [ ] `spectacular snapshot .spectacular/SPEC.md` (before edit)
- [ ] For each picked capability:
  - [ ] Create `.spectacular/specs/<capability>/SPEC.md` with valid frontmatter
  - [ ] Lift the current SPEC.md bullet text + expand to standalone narrative
  - [ ] Add `related:` links to the spec file
  - [ ] Compress the SPEC.md bullet to one-line + link to `specs/<capability>/SPEC.md`
- [ ] Verify SPEC.md still reads cleanly as an index (no half-promoted state)

## M4 — Doctor green

- [ ] `bash cli/spectacular doctor specs` exits 0
- [ ] New spec files appear in the doctor report
- [ ] `bash cli/spectacular doctor frontmatter` validates the new files
- [ ] No new warnings introduced anywhere

## M5 — Ship

- [ ] CHANGELOG entry naming the promoted capabilities
- [ ] Bump version (probably patch — this is restructuring without behavior change)
- [ ] Manifest alignment via git-guard `bump-manifests.sh`
- [ ] Tag + push + GitHub Release
- [ ] `/plugin marketplace update spectacular` (user-triggered)
- [ ] Snapshot PLAN + TASKS, archive request

## Open questions

- [ ] Does each promoted capability get its own snapshot history (`specs/<capability>/SPEC@v1.md`) or just the index?
- [ ] Should ARCHITECTURE.md be updated to mention any promoted capabilities by name, or stay generic?
- [ ] When SPEC.md drops a capability to one line + link, should the bullet still summarize *what* it does or just *that it exists*?
