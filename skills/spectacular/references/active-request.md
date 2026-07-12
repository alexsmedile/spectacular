---
description: Continue active work — session state, task tracking, milestone progress.
when_to_use: Actively working on a request that already exists.
---

# Active Request — Continuing Work

Triggered by: skill detects user is actively implementing a request, or user asks to continue work on a specific slug.

> **@Implementation policy gate.** Before moving a request `planned → active`, run `spectacular policy @Implementation` and follow every active policy. The default blocker is `understand-before-change`: PLAN.md must have a filled `## Understanding` section (or a `UNDERSTANDING.md` with the same three subheads) — satisfy it or stop. See [policy-injection.md](policy-injection.md).

---

## On entering active state

When a request transitions from `planned` → `active`:

1. Update `PLAN.md` frontmatter: `status: active`, `updated: <today>`
2. Create `SESSION.md` if it doesn't exist (see template below)
3. Load the full request context: `PLAN.md`, `TASKS.md`, `SESSION.md`, `specs/index.md`, relevant `specs/<capability>.md`

---

## Superseding a live plan

When findings invalidate part of an active PLAN (a diagnosis corrected, a milestone made obsolete, an estimate disproven):

1. **Prepend** a block: `## SUPERSEDED <date> — <one line: what changed>` containing the corrected understanding.
2. Mark the affected original sections `(superseded <date>, kept for history)` — **never delete disproven content.** The disproof trail — what was believed, what evidence killed it — is the plan's most valuable content for anyone (human or agent) who re-approaches the problem later.
3. Rewrite frontmatter `summary:` to the *current* understanding, ≤2 sentences. The summary reflects now; the body preserves the trail.
4. New design calls that emerge go to `## Decisions` (chose X over Y — because Z), not into prose.

This is the sanctioned form of the pattern healthy projects invent ad hoc ("CORRECTED DIAGNOSIS … superseding §1–§5 below"). One convention, one shape, no hand-rolled layering.

---

## SESSION.md template

```md
---
updated: <today>
---

# Session — <slug>

## Current state
<What's been done so far>

## Active task
<What's being worked on right now>

## Blockers
<Anything blocking progress>

## Next actions
- 
- 
```

SESSION.md is **committed to git** — it's the team's operational record of in-progress work.

Update SESSION.md at natural breakpoints: after a meaningful chunk of work, when blocked, or when handing off.

---

## Task tracking

Monitor `TASKS.md` for completion signals:

- When a task item is checked off, acknowledge it
- When **all items in a group** are checked, note it proactively — this is also the `commit-checkpoint` policy's moment (see `POLICY.md` `@Implementation`): surface a soft reminder to `git commit` the milestone's work before moving on
- When **all items in TASKS.md** are checked → tier-reveal one line: `Next: spectacular advance <slug> to move it to review.` (see `lifecycle.md`)

---

## Context loading during active work

Follow the context-loading table in `.spectacular/AGENTS.md` (authoritative — don't re-derive a read list here). Always prefer targeted per-capability files over the full `specs/` tree; the top-level `specs/index.md` is cheap and always relevant.

---

## Handling blockers

When a blocker is encountered:
1. Note it in `SESSION.md` under `## Blockers`
2. Surface it in the next briefing
3. Propose resolution path if possible
