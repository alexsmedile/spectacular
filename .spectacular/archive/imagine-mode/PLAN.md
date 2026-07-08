---
status: archived
priority: high
owner: alex
updated: 2026-06-28
summary: "/spectacular imagine — imagination-backed planning. Skill renders artifacts (UI/flow/stories/arch) the human reacts to per-part, then derives the spec/plan to deliver the approved flow. Expands the thesis from spec-driven to spec-driven AND imagination-backed."
related:
  - ../../PRD.md
  - ../../roadmaps/index.md
  - ../../ARCHITECTURE.md
  - ../../PERSONAS.md
build: b6
archived: 2026-06-28
---

# Plan — imagine-mode

> **Promoted from idea** `ideas/explore-mode.md` (now at `archive/ideas/explore-mode.md`). The archived idea holds the full thinking (thesis §0, derivation loop §5, all 9 open questions). This PLAN carries the **narrow v1 slice**; the idea is the v2+ backlog and the rationale.

> ✅ **Blocking questions resolved 2026-06-02** (Q1 mode, Q2 layout, Q3 handoff, Q4 approval — see § Open questions). The `## Understanding` slot is filled; the `understand-before-change` gate is satisfied. Remaining open items (Q5/Q6 PRD overlap, Q8 ASCII palette) are v2/M3 and do **not** block `active`.

> 🔍 **Shipped as PREVIEW (2026-06-02, untagged).** The feature is functionally complete (M1–M6 done, dogfooded) and committed/pushed to `main`, but **not released** — no version bump, no tag. It rides to a real release in **v1.15.0** alongside [[visual-layer]]. Until then it's usable on `main` as a preview. Focus moved to v1.13.0 [[cross-request-links]].

## 1. Goal

Add `/spectacular imagine <slug>` — a generative, imagination-backed planning mode that renders artifacts a vision-driven human can see and react to (UI fragments, flow walk, user stories, architecture sketches), then **derives a draft PLAN from the approved vision**. This expands Spectacular's thesis from *spec-driven* to *spec-driven **and** imagination-backed*: specs become accountable to what the human approved, not authored in a vacuum.

## 2. Constraints

- **v1 is request-level only.** `imagine <slug>` operates on a single request's `vision/` folder. The project altitude (around-PRD) is deferred to v2 — it collides with PRD's existing Vision slot (see Open Q5/Q6).
- **v1 derivation = Build only.** At the end, `imagine` drafts a PLAN *from* the approved vision. The Compare/reconcile mode (diff an existing spec against the vision) is the hard half — deferred to v2.
- **`vision` is one `index`-mode soft-folder doc-type** — not multiple doc-types. Stories/UI/arch are fragment *kinds* inside it. Reuses the proven soft-folder plumbing (memory/sessions/ideas/feedback).
- **The spine never balloons.** `VISION.md` holds narrative + a regenerable manifest of fragment links; each fragment is its own file. This is the answer to "many UI fragments."
- **Generative-first, human reacts on parts.** The engine leads with rendered artifacts, not empty-slot prompts. Approval is per-fragment.
- **Derived specs are drafts, never auto-accepted.** Build-mode output is handed to the existing PLAN grill/review for the human gate. `imagine` proposes; the human disposes.
- **No new language deps.** Bash CLI + skill refs, consistent with the rest of Spectacular.
- **Mechanical/agentic split (v1.4.0).** CLI scaffolds the folder + adds fragment files (mechanical); the skill imagines/renders/derives (agentic).

## Understanding

### How it works now

Spectacular is spec-driven: `idea → PRD → SPEC → PLAN → build`. The human reacts to *specification text*. There is no phase that renders see-able artifacts, and no engine that derives specs *from* an approved vision. The existing doc verbs (`grill`/`refine`/`review`) **interrogate** the human slot-by-slot; none **generate** an artifact first or reconcile a spec against it. Per-request docs today are `PLAN.md` + `TASKS.md`; soft-folder DBs (memory/sessions/ideas/feedback) already exist with `index` mode, CLI mutators, and doctor areas — the proven substrate this request reuses.

### What changes

