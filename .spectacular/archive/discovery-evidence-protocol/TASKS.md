---
status: verified
updated: 2026-08-02
related:
  - PLAN.md
---

# Tasks — discovery-evidence-protocol

<!--
  Executable checklist for one request.
  Lives at: .spectacular/requests/<slug>/TASKS.md

  Rules:
  - Group tasks by milestone using `### M<N> — <name>` headings.
  - Flush-left checkboxes are the COUNTED units: `- [ ]` open, `- [x]` done,
    `- [~]` deferred (not-open-not-done; shown separately in progress).
  - Indented `  - [ ]` sub-bullets are allowed as a nested acceptance checklist
    under a task, but are NOT counted — progress counts top-level only, so
    x/total stays comparable across requests.
  - `status:` in frontmatter should match parent PLAN.md.
  - Tasks are owned by the user. Engine never adds/removes/reorders tasks.
-->

## v1

### M1 — Progressive routing gate
- [x] Define the cheapest-answer-first threshold and explicit skip conditions
- [x] Define the research, spike, prototype, and direct-implementation decision table
- [x] Preserve human/AFK authority boundaries and inconclusive-fog behavior
- [x] → check: every route begins with a named uncertainty and unnecessary discovery creates no node

### M2 — Ownership and durability
- [x] Map code destination, owning record, durable output, and cleanup for each discovery mechanism
- [x] Define tracer bullets as approved-spec production execution rather than discovery nodes
- [x] Define artifacts as owned outputs rather than a new entity or catch-all database
- [x] Route technical debt through requests, roadmap candidates, ideas, and decisions without a debt folder
- [x] → check: no concept has two owners, two IDs, or contradictory durability rules

### M3 — Contract and documentation coherence
- [x] Add the canonical protocol and wire the skill router
- [x] Align collection, lifecycle, identity, architecture, system-spec, workflow, and scaffold docs
- [x] Run syntax, version, doctor, full-suite, and diff-integrity checks
- [x] Record docs-impact evidence and verification outcome
- [x] → check: all executable gates pass and user-facing examples match the canonical protocol

## v2 (deferred)

- [~] Activate a standalone `PRT-NNN` workflow only if repeated prototype retrieval needs project-wide identity
- [~] Add a `DEB-NNN` collection only if real projects demonstrate that roadmap/request/idea routing loses debt
