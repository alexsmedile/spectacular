---
type: issue-evidence
source: source-011
github_issue: 17
url: https://github.com/alexsmedile/spectacular/issues/17
actionability: decision-ready
maps_to: [PZL-066, PZL-068, PZL-119, PZL-131]
retrieved: 2026-08-07
---

# GH-017 — Durable goal propagation

## Plain-language focus

Requests and PRs should preserve why a change exists, not only summarize the diff.

## Problem and evidence

PLAN already has Goal and several creation paths already accept outcome summaries, but
presence is not consistently checked and the PR manifest does not reliably carry the why.

## Proposed direction

Mechanically validate non-empty Goal, review its meaning judgmentally, propagate it to PR
intent with an override, and later compare delivered evidence with the stated outcome.

## Relationships and collisions

Shares the PR manifest with #32; implement together or sequence deliberately. #20 consumes
the goal. Goal text and PR-reader one-liners serve related but different audiences.

## Actionability

Ready after replacing the untestable word “meaningful” with mechanical presence/length plus review.

## Refactor relevance

Strengthens the Mission contract regardless of final vocabulary or execution ownership.
