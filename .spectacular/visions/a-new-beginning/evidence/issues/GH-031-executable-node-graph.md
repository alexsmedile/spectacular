---
type: issue-evidence
source: source-011
github_issue: 31
url: https://github.com/alexsmedile/spectacular/issues/31
actionability: needs-joint-grill
maps_to: [PZL-103, PZL-104, PZL-118, PZL-119, PZL-139]
retrieved: 2026-08-07
---

# GH-031 — Executable node graph

## Plain-language focus

Make nodes claimable, dependency-gated units of concurrent work rather than only plan entries.

## Problem and evidence

#20 plans and #24 draws, but neither defines node claiming, concurrency safety, state,
dependency enforcement, or result reintegration.

## Proposed direction

Pending→claimed→running→blocked→done nodes, executable edges, parallel eligibility,
collision safety, result packets, and a human supervision projection.

## Relationships and collisions

Must be grilled with #20. Closed #3 is single-request prior art. Directly conflicts with the
single-agent MVP and may belong in a host orchestrator rather than Spectacular.

## Actionability

Architecture exploration only; no execution substrate is authorized.

## Refactor relevance

Clarifies the strongest graph world so it can be compared honestly with the minimal world.
