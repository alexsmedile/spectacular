---
type: Mission
id: 01a028e9-c0b8-77dc-9dd3-9ff0b46cedca
title: Build the Spectacular effectiveness benchmark suite
status: active
created: "2026-08-22T12:01:56Z"
updated: "2026-08-22T13:01:13Z"
activation:
    at: "2026-08-22T12:01:56Z"
    by: Alex
    fingerprint: sha256:462c98e6ee922498949030314b08cac6de9f3af233ff28ea1707b2cfa3e66189
authority:
    operator:
        - inspect
        - edit-in-scope
        - choose-reversible-implementation
        - run-checks
        - generate-derived-files
        - bounded-repair
        - commit-local
    requires_owner:
        - activate-mission
        - change-outcome-or-completion
        - expand-scope
        - push
        - merge
        - release
        - irreversible-change
        - destructive-data
        - secret-change
baseline:
    branch: refactor/constitutional-skill-kernel
    commit: 14158f92aaa78e1a90f0231afdc52a04f3cca6c3
completion:
    - claim: measurement-contract
      pass_boundary: A versioned benchmark contract separates hard safety gates from task success, routing and context efficiency, interaction cost, and recovery quality; every metric has an observable source and aggregation rule.
      proof_requirement: Table-driven tests load every benchmark case and assert metric definitions, hard-failure precedence, thresholds, held-out classification, and score aggregation without allowing safety failures to be averaged away.
    - claim: uncontaminated-paired-harness
      pass_boundary: The harness materializes complete baseline and candidate skill packages from immutable Git revisions into isolated fixture workspaces, randomizes pair order, records package and model identity, and prevents either run from inspecting the other.
      proof_requirement: Automated harness tests use fake adapters and temporary repositories to prove revision extraction, isolation, randomized order, artifact separation, resumable manifests, and refusal of dirty or ambiguous comparison inputs.
    - claim: discriminating-case-suite
      pass_boundary: Machine-readable behavior and trigger cases cover role bootstrap, exclusive phase routing, progressive disclosure, bounded Runner context, hostile entry contracts, CLI absence and incompatibility, owner batching, review independence, cold recovery, lightweight intake, and a small end-to-end build; plausible baseline behavior can fail each expectation.
      proof_requirement: Case-validation tests require unique IDs, explicit fixtures, objective assertions, hard-failure markers, expected reads and forbidden reads, and separate held-out variants; all shipped cases validate.
    - claim: attributable-results
      pass_boundary: Static, smoke, and full tiers emit human-readable Markdown plus machine-readable JSON containing raw trial references, pass rates, safety failures, context and interaction costs, variance, limitations, and an explicit old-versus-new verdict.
      proof_requirement: Golden report fixtures and scorer tests assert deterministic JSON and Markdown output, baseline-relative host-check handling, per-case regression visibility, and clear inconclusive results when repetitions or trace data are insufficient.
contract:
    fingerprint: sha256:aa2f59e740e9526bacef1dd9999127861836460e5f2f96b5fe05bc86a458ee1a
    ref: Contract:019fe381-5d61-7223-b362-03a5f99a7b10
dependencies: []
gaps: []
objectives:
    - claims:
        - measurement-contract
        - uncontaminated-paired-harness
        - discriminating-case-suite
        - attributable-results
      id: 01a028e9-c0b8-7666-baa6-b5255fd8159d
      outcome: Design, implement, validate, and document the complete Spectacular effectiveness evaluation harness and its initial case catalog.
      ref: O1
      status: pending
outcome: Spectacular skill revisions can be compared reproducibly against an immutable baseline using safety-first behavioral measures, objective context and interaction costs, and evidence-backed reports.
owner: Alex
ref: M14
repair_budget: 2
review: independent
reviews:
    - file: reviews/RV1-independent-review-of-the-spectacular-effectiveness-benchmark-suite.md
      id: 01a02959-5720-7e5e-a08f-76084b80dbd1
      ref: RV1
      verdict: pass
run:
    current_objective: O1
    id: 01a028e9-c0b8-76f5-8902-55404e0e8121
    operator: Alex
    ref: R1
    repairs: 0
    started_at: "2026-08-22T12:01:56Z"
    status: active
scope:
    mechanical:
        - test/evals/spectacular
        - .spectacular/missions
        - CHANGELOG.md
    semantic:
        - Spectacular skill effectiveness measurement
        - paired evaluation methodology
        - benchmark fixtures
        - test-only harness and reports
start_key: sha256:171a63284ccf73804c75d8fb05bdca351c385f168c78c4855d21a48da10b6607
stops:
    - scope-drift
    - public-cli-change
    - benchmark-contamination
    - unverifiable-metric
validation:
    mode: cli
    schema: mission.v2
---
# Mission

Build a test-only empirical evaluation surface. Keep model-driven trials separate from the Go product verification ladder; deterministic harness tests may run in the repository test environment, while paid or stochastic trials require explicit invocation.

## Measurement policy

- Safety and authority violations are hard failures and never averaged away.
- Compare usefulness and cost separately before any composite summary.
- Record actual trace observations when available; label unsupported metrics rather than estimate them silently.
- Use the complete skill package from immutable revisions, not a lone saved SKILL.md.
- Keep behavior evals separate from skill-trigger discovery evals.
- Freeze visible development cases and held-out variants; never tune on held-out results.

## Initial execution

1. Define the case schema, metric dictionary, tier semantics, and acceptance policy.
2. Implement revision materialization, isolated workspace assembly, adapter execution, trace normalization, scoring, and reporting.
3. Add deterministic fake-adapter tests and representative behavior/trigger fixtures.
4. Validate the harness locally, run a bounded smoke comparison when the candidate is immutable, and record limitations when a real model or trace field is unavailable.

## Checkpoints

- After measurement contract: confirm every metric has a source and a failure interpretation.
- After harness tests: confirm baseline/candidate isolation and scorer correctness.
- Before paid model trials: report planned case count, repetitions, model, and expected external cost.
