---
type: MissionPlan
title: Reconcile P11 governance and Contract baselines
owner: Alex
contract:
  ref: Contract:019fe381-5d61-7223-b362-03a5f99a7b10
outcome: P11 begins from one reviewed governance baseline that preserves inherited proof gates, states the real 14-command interface and explicit source input, sequences later Contract changes without rewriting completed Mission bindings, and makes one-Mission-ahead planning the default.
review: independent
completion:
  - claim: inherited-gates-preserved
    pass_boundary: D21 carries forward the complete Objective-scoped Run, frozen-Handoff, clustered-Evidence, history, and reservation model while explicitly restoring D12's at-least-40-percent context reduction and zero-regression conditions.
    proof_requirement: A lineage review compares D12, D17, D21, P11, and the campaign and finds no dropped authority, proof, concurrency, history, or recovery condition and no conflicting active disposition.
  - claim: exact-interface-baseline
    pass_boundary: The generated command registry and every current Contract describe exactly the same 14-command baseline; stale statements that the surface has ten or twelve commands are removed through explicit Contract version bumps rather than editorial amendments.
    proof_requirement: A registry-derived inventory names all 14 commands and a Contract/index audit finds no contradictory count or command list; pure governance checks and diff checks pass without running the Go verification ladder.
  - claim: explicit-retrieval-input-contract
    pass_boundary: The versioned CLI Contract defines optional ordered typed sources lists on Mission and Objective frontmatter as frozen semantic input; automatic traversal reads only the bound Contract and those direct refs, while invocation-added refs remain temporary and no recursive or meaning-inferred traversal is authorized.
    proof_requirement: Contract and schema review identifies exact field ownership, ordering, duplicate handling, valid ref kinds, fingerprint inclusion, and refusal behavior, and finds no second retrieval vocabulary or implicit source edge.
  - claim: deterministic-budget-and-lifecycle-rules
    pass_boundary: Owner activation of M15 freezes spectacular-charter-tokenizer.v1 with byte-exact UTF-8, bundled o200k_base data and a published SHA-256, plus the complete legal Run transition matrix and terminal states proposed by P11; later Missions make no product choice about either rule.
    proof_requirement: Boundary review proves identical bytes and tokenizer data yield identical counts across supported platforms, invalid UTF-8 refuses, tokenizer version changes are explicit, every state has enumerated outgoing edges, and completed and stopped cannot be reopened.
  - claim: safe-contract-transition-map
    pass_boundary: Each M16-M21 candidate identifies a stable bound Contract, any different Contract it may version, the proof required before that version is accepted, and the later Mission that consumes it; no completed Mission fingerprint is repointed.
    proof_requirement: Independent review walks the full transition chain and proves no Mission versions its own bound behavioral Contract, no behavioral change is routed through contract amend, and the CC-projsurf concurrency Gap remains openly documented rather than forcing speculative implementation.
  - claim: progressive-proportional-planning
    pass_boundary: Spectacular prepares one activation-ready Mission at a time; downstream Campaign blocks stay compact and fluid, retained future drafts are explicitly non-authoritative sketches that are not synchronized or reviewed as final plans, owner questions are batched around open semantic choices, and repeated review requires a material semantic repair or newly discovered conflict.
    proof_requirement: A focused behavioral fixture gives the Orchestrator a seven-block Campaign and asserts one complete next-Mission preview, compact downstream blocks, no downstream activation claim, no repeated settled question, and no second review request without a material plan change; preparation guidance states the same boundary directly.
objectives:
  - outcome: Reconcile Decision lineage, the generated command baseline, explicit retrieval input, token counting, and Run transitions.
    claims: [inherited-gates-preserved, exact-interface-baseline, explicit-retrieval-input-contract, deterministic-budget-and-lifecycle-rules]
  - outcome: Freeze the Contract and Mission transition map for the remaining campaign.
    claims: [safe-contract-transition-map]
  - outcome: Make preparation progressive so only the next Mission earns complete planning and review.
    claims: [progressive-proportional-planning]
authority:
  operator: [inspect, edit-in-scope, generate-derived-files, run-checks, commit-local]
  requires_owner: [activate-mission, change-outcome-or-completion, expand-scope, push, merge, release, irreversible-change]
scope:
  mechanical: [.spectacular/contracts, .spectacular/decisions, .spectacular/proposals, .spectacular/campaigns, skills/spectacular/references/prepare.md, test/evals/spectacular]
  semantic: [P11 Decision lineage, current public command inventory, explicit Mission and Objective source input, deterministic charter token counting, legal Run transitions, Contract transition ordering, progressive proportional Campaign planning]
repair_budget: 2
dependencies: [M14 completed with independent review]
gaps: []
stops: [product-code-change, public-command-change, completed-mission-repoint, behavioral-contract-amendment, unresolved-decision-conflict, downstream-plan-freeze, repeated-review-without-material-change]
---

# Mission

This is a governance-only reconciliation slice. It changes no Go code and adds no
command. The independently reviewed output is the baseline the owner accepts before
any context compiler or lifecycle implementation begins.

M16-M21 remain preserved as design sketches. M15 does not maintain or validate
them. After M15 closes, the Orchestrator re-prepares M16 alone from current Evidence.

Contract sequence: M16 binds reconciled CC-missioncli; M17 binds CC-v2prod while
versioning CC-missioncli; M18 binds CC-v2prod while versioning CC-missioncli and
CC-projsurf; M19 binds the new CC-projsurf while versioning CC-v2prod; M20 binds
that CC-v2prod while versioning only CC-missioncli; M21 binds the new CC-missioncli
and versions CC-projsurf only if a deterministic guardrail earns promotion.
