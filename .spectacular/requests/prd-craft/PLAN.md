---
status: verified
priority: high
owner: alex
updated: 2026-05-21
summary: "Effortless, interactive PRD crafting for any kind of project — grill, refine, review"
related:
  - ../../../skills/spectacular/SKILL.md
  - ../doc-writer/PLAN.md
---

# Plan — PRD Craft

## Goal

Ship an interactive PRD-crafting flow inside the Spectacular skill that takes a user from a one-sentence project idea to a usable `.spectacular/PRD.md` in ~10 minutes, without sub-agent pipelines or long upfront interviews. General-purpose: works for coding, content, research, or hardware projects.

## Why

Writing a good PRD from scratch is high-friction. Existing spec-driven frameworks (GSD, BMAD) over-engineer with multi-agent research pipelines. Most projects need a fast, low-ceremony way to go from vibe to a usable PRD — one that's actually referenced downstream instead of abandoned after week one.

Karpathy-aligned: minimum viable PRD first; expand into phases/specs only when complexity demands it.

## Scope

**In scope (v1)**
- 6-slot canonical PRD template (general-purpose, domain-agnostic)
- 5 starter kits: coding, product, content, research, blank
- Three modes: grill (interactive), refine (vibe→spec), review (quality gate)
- Reference docs loaded on demand by SKILL.md
- Routing entries in SKILL.md for `prd`, `prd refine`, `prd review`
- Dogfood: this very PLAN passes the review gate

**Out of scope (v2)**
- Phase splitting (PRD → multiple specs/phases)
- CLI command `spectacular prd` (Bash binary)
- Auto-detect kit from project context
- Quality-gate scoring beyond pass/fail

**Explicit anti-patterns (not v2, never)**
- Per-request PRDs (`requests/<slug>/PRD.md`). PRD is **project-wide** and lives at `.spectacular/PRD.md`. Per-request folders contain `PLAN.md` + `TASKS.md` only. The grill writes to the root PRD, never inside a request folder.

## Approach

Markdown-only build. No CLI, no scripts. Skill reads references on demand and runs the flow itself.

### Architecture

```
skills/spectacular/
├── SKILL.md                          # add prd / prd refine / prd review triggers
├── references/
│   ├── prd-grill.md                  # interactive slot-filling loop
│   ├── prd-refine.md                 # vibe → spec rewrite patterns
│   └── prd-review.md                 # quality gate checklist
└── templates/
    └── prd/
        ├── base.md                   # 6-slot canonical template
        └── kits/
            ├── coding.md             # + STACK section
            ├── product.md            # + user stories, metrics
            ├── content.md            # + audience, distribution
            ├── research.md           # + hypothesis, method
            └── blank.md              # pure base, no extras
```

### The 6 required slots

1. **Problem** — one sentence
2. **Who it's for** — one primary user
3. **What success looks like** — measurable, time-boxed
4. **Non-goals** — explicit exclusions
5. **Constraints** — budget, time, tech, policy
6. **First milestone** — concrete near-term outcome

Optional, prompted only when relevant: stakeholders, risks, open questions, prior art.

### Three modes

| Mode | Trigger | What it does |
|---|---|---|
| Grill | `prd` (no PRD.md) or `prd grill` | Strict slot-order interview, one Q at a time |
| Refine | `prd refine` or detected vibe language | Vibe→spec rewrite proposals with `[NEEDS CLARIFICATION]` markers |
| Review | `prd review` | Punch list against quality gate — no rewrites, user decides |

### When `PRD.md` already exists

The grill never silently overwrites. On `prd` / `prd grill`:

- **Non-empty PRD found** → ask: refine in place (`prd refine`), or start fresh (snapshots current to `PRD@vN.md` first, then grills)?
- **Effectively empty PRD** (only `<PLACEHOLDER>` markers) → grill in place.
- **No PRD** → create from selected kit, then grill.

