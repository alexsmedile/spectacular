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
version: 1.26.1
category: devtools
status: published
tags: [workspace, project-management, context, agents, lifecycle, doc-writing]
---

# Spectacular Skill

AI-native operational workspace for software projects. Lean orchestrator — read this file to understand triggers and routing, then load the relevant reference doc for the actual work.

---

## Trigger detection

### Workspace lifecycle

**Mutation principle (v0.7.0+):** lifecycle mutations go through CLI verbs — never free-form file edits. The CLI is the deterministic mutator; the skill orchestrates, reads, decides, communicates. Manual edits remain available for edge cases the verbs don't cover, but should be the exception. See [[lifecycle]] and individual reference docs for the verb each flow uses.

| User says / context | Route to |
|---|---|
| `/spectacular` with no args | → `references/status.md` (empty workspace → `references/guided-first-run.md`) |
| `spectacular status` | → `references/status.md` |
| `spectacular new <description>` | → `references/new-request.md` (then run `spectacular new <slug>`) |
| `spectacular archive <slug>` | → CLI verb (no skill flow); see [[archive]] for context |
| `spectacular remember this` | → `references/memory.md` (legacy free-text capture) |
| `spectacular remember "<text>"` | → CLI verb; see [[memory-rules]] for entry shape |
| `spectacular decide "<decision>" [--context\|--consequences]` | → CLI verb; see [[decisions-rules]] |
| "record a decision" / "record an ADR" / "architecture decision" | → `spectacular decide`; ADRs live in DECISIONS.md, see [[decisions-rules]] (store-worthy? table) |
| `spectacular session start\|end` | → CLI verb; see [[sessions-rules]] |
| `spectacular idea new <slug>` | → CLI verb; see [[idea-rules]] for entry shape |
| `spectacular idea list` | → CLI verb |
| `spectacular idea promote <slug>` | → CLI verb; scaffolds request, moves source to `archive/ideas/` |
| A bug/quirk/regression is reported (any "why does X do Y", "this is broken") | → **`references/bug-workflow.md`** — check prior fixes first, then decide audit-first vs just-fix. Load this before diagnosing. |
| `spectacular audit new\|list\|resolve` | → CLI verb; bug investigation before a fix. `resolve --into-fix` graduates to a fix (copies all slots). See [[audit-rules]], [[bug-workflow]] |
| `spectacular fix new\|list` | → CLI verb; log a **verified, signed** fix. See [[fixes-rules]], [[bug-workflow]] |
| "record a fix" / "log this fix" / "the bug is fixed and verified" | → `spectacular fix new` once resolved+verified, **with `--signature`**; see [[fixes-rules]] |
| "investigate this bug" / "audit this quirk" before planning | → `spectacular audit new`; see [[audit-rules]] |
| "have we seen this bug before?" / starting to diagnose | → **[[bug-workflow]] Step 0** — grep `.spectacular/fixes/` signatures first (self-learning loop) |
| `spectacular advance <slug>` | → CLI verb (no skill flow); lifecycle move-forward (was `promote`, still an alias); see [[lifecycle]] |
| `spectacular snapshot <file>` | → CLI verb (no skill flow); see [[versioning]] for snapshot rules |
| `spectacular policy [@hook\|<id>\|--principle N\|--json]` | → CLI verb; read the merged policy contract. See [[policy-injection]] for the runtime loop, [[policies-contract]] for the schema |
| Entering any work phase (init/planning/implementation/verification/archive/remember/snapshot/session-end) | → the phase ref doc's **@\<hook\> policy gate** runs `spectacular policy @<hook>` first; see [[policy-injection]] |
| `spectacular touch <file>` | → CLI verb; trivial — just bumps `updated:` |
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

### Read verbs (v1.8.0+) — collapse multi-step inspection to one CLI call

These are read-only — no skill flow needed, no mutation. Always prefer these over walking the filesystem or reading multiple PLAN.md/TASKS.md files manually.

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

The generalized handler matches `spectacular <doc> [<verb>]` where `<doc>` is any doc listed in `references/doc-index.md`. The verb defaults based on the doc's mode and current state.

