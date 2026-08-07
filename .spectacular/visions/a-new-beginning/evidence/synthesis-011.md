---
type: synthesis-checkpoint
checkpoint: synthesis-011
sources: 14
concepts: 171
github_issue_cards: 23
human_dispositions: 0
status: intake-working-synthesis
updated: 2026-08-07
---

# Synthesis 011 — reversible decisions without mixed authority

## What Source 014 changes

Source 014 makes **reversal cost** an explicit part of the refactor method. Every consequential
decision packet should now state impact, reversibility, migration burden, uncertainty, and required
evidence before choosing its review ceremony. This refines risk-calibrated gates and the Decision
Multiplexer's artifact routing; it does not create another lifecycle or document type.

It also supplies a useful governance shape:

```text
authorized owner → affected/expert advice → evidence and dissent → decision or escalation
```

The owner is accountable but not omnipotent. Product value, security, provider operations,
irreversible effects, and policy may have separate legitimate authorities.

## Authority model clarified

The source's proposed single PRD/RFC/ADR pipeline is decomposed into linked authorities:

| Artifact or evidence | Owns | Does not own |
|---|---|---|
| Product intent/Anchors | Problem, users, outcomes, non-goals, business constraints | Implementation architecture |
| Capability Contract | Accepted behavior, invariants, failures, and proof expectations | Current implementation reality |
| Decision record | Chosen mechanism, drivers, alternatives, consequences, decider, and supersession | Product approval by implication |
| Interface schema | Declared machine-readable input/output compatibility | Complete semantic or operational correctness |
| Code/tests/operations | Implemented and observed reality | Intended product truth by itself |

This preserves one authority per fact without pretending one mega-document can be the source of
truth for every layer.

## Architecture posture

“Reversible architecture” becomes an **earned-seam** rule:

- add a narrow boundary when volatility, external ownership, migration cost, or replacement risk is
  demonstrated;
- otherwise prefer the simplest cohesive module that satisfies the accepted contract;
- treat “monolith first,” “boring technology,” and adapters as defaults to test, not laws;
- use disposable spikes only with a predeclared question, discriminator, timebox, and disposal gate;
- mechanically check declared shared-interface compatibility when the interface is stable enough to
  justify schema ownership.

Schema diffs and generated mocks remain projections. They cannot prove behavioral, security, or
operational compatibility and must derive reproducibly from one accepted schema.

## Foundation Plan impact

The twelve sessions and top-20 order remain unchanged, but their contracts are sharper:

- S01 defines the impact/reversibility rubric.
- S03 preserves the product-contract-before-architecture order and layered authorities.
- S05 defines accountable advice, dissent, and escalation.
- S09 includes public-interface schema compatibility.
- S11 identifies earned seams and executable compatibility gates.

The Decision Multiplexer experiment now receives four required routing fields: impact,
reversibility, uncertainty, and evidence threshold.

## Current state

- Sources ingested: 14.
- Atomic concepts: 171.
- GitHub issue evidence cards: 23.
- Human dispositions: 0.
- Promoted specifications: 0.
- Next source: `source-015`.

The source strengthens the decision method. It does not approve an ADR subsystem, a new canonical
PRD template, schema-first development everywhere, or any specific architecture choice.
