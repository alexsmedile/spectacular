---
name: spectacular
description: |
  AI-native operational workspace for software projects. Stop losing context. Start shipping.
  Manages the full lifecycle of a .spectacular/ workspace: reads project state, scaffolds
  requests, manages lifecycle transitions, writes memory, archives completed work, and
  grills/refines/reviews any structured doc (PRD, PLAN, TASKS, PRINCIPLES, POLICY, ARCHITECTURE,
  ROADMAP, STACK, AGENTS, DECISIONS, PERSONAS) plus soft-DB collections (memory, sessions,
  feedback, ideas). Enforces a practice layer (POLICY.md): policies under work-phase hooks gate
  transitions (e.g. understand before implementation, verification before verified).
  Use when: opening /spectacular, scaffolding a request, archiving work, capturing a memory,
  snapshotting a doc, onboarding to a workspace, or building any canonical doc from scratch.
  Triggers: /spectacular, spectacular status|new|archive|advance|next|init|snapshot|remember|policy,
  spectacular <doc> [grill|refine|review], spectacular pack [new|grill|refine|review].
when_to_use: |
  Invoke on any project that has a .spectacular/ directory. Routes to reference docs based on
  the command — never loads full context, always loads minimally and progressively. The
  generalized doc verbs (grill/refine/review) apply to any doc type listed in doc-index.md.
version: 1.21.0
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
| `spectacular remember "<text>"` | → CLI verb (v1.5.0+); see [[memory-rules]] for entry shape |
| `spectacular decide "<decision>" [--context\|--consequences]` | → CLI verb (v1.5.0+; flags v1.8.4+); see [[decisions-rules]] |
| `spectacular session start\|end` | → CLI verb (v1.5.0+); see [[sessions-rules]] |
| `spectacular idea new <slug>` | → CLI verb (v1.7.0+); see [[idea-rules]] for entry shape |
| `spectacular idea list` | → CLI verb (v1.7.0+) |
| `spectacular idea promote <slug>` | → CLI verb (v1.7.0+); scaffolds request, moves source to `archive/ideas/` |
| `spectacular advance <slug>` | → CLI verb (no skill flow); lifecycle move-forward (was `promote`, still an alias); see [[lifecycle]] |
| `spectacular snapshot <file>` | → CLI verb (no skill flow); see [[versioning]] for snapshot rules |
| `spectacular policy [@hook\|<id>\|--principle N\|--json]` | → CLI verb (v1.12.0+); read the merged policy contract. See [[policy-injection]] for the runtime loop, [[policies-contract]] for the schema |
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
| `spectacular convention-pack <verb>` | deprecated alias of `pack <verb>` (old doc-id, still accepted) |

### Public-facing docs (DEPRECATED in v1.2.0 — owned by pageworks)

> Public-facing docs work has moved to the dedicated [pageworks](https://github.com/alexsmedile/pageworks) skill. Spectacular keeps **discovery-only awareness** of `docs/` (folder + manifest presence). Schema, authoring, renderer adapters, and validation all live in pageworks now. The `docs *` verbs were removed in v1.17.0.

| User says | Route to | Handled by |
|---|---|---|
| `spectacular doctor docs` | CLI — discovery only (folder presence, manifest presence, pageworks install hint) | CLI binary |
| User asks "write docs", "create a docs page", "add a tutorial" | → `references/pageworks-handoff.md` § install hint | Skill |
| After `spectacular archive <slug>` with SPEC changes | → CLI prints pageworks-audit hint (suppress with `--no-docs-prompt`) | CLI binary |

When delegating to pageworks, surface the canonical install hint from `references/pageworks-handoff.md`. Never auto-install.

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

### Legacy PRD triggers (backwards compatible)

These map to the generalized handler with `<doc> = prd`. Behavior is identical.

| Legacy trigger | Equivalent | Routes via |
|---|---|---|
| `spectacular prd` | `spectacular prd grill` (if empty) or `spectacular prd review` (if filled) | rules file → grill or review |
| `spectacular prd grill` | same | `prd-rules.md` → `grill.md` |
| `spectacular prd refine` | same | `prd-rules.md` → `refine.md` |
| `spectacular prd review` | same | `prd-rules.md` → `review.md` |

PRD behavior is fully handled by the generic engine (`grill.md` / `refine.md` / `review.md`) driven by `prd-rules.md`. (The pre-v1.4 `prd-grill/refine/review.md` files were removed once superseded — snapshots remain in `versions/`.)

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
