---
type: Decision
id: 01a03099-2a4b-7aae-a16d-3e88502dbc6c
title: A schema field means Spectacular mechanically governs the frontmatter
created_by: Alex
created: "2026-08-23T21:48:59Z"
updated: "2026-08-23T21:48:59Z"
actor: Alex
actor_role: owner
alternatives:
    - keep schema on every entity including unvalidated ones
    - drop schema everywhere and let validators own their rules privately
    - validate an Atlas envelope so its schema claim becomes true
disposition: accepted
question: Which documents carry schema, and what exactly does declaring one promise?
rationale: D23 said every nameable workspace entity carries both type and schema. That overclaimed. Proposals, Decisions, and Contracts carry type without schema, so the rule described two document classes while presenting itself as universal. An adversarial review named the consequence, that a schema nobody enforces is ceremony wearing the clothes of a type system, and it invites tooling to rely on a guarantee that does not exist.
ref: D24-schema-field-mechanically-governs-frontmatter
scope:
    - v2
---
# A schema field means Spectacular mechanically governs the frontmatter

`type:` is universal. Every document that is a thing a reader can name declares
what it is.

`schema:` is not universal. It is a claim with a consequence: **Spectacular
governs this document, and its frontmatter is under mechanical check.** A
document carries `schema:` when a command validates it and refuses on drift. A
document that no command validates does not carry one.

This amends D23, which stated the pairing as universal.

## What the check reaches

Mechanical enforcement covers the frontmatter: required fields, permitted
vocabularies, reference shapes, ordering rules. That is the part a machine can
decide.

The body is not mechanically enforced. Prose is where a document explains itself,
and a validator that graded prose would either be wrong or would flatten the
writing into a form. The body is LLM work: an agent reads it, judges it, and
writes it.

The gap between the two is real and worth naming. A record can carry perfectly
valid frontmatter above a body that contradicts it, and nothing refuses. A future
check may *sniff* the body — heuristically, as a warning rather than a refusal —
to surface that drift. A warning is the honest register for a judgment a machine
cannot make cleanly. It must never become a refusal, because a false refusal on
prose would make the record unwritable.

## Current allocation

Governed, carries `schema:`: Campaign.

Governed by a command but identified through its typed record rather than a
`schema:` string: Mission, Proposal, Decision, Contract, Evidence, Handoff,
Review. These are validated; their schema is the record type itself.

Not governed, `type:` only: Atlas.

Nothing at all: `raw/`, which names no entity.
