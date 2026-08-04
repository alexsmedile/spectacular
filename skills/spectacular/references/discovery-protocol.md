---
description: Progressive routing for direct clarification, research, Vision, technical spikes, prototypes, production tracer bullets, ambiguous experiments, and technical debt.
when_to_use: Choosing whether uncertainty needs Vision reaction, research, a spike, a prototype artifact, a tracer bullet, direct implementation, or no new node.
---

# Discovery Evidence Protocol

Use discovery to remove a named uncertainty, never as automatic ceremony. If direction and implementation path are already clear from code, tests, documentation, or direct user clarification, proceed through the appropriate direct/specification/request route without creating Vision, research, spike, prototype, or finding records.

## Cheapest sufficient answer first

Ask these in order and stop at the first sufficient route:

| Uncertainty | Route | Owner | Code destination | Durable output |
|---|---|---|---|---|
| None: goal, constraints, and implementation path are clear | Direct execution | Approved `SPC` + request | Production branch through normal Git/HITL gates | Code, tests, verification evidence |
| Product, interaction, workflow, or system direction is unsettled and human reaction would change the contract | Vision | `.spectacular/visions/<slug>/VISION.md` | None; optional showable artifacts only | Approved/rejected direction, fragment reactions, linked evidence; approved may derive draft SPC |
| A bounded fact, external behavior, or option comparison is unknown | Research | `RES-NNN` | None; read-only investigation | Sources, evidence, options, `supported|refuted|inconclusive` result |
| Technical feasibility or a risky integration assumption must be tested | Spike | `SPK-NNN` | Disposable `spike/prototype-<SPK-ID>` branch | Hypothesis, experiment, evidence, result; not the code |
| A proposed UI, CLI, workflow, or experience needs to become reactable before its contract is fixed | Prototype fragment | Owning Vision (or feedback/request for an already-built target); `PRT` remains reserved | Mock/sandbox or owner-local artifacts | Human reaction and its decision impact; artifact is not the direction by itself |
| The architecture is chosen but end-to-end production wiring is the next risk | Tracer bullet | Approved `SPC` with `execution_mode: tracer` + request | Production-quality feature branch; retained and extended | Thin integrated code, tests, verification evidence |

Escalation is progressive: direct clarification or inspection first; Vision only for material direction uncertainty; `RES` for bounded facts; `SPK` only when observation must settle feasibility; prototype only when interaction is the evidence; tracer only after an approved SPC authorizes production work.

“Experiment” is deliberately not a command alias. Route it by the question:

- fact or option comparison → research;
- technical feasibility → spike;
- pre-spec human reaction → Vision fragment/prototype;
- learning from built behavior → feedback-loop.

## Entry contract

Every Vision-linked research record, spike, or prototype artifact must state:

- one question or falsifiable hypothesis;
- why existing code, tests, docs, or a direct user answer are insufficient;
- owner and dependency being cleared;
- time/effort boundary and exit condition;
- evidence to preserve;
- where each outcome routes next.

Do not create a node for open-ended “exploration.” Split broad uncertainty into the smallest question whose answer changes a decision, specification, or execution path. `inconclusive` remains fog and must link a narrower successor; it is not permission to continue indefinitely.

## What survives

- **Research:** the `RES` record is already the durable evidence artifact. Do not copy every conclusion into `findings/`.
- **Spike:** delete/archive-isolate experimental code according to [[afk-git-hygiene]]; keep the `SPK` evidence and result.
- **Vision/prototype:** keep the approved/rejected fragments and human observations needed to explain the chosen direction and resulting specification. Mock code is replaced unless explicitly promoted through a newly approved production plan.
- **Tracer bullet:** keep the code. It is a deliberately thin production slice, not a prototype, and must meet the same architecture, security, testing, and review standards as later feature work.

The experimental branch/mock is throwaway; its outcome record is stale-safe history. Apply [[artifact-retention]] before cleanup so knowledge survives while junk code does not.

Research, spikes, vendor evidence, and user interviews may **support** a `DEC`, but never create the choice by themselves. A decision still requires explicit user choice, unambiguous recorded user intent, or the narrow authorized AFK gate in [[lifecycle-contract]]. Unsettled alternatives remain a `QUE` with options, recommendation, and evidence.

`FND-NNN` remains reserved for a future workflow that extracts a reusable, cross-cutting takeaway from several sources or sessions. Until that workflow ships, prefer the originating `RES`, `SPK`, audit, question, decision, memory, or request instead of duplicating the same learning.

## Artifact means output, not entity

“Artifact” is an umbrella term for an owned output: a source list attached to `RES`, logs attached to `SPK`, a mock attached to a request, a benchmark attached to its run, or a verified `DEC`/`SPC`. It does not receive an `ART-NNN` identity and does not justify a project-wide catch-all folder. Put an artifact beside the record or request whose claim it substantiates.

## Technical debt and mocks

A mock is not automatically debt. In a spike or prototype it is disposable evidence and its owner records the cleanup/disposition. A mock that enters production creates an explicit obligation:

| Situation | Route |
|---|---|
| Required remediation is inside the approved active scope | A named request task with verification |
| Concrete remediation is likely soon but not planned yet | Roadmap `candidate` with `target-version: tbd`, then a request when planned |
| Possible improvement has no commitment or schedule | `IDEA-NNN` |
| A deliberate compromise must be understood later | `DEC-NNN` records why, consequences, and links to the cleanup owner |
| Current behavior is defective or regressed | Existing bug workflow: audit/request/fix as warranted |

Do not create a `debt/` collection or `DEB-NNN` merely to hold an issue. Debt needs an execution owner and priority; a passive parallel backlog would duplicate requests and the roadmap. Reconsider a dedicated debt lifecycle only after real project evidence shows these routes lose obligations.

## Boundary checks

- Vision is opt-in and pre-request; it cannot silently expand an active milestone. Tangents route to an idea or future target.
- Disposable code never merges because it “already works.” Production adoption starts from an approved spec and clean implementation request.
- A tracer bullet is the smallest real vertical slice, not a shortcut around production quality.
- Do not create `RES`, `SPK`, prototype artifacts, and `FND` for the same learning by default. Preserve one authoritative record and link downstream decisions/specs to it.
- AFK may run read-only research and explicitly authorized spikes within goal scope; user/business choices and declared HITL boundaries still stop execution.

**Related:** [[research-rules]], [[spike-rules]], [[lifecycle-contract]], [[soft-db-index]], [[wayfinding-sequencer]], and [[afk-git-hygiene]].
