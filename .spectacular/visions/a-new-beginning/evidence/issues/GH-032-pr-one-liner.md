---
type: issue-evidence
source: source-011
github_issue: 32
url: https://github.com/alexsmedile/spectacular/issues/32
actionability: decision-ready
maps_to: [PZL-131, PZL-136]
retrieved: 2026-08-07
---

# GH-032 — Plain-language PR lead

## Plain-language focus

Open every Spectacular-created PR with one sentence explaining what it is about.

## Problem and evidence

The current structured manifest is complete but makes readers encounter bookkeeping before meaning.

## Proposed direction

Prepend a reader-oriented lead distinct from change bullets and durable goal; show it in dry-run.

## Relationships and collisions

Shares the PR-manifest implementation surface with #17. Goal, intent, change summary, and
reader lead should be defined once rather than derived ambiguously from one field.

## Actionability

Decision-ready as a joint #17/#32 packet: choose explicit input versus safe derivation and gate level.

## Refactor relevance

An example of typed handoff fields serving different audiences rather than one overloaded summary.
