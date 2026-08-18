---
type: Proposal
id: 01a00f7b-e046-700f-9b13-ca4b04d03790
ref: P7
title: Amend a signed record without breaking it
status: accepted
created_by: Alex
created: "2026-08-17T11:29:31Z"
updated: "2026-08-18T00:00:00Z"
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

## The same defect, found elsewhere

A declaration nothing verifies is the shape of this whole problem, and the Gap is not the
only instance. `CC-projsurf` declares seven names in `mandatory_validation:`. Four resolve
to nothing in the validator registry:

| Declared | Reality |
|---|---|
| `fallback-fingerprint-coverage` | registered as `frozen-fallbacks`; a naming mismatch |
| `interface-dependency-frozen-target` | implemented inside `objective-dependency-dag`, not separately named |
| `ref-spelling-drift` | not a validator at all — it is a notice, by design, and belongs in a different list |
| `proposal-schema-v2` | `ValidateProposal` exists with no caller outside its own test |

`CC-missioncli` declares fourteen and all fourteen resolve, so this is drift introduced by
newer work rather than a long-standing state. The behaviors mostly exist; what is missing is
any check that a declared name corresponds to something that runs. A Contract can therefore
promise a validation that never executes, and `mission check` passes.

That is the same failure as a Gap written as `blocked_on:` after it was resolved: the record
says something the system does not enforce, and nothing catches the divergence.

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

Open, fix, close, and log — riding the owner gate that already exists at completion rather
than adding a second one. Sketch, not a design:

```
$ spectacular mission complete M9 --by Alex
  will amend CC-projsurf:
    gaps.dead-v1-governance-code -> resolution: closed by M9 as a separate cleanup…
  will re-point contract.fingerprint on M6 M7 M8 M9 (all completed)
  proceed? [y/N]
```

A Mission declares at plan time which Gaps it resolves, so the owner approves the exact
wording at activation and sees it again before it is written. The amendment is a consequence
of completing the work, not an independent act, which is why it needs no command of its own.
Then the Contract is edited, resealed, and the change is recorded — what changed, why, by
whom, when. The log is the point. An amendment that leaves no trace is the thing the
fingerprint exists to prevent; an amendment that leaves a trace is ordinary governance.

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

- **Move *all* Gaps out of Contracts into their own record type.** A Mission could then
  close any Gap without touching the frozen Contract. Declined as the general answer: a
  Gap that qualifies a Contract is part of what the owner accepted, and relocating it to
  dodge the seal is the same weakening as excluding it from the fingerprint, one
  indirection further out. It also costs a migration to solve a problem the status gate
  solves without one. Standalone Gaps are still wanted for Gaps that were never
  Contract-shaped — see the decision below.
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

## The ceremony as decided

Four further decisions settle the shape. Together with the status gate above, they are
what a Mission would build.

**The log lives on the Contract.** An `amendments:` block in the Contract's own
frontmatter, appended once per amendment:

```yaml
amendments:
  - at: "2026-08-17T11:29:31Z"
    by: Alex
    reason: M9 resolved dead-v1-governance-code
    fields: [gaps.dead-v1-governance-code]
    from_fingerprint: sha256:a7ae29b…
    to_fingerprint:   sha256:4f1c8d2…
```

The Contract stays self-describing: a reader sees the agreement and how it got here in
one file, with no second record to find and no chance of the history going missing while
the Contract survives. It grows over time, which is accepted — an agreement that
accumulates its own amendments is the normal shape for a signed document. Note that the
block sits inside the fingerprinted file, so appending to it necessarily changes the
fingerprint; that is not a problem once amendment is a sanctioned operation, but it does
mean the log cannot be written by anything other than the ceremony itself.

**A live Mission's Contract is not amendable.** The ceremony applies only where no bound
Mission is live. A live Mission's binding is precisely what it is being judged against,
and changing it mid-flight is the drift the seal exists to catch. If a Contract is
discovered wrong while a Mission runs, the owner stops the Mission first. This keeps the
surface small and adds no new risk to work in flight — and it avoids the re-activation
trap, where picking up an amended Contract would recompute the activation fingerprint and
invalidate passing reviews as a side effect.

