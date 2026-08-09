---
type: independent-review-handoff
schema_version: spectacular.handoff.v2
handoff_id: H28
mission: M1
mode: read-only-fresh-context-skeptic
authority: advisory-assurance
status: authorized-for-dispatch
review_base_commit: f628e8ca895c47b0e7ee7142b98b89dfb16a1e0f
review_base_tree: 08172f732cc23c2dd5857a8bc1c18de0429a88e4
review_head_commit: 489bd6008e1720e4b0310b999a0bac02c62df6dc
review_head_tree: 911adb7081ea70b10d29539cc76a343c6d58b0cf
review_branch: codex/feat/v2-semantic-substrate
date: 2026-08-09
---

# H28 — M1 independent implementation review

## Outcome

Independently determine whether H27 satisfies the complete M1 charter with fresh primary evidence.
Return ranked findings and `accept | bounce | escalate`. Do not edit, repair, commit, accept M1,
merge, push, or activate M2.

## Binding inputs

- `M1-SEMANTIC-SUBSTRATE-MISSION-CHARTER.md@1.1` — SHA-256
  `a5340b0a63648585a117736d638a7ea0d4de58ae6110ac24e77bcc23babac98f`
- `EXECUTABLE-REFACTOR-PROGRAM-CONTRACT.md@1.6` — SHA-256
  `b96fa5e248e6e1a2c2c8c5d1929e17f94202b2d12e019866358ef7582f53a73d`
- `SHARED-SCAFFOLD-CONTRACT.md@1.0` — SHA-256
  `698997f12972d0b5a186f4d8b8c35753a642cc0454e8ef60f15de81590435d36`
- `EVIDENCE-CLOSURE-CONTINUITY-CONTRACT.md@1.0` — SHA-256
  `7dd763b4fa1a919924e24105790382a51414a0a8ee0222178dee7c9224f11ca9`
- `H27-m1-semantic-substrate-implementation.md` — SHA-256
  `0727b8df84ed712f10920bfc9798f76af98c3d30cf499150930374274737b1b7`

Read `AGENTS.md`, `.spectacular/AGENTS.md`, the binding inputs, the exact H27 diff, all changed code
and tests, `go.mod`/`go.sum`, and only the references necessary to resolve a finding. Treat H27's
return and tests as claims to verify, not authority.

## Launch

Review the exact immutable head `489bd6008e1720e4b0310b999a0bac02c62df6dc` against
`f628e8ca895c47b0e7ee7142b98b89dfb16a1e0f`. A detached clean checkout at the exact head/tree is
expected. Stop on another commit/tree, tracked dirt, missing objects, or input-hash mismatch.

This review is strictly read-only. Do not switch/create a branch, modify source or tests, install or
update dependencies, rewrite `go.sum`, commit, push, open a PR, or call providers. Running checks
that use normal Go caches and disposable test directories is allowed.

## Review rubric

### 1. Authority and scope

- Exactly the authorized paths changed; no v1, root implementation, generated surface, provider,
  cache, CLI/Skill, lifecycle, assessment, reconciliation, release, or migration behavior leaked in.
- Domain, Workspace, and Index each own only their declared invariants; canonical ordering or other
  authority is not duplicated.
- APIs preserve the M1→M2/M3 join without implementing downstream behavior.

### 2. Semantic correctness

- `type` and canonical UUIDv7 `id` are the only universal requirements.
- Proposal/Mission known fields and `source: "Proposal:<UUIDv7>"` are validated deterministically.
- Unknown valid YAML values and the opaque Markdown body survive semantic round-trip.
- Canonicalization and SHA-256 fingerprints ignore only the approved formatting differences and
  change for all meaning changes.
- Exact ID/path lookup and typed relationship validation remain deterministic across discovery
  order and rename.
- Duplicate, broken, and wrong-type states refuse safely and deterministically.

### 3. Persistence and failure safety

- Same-directory replacement cannot corrupt or silently lose the original on the exercised failure
  paths.
- Temporary files, permissions, close/sync/rename behavior, cleanup, and error propagation are
  appropriate for the declared macOS/Linux boundary.
- Tests genuinely induce failure rather than merely asserting a mock or unreachable branch.

### 4. Evidence quality

- Rerun every required command from `v2/` and report exact results.
- Map each charter scenario to primary code plus a meaningful test; identify false-positive,
  overfitted, or missing cases.
- Inspect the two repair claims and confirm deterministic ownership/refusal behavior at final head.
- Confirm final branch cleanliness, commit count, dependency pins, and diff scope.

### 5. Dependency boundary

`go.yaml.in/yaml/v4 v4.0.0-rc.6` is a release candidate while H27 requested stable pinned versions.
Determine explicitly whether:

1. the RC is a blocking authority/envelope deviation;
2. it is an acceptable activation-time Type-2 implementation pin with bounded risk; or
3. it requires an owner decision because no stable YAML v4 release satisfies the approved envelope.

Assess actual module provenance, API maturity, replacement/reversal cost, and whether the code
unnecessarily couples M1 to prerelease behavior. Do not change the dependency.

## Severity and disposition

- `P0 blocker`: data loss, authority violation, false evidence, scope violation, unsafe persistence,
  or core semantic failure.
- `P1 material`: charter mismatch or maintainability/reversal risk that must be corrected before
  M1 acceptance.
- `P2 advisory`: bounded improvement not required for M1 acceptance.

H27 reports both permitted repair rounds consumed. A blocking finding therefore returns `bounce`
or `escalate`; the reviewer must not repair it or invent additional budget.

## Return

Return `spectacular.handoff-return.v2` with exact reviewed commits/trees, verified inputs, read set,
commands and outputs, ranked findings with file/line evidence, scenario/evidence audit, dependency
disposition, scope/authority verdict, assumptions, limitations, and one recommendation:

- `accept`: M1 is eligible for central assessment and owner disposition;
- `bounce`: exact blocking corrections are required, with no repair authority implied;
- `escalate`: a named owner/authority decision is required.

Exactly one next action. Do not accept M1 or authorize M2.
