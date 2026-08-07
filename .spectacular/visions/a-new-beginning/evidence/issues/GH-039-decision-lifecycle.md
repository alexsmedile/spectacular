---
type: issue-evidence
source: source-011
github_issue: 39
url: https://github.com/alexsmedile/spectacular/issues/39
actionability: decision-ready
maps_to: [PZL-076, PZL-100, PZL-145]
retrieved: 2026-08-07
---

# GH-039 — Decision lifecycle vocabulary

## Plain-language focus

Represent rejected and retired decisions without confusing unsettled questions with decisions.

## Problem and evidence

Current records use verified/superseded but cannot clearly mark rejection or deprecation;
`proposed` may violate the QUE/DEC authority split.

## Proposed direction

Keep unresolved alternatives as questions; decide whether accepted is implicit; add rejected
and deprecated only if they improve retrieval beyond alternatives inside an accepted ADR.

## Relationships and collisions

Expanding every entity lifecycle conflicts with the compact Mission/taxonomy goal and needs
usage evidence, not standards imitation alone.

## Actionability

Ready as a small taxonomy decision before any CLI flags or migration.

## Refactor relevance

Tests whether lifecycle precision earns its state cost.
