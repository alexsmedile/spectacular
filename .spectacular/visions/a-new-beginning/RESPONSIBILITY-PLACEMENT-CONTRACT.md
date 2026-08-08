---
type: refactor-foundation-contract
version: 1.0
status: accepted
decision_session: S07
accepted_by: owner
accepted_at: 2026-08-08
central_disposition: accept-with-clarification
upstream:
  - PRODUCT-CONSTITUTION.md@1.0
  - TRUTH-PROVENANCE-FLOOR.md@1.0
  - SUCCESS-EVIDENCE-CONSTITUTION.md@1.0
  - PRODUCT-TRUTH-CONTRACT-MODEL.md@1.0
  - WORK-UNIT-LIFECYCLE-CONTRACT.md@1.0
  - EXECUTION-AUTHORITY-CONTRACT.md@1.0
  - EVIDENCE-CLOSURE-CONTINUITY-CONTRACT.md@1.0
next_session: S08
---

# Responsibility Placement Contract

This accepted S07 contract assigns responsibilities without choosing names,
commands, storage, retrieval, namespaces, or implementation seams.

## Primary owners

| Owner | Responsibility |
|---|---|
| Human | Outcomes, trade-offs, risk appetite, irreversible authority, and Mission/contract disposition |
| Spectacular core | Accepted contracts; Mission semantics; authority envelopes; evidence, assessment, reconciliation; status, continuity, and cold resume |
| Spectacular CLI | Deterministic local validation/projections and authorized scaffolding or migration mechanics |
| Host coding runtime | Bounded planning, edits, tests, retries, and temporary delegated workers |
| Native providers | Git/GitHub/CI/deployment effects, permissions, and attributable receipts |

Spectacular stays complete without companions. No companion, role, runtime, or
provider gains lifecycle or canonical-truth authority by integration.

## Optional companion slate

| Candidate | S07 disposition | Boundary |
|---|---|---|
| Pageworks | Confirmed optional specialist | Owns public-document authoring, schema/rendering, and documentation drift; Spectacular provides source-change context only |
| Bugworks | Confirmed optional specialist | Owns debugging strategy and diagnosis; returns diagnosis/evidence to the owning Mission, not lifecycle mutation or repair closure |
| Specwright | Confirmed optional specialist | Supplies optional adversarial technical-contract review; it is additive and never a prerequisite for Spectacular Capability-Contract quality |
| Decision multiplexer | Deferred experimental idea | High-effort/multi-model brainstorming or coordination may remain a session/skill mechanism until standalone value is proven |
| AI UX profile | Deferred | May later be a decision-multiplexer profile, not a committed product |
| Verifyworks, Wayfinder | Deferred | Neither is MVP-critical or earns separate-product status yet |
| AFK | Responsibility split, not a companion | Authorization/resume stay core; execution belongs to runtime and native Git tooling |

Detailed behavior, state, evidence thresholds, and maintenance of every
companion remain independent future decisions. A confirmed boundary is not an
approved implementation or dependency.

## Roles, modes, adapters, and mechanisms

- Risk-triggered bounded roles: spec architect, security/threat reviewer, and
  request auditor.
- Host-runtime execution profiles: feature, refactor, and migration.
- Provider adaptation plus Spectacular routing: issue triage; external bridges
  remain deferred and read-only/reconcile-first.
- Deterministic drift checks belong at the truth boundary they examine.
- Session handoff, cold resume, and retrospective reconciliation remain core.

## Experimental companion handoff

An optional companion receives a bounded job, stable source references,
authority, scope, permissions, stop conditions, and required evidence. It
returns `succeeded | blocked | failed`, artifact/evidence references, remaining
assumptions, and exactly one safe next action or owner gate.

The handoff is pointer-first and experimental. A companion never mutates
Spectacular lifecycle directly; Spectacular assesses and reconciles an approved
return.

## Extraction test

An extraction earns promotion only after demonstrating a distinct coherent job,
stable external trigger, useful standalone use without Spectacular, owned
substrate or output contract, typed handoff without competing truth, usable
blocked/failed return, optional integration, and acceptable measured
maintenance/context cost. Moving files or reducing router size is not evidence.
