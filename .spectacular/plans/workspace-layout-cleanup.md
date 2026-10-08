---
type: Plan
version: "0.3"
created: "2026-10-07T20:05:13Z"
updated: "2026-10-08T22:16:01.489003Z"
description: Align ontology naming, reading routes, configuration roles, and archive placement.
---

# Workspace layout cleanup

The owner selected ONTOLOGY.md, documented naming/reading conventions, a single
retirement home, JSON generated inventories, and cleanup of stale current context.
D31 records archival authorization and supersedes the historical naming choice.

## Changes

- Rename the existing ontology Anchor, preserving its UUID and incoming soft links.
- Refresh PROJECT, PRODUCT, ARCHITECTURE, STACK, ROADMAP, and manual navigation.
- Retire 20 completed Missions and seven resolved Proposals with authorization,
  resolver where applicable, and canonical source fingerprints. Preserve frozen bindings.
- Retain completed campaign inputs, delivered plans, and byte-preserved legacy
  returns under archive/. P5, P10, and P12 remain open pending scope resolution.
- Remove generated Markdown navigation, the obsolete catalog cache, and redundant
  config defaults. Keep workspace.yaml for governed discovery.
- Store filename, reading-route, config, inventory, and retirement guidance in the
  Skill. Add a repository script for root governed index.json, no new CLI command.

## Verification

- All 123 governed identities remain discoverable. Twenty moved Mission roots
  preserve their original fields/body except added archival provenance; 39 owned
  records and the legacy return packet remain byte-identical.
- All 22 archived Missions and 13 Proposals validate. Fixed the layout validator
  to accept canonical terminal archives while refusing live Missions in archive/.
- Fixed preflight for an absent live Mission container; its explicit all-Mission
  sweep includes archive/missions/. Adjusted the symlink test to create its own target.
- Preflight and quick verification pass. Seven context tests, the index utility
  test, Skill validation, mutable Markdown links, and staged secret scan pass.
- Root index rebuild is deterministic; generated-only writes are idempotent,
  preserve manual INDEX.md, and refuse authored JSON.

The cleanup is included in the 3.2.0 release scope.
