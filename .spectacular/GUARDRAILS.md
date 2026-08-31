# Spectacular project guardrails

@Orient @Prepare @Start @Resume @Run @Assess @Reconcile @Resolve

## @Core Product & Identity Invariants
- Preserve UUID identity, exact revision fingerprints, source drill-down, owner authority, provider boundaries, and recoverable writes.
- Prefer the smallest human-readable workspace structure that keeps those invariants visible.
- Do not add v1 compatibility, generic record/search commands, an authoritative projection, or a second product root.

## @Alignment & Domain Ontology
- All domain actions and entity states must adhere strictly to canonical terms defined in `VOCABULARY.md`. Using Banned Synonyms is an invariant violation.
- Missions must state explicit `Ontology impact` during preparation (D27).

## @Architecture & Pattern Discipline
- Before drafting bespoke implementations, perform an upfront Pattern Pass (D29) surveying standard library idioms, RFCs, and proven reference implementations.

## @Execution & Safety
- Parallel subagents must operate under disjoint `writes:` reservations (D21). Overlapping write perimeters are forbidden.
- Post-mission mistake learnings must be codified directly into this file or `AGENTS.md` upon closure (D29).
