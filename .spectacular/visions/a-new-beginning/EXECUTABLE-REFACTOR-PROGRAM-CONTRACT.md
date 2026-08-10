---
type: refactor-program-contract
contract: executable-refactor-program
version: 1.8
status: accepted
decision_session: scenario-a-acceptance-and-b-c-join
source_handoff: scenario-a-implementation-return
accepted_by: owner
accepted_at: 2026-08-10
central_disposition: accept
supersedes: _snapshots/EXECUTABLE-REFACTOR-PROGRAM-CONTRACT/@v1.7.md
upstream:
  - MVP-SCENARIO-CLI-SEQUENCING-DECISION.md@1.0
  - M1-SEMANTIC-SUBSTRATE-MISSION-CHARTER.md@1.2
  - SHARED-SCAFFOLD-CONTRACT.md@1.0
  - SUCCESS-EVIDENCE-CONSTITUTION.md@1.0
  - PUBLIC-LANGUAGE-AND-INTERFACE-CONTRACT.md@1.0
  - EXECUTION-AUTHORITY-CONTRACT.md@1.0
  - EVIDENCE-CLOSURE-CONTINUITY-CONTRACT.md@1.0
  - SCENARIO-A-COLD-RECOVERY-MISSION-CHARTER.md@1.0
  - GOVERNED-WORK-AND-CLOSURE-CONTRACT.md@1.0
next_action: prepare-b-c-governed-loop-mission
---

# Executable Refactor Program Contract

## Purpose

This v1.8 program accepts Scenario A, joins its mechanical interface with the owner-approved B+C
behavioral contract, adopts a minimal-agent execution default, and separates Skill/runtime
prerequisites from release/distribution. No successor Mission is activated by this record.

## Accepted implementation baselines

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

Scenario A was accepted after bounded repair and integrated without squash:

- final feature head: `3b25f1777c40809be815b3172edebaf588936e99`;
- final feature tree: `25b5a7458a0a5d7020f2ffeeb8179c9508734bdb`;
- independently reviewed product head: `c9efd998e768ac2ba0cdf871acc5368fb35dae05`;
- independently reviewed product tree: `4a99dbb1ecb0a6363d3cf8744c18d0142920b4d8`.

Central verified all Scenario A sidecars, scope, v2 format/module/vet/test/race/build checks, and the
complete 31-file v1 suite before merge. The acceptance receipt is
[`evidence/scenario-a-central-acceptance.md`](evidence/scenario-a-central-acceptance.md).

## Current program

```text
W0 shared-scaffold gate — accepted
  → M1 semantic records + canonical workspace substrate — accepted
  → A cold recovery — accepted
  → BC governed work through closure — next-ready for preparation
  → S Skill and runtime prerequisites
  → R release and distribution
```

The BC Mission contains serial Objectives B then C and one final acceptance boundary. Scenario C
remains mandatory for MVP. Real project migrations and cutovers remain later
project-specific Missions. V2 core contains no v1 compatibility behavior.

## Scenario A — cold recovery (accepted)

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

## BC Mission — bounded governed work through closure

Objective B outcome: fuzzy intent becomes an owner-accepted Proposal, bounded Mission, explicit
authority envelope, expected evidence, and validated Handoff. Guided `/spectacular propose` and
`/spectacular define` own interpretation; noun-first CLI `create` persists only confirmed records.

Proposal base checks, Handoff validation, and Mission transitions are mechanical. Every transition
requires an authorization Decision reference and expected fingerprint.

Objective C outcome: Evidence returns through assessment, owner disposition, authorized Contract
reconciliation, Mission resolution, archival, and a second cold recovery.

Guided Skill operations retain assessment and reconciliation judgment. Mechanical Evidence and
Decision creation, Contract reconciliation, Mission transition, and archival apply only explicit
authority with matching base identity/fingerprint. Archival never proves closure. The accepted
semantic contract is [`GOVERNED-WORK-AND-CLOSURE-CONTRACT.md`](GOVERNED-WORK-AND-CLOSURE-CONTRACT.md).

Objective B is an evidence checkpoint, not a separate owner or central handoff. If its deterministic
checks pass without a contract/authority conflict, the same primary agent continues into C.

## Scenario S — Skill and runtime prerequisites

Build the guided Skill workflows, progressive context compiler, Mission preparation and Autopilot
charter compilation, runtime-neutral Handoff behavior, manifests, and Codex/Claude installation
prerequisites. Validate them through an isolated clean-v2 Spectacular dogfood workspace. This is
project self-dogfooding, not model/server hosting or migration of the live v1 workspace.

## Scenario R — release and distribution

After BC and S pass end to end, complete macOS/Linux native distribution, installer/update and
checksum verification, version alignment, fresh-install/recovery proof, and v1-runtime exclusion.
Windows, generic migration, provider integrations, persistent caches, and broad visual polish remain
deferred until earned.

## Execution controls

Fixed commit-count ceilings are retired. Each Mission declares exact baseline and reviewed
commit/tree, owned paths, prohibited effects, changed invariants, dependency diff, repair budget,
review trigger, recovery point, and one terminal next action. Commits remain coherent review units,
not success metrics.

One primary agent owns decisions, implementation, tests, repairs, commits, evidence, and return by
default. A fresh reviewer is used when consequence requires independence; a cold verifier is used
only when recovery without chat is itself the claim. Additional agents require demonstrably
parallel, loosely coupled work and an explicit context/join benefit. They are never spawned merely
because a role exists. Central orchestration intervenes on scope, product, authority, provider,
irreversible-effect, or exhausted-boundary conflicts—not ordinary corrections.

Before asking the owner, a decision-delta test must show that no accepted contract answers the
question, the answer materially changes behavior/authority/safety/irreversibility, a reversible
Type-2 default is insufficient, and new evidence makes the question live. Otherwise inherit the
accepted decision and continue.

The live control surface is this program plus the current Mission charter. Immutable constitutional
contracts remain authoritative linked inputs and are not restated as mutable progress systems.

## Mechanical interface boundaries

One internal record/projection engine may implement noun-specific commands, but `Record` is not a
primary public noun. Canonical v2 excludes generic `status`, `inspect`, `new`, `advance`, universal
`doctor`, and bare mechanical `assess`, `decide`, or `resolve`. Use scoped `workspace validate`;
correction remains owned by the responsible operation.

## Next gate

Prepare one BC Mission with serial Objectives B and C. Use a short gap audit instead of a new owner
interview. Ask only if a genuine unresolved Type-1 conflict survives the accepted B+C contract and
Scenario A interface join. Do not activate BC merely by updating this program.
