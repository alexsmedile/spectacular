---
type: Decision
id: 01a059d3-3a60-7000-8000-b53d8f8a1921
title: Codify Banned Synonyms Invariant, Architectural Pattern Pass, and Post-Mission Guardrail Feedback Loop
created_by: Alex
created: "2026-08-31T19:35:00Z"
updated: "2026-08-31T19:35:00Z"
actor: Alex
actor_role: owner
alternatives:
    - leave domain vocabulary as a descriptive glossary without banned synonyms
    - allow agents to implement bespoke architectures without an upfront pattern survey
    - create loose ad-hoc LEARNINGS.md files instead of routing invariants to permanent anchors
disposition: accepted
question: How should Spectacular eliminate LLM synonym drift, prevent bespoke wheel reinvention, and systematically codify post-mission failure lessons?
rationale: Drawing from senior AI systems engineering discipline (Chem's 14-step APIV framework), Spectacular strengthens its upfront alignment and closed-loop learning. VOCABULARY.md formalizes explicit Banned Synonyms to prevent state-machine and action drift. Planning incorporates a mandatory Pattern Pass (surveying stdlib idioms, battle-tested OSS, and RFC standards). Mission completion institutes a Mistake Tax routing failure lessons directly to GUARDRAILS.md (domain invariants) and AGENTS.md (tooling rules) without creating redundant governance records.
ref: D29-banned-synonyms-pattern-pass-and-guardrail-feedback-loop
scope:
    - v2
---
# Codify Banned Synonyms Invariant, Architectural Pattern Pass, and Post-Mission Guardrail Feedback Loop

This Decision formalizes three operational disciplines derived from senior-level specification-driven AI engineering:

## 1. Domain Ontology & Banned Synonyms (`VOCABULARY.md`)
- `VOCABULARY.md` (canonical under D25) must define an explicit **Permitted Actions & Banned Synonyms** table and **Permitted Entity States** enumeration.
- Banned synonyms (e.g. using `close`, `finish`, or `resolve` instead of `CompleteMission`) are treated as quality-gate violations in agent prompt compilation and code reviews to eliminate prompt ambiguity and state fragmentation across fresh context windows.

## 2. The Architectural Pattern Pass (Upfront in `prepare.md`)
- Before drafting custom implementations for complex claims or new components, the Orchestrator must execute an upfront **Pattern Pass**:
  - **Track A (Fast Parametric Survey)**: Survey standard library idioms and existing codebase patterns for routine problems.
  - **Track B (Subagent Research Dispatch)**: Dispatch a transient `research` subagent to survey established open-source reference implementations, RFC standards, or battle-tested libraries for novel/complex domains.
- A 3-line **Pattern Census** (Candidate Patterns, Selected Reference Pattern, Rejected Anti-Patterns) is frozen into the Mission body to defend the design against bespoke wheel reinvention ("AI slop").

## 3. Post-Mission Mistake Loop & Guardrail Feedback Loop (`close.md`)
- When a Mission closes having consumed repair budgets, hit tricky regressions, or resolved non-trivial review findings, the Orchestrator must pay the **"Mistake Tax"**:
  - Root-cause domain/system/code failure modes must be codified into permanent, categorized invariants in `.spectacular/GUARDRAILS.md`.
  - Workflow, harness, or command execution errors must be codified into `AGENTS.md`.
  - Architectural trade-offs are recorded as atomic Decisions (`.spectacular/decisions/`).
- Ad-hoc `LEARNINGS.md` or `FIXES.md` records remain forbidden to preserve Spectacular's single-source-of-truth invariants.
