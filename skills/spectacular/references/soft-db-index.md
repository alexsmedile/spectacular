---
description: Canonical routing index for Spectacular's soft-DB collections — role, purpose, structure, and use/boundary rules for each. The skill reads this to pick the right collection.
when_to_use: Deciding where a piece of information belongs (memory vs decision vs fix vs idea…), or explaining what a collection is for.
---

# Soft-DB Index — the collections and how to route between them

Spectacular stores operational knowledge in **soft-DB collections**: folders of `.md` entries with frontmatter, under `.spectacular/`. This is the single routing index — it answers *what each collection is for, when to use it, and when NOT to*. Per-collection detail (entry frontmatter, lifecycle, doctor checks) lives in each `references/<id>-rules.md`; dispatch lives in that file's frontmatter. **This doc is the map; the rules files are the territory.** Wayfinding adds `questions/` as the first explicit open-loop store; research and spike records join this table when their CLI substrates ship.

**What makes something a soft-DB collection:** a folder of individually-addressable `.md` entries (not a single canonical doc), each with frontmatter the skill reads as a signal layer, appended over time and never overwritten in place. All are **committed to git** (team-visible) and iterated via the CLI's `_iter_md <collection>`.

Before routing to any collection, apply [[discovery-protocol]]'s cheapest-sufficient-answer gate. Clear paths create no discovery node; a generic output does not become a new “artifact” entity.

---

## The routing table

| Collection | Role | Purpose — why it exists | Entry / structure | Write verb | Rules |
|---|---|---|---|---|---|
| **`memories/`** | Durable fact | Long-lived operational learning — standing preferences, conventions, "always do X" that outlives any one task | `memories/index.md` index + `M<N>-<slug>.md` entries | `spectacular remember "…"` | [[memory-rules]] |
| **`decisions/`** | Why we chose | ADR log — the *rationale* behind a choice between options, so it isn't relitigated | `decisions/index.md` index + `D<N>-<slug>.md` entries | `spectacular decide "…"` | [[decisions-rules]] |
| **`sessions/`** | Work boundary | Time-log of start/end, auto-linking the decisions + memories captured within | `sessions/index.md` index + `S<N>-<slug>.md` entries | `spectacular session start\|end` | [[sessions-rules]] |
| **`ideas/`** | Local compatibility workbench | Private/offline capture or locally useful refinement of a speculative "might build this". Prefer the actual shared or development destination when known; no execution commitment. | `<slug>.md` entries (no index) | `spectacular idea new` → `promote --to …` | [[idea-rules]] |
| **`visions/`** | Human-approved direction | Material product, experience, workflow, or system direction needs grounded alternatives and concrete human reaction before specification | `<slug>/VISION.md` + `fragments/` + linked `evidence/` | `spectacular imagine` → `vision react|propose|approve` → agentic `derive` | [[vision-rules]], [[imagine]] |
| **`questions/`** | Active ambiguity | A blocker or choice that must remain visible until answered, commonly requiring human product/business judgment | active `questions/QUE-NNN-<slug>.md`; resolved history in `archive/questions/` | `spectacular question new\|resolve` | [[question-rules]] |
| **`research/`** | Evidence gathering | Read-only investigation of facts/options; verified sources clear fog without making a product decision | `research/RES-NNN-<slug>.md` entries | `spectacular research new\|resolve` | [[research-rules]] |
| **`spikes/`** | Feasibility evidence | Human-authorized throwaway experiment; durable output is evidence, not production code | `spikes/SPK-NNN-<slug>.md` entries | `spectacular spike new\|resolve` | [[spike-rules]] |
| **`feedbacks/`** | Prototyping signal | Human-feedback loop during prototyping — "was this the right thing to ship?" | `feedbacks/index.md` index + `<slug>.md` entries | `spectacular feedback-loop new` | [[feedback-rules]] |
| **`audits/`** | Diagnosis-in-progress | Investigate a bug/quirk **before** planning a fix — understand & retrieve the real problem | `audits/index.md` index + `A<N>-<slug>.md` entries | `spectacular audit new` | [[audit-rules]] |
| **`fixes/`** | Self-learning corpus | Verified fixes with a **signature** — "we've seen this, here's how we fix it", reusable by future agents/projects | `fixes/index.md` index + `F<N>-<slug>.md` entries | `spectacular fix new` | [[fixes-rules]] |

---

## Reserved advanced engineering collections

Heavy engineering work may opt into `findings/`, `fixes/`, `bugs/`, `security/`, and `benchmarks/` during init. Reservation stabilizes the folder and canonical prefix without claiming a workflow exists.

