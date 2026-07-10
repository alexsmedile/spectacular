---
status: review
priority: medium
owner: alex
updated: 2026-07-10
build: b25
summary: "Wire repo-explorer, code-reviewer, and test-verifier into the workflow arcs (build-workflow.md, bug-workflow.md) + SKILL triggers so the skill actually dispatches them — as optional, judgment-gated steps"
related:
  - PRD.md
---

# Plan — fleet-arc-wiring

## Goal

Make the three new fleet agents (`repo-explorer`, `code-reviewer`, `test-verifier`) reachable by the
skill. They exist in `agents/` and Claude Code loads them, but no workflow arc routes to them — the
orchestrator has no trigger that says "dispatch this one here." Wire them into the arcs as **optional,
judgment-gated** steps so the skill dispatches them when the change warrants, not on every change.

## Constraints

- **Optional gates, never mandatory.** Each new dispatch is a "consider dispatching when …" step with
  the same worth-it economics as the fan-out gate — not a step every change passes through. A trivial
  one-line fix must not incur two extra dispatches.
- **Preserve the orchestrator/dispatch boundary.** The new agents are discovery/review/verify labor;
  the orchestrator still owns planning + all mutation (ticks, lifecycle, ledger). Reviewer/verifier
  return findings/pass-fail; the orchestrator decides what to do with them.
- **Doc-only change.** No CLI code. Verifiable by `doctor` (links/docs areas) + reading. Matches how
  `spec-builder` was wired (arc doc + SKILL trigger row + doc-index row).
- **Mirror symmetry.** build-workflow and bug-workflow stay structurally parallel — the comparison
  tables in both must stay honest after the edit (build side gains a discover role).
- **No new agent files.** The three agents already exist and conform; this request only wires them.

## Understanding

### How it works now

Two arcs route to agents: `build-workflow.md` → `spec-builder`; `bug-workflow.md` →
`debug-investigator` / `debug-fixer` / `debug-researcher`. `SKILL.md` has trigger rows for both.
The build arc's own comparison table admits the asymmetry: build-side "discover role: none — the plan
*is* the findings." There is **no** review step and **no** independent-verify step in either arc — the
orchestrator self-verifies at Step 3. `repo-explorer`, `code-reviewer`, `test-verifier` are unreferenced.

### What changes

- **`repo-explorer`** → build-workflow gains a **Step 0a (map unfamiliar ground)** before Step 0's
  chain-close: when the orchestrator can't write the Approach because it doesn't understand the
  subsystem, dispatch `repo-explorer` → map → then plan. The build-side mirror of the investigator.
- **`code-reviewer`** → both arcs gain an **optional review** at Step 3 (before recording): when the
  diff is substantial / medium+ blast radius, consider dispatching over the diff → triage findings →
  route fixes. Also offered at the `review → verified` gate over the whole request diff.
- **`test-verifier`** → both arcs gain an **optional arms-length verify** at Step 3: when the builder
  self-reported the pass or blast radius is medium+, dispatch for independent pass/fail instead of the
  orchestrator re-running the check itself. Especially "the agent that built it shouldn't grade it."
- Comparison tables in both arcs updated (build side gains a discover role); `SKILL.md` trigger rows +
  `doc-index.md` rows added so the new routing is discoverable.

### What stays the same

- The orchestrator/dispatch boundary, the closed-brief test, the fan-out economics (≥3 independent),
  the single-threaded ledger. The new steps are *optional* — the default path (self-serve, self-verify)
  is unchanged for trivial work. No CLI, no lifecycle-state change, no new agent definitions.

## Decisions

- Chose **optional judgment-gated wiring** over always-run gates — because a review+verify on every
  change doubles dispatch cost against the fleet's own worth-it economics; the orchestrator decides
  when the change earns arms-length review (2026-07-10).
- Chose to **cut a tracked request** over editing the arcs ad-hoc — multi-doc skill change earns
  ledger provenance (2026-07-10).

## Milestones

- M1 — **repo-explorer wired into build-workflow.** A Step 0a "map unfamiliar ground" precedes the
  chain-close; the build↔bug comparison table's "discover role" row is honest (build side = repo-explorer).
- M2 — **code-reviewer wired as an optional review gate.** Both arcs' Step 3 name a "consider
  dispatching code-reviewer when …" step; the review→verified gate offers a full-diff pass.
- M3 — **test-verifier wired as an optional arms-length verify.** Both arcs' Step 3 name a "consider
  dispatching test-verifier when self-reported / medium+ blast radius" step.
- M4 — **Triggers + index + tables coherent.** `SKILL.md` trigger rows + `doc-index.md` rows added;
  both comparison tables updated; `doctor` (links, docs) green.

## Tasks

See `TASKS.md`.

## Dependencies

- Depends on the three agent definitions already shipped in `agents/` (done — commit `1f98220`).
- No blocking decision open; wiring style + process resolved above.

## Validation

- M1 — build-workflow.md contains a Step 0a dispatching `repo-explorer`; the comparison table's
  discover-role row names it. observable: `grep repo-explorer skills/spectacular/references/build-workflow.md`.
- M2 — both arcs reference `code-reviewer` as an optional Step 3 gate. observable:
  `grep -l code-reviewer` matches both `build-workflow.md` and `bug-workflow.md`.
- M3 — both arcs reference `test-verifier` as an optional arms-length verify. observable: same grep, both files.
- M4 — `SKILL.md` + `doc-index.md` reference all three; `spectacular doctor links` and `doctor docs`
  report 0 errors. run: `./cli/spectacular doctor links docs`.

## Deliverables

- Edited `skills/spectacular/references/build-workflow.md` (Step 0a + Step 3 gates + comparison table).
- Edited `skills/spectacular/references/bug-workflow.md` (Step 3 gates + comparison table).
- New trigger rows in `skills/spectacular/SKILL.md`; new rows in `skills/spectacular/references/doc-index.md`.
- ROADMAP ledger row: build b25 → target version.
