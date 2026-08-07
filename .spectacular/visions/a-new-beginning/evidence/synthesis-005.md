---
type: synthesis-checkpoint
checkpoint: 005
sources: [source-001, source-002, source-003, source-004, source-005, source-006, source-007, source-008, source-009]
concepts: 122
human_dispositions: 0
updated: 2026-08-07
---

# Synthesis checkpoint 005

## What Source 009 clarifies

Source 009 is an expanded explanation of Source 008, not a second vote for the same
architecture. It clarifies four implementable hypotheses:

1. graph edges can be typed handoff packets rather than implicit prompt inheritance;
2. a Mermaid map can compress the decision frontier above atomic Markdown cards;
3. agent evaluation should hold the task constant and treat model plus harness as a system;
4. semantic meaning, logical inference, and closed constraint validation need distinct ownership.

The remaining material strengthens existing cards without increasing concept count.

## How to use these principles in this refactor

### 1. Add a visual projection, not a visual authority

The experimental `decision-map.md` compresses 122 concepts to 16 decision-critical nodes.
Its stable PZL identifiers resolve into the full Markdown cards. This tests whether visual
progressive disclosure improves orientation while preserving provenance.

Success criteria:

- a cold reviewer can identify the next upstream decision without reading the index;
- every shown relationship resolves to card metadata or a recorded collision;
- the view stays below roughly 20 nodes;
- updating a card cannot silently leave the view authoritative but stale.

### 2. Replace agent labels with typed handoffs

Before assigning `planner`, `architect`, `skeptic`, or `reviewer`, define the packet crossing
the boundary: objective, accepted sources, exclusions, evidence, output schema, open questions,
and confidence limits. The role name alone creates no independence or correctness.

### 3. Benchmark the refactor harness

Extend the current cold-start test into a small fixed evaluation deck:

| Task | What it measures |
|---|---|
| Orient and propose a bounded dark-mode change | retrieval volume, route correctness, time to first code action |
| Resume an interrupted Mission from disk | state completeness, baseline verification, safe next action |
| Diagnose a known contract violation | evidence retrieval, hypothesis quality, retry discipline |

Compare current routing, lean routing, and any graph-linked projection with the same model
before changing models or adding a fleet. Record correctness, loaded context, calls, retries,
latency, and cost.

### 4. Earn semantic machinery

Start with stable IDs, frontmatter types, labeled Markdown relationships, and validators.
Only test RDF/OWL/SHACL if a named multi-hop, interoperability, inference, or closed-validation
requirement cannot be expressed safely in that simpler layer.

### 5. Apply deep modules at the correct layer

Atomic concept cards are intentionally small database records. Runtime code should instead
be grouped by cohesive responsibility behind narrow interfaces. “Fewer files” and “deep
modules” are not interchangeable; the refactor must measure coupling and change locality.

### 6. Prototype disputed experiences

Use disposable prototypes for unresolved interactive behavior such as cold-start briefing,
Mission planning, or evidence review. Record what the prototype proves and fails to prove,
then translate the result into a contract. Prototype code never becomes the specification.

### 7. Shift left only on irreversible architecture

The owner should decide execution authority, truth authority, portable-core boundaries,
data/state ownership, and compatibility posture before implementation. File layout and exact
edits should remain derived after code reconnaissance.

## Recommended next review packet

The next decision session should test the experimental decision map and answer one upstream
question only:

> Does Spectacular own autonomous execution, persist a contract for host-owned execution,
> or remain a lifecycle control layer?

That choice determines whether Diamond orchestration, typed node packets, run state, and
independent verifier agents belong inside Spectacular, in companion skills, or in the host.

No product architecture or new runtime has been accepted here. The Mermaid view is a
reversible intake-workflow experiment.
