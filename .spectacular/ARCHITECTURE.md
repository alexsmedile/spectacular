---
type: Anchor
id: 019fe381-5d61-7223-b362-03a5f99a7b15
human_ref: ARCHITECTURE
title: Architecture
direction: Keep canonical Markdown authoritative and derive deterministic indexes and projections from it.
boundaries:
  - The Go domain kernel owns identity, lifecycle, authority, evidence, and refusal invariants.
  - The Skill owns judgment; the CLI owns deterministic validation and persistence.
constraints:
  - Caches and projections are disposable and non-authoritative.
freshness_checked_at: "2026-08-10T00:00:00Z"
freshness_source: .spectacular/workspace.yaml
freshness_source_fingerprint: d8b24fe7cfef0986a4b48e7f4e6dd8c7373b451d4c54bde425a904889539b4d3
freshness_valid_until: "2027-08-10T00:00:00Z"
---
# Architecture

One Go executable operates an inspectable Markdown workspace through explicit
filesystem, Git, runtime, and provider boundaries.
