---
kit: research
version: 2.0
extends: prd
adds-slots:
  - name: Hypothesis
    after: Problem
    required: true
    prompt: |
      What do you currently believe is true? Frame falsifiable — there must be a possible answer that disproves it.
    example: |
      Users in cohort A spend 30% more time on task X than cohort B, due to onboarding-flow differences.
  - name: Method
    after: Deliverable
    required: true
    prompt: |
      How will you investigate? Specific enough that someone else could execute.
      Cover: approach, sources, tools, analysis.
    example: |
      Approach: comparative session analysis + 10 user interviews per cohort.
      Sources: 90 days of session data; recruit interviewees via in-app prompt.
      Tools: Mixpanel for sessions, NotebookLM for interview synthesis.
      Analysis: time-on-task comparison + thematic coding of interview transcripts.
  - name: Decision being informed
    after: Constraints
    required: true
    prompt: |
      What downstream decision does this research feed? Who owns it? When is it due?
      Research without a decision is hobby.
    example: |
      Decision: Whether to invest 6 weeks rebuilding the cohort-B onboarding flow.
      Decision owner: Head of Product.
      Decision deadline: 2026-07-15.
modifies-slots:
  - name: Deliverable
    note: |
      For research: name the artifact + location + length.
      Example: "Decision memo (~1,500 words) saved to _research/onboarding-cohorts/, with 5+ cited sources."
  - name: Goals & success criteria
    note: |
      For research: define the stop condition. "Done when ___" with measurable signal.
      Example: "Done when recommendation is published with confidence rating + 5 cited sources by 2026-06-15."
triggers-docs:
  always: []
  suggested:
    - decisions
description: |
  Investigations, experiments, market research, technical spikes, literature reviews.
  Adds Hypothesis + Method + Decision-being-informed slots.
---

# Research kit

For projects whose output is knowledge, not a shipped artifact.

Use when the project is:
- A technical spike to validate an approach
- A market or competitive analysis
- A user research investigation
- A literature review feeding an architectural decision
- An experiment with measurable hypothesis

Skip this kit when:
- The project ships code (→ `coding`) or content (→ `content`) as primary output
- There's no downstream decision the research feeds — without one, "research" becomes hobby
