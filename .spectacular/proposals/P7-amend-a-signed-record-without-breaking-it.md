---
type: Proposal
id: 01a00f7b-e046-700f-9b13-ca4b04d03790
ref: P7
title: Amend a signed record without breaking it
status: draft
created_by: Alex
created: "2026-08-17T11:29:31Z"
updated: "2026-08-17T11:47:00Z"
scope:
    - v2
target_contract: Contract:01a00a20-63dd-7670-97f1-9eb8e12adc3a
---

# Amend a signed record without breaking it

Exploration for a possible Mission. Nothing here is frozen. The problem statement is
firm and demonstrated; the ceremony shape is a direction, not a design. Anything below
may be dropped, split, or reversed at plan-freeze.

## The problem in one line

A Contract can be signed but not amended, so a Contract that becomes wrong stays wrong.

## Where this came from

Mission M9 resolved a Gap on `CC-projsurf` in fact: `dead-v1-governance-code` asked for
a decision on whether removing unreachable governance code belonged with the Proposal
schema work or with a separate cleanup. M9 made that decision and did the work. But the
Gap is still written as `blocked_on:`, because M9 could not edit the Contract that
holds it.

An attempt to write that resolution down — one field, `blocked_on:` to `resolution:`,
matching the shape a sibling Gap in the same block already uses — produced this:

```
$ spectacular mission check M9
refused stale_fingerprint field contract.fingerprint: bound Contract content changed
```

Not only M9. M6, M7, M8, and M9 all refused simultaneously. Every Mission ever bound to
that Contract. All four are `completed`; none was in flight. The edit was reverted and
the workspace is green, so nothing is currently broken — but the Gap remains
mis-written, and no legal sequence exists that would let it be corrected.

## What the fingerprint gets right

The seal is not the problem and should not be weakened.

A Contract is signed. A Mission's plan is signed. The signature is the owner's
approval, and the fingerprint is what makes that approval verifiable rather than
remembered. `CC-missioncli` states the principle in `required_behavior`:

> Fingerprint the frozen semantic envelope at activation while excluding mutable
> Mission, Objective, Run, and repair progress.

> A changed frozen semantic field invalidates the activation fingerprint; mutable
> progress does not.

That is correct. Keep it.

## What is missing

There is no ceremony for reopening a signed record.

Signed currently means welded shut. The system distinguishes *sealed* from *unsealed*,
but not *sealed* from *deliberately reopened in front of the owner*. Without that third
state, the only ways to correct a wrong Contract are to leave it wrong or to break every
Mission bound to it.

Two consequences worth recording separately, because they may want different fixes:

**The check does not distinguish an honest amendment from a rewrite.** `validateContract`
in `internal/missionbundle/validate.go` hashes the entire Contract file. Every byte —
prose, whitespace, `updated:`, a comma. Recording that a Gap was resolved and rewriting
the outcome produce the same refusal, because the check cannot tell them apart. The
activation check in the same file is field-scoped and can. The Contract check is the
odd one out, and it is odd against a principle its own Contract states.

**The check does not distinguish a completed Mission from a live one.** `validateBaseline`
and `validateActivation` both early-return on `defined` status. `validateContract` has
no status gate at all, so a finished Mission re-verifies its Contract binding as
strictly as a running one. A Mission's binding records which Contract text it was
executed against. Once complete, that is history. Re-hashing it asks a question with no
remaining consequence and produces a refusal with no available correction.

**The refusal promises a path that does not exist.** Its own correction text reads
*"review the Contract delta and explicitly amend or restart the Mission."* There is no
amend command. A reader who follows the instruction finds nothing there.

## What the owner actually wants

Three things, in order: what is wrong, what the fix is, and how to continue. Today the
system supplies the first and stops.

A refusal that names a correction the system cannot perform is worse than a refusal that
admits the dead end, because it costs the reader a search before they learn it is closed.

## The workaround that was assumed to exist does not

M9's own body records the sweep and states that "the Contract is amended after M9
completes, when nothing is bound to it." That sentence describes a window that does not
exist. `validateContract` has no status gate, so a completed Mission verifies its
Contract binding as strictly as a live one — the amendment refuses after completion for
the same reason it refused during. `TODO.md` records the same collision on 2026-08-16,
when a note written into one of these Gaps had to be reverted.

This matters beyond one deferred edit. "Record the resolution in the Mission body and
amend the Contract between Missions" is the current practice and the thing that happens
by default if nobody decides. It cannot work. There is no between.

## Direction — a small ceremony

Open, fix, close, and log. Sketch, not a design:

```
$ spectacular contract amend CC-projsurf --reason "M9 resolved dead-v1-governance-code"
  4 Missions bound: M6 M7 M8 M9 (all completed)
  proceed? [y/N]
```