Snapshot-before-overwrite is automatic when the user chooses "start fresh" — the user does not need to run `spectacular snapshot PRD.md` manually.

### Kit selection

Five kits live in `templates/prd/kits/`. The grill picks one as its **first question**, before slot 1:

| Kit | When to pick |
|---|---|
| `coding` | CLI, library, app, service, SDK. Adds STACK + interfaces. |
| `product` | Consumer/B2B product with user flows. Adds user stories + metrics + distribution. |
| `content` | Course, newsletter, book, video series, docs. Adds audience + format + distribution. |
| `research` | Investigation feeding a downstream decision. Adds hypothesis + method + deliverable. |
| `blank` | None of the above, or when starting with maximum flexibility. Pure 6-slot base. |

**Override mechanism:** if `.spectacular/templates/prd/kits/<kit>.md` exists in the project, it wins over the skill's bundled template. Same filename, project-local takes precedence. (Flag-based override like `prd --kit=coding` is deferred to v2 alongside the CLI.)

### Grill flow (strict slot order)

0. Pre-flight: check for existing PRD (see "When `PRD.md` already exists" above); pick a kit
1. Ask Q for slot 1 (Problem)
2. Receive answer → write to PRD.md immediately
3. Run a **mini-refine** on the answer (inline, pattern-match only — see table in `prd-refine.md`)
4. User accepts/edits → next slot
5. Repeat for slots 2–6
6. After slot 6, run **full review** → show punch list
7. Loop on flagged items until gate passes; if the user can't resolve a flag, insert `[NEEDS CLARIFICATION: <gap>]` and accept it as the exit condition (gate fails-soft with the list still visible)

**Mini-refine vs full review.** Mini-refine = inline pattern check after each answer (vague words, plural users, no number in success). Full review = end-of-grill gate that checks all slots together, plus cross-slot consistency. Both live in `prd-refine.md` / `prd-review.md` — mini is a subset.

### Refine patterns (vibe → spec)

- Vague adjectives (`fast`, `intuitive`, `scalable`) → request measurable replacement
- Plural users (`users`, `customers`) → request the *one* primary user
- Unbounded success (`make it great`) → request 30/60/90-day measurable
- Tech jargon in problem statement → ask for plain-language version
- When can't resolve → insert `[NEEDS CLARIFICATION: <specific gap>]`

### Review gate

PRD passes when:

- All 6 required slots non-empty
- No `[NEEDS CLARIFICATION]` markers
- Success criteria contains at least one number AND one verb AND one date/timeframe
- No vague-word list hits in problem/success/milestone (list lives in prd-review.md)
- Non-goals list is non-empty and specific

## Success criteria

- A user with a 1-sentence project idea runs `/spectacular prd` and walks away with a `.spectacular/PRD.md` that passes the review gate in <15 minutes on a fresh project
- All 6 required slots filled with measurable, specific content
- Zero `[NEEDS CLARIFICATION]` markers remaining
- Dogfood: this project's own `.spectacular/PRD.md` passes `prd review` (PRDs are project-wide — there is no per-request PRD to dogfood against)

## Build steps

1. [x] Scaffold request (PLAN.md, TASKS.md)
2. [x] Write `templates/prd/base.md`
3. [x] Write 5 starter kits
4. [x] Write `references/prd-grill.md`
5. [x] Write `references/prd-refine.md`
6. [x] Write `references/prd-review.md`
7. [x] Add routing entries to `SKILL.md`
8. [x] Verify `references/prd-review.md` includes the vague-word list referenced by grill
9. [ ] Update root `.spectacular/PRD.md` with PRD-vs-PLAN clarifier
10. [ ] Dogfood: run `prd review` on `.spectacular/PRD.md`; snapshot before any edits
11. [ ] Test on a fresh blank project (target: usable PRD in <15 minutes)

## Resolved

