---
status: planned
priority: high
owner: alex
updated: 2026-05-31
summary: "/spectacular imagine — imagination-backed planning. Skill renders artifacts (UI/flow/stories/arch) the human reacts to per-part, then derives the spec/plan to deliver the approved flow. Expands the thesis from spec-driven to spec-driven AND imagination-backed."
related:
  - PRD.md
  - ../../ROADMAP.md
  - ../../ARCHITECTURE.md
  - ../../PERSONAS.md
target_version: v1.15.0
---

# Plan — imagine-mode

> **Promoted from idea** `ideas/explore-mode.md` (now at `archive/ideas/explore-mode.md`). The archived idea holds the full thinking (thesis §0, derivation loop §5, all 9 open questions). This PLAN carries the **narrow v1 slice**; the idea is the v2+ backlog and the rationale.

> ⚠ **OPEN QUESTIONS — resolve before `planned → active`.** This request is scaffolded with the scope intentionally cut small, but several design questions are unresolved. They are listed in § Open questions below and must be settled (some at M1) before implementation. Do **not** silently pick defaults for the ones marked **[blocking]**.

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

<!-- REQUIRED before planned → active by the understand-before-change policy. Filled at M1 once open questions resolve. Left as a stub deliberately — see Open questions. -->

### How it works now

Spectacular is spec-driven: `idea → PRD → SPEC → PLAN → build`. The human reacts to *specification text*. There is no phase that renders see-able artifacts, and no engine that derives specs *from* an approved vision. `grill`/`refine`/`review` interrogate; nothing generates-then-reconciles.

### What changes

*(fill at M1 — depends on Open Q1 mode decision and Q4 handoff decision.)*

### What stays the same

PLAN keeps sole ownership of lifecycle state. PRD/SPEC remain the canonical convergent docs. The soft-folder substrate, frontmatter-as-signal, and snapshot conventions are unchanged. `imagine` is **optional** — small/obvious requests skip it.

## 3. Milestones

- **M1 — Resolve open questions + write the contract.** Settle the blocking open questions (mode, subfolder-vs-flat, handoff); document the `vision` doc-type in a `vision-rules.md` ref + doc-index entry + ARCHITECTURE. Fill the Understanding slot. *(Gate: no code until this lands.)*
- **M2 — `vision/` soft-folder substrate.** `spectacular imagine <slug>` scaffolds `requests/<slug>/vision/` (spine `VISION.md` + `stories/` + `ui/` + `arch/`). CLI mutators add fragments (`vision add <kind> <name>`). Index/manifest regenerable. `doctor vision` area.
- **M3 — Generative render engine.** The skill imagines + renders the spine (end-goal, macro dev phases, flow walk) + ≥1 fragment of each kind (story, ui, arch) in ASCII. Leads with proposed artifacts.
- **M4 — React-on-parts loop.** Per-fragment approval (`approved: true|false|pending` frontmatter); the human approves/redirects/rejects individual fragments; engine regenerates only what's redirected.
- **M5 — Build derivation.** From the approved vision, `imagine` drafts/refines `PLAN.md` (stories→goals, flow→milestone arc, fragments→acceptance surfaces), then hands off to PLAN grill/review. Pre-fills `## Understanding` from the vision.
- **M6 — Dogfood + ship.** Run `imagine` on a real request in this repo; CHANGELOG entry; plugin bump to v1.15.0.

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

> **Signalled per request.** Settle the `[blocking]` ones at M1 before any code. Full discussion lives in `archive/ideas/explore-mode.md` §8.

1. **[blocking] Mode** — distinct `imagine` mode vs reuse `grill-loop`? Leaning **distinct** (generative-first + derivation are real behavioral differences from grill). Decide before M3.
2. **[blocking] Fragment layout** — typed *subfolders* (`stories/` `ui/` `arch/`) vs a flat `fragments/` folder typed by `kind:` frontmatter. User leaned subfolders; confirm before M2 (drives the mutator signature).
3. **[blocking] vision → PLAN handoff** — does `imagine` auto-offer `→ plan` at the end? Does Build derivation pre-fill PLAN `## Understanding`? Decide before M5 (currently assumed yes to both).
4. **Approval substrate** — does per-fragment `approved:` state live in fragment frontmatter, or as `feedback/` entries (reusing [[feedback-loop]])? Affects M4.
5. **PRD overlap (v2 gate)** — at the project altitude, is `vision` a *pre-PRD* doc or a *feedback layer on a PRD draft*? PRD already has a Vision slot. **Must resolve before the v2 project altitude** — not blocking v1 (which is request-only).
6. **PRD positioning copy (v2)** — the thesis shift means `.spectacular/PRD.md`'s Vision section should eventually claim "spec-driven **and** imagination-backed." Does v1 touch PRD positioning, or is it a follow-up? Leaning follow-up.
7. **Derivation trust** — where is the gate on the *derived* PLAN? Assumed: the existing PLAN grill/review (draft is never auto-accepted). Confirm at M5.
8. **ASCII palette** — ship a `templates/vision/` palette (box-drawing chars, screen-frame convention) for render consistency, or let the engine improvise? Leaning ship-a-palette. Decide at M1/M3.
