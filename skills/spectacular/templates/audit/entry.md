<!--
  One audit entry — a bug/quirk investigation BEFORE a fix is planned.
  Written to .spectacular/audit/A<N>.md by `spectacular audit new`.
  Substitution tokens (filled by CLI): <ID> <TITLE> <SEVERITY> <DATE> <PROBLEM> <INTENDED>.
  Other bracketed text is a fill-in hint left in place. See references/audit-rules.md.

  The body follows the bug-fixing skeleton:
    problem → intended behavior → (investigate) → root cause → proposed fix →
    success criteria → disposition.
-->

---
type: audit
status: open
severity: <SEVERITY>
opened: <DATE>
updated: <DATE>
disposition: null
related: []
---

# <ID> — <TITLE>

## Problem
<PROBLEM>

## Intended behavior
<INTENDED>

## Investigation
_(what was checked, what was ruled out)_

## Root cause
_(the actual cause, once found — or "not yet found")_

## Proposed fix
_(the suggested fix — still a proposal at audit stage, not yet applied/verified)_

## Success criteria
_(how we'll know it's fixed — the observable bar the fix must clear)_

## Disposition
_(set on close: requests/<slug> · one-line fix · won't-fix · became fix F<N>)_
