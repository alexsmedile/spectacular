---
type: Mission
id: 01a02b9d-8290-7ec3-8098-be9bcd2e11dd
title: Build and benchmark the bounded charter engine
status: completed
created: "2026-08-22T22:36:31Z"
updated: "2026-08-22T22:40:37Z"
activation:
    at: "2026-08-22T22:36:31Z"
    by: Alex
    fingerprint: sha256:4b0519b664e6be9d6d234a67f6eb4adf6699d322c9ed7e4268284e7dea2e940f
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
baseline:
    branch: codex/m16-charter-engine-benchmark
    commit: f6658f4877f04830ec537e86cb891a9d2bb80f6c
completion:
    - claim: explicit-source-compiler-implemented
      pass_boundary: The internal compiler assembles a 3-layer charter from the bound Contract, Mission/Objective sources, and explicit invocation refs in order without semantic inference or canonical writes.
      proof_requirement: Unit and golden tests verify exact source ordering, duplicate deduplication, missing ref handling, and read-only behavior.
    - claim: charter-tokenizer-and-thresholds-enforced
      pass_boundary: Bundled o200k_base UTF-8 tokenizer counts tokens byte-exactly and enforces standard thresholds (<=1200 pass, 1201-1400 warn, 1401-1440 split notice, >1440 refusal).
      proof_requirement: Tokenizer test suite and boundary test fixtures at 1200, 1201, 1400, 1401, 1440, and 1441 assert exact count and disposition codes.
    - claim: safe-compaction-rules-verified
      pass_boundary: When token limits are approached, safe compaction shortens derived summaries and rationale while strictly preserving claims, authority, stops, and writable perimeters.
      proof_requirement: Compaction fixtures prove frozen fields remain byte-intact while descriptive prose is compacted.
    - claim: context-reduction-benchmark-proven
      pass_boundary: Paired benchmark run against M14 fixtures demonstrates at least 40% reduction in total context ingestion with zero regression in safety, routing, task success, or decision fidelity.
      proof_requirement: Benchmark report outputs paired telemetry with envelope, source, and diagnostic token breakdowns and zero regression status.
completion_record:
    at: "2026-08-22T22:40:37Z"
    authorization: owner supplied --by after schema checks
    by: Alex
    review: RV1
    reviewed_commit: b6d4bf72a70e2d5edf9447ca21bfa1fca6f04171
contract:
    fingerprint: sha256:332adfc32551aba44a790add44aa6db92ed86bafe08799269641657cd9e1db52
    ref: Contract:01a00a20-63dd-7670-97f1-9eb8e12adc3a
dependencies:
    - M15 completed with independent review
gaps: []
objectives:
    - claims:
        - explicit-source-compiler-implemented
      id: 01a02b9d-8290-79f7-a700-a73502385ea7
      outcome: Implement the internal 3-layer charter compiler and source traversal engine (Compiler Pillar).
      ref: O1
      status: implemented
    - claims:
        - charter-tokenizer-and-thresholds-enforced
        - safe-compaction-rules-verified
      id: 01a02b9d-8290-7af1-9298-a4e4bfd5b3fc
      outcome: Implement the o200k_base tokenizer and safe compaction threshold engine (Budget Pillar).
      ref: O2
      status: implemented
    - claims:
        - context-reduction-benchmark-proven
      id: 01a02b9d-8290-724b-8b0c-03bd11099ab5
      outcome: Run the paired effectiveness benchmark and prove >=40% context savings (Proof Pillar).
      ref: O3
      status: implemented
outcome: Unlock sub-1,200 token agent execution that cuts context costs by >=40% and stops prompt escapes by building the internal read-only Context-Sandwich compiler and proving it against the M14 benchmark suite.
owner: Alex
ref: M16
repair_budget: 2
review: independent
reviews:
    - file: reviews/RV1-independent-review-of-m16-charter-compiler-and-effectiveness-benchmark-proof.md
      id: 01a02b9e-5198-782f-b057-902764cc31c4
      ref: RV1
      verdict: pass
run:
    current_objective: O1
    id: 01a02b9d-8290-77aa-a1c3-85b2d5e6c94c
    operator: Alex
    ref: R1
    repairs: 0
    started_at: "2026-08-22T22:36:31Z"
    status: completed
scope:
    mechanical:
        - internal/charter/
        - internal/missionbundle/
        - test/evals/spectacular/
    semantic:
        - Read-only 3-layer charter compilation
        - o200k_base token counting and threshold enforcement
        - Context reduction and safety benchmark proof
start_key: sha256:17f6026c28c21dc19761fd50c54ecf802e7bea956ac68499fef5f217aa36161e
stops:
    - public-command-change
    - canonical-charter-write
    - semantic-source-inference
    - context-reduction-below-40-percent
    - behavioral-regression
validation:
    mode: cli
    schema: mission.v2
---
# Mission: Build and Benchmark the Bounded Charter Engine

## User Superpower (The Hub)
Subagents receive a razor-sharp, sub-1,200 token context envelope containing 100% of the truth they need and 0% of the noise. This eliminates whole-workspace prompt dumping, cuts LLM execution token costs by $\ge 40\%$, and prevents hallucinated scope escapes outside declared target files.

## Technical Pillars (The Spokes)
1. **Compiler Pillar (`internal/charter/`)**: Assembles Layer 1 (Frozen Truth) + Layer 2 (Steering & Decisions) + Layer 3 (Perimeter).
2. **Budget Pillar (`internal/charter/tokenizer/`)**: Bundles byte-exact `o200k_base` counting and enforces the 1200 / 1400 / 1440 token gates with safe compaction.
3. **Proof Pillar (`test/evals/spectacular/`)**: Benchmarks the compiler against M14 fixtures to prove $\ge 40\%$ context reduction with zero behavioral regression.

## Key Deliverables & Actions

### 1. 3-Layer Charter Compiler Package (`internal/charter/`)
- Build the core compiler function `Compile(workspace, missionRef, objectiveRef, extraSources) -> (*Charter, error)`.
- Retrieve sources in strict order: Bound Contract $\to$ Mission `sources:` $\to$ Objective `sources:` $\to$ invocation refs.
- Ensure read-only execution: compiling a charter makes zero canonical disk writes.

### 2. Tokenizer & Compaction Engine (`internal/charter/tokenizer/`)
- Embed `o200k_base` BPE merge table and vocabulary.
- Implement token counting:
  - $\le 1{,}200$ tokens: `Pass`
  - $1{,}201–1{,}400$ tokens: `Warn` + safe compaction
  - $1{,}401–1{,}440$ tokens: `SplitRecommended`
  - $> 1{,}440$ tokens: `Refusal`
- Implement safe compaction to condense Decision rationale while preserving all claims, authority, and file perimeters byte-exact.

### 3. Effectiveness Benchmark Runner (`test/evals/spectacular/`)
- Add paired benchmark cases evaluating full workspace scan vs. compiled context charter.
- Assert $\ge 40\%$ total context ingestion reduction.
- Assert zero regression across task success, routing pass rate, safety checks, and recovery.
