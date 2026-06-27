---
status: archived
priority: medium
owner: alex
updated: 2026-06-28
build: b12
summary: "A reverse gear for lifecycle mutations — `spectacular undo` reverts the last state transition or move (advance/archive/idea promote) so a mis-step doesn't require manual file surgery."
related:
  - PRD.md
  - ../../ROADMAP.md
  - ../../PRINCIPLES.md
archived: 2026-06-28
---

# Plan — lifecycle-undo

> **Origin (2026-06-27):** Spectacular is mutation-heavy and entirely forward-moving.
> `promote/advance`, `archive`, and `idea promote` all change state (and two of them
> physically move directories and rewrite inbound links). There is no reverse gear:
> an accidental `archive` or a premature `advance` currently requires hand-editing
> frontmatter and `git mv`-ing directories back. `snapshot` versions *docs* but does
> nothing for *lifecycle*. For a system this destructive, a trustworthy undo builds
> confidence to use the verbs freely.

## 1. Goal

Add `spectacular undo` that reverts the **most recent lifecycle mutation** (status transition or move) for a request/idea, restoring frontmatter, directory location, and inbound links to their prior state — without the user touching files by hand.

## 2. Constraints

- **CLI-owned (deterministic mutator).** Per the v0.7.0 mutation principle, undo is a CLI verb. The skill may *suggest* it but never performs the reversal itself.
- **Git-aware but not git-dependent.** Must work in a non-git workspace (falls back to plain `mv`), mirroring how `cmd_archive` already branches on `git rev-parse`.
- **Reverse only what Spectacular wrote.** Undo restores Spectacular-managed state (PLAN/TASKS frontmatter, the request dir location, the inbound-link rewrites). It is NOT a general filesystem undo and never touches user content inside the docs.
- **No new persistent state if avoidable.** Prefer deriving the reversal from what's on disk (current status → previous status is deterministic for the linear state machine) over maintaining a separate undo journal — until the move-reversal proves that insufficient.
- **One step back, not a stack** (v1). Multi-level undo is a later milestone if pain surfaces.

## Understanding

### How it works now

Three verbs mutate lifecycle state, in increasing complexity:

1. **`advance`/`promote <slug>`** (`cmd_promote`, cli/spectacular:4460) — pure `fm_set status` on PLAN.md (and TASKS.md). Linear machine: `planned → active → review → verified`. Trivially reversible (set status back one step + `fm_touch`).
2. **`archive <slug>`** (`cmd_archive`, cli/spectacular:4093) — the hard one. It (a) `fm_set status archived` + `archived: <date>`, (b) `git mv requests/<slug> → archive/<slug>`, (c) sed-rewrites inbound `../<slug>/` links to `../../archive/<slug>/` across all *other* requests. Undo must reverse all three.
3. **`idea promote <slug>`** (`cmd_idea`, ~cli/spectacular:2192) — `fm_set status promoted` on the idea, scaffolds a request, moves the idea source to `archive/ideas/`. Undo must restore the source and remove the scaffolded request (or leave it — open question, see M-questions).

There is no record of "what was the last mutation." Status alone tells you current state but archive's move + link-rewrite is not recoverable from frontmatter.

### What changes

- New `cmd_undo` in `cli/spectacular` + dispatch entry + `--help`.
- A lightweight **last-mutation marker**: each of the three mutators writes a one-line breadcrumb (verb + slug + prior-status + timestamp) to `.spectacular/.last-mutation` (gitignored — it's session ephemera, not team state). Undo reads it, reverses, then clears it. This is the minimal state that makes archive's move reversible without a full journal.
- `archive`/`advance`/`idea promote` each gain a 2-3 line "record breadcrumb" call (one shared helper).
- Skill: `lifecycle.md` + `archive.md` add a one-line "↩ revert with `spectacular undo`" hint after the mutation confirmation (tier-reveal, not a new flow).

### What stays the same

- The forward verbs' behavior is untouched — undo is additive.
- `snapshot`/versioning is unrelated and unchanged (docs vs lifecycle are separate axes).
- No undo for doc edits, memory/decision/session entries, or `init` — out of scope.

## 3. Milestones

- **M1 — Undo a status transition.** `spectacular undo` reverts the last `advance`/`promote` (status back one step on PLAN + TASKS, `fm_touch`). Breadcrumb written by `cmd_promote`, read + cleared by `cmd_undo`.
- **M2 — Undo an archive.** Reverse the dir move (`archive/<slug> → requests/<slug>`, git-aware), restore `status` to its pre-archive value, drop the `archived:` field, and reverse the inbound-link rewrites (`../../archive/<slug>/ → ../<slug>/`).
- **M3 — Undo an idea promote.** Restore the idea source from `archive/ideas/`, reset its status, and prompt about the scaffolded request (remove vs keep — user choice).
- **M4 — Guardrails + clarity.** `undo` with no breadcrumb → friendly "nothing to undo." Breadcrumb older than the current HEAD / stale → refuse with explanation. `--dry-run` prints what would be reversed. Confirm before any move.

## 4. Tasks

See `TASKS.md`.

## 5. Dependencies

- Soft dependency on the **naming cleanup** (`advance` rename of `promote`) — if that ships first, undo records the new verb name. Not blocking; undo keys on the operation, not the verb spelling.
- `.gitignore` must cover `.spectacular/.last-mutation` (add in M1).

## 6. Validation

- M1 — `new` → `advance` to active → `undo` → PLAN status is `planned` again; `.last-mutation` cleared. Test in `tests/cli/`.
- M2 — `archive` a verified request → `undo` → dir is back under `requests/`, status restored, inbound links in sibling requests point to `../<slug>/` again. Assert dir location + grep link form.
- M3 — `idea promote` → `undo` → idea source back in `ideas/`, status reset; scaffolded request handled per user choice.
- M4 — `undo` on a fresh workspace prints "nothing to undo" and exits 0; `--dry-run` mutates nothing (assert clean `git status`).

## 7. Deliverables

- `spectacular undo [--dry-run]` CLI verb + `--help`.
- Shared breadcrumb helper + `.last-mutation` writes in the three mutators.
- `.gitignore` entry for `.last-mutation`.
- Skill hints in `lifecycle.md` + `archive.md`.
- `tests/cli/undo.test.sh` covering M1–M4.
- SPEC.md + `specs/cli/SPEC.md` sync at archive time.

## Open questions (resolve during grill)

1. **Idea-promote undo + scaffolded request:** ✅ RESOLVED (2026-06-28) — **prompt, default leave.** undo restores the idea to `ideas/`; the scaffolded request dir is left in place unless the user confirms removal. Removing a dir the user may have edited is the destructive choice — opt-in only.
2. **Breadcrumb vs git:** ✅ RESOLVED (2026-06-28) — **timestamp vs file mtimes.** undo compares `.last-mutation`'s timestamp to the affected files' mtimes and refuses (stale) if any changed after the recorded mutation. Full git reflog/status verification is v2.
3. **Multi-level undo:** ✅ RESOLVED (2026-06-28) — **single-level only.** One breadcrumb; undo reverses the last mutation. A stack is v2 if demand surfaces.
