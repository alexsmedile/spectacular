---
type: Mission
id: 01a007f6-d88c-765a-a69e-e1100c71dabb
title: Make Spectacular launch and execution lean
status: active
created_by: Codex primary implementation session
created: "2026-08-16T00:29:36Z"
updated: "2026-08-16T00:29:55Z"
source: Proposal:01a007f6-d88a-7922-a16e-abd2262feda4
activation_decision: Decision:01a007f6-d893-7dec-8fc8-98639fac28cd
allowed_actions:
    - inspect
    - edit
    - test
    - generate-interface
    - commit-local
    - record-evidence
baseline: 147776293fbd06d9a12093ec103bc940eeb87f7f
budget_units: 18
completion_contract:
    - claim: claim:cli-recovery
      passboundary: Common CLI mistakes accept safe intuitive aliases or return a runnable corrected command, and interface generation defaults to VERSION and canonical paths.
      proofrequirement: Deterministic command and generator tests cover project-dot validation, recovery text, and zero-argument canonical generation.
      reviewlevel: automatic
    - claim: claim:launch-preflight-questions
      passboundary: The Skill checks .spectacular/PROJECT.md before root PROJECT.md, performs the approved read-only launch preflight, proposes one clean next action, and formats owner questions with plain and technical layers plus concrete consequences.
      proofrequirement: Skill conformance tests and a cold launch scenario verify the checklist, discovery order, and question grammar.
      reviewlevel: clustered
    - claim: claim:lean-execution-review
      passboundary: The Skill defines compact cluster grammar, progressive context packets, focused-test then integration/full-gate sequencing, local green commits with one push boundary, the self-host schema rule, and FROST critic invariants.
      proofrequirement: Independent Skill review against the exact tree plus distribution and full repository verification.
      reviewlevel: independent
    - claim: claim:progressive-context
      passboundary: Routine launch and Mission context can render a compact card while retaining exact authoritative pointers for progressive source drill-down.
      proofrequirement: Context and command tests prove compact output, stable source identity, and unchanged full JSON authority.
      reviewlevel: clustered
completion_contract_fingerprint: 09d44b5046feaf15a77d026e8d862b430c78dde041677e7ae03f5683355dd562
current_run: Run:01a007f6-d891-7c3e-adbf-bdd759288182
dependencies:
    - M2 implementation commit 147776293fbd06d9a12093ec103bc940eeb87f7f
    - v2.0.0-rc.2 published baseline
design_sufficiency: sufficient
evidence_claims:
    - claim:launch-preflight-questions
    - claim:progressive-context
    - claim:cli-recovery
    - claim:lean-execution-review
expected_run_fingerprint: c8704ded608d19560b88b9ec640708e1b8c5d60ea61febe4a3056af2c7c5a340
expires_at: "2026-08-19T23:59:59Z"
forbidden_effects:
    - destructive-data
    - merge
    - production-configuration
    - production-release
    - remote-deletion
    - secret-change
    - security-privacy-rights-sensitive
freshness_checked_at: "2026-08-16T00:29:36Z"
freshness_source: .spectacular/workspace.yaml
freshness_source_fingerprint: d8b24fe7cfef0986a4b48e7f4e6dd8c7373b451d4c54bde425a904889539b4d3
freshness_valid_until: "2026-08-17T00:29:36Z"
gaps: []
human_ref: M3
idempotency_key: m3-efficiency-mission
last_authorization: Decision:01a007f6-d892-791f-a2df-77818918985e
last_idempotency_key: m3-efficiency-start
last_transition_input_fingerprint: 7dabd493c09c91904e4e752761a6ef9f558ea73c3dd054fb59d892e7746a90e9
objectives:
    - Objective:01a007f6-d88d-710f-a11f-a0c3613cafaf
    - Objective:01a007f6-d88e-77ce-84fc-dfb369785341
    - Objective:01a007f6-d88f-7b09-802d-7f55bba79a5b
    - Objective:01a007f6-d890-7919-b5ee-d492a5626454
