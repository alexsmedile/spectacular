---
status: planned
updated: 2026-05-24
related:
  - PLAN.md
  - discovery.md
---

# Tasks — substrate-clarity

## M1 — Discover ✅ DONE (2026-05-24)

- [x] Grilling session with user — Decisions 1, 1+2 (revised), 3, 4, 5, 5b
- [x] Constraint surfaced: grill + refine are skill-only
- [x] Locked model captured in `discovery.md`

## M2 — Spec-refine: open decisions + changeset

- [x] Decision 6 — rules-file frontmatter schema locked (see discovery.md)
- [x] Decision 7 — drop "engine" entirely, use "skill" or name verb directly (see discovery.md)
- [x] `spec.md` drafted — full file-by-file changeset, risk register, suggested PR sequencing
- [ ] User reviews `spec.md`, confirms scope before M3 begins
- [ ] Resolve M3 open question: `docs-rules.md` covers two doc-ids (recommend: skip — deprecated)
- [ ] Resolve M5 open question: explicit `[needs-deepening]` marker vs auto-only (recommend: auto-only in v1.4.0)

## M3 — Registry demotion ✅ DONE (2026-05-24)

- [x] `spectacular snapshot skills/spectacular/references/doc-registry.md` → `doc-registry@v1.md`
- [x] Rename `doc-registry.md` → `doc-index.md` (rewrite as human catalog; old file removed)
- [x] Add frontmatter to 7 existing rules files (prd, plan, tasks, roadmap, personas, pack, docs)
- [x] Create 6 minimal new rules files (principles, architecture, stack, agents, spec, decisions)
- [x] Sweep all `doc-registry` → `doc-index` refs across skill (18 files via sed)
- [x] Update SKILL.md routing table, references index, doc-id list (v1.4.0)
- [x] Bump SKILL.md version → 1.4.0-dev
- [x] Update CLAUDE.md repo-structure table
- [x] Update own PLAN.md `related:` path
- [x] `cli/spectacular doctor` exits 0 errors / 0 warnings / 8 info

## M4 — Mode collapse: `reps` → `grill-each` ✅ DONE (2026-05-24)

- [x] roadmap-rules.md: `mode: grill-each` (declared in frontmatter at M3)
- [x] personas-rules.md: `mode: grill-each` (declared in frontmatter at M3)
- [x] PRD / PLAN / convention-pack: `mode: grill` (no migration needed — sugar for grill-wide)
- [x] SKILL.md description + routing + references index updated for v1.4.0
- [x] `references/grill.md` rewritten: mode resolution + sub-mode dispatch + flag override documented
- [x] grill-loop algorithm + heuristic specified in `grill.md` § 3a
- [x] Templates/roadmap/base.md mode-comment updated to `grill-each`
- [x] grill.md examples replaced (correct ROADMAP example, added PERSONAS + grill-loop + stub-with-override)

## M5 — Build grill-loop ✅ DONE at spec level (2026-05-24)

- [x] grill-loop algorithm specified in `references/grill.md` § 3
- [x] Heuristic locked: length < 30 chars OR vague-word hit OR placeholder string OR explicit gate-check fail
- [x] Documented per-doc default = declared mode; flag override (`--loop`) wins per session
- [ ] Manual test on PRD with `--loop` (deferred to M8 verification)
- [ ] Manual test on ROADMAP with `--loop` (deferred to M8 verification)
- [ ] Decide whether any doc should default to `grill-loop` (decision: not in v1.4.0; revisit in v1.4.x based on usage)

## M6 — Agentic/mechanical verb split: CLI redirect + docs ✅ DONE (2026-05-24)

- [x] `cli/spectacular`: doc-verb intercept added before "Unknown command" die
  - [x] Detects `spectacular <doc> grill|refine|review` for any KNOWN_DOCS doc
  - [x] Detects `spectacular <doc>` (no verb) — agentic dispatch handoff
  - [x] Unknown verbs (`<doc> bogus`) get a "must be: grill|refine|review" error
- [x] `KNOWN_DOCS` constant extended with plan + tasks
- [x] Updated "Unknown command" help text to list agentic verbs
- [x] Reused existing `skill_verb_message` infrastructure
- [ ] Update `docs/commands.md` with agentic/mechanical table (deferred to M7 cleanup)
- [ ] Update `docs/installation.md` (deferred to M7 cleanup)

## M7 — Conceptual cleanup pass ✅ DONE (2026-05-24)

- [x] M7.1 — engine word sweep: ~37 architectural occurrences → 0 across 15 files (excluding engineer/engineering domain words)
- [x] M7.2 — Overrides → Rules H1 standardization: all 13 rules files use "Rules" noun
- [x] M7.3 — docs/commands.md: added agentic/mechanical verb table + grill sub-modes section + v1.4.0 doc list
- [x] M7.4 — SPEC.md doc-writing capability rewritten for v1.4.0; ROADMAP forward-looking refs updated; CONTRIBUTING.md new-doc-type guide updated; CHANGELOG history preserved as period-accurate
- [x] M7.5 — Doctor exits 0 errors / 0 warnings / 8 info; remaining doc-registry refs are historical (snapshots, archive, CHANGELOG) or intentional (rename pointer)
- [x] CLAUDE.md request-status updated to "active (M3-M7 done)"
- [x] Verb × mode matrix lives in doc-index.md (decided during M3)

## M8 — Doctor + tests + ship ✅ DONE (2026-05-24)

- [x] `cli/spectacular doctor` exits 0 errors / 0 warnings / 8 info
- [x] `cli/spectacular --version` reports 1.4.0
- [x] Agentic-verb redirect verified: `spectacular roadmap grill` prints friendly redirect
- [x] CHANGELOG [1.4.0] entry written: Breaking (3) / Added (8) / Changed (7) / Fixed (2 — codex G1+G2) / Migration / Process
- [x] Manifests bumped: `.claude-plugin/plugin.json`, `.claude-plugin/marketplace.json` (metadata + plugins[]), `.codex-plugin/plugin.json`, README badge → all 1.4.0
- [x] Manual bumps: `cli/spectacular` SPECTACULAR_VERSION → 1.4.0, SKILL.md frontmatter → 1.4.0
- [x] `check-manifests.sh`: project-level versions aligned (6 locations at 1.4.0)
- [x] Commit `3ca8ec9` — chore: release v1.4.0 — substrate-clarity (46 files, +1554 / -135)
- [x] Tag `v1.4.0` created + pushed
- [x] GitHub release: https://github.com/alexsmedile/spectacular/releases/tag/v1.4.0
- [ ] User-triggered: `/plugin marketplace update spectacular` (can't be done remotely)
- [ ] Archive request (separate step — typically next session)

## Deferred / Open

- [x] **Decision 6** (rules-file schema) — LOCKED during M1 (see discovery.md)
- [ ] **Decision 7** ("engine" rename) — to lock at M2
- [ ] Auto-generate `doc-index.md` from rules-file frontmatter (out of scope for v1.4.0; revisit if maintenance burden surfaces)
- [ ] `grill-loop` heuristic refinement — first ship with simple rule (vague-word + slot-length); improve in v1.4.x if needed
