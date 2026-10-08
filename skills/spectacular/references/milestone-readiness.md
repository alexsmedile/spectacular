# Milestones and readiness

Use this when: defining a substantial delivery increment, launch stage, release
readiness, or an MVP validation experiment. These are reasoning prompts for the
existing plan/requirements owner, not a schema, mandatory checklist, or new gate.

## Define the outcome before the tasks

State who can complete which workflow, in which environment, with what observable
result. Prefer a thin end-to-end slice over separate database/backend/UI phases.
Prerequisite work may be horizontal when necessary; link it to the workflow it
unblocks. A spike closes a question with a decision, not a usable-product claim.
A walking skeleton proves an integrated path, not adoption or market demand.

For each milestone, capture the relevant facts in the existing plan:
- Intended outcome, audience, environment, and scope/non-goals.
- Blocking prerequisites and why each is necessary for this outcome or exposure.
- Applicable gates, observable pass/fail boundaries, evidence, and responsible owner.
- Known limitations, permitted workarounds, unresolved blockers, and next decision.

Use the project's stage names. Define alpha/beta/pilot by audience, exposure,
workflow, support expectations, and exit evidence; names alone prove nothing.
These are independent readiness dimensions, not an automatic stage ladder:

| Dimension | Question | Representative evidence |
|---|---|---|
| Usability | Can the intended user complete the core workflow? | Observed end-to-end use in the declared environment |
| Integrity | Do failures preserve the required invariants? | Relevant error, concurrency, and recovery checks |
| Distribution | Can the recipient install or access this version? | Identified artifact, fresh install/deploy, smoke test |
| Operation | Can someone detect and respond to a material failure? | Appropriate signals, responsible operator, recovery path |
| MVP validation | Did the experiment answer the value hypothesis? | Representative users, measurement, threshold, decision |
| Publication | Can this version reach this audience? | Applicable checks and existing or newly required authority |

## Select indispensable gates from actual constraints

Distinguish gates needed before the first workflow, before real data/users or
irreversible effects, and before wider distribution. Accepted invariants,
repository policy, applicable obligations, and effect permissions remain binding.
Do not infer auth, payments, databases, rate limits, cloud hosting, or automated
deployment from the word MVP. Explain each prerequisite through the failure it
prevents or outcome it enables. Basic repeatable verification should precede
repeated/shared delivery; a disposable spike may use a bounded local check.

A gate needs a reason, scope, pass boundary, evidence source, and owner. Report
its result as passed, failed, unverified, or justified not-applicable in ordinary
prose; these are not governed status fields. Unknown evidence is not a pass.
Carry prior authorizations forward. Passing tests or preparing a release does not
by itself authorize publishing, live-data access, or Mission activation.

## Prioritize the next useful increment

Use only the lenses that resolve the current choice:
- **Vertical slices / walking skeleton:** integrate the smallest useful path early.
- **Risk and learning:** explore the uncertainty that could invalidate that path;
  bound the experiment by a question, effort limit, and decision.
- **MoSCoW:** a Must is necessary for this increment to remain viable; record the
  consequence of omission. Should/Could can wait or use an explicit workaround;
  Won't protects scope for this increment, not forever.
- **Reversibility:** assess actual switching cost, data exposure, dependencies,
  and migration consequences; do not classify every stack choice as irreversible.
- **KISS/YAGNI:** prefer the simplest sufficient design; add speculative capacity
  only for a demonstrated requirement. Safety and integrity are quality needs,
  not optional polish; performance targets follow the actual use and risk.

Treat scoring as a decision aid, not manufactured numerical certainty. Resolve
hard dependencies first; prioritize coherent value/learning within those limits.
No fixed M0–M3 sequence, calendar duration, or infrastructure bundle is required.

## Match proof to the claim

Identify the version/artifact, environment, data/dependency conditions, observation
source and time, and remaining limitations where material. Distinguish simulated,
integrated, real-dependency, and user-observed evidence. A mock cannot prove live
integration; an integration test cannot establish user value. Recheck affected
claims after material changes rather than reusing stale evidence.

For an MVP experiment, state the target user/problem, value hypothesis, smallest
representative intervention, observation window, measure and decision threshold,
and outcome: supported, contradicted, or inconclusive. Thresholds must be agreed
or clearly proposed; do not invent stakeholder acceptance or statistical certainty.
Manual operations are valid when declared and representative of the hypothesis.
An inconclusive result names missing evidence and a bounded next action.

## Reuse owners and route narrowly

The plan owns delivery detail; requirements own needs and acceptance; an existing
roadmap summarizes milestone state. Link evidence and reconcile affected claims
without creating another status catalog. Inspect factual gates independently of
whether the project chooses governed execution. For an explicitly selected
Mission, preserve frozen criteria and retrieve any mechanical input via --schema;
this reference does not change record grammars or certify business acceptance.

Use system-architecture for consequential prerequisites, test-sentinel for
executable proof, rapid-prototyping for bounded uncertainty, and data/interface/
event companions only for effects present in the increment. Do not load the whole
family or create a product-validation skill automatically.
