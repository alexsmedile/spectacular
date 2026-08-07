---
type: issue-evidence
source: source-011
github_issue: 35
url: https://github.com/alexsmedile/spectacular/issues/35
actionability: architecture-dependent
maps_to: [PZL-006, PZL-007, PZL-008, PZL-141]
retrieved: 2026-08-07
---

# GH-035 — Repository collaboration scaffold

## Plain-language focus

Offer standard branch, PR, issue-template, and community conventions for collaborative repositories.

## Problem and evidence

The repo lacks `.github/` templates and branch guidance, but this does not prove every initialized
Spectacular project needs them.

## Proposed direction

Conventional branch categories, PR/issue templates, issue linking, optional CODEOWNERS, and
review-size guidance.

## Relationships and collisions

Conflicts with the earned-scaffold proposal if made default. It may fit a GitHub/coding kit,
convention pack, or companion skill better than core init.

## Actionability

Decide owner and opt-in trigger before implementation; the artifacts themselves are straightforward.

## Refactor relevance

Useful test case for “grow/offer structure when the external collaboration surface exists.”
