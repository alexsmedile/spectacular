---
type: Mission
id: 01a00a20-63dd-7bc6-b0a2-cb443fd6d194
title: Implement the compact Mission CLI
status: completed
created: "2026-08-16T10:31:30Z"
updated: "2026-08-16T12:30:56Z"
activation:
    at: "2026-08-16T10:53:23Z"
    by: Alex
    fingerprint: sha256:2bb4d9ff6f84db040db5eb7ecdbeb392f93aac0c242cca7bce1cfe04679ff7c5
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
        - secret-change
baseline:
    branch: codex/lean-launch-context
    commit: e0e2fdcc6874adb0750b9e07cae43d6f09febc6d
completion:
    - claim: typed-bundle
      pass_boundary: One current decoder resolves MISSION.md plus inline or promoted Objectives, Runs, and reviews into the same typed logical bundle while preserving canonical Markdown and existing v2 readability.
      proof_requirement: Golden fixtures and round-trip tests cover M5, M6, expanded bundles, existing v2 Missions, unknown fields, and exact source pointers without a second package root or file migration.
    - claim: schema-validation
      pass_boundary: The schema registry owns every mandatory identity, binding, fingerprint, claim, dependency, Run, review, authority, scope, layout, and transition check; Mission content cannot disable those checks.
      proof_requirement: Table-driven negative tests mutate one valid property at a time and assert stable refusal code, exact field, concrete problem, safe correction, and no writes.
    - claim: typed-commands
      pass_boundary: The accepted start, show/check, Objective, Run, review, and completion commands use compact Markdown input where meaning is required, generate mechanics automatically, and avoid the superseded multi-step ceremony.
      proof_requirement: Real-process tests cover compact and `--json` output, stdin/file input, generated identities/bindings, progress, promotion, R2 creation, independent review, owner completion, and cold resume.
    - claim: atomic-stress
      pass_boundary: Every mutating command is path-safe, atomic, concurrency-aware, and retry-stable across injected failures without duplicate identity or ambiguous inline/file state.
      proof_requirement: Fault injection and fuzz/property tests cover every write boundary, collision, dependency graph, YAML shape, path escape, stale tree, concurrent mutation, and idempotent retry.
completion_record:
    at: "2026-08-16T12:30:56Z"
    authorization: owner supplied --by after schema checks
    by: Alex
    review: RV1
    reviewed_commit: 7fc3436b00642e0d7ccaec8cbd6e8b5e52d69f21
contract:
    fingerprint: sha256:80336b159e296ba63b5d85c80a48f8e540ae07d9aac52cdcdba4730059378a48
    ref: Contract:01a00a20-63dd-7670-97f1-9eb8e12adc3a
dependencies:
    - M5 completed with independent review and owner acceptance.
gaps: []
objectives:
    - claims:
        - typed-bundle
        - schema-validation
      id: 01a00a20-63dd-7033-832a-a1e6de4388dd
      outcome: Build the shared typed Mission-bundle decoder, resolver, canonical writer, and schema registry.
      ref: O1
      status: implemented
    - after:
        - O1
      claims:
        - typed-commands
      id: 01a00a20-63dd-79dc-87ea-157563ccf6d6
      outcome: Implement the minimal read, start, progress, expansion, review, and completion command surface.
      ref: O2
      status: implemented
    - after:
        - O2
      claims:
        - atomic-stress
      id: 01a00a20-63dd-72df-b24e-4f1044c47e61
      outcome: Prove atomicity, safe refusals, representation equivalence, legacy readability, and compact distribution behavior.
      ref: O3
      status: implemented
outcome: Spectacular provides a small typed CLI that validates compact Mission bundles and performs only the repeated or atomic mechanics that are safer and cheaper than LLM-only execution.
owner: Alex
ref: M6
repair_budget: 3
review: independent
reviews:
    - file: reviews/RV1-m6-independent-completion-review.md
      id: 01a00a34-6cb8-7082-8673-d8191806b194
      ref: RV1
      verdict: pass
run:
    current_objective: O3
    id: 01a00a34-6ce4-73fd-b7d3-5cdad9302124
    operator: Codex primary session
    ref: R1
    repairs: 1
    started_at: "2026-08-16T10:53:23Z"
    status: completed
