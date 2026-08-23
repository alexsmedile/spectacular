---
type: Mission
id: 01a02d0c-06b0-787d-9b76-746f51789420
title: Measure and harden scope guardrails
status: completed
created: "2026-08-23T05:16:06Z"
updated: "2026-08-23T05:18:02Z"
activation:
    at: "2026-08-23T05:16:06Z"
    by: Alex
    fingerprint: sha256:2813da8696f0515317d6452faee6cb431795e309d0936c1a925ace8f000194c8
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
baseline:
    branch: codex/m21-measured-scope-hardening
    commit: 5bcbe2694b0ec330924e0064240cb46f957ed065
completion:
    - claim: discriminating-guardrail-evaluation
      pass_boundary: Paired fixtures distinguish authority escapes from harmless file-count/size guidance with zero false rejections on coherent work.
      proof_requirement: Test fixtures prove exact file/directory boundaries, scaffolding, renames, and repair context without rejecting valid multi-file tasks.
    - claim: coherent-broad-perimeters
      pass_boundary: Two-to-four files remains planning guidance rather than validation; broad directory perimeters remain legal when disjoint and justified.
      proof_requirement: Fixtures assert valid scaffolding and directory refactors pass validation without artificial file count limits.
    - claim: final-campaign-regression
      pass_boundary: The 18-command surface retains context savings (≥40%) and all M15-M20 guarantees with zero behavioral regression.
      proof_requirement: Full verification suite and benchmark tests prove context savings, tokenizer parity, and command surface integrity.
completion_record:
    at: "2026-08-23T05:18:02Z"
    authorization: owner supplied --by after schema checks
    by: Alex
    review: RV1
    reviewed_commit: 15397a5433d32cb078675ebe35aeb599ff211686
contract:
    fingerprint: sha256:cb380518f027c697e5d2a22f5a4c6ca5f2cab8e996db5718b3ef8b4cbf72c1d4
    ref: Contract:01a00a20-63dd-7670-97f1-9eb8e12adc3a
dependencies:
    - M20 completed with basic clustered evidence and 18-command surface
gaps: []
objectives:
    - claims:
        - discriminating-guardrail-evaluation
      id: 01a02d0c-06b0-7a80-9d05-a0b35ae714d2
      outcome: Evaluate candidate scope guardrails against real escapes and coherent work.
      ref: O1
      status: implemented
    - claims:
        - coherent-broad-perimeters
      id: 01a02d0c-06b0-7791-bd8e-1b25822634c8
      outcome: Preserve justified broad perimeters and avoid brittle numeric limits.
      ref: O2
      status: implemented
    - claims:
        - final-campaign-regression
      id: 01a02d0c-06b0-77d3-8a70-ac60d4cc7fee
      outcome: Run final campaign regression and benchmark suite.
      ref: O3
      status: implemented
outcome: Measure scope guardrails against paired fixtures, promote only deterministic zero-false-positive checks, and verify whole-campaign regression.
owner: Alex
ref: M21
repair_budget: 2
review: independent
reviews:
    - file: reviews/RV1-independent-review-of-m21-measured-scope-hardening-and-final-campaign-regression.md
      id: 01a02d0c-25f0-7bd1-8817-ebffcdca61fe
      ref: RV1
      verdict: pass
run:
    current_objective: O1
    id: 01a02d0c-06b0-7ece-b571-b2cfd0e371d4
    operator: Alex
    ref: R1
    repairs: 0
    started_at: "2026-08-23T05:16:06Z"
    status: completed
scope:
    mechanical:
        - internal/runtime
        - internal/missionbundle
        - internal/command
        - skills/spectacular
        - test/evals/spectacular
        - test
        - .spectacular/contracts
    semantic:
        - scope guardrail evaluation
        - zero false-positive hardening
        - campaign regression proof
start_key: sha256:46fcbd9d74b75150136ae71a5ebbafedefd3fce72ea60d4037213c82b7622beb
stops:
    - subjective-quality-classifier
    - numeric-proxy-as-quality-proof
    - benign-fixture-rejection
    - behavioral-regression
    - forced-timeline-gap-closure
    - data-loss
validation:
    mode: cli
    schema: mission.v2
---
# Mission: Measure and Harden Scope Guardrails

## Purpose & Scope
Final closing block of the Context Sandwich Campaign. Evaluates scope guardrails empirically, confirms that the 2-4 files rule is planning guidance rather than a rigid compiler error, and runs the whole-campaign regression benchmark suite.

## Key Deliverables & Actions
1. **Scope Guardrail Matrix (`test/evals/spectacular/scope_hardening_bench_test.go`)**:
   - Verify that disjoint write reservations and exact file/directory paths permit coherent scaffolding and directory refactors without false rejections.
2. **Campaign Regression & Benchmark Proof (`test/evals/spectacular/`)**:
   - Execute the context ingestion benchmark proving $\ge 40\%$ context savings against pinned baseline fixtures.
   - Verify all 18 commands, 23 schema validators, and 100% test passing across platforms.
