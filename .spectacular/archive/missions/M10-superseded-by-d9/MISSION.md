---
type: Mission
id: 01a00f98-0480-7ea2-9f3e-8e3a961aacc6
title: Make a resolved Gap closable and a declared check real
status: resolved
created: "2026-08-17T14:24:36Z"
archive_authorization: Decision:01a0103a-ce5c-72bc-a1d9-fb5e2f3b27ad
archive_input_fingerprint: e212c8a2224cf0b560cf563734a8cdf9e8a73719676a425638017712a0377ea8
updated: "2026-08-17T14:42:21Z"
activation:
    at: "2026-08-17T14:24:36Z"
    by: Alex
    fingerprint: sha256:d6ca09488b06c53c09aae4e7bde6a49ef105339030f2fbc307ddf781805d406a
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
        - change-outcome-or-completion
        - expand-scope
        - push
        - merge
        - release
        - irreversible-change
        - destructive-data
        - secret-change
baseline:
    branch: main
    commit: b55a9a34082f1d543483d0daecdd86daa8807e04
completion:
    - claim: completed-mission-reports-drift
      pass_boundary: The bound-Contract check gates on Mission status. A Mission with status completed does not refuse on `contract.fingerprint`; it reports the drift as a notice and remains valid. A Mission with any other status refuses exactly as it does today, with the same code, field, and correction. The notice names the Contract and states that it changed after completion, so a reader learns the fact without being blocked by it.
      proof_requirement: Table-driven cases cover every Mission status against both a matching and a drifted Contract fingerprint, asserting refusal for live statuses and a notice for completed. A regression case reproduces the failure this Mission exists to fix — amend the `gaps:` block of the Contract that M6 through M9 are bound to, then assert all four report valid with a drift notice. `bash test/verify.sh all` passes.
    - claim: mission-declares-resolved-gaps
      pass_boundary: A Mission plan may declare `resolves_gaps:` as a list of Gap ref and resolution-text pairs naming Gaps on the Contract it is bound to. Each ref is validated for existence on that Contract at plan-freeze; a ref naming a Gap that does not exist refuses with code, field, the available Gap refs, and a safe correction. The field is part of the frozen semantic envelope and is included in the activation fingerprint, so a Mission cannot acquire the authority to amend a Contract after activation and the owner approves the exact resolution wording at the activation gate. Declaring nothing stays the default and adds no ceremony to a Mission that resolves no Gaps.
      proof_requirement: Fixtures cover a valid declaration, a dangling Gap ref, a Gap ref naming a Gap on a different Contract, an empty list, and an absent field. A fingerprint test asserts that adding, removing, reordering, or editing the resolution text changes the activation fingerprint. Existing Mission fixtures with no such field still freeze to their current fingerprints.
    - claim: completion-amends-the-contract
      pass_boundary: 'Mission completion amends the bound Contract as part of the gate it already presents, with no second command and no second owner gate. Before the gate it prints what it will amend — each Gap ref, the resolution text to be written, the fields touched, and the Missions to be re-pointed — so the owner approves a described change rather than an implied one. For each declared Gap, `blocked_on:` becomes `resolution:` with the frozen text, and an entry is appended to the Contract''s `amendments:` block recording time, owner, reason, fields touched, and the fingerprint before and after. Every Mission bound to that Contract has its `contract.fingerprint` re-pointed; no other field of any Mission is written. The amendment may touch only the `gaps:` block and editorial frontmatter and refuses any attempt to reach a field stating what was agreed. The whole amendment is one transaction: the Contract, its log, and every re-pointed Mission all land, or none do.'
      proof_requirement: An end-to-end fixture completes a Mission declaring a resolved Gap and asserts the Gap now reads `resolution:`, the entry survives rather than being deleted, the amendments entry records both fingerprints, and every bound Mission validates clean afterward. A fixture asserts the pre-gate preview names every Gap and Mission the amendment touches. Fault injection at each write boundary proves no partial amendment survives, reusing the existing transaction primitives. Cases cover a live bound Mission elsewhere on the same Contract which refuses, an already-resolved Gap whose resolution text matches which is a no-op, one whose text differs which refuses, and an attempt to amend a non-amendable field which refuses.
    - claim: declared-validators-resolve
      pass_boundary: 'Every name in a Contract''s `mandatory_validation:` resolves to either a registered validator or a declared notice, and the Mission check refuses a Contract declaring a name matching neither. The four currently unresolvable names are resolved at their source: the fallback validator is renamed to the declared `fallback-fingerprint-coverage`; the interface-dependency frozen-target check becomes its own registry entry split out of the dependency-graph validator so its refusal is traceable to the check that produced it; `ref-spelling-drift` is recognized as a notice rather than a validator, which is a category correction and not a change to what was agreed; and `proposal-schema-v2` is satisfied by the proposals-are-checkable claim. No declared name is silently ignored.'
      proof_requirement: A test asserts every declared validation name in every Contract in the workspace resolves, and that an added bogus name fails. Cases cover a name that is a validator, a name that is a notice, and a name that is neither. The renamed and split validators keep their existing refusal codes, fields, and messages, asserted by the existing table-driven mutation cases.
    - claim: proposals-are-checkable
      pass_boundary: 'A `proposal check <ref>` command wires the existing Proposal validator to the CLI, satisfying the required behavior to validate Proposals. It reports valid state, the checks run, and any notices, in both compact and JSON form, matching the Mission check shape. Proposal creation remains forbidden and the forbidden-command test keeps it: Proposals are authored as Markdown wherever the author keeps them — a file, an issue, a chat — and checked, never generated through ceremony. The public registry asserts eleven commands.'
      proof_requirement: The public registry test asserts eleven commands with `proposal check` in position and the same forbidden list, including proposal creation. Real-process tests prove compact and JSON output. The Proposal this Mission derives from validates through the new command and is the fixture. Cases cover a valid Proposal, one missing its ref, one using the legacy ref spelling which passes with a notice, and an absent ref which refuses.
    - claim: contract-version-is-read
      pass_boundary: '`contract_version:` stops being inert. It is validated as a positive integer, required on every Contract, and reported alongside the bound Contract''s ref. A Mission bound to an earlier version is not refused, re-pointed, or migrated: it ran against that version, that is a true fact about it, and the check states it as a notice rather than a problem. The Contract this Mission is bound to is bumped to version 2, because adding a command changes its `command_surface` — a semantic field — and that is what a version bump is for. The history of what changed is the commit history; no changelog, superseded copy, or migration is introduced.'
      proof_requirement: Cases cover a missing version, a non-integer, a zero or negative value, and a Mission bound to a version earlier than the Contract's current one, asserting a notice and a valid result rather than a refusal. A test asserts the reported version matches the bound Contract's field. The version bump is asserted against the command-surface change that motivated it.
    - claim: workflow-states-the-step
      pass_boundary: 'The close reference states that a Mission declares the Gaps it resolved and that completion amends the Contract. The line asserting there is no separate Contract reconciliation command is corrected rather than deleted: there is still no separate command, and it now says why — reconciliation is part of completion. The prepare reference states that a Mission resolving a Gap declares it at plan time. The Skill does not grow: the hand-written Gap sweep prose the previous Mission needed is no longer necessary, and line counts before and after are recorded.'
      proof_requirement: The references are re-read end to end for a Mission that resolves a Gap and one that does not, confirming neither path acquires steps. Line counts before and after are recorded in the Evidence.