A new **distinct `imagine` mode** (Q1) and a `vision` `index`-mode soft-folder doc-type:

- `spectacular imagine <slug>` scaffolds `requests/<slug>/vision/` with **typed subfolders** `stories/` `ui/` `arch/` (Q2) + a spine `VISION.md` (narrative + regenerable manifest).
- The skill **generates-first**: renders the spine (end-goal, macro dev phases, flow walk) + ≥1 ASCII fragment per kind, then the human **reacts per-fragment**. Approval state lives in each fragment's frontmatter as `approved: true|false|pending` (Q4).
- `spectacular vision add <kind> <name>` mutator (kind→folder) lets a vision grow to many fragments; `doctor vision` validates frontmatter + manifest-vs-files drift + dangling persona refs + approval sanity.
- At the end, `imagine` **auto-offers `→ plan`**: derives a draft `PLAN.md` from the approved vision (stories→goals, flow→milestones, fragments→acceptance) and **pre-fills `## Understanding`** from the vision spine (Q3). The draft is never auto-accepted — it flows into the existing PLAN grill/review gate (Q7).

### What stays the same

PLAN keeps sole ownership of lifecycle state. PRD/SPEC remain the canonical convergent docs — untouched in v1 (the PRD-positioning rewrite is v2, Q6). The soft-folder substrate, frontmatter-as-signal, snapshot conventions, and the mechanical(CLI)/agentic(skill) split are unchanged. **v1 derivation is Build-only** — Compare/reconcile (diffing an *existing* spec against a vision) and the **project altitude** are deferred to v2. `imagine` is **optional** — small/obvious requests skip straight to PLAN.

## 3. Milestones

- **M1 — Write the contract.** (Blocking questions resolved 2026-06-02; Understanding filled.) Document the `vision` doc-type: `vision-rules.md` ref (frontmatter, slots, fragment kinds, `imagine` mode behavior) + doc-index entry + ARCHITECTURE section + register `imagine` in the mode taxonomy.
- **M2 — `vision/` soft-folder substrate.** `spectacular imagine <slug>` scaffolds `requests/<slug>/vision/` (spine `VISION.md` + `stories/` + `ui/` + `arch/`). CLI mutators add fragments (`vision add <kind> <name>`). Index/manifest regenerable. `doctor vision` area.
- **M3 — Generative render engine.** The skill imagines + renders the spine (end-goal, macro dev phases, flow walk) + ≥1 fragment of each kind (story, ui, arch) in ASCII. Leads with proposed artifacts.
- **M4 — React-on-parts loop.** Per-fragment approval (`approved: true|false|pending` frontmatter); the human approves/redirects/rejects individual fragments; engine regenerates only what's redirected.
- **M5 — Build derivation.** From the approved vision, `imagine` drafts/refines `PLAN.md` (stories→goals, flow→milestone arc, fragments→acceptance surfaces), then hands off to PLAN grill/review. Pre-fills `## Understanding` from the vision.
- **M6 — Dogfood + ship.** Run `imagine` on a real request in this repo; CHANGELOG entry; plugin bump to v1.15.0. **Co-ships with [[visual-layer]] in v1.15.0** — both are the ASCII-rendering milestone (imagine renders UI/arch fragments; visual-layer renders progress/summary/roadmap). They share the `ascii-render` substrate.

## 4. Tasks

See `TASKS.md`.

## 5. Dependencies

- Reuses the **soft-folder DB substrate** ([[soft-db-substrate]], v1.5.0) — `index`-mode plumbing, CLI mutators, doctor areas.
- Ties into the **feedback substrate** ([[feedback-loop]], v1.6.0) — per-fragment approval is structured human feedback; `imagine` may produce feedback entries.
- Consumes **PERSONAS.md** ([[add-personas-doc]], v1.3.0) when present — stories reference personas.
- Independent of [[cross-request-links]] (v1.13.0) and [[cli-debt-removal]] (v1.14.0); can ship after them.
- **Scope-down policy** ([[../../POLICY.md]] @Planning) directly motivates the narrow v1 cut.

## 6. Validation

