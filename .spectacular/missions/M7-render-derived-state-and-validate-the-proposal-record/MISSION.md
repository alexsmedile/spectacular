---
type: Mission
id: 01a00af1-38c0-7268-9529-5856afc7b2f2
title: Render derived state and validate the Proposal record
status: completed
created: "2026-08-16T14:19:42Z"
updated: "2026-08-16T20:46:01Z"
activation:
    at: "2026-08-16T21:30:00Z"
    by: Alex
    fingerprint: sha256:59d80c9784f7b580ecc6c72eb1d3563f19ca2125b20e12ddd0a67b0580a4168e
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
    branch: m7-derived-state
    commit: 127dac140467a462c3810c85c9ca325c18278a14
completion:
    - claim: state-line
      pass_boundary: Compact Mission and Objective projections lead with one state line stating lifecycle position and a NEXT line stating the next action and who holds it, with every field derived from data the bundle already carries and no canonical field added.
      proof_requirement: Golden fixtures over compact and promoted bundles assert the rendered state line and NEXT line, and assert every field traces to a bundle field; a negative test proves no projection writes to the canonical tree.
    - claim: drift-flags
      pass_boundary: Each frozen completion claim carries named drift flags derived from repairs consumed, verdict state, and fingerprint age, and an unnamed audit target defaults to the most-flagged claim showing the flags that selected it.
      proof_requirement: Table-driven fixtures with known repair counts, verdict states, and fingerprint ages assert the exact flag set, the ranking, and the default selection including tie behavior.
    - claim: authority-table
      pass_boundary: Mission check renders the authority decision table from the vocabularies the authority-vocabulary validator already resolves, and an undeclared verb is refused with code, field, problem, declared vocabulary, and safe correction, without adding any command or noun to the frozen surface.
      proof_requirement: Table-driven lookups cover declared operator verbs, declared owner-gated verbs, and undeclared verbs; a surface test asserts the accepted command list is unchanged.
    - claim: render-equivalence
      pass_boundary: A Mission with inline Objectives and the same Mission with those Objectives promoted produce byte-identical show, graph, and state-line output, and the Objective view falls back from graph to level sets only when the graph would exceed terminal width.
      proof_requirement: Golden tests assert byte equality across both representations for every projection surface, and view fixtures assert the approved notation plus the exact width threshold that selects level sets.
    - claim: proposal-schema
      pass_boundary: A Proposal validates against the compact authored shape of type, id, ref, title, status, created_by, created, updated, and target_contract, with ref required on new records, legacy human_ref decoded and reported as drift rather than refused, and no creation command provided.
      proof_requirement: Every Proposal in the repository validates with legacy fields preserved and human_ref reported as drift; negative tests assert refusal on missing required fields and assert no proposal creation command is accepted.
completion_record:
    at: "2026-08-16T20:46:01Z"
    authorization: owner supplied --by after schema checks
    by: Alex
    review: RV1
    reviewed_commit: 80f69e4e0b481566d52b0e9f8b113e3fc041c604
contract:
    fingerprint: sha256:1ffd39b498b44dce4e77cdf902f5f827bdf40eb2b317573c38a405f9b9ae9a0b
    ref: Contract:01a00aae-8921-7b27-96a9-1a4c175e7dc6
dependencies:
    - M6 completed with independent review and owner acceptance.
gaps: []
objectives:
    - claims:
        - state-line
      id: 01a00af1-38c0-72bf-a3b2-0c8f1595b50d
      outcome: Extract a shared derivation layer that computes lifecycle position, next action, holder, per-claim drift, and Objective readiness from the typed bundle.
      ref: O1
      status: implemented
    - claims:
        - proposal-schema
      id: 01a00af1-38c0-734a-8811-d8950b2a3575
      outcome: Normalize Mission and Proposal ref resolution through one decoder that accepts ref and legacy human_ref and reports the legacy spelling as drift.
      ref: O2
      status: implemented
    - after:
        - O1
      claims:
        - state-line
      id: 01a00af1-38c0-73ab-ac5f-c9e3a8dd875f
      outcome: Render the compact state line and NEXT line on Mission and Objective projections.
      ref: O3
      status: implemented
    - after:
        - O1
      claims:
        - drift-flags
      id: 01a00af1-38c0-7d42-91f4-14299b5941a4
      outcome: Compute per-claim drift flags and default the unnamed audit target to the most-flagged claim.
      ref: O4
      status: implemented
    - after:
        - O1
      claims:
        - authority-table
      id: 01a00af1-38c0-7d9c-bf98-513f520fac13
      outcome: Render the authority decision table in mission check output from the vocabularies the existing validator resolves.
      ref: O5
      status: implemented
    - after:
        - O2
      claims:
        - proposal-schema
      id: 01a00af1-38c0-7597-b0d6-15a85de49676
      outcome: Validate the compact Proposal record and prove the legacy ref spelling is reported rather than refused.
      ref: O6
      status: implemented
    - after:
        - O3
      claims:
        - render-equivalence
      id: 01a00af1-38c0-702a-a2b5-6c1f0235f4b1
      outcome: Render the ASCII Objective graph and level sets, selecting between them by terminal width.
      ref: O7
      status: implemented
    - after:
        - O7
      claims:
        - render-equivalence
      id: 01a00af1-38c0-7803-94b5-9b9c84b38fbd
      outcome: Prove inline and promoted representations render byte-identically across every projection surface.
      ref: O8
      status: implemented
