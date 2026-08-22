---
type: MissionPlan
title: Measure and harden scope guardrails
owner: Alex
contract:
  ref: Contract:01a00a20-63dd-7670-97f1-9eb8e12adc3a
outcome: The Campaign ends with measured evidence about scope failures, only deterministic guardrails that demonstrably avoid rejecting coherent work, and no general quality classifier.
review: independent
completion:
  - claim: discriminating-guardrail-evaluation
    pass_boundary: Paired fixtures distinguish undeclared writes, reservation overlap, dependency escape, unauthorized dependency introduction, and immutable-context loss from harmless excess over file-count or size guidance, with every candidate rule reporting observed false-positive and false-negative cases.
    proof_requirement: Positive and adversarial fixtures isolate each violation and benign overrun across exact files, trailing-directory scaffolding, renames, generated files, repair context, and coherent larger slices; results are attributable and reproducible.
  - claim: earned-deterministic-hardening
    pass_boundary: A candidate rule is promoted only when its input and refusal are deterministic, it catches a demonstrated authority or safety failure, and it rejects no benign fixture; zero promoted rules is a valid successful result and subjective quality remains Orchestrator judgment.
    proof_requirement: For every promoted rule, table-driven tests assert exact inputs, reason code, refusal-before-effect, safe next action, and zero benign-fixture rejection; the final report lists rejected candidate rules and why they were not mechanized.
  - claim: coherent-broad-perimeters
    pass_boundary: Two-to-four files remains planning guidance rather than validation; a frozen exact-file or trailing-directory perimeter may be broader when disjoint and justified, while only a new Handoff can expand it.
    proof_requirement: Scaffolding, refactor, generated-output, and mixed-scope fixtures cover justified broad directories, missing rationale, collisions, new files, rename endpoints, and agent-attempted expansion.
  - claim: final-campaign-regression
    pass_boundary: The exact 18-command product retains M16's context reduction and every M17-M20 safety, history, atomicity, isolation, and Evidence guarantee; no proof graph, timeline feature, semantic quality classifier, or forced Gap closure is introduced.
    proof_requirement: The pinned M14 paired suite plus full Mission verification reports no regression; independent review traces each shipped P11 boundary to code, tests, Contract, Skill, and generated interface and confirms concurrent-run-timelines remains an honest open Gap.
objectives:
  - outcome: Measure candidate scope guardrails against real violations and coherent larger work.
    claims: [discriminating-guardrail-evaluation]
  - outcome: Promote only earned deterministic checks and preserve justified broad perimeters.
    claims: [earned-deterministic-hardening, coherent-broad-perimeters]
  - outcome: Run the final Campaign regression and traceability review.
    claims: [final-campaign-regression]
authority:
  operator: [inspect, edit-in-scope, choose-reversible-implementation, run-checks, generate-derived-files, bounded-repair, commit-local]
  requires_owner: [activate-mission, change-outcome-or-completion, expand-scope, push, merge, release, irreversible-change]
scope:
  mechanical: [internal/runtime, internal/missionbundle, internal/command, skills/spectacular, test/evals/spectacular, test, .spectacular/contracts]
  semantic: [scope-guardrail measurement, deterministic authority refusals only when earned, coherent broad write perimeters, final P11 regression proof]
repair_budget: 2
dependencies: [M20 completed with the exact 18-command surface and accepted Evidence]
gaps: []
stops: [subjective-quality-classifier, numeric-proxy-as-quality-proof, benign-fixture-rejection, agent-authored-scope-expansion, behavioral-regression, forced-timeline-gap-closure, data-loss]
---

# Mission

> **Future Mission sketch.** Preserve as design input. Do not activate, maintain,
> validate, or review as a final plan until its predecessor closes and the
> Orchestrator re-prepares this block from current Evidence.

Measure first. It is acceptable for this Mission to ship no new guardrail when no
candidate meets the deterministic zero-benign-rejection boundary. The existing
`concurrent-run-timelines` Gap stays open; this Mission does not build a feature to
make the governance ledger look complete.
