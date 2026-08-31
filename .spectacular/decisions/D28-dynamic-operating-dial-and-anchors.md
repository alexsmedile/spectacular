---
type: Decision
id: 01a05872-5f05-78c2-a815-62292f017c5d
title: Codify the Dynamic Operating Dial, 5 Foundational Anchors, and Tiered Verification Protocol
created_by: Alex
created: "2026-08-31T15:31:25Z"
updated: "2026-08-31T15:31:25Z"
actor: Alex
actor_role: owner
alternatives:
    - force every mission through heavy multi-agent review bundles
    - rely solely on open-ended prompt loops without structured verification tiers
    - require new custom schema files for every architectural concept
disposition: accepted
question: How should Spectacular balance high-level agentic leverage with low-level direct control across diverse project types?
rationale: Aligning Spectacular with the 20-level Agentic Operating Levels framework and Orca-style orchestration ensures agents ground themselves in hard anchors (types, schemas, state machines) before executing, allows single-agent or multi-agent execution with self-healing test loops, and eliminates token waste through tiered verification and batched reviews.
ref: D28-dynamic-operating-dial-and-anchors
scope:
    - v2
---
# Codify the Dynamic Operating Dial, 5 Foundational Anchors, and Tiered Verification Protocol

This Decision formalizes the synthesis of the Agentic Operating Levels framework, Orca orchestration discipline, and Spectacular v2 governance:

1. **The 5 Foundational Anchors**: Before code execution, missions ground themselves in five core anchors without introducing redundant governance files:
   - *Boundaries & Non-Goals*: Recorded in `.spectacular/PROJECT.md` (`boundaries:`, `constraints:`).
   - *Vocabulary & Domain Ontology*: Recorded in `.spectacular/VOCABULARY.md`.
   - *Invariants & Failure Modes*: Recorded in `.spectacular/GUARDRAILS.md` and `AGENTS.md`.
   - *Data Structures & Schemas*: Kept in project-specific code and type definitions (cited in `.spectacular/contracts/`).
   - *State Machines & Lifecycles*: Visualized as non-governing Mermaid diagrams in `.spectacular/atlas/`.

2. **The Dynamic Operating Dial (`mode:`)**: Missions declare execution posture via frontmatter:
   - `mode: leverage` (Default): High autonomy for familiar/routine tasks. Passing inner-loop test suite (`exit 0`) is primary completion proof.
   - `mode: control`: High precision for out-of-distribution, payments, auth, or schema cutovers. Requires explicit anchor diff checks and dedicated adversarial review.

3. **Tiered Verification Matrix (Zero Duplicate Runs)**:
   - *Tier 1 (Quick / Domain)*: Run by the **Worker** (or solo Orchestrator) on every inner-loop iteration (`verify.sh quick` or domain-scoped unit test).
   - *Tier 0 (Preflight / Lint)*: Run by the **Reviewer** (or automated pre-check gate) to verify tree/AST integrity and contract drift.
   - *Tier 2 (Acceptance)*: Run by the **Orchestrator** at milestone boundaries.
   - *Tier 3 (All / Release)*: Run at the **Owner Gate** prior to final release or mission completion.

4. **Execution Parity & Batched Reviews**:
   - The self-healing execution loop operates identically whether executed solo by the Orchestrator or delegated to a supervised Worker subagent.
   - Routine code tasks avoid per-mission review record sprawl; complex campaigns utilize batched adversarial reviews summarizing multi-mission milestones.
