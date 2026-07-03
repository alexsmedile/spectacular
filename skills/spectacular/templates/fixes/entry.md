<!--
  One fix entry — a FULLY RESOLVED AND VERIFIED bug fix.
  Written to .spectacular/fixes/F<N>.md. Scaffold-only in v1.25.0 (no CLI verb
  yet — hand-write from this template). See references/fixes-rules.md.

  A fix entry is not a symptom log. It's the settled record: what broke, why,
  what we changed, and the check that proves it holds. Log ONLY after verified.
-->

---
type: fix
opened: <DATE>
verified: <DATE>
severity: low | medium | high
from_audit: null
related: []
---

# Fix — <TITLE>

## Bug
<THE SYMPTOM THAT WAS REPORTED>

## Root cause
<WHY IT HAPPENED — the actual mechanism>

## Fix
<WHAT WAS CHANGED — file:line, the shape of the change>

## Verified by
<THE CHECK THAT PROVES IT: a test path, a repro-now-passing, a manual walk>
