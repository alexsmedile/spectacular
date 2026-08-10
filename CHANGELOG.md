# Changelog

## 2.0.0-rc.2 — 2026-08-10

Human-operability correction for the v2 release candidate.

- Replaces flat UUID filenames with named project Anchors and cohesive Mission
  bundles while retaining UUIDv7 identity and SHA-256 revision fingerprints.
- Adds scoped human references for Missions, Objectives, Runs, Checkpoints,
  Evidence, Decisions, Gaps, Handoffs, and Assessments.
- Makes default CLI cards human-first and keeps exact machine data in `--json`.
- Commits deterministic, non-authoritative workspace and Mission indexes.
- Adds atomic whole-bundle Mission archival and a real self-hosted workspace.

RC.1 is superseded because its machine-oriented workspace representation did
not satisfy the human-comprehension contract.

## 2.0.0-rc.1 — 2026-08-10

First externally consumable Spectacular v2 release candidate.

- Breaking: removes the v1 public product surface and compatibility paths; Git
  recovery pointers preserve the frozen history.
- Adds pointer-first retrieval and a governed Mission loop with explicit
  authority, Evidence, assessment, reconciliation, and closure boundaries.
- Ships native four-platform archives and a checksum-verifying installer.

Known limitations: v1 workspaces are unsupported; there is no migration or
compatibility reader; discovery is pointer-driven rather than broad search;
and release publication remains outside the local CLI.