Each doc is described by a rules file at `references/<doc-id>-rules.md`. The rules file's **frontmatter** declares dispatch (mode, slots, template, location, scope, snapshot-on-edit, kit-support). The rules file's **body** declares per-doc prompts and gate checks.

| User says | Route to |
|---|---|
| `spectacular <doc>` (no verb) | → load `references/<doc-id>-rules.md`, resolve mode, dispatch |
| `spectacular <doc> grill` | → `references/grill.md` (with `<doc-id>-rules.md` context) |
| `spectacular <doc> refine` | → `references/refine.md` (with `<doc-id>-rules.md` context) |
| `spectacular <doc> review` | → `references/review.md` (with `<doc-id>-rules.md` context) |

**Registered docs:** the live registry is the set of `references/<doc-id>-rules.md` files — each declares one doc's dispatch + behavior. The authoritative catalog (every doc-id, its mode, and location) is `references/doc-index.md`; the per-capability detail for the engine itself is in `.spectacular/specs/doc-engine/SPEC.md`. Don't maintain a hardcoded id list here — it drifts every time a doc ships.

`spectacular prd [grill|refine|review]` is just this handler with `<doc> = prd` (bare `prd` → grill if empty, else review).

### Where does this belong? — soft-DB routing

When you have a piece of operational knowledge and must decide *which store* it goes in (a fact? a decision? a bug fix? an idea?), route via **`references/soft-db-index.md`** — the canonical index of the seven soft-DB collections (`memory` · `decisions` · `sessions` · `ideas` · `feedback` · `audit` · `fixes`), each with its role, purpose, structure, write verb, and the boundary rule that keeps entries from landing in the wrong collection. Load it whenever the routing isn't obvious. (`requests/` and the canonical docs are **not** collections — soft-db-index says why.)

### Feedback-loop mode (v1.6.0+)

`feedback-loop` is a distinct skill mode for prototyping-stage human-feedback acquisition. **Not** a benchmark or verification pass. See [[feedback-loop]] for the full spec.

| User says / context | Route to |
|---|---|
| `spectacular feedback-loop` (no args) | → `references/feedback-loop.md` — list open entries, offer to start one |
| `spectacular feedback-loop <target>` | → `references/feedback-loop.md` — run the 5-step loop |
| `spectacular feedback-loop new <target>` | → CLI verb (scaffold one entry, status `open`); see [[feedback-rules]] |
| `spectacular feedback-loop list` | → CLI verb |
| `spectacular feedback-loop resolve <slug>` | → CLI verb (close entry, optional auto-promote to memory) |
| `spectacular feedback-loop archive <slug>` | → CLI verb |
| `spectacular feedback grill\|refine\|review` | → generic engine via [[feedback-rules]] (works like any registered doc) |
**Proactive surfacing — three checkpoints only:**
1. **Milestone tick in TASKS.md** — after acknowledging the milestone, may offer: "Want to feedback-loop M<N> before moving on?"
2. **Request status → `review`** — same single-prompt offer scoped to the request as a whole.
3. **End of `spectacular archive <slug>` flow** — may offer: "Anything worth feedback-looping before this leaves the active set?"

Never mid-flow. Never unsolicited. Single short prompt; user accepts or declines.

**Auto-promotion to memory:** when a feedback resolution captures a durable preference signal ("I always want X", "Y is the right default"), the skill must explicitly confirm before writing a memory entry. No silent promotions. Sets `promoted_to:` on the feedback file.

### Imagine mode (v1.15.0+) — imagination-backed planning

`imagine` is a distinct, **generative-first** mode: the skill renders see-able ASCII artifacts (user stories, UI/output mockups, architecture sketches) the human reacts to per-fragment, then **derives a draft PLAN from the approved vision**. This is Spectacular's second planning axis — spec-driven **and** imagination-backed. Unlike `grill` (which interrogates the human slot-by-slot), `imagine` leads with proposed artifacts. Full engine: [[imagine]]. Doc-type rules: [[vision-rules]].

