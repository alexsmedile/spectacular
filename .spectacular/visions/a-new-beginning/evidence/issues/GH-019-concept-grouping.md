---
type: issue-evidence
source: source-011
github_issue: 19
url: https://github.com/alexsmedile/spectacular/issues/19
actionability: needs-scope
maps_to: [PZL-120, PZL-129, PZL-133]
retrieved: 2026-08-07
---

# GH-019 — Concept-grouped review

## Plain-language focus

Review coherent clusters of related information instead of repeatedly rebuilding context item by item.

## Problem and evidence

Strict one-at-a-time review fragments connected decisions and increases user fatigue.

## Proposed direction

Group by explicit metadata, dependency, or domain before considering fuzzy similarity; begin
with one concrete surface such as decision review.

## Relationships and collisions

#12 defines compact records; this issue arranges several. #30 applies the same behavior to a
whole GitHub backlog. The current PZL domains and maps are a live prototype.

## Actionability

Needs one named first surface and grouping key; then it becomes a small reversible change.

## Refactor relevance

Directly validates the puzzle-piece database and domain-map workflow used in this Vision.
