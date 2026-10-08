---
type: Atlas
title: Workspace navigation and durable planning
version: '0.4'
created: '2026-10-07T18:42:01Z'
updated: '2026-10-07T20:07:46Z'
---

# Workspace navigation

This workspace keeps the project’s durable context and existing governed history.
For a realistic active-project example, see the [container guide](../../skills/spectacular/references/workspace.md).

## Start from the question

| Need | Actual home in this repository |
|---|---|
| Current direction | [PROJECT.md](../PROJECT.md) |
| Domain meanings | [ONTOLOGY.md](../ONTOLOGY.md) |
| Planning this change | [workspace-layout-cleanup.md](../plans/workspace-layout-cleanup.md) |
| Earlier implementation and results | [knowledge-first-optional-governance.md](../archive/plans/knowledge-first-optional-governance.md) |
| Lasting relationships | [domain-overview.md](domain-overview.md) |
| Accepted mechanical agreements | contracts/CC-*.md; existing version/amendment rules |
| Formal historical execution | archive/missions/M*-*/M*-*.md and Mission-owned records |
| Public product guidance | [docs/README.md](../../docs/README.md) |

## Metadata profiles

[Profiles](../../skills/spectacular/references/profiles.md) scale optional detail:
minimal uses four basic fields; compact adds a recommended description and is the
soft default; extended adds useful retrieval fields. Raw stays freeform. Governed
schemas are a separate choice. No new file or profile label is required.

## Requirement, plan, spec, Contract

Independently useful needs and acceptance criteria can live in requirements/.
Existing Requirement files in specs/ remain valid; keep one linked owner. See
[requirements guidance](../../skills/spectacular/references/requirements.md) and
[optional drafting methods](../../skills/spectacular/references/drafting-methods.md).
A plan answers what to do and tracks results. Reusable behavior can be extracted
to an optional specs/ file; its type is Spec and it has basic document metadata.
An agreed mechanically bound capability belongs in a Contract, whose supported
amendment/version process applies. Length alone does not change a file’s role.

Use `/spectacular plan` as a Skill route, or `$spectacular plan` in hosts using
that invocation. Native host plan mode owns restrictions and approvals. The
project copy goes in .spectacular/plans/ when writing is permitted; a restricted
session returns an unsaved draft and intended path. No binary plan command exists.

## Audit, Review, Handoff

Ordinary inspection findings can stay in the plan. An independently useful audit
can use optional audits/ with type Audit, inspected revision, scope, findings,
basis, and limitations. A governed Review, Evidence, or Handoff follows the path
reported by its command, normally inside its owning Mission bundle.

Create requirements/, specs/, or audits/ only for useful content. This repository
has one retained Requirement and no need for empty specs/ or audits/. Frozen
execution lives under archive/; manual INDEX.md links current context while generated
index.json inventories governed records. Configuration and filename conventions
live in the [navigation reference](../../skills/spectacular/references/navigation.md).
