---
type: source-card
source: source-014
provided_as: sourceL-reversible-architecture-and-adr-process
retrieved: 2026-08-07
authority: user-provided-unsourced-proposal
status: ingested
scope: [decision-governance, reversibility, spikes, adr, product-contracts, architecture, interface-compatibility]
completeness: supplied-text-complete
---

# Source 014 — Reversibility-calibrated architecture decisions

## Thesis

Spend decision effort according to reversal cost: move quickly on low-consequence choices, use
bounded evidence and explicit ownership for consequential choices, preserve architectural rationale,
and make shared interface contracts mechanically checkable. Product intent should settle before
implementation architecture, while later decisions supersede rather than erase durable history.

## What is new

1. **Reversibility as a routing input:** estimate blast radius, migration burden, and rollback cost
   before selecting decision ceremony or evidence fidelity.
2. **Advice with accountable ownership:** one authorized owner seeks affected expertise, then owns
   the call and its consequences without requiring unanimous agreement.
3. **Problem/contract before mechanism:** settle desired behavior and constraints before recording
   an implementation architecture choice.
4. **Optionality-aware architecture:** prefer simple foundations and introduce seams where an actual
   volatile or external boundary justifies them.
5. **Executable interface compatibility:** schema checks and breaking-change detection can turn
   selected shared contracts into deterministic gates.

## Existing concepts reinforced

- PZL-064 and PZL-083 — compact contracts with outcomes, failures, interfaces, and evidence.
- PZL-085 and PZL-087 — distinct product authority and engineering assurance.
- PZL-095 and PZL-112 — predeclared evidence and disposable uncertainty spikes.
- PZL-111 — ordinary modularity and maintainability remain anti-slop constraints.
- PZL-115 and PZL-155 — consequence-calibrated gates and uncertainty-shaped artifacts.
- PZL-145 and PZL-146 — durable decision lifecycle and authorization attribution.

## Corrections and evidence limits

| Proposal | Intake correction |
|---|---|
| Type 1 versus Type 2 | Reversibility is a spectrum and can change with scale, data, contracts, adoption, and migration state. Use a rubric, not a permanent binary label. |
| Decide Type 2 choices in 30 minutes | Time is an arbitrary proxy. Cap effort proportionally, but require enough evidence for safety, privacy, compatibility, and operational consequences. |
| One-day or one-sprint spike | The timebox follows the cheapest experiment that can discriminate options; predeclare the question, metric, disposal boundary, and stop rule. |
| API payload format as a reversible example | A public or persisted payload can become a one-way door through consumers and stored data. Classification is contextual. |
| Abstract dependencies to preserve optionality | Speculative adapters can add indirection and maintenance cost. Add a seam at a demonstrated volatile/external boundary. |
| Monolith first and boring technology | Useful defaults, not universal laws. Team fit, constraints, failure modes, and evidence determine the choice. |
| One source of truth plus a formal pipeline guarantees alignment | No process guarantees alignment. Product contracts, implementation decisions, schemas, and runtime evidence have distinct authorities that must link without duplication. |
| Put UX, domain, data, KPI, SLO, and requirements into one PRD | This can overload one artifact. Keep product intent compact and promote technical detail into capability contracts or decision records when earned. |
| Never edit an accepted ADR | Preserve material reversals through supersession; permit clearly audited factual/link corrections according to policy rather than duplicating a record for every typo. |
| Disagree and commit | A decision owner needs legitimate delegated authority, affected-party advice, safety escalation, and a visible appeal path; ownership does not override law, policy, or separate authority. |
| Schema-first parallel work prevents blocking | Mocks can accelerate work but drift from implementation. Generated artifacts and compatibility checks must derive from one accepted schema. |
| Example throughput and latency improvements | Illustrative numbers are not evidence for a real architecture choice. |

## Foundation Plan impact

The top-20 order remains stable. Source 014 strengthens S01's decision-effort rubric, S03's layered
truth/contract distinction, S05's human-authority model, S09's public-contract compatibility, and
S11's architecture/migration posture. The Decision Multiplexer should include reversibility and
blast radius in its routing packet.

## New concept pieces

Source 014 adds PZL-167 through PZL-171. Repeated material is merged as provenance into the existing
cards above rather than counted as independent corroboration.