| User says / context | Route to |
|---|---|
| `spectacular imagine <slug>` (bare, in agent) | → `references/imagine.md` — run the render → react → derive loop |
| `spectacular imagine <slug>` (CLI, with slug) | → CLI verb scaffolds `requests/<slug>/vision/` (mechanical), then hands to the skill |
| `spectacular vision add <kind> <name> --slug <s>` | → CLI verb (mechanical fragment mutator); see [[vision-rules]] |
| `spectacular vision grill\|refine\|review` | → generic engine via [[vision-rules]] (manual spine authoring — rare; `imagine` is the default) |
| `spectacular doctor vision` | → `references/doctor.md` (vision area) |

**Scope (v1):** request-level only, **Build-only** derivation (vision → draft PLAN). Compare/reconcile (diff an existing spec against a vision) and the project altitude (`imagine` near PRD) are v2. The derived PLAN is always a **draft** — it flows into the existing PLAN grill/review gate; never auto-accepted.

### Pack-specific verbs (`pack` is the canonical doc-id since v1.19.0)

Packs use a short alias and add a `new <name>` verb (since packs are user-scope, identified by name, not project-singleton):

| User says | Route to |
|---|---|
| `spectacular pack new <name>` | → `references/grill.md` + `pack-rules.md` — pre-flight resolves target `~/.spectacular/packs/<name>/` |
| `spectacular pack new <name> --from <p1>,<p2>` | same + source-ingestion mode active |
| `spectacular pack new <name> --scope project` | same + target `<project>/.spectacular/packs/<name>/` |
| `spectacular pack grill <name>` | → `grill.md` + `pack-rules.md` — resume grill on an existing pack |
| `spectacular pack refine <name>` | → `refine.md` + `pack-rules.md` |
| `spectacular pack review <name>` | → `review.md` + `pack-rules.md` |

### Public-facing docs — owned by pageworks