contract:
    fingerprint: sha256:80336b159e296ba63b5d85c80a48f8e540ae07d9aac52cdcdba4730059378a48
    ref: Contract:01a00a20-63dd-7670-97f1-9eb8e12adc3a
dependencies:
    - M9 completed with a recorded review and owner acceptance.
    - P7 records the problem, the demonstrated refusal, and the owner decisions this Mission implements.
fallbacks:
    - approach: Ship only the status gate on the bound-Contract check, leaving the declaration and the amendment unbuilt.
      invalidated_if: The amendment path proves to need a command of its own, which is a stop and returns the surface question to the owner.
      rejected_because: A completed Mission that reads cleanly still cannot write a resolution, so the Gap stays mis-written and the workflow keeps the hole this Mission exists to close. It is a real fallback because it is the one piece that makes the deferred amendment legal on its own.
    - approach: Amend the Contract through a dedicated command rather than at completion.
      invalidated_if: Amendment is needed outside completion — to correct a Contract no Mission is resolving — in which case the surface question reopens deliberately.
      rejected_because: It invents a second owner gate for a decision the completion gate already covers, and amendment is a consequence of completing work rather than an independent act.
    - approach: Correct the Contract's mandatory_validation names to match what the code already calls them.
      invalidated_if: A declared name turns out to describe a check that should not exist.
      rejected_because: The Contract described the intended checks correctly and the code never named its own work; renaming the declaration would record the drift as the agreement. The exception is the notice miscategorized as a validator, which is a category error and is corrected in the Contract.
