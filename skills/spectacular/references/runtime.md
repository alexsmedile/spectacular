# Runtime & Delegation Micro-Kernel

## 1. Trigger Context
Orchestrator managing multi-step dependency waves, Dispatch Briefs, worktree side sessions, or formal Mission Charters and Handoffs.

## 2. CLI Palette
```bash
# Tier 2: Plain Briefs & Worktrees (Zero CLI Churn)
# Lead provides Dispatch Brief in chat; side worker runs in .worktrees/<slug>

# Tier 3: Governed Charters & Handoffs (When Formal Boundaries Are Required)
spectacular charter <ref>[/<obj>] --json             # Compile context sandwich (≤1200 tok)
spectacular handoff record <ref> draft.md --by <actor> # Record immutable cross-party transfer
```

## 3. The 4 Governance Tiers (`governance:`)
Answers: *"What governs this slice of work?"*
- **`governance: inline` (Tier 0)**: Direct pair programming in `lead-checkout`. Zero governance files.
- **`governance: board` (Tier 1)**: Gated dependency pipeline on `type: WorkBoard`. Lead tracks order & gates.
- **`governance: brief` (Tier 2)**: Temporary teammate in `linked-worktree` with a plain-English Dispatch Brief.
- **`governance: mission` (Tier 3)**: Full immutable `M<N>.md` contract, Handoff, and compiled Charter.

## 4. Physical Execution Workspace Modes
- **`lead-checkout`**: Primary working tree for the Lead Orchestrator and strictly sequential steps.
- **`linked-worktree`**: Isolated git worktree (`.worktrees/<slug>`) on a dedicated branch for one writer.
- **`sandbox`**: Disposable container or experiment branch with zero merge authority.
- **`read-only`**: Non-mutating scout or reviewer thread inspecting diffs.

## 5. The Dispatch Brief (Default for Tier 2)
Dispatched workers receive a self-contained brief:
- **Goal & Prerequisite**: What to build and which upstream interface gate has locked.
- **Physical Workspace**: Worktree path and branch name.
- **Allowed Writable Paths**: Exact file subtrees the worker may touch.
- **Verification Command**: Command proving completion (e.g. `go test -v ./internal/parser/...`).
- **Stop / Blocked Condition**: When to halt and return a blocker receipt.

```markdown
### Dispatch Brief Sample
Goal: Implement recursive AST parser (internal/parser/).
Prerequisite: AST types locked at commit a8f12c (Gate 1 passed).
Workspace: linked-worktree at .worktrees/parser-engine (branch: ast-redesign/parser-engine).
Writable Paths: internal/parser/**
Forbidden Paths: .spectacular/**, cmd/**
Done When: go test -race ./internal/parser/... passes.
If Blocked: Halt, return current diff, and state the unresolved interface constraint.
```

## 6. The Session Lifecycle & Invariants
```text
planned ──► ready ──► active ──► returned ──► integrated ──► verified
                        │
                        └──► blocked / escalated (returns to Lead)
```
- **"Returned ≠ Done"**: A side worker never marks an item complete. It emits a **Return Receipt** (`commit`, `tests_passed`, `diff_stat`, `blockers`). The Lead alone reviews, merges, and verifies.
- **Heartbeat & Leases**: Active reservations release upon an explicit `returned`/`aborted` receipt or after a heartbeat timeout (default: 15m).
- **Conservative Pruning**: Worktrees (`.worktrees/<slug>`) are pruned only after the Lead records `integrated` or `aborted` and confirms uncommitted diffs are safe.

## 7. Negative Constraints (DO NOT)
- **DO NOT** dispatch parallel workers before upstream contract gates are locked.
- **DO NOT** make workers manage `.spectacular/` files or create runs/checkpoints.
- **DO NOT** dispatch subagents with unbounded prompt dumps (keep compiled charters $\le 1{,}200$ tokens).