| Collection | Reserved ID | Intended role | Current behavior |
|---|---|---|---|
| `findings/` | `FND-NNN` | Durable takeaway produced by research, audits, or user sessions | Reserved only; use `fnd1`; ambiguous `f1` is refused |
| `fixes/` | `FIX-NNN` | Padded canonical remediation identity | Existing `F<N>` verified-fix ledger remains active until explicit migration |
| `bugs/` | `BUG-NNN` | Known defect or regression | Reserved only; prefer `bug1` because `b<N>` is also a roadmap build ID |
| `security/` | `SEC-NNN` | Security vulnerability or threat finding | Reserved only; folder is singular `security/`, never `securities/` |
| `benchmarks/` | `BMK-NNN` | Performance, load, or profiling result | Reserved only; not the qualitative feedback loop |

Create these explicitly with `spectacular init --with findings,fixes,bugs,security,benchmarks`. Until a separately approved workflow ships, do not allocate their reserved IDs, infer lifecycle states, or route autonomous writes into them.

---

## When to act — the trigger for each collection

Knowing *where* something belongs (above) is only half of it. The agent must also recognise *the moment to write*. Each collection has a named prompt-moment; act on it **proactively** (surface a short offer), never silently and never autonomously for irreversible writes.

| Collection | Prompt-moment — when the agent should act | Proactive? | Confirm first? |
|---|---|---|---|
| **`memories/`** | A durable lesson surfaces (a non-obvious blocker, a "we should always…" preference, a reusable pattern). Strongest at **archive time** — archive.md Step 3 forces the review. | offer | **yes** — never autonomous |
| **`decisions/`** | A choice is made between real alternatives *and it's architectural / will be re-questioned*. Capture at the moment of choosing, or at archive if it slipped through. | offer | yes |
| **`sessions/`** | Explicit — start when beginning focused work, end when stopping. The one collection with mechanical boundaries. | on request | no |
| **`ideas/`** | The human needs private/offline capture or local refinement of a "maybe we should…". **User-initiated** — the agent captures on ask, doesn't manufacture ideas or duplicate shared discussion. | on request | no (capture is cheap/reversible) |
| **`visions/`** | The destination is materially unsettled and concrete alternatives/fragments would produce better human direction than prose interrogation. Never because work merely has a UI. | suggest when uncertainty is material; user initiates/accepts | **yes** for whole-Vision approval/rejection; reactions follow explicit human input |
| **`questions/`** | Work encounters missing information or a product/business fork that must not be guessed through. Surface at the next session start unless deferred. | yes | resolution: **yes** when `requires_user_input: true` |
| **`research/`** | A bounded external or repository fact is missing and can be gathered read-only. | yes when it clears a declared dependency | no for execution; evidence required to resolve |
| **`spikes/`** | Feasibility is uncertain enough that a throwaway experiment is cheaper than speculative implementation. | propose | **yes** before execution |
| **`feedbacks/`** | Three checkpoints only (never mid-flow): a milestone ticks in TASKS.md, a request enters `review`, or during `archive`. Single short offer. | offer at checkpoints | yes |
| **`audits/`** | A bug is reported whose root cause is unclear, spans multiple sites, or isn't yet reproduced. See [[bug-workflow]] Step 1 — **don't** open an audit for a trivial, understood one-liner. | judge per bug | no (scratchpad) |
| **`fixes/`** | A bug is **resolved AND verified** *and* carries reusable knowledge (non-obvious cause, a recurring class). Also offered at **archive time** if the request fixed a bug. | offer | no (but require `--verified-by`) |
| requests → **`archive/`** | A request reaches `verified` and the human confirms. Then `spectacular archive <slug>` — **never** manual `mv`. | offer | **yes** — irreversible-class |

**The archive checkpoint is where three collections converge** — see [[archive]]: on archiving a request the agent proposes spec-sync, **memory** entries, and (if it was a bug) a **fix** entry / **audit** resolution. That single moment catches the captures that otherwise get forgotten.

**Golden rule for the "when":** reversible + cheap (audit note, session, idea) → just do it on the natural trigger. Irreversible or team-visible-and-permanent (memory, decisions, archive) → **propose, human confirms, then write.** This mirrors the guardrail in `.spectacular/AGENTS.md`.

## Use / boundary rules — the "not this" that prevents mis-routing

The value of separate collections is lost if entries land in the wrong one. Each boundary below is the question to ask when routing:

