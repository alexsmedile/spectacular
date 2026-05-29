---
status: planned
priority: medium
owner: alex
updated: 2026-05-29
summary: "Advisory cross-request awareness — related:/depends-on:/blocks: in PLAN frontmatter, inverse-link resolver, doctor links area; carries a doctor-memory staleness side-rider from FEEDBACKS.md"
related:
  - PRD.md
  - ../../ARCHITECTURE.md
  - ../../ROADMAP.md
target_version: v1.12.0
---

# Plan — cross-request-links

## 1. Goal

Let a request declare its relationships to other requests (`related:` / `depends-on:` / `blocks:`) with inverse links resolved automatically, so an agent picking up a request sees at a glance which other in-flight work it touches — surfacing conflicts as advisory warnings (no locking).

## 2. Constraints

- **Advisory only.** No locking, no blocking, no automatic coordination. Conflict resolution stays human judgment. This is the explicit boundary — don't drift into orchestration.
- **Inverse links are computed, never stored.** `blocks:` on request A shows as `blocked-by:` on request B at read time; never written to B's frontmatter (single source of truth).
- **Dangling references are surfaced, not fatal.** A link to a non-existent/archived slug is a doctor warning, not an error.
- **Reuses the existing `related:` convention.** `depends-on:` / `blocks:` are additive siblings to the `related:` field already in every PLAN.

## 3. Milestones

- M1 — Frontmatter schema extension: `depends-on:` / `blocks:` documented in ARCHITECTURE.md alongside existing `related:`.
- M2 — Inverse-link resolver: given all PLAN frontmatter, compute the bidirectional graph (`blocks` ↔ `blocked-by`) at read time.
- M3 — `doctor links` area: validates link integrity, flags dangling references to missing/archived slugs. **Side-rider (from FEEDBACKS.md 🟢):** add a staleness flag to `doctor memory` — mirror the existing `sessions`/`feedback`/`ideas` convention. Naive age-check is the v1 floor; the valuable contradiction-check (a memory referencing a blocker a later session/decision overturns) is deferred to v2.
- M4 — `status` advisory surface: `spectacular status` shows "⚠ active request `<slug>` is related to active request `<other>`"; `spectacular new` prompts to declare relationships when slug keywords match existing requests.
- M5 — Examples + ship: 2 example projects demonstrate the link graph; CHANGELOG entry; plugin bump to v1.12.0.

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
- M5 — Example projects render a correct link graph; manifests at v1.12.0.

## 7. Deliverables

- ARCHITECTURE.md frontmatter-schema extension (`depends-on:` / `blocks:`)
- Inverse-link resolver (CLI, mechanical)
- `doctor links` area (validate + flag dangling)
- `spectacular status` advisory line + `spectacular new` relationship prompt
- 2 example link-graph demonstrations
- CHANGELOG [1.12.0] entry
