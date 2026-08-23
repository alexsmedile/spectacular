---
type: MissionPlan
title: Measure and harden scope guardrails
owner: Alex
contract:
  ref: Contract:01a00a20-63dd-7670-97f1-9eb8e12adc3a
outcome: Measure scope guardrails against paired fixtures, promote only deterministic zero-false-positive checks, and verify whole-campaign regression.
review: independent
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
objectives:
  - outcome: Evaluate candidate scope guardrails against real escapes and coherent work.
    claims: [discriminating-guardrail-evaluation]
  - outcome: Preserve justified broad perimeters and avoid brittle numeric limits.
    claims: [coherent-broad-perimeters]
  - outcome: Run final campaign regression and benchmark suite.
    claims: [final-campaign-regression]
authority:
  operator: [inspect, edit-in-scope, choose-reversible-implementation, run-checks, generate-derived-files, bounded-repair, commit-local]
  requires_owner: [activate-mission, change-outcome-or-completion, expand-scope, push, merge, release, irreversible-change, destructive-data]
scope:
  mechanical: [internal/runtime, internal/missionbundle, internal/command, skills/spectacular, test/evals/spectacular, test, .spectacular/contracts]
  semantic: [scope guardrail evaluation, zero false-positive hardening, campaign regression proof]
repair_budget: 2
dependencies: [M20 completed with basic clustered evidence and 18-command surface]
gaps: []
stops: [subjective-quality-classifier, numeric-proxy-as-quality-proof, benign-fixture-rejection, behavioral-regression, forced-timeline-gap-closure, data-loss]
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
