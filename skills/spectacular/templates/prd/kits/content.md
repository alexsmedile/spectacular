---
kit: content
version: 2.0
extends: prd
adds-slots:
  - name: Audience
    after: Target users
    required: true
    prompt: |
      Beyond the persona: where does this audience currently get this kind of content?
      Why would they switch?
    example: |
      Currently consuming: Stratechery, Lenny's Newsletter.
      Will switch because: deeper technical examples paired with strategic framing.
  - name: Format
    after: Deliverable
    required: true
    prompt: |
      Concrete production format. Medium, structure, length, cadence.
    example: |
      Medium: text + occasional diagrams.
      Structure: long-form essay with embedded code samples.
      Length: 1,500-2,500 words per issue.
      Cadence: weekly, Tuesdays.
  - name: Distribution
    after: Format
    required: true
    prompt: |
      Where pieces live, how the audience finds them, how they retain.
    example: |
      Home base: own newsletter (Buttondown).
      Discovery: cross-posts to dev.to + HN Show.
      Retention: email subscription, ~50% open rate target.
  - name: Editorial principles
    after: Constraints
    required: false
    prompt: |
      3-5 rules that define the voice and quality bar. What makes a piece "on-brand"?
    example: |
      - No "10 reasons" listicles — every piece has one core argument
      - All code samples runnable as-is
      - Cite primary sources, never secondhand summaries
modifies-slots:
  - name: Deliverable
    note: |
      For content: name the format + cadence + count.
      Example: "Weekly newsletter, 12-issue arc, published Tuesdays."
  - name: Goals & success criteria
    note: |
      For content: measure engagement, not publication.
      Bad: "publish 12 issues". Good: "1,000 subscribers + 40% open rate by episode 8."
triggers-docs:
  always:
    - roadmap
  suggested:
    - principles
    - decisions
description: |
  Courses, newsletters, books, video series, podcasts, documentation projects.
  Adds Audience + Format + Distribution + Editorial principles slots.
---

# Content kit

For content projects where engagement and audience-building matter.

Use when the project ships:
- A newsletter or essay series
- A course or tutorial sequence
- A book or long-form publication
- A video or podcast series
- A documentation project with a defined audience

Skip this kit when:
- The project ships code as the primary artifact (→ `coding`)
- The output is research feeding a single decision, not ongoing publication (→ `research`)
