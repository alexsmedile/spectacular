---
type: Mission
id: 01a0081c-978f-724d-a77b-cf06e3ce2f44
title: Make Spectacular control operations lean
status: active
created_by: Codex primary implementation session
created: "2026-08-16T01:10:12Z"
updated: "2026-08-16T01:10:29Z"
source: Proposal:01a0081c-978f-723e-81dd-bba16b0e82cc
activation_decision: Decision:01a0081c-978f-7246-af3d-dacce64c0958
allowed_actions:
    - inspect
    - edit
    - test
    - generate-interface
    - commit-local
    - record-evidence
    - update-mission-progress
baseline: 87255b9eb297d720a1ec30f8cd60ee55eddf385a
budget_units: 18
completion_contract:
    - claim: claim:batch-autoid
      passboundary: Creation commands accept omitted IDs, stdin input, and atomic batches while preserving UUIDv7, authorization, idempotency, and rollback invariants.
      proofrequirement: Command and governance tests cover generated IDs, stdin, batch success, replay, and rollback/refusal.
      reviewlevel: clustered
    - claim: claim:compact-policy
      passboundary: Routine workflows use compact context first and the core Skill routes rather than carrying detailed execution, review, and audit policy.
      proofrequirement: Canonical Skill tests verify compact-first routing and policy ownership in workflow references.
      reviewlevel: automatic
    - claim: claim:minimal-verification
      passboundary: Guidance and verification support focused changed-scope checks followed by one compact tree-bound full gate, with detailed logs only on failure.
      proofrequirement: Verification-script tests or acceptance checks prove compact success and actionable failure behavior.
      reviewlevel: automatic
    - claim: claim:mission-control
      passboundary: Decision conditions reject invalid vocabulary at creation and one compact Mission progress command advances Objective or emits an implementation-complete owner gate without hand-editing records.
      proofrequirement: Governance and real-process tests cover early refusal, atomic progress, current Objective, and owner-gate continuity.
      reviewlevel: clustered
completion_contract_fingerprint: 0d964e569c529d285be2c481fdee05e4c48e0595f2eec2bf46f0802c1491a87c
current_run: Run:01a0081c-978f-7265-a2ca-edf355c534e8
dependencies:
    - M3 implementation commit 87255b9
    - v2.0.0-rc.2 published baseline
design_sufficiency: sufficient
evidence_claims:
    - claim:compact-policy
    - claim:minimal-verification
    - claim:batch-autoid
    - claim:mission-control
expected_run_fingerprint: ee648ac93e08f29715c2e2862780d2019fc6823509306f0b636792f4f098f217
expires_at: "2026-08-20T23:59:59Z"
forbidden_effects:
    - destructive-data
    - merge
    - production-configuration
    - production-release
    - remote-deletion
    - secret-change
    - security-privacy-rights-sensitive
freshness_checked_at: "2026-08-16T01:10:12Z"
freshness_source: .spectacular/workspace.yaml
freshness_source_fingerprint: d8b24fe7cfef0986a4b48e7f4e6dd8c7373b451d4c54bde425a904889539b4d3
freshness_valid_until: "2026-08-17T01:10:12Z"
gaps: []
human_ref: M4
idempotency_key: m4-lean-control-mission
last_authorization: Decision:01a0081c-978f-7269-8c70-f892b5d2ca27
last_idempotency_key: m4-lean-control-start
last_transition_input_fingerprint: 031eca34a70f320ea67ab60da424a1a70deb863247f2f6438e795fe928e8a0c5
objectives:
    - Objective:01a0081c-978f-7255-9867-0c1935b12355
    - Objective:01a0081c-978f-7259-8bac-5378dc57a9b4
    - Objective:01a0081c-978f-725d-9f2f-b08019b83b2c
