---
description: Prototyping-stage human-feedback acquisition — 5-step loop, distinct from verification/benchmark.
when_to_use: Running a feedback loop on a prototype (PRINCIPLES §9).
---

# Feedback Loop — prototyping-mode human-feedback acquisition

> Spectacular is in prototyping. This mode operationalizes deliberate human-feedback acquisition on the system itself — a strategy for **acquiring knowledge, insights, use-case validation, and durable preferences from the user**.

**Not** a benchmark harness. **Not** a verification pass. **Not** automated grading. The word "evals" is intentionally avoided — it carries HumanEval/MMLU/accuracy-% baggage that pulls the wrong way.

## What this mode does

A `feedback-loop` session is a structured 5-step interaction. The skill is the orchestrator; the user is the signal source.

1. **Pick a target** — recently-shipped behavior, a fuzzy convention, an untested edge of the substrate, or a user-flagged hunch. Targets are concrete (e.g. "grill-each ergonomics on long PRDs"), not abstract (e.g. "is grill good?").
2. **Craft a proposal** — write down the scenario, the variants if comparing two approaches, the question being asked, and what signal we're hoping to acquire. Proposals are the bulk of the work — a well-crafted proposal makes the user's response unambiguous to capture.
3. **Ask the user** — use `AskUserQuestion` with previews when comparing artifacts side-by-side; free-form question when capturing open-ended response. Keep ceremony low: 1-2 questions max per loop.
4. **Capture the response** — write to the feedback file. Include user's verbatim response (or a tight paraphrase), the insight that compounded out of it, and the decision.
5. **Decide next action** — one of: `ship-as-is`, `new-request:<slug>`, `park`, `memory:<entry>`. The decision is the durable output. Record it in frontmatter.

## How feedback-loop differs from sister modes

