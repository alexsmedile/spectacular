---
status: verified
updated: 2026-08-02
related:
  - PLAN.md
---

# Tasks — artifact-retention-contract

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

### M1 — Retention matrix
- [x] Define live, stale-safe, temporary, and throwaway semantics
- [x] Map canonical docs, collections, requests, specs, code/tests, and Git playgrounds
- [x] Define synchronization/check-on-reuse and index compaction rules
- [x] → check: every mapped artifact has one freshness obligation and one terminal disposition

### M2 — Terminal behavior
- [x] Archive resolved questions while preserving answer provenance and dependency satisfaction
- [x] Keep DEC creation optional and choice-gated
- [x] Allow explicit archive-first closure for implemented and abandoned detailed specs
- [x] Preserve hidden recovery refs before disposable branch deletion
- [x] → check: focused lifecycle and wayfinding scenarios prove active context shrinks without history loss

### M3 — Briefing, docs, and verification
- [x] Surface active user questions first in status/onboarding/session flows
- [x] Align roadmap, discovery, lifecycle, architecture, system-spec, and user docs
- [x] Run syntax, version, focused, full-suite, doctor, and diff checks
- [x] Record docs impact and spec delta
- [x] → check: all executable gates pass

## v2 (deferred)

- [~] Add automated freshness timestamps/SLAs only if real projects need enforcement beyond lifecycle-derived rules
- [~] Add remote archive/deletion automation only after a separate explicit security and recovery design
