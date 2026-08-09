---
type: decision-session-handoff
schema_version: spectacular.handoff.v2
handoff_id: H18
session: S11
mode: read-only-interactive-owner-decision
authority: central-orchestration
status: authorized-for-dispatch
expected_parent_commit: d0ff63ceeae9304c0e07e24875002d5b193eb739
expected_parent_tree: 56b5184bd8b7e258bc3447ef496fe4b7826bd547
date: 2026-08-09
---

# H18 — S11 implementation architecture and migration strategy

## Mission

Produce the owner-approved target implementation architecture for clean-break Spectacular v2 and
the isolated v1→v2 migration capsule. Decide how the thirteen accepted contracts become a small,
deep, testable, distributable system without restoring retired v1 structures for convenience.

This is a **decision session**, not an implementation Mission. It must inspect evidence, generate
bounded alternatives, recommend, grill the owner one cluster at a time, and return the accepted
architecture. It may not write product code, migrate data, delete legacy material, create a branch,
commit, open a PR, authorize S12A, or reconcile the central program.

## Baseline gate

The launch prompt supplies the exact dispatch commit and tree containing this handoff. Before any
analysis:

1. verify `git rev-parse HEAD` and `git rev-parse HEAD^{tree}` against that launch baseline;
2. verify the checkout is detached and clean;
3. verify this handoff against its `.sha256` sidecar;
4. verify all thirteen accepted contract files against their sidecars;
5. stop and report the mismatch if any check fails.

The parent before this handoff was commit
`d0ff63ceeae9304c0e07e24875002d5b193eb739`, tree
`56b5184bd8b7e258bc3447ef496fe4b7826bd547`.

## Immutable inputs

| Contract | SHA-256 |
|---|---|
| PRODUCT-CONSTITUTION@1.0 | `99565c58316c4c193fe6108b514b04f664bdee966a840ad2e982ecf580e7dab7` |
| TRUTH-PROVENANCE-FLOOR@1.0 | `49aa9cc1ff2817254bdc542a12c5ed079ff003e2f6f16325822bb1c3387c69c3` |
| SUCCESS-EVIDENCE-CONSTITUTION@1.0 | `a262dcd1a97b35d967b1ec2468d4e0c624b24e352e1b382b1a4b2705898f5547` |
| PRODUCT-TRUTH-CONTRACT-MODEL@1.0 | `e395761f3b9bc0bae0c5816aec95bc5541bac428808b562483719280f8468c8b` |
| WORK-UNIT-LIFECYCLE-CONTRACT@1.0 | `23658bab57bf1551b0ecfbe802be77b8d6401548b3c54d14fcad6c3edb1c205a` |
| EXECUTION-AUTHORITY-CONTRACT@1.0 | `99597d757f5b6cf6293f499e5115ccea1c6eac4e1ec5e684ba17b2f39e4cafc8` |
| EVIDENCE-CLOSURE-CONTINUITY-CONTRACT@1.0 | `7dd763b4fa1a919924e24105790382a51414a0a8ee0222178dee7c9224f11ca9` |
| RESPONSIBILITY-PLACEMENT-CONTRACT@1.0 | `b17c092fcadf7bceb528c103f9e17d1d40e373adf1f13faca91d3c2b6f5711e3` |
| RETRIEVAL-AND-EARNED-WORKSPACE-CONTRACT@1.0 | `e438781d7e485604533e4194840d7eae1a3b9b46d5ce565fe28f469186e101fa` |
| PUBLIC-LANGUAGE-AND-INTERFACE-CONTRACT@1.0 | `9e2b441197e4c30cebca07d3924880872b88d918c582c6d49674ca7bd8dc8e72` |
| CLEAN-BREAK-CUTOVER-AND-RECOVERY-CONTRACT@1.0 | `ec325aca57004d16e8aa6d8b2729c890aa6781893f5f0308441ad179c1a47266` |
| MISSION-PREPARATION-CONTRACT@1.0 | `854d9d2225d954dcb0fd9e93dceb9eb16d3a7e06445ea8653fa82fe2ce6c430b` |
| SUBSYSTEM-SURVIVAL-CONTRACT@1.0 | `abdfc0d684cddcc8a89290893b6e555d43e727f849d1f9af72338a26e46963c4` |

## Required read set

Read progressively, in this order:

1. `AGENTS.md` and `.spectacular/AGENTS.md`;
2. this handoff and its sidecar;
3. the thirteen accepted contracts and sidecars;
4. `evidence/decision-sessions.md` S11 and `evidence/synthesis-027.md`;
5. `evidence/concepts/PZL-031*`, `PZL-035*`, `PZL-042*`, `PZL-046*` through
   `PZL-060*`, `PZL-102*`, `PZL-110*`, `PZL-111*`, `PZL-121*`, `PZL-125*`,
   `PZL-170*`, and `PZL-171*`;
6. only conflict rows and source/synthesis passages directly implicated by those cards;
7. current CLI, installer, plugin manifests, tests, hooks, Skill routing, and workspace files only
   as needed to establish architecture, distribution, maintenance, and migration facts.

Do not preload the archive or the entire evidence workbench. Current implementation is evidence,
not the target architecture.

## Locked constraints

- v2 supports v2 workspaces only. Core has no legacy parser, alias, fallback read, dual write, lazy
  conversion, or in-tree legacy archive.
- v1 freezes at an immutable release/tag. Migration is an explicit, whole-project atomic
  transaction through a removable capsule that never contaminates v2 core.
- Markdown/Git-native canonical state remains inspectable and repairable. Generated indexes and
  views are non-authoritative.
