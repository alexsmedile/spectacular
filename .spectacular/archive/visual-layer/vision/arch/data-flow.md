---
kind: arch
caption: How the aggregate view is computed
approved: true
related: []
updated: 2026-06-02
---

# Architecture — data-flow

How `spectacular progress` (no slug) computes the aggregate. Reuses the existing per-request parser in a loop — no new parsing logic.

```
requests/*/                  ┌─ per request ────────────┐
  ├─ PLAN.md  (status,──────►│ fm_get status, updated   │
  │           updated)       │                          │
  └─ TASKS.md (checkboxes)──►│ _progress parser → ratio │──┐
                             │ last-tick date           │  │
                             └──────────────────────────┘  │
                                                            ▼
                                              ┌─ aggregate ──────────┐
                                              │ collect rows         │
                                              │ sort: active→planned │
                                              │ flag stalled (≥7d)   │
                                              │ count summary header │
                                              └──────────┬───────────┘
                                                         ▼
                                              render table (ui/dashboard)
```

## Notes

- **No new parser.** `_progress` already exists (per-slug); this loops it over `requests/*/`.
- **Stalled** = `status:active && (today − last_tick_date) ≥ threshold`. Last-tick date inferred from the newest `- [x]` line's context or TASKS `updated:`.
- The whole feature is "loop the thing we have + sort + a header" — that's why it's a small MINOR, not a big build. Confirms the request is right-sized.
