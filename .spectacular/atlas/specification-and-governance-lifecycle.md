---
type: Atlas
title: Specification evolution and governance lifecycle
---

# Atlas: Specification evolution and governance lifecycle

## Outcome board

| Actor | Journey step | Desired outcome | Success signal |
| --- | --- | --- | --- |
| Founder / Owner | Brainstorm draft specifications | Iterate rapidly on ideas and requirements without fear of breaking active execution state | Proposals can be freely rewritten across dozens of chat turns |
| Architect / Lead | Settle critical architectural forks | Record permanent rationale and trade-offs that future agents and humans will not re-debate | `spectacular decide` records an immutable UUIDv7 Decision receipt |
| Developer / Agent | Implement observable features | Work against frozen, verifiable claims with exact Git baselines and failable proofs | Mission execution produces attributable Evidence and independent Reviews |
| Team / Maintainer | Discover accepted system truth | Find living subsystem capabilities and system boundaries instantly | Capability Contracts (`contracts/CC-*.md`) reflect current accepted reality |

```mermaid
flowchart LR
  B["💡 Brainstorm & Draft Specs<br><i>(Proposals & Conversations)</i>"] --> D["⚖️ Lock Architectural Forks<br><i>(Decisions: D&lt;N&gt;)</i>"]
  D --> C["📜 Living Subsystem Truth<br><i>(Capability Contracts & Anchors)</i>"]
  C --> M["🚀 Frozen Implementation Slice<br><i>(Missions: M&lt;N&gt;)</i>"]
  M --> E["✅ Verified Evidence & Closure<br><i>(Retire Proposal & Freeze Proof)</i>"]
```

## System board

| Capability | Connection | Implementation boundary | Proof / risk |
| --- | --- | --- | --- |
| Exploratory scratchpad | enables `Brainstorm draft specifications` | `.spectacular/proposals/P<N>.md` (mutable, zero authority) | Zero drift errors during rapid conversational rewrites |
| Atomic architectural rulings | enables `Settle critical architectural forks` | `spectacular decide` (`.spectacular/decisions/D<N>.md`) | Idempotent two-phase transaction journal |
| Living modular specifications | enables `Discover accepted system truth` | `.spectacular/contracts/CC-<name>.md` (`contract_version`) | In-flight Mission editing with gap amendment tracking |
| Frozen execution envelopes | enables `Implement observable features` | `spectacular mission start` (`.spectacular/missions/M<N>/`) | SHA-256 semantic fingerprint and branch guardrail |
| Retired proposal tracking | enables `Verify Evidence & Closure` | `.spectacular/archive/proposals/` (`resolved_by:`) | Transparent historical trace from idea to shipped proof |

```mermaid
flowchart TD
  subgraph Ideation ["1. Exploration & Brainstorming"]
    P["Proposal (P&lt;N&gt;)<br>• Mutable draft spec<br>• Open questions<br>• Alternatives"]
  end

  subgraph Rulings ["2. Architectural Choice"]
    DEC["spectacular decide<br>• Immutable Decision (D&lt;N&gt;)<br>• Attributable owner rationale"]
  end

  subgraph SystemTruth ["3. System Truth"]
    ANC["Anchors (PROJECT.md)<br>• System boundaries<br>• Architecture layers"]
    CC["Capability Contracts (CC-*)<br>• Observable behaviors<br>• Modular invariants"]
  end

  subgraph Execution ["4. Verifiable Delivery"]
    M["Mission (M&lt;N&gt;)<br>• Frozen claims<br>• Failable proof<br>• Evidence & Review"]
  end

  P -->|Owner alignment| DEC
  DEC --> ANC
  DEC --> CC
  CC -->|Authorize change| M
  M -->|Ships & Proves| ARC["Archive Proposal<br>(status: accepted)"]
```

## Links and references

- Product Documentation: [`docs/architecture.md`](../../docs/architecture.md), [`docs/process.md`](../../docs/process.md)
- Agent Skill Guidance: [`skills/spectacular/references/prepare.md`](../../skills/spectacular/references/prepare.md)
- Decisions: `D11-proposal-retirement`, `D15-branch-guardrail-at-activation`, `D24-accepted`
