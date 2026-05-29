---
description: Quality-gate runner — structural checks in CLI, semantic review in skill.
when_to_use: spectacular <doc> review.
---

# Review — generic quality gate runner

Loaded when the user runs `spectacular <doc> review` or implicitly at the end of `grill.md`.

This is the **doc-agnostic** skill. Per-doc gate checks live in rules files (`prd-rules.md`, `plan-rules.md`, etc.).

Review is a **pass/fail quality gate**. It produces a punch list, not a rewrite. The user decides what to fix.

## Core principle

**No score. No rewrite. Just a checklist + actionable findings.**

Review is the safety net at the end of every grill and every refine. It reads the file at the registry's `location:` and reports exactly what blocks it from being "done".

## Base checks (universal — apply to any `mode: grill` doc unless exempted)

| # | Check | How |
|---|---|---|
| 1 | All required slots non-empty | No `<PLACEHOLDER>` markers remain in required slots |
| 2 | No unresolved clarifications | Zero `[NEEDS CLARIFICATION: ...]` markers anywhere |
| 3 | Frontmatter valid | Required fields per registry's frontmatter schema present |

If all base checks pass AND all override checks pass → doc is usable. Exit with a green confirmation.
If any fail → produce a punch list.

## Override checks

Loaded from the rules file referenced in the registry. Each rules file declares:

- Additional gate checks (e.g. PRD's "Goals must contain number+verb+date")
- Slot-specific exemptions (e.g. PRD's "Vision slot exempt from vague-word scan")
- Heuristic rules (e.g. PRD's "Target users must be singular")
- Vague-word lists scoped to specific slots

Override checks run **in addition to** base checks. An rules file cannot disable base checks 1 or 2 (placeholder + clarification are universal); it can only add to them or scope check 4+ to specific slots.

## Detection rules (base)

### Check 1: Placeholders

Look for these unresolved markers anywhere in required slots:
- Any uppercase angle-bracket placeholder (`<ANYTHING>`)
- Slot-specific names from the template (`<PROBLEM>`, `<GOAL>`, `<MILESTONE>`, etc.)

### Check 2: Clarifications

Regex: `\[NEEDS CLARIFICATION:.*?\]` — any match anywhere fails the gate.

### Tokenization (applies to all word-matching checks)

Use hyphen-aware tokenization (`[a-z]+(?:-[a-z]+)*`) — never naïve `\b\w+\b`. Compound identifiers like `smart-init`, `doc-writer`, `kits-as-plugins` must tokenize as single units so they don't trigger vague-word hits on their constituent parts. See `refine.md` § Base patterns for the canonical rule.

### Verb detection (when an override requires verb presence)

Verb checks must accept common inflections (`-s`, `-es`, `-ed`, `-ing`, `-d`) by default, not bare stems only. Real writing uses `passes`, `exercised`, `shipped`, etc. Rules files declaring a verb list should either: enumerate inflections explicitly, or use a stem-with-suffix regex like `\b<stem>(s|es|ed|ing|d)?\b`.

### Check 3: Frontmatter

Required fields per registry (varies by doc type):
- Root docs: `version`, `updated`, `summary`, `related`
- Per-request PLANs: `status`, `priority`, `owner`, `updated`, `summary`
- Per-request TASKS: `status`, `updated`, `related`

Missing required field → fails. Punch list says exactly which field.

## Output format

### When the gate passes

```
✓ <doc> review passed

  <location> is ready.
  All required slots filled, no clarification gaps.

  Next:
    - spectacular new <slug>            scaffold the first request
    - spectacular snapshot <doc>        create a versioned snapshot
```

### When the gate fails (punch list)

```
✗ <doc> review found 4 issues

  Slot 2 — Problem
    L14  vague: "easier" — quantify or replace

  Slot 3 — Target users
    L20  plural primary user: "developers"
         → pick one specific role + situation + constraint

  Slot 4 — Deliverable
    L28  generic only: "a tool" — name the artifact

  Slot 8 — First milestone
    L48  [NEEDS CLARIFICATION: define MVP scope]

  Fix path:
    spectacular <doc> refine        # propose rewrites for all findings
    spectacular <doc> grill         # re-run the slot loop on flagged slots
    (or edit <location> directly and re-run review)
```

## Behavior by mode

- **`mode: grill`** — review runs all base + override checks
- **`mode: append`** — review runs base checks (frontmatter validity) + checks last appended entry conforms to entry template
- **`mode: stub`** — review runs base checks only (placeholder + frontmatter); no slot-completeness check (slots aren't enforced for freeform docs)

## What review does NOT do

- It does not score. Pass or fail.
- It does not rewrite. The user (or `refine`) does that.
- It does not check optional sections unless the rules file declares one as required.
- It does not check across docs. Cross-doc consistency is [[doctor]]'s job (v2).
- It does not warn about length. Length isn't a quality signal here.

## Karpathy alignment

- **Think before coding:** the gate makes assumptions checkable.
- **Simplicity first:** binary checks, no scoring, no NLP.
- **Surgical changes:** review never touches the file, only reads.
- **Goal-driven:** every check has a verifiable signal.

## Related

- [[doc-index]] — the registry the skill consumes
- [[grill]] — runs review as its stop condition
- [[refine]] — produces the rewrites that fix review findings
- [[prd-rules]], [[plan-rules]], [[tasks-rules]] — per-doc check sources
- [[versioning]] — snapshot before any significant rewrite
