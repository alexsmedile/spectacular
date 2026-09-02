# Changelog

## [Unreleased]

## 2.15.0-rc1 — 2026-09-02

### Added

- **Companion Skill: `test-sentinel` (`skills/test-sentinel/`)**: Bundled Staff CI/CD & adversarial test hardening companion skill into Spectacular:
  - 4-Pillar Architecture: Fast Test Pyramid (Tier 0 $\le 1$s hard gate), Cross-Language Determinism Matrix (zero-sleep, dynamic ports, sorted map keys), Regression Shield (`TestM<N>_<slug>`), and Production GitHub Actions CI.
  - Hardened CI Templates (`templates/ci-{go,node,python,rust}.yml`): Action commit SHA pinning, path filtering, shallow clones (`fetch-depth: 1`), and fork 401 guards.
  - Machine-readable receipt schema: `test-sentinel.receipt.v1`.
- **Clean 10-Item Workspace Architecture**:
  - Relocated WAL journals and mutex `.lock` to hidden `.spectacular/.engine/`.
  - Relocated catalog index to on-demand `.spectacular/.cache/catalog.json`.
  - Decoupled standalone PR reviews (`.spectacular/reviews/`) and freeform post-mortems (`.spectacular/retrospectives/`).
  - Documented `.spectacular/raw/` as optional to commit.
- **The Verification Triad (`.spectacular/VOCABULARY.md`)**: Codified formal ontology for Audit, Review, and Assessment, and defined Dynamic Guardrail CLI injection.
- **GitHub Native Integration & Communication Layer (`docs/github-integration.md`)**: Codified 1 PR = 1 Mission model, GitHub Issues $\leftrightarrow$ Proposals/Gaps mapping, GitHub Projects v2, and `gh` CLI orchestration.
- **Dogfooded CI/CD Hardening**: Hardened Spectacular's own `.github/workflows/verify.yml` and `release.yml` with action commit SHA pinning, top-level read-only permissions, path filters, and sub-second Tier 0 preflight fail-fast.

### Added

