---
type: Proposal
id: 01a0105e-6520-7eff-9a8e-5a8a100674ab
ref: P9
title: Frozen Handoff records
status: draft
created_by: Alex
created: "2026-08-17T15:36:56Z"
updated: "2026-08-17T15:36:56Z"
scope:
    - v2
target_contract: Contract:01a00a20-63dd-7670-97f1-9eb8e12adc3a
---

# Frozen Handoff records

Exploration for a possible Mission. Nothing here is frozen. The problem is demonstrated
and the freeze decision is made; the schema and the correction mechanism are directions.
Anything below may be dropped, split, or reversed at plan-freeze.

## The problem in one line

The artifact a fresh agent acts on when its context is emptiest is the only one in the
workflow that is unchecked, unversioned, and stored outside the repository.

## Where this came from

Four handoff documents were written by hand into a scratch directory across two sessions.
One of them was materially wrong.

The M9-to-M10 handoff asserted that a Contract bound only to completed Missions could be
amended — "M7, M8, and M9 are all completed, nothing is bound to it, so the edit is
safe" — and instructed the receiving agent to verify that status before editing. The
agent did verify it, exactly as told. The premise was still false: `validateContract` had
no status gate, so a completed Mission verified its binding as strictly as a live one.
Acting on the handoff refused four Missions at once.

The handoff was not careless. It was specific, it named the command to re-verify with,
and it stated a stop condition. It was wrong about a mechanism, and nothing in the system
could have caught that, because the document was a temp file.

Two more handoffs from the same set remain unexecuted. A third was written for M11's
independent review. None of them is in the repository.

## What already exists

This is not a new record type. It is an unbuilt one.

- `domain.Handoff` is a valid record type (`internal/domain/reference.go`).
- `internal/humanlayout/layout.go` routes it to
  `.spectacular/missions/<mission>/handoffs/H<n>-<shortkey>-<title>.md`, with `H` as the
  ref prefix, `mission:` as the parent field, and a short-key suffix matching Evidence
  and Decisions.
- `internal/discovery/discovery.go` recognizes `handoffs` as a standard collection root.
- `internal/runtime/compiler.go` defines `HandoffContract` with accountability,
  statuses, required return, host-pointer role, and provider boundary.
- `skills/spectacular/references/runtime.md` already specifies what a delegation binds
  to and what the receiver must return.

What is missing is any way to write one and any schema saying what one must carry. The
semantics were designed; the record was never instantiated.

## The inversion worth fixing at the same time

Handoff and Review are the same feature approached from opposite ends, and each has the
half the other lacks:

| | Handoff | Review |
|---|---|---|
| layout rule | yes, in `layout.go` | none |
| creating command | none | `review record` |
| path decided by | the layout system | a hardcoded join in `service.go` |
| collection root in discovery | yes | no |

So `reviews/` lands in the right place by convention inside one function rather than
through the layout system, and `handoffs/` has a layout rule with nothing to invoke it. A
Review record created by any other path would not land in `reviews/`. Whatever fixes
either should probably fix both, and the cheaper direction is to move Review's path into
`layout.go` rather than to hardcode a second one.

## Decided: a Handoff is frozen

A Handoff attests to a state at a moment. One edited after the receiving agent read it is
not the document that was acted on, and the whole value of the record is that a reader
can tell what the sender believed and when.

So it binds the way a Review binds: to a commit, a tree, and an activation fingerprint.
A receiving agent can then answer the question that matters most on arrival — does what I
was handed still describe this workspace — before acting on any of its content.

## The consequence to design for

A frozen Handoff cannot be corrected in place, and handoffs have been wrong. The
receiving agent needs a sanctioned way to say so.

This is the part that must not be left implicit. Silent divergence between a frozen
handoff and reality is worse than a mutable handoff, because the frozen one carries more
authority. Three shapes, and the choice is open:

- **A superseding Handoff** with a `supersedes:` ref. Fits the freeze exactly: the
  original stays as the record of what was believed, and the correction is its own
  attestation. Costs a second record for what may be a one-line correction.
- **An appended correction block** on the original, outside the fingerprinted envelope.
  Cheap and keeps the story in one file, but it makes the record partly mutable and needs
  a rule about what an appended block may contradict.
- **A finding on the Mission.** No Handoff changes at all; the divergence is recorded
  where the work is. Simplest, but a later reader of the handoff has no pointer to the
  correction.

Leaning toward a superseding Handoff, because it is the only one where the correction is
itself attributable and bound to a tree.

## What a Handoff must carry to be checkable

The wrong handoff would have been caught by one field: the distinction between what the
sender verified and what the sender assumed. It asserted a mechanism as fact when it had
only inferred it.

A candidate shape, matching the compact-record practice the Proposal and Mission schemas
already use:

- `mission` and the activation fingerprint it describes
- `reviewed.commit` and `reviewed.tree`, so staleness is answerable on arrival
- sender identity, and the sender's relation to the receiver
- the task, in the receiver's own terms
- **asserted** — what the sender verified, with how it was verified
- **assumed** — what the sender believes but did not check
- stops, and what to return
- `supersedes`, when correcting an earlier Handoff

The asserted-versus-assumed split is the load-bearing part. Everything else is
bookkeeping that `runtime.md` already describes.

## Open questions for the owner

- **Does this need a command, or is a Handoff authored and checked like a Proposal?**
  The command surface is at twelve and growth past it is a stop. A Handoff is written by
  an agent with full context and read by one with none, which argues for validation
  rather than generation — the same argument that gave Proposals `check` and not
  `create`. But a Handoff must bind to a commit, a tree, and a fingerprint, and having a
  human or agent transcribe those by hand is exactly where a stale binding comes from.
- **Which correction mechanism**, from the three above.
- **Is a Handoff scoped to a Mission, or can one exist without one?** The layout rule
  makes `mission:` the parent field, so today it cannot. Session handoffs that span
  Missions, or precede one, have no home under that rule.
- **Does Review's hardcoded path move into `layout.go`** in the same work, or is that a
  separate cleanup?
- **Is `asserted` versus `assumed` checkable at all**, or is it an honesty convention the
  schema can hold a place for but never enforce? A field nothing verifies is the defect
  P7 was about; this one may be worth keeping anyway, because naming the distinction
  changes what a writer writes.

## Not in scope

- Autopilot. The charter is a different record with its own fingerprint and expiry, and
  `runtime.md` already governs it.
- Objective promotion, which already has a command and a home.
- Any change to what a receiver may do. `runtime.md` states it: the receiver never
  changes Mission criteria, declares Evidence sufficient, or gains permission it did not
  have.

## Relationship to current work

M11 is implemented and awaiting independent review. Nothing here touches its scope, and
this Proposal is deliberately unfingerprinted work that cannot collide with a repair.

The handoff written for M11's review is itself an instance of the problem: it lives in a
scratch directory, it names a commit and a tree by hand, and it contains a section
listing what its author is least confident about — which is the asserted-versus-assumed
split, written in prose because no field exists for it.
