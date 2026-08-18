---
type: Proposal
id: 01a0105e-6520-7eff-9a8e-5a8a100674ab
ref: P9
title: Frozen Handoff records
status: accepted
resolved_by: M12
archive_authorization: Decision:01a01241-982d-7a8b-a0c4-799c717abfdd
archive_input_fingerprint: d3e36db09c661b92b58151c6e9589ec29277adbd0661b499e3db0dd4966f7f18
created_by: Alex
created: "2026-08-17T15:36:56Z"
updated: "2026-08-18T00:00:00Z"
scope:
    - v2
target_contract: Contract:01a00a20-63dd-7670-97f1-9eb8e12adc3a
---

# Frozen Handoff records

Exploration for a possible Mission. The five design decisions below are owner-accepted,
but a Proposal carries no authority and they bind only when a Mission plan freezes them.
The problem is demonstrated and the shape is settled; the field names and the slicing may
still change. Anything below may be dropped, split, or reversed at plan-freeze.

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

## Decisions

All five questions are resolved. They are kept here with their reasoning because the
Mission plan will freeze them and the reasoning is what the owner approved.

**1. Does this need a command? (decided: yes — `handoff record`, as a thirteenth command
the owner authorized)**

`review record` is the working precedent, and it answers the staleness worry this question
raised. A `ReviewDraft` is hand-authored with `reviewed.commit` and `reviewed.tree` written
in by the author, and `verifyReviewedGit` then resolves both against the real repository:
a commit that does not exist refuses with `invalid_known_field`, and a tree that does not
belong to that commit refuses with `stale_fingerprint`. So hand-transcription is safe
*when something verifies it*. A Handoff gets the same treatment.

That argues for a command rather than authoring-plus-checking, because the verification is
the point and it needs a caller. It also means the record lands through the layout system
instead of wherever the author happened to save it — which is the actual failure being
fixed here, since all four existing handoffs live in a scratch directory.

This takes the public surface to thirteen. **The owner authorized it, and rejected the
cap's previous wording as a rule.** `AGENTS.md` said growth past twelve was "a stop, not a
judgement call"; it now says adding a command requires owner authorization and the count is
reported rather than defended. Twelve was a proxy for *do not let the surface sprawl
unnoticed*, and stated as an absolute it would have refused a correct command to protect a
number. The general form of that defect is recorded in `FEEDBACKS.md` under "A boundary was
written as a hard stop when it was a proxy for a concern", where this is now the second
logged instance.

**2. Which correction mechanism? (decided: a superseding Handoff)**

A `supersedes:` ref on a new Handoff. It is the only one of the three where the correction
is itself attributable and bound to a tree, which is the same property that made freezing
the original worth doing. The original survives as the record of what the sender believed;
the correction is its own attestation with its own commit, tree, and sender.

The cost is a second record for what may be a one-line fix. That is accepted: the appended
block makes a frozen record partly mutable and needs a rule about what an appendix may
contradict, and a Mission finding leaves a later reader of the Handoff with no pointer to
the correction. A Handoff is never corrected by editing it, the same way a Gap is never
closed by deleting it.

**3. Is a Handoff scoped to a Mission? (decided: yes, scoped)**

`humanlayout/layout.go:169` already makes `mission:` the parent field for Handoff, alongside
Evidence, Decision, Gap, and Assessment. Unscoped session handoffs would need a new layout
branch, and there is no demonstrated instance — all four real handoffs were about a Mission.
If one appears, that is a later change with an example to design against.

**4. Does Review's hardcoded path move into `layout.go`? (decided: yes, same Mission)**

`missionbundle/service.go:540` builds the Review path with a literal
`filepath.Join("reviews", …)`. Handoff has the opposite problem: a layout rule with nothing
to invoke it. Fixing one and not the other preserves exactly the asymmetry this Proposal
names, and a Review record created through any other path would not land in `reviews/`.
The change is small — the layout system already handles this prefix family — and doing it
under a Mission that is adding the Handoff writer is the cheapest moment.

**5. Is `asserted` versus `assumed` checkable? (decided: no — keep it, report it)**

Not enforceable. Nothing can tell whether a sender actually verified what they filed under
`asserted`. The schema holds the distinction and never scores it.

Keep it anyway. The M9-to-M10 handoff failed for precisely this reason: it stated a
mechanism as fact — "nothing is bound to it, so the edit is safe" — when it had only
inferred it, and acting on it refused four Missions. Naming the field changes what a writer
writes, and it gives a receiving agent a list to re-verify before acting.

This is the same instrument as decision 5 in `P10`: report the fact, let judgment handle
it, refuse nothing. It is not the defect `P7` was about — `P7` concerned a Contract
declaring a check that silently never ran, which is a false promise of enforcement. A field
documented as unverified promises nothing.

## Folded in: the two open Gaps on `CC-missioncli`

The owner decided the Mission built from this Proposal also closes both Gaps currently open
on the bound Contract, declaring them under `resolves_gaps:`.

They belong here rather than in their own Mission because they touch the same surface. Both
live in the amendment path, both are textual-rewrite precision bugs, and this Mission is
already moving record paths into `layout.go` and adding a record writer. One review and one
gate covers all of it.

- **`gap-rewrite-matches-by-line`** — `amend.go:51` matches `blocked_on:` with a
  line-anchored regexp that cannot tell it is inside a block scalar. The fix is stated in
  the Gap: track scalar depth while walking the Gap entry, plus the adversarial fixture
  M11's independent review described.
- **`repoint-assumes-one-fingerprint`** — `amend.go:251` re-points a bound Mission with
  `strings.Replace(…, 1)`, rewriting the first occurrence of the old fingerprint anywhere in
  the file. The Gap left the approach undecided between anchoring to the `contract:` block
  and refusing when the fingerprint appears more than once. **Take the refusal**, for the
  reason the Gap itself gives: it is smaller, and it turns a silent corruption into a stated
  problem in a mechanism that rewrites records the owner is not reading. Anchoring can
  follow later if the refusal proves noisy.

Neither closes through `contract amend` alone. The command rewrites `blocked_on:` to
`resolution:`; the Gap closes when the work resolving it lands. `resolves_gaps:` freezes the
resolution wording at activation, and completion refuses while either is still open.

This Mission therefore requires `amend-contract` in `requires_owner`.

## Not in scope

- Autopilot. The charter is a different record with its own fingerprint and expiry, and
  `runtime.md` already governs it.
- Objective promotion, which already has a command and a home.
- Any change to what a receiver may do. `runtime.md` states it: the receiver never
  changes Mission criteria, declares Evidence sufficient, or gains permission it did not
  have.

## Relationship to current work

M11 completed on 2026-08-17 and its independent review is recorded. Nothing here touches
its scope.

The handoff written for M11's review is itself an instance of the problem: it lives in a
scratch directory, it names a commit and a tree by hand, and it contains a section
listing what its author is least confident about — which is the asserted-versus-assumed
split, written in prose because no field exists for it.

`P10` is independent of this and neither blocks the other. The two do share decision 5's
instrument — report the fact and let review judge it — which is now the settled answer in
this repository for a concern that cannot be mechanically scored.
