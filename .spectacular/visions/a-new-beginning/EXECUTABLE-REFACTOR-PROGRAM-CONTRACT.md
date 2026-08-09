---
type: refactor-program-contract
contract: executable-refactor-program
version: 1.0
status: accepted
decision_session: S12B
source_handoff: H21
accepted_by: owner
accepted_at: 2026-08-09
central_disposition: accept
upstream:
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

This accepted S12B contract compiles the approved v2 specifications into five proposed,
strictly-serialized implementation Missions. It is planning authority only: it activates no
Mission, code change, migration, deletion, release, or provider effect. W0 is the only
next-ready action.

## Wave order

```text
W0 shared-scaffold Design Sufficiency gate
  → W1 M1 Canonical semantics
  → W2 M2 Governed Mission loop
  → W3 M3 Guided interface and retrieval
  → W4 M4 Go core, workspace, CLI, and integrity
  → W5 M5 Cutover and recovery capability
```

No implementation Missions run concurrently. A read-only v1 mapping-inventory lane may begin only
after M1 defines target semantics; it joins M5 only after M4 fixes physical representation and an
independent review confirms its completeness and recovery coverage.

## Proposed Mission boundaries

| Mission | Outcome | Depends on | Non-goals | Minimum proof and stop boundary |
|---|---|---|---|---|
| M1 | Semantic records, identities, cardinalities, lifecycle vocabulary, ownership rules, and cold-resume inputs | W0 | transitions, evidence sufficiency, reconciliation, persistence, public grammar | unit/property invariants plus independent boundary review; stop on semantic ambiguity, overlap, drift, or exhausted R2 repair budget |
| M2 | preparation, authorized transitions, envelope revalidation, assessment, owner disposition, reconciliation, and continuity return | M1 | physical schema, public grammar, provider effects, autonomous resolution | transition/refusal proof, fake-provider receipts, provider-neutral Mission proof; stop on authority/evidence/scope/reconciliation conflict or exhausted R2 |
| M3 | Skill workflows, Anchors, cards, Fog, authority spine, retrieval, and guided-interface sources | M1–M2 | command effect classification, Guardrails enforcement, CLI ownership, provider-dependent retrieval | deterministic semantic-operation harness plus cold-entry/interruption/conflict/drill-down/context evaluations; stop on authority leakage, comprehension regression, unstable interface, or exhausted R2 |
| M4 | Go Markdown mechanics, atomic persistence, indexes/projections, registry/CLI, Guardrails, adapters, builds, and distribution wiring | M1–M3 | capsule, legacy parsing/aliases/fallbacks/dual writes, semantic redefinition | round trips, registry/integrity/refusal/projection/adapters/dependency-isolation/macOS-Linux proof; stop on loss, authoritative cache, legacy dependency, registry duplication, scaffold drift, or exhausted R2 |
| M5 | isolated capsule, mapping enforcement, candidate conversion, ambiguity refusal, validation, rollback, and capsule removal proof | M1–M4 + accepted inventory | legacy behavior in core, implicit invocation, real-project migration, ambient release/cutover authority | inventory/candidate/rollback/removal evidence under independent review; stop on unclassified/ambiguous mapping, missing recovery, core dependency, or exhausted R1 |

Each Mission remains a proposal until separately prepared, authorized, and activated. Its branch,
worktree, effects, check commands, reviewer, recovery point, and owner gate are declared by that
Mission's charter, not inferred from this program.

## W0 shared-scaffold gate

Before M1, a proportional Design Sufficiency and Slice Quality review fixes the shared
file/package topology and exclusive ownership of shared/generated surfaces. It must demonstrate:

- M1 semantics, M2 transitions, M3 guided/retrieval use, M4 mechanics, and M5 capsule boundaries
  are distinct and joinable;
- no legacy core dependency, duplicate authority, provider capture, or authoritative projection is
  introduced;
- the first vertical slice has a coherent stop/recovery point; and
- required owner gates and independent review triggers are named.

W0 may refine a proposed Mission or expose a Gap. It cannot activate implementation or decide a
new Type-1 direction.

## Operating and evidence posture

- Routine in-scope local work needs the applicable Mission envelope, not per-edit approval.
- Consequential effects, Contract reconciliation, release/cutover, ambiguous mapping, and reserved
  effects retain owner gates.
- R2 permits at most two hypothesis-changing, scoped repair cycles for M1–M4. M5 uses R1. Retry,
  handoff, or resume never widens authority.
- Every wave pins its accepted commit/tree and evidence receipt. Upstream revision invalidates
  downstream preparation; it is never silently normalized.
- Canonical Markdown remains usable without indexes/projections. Providers supply their own
  receipts and never gain product/lifecycle authority.

## Deferred Type-2 work

Exact packages, paths, files, libraries, encodings, atomic-write/locking mechanics, receipts and
projection serialization, cache thresholds, flags, Skill decomposition, adapters, CI/property-test
tooling, capsule dependencies/signing/reports, Windows delivery, and real-project cutover charters
remain for W0 or individual Mission preparation.
