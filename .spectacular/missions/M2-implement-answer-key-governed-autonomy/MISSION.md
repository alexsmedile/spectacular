---
type: Mission
id: 01a007c2-3d1e-7664-bdb3-42cebb6b1759
title: Implement answer-key governed autonomy
status: completed
created_by: Codex primary implementation session
created: "2026-08-15T23:33:04Z"
updated: "2026-08-16T13:05:00Z"
closure:
    kind: historical
    at: "2026-08-16T13:05:00Z"
    by: Alex
    review: none
    implementation_commit: 147776293fbd06d9a12093ec103bc940eeb87f7f
    merged_commit: c63d1a229ada24ea8914fa0440d185e9aa30cec2
    basis: Owner closed this Mission retroactively. The implementation shipped and is merged to main; the governance record was never advanced during execution.
    not_claimed:
        - No independent review was performed against the reviewed tree.
        - No Evidence records were captured for any claim.
        - Objective and Run records were never advanced past their activation state.
        - The four frozen completion criteria were not individually verified at closure.
source: Proposal:01a007c2-3d1e-7655-a8ae-99d9cc23a35f
activation_decision: Decision:01a007c2-3d1e-765c-b7aa-deb675155b93
allowed_actions:
    - inspect
    - edit
    - test
    - generate-interface
    - assemble-release
    - record-evidence
baseline: d9817144bcfd5248207bcd5e679e6172f518ecf3
budget_units: 24
completion_contract:
    - claim: claim:completion-contract-readiness
      pass_boundary: Preparation freezes one complete criterion per claim, emits only actual readiness failures, and cold Mission recovery succeeds.
      proof_requirement: Focused compiler, persistence, and cold-resume tests plus full verification.
      review_level: independent
    - claim: claim:dependency-aware-delegation
      pass_boundary: Objective graphs are acyclic and live Handoffs carry exact dependencies, disjoint Objective claims, bounded inputs, and explicit returns.
      proof_requirement: Governance tests for DAG and Handoff acceptance and refusal cases.
      review_level: clustered
    - claim: claim:proportional-review
      pass_boundary: Closure accepts automatic, clustered, and independent evidence only at or above each frozen review level.
      proof_requirement: Integration tests for all accepted strengths and downgrade refusals.
      review_level: independent
    - claim: claim:truthful-autopilot-limits
      pass_boundary: Hard caps require measurement and cancellation; observed and unsupported caps cannot claim stronger enforcement.
      proof_requirement: Runtime compiler tests for each enforcement level and false-hard refusal.
      review_level: clustered
completion_contract_fingerprint: 0b0eec80c70dd6f2c9abdca3f1bdb6c77879757f390dcff27d558bffaf29dc17
current_run: Run:01a007c2-3d1e-766c-9f40-633a95f5c731
expected_run_fingerprint: dbae7ab5ca14ef3fe60f45fe9b0790c1e426403c1573d0ca35ccdc3b816fba4f
dependencies:
    - v2.0.0-rc.2 tagged baseline
    - current v2 product Contract
design_sufficiency: sufficient
evidence_claims:
    - claim:completion-contract-readiness
    - claim:proportional-review
    - claim:dependency-aware-delegation
    - claim:truthful-autopilot-limits
expires_at: "2026-08-18T23:59:59Z"
forbidden_effects:
    - destructive-data
    - merge
    - production-configuration
    - production-release
    - remote-deletion
    - secret-change
    - security-privacy-rights-sensitive
