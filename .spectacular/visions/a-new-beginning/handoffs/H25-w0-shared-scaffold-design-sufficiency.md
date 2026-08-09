---
type: decision-session-handoff
schema_version: spectacular.handoff.v2
handoff_id: H25
session: W0
mode: read-only-interactive-owner-decision
authority: central-orchestration
status: authorized-for-dispatch
baseline_commit: e3e35d289c838183cb9cabfebd5616c16b9aac32
baseline_tree: 15fcb8a588bdfdf9421e845aee1d8b71f4127221
date: 2026-08-09
---

# H25 — W0 shared-scaffold Design Sufficiency gate

## Outcome

Obtain the smallest owner-approved shared scaffold that makes M1–M4 joinable without duplicating
authority: repository/build boundary, Go package ownership, canonical/generated surface ownership,
and the first vertical proof seam. Return Design Sufficiency and Slice Quality verdicts. Do not
write code or activate M1.

## Binding inputs

- `EXECUTABLE-REFACTOR-PROGRAM-CONTRACT.md@1.2` — SHA-256
  `79ad885989b0256e1e713ac335c5053254138d0a7553398427d02da77265a8d3`
- `V1-DEPRIORITIZATION-DECISION.md@1.0` — SHA-256
  `53508a986b50b3c9fff99a0641494d697be9f7b1d17b32e650c0dd3952b1b800`
- `IMPLEMENTATION-ARCHITECTURE-AND-MIGRATION-CONTRACT.md@1.0` — SHA-256
  `3a5f3a108e805651c4e0e18aeff97f738c6625f42b95cecd6d94c5cde3728fe3`
- `SPECIFICATION-TOPOLOGY-CONTRACT.md@1.0` — SHA-256
  `d8eb174a5f6aea2356e60ba0acc7a9d3084df53a86e6b9c3694102af5c1d1ed1`
- `MISSION-PREPARATION-CONTRACT.md@1.0` — SHA-256
  `854d9d2225d954dcb0fd9e93dceb9eb16d3a7e06445ea8653fa82fe2ce6c430b`

Read `AGENTS.md`, `.spectacular/AGENTS.md`, the five inputs, the S10/S11/S12A returns, and only the
current source/build paths needed to test scaffold feasibility. Do not load v1 collections for
feature preservation or migration design.

## Decision clusters

Run sequentially. For each cluster, inspect primary evidence, present two or three coherent options
with consequences and one recommendation, ask exactly one explicit owner-decision question, and
stop for the owner's answer. Preserve verbatim and normalized dispositions. Do not infer approval.

### A — Repository and build boundary

Decide the minimal location/build graph for the Go v2 core, Skill source, generated interface, and
tests. V1 code remains outside the v2 dependency and retrieval path; no migration capsule is part
of W0.

### B — Package and authority ownership

Assign one owner to semantic records/canonical workspace, governed operations, deterministic
index/projections, command registry, Guardrails, adapters, and guided Skill workflows. Prefer deep
modules and reject duplicate invariants or a generic framework layer.

### C — Generated surfaces and joins

Fix which source generates CLI dispatch/help/effect metadata and which source supplies Skill-facing
guided semantics. Define the narrow typed joins between M1, M2, and M3 so later work cannot silently
redefine upstream authority.

### D — First vertical proof and verdicts

Choose the smallest M1 proof scaffold that exercises identity + canonical Markdown round trip +
typed relationship + deterministic lookup without implementing M2/M3 behavior. Return:

- Design Sufficiency: `sufficient | needs-evidence | needs-decision`
- Slice Quality: `coherent | too-broad | fragmented | dependency-bound`

## Boundaries

This session is read-only. No files, branch, commit, code, test edits, package installation,
provider effects, migration, v1 repair, P0 work, W0 implementation, or M1 activation. Exact library
selection and low-level filenames remain Type-2 unless they alter a shared authority boundary.

Do not reopen accepted product, semantic, authority, retrieval, public-language, implementation,
or specification contracts without named reversal-grade evidence and an impact audit.

## Return

Return `spectacular.handoff-return.v2` with exact baseline, read set, owner dispositions, approved
scaffold/ownership map, generated-surface contracts, M1 proof seam, verdicts, remaining Type-2
choices, conflicts, evidence strength, and exactly one next action:

`central accept|bounce|escalate; do not activate M1`.
