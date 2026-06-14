---
status: verified
priority: medium
owner: alex
updated: 2026-06-08
summary: "Advisory cross-request awareness — related:/depends-on:/blocks: in PLAN frontmatter, inverse-link resolver, doctor links area; carries a doctor-memory staleness side-rider from FEEDBACKS.md"
related:
  - PRD.md
  - ../../ARCHITECTURE.md
  - ../../ROADMAP.md
build: b5
---

# Plan — cross-request-links

## 1. Goal

Let a request declare its relationships to other requests (`related:` / `depends-on:` / `blocks:`) with inverse links resolved automatically, so an agent picking up a request sees at a glance which other in-flight work it touches — surfacing conflicts as advisory warnings (no locking).

## 2. Constraints

- **Advisory only.** No locking, no blocking, no automatic coordination. Conflict resolution stays human judgment. This is the explicit boundary — don't drift into orchestration.
- **Inverse links are computed, never stored.** `blocks:` on request A shows as `blocked-by:` on request B at read time; never written to B's frontmatter (single source of truth).
- **Dangling references are surfaced, not fatal.** A link to a non-existent/archived slug is a doctor warning, not an error.
- **Reuses the existing `related:` convention.** `depends-on:` / `blocks:` are additive siblings to the `related:` field already in every PLAN.

## Understanding

### How it works now

Every PLAN.md frontmatter already carries `related: [list of slugs/paths]` — a flat, one-directional list with no semantics (it doesn't distinguish "touches" from "depends on" from "blocks") and no inverse resolution (A listing B tells you nothing when you're reading B). `doctor links` already exists and validates that `related:` targets resolve — that's the hook this request extends. There is no `status`-surface awareness of inter-request relationships; an agent picking up a request can't see what else in flight touches it without reading every other PLAN. `fm_get`/`fm_get_list` read frontmatter; the read verbs (`requests`, `request <slug>`) render per-request views.

### What changes

- **Frontmatter schema** gains two additive sibling fields to `related:` — `depends-on:` and `blocks:` (lists of request slugs). Documented in ARCHITECTURE.md.
- **An inverse-link resolver** computes the bidirectional graph at read time: `blocks: [B]` on A means B is `blocked-by: A` — *computed, never written to B*. Same for `depends-on:` ↔ `required-by:`.
- **`doctor links`** extends to validate the new fields + flag dangling references to missing/archived slugs (warning, not error). Side-rider: add a staleness flag to `doctor memory` (age-check floor, mirroring sessions/feedback/ideas).
- **`spectacular status`** gains an advisory line when active requests are related; **`spectacular new`** prompts to declare relationships on slug-keyword match.

### What stays the same

**Advisory only** — no locking, no blocking, no auto-coordination; conflict resolution stays human. The single source of truth is each request's own forward declaration; inverse links are *derived*, never stored (no write-back to other files). `related:` keeps working unchanged — the new fields are additive. Dangling refs stay non-fatal. No change to lifecycle state ownership (PLAN frontmatter) or any existing verb's output shape (advisory lines are additive; `--format json` unaffected).

## Design decisions (2026-06-02)

- **Inverse labels:** `blocks:` ↔ `blocked-by:`, `depends-on:` ↔ `required-by:` (computed, never stored).
- **Path-resolution fix (folded into M3):** the existing `check_links` resolves `related:` targets relative to each file's *own* dir — the real cause of the 7 recurring `related: PRD.md` false warnings (it looks for `requests/<slug>/PRD.md`). Fix: root-aware resolution — bare slugs → `requests/` (then `archive/`); root-doc names (`PRD.md`, `ARCHITECTURE.md`, `ROADMAP.md`, …) → `.spectacular/`. Needed anyway for slug-based `depends-on`/`blocks`.
- **M4 data flow:** `status` is a *skill* verb. CLI computes the graph (mechanical), skill renders the advisory (agentic). The new `spectacular links` verb is the data source.
- **Surfaces (3):** inverse links visible in (1) `spectacular request <slug>` detail, (2) `doctor links`, (3) a new dedicated `spectacular links [--json]` verb (whole-graph dump).
- **Archived dependencies = satisfied:** a `depends-on:` pointing at an *archived* request resolves as met — show `depends-on: X ✓ (shipped)`, not a dangling warning. Only a slug matching *nothing* (active or archived) is dangling.
- **`doctor memory` staleness threshold = 180 days** — memory holds durable facts; the gradient is sessions 4h < feedback 30d < ideas 90d < **memory 180d**. Conservative nudge, not a nag.
- **`spectacular links` (no slug) default = only requests WITH edges** (the actual graph); `--all` includes unlinked requests.
- **Scope:** all 5 milestones ship in the target release (see frontmatter `target_version`).

## 3. Milestones

- M1 — Frontmatter schema extension: `depends-on:` / `blocks:` documented in ARCHITECTURE.md alongside existing `related:` (+ inverse-label table, computed-not-stored rule).
- M2 — Inverse-link resolver (CLI, mechanical): given all PLAN frontmatter, compute the bidirectional graph at read time. Surface in `spectacular request <slug>` detail **and** a new `spectacular links [--json]` verb (whole-graph dump).
- M3 — `doctor links` extension: validate the new fields + flag dangling refs to missing/archived slugs. **Fold in the path-resolution fix** (root-aware target resolution — kills the 7 false `related:` warnings). **Side-rider (FEEDBACKS.md 🟢):** add a staleness flag to `doctor memory` (age-check floor, mirroring sessions/feedback/ideas; contradiction-check deferred to v2).
- M4 — `status` advisory surface: CLI `links` emits the data; skill `status.md` renders "⚠ active request `<slug>` is related to active `<other>`"; `spectacular new` prompts to declare relationships when slug keywords match existing requests.
- M5 — Examples + ship: 2 example projects demonstrate the link graph; CHANGELOG entry; plugin bump to v1.16.0.

## 4. Tasks

See `TASKS.md`.

## 5. Dependencies

- Unblocks the Icebox item "ROADMAP-as-source-of-truth enforcement" (needs this inverse-link registry first).
- Independent of [[verify-walk]] and [[visual-layer]]; can ship in any order relative to them.
- The roadmap already models v2.0.0's `Depends on: v1.15 → v1.16 → v1.17` chain — this request is what lets tooling eventually *read* that chain.

## 6. Validation

- M1 — ARCHITECTURE.md documents all three relationship fields + the computed-not-stored rule.
- M2 — A fixture where A `blocks` B resolves to B showing `blocked-by: A` without B's file being edited.
- M3 — `doctor links` flags a deliberately-dangling reference; passes clean on valid links. `doctor memory` flags a deliberately-aged entry; quiet on fresh entries.
- M4 — Two related active requests trigger the `status` advisory; `new` on a keyword-matching slug prompts for relationships.
- M5 — Example projects render a correct link graph; manifests at target release.

## 7. Deliverables

- ARCHITECTURE.md frontmatter-schema extension (`depends-on:` / `blocks:` + inverse-label table)
- Inverse-link resolver (CLI, mechanical)
- New `spectacular links [<slug>] [--json]` read verb (whole-graph or per-request)
- Inverse links surfaced in `spectacular request <slug>` detail
- `doctor links` extension (new fields + dangling) **+ root-aware path-resolution fix** (kills 7 false `related:` warnings)
- `doctor memory` staleness flag (side-rider)
- `spectacular status` advisory line (CLI data → skill render) + `spectacular new` relationship prompt
- 2 example link-graph demonstrations
- CHANGELOG [1.13.0] entry
