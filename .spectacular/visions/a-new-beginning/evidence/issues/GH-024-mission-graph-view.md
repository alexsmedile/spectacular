---
type: issue-evidence
source: source-011
github_issue: 24
url: https://github.com/alexsmedile/spectacular/issues/24
actionability: split-ready-blocked
maps_to: [PZL-103, PZL-120, PZL-132, PZL-135]
retrieved: 2026-08-07
---

# GH-024 — Stateful Mission graph view

## Plain-language focus

Give humans a compact supervision view over cross-request nodes, dependencies, and state.

## Problem and evidence

Wide agent fan-out is unreadable as concurrent logs or prose; the Mission model to render is
not yet defined.

## Proposed direction

ASCII-first stateful graph with extensible `pending/claimed/in-progress/blocked/done` states;
share emission with #18 while allowing linear and graph layout strategies.

## Relationships and collisions

Renderer design can follow #18 now. Binding and checkpoint surface wait on #20/#31.

## Actionability

Split the renderer experiment from the blocked Mission-state integration.

## Refactor relevance

Tests visual supervision without using a picture as execution authority.