freshness_checked_at: "2026-08-15T23:33:04Z"
freshness_source: .spectacular/workspace.yaml
freshness_source_fingerprint: d8b24fe7cfef0986a4b48e7f4e6dd8c7373b451d4c54bde425a904889539b4d3
freshness_valid_until: "2026-08-16T23:33:04Z"
gaps: []
human_ref: M2
idempotency_key: v21-governed-autonomy-mission
last_authorization: Decision:01a007c2-3d1e-7674-aa7d-95d2b17334dc
last_idempotency_key: v21-governed-autonomy-start
last_transition_input_fingerprint: 5ec3d585fa0c383c5dd101f5915e8afc302dc65d6c613976a53ae2db0607c0e9
objectives:
    - Objective:01a007c2-3d1e-7668-97aa-7f1e0f8274ea
    - Objective:01a007c4-8b0e-7416-9e26-351a0c750e91
    - Objective:01a007c4-8b0e-7455-af99-68a7dd6526b9
    - Objective:01a007c4-8b0e-7459-baab-8935d641d4e1
outcome: Spectacular prepares answer-key Missions, delegates cohesive claim scopes, applies proportional review, and emits truthful Autopilot enforcement limits.
preparation_baseline: d9817144bcfd5248207bcd5e679e6172f518ecf3
preparation_fingerprint: 648efe3b294207604eaf363ce9a551d059b2b05ab8c14368db61f1b190258f0f
preparation_receipt: '{"schema_version":"spectacular.mission-preparation.v2","fingerprint":"648efe3b294207604eaf363ce9a551d059b2b05ab8c14368db61f1b190258f0f","proposal":{"ref":"Proposal:01a007c2-3d1e-7655-a8ae-99d9cc23a35f","fingerprint":"fe2c162257b01de191f729f6f68acf1ed4d8c70d468af9e2a279069c85bb0f31"},"baseline":"d9817144bcfd5248207bcd5e679e6172f518ecf3","direction_sources":[{"ref":"Anchor:019fe381-5d61-7223-b362-03a5f99a7b14","fingerprint":"c45130e86ad48a4edfaebdda67ba4f79e3d51ec0f7c157b2520914ea3cfad93d"},{"ref":"Anchor:019fe381-5d61-7223-b362-03a5f99a7b15","fingerprint":"cdc1d055c63bd234b3c454c364b3ef6b8dc16e50cc77fd5daf0b9621b6f52456"},{"ref":"Contract:019fe381-5d61-7223-b362-03a5f99a7b10","fingerprint":"a64746f0376df57b13cdea0020d974c93a5b1be35fdbac0b5fa4922fbd983625"}],"candidates":[{"name":"governed-autonomy-v2.1","outcome":"Spectacular prepares answer-key Missions, delegates cohesive claim scopes, applies proportional review, and emits truthful Autopilot enforcement limits.","evidence":["Compiler and governance tests cover every frozen criterion.","Real-process acceptance proves recoverability.","The complete release gate passes."],"dependencies":["RC.2 published baseline","current v2 Contract","host capability declarations"],"cancellation_state":"Each completed cluster remains tested and coherent.","reversibility":"Changes remain local to the v2 feature branch until owner disposition.","standalone_coherence":"The clusters form one preparation-to-execution control path.","integration_path":"Preparation, review, Handoffs, Autopilot, Skill guidance, and acceptance coverage.","learning_value":"Tests whether explicit answer keys reduce assumptions and review loops."}],"selected":"governed-autonomy-v2.1","design_sufficiency":"sufficient","design_rationale":"The owner completed an adaptive grill covering authority, criteria, review, delegation, limits, and the stopping predicate.","slice_quality":"coherent","slice_rationale":"Four ordered clusters share one authority model and remain independently testable.","blocking_gaps":[],"completion_criteria":[{"claim":"claim:completion-contract-readiness","pass_boundary":"Preparation freezes one complete criterion per claim, emits only actual readiness failures, and cold Mission recovery succeeds.","proof_requirement":"Focused compiler, persistence, and cold-resume tests plus full verification.","review_level":"independent"},{"claim":"claim:dependency-aware-delegation","pass_boundary":"Objective graphs are acyclic and live Handoffs carry exact dependencies, disjoint Objective claims, bounded inputs, and explicit returns.","proof_requirement":"Governance tests for DAG and Handoff acceptance and refusal cases.","review_level":"clustered"},{"claim":"claim:proportional-review","pass_boundary":"Closure accepts automatic, clustered, and independent evidence only at or above each frozen review level.","proof_requirement":"Integration tests for all accepted strengths and downgrade refusals.","review_level":"independent"},{"claim":"claim:truthful-autopilot-limits","pass_boundary":"Hard caps require measurement and cancellation; observed and unsupported caps cannot claim stronger enforcement.","proof_requirement":"Runtime compiler tests for each enforcement level and false-hard refusal.","review_level":"clustered"}],"stop_conditions":["A change would weaken owner control over semantic success criteria.","A hard limit is requested without measurable cancellation support.","The implementation requires a new record/search compatibility surface or v1 behavior.","The full verification gate cannot pass within the declared repair budget."],"evidence_claims":["claim:completion-contract-readiness","claim:proportional-review","claim:dependency-aware-delegation","claim:truthful-autopilot-limits"],"fresh_until":"2026-08-18T23:59:59Z","ready":true,"unmet_requirements":[],"next":"owner-activation"}'
preparation_sources:
    - Proposal:01a007c2-3d1e-7655-a8ae-99d9cc23a35f@fe2c162257b01de191f729f6f68acf1ed4d8c70d468af9e2a279069c85bb0f31
    - Anchor:019fe381-5d61-7223-b362-03a5f99a7b14@c45130e86ad48a4edfaebdda67ba4f79e3d51ec0f7c157b2520914ea3cfad93d
    - Anchor:019fe381-5d61-7223-b362-03a5f99a7b15@cdc1d055c63bd234b3c454c364b3ef6b8dc16e50cc77fd5daf0b9621b6f52456
    - Contract:019fe381-5d61-7223-b362-03a5f99a7b10@a64746f0376df57b13cdea0020d974c93a5b1be35fdbac0b5fa4922fbd983625
