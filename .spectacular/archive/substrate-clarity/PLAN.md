---
status: archived
phase: spec-refine
priority: medium
owner: alex
updated: 2026-05-26
target_version: v1.4.0
summary: "Rename doc-registry → doc-index (human catalog), move dispatch into rules-file frontmatter, collapse `reps` mode into grill+sub-modes (wide/each/loop), lock agentic vs mechanical verb split"
related:
  - discovery.md
  - ../../ARCHITECTURE.md
  - ../../specs/index.md
  - ../../../skills/spectacular/references/doc-index.md
archived: 2026-05-26
---

# Plan — substrate-clarity

## 1. Goal

Resolve two conceptual confusions in the substrate surfaced by the personas adversarial review (codex findings G1 + G2), with a model locked during the M1 discovery grill.

**Locked conceptual model** (see `discovery.md` for the full session):

1. **Three orthogonal axes** — never collapse: **Phase** (lifecycle), **Verb** (action), **Mode** (doc interaction style).
2. **Grill is a verb and a mode family.** Mode values: `grill` (= grill-wide) / `grill-wide` / `grill-each` / `grill-loop`. No separate sub-mode field — the mode value carries the variant. User overrides via flag (`--wide` / `--each` / `--loop`).
3. **`reps` mode collapses into `grill-each`.** Conceptual modes: grill-family + append + stub + freeform + reference (5 conceptual; 8 values).
4. **Grill and refine are skill-only** — require an LLM. CLI never dispatches them; it redirects with a friendly message. Review is mixed (CLI runs structural checks; skill runs semantic).
5. **Doc-registry demoted to `doc-index.md`** — human catalog only. Dispatch fields move into each rules file's frontmatter. Every doc gets a rules file (consistency over brevity).
6. **Verb × Mode matrix is now defined for every cell** — including `grill × stub` (polite hint), `grill × freeform` (open-ended prompt + on-the-fly slot generation), `refine × append` (ask user scope).

## 2. Constraints

- **No behavior regression on existing docs.** PRD, ROADMAP, PERSONAS, PLAN, etc. must keep working exactly as they do today after the rename + restructure. This is a clarity + naming refactor, not a semantic change.
- **`reps` → `grill` migration must be transparent.** Existing rules files saying `mode: reps` get rewritten to `mode: grill` + `default-sub-mode: grill-each`. Behavior identical. No user action required.
- **CLI surface adds the grill-redirect message but stays mechanically minimal.** `cli/spectacular` learns to recognize agentic verbs and print the redirect; it does not gain new dispatch logic.
- **All rules files gain frontmatter** (currently none have any). Stub docs (PRINCIPLES, ARCHITECTURE, STACK, AGENTS, SPEC, DECISIONS) get newly-created minimal rules files — frontmatter-only is fine.
- **Snapshot everything we touch.** `doc-registry.md` snapshots before rename. SKILL.md snapshots before update. ARCHITECTURE.md if updated. Every existing rules file snapshots before its frontmatter is added.

## 3. Milestones

- **M1 — Discover ✅ DONE** (2026-05-24). Locked model in `discovery.md`. Decisions 1-5b agreed.
- **M2 — Spec-refine: open decisions + changeset.**
  - ✅ Decision 7 locked: drop "engine" entirely, use "skill" or verb name directly.
  - Produce a concrete file-by-file changeset list as a `spec.md` under this request, using the schema locked in discovery.md Decision 6.
- **M3 — Registry demotion.**
  - Snapshot `skills/spectacular/references/doc-registry.md` → `doc-registry@v1.md`.
  - Rename to `doc-index.md`, rewrite content (drop dispatch-contract framing).
  - Move per-doc dispatch YAML into each existing rules file's frontmatter.
  - Create minimal rules files for the 6 stub/append docs (PRINCIPLES, ARCHITECTURE, STACK, AGENTS, SPEC, DECISIONS).
  - Verify skill engine reads dispatch from rules-file frontmatter, not the index.
- **M4 — Mode collapse: `reps` → `grill-each`.**
  - Rewrite all rules files currently declaring `mode: reps` (roadmap, personas) to `mode: grill-each`.
  - Existing grill-mode docs (PRD, PLAN, convention-pack) stay as `mode: grill` (sugar for `grill-wide`) — no migration needed.
  - Update SKILL.md routing + references index.
  - Wire up `--wide` / `--each` / `--loop` flag handling in the grill verb (skill-side).
- **M5 — Build grill-loop.**
  - Implement the wide-then-deep flow in skill engine (grill verb instructions in `references/grill.md` or equivalent).
  - First pass: walk all slots fast (one-line answers ok). Second pass: revisit only slots flagged as vague/incomplete.
  - Test on PRD + ROADMAP.
- **M6 — Agentic/mechanical verb split: CLI redirect + docs.**
  - `cli/spectacular`: detect agentic verbs (`grill`, `refine`) and print redirect message instead of attempting dispatch.
  - Update CLI `--help` to mark agentic verbs ("requires Claude Code / Codex").
  - Update `docs/commands.md` with the agentic/mechanical split.
- **M7 — Conceptual cleanup pass.**
  - Sweep `skills/spectacular/` for stale uses of "generic engine," "mode: reps," confused mode-vs-verb language.
  - Update SKILL.md, ARCHITECTURE.md, AGENTS.md, onboarding.md per the locked model.
  - Update CLAUDE.md repo-structure references.
  - Document the verb × mode matrix in a single canonical location (likely `doc-index.md` or new `verb-mode-matrix.md`).
- **M8 — Doctor + tests + ship.**
  - `spectacular doctor` exits 0 across all areas.
  - Update any test referencing the old registry path, "generic engine," or `mode: reps` wording.
  - CHANGELOG entry under [1.4.0] calling out: rename, mode collapse, grill sub-modes, CLI redirect.
  - Bump manifests via git-guard `bump-manifests.sh 1.4.0`.
  - Commit, tag `v1.4.0`, push, GH release, marketplace update (user-triggered).

## 4. Tasks

See `TASKS.md`.

## 5. Dependencies

None upstream. Downstream: the cleaner conceptual model should make `convention-pack-modules` easier to scope when it eventually opens.

## 6. Open questions

- **Open ergonomics question:** when `grill-loop` revisits slots, how does it decide "vague/incomplete"? Heuristic in the skill, or explicit `[needs-deepening]` markers in the doc? Lock during M5 build.

## 7. Discovery artifacts

- `discovery.md` (locked 2026-05-24) — full session record. Source of truth for what was agreed.
