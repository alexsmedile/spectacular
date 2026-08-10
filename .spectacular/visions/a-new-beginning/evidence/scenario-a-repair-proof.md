---
type: mission-repair-proof-evidence
mission: scenario-a-cold-recovery
version: 1.0
status: ready-for-independent-review
observed_at: 2026-08-10
returned_head: 816c34aec885747c5f4a7d40274a35b34f91becf
returned_tree: 5a79f91dc98949292d8cdef20005017d4a52b347
repair_commit: 151b255e57929949f91c336413b528a8b714fed8
repair_tree: 0e1c2ab7ab4c8200d43731a8282420e934bce274
branch: codex/feat/v2-scenario-a-cold-recovery
repair_rounds_used: 1
next_action: fresh-read-only-independent-review
---

# Scenario A repair proof

## Result

The one authorized repair round closes all three bounced P1 findings without changing the ten-command
Scenario A surface, dependencies, exit classes, effects, or later-scenario boundary.

- Matched command argument shapes now live in the command registry and are validated before workspace
  discovery. Malformed known commands return usage exit `2` outside a workspace.
- Discovery rejects a symlinked `.spectacular` metadata directory and a symlinked or non-regular
  `workspace.yaml` before loading manifest bytes. Existing declared-root and record containment
  checks remain intact.
- Pointer commands are resolved from the command registry. Proposal sources and authoritative
  `current_truth` entries expose exact noun, typed ref, path, and fingerprint, while omitting
  `show_command` because no Proposal operation exists. No `spectacular proposal show` is emitted.

Design Sufficiency remains `sufficient` and Slice Quality remains `coherent`. These were implementation
correctness defects, not new product or architecture choices.

## Exact implementation boundary

- Original program baseline: `243ae3e376e6eb30e59dc28bb691a98dfc6a7b92`
- Original baseline tree: `c0922b35d84617ea82e30c63b6b70729dfa87110`
- Clean bounced return head: `816c34aec885747c5f4a7d40274a35b34f91becf`
- Clean bounced return tree: `5a79f91dc98949292d8cdef20005017d4a52b347`
- Repair authorization commit: `8c86a948e39574abacd349ff0f5e1a56e6c1929d`
- Repair implementation commit: `151b255e57929949f91c336413b528a8b714fed8`
- Repair implementation tree: `0e1c2ab7ab4c8200d43731a8282420e934bce274`
- Branch: `codex/feat/v2-scenario-a-cold-recovery`
- New dependencies: none
- Remote/provider/lifecycle effects: none

The exact `816c34a..151b255` scope contains the repair-delta evidence plus five `v2/` implementation/
test files. The exact original-baseline-to-repair diff contains only the Scenario A charter/evidence
scope and `v2/`; the out-of-scope path query returned no paths. No v1/root CLI, Skill, hook, schema,
test, documentation, or command file changed.

## Focused regression evidence

Tests first reproduced all three failures against the bounced head:

- malformed known commands outside a workspace returned exit `3` / `workspace_not_found`;
- symlinked marker and metadata directory fixtures progressed to manifest/root handling instead of
  returning `path_escape`;
- a Proposal in `current_truth` emitted `spectacular proposal show`.

After the repair, the focused command/discovery suite passes:

```text
go test ./internal/command -run TestKnown...|TestProposal...|TestAllRead...|TestOwnerGate...|TestAdversarial...
ok github.com/alexsmedile/spectacular/v2/internal/command
go test ./internal/discovery -run TestOpenRefusesAuthority...|TestOpenRefusesSymlinkEscape
ok github.com/alexsmedile/spectacular/v2/internal/discovery
```

The real-filesystem marker tests cover a `workspace.yaml` symlink to an external valid manifest and
an entire `.spectacular` symlink to an external directory. Both deterministically refuse with
`path_escape`. Existing external record-symlink proof continues to pass. Existing success, owner-gate,
refusal, adversarial, golden, deterministic lookup/order/projection, and byte/path/mode/mtime
nonmutation tests all pass in the complete suite.

## Required check matrix

| Check | Result |
|---|---|
| `gofmt -l .` from `v2/` | pass; empty output |
| `GOFLAGS=-mod=readonly go mod verify` | pass; all modules verified |
| `GOFLAGS=-mod=readonly go vet ./...` | pass |
| `GOFLAGS=-mod=readonly go test -count=1 ./...` | pass |
| `GOFLAGS=-mod=readonly go test -race -count=1 ./...` | pass |
| `GOFLAGS=-mod=readonly go build ./...` | pass |
| `bash -n cli/spectacular cli/install.sh scripts/hooks/pre-commit` | pass |
| `scripts/hooks/pre-commit --check` | pass; guarded versions consistent |
| `bash tests/run.sh` terminal rerun | pass; 31 test files, 0 failed |
| `git diff --check` | pass |
| exact branch/head/tree and scope-diff queries | pass; branch and paths match the authorized envelope |

Disposable Go caches and the proof binary lived under `/private/tmp`; no generated artifact was
written into the workspace.

## Deterministic cold and self-hosted recovery

A freshly built repaired binary reran the existing deterministic fixture sequence directly, without
a separate verifier. The ten accepted operations completed in ten invocations, all exit `0`, and
emitted 20,983 aggregate JSON bytes. Project orientation remained 3,029 bytes. The output exposed the
same exact Mission → Run → Checkpoint → Evidence → Decision chain and fingerprint-bound `resume`
continuation, and contained no unsupported Proposal command.

The seven-call real `v2/.spectacular/` self-host sequence completed with 16,661 aggregate JSON bytes,
3,029-byte orientation, and the same continuation. Twenty warm self-host orientation processes
measured `min=2.592 ms`, `p95=3.020 ms`, and `max=3.201 ms`. These pass the accepted limits of 12 calls,
64 KiB aggregate JSON, 24 KiB orientation, and 500 ms p95.

The repair does not change the cold command surface or continuation calculus, so the prior independent
cold-agent behavior remains applicable; this run directly reconfirmed its deterministic outputs.

## B+C compatibility finding

Scenario A's envelope continues to satisfy the accepted non-expanding B+C join without redesign:
versioned success/refusal envelopes, stable references, exact source/fingerprint and generation basis,
separate data from any future warning channel, exits `0/2/3`, and a registry that owns dispatch,
argument shapes, help text, effects, schemas, and supported pointer commands. The repaired Proposal
representation strengthens this boundary by omitting false operations while retaining exact source
metadata. No conflict exists. Optional B+C-only mutation/refusal/idempotency fields remain future
additive fields because Scenario A is read-only and B/C remain unauthorized.

## Assumptions and limits

The standard-library pre-read `Lstat` checks prove the non-concurrent filesystem shape exercised by
Scenario A; no concurrent hostile path replacement model is introduced by this delta. Performance is
host-local to the declared fixture and execution host. No acceptance, current-truth reconciliation,
Scenario B authorization, provider effect, dependency, remote operation, or release claim is made.

## Next action

One fresh independent reviewer must inspect the exact proof commit/tree, primary outputs, focused
regressions, complete check matrix, and scope boundary without editing. Then the primary agent records
the review and returns `spectacular.handoff-return.v2` to central orchestration.
