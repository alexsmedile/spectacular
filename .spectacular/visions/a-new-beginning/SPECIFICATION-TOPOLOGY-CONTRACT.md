---
type: refactor-foundation-contract
contract: specification-topology
version: 1.0
status: accepted
decision_session: S12A
source_handoff: H19
accepted_by: owner
accepted_at: 2026-08-09
central_disposition: accept
upstream:
  - PRODUCT-CONSTITUTION.md@1.0
  - TRUTH-PROVENANCE-FLOOR.md@1.0
  - SUCCESS-EVIDENCE-CONSTITUTION.md@1.0
  - PRODUCT-TRUTH-CONTRACT-MODEL.md@1.0
  - WORK-UNIT-LIFECYCLE-CONTRACT.md@1.0
  - EXECUTION-AUTHORITY-CONTRACT.md@1.0
  - EVIDENCE-CLOSURE-CONTINUITY-CONTRACT.md@1.0
  - RESPONSIBILITY-PLACEMENT-CONTRACT.md@1.0
  - RETRIEVAL-AND-EARNED-WORKSPACE-CONTRACT.md@1.0
  - PUBLIC-LANGUAGE-AND-INTERFACE-CONTRACT.md@1.0
  - CLEAN-BREAK-CUTOVER-AND-RECOVERY-CONTRACT.md@1.0
  - MISSION-PREPARATION-CONTRACT.md@1.0
  - SUBSYSTEM-SURVIVAL-CONTRACT.md@1.0
  - IMPLEMENTATION-ARCHITECTURE-AND-MIGRATION-CONTRACT.md@1.0
next_session: S12B
---

# Specification Topology Contract

## Purpose

This accepted S12A contract defines the smallest coherent set of v2 specifications, their
exclusive boundaries, acceptance dependencies, and proof posture. These specifications refine
the foundation contracts only: they are not current Capability Contracts, Proposals, Missions, or
implementation authorization. S12B alone may compile these approved specifications into an
executable refactor program.

## Approved specifications

| Spec | Owns | Explicitly excludes |
|---|---|---|
| 1. Canonical Contracts and Mission Semantics | record meanings, identity, cardinality, lifecycle-state vocabulary, and semantic cold-resume inputs | transition authority, evidence sufficiency, reconciliation, and storage mechanics |
| 2. Governed Mission Loop | preparation and slice verdicts; transition authorization; envelope revalidation; evidence sufficiency; assessment; owner disposition; authorized reconciliation; continuity-return requirements | public grammar and physical schema |
| 3. Guided Workspace Interface and Retrieval | Mission vocabulary, Anchors, guided use, source-backed cards, Fog, authority spine, and authority-correct retrieval | command effect classification and Guardrails enforcement |
| 4. Go Core, Canonical Workspace and Integrity | Go-native canonical Markdown mechanics, command registry, Guardrails parsing and invariant enforcement, deterministic indexes/projections, scoped integrity, and their proof | migration capsule |
| 5. Clean-Break Cutover and Recovery | frozen-v1 boundary, isolated capsule, separate candidate, ambiguity refusal, validation, owner cutover, rollback, and removal proof | in-core legacy behavior and conversion execution authorization |

## Dependency order

```text
1 → 2, 3, 4
2 → 3, 4
3 → 4 (generated interface)
1–4 → 5
```

These are typed specification-acceptance dependencies only. They are neither a scheduler nor a
Mission DAG.

## Boundary and proof corrections

- Spec 1 defines semantic records and lifecycle vocabulary only. Spec 2 is the sole authority for
  authorized transitions, envelope revalidation, evidence sufficiency, assessment, owner
  disposition, reconciliation, and terminal continuity-return requirements. Both use the same
  cold-resume scenario without duplicating authority.
- Each specification declares itself non-authoritative for current product truth and accountable
  delivery work. Current behavior remains in Capability Contracts; candidate change remains in
  Proposals; Mission accountability remains in Missions.
- At least one end-to-end proof remains provider-neutral. Fake-provider tests prove provider
  boundaries and receipts separately.
- Specs 4 and 5 mechanically prove that v2 has no legacy parser, alias, fallback read, dual
  write, lazy conversion, capsule import, or legacy retrieval dependency.
- The reviewed v1 mapping inventory is the complete migration test manifest. Candidate acceptance
  fails for an unclassified item, ambiguous mapping, missing recovery pointer, or legacy
  core/retrieval dependency. Independent review is required for that inventory and for
  candidate-validation, rollback, and capsule-removal evidence; only the owner resolves an
  ambiguous mapping.

## Acceptance posture

Per-spec focused acceptance is supplemented by provider-neutral end-to-end proof of cold recovery,
fuzzy intent to bounded work, and returned evidence through assessment, reconciliation, and fresh
resume. Fake-provider tests and fixed model-plus-harness evaluations provide bounded additional
evidence.

## Deferred Type-2 details

Stable-key/hash/record encoding, receipts and repair representation, flags/card serialization/cache
thresholds, Go libraries/package layout/locking/CI, and capsule signing/release naming remain for
S12B or implementation-level decisions.