preparation_valid_until: "2026-08-18T23:59:59Z"
recovery_point: git branch codex/governed-autonomy at d9817144bcfd5248207bcd5e679e6172f518ecf3
repair_budget: 4
return_destination: owner in this Codex task
scope:
    - v2
slice_quality: coherent
stops:
    - A change would weaken owner control over semantic success criteria.
    - A hard limit is requested without measurable cancellation support.
    - The implementation requires a new record/search compatibility surface or v1 behavior.
    - The full verification gate cannot pass within the declared repair budget.
---
# Mission

## Historical closure

This Mission was closed retroactively on 2026-08-16, after its implementation
had already shipped and merged to `main`. It is recorded as `historical` rather
than `completed` through the normal path because the governance record was never
advanced while the work was being done.

**What shipped.** Commit `1477762` (`feat: add governed autonomy contracts`)
implements the answer-key model: frozen completion criteria per claim, acyclic
Objective dependency graphs with explicit Handoffs, proportional review levels
(automatic, clustered, independent), and truthful Autopilot enforcement
disclosure. The code is on `main` at merge `c63d1a2`, its tests pass, and M3,
M4, M5, and M6 were all built on top of it.

**What was never done.** No Evidence records were captured, no independent
review ran against the reviewed tree, and the Objective and Run records still
carry their activation state. The four frozen criteria in `completion_contract`
above were not individually verified at closure. This record therefore proves
that the work shipped; it does not prove that each frozen criterion was met.

The strongest available evidence that the implementation is sound is indirect:
three later Missions depend on this code and exercise it, and the full
repository verification gate passes on the merged tree.

**Schema note.** This Mission validates at `legacy-v2` (3 checks), not
`mission.v2` (14 checks). The M6 CLI reads it correctly but cannot complete it —
`mission complete` refuses because the legacy schema has no `owner` field, no
frozen activation fingerprint, and no Contract binding. The closure above was
therefore written by hand under owner authority rather than generated by the
CLI.
