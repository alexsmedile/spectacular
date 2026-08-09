---
type: refactor-foundation-contract
contract: implementation-architecture-and-migration
version: 1.0
status: accepted
decision_session: S11
source_handoff: H18
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
next_session: S12A
---

# Implementation Architecture and Migration Contract

## Purpose

This accepted S11 contract fixes the target architecture, distribution boundary, canonical
representation posture, migration capsule, and proof architecture for clean-break Spectacular v2.
It turns the accepted semantic and responsibility contracts into a small, deep, testable system.
It does not authorize implementation, migration, retirement, or S12A specifications.

## Architecture decisions

### Go core and native distribution

The v2 semantic core is Go. It is released as checksummed native binaries for supported macOS and
Linux architectures. End users need no Go runtime, package manager, or shell-specific substrate.
A thin Unix installer may select, download, and verify an artifact, but holds no product semantics.

Windows binaries are deferred, not architecturally excluded. Paths, filesystem access, locking,
process execution, and provider boundaries must remain platform-neutral; unsupported Windows use
must fail honestly rather than rely silently on Unix behavior.

The core will not remain an assembled Bash 3.2 implementation, and Bash will not retain semantic
authority behind a hybrid implementation. Build-time Go versions and dependencies are pinned.

### Deep module map

```text
CLI / generated interface       Guided Skill workflows
          |                              |
          +--------- semantic operations-+
                           |
                      domain kernel
   identity · authority · lifecycle · evidence · refusal
          |                 |                 |
canonical Markdown     index/projections    Guardrails
workspace                                        |
                     explicit adapters
     Git · filesystem/process · host runtime · native providers

v1 workspace → isolated migration capsule → v2 candidate workspace
```

| Boundary | Owns | Must not own |
|---|---|---|
| Domain kernel | identity, authority, lifecycle, evidence, reconciliation, refusal invariants | provider effects or guided judgment |
| Markdown workspace | parse, normalize, validate, atomically persist canonical records | competing truth or cache authority |
| Semantic operations | authorized prepare/start/resume/assess/reconcile/resolve changes | product judgment outside accepted envelopes |
| Index and projections | deterministic retrieval and bounded derived views | acceptance, reconciliation, or conflict resolution |
| Command registry | one mechanical source for dispatch, help, effects, and machine interfaces | a second semantic model |
| Skill kernel | judgment, guided authoring, interpretation, progressively loaded roles | deterministic authority validation |
| Guardrails | select owner prose for events/commands | authority creation or invariant weakening |
| Adapters | external effects and attributable receipts | product or lifecycle authority |
| Migration capsule | explicit v1-to-v2 conversion | a v2 core dependency |

This is not a generic agent-DAG scheduler. Typed canonical relationships and retrieval, Mission/Run
control state, and any future graph analytics are distinct concepts. Graph analytics remain deferred.

## Canonical records and retrieval

Canonical authority is OKF-compatible Markdown with structured frontmatter. Stable identity is
UUIDv7; readable typed references are human handles; slugs are mutable labels, never durable
identity. A fingerprint is SHA-256 over normalized UTF-8/LF canonical content. Typed
relationships are explicit and resolvable. Embedded records retain identity if promoted globally.

Each command builds one deterministic in-memory index from canonical records, then generates only
the requested narrow projection: orientation, Mission Card, Fog/frontier, typed authority spine,
exact lookup/relationships, or scoped integrity results. Every projection and cache is
non-authoritative: deleting it must leave the workspace complete and recoverable from Markdown.

A persistent cache is optional and must earn inclusion through representative measurements. If it
earns inclusion, it is modular, fingerprint-invalidated, atomically replaceable, rebuilt from
canonical state, excluded from agent context, and never preferred to canonical records during a
conflict.

## Command and Guardrails boundary

The command registry generates dispatch, help, effect classification, and machine interfaces from
one source. It must not introduce vocabulary absent from the accepted semantic model.

The deterministic Guardrails parser recognizes these fixed lifecycle events:
`@Orient`, `@Prepare`, `@Start`, `@Resume`, `@Run`, `@Assess`, `@Reconcile`, and `@Resolve`.
It also accepts mechanically valid optional `$domain.verb` selectors. Matched owner-authored prose
is returned verbatim to the Skill. The Skill interprets it, but cannot weaken domain invariants,
enlarge authority, or turn missing authority into permission. Native hooks remain optional and are
never silently installed or changed.

Human gates apply to consequential effects and Mission/Capability-Contract resolution—not every
approved in-scope local edit or file write.

## Migration and cutover

v1 freezes at an immutable annotated final SemVer tag with published source, Skill, and
checksummed CLI artifacts. The original v1 project remains authoritative unless an owner accepts a
validated v2 candidate.

The repository contains one isolated Go migration capsule with its own module, executable,
dependencies, fixtures, tests, release artifact, and separate build graph. The v2 core cannot
import it or invoke it implicitly. It inspects without mutating source, converts only to a separate
candidate worktree/destination, reports mappings and ambiguities, stops for owner disposition on
ambiguous semantics, and produces a candidate fingerprint, validation report, cutover receipt, and
rollback pointer. It supports dry runs and idempotence where applicable; it never merges, deletes,
or rewrites an accepted v1 workspace.

Release order is: authorized v1 stabilization; v1 behavior/migration/recovery proof; final tagged
v1 release and fresh-install verification; capsule release; then v2 distribution only after
candidate, rollback, and capsule-removal proofs pass. Removing the capsule must not affect v2
builds, tests, runtime behavior, or recovery.

## Test and proof architecture

The rebuild uses layered, risk-based tests:

- unit/property tests for identity, parsing, authority, lifecycle, evidence, and refusal;
- generated-interface tests from the command registry;
- golden Markdown normalization and lossless semantic round-trip tests;
- index/projection determinism, freshness, deletion/rebuild, and conflict tests;
- temporary-Git integration tests with fake providers and provider-receipt boundaries;
- full Mission tests from preparation through cold resume;
- migration-capsule fixtures for ambiguity, idempotence, candidate atomicity, rollback, and removal;
- end-to-end acceptance for cold recovery, fuzzy intent to bounded work, and evaluated/reconciled
  return with fresh resume; and
- fixed model/Skill evaluations for correctness, retrieval/context cost, latency, retries, and
  authority compliance.

The first vertical proof is: canonical Markdown and relationships; deterministic index and narrow
views; Proposal → Mission → Run → Evidence → assessment → reconciliation → cold resume; Guardrails
plus a fake provider receipt; then representative v1 conversion with rollback and capsule-free v2
recovery.

## Reconsideration gates

Reconsider the architecture before S12A implementation authorization if it duplicates domain
invariants, makes a cache/projection authoritative, requires a user runtime or provider, restores
legacy parsing/aliases/fallbacks in v2 core, cannot round-trip canonical Markdown, cannot derive
the public mechanical interface from one registry, weakens authority/evidence/cold recovery,
misstates migration information, cannot prove rollback/capsule removal, repeats repair without a
new hypothesis/evidence/narrower action, or fails supported macOS/Linux distribution, startup, or
recovery criteria.

## Deferred Type-2 details

S12A or implementation decides package/file layout, Go/Markdown/YAML/UUIDv7 libraries, normalized
fingerprint exactness, per-platform locking and atomic replacement, projection serialization,
whether and how to persist/shard a cache, cache benchmarks, command flags within accepted grammar,
adapter protocol encoding, CI details, Windows delivery, property-test tooling, fixtures,
performance budgets, and capsule signing/release naming.
