---
type: spectacular-handoff-return
schema_version: spectacular.handoff-return.v2
handoff_id: scenario-a-owner-decision-implementation
handoff_hash: 738817867b750b994b4038b48b99b9e5a5d0dd659da7d040eb5833fe37c7b700
status: complete
mission: scenario-a-cold-recovery
baseline_commit: 243ae3e376e6eb30e59dc28bb691a98dfc6a7b92
baseline_tree: c0922b35d84617ea82e30c63b6b70729dfa87110
bounced_head: 816c34aec885747c5f4a7d40274a35b34f91becf
bounced_tree: 5a79f91dc98949292d8cdef20005017d4a52b347
repair_commit: 151b255e57929949f91c336413b528a8b714fed8
repair_tree: 0e1c2ab7ab4c8200d43731a8282420e934bce274
reviewed_commit: c9efd998e768ac2ba0cdf871acc5368fb35dae05
reviewed_tree: 4a99dbb1ecb0a6363d3cf8744c18d0142920b4d8
evidence_head: ae7e84df900e6747749fdc4dbc2d52bed0094a75
evidence_tree: 6f2f411b900e86fb16682751863d8b2f88d51b3f
branch: codex/feat/v2-scenario-a-cold-recovery
central_disposition_requested: accept
next_action: central-assess-scenario-a-accept-bounce-or-escalate
---

# Scenario A final implementation return

