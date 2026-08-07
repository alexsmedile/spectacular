---
type: source-card
source: source-012
provided_as: sourceX-decision-density-synthesis
retrieved: 2026-08-07
authority: user-provided-unsourced-proposal
status: ingested
scope: [decision-method, prototypes, grilling, logic-harness, verification, wayfinding, skill-ecosystem]
completeness: supplied-text-complete
---

# Source 012 — Compacting decisions through high-fidelity artifacts

## Thesis

Human decisions can sometimes be resolved more efficiently by reacting to comparative, executable,
or stressed artifacts than by inventing behavior through long conversational prose. The source
proposes multi-variant prototypes, targeted grilling, executable logic harnesses, evaluator loops,
and ambiguity-aware wayfinding as a five-command “Decision Density Skill Pack.”

## Strong contributions

1. **Decision-interface selection:** use the artifact that exposes the uncertain property directly.
2. **Deliberate option contrast:** variants should encode materially different trade-offs, not cosmetic
   restyling.
3. **Backend prototyping:** state, recovery, concurrency, and transformation rules can be made
   inspectable through a disposable CLI/simulator, not only UI mockups.
4. **Pre-human quality filtering:** deterministic checks and bounded independent review can remove
   obvious breakage before asking for scarce human attention.
5. **Ambiguity-calibrated routing:** unresolved product/architecture choices need early interaction;
   bounded reversible execution can proceed toward an evidence-backed review point.

## Claims requiring correction

| Source claim | Intake correction |
|---|---|
| “Never ask how this should work in plain text” | Text is appropriate for abstract policy, authority, naming, and trade-offs; use executable artifacts only when they improve fidelity. |
| Always build three variants | Variant count should follow meaningful alternatives and comparison cost; one validated default or two polar options may be better. |
| Resolve 10–20 decisions per pass | The numbers are unsourced and raw decision count is gameable; measure accepted, retained decisions and reversal/error cost. |
| One question per turn and 10–15 decisions in one pass | These are different interaction modes. Preview a coherent option matrix, then traverse blocking branches sequentially where answers are dependent. |
| Continue grilling until 100% alignment | Complete alignment is not observable; stop when remaining ambiguity is non-material, explicitly assumed, deferred, or covered by a stop condition. |
| Critic loops until a high bar is reached | Review-repair must have a rubric, severity threshold, attempt budget, unchanged-finding detection, and escalation state. |
| Evaluators auto-resolve design choices | Reviewers may detect violations and defects; they must not silently decide product trade-offs reserved for the human. |
| Five commands imply five skills | Commands, modes, agents, and skills require different ownership tests; standalone value and substrate ownership must be proven. |

## Capability mapping

| Proposed command | Existing concepts | Refactor interpretation |
|---|---|---|
| `/prototype` | PZL-112, PZL-132, PZL-151 | Discovery mode or companion capability; variants are an earned strategy, not a universal rule. |
| `/grooming` | Current grill engine, PZL-133, PZL-152 | Decision-session interaction pattern; possible specwright capability, not automatically a separate skill. |
| `/harness` | PZL-095, PZL-110, PZL-112, PZL-153 | Disposable logic-evidence artifact; the production harness remains a different concern. |
| `/gauntlet` | PZL-104, PZL-114, PZL-127, PZL-154 | Bounded pre-human evidence gate using deterministic checks and risk-triggered reviewers. |
| `/wayfinder` | PZL-103, PZL-113, PZL-155 | Ambiguity/dependency routing; the current Wayfinder product boundary remains unsettled. |
| `/dd-skills` pack | PZL-065, PZL-123, PZL-156 | Composition proposal to evaluate after responsibility and Mission contracts, not a five-skill authorization. |

## Impact on the Foundation Plan

Source 012 does not reorder the top 20. It strengthens four session methods:

- S01 should consider **validated decisions per unit of human attention** as a secondary metric,
  guarded by decision retention, correctness, and reversal cost.
- S04/S07 must classify prototype, grill, harness, critic, and wayfinding as modes, roles, adapters,
  or skills by responsibility—not command spelling.
- S06 should decide whether a bounded pre-human quality gate protects attention without transferring
  product authority to evaluators.
- Every decision session should choose its interaction artifact from the uncertainty being resolved.

## New concept pieces

Source 012 adds PZL-150 through PZL-156. It also reinforces PZL-104, PZL-110, PZL-112–115,
PZL-121, and PZL-127 without providing independent empirical corroboration.
