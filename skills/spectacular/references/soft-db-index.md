---
description: Canonical routing index for Spectacular's soft-DB collections — role, purpose, structure, and use/boundary rules for each. The skill reads this to pick the right collection.
when_to_use: Deciding where a piece of information belongs (memory vs decision vs fix vs idea…), or explaining what a collection is for.
---

# Soft-DB Index — the collections and how to route between them

Spectacular stores operational knowledge in **soft-DB collections**: folders of `.md` entries with frontmatter, under `.spectacular/`. This is the single routing index — it answers *what each collection is for, when to use it, and when NOT to*. Per-collection detail (entry frontmatter, lifecycle, doctor checks) lives in each `references/<id>-rules.md`; dispatch lives in that file's frontmatter. **This doc is the map; the rules files are the territory.**

**What makes something a soft-DB collection:** a folder of individually-addressable `.md` entries (not a single canonical doc), each with frontmatter the skill reads as a signal layer, appended over time and never overwritten in place. All are **committed to git** (team-visible) and iterated via the CLI's `_iter_md <collection>`.

---

## The routing table

| Collection | Role | Purpose — why it exists | Entry / structure | Write verb | Rules |
|---|---|---|---|---|---|
| **`memory/`** | Durable fact | Long-lived operational learning — standing preferences, conventions, "always do X" that outlives any one task | `MEMORY.md` index + `<slug>.md` entries | `spectacular remember "…"` | [[memory-rules]] |
| **`decisions/`** | Why we chose | ADR log — the *rationale* behind a choice between options, so it isn't relitigated | `DECISIONS.md` (flat prose) **or** `+ decisions/D<N>.md` (index mode) | `spectacular decide "…"` | [[decisions-rules]] |
| **`sessions/`** | Work boundary | Time-log of start/end, auto-linking the decisions + memories captured within | `SESSIONS.md` index + `<slug>.md` entries | `spectacular session start\|end` | [[sessions-rules]] |
| **`ideas/`** | Pre-commitment spark | Speculative "might build this" — no lifecycle, no scope, may never matter again | `ideas/<slug>.md` (folder) | `spectacular idea new` → `promote` | [[idea-rules]] |
| **`feedback/`** | Prototyping signal | Human-feedback loop during prototyping — "was this the right thing to ship?" | `feedback/<slug>.md` (+ request-scoped `requests/<slug>/feedback/`) | `spectacular feedback-loop new` | [[feedback-rules]] |
| **`audit/`** | Diagnosis-in-progress | Investigate a bug/quirk **before** planning a fix — understand & retrieve the real problem | `audit/A<N>.md` (folder, auto-numbered) | `spectacular audit new` | [[audit-rules]] |
| **`fixes/`** | Self-learning corpus | Verified fixes with a **signature** — "we've seen this, here's how we fix it", reusable by future agents/projects | `fixes/F<N>.md` (folder, auto-numbered) | `spectacular fix new` | [[fixes-rules]] |

---

## When to act — the trigger for each collection

Knowing *where* something belongs (above) is only half of it. The agent must also recognise *the moment to write*. Each collection has a named prompt-moment; act on it **proactively** (surface a short offer), never silently and never autonomously for irreversible writes.

| Collection | Prompt-moment — when the agent should act | Proactive? | Confirm first? |
|---|---|---|---|
| **`memory/`** | A durable lesson surfaces (a non-obvious blocker, a "we should always…" preference, a reusable pattern). Strongest at **archive time** — archive.md Step 3 forces the review. | offer | **yes** — never autonomous |
| **`decisions/`** | A choice is made between real alternatives *and it's architectural / will be re-questioned*. Capture at the moment of choosing, or at archive if it slipped through. | offer | yes |
| **`sessions/`** | Explicit — start when beginning focused work, end when stopping. The one collection with mechanical boundaries. | on request | no |
| **`ideas/`** | The human floats a "maybe we should…" with no commitment. **User-initiated** — the agent captures on ask, doesn't manufacture ideas. | on request | no (capture is cheap/reversible) |
| **`feedback/`** | Three checkpoints only (never mid-flow): a milestone ticks in TASKS.md, a request enters `review`, or during `archive`. Single short offer. | offer at checkpoints | yes |
| **`audit/`** | A bug is reported whose root cause is unclear, spans multiple sites, or isn't yet reproduced. See [[bug-workflow]] Step 1 — **don't** open an audit for a trivial, understood one-liner. | judge per bug | no (scratchpad) |
| **`fixes/`** | A bug is **resolved AND verified** *and* carries reusable knowledge (non-obvious cause, a recurring class). Also offered at **archive time** if the request fixed a bug. | offer | no (but require `--verified-by`) |
| requests → **`archive/`** | A request reaches `verified` and the human confirms. Then `spectacular archive <slug>` — **never** manual `mv`. | offer | **yes** — irreversible-class |

