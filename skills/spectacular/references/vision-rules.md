---
doc-id: vision
mode: index
location: .spectacular/requests/<slug>/vision/
entries-dir: .spectacular/requests/<slug>/vision/{stories,ui,arch}/
spine: .spectacular/requests/<slug>/vision/VISION.md
scope: per-request
template: templates/vision/
snapshot-on-edit: false
summary: "Imagination-backed planning artifact — a soft-folder of see-able fragments (stories/ui/arch) plus a narrative spine, produced by the generative `imagine` mode and derived into a draft PLAN."
status: active
---

# VISION Rules — the `vision` doc-type + the `imagine` mode

> **Request-scoped, optional.** A `vision/` folder lives inside one request (`requests/<slug>/vision/`). It is **not** scaffolded by `spectacular new` — only by `spectacular imagine <slug>`. Small/obvious requests skip straight to PLAN; `imagine` is for work whose *shape* the human wants to see before committing to milestones.

> **Thesis.** Spectacular is spec-driven: the human reacts to specification text. `imagine` makes it spec-driven **and imagination-backed** — the skill renders artifacts the human reacts to, then **derives the spec from what the human approved**. Specs become accountable to the approved vision, not authored in a vacuum. Full rationale: `archive/ideas/explore-mode.md` §0.

## Structure — spine + typed subfolders

`vision` is an `index`-mode soft-folder, but with one structural difference from idea/sessions/memory: the index is a **narrative spine** (`VISION.md`), not a bare catalog, and entries live in **typed subfolders**, not one flat dir.

```
requests/<slug>/vision/
├── VISION.md            # SPINE — narrative + regenerable manifest of fragment links
├── stories/             # one user story per file
│   └── 01-first-run.md
├── ui/                  # one UI/output mockup per file  ← scales to MANY fragments
│   └── dashboard.md
└── arch/                # one architecture sketch per file
    └── system.md
```

**Why a folder, not a single file:** a vision can hold many UI screens, several arch views, multiple stories. Cramming them into one `VISION.md` reproduces the wall-of-ASCII problem and blocks per-fragment iteration/diff — the same pressure that pushed memory/sessions/ideas/feedback into soft-folders. The spine never balloons: it holds narrative + a manifest, each fragment is its own file.

**Fragment kind = subfolder location** (resolved Q2). The mutator `spectacular vision add <kind> <name>` maps `kind → folder`. No `kind:` frontmatter field is needed to classify — location is the classifier.

## The spine (`VISION.md`)

Narrative slots + a regenerable manifest:

| Spine slot | Content |
|---|---|
| **End goal** | One paragraph — what the world looks like when this exists. |
| **Macro dev phases** | The big dev-phase *arcs* (NOT milestones): "make it work → observable → self-serve". |
| **Flow walk** | Step-by-step narrative of the user moving through it — the imagined session. |
| **Manifest** | Linked list of every fragment in `stories/` `ui/` `arch/`, one-line caption + `approved:` state each. Regenerated from the fragment files (index mode). |

## Fragment kinds

| Subfolder | Kind | Shape |
|---|---|---|
| `stories/` | user story | `As a <persona>, I want … so that …` + acceptance + implied flow. Pulls personas from `PERSONAS.md` when present. |
| `ui/` | UI / output mockup | ASCII screen / CLI output / artifact the user will see and touch. One file per screen. |
| `arch/` | architecture sketch | ASCII box / structure / data-flow diagram. One file per view. |

**Fragment frontmatter (required shape):**

```yaml
---
kind: story | ui | arch
caption: <one line — feeds the spine manifest>
approved: pending | true | false      # per-fragment human reaction (resolved Q4)
personas: []                          # stories only — slugs into PERSONAS.md
related: []
updated: YYYY-MM-DD
---
```

`approved:` is the **per-fragment reaction state** — the human approves/redirects/rejects one fragment at a time. The derivation step (below) reads only `approved: true` fragments as load-bearing. State lives in the fragment, not in `feedback/` (v1 keeps it self-contained; cross-referencing two locations was the rejected alternative).

## The `imagine` mode

A **distinct mode** (resolved Q1), registered in the doc-index mode taxonomy. It is *not* a `grill` variant: grill **interrogates** the human slot-by-slot; `imagine` **generates** artifacts first, then derives a spec. Two behavioral differences grill-loop doesn't have.

**The loop:**

