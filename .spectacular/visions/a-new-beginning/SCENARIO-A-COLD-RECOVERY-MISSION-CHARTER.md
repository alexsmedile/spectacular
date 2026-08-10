---
type: mission-charter
mission: scenario-a-cold-recovery
version: 1.0
status: active
activation_state: autopilot-active
decision_session: scenario-a-owner-decision-and-implementation
source_authority: evidence/scenario-a-owner-dispositions.md@1.0
activated_by: owner
activated_at: 2026-08-10
central_disposition: pending
baseline_commit: 243ae3e376e6eb30e59dc28bb691a98dfc6a7b92
baseline_tree: c0922b35d84617ea82e30c63b6b70729dfa87110
branch: codex/feat/v2-scenario-a-cold-recovery
repair_budget: 2
next_action: implement-and-independently-review-scenario-a
upstream:
  - EXECUTABLE-REFACTOR-PROGRAM-CONTRACT.md@1.7
  - MVP-SCENARIO-CLI-SEQUENCING-DECISION.md@1.0
  - M1-SEMANTIC-SUBSTRATE-MISSION-CHARTER.md@1.2
  - SHARED-SCAFFOLD-CONTRACT.md@1.0
  - MISSION-PREPARATION-CONTRACT.md@1.0
  - PRODUCT-TRUTH-CONTRACT-MODEL.md@1.0
  - WORK-UNIT-LIFECYCLE-CONTRACT.md@1.0
  - EXECUTION-AUTHORITY-CONTRACT.md@1.0
  - EVIDENCE-CLOSURE-CONTINUITY-CONTRACT.md@1.0
  - RETRIEVAL-AND-EARNED-WORKSPACE-CONTRACT.md@1.0
  - PUBLIC-LANGUAGE-AND-INTERFACE-CONTRACT.md@1.0
  - SUCCESS-EVIDENCE-CONSTITUTION.md@1.0
---

# Scenario A Mission Charter — Cold recovery

## Outcome

Implement the smallest complete cold-recovery slice: a person or agent with no chat history can
enter an explicitly marked v2 workspace, locate authoritative project and Mission context, inspect
every consequential pointer, and receive either one mechanically justified continuation or the
exact owner gate preventing it. The slice is read-only and uses the accepted M1 domain, workspace,
and index primitives without claiming M1 already supplied CLI or lifecycle semantics.

## Understanding

### How it works now

M1 provides canonical Proposal/Mission Markdown parsing, semantic normalization, UUIDv7 identity,
typed Proposal-to-Mission references, SHA-256 fingerprints, exact UUID/path lookup, deterministic
index order, and relationship refusals. It has no command binary, registry, workspace discovery,
project/Mission projection, Gap/Run/Checkpoint/Evidence/Decision read model, safe-continuation
calculus, lifecycle engine, provider behavior, cache, or migration compatibility.

### What changes

Scenario A adds a read-only CLI and one internal record/projection substrate. It discovers only an
explicit v2 workspace; reads canonical anchor and typed recovery records; builds a disposable index;
renders source-backed human and JSON views; validates scopes; resolves exact pointers; and derives a
continuation only from a fresh, conflict-free, Decision-authorized fingerprint chain. All other
cases return a deterministic refusal or an exact owner gate.

### What stays the same

M1 identity, canonical parsing, fingerprints, exact lookup, relationship refusal, and atomic-write
behavior remain authoritative and unchanged in meaning. Canonical Markdown remains complete without
projections. V1 code, commands, skills, tests, schemas, aliases, and workspaces remain untouched and
outside v2 discovery. Product judgment, lifecycle mutation, assessment, reconciliation, provider
effects, caching, migration, release, and Scenario B/C behavior remain absent.

## Exact public contract

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

No public `record`, generic `status`, `inspect`, `new`, `advance`, `doctor`, `assess`, `decide`,
`resolve`, compatibility alias, fallback, or hidden dual path is allowed. Exact record references
are canonical UUIDv7, typed `<Noun>:<UUIDv7>`, or full workspace-relative path. `project` is the
sole special scope. Bare names, slugs, fuzzy matches, cleaned path alternatives, and ambiguity
selection refuse.

## Authority and projection contract

- An authoritative project-anchor document owns only minimal project facts and canonical pointers.
- A Mission record owns its bounded outcome and explicit recovery relationships.
- Gap, Run, Checkpoint, Evidence, and Decision records own only their attributable stored claims.
- Cards, Fog, lists, generation metadata, counts, conflicts, owner gates, and continuation are live,
  disposable projections. Deleting all projections leaves the workspace complete.
