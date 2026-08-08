---
type: refactor-foundation-contract
version: 1.0
status: accepted
decision_session: compatibility-floor
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
next_session: S10
---

# Clean-Break Cutover and Recovery Contract

This accepted checkpoint defines the v1/v2 support boundary, migration
transaction, recovery guarantee, and retirement preconditions. It authorizes no
specific subsystem deletion and does not select implementation paths, schemas,
languages, or module layout.

## Supported population

Spectacular v2 officially supports v2 workspaces only. Its CLI understands only
v2 structures and vocabulary. Core contains no legacy parser, aliases, fallback
reads, dual writes, compatibility branches, or lazy conversion behavior.

An existing v1 workspace is unsupported by the v2 CLI until it has completed
the accepted project-level migration transaction.

## Frozen-v1 boundary

v1 freezes at a final immutable tagged release. That release remains
downloadable, inspectable, and recoverable, but receives no routine fixes,
features, or compatibility maintenance after v2 releases.

An exceptional v1 security patch requires a separate explicit owner decision.
It creates no implied continuing maintenance commitment.

## Project-level migration transaction

Migration is whole-project and atomic. Mixed v1/v2 workspaces and lazy migration
are unsupported.

```text
v1 snapshot
  → separate v2 candidate
  → unsupported/ambiguous-state report
  → clean-v2 validation
  → owner approval
  → cutover
```

The process never edits the accepted v1 source in place. Ambiguous semantic
mappings stop for owner disposition. Unsupported states are reported explicitly
rather than guessed, discarded, or carried forward under misleading v2
semantics.

The original v1 snapshot remains available for rollback until the migrated
project is explicitly accepted. Failed validation or withheld owner approval
leaves the v1 workspace authoritative and prevents cutover.

## Disposable migration capsule

All v1 migration scripts, reference guidance, legacy fixtures, and
capsule-specific tests live in one isolated migration capsule outside v2 core.
The capsule may use explicitly declared and pinned dependencies.

The v2 core must never import the capsule, invoke it implicitly, or depend on it
or its dependencies. Invocation is explicit and project-specific.

Capsule removal is proven by deleting it from a disposable verification checkout
and running the complete v2 behavior and test suite unchanged. Final v1 and v2
Git tags retain its design and implementation history. Each migrated project
retains its v1 snapshot, migration report, validation evidence, owner decision,
cutover receipt, and rollback pointer.

## S10 retirement authority

S10 may apply evidence-gated retirement; this contract grants no blanket
deletion authority. Before removal, every legacy subsystem and artifact is
classified.

Unique product truth, owner decisions, provenance, and historical evidence must
be migrated, retained, or given a durable recovery pointer. The complete v1
implementation remains recoverable through its immutable release/tag and a
verified workspace snapshot.

Legacy material proven redundant, superseded, or recoverable may be absent from
the live v2 tree. Each removal records:

- what replaced it;
- where unique history remains;
- how recovery was verified.

Obsolete code must not remain in v2 search or retrieval paths merely as an
in-tree archive. No deletion is authorized while unique truth, rollback,
replacement, or recovery evidence remains unresolved.

## Reserved implementation decisions

S11 retains the final v1 tag and release sequence; capsule path, implementation
language, dependencies, and invocation; snapshot/candidate formats; validation
command; migration-report and cutover-receipt schemas; rollback procedure; and
project-evidence retention policy. S10 retains subsystem-by-subsystem survival
and retirement dispositions.
