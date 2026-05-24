# PRD Overrides — PRD-specific rules consumed by the generic engine

Loaded by `grill.md` / `refine.md` / `review.md` when the active doc is `prd` (per registry).

This file declares everything PRD-specific. The generic engine handles the rest.

## Kit selection (grill pre-flight only)

PRD is the only doc with `kit-support: true`. Before the slot loop starts, the grill:

1. **Discovers available kits** — scan `templates/prd/kits/*.md`, parse frontmatter, use each kit's `description:` field
2. **Presents the menu** — numbered list of kit names + descriptions (no hardcoded list — discovery is automatic)
3. **Asks the user to pick one**

Example menu (built from kit frontmatter):

> What kind of project is this?
> 1. **blank** — No extras. Pure 8-slot base PRD. Use when no other kit fits.
> 2. **coding** — CLI, library, app, service, SDK. Adds Stack + Interfaces slots; triggers STACK + ARCHITECTURE.
> 3. **content** — Newsletter, course, book, video. Adds Audience + Format + Distribution slots.
> 4. **product** — Consumer or B2B with user flows. Adds User stories + Metrics + Distribution slots.
> 5. **research** — Investigation feeding a decision. Adds Hypothesis + Method + Decision-being-informed slots.

After selection:
1. **Resolve kit file** — project-local override (`.spectacular/templates/prd/kits/<kit>.md`) wins over bundled
2. **Scaffold base PRD** — copy `templates/prd/base.md` to `.spectacular/PRD.md`
3. **Apply kit deltas** — read kit's `adds-slots` and `modifies-slots`; insert added slots at their `after:` positions; layer modify-slot notes onto base prompts
4. **Set frontmatter** — write `kit: <kit-id>` to PRD.md frontmatter (used by review gate)

**v1 constraint:** single-kit-only — exactly one kit per PRD. Multi-kit composition deferred to v2. See [[kits-contract]].

**Kit schema:** see [[kits-contract]] for the full extension contract (adds-slots, modifies-slots, triggers-docs).

## Slot prompts

The engine uses these in the slot loop. One question per slot, short, with one good/bad example to anchor expectations.

**Slot 1 — Vision**
> One paragraph. What is this, philosophically? Why does it exist in the world?
>
> Higher abstraction is OK here — this is the narrative frame, not measurable success.
>
> *Example:* "Spectacular is an AI-native operational workspace for software projects. It helps humans and coding agents maintain coherence across long-running development by separating strategy, current truth, active work, and operational memory."

**Slot 2 — Problem**
> What concrete pain does this solve? One sentence. Who is hurting, in what specific situation, how often.
>
> *Avoid:* "make X better", "improve Y experience"
> *Example:* "Solo devs writing PRDs from scratch waste 2+ hours and end up with vague documents nobody references."

**Slot 3 — Target users**
> Describe **one** primary user. Not a list, not "everyone". A specific role, situation, and constraint.
>
> *Example:* "Solo devs on side projects who use Claude Code and don't have a PM to write specs for them."

**Slot 4 — Deliverable**
> What concretely ships? Name the artifacts.
>
> Distinguish from Vision (the why) and Goals (the how-we-know-it-worked).
>
> *Example software:* "Three layers — Convention (.spectacular/ directory), Skill (/spectacular slash command), CLI (spectacular init)."
> *Example content:* "Weekly newsletter, 12-issue arc, published every Tuesday."

**Slot 5 — Goals & success criteria**
> Measurable. Time-boxed. At least one number, one verb, and one date or timeframe.
>
> *Avoid:* "users love it", "ship fast"
> *Example:* "30 days after launch, 50% of users who run `/spectacular prd` open their PRD.md again within 7 days."

**Slot 6 — Non-goals**
> What are you **not** doing? List 3-5 explicit exclusions you'd push back on if asked to expand.

**Slot 7 — Constraints**
> What's fixed before you start? Budget, time, tech, policy, team.
>
> *Example:* "Markdown-only, no new binaries", "ships before 2026-07-01".

**Slot 8 — First milestone**
> One concrete, demoable outcome that proves this is real. Date-bound.

## Mini-refine patterns

Applied inline by the grill after each answer.

