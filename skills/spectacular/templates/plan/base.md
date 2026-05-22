---
status: planned
priority: medium
owner: <OWNER>
updated: <DATE>
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

<!-- How each milestone is verified. Per-milestone checks. -->

- M1 — <VERIFICATION>
- M2 — <VERIFICATION>
- M3 — <VERIFICATION>

## 7. Deliverables

<!-- Artifacts that ship out of this request. Concrete files, docs, behaviors. -->

- <DELIVERABLE 1>
- <DELIVERABLE 2>