Public-facing docs work lives in the dedicated [pageworks](https://github.com/alexsmedile/pageworks) skill. Spectacular keeps **discovery-only awareness** of `docs/` (folder + manifest presence, surfaced by `doctor docs` and an archive-time audit hint). When the user asks to "write docs / add a page / add a tutorial", route to `references/pageworks-handoff.md` and surface its install hint — never auto-install.

### Verification routing (when writing PLAN.md or moving requests to review)

When grilling, scaffolding, or finalizing a PLAN.md for any request, **route to `references/verify.md`** — the single verification reference (Part 1 the walk · Part 2 the 2-of-6 rule · Part 3 promoting checks to scripts; merged in v1.20.0). Decisions:

| Decision point | Route to |
|---|---|
| Scaffolding a new request (`spectacular new`) | → `verify.md` Part 2 — apply 2-of-6 rule. Default: no VERIFY.md. Add `### Verification` group to TASKS.md or fill PLAN § Validation instead. |
| Grilling/refining a PLAN.md | → `verify.md` Part 2 — confirm 2-of-6 rule result; ask user if VERIFY.md needed |
| Moving request `active → review` | → `lifecycle.md` § Verification artifact detection — pick artifact (VERIFY.md > TASKS Verification > PLAN Validation) |
| Moving request `review → verified` | → **`verify.md` Part 1** — run the interactive validation walk: verify each check by its kind (executable / assertable / judgable / observable / manual), record to VERIFY-LOG, gate the transition. **Never skip.** |
| `spectacular verify <slug>` | → **`verify.md` Part 1** — the validation walk (skill-only; CLI redirects). |
| Automating a shipped scenario | → `verify.md` Part 3 — when to author `tests/verify/<slug>.test.sh`. |

**Critical:** "VERIFY.md is opt-in" refers to *creating the file*, not *performing verification*. Verification always runs against *some* artifact. When VERIFY.md exists it is load-bearing; do not bypass it because it's "optional."

The skill never auto-scaffolds VERIFY.md. It is created only when:
- The 2-of-6 rule triggers during request scaffolding, AND
- The user confirms.

---

## State awareness

Before any action, read frontmatter from:
1. `.spectacular/config.yaml` — project config, naming rules
2. `.spectacular/AGENTS.md` — **authoritative** context-loading rules per task type; follow its table over guessing
3. Root canonical docs — `PRD.md` (intent), `PRINCIPLES.md` (rules), `ARCHITECTURE.md` (structure), `ROADMAP.md` (time), `STACK.md` (host tech), `DECISIONS.md` (ADR log)
4. `SPEC.md` (top-level system spec index) + any `specs/<capability>/SPEC.md` (read frontmatter only unless task needs depth)
5. `requests/*/PLAN.md` — active work (read all frontmatter for status briefing)

Load **only** what the task needs (principle 6 — progressive disclosure). For planning, PRD + PRINCIPLES + DECISIONS. For implementation, STACK + PLAN + TASKS + SPEC + relevant `specs/<capability>/`. For review, VERIFY + RISKS + capability specs. AGENTS.md owns the full table.

Never read `archive/` during normal operation.

---

## Canonical rules (always apply)

- **Never overwrite canonical documents in place** — snapshot first (`PRD@v1.0.md`). See `references/versioning.md`.
- **Lifecycle state** lives in `PLAN.md` frontmatter (`status: planned | active | review | verified`).
- **Capability state** lives in `specs/<capability>/SPEC.md` frontmatter (`status: stable | draft | deprecated`); the top-level `.spectacular/SPEC.md` is the always-on index.
- **Slugs** are kebab-case, skill-derived, user-overridable, uniqueness enforced.
- **Memory** (`spectacular remember this`) writes to `.spectacular/memory/` — git-committed, team-visible. Never to `.claude/` memory.
- Be proactive: surface stale state, propose lifecycle transitions, flag blocked requests.
- **Know when to write to a collection, not just how.** Each soft-DB collection has a named prompt-moment — see the "When to act" trigger table in [[soft-db-index]]. Reversible/cheap writes (audit note, session, idea) happen on their natural trigger; permanent/team-visible writes (memory, decisions, archive) are **proposed, human confirms, then written** — never autonomous. Archive is the convergence point (spec-sync + memory + fix/audit capture); see [[archive]].

### Task tracking — two layers

Spectacular uses **two task trackers at different granularities**:

| Layer | Tool | Purpose |
|---|---|---|
| **Milestones** | On-disk `requests/<slug>/TASKS.md` | Persisted, git-committed, team-visible. Owns the M1/M2/M3… block structure user reads to gauge request state. |
| **Session steps** | Harness `TaskCreate` / `TaskUpdate` | Ephemeral, per-session. Decomposes the *current* milestone into concrete edits/commits/tests. Drives the CLI's live progress UI. |

**When starting non-trivial work** (3+ steps), create harness micro-tasks for the immediate steps and mark them `in_progress` → `completed` as you go. Never copy every TASKS.md line into the harness one-for-one — harness tasks are *finer-grained* than TASKS.md items and exist only for the active session.

The on-disk TASKS.md is updated at session end (or at milestone completion within a session). The harness tracker decays naturally when the session ends. Full convention in `.spectacular/AGENTS.md` § Task tracking.

---

## Output format

Conversational briefing with a minimal embedded table. Never a raw dump. Identify the single highest-priority next action and ask what the user wants to do.

---

## References & templates index

This file deliberately does **not** hand-list every reference doc or template — that list drifts every time a doc ships (see the registry note above, and principle 6). The live sources of truth:

- **Reference docs** — the set of `references/*.md`. Human catalog: `references/doc-index.md`. Each doc's dispatch lives in its own `references/<doc-id>-rules.md` frontmatter.
- **Templates** — the `templates/` tree (canonical bases under `templates/<doc>/base.md`, PRD kits under `templates/prd/kits/`, soft-DB entries under `templates/<collection>/entry.md`). Frontmatter stubs for every file type: `references/scaffold-reference.md`.

Project may override any template by placing files at `.spectacular/templates/<doc>/...` — same filenames, project-local takes precedence.
