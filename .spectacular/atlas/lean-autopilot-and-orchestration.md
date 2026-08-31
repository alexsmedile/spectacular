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
flowchart TD
    subgraph S1["1. Anchor & Spec (Control)"]
        A1["🎯 Intent & PRD"] --> A2["🔒 5 Anchors: Types, Schemas, State Machines, Boundaries, Vocab"]
        A2 --> A3["🚀 Single-File Mission M&lt;N&gt; (≤500t)<br><i>Set Dial: mode: leverage / mode: control</i>"]
    end
    subgraph S2["2. Budget & Dispatch (Leverage)"]
        A3 --> B1["🥪 `spectacular charter` (≤1,200t Context Sandwich)"]
        B1 --> B2["🤖 Solo Execution or Supervised Dispatch"]
    end
    subgraph S3["3. Execute & Self-Heal (Leverage + Safety)"]
        B2 --> C1["⚡ Targeted Code Edits"]
        C1 --> C2{"Tier 1: Quick Tests"}
        C2 -->|Fails| C3["Self-Repair Loop (≤3 tries with error log)"]
        C3 --> C1
        C2 -->|Ambiguous Fork| C4["⚠️ Fail-Fast Escalation ➔ `spectacular decide`"]
        C4 --> C1
    end
    subgraph S4["4. Audit & Close (Control)"]
        C2 -->|Passes| D1["📋 Claim vs. Diff Verification"]
        D1 --> D2["🔍 Tiered / Batched Adversarial Review"]
        D2 --> D3["🏁 `spectacular mission complete`"]
    end
```

---

## 2. System Board

### A. Tiered Verification Matrix (Zero Duplicate Runs)

```mermaid
flowchart LR
    A["<b>Worker / Solo Loop</b><br>Inner code edits"] -->|Fast & local| T1["<b>Tier 1: Quick / Domain</b><br>Unit tests & syntax check (≤5s)"]
    T1 -->|Green| B["<b>Reviewer / Pre-Check</b><br>Audit claims & diffs"]
    B -->|Structural| T0["<b>Tier 0: Preflight / Lint</b><br>AST boundaries & contract drift"]
    T0 -->|Milestones| T2["<b>Tier 2: Acceptance</b><br>Orchestrator e2e suite"]
    T2 -->|Final Gate| T3["<b>Tier 3: All / Release</b><br>Owner release gate"]
```

---

### B. The 3-Tier Layout Judgment Matrix

```mermaid
flowchart TD
    Start["New Work Item"] --> Choice{Layout Tier Selection}

    Choice -->|"90% of Tasks: Routine Code"| T1["<b>Tier 1: Single-File Mission</b><br>• <code>missions/M&lt;N&gt;.md</code> only (≤500t)<br>• Inline deliverables & test boundary<br>• Zero sub-folders"]

    Choice -->|"8% of Tasks: External Proof"| T2["<b>Tier 2: Hybrid Earned</b><br>• <code>M&lt;N&gt;.md</code> + 1 earned record<br>• <code>evidence/</code> for live API receipts<br>• <code>objectives/</code> for parallel worktrees"]

    Choice -->|"2% of Tasks: High Stakes"| T3["<b>Tier 3: Full Governed Bundle</b><br>• Zero-downtime DB cutovers, auth/crypto<br>• Formal <code>checkpoints/</code> & <code>reviews/</code><br>• Complete multi-record directory"]
```

---

### C. Sub-Record Earned Decision Tree

```mermaid
flowchart TD
    Start["Does this task need sub-records?"] --> Q1{"Is it routine local code<br/>where unit tests pass?"}
    Q1 -->|YES| Single["<b>Single-File M&lt;N&gt;.md ONLY</b><br>• Zero sub-folders<br>• Commit + green tests = Proof"]
    Q1 -->|NO| Fork{"What specific external condition exists?"}

    Fork -->|Live external API receipt needed| Ev["<b>evidence/ E&lt;N&gt;.md</b><br>Third-party network receipts"]
    Fork -->|Permanent handoff across users/tools| Ho["<b>handoffs/ H&lt;N&gt;.md</b><br>Transferring project keys"]
    Fork -->|Milestone review or high-stakes auth/DB| Rv["<b>reviews/ RV&lt;N&gt;.md</b><br>Independent reviewer verdict"]
    Fork -->|Multi-agent parallel worktrees| Obj["<b>objectives/</b><br>Split directory for parallel git branches"]
    Fork -->|Multi-day mission requiring owner gate| Cp["<b>checkpoints/</b><br>Durable stop point in long campaign"]
```

| Record Type | **NEVER** Create For (90% Routine) | **ONLY** Create When (The Earned Exception) |
|---|---|---|
| **`evidence/`** | Normal code changes, refactors, or local bug fixes. | **Third-party receipts**: You called Stripe, AWS, or an external API and need durable proof of the HTTP response that local tests cannot reproduce. |
| **`reviews/`** | Individual routine tasks. | **Milestones & High Stakes**: Batch-reviewing a finished 4-mission Campaign milestone, or high-risk auth/payments/zero-downtime DB migrations. |
| **`handoffs/`** | Delegating a task to a worker subagent in the same session. | **Permanent Transfer**: Handing the project keys to another human or switching runtimes (e.g. Claude $\to$ Antigravity). |
| **`objectives/`** | Sequential steps inside a standard single-file mission. | **Parallel Worktrees**: Two different agents working in separate git branches/worktrees simultaneously on the same mission. |
| **`checkpoints/`** | Checking in after an objective finishes (use a simple run note). | **Multi-Day Governance Gates**: A high-stakes mission where work pauses for human/legal sign-off before irreversible operations. |

---

### D. Dual-Lane Orchestration: Supervised Dispatch vs. Full Handoff

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
- Decisions: `D11-proposal-retirement`, `D15-branch-guardrail-at-activation`, `D24-schema-field-mechanically-governs-frontmatter`, `D28-dynamic-operating-dial-and-anchors`
