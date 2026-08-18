---
type: Decision
id: 01a01241-982d-7a8b-a0c4-799c717abfdd
title: Retire an absorbed Proposal to the archive with a stated resolver
created_by: Alex
created: "2026-08-18T00:24:38Z"
updated: "2026-08-18T00:24:38Z"
actor: Alex
actor_role: owner
ref: D11-proposal-retirement
question: What happens to a Proposal once the work it explored has shipped?
disposition: retire-absorbed-proposals
rationale: >-
    A Proposal had a validated status and no end state. Nothing writes the field — there is no
    `proposal accept`, and `proposal check` only validates the value it finds — so a shipped
    Proposal kept reading `draft` until someone edited it by hand. P6, P7, and P9 each sat at
    `draft` after their work landed, which made the live `proposals/` folder describe a backlog
    that no longer existed. Status alone cannot carry the conclusion, because the field records
    what an owner asserted rather than what the workspace can derive: `accepted` says a Proposal
    was agreed to, not which Mission absorbed it. `resolved_by:` states the resolver, so the
    trail survives the record leaving `proposals/` and the status becomes checkable against a
    Mission that exists. Archiving then costs no new concept — the archive already holds dead
    but machine-readable records, admitted under an authorizing Decision and a fingerprint, and
    a retired Proposal is the same kind of thing: kept for its numbering and its reasoning, not
    for reading.
alternatives:
    - leave absorbed Proposals in place at `accepted`, accepting that the live folder grows without bound and mixes open questions with settled ones
    - delete an absorbed Proposal once its Mission completes, accepting that the reasoning behind a shipped design is lost and its ref becomes reusable
    - derive the status from Mission state instead of storing it, which reverses the direction of the record and makes an owner assertion a computed value
authority_basis: Owner reviewed the absent Proposal lifecycle, chose archiving over deletion, and authorized retiring the three absorbed Proposals while keeping P5 live.
authorized_effects:
    - proposal.retire-absorbed
    - archive.admit-proposals
conditions:
    - resolver-named-before-archiving
    - open-proposals-stay-live
scope:
    - v2
targets:
    - Proposal:01a00a98-32b3-7612-b19a-b8ffa479505c
    - Proposal:01a00f7b-e046-700f-9b13-ca4b04d03790
    - Proposal:01a0105e-6520-7eff-9a8e-5a8a100674ab
evidence:
    - Mission:01a00af1-38c0-7268-9529-5856afc7b2f2
    - Mission:01a0102c-a360-71fe-a1be-8e1b010460b2
    - Mission:01a010a6-01b0-7323-a6a6-7bc38c571762
supersedes: ""
---
# Decision

A Proposal is optional, mutable exploration. It had five legal statuses and no defined
ending. Nothing in the system advances one, so `draft` persisted through shipping, and the
live folder stopped distinguishing an open question from a settled one.

## What is decided

A Proposal whose work has shipped is **absorbed**, and an absorbed Proposal is retired:

1. It gains `resolved_by:` naming the Mission that absorbed it, and its status becomes
   `accepted`.
2. It moves to `.spectacular/archive/proposals/`, carrying `archive_authorization:` and
   `archive_input_fingerprint:` like any archived record.
3. Live `proposals/` then holds only Proposals that are still open questions.

`resolved_by:` is written before the move, never after. Archiving a Proposal without a named
resolver is what this Decision exists to prevent: once the record leaves `proposals/`, that
field is the only thing connecting it to the work that answered it.

## What is retired now

| Proposal | Absorbed by | Shipped |
|---|---|---|
| P6 — Condense projection surfaces | M7 | `mission show` states a conclusion and a NEXT line |
| P7 — Amend a signed record without breaking it | M11 | `contract amend`, 2.3.0 |
| P9 — Frozen Handoff records | M12 | `handoff record`, 2.4.0 |

## What stays live

P8 and P10 are unimplemented and remain open questions.

P5 stays live deliberately, and it is the case that shows why status alone was not enough.
Three of its four directions shipped — drift-aimed audits, `fallbacks:`/`invalidated_if:`,
and the `after_interface:` edge — while the fourth, a decaying preflight trace, was never
built. Retiring it would bury a live idea; marking it `accepted` would claim more was agreed
than was. A Proposal is absorbed when the question it asked has been answered, not when most
of it has, so P5 keeps its `draft` status and its absorbed directions are noted in place.

## What this does not do

No command is added. Retirement is an owner act recorded here, consistent with how M1 and M10
entered the archive under D3 and D9. Nothing enforces the archive's admission fields in code;
they are convention, which is why the authorizing Decision is written rather than assumed.