- Every consequential record pointer exposes noun, exact ref, path, fingerprint, and canonical
  noun-first `show_command`. Non-record sources expose exact path and fingerprint.
- A projection exposes sources, explicit `current | stale | unknown` freshness, generation time,
  fingerprint-based generation basis, Gaps, conflicts, and omissions.
- No projection accepts, resolves, reconciles, renews freshness, or changes canonical material.

## Mechanical continuation calculus

A continuation is emitted only when all conditions hold:

1. exactly one relevant Mission/Run/latest-Checkpoint chain resolves;
2. every consequential source and explicit freshness assertion is current;
3. no blocking Gap, broken dependency, duplicate identity, type mismatch, or conflict exists;
4. one attributable Decision explicitly authorizes the proposed operation;
5. the Decision's expected Mission fingerprint equals the current normalized fingerprint; and
6. the Decision's target and operation resolve exactly to the selected Mission chain.

Stored `safe`, `next_action`, status, prose, or completion labels never satisfy these conditions.
If the chain is well-formed but any condition is absent or false, successful orientation emits one
exact owner gate naming the missing authority/freshness/conflict and its source. Malformed,
unreadable, missing, stale-required, conflicting, or ambiguous inputs refuse deterministically.

## Workspace, schema, and output pins

- Discovery ascends from cwd to the nearest `.spectacular/workspace.yaml` whose explicit format is
  v2, then scans only its declared record roots.
- Discovery order is canonical path then stable identity; filesystem enumeration order is irrelevant.
- The internal Go command registry is the sole authority for command, argument, dispatch, effect,
  help, and JSON-schema identifiers. Scenario A persists no generated interface or projection.
- JSON uses versioned `spectacular.<noun>.<operation>.v1` envelopes and
  `spectacular.refusal.v1`. Arrays and object construction are stable; tests inject the clock.
- Human output is a concise card over the same facts and never carries extra authority.
- JSON refusals are written to stdout; human refusals to stderr. Valid cards and owner gates exit
  0, usage exits 2, and workspace/input refusals exit 3.
- No new dependency is authorized unless existing pins and the Go standard library are proven
  insufficient; that proof would be a stop requiring a charter delta.

## Discovery and source refusal rules

Refuse unsupported or absent v2 markers, undeclared roots, path escape, symlink escape, malformed
UTF-8/frontmatter/schema, duplicate paths/identities, missing or wrongly typed targets, invalid
typed references, unknown nouns, noun/ref mismatch, missing records, ambiguous identity, stale
required authority/evidence, conflicting authorities, and invalid continuation fingerprints.
Unknown optional freshness remains explicit and yields an owner gate when consequential; it is
never silently promoted to current.

## Read-only invariant

Every public Scenario A command performs no workspace write, normalization, lock, cache, generated
file, or temporary workspace output. Tests snapshot and compare the complete workspace path set,
file bytes, modes, and mtimes before and after success, owner-gate, and refusal paths. Filesystem
atime is outside the byte-nonmutation guarantee and must not be presented as canonical evidence.

## Owned paths

```text
.spectacular/visions/a-new-beginning/SCENARIO-A-COLD-RECOVERY-MISSION-CHARTER.md
.spectacular/visions/a-new-beginning/SCENARIO-A-COLD-RECOVERY-MISSION-CHARTER.md.sha256
.spectacular/visions/a-new-beginning/evidence/scenario-a-owner-dispositions.md
.spectacular/visions/a-new-beginning/evidence/scenario-a-owner-dispositions.md.sha256
.spectacular/visions/a-new-beginning/evidence/scenario-a-implementation-return.md
v2/.spectacular/
v2/cmd/spectacular/
v2/internal/command/
v2/internal/discovery/
v2/internal/domain/
v2/internal/index/
v2/internal/projection/
v2/internal/workspace/
v2/testdata/scenario-a/
```

Supporting Go tests may live beside owned packages. Existing M1 files may change only when the new
read-only record grammar requires a backward-compatible extension that preserves all M1 proofs.

## Prohibited paths and effects

- No v1/root CLI, Skill, hook, command, schema, test, collection, documentation, or migration edit.
- No program, orchestration, M1 closure, accepted constitutional-contract, central lifecycle, or
  Capability Contract mutation.
- No provider integration, persistent cache, compatibility behavior, alias, fallback, dual write,
  Scenario B transition, Scenario C reconciliation, release, push, PR, merge, publish, deployment,
  remote mutation, destructive cleanup, or secret/configuration effect.

