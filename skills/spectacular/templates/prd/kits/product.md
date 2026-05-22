---
kit: product
version: 2.0
extends: prd
adds-slots:
  - name: User stories
    after: Target users
    required: true
    prompt: |
      List 2-5 user stories in "As a <user>, I can <action>, so that <outcome>" format.
      Each one is one sentence. Focus on the highest-value actions, not edge cases.
    example: |
      - As a freelance designer, I can send client updates, so that I stop context-switching to email.
      - As a freelance designer, I can track project status, so that I see what's overdue at a glance.
  - name: Metrics
    after: Goals & success criteria
    required: true
    prompt: |
      One north-star metric + 2-3 leading indicators. Each must specify how it's measured.
    example: |
      - North star: weekly active designers — measured via session ingest at /api/session/start
      - Leading: project creation rate per week — measured via project_created events
      - Leading: week-4 retention — measured via cohort analysis
  - name: Distribution
    after: Metrics
    required: true
    prompt: |
      How do users find the product? Be honest about acquisition. Include the activation signal.
    example: |
      Primary channel: organic search + designer-community partnerships.
      Activation signal: first project shared with a client.
modifies-slots:
  - name: Target users
    note: |
      Include a job-to-be-done framing alongside the persona. What are they "hiring" the product to do?
  - name: Deliverable
    note: |
      For products: name the surfaces (web app, mobile app, browser extension, etc.).
triggers-docs:
  always:
    - roadmap
  suggested:
    - stack
    - architecture
    - decisions
    - principles
description: |
  Consumer or B2B products with clear user flows.
  Adds User stories + Metrics + Distribution slots. Triggers ROADMAP.md scaffolding.
---

# Product kit

For products with end users and measurable engagement.

Use when the project has:
- Distinct user roles with workflows
- A north-star metric tied to engagement or retention
- A distribution strategy that matters as much as the product itself
- A roadmap with phased feature releases

Skip this kit when:
- The project ships only code/library artifacts (→ `coding`)
- Engagement isn't measurable yet (→ `blank` until clearer)
