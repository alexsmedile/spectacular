---
type: Mission
id: 01a010a6-01b0-7323-a6a6-7bc38c571762
title: Make a handoff a record and a record path one system
status: active
created: "2026-08-17T21:31:06Z"
updated: "2026-08-17T22:32:57Z"
activation:
    at: "2026-08-17T21:31:06Z"
    by: Alex
    fingerprint: sha256:ce103d57cc9322536c9117499ef6b2b6c3b08e2d937ce8b9bbfaa02a2c5f9452
authority:
    operator:
        - inspect
        - edit-in-scope
        - choose-reversible-implementation
        - run-checks
        - generate-derived-files
        - bounded-repair
        - commit-local
    requires_owner:
        - activate-mission
        - amend-contract
        - change-outcome-or-completion
        - expand-scope
        - push
        - merge
        - release
        - irreversible-change
        - destructive-data
        - secret-change
baseline:
    branch: m12-handoff-records
    commit: 26ff72d4c9b6fe9e8273579dc7960a71e33b7b20
completion:
    - claim: handoff-is-a-checkable-record
      pass_boundary: 'A Handoff is a validated record type with a schema. It carries `mission`, `reviewed.commit` and `reviewed.tree`, sender identity and the sender''s relation to the receiver, the task in the receiver''s terms, `asserted` and `assumed` as separate lists, stops, what to return, and an optional `supersedes` ref. `mission` is validated to name a Mission in the workspace and `supersedes` to name an existing Handoff on the same Mission; a dangling ref in either refuses with code, field, problem, and correction. `asserted` and `assumed` are required to be present and are never scored: the schema holds the distinction and no validator judges whether a sender actually verified an asserted item, which is documented at the field rather than left to be discovered.'
      proof_requirement: Table-driven fixtures cover a complete Handoff, one missing each required field, a dangling `mission`, a `supersedes` naming a Handoff on a different Mission, a `supersedes` naming nothing, an empty `assumed` list which is legal, and an absent `asserted` which refuses. A test asserts no validator inspects the content of `asserted` or `assumed`.
    - claim: handoff-record-writes-and-verifies
      pass_boundary: 'A `handoff record <mission-ref> <handoff.md|-> --by <sender>` command writes a Handoff into the Mission bundle and verifies its git binding the way `review record` does: a `reviewed.commit` that does not resolve in the repository refuses with `invalid_known_field`, and a `reviewed.tree` that is not that commit''s tree refuses with `stale_fingerprint`. The write is one atomic transaction over the Handoff and the Mission pointer list; a fault at any write boundary leaves the canonical tree byte-identical. A retry of the same logical record converges on the same identity without duplicates. The Mission carries a `handoffs:` pointer list mirroring `reviews:`. `mission show` reports the Handoffs a Mission carries and, for each, whether its bound tree still matches the working tree. The public command surface is thirteen and the count is stated.'
      proof_requirement: Real-process tests prove compact and JSON output for a successful record and for each refusal. Cases cover a nonexistent commit, a mismatched tree, a Mission that does not exist, and a stale activation fingerprint. Fault injection at each write boundary proves no partial Handoff survives, reusing the existing transaction primitives. The public registry test asserts thirteen commands with `handoff record` in position and the same forbidden list, including `proposal create`.
    - claim: a-handoff-is-corrected-by-superseding-it
      pass_boundary: A Handoff is never corrected by editing it. A correction is a new Handoff carrying `supersedes:` and its own commit, tree, and sender. The superseded Handoff survives unchanged as the record of what its sender believed. `mission show` marks a superseded Handoff as superseded and names the Handoff that replaced it, so a reader arriving at the original is pointed forward. A `supersedes` chain longer than one link resolves to the newest Handoff.
      proof_requirement: A fixture records a Handoff, then a superseding one, and asserts the original file is byte-identical afterward, that both appear in the Mission, and that the original is reported as superseded with a pointer to its replacement. A case covers a two-link chain. A case asserts a Handoff cannot supersede itself.
    - claim: record-paths-resolve-through-one-system
      pass_boundary: The Review record path is produced by `layout.go` rather than by the literal `filepath.Join("reviews", …)` in `missionbundle/service.go`. `Review` gains a layout branch with `RV` as its ref prefix and `mission:` as its parent field, matching the branch `Handoff` already has. Every existing Review in the workspace resolves to the path it occupies today, so no recorded Review moves and no Mission pointer is rewritten. A Review created through any path now lands in `reviews/`.
      proof_requirement: A test asserts the layout system produces the exact path of every Review already in the workspace, byte-for-byte against the recorded pointer. The hardcoded join is asserted absent from `service.go`. Existing review-record fixtures pass unchanged, proving the move is behavior-preserving.
    - claim: gap-rewrite-knows-its-scalars
      pass_boundary: 'The Gap rewrite in the amendment path tracks block-scalar depth while walking a Gap entry, so a line reading `blocked_on:` inside a `problem:` scalar body is not mistaken for the key. The rewrite still operates textually: an amendment''s diff touches the Gap it names and nothing else, and no block scalar elsewhere in the Contract is reflowed. A Gap entry with no `blocked_on:` key at any scalar depth refuses rather than splicing at the wrong place.'
      proof_requirement: 'The adversarial fixture M11''s independent review described: a Contract whose Gap `problem:` block scalar contains the literal text `blocked_on:`, asserting the real key is rewritten and the scalar body is byte-identical afterward. Cases cover the key at the entry''s top level, the decoy inside a scalar, both present, and neither. The existing amendment fixtures pass unchanged.'
    - claim: repointing-refuses-an-ambiguous-fingerprint
      pass_boundary: Re-pointing a bound Mission refuses when the old Contract fingerprint appears more than once in the Mission file. The refusal names the Mission, the fingerprint, and the line of every occurrence, and states that the amendment wrote nothing. A Mission carrying its binding exactly once re-points exactly as it does today. The refusal aborts the whole amendment transaction, so the Contract, its log, and every other Mission are left byte-identical. `--dry-run` reports the ambiguity as a would-refuse without writing.
      proof_requirement: A fixture builds a Mission quoting its own bound fingerprint in prose and asserts the amendment refuses, names both occurrences, and leaves the whole workspace byte-identical. Cases cover exactly one occurrence which succeeds, two occurrences which refuse, and the M9-shaped case of a body quoting three fingerprints none of which is its binding, which succeeds. `--dry-run` is asserted to report the refusal.
    - claim: the-workflow-states-the-handoff
      pass_boundary: The runtime and execute references state that a delegation is recorded as a Handoff, that a Handoff is frozen and corrected by superseding it, and that the receiver re-verifies what the sender filed under `assumed` before acting on it. The existing statement of what a receiver may not do is unchanged. Line counts before and after are recorded; the guidance is stated once and not repeated across references.
      proof_requirement: The references are re-read end to end for a Mission that records a Handoff and one that does not, confirming neither path acquires steps. Line counts before and after are recorded in the Evidence.
