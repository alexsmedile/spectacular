---
status: verified
priority: medium
owner: alex
updated: 2026-07-12
build: b30
summary: "Split verify.md into walk-only core (loads at review→verified, ~half today's 6.2k) + verify-authoring.md (2-of-6 canonical, fold patterns, promoting checks to scripts); patch the TASKS template so v2-deferred items scaffold as - [~]."
related:
  - PRD.md
---

# Plan — verify-split

## Goal

Finish the reference-doc efficiency arc: the `review → verified` walk loads only walk content (verify.md halves), authoring-time content moves to its own doc, and the TASKS template stops scaffolding deferred items as open checkboxes.

## Constraints

- Same recipe as the build/bug and b28 splits: runtime core keeps everything the orchestrator acts on; no behavior loss — every moved section grep-verified present in its new home before the cut.
- verify.md keeps its filename (walk-context `[[verify]]` links in build-workflow, doctor, doctor-substrate stay valid).
- plan-rules.md § 2-of-6 stays the compact copy; its "canonical" note re-points from "verify.md Part 2" to verify-authoring.md.
- The five 2-of-6-context links re-point: lifecycle.md ×3, init-workflow.md, new-request.md → plan-rules § 2-of-6 or verify-authoring as fits the sentence.
- Template patch covers every live copy of the TASKS scaffold (skills/spectacular/templates/tasks/base.md confirmed; check repo-root templates/tasks/ for a second copy).

## Understanding

### How it works now

verify.md (6.2k tok) is three merged parts (v1.20.0): Part 1 the walk (check kinds, walk loop, VERIFY-LOG, gate, retrospective), Part 2 the 2-of-6 rule + fold patterns + VERIFY.md shape, Part 3 promoting checks to scripts. Since b28, scaffold/grill route to plan-rules § 2-of-6 — so Parts 2–3 load pointlessly at the verified gate, the one place verify.md still loads whole. templates/tasks/base.md:42 scaffolds `## v2 (deferred)` items as `- [ ]`, counting them as open; the CLI already renders `- [~]` correctly (proven live on b28: card flipped to `15/15 (+4 def)`).

### What changes

verify.md → Part 1 only (~3k). New verify-authoring.md ← Parts 2–3 (2-of-6 canonical source, when-to-fold patterns, standalone VERIFY.md shape, test-script promotion). SKILL.md's verification-routing rows re-point the authoring-time rows; doc-index.md gains the new doc's row; the five 2-of-6-context links re-point. templates/tasks/base.md's v2 example rows become `- [~]`; tasks-rules.md gains one line ("v2 items use `- [~]` until promoted").

### What stays the same

The walk's content and the `review → verified` gate mechanics; all walk-context `[[verify]]` links; plan-rules' compact table (text unchanged, only its canonical pointer); the CLI (`spectacular verify` still redirects to the skill; progress counting already handles `[~]`).

## Decisions

- verify.md = walk, new authoring doc (grilled 2026-07-12): halves the verified-gate load, walk links stay valid. Rejected thin-index verify.md (pure indirection hop — the exact pattern this arc removes) and "read Part 1 only" notes (section-scoped reads are unenforceable).
- Template fix rides here, not in cli-gate-ergonomics (grilled): templates/tasks/base.md is a skill-side asset — same surface as the reference docs this request edits.

## Milestones

- M1 — verify.md walk-only + verify-authoring.md created; both docs' cross-pointers set
- M2 — Link re-point sweep: SKILL.md rows, doc-index.md, lifecycle.md ×3, init-workflow.md, new-request.md, plan-rules canonical note
- M3 — TASKS template v2 rows scaffold as `- [~]`; tasks-rules.md one-liner

## Tasks

See `TASKS.md`.

## Dependencies

- [[self-healing-optimization]] (b28) — created plan-rules § 2-of-6, which makes Parts 2–3 movable

## Validation

- M1 — run: `wc -c` verify.md ≤ 14000 bytes and verify-authoring.md exists; assertable: every Part 2–3 section heading greps in verify-authoring.md and not in verify.md; walk sections (check kinds, walk loop, VERIFY-LOG shape, gate) grep in verify.md
- M2 — assertable: `grep -rn "verify.md Part 2" skills/spectacular/` → 0 hits outside verify-authoring.md's own history note; the five re-pointed links grep with their new targets; `spectacular doctor links` clean
- M3 — run: `spectacular new tmp-scaffold-test` in a sandbox → TASKS.md v2 rows are `- [~]`; `spectacular request tmp-scaffold-test` progress shows `0/N (+2 def)`; sandbox removed after

## Deliverables

- Slimmed `references/verify.md` + new `references/verify-authoring.md`
- Re-pointed SKILL.md, doc-index.md, lifecycle.md, init-workflow.md, new-request.md, plan-rules.md
- Patched `templates/tasks/base.md` (+ root copy if present) + tasks-rules.md line
- CHANGELOG entry
