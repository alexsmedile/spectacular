---
status: planned
priority: medium
owner: alex
updated: 2026-06-27
build: b14
summary: "Onboarding work: (a) dedup — onboarding.md references status.md for the shared flow, keeping only onboarding-specific intent; (b) guided first-run — an ushered new→prd grill→first request flow for a brand-new (empty) workspace."
related:
  - ../../ARCHITECTURE.md
---

# Plan — onboarding-dedup

> **Origin (2026-06-27):** `onboarding.md` (first-run-on-existing-workspace flow)
> and `status.md` (no-arg briefing) share ~95% of their read sequence: config →
> root docs → requests → SPEC → memory → briefing, even down to "same format as
> status.md". Onboarding is genuinely a distinct moment (taking over mid-project,
> first orientation) and the file should stay — but the duplicated steps drift
> independently. Principle: what is the same, remains the same — referenced once,
> not copied.

## 1. Goal

Two onboarding improvements: (a) make `onboarding.md` reference `status.md` for the shared briefing flow (dedup, single source of truth); (b) add a **guided first-run** — an ushered flow that walks a brand-new user through new→prd grill→first request, instead of making them discover the sequence verb by verb.

**Distinction:** `onboarding.md` (existing) = orient on an *existing* workspace with prior work. Guided first-run (new) = usher through a *fresh/empty* workspace. Different moments; both live here because both are "first contact."

## 2. Constraints

- **Keep `onboarding.md` as a separate file.** It is a distinct flow with distinct intent; this is dedup, not a merge.
- **status.md is the canonical briefing engine.** The shared sequence lives there; onboarding points to it.
- **No behavior change to first-run.** Onboarding still runs the substrate check, uses the takeover tone, and surfaces 2-3 gap observations.
- **Progressive disclosure preserved** — onboarding must still load the minimum needed for a first briefing.

## Understanding

### How it works now

Both docs independently spell out:
1. read config.yaml → 2. read AGENTS.md / root canonical docs → 3. read requests/*/PLAN.md frontmatter → 4. read SPEC.md + capability specs → 5. read memory/ → 6. produce a briefing.

`status.md` adds the substrate-check-on-failure + SPEC drift check. `onboarding.md` adds: substrate check that *always* runs on first invocation, a "taking over mid-project" tone, and a "what to flag on first look" observations table (max 2-3). The briefing format is explicitly "same as status.md".

### What changes

- `onboarding.md` replaces its steps 1-7 read sequence with: "Run the `status.md` read + briefing flow, with these onboarding-specific deltas:" followed by only the deltas (always-run substrate check, takeover tone, gap-observations table).
- `status.md` stays the single owner of the read sequence + briefing format.
- Any genuinely-onboarding-only content (backwards-compat pre-split detection, the example takeover briefing) stays in `onboarding.md`.

### What stays the same

- `status.md` content (it's the source; onboarding now points *to* it).
- First-run behavior end to end.
- The two files both existing.

## 3. Milestones

- **M1 — Identify the shared spine vs onboarding deltas.** Mark which lines of onboarding.md are duplicated status.md flow vs onboarding-specific.
- **M2 — Refactor onboarding.md to reference status.md** for the shared sequence; keep only the deltas + onboarding-specific sections.
- **M3 — Verify existing-workspace onboarding still produces the takeover briefing** and status still works standalone.
- **M4 — Guided first-run flow.** On an empty/new workspace, instead of an empty briefing, usher the user: offer `new` (or PRD grill) as one step, then the first request, then point at `spectacular next`. CLI may provide the entry (`init --walk` or auto-detect empty on `/spectacular`); the skill drives the ushering. One clear path, no verb-discovery burden. Tier-0 onboarding moment.

## 4. Tasks

See `TASKS.md`.

## 5. Dependencies

- Shares a principle with [[rules-files-audit]] (reference-don't-duplicate) but independent.

## 6. Validation

- M2 — onboarding.md no longer restates the 6-step read sequence verbatim; a reader is sent to status.md for it.
- M3 — simulate first invocation: substrate check runs, briefing uses takeover tone + ≤3 observations; `/spectacular` (status) on a warm workspace unchanged.
- M4 — `/spectacular` on an empty workspace ushers (offers new/PRD → first request → `next`) instead of printing an empty briefing; the path is one clear step at a time, no verb menu dumped.

## 7. Deliverables

- Refactored `references/onboarding.md` (deltas + pointer to status.md).
- `references/status.md` confirmed as the single briefing-flow owner (minor edits only if needed to be referenceable).
- Guided first-run flow (skill flow doc + any CLI empty-workspace detection / `init --walk` entry).
