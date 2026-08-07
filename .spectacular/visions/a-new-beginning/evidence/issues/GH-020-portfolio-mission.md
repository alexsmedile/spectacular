---
type: issue-evidence
source: source-011
github_issue: 20
url: https://github.com/alexsmedile/spectacular/issues/20
actionability: needs-joint-grill
maps_to: [PZL-015, PZL-066, PZL-088, PZL-098, PZL-103, PZL-134]
retrieved: 2026-08-07
---

# GH-020 — Portfolio Mission above requests

## Plain-language focus

Let a confirmed goal become a cross-request node plan that the orchestrator navigates.

## Problem and evidence

Spectacular coordinates within requests but lacks a confirmed layer for sequencing several
requests toward one goal. Goal-scoped AFK records are relevant prior art.

## Proposed direction

Model checkpoints, discovery, specs, execution branches, human gates, and outstanding work;
decide whether a Mission is a new artifact, linked requests, or generalized goal/run state.

## Relationships and collisions

Grill jointly with #31 because executable concurrency constrains the data model. #24 waits on
the node model. It conflicts with sources that use Mission as the request replacement.

## Actionability

Not buildable. Decide artifact identity, relationship to PLAN/TASKS, and checkpoint surface first.

## Refactor relevance

This is a product-model fork: portfolio Mission versus one bounded Mission/request.
