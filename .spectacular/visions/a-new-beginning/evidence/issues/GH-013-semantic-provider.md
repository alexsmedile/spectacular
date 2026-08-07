---
type: issue-evidence
source: source-011
github_issue: 13
url: https://github.com/alexsmedile/spectacular/issues/13
actionability: deferred
maps_to: [PZL-090, PZL-106, PZL-121, PZL-124, PZL-130]
retrieved: 2026-08-07
---

# GH-013 — Optional semantic and code-graph provider

## Plain-language focus

Use optional semantic and structural retrieval to suggest related evidence, never to decide state.

## Problem and evidence

Exact lookup may miss terminology variants and multi-hop code impact, but no benchmark yet
shows a failure after the proposed deterministic retrieval improvements.

## Proposed direction

Local-first provider contract, provenance/freshness, bounded hybrid candidates, deterministic
code edges, optional sidecar, explicit cloud egress, and safe fallback to normal CLI reads.

## Relationships and collisions

Hard-sequenced after #11/#12. Converges with graph-memory sources while drawing a stronger
line: semantics cannot choose lifecycle, approval, policy, blockers, or verification.

## Actionability

Park until the deterministic baseline and representative retrieval failures exist.

## Refactor relevance

Defines a possible adapter boundary, not a core dependency or current implementation batch.