- **`memories/` — a standing fact, not an event.** "We always deploy from `main`" → memory. "We deployed today" → not memory (that's a session note). If it has a date and won't recur, it's not memory.
- **`decisions/` — a choice with alternatives, not a fact or a fix.** Records *why A over B*. If there were no real alternatives, it's a fact (`memories/`). If it's a bug resolution, it's `fixes/`.
- **`sessions/` — a time bracket, not content.** Holds *links* to what happened (decisions, memories), not the substance itself.
- **`ideas/` — pre-commitment, not a backlog.** An idea may be grilled, researched, revised, and given a *working* implementation plan while its decisions remain open. It becomes a `request/` only when a human accepts an execution outcome and durable coordination is warranted. ROADMAP's Icebox holds version-tied vision items; ideas have no version.
- **`questions/` — an explicit open loop, not a parked possibility.** It blocks or informs another entity and expects an answer. Out-of-scope inspiration belongs in `ideas/`; resolution archives the QUE, and only a genuine settled choice additionally belongs in `decisions/`.
- **`research/` — evidence, not a choice.** It narrows options and may support a decision; it does not become the ADR itself.
- **`spikes/` — disposable code for knowledge, not implementation.** Preserve evidence/outcome; production work begins only from an approved spec and normal request.
- **`feedbacks/` — post-ship prototyping signal, not a bug.** Answers "right thing to build?" about something already built. A malfunction is a bug (`audits/`/`fixes/`), not feedback. Not a benchmark/eval harness either.
- **`audits/` — diagnosis, not the fix.** Ends in a *disposition*, not a verified change. If root cause is already obvious and the fix is one site, skip the audit entirely (see [[bug-workflow]]).
- **`fixes/` — verified + reusable, not a symptom tracker.** Logged only once resolved *and* verified, and only when it carries reusable knowledge. A typo fix teaches nothing → no entry.

---

## Not a soft-DB collection (adjacent, don't confuse)

| Thing | What it actually is | Why not a collection |
|---|---|---|
| **`visions/`** | An opt-in, pre-request experience-direction workspace: `VISION.md`, reacted fragments, and supporting evidence | It is a bounded approval workspace whose approved result may derive a draft SPC, not an append-only knowledge ledger. |
| **`requests/`** | The **structured request lifecycle** — each `requests/<slug>/` holds PLAN.md (owns `status:`), TASKS.md, and optional VERIFY.md | Not append-only `.md` entries; it's the work-execution unit with lifecycle state, not a knowledge store. It's where a bug *folds into a plan*. |
| **`debugs/`** | Ephemeral machine-readable run traces (`debugs/<slug>/`) | Raw JSON traces generated by automated agent steps, rather than clean, hand-authored, append-only markdown documents. |
| **Canonical docs** (PRD, SPEC, PRINCIPLES, ARCHITECTURE, ROADMAP, STACK, AGENTS, POLICY, PERSONAS) | Single authoritative documents, versioned via snapshots | One file, overwritten-via-snapshot, not a folder of entries. Cataloged in [[doc-index]]. |
| **`archive/`** | Completed requests + historical snapshots | Terminal storage, never read during normal operation. |
| **Prototype artifact** | A Vision/feedback/request-owned mock used for human validation | It inherits its owner's context; `PRT` remains reserved rather than creating a parallel lifecycle. |
| **Tracer bullet** | Approved `SPC` execution mode for a thin production vertical slice | It is retained production code, not discovery evidence or a knowledge collection. |
| **Technical debt** | An execution obligation routed to a request, roadmap candidate, idea, and sometimes a linked decision | A separate `debt/` ledger would duplicate ownership and priority state. |

---

## The bug-lifecycle sub-loop (audit → plan → fix)

Three of these collections form one workflow — Spectacular's self-learning loop for bugs. Full routing (including the "just-fix, no ceremony" fast path) is in [[bug-workflow]]:

```
seen it? (grep fixes/ signatures) → ceremony call (audit vs just-fix)
        → resolve (one-line fix · fold into a request/plan · became fix) → log if reusable (fixes/)
```

- **`audits/`** = understand (diagnose, retrieve the real problem)
- **`requests/`** = plan (strategy + tasks, when the bug needs real planning)
- **`fixes/`** = remember-for-next-time (signed, reusable, cross-project)

---

**Related:** [[doc-index]] (the full doc catalog incl. canonical docs), [[bug-workflow]] (audit/fix routing), each `<id>-rules.md`, and `.spectacular/AGENTS.md` § context-loading (what to read per task).
