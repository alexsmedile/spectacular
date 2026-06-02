---
kind: ui
caption: Stalled-only triage view (--stalled)
approved: true
related: []
updated: 2026-06-02
---

# UI — stalled-filter

`spectacular progress --stalled` — the triage cut. Only active requests with no recent tick.

```
┌─ spectacular progress --stalled ───────────────────────────┐
│ 1 stalled (active, no tick in ≥7d)                          │
├─────────────────────────────────────────────────────────────┤
│ ⚠ visual-layer        active    [██░░░░░░] 25%   last +8d   │
│     last tick: M1 done 2026-05-25 · 0 ticks since           │
└─────────────────────────────────────────────────────────────┘
  ▸ triage:  spectacular request visual-layer
  (nothing stalled? prints "✓ no stalled requests")
```

## Notes

- Same row format as the default view, filtered to `status:active && days_since_last_tick ≥ threshold`.
- Threshold default 7d; overridable `--stalled=14` (mirrors the `--since` flag grammar).
- This is the view that makes `progress` a *triage tool*, not just a status readout.
- Empty state matters: an explicit `✓ no stalled requests` beats a blank table.
