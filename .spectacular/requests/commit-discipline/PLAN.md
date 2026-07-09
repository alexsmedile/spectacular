---
status: planned
priority: medium
owner: alex
updated: 2026-07-06
build: b20
summary: "Add a soft periodic-commit nudge for in-progress code work — Spectacular currently ships zero commit guidance"
related:
  - ../../PRD.md
  - ../../POLICY.md
  - ../../PRINCIPLES.md
---

# Plan — commit-discipline

## Goal

Give Spectacular a soft, non-blocking nudge to commit in-progress code at natural
checkpoints (milestone boundaries and session end), so long agent runs don't end
with hours of uncommitted work — closing a gap that hits Codex harder than Claude.

## Constraints

- **Soft, never blocking.** This is a `warn`-severity nudge, never a `block` gate. Committing is the user's call; Spectacular reminds, it does not gate transitions on a clean tree (P8 — humans decide, agents propose).
- **No git mutation by Spectacular.** The skill/CLI must never run `git commit` itself. It surfaces the nudge; the agent (Claude/Codex) or human decides and acts. Auto-committing user code is out of scope and out of character.
- **Tool-agnostic core, tool-specific delivery.** The POLICY entry is shared prose (both Claude and Codex read the same `skills/`). Any hook wiring must land in both `hooks.json` and `hooks-codex.json` or neither — no silent divergence.
- **Small (P3).** One POLICY entry + at most one wired hook. No new commands, no config surface, no new doc type.
- **Distinct from `snapshot`.** `spectacular snapshot` versions canonical *docs* into `_snapshots/`; it is unrelated to git. This request must not blur the two — the nudge is about `git commit` of source, not doc snapshots.

## Understanding

### How it works now

Spectacular ships **zero** commit guidance (confirmed by full read of POLICY.md,
PRINCIPLES.md, and all skill references):

- **POLICY.md** — 9 work-phase hooks, none about git. The closest is
  `@SessionEnd → summarize-before-handoff` (`severity: warn`), which nudges a
  *prose summary*, not a commit.
- **PRINCIPLES.md** — git appears only as a storage fact ("these files are
  git-committed"), never as an action the agent should take.
- **Plugin hooks** — `hooks/hooks.json` and `hooks/hooks-codex.json` are both
  empty `{"hooks": {}}` stubs. Nothing fires on `Stop` / `SessionEnd`.
- **Claude vs Codex** — treated identically (same `skills/` dir, both empty hook
  stubs). Spectacular gives Codex no extra discipline to offset its tendency to
  forget commits on long jobs.
- The task-tracking model (`SKILL.md` § Task tracking) *assumes* commits happen
  as milestones complete, but never instructs, reminds, or gates on them.

### What changes

- Add one `@Implementation` POLICY entry (`commit-checkpoint`, `severity: warn`):
  nudge a local commit at each milestone boundary.
- Extend the existing `@SessionEnd` guidance so "summarize before handoff" also
  suggests committing outstanding work (or explicitly noting why not).
- (M2, gated) Wire a `Stop`/`SessionEnd` plugin hook — in **both** `hooks.json`
  and `hooks-codex.json` — that emits the reminder when the tree is dirty. Only
  if the prose-only nudge proves insufficient in practice.

### What stays the same

- No lifecycle transition gates on commit state — `planned→active→review→verified`
  still gate only on understanding + verification.
- `snapshot` semantics untouched.
- No auto-commit, ever. No new CLI verb. No config field.

## Milestones

- M1 — **POLICY nudge shipped (prose-only).** A `commit-checkpoint` entry under
  `@Implementation` + extended `@SessionEnd` prose. `spectacular policy @Implementation`
  surfaces it; both Claude and Codex read it. This is the minimum verified slice.
- M2 — **(gated) Wired reminder hook.** A `Stop`/`SessionEnd` hook in both
  `hooks.json` + `hooks-codex.json` that fires a dirty-tree reminder. Gated on M1
  proving insufficient — do not build speculatively (P11).

## Tasks

See `TASKS.md`.

## Dependencies

- None. POLICY.md is always-set and already structured for new work-phase entries;
  `spectacular policy` already reads it. `doctor policies` already validates the
  POLICY schema, so M1 needs no new plumbing.

## Validation

- M1 — `spectacular policy @Implementation` lists `commit-checkpoint` with its
  linked principle; `spectacular policy @SessionEnd` shows the extended guidance;
  `spectacular doctor policies` stays green (valid schema, principle link resolves).
  Judgment check: the prose actually reads as a soft nudge, not a gate.
- M2 — with a deliberately dirty tree, ending a session surfaces the reminder in
  both a Claude and a Codex run; with a clean tree, silence. Hook is present and
  identical in both hook files (`doctor` — or a diff — confirms parity).

## Deliverables

- A `commit-checkpoint` policy entry in `.spectacular/POLICY.md` under `@Implementation`.
- Extended `@SessionEnd` prose in the same file.
- (M2, gated) matching `Stop`/`SessionEnd` entries in `hooks/hooks.json` +
  `hooks/hooks-codex.json`.
- A short note in the relevant skill reference (e.g. `active-request.md` or
  `policy-injection.md`) so the nudge is discoverable, not just enforced.

## Open questions

- **Milestone boundary detection** — does the `@Implementation` nudge key off a
  TASKS.md `- [ ]` → `- [x]` milestone tick (the natural checkpoint), or fire more
  loosely? Lean: tie the prose to milestone completion, matching the existing
  feedback-loop surfacing checkpoints — no new detection machinery.
- **Is M2 ever needed?** If the POLICY prose alone changes agent behavior, the
  wired hook is pure cost. Ship M1, observe, decide. (This is why M2 is gated.)
- **Codex-specific phrasing** — Codex forgets more; does its nudge need to be more
  assertive than Claude's, or is identical prose fine? Lean identical first
  (shared `skills/`), escalate only if Codex keeps forgetting.
