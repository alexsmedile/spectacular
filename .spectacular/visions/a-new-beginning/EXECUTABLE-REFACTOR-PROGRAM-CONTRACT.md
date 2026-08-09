---
type: refactor-program-contract
contract: executable-refactor-program
version: 1.7
status: accepted
decision_session: m1-acceptance-and-scenario-sequencing
source_handoff: H28-final-closure
accepted_by: owner
accepted_at: 2026-08-09
central_disposition: accept
supersedes: _snapshots/EXECUTABLE-REFACTOR-PROGRAM-CONTRACT/@v1.6.md
upstream:
  - MVP-SCENARIO-CLI-SEQUENCING-DECISION.md@1.0
  - M1-SEMANTIC-SUBSTRATE-MISSION-CHARTER.md@1.2
  - SHARED-SCAFFOLD-CONTRACT.md@1.0
  - SUCCESS-EVIDENCE-CONSTITUTION.md@1.0
  - PUBLIC-LANGUAGE-AND-INTERFACE-CONTRACT.md@1.0
  - EXECUTION-AUTHORITY-CONTRACT.md@1.0
  - EVIDENCE-CLOSURE-CONTINUITY-CONTRACT.md@1.0
next_action: prepare-scenario-a-cold-recovery-mission
---

# Executable Refactor Program Contract

## Purpose

This v1.7 program records M1 acceptance and integration, then drives the MVP through the three
constitutional scenarios before release hardening. No later Mission is activated by this record.
Scenario A preparation is the sole next-ready action.

## Accepted implementation baseline

M1 was independently accepted and integrated without squash:

- reviewed feature head: `a488b2efe7828f59724f730b9a590b9a644e6885`;
- reviewed feature tree: `6cce1468c13e51ef007b0e23ed6b5295cdefd87b`;
- central merge commit: `7759461a34ec98a21c6b6e7449ecc0c13a2a87aa`;
- central merge tree: `5c25177b50528e6eb57037700535204170bcb539`.

The merged central tree passed formatting, module verification, vet, race tests, full tests, and
build with `GOFLAGS=-mod=readonly` and a disposable build cache.

M1 owns only semantic Proposal/Mission records, canonical tree-shaped Markdown/YAML persistence,
UUIDv7 identity, typed references, fingerprints, deterministic indexing, and refusal boundaries.
It contains no CLI, registry, workspace discovery, lifecycle engine, reconciliation, or provider
behavior.

## Current program

```text
W0 shared-scaffold gate — accepted
  → M1 semantic records + canonical workspace substrate — accepted and integrated
  → A cold recovery
  → B fuzzy intent to bounded governed work
  → C evidence, disposition, reconciliation, and cold resume
  → R release and distribution hardening
```

Scenario C is mandatory for MVP. Real project migrations and cutovers remain later
project-specific Missions. V2 core contains no v1 compatibility behavior.

## Scenario A — cold recovery

Outcome: a cold person or agent can orient to a project, discover current Missions, resume one
Mission safely, inspect every consequential pointer, and reach one mechanically justified
continuation or the exact unresolved owner gate.

First mechanical slice:

```text
spectacular anchor show project [--json]
spectacular mission list [--json]
spectacular mission show <ref> [--json]
```

Required drill-down uses noun-first `gap`, `run`, `checkpoint`, `evidence`, and `decision`
`list`/`show` operations plus scoped `workspace validate`. `anchor show` must separate authoritative
Anchors, generated projection, Gaps/conflicts, continuation or owner gate, sources, freshness, and
generation basis. It never infers safety from a stored label alone.

## Scenario B — bounded governed work

Outcome: fuzzy intent becomes an owner-accepted Proposal, bounded Mission, explicit authority
envelope, expected evidence, and validated Handoff. Guided `/spectacular propose` and
`/spectacular define` own interpretation; noun-first CLI `create` persists only confirmed records.

Proposal base checks, Handoff validation, and Mission transitions are mechanical. Every transition
requires an authorization Decision reference and expected fingerprint.

## Scenario C — closure and fresh resume

Outcome: Evidence returns through assessment, owner disposition, authorized Contract
reconciliation, Mission resolution, archival, and a second cold recovery.

Guided Skill operations retain assessment and reconciliation judgment. Mechanical Evidence and
Decision creation, Contract reconciliation, Mission transition, and archival apply only explicit
authority with matching base identity/fingerprint. Archival never proves closure.

## Scenario R — release hardening

After A–C pass end to end, complete generated registry/help, Skill packaging, macOS/Linux native
distribution, install/checksum flow, integration/race/recovery proof, self-hosting, and v1-runtime
exclusion. Windows, generic migration, provider integrations, persistent caches, and broad visual
polish remain deferred until earned.

## Execution controls

Fixed commit-count ceilings are retired. Each Mission declares exact baseline and reviewed
commit/tree, owned paths, prohibited effects, changed invariants, dependency diff, repair budget,
review trigger, recovery point, and one terminal next action. Commits remain coherent review units,
not success metrics.

Builder and independent reviewer may exchange bounded repair/re-review messages directly while the
Mission envelope remains unchanged. Central orchestration intervenes on scope, product, authority,
provider, irreversible-effect, or exhausted-boundary conflicts—not on ordinary corrections.

The live control surface is this program plus the current Mission charter. Immutable constitutional
contracts remain authoritative linked inputs and are not restated as mutable progress systems.

## Mechanical interface boundaries

One internal record/projection engine may implement noun-specific commands, but `Record` is not a
primary public noun. Canonical v2 excludes generic `status`, `inspect`, `new`, `advance`, universal
`doctor`, and bare mechanical `assess`, `decide`, or `resolve`. Use scoped `workspace validate`;
correction remains owned by the responsible operation.

## Next gate

Prepare Scenario A only. Its charter must define its smallest vertical slice, exact files and
interfaces, read-only guarantees, JSON schema, source drill-down, cold-actor test, self-hosted test,
authority/refusal rules, and integration boundary with the accepted M1 substrate.
