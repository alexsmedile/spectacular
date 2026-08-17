# Changelog

## 2.2.0 — 2026-08-17

Dead v1 surface removed, derived state reads recorded reviews, and Skill guidance
for branching, worktrees, and execution mode.

- Removes the unreachable v1 context-compiler chain — `internal/context`,
  `internal/projection`, and `internal/guardrails` — as one unit, and
  `internal/index`, the v1 predecessor of `discovery.Workspace.Lookup`. Prunes
  `internal/governance` to the transaction machinery that has live callers.
- Fixes the derived next action to consult the reviews a Mission carries. It
  previously fired on implemented Objectives alone, so a recorded review could not
  retire the instruction to record one. A review bound to a stale activation
  fingerprint keeps asking, since that is the drift the fingerprint exists to catch.
- Surfaces the underlying cause in human refusals. A YAML syntax error reported
  only `invalid_known_field`, sending readers hunting through field names while the
  parser's line number sat unprinted in the JSON envelope.
- Adds Skill guidance for choosing a branch and a worktree by what the job needs,
  including running a Mission session and a feedback session concurrently without
  either destroying the other's work.
- Adds an execution-mode question at activation — autopilot, checkpoints, or a
  named human-in-the-loop moment — so involvement is settled once instead of
  arriving one gate at a time.
- Accounts for every untracked working-tree path with a `.gitignore` rule that
  states its reason, and records `_snapshots/` as local recovery only.

## 2.1.1 — 2026-08-17

Test performance optimizations and contributor verification guidance.

- Optimizes `internal/missionbundle` test fixtures by caching template git repositories and parallelizing rollback subtests.
- Parallelizes 4-platform release archive compilation in `cmd/assemble-release`.
- Retains persistent Go test cache across verification runs and downstream installer/release test scripts.
- Codifies tiered verification guidance (`quick`, `acceptance`, `release`, `all`) in `AGENTS.md`.

## 2.1.0 — 2026-08-16

Governed-autonomy controls for preparation, delegation, review, and runtime limits.

- Adds adaptive preparation diagnostics and a frozen four-field completion criterion for every Mission claim.
- Adds criterion-driven automatic, clustered, and independent review without recursive critic loops.
- Validates Objective dependency DAGs and requires dependency-bound, disjoint Handoff claim scopes with explicit return contracts.
- Adds truthful hard, observed, and unsupported Autopilot caps for wall time, tokens, spend, parallel workers, and repair rounds.
- Fixes cold recovery for newly activated Missions before their first Checkpoint.

## 2.0.0-rc.2 — 2026-08-10

Human-operability correction for the v2 release candidate.

- Replaces flat UUID filenames with named project Anchors and cohesive Mission
  bundles while retaining UUIDv7 identity and SHA-256 revision fingerprints.
- Adds scoped human references for Missions, Objectives, Runs, Checkpoints,
  Evidence, Decisions, Gaps, Handoffs, and Assessments.
- Makes default CLI cards human-first and keeps exact machine data in `--json`.
- Commits deterministic, non-authoritative workspace and Mission indexes.
- Adds atomic whole-bundle Mission archival and a real self-hosted workspace.
- Replaces flat test fixtures with human-layout scenarios and adds a real-binary
  acceptance layer covering cold recovery, executable pointers, governed
  closure, archival, refusals, and zero-mutation reads.
- Fixes stale active indexes/directories after Mission archival, stable bundle
  placement after title changes, and invalid empty optional Evidence fields.

RC.1 is superseded because its machine-oriented workspace representation did
not satisfy the human-comprehension contract.

## 2.0.0-rc.1 — 2026-08-10

First externally consumable Spectacular v2 release candidate.

- Breaking: removes the v1 public product surface and compatibility paths; Git
  recovery pointers preserve the frozen history.
- Adds pointer-first retrieval and a governed Mission loop with explicit
  authority, Evidence, assessment, reconciliation, and closure boundaries.
- Ships native four-platform archives and a checksum-verifying installer.

Known limitations: v1 workspaces are unsupported; there is no migration or
compatibility reader; discovery is pointer-driven rather than broad search;
and release publication remains outside the local CLI.
