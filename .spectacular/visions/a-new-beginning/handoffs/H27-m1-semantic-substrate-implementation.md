---
type: implementation-mission-handoff
schema_version: spectacular.handoff.v2
handoff_id: H27
mission: M1
mode: isolated-implementation
authority: central-orchestration
status: authorized-for-dispatch
content_baseline_commit: b6bf59b95e049f57f75bc10a880ecd786e61fe8e
content_baseline_tree: 3f9124e88cfb4b315825322e4f75cf7e6e15c4d2
required_branch: codex/feat/v2-semantic-substrate
date: 2026-08-09
---

# H27 — M1 semantic substrate implementation

## Outcome

Implement the owner-approved M1 charter as the smallest working Spectacular v2 substrate. Return a
tested local branch ready for independent H28 review. Do not accept M1, activate M2, push, open a
PR, or change provider state.

## Binding inputs

- `M1-SEMANTIC-SUBSTRATE-MISSION-CHARTER.md@1.1` — SHA-256
  `a5340b0a63648585a117736d638a7ea0d4de58ae6110ac24e77bcc23babac98f`
- `EXECUTABLE-REFACTOR-PROGRAM-CONTRACT.md@1.5` — SHA-256
  `1e4027fe7b3eeb0263e5eee0ad946022e3fb09988fa2605caf2508ce8fd6a3c7`
- `SHARED-SCAFFOLD-CONTRACT.md@1.0` — SHA-256
  `698997f12972d0b5a186f4d8b8c35753a642cc0454e8ef60f15de81590435d36`
- `EVIDENCE-CLOSURE-CONTINUITY-CONTRACT.md@1.0` — SHA-256
  `7dd763b4fa1a919924e24105790382a51414a0a8ee0222178dee7c9224f11ca9`

Read `AGENTS.md`, `.spectacular/AGENTS.md`, the complete binding inputs, and H26's reviewed return.
Load no v1 implementation beyond repository rules and build-boundary verification.

## Launch and Git authority

Central will supply the exact dispatch commit and tree in the task prompt. Codex-managed worktrees
may start detached; detached launch is not a conflict when HEAD/tree exactly match the supplied
dispatch baseline and the required branch is absent. In that case create and switch to:

```text
codex/feat/v2-semantic-substrate
```

If already on that branch at the exact baseline, proceed. Stop on any other branch, changed HEAD or
tree, tracked dirt, branch collision, or ambiguous worktree ownership.

Allowed Git effects: create/switch the named local branch and make at most three coherent local
commits. Forbidden: push, PR, merge, tag, release, stash, reset, destructive cleanup, remote
mutation, or changes to another worktree.

## Owned paths

Modify only:

```text
v2/go.mod
v2/go.sum
v2/internal/domain/
v2/internal/workspace/
v2/internal/index/
v2/testdata/m1/
```

Tests live beside their owned packages. Do not touch root/v1 code, existing Skill/docs, canonical
planning contracts, generated interfaces, CI, Guardrails, adapters, caches, or migration surfaces.

## Implementation constraints

- Go module: `github.com/alexsmedile/spectacular/v2`.
- Use the locally available Go 1.26 toolchain and record the exact version in the return.
- Use the standard library except for a UUIDv7 dependency and `go.yaml.in/yaml/v4`; resolve and pin
  exact stable versions in `go.mod`/`go.sum`.
- The only universal record requirements are `type` and canonical UUIDv7 `id`.
- Implement Proposal and Mission grammars only. The M1 Mission relationship is
  `source: "Proposal:<UUIDv7>"`.
- Preserve unknown valid frontmatter values and opaque Markdown body without granting unknown
  fields semantic authority.
- Canonicalize semantic content exactly as the charter specifies; compute SHA-256 fingerprints
  rather than storing them.
- Use same-directory temporary-file replacement with original-file preservation on failure.
- Index deterministically by stable identity and full workspace-relative path; validate typed
  targets after discovery so filesystem order cannot affect results.
- Do not implement lifecycle transitions, authorization, assessment, reconciliation, CLI/Skill
  behavior, fuzzy lookup, Guardrails, providers, persistent caches, release, or migration.

This is one ordered, shared-interface Mission. Build inline and serialize milestones; do not fan out
parallel builders. Independent review is H28 and remains outside this implementation task.

## Milestones and checks

1. Establish `v2/go.mod`, dependencies, Domain types, identity, typed-reference parsing, static
   known-field validation, and refusal errors. Run focused Domain tests.
2. Implement frontmatter/body parsing, semantic canonicalization, fingerprints, and safe atomic
   persistence. Run focused Workspace tests.
3. Implement deterministic exact ID/path lookup and relationship validation. Run focused Index
   tests.
4. Add approved positive/negative fixtures, run the full evidence matrix, inspect the final diff,
   and create coherent checkpoint commits within the three-commit ceiling.

Required final checks from `v2/`:

```text
gofmt -l .
go mod verify
go vet ./...
go test ./internal/domain ./internal/workspace ./internal/index
go test -race ./internal/...
go test ./...
go build ./...
```

The return must map primary test names or exact evidence to every scenario in the Mission charter,
including failed-write preservation. `gofmt -l .` succeeds only with empty output. Race testing may
be marked unsupported only with direct host evidence.

## Repair, recovery, and stops

Initial implementation plus at most two focused repair rounds. Each repair must name a new
hypothesis or materially narrower correction and rerun the narrowest failed check before the full
suite. Recovery is the exact dispatch baseline plus the last green local checkpoint. Preserve
failed evidence and worktree state; never destructively reset.

Stop on any charter stop condition, any required path/dependency expansion, inability to preserve
unknown content or failed-write safety, M2/M3 leakage, exhausted repair budget, or required remote
effect. Ask no product questions unless a true stop condition exposes missing owner authority.

## Completion and return

Completion means the bounded implementation and all self-evidence pass on the final local head.
It does not mean M1 acceptance. Return `spectacular.handoff-return.v2` with:

- exact dispatch baseline and final commit/tree;
- branch, commits, and complete changed-file list;
- Go/toolchain and pinned dependencies;
- every required command and result;
- scenario-to-test/evidence map;
- repair history and recovery point;
- assumptions, limitations, conflicts, and scope deviations;
- confirmation of no push/PR/provider effect;
- exactly one next action: `central review, then dispatch independent H28 review; do not accept M1`.
