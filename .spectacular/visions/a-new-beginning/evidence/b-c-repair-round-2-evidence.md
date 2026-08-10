---
type: implementation-repair-evidence
mission: b-c-governed-loop
version: 1.0
status: ready-for-new-independent-review
recorded_at: 2026-08-10
recorded_by: primary-implementation-owner
reviewed_head: 7ff40764b7f45db97c29550074ff2e64a724011f
reviewed_tree: 3ad0a405dfb14191239334e8335ca37a62e7b6a2
repair_commit: 9333cb4a3045a5da915c21417b1c0dae421c6004
repair_tree: f10c986d3c1977e0f9d6c94a2be3c0189e954640
branch: codex/feat/v2-b-c-governed-loop
repair_round: 2
repair_budget: 2
repair_budget_exhausted: true
recovery_point: 7ff40764b7f45db97c29550074ff2e64a724011f
review_authority: central-orchestration
next_action: central-dispatch-new-fresh-read-only-independent-reviewer
---

# B+C repair round 2 evidence

## Disposition and hypothesis

Central orchestration accepted the independent review's bounce recommendation and authorized the
second and final integrated repair round. The unifying hypothesis was confirmed: semantic
invariants existed in stored fields but were not mechanically revalidated at every filesystem,
authority, Evidence, replay, repair, and closure boundary.

The repair is limited to v2 governed-loop implementation, adversarial tests, this evidence, and the
active charter's versioned status. It performs no provider, migration, release, push, PR, merge,
Skill Mission, or v1 effect and does not claim central acceptance.

## Affected claims and before/after Evidence

| Claim | Before Evidence at reviewed head | After Evidence at repair commit |
|---|---|---|
| Transaction containment and recovery | symlinked parent escaped; crafted journal could target outside; preparation failure stranded artifacts | canonical-parent containment and artifact identity are revalidated before every effect; real-filesystem escape, crafted-journal, cleanup, and interruption regressions pass |
| Closure authority and ordering | a pending Objective, fake reconciliation string, and archive-chosen continuation could reach terminal state | exact receipt provenance/cardinality and resolution authority are required; completed closure requires satisfied Objectives; archive only preserves the Decision-derived continuation |
| Evidence sufficiency | stale support plus omitted contrary canonical Evidence could yield ready assessment | freshness is rechecked at assessment and consequential boundaries; all canonical same-claim Evidence is considered; stale/conflicting/empty-check Evidence blocks only affected material claims |
| Reconciliation replay and atomic public set | idempotency replay did not bind the full set and multi-Contract application lacked a public path | replay digest binds canonical ordered Contract/Proposal/Decision/fingerprint/cardinality semantics; changed input refuses; registry-derived `contract reconcile-set --input` proves public all-or-nothing application |
| Authority and Handoff | role, operation, expiry, and target checks omitted conditions/effects/scope and later Mission expiry/subset dimensions | typed authority validates scope, effects, known conditions, expiry, and target; Handoff validates claims, stops, budget, authoritative inputs, return actions, and same-Mission accountability |
| Governed repair | opaque string attempts and unenforced budget could not prove hypothesis-changing work | Mission/Assessment-owned typed attempts bind claims, hypotheses/evidence/action, actor, before/after Evidence, checks/result, budget consumption, and recovery point; unchanged and over-budget attempts refuse |

## Adversarial regression disposition

- Parent symlink escape: reproduced, then refused without outside write.
- Crafted journal overwrite: reproduced with a journal whose filename matched its key, then refused
  because target/artifact identity and canonical containment do not match.
- Pre-journal cleanup: reproduced stranded candidate/backup artifacts, then proved cleanup preserves
  the original and removes both artifacts.
- Pending Objective + fake reconciliation + archive continuation: reproduced, then each boundary
  independently refuses; only exact Decision-authorized terminal continuity survives.
- Stale support + omitted contrary Evidence: reproduced, then same-claim canonical discovery blocks
  readiness while unrelated claim Evidence remains intact.
- Replay counterexamples: nonexistent Proposal, nonexistent Decision, unrelated fingerprint, changed
  cardinality, or changed set under the same key refuse as non-identical replay.
