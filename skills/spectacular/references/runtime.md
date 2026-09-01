# Runtime & Delegation Micro-Kernel

## 1. Trigger Context
Orchestrator packaging subagent charters, supervised dispatch, or cross-party handoffs.

## 2. CLI Palette
```bash
spectacular charter <ref>[/<obj>] --json             # Compile context sandwich (≤300 tok)
spectacular handoff record <ref> draft.md --by <actor> # Record cross-session transfer
```

## 3. Supervised Dispatch (Zero Framework Lore for Workers)
Dispatched workers receive strictly:
- **Code Paths**: Assigned files (`allowed_changed_paths`).
- **Decision Constraints**: Excerpts from applicable `D<N>` records.
- **Verification Command**: The test command (e.g. `sh tests/check.sh`).
- **Stop Triggers**: Stop immediately on unexpected design fork.

```text
# Sample Dispatch Payload
Task: Implement SQLite Storage Engine (src/storage.py).
Allowed Paths: src/storage.py, src/db.py
Forbidden Paths: .spectacular/**, tests/**
Constraints from D12: Use Python built-in sqlite3. WAL mode enabled.
Verification Command: sh tests/check.sh
Escalation: Halt if schema migrations require external tooling.
```

## 4. The 1-Turn Escalation Protocol
1. **Worker Halts**: Stop immediately on unrecorded fork. Send 1-line query to Orchestrator.
2. **Orchestrator Decides**: Orchestrator runs `spectacular decide --title "<Title>" --disposition "<Choice>"`.
3. **Worker Resumes**: Orchestrator returns decision ID (`D<N>`); worker resumes with locked choice.

## 5. Negative Constraints (DO NOT)
- **DO NOT** make workers manage `.spectacular/` files or create runs/checkpoints.
- **DO NOT** dispatch subagents with large prompt dumps (keep charter $\le 300\text{--}500$ tokens).