outcome: Spectacular carries less default policy and context, verifies once proportionally, and performs common governed mutations in compact atomic commands.
preparation_baseline: 87255b9eb297d720a1ec30f8cd60ee55eddf385a
preparation_fingerprint: cd844ac42550db7e503aed82edf0d8b5c1ae5fbb239114fe890af7a90e0f1374
preparation_receipt: '{"schema_version":"spectacular.mission-preparation.v2","fingerprint":"cd844ac42550db7e503aed82edf0d8b5c1ae5fbb239114fe890af7a90e0f1374","proposal":{"ref":"Proposal:01a0081c-978f-723e-81dd-bba16b0e82cc","fingerprint":"26e39794c70d453624aee0984f6045d3a9efd181055f725b913d275d7eb662f8"},"baseline":"87255b9eb297d720a1ec30f8cd60ee55eddf385a","direction_sources":[{"ref":"Anchor:019fe381-5d61-7223-b362-03a5f99a7b14","fingerprint":"c45130e86ad48a4edfaebdda67ba4f79e3d51ec0f7c157b2520914ea3cfad93d"},{"ref":"Anchor:019fe381-5d61-7223-b362-03a5f99a7b15","fingerprint":"cdc1d055c63bd234b3c454c364b3ef6b8dc16e50cc77fd5daf0b9621b6f52456"},{"ref":"Contract:019fe381-5d61-7223-b362-03a5f99a7b10","fingerprint":"a64746f0376df57b13cdea0020d974c93a5b1be35fdbac0b5fa4922fbd983625"}],"candidates":[{"name":"lean-control-plane","outcome":"Spectacular carries less default policy and context, verifies once proportionally, and performs common governed mutations in compact atomic commands.","evidence":["focused command and governance tests","one full compact verification receipt","canonical Skill conformance"],"dependencies":["M3 implementation tree 87255b9","current v2 Contract"],"cancellation_state":"Each completed cluster remains independently usable.","reversibility":"All product changes remain on the existing unpublished feature branch.","standalone_coherence":"Context, mutation, and progress improvements form one lean control path.","integration_path":"Policy first, mutation ergonomics second, Mission control third, one final gate.","learning_value":"Measures whether governance can remain exact with fewer files, reads, and tool calls."}],"selected":"lean-control-plane","design_sufficiency":"sufficient","design_rationale":"The owner explicitly approved every behavior and requested implementation without another question loop.","slice_quality":"coherent","slice_rationale":"Three dependency-based clusters cover one compact control-plane outcome and can use focused tests.","blocking_gaps":[],"completion_criteria":[{"claim":"claim:batch-autoid","pass_boundary":"Creation commands accept omitted IDs, stdin input, and atomic batches while preserving UUIDv7, authorization, idempotency, and rollback invariants.","proof_requirement":"Command and governance tests cover generated IDs, stdin, batch success, replay, and rollback/refusal.","review_level":"clustered"},{"claim":"claim:compact-policy","pass_boundary":"Routine workflows use compact context first and the core Skill routes rather than carrying detailed execution, review, and audit policy.","proof_requirement":"Canonical Skill tests verify compact-first routing and policy ownership in workflow references.","review_level":"automatic"},{"claim":"claim:minimal-verification","pass_boundary":"Guidance and verification support focused changed-scope checks followed by one compact tree-bound full gate, with detailed logs only on failure.","proof_requirement":"Verification-script tests or acceptance checks prove compact success and actionable failure behavior.","review_level":"automatic"},{"claim":"claim:mission-control","pass_boundary":"Decision conditions reject invalid vocabulary at creation and one compact Mission progress command advances Objective or emits an implementation-complete owner gate without hand-editing records.","proof_requirement":"Governance and real-process tests cover early refusal, atomic progress, current Objective, and owner-gate continuity.","review_level":"clustered"}],"stop_conditions":["A shortcut weakens exact authority, idempotency, atomicity, or evidence attribution.","Batch behavior introduces generic record verbs or ambiguous mutation targets.","The change silently migrates M3 or any active Mission to the new schema.","Compact verification hides a failing check or loses its detailed log."],"evidence_claims":["claim:compact-policy","claim:minimal-verification","claim:batch-autoid","claim:mission-control"],"fresh_until":"2026-08-20T23:59:59Z","ready":true,"unmet_requirements":[],"next":"owner-activation"}'
preparation_sources:
    - Proposal:01a0081c-978f-723e-81dd-bba16b0e82cc@26e39794c70d453624aee0984f6045d3a9efd181055f725b913d275d7eb662f8
    - Anchor:019fe381-5d61-7223-b362-03a5f99a7b14@c45130e86ad48a4edfaebdda67ba4f79e3d51ec0f7c157b2520914ea3cfad93d
    - Anchor:019fe381-5d61-7223-b362-03a5f99a7b15@cdc1d055c63bd234b3c454c364b3ef6b8dc16e50cc77fd5daf0b9621b6f52456
    - Contract:019fe381-5d61-7223-b362-03a5f99a7b10@a64746f0376df57b13cdea0020d974c93a5b1be35fdbac0b5fa4922fbd983625
preparation_valid_until: "2026-08-20T23:59:59Z"
recovery_point: git branch codex/lean-launch-context at 87255b9
repair_budget: 3
return_destination: owner in this Codex task
scope:
    - v2
slice_quality: coherent
stops:
    - A shortcut weakens exact authority, idempotency, atomicity, or evidence attribution.
    - Batch behavior introduces generic record verbs or ambiguous mutation targets.
    - The change silently migrates M3 or any active Mission to the new schema.
    - Compact verification hides a failing check or loses its detailed log.
---
# Mission
