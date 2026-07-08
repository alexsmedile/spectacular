---
status: archived
priority: medium
owner: alex
updated: 2026-05-25
target_version: v1.6.0
summary: "Add a `feedback-loop` mode — a deliberate human-feedback strategy for prototyping-stage spectacular. The skill picks a target, drafts a proposal, asks the user a structured question, captures the response as durable signal, and decides next action. Distinct from VERIFY.md (request-scoped conformance) and `review` (doc-quality pass). Not a benchmark/eval harness — explicitly avoids that framing."
related:
  - ../../PRINCIPLES.md
  - ../../ARCHITECTURE.md
  - ../../roadmaps/index.md
feedback:
  - feedback/2026-05-25-feedback-loop-cli-ergonomics-after-m0-m4.md
archived: 2026-05-25
---

# Plan — feedback-loop

## 1. Goal

Spectacular is in prototyping. Conventions, modes, substrate decisions are being made on intuition and dogfooding — which is fine, but the signal is currently lossy: insights surface mid-conversation, sometimes land in TODO.md, sometimes evaporate. There's no deliberate practice for **acquiring feedback on the system itself**.

Introduce a `feedback-loop` mode that operationalizes this:

1. **Pick a target** — recently-shipped behavior, a fuzzy convention, an untested edge, or a user-flagged hunch.
2. **Craft a proposal** — concrete: scenario, variants if comparing, the question being asked, what signal we're hoping for.
3. **Ask the user** — structured question (AskUserQuestion with previews when comparing artifacts; free-form otherwise).
4. **Capture the response** — written to `.spectacular/feedback/<date>-<slug>.md` or promoted to memory if it's a durable preference.
5. **Decide next action** — ship-it / draft-a-request / park-revisit. Captured in the feedback file's frontmatter.

This is a **strategy for acquiring feedback, knowledge, insights, and use-case validation** — not verification, not a benchmark.

## 2. Constraints

- **No benchmark framing.** The word "evals" is explicitly avoided in user-facing docs, CLI help, frontmatter, and code identifiers. This is a feedback loop, not a benchmark — "evals" carries HumanEval/MMLU/accuracy-% baggage that pulls the wrong way.
- **Orthogonal to VERIFY.md.** VERIFY answers "did we ship what PLAN said?" — request-scoped, confirmatory, terminates at `verified`. Feedback-loop answers "was that the right thing to ship?" — system-scoped, exploratory, never terminates. Both can run on the same change without overlap.
- **Orthogonal to `review` mode.** `review` is a doc-quality pass against PRINCIPLES; feedback-loop is system-fitness probing of behavior in the wild.
- **Lightweight by default.** A feedback session should be possible in 1-2 turns. Heavy ceremony kills the practice.
- **Durable when it matters.** Single-session insights → memory; multi-session investigations → a new request. The mode itself doesn't own long-running state.

## 3. Scope

### In
- New skill mode: `feedback-loop` (canonical name)
- CLI verb aliases: `iterate`, `experiment`, `test`, `probe`, `try` — all route to the same handler
- New ref doc: `skills/spectacular/references/feedback-loop.md`
- New substrate folder: `.spectacular/feedback/` with a documented file shape (`<YYYY-MM-DD>-<slug>.md`)
- Verb-on-docs form: `spectacular prd feedback-loop`, `spectacular plan feedback-loop`, etc. — scopes the probe to a specific doc
- Doctor area: `doctor feedback` — surfaces stale unresolved feedback entries (e.g. status: open for > 30 days)
- Update ARCHITECTURE.md to document `feedback/` folder + lifecycle
- Update PRINCIPLES.md or AGENTS.md with a "feedback ≠ verification ≠ benchmark" note so the distinction is preserved

### Out
- Automated/scheduled feedback prompts (no cron, no "ask me every Monday")
- A feedback-driven changelog or release-notes generator
- Cross-project feedback aggregation
- Anything resembling a quantitative score, rating, or grade
- Integration with external survey tools

## 4. Locked decisions (grilled 2026-05-25)

1. **Lifecycle:** three states — `open` / `resolved` / `parked`. `open` = awaiting response or decision; `resolved` = decision recorded; `parked` = revisit later, not currently actionable. Doctor flags `open` entries older than 30 days.
2. **Proactivity:** **skill surfaces feedback-loop opportunities at request checkpoints** — specifically when a major milestone completes (`PLAN.md` ticks a milestone off), when a request enters `review`, and when a request is archived. Never mid-flow / unsolicited. The user can always decline.
3. **Canonical name:** `feedback-loop` (with hyphen). Aliases `iterate`, `experiment`, `test`, `probe`, `try` route to the same handler but are **hidden routing only** — they work if typed but don't appear in `--help`. Only `feedback-loop` is documented as the official name.
4. **Memory promotion:** **auto-promote durable preferences.** When a feedback entry resolves with a stable preference signal ("I always want X", "Y is the right default"), the skill auto-creates a memory entry of the appropriate type (feedback / user / project) and back-links it from the feedback file via `promoted_to: memory/<slug>.md`. The skill must explicitly confirm the promotion in the loop's final turn — no silent writes.
5. **Archive policy:** **live forever in `.spectacular/feedback/`** unless explicitly archived. `spectacular feedback-loop archive <slug>` is a manual verb. No automatic movement. User curates. Doctor's 30-day stale-open warning is the only ambient signal.
6. **Request back-refs:** **bidirectional within the same request by default.** A feedback entry is born scoped to the active request (if any) and lives at `.spectacular/requests/<slug>/feedback/<date>-<slug>.md`. The request's PLAN.md gets a `feedback:` list pointing to the entries; each entry's frontmatter carries `request: <slug>`. **The loop happens inside the same request** unless the decision is `new-request:<slug>` — then a new request is spawned and the feedback entry records `spawned_request: <slug>`, with the new request's PLAN.md carrying `spawned_by_feedback: <path>`. Out-of-request feedback (probing the system itself) lives at the top-level `.spectacular/feedback/` and has no `request:` field.

