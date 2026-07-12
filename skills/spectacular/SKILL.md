---
name: spectacular
description: |
  AI-native operational workspace for software projects. Stop losing context. Start shipping.
  Manages the full lifecycle of a .spectacular/ workspace: reads project state, scaffolds
  requests, manages lifecycle transitions, writes memory, archives completed work, and
  grills/refines/reviews any structured doc (PRD, PLAN, TASKS, PRINCIPLES, POLICY, ARCHITECTURE,
  ROADMAP, STACK, AGENTS, DECISIONS, PERSONAS) plus soft-DB collections (memory, decisions,
  sessions, ideas, feedback, audit, fixes). Enforces a practice layer (POLICY.md): work-phase hooks gate transitions.
  Use when: opening /spectacular, scaffolding a request, archiving work, capturing a memory,
  snapshotting a doc, onboarding to a workspace, or building any canonical doc from scratch.
  Triggers: /spectacular, spectacular status|new|archive|advance|next|init|snapshot|remember|decide|policy,
  spectacular <doc> [grill|refine|review], spectacular pack [new|grill|refine|review].
when_to_use: |
  Invoke on any project that has a .spectacular/ directory. Routes to reference docs based on
  the command — never loads full context, always loads minimally and progressively. The
  generalized doc verbs (grill/refine/review) apply to any doc type listed in doc-index.md.
version: 1.32.0
category: devtools
status: published
tags: [workspace, project-management, context, agents, lifecycle, doc-writing]
---

# Spectacular Skill

AI-native operational workspace for software projects. Lean orchestrator — read this file to understand triggers and routing, then load the relevant reference doc for the actual work.

---

## Trigger detection

### Workspace lifecycle

