---
type: Atlas
title: Atlas guide
version: "0.2"
created: "2026-08-23T18:24:40+02:00"
updated: "2026-10-07T18:09:01Z"
---

# Atlas: navigable durable context

Atlas holds domain objects and their relationships: entities, capabilities,
journeys, lifecycles, and maps. Diagrams are useful views of that knowledge,
not the only permitted content. Keep each file focused and link related objects.

## Choose the right home

- Anchors explain current project truth and constraints.
- Atlas explains objects, relationships, and how the system fits together.
- Plans describe intended changes and their verification.
- Decisions preserve choices and rationale; mechanical validation is recommended.
- Contracts preserve accepted agreements through mechanical validation and the
  amendment/version pipeline.
- Raw is an unconstrained scratchpad: no metadata, naming, or promotion duties.

Keep project plans out of product docs/. Documentation tools manage that surface.
Prefer ONTOLOGY.md for a new domain anchor; this workspace currently retains
[its existing domain anchor](../VOCABULARY.md) so historical references survive.

## Small, linked objects

Use descriptive types and basic metadata under the
[folder agreement](../../skills/spectacular/knowledge-folders.yaml).
Aim for one useful question or object per file, usually 50–120 lines. Split
independently useful knowledge, not every subsection. Relative Markdown links
and path-qualified wikilinks make relationships navigable.

Manual INDEX.md and generated index.json are valid entry points. Generated
inventories state their scope and never replace manually written navigation.
Atlas does not authorize implementation or acquire a Mission by being read.
Optional metadata/link diagnostics improve context quality without lifecycle gates.
