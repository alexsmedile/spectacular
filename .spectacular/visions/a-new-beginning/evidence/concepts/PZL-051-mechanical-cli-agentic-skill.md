---
type: concept-piece
id: PZL-051
status: captured
domain: execution-authority
sources: [source-005, source-006, source-007]
source_authority: code-audit-proposal
assessment: promising
evidence_status: partial
disposition: pending
depends_on: [PZL-038]
overlaps_with: [PZL-001, PZL-002, PZL-058, PZL-069, PZL-091]
conflicts_with: [PZL-071, PZL-098]
tags: [cli, skill, mechanical, agentic, authority]
updated: 2026-08-07
---

# Mechanical CLI, agentic skill boundary

## Core message

Shell commands should perform deterministic mechanical operations; judgment,
interpretation, and adaptive work should live under the Spectacular skill.

## Value

Gives every entry point a predictable executor, consent model, and test strategy.

## Assumptions

- Hybrid commands can be separated without harming useful shorthand.
- Skill availability is acceptable for every agentic operation.

## Evidence and collisions

Current remember, imagine, status, and spec flows change executor or meaning based
on arguments. Some CLI commands intentionally emit next-step recipes, blurring the line.
Source 006 directly conflicts by placing agentic planning, autonomous execution,
and judgment-heavy closure behind `spectacular mission plan|run|close`.
Source 007 deepens the conflict through a Mission compiler and staged agentic run chain.

## Trade-offs and recommendation

Clarity and testability versus migration and possible loss of ergonomic hybrids.
Build a command-by-command authority matrix before accepting the universal rule.

## Decision

Pending.
