# Refine — generic vibe→spec rewriter and append-mode handler

Loaded when the user runs `spectacular <doc> refine` (full-document refine pass) or implicitly by `grill.md` (mini-refine during slot loop) or when `<doc>` has `mode: append` in the registry.

This is the **doc-agnostic** engine. Per-doc patterns live in rules files (`prd-rules.md`, `plan-rules.md`, etc.).

## Core principle

**Propose, don't impose. Mark unresolved gaps explicitly.**

Refine never silently rewrites the user's words. It either:
1. Proposes a tighter version and asks for acceptance, OR
2. Inserts `[NEEDS CLARIFICATION: <specific gap>]` inline so the gap is visible.

The user is always in the driver's seat.

## Three modes

### Inline mini-refine

Called by `grill.md` after each slot answer. Fast, targeted, **one slot only**. Reads patterns from the doc's rules file (if present) plus base patterns.

### Full-document refine

Called by `spectacular <doc> refine`. Walks the entire file, flags every pattern hit, produces a single diff-style proposal the user reviews in one pass.

```
spectacular <doc> refine
  ↓
Resolve registry entry → load location, overrides
  ↓
Read file at location
  ↓
Scan all sections for patterns (base + override-supplied)
  ↓
Produce a proposal:
  - For each finding: location, what's wrong, proposed rewrite
  ↓
User reviews — accept all / accept selectively / skip
  ↓
If snapshot-on-edit: true, snapshot prior version as <DOC>@vN.md
  ↓
Write accepted changes; insert [NEEDS CLARIFICATION] for skipped
```

### Append mode

Called when registry says `mode: append`. Used for DECISIONS.md and similar log-style docs.

```
spectacular decisions
  ↓
Read template entry (e.g. templates/decisions/entry.md)
  ↓
Ask: title / decision / why / tradeoffs (or whatever the entry template fields are)
  ↓
Build the entry from the template
  ↓
Append to the file at location (do NOT replace existing content)
  ↓
Add a date header if convention requires (e.g. "## YYYY-MM-DD — <title>")
```

Append mode never grills slots and never runs the review gate — the unit of change is one entry, not the whole file.

## Pattern sources

The engine combines two pattern sources:

### Base patterns (universal)

Apply to any doc unless explicitly exempted by the rules file.

- **Vague adjectives** — words like `fast`, `simple`, `intuitive`, `scalable`, `seamless`, `great`, `flexible`, `robust` → propose measurable replacement
- **Empty lists** — list slots with zero or one item, where context implies multiple → propose "add more"
- **Placeholder leftover** — uppercase angle-bracket markers (`<ANYTHING>`) not yet filled → propose answer or `[NEEDS CLARIFICATION]`

**Tokenization rule** (applies to all word-based pattern matching) — preserve hyphenated compounds as single tokens. Use `[a-z]+(?:-[a-z]+)*` (or equivalent) instead of `\b\w+\b`. This prevents false positives where a vague word appears inside a compound identifier (e.g. `smart-init`, `doc-writer`, `kits-as-plugins` should not match the vague-word list even though they contain `smart`, `doc`, `kits`). Compound identifiers are almost always slugs, package names, or hyphenated technical terms — preserving them dramatically reduces false-positive rate.

### Per-doc patterns (rules file)

Loaded from the rules file referenced in the registry. Examples:

- PRD's "plural-user → singular" rule (only valid for the Target users slot)
- PRD's "no number+verb+date in success" check (only valid for Goals slot)
- PLAN's "milestone before tasks" ordering rule
- PLAN's "dependency link validation" (frontmatter `related:` targets must exist)

Rules files declare which slots their patterns apply to. The engine never blindly applies a pattern to the wrong slot.

## The `[NEEDS CLARIFICATION]` convention

Borrowed from spec-kit. Inserted inline at the exact location of the gap:

```markdown
## 5. Goals & success criteria

[NEEDS CLARIFICATION: success criteria has no measurable signal — add a number + timeframe]
```

Rules:
- One marker per distinct gap.
- Specific to the gap — never bare `[NEEDS CLARIFICATION]`.
- The review gate fails as long as any marker exists.
- The user resolves markers by editing in place or re-running `<doc> refine`.

## Output format

### Inline mini-refine (during grill)

```
> "make X writing easier"

⚠️ Vague: "easier" — by what measure?
  Examples:
    - "reduces time to first <doc> from 2 hours to 15 minutes"
    - "produces a <doc> that gets referenced in planning within 7 days"
  Override / pick one / leave as [NEEDS CLARIFICATION]?
```

### Full refine pass

```
Found 4 issues in .spectacular/PRD.md:

  L14  Slot 2 (Problem)  "make X writing easier"
       → vague — propose: "reduces time to first PRD from 2 hours to 15 minutes"

  L20  Slot 3 (Target users)  "for developers"
       → plural — propose: "[NEEDS CLARIFICATION: pick one primary user]"

  L28  Slot 4 (Deliverable)  "a tool"
       → vague deliverable — propose: "[NEEDS CLARIFICATION: name the artifact]"

  L34  Slot 5 (Goals)  no number found
       → propose: "[NEEDS CLARIFICATION: add number + date]"

Apply all / pick specific / cancel?
```

### Append mode

```
spectacular decisions

  Title: Use Bash for the CLI
  Decision: Ship the bootstrap CLI as a single Bash script
  Why: Zero install dependencies, works on macOS/Linux out of the box
  Tradeoffs: No Windows support without WSL; harder to test than TypeScript

  Appending entry to .spectacular/DECISIONS.md ... ✓
```

## What refine does NOT do

- It does not invent content. If the user wrote "fast", refine asks "what does fast mean?" — it does not pick a number.
- It does not reformat or restructure. Only the flagged spans change.
- It does not touch frontmatter unless explicitly asked.
- It does not cross doc boundaries. Refining PRD never touches PLAN.

## Karpathy alignment

- **Think before coding:** propose, don't silently rewrite.
- **Simplicity first:** pattern-match + question generation. No NLP.
- **Surgical changes:** only flagged spans get touched, never adjacent content.
- **Goal-driven:** every refine cycle moves the doc closer to passing review.

## Related

- [[doc-registry]] — the registry the engine consumes
- [[grill]] — calls mini-refine inline during slot loop
- [[review]] — gate that depends on these patterns being resolved
- [[prd-rules]], [[plan-rules]], [[tasks-rules]] — per-doc pattern sources
- [[versioning]] — full refine snapshots prior version when `snapshot-on-edit: true`
