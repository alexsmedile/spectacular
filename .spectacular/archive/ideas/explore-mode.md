---
type: idea
status: promoted
priority: high
owner: alex
updated: 2026-05-31
origin: chat 2026-05-31 — "/spectacular imagine": imagination-backed planning (expands the framework's core thesis)
related:
  - ../../PRD.md
  - ../../roadmaps/index.md
  - ../../ARCHITECTURE.md
  - ../../PERSONAS.md
  - ../../skills/spectacular/references/doc-index.md
  - ../../skills/spectacular/references/plan-rules.md
  - ../../skills/spectacular/references/idea-rules.md
promoted_to: requests/imagine-mode/
---

# Idea — `/spectacular imagine` (imagination-backed planning)

> **Status note (2026-05-31):** captured from chat. This is **not just a new doc-type** — it expands Spectacular's core thesis from *spec-driven development* to **spec-driven and imagination-backed**. The skill renders artifacts a vision-driven human can see and react to, then **derives or reconciles** the PRD/SPEC/PLAN needed to deliver the approved flow. Macro concept confirmed; this doc is raw material for a **narrow, spec-shaped request** (see §6 — per the `scope-down` policy). Treat as thinking, not a plan.

---

## 0. The thesis shift (the headline)

Spectacular today is **spec-driven**: the human reacts to *specifications*. You converge on what to build (PRD → SPEC → PLAN) and imagine the product *from* the spec text.

`/spectacular imagine` makes it **spec-driven AND imagination-backed**: the skill renders **artifacts a vision-driven human can see** (UI fragments, flows, stories, arch sketches), the human gives feedback on *parts*, and the specs are then **back-solved from what the human approved**. Imagination isn't a side sketch — it's an input the specs become accountable to.

```
spec-driven (today):       brief → PRD → SPEC → PLAN → build         (human reacts to spec TEXT)
imagination-backed (new):  brief → IMAGINE (artifacts) → human reacts on PARTS
                                     ↓ derive / reconcile
                                  PRD / SPEC / PLAN accountable to the approved flow
```

The inversion is the point: **artifact-first, spec-derived.** Today specs come first and you imagine from them; with `imagine`, you imagine first and the specs are derived from the approved vision.

## 1. Where it sits — one capability, three entry points

`imagine` is a **mode of working**, not a fixed lifecycle slot (like how `grill`/`refine`/`review` are verbs that apply across docs). Same engine, three insertion points:

1. **After goals + initial brief, before a full PRD** — imagine the product, then write the PRD that delivers it.
2. **As a feedback layer on a PRD draft** — render what the draft *implies*, let the human react to the rendered flow, feed corrections back into the PRD.
3. **Before a request's PLAN** — imagine the slice, then derive/refine the plan (the original "advanced plan mode" framing).

Still fills the divergent gap between unstructured `idea` and convergent `PLAN` — but now also operates *above* the PRD, not just below it.

## 2. Two altitudes (confirmed in chat)

`explore` is **not** request-only. It runs at two altitudes with the **same engine**:

| Altitude | When | Feeds | Output location |
|---|---|---|---|
| **Project** | before/around PRD — "what is this whole product?" | `PRD.md` | `.spectacular/VISION.md` + `.spectacular/stories/` |
| **Request** | before a single request's PLAN — "what is this slice?" | `requests/<slug>/PLAN.md` | `requests/<slug>/VISION.md` + `requests/<slug>/stories/` |

Lifecycle gains an optional divergent stage at both altitudes:

```
project:   [explore → VISION.md + stories/] → PRD → SPEC …
request:   idea → [explore → VISION.md + stories/] → PLAN → active …
```

VISION is **read-only context** once PRD/PLAN exists — it explains the *why behind the shape*, the way PRD explains *why behind the project*. It never owns lifecycle state (PLAN keeps that). It is **optional** — small/obvious requests skip straight to PLAN.

## 3. The structure — `vision` is a soft-folder, not a single file

**Decided 2026-05-31 (full soft-folder).** A vision can hold *many* fragments — multiple UI screens (login, dashboard, empty state, error state, CLI output), multiple arch sketches (system view, data flow, component zoom), multiple stories. Cramming them into one `VISION.md` reproduces the wall-of-ASCII problem and blocks per-fragment iteration/diff. This is the exact pressure that pushed memory/sessions/ideas/feedback out of single files into soft-folder DBs — so `vision` is a soft-folder from day one (no later migration).

```
requests/<slug>/vision/         # (and .spectacular/vision/ at project altitude in v2)
├── VISION.md                   # the SPINE — index/manifest, not a dump
├── stories/
│   ├── 01-first-run.md
│   └── 02-power-user.md
├── ui/                         # many UI fragments, ONE FILE EACH
│   ├── dashboard.md
│   ├── empty-state.md
│   └── cli-output.md
└── arch/
    ├── system.md
    └── data-flow.md
```

### 3a. `vision` (the spine — `index`-mode)

`VISION.md` is the **narrative spine + manifest**, not a container for all the ASCII:

| Spine slot | Produces | Why exploratory |
|---|---|---|
| **End goal** | the macro outcome in one paragraph — what the world looks like when this exists | anchors imagination |
| **Macro dev phases** | the big dev-phase *arcs* (NOT milestones): "make it work → observable → self-serve" | phase thinking, not task thinking |
| **Flow walk** | step-by-step narrative of the user moving through it — the imagined session | "imagining the flow with the user" |
| **Manifest** | linked list of every fragment in `stories/` `ui/` `arch/` with a one-line caption each | keeps the spine scannable at any fragment count |

- **Mode:** `index` for the collection (index regenerated from fragment files; CLI mutators write fragments; doctor checks index-vs-files drift) **+** a generative `explore` engine that renders the fragments and the spine. The `index` plumbing is proven (memory/sessions/ideas/feedback); the generative-first render loop is the new part.
- **Stays scannable** no matter how many fragments — VISION.md never balloons.

### 3b. fragment types (each file is one fragment)

| Subfolder | Fragment | Shape |
|---|---|---|
| `stories/` | user story | `As a <persona>, I want … so that …` + acceptance + implied flow. Pulls personas from `PERSONAS.md`. |
| `ui/` | UI / output mockup | ASCII screen / CLI output / artifact the user will see and touch. **One file per screen** — this is the answer to "many UI fragments." |
| `arch/` | architecture sketch | ASCII box/structure/data-flow diagram. One file per view. |

- All fragment files share a common frontmatter (`kind:`, `caption:`, `personas:`/`related:`) so the spine's manifest is regenerable and doctor can validate links.
- Lives at both altitudes: `.spectacular/vision/` (product, v2) and `requests/<slug>/vision/` (slice, v1).
- **Open (was Q2, now narrowed):** are `stories` / `ui` / `arch` *subfolders* (typed by location) or a flat `fragments/` folder typed by a `kind:` frontmatter field? Subfolders = clearer taxonomy; flat = one mutator + filter-by-kind. Resolve at spec time. *(User leaned full soft-folder with subfolders.)*

## 4. ASCII artifact taxonomy (all four confirmed wanted)

1. **Architecture diagrams** — box/structure, components, data flow.
2. **UI / output mockups** — screens, CLI output, generated artifacts the user will see.
3. **User stories + flow walk** — as-a/I-want + narrative step-through of the imagined session.
4. **Macro dev phases** — phase arcs, not milestones.

Open: do we template ASCII palettes (box-drawing chars, a "screen frame" convention) so renders are consistent? Probably yes — a small `templates/vision/` with ASCII scaffolds.

## 5. The derivation loop — the NEW hard part (artifacts → specs)

This is the heart of the thesis shift and the genuinely new capability. The earlier framing only *produced* artifacts; `imagine` must **close the loop**: approved artifacts → derive/reconcile PRD/SPEC/PLAN so the spec is accountable to the imagined flow.

Two derivation modes:

| Mode | Precondition | Behavior |
|---|---|---|
| **Build** | no PRD/PLAN yet | generate a *draft* PRD/PLAN **from** the approved vision — stories become goals/requirements, the flow walk becomes the milestone arc, UI fragments become acceptance surfaces. |
| **Compare / reconcile** | PRD/PLAN already exists | **diff** the spec against the approved vision; surface gaps where the spec doesn't yet deliver the imagined flow. e.g. "the `dashboard` fragment you approved has no requirement in PRD §3"; "the `empty-state` flow isn't covered by any milestone in PLAN." |

The fragment soft-folder (§3) is just the **substrate** `imagine` speaks through; this derivation/reconciliation engine is the new, non-trivial part. It is where "imagination backed in" becomes real — the specs are *checked against* what the human reacted to and approved, not authored in a vacuum.

**Feedback granularity:** the human reacts on *parts* (approve this UI fragment, redirect that flow step, reject this story) — not the whole vision at once. So the loop is per-fragment, and approval state likely lives in fragment frontmatter (`approved: true|false|pending`) so derivation knows which artifacts are load-bearing. This dovetails with the existing `feedback/` substrate (v1.6.0) — `imagine` may *be* a structured producer of feedback entries.

## 6. Interaction character (confirmed: generative-first)

The agent **imagines and renders** stories / flow / diagrams / ASCII UI up front; the user reacts and redirects. Contrast with grill (agent interrogates slot-by-slot). This "imagination-backed" feel is the headline — the engine must lead with a proposed artifact, not an empty prompt.

Implication for the engine: `imagine` can't just be a thin `grill` alias. It needs a **render → react-on-parts → regenerate/deepen → derive** loop. Closest existing primitive is `grill-loop` (wide then narrow) but inverted toward generation, and extended with the derivation step (§5).

## 7. Narrow first slice (apply scope-down before building)

Per the `scope-down` policy: name the smallest high-impact slice, push the rest to ROADMAP `v2+`. The thesis is big; the **first shippable slice must be small** or it won't close.

**Candidate v1 (smallest that proves "imagination-backed"):**
- **One entry point: request-level.** `spectacular imagine <slug>` produces the `requests/<slug>/vision/` soft-folder (spine + `stories/` + `ui/` + `arch/`).
- `vision` registered as one `index`-mode soft-folder doc-type (doc-index, rules file, templates, doctor area). Fragments typed by subfolder/`kind:`.
- Generative-first engine: render spine (end-goal + macro-phases + flow-walk) + ≥1 fragment of each kind, then react-on-parts per fragment (`approved:` frontmatter).
- **Derivation = Build only.** At the end, `imagine` drafts/refines `PLAN.md` *from* the approved vision. **Defer Compare/reconcile to v2** — the diff engine is the hard half; prove Build first.
- CLI mutators to add fragments so a vision grows without hand-editing the manifest.

**Deferred to v2+ (roadmaps/index), in rough order:**
- **Compare/reconcile derivation** (§5) — diff an existing PRD/PLAN against the vision. The valuable-but-hard half.
- **Project altitude** — `imagine` before/around PRD (§1 entry points 1 & 2), output at `.spectacular/vision/`. Needs the PRD-overlap question (§8 Q5) resolved first.
- Rich diagram *types* beyond box/flow (sequence, state).
- ASCII palette/templating system if hand-rolled renders prove inconsistent.
- Auto-promotion of fragments → tasks (stays human judgment, mirrors idea non-routing rule).

## 8. Open questions

1. **Mode:** distinct `imagine` mode vs `grill-loop` reuse? (Leaning distinct — generative-first + the derivation step are real behavioral differences.)
2. ~~**Doc-type count:** is `story` worth a full index-doc-type, or stories-as-section?~~ **RESOLVED 2026-05-31:** the "many UI fragments" pressure settled it — `vision` is a full `index`-mode soft-folder (spine + stories/ui/arch fragments); stories/UI/arch are fragment *kinds* inside it, not separate doc-types. See §3.
3. **Altitude routing:** if v1 is request-only, does the CLI verb signature already reserve a `--project` flag for the v2 project-level altitude?
4. **vision → PLAN handoff:** does `imagine` auto-offer `→ plan` at the end? Does the Build derivation pre-fill PLAN's `## Understanding` slot from the approved vision?
5. **PRD relationship at project altitude:** is project-`vision` a *pre-PRD* doc (feeds PRD) or a *feedback layer* on a PRD draft (§1 entry points 1 vs 2)? PRD already has a Vision slot — real overlap risk. Must resolve before building the project altitude. (Why v1 is request-only.)
6. **PRD positioning copy:** the thesis shift (§0) means `.spectacular/PRD.md`'s own Vision section (currently coherence/context-management only, no imagination concept) should eventually be updated to claim "spec-driven **and** imagination-backed." Does the v1 request touch the PRD positioning, or is that a follow-up?
7. **Derivation trust:** how much does the human review the *derived* PRD/PLAN vs trust it? Build-mode output is a draft to grill, never auto-accepted — but where's the gate? (Likely the existing PLAN grill/review.)
8. **Snapshot policy:** `snapshot-on-edit: false` — consistent with other `index`-mode soft-folders. (Settled by the soft-folder decision.)
9. **doctor area:** what does `doctor vision` check? (fragment frontmatter, dangling persona refs, spine-manifest-vs-files drift, approval-state sanity.)

## 9. Why park it (not cut a request yet)

This is a **thesis-level expansion** (spec-driven → spec-driven + imagination-backed), not a bolt-on. The full concept = new generative mode + soft-folder substrate + the derivation/reconciliation engine + two altitudes + a PRD-positioning change. Far more than one clean milestone. Parking lets the concept settle and the narrow slice (§7 — request-level, Build-only derivation) get pressure-tested before it becomes a `requests/<slug>/PLAN.md`. When ready: `spectacular idea promote explore-mode` cuts the request from the v1 slice.
