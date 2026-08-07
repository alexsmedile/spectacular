---
type: concept-piece
id: PZL-052
status: captured
domain: cli-surface
sources: [source-005]
source_authority: code-audit-proposal
assessment: strong
evidence_status: supported
disposition: pending
depends_on: []
overlaps_with: [PZL-058, PZL-060]
conflicts_with: []
tags: [aliases, feedback, cli, deprecation]
updated: 2026-08-07
---

# Remove surprising aliases

## Core message

Remove generic aliases such as `iterate`, `test`, `probe`, and `try` that silently
dispatch to feedback, and deprecate the redundant `feedback-loop` spelling.

## Value

Avoids collisions with ordinary developer vocabulary and shrinks memorized grammar.

## Assumptions

- Alias usage is low or migratable.
- Feedback remains directly discoverable under a clear noun.

## Evidence and collisions

All named aliases exist in the dispatcher and share feedback behavior. Usage and
external automation have not yet been measured.

## Trade-offs and recommendation

Clearer semantics versus compatibility churn. Strong removal candidate with usage
search, warnings, and exact replacements.

## Decision

Pending.
