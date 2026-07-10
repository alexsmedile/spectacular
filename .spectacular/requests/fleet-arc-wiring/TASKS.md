---
status: review
updated: 2026-07-10
related:
  - PLAN.md
---

# Tasks — fleet-arc-wiring

<!--
  Executable checklist for one request.
  Grouped by milestone. Flush-left checkboxes are the counted units.
-->

## v1

### M1 — repo-explorer wired into build-workflow
- [x] Add a **Step 0a — map unfamiliar ground** to `build-workflow.md` before Step 0: when the orchestrator can't write the Approach because the subsystem is unfamiliar, dispatch `repo-explorer`, read its map, then plan
- [x] Update the build↔bug comparison table's "Discover role" row: build side = `repo-explorer` (was "none — the plan *is* the findings"), with the nuance that it maps *for planning*, not diagnoses
- [x] → check: `grep repo-explorer skills/spectacular/references/build-workflow.md` matches; the table row is honest

### M2 — code-reviewer wired as an optional review gate
- [x] Add an optional "consider dispatching `code-reviewer`" step at `build-workflow.md` Step 3 (before ticking): when the milestone diff is substantial / medium+ blast radius → review → triage findings → route fixes to fixer/builder
- [x] Add the mirror step at `bug-workflow.md` Step 3 (before logging the fix)
- [x] Note the review→verified gate option: a full-request-diff `code-reviewer` pass at the transition
- [x] → check: `grep -l code-reviewer` matches both build-workflow.md and bug-workflow.md

### M3 — test-verifier wired as an optional arms-length verify
- [x] Add an optional "consider dispatching `test-verifier`" step at `build-workflow.md` Step 3: when the builder self-reported the pass OR blast radius is medium+, dispatch for independent pass/fail instead of self-re-running
- [x] Add the mirror step at `bug-workflow.md` Step 3 (independent verify of the fix)
- [x] Capture the "agent that built it shouldn't grade it" rationale in one line
- [x] → check: `grep -l test-verifier` matches both arcs

### M4 — Triggers + index + tables coherent
- [x] Add `SKILL.md` trigger row(s) so the new dispatch points are discoverable
- [x] Add `doc-index.md` rows for any new reference surface (or note the arcs already cover it)
- [x] Confirm both comparison tables (build-workflow ↔ bug-workflow) are symmetric after the edits
- [x] → check: `./cli/spectacular doctor links docs` → 0 errors
- [x] ROADMAP ledger row: build b25 → target version

## v2 (deferred)

- [~] Consider whether repo-explorer also earns a bug-workflow slot (map-before-investigate on an unfamiliar subsystem) — deferred until the build-side use proves it