- **Inline mini-refine vs on-demand full refine** → both. Mini runs inline after each grill answer (pattern-match only); full review runs at end-of-grill and on `prd review`.
- **Kit selection** → grill asks as first question; project-local override via `.spectacular/templates/prd/kits/<kit>.md`.
- **Existing PRD behavior** → ask refine-vs-fresh; "fresh" auto-snapshots to `PRD@vN.md` first.
- **Per-request PRDs** → not v2, never. Anti-pattern. Requests use PLAN + TASKS only.

## Open questions

- Auto-detect kit from existing project context (vs. always asking)? Deferred to v2.
- Should the grill ever offer to draft slots from existing files (`STACK.md`, `README.md`, etc.) instead of asking cold? Deferred — risks crossing the line into research-pipeline territory the scope explicitly rejects.

---

## v1.1 — Slot alignment + kit-prep

Added 2026-05-21 after auditing prd-craft against the post-rework canonical docs (PRD v2.0 / ARCHITECTURE v1.0). Slot names diverged from the root PRD; Vision + Deliverable sections live in the root PRD but aren't grill slots.

### Why now

The dogfood task ("run `prd review` against root PRD.md") can't pass cleanly while the grill's slot names and slot set differ from what the canonical root PRD uses. v1.1 fixes the alignment before dogfood runs. Kit refactor is delegated to [[kits-as-plugins]] — v1.1 only **prepares** the base for that refactor.

### Scope (v1.1)

**Slot rename — align with root PRD.md**
- Slot 2: "Who it's for" → **"Target users"**
- Slot 3: "What success looks like" → **"Goals & success criteria"**

**Slot additions — bring base in line with root PRD shape**
- Add slot 0: **"Vision"** — the philosophical "what this is" (one paragraph, narrative)
- Add slot: **"Deliverable"** — what concretely ships (between Target users and Goals & success criteria, per root PRD order)

**New base slot order (8 slots)**
1. Vision
2. Problem
3. Target users
4. Deliverable
5. Goals & success criteria
6. Non-goals
7. Constraints
8. First milestone

**Files touched**
- `templates/prd/base.md` — rename + add 2 slots
- `references/prd-grill.md` — update slot loop count + names + intro
- `references/prd-review.md` — update gate checklist to 8 slots; vague-word list unchanged
- `references/prd-refine.md` — update slot references in patterns

**Not in v1.1 — punted to [[kits-as-plugins]]**
- Refactoring kits to diff-only format
- Kits declaring `triggers-docs`
- Composition rules

### Anti-pattern

Don't refactor kits in v1.1. Keep them as standalone-with-extras temporarily; [[kits-as-plugins]] does the formal refactor against the now-stable 8-slot base.

### Success criteria (v1.1)

- Root `.spectacular/PRD.md` passes `prd review` against the updated 8-slot gate (snapshot to `PRD@v2.0.md` first if any edits land)
- A grill run on a throwaway project produces a PRD whose slot names match root PRD verbatim
- All 5 existing kits still load through grill without error (even though refactor is deferred)

### v1.1 dependency for downstream requests

Implementation order locked 2026-05-21:

```
prd-craft v1.1  →  doc-writer  →  kits-as-plugins  →  smart-init
                                   ↘                   ↘
                                    doctor (can land in parallel with smart-init)
```

- **[[doc-writer]]** waits for v1.1 — generalizes grill/refine/review into a shared engine + adds the doc registry. PRD becomes the first registry entry, modeled against the v1.1 8-slot base.
- **[[kits-as-plugins]]** waits for doc-writer — kits become registry-aware deltas, declaring `triggers-docs` for the engine to consume.
- **[[smart-init]]** waits for kits-as-plugins — CLI consumes registry + kit `triggers-docs` to drive selective scaffolding.
- **[[doctor]]** waits for doc-writer — needs the registry to know what "correct" scaffold looks like; can land in parallel with smart-init.

v1.1 itself is intentionally narrow: slot rename + add Vision/Deliverable. The engine generalization is **not** in scope here.
