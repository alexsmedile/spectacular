---
type: issue-evidence
source: source-011
github_issue: 53
url: https://github.com/alexsmedile/spectacular/issues/53
actionability: decision-ready
maps_to: [PZL-051, PZL-054, PZL-149]
retrieved: 2026-08-07
---

# GH-053 — Proposal-only capture surface

## Plain-language focus

Choose the smallest local capture interface that recommends a destination without publishing anything.

## Problem and evidence

SPC-007 is approved, but extending private IDEA, shipping guidance only, or coupling to one
provider have different product consequences.

## Proposed direction

Local `capture propose`: accept a redacted signal, render destination recommendation and a
copyable envelope, perform no network mutation, and let later adapters render provider drafts.

## Relationships and collisions

Strongly aligns with mechanical/agentic and native-provider boundaries. May become the handoff
surface to issue-triage or external capture companion skills.

## Actionability

Ready for owner choice among the recorded alternatives; implementation may then derive from SPC-007.

## Refactor relevance

Concrete thin slice for proposal-versus-publication authority.
