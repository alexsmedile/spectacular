---
type: derived-view
status: experimental
authority: none
source_of_truth: concepts/PZL-*.md
checkpoint: synthesis-006
updated: 2026-08-07
---

# Spectacular abstraction map

This view applies Source 010's abstraction question to Spectacular. It is an audit
projection, not accepted architecture. Linked PZL cards retain authority.

```mermaid
flowchart TB
  subgraph INTENT["1 · Intent — why and what"]
    ANCHORS["PZL-012 · Anchors"]
    CONTRACTS["PZL-013 · Capability contracts"]
    MISSION["PZL-066 · Mission delta"]
  end

  subgraph ORCH["2 · Orchestration — judgment and bounded context"]
    ROUTER["PZL-001 · Lean skill router"]
    CONTEXT["PZL-070 · Run-context compiler"]
    HANDOFF["PZL-119 · Typed handoff"]
  end

  subgraph MECH["3 · Mechanism — deterministic operations"]
    HELP["PZL-002 · CLI help authority"]
    REGISTRY["PZL-055 · Command registry"]
    HARNESS["PZL-110 · Deterministic harness"]
  end

  subgraph TRUTH["4 · Truth and evidence"]
    PRODUCT["PZL-019 · Product truth"]
    EVIDENCE["PZL-073 · Clause evidence"]
    STATE["PZL-118 · Resume state"]
  end

  subgraph ADAPTERS["5 · External adapters"]
    PROVIDER["PZL-054 · Native provider boundary"]
    RUNTIME["PZL-069 · Host agent runtime"]
    SEMANTIC["PZL-124 · Semantic query adapter"]
  end

  ANCHORS --> CONTRACTS --> MISSION
  MISSION --> ROUTER --> CONTEXT --> HANDOFF
  HANDOFF --> HARNESS
  HELP --> REGISTRY --> HARNESS
  PRODUCT --> CONTRACTS
  HARNESS --> EVIDENCE --> STATE
  CONTEXT -. "references" .-> PRODUCT
  HARNESS --> PROVIDER
  HARNESS --> RUNTIME
  HARNESS -. "optional" .-> SEMANTIC
```

## Boundary audit

| Boundary | Required interface | Current question |
|---|---|---|
| Intent → orchestration | accepted outcome, authority, constraints, non-goals | Does Mission/request/spec ownership overlap? |
| Orchestration → mechanism | one mechanical operation with validated parameters | Does the skill document CLI behavior the CLI should own? |
| Mechanism → truth | explicit read/write contract, schema, provenance, failure | Does the CLI create or reinterpret product decisions? |
| Truth → orchestration | bounded authoritative projection plus freshness | Can status/run context omit or shadow canonical evidence? |
| Mechanism → adapters | narrow capability, permission, and mutation boundary | Are Git/GitHub and host-agent operations duplicated internally? |

## Audit rule

For each command, reference, agent, and collection, record its owning layer, callers,
authoritative inputs, guaranteed output, deterministic checks, failure mode, and drill-down
path. A component that owns two layers or duplicates authority becomes a refactor candidate;
crossing a layer is not itself a defect.
