---
type: Contract
id: 01a00a20-63dd-7670-97f1-9eb8e12adc3a
ref: CC-missioncli
title: Spectacular mechanical Mission CLI
status: current
owner: Alex
created: "2026-08-16T10:31:30Z"
updated: "2026-08-16T10:38:58Z"
contract_version: "5"

purpose: Supply deterministic Mission mechanics where exact repeated enforcement saves LLM work, tokens, and recovery cost.
outcome: The plan supplies meaning, the CLI supplies invariants, and both operate one canonical compact Mission bundle without duplicate vocabularies or ceremony.

applies_when:
  - A compact Mission is checked, started, resumed, expanded, reviewed, progressed, or completed.
  - A repeated invariant or multi-file transition requires deterministic enforcement.
does_not_apply_when:
  - Exploring fuzzy intent, choosing product behavior, wording semantic criteria, decomposing contextual work, or writing explanatory prose.
  - A supported command would cost more interpretation or ceremony than directly drafting one canonical Markdown file and checking it afterward.
does_not_provide:
  - Owner judgment, product meaning, proof sufficiency, reviewer independence by assertion, provider permission, or autonomous runtime execution.

required_behavior:
  - Decode one typed Mission bundle shared by discovery, validation, projections, and mutations.
  - Treat `MISSION.md` as the bundle entry point and resolve inline or promoted Objectives and Runs to one logical representation.
  - Support optional ordered typed `sources:` references in Mission and Objective frontmatter as frozen semantic retrieval input.
  - Let Mission plans supply semantic fields while mechanics generate UUIDv7 identities, refs, timestamps, Contract fingerprints, Git baseline, activation, retry identity, and canonical paths.
  - Accept a Mission plan from Markdown file or stdin for routine start; do not require a parallel JSON authoring vocabulary.
  - Treat that plan as command input, not a required persisted Proposal or a second canonical artifact.
  - Own mandatory validation in the schema registry; a Mission may add checks but cannot remove core invariants.
  - Fingerprint the frozen semantic envelope at activation while excluding mutable Mission, Objective, Run, and repair progress.
  - Preserve UUID/ref identity across Objective promotion and Run materialization.
  - Record earned reviews under the Mission with reviewed tree, reviewer identity and operator relation, claim verdicts, findings, limitations, and time.
  - Require actual independence when `review: independent`; a fresh session or self-declaration alone is insufficient.
  - Perform start, promotion, Run creation, review recording, and completion as atomic typed transitions with safe retry and concurrency refusal.
  - Return compact success receipts and exact refusals containing code, field, problem, and safe correction.
  - Keep read-only `show` and `check` free of canonical writes.
  - Complete through frozen criteria, required proof/review, applicable specification edits, and one attributable owner gate without separate Contract reconciliation.
  - Refuse path traversal, symlink escape, duplicate identity/ref, stale bindings, invalid dependency graphs, ambiguous inline/file state, unauthorized effects, and unsupported transitions.
  - Avoid generic record mutation, mandatory Proposal, preparation receipt, lifecycle Decision, per-Mission index, and compatibility package roots.

command_surface:
  - mission start <plan.md|->
  - mission show <ref>
  - mission check <ref>
  - objective show <mission-ref>/<objective-ref>
  - objective promote <mission-ref>/<objective-ref>
  - objective finish <mission-ref>/<objective-ref>
  - run show <mission-ref>/<run-ref>
  - run start <mission-ref>[/<objective-ref>] --title <title>
  - run transition <target-ref> --to <state> --by <actor> --reason <text>
  - review record <mission-ref> <review.md|->
  - handoff record <mission-ref> <handoff.md|-> --by <sender>
  - mission complete <ref> --by <owner>
  - proposal check <ref>
  - campaign check <campaign-path>
  - contract amend <contract-ref> --gap <gap-ref> --by <owner>
  - charter <mission-ref>/<objective-ref> [sources...]
  - decide <decision.md|->