### Final file layout

```
.spectacular/
├── feedback/                                  # System-level feedback (no request scope)
│   ├── 2026-05-25-grill-each-feels-heavy.md
│   └── ...
└── requests/
    └── <slug>/
        ├── PLAN.md
        ├── TASKS.md
        └── feedback/                          # Request-scoped feedback (the common case)
            ├── 2026-05-26-m2-skill-routing.md
            └── ...
```

## 5. File shape

`.spectacular/feedback/2026-05-25-grill-each-feels-heavy.md` (top-level) or `.spectacular/requests/<slug>/feedback/2026-05-25-...md` (request-scoped):

```markdown
---
type: feedback
target: grill-each mode
scope: skill                              # skill | substrate | convention | doc-type | request
status: open                              # open | resolved | parked
opened: 2026-05-25
resolved:                                 # date when status flips to resolved
proposal_summary: "Compare grill-wide vs grill-each on a 10-slot PRD"
next_action: tbd                          # ship-as-is | new-request:<slug> | park | memory:<entry>
request: <slug>|null                      # set when feedback is request-scoped
spawned_request: null                     # set when next_action = new-request:<slug>
promoted_to: null                         # path to memory entry if auto-promoted
related: []
---

# Feedback — grill-each feels heavy on long PRDs

## Target
grill-each mode (introduced v1.4.0)

## Hypothesis / hunch
On PRDs with >8 slots, grill-each becomes tedious — user has to answer N slot-level questions in sequence with no escape hatch.

## Proposal
[Two grilled PRDs on the same input — A grilled wide, B grilled each-slot. Both attached as artifacts.]

## Question asked
"Which lands closer to what you wanted, and where did each one feel wrong?"

## User response
[captured verbatim or summarized]

## Insight
[what we learned that's durable]

## Decision
[next_action explained in prose]
```

## 6. Milestones

### M0 — Grill the request ✅
- All 7 open questions resolved (see §4)
- Frontmatter schema locked
- Command surface locked: `feedback-loop` canonical, aliases hidden
- Target version: v1.6.0 (alongside `memory-protocols` and `snapshot-tidy`)

### M1 — Substrate + ref doc
- Create `skills/spectacular/references/feedback-loop.md` with the full mode spec
- Add `.spectacular/feedback/` to ARCHITECTURE.md
- Add the PRINCIPLES.md note distinguishing feedback / verification / benchmark
- Add `feedback/` to the `init` scaffold (empty folder + `.gitkeep` or a README)

### M2 — Skill mode
- Implement `feedback-loop` mode in SKILL.md routing
- Wire the 5-step loop (target → proposal → ask → capture → decide)
- Wire verb-on-doc form (`spectacular prd feedback-loop` scopes the probe)
- Aliases (`iterate`, `experiment`, `test`, `probe`, `try`) route to the same handler

### M3 — CLI surface (if any CLI work needed)
- `spectacular feedback-loop new <target>` — scaffold a feedback file
- `spectacular feedback-loop list` — show open feedback entries
- `spectacular feedback-loop resolve <slug>` — close a feedback entry with a decision
- Aliases plumbed at CLI layer too

### M4 — Doctor area
- `doctor feedback` — flag entries with `status: open` older than 30 days
- `--fix` n/a (judgment-only area, no mechanical repair)
- Update `references/doctor.md`

### M5 — Dogfood
- Run 2-3 real feedback-loop sessions on already-shipped behavior (candidates: grill-each ergonomics, soft-db memory mutators, pack discovery UX)
- Capture them in `.spectacular/feedback/`
- Verify the loop feels lightweight in practice; iterate the ref doc based on friction

### M6 — Release
- CHANGELOG entry under Added
- Bump manifests
- Mention in README's "What's in spectacular" section if appropriate

## 7. Non-goals (locked)

- Benchmarks, scores, accuracy metrics, or any quantitative grading
- Automation that asks the user without being invoked
- A separate "feedback dashboard" — folder listing + doctor is enough
- Replacing or subsuming `review` or `VERIFY.md`
- Cross-project / cross-repo feedback aggregation

## 8. Dependencies

- None hard. Composes with everything; doesn't gate or get gated by other planned requests.
- Soft synergy with `memory-protocols` (v1.6.0) — feedback-to-memory promotion is cleaner if memory protocols are in place first. Not a blocker; can ship before and refine after.

## 9. References

- `.spectacular/PRINCIPLES.md` — guardrails for what gets shipped
- `skills/spectacular/references/verification.md` — VERIFY.md contract (sister concept, orthogonal axis)
- `skills/spectacular/references/review.md` — review mode (different axis)
- `TODO.md` § "Feedback loop (prototyping mode)" — original capture
