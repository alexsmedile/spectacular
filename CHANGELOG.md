# Changelog

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
