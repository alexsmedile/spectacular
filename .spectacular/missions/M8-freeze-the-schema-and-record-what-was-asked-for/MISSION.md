---
type: Mission
id: 01a00c59-6bb0-7be4-89cf-07e157ba6b5c
title: Freeze the schema and record what was asked for
status: active
created: "2026-08-16T20:59:04Z"
updated: "2026-08-16T20:59:04Z"
activation:
    at: "2026-08-16T20:59:04Z"
    by: Alex
    fingerprint: sha256:db565ea862ca28817d2626942fb15231fecf34b208afd9253e39df597ff281bb
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
    branch: m8-frozen-schema
    commit: fe165b2de654397eda3b9e1dd04755784fed1de6
completion:
    - claim: frozen-fallbacks
      pass_boundary: A Mission records fallbacks inside the activation fingerprint, so a fallback cannot be introduced mid-Run to escape a stop, and repair exhaustion returns every recorded fallback alongside the failure that consumed the budget, marking any recommendation as a recommendation rather than a match.
      proof_requirement: Fingerprint tests prove a mutated fallback invalidates activation and mutable Run progress does not; exhaustion fixtures assert the full fallback set is surfaced with the consuming failure, and that a ranked recommendation never suppresses the alternatives.
    - claim: interface-edge-split
      pass_boundary: An Objective declares artifact dependencies and interface dependencies as distinct edge kinds, and an interface-only dependency onto an unfrozen target is refused with code, field, problem, and safe correction.
      proof_requirement: Table-driven fixtures assert both edge kinds validate independently, that a cycle across mixed edge kinds is refused, and that an interface dependency on an unfrozen target names the unfrozen target in its refusal.
    - claim: mission-order
      pass_boundary: Mission order is declared as typed Mission refs resolved through the existing decoder, activation ahead of an incomplete predecessor is refused, and both ref spellings resolve without rewriting any completed Mission.
      proof_requirement: The real repository Missions M2 through M8 resolve as one order despite mixed ref and human_ref spellings; negative tests assert refusal on a dangling Mission ref, on a Mission-level cycle, and on activation ahead of an incomplete predecessor.
    - claim: graph-edge-kinds
      pass_boundary: The Objective graph draws both edge kinds distinguishably and the Mission timeline draws order edges, with inline and promoted bundles still rendering byte-identically on every surface.
      proof_requirement: Golden fixtures assert each edge kind renders with its own notation, that the M7 equivalence property still holds across both representations, and that a Mission with only artifact edges renders exactly as it does today.
    - claim: request-fidelity
      pass_boundary: A Mission records the request that produced it as a source, a capture time, and a list of distinct asks, each carrying a disposition of covered, deferred, or declined; a covered ask names completion claims that exist, a deferred or declined ask carries a stated reason, and completion refuses while any ask is undispositioned. The request text sits outside the activation fingerprint so it can be sharpened; the dispositions sit inside it so a covered ask cannot be silently relabelled after activation.
      proof_requirement: Fixtures assert refusal on an undispositioned ask, on a covered ask naming a claim that does not exist, and on a deferred ask with no reason; fingerprint tests prove that editing request text preserves activation while changing any disposition invalidates it; a negative test asserts the validator never infers a disposition the author did not write.
contract:
    fingerprint: sha256:1ffd39b498b44dce4e77cdf902f5f827bdf40eb2b317573c38a405f9b9ae9a0b
    ref: Contract:01a00aae-8921-7b27-96a9-1a4c175e7dc6
dependencies:
    - M7 completed with a recorded review and owner acceptance.
gaps: []
objectives:
    - claims:
        - frozen-fallbacks
      id: 01a00c59-6bb0-7ab9-8fc5-e2dcbeb8cc08
      outcome: Add fallbacks to the Mission model and the activation fingerprint, and return them at repair exhaustion.
      ref: O1
      status: implemented
    - claims:
        - interface-edge-split
      id: 01a00c59-6bb0-7c9f-a854-eda2c935a602
      outcome: Split Objective dependencies into artifact and interface edge kinds in the model and the DAG validator.
      ref: O2
      status: implemented
    - after:
        - O2
      claims:
        - mission-order
      id: 01a00c59-6bb0-7745-957a-5f0bd946071a
      outcome: Declare Mission order as typed refs resolved through the existing decoder and validate the Mission-level graph.
      ref: O3
      status: implemented
    - claims:
        - request-fidelity
      id: 01a00c59-6bb0-7c81-8e61-fb60b1bc5187
      outcome: Add the request record with per-ask dispositions and the request-coverage validator.
      ref: O4
      status: implemented
    - after:
        - O2
        - O3
      claims:
        - graph-edge-kinds
      id: 01a00c59-6bb0-71ee-ad32-37613afd93eb
      outcome: Draw both Objective edge kinds and Mission order edges without breaking representation equivalence.
      ref: O5
      status: implemented
