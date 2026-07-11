# Changelog

All notable changes are documented here.
Format: [Keep a Changelog](https://keepachangelog.com/en/1.0.0/)

---

## [Unreleased]

### Fixed

- **`spectacular new` no longer re-issues a duplicate build id when `config.last_build` drifts behind the roadmap ledger.** The counter previously trusted only `config.yaml`; if a request was slotted into the ledger without `last_build` catching up, `new` handed out a `bN` the ledger already used. It now reconciles to `max(config.last_build, ledger_max) + 1`, self-healing regardless of how the drift arose. Regression test + fix log `F4`.

## [1.30.1] — 2026-07-11

### Rule-layer legibility for mixed-model agents

A patch pass making `PRINCIPLES.md` + `POLICY.md` legible to lower-intelligence
agents that pattern-match surface tokens instead of resolving cross-references —
without softening anything for frontier models, which read the added redundancy
as confirmation. All changes are additive prose; no severity logic changed.

### Changed — POLICY.md (v1.4 → v1.5) + PRINCIPLES.md (v1.2 → v1.3)

- **Blocking gates are now visually unmissable.** The 4 block-severity policies (`understand-before-change`, `verification-present`, `confirm-before-write`, `snapshot-before-overwrite`) carry a `⛔ **BLOCKING**` body line naming the exact consequence. The marker sits on a body line, *not* the heading — the heading is the parsed policy id, so a suffix would corrupt `_policy_lookup` + config overrides. Warns stay unmarked so the marker means something.
- **Excuse/Reality tables extended to the highest-temptation warns** (`scope-down`, `build-order`, `earn-the-verification`, `understand-before-change`) — the pattern previously reserved to `@Debugging`. Each names the *specific* rationalization an agent reaches for, plus a Law and a Red-flag line.
- **Legitimate-override clauses** (`**Override:**`) on 3 warns, so a frontier model knows when a skip is policy-sanctioned rather than guessing.
- **`@Debugging` `Law:` lines now inline their mechanism** (`spectacular fix list`, grep-the-callers) instead of asserting the rule and leaving the how one bullet away.
- **PRINCIPLES #10 vs #11** each gained a scope-vs-sequence contrast box so the two restraint principles stop collapsing into "do less" for weaker readers.
- `_policy_principle_line` (CLI) now skips a leading blockquote so the new #10/#11 boxes don't leak into the rendered principle summary. `spectacular policy` already renders the block/warn split natively (`⛔` vs `·`).

### Added — `stance-layer` request (b26, planned)

- Records the six-note improvement triage + all grill decisions in `ideas/stance-layer.md`, and cuts the `stance-layer` request: an `architectural-stance` `@Planning` warn policy (fires only on a real architectural fork; offers `spectacular decide`, never forces) + a `grade` label for PLAN frontmatter (`prototype | mvp | standard | production`) that `status` surfaces so prototype-green never reads as production-green. The severity *dial* was deliberately rejected as over-engineering. Planned, not built.

### Fixed — CLI snapshot path scoping

- Snapshot ineligibility for PLAN files is now scoped to request plans only (`requests/*/PLAN.md`); a `PLAN.md` elsewhere (e.g. `docs/PLAN.md`) falls through to the generic "not a registered canonical doc" path. Tightens both `diagnose_path_error` and the snapshot guard, with a covering mutator test scenario.
- Follows `78a4af0`'s path-error-suggestion work for CLI mutators and the touch/snapshot contract split.

## [1.30.0] — 2026-07-10

### The Agent Fleet update

v1.30.0 turns Spectacular's agent story from a single debug fleet into a coherent
**8-agent fleet** spanning *discover / apply / review / verify* across the *fix* and
*build* directions, all obeying one boundary — closed handoff in, findings/diff/pass-fail
out, **the orchestrator is the only mutator**. Agent defs now live in `agents/` at the
plugin root (source of truth; `.claude/agents/*.md` are relative symlinks), all conform
to the Claude Code subagent spec, and every agent is reachable by the skill via the
workflow arcs or the spec-sync flow.

### Added — `spec-reviewer` agent (the specs' guardian)

- **A read-only guardian of the spec files** (`specs/index.md` + `specs/*.md`), the doc-review analog of `code-reviewer`. Its primary job is keeping the spec **true to the code**: alongside spec-writing quality (well-formed, concrete, no vague theme-gestures), it runs a **currency** check — for each capability the spec claims, it greps the code/CLI/tests for evidence the claim still holds and cross-checks the roadmap/CHANGELOG ledger, classifying drift as **stale** (spec claims a capability the code lost/changed), **gap** (shipped code the spec omits), or **premature** (spec claims something unshipped). Every currency finding cites its evidence, so "this is stale" is a fact, not an opinion. Coherence-vs-intent (contradicts PRD/PRINCIPLES/personas) is a secondary axis; user-fitness is an optional low-confidence NOTE, not a verdict.
- Returns a ranked punch list, never rewrites. Preserves the fleet boundary: `grill` (interactive) and `refine` (mutation) stay on the orchestrator/main thread; `spec-reviewer` is the read-only gate that tells the orchestrator whether either is needed. Wired into the `spec-sync`/archive flow as an optional currency pass before a SPEC-DELTA is written. Dogfooded on this repo's own `specs/index.md` — caught 6 real drift issues (5 unmentioned agents, a shipped-but-described-as-future check, a `hold:` gap, `debug/`→`debugs/`, `memory/`→`memories/`, and a lagging CLI version constant), all fixed + re-verified clean.

### Added — the three new agents wired into the workflow arcs (fleet-arc-wiring, b25)

- **The skill now dispatches `repo-explorer`, `code-reviewer`, and `test-verifier`** — they existed in `agents/` but no workflow arc routed to them, so the orchestrator never reached for them. Wired as **optional, judgment-gated** steps (same worth-it economics as the fan-out gate — not a step every change passes through):
  - **`repo-explorer`** → `build-workflow.md` gains a **Step 0a — map unfamiliar ground**: when the orchestrator can't write a milestone's Approach because the subsystem is unfamiliar, dispatch the explorer for a structured map, then plan. The build-side mirror of the debug Investigator; the build↔bug comparison table's "discover role" is now honest.
  - **`code-reviewer`** → an optional review at Step 3 of **both** arcs (before ticking / before logging a fix), and at the `review → verified` gate over the full request diff. Returns findings; the orchestrator triages and routes fixes.
  - **`test-verifier`** → an optional arms-length verify at Step 3 of both arcs — dispatched when a builder/fixer *self-reported* the pass or blast radius is medium+ ("the agent that built it shouldn't be the only one to grade it").
- `SKILL.md` trigger rows + `doc-index.md` updated so the routing is discoverable. Doc-only change; `doctor` 0 errors. Tracked as request `fleet-arc-wiring` (b25).

### Added — three new fleet agents (repo-explorer, code-reviewer, test-verifier)

- **The agent fleet now spans a discover / apply / review grid across both directions (fix / build), plus specialists.** Three agents added to `agents/`, each conforming to the Claude Code subagent spec (aliased model, scannable `description`, restricted tools, no plugin-unsupported fields) and symlinked into `.claude/agents/`:
  - **`repo-explorer`** (read-only, opus) — the build-side analog of `debug-investigator`. Maps an unfamiliar subsystem before a milestone is planned: entry points, the sibling pattern to mirror, integration seams, blast radius — returns a structured map with `file:line` anchors. Illuminates; never plans the change.
  - **`code-reviewer`** (read-only, opus) — reviews a bounded diff across five lenses (correctness, structure, security, perf, dead-code), **all by default or a `FOCUS` subset**. Returns severity-ranked findings with fix *direction*, never the fix. Folds five waterfall-roster roles into one lens-parametric agent.
  - **`test-verifier`** (apply-only, sonnet) — independently confirms a change: runs a named check, or writes a test to a *closed* behavioural spec in the project's own test style, and reports honest pass/fail with real output. `Write`/`Edit` touch test files only — never the code under test, never the ledger. Bounces when verifying would become deciding.
- Design rationale (which of the 14 archived `_archive/agent-fleet/` roles to adopt vs cut) is captured in the fleet decision doc: most archived roles were either the orchestrator itself, duplicates of `spec-builder`, or review *lenses* masquerading as standing agents — folded into `code-reviewer`. Repo maps in `CLAUDE.md` + `AGENTS.md` updated; all 7 agents validate.

### Changed — agents/ is now the source of truth (plugin-root convention)

- **The four fleet agents moved from `.claude/agents/` to a root `agents/` directory** (`debug-investigator`, `debug-fixer`, `debug-researcher`, `spec-builder`), matching the plugin-root convention (agents belong at the plugin root alongside `skills/`, not in the Claude-local dir). `.claude/agents/*.md` are now **relative symlinks** (`../../agents/<name>.md`, git-tracked as mode 120000) so Claude Code still loads them. Edit the real files in `agents/`. Repo maps in `CLAUDE.md` + `AGENTS.md` updated; plugin validation passes.

### Added — `hold:` request modifier (deferred/blocked without polluting the lifecycle)

- **A request can now carry an orthogonal `hold:` field** (`deferred`, `blocked`, or any short reason) that pauses it *without* adding a sixth lifecycle status. The five-state chain (`planned → active → review → verified → archived`) stays pure; the hold is a modifier on whatever stage the request is actually in. `spectacular status` surfaces it everywhere — fleet column shows `planned(deferred)`, the card shows `(hold: deferred)`, and `--json` gains a `"hold"` field (the agent contract). `spectacular advance` **refuses** to move a held request until the field is cleared (deleted or set to `none`). Sort rank keys off the raw lifecycle status, so a held request still sorts by its real stage.
- This supersedes the earlier expedient of setting `status: parked` on a request — which duplicated the `parked` value already used by the `ideas`/`feedback` collections and broke the documented five-state invariant. `commit-discipline` was migrated to `status: planned` + `hold: deferred`. Documented in `ARCHITECTURE.md` (PLAN frontmatter schema) and `lifecycle.md`; 5 new assertions in `tests/cli/status.test.sh` (23/23).

### Changed — reconciled builder-agent (b21) ledger drift

- **`builder-agent`'s TASKS/PLAN were reconciled to reality and advanced `active → review`.** M1 (the `spec-builder` agent) and M2 (`build-workflow.md` orchestrator arc) shipped 2026-07-06 and were verified live (built `commit-discipline` M1; a deliberately-vague brief bounced), but the checkboxes were never ticked and the PLAN still named the working title `debug-builder.md`. Fixed: M1/M2 tasks checked, the `debug-builder → spec-builder` rename recorded, frontmatter + roadmap row updated. The only open work is the ≥3-parallel fan-out walkthrough (the `review → verified` item) and the gated M3 trace/CLI signal.

### Added — `doctor specs` frontmatter schema check (spec-audit-mode, b11)

- **`doctor specs` now validates the frontmatter of each flat capability spec** (`specs/*.md`, excluding `index.md`): required keys `status, updated, summary, related` present; `updated` in ISO `YYYY-MM-DD`; `status` in the closed enum `draft | published | deprecated`; and `version:` present **iff** `status: published` (a published spec is a contract and must be versioned; a draft has nothing to version). All findings are `warning`-class, mirroring the `frontmatter` area's severity model. Mechanical only — no semantic matching, no false positives.
- The capability-spec required set **deliberately differs** from the root-anchor set in the `frontmatter` area (`version, updated, summary`): specs make `version` conditional and add `related`. `related:` *resolution* is left to the `links` area — not re-checked here. `index.md` is skipped (catalog doc-class, not a capability).
- **7 new scenarios** in `tests/cli/specs.test.sh` (one per rule, both branches of conditional-`version`, plus an `index.md`-skip regression guard). Full suite: 36/36.
- **Pivot note:** `spec-audit-mode` (b11) originally scoped a semantic coverage audit (orphan bullets/files, stale specs) written against the pre-OKF `specs/<slug>/SPEC.md` layout; that design never settled and its paths went stale after the OKF flattening. Regrilled into this mechanical schema check — the semantic audit was dropped, not deferred.

## [1.29.0] — 2026-07-09

### Added — deterministic `spectacular status` fleet view (status-fleet-view, b23)

- **`spectacular status` renders the active-request fleet from PLAN/TASKS** — an aligned table (slug · status · priority · build · progress · milestone/goal) built directly from each `requests/<slug>/PLAN.md` frontmatter plus grep-safe body signals, sorted active→planned. No doc hand-caches that table anymore. `spectacular status <slug>` prints a single request card (goal · summary · progress · current milestone · `related:` deps · stale flag); `spectacular status --json` emits the machine-readable fleet — the **agent opt-in contract** for discovering requests without opening files.
- **Body signals**: the `## Goal` line, `x/total` top-level task progress (`- [~]` counted separately as deferred, e.g. `5/8 (+1 def)`; indented sub-bullets not counted), and the current milestone (first `### M<N>` with an open top-level task).
- **`tests/cli/status.test.sh`** — covers the fleet table (sort + columns), body-signal edge cases (indented subtasks ignored, `[~]` deferred, milestone advance), `--json` validity, and the unknown-slug error path.

### Changed — enforced PLAN/TASKS structure

- **PLAN sections are now canonical + unnumbered** (`## Goal / ## Constraints / ## Milestones / ## Tasks / ## Dependencies / ## Validation / ## Deliverables`, in order; extra sections allowed between). The PLAN template, `plan-rules.md`, and `tasks-rules.md` were updated; the legacy numbered form (`## 1. Goal`) is now a fixable error.
- **`doctor` (lifecycle area) enforces the structure on active requests** — errors on missing/mis-ordered required PLAN sections, missing `### M` milestone headings, and malformed checkbox states; `.spectacular/archive/` is skipped. `doctor --fix` de-numbers legacy PLAN headings. `- [~]` is now a documented deferred checkbox state.
- **CLAUDE.md Active Requests table retired** → a one-line `spectacular status` pointer (it had already drifted once). `status.md` now consumes `status --json` and layers its judgment/signal detection on top. The repo's own 6 active requests were converted to the canonical schema.
- **Fixed a latent BSD-sed bug** in the milestone-drift check — `\?` in BRE silently failed on macOS; switched the range greps to `sed -nE` (ERE).

## [1.28.1] — 2026-07-09

### Fixed — v0.6 → v2.0 (OKF) migration link-rewriter

- **`migration_apply_v06_to_v20` link rewrite rebuilt** — the old step used over-broad bash `${//}` substring substitution: it injected doubled parens/backslashes into wikilinks and matched inside prose words (`SPEC` → `specs/index` inside `SPECIFICATION`; `debugging` → `debugs/s/ging`), corrupting ~10 live docs. Replaced with one anchored, link-only Python pass that rewrites only inside `[[…]]`/`(…)` targets and YAML `related:`/`depends-on:`/`blocks:` fields — never bare prose — and is idempotent (matches the old target only, so re-runs are no-ops).
- **`related:` rewriting is now depth-aware** — a relocated file's relative targets resolve at its new location (a flattened `specs/<cap>/SPEC.md` → `specs/<cap>.md` no longer overshoots with `../../ARCHITECTURE.md`). Closed all 13 `doctor links` warnings on the migrated workspace.
- **Decision renames use `_resolve_slug_collision`** — two decisions with colliding H1 slugs no longer clobber each other via an unchecked `mv`.
- **Memory `M<N>` numbering is stable across re-runs** — already-prefixed entries keep their number; fresh entries continue past the current max (was: renumbered every run by date sort, silently breaking links).
- **Memory renames are collision-safe** — the two-pass rename now de-dups colliding final bases instead of silently clobbering. Two malformed inputs that normalize to the same `M<N>-<slug>` (e.g. `M1-Foo.md` + `M1-foo.md`) become `M1-foo.md` + `M1-foo-2.md`, both bodies preserved, with a `stderr` warning — restoring the guard the decisions path already had. Uses a bash-3.2-safe membership string (macOS system bash has no associative arrays).
- **H1-slug derivation no longer over-strips** — a hyphenated decision title survives intact (only a leading `D<N> —` label is stripped), producing fuller slugs (e.g. `D6-remove-docs-verbs-deprecation-notice`).
- **`migrations.log` no longer double-logs** — a forced `spectacular migrate --from <ver>` re-run on an already-migrated workspace re-applies idempotently without appending a duplicate history line.

### Changed

- **Skill references + templates synced to the v2.0 (OKF) layout** — `location:` frontmatter and path references across ~20 skill docs/templates now point at the new layout (`SPEC.md` → `specs/index.md`, `ROADMAP.md` → `roadmaps/index.md`, `DECISIONS.md`/`MEMORY.md`/`SESSIONS.md` → their `<dir>/index.md`, `specs/<cap>/SPEC.md` → `specs/<cap>.md`, singular collection dirs → plural), and two broken wikilinks (`[[D7]]`, `[[specs/roadmap/SPEC]]`) were repaired.
- **This repo's `.spectacular/` workspace migrated to v2.0** — plural collection dirs, root index files relocated into their folders as `index.md`, capability specs flattened, decision/memory entries sequentially prefixed. `doctor`: 0 errors, 0 link warnings.

### Added

- **`migrate.test.sh` v0.6 → v2.0 coverage** — a real OKF migration scenario asserting every transform + depth-correct link/`related:` rewrites + prose-safety, plus a byte-identical idempotency scenario (forces a re-run, asserts the tree is unchanged). Previously the migration's behavior was untested (only its registry entry was checked).
- **Prose-safety + collision regression tests** — explicit assertions that the relaxed link regexes leave prose stems untouched (`(the ROADMAP is authoritative)`, `foo(SPEC_VERSION)` survive verbatim), plus a memory-collision scenario proving two entries that normalize to the same slug both survive as distinct files with bodies intact. 75/75 `migrate.test.sh`, 15/15 suite areas.

## [1.28.0] — 2026-07-07

### Added — archive closure gate (b22)

- **`spectacular archive` now runs a closure gate** — three mechanical checks block a drifting archive: (1) `tasks` — every TASKS.md milestone box is `[x]`, or `[~]` with a ` — <reason>`; (2) `verify` — if VERIFY.md exists, VERIFY-LOG.md has a passed (`✅`) walk row; (3) `spec` — a `SPEC-DELTA.md` declaring spec impact exists (or `NONE — <why>`). Each block is overridable **once, explicitly** with `--override <check> --reason "<text>"`, which records `archive_overrides:` (a `{check, reason, date}` list) into the archived PLAN — auditable, never a silent bypass. `--force` is unchanged: it clears the status gate only, never a closure check. Closes the corpus's weakest metric (anti-drift) at the lifecycle tail. Fable review #1.
- **Structured spec deltas** — `spec-sync.md` proposals are now `### ADDED` / `### MODIFIED` (`"<current bullet>" -> "<replacement>"`) / `### REMOVED` blocks written to `SPEC-DELTA.md`, mechanically mergeable and machine-validatable. `doctor specs` gains a **delta-integrity check**: MODIFIED/REMOVED must quote a bullet that exists in the target; ADDED must not duplicate one — the primary drift signal, with the old date heuristic kept as a backstop. Fable review #2 (OpenSpec-inspired).

### Fixed

- **`fm_unset` block-list orphaning (surfaced by undo)** — `spectacular undo` after an overridden archive now drops the whole `archive_overrides:` block cleanly via a new `_fm_unset_block` helper; previously the list items were re-parented under the preceding `related:` list (valid YAML, wrong data).

### Changed

- **SPEC.md synced** — the archive-closure-gate request archived through its own gate; its `SPEC-DELTA.md` merged the closure-gate + delta-integrity behavior into `SPEC.md` (Lifecycle + Substrate doctor bullets), the first end-to-end use of the new delta-sync flow.

## [1.27.0] — 2026-07-07

### Fixed (fable review W1 — guidance contradictions)

- **`target_version:` fully retired** — roadmap-rules Slot 6 + full-tier unlock now derive from the ledger (`build → target-version`), the roadmap template comments match, and the broken `idea promote --target-version` flag (passed a flag `spectacular new` never accepted) is removed. The plan-rules ban was already correct; everything now agrees with it.
- **new-request.md de-forked** — its embedded PLAN/TASKS template blocks (which had drifted from `templates/*/base.md` and even failed tasks-rules' own heading check) are replaced with pointers to the canonical templates.
- **review.md aligned with per-doc schemas** — base check 3 defers to each doc's rules file (PLAN requires `status`/`updated`/`summary`, not `priority`/`owner`); the ambiguous "check 4+" wording clarified (base checks 1–3 immune, override checks numbered 4+).
- **TASKS `status:` mirror acknowledged** — SKILL.md/CLAUDE.md no longer claim lifecycle state is "never duplicated"; TASKS mirrors PLAN for skim tooling, PLAN is authoritative, doctor repairs drift.
- **CLI fallback scaffold** emits the full PLAN shape (placeholders + `## Understanding` + `## Decisions`) — a template-less workspace no longer produces PLANs the placeholder check can't catch.

### Changed (fable review W2 — sharpened spec/plan/strategy rules)

- **PLAN template gains `## Decisions`** (chose X over Y — because Z; rejected alternatives stay) — the destination decisions-rules always pointed at now exists; plan review check 10 verifies entries name an alternative.
- **Plan-time falsifiability** — PLAN § Validation checks must state an authority (run/assert/judge/observable); aspiration verbs fail; mini-refine pattern + strengthened check 7 (plan-rules, template comment).
- **PLAN Goal → PRD traceability** — new check 11 + mini-refine: a Goal that only re-words the request summary fails.
- **Supersession convention** (active-request.md) — sanctioned `## SUPERSEDED <date>` block; disproven content is never deleted; formalizes the pattern healthy projects invented ad hoc.
- **Autopilot passes the gate** (new-request.md) — a skill-drafted PLAN runs `plan review` and shows the punch list before confirmation.
- **Lifecycle vocabulary lock** (lifecycle.md) — `status:` is vocab-only; intermediate intent goes in `note:`; kills invented states like `fixed-pending-verify`.
- **Evidence-before-questions** (grill.md pre-flight) — code-touching grills cite `path:line` in the first technical question; never ask what the code answers.
- **Bug workflow: disproof ledger + 3-strikes** — investigator returns list ruled-out hypotheses with killing evidence (`ruled_out` field in debug-trace schema, copied to audit entries); after 3 failed fixes, question the architecture with the human.
- **POLICY armor** (template + this repo's POLICY.md v1.4) — the 4 block policies and 4 @Debugging policies each gain a one-line Law, an Excuse/Reality table, and red-flag self-checks; severities unchanged.
- **TASKS template acceptance stubs** (`→ check:` per milestone) and **brief placeholder ban** (build-workflow: "TBD"/"appropriate error handling"/"similar to M<N>" are brief failures).
- **SPEC index check** (spec-rules) — a Capabilities bullet >2 lines or with sub-bullets fails review; spec-sync bullets prefer observable SHALL-strength phrasing, scenarios live in capability specs.
- **Verify: coherence pass + retrospective default-on** — the walk confirms PLAN Decisions actually shipped (advisory); the retrospective question is asked once instead of "optionally".

### Added

- **`archive-closure-gate` request scaffolded (b22, planned)** — archive blocks (with recorded override) on unticked TASKS / unwalked VERIFY / missing spec delta; spec-sync becomes ADDED/MODIFIED/REMOVED deltas merged mechanically. The fable review's highest-leverage change, routed through the request lifecycle.
- **Fable review report** at `docs/reviews/fable-spec-quality-review.md` — corpus audit (18 artifacts, 5 workspaces), guidance audit, comparative research (superpowers/gstack/OpenSpec), ranked 15-item change list.
- **A/B validation of the review changes** — 5 synthetic benchmarks (plan-drafting, flawed-plan review, debug orchestration, plan supersession, verify coherence) run blind against main: branch 17/18 traps vs main 9/18 by-the-book. Report §7; re-runnable scenarios in `docs/reviews/ab-scenarios/`.

## [1.26.3] — 2026-07-06

### Added

- **Workspace provenance tracking** — `config.yaml` now carries `created_with` (the spectacular version that scaffolded the workspace, write-once) and `last_touched_with` (the version that last structurally touched it, bumped by `migrate` and `doctor --fix`), alongside the existing `workspace_schema`. Answers "what version created this workspace, and what version last edited it?" — previously unanswerable. Workspaces predating this get `created_with: "unknown"` backfilled by `doctor --fix` (never a false claim of the current version). See `docs/configuration.md` § workspace_schema + provenance.
- **`doctor workspace` schema-drift check** — emits a warning when `workspace_schema` is behind the CLI's expected version (→ run `migrate`) or ahead of it (→ update the CLI). The logic existed only behind `status --against-latest`; now a routine `doctor` run surfaces it. Covered by `scenario_21_schema_behind_warning`.
- **Migration history log** — `migrate` appends to `.spectacular/migrations.log` (`<date>  <from> → <to>  (spectacular <version>)`), making upgrade provenance auditable rather than merely inferable from the current schema.

### Fixed

- **Stale CLI version string** — `SPECTACULAR_VERSION` was hardcoded to `1.24.0` while the plugin manifests shipped `1.26.2`; bumped to `1.26.3` so `spectacular --version` and the provenance stamps report the truth.

## [1.26.2] — 2026-07-06

### Added

- **`doctor lifecycle` milestone-label alignment check** — flags `M<N>` label drift between a request's `TASKS.md`, `PLAN.md` §3 Milestones, and §6 Validation (advisory `judgment` warning, never blocks). Also flags a non-standard milestone prefix (e.g. `G1` instead of `M1`), but falls back to matching by milestone **name** before declaring a real chain-break — so a relettered-but-same-named milestone doesn't false-positive. Closes a gap surfaced while scoping a "Builder agent" idea (`.spectacular/ideas/coding-agents.md`): the task-row → milestone-block → plan-section chain is only reliably walkable when IDs/names actually agree, and nothing checked that before.
- **ID-namespace convention documented** — `ARCHITECTURE.md` now has a table of Spectacular's existing single-letter + number ID families (`M<N>` milestones, `D<N>` decisions, `F<N>` fixes, `b<N>` roadmap builds, `A<N>` debug findings), written down for the first time instead of living as implicit tribal convention.
- **`codex-agent` second opinion in `bug-workflow.md`** — an optional read-only cross-check during Step 2b, using a different reasoning model than the `debug-investigator`. Gated by a 3-trigger table (cross-cutting/high blast-radius fix, already looped once with `needs-more-context`, or a low-confidence root cause) so it's reached for deliberately, not reflexively — the routine `root-cause-found` case skips it.

### Fixed

- **`spec-audit-mode` request's `TASKS.md`** — was still the unfilled scaffold template (3 generic milestones vs `PLAN.md`'s 4 real ones); filled in with the actual milestone checklist.

## [1.26.1] — 2026-07-06

### Fixed

- **Project docs synced to v1.26.0** — `CLAUDE.md`, `AGENTS.md`, and `README.md` still described pre-debug-fleet state: stale `(Unreleased)` tags on `bug-workflow.md`/`debug-trace.md` references, `AGENTS.md`'s Testing section claiming "no formal test suite exists" (false since `tests/{cli,pipeline,agents}` shipped), and README's `doctor <area>` list + workspace tree missing `debug`/`vision`/`decisions`/`roadmap` and the `debug/`/`audit/`/`fixes/` collections.
- **`.spectacular/SPEC.md` synced to v1.26.0** — the capabilities index was 9 releases stale (claimed "as of v1.17.x") because v1.25.0 and v1.26.0 shipped via direct commits rather than the request lifecycle, so the doctor's spec-drift check (compares `SPEC.md updated:` against the newest *archived request*) had nothing to trigger on. Added the missing `audit/`+`fixes/` soft-DB collections, debug agent fleet, and `doctor debug` bullets; synced the doctor-area list.

## [1.26.0] — 2026-07-05

### Added

- **Debug agent fleet (`.claude/agents/`)** — three read-only-or-closed-contract subagents the orchestrator delegates to during a bug flow: **`debug-investigator`** (discovers *where* + *why* on an open bug, returns ranked findings + plausible-solution space, never prescribes the literal edit), **`debug-fixer`** (applies one *closed* five-slot brief under an apply-only contract — smallest faithful diff, local style, operation-care gradient add<edit≈patch<delete, risk-scaled verify — and bounces the moment execution turns to judgment), and **`debug-researcher`** (searches forums/docs/issues for known-external bugs, returns a cited verdict). Each writes only its own trace artifact; none writes the ledger.
- **Debug trace schema (`references/debug-trace.md`)** — one folder per live job under `.spectacular/debug/<job-slug>/`, one JSON artifact per agent turn (`job.json` spine + `investigation.json` / `research/` / `fixes/` leaves + `outcome.json`). `job.json` carries `symptom_class` (test_failure · runtime_error · wrong_behavior · build_error · performance · unknown) and the persisted investigation `brief`; `fix-NN.json` carries `changed[]`, `test`, and `risk`. debug/ = the raw pipeline (kept as trace, never pruned); audit/ + fixes/ = the two distilled summaries earned at resolution.
- **Principle 11 — "Earn each step: no rockets without the launchpad"** (`PRINCIPLES.md`) — the sequence complement to principle 10: build in the *right order*, never pour effort into an impressive far step while the near step it depends on is still missing.
- **`@Debugging` policy set (`POLICY.md`)** — `check-prior-fixes`, `ceremony-matches-uncertainty`, `fix-root-not-symptom`, `log-only-verified-reusable`, and `use-audit-fix-verbs` (write audit/fix entries via their verbs, never by hand — `prefer-cli-mutator` applied to the debugging phase). Plus new `@Implementation` / `@Planning` policies backing principle 11: `milestones-in-build-order`, `build-order`, `earn-the-verification`, `prefer-cli-mutator`.
- **Pipeline test suite (`tests/pipeline/`)** — seven integration runbooks with real bug fixtures exercising the orchestrator's choreography end-to-end: resolve→ledger graduation, fix-needs-a-request routing, Researcher live run, concurrent-Fixer disjoint trace writes, the full pipeline capstone, the **won't-fix disposition** (P6 — a real one-line-fixable bug the orchestrator deliberately declines: deprecated path, frozen consumer), and the **inbound safety valve** (P7 — an honest under-determined Investigator return that the orchestrator's well-formedness + symmetry backstop must catch without fabricating a fix). Plus `tests/agents/` judgment fixtures for the fixer/investigator (incl. an unsafe-delete bounce case).
- **`spectacular fix new --debug-job <slug>`** — writes a `debug_job:` back-link into `F<N>` frontmatter, ending the orchestrator's manual hand-stamp of the trace link. Renders `null` when omitted; `--into-fix` defaults it cleanly.
- **`spectacular doctor debug`** — validates hand-written `debug/<slug>` trace spines against the schema enums (`job.json` `status`, `outcome.json` `disposition`) and the invariant that a `wont-fix`/`folded-into-request` job logs no `F<N>`. Guards the closed-enum drift an LLM orchestrator makes writing these JSON files by hand (e.g. leaking a `reason` value like `needs-more-context` into the `status` slot) — catches it at check time instead of a resume failure.

### Changed

- **Bug workflow (`references/bug-workflow.md`)** — rewritten as the orchestrator's full arc: a top-of-doc map (Steps 0→3 with the new **Step 1c "open the job"** and a `status`→step resume crosswalk), a Step 1b fan-out decision table (≥3 independent closed disjoint-file fixes → fan out; else self-serve), same-file serialization rules (serialize inline, never parallel, no git branches), the investigation-brief quality bar, the block↔findings symmetry check, malformed-return backstops, and the channel mechanic (returned block = Agent-tool result; JSON = durable copy). `needs-reproduction` routes to the orchestrator inline / the user — no Reproducer agent (deferred to the Heisenbug case). Step 2b now has **two disposition forks** for findings that don't become a fix: "can the findings even close?" (→ `folded-into-request`) and "should the fix even be applied?" (→ `wont-fix`) — the two ways a debug job closes without a fix landing, neither logging an `F<N>`.
- **CLI policy-hook validation (`cli/spectacular`)** — the nine work-phase hooks (`@Init @Planning @Implementation @Debugging @Verification @Archive @Remember @Snapshot @SessionEnd`) are now the validated set for custom policies.

## [1.25.0] — 2026-07-04

### Added

- **`audit/` + `fixes/` soft-DB collections** — two new bug-lifecycle collections (index-only, auto-numbered `A<N>` / `F<N>`). `audit/` is the diagnosis scratchpad *before* a fix is planned; `fixes/` is the verified-fix log written only once a bug is resolved **and** verified. Both ride the existing collection machinery (templates + `<id>-rules.md` + `_iter_md`).
- **`spectacular audit new|list|resolve`** — scaffold, list, and close bug investigations. `audit resolve <A> --into-fix` **graduates** an audit into a fix entry, copying every matching slot forward (Problem, Intended behavior, Root cause, Proposed fix, Success criteria) and setting `from_audit: A<N>`.
- **`spectacular fix new|list`** — log a verified fix. Flags: `--problem/--intended/--cause/--fix/--criteria`, `--verified-by`, `--signature`, `--from-audit` (validated). Omitting `--verified-by` warns and marks the entry `verified: null` (a soft gate — a draft, not a trusted fix).
- **Bug-fixing skeleton in both entry schemas** — every audit/fix entry now carries `problem → intended behavior → root cause → fix → success criteria`. `fixes/` adds **Verified by** (the evidence, distinct from the success-criteria bar) and a searchable **Signature** field.
- **Self-learning loop (`references/bug-workflow.md`)** — the skill checks `.spectacular/fixes/` signatures *before* diagnosing a new bug ("have we fixed this before?"), applies a lightweight audit-first-vs-just-fix heuristic (no ceremony on one-liners), and logs a signed fix only when it carries reusable knowledge. Per-project today; corpus designed to later pool/export across projects.
- **`references/soft-db-index.md`** — canonical routing index for all 7 soft-DB collections (memory · decisions · sessions · ideas · feedback · audit · fixes): role, purpose, structure, write verb, and the boundary rule that prevents mis-routing. Clarifies that `requests/` and canonical docs are *not* collections.

### Changed

- **SKILL.md** — added routing for bug reports (→ `bug-workflow.md`), the audit/fix verbs, and a "Where does this belong?" soft-DB routing section (→ `soft-db-index.md`). The `description` now lists all seven collections (trimmed to 949 chars, under the doctor 1000-char warn band). Ponytail trim carried over from the same session: −20 lines (dead PRD-legacy table, inline version stamps).
- **`decide` flat-mode fixes** — `summary` now counts decisions in flat mode (was reporting 0 when ADRs are prose blocks in `DECISIONS.md`); flat-mode `decide` returns exit 0 on success (was returning 1 via a trailing `[[ ]] &&` short-circuit despite persisting). Logged as `fixes/F1`–`F3`; guarded by `tests/cli/decide.test.sh`.

## [1.24.0] — 2026-06-30

### Added

- **Snapshot version coupling** (snapshot-retention b16) — `spectacular snapshot` now names the snapshot for the version the copied content **is** (a doc at `version: 1.3` → `_snapshots/PRD/@v1.3.md`), *then* bumps the live doc to `1.4`. The `@v` label and the `version:` field can no longer drift. Docs without a `version:` field (e.g. the newly snapshot-able `DESIGN.md`) use a plain `@v<N>` counter and are not version-bumped.
- **Tiered snapshot retention + `spectacular snapshot prune`** — generational retention keeps the union of three tiers: origin (`@v1`), periodic (newest per `month`/`week` bucket, keyed off `updated:` frontmatter dates), and recent (newest `keep`, default 3). `snapshot prune` removes the rest — `git rm` if tracked (history holds it), else moved to `.spectacular/.trash/`; dry-run by default, `--apply` to perform. Bounds a doc to ≈ `1 + periods_alive + keep` snapshots instead of unbounded growth.
- **`snapshots:` config block** — `folder` (store dir, default `_snapshots`), `keep` (default 3), `period` (`month`|`week`|`off`, default `month`), `gitignore` (default `false`). All optional with sane defaults. See `docs/configuration.md`.
- **`DESIGN.md` is snapshot-able** — added to the canonical allowlist (the allowlist stays closed otherwise).
- **`doctor snapshots` retention/migration/gitignore checks** — flags prunable accumulation (info), the `snapshots/` → `_snapshots/` folder rename (warning + `--fix`), and `.gitignore` drift vs `snapshots.gitignore` (warning + `--fix`). Gap detection now skips dirs that mix counter and version names (a b16 transition guard — no false positives).

### Changed

- **Snapshot store renamed `snapshots/` → `_snapshots/`** (configurable via `snapshots.folder`). The `_` prefix marks it a non-content layer, consistent with `_archive/`. `doctor --fix snapshots` migrates an existing `snapshots/` dir losslessly (git-mv when tracked). This repo was dogfood-migrated (18 snapshots). The 15 hardcoded `snapshots/` paths in the CLI collapsed to a single config-resolved `$snap_root`.

## [1.23.3] — 2026-06-30

### Changed

- **Closed `skill-desc-length-check` (b10 — verified + archived).** The doctor sub-check (`check_skill_desc_len`), shared `scripts/check-skill-desc.sh` helper, and the `pre-commit-wrapper` guard were all verified live (doctor.test scenario 17: 53/53 assertions; guard fires on every commit). The awk parser is intentionally duplicated between the CLI binary and the helper — `cli/install.sh` ships only the binary, so the installed doctor cannot source `scripts/`; the two copies are kept byte-identical in lock-step.
- **ROADMAP/CLAUDE.md drift reconciled.** Ledger slots `b4 → v1.23.2` and `b10 → v1.23.3` (both shipped). `CLAUDE.md`'s active-requests table no longer lists shipped/archived work (`roadmap-ledger`, `decisions-index`, `cli-debt-removal`, `skill-desc-length-check`, `roadmap-contract-docs`, `roadmap-pruning` removed); the remaining active set is ordered by priority.
- **`ROADMAP.md` pruned via `roadmap migrate`** — the v1.20.0 shipped block moved to `.spectacular/roadmap/v1.20.0.md` behind the `## Shipped` index (keep newest 3 inline); `doctor roadmap` now fully green.

## [1.23.2] — 2026-06-30

### Removed

- **Dead `skills/spectacular/templates/docs/` directory** (`docs.yaml.tmpl`, `index.md.tmpl`, `page.md.tmpl`). Residue the `cli-debt-removal` cleanup missed: the templates scaffolded the removed `docs export` machinery, `docs.yaml.tmpl` still pointed at the deleted `references/docs-renderer-adapters.md` and documented the removed `docs export <renderer>` verb, and nothing in the CLI loaded any of them (`doctor docs` directs users to `pageworks init` to scaffold `docs.yaml`). Closes the `cli-debt-removal` request (b4 — verified + archived).

## [1.23.1] — 2026-06-30

### Changed

- **Contract-prep ladder un-pinned to `target: tbd`** in `ROADMAP.md` — the v2.0.0 ledger entries (①→②→③) no longer carry concrete version numbers, so reslotting a near-term build never forces a cascade-renumber of the runway.
- **Roadmap renumber anti-pattern documented** in `roadmap-rules.md` — keep unstarted runway `tbd`, treat a reslot as a one-cell edit (not a renumber), and use `## Label *(target: tbd)*` headers until a version is actually pinned.

## [1.23.0] — 2026-06-29

### Added

- **`spectacular roadmap migrate [--dry-run] [--keep N]`** — index-mode shipped-history scaling for `ROADMAP.md` (roadmap-pruning b18). Moves shipped per-version prose blocks into per-version files (`.spectacular/roadmap/v<X.Y.Z>.md`) behind a `## Shipped` index, keeping the most-recent **N (default 3)** shipped blocks inline. Only blocks whose own `**Status:**` is `shipped` move; planned/active/vision blocks stay. Snapshot-safe (writes per-version files before rewriting ROADMAP.md — no data loss on a partial run), idempotent, dry-run by default. Mirrors the decisions-index pattern. Bounds ROADMAP.md's agent-context cost as history grows (this repo: 528 → 410 lines).
- **`doctor roadmap` area** — index-mode integrity: orphan `## Shipped` index lines (no matching file), stale per-version files (no index line), and an info nudge when shipped blocks beyond the keep-window are still inline (prunable via `roadmap migrate`). Flat mode emits only the nudge.
- **Roadmap ledger documentation** (roadmap-contract-docs b17) — the build-id → version model is now specified, not just architecture-noted: `specs/roadmap/SPEC.md` gained a ledger section (build ids, `target-version` single-source, `tbd` sentinel, ledger-status-vs-request-lifecycle); `docs/versioning.md` gained "The roadmap ledger" walkthrough; `docs/configuration.md` documents `last_build:`; `docs/commands.md` gained a `spectacular roadmap` section.
- **ADR discoverability** — `decisions-rules.md` now carries a "store-worthy decision?" routing table and an explicit "ADRs live in DECISIONS.md — don't create `docs/adr/`" callout. `doc-index.md` and `SKILL.md` triggers now grep-match "ADR" / "architecture decision" → `spectacular decide`.

### Changed

- **`target-version: tbd`** is now a documented ledger sentinel ("slotted but not version-pinned yet"), distinct from a `<TBD>` placeholder. `roadmap-rules.md`'s placeholder check is scoped to prose slots so ledger `tbd` is no longer falsely flagged.
- **ROADMAP.md** dogfooded the new index mode: 7 oldest shipped blocks moved to `.spectacular/roadmap/`, the "Recently shipped" CHANGELOG-mirror section removed, stale reconciliation notes pruned, `roadmap-overrides` → `roadmap-rules` references fixed.

## [1.22.0] — 2026-06-28

### Added

- **`spectacular undo`** — a reverse gear for lifecycle mutations (lifecycle-undo b12 → v1.22.0). Reverses the most recent `advance`, `archive`, or `idea promote`:
  - **advance** → status back one step on PLAN + TASKS.
  - **archive** → moves the dir back to `requests/`, restores the pre-archive status, drops the `archived:` field, and reverses the inbound `../../archive/<slug>/` link rewrites across sibling requests (git-aware, plain-`mv` fallback).
  - **idea promote** → restores the idea source to `ideas/`, resets its status, drops `promoted_to:`; the scaffolded request dir is **left in place** unless the user confirms removal (decision D9).
  - **Single-level** (one `.last-mutation` breadcrumb, gitignored). Refuses on a **stale breadcrumb** (any affected file modified after the recorded mutation — timestamp-vs-mtime guard). `--dry-run` previews without mutating. "Nothing to undo" exits 0.
  - Each mutator (`cmd_promote`, `cmd_archive`, `cmd_idea_promote`) writes the breadcrumb; `cmd_undo` reads, reverses, and clears it. New `fm_unset` frontmatter helper. Skill tier-reveal hints added to `lifecycle.md` + `archive.md`. `tests/cli/undo.test.sh` (30 assertions).

## [1.21.0] — 2026-06-28

### Changed

- **`onboarding.md` deduped against `status.md`** (onboarding-dedup b14 → v1.21.0). Onboarding no longer restates the ~95%-shared read+briefing sequence — it now says "run the status.md flow, with these deltas" and keeps only what's onboarding-specific (always-run substrate check, takeover tone, first-look observations table, pre-split detection, example briefing). `status.md` is the single owner of the read sequence; one source of truth, no independent drift.

### Added

- **Guided first-run** (`references/guided-first-run.md`) — when `/spectacular` hits a fresh/empty workspace (init ran, no requests), the skill ushers the user new → optional PRD grill → first request → `spectacular next`, **one step at a time, never dumping the verb surface**, instead of printing an empty briefing. Routing wired into both `status.md` and `onboarding.md`; the empty-vs-existing distinction is explicit (onboarding = existing project with prior work; guided first-run = blank slate). Skill-driven — no new CLI flag (an optional `init --walk` is left for later).

## [1.20.0] — 2026-06-28

### Changed

- **Skill-reference doc sprawl reduced** (rules-files-audit b13 → v1.20.0):
  - The 5 boilerplate-only `<doc>-rules.md` stub bodies (architecture, principles, stack, agents — and the verb list in spec) now point to a single **"Stub default behavior"** section in `doc-index.md` instead of each restating the grill/refine/review default. Frontmatter (the engine's dispatch) is untouched. `agents` keeps its top-level-`AGENTS.md` note; `spec` keeps its index role + archive-sync `review` override; `tasks` is unchanged (it carries a real body, mislabeled `mode: stub`). Decision recorded as D8.
  - **Verify-doc trio merged into one `verify.md`** — the former `verification.md` (2-of-6 rule) becomes Part 2 and `verify-tests.md` (promoting checks to scripts) becomes Part 3 of the existing walk doc (Part 1). All inbound `[[verification]]`/`[[verify-tests]]` wikilinks, path references, and SKILL.md routing updated; the two source files removed. Three files for one concept → one reference with three labelled parts.

## [1.19.0] — 2026-06-28

### Added

- **`spectacular advance`** — lifecycle move-forward verb (renamed from `promote`; `promote` stays as a deprecated alias that prints a one-line notice). Frees `promote` to read unambiguously as `idea promote`. (naming-coherence b15 → v1.19.0)
- **`spectacular next`** — read-only verb that prints the single highest-priority next action (active → review → planned → empty-workspace usher). Mutates nothing.
- **Tier-reveal suggestions** in skill flow docs — one-line "next step" hints after scaffolding (`new`) and at lifecycle checkpoints, never mid-flow.
- `doc-id-aliases:` support sketch on `pack-rules.md` so the renamed `pack` doc-id keeps `convention-pack` as a back-compat alias.

### Changed

- **`feedback`** is now the canonical verb (was `feedback-loop`); `feedback-loop` joins the hidden aliases. `feedback-rules.md` doc name unchanged.
- **Pack doc-id renamed `convention-pack` → `pack`** to match the `pack` CLI verb and `pack-rules.md`. The old id is still accepted. Updated doc-index.md, pack-rules.md, SKILL.md, grill.md.
- `spectacular pack new|grill|refine|review` now redirect to the skill (documented but previously died as "unknown pack verb").
- Lifecycle/verify/troubleshooting/scaffold/commands docs updated to teach `advance` and the correct `idea promote` form.

### Fixed

- Latent bug: backtick command-substitution in the feedback usage heredoc — bare `spectacular feedback` actually executed `spectacular remember` while printing help. Heredoc switched to a quoted delimiter (`<<'EOF'`).

## [1.18.1] — 2026-06-27

### Fixed

- `is_canonical_doc` had an unreachable duplicate `SPEC.md` case arm, so per-capability `specs/<cap>/SPEC.md` files were never recognized as canonical and **could not be snapshotted**. The `SPEC.md` arm now matches both the top-level doc and capability specs.
- Dead, misleadingly-commented `exit 0` after the `doctor` dispatch replaced with `exit $?` so a real exit code propagates if `doctor` ever returns.

### Changed

- Internal cleanup (no behavior change): collapsed four byte-identical collection walkers (`_idea`/`_decision`/`_memory`/`_session_iter_all`) into one `_iter_md <subdir>`, and merged `kit_triggers_always`/`kit_triggers_suggested` into a single `kit_triggers <kit> <which>`. Net −48 lines in `cli/spectacular`.

## [1.18.0] — 2026-06-26

### Added

- **SPEC.md drift check** in `spectacular doctor specs` — warns (`⚠️`) when `SPEC.md`'s `updated` date predates the newest archived request, signalling a likely missed spec-sync. Surfaced in `/spectacular status` and routed to the skill's spec-sync flow for content reconciliation. Date heuristic ("may be stale"), not a content diff.
- New planned request `spec-audit-mode` — content-aware spec audit (orphan capability bullets, orphan spec files, stale per-capability specs) building on the drift heuristic.

### Fixed

- `spectacular summary` crashed with `_active_plan_slugs[@]: unbound variable` when no request had `status: active` (empty array under `set -u`). Both dependency-graph loops now guard the empty array.

### Docs

- `docs/commands.md` doctor area list refreshed (was missing `docs`, `personas`, `memory`, `sessions`, `feedback`, `ideas`, `policies`) and documents the new `specs` drift check.
- `SPEC.md` synced: added the previously-unrecorded `policy-engine` capability bullet (shipped v1.12.0) and the `policies` doctor area; cleared the live drift the new check flagged.

## [1.17.2] — 2026-06-20

### Fixed

- Skill `description` trimmed from 1146 → 986 chars so it loads under Codex's 1024-char limit (Claude Code's 1536 limit had masked the issue). No triggers or doc names removed — only wording tightened.

## [1.17.1] — 2026-06-20

### Fixed

- `DECISIONS.md` missing `version` and `summary` frontmatter fields (doctor frontmatter warnings)
- Stale `related:` link in `ideas/memory-protocols.md` pointing to removed `requests/soft-db-substrate/PLAN.md` (doctor links warning)

## [1.17.0] — 2026-06-16

### Added

- **Roadmap ledger — single source of truth for `build → version` mapping.** Every request gets a stable build id (`b<N>`) and one row in the `## Roadmap ledger` table in `ROADMAP.md`. The `target-version` column is the only place a version number is written; everything else uses slug or build id. `spectacular new` stamps `build: bN` on new PLAN frontmatter and increments `last_build:` in `config.yaml`.
- **`spectacular roadmap` reads from the ledger (v1.17.0+).** The render verb now parses the ledger table as its data source instead of prose version blocks — a one-row edit to reslot a request is reflected immediately in `spectacular roadmap` output. JSON output gains `build` and `slug` fields.
- **`DECISIONS.md` index mode.** When a project's decisions file outgrows flat-file scale (~50+ entries), agents can split it into a cheap one-liner index + per-entry files in `decisions/`. Detected by presence of `decisions/` subfolder — flat mode remains fully valid and backwards-compatible.
- **`spectacular decisions migrate [--dry-run]`.** One-shot verb: reads flat `DECISIONS.md`, extracts each `## YYYY-MM-DD —` block into `decisions/D<N>.md`, rewrites root as one-liner index. `--dry-run` previews without writing. Idempotent if already migrated.
- **`spectacular decide` writes index mode automatically.** When `decisions/` folder exists, `decide` writes full ADR prose to `decisions/D<N>.md` and appends one index line to `DECISIONS.md`. Flat-mode behavior unchanged. Auto-numbers from the highest D<N> in the folder.
- **`doctor decisions` area.** Index-mode consistency checks: mode consistency (no prose in index), orphan index lines, stale per-entry files, sequential D-numbering (gaps → warning, duplicates → error).
- **`decisions-rules.md` updated.** `mode: index | flat`, canonical index line format, per-entry file format, detection rule, and agent read pattern all documented.

### Removed

- **`spectacular docs init|export|new|review|status` verbs.** Removed after being deprecated with in-product banners since v1.2.0 (4+ releases). Public docs work lives in the [`pageworks`](https://github.com/alexsmedile/pageworks) skill; `doctor docs` (discovery-only) remains.
- **`deprecation_notice()` banner machinery.** No longer needed once the deprecated verbs are gone.
- **`--global` init flag.** Deprecated alias for `--skill-scope global`; use `--skill-scope global` directly.
- **`docs-contract.md`, `docs-rules.md`, `docs-renderer-adapters.md` reference docs.** Canonical versions live at `pageworks/references/`.
- **Deprecated `docs-manifest` and `docs-page` entries from `doc-index.md`.**

> **MINOR classification rationale (D6):** Removal of banner-warned, continuously-telegraphed deprecated surface is treated as MINOR per the project's versioning convention. No undeprecated behavior changes. `pageworks` is the documented, available replacement.

## [1.16.0] — 2026-06-08

### Added

- **Cross-request link schema (`depends-on:` / `blocks:` in PLAN frontmatter).** Two additive sibling fields to `related:` — `depends-on: [slug]` (A cannot ship before B) and `blocks: [slug]` (A must ship before B can proceed). Documented in ARCHITECTURE.md alongside the inverse-label table and computed-not-stored rule. Advisory only — no locking, no auto-blocking.
- **Inverse-link resolver.** At read time, the CLI computes the bidirectional graph from all forward declarations: `blocks: [B]` on A surfaces as `blocked-by: A` on B; `depends-on: [B]` surfaces as `required-by: A` on B. Inverses are never written to disk — single source of truth stays the declaring request.
- **`spectacular links [<slug>] [--json] [--all]`** — new read verb. Shows the whole-graph dump (default: only requests with edges; `--all` includes unlinked). Per-request view with `<slug>`. `--json` emits `{graph: [{slug, depends_on, blocks, related, required_by, blocked_by}]}`.
- **`spectacular request <slug>` gains a Links section.** When a request has any declared or computed edges, they appear below the progress bars in the detail view.
- **`spectacular summary` link advisory.** When active requests have declared edges, a compact "Active links:" section surfaces ordering dependencies at a glance (advisory, non-blocking).
- **`spectacular new` relationship prompt.** After scaffolding a new request, if existing active/planned requests share keyword overlap with the new slug, a hint to declare `depends-on:`/`blocks:`/`related:` is printed.
- **`doctor links` root-aware path resolution.** `related:` targets that are bare root-doc filenames (`PRD.md`, `ARCHITECTURE.md`, `ROADMAP.md`, etc.) now resolve against `.spectacular/` rather than the declaring file's own directory — eliminating 7 false "not found" warnings for canonical doc references.
- **`doctor links` validates `depends-on:` and `blocks:`.** Slug targets are checked against `requests/` and `archive/`; archived = satisfied (shows `✓ (shipped)`); unknown slug = warning. All three link fields validated in one pass.
- **`doctor memory` staleness flag.** Memory entries older than 180 days trigger a warning to review and prune — conservative nudge, not a nag. Gradient: sessions 4h < feedback 30d < ideas 90d < memory 180d.
- **Example link graphs in `tests/cli/links.test.sh`.** Two example scenarios: (A) `depends-on` + inverse `required-by`; (B) `blocks` + `blocked-by` + archived dep resolved as satisfied + dangling slug flagged by `doctor links`.

## [1.15.0] — 2026-06-07

### Added

- **Visual layer — ASCII rendering (`_ascii_bar`, `_ascii_box`, `_ascii_color_enabled`).** A shared rendering helper layer in the CLI: `_ascii_bar <done> <total> [<width>]` fills `█░` (TTY/color) or `#.` (plain) with a percentage; `_ascii_box <title> [lines...]` draws a left-border box. Both degrade cleanly via `NO_COLOR=1` or non-TTY stdout (no escape codes, no block characters, no color). One helper reused by every visual surface — no per-command bespoke rendering.
- **Visual `progress <slug>` render.** Milestone bars with percentage and done/total count; a roll-up overall bar at the bottom. Completed milestones show `✓`; `--format json` output is byte-identical to v1.14.x.
- **Visual `summary` dashboard.** Request-state counts (planned/active/review/verified) rendered as proportional mini bars (width 10). Zero-count states are omitted. `--json` unchanged.
- **`spectacular roadmap` CLI verb.** Renders `ROADMAP.md` as a version arc grouped by tier (`full` → Runway · `themed` → Major · `vision` → Vision). Status indicators: `✓` shipped · `▶` active · `·` planned. Shipped versions hidden by default; `--all` includes them. `--json` emits an array of `{version, title, tier, status}` objects.
- **ASCII app-UI mockup block format.** Documented convention for dropping a renderable mockup into a request `PLAN.md` or `SPEC.md`: fenced code block with language tag `mockup`, ≤ 64-char lines, `[square brackets]` for actions, `[____]` for input fields. Used by the skill during `imagine` to propose UI artifacts for human approval. Full spec: `docs/visual-conventions.md`.
- **`docs/visual-conventions.md`.** Public-facing doc covering: bar fill conventions (block vs plain), summary dashboard layout, roadmap arc tier legend, and the mockup block format with a real example. Registered in `docs/docs.yaml`.
- **`_progress_text` / `_progress_json` milestone-header fix.** Both helpers now match `### M` (H3) in addition to `## M` (H2) — all real TASKS.md files use H3 milestone headers; the old `## M` pattern silently produced empty output for every request.

### Imagine mode (v1.15.0 co-ship)

- **`/spectacular imagine <slug>` — imagination-backed planning.** Generative-first mode: renders see-able ASCII artifacts (user stories, UI mockups, architecture sketches) the human reacts to per-fragment, then derives a draft PLAN from the approved vision. Expands Spectacular's thesis from spec-driven to *spec-driven AND imagination-backed*. Full engine: `references/imagine.md`.
- **`vision/` soft-folder substrate.** `requests/<slug>/vision/` holds a `VISION.md` spine + typed subfolders (`stories/`, `ui/`, `arch/`). `spectacular imagine <slug>` scaffolds it; `spectacular vision add <kind> <name>` is the mechanical fragment mutator. Manifest regenerates from fragment files.
- **`doctor vision` area.** Fragment frontmatter check + kind/subfolder match + manifest drift (with `--fix`) + dangling persona refs + approval progress.
- **`references/imagine.md`.** Full render→react→derive loop spec: generative rendering (step 1), per-fragment approve/redirect/reject (step 2), approved vision → draft PLAN derivation (step 3). Draft never auto-accepted — hands off to PLAN grill/review.
- **`references/vision-rules.md`.** Doc-type rules for `vision`: frontmatter schema, fragment kinds, `imagine` dispatch mode, spine/subfolder structure.

## [1.12.2] — 2026-05-31

### Added
- **`scope-down` + Principle 10 now ship in the default scaffold.** Every `spectacular init` scaffolds the `scope-down` policy (`@Planning`, warn, `principle: 10`) and, when PRINCIPLES.md is scaffolded, Principle 10 *Build the smallest verified slice, full scope in mind*. Propagated into `doc_policy()` + `doc_principles()` (CLI), `templates/policy/base.md`, and `scaffold-reference.md` — v1.12.1 added them only to this repo's workspace.

### Fixed
- **Backfilled Principle 9 into the default principles scaffold.** `doc_principles()` and `scaffold-reference.md` shipped 8 principles; Principle 9 (*Feedback ≠ verification ≠ benchmark*, added in v1.6.0) was never propagated. Added so the scaffold is contiguous 1→10 and `scope-down`'s `principle: 10` link resolves on a fresh `init --with principles`.

## [1.12.1] — 2026-05-31

### Added
- **`scope-down` policy (`@Planning`, warn) + Principle 10.** Practice layer: before fixing milestones, name the smallest high-impact slice that delivers the core value now and defer the rest to ROADMAP `v2+`; flag speculative generality and features without a current need. Theory: Principle 10 — *Build the smallest verified slice, full scope in mind* — counters the agent build-everything reflex (build less, as a finished block, future-proof). `warn`, not `block`: scope is a human judgment call (Principle 8). Added to this repo's workspace `.spectacular/POLICY.md` + `PRINCIPLES.md`; not yet propagated to the default `init` scaffold.

## [1.12.0] — 2026-05-31

### Added
- **POLICY.md — the practice layer (`references/policies-contract.md`).** A new always-set canonical doc, the operational sibling to `PRINCIPLES.md`: PRINCIPLES is *theory* (the why, optional), POLICY is *practice* (the how-we-actually-work, the floor). Policies are filed under named **work-phase hooks** and the skill retrieves only the active hook's policies on entering a phase — progressive disclosure (Principle 6) applied to the rule layer. Deliberately asymmetric with optional PRINCIPLES: every `spectacular init` scaffolds POLICY.md with 8 prefilled defaults.
- **8 work-phase hooks (`@` reads "at").** Spine: `@Init`, `@Planning`, `@Implementation`, `@Verification`, `@Archive`. Moments: `@Remember`, `@Snapshot`, `@SessionEnd`. The before/after verb lives in the *policy name* (`understand-before-change`), never the hook.
- **8 prefilled default policies — 4 block · 4 warn.** Block: `understand-before-change` (@Implementation), `verification-present` (@Verification, absorbs verify-walk's gate), `confirm-before-write` (@Remember), `snapshot-before-overwrite` (@Snapshot). Warn: `scaffold-contract` (@Init), `request-shape` (@Planning), `spec-sync`+`memory-propose` (@Archive), `summarize-before-handoff` (@SessionEnd). **Severity is opt-in to blocking** — a policy hard-stops only if it explicitly declares `severity: block`; absent/warn/unrecognized → surface-and-continue (no policy accidentally blocks).
- **`spectacular policy` verb — 5 forms.** `policy` (all, grouped by hook), `policy @<hook>` (one phase's policies + each linked principle's heading and one line), `policy <id>` (one policy, full text + principle), `policy --principle N` (reverse: which policies enforce principle N), `policy --json` (machine form). Skim-by-default, matching the read-verbs convention; merges POLICY.md (definition) with `config.yaml` overrides.
- **Injection loop + phase gate blocks (`references/policy-injection.md`).** Each phase reference doc (`init-workflow`, `new-request`, `active-request`/`lifecycle`, `verification`/`lifecycle`, `archive`, `memory`, `versioning`, `sessions-rules`) opens with a 2-line `@<hook> policy gate` instructing the skill to run `spectacular policy @<hook>` first. The ref doc *is* the phase boundary — no event bus, no `hooks.json` wiring (skill-native; works in bare-CLI and installed-plugin sessions alike).
- **`## Understanding` PLAN slot.** Optional authoring slot (`How it works now` / `What changes` / `What stays the same`) required before `planned → active` by `understand-before-change`; escalates to a dedicated `requests/<slug>/UNDERSTANDING.md` for large requests (satisfied by either — the VERIFY.md 2-of-N pattern). No `ANALYSIS.md`.
- **`doctor policies` area.** Mechanical structure check (POLICY.md present + frontmatter; every blocker has a `check:`; severities are `block|warn`; hooks are from the locked 8; no orphan sections) plus the `understand-before-change` presence-check on every active request. `--fix` re-scaffolds a missing/empty POLICY.md.
- **`config.yaml` `policies:` override layer.** Per-policy `enabled` / `severity` overrides and custom-policy registration (declare a `hook:` to add your own). POLICY.md is the source of truth; config tunes it — layers, not competing copies. Commented stanza shipped in the init scaffold. Scope is config-only in v1 (4-tier precedence deferred to v2).

### Changed
- **Always-set grows to five docs** — `prd spec config agents policy`. POLICY.md joins the always-set scaffold and the `doctor workspace` always-set check; `doctor --fix` re-scaffolds it like any other always-set file.
- **`PRINCIPLES.md` stays optional** — explicitly the asymmetry: theory is optional reading, practice (POLICY) is the operational floor. The optional `principle: N` tag links a policy back to the principle it enforces.

### Notes
- Self-dogfooded: this repo's `.spectacular/POLICY.md` was scaffolded via `doctor --fix`, and the `policy-engine` request itself was promoted `planned → active` through the `understand-before-change` gate with a filled `## Understanding`.
- No harness `hooks.json` wiring in v1 — enforcement is skill-side + doctor; kernel-level locks are the v2 upgrade path. verify-walk is absorbed as the `verification-present` policy but not refactored onto the engine in this release.

## [1.11.0] — 2026-05-30

### Added
- **Validation walk — `spectacular verify <slug>` (`references/verify.md`).** A skill-side interactive ritual that moves a request `review → verified` by *running* its checks, not just claiming them. Closes PRINCIPLES.md Principle 7 (the validation layer). Skill-only; the CLI redirects.
- **Typed verification checks — verification is multi-authority, not one thing.** Five kinds along a deterministic → judgment → human spine, each verified by its own authority: `executable` (`` `run: <cmd>` `` → exit code), `assertable` (`{assert}` → agent checks a binary property), `judgable` (`{judge}` → LLM reasons over artifacts), `observable` (`{observable}`, the default → human looks), `manual` (`{manual}` → human acts then confirms). Tags work **inline per-line** or **section-grouped** (`## Title {kind}`, absolute); executable checks confirm-before-run with a batch-allow option.
- **`VERIFY-LOG.md` — append-only walk audit trail.** Each walk records every check with the `[kind]` that confirmed it (evidence, reasoning, exit codes), so verification becomes a recorded event, not just a checkbox state. Stubs for VERIFY.md (typed) + VERIFY-LOG.md added to `scaffold-reference.md`.
- **`spectacular verify` CLI redirect + docs.** `verify` dispatches as a skill-only verb (terminal prints a redirect); `docs/commands.md` documents the verb + the kind taxonomy.

### Changed
- **`spectacular archive` warns on verified-without-a-walk.** If a request is `verified` but has no `VERIFY-LOG.md`, archive emits an advisory (non-blocking) note — it was flipped verified without running the walk.
- **Lifecycle routing updated.** SKILL.md routes `review → verified` through the verify walk; `verification.md` (where checks live) and `verify.md` (how they're walked) are cross-linked as the two halves of one system.

> Dogfood: verify-walk was verified by its own mechanism — 10 typed checks across all five kinds, walked end-to-end, producing the first VERIFY-LOG. See `.spectacular/archive/verify-walk/`.

## [1.10.0] — 2026-05-29

### Added
- **Per-capability specs for the two densest capabilities.** `.spectacular/specs/doc-engine/SPEC.md` and `.spectacular/specs/roadmap/SPEC.md` promote the registry-driven doc engine and the structured-roadmap artifact out of the cramped `SPEC.md` index into standalone specs. The doc-engine spec documents the full mode taxonomy (now correctly **9 modes**, including the previously-undocumented `index` soft-DB mode), a drift-proof by-scope registry (no hardcoded doc count), and carries an inline design-decision log.
- **Self-describing skill reference docs + a catalog script.** Every `skills/spectacular/references/*.md` now carries `description` + `when_to_use` frontmatter (mirroring the `SKILL.md` field convention); `scripts/catalog.sh` renders the catalog from that frontmatter — `--when` (with load triggers), `--missing` (lint for undocumented docs), `--json`. The catalog is self-maintaining; `SKILL.md`'s reference-loading table remains the authoritative routing source.

### Changed
- **`SPEC.md` index entries compressed.** The doc-engine and structured-roadmap bullets collapse to one line + a link to their new capability specs; the index stays terse. Pre-edit snapshot at `snapshots/SPEC/@v2.md`.
- **README repositioned around spec-driven development.** New thesis ("No spec. No plan. No clue." / "Agents build. Humans decide."), a colorful 6-benefit SVG grid, an extensible "Works well with" ecosystem block, and an updated banner.
- **`SKILL.md` registered-docs reference de-hardcoded.** Dropped the stale "18 doc IDs (v1.7.0)" list — the live registry is the `references/*-rules.md` set, catalogued in `doc-index.md`. The skill description now also names the soft-DB collections (memory, sessions, feedback, ideas).

### Fixed
- **`index`-mode docs no longer fall through the grill router.** `grill.md` listed only `append`/`stub`/`freeform`/`reference` as non-grill routes, so a soft-DB doc (memory/sessions/feedback/idea) reaching the grill flow had no route. It now redirects to the doc's CLI mutator.

### Removed
- **Legacy `prd-grill.md` / `prd-refine.md` / `prd-review.md` reference docs.** Superseded by the generic engine (`grill`/`refine`/`review` + `prd-rules.md`) since v1.4.0; the 546 lines of duplicated routing are gone (snapshots preserved in `versions/`). The orphaned `doc-registry@v1.md` snapshot (renamed `doc-registry` → `doc-index` in v1.4.0) moved out of `references/` into `versions/`.

---

## [1.9.0] — 2026-05-29

### Added
- **Versioning convention doc (`docs/versioning.md`).** Codifies how Spectacular versions itself: SemVer as the canonical scheme with a Spectacular-specific breaking-change trigger (renamed/removed verb or flag, changed invocation syntax, `.spectacular/` file-contract break = MAJOR); default-silent mechanical increments with an ask-first rule only for a probable MAJOR or a roadmap-pinned milestone; the single-canonical-version-source rule across all 7 version-bearing locations (flagging the `cli/spectacular` + `SKILL.md` manual-bump drift point); the standard `-alpha/-beta/-rc` pre-release ladder; and two opt-in, roadmap-only marketing layers — a pinned milestone number and a full Apple-style major-line **release arc** (`X.0` launch → staged `X.x` features → terminal stable before `X+1.0`). Registered in `docs/docs.yaml`.

## [1.8.4] — 2026-05-29

### Fixed
- **`spectacular remember` (and every template-backed verb) failed with "memory entry template not found" on the canonical symlinked install.** `SCRIPT_DIR` was computed from the symlink path (`~/.local/bin`) instead of its target, so the bundled-template fallback looked in `~/.local/skills/` which never exists. `SCRIPT_DIR` now resolves symlinks, and `_resolve_template` gained a scope-independent `~/.agents/skills/spectacular/templates/` fallback.
- **`spectacular new` silently scaffolded empty PLAN.md / TASKS.md** when run via the symlinked install — same `SCRIPT_DIR` root cause. Now renders full template content.
- **`spectacular decide` wrote dead empty sections.** The verb dumped all text into `**Decision:**` and left the other ADR sections blank. Root cause was schema drift between the CLI/templates (`Decision/Why/Tradeoffs`) and the rules doc (`Context/Decision/Consequences`).

### Added
- **`spectacular decide --context "..." --consequences "..."`** — populate those ADR sections at write time. The positional argument fills `**Decision:**`; omitted sections are emitted as empty headers to fill in later, never invented from the decision text.
- **Mutator failures now print a manual-recovery path.** A new `die_recover` helper emits the error plus a `→ Manual recovery:` hint (target file + frontmatter shape) so a broken template-backed verb is never a dead end. Wired into `remember`, `session start`, `idea new`, and `feedback-loop new`.

### Changed
- **Canonical ADR schema is now `Context / Decision / Consequences`** (Michael Nygard shape) across the CLI inline entry, `templates/decisions/entry.md`, and the embedded `doc_decisions` scaffold — reconciling the prior drift against `decisions-rules.md`.

## [1.8.3] — 2026-05-29

### Added
- **`spectacular init` now detects existing installs and skips redundant skill copies.** Before installing, init scans every place spectacular could already be available: the current project, parent directories up the worktree, the global user scope (`~/.agents` + `~/.claude`), and plugin installs (Claude Code `~/.claude/plugins/cache/`, Codex `~/.codex/plugins/cache/`, Gemini `~/.gemini/extensions/`). If any are found it warns, lists each location, and defaults to **not** installing a duplicate — the `.spectacular/` scaffold still proceeds.
- **`--skill-scope <project|global|none>`** flag for `init` to control where (or whether) the skill is installed. `project` = `./.agents` + `./.claude` (the prior default), `global` = `~/.agents` + `~/.claude`, `none` = scaffold only. When unset, init auto-resolves: skip if already available, else project. `--no-skill` is an alias for `--skill-scope none`; `--global` is now a deprecated alias for `--skill-scope global`.
- Interactive `init -i` gains a `none` option in the skill-scope prompt and defaults it to the detected value.

### Fixed
- **`spectacular decide "..." --dry-run` no longer creates `DECISIONS.md` as a side effect.** The bootstrap of a missing `DECISIONS.md` now runs only on a real write; the dry-run path previews `would create` + `would append` and writes nothing to disk.
- **Stale version constant:** `SPECTACULAR_VERSION` still read `1.8.1` despite the existing `v1.8.2` tag (the constant wasn't bumped in the v1.8.2 chore commit). Now `1.8.3`.

## [1.8.1] — 2026-05-26

### Fixed
- **Scaffold bug:** `templates/architecture/base.md` showed `current/` in the `.spectacular/` tree diagram instead of `specs/`, baking a stale (pre-v0.5) convention into every newly-initialized ARCHITECTURE.md. Now reflects the canonical `specs/` layout. Existing projects: `spectacular doctor` already flags legacy `current/` and `spectacular migrate` renames it — no action needed beyond editing the doc.

## [1.8.0] — 2026-05-26

### Added
- **Read-verb family (11 new top-level CLI verbs)** designed to collapse multi-step agent reads into single deterministic calls. Read-only — no state mutation.
- **`spectacular requests`** — list requests with `--status <s>`, `--active` (alias for `--status active`), `--since <Nd|Nh|Nw>`, `--limit N` (default 20), `--all`, `--json`. Default table view shows slug/status/priority/target/updated/summary.
- **`spectacular request <slug>`** — detail view (skim by default: frontmatter + section outline + milestone progress; `--full` for raw PLAN.md). Falls back to `archive/<slug>/PLAN.md` if not in `requests/`.
- **`spectacular decisions`** / **`spectacular decision <slug>`** — list and inspect entries from `.spectacular/decisions/`. Filters: `--tag`, `--since`. Detail view shows frontmatter + outline.
- **`spectacular memories`** / **`spectacular memory <slug>`** — same pattern for `.spectacular/memory/`.
- **`spectacular sessions`** / **`spectacular sessions show <slug>`** — read sessions; `show <slug>` is the detail subverb (avoids collision with the existing `session start|end` mutators). `--status open|closed|all`.
- **`spectacular show <doctype>`** — dump a canonical doc (`prd|spec|principles|architecture|roadmap|stack|agents|decisions|memory|sessions|personas`). `--section <name>` filters to one H2; `--json` returns `{path, content}`.
- **`spectacular summary`** — one-page workspace overview: project name + request counts by status + decisions/memories/sessions/ideas/feedback counts. `--json` for machine. Aggregates by calling the list verbs internally.
- **`spectacular progress <slug>`** — milestone tick rate parsed from TASKS.md. Returns `M1: 8/8 ✓, M2: 3/5, ...`. `--json` for machine.
- **`spectacular paths`** — JSON map of conventional workspace paths (PRD, SPEC, requests_dir, memory_dir, etc.). Default JSON; `--text` for human. Lets tools locate files without hardcoding.
- **Universal flags across all list verbs:** `--status`, `--since`, `--limit` (default 20), `--all`, `--json`. Default limit prevents context overflow; `--all` overrides.
- **Skim-by-default detail views** — single short prompt instead of full file content. `--full` always available for raw dumps.
- **Internal helpers:** `_parse_since` (duration parser), `_date_to_epoch`, `_json_escape`, `_request_iter_all`, `_decision_iter_all`, `_memory_iter_all`, `_session_iter_all`, `_skim_file`, `_progress_text`, `_progress_json`.

### Changed
- `skills/spectacular/SKILL.md` — new "Read verbs (v1.8.0+)" routing block with 12 trigger rows + cold-start pattern guidance. Version frontmatter to 1.8.0.
- `cli/spectacular` top-level `--help` adds a "Read verbs" section.
- `README.md` CLI reference adds all new verbs.

### Notes
- **Grammar locked (2026-05-26):** plural noun = list; singular noun + slug = detail; bare verb = high-frequency action on implicit object (`new`, `archive`, `promote`, etc.); noun + subcommand = multi-verb lifecycle (`session start|end`, `idea new|list|promote`, ...). Zero breaking changes — every existing verb stays. The new read verbs codify the pragmatic mixed grammar that was already emerging.
- **Pure verb-first was considered and rejected**: `spectacular new request|idea|session|feedback|...` would overload `new` across unrelated lifecycles, the same trap `git` solves with `git <object> <verb>`. The current mixed grammar matches real-world CLI conventions (git/gh/kubectl all blend bare verbs with noun-namespaces).
- **`spectacular summary` as cold-start primitive**: agents should prefer `summary → requests --active → request <slug>` (three calls) over walking the filesystem. Documented in SKILL.md.
- **Why hardcoded doctype list in `show`**: doc-types are rare additions; dynamic registry dispatch adds parsing cost on every call. The hardcoded switch is faster and the v1.7.0 doc-index.md already serves as the catalog.

---

## [1.7.0] — 2026-05-26

### Added
- **New doc-type `idea`** registered as `mode: index`. Lives at `.spectacular/ideas/`. Promotes the previously folder-only `ideas/` convention to a first-class doc-type with rules file, template, doctor area, and CLI verbs. Dispatch via `skills/spectacular/references/idea-rules.md`. No top-level `IDEAS.md` index file — folder listing is canonical.
- **CLI verbs: `spectacular idea new|list|promote`** — scaffold an idea entry (status `parked`), list across the folder (with `--status` filter), promote to a full request (scaffolds via `cmd_new`, sets `promoted_to:`, moves source to `archive/ideas/`). Verb surface mirrors `feedback-loop`.
- **Doctor area: `ideas`** (judgment-only, no `--fix`). Flags: required frontmatter missing, `status: exploring` entries older than 90 days, orphan `promoted` entries still living in `.spectacular/ideas/` instead of `archive/ideas/`, unknown status values. `DOC_AREAS` count grows from 14 to 15.
- **Status lifecycle for ideas:** `parked` (captured, not actively shaping) → `exploring` (actively thinking) → `promoted` (became a request). Promotion is explicit and one-way via the CLI verb.
- **New template: `templates/idea/base.md`** — frontmatter stub + 4 required body sections (Hypothesis / Context / Open questions / Promoted-to).
- **`spectacular idea promote <slug>`** annotates the new request's PLAN.md with a banner pointing back at `archive/ideas/<slug>.md` so the original content stays discoverable.

### Changed
- `SKILL.md` routing table adds three idea trigger rows. Doc IDs registered string bumped to v1.7.0 (`idea` added). Version frontmatter to 1.7.0.
- `doc-index.md` adds the `idea` catalog row under "Project-wide canonical docs".
- `references/doctor-areas.md` documents the new `ideas` area (4 check types, judgment-only rationale).
- `.spectacular/ARCHITECTURE.md` § Ideas layer adds CLI verb cross-references to `idea-rules` and `doctor-areas`.
- `cli/spectacular` top-level `--help` adds `idea <sub>` row; doctor area list adds `feedback` + `ideas` rows previously missing from `doctor_usage`.

### Notes
- v1.7.0 closes the gap where `ideas/` existed in ARCHITECTURE.md as a documented convention but the skill couldn't create, list, doctor-check, or promote ideas — every interaction required manual file edits. The `feedback` doc-type (v1.6.0) served as the canonical pattern; the `idea` shape mirrors it deliberately.
- Status enum was decided as `parked|exploring|promoted` (three states, locked design 2026-05-26). Two-state and free-form variants were considered and rejected — three states preserve a "currently shaping" signal that doctor can act on while keeping the surface narrow.
- The "abandon a request into an idea" flow (used manually 2026-05-26 to convert the `memory-protocols` request into `.spectacular/ideas/memory-protocols.md`) is explicitly out of scope. It's a one-off corner case, not a recurring pattern worth verbifying.

---

## [1.6.0] — 2026-05-25

### Added
- **New mode: `feedback-loop`** — prototyping-stage human-feedback acquisition. A 5-step interaction (pick target → craft proposal → ask user → capture response → decide next action) for probing system fitness. Not a benchmark, not verification, not `review` — a distinct axis orthogonal to all three. Full spec in `skills/spectacular/references/feedback-loop.md`.
- **New doc-type `feedback`** registered as `mode: index`. Lives at two locations: system-level (`.spectacular/feedback/`) and request-scoped (`.spectacular/requests/<slug>/feedback/`). Dispatch via `feedback-rules.md`. No top-level `FEEDBACK.md` index file — folder listing is canonical.
- **CLI verbs: `spectacular feedback-loop new|list|resolve|archive`** — scaffold an entry, list across both locations, close with a decision, manually curate to archive. `--request <slug>` scopes to a request; `--next-action` is required on resolve (no silent `tbd` resolutions).
- **Hidden alias routing** — `iterate`, `experiment`, `test`, `probe`, `try` all route to `cmd_feedback_loop`. Not shown in `--help` per the contract — only `feedback-loop` is documented as the official mode name. Hidden aliases give ergonomic short forms without cluttering discovery.
- **Doctor area: `feedback`** (judgment-only, no `--fix`). Scans both feedback locations. Flags: required frontmatter missing, `status: open` entries older than 30 days, orphan back-refs (`request:` field pointing to a missing request folder). `DOC_AREAS` count grows from 13 to 14.
- **PRINCIPLES.md §9 — "Feedback ≠ verification ≠ benchmark"** — codifies the three-axis distinction so future work doesn't conflate them. Explicit guard against the word "evals" (carries HumanEval/MMLU baggage that pulls the wrong way).
- **Proactive-surfacing contract** — the skill may offer a feedback-loop session at exactly three checkpoints (milestone tick, request status → `review`, end of archive flow). Never mid-flow. Never unsolicited. Single short prompt; user accepts or declines.
- **Bidirectional back-refs** — feedback entries scoped to a request carry `request: <slug>` in frontmatter; the request's PLAN.md gets a `feedback:` list. When a feedback resolution spawns a new request, the spawned request's PLAN.md gets `spawned_by_feedback:` pointing back.
- **Auto-promotion to memory contract** — when a feedback resolution captures a durable preference, the skill explicitly confirms before writing a memory entry. CLI flag `--promote-hint` prints the suggested `spectacular remember` command rather than silently writing (memory promotion is a judgment call that needs an LLM in the loop).
- **New template: `templates/feedback/entry.md`** — frontmatter stub + 7 required body sections (Target / Hypothesis / Proposal / Question asked / User response / Insight / Decision).

### Changed
- `SKILL.md` routing table adds the feedback-loop mode block + alias routing + three-checkpoint surfacing rules. Doc IDs registered string bumped to v1.6.0 (`feedback` added).
- `doc-index.md` adds the `feedback` catalog row.
- `ARCHITECTURE.md` documents both feedback folder locations (system-level and per-request) in the directory trees.
- `doctor-areas.md` documents the new `feedback` area with check matrix.
- `top_usage` lists `feedback-loop <sub>` under CLI verbs; help-string area count bumped from "10 areas" to "14 areas" (matches actual `DOC_AREAS`).
- `doctor_parse_args` recognizes `feedback` as a scoped area.

### Notes
- **Dogfooded the mode on itself before shipping.** A request-scoped feedback-loop session on the feedback-loop CLI surfaced 4 ergonomic issues — all fixed in the same release: `list` drops the redundant DATE column and truncates slugs >32 chars to 29+`...`; `--promote` renamed `--promote-hint` (honest about advisory behavior); `resolve` requires `--next-action` (clear error guides to `park` if undecided). See `.spectacular/requests/feedback-loop/feedback/2026-05-25-feedback-loop-cli-ergonomics-after-m0-m4.md`.
- The mode is explicitly **prototyping infrastructure** — not a benchmark, not automated grading. Feedback compounds across sessions as durable insight; promotion to memory is the way feedback graduates into preferences.
- Composes cleanly with `memory-protocols` (planned v1.6.x+) — auto-promotion will get smarter as memory protocols formalize.

---

## [1.5.0] — 2026-05-25

### Added
- **Two new doc-types: `memory` and `sessions`** registered as `mode: index` — soft-folder databases with an index file (`MEMORY.md` / `SESSIONS.md`) regenerated from per-entry markdown files in `memory/` / `sessions/`. Frontmatter-driven dispatch via `memory-rules.md` and `sessions-rules.md` plugs into the v1.4.0 doc-writing substrate with no special-casing.
- **CLI mutator: `spectacular remember "<text>" [--tag a,b] [--dry-run]`** — writes one memory entry with auto-derived slug + summary, regenerates `MEMORY.md` index. Bare `spectacular remember` / `remember this` still routes to the skill flow (backwards-compatible).
- **CLI mutator: `spectacular decide "<text>" [--dry-run]`** — appends one ADR-style entry to `DECISIONS.md`. Auto-derives a title from the first ~6 words.
- **CLI mutator: `spectacular session start|end`** — opens/closes a working session entry. Lifecycle invariant enforced: at most one session can be open at a time.
- **Auto-session linkage** — when a session is open, `decide` and `remember` set `session: <slug>` in the new entry's frontmatter. At `session end`, the writer scans `DECISIONS.md` + `memory/*.md` for matching `session:` fields and appends "Linked decisions" / "Linked memories" sections to the session body, plus recomputes `decisions_count` / `memories_count`.
- **New `mode: index`** taxonomy entry in `doc-index.md`. The index file is regenerated from `entries-dir/`; CLI mutators write entries; agentic verbs (grill/refine/review) operate on the collection.
- **Doctor areas: `memory` and `sessions`** — index ↔ entries drift detection, frontmatter validation, lifecycle-invariant check (≤1 open), 4h stale-open-session warning.
- **Coding kit triggers** — `--kit coding` now scaffolds `MEMORY.md` + `SESSIONS.md` alongside the existing `DECISIONS.md` + `STACK.md` + `ARCHITECTURE.md`.
- **New templates: `templates/memory/entry.md`, `templates/sessions/entry.md`** — frontmatter schemas with `type`, `summary`, `tags`, `session` (memory) and `type`, `status`, `start_date`, `end_date`, `decisions_count`, `memories_count` (sessions).
- **Snapshot tidy (M1–M3 of the `snapshot-tidy` request)** — versioned snapshots now live in a dedicated `.spectacular/snapshots/<DOC>/@v<N>.md` tree (one folder per canonical doc, uppercase preserved, `@v` retained in filename). Sub-doc snapshots mirror their path: `specs/cli/SPEC.md` → `snapshots/specs/cli/SPEC/@v1.0.md`.
- **Doctor `snapshots` area extended** — warns on legacy root-level `*@v*.md` files with the target path in the fix hint. `spectacular doctor --fix snapshots` migrates them via `git mv` (or plain `mv` when untracked).

### Changed
- `decisions-rules.md` documents the new `spectacular decide` CLI verb and the optional `Session:` link field appended to entries when a session is open.
- `SKILL.md` routing table adds the three CLI mutators; Doc IDs registered string bumped to v1.5.0 (`memory`, `sessions` added).
- `doc-index.md` adds two catalog rows + new `index` row in the mode taxonomy + a new column in the verb × mode matrix.
- **`spectacular snapshot <file>`** writes to the new `snapshots/<DOC>/@v<N>.md` location; reads from both new and legacy locations when computing the next N (back-compat preserved). Auto-creates the target directory.
- **`versioning.md` + `ARCHITECTURE.md`** updated for the new snapshot layout; migration notes inline.
- Dogfood: 11 snapshots in this repo (PRD ×4, ROADMAP ×4, AGENTS, ARCHITECTURE, SPEC) migrated from `.spectacular/` root to `.spectacular/snapshots/<DOC>/`.

### Fixed
- `_summary_from_text` helper no longer crashes on UTF-8 multibyte chars (em-dash, etc.) — switched from awk char-iteration to sed-based extraction.
- Session-end body builder no longer silently exits under `set -o pipefail` when no matching `## ` headers found — grep chain wrapped in `|| true` and switched to awk for the ADR scan.

### Notes
- **Migration of flat `DECISIONS.md` → `decisions/<slug>.md` folder shape is deferred to v1.6.x** alongside query verbs (`spectacular decisions --7d`, `spectacular recall`, `spectacular sessions`). v1.5.0 leaves `DECISIONS.md` in its existing flat format.
- The frontmatter schema (`type`, `tags`, `summary`, `related`, `session`) is **RAG-ready** — future embedding/retrieval layers can read these fields without a schema change.
- This is the foundation block for the planned v1.5.x → v1.7.x memory line on the roadmap. Research synthesis lives in `_research/agent-memory/REPORT.md` (NotebookLM + scrapekit consolidation of mem0, Letta/MemGPT, Graphiti/Zep, Cognee, Anthropic Memory Tool, Cline Memory Bank, Cursor Rules patterns).

---

## [1.4.0] — 2026-05-24

### Breaking
- **`mode: reps` removed** from the substrate. Migrated to `mode: grill-each` (per-block grill walk). Existing rules files in `skills/spectacular/references/` are auto-migrated. Custom packs declaring `mode: reps` need a one-line update.
- **`doc-registry.md` renamed to `doc-index.md`** and reframed as a human-readable catalog. Dispatch (mode, slots, template, location, scope, snapshot-on-edit, kit-support) now lives in each `<doc-id>-rules.md` file's **frontmatter**, not in the index. Snapshot preserved at `doc-registry@v1.md`.
- **Top-level "engine" terminology dropped** across the skill (~37 architectural occurrences). The shared verbs are now described as skill flows (grill / refine / review) — not "the engine".

### Added
- **Grill sub-modes** — `mode: grill` is now an alias for `grill-wide` (single broad pass); new values `grill-wide` / `grill-each` / `grill-loop` describe interaction shape. `grill-loop` is a new wide-then-deep style (fast pass with short answers, then revisit slots flagged vague/incomplete via the heuristic: length < 30, vague-word match, placeholder string, or gate-check fail).
- **Flag override** — `spectacular <doc> grill --wide | --each | --loop` forces a sub-mode for this session, overriding the doc's declared mode.
- **Per-doc rules files for the 6 implicit docs** — `principles-rules.md`, `architecture-rules.md`, `stack-rules.md`, `agents-rules.md`, `spec-rules.md`, `decisions-rules.md`. Every registered doc now has a rules file (consistency over brevity).
- **Frontmatter schema on every rules file** — 4 required fields (`doc-id`, `mode`, `location`, `scope`) + mode-conditional (`template`, `slots`, `kit-support`, `snapshot-on-edit`) + 3 optional (`summary`, `version`, `status`). Strict — `doctor frontmatter` will validate.
- **CLI agentic-verb redirect** — typing `spectacular <doc> grill | refine | review` at terminal prints a friendly redirect to run inside Claude Code or Codex. Agentic verbs require an LLM; mechanical verbs (`new`, `archive`, `snapshot`, `init`, `doctor`, `pack`, `migrate`) continue to run in CLI.
- **Verb × mode matrix** documented in `doc-index.md` — defines behavior for every cell including the previously undefined ones: `grill × stub` (polite hint + optional `--wide` override), `grill × freeform` (open-ended prompt; skill infers slot list on the fly), `refine × append` (user picks scope: latest / all / pick).
- **PLAN.md `phase:` axis** — Phase (lifecycle), Verb (action), Mode (doc shape) are now treated as three orthogonal axes, never collapsed. Phase taxonomy: `discover / spec-refine / mvp / iterate / test / release-prep / release`.
- **`KNOWN_DOCS` extended** — CLI now recognizes `plan` and `tasks` for doc-verb dispatch and redirect.

### Changed
- **`grill.md` rewritten** — mode resolution section, sub-mode dispatch logic, flag-override behavior, grill-loop algorithm + vagueness heuristic, 7 worked examples (PRD wide, PLAN wide, ROADMAP each, PERSONAS each, PRD loop override, DECISIONS append, AGENTS stub).
- **Rules-file H1s standardized** — "X Overrides" → "X Rules" across prd, plan, tasks, roadmap, pack, docs. The word "overrides" was misleading: rules files don't override anything, they declare per-doc behavior consumed by the shared skill flows.
- **SKILL.md routing table + references index** — updated for v1.4.0 (15 reference rows for the doc-writing layer; doc-id list bumped to 14; references-table groups rules files together).
- **`doc-index.md`** is now a human catalog, not a dispatch contract. Sections: project-wide / per-request / user-scope / public-facing (deprecated) / skill-internal. Mode taxonomy table + verb × mode matrix included.
- **`docs/commands.md`** — added agentic vs mechanical verb table; grill sub-modes section; v1.4.0 doc list.
- **`.spectacular/SPEC.md`** — Doc-writing capability bullet rewritten for v1.4.0 (rules-files-as-dispatch, verb taxonomy, agentic/mechanical split).
- **`CONTRIBUTING.md`** — new-doc-type contribution guide updated to point at rules-file + template + catalog row pattern.

### Fixed (from substrate audit — codex G1 + G2 findings)
- **G1: registry's "no code changes" claim is now honest.** Pre-v1.4.0 `doc-registry.md` claimed adding a new doc required no code changes, but the CLI's `KNOWN_DOCS` constant and dispatch logic needed editing. Now: CLI reads only catalog fields (doc-id, location, mode) for `init` / `doctor`; agentic dispatch lives entirely in the skill and reads rules-file frontmatter. New doc = rules file + template + catalog row. No CLI edits.
- **G2: "generic engine" overclaim removed.** The skill flows (grill / refine / review) aren't one unified engine — they're verb-specific behaviors that dispatch on mode. Docs now describe them honestly.

### Migration notes
- **No user action required for in-repo workspaces.** Rules files in `.spectacular/` are bundled by the skill, not user-authored. They migrate transparently when the user updates the skill.
- **Custom packs / project-local rules files** using `mode: reps` should update to `mode: grill-each`. `doctor packs` will warn.
- **External tooling reading `doc-registry.md`** by path needs to point at `doc-index.md`. A snapshot is preserved at `references/doc-registry@v1.md`.

### Process
- substrate-clarity request shipped via M1 discovery grill + M2 spec-refine + M3-M8 build. 7 decisions locked during the discovery session (see `.spectacular/archive/substrate-clarity-v1.4.0/discovery.md` post-archive).

---

## [1.3.0] — 2026-05-24

### Added
- **PERSONAS.md as opt-in canonical doc** — proto-audience profiles + user stories. Each persona is ~6-10 lines (Who / Wants to / Pain / Stories / Not for). Grill walks the 5 slots per persona, then asks "add another?" — same shape as ROADMAP version blocks. Triggered by `product` + `content` PRD kits via `triggers-docs.always`, or opt in on any kit with `spectacular init --with personas`.
- **New `personas` doctor area** — validates frontmatter, counts persona blocks, per-persona story counts, soft-warns at >5 personas. Absent → info note (one-time skip).
- **`personas-rules.md` reference** — slot definitions, grill prompts, vague-word lists, anti-patterns (no JTBD framework, no demographics, no >5 personas), gate checks. JTBD framework apparatus explicitly out of scope.
- **`templates/personas/base.md`** — template with 5 slots.
- **`.spectacular/PERSONAS.md`** — dogfood example for this repo (3 personas: Solo OSS maintainer, Small-team tech lead, Tool builder using AI agents to build AI tools).

### Changed (breaking — naming-only, no schema removal)
- **Mode renames in doc-registry:**
  - `structured` → `reps` (grill walk repeats per block; ROADMAP, PERSONAS)
  - `freeform` → `stub` (scaffold + exit, no walk; STACK, AGENTS, PRINCIPLES, ARCHITECTURE, SPEC)
  - New `freeform` reserved for "agent has full creative liberty over structure" (no template, no slots) — no docs in v1 use this mode
- **Per-doc rules file rename:** `*-overrides.md` → `*-rules.md` (7 files: prd, plan, tasks, roadmap, docs, pack, personas). The registry field renamed from `overrides:` to `rules:`. Rationale: these files don't override anything — they're per-doc rules (grill prompts, vague-word lists, gate checks). "Generic engine + per-doc rules" is the accurate model.
- **`doc-registry.md` schema docstring + Adding-a-new-doc-type section** updated for new names.
- **`SKILL.md` v1.3.0:** doc IDs registered list now includes `personas` (14 total); references index lists `personas-rules.md` + `roadmap-rules.md`; schema field text updated.

### Fixed (from codex adversarial review)
- **Personas template `<DATE>` → `<today>`** placeholder consistency with PRD template convention.
- **Doctor per-persona story counting** — was project-wide grep; now awk walks blocks and reports which personas have zero stories.

### Migration notes
Existing projects do not need a migration. The renames are agent-facing only:
- Existing `.spectacular/` workspaces continue working unchanged (no override-named files exist in user projects)
- Existing PRD/PLAN/TASKS/ROADMAP workflows continue working — only the underlying registry mode strings changed
- If you've authored project-local kit files or rules files, rename `<id>-overrides.md` → `<id>-rules.md` and update `overrides:` → `rules:` in any referenced YAML

### Test changes
- All 9 existing test files continue passing (no test file required renaming)
- Personas doctor wiring validated via repo's own PERSONAS.md (3 personas, 10 stories, 0 errors)

---

## [1.2.1] — 2026-05-24

### Changed
- `.spectacular/SPEC.md`: "Public docs surface" capability bullet rewritten to reflect v1.2.0 deprecation and pageworks handoff; summary line bumped to v1.2.0.
- `.spectacular/ROADMAP.md`: prepended 6 missing "Recently shipped" entries (v0.7.1, v0.7.2, v1.0.0, v1.0.1, v1.1.0, v1.2.0).
- `README.md`: version badge bumped; skill commands table refreshed (added `promote`, `touch`; archive description references SPEC.md/specs/ sync instead of legacy `current/`); deprecation banner for `docs *` verbs added; doctor area list now includes `specs` and `docs`; CLI reference now lists `spectacular migrate`; new "Pairing with pageworks" section.

### Notes
- Docs-only patch. No behavioral changes to CLI or skill.

---

## [1.2.0] — 2026-05-23

### Deprecated
- **Public-facing docs surface** moved to a new standalone skill: [pageworks](https://github.com/alexsmedile/pageworks). The verbs `spectacular docs init`, `spectacular docs export`, `spectacular docs new`, `spectacular docs review`, and `spectacular docs status` still work but now emit a deprecation banner to stderr pointing at the pageworks equivalent. Removal target: spectacular v2.0.0.
- References `docs-contract.md`, `docs-overrides.md`, and `docs-renderer-adapters.md` gained a `> DEPRECATED in v1.2.0` banner at the top. They remain loaded for backward compatibility; canonical versions live in pageworks (`references/contract.md`, `references/authoring.md`, `references/renderers.md`).

### Changed
- `spectacular doctor docs` slimmed from ~190 lines of validation (schema + frontmatter + orphans + renderers) to ~25 lines of **discovery-only** checks: docs/ folder presence, docs.yaml manifest presence, and a pageworks-install hint when the pageworks CLI is not in PATH. Full validation lives in pageworks (`pageworks doctor`).
- `spectacular docs --help` rewritten with the deprecation banner and per-verb migration hints.
- Top-level usage (`spectacular --help`) marks `docs <...>` as DEPRECATED.
- Skill SKILL.md: "Public-facing docs verbs" section rewritten to route to the new pageworks-handoff reference. References index updated to flag the deprecated docs-* files.
- `.spectacular/AGENTS.md`: new "Don't write into docs/" rule + new "Skill boundary — spectacular vs pageworks" section.

### Added
- New reference: `skills/spectacular/references/pageworks-handoff.md` — when and how spectacular delegates public-doc work to pageworks. Canonical install hint phrasing, archive-time prompt mechanics, status-briefing reference, anti-patterns.
- `spectacular archive <slug>` now prints a one-time hint after archive completion when:
  - `docs/` folder exists in the project, AND
  - the archived request's PLAN/TASKS/VERIFY references SPEC.md, specs/, ARCHITECTURE.md, or PRD.md
  The hint suggests `pageworks audit` (when pageworks is installed) or surfaces the install URL (when not). Suppress per-call with `--no-docs-prompt` or per-project with `docs.prompt_on_archive: false` in `.spectacular/config.yaml`.
- `--no-docs-prompt` flag for `spectacular archive`.
- Tests: `tests/cli/archive-pageworks-prompt.test.sh` — 5 scenarios covering the prompt fires/suppresses correctly across docs/ presence + spec-reference detection + flag + config.

### Test changes
- `tests/cli/docs.test.sh` rewritten to assert the v1.2.0 deprecation state: verbs still work, banners emitted, doctor docs is discovery-only, skill verbs refused with deprecation. Previous deep-validation scenarios moved to pageworks's `tests/cli/doctor.test.sh`.
- `tests/cli/docs-export.test.sh` scenarios 13–16 (doctor renderers validation) moved to pageworks's test suite. The remaining 12 scenarios (CLI export behavior) stay — `spectacular docs export` still works for backward compatibility.

### Migration guide for users on spectacular v1.1.0
Nothing breaks immediately. Three paths:
1. **Keep using `spectacular docs ...`** — works through v1.x, removed in v2.0.0.
2. **Migrate to pageworks** — `curl -fsSL https://raw.githubusercontent.com/alexsmedile/pageworks/main/cli/install.sh | bash`. Then use `pageworks <verb>` everywhere you used `spectacular docs <verb>`. Same schema, same behavior, plus new Diátaxis page templates, prose patterns reference, maintenance/drift detection groundwork.
3. **Mute the deprecation banner project-wide** by setting `docs.prompt_on_archive: false` in `.spectacular/config.yaml` (only affects the archive-time prompt; per-verb banners still fire).

---

## [1.1.0] — 2026-05-23

### Added
- `spectacular docs export <renderer>` — generates renderer configs from `docs/docs.yaml`. Two adapters ship: `mkdocs` (writes `mkdocs.yml`, Material theme defaults) and `docusaurus` (writes `docusaurus.config.js` + `sidebars.js`).
- Both adapters also write `.github/workflows/docs.yml` for GitHub Pages deployment (idempotent, opt out with `--no-workflow`).
- `--force` overwrites generated files; default re-runs are safe and report skipped targets.
- `// spectacular: do-not-overwrite` magic comment pins manual edits — respected even with `--force`.
- `--out <path>` overrides default output location (default: alongside `docs/` at repo root).
- Optional `renderers:` block in `docs.yaml` carries per-renderer hints (`theme`, `primary`, `scheme`, `repo_url`, `edit_uri`, `organizationName`, `projectName`, `preset`).
- Doctor `docs` area validates the `renderers:` block: pass on recognized renderer keys (`mkdocs`, `docusaurus`), warning on unknown keys (typo guard), error on scalar or list shapes.
- Reference doc: `skills/spectacular/references/docs-renderer-adapters.md` — mapping tables, frontmatter translation, GitHub Pages boilerplate, contributing-a-renderer contract.
- `tests/cli/docs-export.test.sh` — 16 scenarios covering both adapters, idempotency, `--force`, pinning, error paths, doctor checks.

### Changed
- `spectacular docs --help` and top-level usage now list `export` and the two shipped renderers.
- `docs init` scaffolds a commented `renderers:` example in the generated `docs.yaml` (and template).
- `docs-contract.md` extended with the full `renderers:` block schema, recognized renderer table, and doctor validation rules.

### Notes
- Mintlify and Fumadocs adapters are not shipped — community-contributable via the contract in `docs-renderer-adapters.md` § Contributing a renderer.
- Empty sections (declared with `pages: []`) are dropped from generated nav for both adapters — keeps `mkdocs.yml` valid and Docusaurus sidebars uncluttered.

---

## [1.0.1] — 2026-05-23

### Added
- `spectacular --version` / `-v` / `version` prints the CLI version.
- `spectacular` / `spectacular help` / `spectacular --help` / `-h` (no subcommand) prints a top-level usage with the full verb list. `spectacular init --help` still prints the init-specific flag reference.

### Fixed
- `.codex-plugin/plugin.json` version drifted to `0.6.0` during the v1.0.0 cut; realigned to match the Claude plugin manifest.

---

## [1.0.0] — 2026-05-23 — first stable release

**Spectacular reaches v1.** The surface developed across v0.6.0 → v0.7.5 is now frozen as the stable contract. Future changes follow strict semver: breaking changes require a major bump; new capabilities require a minor; fixes ship as patch releases.

This release is a **tag, not new code** — it captures the dev arc's endpoint. All capabilities below are inherited from the v0.x lineage; per-version detail lives in the entries that follow.

### v1 surface (what's stable)

**Workspace foundation**
- `.spectacular/` directory contract: PRD, SPEC, config, AGENTS, requests/, specs/ as always-set
- Frontmatter as the signal layer — skill reads frontmatter, never full file bodies, for briefings
- 5-state lifecycle (planned → active → review → verified → archived) stored in PLAN.md frontmatter
- Convention packs (4-tier scope: project-local → user → app-store → bundled; 3 modes: suggest / scaffold / enforce; 6 rule categories)
- Snapshot-before-edit on canonical docs (`<FILE>@vN.md`)
- Memory (`spectacular remember this`) writes to `.spectacular/memory/`, git-committed
- Workspace schema versioning (`workspace_schema:` field) + migration registry + `spectacular migrate` verb

**Doc-writing engine**
- One generic grill / refine / review engine + registry-driven doc types
- 13 registered docs: prd, spec, plan, tasks, principles, architecture, roadmap, stack, agents, decisions, convention-pack, docs-manifest, docs-page
- Kit system: PRD has 5 kits (blank, coding, content, product, research) with extension contract
- Structured ROADMAP with 3 precision tiers (full / themed / vision), 9-phase chain, Icebox section
- 2-of-6 verification rule (when VERIFY.md is required vs. folded into PLAN/TASKS)

**CLI surface** (`spectacular`)
- `init` (kit + flags + interactive mode + --update for skill refresh)
- `doctor` (10 areas: skill, workspace, frontmatter, snapshots, links, lifecycle, kits, conventions, specs, docs; --fix for mechanical repairs; precondition for destructive verbs)
- `pack` (list, install, show, remove, with 4-tier scope resolution)
- `docs` (init, doctor passes)
- `migrate` (--dry-run, --to, --from, --list; registry-driven)
- **Mutator verbs**: `new`, `promote`, `snapshot`, `archive`, `touch` — CLI mutates atomically; skill orchestrates
- `status --against-latest` (workspace-schema verdict)
- `status --since <date>` (activity report; YYYY-MM-DD / 7d / 2w / 1m / yesterday)

**Skill surface** (`/spectacular`)
- Lean orchestrator (~120 lines) + on-demand reference loading
- Routing for status, new, archive, doctor, migrate, doc verbs, memory, snapshot, lifecycle transitions
- Skill walks judgment-requiring flows; CLI handles mechanical mutations

**Distribution**
- Claude Code plugin marketplace (`alexsmedile/spectacular`)
- Codex plugin marketplace
- CLI via curl one-liner installer

### Testing

- 7 test files, 216 asserts
- `init.test.sh`: 66 asserts
- `doctor.test.sh`: 48 asserts
- `migrate.test.sh`: 39 asserts
- `mutator.test.sh`: 59 asserts
- `pack.test.sh`: 48 asserts
- `conventions.test.sh`: 18 asserts
- `specs.test.sh`: 29 asserts

### Semver discipline starts here

- Patch (1.0.x): bug fixes, doc clarifications, internal refactors that preserve all public contracts
- Minor (1.x.0): new CLI verbs, new doctor areas, new registered doc types, new kits, new convention-pack rule categories
- Major (2.0.0): any change that breaks an existing workspace's ability to load, validate, or migrate

The `workspace_schema:` field tracks structural shape independently of CLI version. A v2.0.0 CLI with `workspace_schema: "1.0"` workspaces will run the migration chain on first invocation.

### Migration from v0.x

Nothing to do. Anyone on any v0.6.0+ CLI keeps working — v1.0.0 is the same surface with a fresh version label. If you're on a pre-v0.6.0 workspace, run `spectacular migrate` to bring your workspace schema forward.

### Pre-1.0 development arc

The history that produced v1.0.0 is preserved below as per-version entries. Each entry documents the capability that shipped in that version and is the place to look for "when did X land?" questions. Future releases (v1.0.1+, v1.1.0+) will continue this entry pattern.

---

## [0.7.5] — 2026-05-23

### Added — `spectacular status --since <date>` activity report

Closes task #57. Lists requests + canonical doc changes since a cutoff date.

- **Frontmatter-only scan** — reads `updated:` fields from `requests/*/PLAN.md`, `archive/*/PLAN.md`, and `.spectacular/*.md` (skipping snapshot files). Lexicographic YYYY-MM-DD compare against cutoff. No git involvement (matches Spectacular's "frontmatter is the signal layer" principle).
- **Date input formats**:
  - `--since 2026-05-20` — absolute YYYY-MM-DD
  - `--since 7d` / `--since 2w` / `--since 1m` — relative (days / weeks / months-approx-30d)
  - `--since yesterday` — keyword
  - `--since=<date>` — equals form also accepted
- **BSD + GNU date compatibility** — uses macOS `date -v-Nd` when available, falls back to GNU `date -d "N days ago"`.
- **Output groups requests by status** — archived / verified / review / active / planned, in lifecycle order. Each group only renders if non-empty. Empty bucket renders `(none)`.
- **Canonical docs section** lists `.spectacular/*.md` files (excluding snapshot @v* files) with their `updated:` date.

Exit codes: `0` success (including empty buckets); `1` outside a workspace; `2` missing or unparseable `--since` argument.

### Implementation note

Intercepts `--since` in the same `status` dispatch path as `--against-latest` (introduced v0.6.1). Both are CLI-mechanical flags on the otherwise-skill-owned `status` verb. Plain `spectacular status` (no flags) still routes to the skill stub for AI-driven briefing.

### Testing

- 7 test files, all green
- `init.test.sh` +1 scenario (status --since: 7 asserts covering absolute/relative/equals/empty/missing-arg/bad-format/outside-workspace) — 66 asserts (was 59)
- Plugin bumped 0.7.4 → 0.7.5

### Closes

- Task #57 — `spectacular status --since <date>`

---

## [0.7.4] — 2026-05-23

### Added — Doctor precondition for destructive verbs + VERIFY-as-tests convention

Two safety/process improvements that close out long-standing pending tasks (#55, #56).

**Doctor precondition on `archive` + `migrate` + `promote --archive`:**
- Before running destructive ops, CLI runs a scoped doctor check. Errors block (exit 1); warnings/info pass through silently.
- `archive` runs `workspace + specs + links` — full structural validation. Refuses if any errors found.
- `migrate` runs `links` only — workspace/specs would self-trigger on the v0.4-shape drift migrate is designed to repair. Dry-run skips the check entirely.
- `promote --archive` chain passes `--skip-doctor` through to the chained `cmd_archive`.
- Bypass: `--skip-doctor` flag on each verb. Caller responsibility to know what they're skipping.
- Shared helper `_doctor_precondition` for consistent behavior; thin wrappers `require_doctor_clean` (3-area) and `require_doctor_clean_links` (links-only).

**VERIFY-as-tests convention (documentation-only):**
- New `skills/spectacular/references/verify-tests.md` documents the `tests/verify/<slug>.test.sh` pattern: when to author one (multi-verb workflows not covered by `tests/cli/` suites), template, naming convention.
- No backfill: existing 7 `tests/cli/` suites already cover most of what archived VERIFY.md files would script. Pattern reserved for new requests that ship behavior not already in an area-level test suite.
- New VERIFY-tagging convention: `[x] mechanically verified` (auto-checkable) vs `[x] manually verified` (interactive UX, judgment calls). Lets future-agents grep to know which scenarios have a safety net.
- Wired into SKILL.md routing.

### Testing

- 7 test files, all green
- `mutator.test.sh` +2 scenarios (doctor precondition + clean passes through) — 59 asserts (was 49)
- Plugin bumped 0.7.3 → 0.7.4

### Closes

- Task #55 — Promote VERIFY scenarios → tests/verify/<slug>.sh (convention-doc only per scope cut)
- Task #56 — Doctor-as-precondition for destructive ops (archive + migrate + promote --archive)

---

## [0.7.3] — 2026-05-23

### Fixed — Unknown `convention_pack.mode` values fall back to `suggest` with info note

Previously, unrecognized `mode:` values in `config.yaml` (e.g. `strict`) silently passed through and were treated as scaffold-equivalent (warnings, not errors). Behavior diverged from the spec, which says: validate against `{suggest, scaffold, enforce}`, fall back to `suggest` if invalid, emit info note.

- New helper `config_pack_mode_raw()` returns the raw config value verbatim
- `config_pack_mode()` now validates and falls back to `suggest` for unknown values
- `check_conventions` emits info line when the raw value doesn't match the allow-list: "unknown convention_pack.mode 'strict' (valid: suggest / scaffold / enforce) — falling back to 'suggest'"
- Valid modes (`suggest`, `scaffold`, `enforce`) and absent `mode:` field both stay silent (latter defaults to `suggest`)

Found during S18 verification in `convention-pack-application` (task #62).

### Testing

- 7 test files, all green
- `pack.test.sh` +1 scenario (unknown mode fallback) — 48 asserts (was 44)
- Plugin bumped 0.7.2 → 0.7.3

---

## [0.7.2] — 2026-05-23

### Added — Roadmap-research findings applied (Outcome slot, Icebox, gate guardrails)

Tightens the v0.7.1 structured ROADMAP based on convergent advice from Pichler (GO Product Roadmap), Torres (Opportunity Solution Trees), Cagan (SVPG), Gilad (GIST), Productside, and beginner-tool onboarding patterns (GitHub Projects, Trello, Aha!). Closes the gap between v0.7.1 and what dominant roadmap frameworks recommend without ballooning the slot set.

**Structural changes:**
- **`Outcome:` slot added** between Phase and Scope-in. Required for `full` + `themed` tiers; absent for `vision` (Direction covers). One paragraph: "what business or product outcome does this version move?" Pichler/Torres/Cagan/Gilad convergent: the #1 missing slot in feature-list-shaped roadmaps. Forces goal-before-features discipline.
- **"Bucket list" → "Icebox"** rename across template, overrides doc, and live ROADMAP. Convergent dev-tool idiom (GitHub Projects, Pivotal Tracker, Linear; GIST's "Idea Bank"). Distinguishes "unbound idea" from "planned but vague" (which is what `vision`-tier blocks are for).

**Review gate additions** (all warnings/info — preserves recommend-not-enforce stance):
- **Date guards extended to themed/vision blocks** (gate check 12) — Cagan's "#1 sin." Scans entire block for `YYYY-MM-DD`, `Q[1-4] YYYY`, `MMM YYYY`. Warning, not error.
- **Outcome required by tier** (gate check 16) — full + themed must have Outcome; vision must not (Direction covers).
- **Full-tier row count** (gate check 17) — tiered: silent ≤7 (sweet spot per Cagan), info 8-10 ("consider demoting older versions to themed"), warning 11+ ("roadmap-as-backlog anti-pattern").
- **Scope-out push** (gate check 18) — when Scope-in has ≥4 items and Scope-out is empty, warning ("every item you add implies others you're not building" — Productside).

**Phase taxonomy extension:**
- **Meta-phase aliases** — `Phase:` field accepts both individual values (`mvp`, `release-prep`) AND coarser meta-phase values (`discover`, `build`, `release`). Coexist. Document the "start coarse, refine as work crystallizes" rule. Maps Cagan's discovery-vs-delivery split onto our 9-phase chain.

**Documentation-only additions** (no automation):
- **Beginner pattern** in `roadmap-overrides.md` — start at vision tier (one paragraph), graduate to themed when 2nd version exists, unlock full when first request links via `target_version:`. Mirrors GitHub Projects/Trello/Notion progressive-disclosure onboarding.
- **Icebox-promotion ritual** — explicit 4-step walk (pick item → choose version → choose tier → fill slots → delete from Icebox). Skill executes on `/spectacular roadmap` invocation. No new CLI verb; manual ritual is the point.

**Doctor extension:**
- `check_workspace` flags pre-v0.7.2 "Bucket list" heading in ROADMAP, suggests Icebox rename. Mechanical fix tag (sed substitution). Silent once renamed.

**Dogfood:** live `.spectacular/ROADMAP.md` snapshotted (`ROADMAP@v2.md`) then updated with Outcome paragraphs on v0.7.1 + v0.7.x + v0.11.x + v1.0.0, renamed Bucket list → Icebox, and added the 4 deferred research items to the Icebox.

### Out of scope (deferred per interview)

- Confidence rating per row (GIST/ProductBoard) — overlaps tier
- Audience field (Pichler internal-vs-external) — over-engineered for solo/small-team
- Opportunity-Solution-Tree as separate doc type (Torres) — heavyweight
- ICE/RICE scoring for icebox items (GIST signature) — convention-pack territory
- CLI verb for icebox promotion — manual via skill is enough
- `--beginner` flag — doc-only enough

### Testing

- 7 test files, all green
- doctor.test.sh: +1 scenario (Bucket-list → Icebox info) — 48 asserts (was 47)
- Plugin bumped 0.7.1 → 0.7.2

---

## [0.7.1] — 2026-05-23

### Added — Structured ROADMAP with precision tiers + 9-phase chain

ROADMAP graduates from `mode: freeform` to `mode: structured`. Each version block has a precision tier (`full | themed | vision`) that controls which slots are required, capturing the natural precision gradient — active work is detailed; long-term direction is intentionally fuzzy.

- **9-phase chain** for version progression: `intent → discover → prototype → spec-refine → mvp → iterate → test → release-prep → release`. Skill recommends the next phase; user can skip with reason. Skips recorded explicitly in `Phase:` field (e.g. `Phase: spec-refine (skipped: discover, prototype)`).
- **Prototype phase broadened** — any artifact produced to validate a decision against real tooling or against the user counts: data/schema drafts run through parsers, fake datasets tested against downstream scripts, mock API responses, ASCII wireframes, interactive mocks, sample CLI output. The artifact isn't the deliverable; the decision it informs is.
- **Three precision tiers**:
  - `full` — Status, Phase, Scope (in), Scope (out), Exit criteria, Linked requests. Use for active + near-term planned versions.
  - `themed` — Status, Phase, Themes (list), Exit criteria (directional). Use for mid-term (2-3 versions out).
  - `vision` — Status, Direction (free-text paragraph). Use for long-term + speculative.
- **Bucket list section** at the end of `ROADMAP.md` for ideas not yet tied to any version. Promoting an item: pick a version label + `Tier: vision` minimum.
- **`spectacular new --target-version <ver>`** flag — adds `target_version:` to PLAN.md frontmatter. Used by `spectacular roadmap refine` to autopopulate Linked requests in matching version blocks.
- **Doctor extension** — `check_workspace` flags pre-v0.7.1 freeform ROADMAP shape with info line (never error). Skill walk to migrate available via `spectacular roadmap grill`.

### Reference + template

- `skills/spectacular/references/roadmap-overrides.md` — tier-aware slot prompts (full/themed/vision variants), mini-refine patterns, vibe→spec rewrite tables, 15-check review gate
- `skills/spectacular/templates/roadmap/base.md` — structured template showing all 3 tiers + bucket list
- `skills/spectacular/references/doc-registry.md` — ROADMAP entry switched `mode: freeform` → `mode: structured`; mode reference table extended with `structured` and `reference` modes

### Dogfood

`.spectacular/ROADMAP.md` rewritten against new shape (snapshotted as `ROADMAP@v1.md`). Current versions tiered as: v0.7.1 full, v0.7.x/v0.11.x/v1.0.0 themed, v2.x/v3+ vision. Bucket list populated with previously-roadmapped ideas that no longer fit a specific version (hook automation, multi-agent orchestration, burndown viz, etc.).

### Testing

- 7 test files, all green
- mutator.test.sh: +1 scenario (target_version field) — 49 asserts (was 48)
- doctor.test.sh: +1 scenario (ROADMAP shape detection: old triggers info, new is silent) — 47 asserts (was 45)
- Plugin bumped 0.7.0 → 0.7.1

---

## [0.7.0] — 2026-05-23

### Added — CLI mutator verbs (skill orchestrates, CLI mutates)

Five CLI verbs replace skill-side manual file edits for the most common lifecycle mutations. Establishes the **mutation principle**: lifecycle changes go through CLI verbs; the skill orchestrates (reads, decides, communicates); the CLI mutates (atomically, deterministically, tested). Manual file edits remain available for edge cases but become the exception.

- **`spectacular new <slug> [--summary ...] [--status planned|active|review] [--priority low|medium|high]`** — scaffolds `.spectacular/requests/<slug>/PLAN.md` + `TASKS.md` from templates with frontmatter prefilled. Validates slug (kebab-case, max 64 chars). Refuses duplicates in `requests/` or `archive/`.
- **`spectacular promote <slug> [--to <state>] [--force] [--archive]`** — advances request through the lifecycle (planned → active → review → verified). Refuses backward transitions without `--force`. Mutates PLAN.md + TASKS.md `status:` + `updated:` atomically. `--archive` chains into archive after promoting to verified.
- **`spectacular snapshot <file> [--major]`** — snapshots a canonical doc to `<base>@v<N>.md`, bumps `version:` field, sets `updated:`. Refuses non-canonical files. Idempotent: compares body (frontmatter-stripped) against latest snapshot; exits cleanly if unchanged.
- **`spectacular archive <slug> [--force]`** — moves request to `.spectacular/archive/<slug>/`, rewrites every inbound `related:` link in other request files (`../<slug>/...` → `../../archive/<slug>/...`), sets PLAN frontmatter `status: archived` + `archived: <today>`. Refuses unless status is `verified` or `review`.
- **`spectacular touch <file>`** — sets frontmatter `updated:` to today. Idempotent. Refuses files without frontmatter.

### Shared infrastructure

- New frontmatter helpers (`fm_get`, `fm_set`, `fm_touch`, `fm_add_to_list`) — single shared implementation used by all 5 verbs. Single-source-of-truth for YAML rewriting.
- `is_canonical_doc` helper — gates `snapshot` to registered canonical files only.

### Skill instruction sync

To resist drift, both surfaces updated:

- `SKILL.md` routing table: each verb points to its CLI verb (not a reference doc — the CLI is the runtime). Top-level **mutation principle** stated explicitly.
- Reference docs rewritten to instruct CLI verb usage:
  - `references/new-request.md` → `spectacular new <slug>`
  - `references/archive.md` → `spectacular archive <slug>`
  - `references/lifecycle.md` → `spectacular promote <slug>` (state transitions)
  - `references/versioning.md` → `spectacular snapshot <file>`

### Absorbs

- Task #54 (archive --check / auto-rewrite related: paths) — same behavior shipped as the standard `archive` verb behavior.

### Testing

- 7 test files, all green; mutator.test.sh adds 48 asserts across 7 scenarios covering all 5 verbs + the promote-archive combo + --help flags
- init.test.sh scenario_10 updated — only status + remember remain as skill stubs (archive/new/snapshot/promote no longer stubs)

### Why this matters

Before v0.7.0, lifecycle mutations happened via the skill writing free-form file edits. Two agents would do the same thing differently; frontmatter parsers drifted; archive link-rewriting was easy to forget. The CLI verbs collapse this: one implementation, deterministic output, tested. Skill instruction sync ensures the skill actually calls the verbs instead of falling back to manual edits.

Remaining as skill flows (intentionally): `status` (briefing requires AI judgment), `remember` (memory entries require user-context distillation).

---

## [0.6.2] — 2026-05-23

### Added — Workspace migrations: registry pattern + judgment skill walk

Stage 2 of workspace-migrations (Stage 1 shipped in v0.6.1). Replaces the v0.6.1 hardcoded migration list with a proper registry + adds judgment-migration support via the skill.

- **Migration registry** — migrations now live as .md files under `skills/spectacular/references/migrations/v<from>-to-v<to>.md`. Frontmatter declares `(id, from, to, mechanical, reversible, apply-fn, affects)`; body documents detection rule, steps, rollback, validation. Full schema: `skills/spectacular/references/migrations-contract.md`.
- **CLI loader** — `cmd_migrate` now scans the registry, sorts by `from` semver, resolves `apply-fn` to bash functions in `cli/spectacular`. Adding a new migration = one .md file + one bash function (mechanical) OR one skill-walk section (judgment).
- **`spectacular migrate --to <ver>`** — migrate up to a specific schema version (default: latest).
- **`spectacular migrate --from <ver>`** — re-run starting from a specific schema (for repair).
- **`spectacular migrate --list`** — show all registered migrations with descriptions.
- **`/spectacular migrate` skill walk** — `skills/spectacular/references/migrate.md` defines the flow for judgment migrations: snapshot-before-edit on affected canonical docs, y/n/q per step, validation phase, audit-trail memory write. No judgment migrations ship in v0.6.2; the skill flow is scaffolding for future migrations that need it.
- **Chain validation in doctor** — `check_kits` validates the migration registry: no gaps in the chain, no duplicate `(from, to)` edges, every `apply-fn` resolves to a defined bash function, reversible migrations have reverse functions. Reports as part of the `kits` area; no new doctor area needed.
- **Downgrade refused** — `migrate --to <older-version>` exits non-zero. Bidirectional migration is explicitly out of scope.

### Why this matters

In v0.6.1 the migration logic was hardcoded directly in `cmd_migrate` with two if-branches. Easy to ship but doesn't scale: every new migration would have meant editing `cmd_migrate` plus the bash function. The registry pattern collapses this — maintainers add a .md spec + a bash function; the loader handles dispatch + ordering + chain validation automatically.

### Migration

For consumers, no action needed. v0.6.2 reads the same `workspace_schema:` field that v0.6.1 wrote. The behavior of `spectacular migrate` (no flags) is identical.

### Testing

- 6 test files, all green
- migrate.test.sh: 39 asserts (+12 for `--list`, `--to`, downgrade refusal)
- doctor.test.sh: 45 asserts (+3 for chain validation pass + gap detection)
- specs.test.sh: 29 asserts (unchanged)
- init.test.sh: 63 asserts (unchanged)
- conventions.test.sh: 18 asserts (unchanged)
- pack.test.sh: 44 asserts (unchanged)

### What's still in v0.6.1 territory (no changes here)

`workspace_schema:` field, `spectacular status --against-latest`, doctor flat-contract-docs support, v0.6+ scaffold suggestion — all from v0.6.1, unchanged. Stage 2 only touches the registry/loader/skill-walk surfaces.

---

## [0.6.1] — 2026-05-23

### Added — Workspace migrations + scaffold discoverability

Bridge the gap between Spectacular versions for already-scaffolded workspaces. Triggered by a real audit on the [Octopus](https://github.com/alexsmedile/octopus) repo where a v0.1.x-shape workspace had no discoverability path to v0.6 conventions.

- **`spectacular migrate`** — applies pending workspace-schema migrations to bring `.spectacular/` to the shape this CLI expects. Idempotent.
  - `--dry-run` lists planned migrations without writing
  - Two backfilled migrations: **v0.4 → v0.5** (rename `current/` → `specs/`, preserve contents) + **v0.5 → v0.6** (ensure `specs/` exists as always-set)
- **`workspace_schema:` field** — new top-level key in `config.yaml`. Records the structural version. `init` writes `"0.6"` on fresh workspaces; absent value treated as `"0.4"`.
- **`spectacular status --against-latest`** — one-line discoverability check. CLI-mechanical (no skill); the rest of `status` stays in the skill.
- **Flat contract docs in `specs/`** — top-level `.md` files in `specs/` are valid alongside per-capability subfolders. Pattern for projects whose primary truth is on-disk contracts (e.g. `SCHEMA-TASK.md`, `AXIS-MODEL.md`).
- **v0.6+ scaffold suggestion** — `doctor workspace` surfaces missing PRINCIPLES/ARCH/ROADMAP as one info line with the exact `init --with` command. Silent when all present.

Stage 2 (v0.6.2) will move the migration list from hardcoded into a `references/migrations/` registry, add `--to`/`--from` flags, and bring `/spectacular migrate` skill walk for judgment migrations.

---

## [0.6.0] — 2026-05-23

### Added — Public docs as a first-class surface

The `docs/` tree is now a first-class Spectacular surface, sibling to `.spectacular/` (workspace) and `specs/` (system truth). Flat, opinionated, renderer-agnostic — single `docs.yaml` nav manifest, page-level frontmatter, CLI scaffold + doctor validation + skill verbs for interactive authoring.

**Audience clarification:** the spec/doc boundary lives at the *folder* level, not the page level. `docs/` is for users + agents consuming the product; `specs/` is for devs + coding agents building it. Per-page `audience` would be ceremony — no such field.

- **`spectacular docs init [--minimal]`** — CLI subcommand. Scaffolds `docs/docs.yaml` + `index.md` + 3 default sections (`getting-started`, `guides`, `reference`) with placeholder pages. `--minimal` skips the sections and ships just docs.yaml + index. Idempotent; re-running fills empty stubs without overwriting content.
- **`spectacular doctor docs`** — substrate validation: docs.yaml parseable, declared pages exist, no orphan files, required frontmatter present (`title`, `description`, `section`, `status`, `updated`). Supports both sectioned trees (`docs/<section>/<page>.md`) and flat-tree extras (`docs/<slug>.md` registered via `docs.yaml extras:`).
- **`doctor docs --fix`** — mechanical: injects frontmatter stubs into pages missing them (delimiter or individual fields). Title defaults to slug, status to `draft`, updated to today.
- **Skill verbs** (registry-driven via the existing engine):
  - `spectacular docs new <page>` — scaffolds a page, prompts for section if omitted, updates docs.yaml
  - `spectacular docs new --section <name>` — declares a new section + scaffolds the dir
  - `spectacular docs review` — quality gate (same checks as doctor)
  - `spectacular docs status` — briefing scoped to docs/
- **`references/docs-contract.md`** — schema spec: folder shape, docs.yaml manifest, page frontmatter contract, validation rules, anti-patterns. Documents the spec-vs-doc boundary at folder level.
- **`references/docs-overrides.md`** — engine rules: `docs new` flow with section-prompt UX, `docs review` gate checks, `docs status` briefing format. Vibe→spec patterns for future `docs refine` (deferred to v2).
- **Doc registry** — `docs-manifest` and `docs-page` entries registered. Doc IDs registered count rises from 11 to 13.
- **`templates/docs/`** — `docs.yaml.tmpl`, `index.md.tmpl`, `page.md.tmpl` for the engine + CLI.
- **`tests/cli/docs.test.sh`** — 12 scenarios, 38 asserts covering init (default + minimal + idempotent), doctor (skip / clean / missing-declared / orphan / missing-frontmatter / --fix injection / extras), skill-verb refusal by CLI, help text.

### Dogfood

- **This repo's own `docs/`** migrated to v0.6.0 shape: `docs.yaml` authored with three sections + 5 existing pages registered as `extras:` (flat-tree preservation — moving files would break README links). Each of `workflow.md`, `commands.md`, `configuration.md`, `scaffold.md`, `troubleshooting.md` got proper frontmatter (title, description, section, status, since, updated). Doctor `docs` clean — 0 errors / 0 warnings.

### Out of scope (deferred to `public-docs-advanced` v0.6.2+)

- Renderer adapters (Mintlify / Docusaurus / Fumadocs / MkDocs export)
- Versioned docs snapshots (`docs/versioned/v<x.y.z>/`)
- `docs sync-from-spec` (spec ↔ doc sync flow)
- Convention-pack `docs-layout` rule category

These ship only when real-world demand surfaces (per the same activation-trigger pattern as `convention-pack-modules`).

### Quality

- 191 asserts pass across 5 test files (init, doctor, pack, specs, docs)
- Pre-commit version-consistency check green across 7 sources
- Doctor on this repo: 0 errors / 0 warnings on all 10 areas

---

## [0.5.0] — 2026-05-23

### Breaking — `current/` folder renamed to `specs/` + new `SPEC.md` index

The legacy `.spectacular/current/` folder convention is replaced by a two-part surface: `.spectacular/SPEC.md` (always-on, present-tense index of what's built) plus an optional `.spectacular/specs/` folder for per-capability detail when a SPEC.md bullet outgrows one line.

- **Why:** "current" is a temporal word, not a content word — agents kept mis-routing it as recency state. "spec" is the industry term and ties cleanly to what the layer actually holds. The TODO had this flagged since v0.3.0; v0.5.0 ships the migration.
- **`SPEC.md` is always-on**, scaffolded by every init. Per-capability `specs/<capability>/SPEC.md` files are **optional** — only break out when the bullet in SPEC.md outgrows one line. Small projects ship with one file.
- **Mechanical migration via doctor**: `spectacular doctor specs --fix` renames any legacy `current/` → `specs/`, preserving contents. Conflict case (both dirs present) raises an error and refuses auto-fix.

### Added

- **`SPEC.md` doc type** registered in `doc-registry.md` — uses the generic engine, mode `freeform`, snapshot-on-edit, scope `project-wide`.
- **`templates/spec/base.md`** — index-style template with "What this system is" + "Capabilities" sections.
- **`doc_spec()` writer** in `cli/spectacular` — scaffolds SPEC.md as part of the always-set.
- **Doctor `specs` area** — validates SPEC.md presence/parseability, specs/ dir presence, per-capability SPEC.md frontmatter, legacy current/ migration detection, conflict detection.
- **`spectacular spec` skill verb** (via existing registry-driven engine — `grill`, `refine`, `review`).
- **`references/spec-sync.md`** — renamed from `current-sync.md`, updated to drive SPEC.md bullet edits + specs/ creation during archive flow.
- **`tests/cli/specs.test.sh`** — 8 scenarios, 25 asserts covering fresh init, kit init, doctor on clean v0.5.0 workspace, legacy detection, mechanical migration, conflict refusal, per-capability validation, re-init non-destructive.

### Changed

- **Always-set bumped from 5 → 6 files**: PRD.md, **SPEC.md**, config.yaml, `<agents-file>`, requests/, **specs/** (replacing `current/`).
- **`doc_agents()` template rewrite** — new build's `.spectacular/AGENTS.md` now leads with the four-layer model (Intent/Truth/Work/Memory), documents two-layer task tracking, and consistently uses `SPEC.md` + `specs/<capability>/SPEC.md` references throughout. Mirrors landed in `templates/agents/base.md`.
- **Doctor re-run dispatcher fix** — `conventions` area was missing from the `--fix` re-run loop, causing pack-driven gitignore fixes to not refresh detection state. Both `conventions` and the new `specs` area are now in both dispatcher loops.
- All `current/<capability>` references swept through skill references, `.spectacular/` live workspace, `docs/`, `README.md`, `CLAUDE.md` — replaced with `specs/<capability>/SPEC.md` per the new convention. Migration callouts retained where users upgrading from v0.4.x need them.
- Root `AGENTS.md` updated to point at `.spectacular/SPEC.md` for system-truth queries and clarify per-request loading discipline.
- **codex-plugin `longDescription`** — replaced "current truth" framing with "system spec" framing to match the new convention.

### Migration from v0.4.x

```bash
spectacular doctor specs        # detect legacy current/
spectacular doctor specs --fix  # rename current/ → specs/, preserve contents
spectacular init                # fills in SPEC.md if missing
```

Workspaces with both `current/` and `specs/` present require manual merge — doctor refuses auto-fix.

---

## [0.4.0] — 2026-05-23

### Added — Convention Pack system

A new opt-in layer for declaring repo-shape opinions: naming rules, folder taxonomy, required root files, gitignore defaults, file-placement rules, project-type scaffolds. Packs are mini-skills (folder + `pack.md` + `templates/` + `references/`), distributable via four scope locations (project-local → user → app-store → bundled, in precedence order).

- **`packs-contract.md`** — full schema spec covering 6 rule categories (`naming`, `taxonomy`, `root-files`, `gitignore`, `file-placement`, `project-types`), pack folder shape, 4-tier scope precedence, single-pack-only v1 (multi-pack composition in [convention-pack-modules](.spectacular/requests/convention-pack-modules/)).
- **`pack-overrides.md`** — pack-specific grill rules: 7 slot prompts, source-ingestion (`--from`), 8 mini-refine patterns, vibe→spec rewrite tables, review gate checks, reserved pack-id enforcement.
- **`spectacular pack` CLI subcommand** — `list` (shows all 4 scope locations), `install <name> [--from <path>]` (copies pack to `~/.spectacular/packs/<name>/`), `remove <name> [--force]` (refuses bundled/app-store/project-local without `--force`), `show <name>` (prints scope + frontmatter).
- **`config.yaml` `convention_pack:` block** — declares active pack per repo with `source`, `mode` (suggest|scaffold|enforce), and reserved `overrides` field.
- **Init wiring** — when a pack is declared with `mode: scaffold` or `enforce`, init appends the pack's `gitignore.always-add` entries (deduplicated). Always-set wins on conflicts; pack scaffold is purely additive.
- **Doctor `conventions` area** — validates pack source resolves; in `scaffold` mode flags gitignore drift as warnings (exit 1); in `enforce` mode escalates to errors (exit 2); `suggest` mode skips drift checks. `--fix` mechanically appends missing pack-declared gitignore entries.
- **Bundled `minimal` pack** — ships at `skills/spectacular/templates/packs/minimal/`. Enforces only README contract + `.gitignore` baseline. The implicit default when no other pack is declared.
- **App-store `alex-default` pack** — ships at `packs/alex-default/`. Fully-opinionated pack encoding kebab-case naming with role suffixes, mono-collection detection, 8 project-type scaffolds, full `.gitignore` baseline + language-specific blocks for Python/Node/Go.
- **`tests/cli/pack.test.sh`** — 12 scenarios, 44 asserts covering list/install/remove/show, init wiring, doctor across all 3 modes, mechanical fix repair, scope precedence, `--from` install, error paths.

### Added — Workflow conventions

- **Two-layer task tracking convention** documented in `.spectacular/AGENTS.md` and skill's `SKILL.md`: harness `TaskCreate`/`TaskUpdate` = ephemeral session micro-tracker (drives CLI live progress UI); on-disk `requests/<slug>/TASKS.md` = persistent milestone blocks. Anti-pattern: one-for-one duplication. Resolves the recurring "task tools haven't been used" warning by giving it a real role.

### Changed

- `references/init-workflow.md` — § "Convention packs (v0.4.0+)" added with 3-mode behavior table + 4-tier precedence table.
- `references/doctor.md` — `conventions` check area added; severity-per-mode table.
- `references/new-request.md` — `artifacts/` directory consults active pack's `file-placement.request-artifacts:` rule.
- `references/doc-registry.md` — `convention-pack` entry registered with new `scope: user` value (packs live under `$HOME`, not per-project).
- `cli/spectacular` — new `SCRIPT_DIR` constant; `pack` subcommand sibling of `init` + `doctor`; `check_conventions()` doctor area; `pack_apply_scaffold()` init hook; `config_pack_source/mode()` awk parsers.
- `cli/spectacular@v0.3.1` snapshot captured before v0.4.0 work.

### Convention pack chain (3 requests)

- `convention-pack-schema` (verified) — locks the schema + ships bundled `minimal`
- `convention-pack-fabricator` (review) — pack-overrides + alex-default dogfood; live grill walkthroughs remain for full signoff
- `convention-pack-application` (review) — CLI + init + doctor wiring; live three-mode + cross-machine scenarios remain

### Planned

- **`convention-pack-modules`** (planned, v0.5.0+) — split monolithic packs into composable rule-category modules. Stays planned until composition pain surfaces from v1 use.

---

## [0.3.1] — 2026-05-23

### Added

- **`spectacular doctor`** — environment/infrastructure self-check. CLI detects substrate drift across 7 areas (`skill`, `workspace`, `frontmatter`, `snapshots`, `links`, `lifecycle`, `kits`); skill handles judgment-requiring repairs via `/spectacular doctor --fix`. Severity model: ✅ pass / ⚠️ warning / ❌ error / ℹ️ info. Exit codes: 0 clean, 1 warnings, 2 errors. `--format text|json` for human + machine consumption. `--fix` applies content-free mechanical repairs (`.gitignore` append, missing always-set dirs, dangling symlinks, re-stub missing canonical files).
- **Skill-invoked doctor subsets** — `status.md`, `grill.md`, `onboarding.md`, `lifecycle.md` auto-run scoped doctor checks when substrate failures block their operation; findings surface inline.
- **`references/doctor.md`** — full spec: check definitions per area, severity model, report format (text + JSON), CLI vs skill split, repair-flow walkthrough with worked examples, anti-patterns.
- **`tests/cli/doctor.test.sh`** — 11 scenarios, 33 asserts covering detect, mechanical fix, scoped areas, JSON output, regression.

### Changed

- Smart-init's diagnostic placeholder (`"run diagnostics via spectacular doctor once available"`) replaced with explicit area pointer: `"run \`spectacular doctor frontmatter\` for details"`.
- `references/SKILL.md` routing extended with doctor triggers + skill-invoked-subset note.

### Workspace cleanup (from doctor dogfood)

- Fixed broken `related:` paths in `prd-craft/PLAN.md` + `repo-conventions/PLAN.md` (file-relative paths instead of repo-root-relative).
- Parked `cli-bootstrap` request (`active → planned`) since smart-init shipped at v0.3.0.
- Documented `PRD@v1.1.md` snapshot gap in `DECISIONS.md` (version bump skipped during canonical-docs-rework; no git trace exists).

---

## [0.3.0] — 2026-05-22

### Changed (breaking)

- **`spectacular init` scaffolds only what the project needs**, not all 7 root docs. Default = 5-file always-set (`PRD.md`, `config.yaml`, `<agents-file>`, `requests/`, `current/`). Extra docs come from the selected kit's `triggers-docs.always` and `triggers-docs.suggested` lists, or via explicit `--with <doc1,doc2>` flag.

### Added

- **PRD 8-slot shape** — base PRD template + all 5 kits versioned to v1.1 (Vision / Problem / Target users / Deliverable / Goals & success criteria / Non-goals / Constraints / First milestone).
- **Doc-writer engine** — generic `grill.md` / `refine.md` / `review.md` references that consume `doc-registry.md` to handle any doc type. PRD-specific logic moved to `prd-overrides.md`. Per-doc overrides for PLAN + TASKS. New templates for plan/tasks/principles/architecture/roadmap/stack/agents/decisions.
- **Kits as diff-only plugins** — kits now declare `adds-slots`, `modifies-slots`, `triggers-docs.always`, `triggers-docs.suggested` via frontmatter. Single-kit-only in v1. Documented in `references/kits-contract.md`. All 5 bundled kits (`blank`, `coding`, `content`, `product`, `research`) refactored to diff format.
- **Verification convention** — `references/verification.md` formalizes when VERIFY.md is needed (2-of-6 rule) vs folded into PLAN § Validation or TASKS § Verification. "Opt-in" refers to file scaffolding only; verification itself is mandatory before any `verified` transition.
- **CLI flags** — `--kit <name>`, `--with <doc1,doc2>`, `--minimal` flags for `spectacular init`. Backwards-compatible with existing `--name`, `--summary`, `--agents-file`, `--global`, `--update`, `-i` flags.
- **Pre-flight non-overwrite** — init detects existing/empty/malformed files and skips/fills/diagnoses without ever overwriting. Generic "run diagnostics via `spectacular doctor`" message emitted for malformed cases (to be replaced when doctor ships).
- **Test harness** — `tests/run.sh` discovers and runs `tests/**/*.test.sh`. First test suite at `tests/cli/init.test.sh` covers 6 smart-init scenarios (41 asserts).

### Anti-patterns formalized

- Per-doc skills (e.g. one skill per doc type) — superseded by all-in-one `/spectacular` with registry-driven verbs.
- `--force` flag — explicitly rejected. Re-init never overwrites; to regenerate a stub, delete the file first.
- Project-type inference in init — bare init uses `blank` kit unconditionally. Auto-detection deferred to v2.
- Skipping verification because "VERIFY.md is optional" — opt-in refers to the file, not the practice. Every request reaches `verified` through some artifact (VERIFY.md > TASKS § Verification > PLAN § Validation).

### Upgrading from v0.2.x

`spectacular init` on a v0.2.x workspace is safe — pre-flight skips every existing file. To add docs the v0.3.0 init no longer scaffolds by default (PRINCIPLES, ARCHITECTURE, etc.), run `spectacular init --with <docs>` or `spectacular init --kit <kit>`.

CLI snapshot preserved at `cli/spectacular@v0.2.0` for reference.

---

## [0.2.0] — 2026-05-21

### Changed

- **Canonical docs split.** The original 896-line `.spectacular/PRD.md` was split into four focused root docs:
  - `PRD.md` — product intent (now 121 lines, 6-slot shape: problem / who / success / non-goals / constraints / milestone)
  - `PRINCIPLES.md` — 8 operating principles, each with a runtime enforcement hook
  - `ARCHITECTURE.md` — workspace structure, frontmatter conventions, lifecycle, versioning
  - `ROADMAP.md` — versioned future work (v1 / v2 / v3+)
- **AGENTS.md rewritten** as the in-folder onboarding doc for any agent landing in `.spectacular/`. Authoritative source for per-task context loading rules.
- **PLAN.md template upgraded** to the 7-slot decomposition: goal / why / constraints / milestones / tasks / dependencies / validation / deliverables.
- **CLI scaffolds the full 7-doc root layer** on every `spectacular init`. PRD stub uses the new 6-slot shape; AGENTS stub uses the onboarding shape; new PRINCIPLES / ARCHITECTURE / ROADMAP stubs included.
- **Skill references aligned** with the new doc set — `status.md`, `onboarding.md`, `init-workflow.md`, `scaffold-reference.md`, `new-request.md`, `versioning.md`, and SKILL.md state-awareness all updated.
- **Project docs aligned** — README, CLAUDE.md, `docs/scaffold.md`, `docs/configuration.md`, `docs/commands.md`, `docs/troubleshooting.md`, `docs/workflow.md` all reflect the new 7-doc canonical set.

### Added

- `prd / prd refine / prd review` skill triggers — interactive PRD building with 5 kits (coding / product / content / research / blank), vibe→spec refine patterns, and a pass/fail quality gate.
- `requests/prd-craft/` and `requests/canonical-docs-rework/` — tracking artifacts for the v0.2.0 work.
- Snapshot history preserved: `PRD@v1.3.md`, `AGENTS@v1.0.md`, request PLAN/TASKS `@v1.0.md` and `@v1.1.md`.

### Anti-patterns formalized

- **Never create `requests/<slug>/PRD.md`.** Product intent is project-wide and lives at `.spectacular/PRD.md`. Per-request folders use `PLAN.md` + `TASKS.md` only.

---

## [0.1.1] — 2026-05-11

### Fixed

- `cli/spectacular`: progress log in `download_and_install_skill` redirected to stderr — was polluting `skills.lock` with the "Fetching skill..." line
- `cli/spectacular`: skill ref resolution now tries releases API first, tags API second, `main` as final fallback — previously fell straight to `main` when no GitHub release existed

---

## [0.1.0] — 2026-05-11

### Added

- `cli/spectacular` — Bash CLI binary (`spectacular init`) with zero-prompt default, `-i` interactive mode, and flags: `--name`, `--summary`, `--agents-file`, `--global`, `--update`
- `cli/install.sh` — curl-installable installer; places binary at `~/.local/bin/spectacular`
- `skills/spectacular/` — `/spectacular` Claude Code slash command; lean SKILL.md orchestrator routing to `references/` subdocs
- `.spectacular/config.yaml` — `agents.file` key (override primary agents file) and `agents.tool_overrides` map (per-tool supplementary files, e.g. `claude: CLAUDE.md`)
- `skills.lock` — CLI-written lockfile tracking installed skill ref, SHA, and source URL
- Skill install targets both `.agents/skills/spectacular/` (source) and `.claude/skills/spectacular/` (symlink) for multi-tool compatibility
- `CLAUDE.md` — project guidance for Claude Code