| Pattern | Slots scope | Trigger | Proposed action |
|---|---|---|---|
| Vague adjective | 2, 5, 8 | Vague-word list hit | "What does '<word>' mean concretely? A number or comparison?" |
| Plural user | 3 only | `users`, `customers`, `developers`, `people` without qualifier | "Pick **one** primary user. Who's the most important?" |
| Unbounded success | 5 only | No number AND no date | "Add a number and a timeframe. Example: 'X by date Y'." |
| Tech jargon in problem | 2 only | `microservices`, `embeddings`, `framework`, similar | "Restate in plain language — what's the user-visible pain?" |
| Empty exclusion | 6 only | < 2 items | "What would you push back on if someone tried to expand scope?" |
| Vague deliverable | 4 only | Generic-only: `tool`, `system`, `platform`, `framework` without concrete artifact | "Name the concrete artifact — binary, library, doc, format." |

## Mini-refine exemptions

**Slot 1 (Vision) is fully exempt from mini-refine** — narrative abstraction is expected; the measurable slots (5, 8) are where precision matters.

**Slot 4 (Deliverable) is exempt from the generic vague-word scan** but has its own concrete-artifact check (see "Vague deliverable" pattern above).

## Vibe → spec rewrite tables (refine mode)

### Vague adjectives

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

### Plural users → singular

| Vibe | Spec |
|---|---|
| "developers" | "Solo devs using Claude Code on side projects, no PM, low budget" |
| "users" | "[NEEDS CLARIFICATION: pick one primary user with role + situation + constraint]" |
| "everyone" | "[NEEDS CLARIFICATION: pick one primary user — 'everyone' = no one]" |
| "small teams" | "5-15 person product teams shipping weekly, no dedicated PM" |

### Unbounded success → number + verb + date

| Vibe | Spec |
|---|---|
| "users love it" | "By 2026-09-01, NPS > 40 across 200+ respondents" |
| "ship fast" | "Demo-ready by 2026-07-01; v1 in users' hands by 2026-08-15" |
| "drives adoption" | "1,000 weekly active users within 90 days of launch" |
| "high quality" | "[NEEDS CLARIFICATION: what's the quality bar — bug rate, test coverage, NPS?]" |

### Tech jargon in problem

| Vibe | Spec |
|---|---|
| "no embeddings layer yet" | "User searches return irrelevant results — they give up after 2 queries" |
| "lack of microservices" | "[NEEDS CLARIFICATION: what user-visible pain does the missing architecture cause?]" |
| "no CI/CD" | "Deploys take 2 hours and break weekly — team avoids releasing on Fridays" |

### Vague deliverables

| Vibe | Spec |
|---|---|
| "a tool" | "[NEEDS CLARIFICATION: name the artifact — CLI binary, library, doc, plugin?]" |
| "a system" | "[NEEDS CLARIFICATION: name the components that ship]" |
| "a platform" | "[NEEDS CLARIFICATION: what concretely ships — API, dashboard, SDK?]" |
| "a framework" | "Three layers — Convention (markdown directory), Skill (/spectacular slash command), CLI (spectacular init)" |

### Soft constraints → hard constraints

| Vibe | Spec |
|---|---|
| "limited time" | "Ships by 2026-08-15; ~40 dev-hours/week available" |
| "small team" | "1 engineer + occasional design help (3 hrs/week)" |
| "tight budget" | "$500 total infrastructure spend through Q3" |

### Vague milestones

| Vibe | Spec |
|---|---|
| "MVP" | "[NEEDS CLARIFICATION: what does MVP demo look like? Which user can do which thing?]" |
| "beta" | "10 invited users complete the core flow without bugs by 2026-07-15" |
| "v1" | "Tagged v1.0 release, install path documented, 1 external user onboarded" |

## Review gate checks (in addition to base)

| # | Check | How |
|---|---|---|
| 4 | Success criteria is measurable | Slot 5 contains ≥1 number AND ≥1 verb (any inflection — see Verb detection) AND ≥1 date/timeframe |
| 5 | No vague-word hits in critical slots | Slots 2, 5, 8 contain none of the vague-word list (with hyphen-aware tokenization — see Tokenization rule). Vision (1) and Deliverable (4) exempt. |
| 6 | Non-goals are specific | Slot 6 has ≥2 items, none of which are vague meta-terms |
| 7 | Single primary user | Slot 3 reads as a singular role + situation + constraint |
| 8 | Deliverable names artifacts | Slot 4 contains at least one concrete noun (binary, library, doc, plugin, format) — not just generic categories |
| 9 | Vision is non-trivial | Slot 1 is at least 2 sentences OR 30+ words. No vague-word check. |
| 10 | Kit-required slots filled | If PRD frontmatter has `kit: <name>`, every kit slot with `required: true` must be non-empty. See "Kit-aware gate" below. |

### Kit-aware gate

After base + checks 4-9 run, the engine reads the PRD's frontmatter for `kit: <name>`:

