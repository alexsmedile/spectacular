<!--
  One audit entry — a bug/quirk investigation BEFORE a fix is planned.
  Written to .spectacular/audit/A<N>.md. Scaffold-only in v1.25.0 (no CLI verb
  yet — hand-write from this template). See references/audit-rules.md.

  audit/ is the diagnostic scratchpad between "something's off" and "here's the
  plan". Its exit is a disposition: fold into a request, one-line fix, won't-fix,
  or became a fix entry F<N>.
-->

---
type: audit
status: open
severity: low | medium | high
opened: <DATE>
updated: <DATE>
disposition: null
related: []
---

# Audit — <TITLE>

## Symptom
<WHAT WAS OBSERVED — the report, not the cause>

## Investigation
<WHAT WAS CHECKED, WHAT WAS RULED OUT>

## Root cause
<THE ACTUAL CAUSE, ONCE FOUND — or "not yet found">

## Disposition
<ONE OF: folded into requests/<slug> · one-line fix · won't-fix (why) · became fix F<N>>