outcome: Compact projections state conclusions instead of records, audits aim at observed drift instead of memory, and a Proposal validates against the shape it is actually authored in.
owner: Alex
ref: M7
repair_budget: 3
review: clustered
reviews:
    - file: reviews/RV1-clustered-review-of-m7-derived-state-drift-authority-equivalence-and-proposal-schema.md
      id: 01a00c7b-43c0-7d36-94d6-206c32b17da5
      ref: RV1
      verdict: pass
run:
    current_objective: O1
    id: 01a00af1-38c0-7af4-9f24-a3074277d3c2
    operator: Alex
    ref: R1
    repairs: 0
    started_at: "2026-08-16T14:19:42Z"
    status: completed
scope:
    mechanical:
        - cmd/spectacular/
        - internal/
        - test/
        - skills/spectacular/generated/
        - .spectacular/
    semantic:
        - Derivation of lifecycle position, next action, drift, readiness, and authority answers from state the bundle already carries.
        - ASCII rendering of Objective structure, with representation equivalence between inline and promoted bundles.
        - Validation of the compact Proposal record and normalized ref resolution.
start_key: sha256:997654575e4f12fa346d28e9a51fa54496c54962df4ca2aa278c42c852e072cf
stops:
    - A projection writes to the canonical tree, caches a rendered view, or requires a new canonical field to render.
    - The accepted command surface grows beyond the ten commands CC-missioncli enumerates, or a new noun is introduced.
    - A state line, drift flag, or readiness conclusion disagrees with the record it summarizes.
    - Inline and promoted representations of the same bundle render differently on any surface.
    - Work reaches for dependency shape, fallbacks, or Mission order, which are frozen to M8.
validation:
    mode: cli
    schema: mission.v2
---
# Render derived state and validate the Proposal record

This Mission makes state the bundle already holds legible, and records the Proposal
shape that practice already uses. It adds no canonical Mission field and no command.

## Why no command is added

`CC-missioncli` enumerates a closed ten-command surface, and M6 names a growing surface
as a stop. Every surface added here is output on an existing command or a flag on one.
The authority answer in particular is free: `mission check` already runs the
`authority-vocabulary` validator, which resolves the full operator and requires-owner
vocabularies and then discards them. Rendering that table changes what is printed, not
what is computed or what can be invoked.

## Why the ref decoder is here and not in M8

M8 validates Mission order against a list of Mission refs. That is impossible while
M2, M3, and M4 say `human_ref:` and M5 and M6 say `ref:`. One decoder that accepts both
and reports the legacy spelling as drift closes the gap without rewriting any completed
Mission's frontmatter. `nextMissionRef` already reads both spellings; this extends that
tolerance to the validation path so M8 has one spelling to compare.

## Why this Mission completed under clustered rather than independent review

M7 was planned for independent review, and two independent reviewers did audit it
at commit `b0a6c5d`. Their findings are recorded below and were repaired. But the
repairs, and the two boundary amendments those reviewers prompted, landed after
they finished. No independent reviewer has seen the amended boundary or the fixes.

Recording an independent review bound to the current activation fingerprint would
assert that a reviewer passed work they never read. The owner chose to downgrade
the review mode to `clustered` instead, which states plainly that the final
verification was the operator's own. The weaker claim is the true one.

What the reviewers did establish still stands: twelve of fourteen mutations were
caught at `b0a6c5d`, and the two survivors are now caught by tests added in
response. What is not established is independent confirmation of that repair.

## Two completion criteria were narrowed after review

Independent review found two frozen criteria naming inputs the record cannot
supply. The owner narrowed both rather than implementing against a shape that does
not exist.

`drift-flags` named "evidence age and freshness". `Reviewer.Evidence` is a list of
bare references with no timestamps, so evidence age is not derivable. `Review.Created`
exists but measures when a review was recorded, not how old its evidence is;
implementing against it would have produced a flag whose name misdescribed what it
measured. The other three inputs — repairs consumed, verdict state, fingerprint age —
are implemented and tested.

`proposal-schema` named `scope` among the required fields. No Proposal in this
repository is refused for lacking it, and P1 does not carry it at all. Requiring it
would have invalidated an accepted record to satisfy a field the authored shape never
adopted. The same criterion also described P5 and P6 as using the current `ref:`
spelling; all six Proposals use `human_ref:`, and the text now matches the record.

## What is deliberately out of scope

Fallbacks, the artifact-versus-interface edge split, and Mission order are M8. They all
add or alter canonical fields and belong behind one freeze boundary. The Objective graph
rendered here draws the single `after:` edge that exists today; when the split lands, the
same renderer gains a second edge kind without changing what it already draws.