- **4-Skill Companion Suite (`spectacular/skills/`)**: Shipped three official companion micro-kernels providing full-lifecycle domain engineering alongside Spectacular's mission governance:
  - **`system-architecture`**: Bounded contexts, service boundaries, C4 diagrams-as-code, distributed trade-offs, and Architecture Decision Records (ADRs).
  - **`data-modeling`**: 3-tier data modeling (Conceptual Chen $\to$ Logical Crow's Foot ERD $\to$ Physical SQL DDL/ORMs), PostgreSQL/MySQL/SQLite dialect conversions, Prisma/Drizzle support, composite index ordering, and zero-downtime expand/contract migrations.
  - **`rapid-prototyping`**: Fast 3-option tracer fragment exploration (A, B, C) across a 5-tier fidelity ladder (Atom $\to$ Integration) with persistent decision ledgers (`Locked`, `Open`, `Level`, `Lineage`).
- **Atlas Domain Map for Companion Skills (`.spectacular/atlas/companion-skills-and-domain-execution.md`)**: Visualized the 4-Pillar architecture, mutual exclusivity matrix, and expansion handoff state machine.
- **Strict Expansion Handoff Contracts**: Established clean cross-skill delegation boundaries preventing instruction bloat across domain borders.
- **Multi-Harness Global Distribution**: Configured and verified canonical symlinks across Claude Code (`~/.claude/skills/`), Codex (`~/.agents/skills/`), OpenCode (`~/.config/opencode/skills/`), Cursor (`~/.cursor/skills/`), Antigravity/Gemini (`~/.gemini/config/skills/`), and the canonical library root (`~/skills/`).

## 2.13.0 — 2026-09-01

### Added

- **Mechanical Perimeter Guard & Zero Wasted Work (`spectacular guard`)**: Built-in OS watchdog monitoring subagent file operations, enforcing strict `writes_paths` perimeters, auto-purging rogue files on escape, and preserving valid progress with surgical quarantine.
- **Token Efficiency & Process Isolation**: Cross-platform process isolation and snapshot diffing for autonomous subagents.


### Added

- **Unified 4-Channel Input Protocol across CLI**: Structured commands (`mission start`, `decide`, `review record`, `handoff record`, `evidence record`, `run transition`) now support piped stdin (`-`), direct inline JSON strings, `--data` / `--json-payload` flags, and direct CLI argument flags.
- **Zero-Scratchpad In-Memory Recording**: Eliminated temp files in `/tmp` for decision, review, handoff, and evidence drafting with transparent in-memory validation and atomic commits.
- **Active-First `mission list`**: `spectacular mission list` defaults to live/in-flight missions; added `--all` flag to query the complete historical ledger.
- **Terminal State Recognition in Lifecycle Derivation**: `derive.go` now handles `resolved`, `superseded`, and `cancelled` missions as finished terminal states, resolving ghost operator assignments.
- **Architectural Decision D29**: Formally codified the Banned Synonyms Invariant, Architectural Pattern Pass, and Post-Mission Guardrail Feedback Loop.
- **Permitted Actions & Banned Synonyms Invariant**: Formalized domain ontology schema in `VOCABULARY.md` with explicit canonical actions, permitted entity states, and banned synonyms to eliminate LLM drift across fresh context windows.
- **Architectural Pattern Pass & Pattern Census**: Added upfront 2-track pattern survey (Fast Parametric Survey vs Subagent Research Dispatch) in `prepare.md` to prevent bespoke wheel reinvention and freeze a 3-line Pattern Census in Mission rationales.
- **Post-Mission "Mistake Tax" Feedback Loop**: Codified post-mission learning triage in `close.md` routing failure root causes into `.spectacular/GUARDRAILS.md` (domain invariants) or `AGENTS.md` (tooling rules) without creating redundant governance files.
- **Single-Writer Vocabulary Authority**: Clarified in `mission-anatomy.md` and `runtime.md` that `VOCABULARY.md` is strictly single-writer (Owner/Orchestrator only) with worker subagents operating as read-only consumers.
- **1-Claim Adversarial Hunter Pass**: Codified the Hunter review protocol in `skills/spectacular/references/close.md` for dedicated, read-only subagent falsification passes targeting single high-risk invariants (security, concurrency, migrations).
- **Reachable Surface Filter**: Added structured `reachable_surface` (`cli`, `api`, `contract`, `none-phantom`) to `ReviewDraft` finding schema to immediately drop unreachable phantom hallucinations before consuming repair budgets.
- **Failing-Test-First Rule Guidance**: Established workflow recommendation requiring authors to write and commit a failing deterministic regression test for accepted adversarial findings before implementing fixes and closing Missions.
- **Productivity Genesis Eval Fixtures (MX-06..MX-08)**: Added 3 new full-tier benchmark fixtures (SQLite task CLI genesis, concurrent webhook retry service, and contradictory legacy schema migration) expanding benchmark catalog from 6 to 9 cases.

### Changed

- **Telemetry Normalization**: Updated `agy-adapter.sh` to emit certified observation, usage, and tool_call traces with normalized heuristic token counts, and documented telemetry normalization in `trace.go`.
- **Ephemeral Port Allocation**: Updated `mode-c6` webhook eval fixture verifier to bind dynamic ephemeral ports preventing port collision during concurrent test runs.
- **Cognitive Mode Stance**: Clarified that `mode:` represents prompt-level cognitive posture (defaulting to `leverage`) without requiring mechanical enforcement.

## 2.11.0 — 2026-08-31

### Added

- **The 5 Foundational Anchors**: Formally standardized the grounding layers across `.spectacular/PROJECT.md` (Boundaries/Non-Goals), `VOCABULARY.md` (Domain Ontology), `GUARDRAILS.md` (Invariants/Failure Modes), codebase types/schemas (Data Contracts), and `atlas/` (State Machines/Lifecycles).
- **Dynamic Operating Dial (`mode:`)**: Added explicit execution posture selection to Mission frontmatter (`mode: leverage` for high-autonomy routine tasks; `mode: control` for high-risk auth/payments/DB cutovers).
- **Tiered Verification Matrix (Zero Duplicate Runs)**: Structured test tier execution strictly by agent role (Worker = Tier 1 Quick/Domain; Reviewer = Tier 0 Preflight/Lint; Orchestrator = Tier 2 Acceptance; Owner Gate = Tier 3 All).
- **Batched & Tiered Adversarial Reviews**: Codified policy eliminating per-task review file sprawl by allowing multi-mission campaigns to batch reviews at milestone gates.
- **Architectural Decision D28**: Atomically recorded immutable decision `D28-accepted` formalizing the operating levels, foundational anchors, and verification protocol.

## 2.10.0 — 2026-08-31

### Added

- `spectacular init [--name <project>]` (Command #23 in CLI registry): Safe, conflict-free initialization command creating `.spectacular/workspace.yaml` and `.spectacular/PROJECT.md` without overwriting existing files.
- **The Lean 3-Layer Autopilot Model**: Reduces governance to Living Ground Truth & Decisions (`PROJECT.md` + `decisions/`), Topological Flight Plan (`campaigns/flight-plan.md`), and Single-File Execution Envelopes (`missions/M<N>.md`, $\le 500$ tokens) with zero sub-record sprawl.
- **3-Tier Layout Judgment Protocol**: Standardized classification matrix dividing work into Single-File (90% default), Hybrid Earned (~8%), and Full Governed Bundle (~2%).
- **Dual-Lane Orchestration**: Formalized distinction between Supervised Dispatch ($\le 300$-token charter, 0 files written, tests pass = proof) and Full Ownership Handoff (`spectacular handoff record`).
- **The Escalation & Decision Gate Protocol**: Subagents halt on ambiguity and escalate to Orchestrator to lock choices via `spectacular decide` (`D<N>`).
- **The 3-Tier Question Escalator**: Non-blocking Optimistic Consent (Tier 1), Structured Numbered/Lettered Batch Cards (`1. Q ➔ A, B, C (Recommended)`) (Tier 2), and Trade-off Spectrums / Interactive Modals (Tier 3).
- **Reviewer Role Hygiene (Observe ≠ Act)**: Verifiers inspect primary evidence and report structured verdicts/findings; never apply drive-by fixes to code.
- **Auto-Default Identifiers**: Auto-resolution of `--by`, `--operator`, and `--from` from workspace configuration and Git config.

## 2.9.0 — 2026-08-31

### Added

- Support optional host-harness thread linkage (`runtime_pointer:` with `harness`, `thread_id`, and `workspace_mode`) in `Handoff` and `HandoffDraft` for advanced multi-agent workflows (Claude Code `invoke_subagent`, Codex thread runs).
- Strict vocabulary enforcement for `workspace_mode` (`share`, `branch`, `inherit`) in Handoff validation and recording.
- Comprehensive runtime documentation in `skills/spectacular/references/runtime.md` for thread and subagent orchestration lifecycles.

## 2.8.2 — 2026-08-30

### Changed

- `.spectacular/missions/index.md` filters out subordinate bundle records (Objectives,
  Runs, Evidence, Handoffs, Reviews) to present a clean discovery table of top-level
  Missions only.
- Fully qualify subordinate `Evidence`, `Handoff`, and `Checkpoint` references with
  their parent Mission scope (`M<N>/...`) in `.spectacular/catalog.md`.

## 2.8.1 — 2026-08-30

### Fixed

- `spectacular evidence record <ref> -` preserves the standard input indicator (`-`)
  rather than clearing the draft target path, allowing evidence drafts to be piped
  via stdin.
- `cmd/release-smoke` activates Missions with `--create-branch`, aligning acceptance
  smoke tests with the activation branch guardrails (D15).

## 2.8.0 — 2026-08-30

### Added

- `spectacular contract create <ref> [--title <title>]` scaffolds a new typed
  Contract with UUIDv7 identity, timestamps, and layout.
- `spectacular mission list [--status <status>] [--json]` (Command #22) provides
  a read-only tabular or JSON overview of all discovered Missions, their status,
  current holder, and next action without directory browsing.
- `spectacular mission amend-scope <ref> --add <paths> --by <owner> [--reason <text>] [--dry-run]`
  allows owner-authorized expansion of the mechanical scope boundary, with
  `--dry-run` to preview fingerprint changes before mutation.
- `spectacular mission close <ref> --by <owner>` finishes all in-flight objectives
  and marks an automatic-review Mission completed in a single command.
- `spectacular mission start` includes branch guardrails (`--create-branch` to auto-scaffold
  `feat/<ref>-<slug>`, `--allow-main` to override) and refuses in-place activation on `main`/`master`.
- `spectacular evidence record <ref> --from <test-output.json>` automatically derives
  verifiable test checks, commit hashes, and tree hashes from structured test receipts.
- `spectacular campaign check` dynamically computes and projects `live_state`
  (`planned`, `active`, `complete`, `blocked`) from linked Mission states.
- `spectacular run start <mission-ref>/<objective-ref>` automatically promotes inline
  objectives into standalone objective documents in the same transaction.

### Changed

- `resolveContract` and `validateContract` accept human-friendly contract references
  (`CC-*`) interchangeably with canonical UUIDv7 identifiers.
- `validateScope` validates git working tree modifications against the frozen mechanical
  scope envelope on active Mission branches.
- `mission show` accurately derives `NEXT: close mission (ready for completion)` for
  `review: automatic` missions once all objectives are implemented.

### Fixed

- `contract amend` refused to close a Gap recorded with only a `problem:`
  statement, because the rewrite replaces an existing `blocked_on:` key and
  could not add one. The only way forward was to hand-edit a bound Contract to
  insert a placeholder, which is exactly what an amendment exists to prevent.
  A Gap with no `blocked_on:` now has its `resolution:` appended to the entry it
  names, aligned with its sibling keys, leaving neighbouring Gaps and every
  scalar body byte-identical.
- Closed `mutation-lock-is-unix-only` on the mechanical CLI Contract through
  that path. The Gap's problem statement survives with a stated resolution, as
  a closed Gap should.

## 2.7.2 — 2026-08-24

### Fixed

- The exclusive mutation lock called `syscall.Flock`, whose symbols do not exist
  on Windows, so `internal/missionbundle` did not compile there and every
  mutating command was unavailable. The kernel call is now split behind
  `lockFile`/`unlockFile`: `service_unix.go` keeps `LOCK_EX|LOCK_NB` and
  `service_windows.go` uses `LockFileEx` with
  `LOCKFILE_EXCLUSIVE_LOCK|LOCKFILE_FAIL_IMMEDIATELY`. Immediate-refusal
  semantics are identical on both platforms — a second mutation refuses rather
  than queueing — and all six `GOOS/GOARCH` targets compile, so the
  cross-platform acceptance matrix exercises behaviour instead of failing on
  compilation (`mutation-lock-is-unix-only`).

## 2.7.1 — 2026-08-24

### Added

- `.spectacular/VOCABULARY.md` is an earned Anchor holding the canonical domain
  ontology and ubiquitous language: an alphabetical glossary index, then bounded
  contexts, objects, relationships with cardinality and ownership, actions and
  events, invariants and policies, implementation mappings, and semantic gaps
  (`D25-accepted`).
- `.spectacular/atlas/domain-overview.md` is the visual projection of that
  ontology. It renders bounded contexts as clusters using a fixed node legend;
  relationships are labelled edges, never nodes. The Atlas grants no authority
  and the Vocabulary wins if the two disagree (`D26-accepted`).
- Genesis and Mission guidance now formulate domain nouns before implementation
  verbs and make ontology impact explicit in a Mission body — as guidance, not a
  new frontmatter field, validator, or claim of mechanical enforcement
  (`D27-accepted`, `P13`).

### Changed

- A documented cardinality is not an enforced invariant unless its
  implementation mapping names a schema, code rule, API contract, or test.

## 2.7.0 — 2026-08-24

### Changed

- Workspace entities identify themselves with `type:` and version their rules
  with `schema:`, replacing the combined `atlas_schema:` and `campaign_schema:`
  keys. Non-governance is a property of discovery skip-listing, not of hiding an
  entity's type. The Campaign schema moves to `spectacular.campaign.v2`; there is
  no compatibility reader (`D23-accepted`).
- `spectacular proposal check --schema` and `spectacular campaign check --schema`
  now return the frontmatter template an agent should write, alongside the output
  schema they already reported. Every published template is round tripped through
  the validator that emitted it, so a template can no longer name a field the
  parser does not read.
- A `schema:` field states that Spectacular governs the document and its
  frontmatter is under mechanical check; a document no command validates does not
  carry one. `type:` stays universal. An Atlas therefore declares `type:` alone
  (`D24-accepted`, amending `D23-accepted`).

### Added

- `.spectacular/atlas/` holds Atlas planning maps: top-view projections of a user
  journey or business system, used at the thinking-init stage and later attached
  to a Proposal or Contract. An Atlas is skip-listed in discovery, grants no
  authority, and no command validates it.
- `.spectacular/raw/` is the unstructured sketchpad stage: gitignored,
  skip-listed in discovery, no frontmatter and no entity. A stray note there can
  never refuse a command.
- Proposals and Contracts accept an optional `atlas:` attachment pointing at the
  top-view map they stand on, matching what Campaigns already do.
- `cmd/release-smoke` now drives six mutating commands that no end-to-end test
  reached before: `evidence record`, `handoff record`, `run transition`,
  `decide`, and `contract amend`, alongside the compact Mission loop it already
  covered. Each mutating step asserts that the record it reported is on disk, so
  a command that returns a typed success without writing anything fails the
  release gate rather than a user's repository.
- The `scenario-b-c` fixture contract carries an open Gap so the amendment path
  can be exercised through the declaring-Mission exemption rather than bypassed.

### Fixed

- The doctor compatibility test read a hardcoded version rather than `VERSION`,
  so it passed until a release bumped the version and then failed inside the
  release workflow itself. The release manifest bumper does not edit Go test
  files, so nothing caught it earlier.

## 2.6.0 — 2026-08-23

### Added

- Skill guidance for simpler execution: the seven-rung reuse ladder with a
  non-negotiable preserve-list when planning Mission work (`prepare.md`), a
  diff-level traceability test that records drive-by edits as findings instead
  of repairing them (`close.md`), and a no-mission lane for owner-approved
  micro-tasks (`prepare.md`).

### Changed

- Rewrote the main guides in plainer language, with a clearer focus on what a
  user needs to do and why.

## 2.6.0-rc1 — 2026-08-22

- Adds a test-only Spectacular effectiveness benchmark under `test/evals/spectacular/`
  with immutable old-versus-new package materialization, randomized isolated pairs,
  safety-first scoring, observed token/tool costs, paired-noise statistics,
  resumable trial manifests, 23 behavior and trigger cases (including adversarial
  authority and prompt-injection cases), held-out variants, grader mutation tests,
  and JSON plus Markdown reports.
- Hardens that benchmark with native Codex and Claude trace normalization,
  offline adapter certification fixtures, self-report rejection, failed-trial and
  resume integrity, explicit measurement/comparison/readiness verdicts, total
  cost-per-success reporting, model-call and per-trial budget guards, and a
  six-case productivity frontier comparing native direct, native planning,
  canonical Markdown workspace, and full Spectacular modes.

- Adds `-h` and `--help` flag interception across all subcommands. For document-input
  commands (`mission start`, `review record`, `handoff record`), `--help` emits the syntax
  and an annotated, minimal valid YAML frontmatter starter skeleton directly to stdout,
  eliminating the need to crawl reference docs for template schemas.
- Adds `--schema` inspection flag to all commands, returning machine-readable JSON schema
  metadata, input types (e.g. `MissionPlan`, `ReviewDraft`, `HandoffDraft`), and output types.
- Auto-derives `reviewed.tree` from `commit^{tree}` in `ReviewDraft` and `HandoffDraft`
  when omitted, removing the redundant requirement to manually compute and copy `git rev-parse HEAD^{tree}`.
- Reframes the public command surface rule in `AGENTS.md` to require explicit owner authorization
  when modifying or introducing commands, rather than enforcing a rigid static cap.
- Refactors `skills/spectacular/SKILL.md` into a compact, role-first constitutional
  kernel (83 body lines, 62% fewer words) with a 4-in-1 role matrix
  (`Orchestrator`, `Runner`, `Reviewer`, `Autopilot`), a 3-state mechanical mode invariant,
  and a closed primary phase router.
- Extracts detailed procedural guidance into dedicated progressively disclosed references,
  adding `references/reduced-mode.md` (CLI fallbacks and installation) and
  `references/owner-guidance.md` (question formats, authorization holding, and batching).
- Standardizes uniform `Use this when:` activation triggers across all reference documents.

## 2.5.0 — 2026-08-18

- Detects the CLI before doing anything and states the reduced mode when it is absent.
  The Skill previously assumed the binary was present and had no way to tell a user it
  was not. Reading, explaining, and drafting stay available without it; starting,
  promoting, completing, amending, and recording do not, because those produce
  fingerprints and transactional writes. The Skill is instructed never to emulate a
  missing CLI — no hand-written record that a command owns, no invented fingerprint, and
  no plain file write described as atomic.
- Bundles read-only fallbacks with the Skill so a host without the binary is usable
  rather than blind. Three POSIX shell scripts need no toolchain (`doctor.sh`,
  `orient.sh`, `where.sh`), and two Node helpers parse frontmatter properly
  (`show.mjs`, `check.mjs`). All read and report only; none writes, computes a
  fingerprint, or verifies one, and each says so in its own output.
- Publishes through the portable Agent Plugins manifest, making the repository
  installable on ChatGPT, Cursor, GitHub Copilot, Kiro, and VS Code from one root
  `plugin.json`. Claude Code and Codex keep their vendor manifests. Only `skills/` and
  `mcp.json` travel under the standard, so a conformant host installs the method and the
  binary is installed alongside it.
- Adds human-facing product documentation under `docs/`: a quickstart that runs one
  Mission end to end, an architecture page, a process page with a lifecycle diagram, and
  an installation guide covering both halves and both update paths.

- Adds Abstract Model Profiles (`reasoning`, `fast-code`, `strict-verifier`) to semantic
  work descriptions, decoupling agent roles from host-specific model flags. High-reasoning
  models handle Genesis, scoping, and audits; fluent coding models execute routine worker
  sweeps; and clean-context verifiers conduct adversarial validation.
- Adds Dual-Path Independent Review workflows for `review level: independent`:
  Path A dispatches clean-context in-harness subagents with the Git commit/tree hash
  and FROST checklist directly to `reviews/`; Path B generates a copy-pasteable review
  prompt in `handoffs/review-handoff-prompt.md` to evaluate against external models
  (OpenAI o3, DeepSeek-R1, external sessions, peer developers) and record via CLI.
- Teaches One-Shot Genesis to auto-detect intake starter PRDs (`./PRD.md`,
  `scratch/PRD.tmp.md`, or outputs from the companion `write-prd` skill) and losslessly
  digest all 8 foundational PRD dimensions into Core Anchors, On-Demand Anchors, and
  the `M1-bootstrap` Mission plan without interactive interview fatigue.
- Clarifies the directory boundary between governed delegation handoffs and review
  prompts (`.spectacular/missions/<slug>/handoffs/`), formal review verdicts (`reviews/`),
  and ephemeral kickoff/scratchpad files (`scratch/` or project root).
- Defines what happens to a Proposal once its work has shipped. A Proposal had a
  validated status and no end state: nothing advances the field, so an absorbed Proposal
  read `draft` indefinitely and the live folder described a backlog that no longer
  existed. An absorbed Proposal now names its resolver in `resolved_by:` and is retired
  to `.spectacular/archive/proposals/` under an authorizing Decision and a fingerprint,
  the same admission rule an archived Mission follows. `resolved_by:` is written before
  the move, because once the record leaves `proposals/` that field is the only thing
  tying it to the work that answered it. P6, P7, and P9 were retired under
  `D11-proposal-retirement`; P5 stays live, having shipped three of its four directions.
  No command was added — retirement is an owner act.

## 2.4.0 — 2026-08-18

A delegation is now a record rather than a chat message, record paths resolve through
one system instead of two, and every record names itself.

- Adds `handoff record <mission-ref> <handoff.md|-> --by <sender>`, which files a
  delegation into the Mission bundle and binds the exact commit and tree it was sent
  against, verified against the repository. A Handoff separates `asserted` — what the
  sender verified — from `assumed`, what the sender is carrying on trust; both are
  required, an empty list is a legal statement and an absent one is not. Neither is ever
  scored: the split records a claim its sender signs, and the receiver re-verifies
  everything under `assumed` before acting on it.
- A recorded Handoff is frozen and corrected by recording a new one carrying
  `supersedes:`. The original survives as what its sender believed at the time, and
  `mission show` points a reader of a superseded record forward to the one that is
  current. A Handoff is never edited in place.
- Resolves the Review record path through the layout system rather than a hardcoded
  join, so Reviews and Handoffs are placed by one mechanism. No recorded Review moved.
- The command surface goes to thirteen, authorized by the owner. `proposal create`
  stays forbidden.
- Restates the command-surface rule. Adding a command now requires owner authorization
  and the count is reported rather than defended, replacing a hard stop at twelve. The
  number was a proxy for keeping the surface from sprawling unnoticed, and as an absolute
  it would have refused a correct command in order to protect a count.
- Fixes the Gap rewrite mistaking text for a key. It located `blocked_on:` by walking
  lines without tracking block-scalar depth, so a Gap whose `problem:` scalar contained
  the literal text would have had the resolution spliced into the middle of a sentence
  while the real key stayed untouched. The textual approach is kept deliberately, so an
  amendment's diff still touches only what it changed.
- Re-pointing a bound Mission now refuses when the old Contract fingerprint appears more
  than once in the Mission file, naming the Mission, the fingerprint, and every
  occurrence, rather than rewriting the first one and leaving the real binding stale.
  The refusal was chosen over anchoring to the `contract:` block because it turns a
  silent corruption into a stated problem.
- Amendment no longer refuses while the Mission that declared the Gap is live. Completion
  refuses until a declared Gap closes and amendment refused while the declaring Mission
  was live, which deadlocked the first Mission ever to declare one. A Mission closing the
  Gap whose wording it froze at its own activation gate is exempt; an unrelated live bound
  Mission still blocks, and an owner `--resolution` override is never exempt because its
  wording was typed at a prompt rather than approved at an activation gate.
- Re-points only the live Mission. Amendments previously rewrote `contract.fingerprint`
  on every bound Mission, completed and archived included. A completed Mission's binding
  is the historical fact of which agreement it was executed against, and rewriting it
  replaces that fact with today's answer — destroying the record re-pointing was meant to
  protect. The stale binding was never a defect: `mission check` reports it as a
  contract-drift notice, the Mission stays valid, and `git log -S <fingerprint>` recovers
  the exact Contract text in force. A live Mission still re-points, because its binding is
  a statement about the present.
- Names a Mission or Run record by its own reference: `MISSION.md` and `RUN.md` become
  `<ref>-<slug>.md`, so a record is identifiable from a search result, an editor tab, or a
  diff without its parent directory supplying the context. The layout system resolves the
  scoped name, and a Run resolves to `runs/<run-ref>-<slug>/<run-ref>-<slug>.md`. Every
  existing Mission and Run, live and archived, was renamed with its history preserved.

## 2.3.0 — 2026-08-17

A signed Contract can now be amended, so a Contract that becomes wrong no longer
stays wrong, and every check a Contract declares either runs or fails loudly.

- Adds `contract amend <contract-ref> --gap <gap-ref> --by <owner>`, which closes one
  Gap by rewriting `blocked_on:` to `resolution:` and re-points the bound Missions, all
  as one recoverable transaction. `--dry-run` prints the resolution text, both
  fingerprints, and every Mission that would be re-pointed, and writes nothing. The
  amendment reaches the `gaps:` block and editorial fields only, and refuses while any
  bound Mission is live. A Gap is never closed by deleting it.
- Records each amendment in an append-only companion file beside the Contract, carrying
  time, owner, the Mission that declared the resolution, and the Contract fingerprint
  before and after. The log lives outside the Contract so an entry can record the
  fingerprint its own amendment produced.
- Stops refusing a completed Mission over a Contract that changed after it completed.
  The bound-Contract check previously gated on nothing, so amending a Contract refused
  every Mission ever bound to it — with no legal correction, because a completed Mission
  is never rewritten to satisfy a fingerprint. A completed Mission now reports the change
  as a notice; a live one still refuses.
- Adds `resolves_gaps:` to the Mission plan, naming the Gaps a Mission closes and the
  resolution text it will write. Both are frozen in the activation fingerprint, so the
  owner approves the exact wording at activation and a Mission cannot acquire the
  authority to amend a Contract afterwards. Completion refuses while a declared Gap is
  still open, and names the command that closes it.
- Adds `proposal check <ref>`, wiring the Proposal validator that existed with no caller.
  Proposals stay authored as Markdown wherever the author keeps them and are checked
  rather than generated; `proposal create` remains forbidden.
- Fixes four validation names that a Contract declared and nothing ran, so a Contract
  could promise a check that never executed. `frozen-fallbacks` is renamed to the
  declared `fallback-fingerprint-coverage`, the interface-dependency frozen-target check
  is split into its own registered validator, `ref-spelling-drift` is recognized as a
  notice rather than a validator, and `proposal-schema-v2` resolves to the newly wired
  Proposal check. Every declared name in every Contract now resolves, and one that
  matches nothing refuses.
- Reads `contract_version:`, which was carried by every Contract and used by nothing. It
  is validated as a positive integer and reported on `mission check`. A Mission bound to
  an earlier version is not refused or migrated: it ran against that version, and that is
  a true fact about it.
- Adds `amend-contract` to the owner authority vocabulary. A Mission declaring
  `resolves_gaps:` must list it under `requires_owner`, so the authority to change a
  signed agreement is declared rather than implied.
- Grows the command surface from ten to twelve, deliberately. Growth past twelve is a
  stop.

## 2.2.0 — 2026-08-17

Dead v1 surface removed, derived state reads recorded reviews, and Skill guidance
for branching, worktrees, and execution mode.

- Removes the unreachable v1 context-compiler chain — `internal/context`,
  `internal/projection`, and `internal/guardrails` — as one unit, and
  `internal/index`, the v1 predecessor of `discovery.Workspace.Lookup`. Prunes
  `internal/governance` to the transaction machinery that has live callers.
- Fixes the derived next action to consult the reviews a Mission carries. It
  previously fired on implemented Objectives alone, so a recorded review could not
  retire the instruction to record one. A review bound to a stale activation
  fingerprint keeps asking, since that is the drift the fingerprint exists to catch.
- Surfaces the underlying cause in human refusals. A YAML syntax error reported
  only `invalid_known_field`, sending readers hunting through field names while the
  parser's line number sat unprinted in the JSON envelope.
- Adds Skill guidance for choosing a branch and a worktree by what the job needs,
  including running a Mission session and a feedback session concurrently without
  either destroying the other's work.
- Adds an execution-mode question at activation — autopilot, checkpoints, or a
  named human-in-the-loop moment — so involvement is settled once instead of
  arriving one gate at a time.
- Accounts for every untracked working-tree path with a `.gitignore` rule that
  states its reason, and records `_snapshots/` as local recovery only.

## 2.1.1 — 2026-08-17

Test performance optimizations and contributor verification guidance.

- Optimizes `internal/missionbundle` test fixtures by caching template git repositories and parallelizing rollback subtests.
- Parallelizes 4-platform release archive compilation in `cmd/assemble-release`.
- Retains persistent Go test cache across verification runs and downstream installer/release test scripts.
- Codifies tiered verification guidance (`quick`, `acceptance`, `release`, `all`) in `AGENTS.md`.

## 2.1.0 — 2026-08-16

Governed-autonomy controls for preparation, delegation, review, and runtime limits.

- Adds adaptive preparation diagnostics and a frozen four-field completion criterion for every Mission claim.
- Adds criterion-driven automatic, clustered, and independent review without recursive critic loops.
- Validates Objective dependency DAGs and requires dependency-bound, disjoint Handoff claim scopes with explicit return contracts.
- Adds truthful hard, observed, and unsupported Autopilot caps for wall time, tokens, spend, parallel workers, and repair rounds.
- Fixes cold recovery for newly activated Missions before their first Checkpoint.

## 2.0.0-rc.2 — 2026-08-10

Human-operability correction for the v2 release candidate.

- Replaces flat UUID filenames with named project Anchors and cohesive Mission
  bundles while retaining UUIDv7 identity and SHA-256 revision fingerprints.
- Adds scoped human references for Missions, Objectives, Runs, Checkpoints,
  Evidence, Decisions, Gaps, Handoffs, and Assessments.
- Makes default CLI cards human-first and keeps exact machine data in `--json`.
- Commits deterministic, non-authoritative workspace and Mission indexes.
- Adds atomic whole-bundle Mission archival and a real self-hosted workspace.
- Replaces flat test fixtures with human-layout scenarios and adds a real-binary
  acceptance layer covering cold recovery, executable pointers, governed
  closure, archival, refusals, and zero-mutation reads.
- Fixes stale active indexes/directories after Mission archival, stable bundle
  placement after title changes, and invalid empty optional Evidence fields.

RC.1 is superseded because its machine-oriented workspace representation did
not satisfy the human-comprehension contract.

## 2.0.0-rc.1 — 2026-08-10

First externally consumable Spectacular v2 release candidate.

- Breaking: removes the v1 public product surface and compatibility paths; Git
  recovery pointers preserve the frozen history.
- Adds pointer-first retrieval and a governed Mission loop with explicit
  authority, Evidence, assessment, reconciliation, and closure boundaries.
- Ships native four-platform archives and a checksum-verifying installer.

Known limitations: v1 workspaces are unsupported; there is no migration or
compatibility reader; discovery is pointer-driven rather than broad search;
and release publication remains outside the local CLI.
