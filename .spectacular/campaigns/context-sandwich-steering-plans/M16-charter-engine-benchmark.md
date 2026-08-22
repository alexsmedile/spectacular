---
type: MissionPlan
title: Build and benchmark the bounded charter engine
owner: Alex
contract:
  ref: Contract:01a00a20-63dd-7670-97f1-9eb8e12adc3a
outcome: Spectacular has an internal read-only compiler that retrieves only declared governance sources for one Objective and proves material context savings without weakening behavior.
review: independent
completion:
  - claim: deterministic-charter-engine
    pass_boundary: The internal compiler retrieves the bound Contract, then typed refs in Mission and Objective sources lists in declaration order with duplicate removal, then explicit invocation-added refs; it follows no prose link, recursive citation, or semantic guess and performs no canonical write.
    proof_requirement: Table and golden fixtures cover ordering, duplicate removal, missing and malformed refs, conflicting refs, explicit additions, source attribution, byte-stable identical inputs, forbidden traversal, no whole-ledger body scan, and byte-identical canonical state.
  - claim: bounded-governance-envelope
    pass_boundary: spectacular-charter-tokenizer.v1 counts byte-exact valid UTF-8 with its bundled o200k_base data and published digest; safe compaction returns normally at at most 1200 tokens, warns from 1201 through 1400, strongly recommends a split and requires explicit Orchestrator approval from 1401 through 1440, and refuses above 1440 without omitting authority, stops, claims, proof, or writable scope.
    proof_requirement: Cross-platform tokenizer goldens, digest and invalid-UTF-8 tests, plus boundary fixtures at 1200, 1201, 1400, 1401, 1440, and 1441 assert exact count, disposition, and preserved frozen fields; compaction fixtures prove only derived summaries are shortened.
  - claim: measured-context-economy
    pass_boundary: Against M14's immutable paired baseline, total context ingestion falls by at least 40 percent with no regression in safety, task success, recovery, routing, or Decision fidelity; envelope, named-source, and repair tokens are reported separately.
    proof_requirement: Trials pin M14 reviewed commit 3a5bab29ed1941a9d2b9873e63324a5bbf29620d and recorded baseline 14158f92aaa78e1a90f0231afdc52a04f3cca6c3, plus identical model, seed, case set, pair order policy, and repetition policy; missing, contaminated, or regressing telemetry returns INCONCLUSIVE or REGRESSION, never pass.
objectives:
  - outcome: Define and implement explicit source traversal and the three-layer internal compiler.
    claims: [deterministic-charter-engine]
  - outcome: Implement safe compaction and the owner-set token dispositions.
    claims: [bounded-governance-envelope]
  - outcome: Extend the immutable M14 harness and run the paired proof gate.
    claims: [measured-context-economy]
authority:
  operator: [inspect, edit-in-scope, choose-reversible-implementation, run-checks, generate-derived-files, bounded-repair, commit-local]
  requires_owner: [activate-mission, change-outcome-or-completion, expand-scope, push, merge, release, irreversible-change]
scope:
  mechanical: [internal/runtime, internal/missionbundle, test/evals/spectacular, test]
  semantic: [read-only Objective charter compilation, explicit source retrieval, context-economy benchmark]
repair_budget: 2
dependencies: [M15 accepted; M14 benchmark fixtures remain immutable]
gaps: []
stops: [public-command-change, canonical-charter-write, semantic-source-inference, context-reduction-below-40-percent, behavioral-regression, benchmark-contamination]
---

# Mission

> **Future Mission sketch.** Preserve as design input. Do not activate, maintain,
> validate, or review as a final plan until M15 closes and the Orchestrator
> re-prepares this block from current Evidence.

Build the engine behind a private boundary first. This Mission adds neither
`charter` nor `decide`. Its Evidence determines whether public exposure is allowed.
The charter remains temporary retrieval output, not the complete frozen Handoff.
