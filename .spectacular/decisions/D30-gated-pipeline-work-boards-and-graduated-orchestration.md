---
type: Decision
id: 01a07c42-8e10-7e4a-9310-482910fa91b1
title: Adopt Gated Pipeline Orchestration, Non-Governing Work Boards, and Graduated Governance
created_by: Alex
created: "2026-09-03T01:20:00Z"
updated: "2026-09-03T01:20:00Z"
actor: Alex
actor_role: owner
ref: D30-gated-pipeline-work-boards-and-graduated-orchestration
question: How should Spectacular manage multi-step pipelines, side sessions, and worktree isolation without forcing heavy Mission bureaucracy onto routine tasks?
disposition: gated-pipeline-and-graduated-governance
rationale: >-
    Multi-agent software engineering fails primarily due to context fragmentation, premature parallelism,
    and coordination overhead before interfaces stabilize. Spectacular operates as a calm Project Manager:
    multi-step work progresses in sequential dependency waves by default, and parallel side sessions in isolated
    Git worktrees are earned only after upstream contract gates pass. Governance is graduated across four distinct
    tiers (inline, board, brief, mission) to match ceremony directly to risk.
alternatives:
    - force every task into a formal Mission envelope, which introduces friction and token overhead for routine work
    - unconstrained multi-agent swarm dispatch, which causes merge conflicts, interface hallucinations, and lock collisions
authority_basis: Owner explicitly approved the 4-tier governance ladder (inline, board, brief, mission), non-governing Work Boards, Dispatch Brief defaults, and conservative worktree cleanup rules.
authorized_effects:
    - prompt.skill-guidance-evolution
    - docs.process-update
scope:
    - v2
supersedes: ""
---

# Adopt Gated Pipeline Orchestration, Non-Governing Work Boards, and Graduated Governance

## Context & Problem
Multi-turn software development with agents requires balancing speed for routine code against rigorous containment for high-consequence changes. Tying all structured execution strictly to formal Mission bundles (`M<N>.md` / `Handoff`) created unnecessary ceremony for everyday multi-step features. Conversely, unconstrained swarm dispatch caused merge collisions, context dilution, and premature parallel coding before interface boundaries stabilized.

## Decision

1. **The 4 Graduated Governance Tiers (`governance:`)**:
   - **Tier 0 (`governance: inline`)**: Fast, single-context edits directly in the Lead session. Zero governance records.
   - **Tier 1 (`governance: board`)**: Gated multi-step pipelines tracked via non-governing `type: WorkBoard` projections.
   - **Tier 2 (`governance: brief`)**: Bounded side sessions in dedicated `linked-worktree` checkouts, guided by plain-English Dispatch Briefs and tracked by standalone Session records.
   - **Tier 3 (`governance: mission`)**: Full immutable contracts, formal Reviews, Evidence packages, and mechanical Charters for high-stakes, irreversible, or multi-party milestones.

2. **The Gated Wave Invariant (Sequential by Default)**:
   - Work progresses in sequential dependency waves in the Lead checkout unless tasks have separate inputs, disjoint write scopes, and locked upstream interface contracts.
   - A side session is dispatched only when its prerequisite interface gate has locked and its output is independently testable.

3. **Workspace Physical Modes**:
   - `lead-checkout`: Primary working tree for the Lead session and sequential steps.
   - `linked-worktree`: Physically isolated Git worktree (`.worktrees/<slug>`) on a dedicated branch for a single writer.
   - `sandbox`: Isolated container or disposable branch with zero merge authority.
   - `read-only`: Non-mutating scout or auditor thread inspecting diffs.

4. **Session Lifecycle & "Returned ≠ Done"**:
   - Side sessions follow: `planned → ready → active → blocked | returned → integrated → verified`.
   - A side session never marks an item done; it emits a Return Receipt containing commit SHA, test results, and diff stats.
   - The Lead Orchestrator holds sole accountability for reviewing diffs, integrating branches, and executing project-wide verification suites.
   - Worktrees are pruned conservatively only after integration or deliberate abort confirmation.

5. **Lease & Lock Expiration**:
   - Write reservations held by active sessions expire upon explicit `returned`/`aborted` receipt or after a configured heartbeat timeout. Silent crashes never permanently lock repository paths.
