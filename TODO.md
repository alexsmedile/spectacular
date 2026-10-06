---
schema: make-a-change/todo/v1
extensions:
  - "octopus:all"
---

# TODO

- follow-up requests suggestions
- catch new plans/ideas suggest wiring plan in spectacular
- add numbering to requests so they appear in order in requests/ folder. number is also id. — **see DECISIONS 2026-05-23 (rejected for v1)**
- when phase ends, verify files, then auto-continue, user expects that spectacular takes task till the end with no interruption
- implement interview mode like grill-me?
- subagents management
- tak inspo from gsd, superpowers
- how do we integrate with octopus tasks/sessions?
- workflow preview (ascii diagram)
- headless-verifier / ui-checker: automated headless visual verification (e.g. Playwright multi-viewport screenshots + axe-core a11y scan) captured directly into `.spectacular/evidence/`
- github-actions: event-driven / CI triggers (e.g. PR open -> `spectacular mission check` & `verify.sh acceptance`, tag `needs-contract` -> auto-draft Proposal/Contract)
- loop telemetry & metrics: capture attempt counts, token/spend budget metrics, test duration, and intermediate failure traces in Evidence records

## Future Companion Skills Ideas (Full-Lifecycle Software Factory)

- **`ui-craftsman`** (Frontend & Design Systems): Accessible (WCAG AA) components, design tokens (Tailwind/CSS vars), state machines (Zustand/XState), optimistic UI, form validation (Zod), and micro-interactions. Direct handoff from `rapid-prototyping` Level 3.
- **`api-design`** (Interface & Wire Contracts): OpenAPI 3.1, gRPC / Protobuf, GraphQL, RFC 7807 Problem Details error envelopes, cursor pagination, and webhook signing/retry specs. Bridges `system-architecture` and `data-modeling`.
- **`test-engineer` / `adversarial-qa`** (Proof & Falsification): TDD failure-test-first enforcement, dedicated Hunter adversarial passes (concurrency races, auth edge cases, injection), property-based testing, and contract fixtures.
- **`cloud-ops` / `infra-architect`** (Deployment & Operations): Multi-stage Dockerfiles, Terraform/OpenTofu IaC, hardened GitHub Actions CI/CD matrix, and OpenTelemetry observability wiring.

## Matrix proposal loop — standalone project name ideas

- `align-loop`
- `design-relay`
- `choicecraft`
- `progressive-canvas`
- `design-trident`

## Reserved Wayfinding entities

- `PRT` — reserved for prototype artifacts linked from spikes. Do not create a standalone `prototypes/` collection until artifact storage, retention, and promotion have a concrete consumer.
- `TSK` — reserved for durable task identity. Do not assign IDs or create a `tasks/` collection until its relationship with request `TASKS.md` and harness/session tasks is deliberately designed.

## Lifecycle contract follow-ups

- [ ] Review the explicit lifecycle-migration preview for legacy memories `M1`, `M2`, and `M3`; only add `status: active` after user confirmation. Until then, they remain legacy-readable and `spectacular doctor` intentionally reports three warnings.

- [ ] Install Pageworks and run its audit against the lifecycle-related public documentation (`README.md`, `docs/commands.md`, `docs/workflow.md`, and `docs/scaffold.md`). The files were manually aligned, but have not received Pageworks validation.

- [ ] Decide whether `TODO.md` should remain local-only and gitignored. If so, mirror every durable team-visible follow-up into a tracked Spectacular idea, question, roadmap entry, or request so pushing a branch cannot lose it.
- [ ] Confirm that ignored canonical snapshots are intentionally local recovery artifacts. If any snapshot must travel with a review or release, define an explicit export/evidence mechanism rather than silently force-adding `_snapshots/`.

## Next session — distribution and onboarding (2026-08-18)

Priority order. The CLI stays a separate install; the Skill travels via the Agent
Plugins standard and must onboard the user to the binary rather than pretend it is
present. See the session decision below on why TS scripts were not chosen.

**P1 — DONE 2026-08-18 (`8721d4e`) — CLI detection and onboarding in `SKILL.md`.** The Skill has no way to tell
whether the binary exists, so it cannot say "not installed, run this." Detect it,
and when absent give the exact install command and state plainly what is
unavailable: no fingerprints, no atomic writes, no governed mutation. A silent
fallback is worse than a refusal — a user who believes they have freeze points and
does not is worse off than one who knows. This is the deliverable that makes the
split honest.

**P2 — DONE 2026-08-18 (`8721d4e`) — `docs/installation.md`.** One page covering both halves and both update
paths: Skill via Agent Plugins / Claude marketplace / `npx skills`, CLI via
`install/install.sh`. The install surface already exists and is complete
(`install | update | select | rollback | uninstall | recover`, with `--version`
and platform flags) — this is documentation, not engineering. Link from `README.md`
and `docs/README.md`.

