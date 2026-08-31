---
type: Decision
id: 01a030b9-838f-7e34-b25f-cc3e1451dff3
title: Make ontology impact explicit in planning without adding a Mission field
created_by: Alex
created: "2026-08-23T22:24:19Z"
updated: "2026-08-23T22:24:19Z"
actor: Alex
actor_role: owner
alternatives:
    - add an ontology-impact frontmatter field and validator now
    - leave ontology changes implicit in implementation work
    - require a graph database before recording domain changes
disposition: accepted
question: How should a Mission account for semantic-model changes before Spectacular has an ontology conformance checker?
rationale: A new typed field or CLI command would create premature ceremony and falsely imply that prose-level ontology completeness can be mechanically decided. Existing Mission sources and frozen body text can make semantic impact visible now, while retaining a future option for narrow warning-only checks grounded in real project use.
ref: D27-ontology-impact-explicit-in-planning
scope:
    - v2
---
# Make ontology impact explicit in planning without adding a Mission field

When a Mission changes a public behaviour, data model, permission, workflow, or
lifecycle, its plan names the affected Vocabulary concepts and cites the
`VOCABULARY` Anchor through existing `sources:` or the Mission body. A Mission
with no semantic-model effect says `Ontology impact: none` and gives a short
reason.

This is guidance, not a new frontmatter field, hard validator, or claim that
all declared ontology relationships are enforced. A later owner Decision may
authorize a warning-only conformance check after several projects reveal stable,
objective checks worth automating.
