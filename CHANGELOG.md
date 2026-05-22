# Changelog

All notable changes are documented here.
Format: [Keep a Changelog](https://keepachangelog.com/en/1.0.0/)

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
