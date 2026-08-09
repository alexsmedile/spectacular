---
type: decision-session-handoff
schema_version: spectacular.handoff.v2
handoff_id: H19
session: S12A
mode: read-only-interactive-owner-decision
authority: central-orchestration
status: authorized-for-dispatch
expected_parent_commit: 149e5ef76e8e494d92b138620a21e2b7c5a2d7de
expected_parent_tree: 7e81ac5af8c6ea149c643070116c6411409bc9aa
date: 2026-08-09
---

# H19 — S12A specification topology and approval

## Mission

Produce the smallest coherent, owner-approved specification set for the v2 rebuild. It must turn
the fourteen accepted foundation contracts into bounded, independently reviewable specifications
with clear acceptance/proof, dependencies, material Gaps, and a skeptical check. It must not turn
those specifications into implementation Missions, code, migration, deletion, or release work.

This is a **decision session**, not an implementation Mission. It investigates, proposes high-
density alternatives, grills the owner one cluster at a time, and returns the chosen specification
topology. It may not edit files, create a branch/commit/PR, use a provider, create a Mission,
authorize S12B, or reconcile the central program.

## Baseline gate

The launch prompt supplies the exact dispatch commit and tree containing this handoff. Before any
analysis:

1. verify `git rev-parse HEAD` and `git rev-parse HEAD^{tree}` against that launch baseline;
2. verify the checkout is detached and clean;
3. verify this handoff against its `.sha256` sidecar;
4. verify all fourteen accepted contracts against their sidecars;
5. stop and report any mismatch.

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
| IMPLEMENTATION-ARCHITECTURE-AND-MIGRATION-CONTRACT@1.0 | `3a5f3a108e805651c4e0e18aeff97f738c6625f42b95cecd6d94c5cde3728fe3` |

## Required progressive read set

1. `AGENTS.md`, `.spectacular/AGENTS.md`, this handoff, and all sidecars;
2. all fourteen accepted contracts, but only their sections relevant to each cluster;
3. `evidence/decision-sessions.md` S12A, `FOUNDATION-PLAN.md`, and `ORCHESTRATION.md`;
4. PZL-023, 029–033, 036, 078–079, 102, and 121, plus directly implicated conflicts/syntheses;
5. the accepted architecture's test/proof and migration boundaries;
6. only relevant current source/templates/tests when needed to expose an implementation constraint
   or candidate legacy ambiguity. Current v1 artifacts are evidence, never a v2 spec boundary.

Never preload the archive or the full evidence workbench.

## Locked constraints

- A Capability Contract remains current behavioral authority; a Proposal remains a candidate
  change; a Mission remains the sole accountable delivery unit. A specification cannot become a
  competing lifecycle or current-truth authority.
- v2 is a clean break: no core compatibility parser, aliases, fallbacks, dual writes, lazy
  conversion, historical archive tree, or resident migration registry.
- Go-native architecture, Markdown authority, non-authoritative projections, the command registry,
  Guardrails boundary, and isolated migration capsule are accepted architecture—not open options.
- S12A produces no implementation Mission/request. Only S12B compiles approved specifications into
  Missions after a separate central authorization.
- Prefer a few deep specifications. A spec survives only if it owns a distinct acceptance question,
  implementation boundary, change cadence, or risk/proof surface that simpler sections cannot
  safely carry.
- Do not reintroduce retired v1 collections, universal document engines, fixed agent fleets,
  generic scheduler, broad doctor repair, fixed POLICY, Wayfinding, or AFK as independent systems.
- A missing material fact becomes an explicit Gap/assumption/decision gate; do not invent it.

## Interactive protocol — mandatory

Do not jump to a final specification set or a return packet.

For each cluster:

1. inspect the stated evidence;
2. present two or three coherent alternatives with consequences;
3. recommend one with explicit rationale and reversible/deferred details;
4. ask **one** explicit owner-decision question;
5. **stop and wait** for Alex;
6. record the verbatim answer and normalized disposition before advancing.

The owner decides boundaries, priorities, promised behavior, risk acceptance, and spec approval.
Do not infer acceptance from silence or from an implementation preference. You may propose prose
and small diagrams, but state what remains a recommendation.

## Decision clusters

### A — Topology: the smallest set of specs

Map all accepted contracts to candidate product specifications. Compare at least two topology
options—for example, a small vertical-slice set versus a larger concern-by-concern decomposition.
For every candidate, state the unique acceptance question it owns, why it cannot be only a section
of another spec, what it deliberately excludes, and its primary dependencies. Use a compact matrix.

Ask the owner to select or modify the recommended specification set. Stop.

### B — Contract and interface boundaries

For the accepted topology, draft a one-page skeleton for each specification: outcome, non-goals,
affected accepted contracts, promised/forbidden behavior, material interfaces and failure modes,
continuity/authority boundaries, migration implications, and proof obligations. Identify whether
the specification is foundational, vertical, or a dedicated migration/recovery concern.

Show a typed dependency/acceptance map. A relationship is not a scheduler. Ask the owner to approve
or change the boundaries and order. Stop.

### C — Proof, risk, and design sufficiency

For each proposed spec, define the minimal acceptance tests, evidence, independent-review triggers,
and material Gaps/assumptions. Apply the accepted Design Sufficiency and Slice Quality verdicts;
do not create a separate design lifecycle. Ensure capability/authority/evidence/reconciliation/cold
resume/migration claims have appropriate proof without repeating all foundation text.

Ask the owner to accept the proof/risk posture or resolve one material Gap. Stop.

### D — Skeptical review and explicit approval

Run a fresh, adversarial self-review against all accepted contracts: duplication, hidden competing
authority, weak recovery, provider capture, compatibility leakage, unbounded scope, insufficient
proof, and spec dependency conflict. Present only material findings, fixes, and any required owner
choice. Do not grade prior authoring conversationally; inspect the proposed artifacts/claims.

Ask the owner whether to approve each revised specification, bounce specific ones, or escalate an
unresolved Type-1 conflict. Stop.

## Required return

Only after every cluster has an explicit owner disposition, return a complete
`spectacular.handoff-return.v2` packet containing:

- H19/S12A status, exact baseline/tree, clean-state result, and all input-hash verification;
- owner dispositions verbatim and normalized by cluster;
- each proposed specification with unique owned question, scope/non-goals, dependencies,
  acceptance/proof, material constraints/Gaps, and approval state;
- the typed dependency and acceptance topology, with no Mission/request;
- skeptical-review findings and their resolution/escalation;
- a traceability table from accepted contracts to specifications;
- explicit unresolved Type-2 implementation choices; and
- `next_action: central accept|bounce|escalate; do not authorize S12B`.

Do not write canonical specs or return artifacts in this read-only session. Central orchestration
alone accepts, records, and authorizes subsequent work.
