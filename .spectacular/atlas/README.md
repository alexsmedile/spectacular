---
type: Atlas
title: Atlas guide
version: '0.4'
created: '2026-08-23T18:24:40+02:00'
updated: '2026-10-07T20:07:46Z'
---

# Atlas: navigable durable context

Atlas holds domain objects and relationships: entities, capabilities, journeys,
lifecycles, and maps. Keep each file focused and link related objects.

## Entry points

- [Workspace navigation](workspace-navigation.md): requirements, plans, specs, audits, and
  governed record ownership, with links to current work.
- [Domain overview](domain-overview.md): ordinary context and optional governance.
- [Specification lifecycle](specification-and-governance-lifecycle.md): how a
  draft, plan, spec, or Contract changes purpose.
- [Governance choice](governance-tiers-and-cadence.md): explicit selection.
- [Session continuity](session-orchestration-and-lifecycle.md): resume and dispatch.

Anchors define current project truth. Plans own intended work and results;
Atlas explains lasting relationships. Decisions preserve choices. Contracts keep
their mechanical validation and amendment/version pipeline. Raw is unconstrained.

ONTOLOGY.md is this project’s domain Anchor; preserve an existing ONTOLOGY.md
in projects that already use it. Product docs/ stays outside this context store.

## Small, linked objects

Follow the [folder agreement](../../skills/spectacular/knowledge-folders.yaml).
Aim for one useful object or question per file, usually 50–120 lines. Split by
independent meaning, domain, code module, phase, or delivered tranche. Link
reusable context instead of copying it into each plan.

Manual INDEX.md and generated index.json are optional entry points. Generated
inventories state their scope and preserve authored navigation. Atlas does not
activate a Mission or authorize implementation merely by being read.
