---
type: Atlas
title: Spectacular domain overview
---

# Atlas: Spectacular domain overview

This is the visual companion to [VOCABULARY.md](../VOCABULARY.md). It shows the
whole domain at a legible level; the Vocabulary defines every concept and wins
if the two documents differ.

## Domain board

```mermaid
flowchart LR
  subgraph Definition[Project definition]
    Project[Entity: Project]
    Anchor[Entity: Anchor]
    Vocabulary[Entity: Vocabulary Anchor]
    Contract[Entity: Contract]
  end

  subgraph Planning[Planning]
    Proposal[Entity: Proposal]
    Atlas[Entity: Atlas]
  end

  subgraph Execution[Governed execution]
    Owner[Actor: Owner]
    Mission[Entity: Mission]
    Objective[Entity: Objective]
    Run[Entity: Run]
    Authority[Policy: Authority]
    Start[Action: Start Mission]
    Activated[Event: Mission activated]
  end

  subgraph Proof[Proof and continuity]
    Decision[Entity: Decision]
    Decide[Action: Decide]
    Evidence[Entity: Evidence]
    Review[Entity: Review]
    Handoff[Entity: Handoff]
    Gap[Entity: Gap]
  end

  Owner -->|performs| Decide
  Decide -->|emits| Decision
  Owner -->|requests| Start
  Start -->|emits| Activated
  Project -->|has 1..*| Anchor
  Project -->|has 0..1| Vocabulary
  Project -->|has 0..*| Contract
  Mission -->|governed_by| Contract
  Proposal -->|references| Atlas
  Vocabulary -->|references| Atlas
  Mission -->|contains 1..*| Objective
  Objective -->|has 0..*| Run
  Mission -->|governed_by| Authority
  Mission -->|has 0..*| Evidence
  Mission -->|has 0..*| Review
  Mission -->|has 0..*| Handoff
  Mission -->|has 0..*| Gap
```

## Legend

| Type | Rendered as | Purpose |
| --- | --- | --- |
| Bounded Context | Mermaid subgraph | Names the context in which terms are interpreted. |
| Actor | `Actor:` node | Person, organisation, or agent that participates. |
| Entity | `Entity:` node | A thing with identity and lifecycle. |
| Value Object | `Value:` node | A defined value without independent identity. |
| Action | `Action:` node | A meaningful permitted operation. |
| Event | `Event:` node | A fact that occurred. |
| Policy / Invariant | `Policy:` node | A rule governing states or actions. |
| External System | `External:` node | A dependency outside the domain. |

Relationships are labelled edges, never nodes. Use the default labels defined
in `VOCABULARY.md`; place `1`, `0..1`, `1..*`, or `0..*` on an edge when the
cardinality is important. No External System is part of this overview today.

## Open questions

- After several projects use this model, which ontology-impact checks are stable
  enough to become warning-only conformance checks?