outcome: Spectacular starts from a compact preflight, asks dual-layer questions, loads context progressively, recovers from CLI mistakes directly, and executes coherent clusters with FROST review.
preparation_baseline: 147776293fbd06d9a12093ec103bc940eeb87f7f
preparation_fingerprint: 0c84e87d40003c6f833c8b220cd0bbc8cb1383f0f5fc1f0025a8be111df3ff28
preparation_receipt: '{"schema_version":"spectacular.mission-preparation.v2","fingerprint":"0c84e87d40003c6f833c8b220cd0bbc8cb1383f0f5fc1f0025a8be111df3ff28","proposal":{"ref":"Proposal:01a007f6-d88a-7922-a16e-abd2262feda4","fingerprint":"8d6532a921265fba47d710e8f22c29941751ee60de3fe61b008734344832039a"},"baseline":"147776293fbd06d9a12093ec103bc940eeb87f7f","direction_sources":[{"ref":"Anchor:019fe381-5d61-7223-b362-03a5f99a7b14","fingerprint":"c45130e86ad48a4edfaebdda67ba4f79e3d51ec0f7c157b2520914ea3cfad93d"},{"ref":"Anchor:019fe381-5d61-7223-b362-03a5f99a7b15","fingerprint":"cdc1d055c63bd234b3c454c364b3ef6b8dc16e50cc77fd5daf0b9621b6f52456"},{"ref":"Contract:019fe381-5d61-7223-b362-03a5f99a7b10","fingerprint":"a64746f0376df57b13cdea0020d974c93a5b1be35fdbac0b5fa4922fbd983625"}],"candidates":[{"name":"lean-launch-context","outcome":"Spectacular starts from a compact preflight, asks dual-layer questions, loads context progressively, recovers from CLI mistakes directly, and executes coherent clusters with FROST review.","evidence":["Focused tests cover compact context, aliases, and corrective refusals.","Skill verification covers preflight, question, schema, cluster, and FROST guidance.","The complete repository verification gate passes."],"dependencies":["M2 implementation commit 1477762","current v2 Contract","current generated mechanical interface"],"cancellation_state":"Each green cluster remains a coherent local commit; later clusters may be omitted without weakening earlier behavior.","reversibility":"All changes remain on a stacked local feature branch until the integrated review boundary.","standalone_coherence":"Preflight, context, CLI recovery, and execution guidance form one token-efficiency path.","integration_path":"Guidance first, then compact context, CLI corrections, lean execution policy, generation, and acceptance.","learning_value":"Measures whether Spectacular can preserve governance while reducing context reloads and avoidable tool operations."}],"selected":"lean-launch-context","design_sufficiency":"sufficient","design_rationale":"The owner approved the retrospective findings, communication contract, preflight checklist, progressive disclosure model, CLI improvements, commit strategy, and FROST principles.","slice_quality":"coherent","slice_rationale":"Four cheap ordered clusters share a single launch-to-review efficiency objective and have focused checks.","blocking_gaps":[],"completion_criteria":[{"claim":"claim:cli-recovery","pass_boundary":"Common CLI mistakes accept safe intuitive aliases or return a runnable corrected command, and interface generation defaults to VERSION and canonical paths.","proof_requirement":"Deterministic command and generator tests cover project-dot validation, recovery text, and zero-argument canonical generation.","review_level":"automatic"},{"claim":"claim:launch-preflight-questions","pass_boundary":"The Skill checks .spectacular/PROJECT.md before root PROJECT.md, performs the approved read-only launch preflight, proposes one clean next action, and formats owner questions with plain and technical layers plus concrete consequences.","proof_requirement":"Skill conformance tests and a cold launch scenario verify the checklist, discovery order, and question grammar.","review_level":"clustered"},{"claim":"claim:lean-execution-review","pass_boundary":"The Skill defines compact cluster grammar, progressive context packets, focused-test then integration/full-gate sequencing, local green commits with one push boundary, the self-host schema rule, and FROST critic invariants.","proof_requirement":"Independent Skill review against the exact tree plus distribution and full repository verification.","review_level":"independent"},{"claim":"claim:progressive-context","pass_boundary":"Routine launch and Mission context can render a compact card while retaining exact authoritative pointers for progressive source drill-down.","proof_requirement":"Context and command tests prove compact output, stable source identity, and unchanged full JSON authority.","review_level":"clustered"}],"stop_conditions":["A change would weaken owner authority, evidence attribution, or cold recovery.","Compact output would become authoritative or omit exact source pointers.","A CLI alias would make a mutating target ambiguous.","The implementation would introduce v1 compatibility, generic record/search verbs, or automatic active-Mission migration."],"evidence_claims":["claim:launch-preflight-questions","claim:progressive-context","claim:cli-recovery","claim:lean-execution-review"],"fresh_until":"2026-08-19T23:59:59Z","ready":true,"unmet_requirements":[],"next":"owner-activation"}'
preparation_sources:
    - Proposal:01a007f6-d88a-7922-a16e-abd2262feda4@8d6532a921265fba47d710e8f22c29941751ee60de3fe61b008734344832039a
    - Anchor:019fe381-5d61-7223-b362-03a5f99a7b14@c45130e86ad48a4edfaebdda67ba4f79e3d51ec0f7c157b2520914ea3cfad93d
    - Anchor:019fe381-5d61-7223-b362-03a5f99a7b15@cdc1d055c63bd234b3c454c364b3ef6b8dc16e50cc77fd5daf0b9621b6f52456
    - Contract:019fe381-5d61-7223-b362-03a5f99a7b10@a64746f0376df57b13cdea0020d974c93a5b1be35fdbac0b5fa4922fbd983625
preparation_valid_until: "2026-08-19T23:59:59Z"
recovery_point: git branch codex/lean-launch-context at 147776293fbd06d9a12093ec103bc940eeb87f7f
repair_budget: 3
return_destination: owner in this Codex task
scope:
    - v2
slice_quality: coherent
stops:
    - A change would weaken owner authority, evidence attribution, or cold recovery.
    - Compact output would become authoritative or omit exact source pointers.
    - A CLI alias would make a mutating target ambiguous.
    - The implementation would introduce v1 compatibility, generic record/search verbs, or automatic active-Mission migration.
---
# Mission
