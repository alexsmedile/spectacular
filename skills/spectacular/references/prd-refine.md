# PRD Refine — vibe → spec rewrite patterns

Loaded when the user runs:
- `prd refine` — full-document refine pass on an existing PRD.md
- Implicitly by `prd-grill` during the slot loop (mini-refine, one slot at a time)

Refine is the **transformation engine** that turns conversational/vague language into specific, measurable, actionable PRD content.

## Core principle

**Propose, don't impose. Mark unresolved gaps explicitly.**

Refine never silently rewrites the user's words. It either:
1. Proposes a tighter version and asks for acceptance, OR
2. Inserts `[NEEDS CLARIFICATION: <specific gap>]` inline so the gap is visible.

The user is always in the driver's seat.

## Two modes

### Inline mini-refine

Called by `prd-grill` after each slot answer. Fast, targeted, **one slot only**. See the pattern table below.

### Full-document refine

Called by `prd refine`. Walks the entire PRD.md, flags every pattern hit, produces a single diff-style proposal the user reviews in one pass.

```
spectacular prd refine
  ↓
Read .spectacular/PRD.md
  ↓
Scan all sections for patterns
  ↓
Produce a proposal:
  - For each finding: location, what's wrong, proposed rewrite
  ↓
User reviews — accept all / accept selectively / skip
  ↓
Write accepted changes; insert [NEEDS CLARIFICATION] for skipped
  ↓
Snapshot prior version as PRD@vN.md
```

## Vibe → spec pattern table

Patterns listed in priority order (highest signal first):

### 1. Vague adjectives → measurable replacement

| Vibe | Spec |
|---|---|
| "fast" | "p95 latency < 200ms on M1 laptop" |
| "intuitive" | "new user completes first task in < 60s without docs" |
| "scalable" | "handles 10x current load without architectural changes" |
| "simple" | "one command, no flags, no config file required" |
| "seamless" | "[NEEDS CLARIFICATION: which transitions feel rough today?]" |
| "great UX" | "[NEEDS CLARIFICATION: what user-visible behavior signals 'great'?]" |
| "robust" | "passes <specific failure scenario list>" |
| "flexible" | "supports <specific extension points>" |

Vague-word list (full): `fast`, `slow`, `intuitive`, `simple`, `easy`, `scalable`, `seamless`, `great`, `good`, `nice`, `clean`, `flexible`, `robust`, `powerful`, `elegant`, `modern`, `smart`, `lightweight`, `heavyweight`, `solid`, `polished`.

### 2. Plural-user → single primary user

In **slot 2 (Who it's for)** only:

| Vibe | Spec |
|---|---|
| "developers" | "Solo devs using Claude Code on side projects, no PM, low budget" |
| "users" | "[NEEDS CLARIFICATION: pick one primary user with role + situation + constraint]" |
| "everyone" | "[NEEDS CLARIFICATION: pick one primary user — 'everyone' = no one]" |
| "small teams" | "5-15 person product teams shipping weekly, no dedicated PM" |

Other slots can have plurals (non-goals lists, stakeholders, etc.) — only refine plural primary users in slot 2.

### 3. Unbounded success → number + verb + date

In **slot 3 (Success)** only. Required signals:
- At least one **number** (or quantifiable comparison)
- At least one **verb** (action / observable outcome)
- At least one **timeframe** (date, "N days/weeks", or named event)

| Vibe | Spec |
|---|---|
| "users love it" | "By 2026-09-01, NPS > 40 across 200+ respondents" |
| "ship fast" | "Demo-ready by 2026-07-01; v1 in users' hands by 2026-08-15" |
| "drives adoption" | "1,000 weekly active users within 90 days of launch" |
| "high quality" | "[NEEDS CLARIFICATION: what's the quality bar — bug rate, test coverage, NPS?]" |

### 4. Tech jargon in problem statement

In **slot 1 (Problem)** only. The problem must be user-visible, not implementation-flavored.

| Vibe | Spec |
|---|---|
| "no embeddings layer yet" | "User searches return irrelevant results — they give up after 2 queries" |
| "lack of microservices" | "[NEEDS CLARIFICATION: what user-visible pain does the missing architecture cause?]" |
| "no CI/CD" | "Deploys take 2 hours and break weekly — team avoids releasing on Fridays" |

### 5. Empty / weak non-goals

In **slot 4** only.

If non-goals has fewer than 2 items, ask: "What would you push back on if a stakeholder tried to expand scope tomorrow?"

If non-goals are themselves vague (e.g. "scope creep", "feature bloat"), refine to concrete examples.

### 6. Soft constraints → hard constraints

In **slot 5**.

| Vibe | Spec |
|---|---|
| "limited time" | "Ships by 2026-08-15; ~40 dev-hours/week available" |
| "small team" | "1 engineer + occasional design help (3 hrs/week)" |
| "tight budget" | "$500 total infrastructure spend through Q3" |

### 7. Vague milestones

In **slot 6**.

| Vibe | Spec |
|---|---|
| "MVP" | "[NEEDS CLARIFICATION: what does MVP demo look like? Which user can do which thing?]" |
| "beta" | "10 invited users complete the core flow without bugs by 2026-07-15" |
| "v1" | "Tagged v1.0 release, install path documented, 1 external user onboarded" |

## The `[NEEDS CLARIFICATION]` convention

Borrowed from spec-kit. Inserted inline at the exact location of the gap:

```markdown
## 3. What success looks like

[NEEDS CLARIFICATION: success criteria has no measurable signal — add a number + timeframe]
```

Rules:
- One marker per distinct gap.
- Specific to the gap — never bare `[NEEDS CLARIFICATION]`.
- The review gate fails as long as any marker exists.
- The user resolves markers by editing in place or re-running `prd refine`.

## Output format

### Inline mini-refine (during grill)

```
> "make PRD writing easier"

⚠️ Vague: "easier" — by what measure?
  Examples:
    - "reduces time to first PRD from 2 hours to 15 minutes"
    - "produces a PRD that gets referenced in planning within 7 days"
  Override / pick one / leave as [NEEDS CLARIFICATION]?
```

### Full refine pass

```
Found 5 issues in .spectacular/PRD.md:

  L12  Slot 1  "make PRD writing easier"
       → vague — propose: "reduces time to first PRD from 2 hours to 15 minutes"

  L18  Slot 2  "for developers"
       → plural — propose: "[NEEDS CLARIFICATION: pick one primary user]"

  L24  Slot 3  no number found
       → add timeframe — propose: "[NEEDS CLARIFICATION: add number + date]"

  L31  Slot 4  only 1 non-goal listed
       → propose: "What else would you push back on?"

  L40  Slot 6  "MVP"
       → vague — propose: "[NEEDS CLARIFICATION: define MVP scope]"

Apply all / pick specific / cancel?
```

## What refine does NOT do

- It does not invent content. If the user wrote "fast", refine asks "what does fast mean?" — it does not pick a number.
- It does not reformat or restructure. Only the flagged spans change.
- It does not touch frontmatter unless explicitly asked.
- It does not run on PLAN.md or TASKS.md — those are scope of separate refinement passes (future).

## Karpathy alignment

- **Think before coding:** propose, don't silently rewrite.
- **Simplicity first:** ~7 patterns total, no NLP, just pattern matching + question generation.
- **Surgical changes:** only flagged spans get touched, never adjacent content.
- **Goal-driven:** every refine cycle moves the PRD closer to passing the review gate.

## Related

- [[prd-grill]] — calls mini-refine inline during slot loop
- [[prd-review]] — gate that depends on these patterns being resolved
- [[versioning]] — full refine snapshots prior version as `PRD@vN.md`
