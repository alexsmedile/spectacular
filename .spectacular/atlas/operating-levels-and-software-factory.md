---
type: Atlas
title: Operating levels and software factory architecture
---

# Atlas: Operating Levels and Software Factory Architecture

This Atlas projects the operational topology of Spectacular's **20-Level Agentic Operating Model**, **Dynamic Operating Dial** (`mode: leverage` vs. `mode: control`), **5 Foundational Anchors**, and **Trunk-First Multi-Agent Collaboration**.

---

## 1. Outcome Board

| Actor | Journey Step | Desired Outcome | Success Signal |
|---|---|---|---|
| **Human Architect** | Ground domain & boundaries | Settle hard data models and invariants before code execution | Schemas, types, and decisions committed upstream to `main` |
| **Orchestrator** | Choose leverage vs. control | Set mission posture based on task risk and distribution | `mode: leverage` (speed) or `mode: control` (audit) declared |
| **Worker / Solo** | Self-heal against tools | Implement features and fix test failures autonomously | Fast Tier 1 unit test pass (`exit 0`) with zero prompt loops |
| **Reviewer** | Independent claim audit | Inspect diffs and primary evidence without modifying code | Structured FROST review verdict recorded (`pass`/`fail`) |

```mermaid
flowchart TD
    subgraph Pyramid["The 5 Operating Tiers (20 Levels)"]
        G5["<b>Group 5: Agentic System (Levels 17–20)</b><br>Product · Agent · AI Developer Workflow (ADW) · Software Factory"]
        G4["<b>Group 4: Delivery & Intent (Levels 13–16)</b><br>Application · Repository · Plan · Media-Rich Documentation"]
        G3["<b>Group 3: Data & Execution (Levels 09–12)</b><br>Database Table · Database · CLI Tools · Deterministic Scripts"]
        G2["<b>Group 2: Code Structure (Levels 06–08)</b><br>File · Module · Directory Boundaries"]
        G1["<b>Group 1: Code Primitives (Levels 01–05)</b><br>Line · Block · Function · Type Contract · Class"]
    end

    G5 -->|Steers & Orchestrates| G4
    G4 -->|Compiles Context Sandwich| G3
    G3 -->|Executes deterministic tools against| G2
    G2 -->|Enforces modular AST boundaries on| G1

    style G5 fill:#2d3748,stroke:#4a5568,color:#fff
    style G4 fill:#2b6cb0,stroke:#3182ce,color:#fff
    style G3 fill:#2c7a7b,stroke:#319795,color:#fff
    style G2 fill:#2f855a,stroke:#38a169,color:#fff
    style G1 fill:#744210,stroke:#975a16,color:#fff
```

---

## 2. System Board

### A. The Dynamic Operating Dial: System Leverage vs. Direct Control

```mermaid
flowchart LR
    subgraph Leverage["🚀 Move UP: System Leverage"]
        L1["<b>When to choose:</b><br>• Familiar domain & repeatable patterns<br>• Strong test suite with clear pass boundary<br>• Routine features, refactors, bug fixes"]
        L2["<b>Spectacular Posture:</b><br>• <code>mode: leverage</code> (Default)<br>• Single-file mission (<code>M&lt;N&gt;.md</code>)<br>• Worker self-healing test loop<br>• Zero review record sprawl"]
    end

    subgraph Control["🔍 Move DOWN: Direct Control"]
        C1["<b>When to choose:</b><br>• Out-of-Distribution (OOD) logic<br>• High-risk auth, payments, zero-downtime DB<br>• Ambiguous failure traces or security cutovers"]
        C2["<b>Spectacular Posture:</b><br>• <code>mode: control</code><br>• Hard schema & type contract inspection<br>• Explicit step-by-step owner gates<br>• Dedicated independent FROST review"]
    end

    Leverage <-->|Dynamic Range Slider| Control
```

---

### B. The 5 Foundational Anchors Mapping

Before autonomous execution, Spectacular grounds the rules of physics across 5 canonical project surfaces:

```mermaid
flowchart TD
    subgraph Anchors["5 Core Foundational Anchors"]
        A1["<b>1. Boundaries & Non-Goals</b><br><code>.spectacular/PROJECT.md</code><br><i>Explicit constraints on what NOT to do</i>"]
        A2["<b>2. Vocabulary & Ontology</b><br><code>.spectacular/VOCABULARY.md</code><br><i>Canonical ubiquitous language</i>"]
        A3["<b>3. Invariants & Safety</b><br><code>.spectacular/GUARDRAILS.md</code><br><i>Non-negotiable architectural rules</i>"]
        A4["<b>4. Data Contracts & Schemas</b><br><code>codebase types / contracts/</code><br><i>Database tables, structs, API shapes</i>"]
        A5["<b>5. State Machines & Lifecycles</b><br><code>.spectacular/atlas/</code><br><i>Visual Mermaid state diagrams</i>"]
    end
```

---

### C. Trunk-First Collaboration & Worktree Isolation

```mermaid
flowchart TD
    subgraph Trunk["👑 Default Branch (`main`): Living Ground Truth"]
        Dec["⚖️ Decisions (<code>D1..DN</code>)<br><i>Decided upstream first</i>"]
        Anch["🔒 Core Anchors (<code>PROJECT.md</code>, <code>VOCABULARY.md</code>)"]
        Flight["🗺️ Flight Plan (<code>campaigns/</code>)"]
    end

    Trunk -->|<code>git worktree add</code>| WT1["🌲 Worktree 1 (Agent A)<br><b>Branch: feat/M29-auth</b><br>• Isolated working copy<br>• Runs Tier 1 tests locally<br>• Emits worker_done"]
    Trunk -->|<code>git worktree add</code>| WT2["🌲 Worktree 2 (Agent B)<br><b>Branch: feat/M30-db-pool</b><br>• Isolated working copy<br>• Zero file locking with Agent A<br>• Emits worker_done"]

    WT1 -->|Passes tests + Review| PR1["Pull Request / Fast-Forward Merge"]
    WT2 -->|Passes tests + Review| PR2["Pull Request / Fast-Forward Merge"]

    PR1 --> Trunk
    PR2 --> Trunk
```

---

## 3. Links and References

- Product Documentation: [`docs/architecture.md`](../../docs/architecture.md), [`docs/process.md`](../../docs/process.md)
- Skill References: [`skills/spectacular/SKILL.md`](../../skills/spectacular/SKILL.md), [`skills/spectacular/references/prepare.md`](../../skills/spectacular/references/prepare.md), [`skills/spectacular/references/execute.md`](../../skills/spectacular/references/execute.md), [`skills/spectacular/references/close.md`](../../skills/spectacular/references/close.md)
- Key Decisions: [`D28-dynamic-operating-dial-and-anchors`](../decisions/D28-dynamic-operating-dial-and-anchors.md), [`D15-branch-guardrail-at-activation`](../decisions/D15-branch-guardrail-at-activation.md), [`D16-auto-prune-merged-worktrees`](../decisions/D16-auto-prune-merged-worktrees.md)
