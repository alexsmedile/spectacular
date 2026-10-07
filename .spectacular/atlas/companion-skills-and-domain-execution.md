---
type: Atlas
title: Companion skills and domain execution topology
version: '0.3'
created: '2026-09-01T13:09:35+02:00'
updated: '2026-10-07T18:49:34Z'
---

# Companion skills and domain execution

Spectacular owns durable context placement and optionally selected governance.
Domain skills own their specialized work. Choose one only when it changes the
quality or predictability of the task; no starting pack is mandatory.

| Need | Relevant domain skill | Useful output |
|---|---|---|
| Service boundaries and system trade-offs | system-architecture | Focused architecture choice and linked model |
| Database entities, schemas, indexing, migrations | data-modeling | Data model, DDL, and migration approach |
| Several viable design options | rapid-prototyping | Bounded alternatives and inspected evidence |
| Human-facing product docs in root docs/ | pageworks | Product documentation in its own structure |
| Skill behavior and runtime instructions | skill-forge | Reviewed Skill route and disclosed guidance |

A small-file convention alone does not require every domain skill to be reviewed.
Review the owning Skill when its behavior changes; use domain review when the
substantive architecture, schema, or design actually changes.

## Keep the output in its correct home

The implementation approach goes in plans/. Reusable behavior can go in optional
specs/. Domain objects and relationships go in Atlas. Choices go in Decisions.
Public tutorials and product reference belong in root docs/. Mechanically bound
agreements remain Contracts and follow their amendment/version pipeline.

A domain skill does not activate a Mission by being invoked. Authorized delegation
can use a compact task envelope, disjoint writes, and an acceptance check. Existing
governed workers use their named Mission/Objective/Charter; others follow the task.

[Workspace navigation](workspace-navigation.md) ·
[Domain overview](domain-overview.md) · [Modular delivery](operating-levels-and-software-factory.md)
