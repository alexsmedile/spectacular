---
type: implementation-evidence
mission: b-c-governed-loop
status: ready-for-independent-review
recorded_at: 2026-08-10
recorded_by: primary-implementation-owner
baseline_commit: 08eca1008d8ae16519f828d12112621ab6a506e9
baseline_tree: 5cb75fecc4251ff781e3489ff344495d4d401ed9
implementation_commit: e247f0de0686a30d904a02e49e55bd1d88061eb8
implementation_tree: 26f7009e5232522637c932d1e1b896a0c3670811
branch: codex/feat/v2-b-c-governed-loop
review_authority: central-orchestration
---

# B+C primary implementation evidence

## Outcome

The isolated B+C implementation completes the provider-neutral governed loop without changing v1,
granting the CLI assessment/owner authority, or performing provider effects. It stops at the fresh
independent-review boundary and does not claim central acceptance.

Proposal and Capability Contract remain distinct records and authorities. A Proposal persists the
accepted, exact base-bound delta plus a disposable complete candidate. Reconciliation consumes the
current Contract, accepted Proposal, ready Assessment, explicit owner Decision, and expected
fingerprint to create the next complete Contract version and an inspectable Evidence receipt.

## Material Type-2 choices

- Confirmed create inputs use strict single-value JSON files; unknown fields and trailing values
  refuse. Guided Skill work retains interpretation, slicing, assessment, and Decision authorship.
- All governed identities are UUIDv7. Creates are atomic and identity-idempotent across later
  processes; transitions, reconciliation, and archival also return exact replays.
- Multi-file and multi-Contract changes use a recoverable synced-file transaction journal under the
  disposable workspace metadata. Recovery restores every original before a later mutation.
- A deterministic UUIDv7 Evidence record is the durable atomic-reconciliation receipt. Prior
  Contract bytes remain under versioned Contract history.
- Archival is logical and atomic with Anchor truth replacement: the current Contract becomes
  authoritative truth while `last_closed_mission` preserves terminal continuity for cold recovery.
- Scenario A exit classes remain `0` success, `2` usage, and `3` refusal. Refusals add expected/
  actual state where applicable, zero-mutation status, and a mechanical recovery action.

## Claim mapping

| Claim | Primary evidence |
|---|---|
| B1 Proposal is rigorous but never current Contract truth | `TestProviderNeutralGovernedLoopAndSecondColdResume`; Proposal view marks candidate non-authoritative; Scenario B+C fixture Anchor names only Contract truth |
| B2 Creation and lifecycle are exact, authorized, atomic, and idempotent | governance service tests for stale base, expired authority, later-process replay, transition replay, and refusal-before-write |
| B3 Handoff is bounded, immutable, runtime-neutral, and non-accountability-transferring | envelope subset refusal, validation proof/non-proof set, immutable return, supersession refusal, replacement-runtime return |
| C1 Evidence is claim-mapped and assessment remains explicit guided output | direct claim/check mapping, executor-authored independent-review requirement, CLI operation `assessment record` rather than `assess` |
| C2 Decision-bound reconciliation is all-or-nothing and inspectable | single reconciliation, shared-receipt multi-Contract reconciliation, stale set zero-write digest check, injected interruption rollback, version history |
| C3 Resolution and archival are ordered and terminally recoverable | premature resolution refusal, ready Assessment requirement, reconciliation/zero-delta gate, abandoned/superseded cases, archive replay |
| C4 A second cold actor can recover without chat/runtime state | fixture loop reopens discovery and projection from cwd only; recovers current Contract, archived Mission, Proposal, Assessment, Decisions, Evidence receipt, and one terminal continuation |
| A1 Scenario A behavior remains compatible | original command tests, exact golden digest, deterministic lookup/refusal and byte-nonmutation tests all pass |

## Verification

All commands exited successfully at implementation commit `e247f0d`:

```text
cd v2
gofmt -w internal
GOCACHE=/tmp/spectacular-bc-gocache GOFLAGS=-mod=readonly go test ./...
GOCACHE=/tmp/spectacular-bc-gocache GOFLAGS=-mod=readonly go mod verify
GOCACHE=/tmp/spectacular-bc-gocache GOFLAGS=-mod=readonly go vet ./...
GOCACHE=/tmp/spectacular-bc-gocache GOFLAGS=-mod=readonly go test -race ./...
GOCACHE=/tmp/spectacular-bc-gocache GOFLAGS=-mod=readonly go build -o /tmp/spectacular-bc-build/spectacular ./cmd/spectacular

cd ..
bash -n cli/spectacular cli/install.sh scripts/hooks/pre-commit
scripts/hooks/pre-commit --check
bash tests/run.sh
git diff --check
```

Focused proof includes Proposal base/collision/idempotency, authority expiry, Handoff subset and
supersession, executor-dependent Evidence blocking, assessment and closure refusal, single/multi
Contract reconciliation, transaction interruption rollback, zero-delta abandoned/superseded
outcomes, terminal archival, CLI strict-input/readback, and second cold resume. The full v1 suite
ran all 31 test files with no failures; the version guard reported every guarded value consistent.

## Adversarial review

- Stale Proposal base: refuses before mutation with expected/actual fingerprints and recovery.
- Conflicting Decision authority: superseded Decisions/Handoffs refuse; no latest-wins behavior.
- Interrupted Run or expired envelope: expiry refuses; transaction interruption restores all files.
- Incomplete/executor Evidence: missing checks, conflicts, unknown classification, or unreviewed
  executor authorship blocks only the affected claim from owner readiness.
- Multi-Contract partial failure: full set validates before install; stale member leaves both current
  Contract file digests unchanged; injected install interruption rolls both back.
- Abandoned zero-delta Mission: explicit `no-contract-delta` owner disposition permits resolution
  and archival only after ready Assessment and terminal continuity; superseded is equivalent.
- Cold actor with no chat: current Contract and historical closure chain recover from canonical
  records with one justified continuation.
- Host/runtime replacement: host pointer stays non-authoritative and a replacement runtime can
  validate/return the bounded Handoff without acquiring Mission accountability.

## Repair accounting and gaps

One hypothesis-changing repair round was used: the first integrated regression showed that the
initial governance-Decision discriminator softened Scenario A's malformed-resume refusal. The
discriminator was changed to legacy Mission linkage, restoring strict Scenario A behavior while
keeping governance Decisions additive. Subsequent focused, full, race, build, and v1 checks pass.

No known implementation gap remains within the B+C charter. Skill packaging, release mechanics,
providers, migration, push/PR/merge, and central acceptance remain intentionally out of scope.

## Owner gate

Central orchestration must dispatch one fresh read-only independent reviewer against the final head.