**P3 — Fallback scripts in the Skill bundle.** Ship read-only helpers at
`skills/spectacular/scripts/` so a Skill-only host is usable rather than blind,
and point step 0 at them when `spectacular --version` fails.

Three tiers, best available wins: **CLI → TS helpers → shell floor.**

*Shell floor — no toolchain, runs anywhere the Skill lands:*

- `orient.sh` — locate `.spectacular/`, list Missions with `status:`, name the live one
- `where.sh <ref>` — resolve a ref to its record path (pure globbing)
- `doctor.sh` — is the CLI present, is this a workspace, is the tree clean

*TypeScript — needs Node, but parses correctly:*

- `show.ts <ref>` — parse frontmatter, print state and outcome
- `check.ts <ref>` — structural validation: required fields, referenced paths exist
- `graph.ts <mission>` — Objective DAG, `after:` vs `after_interface:`, what is startable
- `gaps.ts` — every open Gap across Contracts with its `blocked_on:`
- `archive.ts` — retired records with `resolved_by:` and authorization
- `handoff-verify.ts` — re-check a Handoff's recorded commit and tree against the repo

**Why the split.** Shell may only touch flat fields. Spectacular's frontmatter has
nested blocks (`activation:`, `authority.operator[]`), block scalars, and lists, and
grepping those is the exact defect already hit once: the Gap rewrite located
`blocked_on:` by walking lines without tracking block-scalar depth and would have
spliced a resolution into the middle of a sentence. Anything that traverses nested
YAML gets a real parser.

**Hard boundary: read and report only.** No writes, no fingerprints, no
transactions. `handoff-verify.ts` is the edge worth watching — *verifying* an
existing fingerprint is legal, *producing* one is not. `mission start`,
`objective promote|finish`, `run start`, `review record`, `handoff record`,
`mission complete`, and `contract amend` stay CLI-only.

**Build order.** Shell floor first: no toolchain, immediate value, zero risk. Then
`show.ts` + `check.ts` as a deliberate probe of the TS path — if parsing and
validating records in TS is clean at this scale, the rewrite in the open decision
below is plausible; if it is already awkward, that is learned cheaply.

**P4 — DONE 2026-08-18 (M13) — `plugin.json` version drift.** Checked in `test/verify.sh` `manifest_checks` so the gate catches any mismatch against `VERSION`.

**P5 — verify the Antigravity install path.** Moving Antigravity data under
`extensions["com.google.antigravity"]` follows the current spec, but skizl warns
that convergence is expected rather than confirmed. Untested for this repo.

### OPEN DECISION — owner to settle: Go CLI + Skill, or a TS-only standalone skill

Not decided. Two coherent products, and the choice is about what Spectacular *is*,
not about what is technically possible.

**A — Split (current).** Go CLI installed separately; Skill travels via Agent
Plugins. Ships today, no rewrite, keeps a single static 5.1MB binary with no
runtime dependency. Cost: a conformant host (Cursor, Copilot, ChatGPT) gets
guidance without governed execution, and the user must install a binary.

**B — Replace with TS.** Rewrite the mechanical layer in TypeScript so
`skills/spectacular/` is a genuinely standalone folder: `SKILL.md` plus `scripts/`
under `npx`, installable on all six Agent Plugins hosts as one unit, no separate
install step. Cost: ~7,700 lines to port, atomicity and transactions re-proven in
Node, and a Node toolchain becomes a dependency.

**Correction to an earlier note in this file's history:** the claim that TS *cannot*
do what Go does was wrong. Node can do atomic writes, SHA-256, YAML, and typed
validation. The real costs are rewrite scope and byte-identical serialization — not
capability. Note also that much of the Go is the proof harness, not the runtime:
`internal/missionbundle` is 4,538 product lines against 22 test files.

**Check this before choosing B:** whether a Node serializer can reproduce the
fingerprints already recorded in `.spectacular/`. If existing record hashes cannot
be reproduced byte-for-byte, every freeze point in the archive breaks, and B becomes
a migration with a compatibility boundary rather than a clean rewrite. This is the
single question that decides whether B is cheap or expensive.

### Context: why scripts *alongside* Go was rejected

The mechanical layer is ~9,000 lines of Go producing a 5.1MB binary, and what it
owns is exactly what an agent does unreliably: atomic transactions proven by
fault injection at every write boundary, SHA-256 freeze points, typed validation
against a frozen schema, and byte-identical canonical writes. Porting it to
scripts inside `skills/spectacular/scripts/` would mean either reimplementing all
of that in a second language and keeping the two byte-identical, or quietly
dropping the guarantees. Both are worse than a separate install.

