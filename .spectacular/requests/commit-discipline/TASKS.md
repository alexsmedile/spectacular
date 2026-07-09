---
status: parked
updated: 2026-07-09
related:
  - PLAN.md
---

# Tasks — commit-discipline

<!--
  Executable checklist for one request.
  Lives at: .spectacular/requests/<slug>/TASKS.md

  Rules:
  - Group tasks by milestone using `### M<N> — <name>` headings.
  - Use `- [ ]` for open, `- [x]` for done. No other bullet syntax.
  - `status:` in frontmatter should match parent PLAN.md.
  - Tasks are owned by the user. Engine never adds/removes/reorders tasks.
-->

## v1

### M1 — POLICY nudge shipped (prose-only)
- [x] Add `commit-checkpoint` entry under `@Implementation` in `.spectacular/POLICY.md` (severity: warn, linked principle, soft-nudge prose tied to milestone completion)
- [x] Extend `@SessionEnd` prose so summarize-before-handoff also suggests committing outstanding work (or noting why not)
- [x] Add a discoverability note in `active-request.md` (or `policy-injection.md`) so the nudge is visible, not just enforced
- [x] Verify: `spectacular policy @Implementation` lists it; `spectacular doctor policies` stays green; prose reads as nudge not gate

### M2 — (gated) Wired reminder hook
> Do not start until M1 has shipped and proven insufficient in practice (P11 — earn the step).
- [ ] Add a `Stop`/`SessionEnd` hook entry to `hooks/hooks.json` (dirty-tree reminder)
- [ ] Add the identical entry to `hooks/hooks-codex.json` (parity — no divergence)
- [ ] Verify: dirty tree → reminder fires in both a Claude and a Codex run; clean tree → silent; hook files match