- Expired envelope and overbroad Handoff return: later transitions/validation refuse.
- Multi-Contract partial failure: set validation and transaction rollback preserve every prior
  Contract; the public two-Contract provider-neutral path succeeds atomically.

## Typed implementation choices

- Transaction journals store workspace-relative paths and deterministic key-derived artifact names;
  effects resolve the canonical workspace root and nearest real parent immediately before mutation.
- Reconciliation sets use one strict JSON input with canonical Contract ordering under the noun-first
  registry operation `contract reconcile-set`.
- Authority, Evidence, reconciliation receipt, terminal packet, Objective state, Mission/Handoff
  envelopes, and repair attempts have typed parsers at consequential gates. Malformed known fields
  no longer degrade to empty authority.
- Repair attempts remain Assessment-owned and consume the Mission's numeric repair budget; no
  generic repair subsystem was introduced.

## Checks and result

All v2 and compatibility checks exited successfully after final repair commit
`9333cb4a3045a5da915c21417b1c0dae421c6004`, except for the disclosed frozen-v1 baseline failure
below:

```text
cd v2
gofmt -w internal
GOCACHE=/tmp/spectacular-go-cache GOMODCACHE=/tmp/spectacular-go-mod GOFLAGS=-mod=readonly go mod verify
GOCACHE=/tmp/spectacular-go-cache GOMODCACHE=/tmp/spectacular-go-mod GOFLAGS=-mod=readonly go vet ./...
GOCACHE=/tmp/spectacular-go-cache GOMODCACHE=/tmp/spectacular-go-mod GOFLAGS=-mod=readonly go test ./...
GOCACHE=/tmp/spectacular-go-cache GOMODCACHE=/tmp/spectacular-go-mod GOFLAGS=-mod=readonly go test -race ./...
GOCACHE=/tmp/spectacular-go-cache GOMODCACHE=/tmp/spectacular-go-mod GOFLAGS=-mod=readonly go build -o /tmp/spectacular-bc-repair-build/spectacular ./cmd/spectacular

cd ..
bash -n cli/spectacular cli/install.sh scripts/hooks/pre-commit
scripts/hooks/pre-commit --check
bash tests/run.sh
git diff --check
```

Focused tests cover the independent counterexamples, provider-neutral B+C cold resume, Scenario A
golden/nonmutation compatibility, transaction interruption/rollback, public atomic reconciliation,
and exact refusal classes. Full v2, race, and build checks pass. The version guard reports 1.37.3
consistently. The frozen-v1 suite reports 30 test files passing and one pre-existing assertion
failure in `tests/cli/wayfinding-sequencer.test.sh`: `spike outranks research at equal priority`.
The unchanged v1 discovery writer stores `kind:` while the sequencer reads `type:`, so both nodes
receive the fallback rank and research wins by iteration order. The required baseline and current
`cli/spectacular` bytes have the same SHA-256
`468e3702fd8a1eca48c166c959277201e0ab88429ab96a871ab76b329b8ea8dc`; the B+C scope changes
neither v1 code nor this test. Repair is prohibited by this Mission and remains a disclosed
out-of-scope baseline gap.

## Repair accounting and recovery

This round consumed repair budget unit 2 of 2. The hypothesis changed from treating stored semantic
fields as sufficient to requiring typed revalidation at every consequential consumer. Before
Evidence is the reproduced counterexample suite at reviewed head `7ff4076`; after Evidence is the
same committed suite passing at `9333cb4`. The non-destructive recovery point is the reviewed head;
prior Contract versions and originals remain preserved by the repaired transaction protocol.

No known in-charter implementation gap remains. The frozen-v1 sequencer mismatch above is an
unresolved required-check exception whose repair would violate the explicit no-v1 boundary. The
repair budget is exhausted, so that exception and any new finding belong to central disposition
rather than another self-directed repair.

## Owner gate

Central orchestration must dispatch a new fresh read-only independent reviewer against the final
documented head. This evidence does not authorize acceptance or any Skill Mission.