1. **Render** — the skill imagines and renders the spine (end-goal, macro phases, flow walk) + ≥1 ASCII fragment per kind. Leads with proposed artifacts, never empty-slot prompts.
2. **React on parts** — the human approves / redirects / rejects *individual* fragments. Approval written to fragment frontmatter; only redirected fragments are regenerated (not the whole vision).
3. **Derive (Build)** — when the human is satisfied, `imagine` **auto-offers `→ plan`** (resolved Q3): it drafts a `PLAN.md` from the approved vision —
   - approved **stories** → PLAN Goal + per-milestone outcomes
   - the **flow walk** → the milestone arc
   - approved **ui/arch fragments** → acceptance surfaces in Validation
   - and **pre-fills PLAN's `## Understanding`** from the vision spine.
4. **Gate** — the derived PLAN is a **draft, never auto-accepted**. It flows into the existing PLAN grill/review for the human gate (Q7). `imagine` proposes; the human disposes.

**v1 derivation is Build-only.** Compare/reconcile (diffing an *existing* PRD/PLAN against a vision to surface gaps) is deferred to v2, as is the **project altitude** (`imagine` around PRD, output at `.spectacular/vision/`) — gated on the PRD-overlap question (Q5/Q6). v1 is request-level only.

## Verbs

- **`imagine`** (the mode) — render → react → derive. The headline verb. `spectacular imagine <slug>`.
- **`grill`** — fall back to slot-by-slot if the human wants to author the spine manually instead of generatively. Rare; `imagine` is the default entry.
- **`refine`** — rewrite a vague spine slot or fragment caption into a specific one.
- **`review`** — validate fragment frontmatter shape + spine-manifest-vs-files drift + approval sanity (see doctor).

**Mutator verbs (CLI, not skill — mechanical/agentic split):**
- `spectacular imagine <slug>` — scaffold `vision/` (spine + empty subfolders) and enter the mode. Mechanical scaffold; the skill does the generative render.
- `spectacular vision add <kind> <name>` — write one fragment file (`kind → folder`), `approved: pending`, update the spine manifest.

## Doctor area

`spectacular doctor vision` — flags:

| Check | Severity | Condition |
|---|---|---|
| Manifest drift | warning | spine `VISION.md` manifest lists a fragment that no longer exists, or a fragment file absent from the manifest |
| Missing fragment frontmatter | warning | `kind`, `caption`, or `approved` absent/empty on a fragment |
| Unknown kind | warning | `kind:` not one of `story\|ui\|arch`, or fragment in a subfolder mismatching its `kind:` |
| Dangling persona ref | warning | a story's `personas:` slug not found in `PERSONAS.md` |
| Approval not advanced | info | request is `active` but all fragments still `approved: pending` (vision never reacted-to) |

Mechanical fixes (`--fix`): regenerate the spine manifest from fragment files. Approval state and content are judgment — never auto-set.

## Lifecycle fit

```
idea / brief ──► [ imagine → vision/ + draft PLAN ] ──► PLAN (grill/review) ──► active ──► …
```

The `vision/` folder becomes **read-only context** once PLAN exists — it explains the *why behind the shape*, the way PRD explains *why behind the project*. It never owns lifecycle state (PLAN keeps that). `snapshot-on-edit: false` — fragments are scratch-that-graduates, consistent with the other soft-folders.

## What this is **not**

- **Not a PRD.** PRD is project-wide convergent intent. A vision is per-request, divergent, and feeds a single PLAN. (The project-altitude vision that *would* sit near PRD is v2 — and its PRD overlap is unresolved, Q5.)
- **Not a PLAN.** PLAN is the convergent decomposition. A vision is the divergent imagination that PLAN is *derived from*. They coexist; PLAN holds lifecycle state.
- **Not feedback.** `feedback/` answers "was that the right thing to ship?" post-build. A vision's `approved:` is pre-build reaction to imagined artifacts.
- **Not auto-routed.** Like ideas, nothing in `vision/` is acted on automatically. The human triggers `imagine`; the derived PLAN is always a draft to gate.

**Related:** [[idea-rules]] (the upstream scratchpad), [[plan-rules]] (the derived doc), [[new-request]] (handoff target), [[personas-rules]] (story personas), [[feedback-rules]] (post-ship vs pre-build), [[doc-index]] (mode taxonomy), [[scaffold-reference]], [[grill]] / [[refine]] / [[review]] (generic engine).
