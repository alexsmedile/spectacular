---
status: review
updated: 2026-07-10
related:
  - PLAN.md
  - TASKS.md
---

# Session — fleet-arc-wiring

## Current state

Cut (b25) and built end-to-end in one session (2026-07-10). All four milestones shipped as doc-only
edits to the two workflow arcs + SKILL/doc-index. Wiring style: **optional, judgment-gated** dispatch
(the orchestrator MAY reach for the new agents when the change warrants — same worth-it economics as
fan-out), decided before building.

## What shipped

**M1 — repo-explorer wired into build-workflow. ✅**
- New **Step 0a — map unfamiliar ground** before the chain-close: when the orchestrator can't write
  the Approach because the subsystem is unfamiliar, dispatch `repo-explorer` → map → then plan. The
  build-side mirror of the Investigator.
- Arc header step-list + the build↔bug comparison table's "Discover role" row updated (build side is
  now `repo-explorer`, *for planning* — no longer "none").

**M2 — code-reviewer wired as an optional review gate. ✅**
- build-workflow Step 2a (new): consider dispatching `code-reviewer` over a substantial/medium+ diff
  before ticking → triage findings → route fixes.
- Mirror in bug-workflow Step 3 (before a fix is called resolved).
- review→verified gate note: a full-request-diff pass is the highest-leverage place to spend it.

**M3 — test-verifier wired as an optional arms-length verify. ✅**
- Both arcs' Step 3: when the builder/fixer self-reported the pass or blast radius is medium+,
  dispatch `test-verifier` for independent pass/fail — "the agent that built it shouldn't grade it."

**M4 — triggers + index + tables coherent. ✅**
- `SKILL.md` build-workflow + bug-workflow trigger rows name the new optional fleet.
- `doc-index.md` build-workflow row updated.
- Both comparison tables symmetric; `doctor links docs` → 0 errors.

## Verification

- M1: `grep repo-explorer build-workflow.md` → 5 refs incl. the table row. ✓
- M2/M3: `code-reviewer` + `test-verifier` present in **both** arcs. ✓
- M4: `./cli/spectacular doctor links docs` → 0 errors; full `doctor` → 0 errors. ✓

## Next

Ready for `review → verified`. The one deferred item (does repo-explorer also earn a *bug*-workflow
map-before-investigate slot?) is parked as v2 (`- [~]`), pending real build-side use.