The Agent Plugins standard covers `skills/` and `mcp.json` only, so `cmd/`,
`install/`, and `.spectacular/` were never portable. The split already exists —
the only choice is whether it is stated or accidental.

## Open after v2.4.0 (2026-08-18)

Every Mission M1–M12 is completed and v2.4.0 is tagged. What remains is proposal-stage
work and one contract Gap. Nothing here blocks using the product.

### Undecided proposals

- [ ] **P8 — mechanical Git branch guardrail.** `mission start` still accepts activation
  on `main`/`master`, so a multi-step Mission can destroy the review and isolation
  boundary it depends on. Demonstrated problem, no implementation. The highest-value
  open item — this release needed a manual `--allow-main` override to cut, which is the
  same hole from the other side.
- [ ] **P10 — preparation judgment checkpoint.** A Mission can freeze and activate with
  nobody having asked whether the approach is understood or the slice correctly sized.
  Five design decisions are owner-accepted; they bind only when a Mission freezes them.
- [ ] **P5 — a trace of the last preflight, with decay.** Three of P5's four directions
  shipped (drift-aimed audits, `fallbacks:`/`invalidated_if:`, `after_interface:`); only
  the preflight trace remains. It is the cheapest of the four and unblocked. P5 stays live
  and annotated in place rather than retired — see `D11-proposal-retirement`.

### Proposal lifecycle — settled

`D11-proposal-retirement` defines what happens when a Proposal is done: it gains
`resolved_by:` naming the Mission that absorbed it, moves to
`.spectacular/archive/proposals/` with `archive_authorization:` and a fingerprint, and
leaves `proposals/` holding only open questions. P6, P7, and P9 were retired under it.

- [ ] P1–P4 read `accepted` and shipped long ago, but predate `resolved_by:` and were not
  in D11's targets. Retiring them needs their absorbing Missions identified and a Decision
  authorizing the move. Low urgency; they are correct where they are, just not yet retired.

### Open Contract Gap

- [ ] `concurrent-run-timelines` on `CC-projsurf` stays open with its original reason:
  timelines across concurrently live Runs are not renderable because the Run model
  permits exactly one live Run. Genuinely blocked on a Run-model change touching run
  start, fingerprints, atomicity, and review boundaries. Correct to leave open — it is a
  stated limit, not a defect.

## Multi-Agent Architecture & Intake Enhancements (Post-v2.4.0)

### 1. `write-prd` / Intake PRD Authoring Skill
- [ ] Dedicated intake skill (or mode) for 0-to-1 project kickoff.

- [ ] Uses a high-reasoning (thinking) model in a single-shot generation to draft a comprehensive starter PRD into a temporary file (`scratch/PRD.tmp.md` or `.spectacular/PRD.tmp.md`).

- [ ] Captures all 8 foundational v1 dimensions in one pass:
  1. Vision & Problem statement
  2. Target Users & Personas (Primary vs Secondary)
  3. Deliverables & System Layers
  4. Stack & Constraints
  5. Strict Non-Goals (Scope boundaries)
  6. Ubiquitous Language & Core State Machines
  7. Milestone Arc & Phasing
  8. Measurable Success Criteria
- [ ] Acts as ephemeral launchpad directly digested by Spectacular's One-Shot Genesis into Core Anchors (`PROJECT.md`, `STACK.md`, `ARCHITECTURE.md`), On-Demand Anchors (`PRODUCT.md`, `VOCABULARY.md`), and `M1-bootstrap` claims.

### 2. Model Profile Abstractions for Roles (`reasoning`, `fast-code`, `strict-verifier`)
- [ ] Define abstract semantic model profiles in Mission/Objective/Run frontmatter (`profile: fast-code | reasoning | strict-verifier`):
  - `reasoning`: Orchestration, Genesis, Campaign planning, FROST audits (e.g. Claude Sonnet w/ Thinking, Gemini Pro, o1/o3).
  - `fast-code`: Bounded worker execution, routine file edits, test-running (e.g. Gemini Flash, Claude Haiku, GPT-4o-mini).
  - `strict-verifier`: Adversarial validation, clean context, strict instruction following.
- [ ] Teach the Skill to map abstract profiles to host runtime tool invocations (e.g. `invoke_subagent` model parameters: `pro` vs `flash` in Antigravity; `sonnet` vs `haiku` in Claude; `o3-mini` vs `gpt-4o-mini` in OpenAI; droids in Goose).

