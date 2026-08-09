---
type: refactor-program-contract
contract: corrected-executable-refactor-program
version: 1.0
status: accepted
decision_session: S12B-R1
source_handoff: H21-R1
accepted_by: owner
accepted_at: 2026-08-09
central_disposition: accept
supersedes: EXECUTABLE-REFACTOR-PROGRAM-CONTRACT.md@1.0
upstream:
  - SPECIFICATION-TOPOLOGY-CONTRACT.md@1.0
  - IMPLEMENTATION-ARCHITECTURE-AND-MIGRATION-CONTRACT.md@1.0
  - SUBSYSTEM-SURVIVAL-CONTRACT.md@1.0
  - MISSION-PREPARATION-CONTRACT.md@1.0
  - EXECUTION-AUTHORITY-CONTRACT.md@1.0
  - EVIDENCE-CLOSURE-CONTINUITY-CONTRACT.md@1.0
next_action: W0-shared-scaffold-design-sufficiency
---

# Corrected Executable Refactor Program Contract

## Purpose

This accepted S12B-R1 contract supersedes only the original H21 implementation-Mission slicing.
All foundation contracts and the approved S12A specification topology remain intact. The original
H21 program remains historical evidence; this is the current program authority. It activates no
code, migration, deletion, release, provider effect, W0 gate, or implementation Mission.

## P0 safety disposition

P0 is required. On H21-R1's verified baseline, PZL-047 and PZL-048 were still pending: Wayfinder
directly reads legacy `type`, and ordinary cleanup deletes matching remote branches. No immutable
accepted evidence proved these v1 defects fixed and reconciled.

P0 restores only those v1 contracts: a shared canonical `kind` reader with required legacy `type`
fallback in affected readers, and separately consented remote deletion so local cleanup never
implies a remote deletion. It requires focused regressions, reference reconciliation, and
independent review. It cannot adopt v2 semantics, publish/freeze v1, migrate a project, or make a
real provider effect.

## Current program

```text
P0 v1 safety stabilization
  → W0 shared-scaffold Design Sufficiency gate
  → M1 semantic records + canonical workspace substrate
  → M2 durable governed Mission loop
  → M3 guided Skill + registry-generated CLI + retrieval/integrity
  → M4 frozen v1 release + final reviewed mapping inventory
  → M5 isolated capsule + v2 release readiness
```

The program is strictly serialized. No implementation Missions run concurrently. A read-only v1
mapping inventory may begin after M1 fixes target semantics, but joins M4 only after independent
review establishes complete classification, recovery pointers, and no unresolved ambiguity. Real
project migrations and cutovers are later project-specific Missions, outside M5.

## Mission charter floor

| Stage | Bounded outcome | Required evidence and stop |
|---|---|---|
| P0 | narrow v1 safety repair | focused regressions, reference reconciliation, independent review; stop on wider v1 scope, ambiguous consent, or failed check |
| W0 | shared package/file/generated-surface ownership and first-slice joinability | Design Sufficiency and Slice Quality verdicts; stop on duplicate authority, legacy leakage, or unresolved shared boundary |
| M1 | semantic records and canonical Markdown substrate | semantic/property/round-trip proof; stop on ambiguity, data loss, or overlap |
| M2 | governed lifecycle, assessment, and reconciliation | transition/refusal/provider-neutral proof; stop on authority/evidence/reconciliation conflict |
| M3 | guided Skill, registry CLI, retrieval, cards/Fog, and integrity | cold-entry/interruption/conflict/drill-down/context plus registry/integrity proof; stop on authority leakage or seam conflict |
| M4 | frozen v1 and reviewed mapping inventory | behavior/recovery/fresh-install and independent inventory proof; stop on unclassified/ambiguous mapping or missing recovery pointer |
| M5 | isolated capsule and v2 release readiness | candidate/idempotence/rollback/capsule-removal proof; stop on core dependency, invalid candidate, or missing readiness evidence |

Each stage requires a separate Mission charter before activation: exact scope/effects, baseline,
branch/worktree choice, checks, retry budget, review trigger, owner gates, recovery point, and
terminal continuity return. Retry or resume never widens authority.

## Deferred Type-2 details

Record encoding/stable keys, receipt/repair representation, card/cache serialization thresholds, Go
package/library/locking/CI choices, and capsule signing/release naming remain for W0 or individual
Mission preparation.
