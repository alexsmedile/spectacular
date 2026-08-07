---
type: issue-evidence
source: source-011
github_issue: 12
url: https://github.com/alexsmedile/spectacular/issues/12
actionability: coordination-blocked
maps_to: [PZL-002, PZL-010, PZL-053, PZL-055, PZL-123, PZL-129]
retrieved: 2026-08-07
---

# GH-012 — Soft-DB projections and CLI-first retrieval

## Plain-language focus

Let deterministic CLI projections select records and decisions before agents read bodies.

## Problem and evidence

Frontmatter is already the signal layer, but views and escalation rules vary across entity
types; derived indexes and permissive YAML parsing can silently drift.

## Proposed direction

Decision-capable list/brief/full projections, explicit drill-down, one compiled briefing,
small scalar filters, a strict indexed-frontmatter subset, deterministic rebuild, and caches
only after correctness.

## Relationships and collisions

Shares measurement and briefing schema with #11. #19 should influence multi-record grouping.
#13 must wait until this deterministic authority is stable.

## Actionability

Do not dispatch independently until #11/#12 own one briefing command and field schema.

## Refactor relevance

Provides the replacement contract required before deleting router instructions or collections.
