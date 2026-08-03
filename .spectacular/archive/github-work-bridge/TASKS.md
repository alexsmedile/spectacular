---
status: verified
updated: 2026-08-03
related:
  - PLAN.md
source_spec: SPC-003
source_type: spec
source_ref: SPC-003
---

# Tasks — github-work-bridge

## v1

### M1 — Routing and readiness
- [x] Implement and test `direct | request | spec-first` classification plus the concise Issue-readiness card.
  - [x] Interpret repository-specific Issue types and labels by meaning.
  - [x] Fail closed for missing authority, unresolved consequential choices, and protected security input.

### M2 — Issue and goal provenance
- [x] Extend request creation and retrieval for Issue/goal sources without regressing approved-spec provenance.
  - [x] Store canonical source references without copying remote bodies or comments.
  - [x] Keep compact PLAN/TASKS artifacts and orchestrator-only lifecycle ownership.

### M3 — PR handoff and reconciliation
- [x] Generalize draft/ready PR handoff and implement read-only reconciliation.
  - [x] Render the PR integration manifest and correct `Fixes`/`Refs` relationships.
  - [x] Require current-head verification and explicit confirmation before ready; never merge.
  - [x] Report broken/stale links and remote ambiguity without bidirectional mutation.

### M4 — Documentation and verification
- [x] Update the skill and user-facing docs, dogfood all three routes, and complete focused plus full verification.
  - [x] Preserve raw `gh` escape hatches and document deferred managed/security/Projects scope.
  - [x] Record docs-impact and verification evidence before lifecycle completion.
