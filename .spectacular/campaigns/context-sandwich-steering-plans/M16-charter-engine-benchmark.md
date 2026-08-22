---
type: MissionPlan
title: Build and benchmark the bounded charter engine
owner: Alex
contract:
  ref: Contract:01a00a20-63dd-7670-97f1-9eb8e12adc3a
outcome: Unlock sub-1,200 token agent execution that cuts context costs by >=40% and stops prompt escapes by building the internal read-only Context-Sandwich compiler and proving it against the M14 benchmark suite.
review: independent
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

objectives:
  - outcome: Implement the internal 3-layer charter compiler and source traversal engine (Compiler Pillar).
    claims: [explicit-source-compiler-implemented]
  - outcome: Implement the o200k_base tokenizer and safe compaction threshold engine (Budget Pillar).
    claims: [charter-tokenizer-and-thresholds-enforced, safe-compaction-rules-verified]
  - outcome: Run the paired effectiveness benchmark and prove >=40% context savings (Proof Pillar).
    claims: [context-reduction-benchmark-proven]

authority:
  operator: [inspect, edit-in-scope, choose-reversible-implementation, run-checks, generate-derived-files, bounded-repair, commit-local]
  requires_owner: [activate-mission, change-outcome-or-completion, expand-scope, push, merge, release, irreversible-change]
scope:
  mechanical:
    - internal/charter/
    - internal/missionbundle/
    - test/evals/spectacular/
  semantic:
    - Read-only 3-layer charter compilation
    - o200k_base token counting and threshold enforcement
    - Context reduction and safety benchmark proof
repair_budget: 2
dependencies: [M15 completed with independent review]
gaps: []
stops:
  - public-command-change
  - canonical-charter-write
  - semantic-source-inference
  - context-reduction-below-40-percent
  - behavioral-regression
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
