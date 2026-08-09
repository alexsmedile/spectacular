---
type: mission-preparation-handoff
schema_version: spectacular.handoff.v2
handoff_id: H26
session: M1-preparation
mode: read-only-interactive-owner-decision
authority: central-orchestration
status: authorized-for-dispatch
content_baseline_commit: 553817d11e02b37564e599462bfae9140b44ba58
content_baseline_tree: b8dfd3b5773ef429e30898623c7503e2a612649d
date: 2026-08-09
---

# H26 — M1 semantic substrate preparation

## Outcome

Prepare the exact activation charter for M1: semantic records plus the canonical Markdown
workspace substrate. Resolve only material remaining choices, define the smallest vertical slice,
compile its authority/effects/evidence/recovery envelope, and obtain the owner's explicit approval.
Do not create `v2/`, write Go code, create an implementation branch, install dependencies, or
activate M1.

## Binding inputs

- `SHARED-SCAFFOLD-CONTRACT.md@1.0` — SHA-256
  `698997f12972d0b5a186f4d8b8c35753a642cc0454e8ef60f15de81590435d36`
- `EXECUTABLE-REFACTOR-PROGRAM-CONTRACT.md@1.3` — SHA-256
  `f9322b26e070882e50eac64dad4e70eb28f2b5441408c635d2fe87f56d718b11`
- `MISSION-PREPARATION-CONTRACT.md@1.0` — SHA-256
  `854d9d2225d954dcb0fd9e93dceb9eb16d3a7e06445ea8653fa82fe2ce6c430b`
- `IMPLEMENTATION-ARCHITECTURE-AND-MIGRATION-CONTRACT.md@1.0` — SHA-256
  `3a5f3a108e805651c4e0e18aeff97f738c6625f42b95cecd6d94c5cde3728fe3`
- `SPECIFICATION-TOPOLOGY-CONTRACT.md@1.0` — SHA-256
  `d8eb174a5f6aea2356e60ba0acc7a9d3084df53a86e6b9c3694102af5c1d1ed1`
- `PRODUCT-TRUTH-CONTRACT-MODEL.md@1.0` — SHA-256
  `e395761f3b9bc0bae0c5816aec95bc5541bac428808b562483719280f8468c8b`
- `WORK-UNIT-LIFECYCLE-CONTRACT.md@1.0` — SHA-256
  `23658bab57bf1551b0ecfbe802be77b8d6401548b3c54d14fcad6c3edb1c205a`

Read `AGENTS.md`, `.spectacular/AGENTS.md`, the binding inputs, H25's reviewed return, and only the
current repository/build surfaces required to prepare a safe new `v2/` implementation boundary.
Do not treat v1 implementation structure as a preservation requirement.

## Required interaction protocol

This is an interactive owner-decision task, not a findings-only investigation.

1. Verify the dispatch baseline and binding hashes.
2. Research the smallest credible technical options using primary repository evidence and official
   upstream documentation only when a dependency or platform fact matters.
3. Brief the owner on the decision map before requesting disposition.
4. Work through the clusters below sequentially. Ask exactly one explicit owner question at a time,
   provide two or three coherent options only where a real choice exists, state consequences and a
   recommendation, and wait for the answer.
5. Do not ask the owner to choose low-level Type-2 details that do not affect product behavior,
   authority, portability, recovery, or the M1/M2 join. Recommend and record those instead.
6. Do not jump directly to a handoff return. The final activation question is mandatory, and the
   owner's approval may not be inferred from earlier W0 or S12B dispositions.

## Decision clusters

### A — M1 slice and acceptance demonstration

Confirm the exact bounded outcome: the minimum Proposal → Mission linked-record slice that proves
identity, canonical Markdown round-trip, typed relationships, fingerprints, deterministic lookup,
and the declared refusal cases. Preserve W0's exclusions and prevent M2/M3 behavior from entering
M1.

### B — Canonical record mechanics

Recommend the minimum mechanically enforceable record envelope, normalization/fingerprint rule,
identity strategy, relationship encoding, and dependency posture. Escalate only choices that are
hard to reverse or alter the accepted authority/recovery contract. Plain Markdown must remain
inspectable and authoritative; indexes and generated output remain disposable.

### C — Implementation and proof envelope

Compile proposed M1 paths/packages, the implementation branch/worktree rule, allowed local effects,
commit permission, forbidden provider effects, focused and integrated checks, independent boundary
review, repair budget, recovery point, and stop conditions. The proposed implementation branch is
`codex/feat/v2-semantic-substrate`; verify that it is available but do not create it.

### D — Final charter approval

Present one compact corrected Mission charter with:

- outcome and non-goals;
- exact owned paths and shared joins;
- authority and effect envelope;
- Objectives or milestones sized for one reviewable Mission;
- evidence and independent-review requirements;
- retry budget, recovery point, stop conditions, and terminal return;
- Design Sufficiency and Slice Quality verdicts.

Ask the owner explicitly to `approve | revise | reject` this charter. Approval authorizes central
orchestration to activate M1 later; it does not itself activate M1.

## Boundaries

This task is read-only. No repository write, branch, commit, package installation, generated output,
provider effect, v1 repair, compatibility logic, migration work, W0 revision, or M1 activation.
Do not reopen accepted foundation or scaffold contracts without named reversal-grade evidence and a
downstream impact audit.

M1 excludes governed lifecycle transitions, authorization, evidence sufficiency, assessment,
reconciliation, public CLI/Skill behavior, Guardrails, real providers, persistent caches, release,
and migration.

## Return

Return `spectacular.handoff-return.v2` with the exact dispatch baseline, verified inputs, read set,
verbatim and normalized owner dispositions, approved M1 charter, technical recommendations and
deferred Type-2 choices, proposed branch/effects, evidence contract, review/retry/recovery/stops,
verdicts, conflicts, and exactly one next action:

`central accept|bounce|escalate; do not activate M1 from this return alone`.
