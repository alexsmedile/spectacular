---
type: Decision
id: 01a029c5-9295-7ffb-84c5-6b6c095d3df2
title: Adopt Git worktree isolation and bounded context-sandwich compilation for async execution
created_by: Alex
created: "2026-08-22T16:00:10Z"
updated: "2026-08-22T17:25:00Z"
actor: Alex
actor_role: owner
ref: D12-isolation-and-context-compilation
question: How should concurrent and asynchronous execution boundaries and context envelopes be structured in Spectacular?
disposition: worktree-isolation-and-bounded-charters
rationale: >-
    To enable high-velocity human steering where an owner can live-prompt and settle new
    decisions while background subagents build confirmed slices, execution boundaries must be
    physically isolated and context envelopes strictly token-bounded. Physical Git worktrees
    (`git worktree add`) eliminate concurrent branch mutation and dirty-tree lock collisions on
    the main working directory. Simultaneously, compiling a 3-layer Context Sandwich (~1,200 tokens)
    per Objective ensures workers receive frozen truth, relevant prior decisions, and exact target
    files without wasting tokens or suffering attention dilution from whole-workspace scans.
alternatives:
    - in-place branch hopping within a single shared working directory, which causes file lock collisions and interrupts live steering
    - dumping unbounded repository context into worker prompts, which causes context bloat and hallucinated scope escapes
authority_basis: Owner explicitly confirmed the physical Git worktree isolation model and bounded on-demand Context Sandwich compilation protocol for P11 exploration and the M15–M17 campaign.
authorized_effects:
    - command.surface-growth
conditions:
    - command-surface-growth-14-to-16
    - demonstrated-context-reduction-gte-40-pct
    - zero-regression-on-paired-benchmark-suite
scope:
    - v2
targets:
    - Proposal:01a029be-b7d3-703c-a7ee-50c6b8bae3a2
supersedes: ""
---

# Adopt Git worktree isolation and bounded context-sandwich compilation for async execution

## Context & Problem
Multi-turn software development with autonomous agents requires a safe, non-blocking way to dispatch confirmed implementation work while continuing live conversation and steering. Without physical directory isolation, concurrent tasks collide on branch checkouts or dirty working trees. Without bounded context envelopes, agents dilute attention and inflate costs scanning irrelevant codebase directories.

## Decision
1. **Physical Isolation via Git Worktrees**:
   - Every asynchronous or background Runner executes inside a dedicated Git worktree (`git worktree add ../<repo>-<purpose> -b <branch>`).
   - The primary workspace where the owner and Orchestrator interact remains clean and non-blocked.
   - Concurrently dispatched Objectives must have disjoint file perimeters.
2. **On-Demand Context-Sandwich Compilation**:
   - Delegated workers receive a compact (~1,200 token baseline) 3-layer charter (Top: Frozen Truth · Middle: Steering & Decisions · Bottom: Target Files & Stop Command).
   - Pruning follows a deterministic 3-tier drop hierarchy when content overflows.
   - Bounded repairs allow an elastic window up to ~3,500 tokens for diagnostic traces.
   - Whole-workspace tree scanning by delegated subagents is strictly forbidden.
3. **Pipelined Dispatch**:
   - As soon as an Objective's pass boundary and inputs are confirmed, work dispatches immediately without waiting for the entire campaign to be settled.
4. **Targeted Mechanical Interface Growth (`spectacular charter` & `spectacular decide`)**:
   - Owner pre-authorizes introducing two new CLI commands (`spectacular charter <ref>/<objective>` and `spectacular decide`), growing the public mechanical command surface from **14 to 16**.
   - Activation of this surface growth is conditioned on a dedicated Mission (M15) demonstrating a $\ge 40\%$ reduction in token consumption compared to full directory scans and passing the paired behavioral benchmark suite without regression.
