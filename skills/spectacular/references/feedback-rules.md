---
doc-id: feedback
mode: index
location: .spectacular/feedbacks/
entries-dir: .spectacular/feedbacks/
alt-location: .spectacular/requests/<slug>/feedbacks/
scope: project-wide-and-request-scoped
template: templates/feedback/entry.md
snapshot-on-edit: false
summary: "Prototyping-mode human-feedback loop — capture insights, durable preferences, and use-case validation. NOT a benchmark/eval harness."
status: active
---

# Feedback Rules

Soft-folder database of prototyping-mode feedback entries. There is **no top-level `FEEDBACK.md` index file** — folder listing is the canonical view. Each entry is a fully self-contained markdown file.

**Mode: `index`** (no regenerated index file in v1.6.0). Entry files at `entries-dir`. May also live request-scoped under `requests/<slug>/feedbacks/`.

**Verbs:**
- `grill` → run a full 5-step feedback loop interactively (pick target → craft proposal → ask user → capture → decide). The primary user-facing verb. Full spec in [[feedback-loop]].
- `refine` → reopen a `resolved` entry and revise the insight or decision; rare.
- `review` → validate entry frontmatter across all entries; flag stale `open` (> 30 days), orphan back-refs, missing required fields.

**Mutator verbs (CLI, not skill):**
- `spectacular feedback-loop new <target>` — scaffold one entry with stub frontmatter + section headers (status `open`)
- `spectacular feedback-loop list [--status <state>]` — list entries across both locations
- `spectacular feedback-loop resolve <slug> --next-action <action>` — close entry, optionally auto-promote to memory
- `spectacular feedback-loop archive <slug>` — move to `.spectacular/archive/feedbacks/<year>/`

**Aliases (hidden):** `iterate`, `experiment`, `test`, `probe`, `try` — accepted by CLI dispatch, route to `cmd_feedback_loop`, not documented in `--help`.

**Snapshot-on-edit: false** — entries are durable insight records, not versioned canonical docs.

**Entry frontmatter (required shape):**

```yaml
---
type: feedback
target: <short description>
scope: skill | substrate | convention | doc-type | request
status: open | resolved | parked
opened: YYYY-MM-DD
resolved: YYYY-MM-DD | null
proposal_summary: "<one-line>"
next_action: ship-as-is | new-request:<slug> | park | memory:<entry> | tbd
request: <slug> | null
spawned_request: <slug> | null
promoted_to: memory/<slug>.md | null
related: []
---
```

**Required body sections:** Target, Hypothesis / hunch, Proposal, Question asked, User response, Insight, Decision. See [[feedback-loop]] for full template + examples.

**Proactive surfacing rules:** Skill may surface a feedback-loop offer at three checkpoints only:
1. Milestone completion (`TASKS.md` tick)
2. Request status flip to `review`
3. End of `spectacular archive <slug>` flow

Never mid-flow. Never unsolicited. Single short prompt; user accepts or declines.

**Auto-promotion to memory:** When resolution captures a durable preference signal, the skill explicitly confirms in its closing turn and (on accept) calls `spectacular remember` with the distilled preference, tagging `feedback,<scope>`. Sets `promoted_to:` on the feedback entry. No silent writes.

**Doctor area:** `spectacular doctor feedback` is judgment-only (no `--fix`). Flags:
- `open` entries with `opened` > 30 days → warning
- Missing required frontmatter fields → warning
- Orphan back-refs (PLAN.md mentions a feedback file that doesn't exist, or vice versa) → warning

**Not** what this is:
- Not a benchmark or evals harness
- Not VERIFY.md (request-scoped conformance pass)
- Not `review` (doc-quality gate against PRINCIPLES)
- Not a substitute for memory — feedback is the *acquisition mechanism*, memory is the *durable store* for preferences that emerge

**Relationship to VERIFY.md:** orthogonal. VERIFY answers "did we ship what PLAN said?" — confirmatory, terminates at `verified`. Feedback-loop answers "was that the right thing to ship?" — exploratory, never terminates. Both can run on the same request without overlap.

**Related:** [[feedback-loop]], [[verify]], [[review]], [[memory-rules]], [[doc-index]], [[archive]]
