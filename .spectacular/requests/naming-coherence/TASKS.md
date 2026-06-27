---
status: verified
updated: 2026-06-28
related:
  - PLAN.md
---

# Tasks — naming-coherence

## v1

### M1 — advance rename (lifecycle promote → advance)
- [x] Add `advance` dispatch → `cmd_promote`
- [x] Keep `promote <slug>` as alias with one-line deprecation notice
- [x] Update `--help`, docs/commands.md, SKILL.md routing (lifecycle.md, verify.md, troubleshooting.md too)
- [x] Tests: both `advance` and `promote` paths (scenario 7b)

### M2 — feedback rename (feedback-loop → feedback)
- [x] Rename verb to `feedback`; `feedback-loop` → hidden alias
- [x] Update skill refs + CLI + docs
- [x] Leave `feedback-rules.md` doc name unchanged
- [x] Fixed latent bug: backtick command-substitution in feedback usage heredoc (`<<EOF` → `<<'EOF'`)

### M3 — drop convention-pack form
- [x] Make `pack` canonical: doc-id `convention-pack` → `pack` (pack-rules.md, doc-index.md, SKILL.md, grill.md)
- [x] `convention-pack` → recognized alias (`doc-id-aliases:`), flagged for removal next release
- [x] Fixed gap: `pack new|grill|refine|review` now redirect to skill (were dying as "unknown pack verb")

### M4 — spectacular next
- [x] Add read-only `cmd_next` → prints single highest-priority next action
- [x] Wire into dispatch + help
- [x] Test: active/planned/empty workspace, rejects args, mutates nothing (scenario 7c)

### M5 — tier-reveal suggestions (skill)
- [x] new-request.md ends with grill suggestion
- [x] active-request.md (advance→review) / lifecycle.md (verb-name note) suggest next at checkpoints
- [x] archive.md already suggests spec-sync (pre-existing)
- [x] Rule: one suggestion max, never mid-flow

## v2 (deferred)
- [ ] Remove `convention-pack` alias entirely
- [ ] Tiered presentation in human docs/
