---
type: mission-proof-evidence
mission: scenario-a-cold-recovery
version: 1.0
status: ready-for-independent-review
observed_at: 2026-08-10
baseline_commit: 243ae3e376e6eb30e59dc28bb691a98dfc6a7b92
implementation_commit: 231a3775bfe498135dfba60fdaa2a82ad9ff57f0
implementation_tree: a8b48ebd0bb95154c53fec75efb5a0d4eec49e57
branch: codex/feat/v2-scenario-a-cold-recovery
repair_rounds_used: 2
next_action: independent-review-final-implementation-and-proof
---

# Scenario A primary proof

## Result

The committed v2-only implementation provides the ten accepted noun-first read operations through
one Go registry, discovers only explicitly marked v2 workspaces, separates authoritative project
Anchors from disposable projections, resolves every consequential pointer exactly, and emits either
one fingerprint-bound `resume` continuation or an exact owner gate. It performs no lifecycle or
provider mutation and contains no v1 compatibility path.

## Exact implementation boundary

- Baseline: `243ae3e376e6eb30e59dc28bb691a98dfc6a7b92`
- Baseline tree: `c0922b35d84617ea82e30c63b6b70729dfa87110`
- Implementation commit: `231a3775bfe498135dfba60fdaa2a82ad9ff57f0`
- Implementation tree: `a8b48ebd0bb95154c53fec75efb5a0d4eec49e57`
- Branch: `codex/feat/v2-scenario-a-cold-recovery`
- Changed product paths: `v2/` only
- New dependencies: none
- Remote/provider effects: none

## Mechanical and adversarial evidence

Go tests cover all ten registry entries and schemas, dispatch parity, exact UUIDv7/typed/path lookup,
deterministic fixed-clock JSON and its exact digest, source/show-command drill-down, derived
freshness with expected/actual source fingerprints, authorized continuation, missing-authority
owner gate, blocking/multiple-authority conflict, malformed/missing/changing sources, stale detail,
duplicate identity, noun/scope/usage refusals, path and symlink escape, deterministic discovery, and
success/gate/refusal byte-mode-mtime nonmutation.

The locked Anchor JSON golden digest is
`9c2d6d213bd8e66c98643a0a0ea1009b93ecfa9dc1ef269445399284396da937` for the
fixed `2026-08-10T10:00:00Z` clock and committed adversarial fixture.

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
| `bash tests/run.sh` final rerun | pass; 31 test files, 0 failed |

The first full v1 suite run had one isolated `wayfinding-sequencer` assertion failure; its immediate
targeted rerun passed 15/15, and the complete final rerun passed all 31 files. No v1 source or test
file changed.

## Self-hosted scenario

The built non-race binary ran against the committed real `v2/.spectacular/` workspace. Seven
source-drill-down calls (`anchor`, Mission list/show, Run, Checkpoint, Evidence, Decision) produced
16,517 aggregate JSON bytes; the project orientation card was 3,029 bytes. The card exposed a
current, conflict-free exact `resume` continuation for
`Run:0198a1a0-0000-7000-8000-000000000005`, authorized by
`Decision:0198a1a0-0000-7000-8000-000000000008` against the current Mission fingerprint.

Twenty warm project-orientation runs measured:

```text
runs=20 min_ms=5.380 median_ms=5.620 p95_ms=6.216 max_ms=6.788
```

This passes the accepted limits of 12 calls, 64 KiB aggregate JSON, 24 KiB project orientation,
and 500 ms p95.

## Fresh cold-agent evaluation

Verifier `/root/scenario_a_cold_verifier`, started without prior conversation context, used ten CLI
invocations and nine JSON responses totaling 16,056 bytes. It used no manual repository walk,
validated all eight records, followed the Mission → Run → Checkpoint → Evidence → Decision chain,
identified the visible non-blocking Gap, and reached the same authorized `resume` continuation.

Its initial v1-shaped guess `status --brief --json` returned the designed usage refusal (exit 2),
after which registry-derived help exposed the noun-first v2 interface. This was a command-discovery
miss, not an authority error or fallback: no unsupported command ran, no state changed, and no
unsafe continuation was inferred. The actor stayed within all call/context limits.

## Repair history

1. Repair round 1 replaced trusted stored freshness labels with dated, source-backed derived
   freshness; separated authoritative project Anchors from projection; fixed schema/refusal
   metadata and Anchor drill-down; and required attributable exact `resume` Decisions.
2. Repair round 2 bound freshness to expected/actual source fingerprints, made the registry's
   operation binding authoritative for dispatch, refused duplicate JSON flags and stale detail,
   and added the remaining deterministic adversarial/nonmutation/golden proofs.

Both rounds used a narrower hypothesis and reran targeted then complete checks. The declared repair
budget is exhausted; unresolved blocking independent-review findings must therefore bounce or
escalate rather than trigger another implementation repair.

## B+C join compatibility finding

Scenario A is compatible with the accepted non-expanding B+C join without redesign:

- success and refusal outputs are versioned typed envelopes;
- stable refs, canonical paths, source fingerprints, generation basis, and expected/actual
  freshness fingerprints are explicit;
- the registry owns command identity, arguments, operation binding, help enumeration, effects, and
  schema IDs;
- exit semantics remain exactly A's `0 / 2 / 3`;
- data contains no warning channel, so later warnings can remain separate; and
- the typed refusal envelope can gain optional affected-ref, expected/actual operation fingerprint,
  zero-mutation, and recovery/owner-gate fields additively.

No semantic conflict was found. The only deliberate provisional boundary is that Scenario A does
not yet populate B+C-only idempotency or mutation-recovery fields because it performs no mutation
and B/C are unauthorized.

## Facts, assumptions, and limits

Facts: all current source fingerprints are calculated from canonical record content or exact raw
source bytes; projections are in-memory only; every public operation is declared read-only; no
write/persistence function is reachable from command dispatch; and `git diff` contains no v1 path.

Assumptions: the accepted 500 ms threshold is host-local and was evaluated on the execution host;
future operation/result schemas may extend optional fields while preserving their versioned
envelope contract. No provider, external owner, migration, cache, Scenario B transition, Scenario C
reconciliation, or release claim is made.