**Whole-file hashing stays.** Field-scoped hashing for Contracts is not needed once
amendment is explicit: the owner approved the delta at the gate, so the check no longer
has to classify it. Whole-file is simpler and catches everything. This leaves the
asymmetry with the activation fingerprint in place deliberately — the two checks answer
different questions, and only the activation one needs to distinguish semantic change
from mutable progress. Revisit only if a real case demands it.

**An amendment may reach Gaps and editorial fields only.** Closing a Gap, correcting
prose, and bumping `updated:` are amendable. Changing `purpose`, `outcome`,
`required_behavior`, or any other field that states what was agreed requires a new
Contract version instead. Without this limit, "amendment" becomes "rewrite with extra
steps" — the reason string is the only thing distinguishing the two, and nothing
validates a reason string. The limit is what keeps the ceremony honest.

**Versioning is `contract_version:` plus Git.** The field already exists on every
Contract and is the right place to carry the version. Semantic change bumps it; the
history of what changed is the commit history, which already records every Contract edit
with its diff, author, time, and message. No new versioning machinery, no stored
changelog, no superseded-Contract copies.

This settles the division cleanly. `amendments:` records the small, in-place corrections
the ceremony permits — the ones that must not silently rewrite an agreement. A
`contract_version:` bump marks a semantic change, and Git carries its story. Note that no
Go code reads `contract_version:` today; it is inert metadata, and `CC-v2prod` is already
at `"2"` from an ordinary commit with no mechanism behind it. Any Mission that acts on
versioning should make the field meaningful rather than assume it already is.

**Gaps may be repository-wide or Mission-local.** `.spectacular/gaps/` as a collection
root serves the first; a Mission-local Gap belongs with its Mission. This is why the
scaffolding exists and why it should stay: not every Gap is a limitation of a Contract.
A Gap discovered during execution that concerns only that Mission's work has no business
in a signed agreement, and forcing it there is part of what created the problem this
Proposal describes.

That does not reverse the declined approach above. Gaps that qualify a Contract stay
embedded in it, because they are part of what the owner accepted, and the ceremony is
what closes them. What the standalone form adds is a home for the Gaps that were never
Contract-shaped to begin with.

**A Mission bound to an earlier Contract version is simply outdated.** It ran against that
version; that is a true fact about it and not a problem to solve. No migration, no
superseded copy, no re-pointing on a version bump — `mission check` states the version and
moves on. This is what keeps `contract_version:` cheap enough to be worth wiring: the field
records which agreement the work was done under, and history answers everything else.

**`contract_version:` should be wired rather than left inert.** It is currently read by
nothing, which is why `CC-v2prod` reaching "2" meant nothing mechanically. If versioning is
where semantic change goes, the field has to be validated and reported, or the rule routes
change into a field no reader can trust.

**Proposal is a primitive and gets its own command.** `CC-projsurf` already requires
Proposals to be validated — *"Validate Proposals without providing a creation command;
Proposals are authored as Markdown and checked, not generated through ceremony"* — and the
forbidden-command test names `proposal create`, never `proposal check`. `ValidateProposal`
exists in `internal/missionbundle/proposal.go` with no caller outside its own test, so the
required behavior was built and never wired. This Proposal has therefore never been
validated by anything.

The surface goes to eleven commands deliberately. A Proposal is the one record type that
need not live in the workspace at all — it can sit in an issue, a chat, or a scratch file —
which is precisely why it needs `check` and not `create`: the author brings a Proposal from
wherever it lives and the tool says whether it is well-formed. Creation ceremony is what
would make Proposals mandatory, and that is the thing declined.

## Still open

- **Which Gaps belong where.** Repository-wide, Mission-local, and Contract-embedded are
  three homes with no stated rule for choosing between them. The rule matters more than
  the mechanism: without it, the standalone form becomes an escape hatch for Gaps that
  should have qualified a Contract and been subject to the ceremony.
- **`Gap` is already a record type with nothing using it.** `internal/domain/reference.go`
  lists `Gap` among the valid record types and `internal/discovery/discovery.go` already
  recognizes `.spectacular/gaps/` as a standard collection root, but no code reads either.
  The scaffolding is now wanted rather than suspect — it needs a schema, a validator, and
  a reference form before it is usable.

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