**Mutation principle (v0.7.0+):** lifecycle mutations go through CLI verbs — never free-form file edits (manual edits only for edge cases the verbs don't cover). See [[lifecycle]].

| User says / context | Route to |
|---|---|
| `/spectacular` with no args | → `references/status.md` (empty workspace → `references/guided-first-run.md`) |
| `spectacular status` | → `references/status.md` |
| `spectacular new <description>` | → `references/new-request.md` (then run `spectacular new <slug>`) |
| `spectacular archive <slug>` | → CLI verb; see [[archive]] (its spec-sync step may dispatch `spec-reviewer` — [[spec-sync]]) |
| `spectacular remember this` | → `references/memory.md` (legacy free-text capture) |
| `spectacular remember "<text>"` | → CLI verb; see [[memory-rules]] for entry shape |
| `spectacular decide "<decision>" [--context\|--consequences]` | → CLI verb; see [[decisions-rules]] |
| "record a decision" / "record an ADR" / "architecture decision" | → `spectacular decide`; ADRs live in decisions/index.md, see [[decisions-rules]] (store-worthy? table) |
| `spectacular session start\|end` | → CLI verb; see [[sessions-rules]] |
| `spectacular idea new <slug>` | → CLI verb; see [[idea-rules]] for entry shape |
| `spectacular idea list` | → CLI verb |
| `spectacular idea promote <slug>` | → CLI verb; scaffolds request, moves source to `archive/ideas/` |
| A bug/quirk/regression is reported (any "why does X do Y", "this is broken") | → **`references/bug-workflow.md`** — load before diagnosing; routes the debug fleet + the ceremony/fan-out gates. (Rationale: `bug-workflow-doctrine.md`, only if a routing call is uncertain.) |
| `spectacular audit new\|list\|resolve` | → CLI verb; bug investigation before a fix. `resolve --into-fix` graduates to a fix (copies all slots). See [[audit-rules]], [[bug-workflow]] |
| `spectacular fix new\|list` | → CLI verb; log a **verified, signed** fix. See [[fixes-rules]], [[bug-workflow]] |
| "record a fix" / "log this fix" / "the bug is fixed and verified" | → `spectacular fix new` once resolved+verified, **with `--signature`**; see [[fixes-rules]] |
| "investigate this bug" / "audit this quirk" before planning | → `spectacular audit new`; see [[audit-rules]] |
| "have we seen this bug before?" / starting to diagnose | → **[[bug-workflow]] Step 0** — grep `.spectacular/fixes/` signatures first (self-learning loop) |
| `spectacular advance <slug>` | → CLI verb (no skill flow); lifecycle move-forward (was `promote`, still an alias); see [[lifecycle]] |
| `spectacular snapshot <file>` | → CLI verb (no skill flow); see [[versioning]] for snapshot rules. Requires a literal path relative to working directory (canonical docs only). |
| `spectacular policy [@hook\|<id>\|--principle N\|--json]` | → CLI verb; read the merged policy contract. See [[policy-injection]] for the runtime loop, [[policies-contract]] for the schema |
| Entering any work phase (init/planning/implementation/verification/archive/remember/snapshot/session-end) | → the phase ref doc's **@\<hook\> policy gate** runs `spectacular policy @<hook>` first; see [[policy-injection]] |
| `spectacular touch <file>` | → CLI verb; trivial — just bumps `updated:`. Requires a literal path relative to working directory, not a slug. |
| First invocation on existing `.spectacular/` project *with prior work* | → `references/onboarding.md` |
| First invocation on a *fresh/empty* `.spectacular/` (init ran, no requests) | → `references/guided-first-run.md` |
| `spectacular init` (CLI context) | → `references/init-workflow.md` |
| `spectacular doctor` / `spectacular doctor <area>` | → `references/doctor.md` (lean entry) |
| `/spectacular doctor --fix` (judgment walk) | → `references/doctor-repair.md` |
| Explain a finding or area check | → `references/doctor-areas.md` |
| Skill operation hits substrate failure (rules file won't parse, kit malformed, etc.) | → `references/doctor-substrate.md` |
| `spectacular migrate [--dry-run\|--list]` | → CLI verb. Mechanical apply of pending schema migrations. |
| `/spectacular migrate` (walk judgment migrations) | → `references/migrate.md` |
| Explain a migration spec or contract | → `references/migrations-contract.md` |
| Actively working on a request | → `references/active-request.md` |
| Implementing a milestone — decide build-inline vs dispatch a `spec-builder` | → **`references/build-workflow.md`** — the closed-brief chain, the inline-vs-dispatch gate, the build fleet. (Rationale: `build-workflow-doctrine.md`, only if a routing call is uncertain.) |

### Read verbs (v1.8.0+) — read-only, no skill flow

Always prefer these over walking the filesystem or hand-reading multiple PLAN/TASKS files.

| User says / context | Route to |
|---|---|
| `spectacular requests [--active\|--status\|--since\|--json]` | → CLI verb. Lists requests with frontmatter view. |
| `spectacular request <slug>` | → CLI verb. Skim view of one request (frontmatter + outline + milestone progress). `--full` for raw. |
| `spectacular decisions [--tag\|--since\|--json]` | → CLI verb. Lists decisions. |
| `spectacular decision <slug>` | → CLI verb. Skim view of one decision. |
| `spectacular memories [--tag\|--since\|--json]` | → CLI verb. Lists memory entries. |
| `spectacular memory <slug>` | → CLI verb. Skim view of one memory. |
| `spectacular sessions [--status\|--since\|--json]` | → CLI verb. Lists sessions (read-only — distinct from `session start\|end` mutators). |
| `spectacular sessions show <slug>` | → CLI verb. Skim view of one session. |
| `spectacular show <doctype>` | → CLI verb. Dumps a canonical doc (prd/spec/principles/...). `--section <name>` filters to one H2. |
| `spectacular summary` | → CLI verb. One-page workspace overview (counts + active requests). Cheap cold-start. |
| `spectacular progress <slug>` | → CLI verb. Milestone tick rate parsed from TASKS.md. |
| `spectacular paths` | → CLI verb. JSON map of conventional paths. Use when locating files programmatically. |

**Universal flags:** `--status <s>`, `--since <Nd\|Nh\|Nw>`, `--limit N` (default 20), `--all`, `--json`. Detail verbs add `--full` to bypass skim mode.

**Cold-start pattern:** prefer `spectacular summary` first → `spectacular requests --active` for context → `spectacular request <slug>` for the one you'll work on. Three CLI calls beats walking the filesystem.

### Doc-writing (generalized — works for any registered doc)

`spectacular <doc> [<verb>]` works for any doc in `references/doc-index.md`; the verb defaults from the doc's mode + state. Each doc's rules file (`references/<doc-id>-rules.md`) declares dispatch in frontmatter, prompts + gate checks in the body.

| User says | Route to |
|---|---|
| `spectacular <doc>` (no verb) | → load `references/<doc-id>-rules.md`, resolve mode, dispatch |
| `spectacular <doc> grill` | → `references/grill.md` (with `<doc-id>-rules.md` context) |
| `spectacular <doc> refine` | → `references/refine.md` (with `<doc-id>-rules.md` context) |
| `spectacular <doc> review` | → `references/review.md` (with `<doc-id>-rules.md` context) |

**Registered docs:** the live registry is the set of `references/<doc-id>-rules.md` files; the catalog is `references/doc-index.md`. No hardcoded id list here — it drifts. (`spectacular prd …` is just this handler with `<doc> = prd`; bare `prd` → grill if empty, else review.)

### Where does this belong? — soft-DB routing

Deciding *which store* a piece of knowledge goes in (fact? decision? fix? idea?) → **`references/soft-db-index.md`**, the canonical index of the seven collections (`memory` · `decisions` · `sessions` · `ideas` · `feedback` · `audit` · `fixes`). Load it whenever the routing isn't obvious.

### Feedback-loop mode (v1.6.0+)

`feedback-loop` is a distinct skill mode for prototyping-stage human-feedback acquisition. **Not** a benchmark or verification pass.

| User says / context | Route to |
|---|---|
| `spectacular feedback-loop` (no args) | → `references/feedback-loop.md` — list open entries, offer to start one |
| `spectacular feedback-loop <target>` | → `references/feedback-loop.md` — run the 5-step loop |
| `spectacular feedback-loop new <target>` | → CLI verb (scaffold one entry, status `open`); see [[feedback-rules]] |
| `spectacular feedback-loop list` | → CLI verb |
| `spectacular feedback-loop resolve <slug>` | → CLI verb (close entry, optional auto-promote to memory) |
| `spectacular feedback-loop archive <slug>` | → CLI verb |
| `spectacular feedback grill\|refine\|review` | → generic engine via [[feedback-rules]] (works like any registered doc) |

Proactive-surfacing rules (three checkpoints only, never mid-flow) and the memory auto-promotion contract live in [[feedback-loop]] — loaded whenever the mode runs.

### Imagine mode (v1.15.0+) — imagination-backed planning

`imagine` is a **generative-first** mode: render see-able ASCII artifacts the human reacts to per-fragment, then derive a draft PLAN from the approved vision. Full engine + v1 scope rules: [[imagine]]. Doc-type rules: [[vision-rules]].

| User says / context | Route to |
|---|---|
| `spectacular imagine <slug>` (bare, in agent) | → `references/imagine.md` — run the render → react → derive loop |
| `spectacular imagine <slug>` (CLI, with slug) | → CLI verb scaffolds `requests/<slug>/vision/` (mechanical), then hands to the skill |
| `spectacular vision add <kind> <name> --slug <s>` | → CLI verb (mechanical fragment mutator); see [[vision-rules]] |
| `spectacular vision grill\|refine\|review` | → generic engine via [[vision-rules]] (manual spine authoring — rare; `imagine` is the default) |
| `spectacular doctor vision` | → `references/doctor.md` (vision area) |

### Pack-specific verbs (`pack` is the canonical doc-id since v1.19.0)

Packs add a `new <name>` verb (user-scope, identified by name):

| User says | Route to |
|---|---|
| `spectacular pack new <name>` | → `references/grill.md` + `pack-rules.md` — pre-flight resolves target `~/.spectacular/packs/<name>/` |
| `spectacular pack new <name> --from <p1>,<p2>` | same + source-ingestion mode active |
| `spectacular pack new <name> --scope project` | same + target `<project>/.spectacular/packs/<name>/` |
| `spectacular pack grill <name>` | → `grill.md` + `pack-rules.md` — resume grill on an existing pack |
| `spectacular pack refine <name>` | → `refine.md` + `pack-rules.md` |
| `spectacular pack review <name>` | → `review.md` + `pack-rules.md` |

### Public-facing docs — owned by pageworks

"Write docs / add a page / add a tutorial" → `references/pageworks-handoff.md` (surface its install hint — never auto-install). Spectacular keeps discovery-only awareness of `docs/`.

### Verification routing (when writing PLAN.md or moving requests to review)

| Decision point | Route to |
|---|---|
| Scaffolding a new request (`spectacular new`) | → **[[plan-rules]] § 2-of-6 rule** (compact table; canonical: verify.md Part 2). Default: no VERIFY.md — `### Verification` group in TASKS.md or PLAN § Validation instead. |
| Grilling/refining a PLAN.md | → **[[plan-rules]] § 2-of-6 rule** — confirm result; ask user if VERIFY.md needed |
| Moving request `active → review` | → `lifecycle.md` § Verification artifact detection — pick artifact (VERIFY.md > TASKS Verification > PLAN Validation) |
| Moving request `review → verified` | → **`verify.md` Part 1** — the interactive validation walk, record to VERIFY-LOG, gate the transition. **Never skip.** |
| `spectacular verify <slug>` | → **`verify.md` Part 1** — the validation walk (skill-only; CLI redirects). |
| Automating a shipped scenario | → `verify.md` Part 3 — when to author `tests/verify/<slug>.test.sh`. |

Verification always runs against *some* artifact — "VERIFY.md is opt-in" means the *file*, never the act; the full doctrine lives in `verify.md` § "Verification always happens" and loads with the walk.

---

## State awareness

Load **only** what the task needs (principle 6 — progressive disclosure). Two authorities, no third list:

- **What to load per task type** — `.spectacular/AGENTS.md`'s context-loading table is authoritative; follow it over guessing or re-deriving a read list.
- **How to read state** — prefer the read verbs (§ Cold-start pattern above: `summary` → `requests --active` → `request <slug>`) over walking the filesystem; the flow docs (`status.md`, `active-request.md`) own their own read steps.

Never read `archive/` during normal operation.

---

## Canonical rules (always apply)

- **Never overwrite canonical documents in place** — snapshot first (`PRD@v1.0.md`). See `references/versioning.md`.
- **Lifecycle state** lives in `PLAN.md` frontmatter (`status: planned | active | review | verified`). TASKS.md mirrors it for skim tooling; PLAN is authoritative — `doctor` repairs drift.
- **Capability state** lives in `specs/<capability>.md` frontmatter (`status: stable | draft | deprecated`); the top-level `.spectacular/specs/index.md` is the always-on index.
- **Slugs** are kebab-case, skill-derived, user-overridable, uniqueness enforced.
- **Memory** (`spectacular remember this`) writes to `.spectacular/memories/` — git-committed, team-visible. Never to `.claude/` memory.
- Be proactive: surface stale state, propose lifecycle transitions, flag blocked requests.
- **Know when to write to a collection, not just how** — the "When to act" trigger table in [[soft-db-index]]. Cheap/reversible writes on their natural trigger; permanent/team-visible writes (memory, decisions, archive) are proposed → human confirms → written, never autonomous.

### Task tracking — two layers

On-disk `requests/<slug>/TASKS.md` owns milestones (persistent, team-visible); harness `TaskCreate`/`TaskUpdate` owns ephemeral session micro-steps (finer-grained — never a one-for-one copy of TASKS.md lines). Full convention: `.spectacular/AGENTS.md` § Task tracking.

---

## Output format

Conversational briefing with a minimal embedded table. Never a raw dump. Identify the single highest-priority next action and ask what the user wants to do.

---

## References & templates index

No hand-list here — it drifts. Reference docs: `references/*.md`, cataloged in `references/doc-index.md`. Templates: the `templates/` tree; frontmatter stubs in `references/scaffold-reference.md`. Projects may override any template at `.spectacular/templates/<doc>/...` (same filenames, project-local wins).