### 3. Dual-Path Independent Review Workflow
- [ ] Provide two explicit paths at the review gate when `review level: independent` is required:
  - **Path A (In-Harness Subagent)**: Skill automatically dispatches a fresh subagent with clean context, binding Git commit & tree SHA and the FROST review checklist, writing `ReviewDraft` directly to disk.
  - **Path B (External Model / Human Handoff)**: Skill generates a self-contained, copy-pasteable Markdown review prompt (`.spectacular/missions/<slug>/reviews/review-handoff-prompt.md` or printed to chat) for testing in ChatGPT, DeepSeek, external Claude session, or peer developer, returning the output into `spectacular review record <mission-ref> -`.

### 4. Handoff Directory Architecture Resolution
- [ ] Formalize directory locations for delegation and review handoffs:
  - **Canonical Governed Handoffs**: `.spectacular/missions/<slug>/handoffs/H<n>-<shortkey>-<slug>.md` (created via `spectacular handoff record`, binding Git commit/tree and `asserted`/`assumed` lists).
  - **Review Handoffs & Prompts**: `.spectacular/missions/<slug>/reviews/` (e.g. `review-handoff-prompt.md` and recorded `ReviewDraft` files).
  - **Ephemeral Scratch Drafts**: `scratch/` (temporary, non-governed workspace files).

### 5. Headless Browser & User Testing Evidence Recipes
- [ ] Reference doc (`skills/spectacular/references/testing-recipes.md`) documenting headless browser validation (Playwright / Puppeteer / computer-use agents).

- [ ] Standardized conventions for depositing attributable UI QA receipts (screenshots, interaction traces, HAR logs, CLI exit codes) into `.spectacular/missions/<slug>/evidence/`.

### 6. Autopilot Architecture Upgrade
- [ ] Dedicated review and design session for long-running unattended autonomy:
  - Multi-objective autonomous loops
  - Hard resource enforcement and cancellation mechanisms
  - Heartbeat/timeout monitors
  - Automated fallback and rollback points upon stop condition triggers

## Configurable workspace folder name

Currently the workspace directory is hardcoded as `.spectacular/`. Let users configure it at `init` time (and persist the choice) so it fits the host project's naming preference.

### Options to support

- **Default** — `.spectacular/` (status quo; brand-visible, unambiguous)
- **Project-scoped** — `.<project-name>/` (e.g. `.pageworks/`, `.octopus/`) — matches the host repo identity, feels native
- **Generic** — `.specs/` (terse, neutral, reads well in `ls -la`)
- **Custom** — arbitrary user-supplied name (validate: leading dot, lowercase, no slashes)

### CLI surface

```
spectacular init --workspace-dir .specs
spectacular init --workspace-dir .myproject
spectacular init                                # defaults to .spectacular/
spectacular init -i                             # interactive prompt offers the 3 presets + custom
```

### Where the choice gets persisted

- `.<chosen>/config.yaml` gets a top-level `workspace_dir: .specs` field so tooling can find itself
- Optionally a tiny pointer file at repo root (`.spectacular-pointer` or similar) so the skill/CLI can discover the workspace regardless of name — OR the CLI/skill scans for any `*/config.yaml` matching the spectacular schema
- `spectacular doctor` learns to detect mis-named or orphaned workspaces

### Open questions

- Discovery: pointer file vs schema-sniffing — pointer is simpler, sniffing avoids extra file
- Migration: `spectacular migrate-workspace --to .specs` command to rename existing workspaces?
- Multiple workspaces in one repo? (e.g. monorepo where each package has its own.) Probably out of scope for v1 of this feature
- Should packs/kits be able to *recommend* a workspace name? (e.g. a `pageworks` pack defaults to `.pageworks/`)
- Doctor + onboarding refs all currently say ".spectacular/" — need a templating pass so docs reflect the chosen name

## Feedback loop (prototyping mode)

Spectacular is still in prototyping. Avoid the word "evals" here — it implies benchmarks, accuracy scores, automated grading. That's not what this is. What we need is a **deliberate human-feedback loop**: the skill (or a sub-mode) decides what's worth probing, drafts proposals, runs them past the user, and captures the response as durable signal.

This is a **strategy for acquiring feedback, knowledge, insights, and use-case validation** — not a verification harness and not a benchmark. Treat it as exploratory.

### How this differs from VERIFY.md

- **VERIFY.md** — request-scoped, confirmatory: "we said X, did we ship X?" Closed-ended, terminates at `verified`, archived with the request.
- **Feedback loop** — system-scoped, exploratory: "we shipped X, was X the right thing to ship?" Open-ended, compounds across sessions, lives outside any single request.
- They probe orthogonal axes: **conformance to plan** vs **fitness for purpose**. VERIFY can pass while feedback reveals we built the wrong thing.

