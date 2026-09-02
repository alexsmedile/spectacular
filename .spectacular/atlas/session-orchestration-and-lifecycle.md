---
type: Atlas
title: Session orchestration and lifecycle
---

# Atlas: Session orchestration and lifecycle

This is the visual companion to [D30](../decisions/D30-gated-pipeline-work-boards-and-graduated-orchestration.md) and [runtime.md](../../skills/spectacular/references/runtime.md). It models how work progresses across sequential waves and temporary side sessions.

## 1. The Gated Wave Flow

Sequential execution in the Lead session is the baseline. Parallel side sessions in isolated Git worktrees are earned only after interface/contract gates lock.

```mermaid
flowchart TD
  subgraph Lead[Primary Thread / Lead Checkout]
    Frame[1. Frame Outcome] --> Contract[2. Define Contract & Types]
    Contract --> Gate{Decision Gate: Contract Locked?}
    Gate -- No --> Refine[Refine Contract] --> Gate
    Gate -- Yes --> Dispatch[3. Dispatch Side Sessions]
    Integrate[4. Integrate Branches] --> Verify[5. Project-Wide Verification]
    Verify --> Complete([Done & Verified])
  end

  subgraph SideWorkers[Earned Side Sessions]
    Dispatch --> WorkerA[Teammate A: Implement Parser Engine<br/>.worktrees/parser-engine]
    Dispatch --> WorkerB[Teammate B: Build Tokenizer<br/>.worktrees/tokenizer]
    WorkerA --> ReceiptA[Return Receipt A]
    WorkerB --> ReceiptB[Return Receipt B]
    ReceiptA --> Integrate
    ReceiptB --> Integrate
  end
```

## 2. The 7-State Session State Machine

A side session never marks an item done. It only returns code, diffs, and test receipts. The Lead alone integrates and verifies.

```mermaid
stateDiagram-v2
  [*] --> planned
  planned --> ready: Upstream gate locked
  ready --> active: Dispatched to worktree
  
  state active {
    [*] --> executing
    executing --> executing: Heartbeat ping
  }

  active --> returned: Task completed (Receipt emitted)
  active --> blocked: Architectural roadblock hit
  active --> aborted: Session cancelled / timed out
  
  blocked --> ready: Lead resolves blocker
  blocked --> aborted: Abandoned
  
  returned --> integrated: Lead reviews diff & merges branch
  integrated --> verified: Project verification green
  
  verified --> [*]: Worktree pruned safely
  aborted --> [*]: Worktree pruned safely
```

## 3. Physical Workspace Modes

```mermaid
classDiagram
  class PhysicalWorkspaceModes {
    +lead-checkout: Primary working tree for Lead session
    +linked-worktree: Isolated git worktree for single writer (.worktrees/<slug>)
    +sandbox: Disposable container or experiment branch
    +read-only: Non-mutating scout or auditor thread
  }
```