- **No kit declared** — gate stops here. Base-only PRDs are valid; no kit checks run.
- **Kit declared** — load `templates/prd/kits/<name>.md` (project-local override wins). For each entry in the kit's `adds-slots`:
  - If `required: true` — slot must exist and be non-empty in the PRD, else fail with "Kit '<name>' requires '<slot>' but it is empty/missing"
  - If `required: false` — skip; empty optional slots are fine
- Kit `modifies-slots` are **not** gate checks — they're prompt guidance applied at grill time, not validation criteria.

Universal base checks (placeholder, clarification, frontmatter) ALWAYS run regardless of kit. Kits never disable base checks.

### Vague-word list (slots 2, 5, 8 only)

`fast`, `slow`, `intuitive`, `simple`, `easy`, `scalable`, `seamless`, `great`, `good`, `nice`, `clean`, `flexible`, `robust`, `powerful`, `elegant`, `modern`, `smart`, `lightweight`, `heavyweight`, `solid`, `polished`.

**Tokenization rule** — vague-word matching must preserve hyphenated compounds as single tokens. Tokenize with `[a-z]+(?:-[a-z]+)*` (or equivalent), not `\b\w+\b`. Examples:
- `smart-init` → one token `smart-init` (does NOT hit `smart`)
- `doc-writer` → one token `doc-writer`
- `kits-as-plugins` → one token `kits-as-plugins`
- `the smart approach` → tokens `the`, `smart`, `approach` (DOES hit `smart`)

Compound identifiers (slugs, kebab-case names, hyphenated terms) are almost never vague-language uses — preserving them avoids false positives on request slugs, package names, and similar.

### Non-goals meta-term blocklist (slot 6)

`scope creep`, `feature bloat`, `over-engineering`, `complexity`, `tech debt` (without context).

These describe what you *don't want*, not what you're *excluding*. Push for concrete: "Not a multi-agent pipeline" beats "scope creep".

### Singular-user heuristic (slot 3)

Fails if:
- First sentence uses a bare plural noun ("Users", "Customers", "Developers", "Teams") without a qualifier
- Slot lists multiple distinct user types (e.g. "Both PMs and engineers")

Heuristic, not strict — user can override during grill.

### Verb detection (slot 5)

Slot 5 must contain an action verb. Match must accept **common inflections**, not just the bare stem, since real writing rarely uses bare infinitive form.

Approach (any of):
- **Stem-with-suffix regex** — `\b<stem>(s|es|ed|ing|d)?\b` for each verb in the list
- **Inflection enumeration** — list common inflections explicitly (`pass`, `passes`, `passed`, `passing`)
- **Lemmatizer** — overkill for v1; defer to v2 if false-positive/false-negative rate is high

Reference verb list (extend as needed): `open`, `complete`, `return`, `reach`, `ship`, `launch`, `pass`, `exercise`, `produce`, `run`, `write`, `verify`, `deliver`, `release`, `validate`, `test`, `confirm`, `demonstrate`, `achieve`, `hit`, `meet`.

Example matches that must pass:
- "the PRD **passes** review" ✓ (`pass` + s)
- "3 doc types are **exercised**" ✓ (`exercise` + d)
- "users **complete** the flow" ✓ (bare stem)
- "we **shipped** v1" ✓ (`ship` + ed)

### Concrete-deliverable heuristic (slot 4)

Fails if the text is only generic-category nouns: `tool`, `system`, `platform`, `framework`, `solution`, `service` — without naming a concrete artifact.

Pass: "CLI binary + Claude skill + Bash installer" (three concrete artifacts).
Pass: "Weekly newsletter, 12 issues, published Tuesdays" (format + cadence).
Fail: "a tool to help with PRDs" (no artifact named).

A deliverable can use generic words *alongside* concrete artifacts — `"A CLI tool (Bash binary)"` passes; `"A tool"` fails.

### Vision non-trivial heuristic (slot 1)

Fails only if Vision is a single short sentence under 30 words and not split into 2+ sentences. Vision tolerates abstraction; this check ensures it's actually written, not skipped.

## Related

- [[doc-registry]] — registry entry referencing this file
- [[grill]] — consumes the slot prompts + mini-refine patterns from this file
- [[refine]] — consumes the vibe→spec tables for full refine passes
- [[review]] — consumes the gate checks from this file
- [[prd-grill]] — legacy reference (now superseded by `grill.md` + this file; kept for v1 backwards compat during migration)
- [[prd-refine]] — legacy reference
- [[prd-review]] — legacy reference
- [[kits-as-plugins]] — how kits extend this doc's slots
