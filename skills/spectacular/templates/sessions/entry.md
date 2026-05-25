<!--
  Single session entry. Written to .spectacular/sessions/<DATE>-<SLUG>.md by `spectacular session start`.
  Slug derives from --tag (first tag) or "work" if none provided.

  status: open  → set by `spectacular session start`
  status: closed + end_date + counts → set by `spectacular session end`
  Linked decisions/memories appended at session end by scanning entries with matching `session:` field.
-->

---
type: session
status: open
start_date: <START_DATE>
end_date: null
tags: <TAGS>
related: []
summary: "<SUMMARY>"
decisions_count: 0
memories_count: 0
---

# Session — <SLUG>

<NOTES>

<!-- Linked decisions + memories appended automatically by `spectacular session end` -->