outcome: The record carries frozen fallbacks, two Objective edge kinds, typed Mission order, and the request that produced the plan, so a reader can tell what was asked for and what was dropped.
owner: Alex
ref: M8
repair_budget: 3
request:
    source: chat, session opening
    captured_at: "2026-08-16T20:59:04Z"
    asks:
        - ask: Freeze fallbacks into the activation fingerprint and return them on repair exhaustion
          disposition: covered
          claims:
              - frozen-fallbacks
        - ask: Split Objective dependencies into artifact and interface edge kinds
          disposition: covered
          claims:
              - interface-edge-split
        - ask: Declare Mission order as typed refs and validate the Mission graph
          disposition: covered
          claims:
              - mission-order
        - ask: Draw Objective edge kinds and Mission order edges in graph projections
          disposition: covered
          claims:
              - graph-edge-kinds
        - ask: Record what was asked for with per-ask dispositions and request coverage validation
          disposition: covered
          claims:
              - request-fidelity
        - ask: Review open gaps and dead weight and perform a repository cleanup
          disposition: deferred
          reason: Needs its own cleanup Mission; outside M8 schema freeze boundary
review: independent
run:
    current_objective: O5
    id: 01a00c59-6bb0-788c-bcbe-9ed2cd017523
    operator: Alex
    ref: R1
    repairs: 0
    started_at: "2026-08-16T20:59:04Z"
    status: active
scope:
    mechanical:
        - cmd/spectacular/
        - internal/
        - test/
        - skills/spectacular/generated/
        - .spectacular/
    semantic:
        - 'Canonical schema changes: fallbacks, Objective edge kinds, Mission order, and the request record.'
        - Validation of two Objective edge kinds, a Mission-level order graph, and per-ask request coverage.
        - Rendering of the new edge kinds without breaking inline and promoted equivalence.
start_key: sha256:ebbe74c123d97060740fcffef367c0ee96c99f55ce638743062ba958304cec89
stops:
    - The accepted command surface grows beyond the ten commands CC-missioncli enumerates, or a new noun is introduced.
    - A projection writes to the canonical tree or caches a rendered view.
    - Inline and promoted representations of the same bundle render differently on any surface.
    - The request-coverage validator infers a disposition, matches ask text to claim text, or otherwise asserts a semantic judgment as a mechanical result.
    - A completed Mission's frontmatter is rewritten to satisfy a new field.
validation:
    mode: cli
    schema: mission.v2
---
# Freeze the schema and record what was asked for

Every change here adds or alters a canonical field, which is why they share one
freeze boundary. `fallbacks` enters the activation fingerprint, `after_interface:`
enters the Objective model, `after_mission:` enters the Mission model, `request:`
enters the Mission model, and the rendering cannot be frozen before the edge kinds
it draws are decided.

## Why request fidelity is here

A Mission's `outcome` is the agent's interpretation of what a human asked for.
Once frozen, nothing checks the interpretation against its source. If the request
was misread at planning time, every downstream gate faithfully enforces the
misreading: the completion criteria, the review, and the fingerprint are all
airtight around a possibly wrong premise.

This is not hypothetical. The session that produced M7 opened with a request to
"review the open gaps and dead weight and perform a cleanup" after M7 and M8. That
ask never became a Mission, never became a Gap, and survived only in chat
scrollback. Nothing in the record noticed, because nothing was watching.

## What is mechanical and what is not

The validator checks structure only: that every ask carries a disposition, that a
covered ask names claims that exist, and that a deferred or declined ask carries a
reason. Whether a claim genuinely answers an ask is a semantic judgment, and it
stays with review, where judgment already lives.

This split is deliberate. A validator that tried to match ask text against claim
text would produce false refusals on correct plans and false confidence on
plausible-sounding ones. The mechanical layer makes a dropped ask impossible to
hide; the review layer decides whether a covered ask is honestly covered.

## Why the fingerprint is split

The request text sits outside the activation fingerprint. An owner who sharpens an
ask mid-Mission should not invalidate activation — punishing clarification is
exactly the wrong incentive.

The dispositions sit inside it. Relabelling an ask from covered to deferred after
activation is precisely the drift this claim exists to catch, and it should cost a
re-activation and an owner signature.

## Claim five is severable

`request-fidelity` shares this Mission's freeze boundary because it adds a
canonical field, but it does not depend on the other four claims and nothing in
them depends on it. If it turns out to need more design than this plan assumes, it
can be lifted into its own Mission without disturbing the schema work.
