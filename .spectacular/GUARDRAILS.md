---
type: Anchor
version: '0.3'
created: '2026-08-10T21:14:12+02:00'
updated: '2026-10-07T20:07:46Z'
---

# Spectacular project guardrails

@Orient @Prepare @Start @Resume @Run @Assess @Reconcile @Resolve

## @Ordinary Context
- Start authorized work from linked context or plans; no Mission is inferred from folder presence or risk.
- Store project plans in plans/; native host plan-mode restrictions still apply.
- Use optional specs/ and audits/ only for independently useful content.

## @Core Product & Identity Invariants
- Preserve UUID identity, exact revision fingerprints, source drill-down, owner authority, provider boundaries, and recoverable writes.
- Prefer the smallest human-readable workspace structure that keeps those invariants visible.
- Do not add v1 compatibility, generic record/search commands, an authoritative projection, or a second product root.

## @Alignment & Domain Ontology
- All domain actions and entity states must adhere strictly to canonical terms defined in `ONTOLOGY.md`. Using Banned Synonyms is an invariant violation.
- Missions must state explicit `Ontology impact` during preparation (D27).

## @Architecture & Pattern Discipline
- Before drafting bespoke implementations, perform an upfront Pattern Pass (D29) surveying standard library idioms, RFCs, and proven reference implementations.

## @Execution & Safety
- Parallel subagents must operate under disjoint `writes:` reservations (D21). Overlapping write perimeters are forbidden.
- Retain independently useful learnings where they belong; update guardrails only when an accepted constraint changes.
