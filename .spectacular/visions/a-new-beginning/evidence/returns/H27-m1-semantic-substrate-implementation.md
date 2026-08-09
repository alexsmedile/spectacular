---
type: handoff-return
schema_version: spectacular.handoff-return.v2
handoff_id: H27
mission: M1
status: implementation-complete
central_disposition: proceed-to-independent-review
acceptance: pending-H28
branch: codex/feat/v2-semantic-substrate
dispatch_commit: f628e8ca895c47b0e7ee7142b98b89dfb16a1e0f
dispatch_tree: 08172f732cc23c2dd5857a8bc1c18de0429a88e4
final_commit: 489bd6008e1720e4b0310b999a0bac02c62df6dc
final_tree: 911adb7081ea70b10d29539cc76a343c6d58b0cf
commit_count: 3
worktree_clean: true
date: 2026-08-09
---

# H27 reviewed return — M1 implementation

## Central assessment

H27 is complete enough to enter independent review. Central orchestration reproduced the final
branch/tree, three-commit history, owned-path-only diff, clean worktree, `go mod verify`, and
`go test ./...`. This is not M1 acceptance and does not authorize M2.

## Implementation result

- Added the isolated `github.com/alexsmedile/spectacular/v2` Go module.
- Implemented Domain identity, Proposal/Mission semantics, typed references, deterministic
  refusal ordering, canonical UUIDv7 generation and validation.
- Implemented Markdown/frontmatter parsing, unknown-property and body preservation, semantic
  canonicalization, fingerprints, and atomic replacement.
- Implemented deterministic identity/path lookup and post-discovery relationship validation.
- Added positive and negative fixtures plus focused and integrated tests for every charter scenario.

Final self-evidence passed: formatting, module verification, vet, focused tests, race tests, full
tests, and build. The exact scenario-to-test map and command results are preserved in the H27 task
return.

## Dependencies and named review question

- Go `1.26.5` on `darwin/arm64`
- `github.com/google/uuid v1.6.0`
- `go.yaml.in/yaml/v4 v4.0.0-rc.6`

The YAML dependency is an official v4 release candidate. H27's instruction requested stable pinned
versions, so H28 must determine whether this is a blocking envelope deviation, an acceptable
activation-time Type-2 pin, or requires owner escalation. Central orchestration does not silently
normalize it.

## Repair and authority

H27 used the two permitted focused repair rounds: it removed duplicate canonical-order authority
from Domain, then replaced map-based validation with deterministic grammar ordering. The repair
budget is consumed. No push, PR, merge, tag, release, provider change, v1/root mutation, or M2
activation occurred.

`next_action: H28 independent review; do not accept M1`
