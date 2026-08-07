---
type: synthesis-checkpoint
checkpoint: 006
sources: [source-001, source-002, source-003, source-004, source-005, source-006, source-007, source-008, source-009, source-010]
concepts: 127
human_dispositions: 0
updated: 2026-08-07
---

# Synthesis checkpoint 006

## What Source 010 adds

Source 010 reframes abstraction as context-budget architecture. Its useful contribution is
not another mandate for a multi-agent Diamond; it is an audit method for deciding what each
Spectacular layer should know, own, expose, validate, and hide.

The proposed audit layers are:

1. intent and accepted outcomes;
2. agentic orchestration and bounded context;
3. deterministic CLI mechanisms and gates;
4. canonical truth, evidence, and resume state;
5. external runtimes, providers, databases, and companion skills.

The experimental `abstraction-map.md` makes those boundaries reviewable without accepting
them as the final architecture.

## Corrections that protect the anti-slop goal

### Encapsulation is not absolution

A narrow interface is valuable, but agent-written internals still need security, cohesion,
observability, performance, testability, and maintenance evidence. “Messy inside the box”
would institutionalize exactly the hidden debt this refactor is intended to remove.

### Isolation is not ignorance

An orchestrator needs enough system and risk context to decompose safely. A specialist needs
the local purpose, integration contract, authority, and stop conditions—not only procedural
instructions. Context should be role-bounded, never purpose-free.

### Compression is not deletion

A handoff should omit conversation and irrelevant failed trials, but retain negative evidence
that prevents repetition or changes the next hypothesis. A 300-token number is suitable as an
experiment target, not a completeness rule.

### Review requires a terminal state

A review-repair loop must end as accepted, accepted-risk, blocked, or failed. Subjective
Skeptic satisfaction cannot be the gate. Findings need contract linkage, severity,
deduplication, retry budget, and an escalation packet.

## Concrete refactor enhancements

1. Run the abstraction audit before choosing new folders or companion skills.
2. Classify every current agent and command by one owning layer; record deliberate crossings.
3. Specify one typed handoff capsule and test it on cold resume and independent review.
4. Measure a soft handoff budget for completeness rather than adopting 300 tokens blindly.
5. Add internal-quality evidence to any grey-box implementation contract.
6. Compose deterministic gates, independent review, and hypothesis retry into one bounded
   state machine rather than three overlapping workflow systems.
7. Keep semantic database adapters outside the portable core until a real query abstraction
   requirement appears.
8. Use the abstraction map to decide companion-skill boundaries: a companion should own a
   distinct job, state substrate, and external interface—not merely reduce SKILL.md size.

## Companion-skill implication

The abstraction audit gives a sharper extraction test for bugworks, specwright, verifyworks,
Wayfinder, and AFK candidates:

- Does the candidate own a complete responsibility behind a small stable interface?
- Can it operate from a typed handoff without importing Spectacular's internal state model?
- Does it own its evidence and lifecycle without duplicating canonical project truth?
- Can Spectacular reference its result as a file or schema contract rather than invoking its
  internal workflow?
- Does standalone use remain coherent?

An extraction that fails these tests is code movement, not an abstraction boundary.

## Recommended decision sequence

1. Decide execution authority and portable-core boundary.
2. Approve or revise the five-layer abstraction audit model.
3. Classify the existing surface and candidate companion skills by owning layer.
4. Define typed handoff and resume contracts.
5. Define bounded verification and retry termination.
6. Only then choose skill extraction, CLI modularization, graph topology, or storage technology.

No abstraction map, agent hierarchy, extraction, or review loop is approved by this checkpoint.
