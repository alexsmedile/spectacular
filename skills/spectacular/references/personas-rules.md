---
doc-id: personas
mode: grill-each
location: .spectacular/PERSONAS.md
scope: project-wide
template: templates/personas/base.md
slots: [Who, Wants to, Pain, Stories, Not for]
snapshot-on-edit: true
summary: "Opt-in audience profiles + user stories — Who / Wants to / Pain / Stories / Not for"
status: active
---

# PERSONAS Rules

Per-doc rules consumed by `grill.md` / `refine.md` / `review.md` when the user invokes `spectacular personas <verb>`.

PERSONAS.md is an **opt-in, `grill-each`-mode** doc. The agent walks all 5 slots once per persona block, then asks "add another persona? [y/n]" and repeats until done. Same shape as ROADMAP (per-version-block slot walks).

## Slot definitions (per persona block)

| Slot | Required | Form | Notes |
|---|---|---|---|
| `Who` | ✓ | One sentence | Role + context. No demographics, no "everyone". |
| `Wants to` | ✓ | One sentence | Outcome-focused. The thing they're trying to accomplish. Not the product feature. |
| `Pain` | ✓ | 1-2 bullets | What's friction-y today, before this product/feature exists. |
| `Stories` | ✓ | List, "As X, I want Y, so Z" format | Minimum 1 story per persona. Each story is one sentence. |
| `Not for` | — | 1 sentence (optional) | Who this persona explicitly excludes. Helps scope decisions. |

## Grill prompts

When grilling PERSONAS.md, the skill walks each persona block and asks:

### For each persona block

1. **Who is this?** "Give me one concrete sentence — role + context. Not a list, not a demographic, not 'everyone'."
2. **What are they trying to accomplish?** "One outcome. Not 'use the feature' — the *result* the feature enables."
3. **What's painful today?** "1-2 friction points. Be concrete. Avoid 'it's hard' / 'it's slow' without a why."
4. **Give me the top 2-3 stories.** "Use 'As X, I want Y, so Z.' Each is one sentence. Focus on the highest-value behaviors, not edge cases."
5. **Anyone this is explicitly NOT for?** *(optional)* "Skip if you don't have a sharp answer. Force-filling this slot rots faster than leaving it empty."

### Across the whole file

6. **How many personas do you have?** "If >5, you probably need to merge or cut. Ask: would removing one change a build decision? If no, cut it."
7. **Any persona without a story?** "A persona without stories is decoration. Either give it stories or remove it."
8. **Any duplicate stories across personas?** "Stories belong to the persona that benefits *most*. If 3 personas share a story, one of them owns it."

## Vague-word list (slot-specific)

| Slot | Vague words → push back |
|---|---|
| `Who` | "users", "everyone", "people", "developers" (without context), "teams" (without size/type) |
| `Wants to` | "use", "leverage", "engage with", "experience" (these are feature-words, not outcomes) |
| `Pain` | "it's hard", "it's slow", "bad UX", "lots of friction" (no concrete cause) |
| `Stories` | Stories missing the "so Z" outcome — these are tasks, not stories |

## Anti-patterns

- **JTBD methodology apparatus.** No job stories ("When _, I want to _, so I can _"), no functional/emotional/social decomposition, no outcome statements. A single "Wants to" line per persona is the design — *deliberately* simpler than JTBD frameworks. If you need JTBD depth, write a separate research doc; don't bloat PERSONAS.md.
- **Demographics.** No age, income, gender, location, education level — none of these drive build decisions and they age badly.
- **Photos / quotes / scenarios.** This is a reference doc, not a UX research artifact.
- **>5 personas.** Past 5, personas become decoration. Cut or merge.
- **Persona without stories.** Inverse: a persona that exists only to be listed. Either add stories or remove the persona.
- **Stories without a clear "so Z" outcome.** These are tasks, not stories — they belong in PLAN.md or TASKS.md.
- **Duplicating PRD § Target users verbatim.** PERSONAS.md is the *deepening* of that bullet, not a copy. PRD stays terse (one primary user); PERSONAS expands.

## Custom gate checks

For `spectacular personas review`:

1. ✅ **Frontmatter valid** — `version`, `updated`, `summary`, `related:` present
2. ✅ **At least one persona block** — file isn't empty
3. ✅ **Each persona has all required slots** — Who, Wants to, Pain, Stories
4. ✅ **Each persona has ≥1 story** — empty Stories list fails
5. ✅ **Each story matches "As X, I want Y, so Z" shape** — soft check; warn if a story is missing "so" or "want"
6. ⚠ **Persona count check** — info if 0; warn if >5
7. ⚠ **Vague-word check** — soft warn if any slot contains words from the vague-word list

## Refine patterns

When refining a near-empty PERSONAS.md:

- **Bare role becomes Who** — "designer" → ask: "freelance? in-house? at what company size?"
- **Feature description becomes Wants to** — "use Spectacular" → ask: "what's the outcome they're chasing by using it?"
- **Generic pain becomes specific** — "it's slow" → ask: "slow in what specific moment?"
- **Tasks become stories** — "log in" → push: "what does logging in unlock? Frame as 'As X, I want Y, so Z.'"

## Example output (passing)

A PERSONAS.md with 3 personas, each fully filled, ~80 lines total:

```
✓ frontmatter valid
✓ 3 persona blocks
✓ all required slots filled (Who, Wants to, Pain, Stories)
✓ 8 total stories, all in "As X, I want Y, so Z" form
✓ 1 persona has "Not for" slot filled (optional)
⚠ "Wants to" for persona 2 uses "leverage" — consider an outcome word
```

## Related

- [[grill]] — generic interactive slot-filler
- [[review]] — generic quality gate runner
- [[doc-index]] — registry entry for `personas`
- [[prd-rules]] — reference example of rules file shape
- [[scaffold-reference]] § PERSONAS.md — template stub
