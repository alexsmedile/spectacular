---
id: DEC-016
type: decision
status: verified
origin: user-derived
derived_from: "User confirmations in request-workflow interview, 2026-08-02"
supersedes: ""
evidence: "user direction"
tags: [decisions, indexing, retention, docs, context]
updated: 2026-08-02
---

# DEC-016 — Compact decision history with rolling fifty, ten, and fifty-entry tiers after the first fifty decisions

**Context:**
Decision history must remain quickly retrievable without forcing agents to load every ADR body as the corpus grows.

**Evidence:**
User direction recorded by explicit command.

**Alternatives considered:**
Recorded in the cited user context or decision discussion.

**Decision:**
Compact decision history with rolling fifty, ten, and fifty-entry tiers after the first fifty decisions.

**Consequences:**
Keep the newest 50 decisions individual; compact the preceding 50 into blocks of 10; compact older history into blocks of 50. Strong frontmatter and index summaries preserve block discovery and selective retrieval. #decisions #indexing #retention #docs #context
