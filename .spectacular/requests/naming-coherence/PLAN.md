---
status: verified
priority: medium
owner: alex
updated: 2026-06-28
build: b15
summary: "Resolve verb naming inconsistencies (feedback-loop→feedback, lifecycle promote→advance, drop convention-pack alias) and add contextual next-step suggestion (tier-reveal) plus a `spectacular next` verb."
related:
  - ../../PRINCIPLES.md
  - ../../SPEC.md
---

# Plan — naming-coherence

> **Origin (2026-06-27):** The verb surface grew organically to ~34 verbs with
> three naming frictions: (1) `feedback-loop`/`feedback`/`feedback-rules` — three
> names for one concept; (2) `promote` is overloaded — `idea promote` (elevate)
> and lifecycle `promote <slug>` (move forward through states) are different
> operations sharing a verb; (3) `convention-pack` and `pack` are documented as
> equivalent — two names, one thing. Separately, the system orbits "what's the
> next action" but has no verb that just says it, and never teaches the next
> probable verb contextually. This request fixes naming + adds the discovery aids.

## 1. Goal

Make the verb surface coherent: one name per concept, a forward-lifecycle verb that reads correctly, a single `spectacular next`, and contextual next-step suggestions during flows — without breaking existing muscle memory.

## 2. Constraints

- **Backwards compatible.** Renamed verbs keep their old spelling as a deprecation alias (hidden from help, prints a one-line "renamed to X" notice). No hard breaks.
- **Tier-reveal is a SKILL change, not CLI.** The CLI stays the deterministic mutator; suggesting the next probable verb is skill behavior. Reuse the existing checkpoint machinery (the feedback-loop spec already defines milestone-tick / status→review / archive checkpoints).
- **`spectacular next` is read-only.** It surfaces the single highest-priority next action (the briefing's CTA) — it never mutates.
- **commands.md (skill) keeps ALL verbs.** Only the human `docs/` get tiered presentation (separate concern; may be a follow-up).

## Understanding

### How it works now

- **feedback:** verb is `feedback-loop` (+ 5 now-removed aliases); doc is `feedback-rules.md`; "loop" is baked into the verb name though it's really a mode. Surface: ~8 files.
- **promote:** `cmd_promote` (cli/spectacular:4460) drives the forward state machine `planned→active→review→verified`; `idea promote` is a separate elevate-to-request op. Both spelled `promote`.
- **convention-pack:** full doc-id form, documented as equivalent to `pack <verb>`. Surface: ~9 files.
- **next action:** status.md/onboarding.md compute "the single highest-priority next action" but there's no verb to print it on demand, and flows don't suggest the next probable verb when they end.

### What changes

- **Rename `feedback-loop` → `feedback`** (verb). `feedback` becomes the verb, "loop" is the mode/behavior, `feedback-rules.md` stays as the doc. Old `feedback-loop` → hidden alias.
- **Rename lifecycle `promote <slug>` → `advance <slug>`** (forward through states). `idea promote` keeps `promote` (namespaced, unambiguous). Old `promote <slug>` → hidden alias with notice.
- **Drop `convention-pack`** as a documented form; `pack` only. Keep `convention-pack` as a silent alias one release, then remove (coordinate with cli-debt-removal philosophy).
- **Add `spectacular next`** — prints the briefing's single highest-priority next action, nothing else.
- **Tier-reveal (skill):** flow docs (new-request.md, active-request.md, lifecycle.md, archive.md) end with one contextual "probable next" suggestion, reusing the checkpoint hooks.

### What stays the same

- All operations' behavior — only names + additive discovery.
- `idea promote`, `snapshot`, and every non-renamed verb.
- CLI remains the only mutator; skill remains the suggester.

## 3. Milestones

- **M1 — `advance` rename.** `cmd_promote` reachable as `advance`; `promote <slug>` aliases it with a deprecation notice. `--help`, docs/commands.md, SKILL.md routing updated. Tests cover both spellings.
- **M2 — `feedback` rename.** `feedback-loop` → `feedback` across skill refs + CLI; old form hidden alias. `feedback-rules.md` unchanged.
- **M3 — Drop `convention-pack` form.** `pack` is canonical everywhere; `convention-pack` silent alias + flagged for removal.
- **M4 — `spectacular next` verb.** Read-only; prints the single next action. Wired into dispatch + help.
- **M5 — Tier-reveal suggestions.** Flow docs end with one contextual next-step hint via existing checkpoints.

## 4. Tasks

See `TASKS.md`.

## 5. Dependencies

- Soft link to [[lifecycle-undo]] — undo records the operation; if `advance` lands first, undo references the new name. Not blocking (undo keys on operation, not spelling).
- Aligns with cli-debt-removal's deprecation-then-remove pattern for the `convention-pack` alias lifecycle.

## 6. Validation

- M1/M2/M3 — both old and new spellings work; old prints a deprecation notice (where chosen) and is absent from help; `tests/cli/` assert both paths.
- M4 — `spectacular next` on a workspace with an active request prints that request's next action; on empty workspace prints the "start something" CTA; mutates nothing (clean `git status`).
- M5 — after `spectacular new`, the flow surfaces a single grill suggestion; at status→review, a verify suggestion. Never more than one, never mid-flow.

## 7. Deliverables

- Renamed verbs + deprecation aliases (`advance`, `feedback`, `pack`-only).
- New `spectacular next` CLI verb + help.
- Tier-reveal hints in 4 skill flow docs.
- Updated docs/commands.md + SKILL.md routing tables.
- `tests/cli/` coverage for renames + `next`.
- SPEC + specs/cli SPEC sync.