gaps: []
objectives:
    - claims:
        - completed-mission-reports-drift
      id: 01a00f98-0480-7c24-ba07-8a9ce8bffd69
      outcome: Gate the bound-Contract check on Mission status so a completed Mission reports drift instead of refusing on it.
      ref: O1
      status: implemented
    - claims:
        - mission-declares-resolved-gaps
      id: 01a00f98-0480-70f2-b7ad-1f949dfe6793
      outcome: Add a frozen resolves_gaps declaration validated against the bound Contract's Gap refs at plan-freeze.
      ref: O2
      status: implemented
    - after:
        - O1
        - O2
      claims:
        - completion-amends-the-contract
      id: 01a00f98-0480-7252-88d2-69b3d22245a2
      outcome: Amend the bound Contract at completion in one transaction, previewing the change before the owner gate and re-pointing every bound Mission.
      ref: O3
      status: pending
    - claims:
        - declared-validators-resolve
        - proposals-are-checkable
      id: 01a00f98-0480-7c63-b4d3-113b51ab4d46
      outcome: Make every declared validation name resolve, and wire the existing Proposal validator to the CLI as an eleventh command.
      ref: O4
      status: pending
    - after:
        - O4
      claims:
        - contract-version-is-read
      id: 01a00f98-0480-7afc-bf84-5ba6c6b51775
      outcome: Validate and report contract_version, and bump the bound Contract to version 2 for the command it gains.
      ref: O5
      status: pending
    - after:
        - O3
        - O5
      claims:
        - workflow-states-the-step
      id: 01a00f98-0480-76b1-8d49-e0af550d14ef
      outcome: State in the workflow that a Mission declares resolved Gaps and that completion amends the Contract.
      ref: O6
      status: pending
outcome: A Mission that resolves a Gap closes it at completion in one owner gate, a completed Mission reads its Contract without refusing on a change it cannot cause, and every check a Contract declares either runs or fails loudly.
owner: Alex
ref: M10
repair_budget: 3
review: independent
run:
    current_objective: O1
    id: 01a00f98-0480-7ca7-bfcc-cadce5b19e56
    operator: Alex
    ref: R1
    repairs: 0
    started_at: "2026-08-17T14:24:36Z"
    status: active
scope:
    mechanical:
        - cmd/spectacular/
        - internal/
        - test/
        - skills/spectacular/
        - .spectacular/
    semantic:
        - A status gate on the bound-Contract check so a completed Mission reports drift instead of refusing.
        - A frozen `resolves_gaps:` declaration naming Gaps a Mission closes and the resolution text it will write.
        - Amendment of the bound Contract at completion, with an append-only amendments log and re-pointed Mission bindings.
        - Resolution of every declared `mandatory_validation:` name to a registered validator or a declared notice.
        - A `proposal check` command wiring the existing Proposal validator to the CLI.
        - Validation and reporting of `contract_version:`, and the bump of CC-missioncli to version 2.
        - Workflow guidance stating that a Mission declares resolved Gaps and completion amends the Contract.
