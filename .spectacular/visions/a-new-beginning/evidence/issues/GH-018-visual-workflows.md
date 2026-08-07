---
type: issue-evidence
source: source-011
github_issue: 18
url: https://github.com/alexsmedile/spectacular/issues/18
actionability: implementation-candidate
maps_to: [PZL-120, PZL-132]
retrieved: 2026-08-07
---

# GH-018 — Visual workflow and schema proposals

## Plain-language focus

Show proposed structure visually so users can review relationships faster than prose alone.

## Problem and evidence

The owner reports lower comprehension and higher fatigue from long workflow/schema prose.

## Proposed direction

Start with PLAN-oriented ASCII diagrams, keep an internal model capable of later Mermaid
emission and state annotation, and retain prose for nuance rather than drawing everything.

## Relationships and collisions

Renderer prototype for #24. The intake's Mermaid maps already provide evidence for the
derived-view approach but do not validate a product renderer.

## Actionability

Narrow prototype is ready; choose one PLAN fixture and comprehension/maintenance measures.

## Refactor relevance

Potential progressive-disclosure layer, not a reason to add diagrams to every artifact.
