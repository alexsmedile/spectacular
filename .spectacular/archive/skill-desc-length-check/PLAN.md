---
status: archived
priority: medium
owner: alex
updated: 2026-06-30
build: b10
summary: "Doctor + pre-commit check flagging skill SKILL.md descriptions over Codex's 1024-char limit"
related:
  - PRD.md
archived: 2026-06-30
---

# Plan — skill-desc-length-check

## 1. Goal

Detect when a skill's `SKILL.md` `description` exceeds Codex's hard 1024-char limit — at `doctor` time and (optionally) at commit time — so a release can never silently ship a description Codex refuses to load.

## 2. Constraints

- Codex's limit is **1024 chars** for `description`; Claude Code's is **1536**. The check must use 1024 as the failing threshold so the stricter runtime governs. (The gap is exactly why this regresses silently — Claude Code loads fine while Codex skips the skill.)
- **M1 resolved:** Codex measures `description` **alone**, not the concatenation. Proof: the v1.17.2 patch dropped `description` from 1146 → 986 while `description + when_to_use` stayed at 1253 (> 1024) — and the error cleared. If Codex measured the combined block, 1253 would still have failed. So the check gates on `description` alone, **error >1024 / warning >1000**. `when_to_use` is not counted. (Warn band set to 1000 so the deliberate ~986 steady state is a clean pass — see [[D7]].)
- Must handle YAML literal-block (`|`) descriptions — the indented multi-line form spectacular itself uses — not just single-line `description:`.
- Pure-bash, no new deps (consistent with the existing `cli/spectacular` doctor areas).
- Repo-shape-aware: this repo's skill is at `skills/spectacular/SKILL.md`, but the check should generalize to any `**/SKILL.md` it's pointed at.

## Understanding

### How it works now

`spectacular doctor skill` (`check_skill()` in `cli/spectacular`, ~line 6498) checks install path, symlink resolution, and `skills.lock` presence. It does **not** inspect `SKILL.md` frontmatter content. There is no length guard anywhere in the doctor suite or in `scripts/hooks/pre-commit`. The v1.17.2 release was needed precisely because the 1146-char description shipped undetected.

### What changes

- Add a frontmatter-length sub-check to the doctor `skill` area: parse the `description` (+ `when_to_use`) YAML literal block, measure char count, emit `warning` near the limit and `error` over 1024.
- Optionally wire the same logic into `scripts/hooks/pre-commit` (git-guard managed) so it fails the commit before a bad description can be tagged.
- Decide whether the shared measurement logic lives inline in `check_skill()` or in a small `scripts/` helper both the doctor and the hook can source (DRY).

### What stays the same

- Existing `check_skill()` install/symlink/lock checks are untouched — this is additive.
- No change to `description` *content* policy or to the skill schema. This only measures and flags.
- Claude Code's 1536 limit is irrelevant to the threshold; we gate on 1024.

## 3. Milestones

- M1 — **Threshold confirmed.** Verify exactly what Codex measures (description alone vs description+when_to_use) and at what cap; document the answer. Reproduce the v1.17.2 case as the baseline.
- M2 — **Doctor check.** `spectacular doctor skill` flags descriptions over 1024 (error) and near it, >1000 (warning), correctly parsing YAML literal blocks. Shows actual char count in the finding.
- M3 — **Pre-commit guard (optional).** Same measurement runs in `scripts/hooks/pre-commit`, failing the commit on any `SKILL.md` over the limit. Shared logic, no duplication.

## 4. Tasks

See `TASKS.md`.

## 5. Dependencies

- None blocking. Touches `cli/spectacular` (`check_skill`) and optionally `scripts/hooks/pre-commit` (git-guard territory — coordinate so the guard isn't clobbered on next hook regen).

## 6. Validation

- M1 — A written note (in this PLAN or DECISIONS) stating Codex's measured field + cap, with the v1.17.2 numbers as evidence.
- M2 — Point doctor at a fixture SKILL.md >1024 → `error`; at one ~1010 → `warning`; at one ≤1000 → `pass`. The live `skills/spectacular/SKILL.md` (986) → pass.
- M3 — Staging a SKILL.md edit pushing description >1024 → `git commit` rejected with a clear message; trimming back → commit succeeds.

## 7. Deliverables

- A description-length sub-check in the doctor `skill` area, showing the char count.
- (Optional) a matching pre-commit guard.
- Shared measurement helper if both consume it.
