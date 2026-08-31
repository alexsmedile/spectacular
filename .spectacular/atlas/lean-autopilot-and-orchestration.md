---
type: Atlas
title: Lean autopilot and orchestration architecture
---

# Atlas: Lean Autopilot and Orchestration Architecture

This Atlas projects the operational topology of Spectacular's **Lean 3-Layer Autopilot Model**, **3-Tier Layout Matrix**, **Dual-Lane Orchestration**, and **3-Tier Question Escalator**.

---

## 1. Outcome Board

| Actor | Journey Step | Desired Outcome | Success Signal |
|---|---|---|---|
| **Human Owner** | Steer project & bulk-decide | Settle architectural forks once with zero repetitive questions | `spectacular decide` records permanent `D<N>` rulings |
| **Orchestrator** | Plan and supervise roadmap | Sequence 4–8 milestone blocks and dispatch workers | Flight plan maintains topological order with zero file clutter |
| **Worker Subagent** | Execute bounded mission | Implement code and pass tests autonomously in autopilot | Tests pass (`exit 0`) + Git commit with zero sub-record sprawl |
| **Independent Reviewer** | Verify high-stakes claims | Inspect diffs and primary evidence without modifying code | Reviewer returns clean FROST verdict (`pass`/`fail` + findings) |

```mermaid
flowchart LR
    A["💡 1. Bulk Ideate & Decide<br><i>`spectacular decide` (D1..DN)</i>"] --> B["🗺️ 2. Flight Plan Roadmap<br><i>.spectacular/campaigns/flight-plan.md</i>"]
    B --> C["🚀 3. Single-File Mission<br><i>`spectacular mission start` (≤500t)</i>"]
    C --> D["🤖 4. Supervised Subagent<br><i>Autopilot execution (≤300t charter)</i>"]
    D --> E["✅ 5. Verified Gate & Close<br><i>Tests pass = Proof · Owner completes</i>"]
```

---

## 2. System Board

### A. The 3-Tier Layout Judgment Matrix

```mermaid
flowchart TD
    Start["New Work Item"] --> Choice{Layout Tier Selection}

    Choice -->|"90% of Tasks: Routine Code"| T1["<b>Tier 1: Single-File Mission</b><br>• <code>missions/M&lt;N&gt;.md</code> only (≤500t)<br>• Inline deliverables & test boundary<br>• Zero sub-folders"]

    Choice -->|"8% of Tasks: External Proof"| T2["<b>Tier 2: Hybrid Earned</b><br>• <code>M&lt;N&gt;.md</code> + 1 earned record<br>• <code>evidence/</code> for live API receipts<br>• <code>objectives/</code> for parallel worktrees"]

    Choice -->|"2% of Tasks: High Stakes"| T3["<b>Tier 3: Full Governed Bundle</b><br>• Zero-downtime DB cutovers, auth/crypto<br>• Formal <code>checkpoints/</code> & <code>reviews/</code><br>• Complete multi-record directory"]
```

---

### B. Dual-Lane Orchestration: Supervised Dispatch vs. Full Handoff

```mermaid
flowchart TD
    subgraph LaneA ["Lane A: Supervised Subagent Dispatch (90% Default)"]
        O1["👑 Live Orchestrator Session<br><i>Retains Mission lifecycle ownership</i>"]
        W1["🤖 Autonomous Worker Subagent<br><i>Receives compact ≤300-token charter</i>"]
        R1["⚡ Reactive Execution<br><i>Runs tests, commits locally, emits worker_done</i>"]
        O1 -->|Dispatches via host channel| W1
        W1 -->|Passes tests & commits| R1
        R1 -->|worker_done signal| O1
    end

    subgraph LaneB ["Lane B: Full Ownership Handoff (10% Transfer)"]
        O2["👑 Outgoing Operator / Harness"]
        H2["📜 `spectacular handoff record`<br><i>Immutable Git-bound contract (asserted vs assumed)</i>"]
        W2["🎯 Incoming Operator / New Harness<br><i>Takes full Mission ownership; sender exits</i>"]
        O2 -->|Records formal handoff| H2
        H2 -->|Transfers ownership| W2
    end
```

---

### C. The Fail-Fast Escalation & Decision Gate Protocol

```mermaid
flowchart TD
    Worker["🤖 Worker Subagent executing in Autopilot"] --> Check{Discovers unrecorded choice or interface conflict?}
    Check -->|No| Continue["Continue implementation & run test suite"]
    Check -->|Yes: Fail-Fast Stop| Escalate["⚠️ Send Escalation Message<br><i>Names question, options, and recommended default</i>"]
    Escalate --> Orch["👑 Orchestrator & Owner Decision Gate"]
    Orch --> Decide["⚖️ `spectacular decide`<br><i>Records permanent ruling D&lt;N&gt;</i>"]
    Decide --> Resume["🔄 Resume Worker with locked Decision ID D&lt;N&gt;"]
    Resume --> Continue
```

---

### D. The 3-Tier Question Escalator

```mermaid
flowchart LR
    subgraph T1 ["Tier 1: Optimistic Consent"]
        Q1["⚡ 1-Line Statement of Intent<br><i>'Proceeding with X (reason) unless you prefer Y'</i><br><b>Non-blocking</b>"]
    end

    subgraph T2 ["Tier 2: Structured Batch Cards"]
        Q2["📋 Numbered Qs + Lettered Options<br><i>1. Question ➔ A, B, C (Recommended default)</i><br><b>Batch shorthand replies ('A, B, A')</b>"]
    end

    subgraph T3 ["Tier 3: Spectrum & Modals"]
        Q3["🧭 Trade-Off Spectrums & Modals<br><i>Frames competing design axes or interactive UI</i><br><b>Deep ambiguous forks</b>"]
    end

    T1 -->|Higher Stakes| T2
    T2 -->|Deep Ambiguity| T3
```

---

## 3. Links and References

- Product Documentation: [`docs/architecture.md`](../../docs/architecture.md), [`docs/process.md`](../../docs/process.md), [`docs/quickstart.md`](../../docs/quickstart.md)
- Skill References: [`skills/spectacular/references/runtime.md`](../../skills/spectacular/references/runtime.md), [`skills/spectacular/references/owner-guidance.md`](../../skills/spectacular/references/owner-guidance.md), [`skills/spectacular/references/prepare.md`](../../skills/spectacular/references/prepare.md)
- Decisions: `D11-proposal-retirement`, `D15-branch-guardrail-at-activation`, `D24-accepted`
