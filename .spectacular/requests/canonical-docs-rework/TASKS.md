---
status: verified
updated: 2026-05-21
related:
  - PLAN.md
  - ../../PRD.md
---

# Tasks — Canonical Docs Rework

All milestones complete. See PLAN.md for context.

## Milestone 1 — PRD reshaped ✅

- [x] Snapshot current `PRD.md` to `PRD@v1.3.md`
- [x] Draft new PRD.md sections in order
- [x] Write "Problem" section
- [x] Write "Constraints" section
- [x] Write "First milestone" section
- [x] Write "Principles (summary)" section
- [x] Write "Related docs" section
- [x] Cross-check routing map — every section has a new home
- [x] Verify line count <250 — landed at **121 lines**

## Milestone 2 — PRINCIPLES.md ✅

- [x] Create `.spectacular/PRINCIPLES.md` with `version: 1.0` frontmatter
- [x] Move principles 1–5 from old PRD
- [x] Add principle 6: Progressive disclosure (+ enforcement hook)
- [x] Add principle 7: Three layers — intent / execution / validation (+ enforcement hook)
- [x] Add principle 8: Humans decide, agents propose (+ enforcement hook)
- [x] For principles 1–5, add `How the skill enforces this:` sub-bullet to each
- [x] Move "Agent Design Principles" from old PRD § 22
- [x] Move "Anti-Entropy Rules" + "Most important insight" from old PRD § 24
- [x] Add Related docs pointer back to PRD

## Milestone 3 — ARCHITECTURE.md ✅

- [x] Create `.spectacular/ARCHITECTURE.md` with `version: 1.0` frontmatter
- [x] § Layout
- [x] § Root layer (with explicit STACK vs ARCHITECTURE distinction)
- [x] § Configuration (config.yaml schema)
- [x] § Frontmatter conventions
- [x] § Ideas / Current / Requests / Skills / Memory / Archive layers
- [x] § Request files — PLAN/TASKS/SESSION/RISKS/VERIFY + PRD-vs-PLAN
- [x] § Lifecycle
- [x] § Versioning
- [x] § Init flow (summary; pointer to cli-bootstrap)
- [x] Add Related docs pointer back to PRD + PRINCIPLES + AGENTS

## Milestone 4 — ROADMAP.md ✅

- [x] Create `.spectacular/ROADMAP.md` with `version: 1.0` frontmatter
- [x] § v1 (current) — 1-line status per item
- [x] § v2 — Workflows layer
- [x] § v2 — Workspaces
- [x] § v2 — Nested workspaces
- [x] § v2 — Multi-agent orchestration
- [x] § v2 — Hook-driven automation
- [x] § v3+ — Context orchestration / Repository operating system
- [x] Add Related docs pointer back to PRD

## Milestone 5 — AGENTS.md rewrite ✅

- [x] Snapshot current `AGENTS.md` to `AGENTS@v1.0.md`
- [x] § What this folder is
- [x] § How to operate
- [x] § Context loading by task (retrieval principles from old PRD § 23)
- [x] § Available skills
- [x] § Creating requests
- [x] § Don'ts
- [x] § Handoff conventions

## Milestone 6 — Cross-links + decision log ✅

- [x] PRD § Related docs points to all 4 companion docs + STACK + DECISIONS
- [x] Each companion doc has `related:` frontmatter + tail Related section
- [x] DECISIONS.md entry added — 2026-05-21 split
- [x] Updated `requests/prd-craft/TASKS.md` task superseded note (clarifier now in ARCHITECTURE.md § Request files)

## Cleanup (deferred — not blocking)

- [ ] Investigate missing `PRD@v1.1.md` snapshot (v1.0, v1.2, v1.3 exist; v1.1 gap)
- [ ] Move this request to `archive/canonical-docs-rework/` (on next archive pass)

## Final state

| File | Lines | Target | Status |
|---|---|---|---|
| `PRD.md` | 121 | <250 | ✅ |
| `PRINCIPLES.md` | 166 | ~150 | ✅ |
| `ARCHITECTURE.md` | 474 | 400–500 | ✅ |
| `ROADMAP.md` | 129 | ~100 | ✅ |
| `AGENTS.md` | 103 | ~80 | ✅ |
| Total | 993 | — | distributed |

Snapshots preserved: `PRD@v1.3.md`, `AGENTS@v1.0.md`.
