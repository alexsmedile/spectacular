---
type: post-budget-correction-receipt
mission: b-c-governed-loop
version: 1.0
status: awaiting-targeted-independent-verification
recorded_at: 2026-08-10
recorded_by: primary-implementation-owner
authorized_by: central-orchestration-under-owner-autonomous-continuation
authority_kind: narrow-existing-invariant-correction
correction_baseline_commit: 2af2e5f0c1594bc106c8e0ce57ea44a69744e7bd
correction_baseline_tree: 4b30c7d00ec3f71fa585970aadc6436c2de01aec
correction_commit: 3164496446675a683cf5860d56e09647189dd499
correction_tree: a414c04487919184c5c7fa11ebb56dcd8bccfd8d
branch: codex/feat/v2-b-c-governed-loop
repair_budget: 2
repair_budget_exhausted: true
repair_budget_reset: false
recovery_point: 2af2e5f0c1594bc106c8e0ce57ea44a69744e7bd
next_action: targeted-fresh-independent-verification-of-this-correction
---

# B+C post-budget transaction correction receipt

## Finding and exception basis

Independent re-review found one remaining direct P1 in
`v2/internal/governance/transaction.go`: final installation validated candidate and target paths,
then called ordinary `os.Rename`. An adversary could replace the validated target parent with a
symlink before pathname resolution, redirecting the candidate outside the canonical workspace.

Central orchestration authorized only the smallest correction of this existing charter invariant
under the owner's instruction to continue autonomously. This is not a third general repair round,
does not reopen B+C semantics, and grants no v1, provider, Skill, release, migration, push, PR, or
merge authority.

## Hypothesis and changed paths

Hypothesis: pre-validation can produce precise refusals but cannot safely authorize a later
pathname-based effect. A workspace-root file descriptor or platform handle must anchor both path
resolution and the consequential operation. The same mechanism applies to preparation writes,
final rename, rollback restoration, recovery reads/removals, and cleanup removals.

Changed paths:

- `v2/internal/governance/transaction.go`
- `v2/internal/governance/governance_test.go`
- this receipt, its SHA-256 sidecar, and the versioned Mission charter/snapshot metadata

The implementation uses Go 1.26 `os.Root`, which provides race-resistant rooted operations on
macOS/Linux and a native Windows handle boundary. No shell behavior, Unix-only syscall package, or
silent Windows assumption was introduced.

## Direct before/after evidence

Before correction, `TestTransactionInstallResistsValidatedParentSwap` deterministically moved the
validated `records/` parent aside, replaced it with a symlink to a disposable outside directory,
and observed the old `os.Rename` install `candidate` outside the workspace. The test failed with:

```text
parent swap unexpectedly installed outside the rooted workspace
```

After correction, the same hook runs between pre-validation and install. Rooted rename refuses
with `path_escape`; the outside `contract.md` remains absent and the displaced original remains
byte-identical to `original`. All transaction reads and effects now share the pinned root boundary;
only diagnostic `Lstat` calls remain pathname-based, and safety does not depend on them.

## Checks and result

All authorized checks passed at correction commit `3164496446675a683cf5860d56e09647189dd499`:

```text
go test ./internal/governance -run '<transaction cluster>' -count=50
go test -race ./internal/governance -run '<transaction cluster>' -count=10
go mod verify
go vet ./...
go test ./...
go test -race ./...
go build -o /tmp/spectacular-post-budget-build/spectacular ./cmd/spectacular
GOOS=windows GOARCH=amd64 go test -c ./internal/governance
go test ./internal/command -run '<Scenario A focused set>' -count=1
bash -n cli/spectacular cli/install.sh scripts/hooks/pre-commit
scripts/hooks/pre-commit --check
git diff --check
```

The focused transaction cluster passed 50 ordinary and 10 race-enabled repetitions. Full Go,
Scenario A, Bash syntax, version consistency, native build, and Windows cross-compile checks pass.
The known frozen-v1 `kind`/`type` sequencer failure remains unchanged and out of scope; no v1 path
changed.

## Result and recovery

Result: the reported parent-swap escape is closed, and equivalent consequential transaction
effects share the same rooted mechanism. Journal identity, original-byte preservation, atomic
multi-file behavior, cleanup, rollback, recovery, public CLI behavior, and Scenario A remain
intact.

The non-destructive recovery point is correction baseline
`2af2e5f0c1594bc106c8e0ce57ea44a69744e7bd`. General repair budget remains exhausted at 2/2.
This receipt does not accept B+C.