| | `feedback-loop` | `verify` (VERIFY.md) | `review` |
|---|---|---|---|
| **Question** | "Was that the right thing to ship?" | "Did we ship what PLAN said?" | "Does this doc meet PRINCIPLES?" |
| **Scope** | The system (or a request's behavior in the wild) | A single request | A single doc |
| **Mode** | Exploratory, open-ended | Confirmatory, closed-ended | Confirmatory, closed-ended |
| **Output** | Insight + decision (qualitative) | Pass/fail assertions | Gate-check findings |
| **Terminates?** | Never — same target can be probed again later | Yes — at `verified` | Yes — at next refine |
| **Lives where?** | `.spectacular/feedback/` or `.spectacular/requests/<slug>/feedback/` | `.spectacular/requests/<slug>/VERIFY.md` | In-place against the doc |

`feedback-loop` can pass while `verify` fails (we built the wrong thing correctly). `verify` can pass while `feedback-loop` reveals trouble (we built the right thing but it doesn't work in practice). They probe orthogonal axes.

## File shape

Feedback files live at one of two locations:

- **Request-scoped (common case):** `.spectacular/requests/<slug>/feedback/<YYYY-MM-DD>-<short-slug>.md`
- **System-level:** `.spectacular/feedback/<YYYY-MM-DD>-<short-slug>.md`

Frontmatter:

```yaml
---
type: feedback
target: <short description of what's being probed>
scope: skill | substrate | convention | doc-type | request
status: open | resolved | parked
opened: YYYY-MM-DD
resolved: YYYY-MM-DD | null
proposal_summary: "<one-line summary of what was proposed to the user>"
next_action: ship-as-is | new-request:<slug> | park | memory:<entry> | tbd
request: <slug> | null              # set when feedback lives under requests/<slug>/feedback/
spawned_request: <slug> | null      # set when next_action triggers a new request
promoted_to: memory/<slug>.md | null # set when feedback resolves to a durable preference
related: []
---
```

Body sections (all required):

```markdown
## Target
<what is being probed — be concrete>

## Hypothesis / hunch
<the starting suspicion, or "no prior, exploratory">

## Proposal
<the artifacts, scenarios, or variants shown to the user>

## Question asked
<the exact question put to the user>

## User response
<verbatim or tight paraphrase>

## Insight
<what we learned that's durable — distinct from the response itself>

## Decision
<next_action in prose, including any auto-promoted memory entry>
```

## Mode behavior

**Mode in doc-index taxonomy:** `feedback-loop` is a registered doc type (`doc-id: feedback`). It uses **mode `index`** — like `memory` and `sessions`, the canonical location is a folder of entries, not a single rolling document. There is no top-level `FEEDBACK.md` index file in v1.6.0; folder listing is enough.

**Verbs:**
- `grill` → run one full 5-step loop interactively (this is the primary verb)
- `refine` → reopen a `resolved` entry to revise the insight or decision (rare)
- `review` → validate frontmatter shape across all entries; flag stale `open` entries

**Mutator verb (CLI):**
- `spectacular feedback-loop new <target>` — scaffold a new entry (stub frontmatter + section headers, status `open`)
- `spectacular feedback-loop list [--status open|resolved|parked]` — list entries
- `spectacular feedback-loop resolve <slug> --next-action <action>` — close an entry with a decision; auto-promotes to memory if action signals durable preference
- `spectacular feedback-loop archive <slug>` — move to `.spectacular/archive/feedback/<year>/`

**Hidden aliases (route to `feedback-loop`):** `iterate`, `experiment`, `test`, `probe`, `try`. These work if typed but do not appear in `--help` output. Only `feedback-loop` is documented as the official mode name.

## Proactive surfacing

The skill **does not** spontaneously interrupt the user with feedback prompts. It surfaces feedback-loop opportunities **only at request checkpoints**:

- **Milestone completion** — when the user ticks a milestone in `TASKS.md`, the skill may offer: "Want to feedback-loop M<N> before moving on?"
- **Request enters `review`** — same offer, scoped to the request as a whole.
- **Archive flow** — at the end of `spectacular archive <slug>`, the skill may offer: "Anything worth feedback-looping before this leaves the active set?"

The user can always decline. The skill should make these offers low-friction — single short prompt, accept/decline, move on. Never two-step interrogation at a checkpoint.

The skill **may not** auto-surface feedback opportunities outside these three checkpoints. No "haven't probed X in N weeks" suggestions. No mid-flow nags.

## Auto-promotion to memory

When a feedback loop resolves with a **durable preference signal** — phrases like "I always want X", "Y is the right default", "never do Z" — the skill creates a corresponding memory entry and back-links:

1. Skill explicitly confirms the promotion in its closing turn: "This sounds like a durable preference — I'll write it to memory as well. OK?"
2. On confirm, skill calls `spectacular remember "<distilled preference>" --tag feedback,<scope>`.
3. Feedback file's frontmatter gets `promoted_to: memory/<slug>.md`.

No silent promotions. If the user declines, the feedback file stays as the only record.

Single-loop insights → feedback file only. Multi-loop convergent insights (same preference surfaces 2+ times) → strong candidate for promotion regardless.

## Doctor area

`spectacular doctor feedback` (judgment-only, no `--fix`):

- Scan `.spectacular/feedback/*.md` and `.spectacular/requests/*/feedback/*.md`
- For each entry with `status: open`:
  - If `opened` is more than 30 days ago → `warning`: "open feedback entry stale (opened YYYY-MM-DD)"
  - Suggested action: "resolve with `spectacular feedback-loop resolve <slug>` or mark `parked`"
- For each entry: validate frontmatter shape (`type: feedback`, `target`, `status`, `opened` required)
- Report orphan back-refs (request frontmatter mentions a feedback file that doesn't exist, or vice versa)

## Examples

### Example 1 — request-scoped feedback at milestone checkpoint

User just ticked M2 in `.spectacular/requests/grill-each-rewrite/TASKS.md`. Skill offers:

> M2 done — want to feedback-loop the grill-each rewrite before moving on?

User accepts. Skill picks target ("grill-each rewrite UX on long PRDs"), drafts proposal (run rewrite on the live PRD + show a before/after slot comparison), asks one question via AskUserQuestion with previews. Captures response to `.spectacular/requests/grill-each-rewrite/feedback/2026-05-25-m2-grill-each-ux.md`. Decision: `ship-as-is`. PLAN.md's frontmatter gets `feedback: [feedback/2026-05-25-m2-grill-each-ux.md]`.

### Example 2 — system-level feedback spawning a new request

During the same loop, user reveals "actually the bigger issue is that grill-loop never narrows down — it keeps re-asking the same slots". Skill captures this as a *separate* feedback entry at `.spectacular/feedback/2026-05-25-grill-loop-doesnt-narrow.md` (system-level, no request scope), decision `new-request:grill-loop-narrowing`. New request scaffolded; its PLAN.md frontmatter gets `spawned_by_feedback: ../../feedback/2026-05-25-grill-loop-doesnt-narrow.md`.

### Example 3 — auto-promoted to memory

User during loop says "I always want grill-each to skip empty optional slots without asking." Skill confirms: "Should I write this as a durable preference?" → user accepts. Skill runs `spectacular remember "grill-each should skip empty optional slots without asking" --tag feedback,grill`. Feedback file gets `promoted_to: memory/grill-each-skip-empty-optional-slots.md`. Both files now cross-reference.

## Related references

- [[verification]] — sister contract on the conformance axis
- [[review]] — sister contract on the doc-quality axis
- [[memory-rules]] — promotion target shape
- [[archive]] — feedback files are preserved during request archive
- [[doc-index]] — registry entry

## Implementation contract (for the CLI)

- Folder creation is lazy: `mkdir -p` only on first `feedback-loop new`
- Slug derivation: same `_slug_from_text` helper used by `remember`
- Index regeneration: not applicable — there is no top-level FEEDBACK.md file in v1.6.0
- Doctor: judgment-only area; no mechanical `--fix`
- Aliases plumbed at top-level dispatch in `cli/spectacular`; route to single `cmd_feedback_loop` handler
