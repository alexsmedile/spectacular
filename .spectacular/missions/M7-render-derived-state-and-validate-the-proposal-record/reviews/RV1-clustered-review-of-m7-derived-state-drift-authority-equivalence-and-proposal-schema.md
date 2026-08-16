---
type: Review
id: 01a00c7b-43c0-7d36-94d6-206c32b17da5
title: Clustered review of M7 derived state, drift, authority, equivalence, and Proposal schema
status: passed
created: "2026-08-16T20:42:56Z"
claims:
    - claim: state-line
      verdict: pass
    - claim: drift-flags
      verdict: pass
    - claim: authority-table
      verdict: pass
    - claim: render-equivalence
      verdict: pass
    - claim: proposal-schema
      verdict: pass
findings: []
limitations: []
mission: M7
ref: RV1
reviewed:
    activation_fingerprint: sha256:59d80c9784f7b580ecc6c72eb1d3563f19ca2125b20e12ddd0a67b0580a4168e
    commit: 80f69e4e0b481566d52b0e9f8b113e3fc041c604
    tree: ab4738fde1ffe2ad960a709cd421af3b21221cc0
reviewer:
    actor: Claude Opus 5 (operator)
    evidence:
        - ad9e3267820cfe58b
        - a31fe026488ab8694
    implemented_reviewed_scope: true
    independence_basis: none; this is a clustered review recorded by the implementing operator, not an independent one
    operator: Claude Opus 5 (operator)
    relation_to_operator: same-actor
---
# Clustered review of M7

## What this record is, and is not

This is a **clustered** review recorded by the operator who implemented the work.
It is not independent, and the frontmatter says so: `actor` and `operator` are the
same, and `implemented_reviewed_scope` is `true`.

Two genuinely independent reviewers (task IDs above) audited M7 at commit
`b0a6c5d`, under the pre-amendment boundary. Their findings are summarized below
and were repaired. They have not seen the repairs or the amended boundary. The
Mission's review mode was downgraded from `independent` to `clustered` rather
than binding their names to work they never read.

## What the independent reviewers established at b0a6c5d

Twelve of fourteen mutations were caught. Reviewer 1 covered state-line,
drift-flags, and authority-table; reviewer 2 covered render-equivalence and
proposal-schema. Both restored the tree and verified a clean build.

Four findings were raised.

**Finding 1 — `show` equivalence untested.** The frozen boundary names `show`
first, but the equivalence test compared derived structures and never called the
human renderer. A mutation leaking a promoted Objective's file path into
`mission show` output passed the entire suite.

**Finding 2 — width threshold not tested at its boundary.** The test sampled
widths 100 and 12; the real threshold was 25. An off-by-one in the selection rule
was invisible.

**Finding 3 — `scope` named but not enforced.** `proposalRequiredFields` omitted
`scope`. P1 carries no `scope:` and validated with zero notices.

**Finding 4 — evidence age named but not implemented.** `drift.go` mentioned
evidence only in a comment. `Reviewer.Evidence` is a list of bare references with
no timestamps, so evidence age is not derivable from the record.

## What was done in response

Findings 1 and 2 were proof gaps and were closed with tests
(`801c26b`). `TestShowRendersIdenticallyAcrossInlineAndPromotedObjectives` drives
the real CLI across five human surfaces after promoting a root and a leaf.
`TestGraphSelectionIsExactAtTheWidthThreshold` derives the threshold at runtime
and asserts both sides of it.

Findings 3 and 4 were boundary defects, not implementation defects: both named
inputs the record cannot supply. The owner narrowed both criteria (`868f30d`)
rather than implementing against a shape that does not exist. The
`proposal-schema` proof requirement was also corrected — it described P5 and P6
as using the current `ref:` spelling, but all six Proposals use `human_ref:`.

## Verification performed for this review

Each repair was verified by replaying the reviewers' own surviving mutants:

| Mutation | Result |
|---|---|
| Leak promoted file into `renderHuman` output | caught, names the failing command |
| Graph width threshold `<= width+1` | caught |
| Graph width threshold `< width` | caught |

Additional probes: all six Proposals validate against the amended schema with
`ref-spelling-drift` reported on every one; dropping `target_contract`,
`created_by`, or `status` from the required set is still caught, so narrowing the
boundary did not weaken refusals; the command surface remains exactly ten commands
with no Proposal verb.

`go build ./...`, `go vet ./...`, and `go test ./... -count=1` are clean across all
fifteen packages. `mission check M7` reports `valid=true` over fourteen checks.

## The limitation of this record

The mutation replays above were run by the same actor who wrote the repairs. That
is precisely the arrangement independent review exists to avoid. A reader should
treat the b0a6c5d audit as independently established and everything after it as
operator-verified only.

If independent confirmation of the repairs matters for this Mission's downstream
use, the honest remedy is a fresh independent review at this commit, not a
stronger reading of this record.
