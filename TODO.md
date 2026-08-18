- follow-up requests suggestions
- catch new plans/ideas suggest wiring plan in spectacular
- add numbering to requests so they appear in order in requests/ folder. number is also id. — **see DECISIONS 2026-05-23 (rejected for v1)**
- when phase ends, verify files, then auto-continue, user expects that spectacular takes task till the end with no interruption
- implement interview mode like grill-me?
- subagents management
- tak inspo from gsd, superpowers
- how do we integrate with octopus tasks/sessions?
- workflow preview (ascii diagram)

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
