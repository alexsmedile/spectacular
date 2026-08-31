---
type: Decision
id: 01a03008-1c7a-7c6b-bd0e-ea2f33df6be4
title: Identify workspace entities with type and version them with schema
created_by: Alex
created: "2026-08-23T19:10:33Z"
updated: "2026-08-23T19:10:33Z"
actor: Alex
actor_role: owner
alternatives:
    - keep the combined atlas_schema and campaign_schema keys
    - add a compatibility reader accepting both spellings
    - give every non-governed document a validating CLI command
disposition: accepted
question: How should non-governed workspace documents identify themselves, and where does unstructured thinking live?
rationale: An entity's identity and its validation rules are two facts, not one. Collapsing them into a single atlas_schema or campaign_schema key made non-governance a property of the frontmatter, when it is really a property of discovery skip-listing. Separating type from schema lets a reader and an agent name what a file is without asking whether the CLI validates it, and lets the schema version move independently of the entity.
ref: D23-workspace-entities-type-and-schema
scope:
    - v2
---
# Identify workspace entities with type and version them with schema

Every workspace entity that is a thing a reader can name carries two frontmatter
fields:

- `type:` states what the entity is: `Atlas`, `Campaign`, `Proposal`.
- `schema:` states which rules validate it: `spectacular.campaign.v2`.

Governance is not expressed in frontmatter. A document is non-governed because
its directory is skip-listed in discovery, not because it hides its type.

`raw/` is the exception that proves the rule: it holds unstructured thinking, it
is gitignored, it is skip-listed, and nothing in it is an entity. It carries no
frontmatter because there is nothing there to name.

The combined `campaign_schema:` and `atlas_schema:` keys are replaced rather than
read alongside the new spelling. This repository is the only consumer, and
AGENTS.md forbids compatibility readers. The Campaign schema moves to
`spectacular.campaign.v2` in the same cut.

An Atlas gains a `type:` but no validator. It is a top-view map used at the
thinking-init stage and later attached to Proposals and Contracts; it has nothing
mechanical to check, so it gets no command. Its frontmatter template lives in
`.spectacular/atlas/README.md` rather than behind a CLI slot.

A Proposal remains the first official Spectacular document and keeps full record
frontmatter. Making it minimal would make it a heavier Atlas rather than a
lighter Proposal.
