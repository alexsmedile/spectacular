---
type: refactor-foundation-contract
version: 1.0
status: accepted
decision_session: S02
accepted_by: owner
accepted_at: 2026-08-08
central_disposition: accept
upstream:
  - PRODUCT-CONSTITUTION.md@1.0
  - TRUTH-PROVENANCE-FLOOR.md@1.0
next_session: S03B
---

# Refactor Success and Evidence Constitution

This is the accepted S02 ruler for the `a-new-beginning` rebuild. It defines
how later proposals will be evaluated. It does not evaluate any subsystem or
choose the target ontology, packaging, interface, scaffold, implementation, or
migration architecture.

## 1. Protected-loop gates come first

A proposal cannot compensate for breaking Spectacular's accepted identity by
being smaller, faster, or more popular. Before comparison, it must preserve:

- accepted intent and boundaries;
- declared authority and human decision gates;
- bounded, accountable work and handoffs;
- evidence-backed acceptance;
- authorized reconciliation of accepted results;
- cold recovery to a safe next action or an explicit unresolved gate.

Failure of any relevant gate is disqualifying until repaired or explicitly
reframed through the appropriate upstream decision process.

## 2. Passing proposals use a lightweight comparative review

There is no universal numeric score or fixed weighting. For each proposal:

1. State the intended gain in plain language.
2. Confirm the protected-loop gates remain intact.
3. Check for material regressions in:
   - owner comprehension and control;
   - cold recovery and resumability;
   - cold-start and context cost;
   - maintainability, inspection, testing, and repair;
   - compatibility and recoverability;
   - practical owner value;
   - human attention and decision burden.
4. Require evidence proportionate to the consequence and reversibility.

Maintainability failure is a real failure even when immediate behavior works.
No gain in one dimension silently offsets an unexamined material regression in
another.

The owner's initial acceptance of “gates first, then a weighted scorecard” was
explicitly refined during S02: gates remain, but the downstream comparison is
qualitative and evidence-backed rather than numerically weighted.

## 3. Evidence must fit the claim

| Claim | Strong practical evidence | Insufficient alone |
|---|---|---|
| Protected behavior works | End-to-end scenario plus focused behavioral checks | File presence or an agent completion claim |
| Cold recovery improved | A cold human or agent performs the same recovery task | Token count |
| Context cost fell safely | Paired task comparison plus correctness and retrieval-miss checks | Lines, tokens, or reference count |
| Maintenance improved | A representative inspection/change/repair exercise plus validation | Fewer files or less code |
| Compatibility survives | A declared history, interface, or recovery path is exercised | A compatibility promise |
| Owner value or comprehension improved | Observation of representative owners performing the task | Usage count or model opinion |
| A subsystem should survive | Distinct protected value, credible need, alternatives, and sustainable cost | Popularity or frequency alone |
| A judgment is preferable | Attributed expert or owner judgment with assumptions and trade-offs | Unattributed advice or confidence |

Consequential claims identify their source and scope and label the evidence as
direct, proxy, attributed judgment, or unknown under the accepted truth and
provenance floor.

## 4. Anti-gaming rules

- No file, line, token, command, step, test, record, or use count is a success
  claim by itself.
- Reductions pair the attractive count with relevant behavior and recovery
  checks.
- Changes to a route or workflow exercise a representative end-to-end outcome,
  not only structural checks.
- Usage is distinct from value and declares its population, observation window,
  and self-hosting versus external-owner scope.
- Missed context, authority errors, confusion, unsafe continuation, failed
  recovery, or lasting maintenance debt are disconfirming evidence. They cannot
  be offset by a better vanity metric.
- Model and harness effects are separated in comparisons by holding the task and
  evidence criteria stable where practical.

## 5. Reversibility controls decision investment

### Type 1 — hard to reverse or high consequence

Requires explicit owner review, a concise rationale, stated evidence, and a
bounded spike when material uncertainty cannot be resolved more cheaply.
Examples include accepted behavior, public interfaces with consumers, stored
data, compatibility commitments, authority boundaries, and recovery-critical
architecture.

### Type 2 — reversible and contained

Is timeboxed, tested at the affected boundary, and recorded only enough to
prevent confusion. It has no material external, durable, or difficult-to-repair
effect.

Classification is contextual and may change with adoption, stored data,
consumers, or operational dependence. Security, privacy, rights, authority, and
irreversible effects always retain their own approval requirements.

## 6. Subsystem survival is predeclared

A subsystem earns continued existence only when it shows, proportionately to
its cost and consequence:

1. a distinct job protecting the governing loop or a clearly different owner need;
2. a credible benefit supported by appropriately labeled evidence;
3. no simpler adequate replacement against the same representative scenario;
4. sustainable maintenance, context, complexity, and recovery cost;
5. a safe disposition path preserving required history, evidence, and recovery.

Later S10 outcomes remain neutral:

- **Keep:** distinct value at sustainable cost.
- **Simplify:** the job remains but its current shape costs too much.
- **Merge or extract:** the job remains distinct but its placement is wrong.
- **Retire:** no distinct protected value or credible need remains, and removal
  or replacement is safe.

No subsystem is evaluated by this S02 contract itself.

## 7. Prototypes and spikes are disposable experiments

Use a prototype or spike when a material uncertainty can change a Type 1 choice
and inspection or existing evidence cannot answer it cheaply.

Declare before work:

- **Question:** which uncertainty is being resolved?
- **Success observation:** what result supports or rejects each option?
- **Boundary:** what is excluded so the experiment cannot become shadow production?
- **Timebox:** what proportionate limit applies?
- **Disposition:** where evidence goes and whether the experiment is discarded
  or separately promoted through normal review.

A prototype may compress a connected set of uncertainties into one human review
pass. The number of questions is descriptive, never a performance target.
Prototype code and visuals are evidence surfaces, not accepted behavior,
production architecture, security proof, or maintainability proof.

## 8. Minimum rebuilt-product acceptance scenarios

### Scenario A — cold recovery

A person or agent with no chat history locates accepted direction, current
state, relevant evidence, blockers, and either one safe next action or the exact
unresolved decision. Every consequential conclusion drills down to its source.

### Scenario B — fuzzy intent to bounded accountable work

An owner brings an incomplete change idea. The system exposes uncertainty and
boundaries, records required owner acceptance, defines expected evidence, and
produces a bounded handoff with authority, responsibility, dependencies, stop
conditions, and return expectations.

### Scenario C — return, evaluate, reconcile, and resume

An execution actor returns a result with evidence. The system evaluates it
against accepted expectations, records the authorized disposition, reconciles
accepted changes with current project truth, and supports a fresh cold recovery.

These scenarios are provider-neutral. Record outcome quality, missed context or
authority errors, recovery failures, and context/attention cost—not merely
whether completion eventually occurred.

## Evidence still required later

Representative project/task cohorts, scenario protocols, comparative baselines,
external-owner outcomes, and claim-specific thresholds remain unknown until a
later decision or experiment requires them. Their absence must be exposed; it
does not authorize invented precision.

## Exit condition

S03B may now define the full truth hierarchy and contract model using a stable,
non-gameable evaluation method. S10 must use this ruler after the intervening
contracts and compatibility-floor checkpoint are accepted.
