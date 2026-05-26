---
status: review
updated: 2026-05-25
related:
  - PLAN.md
---

# Tasks — feedback-loop

## M0 — Grill the request ✅

- [x] Resolve PLAN §3 open questions (all 7 decisions locked in §4)
- [x] Lock feedback file frontmatter schema (v1.6.0 shape locked)
- [x] Decide canonical command: `feedback-loop` (with hyphen)
- [x] Decide aliases as hidden routing (not first-class)
- [x] Set target version: v1.6.0

## M1 — Substrate + ref doc ✅

- [x] Write `skills/spectacular/references/feedback-loop.md` (mode spec, 5-step loop, file shape, examples)
- [x] Write `skills/spectacular/references/feedback-rules.md` (dispatch + behavior for the engine)
- [x] Add `skills/spectacular/templates/feedback/entry.md` template
- [x] Add `feedback` row to `references/doc-index.md` catalog
- [x] Add `.spectacular/feedback/` + `requests/<slug>/feedback/` to `ARCHITECTURE.md` directory trees
- [x] Add PRINCIPLES.md §9 "Feedback ≠ verification ≠ benchmark"

## M2 — Skill mode ✅

- [x] Add feedback-loop trigger routing + proactive-surfacing rules to SKILL.md
- [x] Add aliases (`iterate`, `experiment`, `test`, `probe`, `try`) as hidden routing in SKILL.md
- [x] Document three-checkpoint proactive offer policy (milestone tick / review entry / archive)
- [x] Doc IDs registered list bumped to v1.6.0 with `feedback`

## M3 — CLI surface ✅

- [x] `cmd_feedback_loop` orchestrator with `new`, `list`, `resolve`, `archive` subverbs
- [x] Helpers: `_feedback_dir`, `_feedback_iter_all`, `_feedback_find`
- [x] Top-level dispatch routes `feedback-loop` + 5 hidden aliases to same handler
- [x] `--request <slug>` flag scopes to request folder
- [x] `--next-action` on resolve + `--promote` (advisory; no silent memory writes)
- [x] `top_usage` shows `feedback-loop` only — aliases stay hidden per contract
- [x] Smoke-tested: new / list / list --status / resolve / archive / alias `probe new` / alias `try archive` / request-scoped

## M4 — Doctor area ✅

- [x] Add `feedback` to `DOC_AREAS` + `doctor_parse_args` accepted areas
- [x] Implement `check_feedback`: stale-open (>30d), malformed frontmatter, orphan back-refs
- [x] Wire into `run_areas` dispatcher
- [x] Update `references/doctor-areas.md` with feedback area docs
- [x] Bumped help string from "10 areas" to "14 areas" (matches actual DOC_AREAS count)
- [x] Smoke-tested: empty / one-entry / stale-entry scenarios

## M5 — Dogfood (in progress)

- [x] **Run feedback-loop session on feedback-loop itself** (CLI ergonomics after M0-M4) — resolved `ship-as-is`. See `feedback/2026-05-25-feedback-loop-cli-ergonomics-after-m0-m4.md`
- [x] Apply the 4 fixes surfaced by the dogfood pass:
  - [x] `list` drops DATE column + truncates slugs >32 chars to 29+`...`
  - [x] `--promote` renamed to `--promote-hint` (honest naming; no silent memory writes)
  - [x] `resolve` requires `--next-action` (error if missing); guides toward `park` if undecided
  - [x] Help text updated for `--promote-hint`
- [x] PLAN.md gets `feedback:` back-ref to the entry (per the bidirectional contract)
- [x] Doctor confirms clean (0 errors / 0 warnings / 0 info — 1 entry, 0 open)
- [ ] Run feedback-loop session on `grill-each` ergonomics (long PRDs) — deferred, separate from this request
- [ ] Run feedback-loop session on soft-db memory mutators (`remember`, `decide`) — deferred
- [ ] Run feedback-loop session on pack discovery UX — deferred
- [ ] Validate proactive-surfacing offers at the three checkpoints — needs a *separate* request to dogfood (this one was checkpoint-triggered manually)

## M6 — Release (ready to ship)

- [x] Bump manifests to v1.6.0 (SKILL.md, cli/spectacular, .claude-plugin/plugin.json, .codex-plugin/plugin.json, .claude-plugin/marketplace.json ×2)
- [x] CHANGELOG entry under Added (feedback-loop mode, doc-type, CLI verbs, doctor area, PRINCIPLES §9, aliases, surfacing, back-refs, promotion, template) and Changed (SKILL routing, doc-index row, ARCHITECTURE trees, doctor-areas docs, top_usage)
- [x] README mention: workspace tree (`feedback/`), skill commands table, CLI reference block
- [x] Final doctor pass — feedback area clean (1 entry, 0 open)
- [ ] Tag + push via `/release` (user-driven)
- [ ] Flip request status to `review`, then archive after release (will trigger a feedback-loop offer per checkpoint #3 🙂)

## Status

Status flips to `review` — implementation complete, dogfood passed, release artifacts ready. Awaiting user to run `/release` for the v1.6.0 tag.
