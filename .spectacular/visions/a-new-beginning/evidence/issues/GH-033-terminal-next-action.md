---
type: issue-evidence
source: source-011
github_issue: 33
url: https://github.com/alexsmedile/spectacular/issues/33
actionability: decision-ready
maps_to: [PZL-075, PZL-118, PZL-140]
retrieved: 2026-08-07
---

# GH-033 — Terminal next-action contract

## Plain-language focus

End true completion or blocking reports with one concrete continuation or unblocker.

## Problem and evidence

Spectacular can compute next actions but users repeatedly ask “now what?” after gates and work.

## Proposed direction

For terminal output, name one runnable action or explicit decision; blocked output names the
unblocker. Do not turn normal checkpoints into premature stopping or autonomous mutation.

## Relationships and collisions

Reuses status/wayfinding primitives but must respect the execution state machine: ongoing
approved work continues instead of ending merely to announce a next step.

## Actionability

Ready for a behavior-contract decision before choosing skill versus CLI implementation.

## Refactor relevance

Strengthens cold resume and handoff without preserving redundant `next` commands automatically.
