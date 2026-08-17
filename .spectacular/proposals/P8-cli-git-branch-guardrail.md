---
type: Proposal
id: 01a01048-2a69-74d6-a5d5-86d8daf018e5
ref: P8
title: Add mechanical Git branch guardrails to mission start
status: draft
created_by: Alex
created: "2026-08-17T15:15:00Z"
updated: "2026-08-17T15:15:00Z"
scope:
    - v2
target_contract: Contract:01a00a20-63dd-7670-97f1-9eb8e12adc3a
---

# Add mechanical Git branch guardrails to mission start

Exploration for a possible Mission. Nothing here is frozen. The problem statement is
firm and demonstrated; the mechanical behavior is a proposed direction, not an active
Mission.

## The problem in one line

`spectacular mission start` silently accepts activation on default branches (`main`/`master`),
allowing multi-step Missions to destroy review, merge, and isolation boundaries.

## Where this came from (The M9 Incident)

During Mission M9, an agent activated, executed, and committed ~4,300 lines of deletions
directly on `main`. No branch was created, no merge point existed, and no clean review
boundary remained. The owner had to catch it retrospectively:

> "you did not start mission 9 on a branch, which I believe was the wrong move."

Retrofitting the branch required `git branch <name> <sha>` followed by `git reset --hard <base>`,
which wiped uncommitted working-tree state and was only recoverable because commits had
already been recorded.

The root cause was that `execute.md` provided prose guidance on *which* branch to choose among
existing branches, but the mechanical layer (`cmd/spectacular` / `internal/missionbundle`)
silently recorded `baseline: branch: main` without objection. The system watched a Mission
freeze onto the default branch and said nothing.

## The Core Philosophy: CLI Invariant Floor vs Skill Orchestration

Spectacular divides responsibilities clearly:
- **The Go binary** handles invariants, hashes, transitions, and JSON projections.
- **The Skill** handles judgment, planning, problem-solving, and human interactions.
- **The filesystem (Canonical Markdown)** remains human-readable without running proprietary daemons or databases.

When a rule is exact, cheap to check, and failure is expensive (e.g. destructive git state or
broken merge boundaries), the CLI should provide the mechanical safety floor.

Leaving branch safety purely to prompt guidance forces LLMs to re-derive the negative constraint
on every run. A mechanical refusal turns a costly operational mistake into an instant,
runnable correction.

## Proposed Mechanical Behavior

### 1. Refusal on Default Branches by Default

When `spectacular mission start <plan.md>` is invoked, the command inspects the active Git branch:

- If the active branch is a default branch (`main`, `master`, `trunk`, or default remote HEAD),
  `mission start` refuses activation:
  ```text
  refused default_branch_activation field baseline.branch: multi-step Missions must not be activated directly on default branch 'main'
  correction: create and switch to a dedicated branch (e.g. `git checkout -b <mission-slug>`) or pass `--allow-default-branch` for an intentional quick patch
  ```

- Refusal JSON payload:
  ```json
  {
    "valid": false,
    "code": "default_branch_activation",
    "field": "baseline.branch",
    "problem": "cannot activate mission on default branch 'main'",
    "correction": "git checkout -b <slug> (or pass --allow-default-branch)"
  }
  ```

### 2. The Explicit Escape Hatch: `--allow-default-branch`

An owner or operator may legitimately need to execute an emergency one-step patch directly on `main`
when no concurrent session exists.

To allow this without weakening the default invariant:
- Add an explicit flag: `spectacular mission start plan.md --allow-default-branch`
- When passed, activation succeeds on `main` and records `baseline.branch: main` with an explicit
  notice or metadata indicating deliberate default-branch execution.

### 3. Preserving the 11-Command Surface

This proposal does **not** add a new command or noun. It tightens the validation and transition
logic of `mission start`, which is already part of the frozen command surface in `CC-missioncli`.

## Alternatives Considered & Declined

1. **Leave it strictly to the Skill.**
   - *Rejected because:* Proved brittle in practice. LLMs under dense context miss the negative
     instruction. The cost of a failure is polluted git history and potential tree loss on reset.

2. **Have the CLI automatically create a git branch.**
   - *Rejected because:* The CLI should not perform side-effecting git operations invisibly
     without explicit operator consent and naming intent. Git branch creation is Skill/Operator
     domain; the CLI only validates the boundary.

3. **Strictly forbid default branch activation with zero exceptions.**
   - *Rejected because:* Single-file emergency fixes or standalone workspace updates legitimately
     occur on `main` when authorized by the owner. An explicit flag (`--allow-default-branch`)
     preserves operator agency while preventing accidental drift.

## Verification & Proof Requirements (for future Mission)

- Table-driven unit tests in `internal/missionbundle` asserting refusal on `main`, `master`, and `trunk`.
- Tests asserting activation succeeds on feature branches (e.g., `feat/m10-gap-lifecycle`).
- Tests asserting activation succeeds on `main` when `--allow-default-branch` is explicitly provided.
- Regression fixture verifying the refusal output matches standard Spectacular error envelope format.
