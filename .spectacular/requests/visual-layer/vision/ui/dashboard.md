---
kind: ui
caption: Workspace progress table (default view)
approved: true
related: []
updated: 2026-06-02
---

# UI — dashboard

The default `spectacular progress` (no slug) output — the whole workspace in one table.

```
┌─ spectacular progress ─────────────────────────────────────┐
│ 2 active · 4 planned · 1 stalled            updated 2w span │
├─────────────────────────────────────────────────────────────┤
│ ● imagine-mode        active    [██████░░] 75%   2026-06-02 │
│ ⚠ visual-layer        active    [██░░░░░░] 25%   2026-05-25 │ ← stalled (8d)
│ ○ cross-request-links planned   [░░░░░░░░]  0%   2026-05-29 │
│ ○ cli-debt-removal    planned   [░░░░░░░░]  0%   2026-05-29 │
│ ○ roadmap-ledger      planned   [░░░░░░░░]  0%   2026-05-20 │
│ ○ progress-view       planned   [░░░░░░░░]  0%   2026-06-02 │
└─────────────────────────────────────────────────────────────┘
  ▸ dig in:  spectacular request <slug>
```

## Notes

- **Sort:** active first (by last-touched desc), then planned, then review/verified dimmed.
- **Markers:** `●` active · `○` planned · `⚠` stalled (active + no tick in ≥7d) · `✓` verified (dimmed).
- **Bar:** 8 cells, filled = checked-task ratio from TASKS.md (reuses the existing `_progress` parser).
- **Header:** the one-line triage summary — counts + the stalled count is the hook.
- Width held to ~62 chars so it reads in a narrow terminal.
- Open question for the human: is the last-touched **date** the right 4th column, or would **last-touched age** (`8d`, `2w`) read faster for the stall-spotting use case?