contract:
    fingerprint: sha256:6c19cc655e1cce0aa6bb8f9c227b0d8b1c353472f32f9c7155ac4aa665fefbf3
    ref: Contract:01a00a20-63dd-7670-97f1-9eb8e12adc3a
dependencies:
    - M11 completed with a recorded independent review and owner acceptance; the amendment path this Mission repairs is what M11 built.
    - P9 records the problem, the demonstrated failure, and the five owner decisions this Mission freezes.
    - Both Gaps this Mission resolves are open on CC-missioncli and were confirmed by M11's independent review.
fallbacks:
    - approach: Ship the Handoff schema and a check, with no command, leaving authorship to hand-written files validated afterward.
      invalidated_if: Hand-authored Handoffs land outside the layout system, which is the failure this Mission exists to fix.
      rejected_because: The verification of the git binding is the value, and a verifier needs a caller. Without a command the record lands wherever the author saved it, which is exactly how all four existing handoffs ended up in a scratch directory.
    - approach: Close the two Gaps in a Mission of their own, separate from the Handoff work.
      invalidated_if: The Gap fixes turn out to need a change to the amendment transaction that the Handoff writer would conflict with.
      rejected_because: Both Gaps live in the amendment path and both are textual-rewrite precision bugs. One review and one gate covers them alongside a Mission already moving record paths and adding a record writer.
    - approach: Anchor the re-pointing replacement to the `contract:` block instead of refusing on ambiguity.
      invalidated_if: The refusal fires on a legitimate Mission during this Mission's own verification, which would mean the ambiguity is common rather than pathological.
      rejected_because: The Gap itself argues the refusal is smaller and turns a silent corruption into a stated problem. Anchoring is a larger change to a mechanism that rewrites records the owner is not reading, and it stays available if the refusal proves noisy.