- Skill owns judgment; deterministic code owns mechanical validation, references, projections,
  authorized transitions, and scaffolding.
- Host runtimes execute; native providers perform and attest external effects.
- S09's vocabulary and Skill/CLI grammar are fixed inputs, not topics to rename.
- S10's survival decisions are fixed inputs. Do not resurrect v1 collections, fixed POLICY,
  universal doctor repair, alternate VERIFY authority, fixed agent fleet, Wayfinding, or AFK as
  independent subsystems.
- Deep modules and simple interfaces are preferred, but language choice must follow evidence—not a
  slogan or rewrite instinct.
- The distributable product must remain practical for macOS and agent environments. Do not assume
  a new runtime or package manager is acceptable without presenting its installation and recovery
  consequences to the owner.
- No compatibility layer in v2 core. Migration helpers live only in the disposable capsule.

## Interactive protocol — mandatory

This session must not jump directly to a return packet.

For each cluster below:

1. inspect the named evidence;
2. present two or three coherent options and their consequences;
3. identify the recommended option and why;
4. ask one explicit owner decision question;
5. **stop and wait for the owner**;
6. record the owner's verbatim and normalized disposition before advancing.

Ask only decisions that materially change architecture, dependency, distribution, recovery, or
long-term maintainability. Resolve reversible implementation details through a recommended default
and clearly mark them as S12/implementation details. Never infer owner acceptance from silence,
tone, or an earlier adjacent decision.

## Decision clusters

### A — Language, build, and distribution posture

Compare at least:

- a deeply modular Bash 3.2 source assembled into one distributable script;
- a clean implementation in another language/runtime with explicit packaging consequences;
- a justified hybrid where a thin bootstrap/distribution shim invokes a deeper implementation.

Measure current Bash surface, coupling, tests, installer/plugin constraints, startup behavior,
portability, contributor ergonomics, typed-data needs, and recovery. Decide the source language,
release artifact, runtime/dependency promise, and explicit conditions that would invalidate the
choice. “Rewrite because the file is large” and “keep Bash because it already exists” are both
insufficient.

### B — Deep module and adapter architecture

Define a small module map that realizes the accepted responsibilities without mirroring every noun
as a module. Cover:

- semantic/domain core and authority invariants;
- canonical Markdown storage and identity/reference resolution;
- Mission lifecycle, preparation, evidence/assessment, and reconciliation;
- deterministic projections/indexes and scoped integrity;
- command registry and generated help/interface tests;
- guided-authoring Skill/kernel and progressively loaded role contracts;
- provider/runtime adapters;
- strict v2-core versus migration-capsule boundary.

Ask the owner to choose among genuinely different architectural shapes, not filenames.

### C — Canonical representation and deterministic substrate

Decide the minimum machine-readable representation needed for Project Anchors, Contracts,
Proposals, Missions/Runs, embedded/global earned records, Receipts, Guardrails selectors, generated
indexes, and promotion continuity. Compare Markdown/frontmatter-only, a deterministic local index
derived from Markdown, and any stronger substrate only if evidence earns it. Preserve Markdown as
canonical authority and make deletion/rebuild behavior explicit.

Define stable identity versus readable slug/version/fingerprint, query/rebuild boundaries,
freshness/conflict behavior, and what remains a Skill interpretation rather than a schema rule.

### D — Migration capsule, cutover, and recovery

Choose the capsule's architectural boundary, implementation language/dependencies, explicit
invocation, report/validation/receipt outputs, rollback procedure, and removal proof. Define the
final v1 freeze/tag sequence and the evidence required before a project cuts over. Ambiguous
semantic mappings must stop for owner disposition; the capsule never edits accepted v1 in place.

### E — Test architecture and no-rewrite gates

Define tests and evaluation for:

- semantic and lifecycle invariants;
- command registry/help/interface consistency;
- Markdown parsing, identity, links, freshness, projection rebuilds, and integrity;
- Mission preparation, authority, evidence, closure, and cold resume;
- provider adapter contracts without granting provider authority;
- migration fixtures, atomic cutover, rollback, and capsule-removal proof;
- the three accepted rebuilt-product scenarios and representative context/recovery measurements.

State explicit no-rewrite/stop criteria, vertical-slice proof order, and conditions that require the
architecture to be reconsidered before S12A.

## Required output

Return one `spectacular.handoff-return.v2` packet with:

- `handoff_id: H18`, `session: S11`, and `status: complete`;
- verified dispatch commit/tree, clean state, handoff hash, and all thirteen contract hashes;
- exact read set and measurements used;
- owner dispositions, preserving verbatim wording and normalized meaning;
- chosen language/build/distribution contract and rejected alternatives;
- a compact module/adapter diagram and responsibility table;
- canonical representation, identity, index, and rebuild contract;
- Project Guardrails parser/event/command boundary;
- migration capsule, cutover, rollback, recovery, and capsule-removal contract;
- test architecture, executable interface gates, no-rewrite criteria, and vertical-slice proof;
- open Type-2 details explicitly reserved for S12A or implementation;
- facts, assumptions, conflicts, scope deviations, and labeled evidence;
- exactly one next action: central `accept | bounce | escalate`.

The final return must say explicitly: **do not authorize S12A from this return alone**.

## Stop conditions

Stop and ask the owner when:

- a decision would add a mandatory runtime/provider/dependency;
- evidence cannot distinguish architecture options;
- the proposed design conflicts with an accepted contract;
- a choice would reintroduce legacy behavior into v2 core;
- a migration ambiguity could discard or misstate owner truth;
- a destructive or implementation action would be required;
- the owner has not explicitly disposed every Type-1 cluster.
