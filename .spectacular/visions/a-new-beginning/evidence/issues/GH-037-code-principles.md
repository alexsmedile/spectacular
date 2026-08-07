---
type: issue-evidence
source: source-011
github_issue: 37
url: https://github.com/alexsmedile/spectacular/issues/37
actionability: needs-decision
maps_to: [PZL-010, PZL-111, PZL-125, PZL-143]
retrieved: 2026-08-07
---

# GH-037 — Code-quality principles guidance

## Plain-language focus

Give builders and reviewers a shared compact vocabulary for maintainability trade-offs.

## Problem and evidence

Project principles focus on Spectacular operations, while code-quality judgment is implicit.
Generic DRY/SOLID definitions are already model prior and can add prompt cost without behavior.

## Proposed direction

Prefer a short repository-specific review rubric tied to observed failure modes and checks;
link external generic definitions rather than reproducing a textbook glossary.

## Relationships and collisions

Directly supports anti-slop PZL-111/125 but collides with #11's context-reduction goal if
always loaded. Placement determines whether it is useful or inert.

## Actionability

Decide whether the need is vocabulary, enforceable rubric, or both; then choose canonical owner.

## Refactor relevance

An important warning that “more principles” is not equivalent to better agent behavior.
