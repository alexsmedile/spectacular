---
type: refactor-program-contract
contract: executable-refactor-program
version: 1.2
status: accepted
decision_session: owner-v1-deprioritization
source_handoff: V1-DEPRIORITIZATION-DECISION.md
accepted_by: owner
accepted_at: 2026-08-09
central_disposition: accept
supersedes: _snapshots/EXECUTABLE-REFACTOR-PROGRAM-CONTRACT/@v1.1.md
upstream:
  - V1-DEPRIORITIZATION-DECISION.md@1.0
  - SPECIFICATION-TOPOLOGY-CONTRACT.md@1.0
  - IMPLEMENTATION-ARCHITECTURE-AND-MIGRATION-CONTRACT.md@1.0
  - SUBSYSTEM-SURVIVAL-CONTRACT.md@1.0
  - MISSION-PREPARATION-CONTRACT.md@1.0
  - EXECUTION-AUTHORITY-CONTRACT.md@1.0
  - EVIDENCE-CLOSURE-CONTINUITY-CONTRACT.md@1.0
next_action: W0-shared-scaffold-design-sufficiency
---

# Executable Refactor Program Contract

## Purpose

This v1.2 program incorporates the owner's scoped decision to stop investing in v1. The approved
S12A topology and v2 foundation remain intact. Earlier program versions are preserved in the
snapshot tree; this unversioned file is the sole current program authority. It activates no code,
migration, deletion, release, provider effect, W0, or implementation Mission.

## V1 disposition

P0 is abandoned and its branch is not merged. No repair, Pageworks reconciliation, new final-v1
release, universal mapping inventory, or generic migration capsule blocks v2. Existing Git history
and published tags remain recovery evidence, not supported product behavior. The controlling
scope override is [`V1-DEPRIORITIZATION-DECISION.md`](V1-DEPRIORITIZATION-DECISION.md).

Real migrations are later project-specific Missions. They may earn isolated scripts, reference
guidance, or a disposable capsule, but v2 core contains no compatibility behavior.

## Current program

```text
W0 shared-scaffold Design Sufficiency gate
  → M1 semantic records + canonical workspace substrate
  → M2 durable governed Mission loop
  → M3 guided Skill + registry-generated CLI + retrieval/integrity
  → M4 clean-v2 release readiness
```

The program is strictly serialized and W0 is the sole next-ready gate. No implementation Missions
run concurrently. Real project migrations and cutovers are later project-specific Missions and do
not block clean-v2 release readiness.

## Mission charter floor

| Stage | Bounded outcome | Required evidence and stop |
|---|---|---|
| W0 | shared package/file/generated-surface ownership and first-slice joinability | Design Sufficiency and Slice Quality verdicts; stop on duplicate authority, legacy leakage, or unresolved shared boundary |
| M1 | semantic records and canonical Markdown substrate | semantic/property/round-trip proof; stop on ambiguity, data loss, or overlap |
| M2 | governed lifecycle, assessment, and reconciliation | transition/refusal/provider-neutral proof; stop on authority/evidence/reconciliation conflict |
| M3 | guided Skill, registry CLI, retrieval, cards/Fog, and integrity | cold-entry/interruption/conflict/drill-down/context plus registry/integrity proof; stop on authority leakage or seam conflict |
| M4 | clean-v2 release readiness | native build/distribution, provider-neutral end-to-end acceptance, install/recovery, and legacy-exclusion proof; stop on authority leakage, legacy dependency, unsupported release claim, or missing recovery evidence |

Each stage requires a separate Mission charter before activation: exact scope/effects, baseline,
branch/worktree choice, checks, retry budget, review trigger, owner gates, recovery point, and
terminal continuity return. Retry or resume never widens authority.

## Deferred Type-2 details

Record encoding/stable keys, receipt/repair representation, card/cache serialization thresholds, Go
package/library/locking/CI choices and release signing remain for W0 or individual Mission
preparation. Migration-tool details belong only to a later project-specific migration Mission.
