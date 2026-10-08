---
name: system-architecture
description: >-
  Architect software systems by exploring options, making focused architecture decisions, designing target structures,
  reviewing proposals, documenting C4 diagrams, planning evolutionary migrations, and recording ADRs.
  Triggers on "system architecture", "architecture review", "C4 diagram", "ADR", "service boundaries", "bounded context",
  "distributed trade-offs", "modernization plan", or "system design".
  Do not invoke for local coding refactors, isolated bug fixes, database schema DDL, or UI styling without architectural consequence.
---

# System Architecture

Make architecture decisions traceable to business outcomes, measurable quality attributes, constraints, and verified evidence.

## Route Matrix

| Route | Required Workflow Steps | Deliverable | Complete When |
|---|---|---|---|
| **Explore** | 1–3, 7 | Trade-off matrix & options | Drivers explicit; all viable options compared against same drivers. Read [decision guide](references/decision-guide.md). |
| **Decide** | 1–3, 7 | Single ADR recommendation | Evidenced choice with drivers, alternatives, and validation step. Read [templates](references/deliverable-templates.md). |
| **Design** | 1–8 | Proportional target architecture | All selected steps satisfy completion criteria; passes quality gate. Read [architecture method](references/architecture-method.md). |
| **Review** | 1, assessed 2–6 | Findings & verdict | All areas cite verified evidence or marked unsupported with severity. Read [review template](references/deliverable-templates.md). |
| **Document** | 1, 7, (8 if migrating) | ADR, arch doc, diagram set | Consistent naming, decision status preserved, zero contradictions. Read [templates](references/deliverable-templates.md). |
| **Evolve** | 1–3, 7–8 | Phased migration plan | Each phase has compatibility, rollback, and exit criteria. Read [visual delta](references/visual-communication.md). |
| **Explain visually** | Audience + boundary + checks | Explanatory / technical visuals | Answers one named question; verified against inspected evidence. Read [diagram patterns](references/diagram-patterns.md). |

## Core Workflow Steps

1. **Frame the problem:** Business outcomes, actors, system boundaries, constraints, non-goals, measurable quality scenarios.
2. **Model domain & ownership:** Capabilities, bounded contexts, invariants, authoritative owners (distinguish logical vs deployment). For module design/review, apply the [module boundary review](references/architecture-method.md#module-boundary-review); distinguish declared, implemented, and verified boundaries.
3. **Compare system shapes:** Evaluate simplest viable shapes against ranked drivers (delivery, operational cost, failure modes).
4. **Trace runtime behavior:** Happy path + failure modes (timeouts, retries, idempotency, backpressure, degraded modes, recovery).
5. **Design data ownership:** Transaction boundaries, consistency semantics, authoritative writers.
6. **Design operations & security:** Trust boundaries, IAM, secrets, encryption, observability, capacity, deployment, disaster recovery.
7. **Record decisions & uncertainty:** ADR capturing chosen option, drivers, consequences, confidence, and open validation questions.
8. **Plan evolution:** Thin delivery slices with backward compatibility, rollback steps, and verified acceptance gates. End a delivery design with the minimum path to the first real workflow, its dependencies and actual authorization boundaries; distinguish technical feasibility from usable-product acceptance.

## C4 Abstraction Levels & Quick Pattern

| Level | Boundary Scope | Audience | Diagrams-as-Code Syntax |
|---|---|---|---|
| **L1: Context** | System of interest + external actors/systems | Everyone | `flowchart LR; user["User"] --> sys["System"] --> ext["Ext API"]` |
| **L2: Container** | Deployable units, datastores, queues | Engineers & Ops | `subgraph sys; web["App"] --> api["API"] --> db[("DB")]; end` |
| **L3: Component** | Internal modules & controllers | Module Owners | Class/module interaction within a single container |

```mermaid
flowchart LR
    customer["Person: Customer"]
    subgraph platform["Software System: Platform"]
        web["Container: Web App"]
        api["Container: Core API"]
        db[("Container: Database")]
    end
    customer -->|"HTTPS"| web -->|"JSON/HTTPS"| api -->|"SQL/TLS"| db
```

## Module responsibilities and evidence

For each material module, identify its reason to change, owned data and writers,
public operations/read views, hidden details, allowed dependency direction,
external effects, and transaction invariants. Attach an executable check or state
that enforcement is absent. Keep shared types minimal; use consumer-defined ports
where they remove a concrete dependency. Preserve coordinated transactions even
when implementations move into separate files. File size is a signal, not proof
of a failed boundary.

Use existing project owners and artifacts. Pass only the relevant boundary facts
to companion skills; do not duplicate their procedures or load the whole family.
Software interface contracts do not activate governed Spectacular Contracts or
Missions.

## Delivery prerequisites

For substantial delivery or launch design, consult the relevant prompts in
[Spectacular readiness](../spectacular/references/milestone-readiness.md). Classify
prerequisites by first workflow, real exposure, and wider distribution; explain
the failure each prevents. Choose the simplest integrated path. Do not infer auth,
queues, multi-region infrastructure, or deployment gates from an MVP label.

## Expansion Handoffs

| Out-of-Scope Need | Action / Delegate |
|---|---|
| Executable dependency, isolation, compatibility, and transaction proof | Use `test-sentinel` with the boundary and expected refusal |
| Typed operations across module or external interfaces | Use `action-contracts` with the owner, consumers, and effects |
| Physical DDL, ER diagrams, indexing, zero-downtime migrations | Invoke `data-modeling` companion skill |
| 3-option tracer spike on ambiguous UI/architecture variants | Invoke `rapid-prototyping` companion skill |
| Durable delivery context or explicitly requested governance | Use `spectacular` ordinary plans by default; Mission procedures only on owner opt-in |

## Core Invariants & Negative Constraints

- **DO NOT make architectural claims without inspecting evidence.** In codebases, inspect source, configs, schemas, and telemetry before asserting state; label unverified statements as assumptions.
- **DO NOT produce monolithic all-or-nothing designs.** Structure architectures in independently verifiable, reversible slices.
- **DO NOT mix architectural decisions with ungrounded speculation.** Explicitly separate *Facts*, *Assumptions*, *Decisions*, and *Open Questions*.
- **DO NOT add unnecessary distributed complexity.** Default to the simplest monolithic or modular architecture unless scale, team boundaries, or isolation constraints force decoupling.