The amendment is owner-authorized at the gate, which is the same approval the signature
represented in the first place. Then the Contract is edited, resealed, and the change is
recorded — what changed, why, by whom, when. The log is the point. An amendment that
leaves no trace is the thing the fingerprint exists to prevent; an amendment that leaves
a trace is ordinary governance.

Bound Missions have their `contract.fingerprint` re-pointed to the new value.

**Re-pointing a Mission is not reopening it.** This distinction is load-bearing and
should survive into any Contract that comes from this Proposal. A completed Mission's
plan — its outcome, claims, boundaries, authority — stays frozen and is never rewritten.
What updates is only the pointer recording which Contract text it was bound to. The
sibling Gap `mission-ref-frontmatter-drift` already establishes the principle in the
narrower case: *"Completed Missions are not rewritten."*

A closed Mission is closed. A live Mission has room to change. A Contract is a living
agreement and accumulates amendments.

## Decided by the owner

**A completed Mission reports Contract drift instead of refusing on it.**

`validateContract` gains a status gate, matching `validateBaseline` and
`validateActivation`, which both already early-return rather than checking a boundary
that no longer constrains anything. For a completed Mission the bound Contract's
fingerprint is history: it records which Contract text the work was executed against,
and re-hashing it asks a question with no remaining consequence and no available
correction.

Drift is still surfaced. A reader deserves to know the Contract moved after completion;
what they do not deserve is a refusal with no legal fix. The existing `Notices` channel
carries this — `ref-spelling-drift` in `internal/missionbundle/drift.go` already
establishes report-not-refuse for exactly this shape, so the mechanism does not need to
be invented.

This is the smallest change that makes the deferred amendment legal. It does not weaken
the seal for work in flight: a live Mission still refuses, because there the Contract
genuinely constrains what is being judged.

## Approaches considered and declined

- **Move Gaps out of Contracts into their own record type.** A Mission could then close a
  Gap without touching the frozen Contract. Declined: it adds a noun, and
  `CC-missioncli`'s ten-command surface is closed with growth named as a stop. It also
  costs a migration to solve a problem the status gate solves without one.
- **Exclude the `gaps:` block from the Contract fingerprint, the way Run progress is
  excluded.** Cheapest to implement and the most dangerous. A Gap is part of what the
  owner accepted; excluding it would let a live Mission quietly retire a blocker it was
  supposed to respect. Declined explicitly — the seal is not weakened to make this easy.
- **Codify the current workaround as a named between-Missions phase.** Declined because
  it does not work, not because it is unattractive: there is no window between Missions,
  as recorded above.
- **Re-activate a bound Mission to pick up the amended Contract.** Declined: re-activation
  recomputes the activation fingerprint, which reviews bind to, so it would invalidate
  passing reviews as a side effect of an unrelated Contract edit. Re-pointing
  `contract.fingerprint` is a narrower operation and leaves the frozen envelope alone.

## Open questions for the owner
- **Is field-scoped hashing needed for Contracts, or does the ceremony make it
  unnecessary?** If an amendment is explicit and logged, whole-file hashing may be fine —
  the owner approved the delta, so the check no longer needs to classify it.
- **Where does the amendment log live?** On the Contract itself, in `.spectacular/decisions/`,
  or in the transaction record. A Contract that carries its own amendment history is
  self-describing but grows; a separate log keeps the Contract stable but splits the story.
- **Does the same ceremony apply to a live Mission's plan?** Amending a frozen plan
  mid-Run is a different risk than amending a Contract, and may deserve a stricter gate
  or an outright refusal.
- **Does amendment need a scope limit?** An amendment that may touch any field is a
  rewrite with extra steps. It may be worth naming which fields an amendment may reach.

## Note on where this is written

This Proposal exists here rather than as a Gap on `CC-missioncli` because of the problem
it describes. The natural home for the finding is the Contract that owns the freeze
model. Adding a `gaps:` entry there would break every Mission bound to it — the same
refusal, for the same reason.

The bug prevents recording the bug in the place it belongs. A Proposal is unfingerprinted
and freely editable, so it is the one surface the problem does not block.

That is itself an argument for the ceremony.

## Relationship to current work

`CC-projsurf`'s `dead-v1-governance-code` Gap stays written as `blocked_on:` until this
is resolved. That is a known, deliberate inaccuracy, not an oversight — the record
cannot currently be corrected. Anyone reading that Gap should read this Proposal
alongside it.

The other two open Gaps on that Contract, `lifecycle-diagram-ungenerated` and
`concurrent-run-timelines`, are genuinely open and unaffected.
