---
type: concept-piece
id: PZL-013
status: captured
domain: specification-model
sources: [source-002, source-006, source-007]
source_authority: proposal
assessment: promising
evidence_status: supported
disposition: pending
depends_on: []
overlaps_with: [PZL-003, PZL-006, PZL-064, PZL-065, PZL-080, PZL-081]
conflicts_with: []
tags: [specifications, contracts, capability, intent]
updated: 2026-08-07
---

# Capability Contracts

## Core message

Name the convergent capability-level specification a Capability Contract to make
its behavioral promise and approval role explicit without claiming it is runtime truth.

## Value

Could distinguish accepted behavioral intent from exploratory ideas, execution
plans, and implementation evidence.

## Assumptions

- “Contract” communicates shared agreement rather than legal rigidity.
- Every spec represents a capability rather than a schema, axis, migration, or policy.

## Evidence and collisions

Current specs already function as implementation contracts and requests persist a
`contract` UUID. The exact phrase is absent. Flat SCHEMA and AXIS contract docs
show that not every specification is naturally capability-shaped. Source 006
supplies a concrete end-to-end schema and makes capability contracts the central
truth updated at mission closure.
Source 007 keeps the capability contract but explicitly separates it from component,
interface, state, operational, and policy contracts linked through a graph.

## Trade-offs and recommendation

Clearer authority versus a potentially narrower taxonomy and migration churn.
Compare “Capability Contract” as a user-facing definition, a subtype, and a full
rename before selecting scope. Promising clarification.

## Decision

Pending.
