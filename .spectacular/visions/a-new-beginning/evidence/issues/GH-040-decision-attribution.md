---
type: issue-evidence
source: source-011
github_issue: 40
url: https://github.com/alexsmedile/spectacular/issues/40
actionability: decision-ready
maps_to: [PZL-086, PZL-115, PZL-146]
retrieved: 2026-08-07
---

# GH-040 — Decision accountability attribution

## Plain-language focus

Distinguish who authorized a decision from where its proposal or evidence originated.

## Problem and evidence

`origin` and `derived_from` preserve provenance but not accountable approval; multi-human
workspaces make that distinction useful.

## Proposed direction

Optional human/role deciders or approvers, with agent authorship kept as provenance rather
than implying an agent holds product authority.

## Relationships and collisions

Converges with provenance and human-gate concepts. A free-text person list creates privacy,
identity, and staleness concerns that the issue does not resolve.

## Actionability

Decide whether attribution is needed now and define human/role versus agent semantics first.

## Refactor relevance

Important if multi-agent output expands, but low value for a solo portable minimum.
