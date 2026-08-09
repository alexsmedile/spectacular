---
type: shared-scaffold-contract
contract: shared-scaffold
version: 1.0
status: accepted
decision_session: W0
source_handoff: H25
accepted_by: owner
accepted_at: 2026-08-09
central_disposition: accept
upstream:
  - EXECUTABLE-REFACTOR-PROGRAM-CONTRACT.md@1.2
  - IMPLEMENTATION-ARCHITECTURE-AND-MIGRATION-CONTRACT.md@1.0
  - SPECIFICATION-TOPOLOGY-CONTRACT.md@1.0
  - SUBSYSTEM-SURVIVAL-CONTRACT.md@1.0
  - MISSION-PREPARATION-CONTRACT.md@1.0
next_action: M1-preparation
---

# Shared Scaffold Contract

## Purpose

This contract fixes the v2 scaffold boundary, exclusive module ownership, generated-surface
authority, Mission joins, and the first proof seam. It is the accepted W0 result. It authorizes M1
preparation only; it does not activate M1 or authorize code, release, provider, or migration effects.

## V2 boundary

- `v2/` is the sole new v2 product boundary.
- It contains one Go module and build graph, the v2 core and CLI, canonical v2 Skill source,
  generated-interface area, and v2 tests.
- V1 CLI, Skill, tests, and collections remain outside all v2 import, build, test-discovery, and
  retrieval paths.
- Release assembly may package the v2 Skill and generated interface, but generated or staged output
  is never canonical.
- W0 and v2 core contain no compatibility behavior or migration capsule.

## Exclusive ownership

| Owner | Exclusive authority |
|---|---|
| Domain | Semantic records, identity, typed relationships, authority, lifecycle vocabulary, evidence, reconciliation invariants, and refusals |
| Workspace | Canonical Markdown parsing, normalization, validation, and atomic persistence |
| Operations | Governed Mission actions, authorization, transitions, assessment, reconciliation, and effect boundaries |
| Index | Deterministic lookup and disposable projections |
| Command registry | Commands, arguments, dispatch bindings, help, effect classification, and machine interfaces |
| Guardrails | Mechanical selection of owner-written guidance without creating authority |
| Adapters | Filesystem, Git, process, runtime, and provider effects with attributable receipts |
| Guided Skill | Judgment, guided authoring, interpretation, and progressive context loading |

Dependencies point inward. No downstream owner may redefine Domain invariants.

## Generated surfaces

The Go command registry is canonical for mechanical command facts. CLI dispatch, help, effect
metadata, and the Skill-facing machine catalog are reproducibly generated from it and never edited
as independent authority. Human-written Skill source is canonical for guided semantics and
judgment; it does not duplicate arguments, effect classifications, or lifecycle rules.

## Typed Mission joins

- **M1 → M2:** typed records, stable identities, relationships, fingerprints, and canonical
  workspace access.
- **M2 → M3:** governed operation requests, results, refusals, owner gates, effect intentions, and
  receipt boundaries.
- **M3 → M1/M2:** reads use M1 lookup; governed changes use M2 operations. M3 never mutates
  canonical state directly.

## M1 proof seam

M1 proves the substrate with two minimal linked records, representatively Proposal → Mission:

- stable identity creation and preservation;
- typed identity relationships independent of mutable slugs;
- canonical Markdown semantic round-trip;
- stable post-normalization fingerprint;
- identical index results independent of discovery order;
- exact identity and relationship lookup;
- refusal of duplicate identities and broken or wrongly typed relationships.

M1 excludes lifecycle operations, authorization, assessment, reconciliation, CLI generation, Skill
behavior, Guardrails, providers, caches, and migration.

## Gate verdicts

- Design Sufficiency: `sufficient`
- Slice Quality: `coherent`

## Deferred Type-2 choices

Exact paths, package/API names, Go version and dependencies, Markdown/YAML and UUIDv7 libraries,
canonical encoding, fingerprint mechanics, atomic-write and locking mechanics, generator details,
test tooling, adapter encoding, release assembly/signing, and persistent-cache thresholds remain
for Mission preparation or later evidence-gated decisions.
