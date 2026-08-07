---
type: issue-evidence
source: source-011
github_issue: 36
url: https://github.com/alexsmedile/spectacular/issues/36
actionability: architecture-dependent
maps_to: [PZL-054, PZL-110, PZL-142]
retrieved: 2026-08-07
---

# GH-036 — Conventional-commit enforcement

## Plain-language focus

Turn a documented commit convention into an optional deterministic repository check.

## Problem and evidence

Commit guidance is advisory and recent history was reported as partially noncompliant; CI and
SemVer mapping are absent.

## Proposed direction

Validate PR titles or commits in CI and document type-to-version semantics without silently
automating releases.

## Relationships and collisions

Shares `.github/` ownership with #35. External CI is a provider/repository concern, not
necessarily Spectacular's portable core.

## Actionability

Choose PR-title versus every-commit enforcement, opt-in location, and release boundary first.

## Refactor relevance

Strong deterministic-harness example that may belong in a coding convention pack.
