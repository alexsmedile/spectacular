---
type: Anchor
id: 019fe381-5d61-7223-b362-03a5f99a7b16
human_ref: STACK
title: Stack
direction: Ship native governed tooling and modular Skill guidance while keeping ordinary
  context runtime-independent.
boundaries:
- macOS and Linux on amd64 and arm64 are release binary targets.
- Ordinary Markdown context requires no language runtime.
- Optional diagnostics and repository developer utilities declare their dependencies.
constraints:
- Canonical content is UTF-8 Markdown with YAML frontmatter.
freshness_checked_at: '2026-08-10T00:00:00Z'
freshness_source: .spectacular/workspace.yaml
freshness_source_fingerprint: d8b24fe7cfef0986a4b48e7f4e6dd8c7373b451d4c54bde425a904889539b4d3
freshness_valid_until: '2027-08-10T00:00:00Z'
version: '0.1'
updated: '2026-10-07T20:05:13Z'
created: '2026-05-11T14:46:49+02:00'
---
# Stack

Go implements the governed CLI; Git preserves revision history. Markdown with YAML
frontmatter stores typed context. Progressive Skill references guide ordinary work
and optional governance, with bundled architecture, modeling, and prototyping Skills.

Ordinary context needs no CLI, Go, Python, or workspace manifest. Distributed native
CLI archives cover macOS/Linux on amd64/arm64. Optional context diagnostics require
Python 3.9+ and PyYAML. Repository index rebuilding requires Go and uses existing
Discovery and JSON projection code; it does not add a public CLI command.

[Workspace configuration and navigation](../skills/spectacular/references/navigation.md)
explains the manifest, optional overrides, and generated inventory.
