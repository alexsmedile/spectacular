---
type: MissionPlan
title: Enforce isolated frozen-Handoff dispatch
owner: Alex
contract:
  ref: Contract:01a00aae-8921-7b27-96a9-1a4c175e7dc6
outcome: Orchestrators can dispatch eligible Objectives across active Missions only when frozen Handoffs have disjoint writable reservations and native Git isolation is already safe.
review: independent
completion:
  - claim: frozen-write-reservations
    pass_boundary: A Handoff freezes Mission, Objective, Run, sources, and writes; writes accepts repository-relative exact files or trailing-slash directory subtrees, forbids globs and parent traversal, requires both rename endpoints, and refuses overlap with every reserving Handoff across active Missions.
    proof_requirement: Adversarial path and overlap matrices cover files, nested directories, siblings, new scaffolded files, renames, normalization, symlink escape, paused and blocked reservations, cross-Mission conflicts, and atomic refusal-before-dispatch.
  - claim: native-git-inspection
    pass_boundary: mission start and dispatch inspect that work occurs on a clean non-default branch with no interrupted Git operation; Spectacular never creates, switches, stashes, merges, pushes, removes, or force-removes branches or worktrees.
    proof_requirement: Real temporary repositories cover clean and dirty trees, default and non-default branches, merge/rebase/cherry-pick interruptions, detached or ambiguous state, and command-spy assertions that no forbidden Git effect occurs.
  - claim: preserved-stop-and-cleanup
    pass_boundary: Scope escape or dispatch failure stops with reason and next action while preserving work; clean merged and tested Run or Objective worktrees may be proposed for cleanup, but every removal requires owner confirmation and remaining cleanup is a Mission-closure responsibility.
    proof_requirement: Failure injection proves zero data loss and retained refs for timeout, scope escape, repair exhaustion, conflict, and dirty work; cleanup fixtures prove no proposal before merge and tests and no effect without owner confirmation.
  - claim: project-surface-transition
    pass_boundary: The already versioned CC-projsurf behavior for concurrent Objective Runs and write reservations is implemented while its concurrent-run-timelines Gap remains explicitly open for M20 proof and M21 closure; CC-v2prod is versioned for frozen-Handoff isolation, and P8 is retired only after all of its isolation questions are demonstrably answered.
    proof_requirement: Contract and Proposal lineage review finds no silent Gap deletion, no completed binding rewrite, and correct D11 retirement metadata only after accepted implementation Evidence.
objectives:
  - outcome: Freeze and enforce exact writable reservations across active Missions.
    claims: [frozen-write-reservations]
  - outcome: Inspect native Git isolation and preserve every refusal and stop path.
    claims: [native-git-inspection, preserved-stop-and-cleanup]
  - outcome: Version product Handoff behavior and retire answered isolation planning.
    claims: [project-surface-transition]
authority:
  operator: [inspect, edit-in-scope, choose-reversible-implementation, run-checks, generate-derived-files, bounded-repair, commit-local]
  requires_owner: [activate-mission, change-outcome-or-completion, expand-scope, push, merge, release, irreversible-change, destructive-data]
scope:
  mechanical: [internal/runtime, internal/missionbundle, internal/command, skills/spectacular, test, .spectacular/contracts, .spectacular/proposals, .spectacular/archive/proposals]
  semantic: [frozen Handoff write reservations, cross-Mission dispatch eligibility, native-Git inspection, owner-confirmed cleanup, P8 retirement]
repair_budget: 2
dependencies: [M18 completed with Objective-scoped Run reservations]
gaps: []
stops: [git-mutation-by-spectacular, write-reservation-overlap, scope-escape, dirty-or-unmerged-cleanup, unconfirmed-cleanup, data-loss, silent-gap-closure]
---

# Mission

> **Future Mission sketch.** Preserve as design input. Do not activate, maintain,
> validate, or review as a final plan until its predecessor closes and the
> Orchestrator re-prepares this block from current Evidence.

Concurrency is a planning result, not an autonomous scheduler. There is no hard
numeric ceiling; explicit dependencies and frozen disjoint writes decide eligibility.
