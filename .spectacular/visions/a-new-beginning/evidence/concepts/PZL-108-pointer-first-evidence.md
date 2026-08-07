---
type: concept-piece
id: PZL-108
status: captured
domain: evidence-storage
sources: [source-008, source-009, source-010]
source_authority: unsourced-protocol-synthesis
assessment: strong
evidence_status: supported
disposition: pending
depends_on: [PZL-073]
overlaps_with: [PZL-075, PZL-091]
conflicts_with: []
tags: [artifacts, pointers, logs, context-budget]
updated: 2026-08-07
---

# Pointer-first externalized evidence

## Core message

Keep durable summaries and stable pointers in active state; store large logs, traces, and
tool output as retrievable artifacts outside the prompt.

## Value

Preserves inspectability and resume evidence without repeatedly loading high-volume output.

## Assumptions

- Artifact identity, retention, redaction, and access are defined.
- The summary never substitutes for required raw evidence.

## Evidence and collisions

Spectacular already externalizes records and verification logs. Source 008 overreaches by
implying all raw output should be Markdown: logs can contain secrets, personal data, huge
payloads, binary material, or transient machine paths.

## Trade-offs and recommendation

Pointers protect context budgets but can rot or conceal evidence. Standardize artifact
manifests and local/redacted storage classes before adding another committed corpus.

## Decision

Pending.
