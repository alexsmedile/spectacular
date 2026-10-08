---
type: Anchor
id: 019fe381-5d61-7223-b362-03a5f99a7b15
human_ref: ARCHITECTURE
title: Architecture
direction: Keep canonical Markdown authoritative and derive deterministic indexes
  and projections from it.
boundaries:
- The Go domain kernel owns identity, lifecycle, authority, evidence, and refusal
  invariants.
- The Skill owns judgment; the CLI owns deterministic validation and persistence.
constraints:
- Caches and projections are disposable and non-authoritative.
freshness_checked_at: '2026-08-10T00:00:00Z'
freshness_source: .spectacular/workspace.yaml
freshness_source_fingerprint: d8b24fe7cfef0986a4b48e7f4e6dd8c7373b451d4c54bde425a904889539b4d3
freshness_valid_until: '2027-08-10T00:00:00Z'
version: '0.1'
updated: '2026-10-07T20:05:13Z'
created: '2026-05-21T18:05:07+02:00'
---
# Architecture

Canonical Markdown remains understandable directly. The Skill owns context
placement, judgment, and optional drafting methods. The Go kernel owns governed
identity, lifecycle, authority, proof, validation, and recoverable persistence.

Manual INDEX.md owns curated reading routes. Generated index.json inventories
only governed records and can be rebuilt; it is not project truth. Optional
collection JSON files are filtered projections of that same graph, not another
catalog. Cache files never supply authority.

workspace.yaml locates governed records; config.yaml is an optional override.
Neither file is needed for ordinary Markdown work. Root docs/ remains a separate
human-facing documentation surface. See the [navigation reference](../skills/spectacular/references/navigation.md).
