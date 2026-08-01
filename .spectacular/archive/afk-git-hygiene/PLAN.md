---
status: archived
priority: medium
owner: alex
updated: 2026-08-02
build: b34
summary: "Add human-gated AFK branch conventions, spike playground hygiene, archival cleanup, and verified pull-request handoff"
related:
  - PRD.md
depends-on:
  - wayfinding-contract
  - wayfinding-sequencer
archived: 2026-08-02
---

# Plan — afk-git-hygiene

<!--
  Canonical 7-slot PLAN template for a single request.
  Lives at: .spectacular/requests/<slug>/PLAN.md

  Rules:
  - PLAN is per-request. PRD is project-wide. Never put a PRD inside requests/.
  - This file's frontmatter `status:` is the single source of lifecycle state for the request.
  - The 7 required sections must appear IN ORDER, unnumbered:
      ## Goal, ## Constraints, ## Milestones, ## Tasks, ## Dependencies, ## Validation, ## Deliverables
    Extra sections (## Understanding, ## Decisions, request-specific) may appear
    BETWEEN them; doctor enforces the required set's presence + order, not a closed list.
  - All 7 required sections must be filled before this PLAN is considered usable.
  - Replace every <placeholder> with concrete content.
-->

## Goal

<!-- One sentence. What does this request change? -->
<!-- Compress the request's intent. Aligns with PRD's Vision or Goals — this is a slice, not a restatement. -->

Let authorized AFK discovery and execution use isolated, recoverable Git branches and hand verified work to humans without risking primary branches or deleting evidence.

## Constraints

<!-- What's fixed before you start? Inherited from PRD/STACK/PRINCIPLES + request-specific limits. -->

- Branch creation is proposed by default; autonomous creation requires explicit AFK authorization or project configuration.
- Never commit failing or unverified spike code to a primary branch.
- Merge, breaking API/schema changes, remote branch deletion, and product/business question resolution remain HITL gates.
- Archive is preferred to deletion; cleanup is dry-run first and records recoverability.
- Branch naming must integrate with host-agent branch-prefix requirements rather than overriding them blindly.

<!-- ## Understanding and ## Decisions below are OPTIONAL extra sections,
     allowed between Constraints and Milestones. -->

## Understanding

<!--
  OPTIONAL authoring slot, but REQUIRED before `planned → active` by the
  `understand-before-change` policy (@Implementation). Fill it here for a
  typical request; escalate to a dedicated requests/<slug>/UNDERSTANDING.md
  (same three subheads) for large ones — the policy is satisfied by EITHER.
  Not one of the 7 required authoring slots; it gates implementation, not planning.
-->

### How it works now

Spectacular deliberately avoids branch/worktree orchestration and assumes one mutator. It has no AFK authorization contract, spike isolation convention, cleanup ledger, or verified PR handoff.

### What changes

Add an opt-in Git operating policy for draft specs, spikes, forks, and confirmed execution; provide guarded setup/status/cleanup/PR handoff flows; archive abandoned evidence before branch deletion.

### What stays the same

Default interactive work stays branch-agnostic, Spectacular does not replace Git, the orchestrator remains the only lifecycle mutator, and no primary-branch merge happens without the human.

## Decisions

<!--
  Design calls made inside this request. Format: chose X over Y — because Z.
  Rejected alternatives stay listed; deleting them re-litigates them later.
  Project-wide calls go to DECISIONS.md via `spectacular decide` instead
  (see decisions-rules.md routing table). Empty is fine — no decisions yet.
-->

- Chose propose-by-default with explicit AFK authorization over unconditional branch mutation because repositories have host-specific policies and user-owned worktrees.
- Chose archive-before-delete over automatic cleanup because spike evidence and abandoned alternatives may remain valuable.
- Chose separate branch classes (`spec/draft-*`, `spike/prototype-*`, `fork/idea-*`, `feat/v*-*`) over one generic AFK prefix because intent determines verification and cleanup behavior.

## Milestones

<!-- Ordered, demoable checkpoints. Outcomes, not tasks. -->
<!-- 3-7 milestones for a typical request. Each is something someone can see working. -->

- M1 — A project can inspect and opt into AFK Git policy without changing its repository.
- M2 — Draft specs, spikes, forks, and confirmed execution receive safe branch proposals and clean-worktree preflight.
- M3 — Completed or abandoned playground work archives evidence and produces a dry-run cleanup plan.
- M4 — Verified execution can open a consistently titled human-review PR while merge and breaking-change gates remain closed.

## Tasks

<!-- Pointer. The executable checklist lives in TASKS.md, grouped by milestone. -->

See `TASKS.md`.

## Dependencies

<!-- Other requests, skills, blocking decisions. Use [[request-slug]] notation. -->

- Hard dependency on [[wayfinding-contract]] and [[wayfinding-sequencer]].
- Host Git provider/agent branch rules and existing verification artifacts.

## Validation

<!--
  How each milestone is verified. Per-milestone checks.
  Each check states its AUTHORITY: a run: command, an assertable property,
  a judgable artifact, or a human-observable behavior (see verify.md kinds).
  A check with no authority can't fail. Aspiration verbs (improve, enhance,
  optimize, handle gracefully) are not checks.
-->

- M1 — assert: default command is read-only and prints the exact authorization/config needed before mutation.
- M2 — run: repository fixtures reject dirty/primary-branch spike writes and emit policy-compatible names.
- M3 — run: cleanup defaults to dry-run, archives metadata/evidence, and never deletes an unmerged branch without confirmation.
- M4 — run: `bash tests/cli/afk-git-hygiene.test.sh` proves the exact PR title, verified-spec/test gates, no merge call, and separate breaking-change approval.

## Deliverables

<!-- Artifacts that ship out of this request. Concrete files, docs, behaviors. -->

- AFK Git policy contract and configuration schema.
- Branch proposal/preflight and playground lifecycle commands.
- Archive-first cleanup ledger and dry-run behavior.
- Verified PR handoff integration and safety tests.