gaps: []
objectives:
    - claims:
        - handoff-is-a-checkable-record
      id: 01a010a6-01b0-77c3-81ef-c655a68b91c8
      outcome: Add the Handoff schema and its validation, including the asserted/assumed split, the supersedes ref, and reference integrity against the Mission and the superseded Handoff.
      ref: O1
      status: implemented
    - claims:
        - record-paths-resolve-through-one-system
      id: 01a010a6-01b0-7aa0-91eb-e4d065a567d9
      outcome: Give Review a layout branch so its path is produced by layout.go, and remove the hardcoded join from service.go without moving any recorded Review.
      ref: O2
      status: implemented
    - after:
        - O1
        - O2
      claims:
        - handoff-record-writes-and-verifies
      id: 01a010a6-01b0-770f-8ea6-d475a9a6d44b
      outcome: Add the handoff record command, writing a Handoff through the layout system in one atomic transaction and verifying its git binding against the real repository.
      ref: O3
      status: implemented
    - after:
        - O3
      claims:
        - a-handoff-is-corrected-by-superseding-it
      id: 01a010a6-01b0-78c0-b735-f7396847cc8e
      outcome: Make a Handoff correctable only by superseding it, and report the supersession where a reader of the original will see it.
      ref: O4
      status: implemented
    - claims:
        - gap-rewrite-knows-its-scalars
      id: 01a010a6-01b0-7289-bf41-e482299bd784
      outcome: Make the Gap rewrite block-scalar aware, closing gap-rewrite-matches-by-line.
      ref: O5
      status: implemented
    - claims:
        - repointing-refuses-an-ambiguous-fingerprint
      id: 01a010a6-01b0-78b7-895a-e22dc1448cef
      outcome: Refuse re-pointing when the old fingerprint is ambiguous, closing repoint-assumes-one-fingerprint.
      ref: O6
      status: implemented
    - after:
        - O4
        - O6
      claims:
        - the-workflow-states-the-handoff
      id: 01a010a6-01b0-7bf6-b633-0d178e951189
      outcome: State in the workflow that a delegation is recorded as a Handoff, that it is frozen, and that the receiver re-verifies what was assumed.
      ref: O7
      status: pending
outcome: A handoff between agents is a checkable, tree-bound record inside the workspace instead of a temp file, every record type resolves its path through one layout system, and both open Gaps in the Contract amendment path are closed.
owner: Alex
ref: M12
repair_budget: 3
resolves_gaps:
    - gap: gap-rewrite-matches-by-line
      resolution: The Gap rewrite tracks block-scalar depth while walking a Gap entry, so a `blocked_on:` appearing inside a scalar body is not mistaken for the key. The textual approach is kept deliberately, so an amendment's diff still touches only what it changed. An adversarial fixture carrying the literal text inside a `problem:` scalar asserts the correct key is rewritten and the scalar body is left byte-identical.
    - gap: repoint-assumes-one-fingerprint
      resolution: Re-pointing a bound Mission refuses when the old Contract fingerprint appears more than once in the Mission file, naming the Mission, the fingerprint, and every occurrence, rather than rewriting the first one. The refusal was chosen over anchoring to the `contract:` block because it is smaller and turns a silent corruption into a stated problem in a mechanism that rewrites records the owner is not reading. Anchoring remains available later if the refusal proves noisy.
review: independent
run:
    current_objective: O1
    id: 01a010a6-01b0-7320-acc2-5c695bec2843
    operator: Alex
    ref: R1
    repairs: 0
    started_at: "2026-08-17T21:31:06Z"
    status: active
scope:
    mechanical:
        - cmd/spectacular/
        - internal/
        - test/
        - skills/spectacular/
        - .spectacular/
    semantic:
        - A Handoff schema carrying the mission, a verified git binding, sender identity, the task, an asserted/assumed split, stops, a return contract, and an optional supersedes ref.
        - A `handoff record` command that writes a Handoff through the layout system and verifies its git binding against the real repository.
        - Correction of a Handoff by superseding it, never by editing it.
        - Resolution of the Review record path through `layout.go` rather than a hardcoded join.
        - Block-scalar-aware Gap rewriting in the Contract amendment path.
        - A refusal when re-pointing a bound Mission finds the old fingerprint more than once.
start_key: sha256:0c7e6ff3e3fa09516f1d442f0d50b294a0fb5aa86a11ae04652c9e172d19fef7
stops:
    - The command surface would grow beyond the thirteen commands this Mission establishes.
    - A Handoff would need to be editable in place for the correction path to be usable, which would mean the freeze decision is wrong and returns to the owner.
    - An unscoped Handoff belonging to no Mission is required, which would need a new layout branch and is a decision the owner made the other way.
    - Moving Review's path into layout.go would move any Review already recorded in the workspace.
    - Validating `asserted` or `assumed` content becomes necessary for the record to be useful, which would contradict the decision to report and never score it.
    - The Gap rewrite requires decoding and re-emitting canonical YAML, which would reflow block scalars the amendment did not change.
    - Closing either Gap requires reaching a Contract field beyond the `gaps:` block.
