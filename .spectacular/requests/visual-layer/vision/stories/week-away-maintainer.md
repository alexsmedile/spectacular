---
kind: story
caption: Maintainer returns after a week
approved: true
personas: [Solo OSS maintainer]
related: []
updated: 2026-06-02
---

# Story — week-away-maintainer

**As a** Solo OSS maintainer, **I want** one command that shows how far along *every* request is, **so that** after time away I can see what's in flight and what's stalled without opening six TASKS.md files.

## Acceptance

- `spectacular progress` with **no slug** prints a table of all non-archived requests.
- Each row shows status, a milestone bar, and a recency signal.
- Stalled active requests (no tick in ≥7d) are visually flagged.
- The view fits one screen and reads in ≤3 seconds.

## Implied flow

Returns after a week → `spectacular progress` → scans the table → spots the `⚠` row → `spectacular request visual-layer` to dig in. See `ui/dashboard.md` for the rendered table this story implies.
