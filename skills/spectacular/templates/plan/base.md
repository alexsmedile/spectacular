---
status: planned
priority: medium
owner: <OWNER>
updated: <DATE>
build: <BUILD>
summary: "<one-sentence description of what this request changes>"
related:
  - PRD.md
---

# Plan — <Request Title>

<!--
  Canonical 7-slot PLAN template for a single request.
  Lives at: .spectacular/requests/<slug>/PLAN.md

  Rules:
  - PLAN is per-request. PRD is project-wide. Never put a PRD inside requests/.
  - This file's frontmatter `status:` is the single source of lifecycle state for the request.
  - All 7 required slots must be filled before this PLAN is considered usable.
  - Replace every <placeholder> with concrete content.
-->

## 1. Goal

<!-- One sentence. What does this request change? -->
<!-- Compress the request's intent. Aligns with PRD's Vision or Goals — this is a slice, not a restatement. -->

<GOAL>

## 2. Constraints

<!-- What's fixed before you start? Inherited from PRD/STACK/PRINCIPLES + request-specific limits. -->

- <CONSTRAINT 1>
- <CONSTRAINT 2>

## Understanding

<!--
  OPTIONAL authoring slot, but REQUIRED before `planned → active` by the
  `understand-before-change` policy (@Implementation). Fill it here for a
  typical request; escalate to a dedicated requests/<slug>/UNDERSTANDING.md
  (same three subheads) for large ones — the policy is satisfied by EITHER.
  Not one of the 7 required authoring slots; it gates implementation, not planning.
-->

### How it works now

<!-- The current behavior/structure this request touches. -->

### What changes

<!-- The specific surfaces this request modifies. -->

### What stays the same

<!-- The boundary — what this change deliberately leaves alone. -->

## Decisions

<!--
  Design calls made inside this request. Format: chose X over Y — because Z.
  Rejected alternatives stay listed; deleting them re-litigates them later.
  Project-wide calls go to DECISIONS.md via `spectacular decide` instead
  (see decisions-rules.md routing table). Empty is fine — no decisions yet.
-->

- <DECISION — chose X over Y because Z>

## 3. Milestones

<!-- Ordered, demoable checkpoints. Outcomes, not tasks. -->
<!-- 3-7 milestones for a typical request. Each is something someone can see working. -->

- M1 — <DEMOABLE OUTCOME>
- M2 — <DEMOABLE OUTCOME>
- M3 — <DEMOABLE OUTCOME>

## 4. Tasks

<!-- Pointer. The executable checklist lives in TASKS.md, grouped by milestone. -->

See `TASKS.md`.

## 5. Dependencies

<!-- Other requests, skills, blocking decisions. Use [[request-slug]] notation. -->

- <DEPENDENCY 1>
- <DEPENDENCY 2>

## 6. Validation

<!--
  How each milestone is verified. Per-milestone checks.
  Each check states its AUTHORITY: a run: command, an assertable property,
  a judgable artifact, or a human-observable behavior (see verify.md kinds).
  A check with no authority can't fail. Aspiration verbs (improve, enhance,
  optimize, handle gracefully) are not checks.
-->

- M1 — <VERIFICATION — e.g. "run: tests/foo.test.sh exits 0" or "observable: X appears in Y">
- M2 — <VERIFICATION>
- M3 — <VERIFICATION>

## 7. Deliverables

<!-- Artifacts that ship out of this request. Concrete files, docs, behaviors. -->

- <DELIVERABLE 1>
- <DELIVERABLE 2>
