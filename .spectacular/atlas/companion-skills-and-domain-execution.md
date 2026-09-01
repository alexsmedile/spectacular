---
type: Atlas
title: Companion skills and domain execution topology
---

# Atlas: Companion Skills and Domain Execution Topology

This Atlas maps Spectacular's **4-Skill Starting Pack**, the separation between **Mission Governance** and **Domain Execution**, and the deterministic **Expansion Handoff State Machine**.

---

## 1. Outcome Board

| Actor | Journey Step | Desired Outcome | Success Signal |
|---|---|---|---|
| **Mission Lead** | Frame problem & bounds | Govern scope, contracts, and proof without domain bloat | Single-file Mission (`M<N>.md`) with clear capability boundaries |
| **System Architect** | Design topology & ADRs | Make system shapes and trade-offs traceable to NFRs | C4 Context/Container models + verified ADR records |
| **Data Modeler** | Design schemas & migrations | Generate normalized schemas, ERDs, and zero-downtime DDL | Crow's Foot ERD + dialect DDL + non-blocking Expand/Contract plan |
| **Prototyper** | De-risk ambiguous choices | Settle single open axes across 3 tracer options (A/B/C) | Decision ledger locked on selected lineage (`root → B → A`) |

```mermaid
flowchart LR
    M[Mission Initialized] --> P{Ambiguous Direction?}
    P -->|Yes: 3 Tracer Spike| RP[rapid-prototyping]
    P -->|No| SA[system-architecture]
    RP -->|Lineage Locked| SA
    SA -->|Data Ownership Settled| DM[data-modeling]
    DM -->|Schema & DDL Verified| E[Governed Execution & Proof]
```

---

## 2. System Board

### A. The 4-Pillar Starting Pack

```mermaid
flowchart TD
    subgraph Governance ["Governance & Orchestration Spine"]
        SPEC["<b>spectacular</b><br>• Missions & Flight Plans<br>• Contracts & Objectives<br>• Verification Receipts & FROST Reviews"]
    end

    subgraph DomainExecution ["Domain Execution Micro-Kernels"]
        SA["<b>system-architecture</b><br>• C4 Context & Container Models<br>• Service & Bounded Contexts<br>• Architecture Decision Records (ADRs)"]
        
        DM["<b>data-modeling</b><br>• Conceptual (Chen) & Logical (Crow's Foot)<br>• Physical DDL (Postgres/MySQL/SQLite)<br>• Zero-Downtime Migrations & Indexes"]
        
        RP["<b>rapid-prototyping</b><br>• 3 Tracer Fragments (A, B, C)<br>• 5-Tier Fidelity Ladder (Atom &rarr; System)<br>• Decision Ledger (Locked, Open, Level, Lineage)"]
    end

    SPEC -->|Dispatches Subagent / Delegates| SA
    SPEC -->|Dispatches Subagent / Delegates| DM
    SPEC -->|Dispatches Subagent / Delegates| RP
    
    SA <-->|Expansion Handoff: DDL & ERD| DM
    SA <-->|Expansion Handoff: 3-Option Spike| RP
```

---

### B. Trigger & Boundary Matrix (Zero Overlap)

| Skill | Primary Operational Domain | Quoted Trigger Keywords | Negative Constraints (`DO NOT`) |
|---|---|---|---|
| **`spectacular`** | Governance, Mission lifecycle, and flight plans | `"start mission"`, `"spectacular decide"`, `"flight plan"`, `"autopilot"`, `"complete mission"`, `$spectacular` | Do not invoke for ungrounded chat, generic planning, or ordinary git ops. |
| **`system-architecture`** | System topology, C4 models, and trade-off ADRs | `"system architecture"`, `"architecture review"`, `"C4 diagram"`, `"ADR"`, `"service boundaries"`, `"bounded context"` | Do not invoke for local code refactors, database DDL, or isolated UI styling. |
| **`data-modeling`** | ERDs, SQL DDL, indexes, and schema migrations | `"ER diagram"`, `"ERD"`, `"database schema"`, `"data model"`, `"SQL DDL"`, `"database migration"`, `"Crow's Foot"` | Do not invoke for simple data inspection or high-level topology without entities. |
| **`rapid-prototyping`** | 3-option tracer spikes and progressive fidelity growth | `"prototype"`, `"explore options"`, `"tracer fragment"`, `"matrix prototype"`, `"compare designs"`, `"fidelity ladder"` | Do not invoke for single-path tasks, simple bugs, or open-ended brainstorming. |

---

## 3. Domain Board

```mermaid
flowchart LR
    subgraph GovernanceContext [Context: Governed Execution]
        Mission[Entity: Mission]
        Contract[Entity: Contract]
        Objective[Entity: Objective]
    end

    subgraph DesignContext [Context: System Design]
        Topology[Entity: System Topology]
        ADR[Entity: Decision Record]
        C4[Value: C4 Model]
    end

    subgraph DataContext [Context: Data Modeling]
        ERD[Entity: ER Diagram]
        DDL[Value: Physical DDL]
        Migration[Entity: Migration Plan]
    end

    subgraph SpikeContext [Context: Rapid Prototyping]
        Tracer[Entity: Tracer Fragment]
        Ledger[Value: Decision Ledger]
    end

    Mission -->|governed_by| Contract
    Mission -->|contains 1..*| Objective
    
    Objective -->|delegates_to 0..1| Topology
    Objective -->|delegates_to 0..1| ERD
    Objective -->|delegates_to 0..1| Tracer
    
    Topology -->|has 0..*| ADR
    Topology -->|references| C4
    
    ERD -->|implemented_by| DDL
    ERD -->|migrated_by| Migration
    
    Tracer -->|governed_by| Ledger
```

---

## 4. Links and Open Questions

- **Companion Skill Sources:** `skills/system-architecture/`, `skills/data-modeling/`, `skills/rapid-prototyping/`
- **Related Atlases:**
  - [operating-levels-and-software-factory.md](operating-levels-and-software-factory.md)
  - [domain-overview.md](domain-overview.md)
  - [lean-autopilot-and-orchestration.md](lean-autopilot-and-orchestration.md)
- **Open Question:** As UI engineering companion skills (e.g. `ui-craftsman` or `component-architect`) are introduced, should they nest directly under `rapid-prototyping` Level 2–4 or form a 5th standalone pillar?