mandatory_validation:
  - contract-version
  - resolved-gap-integrity
  - yaml-schema
  - uuidv7-identity
  - reference-integrity
  - contract-binding
  - baseline-binding
  - activation-fingerprint
  - completion-claim-coverage
  - objective-dependency-dag
  - run-state
  - review-independence
  - authority-vocabulary
  - mechanical-scope
  - safe-file-layout
  - transition-atomicity

stress_properties:
  - A failed validation or transition leaves the canonical tree byte-identical.
  - Retrying the same logical start or transition converges on the same identities without duplicates.
  - Inline and promoted representations produce the same show, dependency, claim, and completion results.
  - A changed frozen semantic field invalidates the activation fingerprint; mutable progress does not.
  - Every malformed field and illegal transition produces a stable typed refusal with one actionable correction.
  - Existing v2 Missions remain inspectable through one current decoder without rewriting their files.

conformance_checks:
  - Golden M5 and planned M6 fixtures decode through the same typed bundle model.
  - Table-driven mutations cover each mandatory validation domain and exact refusal.
  - Fault injection covers every write boundary in start, promote, Run start, review record, and complete.
  - Fuzzing covers YAML shape, identity/ref collisions, dependency graphs, paths, and round-trip preservation.
  - Real-process tests prove compact default output and machine-readable `--json` output without making JSON the routine authoring format.

gaps:
  - ref: gap-rewrite-matches-by-line
    problem: >-
      The Gap rewrite in the amendment path locates a Gap by ref and then walks forward for a line
      matching `blocked_on:`, splicing that key and its indented continuation. The textual approach is
      deliberate — decoding and re-emitting canonical YAML would reflow every block scalar in the
      Contract, and an amendment whose diff touches prose it did not change is not reviewable — but the
      match does not know when it is inside a scalar body. A Gap whose `problem:` is itself a block
      scalar containing the literal text `blocked_on:` would collide. The amendable-field guard limits
      the blast radius to the `gaps:` block, which is exactly what an amendment may change, so it does
      not catch this. Confirmed as a documented limitation by M11's independent review.
    resolution: >-
        The Gap rewrite tracks block-scalar depth while walking a Gap entry, so a `blocked_on:`
        appearing inside a scalar body is not mistaken for the key. The textual approach is kept
        deliberately, so an amendment's diff still touches only what it changed. An adversarial
        fixture carrying the literal text inside a `problem:` scalar asserts the correct key is
        rewritten and the scalar body is left byte-identical.
  - ref: repoint-assumes-one-fingerprint
    problem: >-
      Re-pointing a bound Mission replaces the first occurrence of the old Contract fingerprint in the
      raw Mission file. Every Mission in the workspace carries its binding exactly once, so the
      assumption holds today and nothing enforces it. A Mission that quoted its own bound fingerprint in
      prose, a pass boundary, or a rejected approach could have the wrong occurrence rewritten and would
      still parse. M9 shows the shape is plausible: its body quotes a stale-fingerprint refusal
      containing three distinct fingerprint values, none of them its own binding.
    resolution: >-
        Re-pointing a bound Mission refuses when the old Contract fingerprint appears more than
        once in the Mission file, naming the Mission, the fingerprint, and every occurrence,
        rather than rewriting the first one. The refusal was chosen over anchoring to the
        `contract:` block because it is smaller and turns a silent corruption into a stated
        problem in a mechanism that rewrites records the owner is not reading. Anchoring remains
        available later if the refusal proves noisy.
---
# Spectacular mechanical Mission CLI

## Division of work

LLMs are fastest and most useful when interpreting intent, drafting concise
canonical Markdown, forming criteria, choosing coherent Objectives, and solving
contextual problems. Mechanical tooling is most useful when the same exact rule
must be enforced repeatedly or when partial failure would corrupt state.

The CLI therefore validates and performs typed transitions. It does not force
all Mission authorship through a large transport schema, invent owner intent, or
replace semantic review.

## Refusal contract

Every refusal answers four things:

```text
code · exact field/target · concrete problem · safe correction
```

A refusal does not mutate canonical files. Detailed diagnostics are available
with `--json`; the default remains a compact, directly runnable explanation.