## Dependencies and effects

The exact baseline is commit `243ae3e376e6eb30e59dc28bb691a98dfc6a7b92`, tree
`c0922b35d84617ea82e30c63b6b70729dfa87110`, on branch
`codex/feat/v2-scenario-a-cold-recovery`. The integrated M1 merge `7759461a…` is its direct parent.
The work is serialized across shared command, schema, registry, and test surfaces. Effects are
limited to owned local edits, disposable local test/build output, coherent local commits, and
temporary bounded Investigator/Builder/Reviewer/Verifier agents.

## Milestones

1. Lock the owner dispositions, Mission envelope, schemas, refusals, and tests.
2. Implement v2 discovery, typed read records, exact lookup, registry, projection, and CLI.
3. Prove adversarial cold recovery, source drill-down, continuation/gate behavior, nonmutation,
   deterministic projection, and real self-hosting.
4. Run the complete Go and project charter matrix, obtain independent review, repair within budget,
   and return to central orchestration.

## Proof contract

Required automated evidence includes deterministic discovery/order/lookup/JSON/projection; every
pointer's source drill-down; mechanically authorized continuation and exact owner gates; malformed,
stale, missing, conflicting, and ambiguous refusals; successful and refused byte nonmutation;
disposable projections; explicit absence of v1 discovery/import/change; and registry/CLI parity.

The cold-agent evaluation supplies only cwd and “recover and resume safely,” permits at most 12 CLI
invocations and 64 KiB aggregate JSON, caps project orientation JSON at 24 KiB, permits no manual
full-tree walk or authority error, and requires every consequential conclusion to cite a direct
source. The real `v2/.spectacular/` self-hosted run must reach the same justified continuation or
owner gate. Each command must remain under 500 ms p95 over 20 warm runs on the execution host.

Required final checks from `v2/`:

```text
gofmt -l .
GOFLAGS=-mod=readonly go mod verify
GOFLAGS=-mod=readonly go vet ./...
GOFLAGS=-mod=readonly go test -race ./...
GOFLAGS=-mod=readonly go test ./...
GOFLAGS=-mod=readonly go build ./...
```

Run the repository's applicable charter/scenario/model-harness tests without modifying v1. A fresh
reviewer who did not implement the work or author disputed evidence must inspect the final commit,
tree, primary outputs, contracts, fixtures, nonmutation proof, self-hosted result, and complete
check matrix. A separate verifier reruns primary checks against the reviewed tree.

## Design and slice gate

- Design Sufficiency: `sufficient` — authority, source, schema, failure, recovery, proof, and
  integration boundaries are explicit; no unresolved product or architecture choice remains.
- Slice Quality: `coherent` — the slice produces one independently reviewable cold-recovery outcome
  and leaves a complete read-only v2 surface if later scenarios are cancelled.

These verdicts are revalidated after investigation and before implementation begins. A new material
finding returns the affected path to preparation rather than being normalized silently.

## Autopilot authority and repair

Autopilot is explicitly active within this charter. It may inspect, edit owned paths, run checks,
build, use disposable fixtures, commit coherently, and dispatch the named temporary roles. Builders
and reviewers may converge on bounded findings without another owner gate while the envelope is
unchanged.

The repair budget is two focused rounds. Each requires a new hypothesis, new evidence, or materially
narrower correction, followed by the narrowest relevant check and then the full required matrix.
The recovery point is the verified baseline plus the last green coherent local commit. Preserve
failed evidence; never reset, stash, overwrite unrelated changes, or broaden the charter.

## Stop conditions

Stop only for a new product/architecture choice; material scope, public-contract, authority,
baseline, policy, or risk delta; need for a new dependency; forbidden or irreversible effect;
unresolved required check after bounded diagnosis; unexpected v1/central-lifecycle impact;
review-blocking finding outside the envelope; or exhausted repair budget. Routine implementation
choices, red tests, and bounded review fixes are not owner gates.

## Cold-resume packet

A cold executor reads this charter, its owner-decision receipt, the v1.7 program, exact Git
baseline/final commits, and the latest implementation return. It revalidates branch, dirty state,
upstream sidecars, repair budget, and required checks. The sole terminal next action is return a
`spectacular.handoff-return.v2` packet to central orchestration for `accept | bounce | escalate`.
The executor may not claim Scenario A accepted, resolve it, reconcile current truth, or authorize
Scenario B.
