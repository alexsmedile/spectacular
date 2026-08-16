# Propose and define

## Propose

Interpret fuzzy intent against current Capability Contract truth. Draft a base-bound Contract
delta with target identity/version/fingerprint, exact additions/modifications/removals, rationale,
scope, assumptions, and Gaps. Show it to the owner. Only after confirmation invoke `proposal
create`; acceptance remains an owner Decision and never promotes the Proposal into current truth.

## Define

Compare one to three proportional candidate slices. For each, cover observable outcome,
independent proof, coherence, dependency burden, reversibility, learning value, integration path,
and the state left if later work is cancelled.

Record two separate verdicts:

- Design Sufficiency: `sufficient | needs-evidence | needs-decision`
- Slice Quality: `coherent | too-broad | fragmented | dependency-bound`

Compile the candidate with `mission prepare --input <json-file> --json`. A receipt is ready only
when the verdicts are `sufficient` and `coherent`, no blocking Gap remains, its exact source
fingerprints and baseline still match, and freshness has not expired. Material discovery returns
the affected path here and produces a new receipt; it does not create a Design lifecycle state or
mandatory design document.

Freeze one minimal completion criterion per evidence claim before Mission creation: exact claim,
pass boundary, proof requirement, and `automatic | clustered | independent` review level. The
receipt reports `unmet_requirements` and selects `adaptive-grill` only while a verdict, blocking
Gap, coherent slice, or criterion remains unresolved. Ask plain, decision-sized questions about
only those failures; make silent product assumptions explicit Decisions or Gaps. A receipt with no
unmet requirements goes directly to owner activation without a mandatory interview.

For each unresolved fork, ask in this compact grammar:

`plain outcome -> technical basis -> {action -> consequence} -> recommended default`.

Use concrete verbs such as `publish RC.2`, `revise the criterion`, or `leave the Mission defined`;
do not substitute internal labels such as `accept`, `bounce`, or `escalate` unless the same line
defines the exact state change. Once the owner accepts a recurring pattern, carry it as a default
and ask only for remaining variables.

Present the complete Mission envelope once: outcome, Objectives, frozen completion contract, scope,
authority/effects, evidence claims, dependencies/Gaps, checks, budgets, expiry, stops, recovery,
return, and preparation receipt. Only owner approval permits `mission create` and activation.