start_key: sha256:7e90d373c673c85ea4d3ed5b7cbfc2418a36a82c041bcfd2bc9461303768eb0b
stops:
    - The accepted command surface grows beyond the eleven commands this Mission establishes, or a noun beyond Proposal is introduced.
    - A Gap is closed by deletion rather than by a named resolution.
    - An amendment needs to reach a semantic field to be useful, which would mean the editorial and semantic split is wrong and is a Contract question.
    - Re-pointing a bound Mission requires touching its activation fingerprint, which would mean re-pointing is re-activation and the distinction this Mission rests on fails.
    - A live Mission is found bound to any Contract this Mission would amend.
    - Wiring contract_version requires migration, superseded copies, or any handling of a bound Mission beyond reporting its version.
    - A completed Mission's frontmatter is rewritten beyond re-pointing its contract fingerprint.
    - A projection writes to the canonical tree or caches a rendered view.
validation:
    mode: cli
    schema: mission.v2
---
# Make a resolved Gap closable and a declared check real

## Origin

Mission M9 resolved the `dead-v1-governance-code` Gap in fact and could not write the
resolution down. Editing the Contract that holds it refuses with `stale_fingerprint` for
every Mission bound to it, and a status gate is absent from the bound-Contract check, so
completed Missions refuse exactly as live ones do. M9's own body records an amendment
deferred to "after M9 completes, when nothing is bound to it" — a window that does not
exist.

`P7` records the problem, the reproduction, the approaches declined, and the owner
decisions this Mission implements. Reproducing the refusal broke M6, M7, M8, and M9
simultaneously; the edit was reverted and the workspace left green, so the Gap is still
written as `blocked_on:` today.

## Why the declared-check work is in the same Mission

Designing the Gap mechanism surfaced the same defect in a second place. The Contract this
Mission is bound to declares fourteen validation names and all fourteen resolve; the
Contract holding the Gaps declares seven and four resolve to nothing. One is a rename, one
is folded inside another validator, one is a notice miscategorized as a validator, and one
is a validator built and never wired to the CLI.

A record that says something the system does not enforce is the same failure as a Gap
written as `blocked_on:` after it was resolved. Fixing the Gap instance while leaving four
unresolvable declarations would close one case and leave the class open.

## Execution plan

O1 is the smallest independent change and makes the deferred amendment legal on its own.
O2 adds the declaration, inert until O3 reads it. O3 applies the amendment at completion,
which needs both the field to read and the relaxed check so re-pointed Missions validate.

O4 is independent of the amendment path: it renames one validator, splits another out of
the dependency-graph validator, recategorizes the notice, and wires the Proposal validator
to the CLI. O5 follows O4 because the version bump is motivated by the command O4 adds. O6
is last so the documented workflow matches shipped behavior.

## The command surface goes to eleven

The previous Mission's stop named growth beyond ten commands. This Mission grows it to
eleven deliberately and states the case rather than smuggling it.

The Contract holding the Gaps already requires Proposals to be validated without providing
a creation command, and the forbidden-command test names proposal creation, never proposal
checking. The validator exists with no caller outside its own test, so the required
behavior was built and never wired — this Mission's own Proposal has never been validated
by anything.

A Proposal is the one record type that need not live in the workspace at all: it can sit in
an issue, a chat, or a scratch file. That is precisely why it needs `check` and not
`create` — the author brings a Proposal from wherever it lives and the tool says whether it
is well-formed. Creation ceremony is what would make Proposals mandatory, and that stays
declined. Growth beyond eleven remains a stop.

## What this Mission does not touch

Standalone and Mission-local Gap records stay unbuilt; the `Gap` record type and its
collection root exist unused, and the rule for which Gaps live where is an open question in
`P7` that must not be settled as a side effect. Amending a Contract bound to a live Mission
stays refused. Whole-file Contract hashing stays. The other Contract fields nothing reads
stay descriptive — `contract_version:` is wired because a decision needs it, not as part of
a sweep. Unused record types and their empty directories, and the missing way to regenerate
a hand-authored record's index, are recorded elsewhere and are not cleanup for this Mission.

## Review

Independent. The amendment path writes to canonical records across multiple files in one
transaction, and the reviewer must not be the agent that implemented it.