### Shape of the loop

1. **Pick a target** — a recent change, a fuzzy convention, an untested edge of the substrate, or a hypothesis ("does grill-each actually feel better than grill-wide on PRDs > 10 slots?").
2. **Craft a proposal** — concrete: a scenario, two variants, the question being asked, the expected signal. Not "what do you think of grill?" but "here are PRD A (grilled wide) and PRD B (grilled each-slot) on the same input — which lands closer to what you wanted, and why?"
3. **Ask the user** — surface the proposal as a structured question (AskUserQuestion with previews when comparing artifacts; free-form otherwise).
4. **Capture the response** — write to `.spectacular/feedback/<date>-<slug>.md` (or memory entry if it's a durable preference). Tag with the area being probed (substrate, grill, packs, doctor, etc.).
5. **Decide next action** — sometimes the answer is "ship it", sometimes "draft a request", sometimes "park, revisit after N more sessions".

### What this is NOT

- Not a test suite. No assertions, no pass/fail counts.
- Not a replacement for `doctor` (which is mechanical substrate checks).
- Not the existing `review` mode (which is a doc-quality pass against principles).

### Open questions

- Canonical mode name: **`feedback-loop`**. Accepted aliases (all route to the same mode): `iterate`, `experiment`, `test`, `probe`, `try`. Also support as a verb on existing docs (`spectacular prd feedback-loop`).
- Where does captured feedback live? `feedback/` folder vs memory entries vs request-scoped notes — likely all three depending on durability.
- How does the skill *decide* what to probe? Heuristics: recently-changed refs, low-signal areas (no feedback in N sessions), user-flagged hunches.
- Cadence: ad-hoc only, or a periodic "eval session" prompt?
- Should proposals be stored even when not yet asked? (Backlog of feedback prompts.)
- Relationship to `grill-me` skill — overlap is real; grill-me interrogates a plan, this interrogates the system itself.

## Done

- ~~clarify distinction specs/ vs docs/~~ — shipped v0.5.0 (`spec-rename`) + v0.6.0 (`public-docs-foundation`). v2 capabilities tracked in `public-docs-advanced` (gated on real demand).
- ~~rename current/ → specs/ + SPEC.md~~ — shipped v0.5.0 in `spec-rename`. Auto-migration via `doctor specs --fix`.
- ~~archive approved/reviewed requests~~ — shipped. `spectacular archive <slug>` flow + auto-detection on `verified` status. Dogfooded this session (spec-rename + public-docs-foundation archived).
- ~~initial request / PRD management~~ — shipped. `init` scaffolds PRD.md, full engine (`prd grill|refine|review`) + 5 kits (blank/coding/content/product/research) + 8-slot base template. No open question remaining.
- ~~verification convention (when VERIFY.md is needed vs PLAN/TASKS fold-in)~~ — shipped 2026-05-22 in `references/verification.md` + lifecycle.md + ARCHITECTURE.md v1.1 + new-request.md + SKILL.md routing. 2-of-6 rule locked.
- ~~prd-craft v1.1 (8-slot base)~~ — verified
- ~~doc-writer (registry + engine + 8 templates)~~ — verified
- ~~kits-as-plugins (diff-only kit contract)~~ — verified
- ~~smart-init (CLI v0.3.0 — always-set + kit-driven + flags + pre-flight + tests/)~~ — verified 2026-05-22 via VERIFY.md walkthrough; 50/50 asserts across 8 scenarios; first request to exercise the 2-of-6 rule and ship a VERIFY.md

## Overrides simplification (parked)

User flagged that `prd-overrides.md` / `plan-overrides.md` / `tasks-overrides.md` may be redundant. Three alternatives discussed:
- **(a)** Eliminate overrides; push everything into templates + registry (preferred)
- **(b)** One file per concern (slot-prompts.md / gate-checks.md / vibe-patterns.md) not per doc
- **(c)** Keep only when complexity warrants — likely delete tasks-overrides + plan-overrides, keep prd-overrides

Open as request `overrides-cleanup` when ready to refactor.

## Host repo structure conventions (not just .spectacular/)

Spectacular currently only opinionated about `.spectacular/` itself. It should also have **opinions about the surrounding repo** so init/new-request can scaffold or suggest the right layout.

### Standard folders (preferred conventions)

- `src/` — source code (default for code projects)
- `scripts/` — utility scripts (preferred over loose root scripts; see global CLAUDE.md)
- `tests/` or `test/` — match language convention (Python: `tests/`, Node: `test/`, Go: `_test.go` co-located)
- `docs/` — human-facing documentation
- `examples/` — runnable examples
- `assets/` — static media
- `_research/` — research artifacts (NotebookLM exports, source dumps, query logs)
- `_archive/` — archived/old content (gitignored by default, see global CLAUDE.md)
- `_backups/` — timestamped backups (gitignored by default)

### Root files

- `README.md` — human-facing intro
- `AGENTS.md` — root-level agent guidance (governs over README for agents)
- `CLAUDE.md` — Claude-specific (often symlink to AGENTS.md, or scoped variant)
- `CHANGELOG.md` — versioned changes
- `LICENSE` — license file
- `.gitignore` — must include `_archive/`, `_archived/`, `_backup/`, `_backups/`, `.spectacular.local/`, tool-generated hidden dirs (`.scrapekit/`, `.playwright-mcp/`) by opt-in

### Naming preferences

- Folders: `kebab-case` for projects, `snake_case` for Python packages
- `_archived/` → prefer rename to `_archive/` (shorter)
- `_backup/` → prefer rename to `_backups/` (plural)
- Database folders: `<name>_db/` suffix (vault convention)

### Project-type-aware scaffolds

`spectacular init` should detect or ask project type and scaffold accordingly:

| Type | Adds | Notes |
|---|---|---|
| `cli` | `src/`, `tests/`, `bin/`, `install.sh` | Bash or compiled binary |
| `library` | `src/`, `tests/`, `examples/`, `docs/` | Language-shaped |
| `webapp` | `src/`, `public/`, `tests/`, `.env.example` | Add framework-specific later |
| `cli-tool` | `cli/`, `scripts/`, `README.md` | Mirrors spectacular itself |
| `skill` | `SKILL.md`, `references/`, `templates/`, `scripts/` | Standard skill scaffold |
| `plugin` | `.claude-plugin/`, `skills/`, `agents/`, `commands/` | Standard plugin scaffold |
| `content` | `articles/`, `_research/`, `assets/`, `drafts/` | For newsletters, books, courses |
| `research` | `_research/`, `notebooks/`, `data/`, `reports/` | For investigations |
| `vault` | Obsidian-style: `core/`, `data/`, `projects/`, `spaces/`, `home/`, `inbox/`, `assets/` | See vault/CLAUDE.md |

### File placement rules (where to put new files)

When creating any new file, follow:

1. **Scripts** → `scripts/` (never root, unless single-file project)
2. **Docs** → `docs/` (architecture, guides, contributor docs)
3. **Reference docs** for skills → `references/` inside the skill folder
4. **Research artifacts** → `_research/` (NotebookLM exports, query logs, source dumps)
5. **Backups** → `_backups/` (always gitignored)
6. **Generated/cached** → `.cache/` or hidden tool dirs (always gitignored)
7. **Sensitive data** → `.env.local`, `.spectacular.local/`, never committed
8. **Large files** (>5MB) → flag to user, never commit silently
9. **Temporary work** → `scratch/` or `_tmp/` (gitignored)

### Where it should be enforced

- `spectacular init` — scaffold the right folders for the project type
- `spectacular new <slug>` — when a request creates artifacts, route them correctly (e.g. research → `_research/<slug>/`, screenshots → `requests/<slug>/artifacts/screenshots/`)
- File-placement reference doc — `skills/spectacular/references/repo-layout.md` for the skill to load on demand
- A `repo-scaffold` command — `spectacular scaffold <type>` to retrofit an existing repo

### Open questions

- Should this be a **separate skill** (`repo-scaffold`) or **baked into spectacular**?
- How opinionated? Suggest vs enforce? (Probably suggest — show diff, ask before creating)
- How to detect project type when not specified? (Read `package.json`, `pyproject.toml`, presence of `SKILL.md`, `.claude-plugin/`, etc.)

## Snapshot cleanup / retention (anti-bloat)

Snapshots accumulate forever — every canonical-doc edit can leave a new
`@v<N>.md` under `.spectacular/snapshots/<DOC>/`. With many docs (now including
per-capability `specs/<cap>/SPEC/`) and a long-running project, this bloats the
tree with files no one reads. Need a way to prune old snapshots.

### What to build

- **`spectacular doctor snapshots`** gains a retention check: flag (info/warning)
  when a doc has more than X snapshots, or snapshots older than Y days.
- **Auto-clean** via `spectacular doctor --fix snapshots` (or a dedicated
  `spectacular snapshots prune`): keep the most-recent N per doc (and/or anything
  newer than Y), delete the rest. Never touch the live canonical file.
- Retention policy configurable in `config.yaml` (e.g.
  `snapshots: { keep: 3, max_age_days: 180 }`). **Default: keep 3** per doc,
  auto-clean older. Prune by **highest @vN** (filesystem mtime drifts on
  clone/restore), not by file date.
- Always show what would be deleted first (dry-run / confirm) — snapshots are
  history; deletion is destructive. Possibly move-to-`.trash/` rather than `rm`.

### Notes from the v1.22 audit of the snapshot system

- **Folder:** `.spectacular/snapshots/<DOC>/@v<N>.md` — historical copies; the
  unversioned file is always current. `cmd_snapshot` copies + bumps `version:`.
- **Handles new spec files?** Yes — `is_canonical_doc` recognizes
  `specs/<cap>/SPEC.md` (fixed v1.18.1); snapshot path mirrors sub-paths
  (`specs/cli/SPEC.md` → `snapshots/specs/cli/SPEC/@v<N>.md`).
- **Modular?** Driven by one `is_canonical_doc` allowlist + path derivation —
  adding a doc type is ~one line.
- **Multi-version?** `@v1, @v2, …` integer sequence (also accepts `X.Y`); next-N
  inferred by scanning; idempotent (no-op when body unchanged, frontmatter excluded).
- **`@vN` vs `version:` drift (by design):** the `@vN` filename is a plain
  snapshot counter (`max(@vN)+1`), while `version:` is a MAJOR.MINOR field bumped
  `minor+1` (or `major+1` with `--major`). They start aligned (`@v1` ↔ `1.0`) but
  diverge on any `--major` bump or hand-set version — `@v3.md` can hold version
  `2.0`. Filename counts snapshots; frontmatter tracks semantic version. Retention
  must key off `@vN` (the counter), never `version:`.
- **Reliable?** Mostly — idempotence + doctor gap/legacy-layout checks. **Gaps:**
  (1) no retention/cleanup (this TODO); (2) stray `.DS_Store` sat in `snapshots/`
  (cleaned 2026-06-28); (3) inconsistent version schemes in the wild
  (`@v1.md` vs `@v1.0.md`) — parser tolerates both but it's untidy.

### Open questions

- Retention by **count**, **age**, or **both**? (**Resolved 2026-06-28:** both —
  count is the primary knob, default keep 3.)
- Keep a **floor** (always retain `@v1` as the origin + last N)? Probably yes.
- `rm` vs move-to-`.trash/`? (Lean: `.trash/` — snapshots are history, deletion should be recoverable.)
- Should `git` already cover this? (Snapshots duplicate what git history holds — worth asking whether the whole snapshot mechanism earns its keep vs `git show <rev>:<file>`. Bigger question; retention is the cheap win regardless.)
- ~~Should `.spectacular/STACK.md` capture the repo conventions per project? Or live separately as `CONVENTIONS.md`?~~ **Resolved (2026-05-21):** CONVENTIONS folded into `ARCHITECTURE.md` (frontmatter schema + versioning + lifecycle) as part of canonical-docs-rework. STACK.md remains for host-project tech only.

## Roadmap-reserved build IDs during request creation

Dogfooding `SPC-003` exposed that `request new` always allocates after the highest
request/config build, even when the same slug already owns a candidate/active
roadmap row. The `github-work-bridge` request was correctly restored to reserved
`b40`, but the allocator temporarily advanced `last_build` to unused `b42`.

- Teach request creation to reuse the exact existing roadmap build when the slug
  matches one unowned candidate/active row.
- Refuse ambiguity or a build already owned by another request.
- Advance `last_build` only when a genuinely new build ID is allocated.
- Add a regression scenario covering a reserved row plus another later request.

## Process learnings — P5/P6 merge and M7 planning session (2026-08-16)

- **Show a brief plan before activating a Mission.** `mission start` freezes and
  activates in one step (`service.go:173`). The owner gate is real but arrives with
  no preview. The Skill should require a short plan summary — title, claim names,
  Objective graph, stops — before the call, and say explicitly that everything else
  is in the file.
- **Read a real record before hand-authoring one.** Guessing the plan shape cost
  five refusals in one sitting (`owner`, `contract`, `validation`, `scope`, `run`),
  every one a wrong nested shape. `scope:` is `{mechanical, semantic}`, not a list.
  Reading M6's frontmatter first costs one call and avoids all of it.
- **Concurrent sessions on separate Proposals worked well.** P5 and P6 were written
  in parallel and audited the merged Contract independently. Both found real defects
  the merging session missed. Worth making a named pattern rather than an accident.
- **An audit that corrects the auditor is the useful kind.** P5 cited
  `.last-mutation` as precedent for its own proposal, then checked and found it
  abandoned — which reversed its own argument. Verifying a cited precedent should be
  explicit in the review step.
- **Check the frozen command surface before proposing a command.** CC-missioncli
  enumerates ten commands and M6 stops on growth. A proposed eleventh survived into
  a Contract draft before an audit caught it. Cheap check, expensive miss.
- **Regenerate `.spectacular/index.md` after adding records by hand.** Adding the
  Contract and two Proposals left `TestSelfHostedIndexesAreRebuildableCollectionCaches`
  failing. It self-heals on the next mutating command, but a hand-authored record
  leaves the tree red until then.

## Dead v1 code (removed by M9/O1, 2026-08-17)

Resolved. A dependency walk from `cmd/spectacular` proved four packages absent
from the main package's transitive closure. Three were deleted as one unit;
`internal/index` was found during the same walk and is recorded below.

- `internal/context`, `internal/projection`, `internal/guardrails` — v1's context
  compiler. `guardrails` supplied declared guidance, `projection` built cards and
  pointers over the workspace, and `context` assembled them into a bounded,
  fingerprinted `Bundle` answering "what should be loaded right now". Deleted
  together: `compiler.go` imported the other two, so removing any one alone broke
  the build.
- `internal/governance` — **retained**. It is reachable from main. Only the
  unreachable `ProposalInput`, `CreateProposal`, and `candidate_*` members were
  pruned; `ApplyTransaction`, `FileChange`, and `RecoverTransactions` are live in
  `internal/command` and `internal/missionbundle/service.go`. The earlier note
  here overstated this as a whole-package removal.
- The original note said `projection` had "no test files" and was unreferenced.
  Both were wrong in detail: it had a live importer in `internal/context`, and the
  chain carried tests. The reachability question is the one that decides deletion,
  not the grep.

### Capabilities lost with the context compiler

Recorded before deletion, per M9's stop on discarding a capability without naming
it. Git history holds the implementation; these are the ideas worth reimplementing
against the v2 model if they earn a Mission:

- **Conflict reporting.** The Bundle named what it could not reconcile. No v2
  surface reports its own internal disagreements.
- **Omission reporting.** The Bundle named what it deliberately left out. v2
  states limits nowhere.
- **Loaded versus available record counts.** The Bundle reported loading twelve of
  forty records, making a bounded-context claim checkable rather than asserted.

The discipline itself survived the rewrite: the compiler's package comment — its
output "is a disposable projection and never owns Mission or Contract truth" — is
the same rule `Bundle.Derive()` follows. v2 reached it more cheaply by deriving
state on read, beside the Bundle it reads.

### `internal/index` — removed 2026-08-17, on owner approval

Found during M9's dependency walk and left in place at the time: the Mission's
frozen scope named three packages, and adding a fourth is `expand-scope`. Removed
immediately after M9 completed, on the owner's explicit approval, rather than
carried as a standing follow-up.

It was the v1 predecessor of `discovery.Workspace.Lookup` — an in-memory record
index keyed by ID and workspace path, with sorted iteration and defensive cloning
on read. Zero importers, not even from a test outside itself; 8 tests exercising
only its own surface.

Nothing was salvaged. `discovery` already provides the lookup this package
existed for, and unlike the context compiler it carried no capability the current
system lacks. The defensive `cloneEntry`/`cloneRecord` pattern is the one idea
worth remembering: it returned copies so a caller could not mutate indexed state
through a read. `discovery` should be checked against that property if it is ever
found to hand out shared structures.

## Noticed 2026-08-23

- [ ] **Windows: `syscall.Flock` does not compile** — `internal/missionbundle/service.go:955`
  uses `Flock`/`LOCK_EX`/`LOCK_NB`/`LOCK_UN`, which are Unix-only. The
  `Cross-Platform Acceptance (windows-latest)` job has failed on every main run since the
  matrix was added in `7404291`; ubuntu and macos pass. Needs a `_unix.go` / `_windows.go`
  split (`LockFileEx` on Windows). This is the same blocker that parked the WASM experiment.
  Recorded as Gap `mutation-lock-is-unix-only` on the mechanical CLI Contract.
- [ ] **`release-proof` skips whenever any acceptance job fails** — it declares
  `needs: [unit-and-race, acceptance]`, so one red matrix leg hides the release gate entirely.
  Consider whether the Windows leg should block it.
- [ ] **Node 20 deprecation** — `actions/checkout@v4`, `actions/setup-go@v5`,
  `actions/upload-artifact@v4` are being force-run on Node 24. Warnings only today.
- [ ] **Proposal index sorts refs as strings** — reads `P1, P10, P11, P12, P2`. Natural sort
  would fix it; the index is non-authoritative so this is cosmetic.
- [ ] **`git gc` is now safe** — `salvage/p11-exploration` tags the one dangling commit worth
  keeping. 751 loose objects at last count.