```yaml
return:
  schema_version: spectacular.handoff-return.v2
  handoff_id: scenario-a-owner-decision-implementation
  handoff_hash: 738817867b750b994b4038b48b99b9e5a5d0dd659da7d040eb5833fe37c7b700
  status: complete
  baseline:
    commit: 243ae3e376e6eb30e59dc28bb691a98dfc6a7b92
    tree: c0922b35d84617ea82e30c63b6b70729dfa87110
    dirty_state: clean
  bounced_return:
    commit: 816c34aec885747c5f4a7d40274a35b34f91becf
    tree: 5a79f91dc98949292d8cdef20005017d4a52b347
  repair:
    commit: 151b255e57929949f91c336413b528a8b714fed8
    tree: 0e1c2ab7ab4c8200d43731a8282420e934bce274
  final_review_target:
    commit: c9efd998e768ac2ba0cdf871acc5368fb35dae05
    tree: 4a99dbb1ecb0a6363d3cf8744c18d0142920b4d8
  review_evidence_head:
    commit: ae7e84df900e6747749fdc4dbc2d52bed0094a75
    tree: 6f2f411b900e86fb16682751863d8b2f88d51b3f
  branch: codex/feat/v2-scenario-a-cold-recovery
  input_refs:
    - EXECUTABLE-REFACTOR-PROGRAM-CONTRACT.md@1.7
    - MVP-SCENARIO-CLI-SEQUENCING-DECISION.md@1.0
    - M1-SEMANTIC-SUBSTRATE-MISSION-CHARTER.md@1.2
    - SCENARIO-A-COLD-RECOVERY-MISSION-CHARTER.md@1.0
    - evidence/scenario-a-owner-dispositions.md@1.0
    - evidence/scenario-a-repair-delta.md@1.0
  upstream_contracts:
    - SHARED-SCAFFOLD-CONTRACT.md@1.0
    - MISSION-PREPARATION-CONTRACT.md@1.0
    - PRODUCT-TRUTH-CONTRACT-MODEL.md@1.0
    - WORK-UNIT-LIFECYCLE-CONTRACT.md@1.0
    - EXECUTION-AUTHORITY-CONTRACT.md@1.0
    - EVIDENCE-CLOSURE-CONTINUITY-CONTRACT.md@1.0
    - RETRIEVAL-AND-EARNED-WORKSPACE-CONTRACT.md@1.0
    - PUBLIC-LANGUAGE-AND-INTERFACE-CONTRACT.md@1.0
    - SUCCESS-EVIDENCE-CONSTITUTION.md@1.0
  reviewer: /root/scenario_a_repair_reviewer
  read_set:
    - exact original-baseline, bounced-head, repair, and reviewed-target diffs
    - all changed Scenario A v2 code, tests, fixtures, and self-host workspace
    - Scenario A charter, owner dispositions, design, prior proof/review, repair delta/proof
    - accepted program, M1 charter, and relevant immutable contracts
    - direct focused, full-matrix, cold, self-host, performance, nonmutation, and scope evidence
  result: >-
    Scenario A implements the accepted ten-command, read-only cold-recovery slice. The owner-authorized
    bounce repair closed all three prior P1 defects: known-command usage is decided before discovery,
    authority markers cannot escape through symlinks before parsing, and Proposal sources retain exact
    drill-down metadata without a fabricated command. All required proof is green and one fresh read-only
    reviewer returned accept with no P0, P1, or P2 findings against the exact reviewed commit/tree.
  decisions:
    - Owner accepted option A for all eight original Scenario A decision clusters.
    - Central orchestration bounced the first return and authorized exactly one cohesive repair round for three named P1 findings.
    - The primary agent owned diagnosis, edits, tests, commits, evidence, and return; exactly one fresh read-only reviewer was used.
    - Proposal record sources expose noun, typed ref, canonical path, and fingerprint but omit show_command because the registry has no Proposal operation.
  facts:
    - The registry exposes exactly ten noun-first read operations and exit classes 0/2/3.
    - Argument shapes, dispatch, help, effects, schemas, and supported pointer commands derive from the registry.
    - The required Go format/module/vet/test/race/build matrix passes.
    - Bash syntax, version guard, and the terminal full v1 suite pass; 31 files, 0 failed.
    - Focused real-filesystem and adversarial regressions for all three bounced findings pass.
    - Existing deterministic, golden, lookup/order/projection, refusal, owner-gate, and byte/path/mode/mtime nonmutation tests pass.
    - Fresh direct cold recovery used 10 calls, 20,983 JSON bytes, and a 3,029-byte orientation card.
    - Primary self-host recovery used 7 calls and 16,661 bytes; independent review used 7 calls and 14,426 bytes; both found the same exact resume continuation.
    - Primary 20-run p95 was 3.020 ms and independent p95 was 3.134 ms, below 500 ms.
    - The original-baseline diff contains only Scenario A charter/evidence and v2 paths; no v1 path changed.
    - No new dependency, provider effect, cache, migration, compatibility path, release, Scenario B transition, or Scenario C reconciliation exists.
  assumptions:
    - Pre-read Lstat containment evidence covers the declared non-concurrent filesystem model; no hostile concurrent path-replacement claim is made.
    - Performance is host-local to the declared fixture and execution host.
    - B+C-only mutation, recovery, authorization, and idempotency fields remain future additive fields.
  artifacts:
    - evidence/scenario-a-repair-delta.md
    - evidence/scenario-a-repair-proof.md
    - evidence/scenario-a-repair-independent-review.md
    - commit:8c86a948e39574abacd349ff0f5e1a56e6c1929d
    - commit:151b255e57929949f91c336413b528a8b714fed8
    - commit:c9efd998e768ac2ba0cdf871acc5368fb35dae05
    - commit:ae7e84df900e6747749fdc4dbc2d52bed0094a75
  evidence:
    - All Scenario A repair, accepted-program, M1, and relevant immutable input sidecars verified.
    - Primary repair methods/results are recorded in evidence/scenario-a-repair-proof.md.
    - Fresh reviewer methods, exact findings, independent outputs, and accept verdict are recorded in evidence/scenario-a-repair-independent-review.md.
  conflicts: []
  scope_deviations: []
  recovery_point: >-
    Original baseline 243ae3e plus the clean reviewed proof target c9efd99 and review-evidence head
    ae7e84d; all prior bounce evidence is preserved without reset, stash, merge, push, or cleanup.
  central_disposition_requested: accept
  next_action: >-
    Central orchestration must assess this exact packet and choose accept, bounce, or escalate for
    Scenario A. Do not treat this request as self-acceptance and do not authorize Scenario B unless
    central separately activates it after Scenario A disposition.
```

## Exact command contract

```text
spectacular anchor show project [--json]
spectacular mission list [--json]
spectacular mission show <ref> [--json]
spectacular gap list --scope <ref> [--json]
spectacular gap show <ref> [--json]
spectacular run show <ref> [--json]
spectacular checkpoint show <ref> [--json]
spectacular evidence show <ref> [--json]
spectacular decision show <ref> [--json]
spectacular workspace validate <scope> [--json]
```

No public `record`, Proposal command, generic alias, compatibility/fallback path, write operation, or
later-scenario behavior exists.

## Compatibility/join finding

Scenario A's implemented envelope serves the accepted B+C invariants without redesign and with no
conflict: versioned results/refusals, stable refs, source/fingerprint and generation basis, registry-
owned public mechanics, warnings separable from data, and additive space for later refusal and
optimistic-concurrency fields. Scenario A remains read-only and does not implement B/C semantics.

## Prohibited conclusion

This return requests central `accept`; it does not itself accept or resolve Scenario A, mutate central
lifecycle, reconcile current truth, or authorize Scenario B.