- **M1** — `vision-rules.md` + doc-index entry + ARCHITECTURE section exist; all `[blocking]` open questions marked resolved with a recorded decision; Understanding slot filled.
- **M2** — `spectacular imagine <slug>` creates the `vision/` tree; `vision add ui dashboard` writes `vision/ui/dashboard.md` and the spine manifest updates; `doctor vision` passes clean and flags a deliberately-broken manifest link.
- **M3** — On a fixture request, the engine renders a spine + one story + one ASCII UI fragment + one ASCII arch sketch without manual authoring.
- **M4** — Approving one fragment and redirecting another regenerates only the redirected one; approval state persists in frontmatter.
- **M5** — A vision with 2 approved stories + 1 UI fragment yields a draft PLAN whose Goal/Milestones/Validation visibly trace to those artifacts; PLAN `## Understanding` is pre-filled.
- **M6** — A real in-repo request goes idea/brief → `imagine` → derived PLAN; test suite green; manifests at v1.15.0.

## 7. Deliverables

- `spectacular imagine <slug>` CLI verb + `vision add <kind> <name>` mutator
- `vision` doc-type: `vision-rules.md`, doc-index entry, `templates/vision/` (spine + fragment scaffolds + ASCII palette)
- `doctor vision` area (fragment frontmatter, manifest-vs-files drift, dangling persona refs, approval-state sanity)
- Generative render engine (skill ref — the `imagine` mode behavior)
- Build-derivation engine: approved vision → draft PLAN + pre-filled Understanding
- ARCHITECTURE.md section for the `vision/` substrate + the imagination-backed thesis
- CHANGELOG [1.15.0] entry; plugin bump to v1.15.0
- Dogfood artifact: one real request planned via `imagine`

## Open questions

> **Signalled per request.** The three `[blocking]` ones + Q4 were **resolved 2026-06-02** (see below). Full original discussion lives in `archive/ideas/explore-mode.md` §8.

1. ~~**[blocking] Mode**~~ **RESOLVED 2026-06-02 → distinct `imagine` mode.** grill interrogates slot-by-slot; `imagine` leads with rendered artifacts (generative-first) then derives a spec — two behavioral differences grill-loop doesn't have. New mode in the taxonomy.
2. ~~**[blocking] Fragment layout**~~ **RESOLVED 2026-06-02 → typed subfolders** `vision/stories/` `vision/ui/` `vision/arch/`. Kind = location. Mutator `vision add <kind> <name>` maps kind→folder; the spine manifest groups by folder.
3. ~~**[blocking] vision → PLAN handoff**~~ **RESOLVED 2026-06-02 → auto-offer `→ plan` AND pre-fill `## Understanding`** from the vision spine. Closes the imagination-backed loop. Derived PLAN is a draft — still goes through PLAN grill/review (never auto-accepted).
4. ~~**Approval substrate**~~ **RESOLVED 2026-06-02 → fragment frontmatter** `approved: true|false|pending`. Self-contained, travels with the fragment, simple for derivation to read; doctor sanity-checks it. (Not reusing `feedback/` entries in v1 — avoids cross-referencing two locations.)
5. **PRD overlap (v2 gate)** — at the project altitude, is `vision` a *pre-PRD* doc or a *feedback layer on a PRD draft*? PRD already has a Vision slot. **Must resolve before the v2 project altitude** — not blocking v1 (which is request-only). *Still open — v2.*
6. **PRD positioning copy (v2)** — the thesis shift means `.spectacular/PRD.md`'s Vision section should eventually claim "spec-driven **and** imagination-backed." Follow-up, not v1. *Still open — v2.*
7. **Derivation trust** — gate on the *derived* PLAN = the existing PLAN grill/review (draft never auto-accepted). **Settled by Q3** — confirm in practice at M5.
8. ~~**ASCII palette**~~ **RESOLVED 2026-06-02 (M3) → ship a light convention** in [[imagine]] (box-drawing frames, `───►` flow, `[ x ]`/`‹ x ›` markers, ≤64-char width). Readability over strict grammar. Templates already carry box-drawing scaffolds.