**The archive checkpoint is where three collections converge** — see [[archive]]: on archiving a request the agent proposes spec-sync, **memory** entries, and (if it was a bug) a **fix** entry / **audit** resolution. That single moment catches the captures that otherwise get forgotten.

**Golden rule for the "when":** reversible + cheap (audit note, session, idea) → just do it on the natural trigger. Irreversible or team-visible-and-permanent (memory, decisions, archive) → **propose, human confirms, then write.** This mirrors the guardrail in `.spectacular/AGENTS.md`.

## Use / boundary rules — the "not this" that prevents mis-routing

The value of separate collections is lost if entries land in the wrong one. Each boundary below is the question to ask when routing:

- **`memory/` — a standing fact, not an event.** "We always deploy from `main`" → memory. "We deployed today" → not memory (that's a session note). If it has a date and won't recur, it's not memory.
- **`decisions/` — a choice with alternatives, not a fact or a fix.** Records *why A over B*. If there were no real alternatives, it's a fact (`memory/`). If it's a bug resolution, it's `fixes/`.
- **`sessions/` — a time bracket, not content.** Holds *links* to what happened (decisions, memories), not the substance itself.
- **`ideas/` — pre-commitment, not a backlog.** No scope, no exit criteria. The moment it gains a plan it becomes a `request/`. ROADMAP's Icebox holds version-tied vision items; ideas have no version.
- **`feedback/` — post-ship prototyping signal, not a bug.** Answers "right thing to build?" about something already built. A malfunction is a bug (`audit/`/`fixes/`), not feedback. Not a benchmark/eval harness either.
- **`audit/` — diagnosis, not the fix.** Ends in a *disposition*, not a verified change. If root cause is already obvious and the fix is one site, skip the audit entirely (see [[bug-workflow]]).
- **`fixes/` — verified + reusable, not a symptom tracker.** Logged only once resolved *and* verified, and only when it carries reusable knowledge. A typo fix teaches nothing → no entry.

---

## Not a soft-DB collection (adjacent, don't confuse)

| Thing | What it actually is | Why not a collection |
|---|---|---|
| **`requests/`** | The **structured request lifecycle** — each `requests/<slug>/` holds PLAN.md (owns `status:`), TASKS.md, optional VERIFY.md, vision/ | Not append-only `.md` entries; it's the work-execution unit with lifecycle state, not a knowledge store. It's where a bug *folds into a plan*. |
| **Canonical docs** (PRD, SPEC, PRINCIPLES, ARCHITECTURE, ROADMAP, STACK, AGENTS, POLICY, PERSONAS) | Single authoritative documents, versioned via snapshots | One file, overwritten-via-snapshot, not a folder of entries. Cataloged in [[doc-index]]. |
| **`archive/`** | Completed requests + historical snapshots | Terminal storage, never read during normal operation. |

---

## The bug-lifecycle sub-loop (audit → plan → fix)

Three of these collections form one workflow — Spectacular's self-learning loop for bugs. Full routing (including the "just-fix, no ceremony" fast path) is in [[bug-workflow]]:

```
seen it? (grep fixes/ signatures) → ceremony call (audit vs just-fix)
        → resolve (one-line fix · fold into a request/plan · became fix) → log if reusable (fixes/)
```

- **`audit/`** = understand (diagnose, retrieve the real problem)
- **`requests/`** = plan (strategy + tasks, when the bug needs real planning)
- **`fixes/`** = remember-for-next-time (signed, reusable, cross-project)

---

**Related:** [[doc-index]] (the full doc catalog incl. canonical docs), [[bug-workflow]] (audit/fix routing), each `<id>-rules.md`, and `.spectacular/AGENTS.md` § context-loading (what to read per task).