validation:
    mode: cli
    schema: mission.v2
---
# Make a handoff a record and a record path one system

## Origin

The artifact a fresh agent acts on when its context is emptiest is the only one in the
workflow that is unchecked, unversioned, and stored outside the repository.

The M9-to-M10 handoff asserted that a Contract bound only to completed Missions could be
amended — "M7, M8, and M9 are all completed, nothing is bound to it, so the edit is safe" —
and instructed the receiving agent to verify that status before editing. The agent did
verify it, exactly as told. The premise was still false: `validateContract` had no status
gate, so a completed Mission verified its binding as strictly as a live one. Acting on the
handoff refused four Missions at once.

The handoff was not careless. It was specific, it named the command to re-verify with, and
it stated a stop condition. It was wrong about a mechanism, and nothing in the system could
have caught that, because the document was a temp file. Four such handoffs exist today, all
in a scratch directory, two of them never executed.

`P9` records the problem and the five owner decisions this Mission freezes.

## This is not a new record type

`domain.Handoff` is already a valid record type. `humanlayout/layout.go` already routes it
to `.spectacular/missions/<mission>/handoffs/H<n>-<shortkey>-<title>.md` with `mission:` as
the parent. `discovery.go` already recognizes `handoffs` as a collection root.
`runtime/compiler.go` already defines `HandoffContract`. `references/runtime.md` already
specifies what a delegation binds to and what a receiver must return.

The semantics were designed and the record was never instantiated. What is missing is a
schema saying what one must carry, and any way to write one.

## Why the Review path moves in the same Mission

Handoff and Review are the same feature approached from opposite ends, and each has the half
the other lacks. Handoff has a layout rule with nothing to invoke it. Review has a creating
command whose path comes from a literal `filepath.Join("reviews", …)` inside
`missionbundle/service.go`, so a Review created through any other path would not land in
`reviews/`.

Fixing one and not the other preserves exactly that asymmetry. The layout system already
handles this prefix family, so the move is small, and the moment a Mission is adding a
record writer is the cheapest moment to make both go through one system.

## Why the two Gaps are here

Both open Gaps on the bound Contract live in the amendment path, and both are
textual-rewrite precision bugs in a mechanism that rewrites records the owner is not
reading. This Mission is already moving record paths and adding a record writer, so one
review and one owner gate covers all of it. `resolves_gaps:` freezes the resolution wording
at activation and completion refuses while either Gap is still open.

Neither closes through `contract amend` alone. That command rewrites `blocked_on:` to
`resolution:`; the Gap closes when the work resolving it lands.

## Execution plan

O1 and O2 are independent and can run in either order. O1 adds the schema, inert until O3
writes one. O2 moves the Review path, which O3 then relies on to place a Handoff through the
same system rather than adding a second hardcoded join.

O3 needs both: the schema to validate and the layout system to place the file. O4 follows O3
because supersession is a property of records that exist.

O5 and O6 are independent of everything above and of each other. They are the two Gap fixes,
each confined to one function in `internal/missionbundle/amend.go`.

O7 is last so the documented workflow matches shipped behavior.

## The command surface goes to thirteen

Twelve before, thirteen after. The owner authorized the growth and rejected the previous
wording of the cap, which stated growth past twelve as a hard stop rather than a matter for
owner authorization. `AGENTS.md` now says adding a command requires owner authorization and
the count is reported rather than defended.

`handoff record` exists because the verification is the point. A `ReviewDraft` is
hand-authored with `reviewed.commit` and `reviewed.tree` written in by the author, and
`verifyReviewedGit` resolves both against the real repository. Hand-transcription is safe
when something verifies it, and the verifier needs a caller. Without one, the record also
lands wherever the author saved it, which is the failure being fixed.

`proposal create` stays forbidden.

## What this Mission does not touch

Autopilot charters, which are a different record with their own fingerprint and expiry.
Objective promotion, which already has a command and a home. What a receiver may do:
`runtime.md` already states that a receiver never changes Mission criteria, declares
Evidence sufficient, or gains permission it did not have, and that is unchanged.

Unscoped Handoffs belonging to no Mission stay unbuilt. All four real handoffs were about a
Mission, and there is no demonstrated instance to design against.

## Review

Independent. This Mission adds a public command, moves a shared record-path interface, and
modifies a mechanism that rewrites canonical records in a transaction. The reviewer must not
be the agent that implemented it.