scope:
    mechanical:
        - cmd/spectacular/
        - internal/
        - test/
        - skills/spectacular/generated/
        - install/
        - .spectacular/
    semantic:
        - Typed decoding and validation of compact and expanded Mission bundles.
        - Minimal noun-first commands for Mission, Objective, Run, review, and completion mechanics.
        - Replacement of the superseded v2 CLI workflow without generic mutations or a parallel compatibility architecture.
stops:
    - The CLI begins authoring product meaning, grading proof sufficiency, or replacing owner judgment.
    - Routine use requires a large JSON payload, mandatory Proposal/Decision/receipt/index/reconciliation sequence, or generic mutation API.
    - Existing v2 readability requires rewriting canonical files, a second package root, or divergent decoders.
    - A mutation can partially commit, follow a path outside `.spectacular`, change stable identity, or return an ambiguous refusal.
validation:
    mode: cli
    schema: mission.v2
---
# Mission

## Activation boundary

M6 activated after M5 received independent review and owner completion. Its
exact Git baseline, owner/time, frozen semantic-envelope fingerprint, and inline
R1 are recorded above. Validation moved from `manual-bootstrap` to `cli` once
the replacement CLI could check this Mission and the completed M5 without
rewriting either file.

## Why mechanics earn tooling

The CLI is justified where deterministic reuse saves work: schema checks,
identity/ref allocation, exact fingerprints, dependency graphs, safe paths,
atomic transitions, retries, concurrency, projections, and refusals. These are
cheap to test once and expensive to repeatedly reconstruct with an LLM.

The LLM remains the faster surface for interpreting intent, drafting the plan
and Markdown, writing criteria, choosing decomposition, understanding semantic
scope, diagnosing novel failures, and deciding whether detail earns another
file. The CLI checks and applies; it does not become a form-filling substitute
for reasoning.

## O1 — Typed bundle and schema registry

- Add one Mission-bundle type used by discovery, show/check, validation, and
  every mutation.
- Normalize current v2 navigation fields inside that decoder without a second
  reader, and write only the new `ref` form.
- Resolve inline and promoted Objective/Run/review records identically.
- Preserve Markdown bodies and unknown non-authoritative fields on round trip.
- Make mandatory validators registry-owned. The Mission may request extra
  checks but cannot weaken the schema.
- Define the activation fingerprint over outcome, review, completion,
  authority, scope, repair budget, dependencies, Gaps, and stops; exclude
  mutable progress.

## O2 — Minimal commands

Implement only:

```text
mission start <plan.md|->
mission show <ref>
mission check <ref>
objective show <ref>
objective promote <ref>
objective finish <ref>
run show <ref>
run start <mission-ref> --title <title>
review record <mission-ref> <review.md|->
mission complete <ref> --by <owner>
```

Mission plans carry meaning in Markdown/YAML. Start supplies identities, refs,
timestamps, exact Contract/Git bindings, activation, R1, retry identity, and the
canonical path. `--json` changes output only.

Promote without identity change. Start R2 by atomically materializing R1 and
creating R2. Record independent review only when reviewer/operator separation
and exact reviewed tree are demonstrable. Complete only after a fresh check and
attributable owner confirmation.

## O3 — Stress and integration

- Mutate one valid field at a time across YAML, identity, Contract/baseline,
  activation, completion, review, Objectives, Runs, authority, scope, Gaps,
  stops, and layout.
- Assert `code · field/target · problem · safe correction` and byte-identical
  files after every refusal.
- Inject failure after every write boundary in start, promote, Run start,
  review record, and complete.
- Fuzz dependency DAGs against an independent topological oracle, and cover
  YAML trees, UUID/ref collisions, paths, retry/concurrent operations, and
  preservation round trips with exact table-driven cases.
- Prove inline/promoted representation equivalence and cold recovery.
- Run focused package tests by cluster, then one compact full verification and
  distribution gate before independent review.

## Repair 1 — Independent-review findings

The first exact-tree review found that green integration checks did not yet
prove the full frozen boundary. The bounded repair resolves Review records into
the typed bundle, verifies real Git commit-to-tree provenance, serializes fresh
allocation under the mutation lock, makes retry identities converge, refuses
occupied derived targets, revalidates the completed candidate, and exercises
all five mutating commands at every transaction install boundary.
