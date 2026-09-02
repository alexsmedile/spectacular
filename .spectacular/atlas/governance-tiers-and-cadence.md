---
type: Atlas
title: Governance tiers and cadence
---

# Atlas: Governance tiers and cadence

This is the visual companion to [D30](../decisions/D30-gated-pipeline-work-boards-and-graduated-orchestration.md). It visualizes the 4-tier Graduated Governance Ladder and answers: *"What governs this slice of work?"*

## 1. The 4-Tier Graduated Governance Ladder

```mermaid
quadrantChart
    title Ceremony vs Consequence
    x-axis Low Ceremony --> High Ceremony
    y-axis Routine Task --> Critical Milestone
    quadrant-1 Tier 3: Mission
    quadrant-2 Tier 2: Dispatch
    quadrant-3 Tier 0: Inline
    quadrant-4 Tier 1: Wave
    "Direct Chat Edits": [0.15, 0.2]
    "WorkBoard Pipeline": [0.4, 0.4]
    "Side Session in Worktree": [0.65, 0.6]
    "Governed Mission Bundle": [0.85, 0.85]
```

## 2. Tier Comparison & Properties

```mermaid
flowchart LR
  subgraph T0["Tier 0: governance: inline"]
    direction TB
    T0_Desc["Direct pair programming<br/>Run in: lead-checkout<br/>Artifacts: None (chat only)"]
  end

  subgraph T1["Tier 1: governance: board"]
    direction TB
    T1_Desc["Gated dependency pipeline<br/>Run in: lead-checkout<br/>Artifacts: type: WorkBoard"]
  end

  subgraph T2["Tier 2: governance: brief"]
    direction TB
    T2_Desc["Temporary teammate in worktree<br/>Run in: linked-worktree<br/>Artifacts: Dispatch Brief + Session"]
  end

  subgraph T3["Tier 3: governance: mission"]
    direction TB
    T3_Desc["Sovereign immutable milestone<br/>Run in: linked-worktree / sandbox<br/>Artifacts: M<N>.md + Handoff + Charter"]
  end

  T0 --> T1 --> T2 --> T3
```
