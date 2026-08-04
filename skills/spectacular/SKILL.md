---
name: spectacular
description: |
  AI-native operational workspace for software projects. Operates a .spectacular/ workspace:
  retrieves state, scaffolds and advances requests, runs lifecycle gates, captures durable records,
  and grills/refines/reviews canonical documents. Use only for durable workspace work—requests,
  specs, decisions, ideas, lifecycle, or canonical docs—not merely because a repo has
  .spectacular/. A bounded code/docs/configuration change that can ship as one PR without durable
  planning is a direct change. For natural-language work requests, route intent before drafting an
  SPC or request. Triggers: /spectacular; spectacular status|new|archive|advance|next|init|snapshot|
  remember|decide|policy; spectacular <doc> [grill|refine|review]; spectacular pack [new|grill|
  refine|review].
when_to_use: |
  Invoke when the user explicitly requests a Spectacular operation, or when the task needs durable
  workspace context, planning, lifecycle, or canonical documentation. First distinguish a direct
  PR-shaped change from workspace work; then determine whether the workspace work needs a new SPC.
  Routes to reference docs based on the command — never loads full context, always loads minimally
  and progressively. The generalized doc verbs (grill/refine/review) apply to any doc type listed
  in doc-index.md.
version: 1.37.0
category: devtools
status: published
tags: [workspace, project-management, context, agents, lifecycle, doc-writing]
---

# Spectacular Skill

AI-native operational workspace for software projects. Lean orchestrator — read this file to understand triggers and routing, then load the relevant reference doc for the actual work.

---

## First decision — does this need Spectacular?

Before interpreting an ordinary user request as workspace work, run the short
route in [[intent-routing]]. **A `.spectacular/` directory is context, not an
automatic trigger.** Choose the smallest fitting route:

| User need | Smallest owner / artifact | Does it create durable workspace state? |
|---|---|---|
| Read the codebase or inspect current workspace status | Read code directly; use `spectacular summary`, `status`, or `request <slug>` only when the user asks for workspace state | No |
| Bounded code, docs, or configuration change with clear outcome, likely files, and check | Direct PR-shaped change plus an in-chat Codex plan if it has several steps | No |
| Work already owned by a live request | Resume that request; its PLAN/TASKS remain authoritative | No new record |
| New behavior, contract, product/architecture decision, or multi-session/dependent implementation boundary | New SPC candidate, then a request only after SPC approval | Yes, after explicit confirmation |
| Missing fact, feasibility result, business choice, or future commitment | `RES`, `SPK`, `QUE`, or `IDEA` | Yes, only when that open loop needs to survive the chat |

`PLAN.md` and `TASKS.md` are durable request artifacts: PLAN owns cross-session
scope; TASKS owns milestone progress. A Codex/harness plan is an ephemeral
in-chat execution checklist, and subagent briefs are narrower still. Neither
creates nor requires a Spectacular request. “Implement this plan” means execute
the plan already present in chat or named by the user; ask which plan only when
its scope is unavailable or materially ambiguous.

An explicit terminal `spectacular spec new <slug> --summary <text>` is a
user-directed mechanical write; it does not infer intent. For a natural-language
“draft a spec” request, the intent receipt is mandatory before invoking that CLI
verb.

---

## Trigger detection

### Workspace lifecycle

**Mutation principle (v0.7.0+):** lifecycle mutations go through CLI verbs — never free-form file edits (manual edits only for edge cases the verbs don't cover). See [[lifecycle]].

| User says / context | Route to |
|---|---|
| An ordinary change/build/plan request that does not explicitly name a Spectacular operation | → [[intent-routing]] first: decide read-only orientation vs direct PR-shaped change vs existing request vs new SPC vs an open-loop record. |
| `/spectacular` with no args | → `references/status.md` (empty workspace → `references/guided-first-run.md`) |
| `spectacular status` | → `references/status.md` |
| `spectacular request new [<slug>] --from <SPC>` | → CLI verb; approved-spec PLAN/TASKS scaffold; see [[request-workflow]] and [[new-request]] (`spectacular new` remains the free-form compatibility alias) |
| `spectacular request new <slug> --from-issue <ref> --summary <outcome> --sensitivity <class>` / `--from-goal <ref>` | → CLI verb; lean durable coordination for an already-defined destination, without manufacturing an SPC; see [[github-work-bridge]] and [[request-workflow]] |
| `/spectacular act <SPC>` / `/spectacular <SPC>` / compatibility `/spectacular spec act <SPC>` | → [[request-workflow]] approved-spec handoff: resolve one request, run gates, activate with provenance, compile `--brief`, initialize native planning, implement |
| `spectacular archive <slug>` | → CLI verb; see [[archive]] (its spec-sync step may dispatch `spec-reviewer` — [[spec-sync]]) |
| `spectacular remember this` | → `references/memory.md` (legacy free-text capture) |
| `spectacular remember "<text>"` | → CLI verb; see [[memory-rules]] for entry shape |
| `spectacular decide "<decision>" [--context\|--consequences]` | → CLI verb; see [[decisions-rules]] |
| "record a decision" / "record an ADR" / "architecture decision" | → `spectacular decide`; ADRs live in decisions/index.md, see [[decisions-rules]] (store-worthy? table) |
| `spectacular session start\|end` | → CLI verb; see [[sessions-rules]] |
| `spectacular idea new <slug>` | → CLI verb; see [[idea-rules]] for entry shape |
| `spectacular idea list` | → CLI verb |
| `spectacular idea promote <slug>` | → CLI verb; scaffolds request, moves source to `archive/ideas/` |
| `spectacular question new\|list\|resolve` | → CLI verb; active human blockers live in `questions/`; see [[question-rules]] and [[canonical-ids]] |
| `spectacular research new\|list\|resolve` | → CLI verb; read-only evidence discovery in `research/`; see [[research-rules]] |
| `spectacular spike new\|list\|resolve` | → CLI verb; human-authorized feasibility discovery in `spikes/`; see [[spike-rules]] |
| Choosing research vs spike vs prototype vs tracer bullet, or routing technical debt | → `references/discovery-protocol.md`; use the cheapest sufficient answer and avoid creating redundant nodes |
| Deciding what must stay current, may remain stale, should archive, or may be deleted | → `references/artifact-retention.md`; derive live/stale-safe/temporary/throwaway from entity, status, and path |
| “Draft/create a new spec” / `/spectacular spec draft` | → [[intent-routing]] intent receipt, then `spectacular spec new`; see [[spec-lifecycle]]. Never infer a new SPC from nearby docs impact or repository context. |
| `spectacular spec new\|list\|approve\|act\|implement\|deprecate\|archive` | → CLI verbs; `spec new` is mechanical only and assumes the user already chose the route and confirmed the summary. Evidence-gated lifecycle follows; `confirm` remains an alias for `approve`; see [[spec-lifecycle]], [[lifecycle-contract]], and [[canonical-ids]] |
| `spectacular wayfind status\|next\|order\|resolve\|defer\|resume\|path\|route` | → CLI verb; strict dependency-first fog/frontier sequencing and durable open-loop control; see [[wayfinding-sequencer]] |
| “Park this idea” | → `spectacular wayfind route "park this idea" <slug>`; creates an `IDEA-NNN` record and does not alter active milestone scope |
| “Put it on ice” / “Icebox” | → `spectacular wayfind route icebox <id> --reason <why>`; durable `status: deferred` |
| “Find your way to <destination>” | → `spectacular wayfind path <id>` first; then resolve its dependency-ready discovery path without bypassing gates |
| “Act on goal <target>” | → resolve the approved SPC, then run `/spectacular act <SPC>` per [[request-workflow]] |
| `spectacular afk run\|status\|configure\|propose\|preflight\|start\|cleanup\|pr` | → CLI verbs; durable goal authorization plus opt-in, dry-run-first Git isolation and verified PR handoff; see [[afk-git-hygiene]] and [[lifecycle-contract]] |
| `spectacular github triage <issue>` | → [[github-work-bridge]] agentic readiness card and `direct | request | spec-first` route; assignment/labels are evidence, not authorization |
| `spectacular github pr open\|ready` / `github reconcile` | → CLI verbs; draft PR integration manifest, current-head ready gate, and read-only discrepancy report; see [[github-work-bridge]] |
| User authorizes AFK work | → inspect `spectacular afk status`; branch mutation still requires enabled project config plus explicit apply; merge/remote deletion remain HITL |
| A built-in `/goal` begins in a Spectacular workspace | → create/resume a durable goal-scoped AFK run for that goal; keep it active until completion/cancellation or a declared/unexpected HITL gate |
| A user refers to `D1`, `Q1`, `R1`, `SPK1`, `S1`, or another entity alias | → normalize via `spectacular id resolve`; persist the canonical ID; see [[canonical-ids]] |
| A bug/quirk/regression is reported (any "why does X do Y", "this is broken") | → **`references/bug-workflow.md`** — load before diagnosing; routes the debug fleet + the ceremony/fan-out gates. (Rationale: `bug-workflow-doctrine.md`, only if a routing call is uncertain.) |
| `spectacular audit new\|list\|resolve` | → CLI verb; bug investigation before a fix. `resolve --into-fix` graduates to a fix (copies all slots). See [[audit-rules]], [[bug-workflow]] |
| `spectacular fix new\|list` | → CLI verb; log a **verified, signed** fix. See [[fixes-rules]], [[bug-workflow]] |
| "record a fix" / "log this fix" / "the bug is fixed and verified" | → `spectacular fix new` once resolved+verified, **with `--signature`**; see [[fixes-rules]] |
| "investigate this bug" / "audit this quirk" before planning | → `spectacular audit new`; see [[audit-rules]] |
| "have we seen this bug before?" / starting to diagnose | → **[[bug-workflow]] Step 0** — grep `.spectacular/fixes/` signatures first (self-learning loop) |
| `spectacular request advance <slug>` | → CLI verb; lifecycle move-forward (`advance` remains an alias); review→verified is evidence-gated; see [[request-workflow]] and [[lifecycle]] |
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
| Actively working on, resuming, or retrieving context for a request | → `references/request-workflow.md`, then `references/active-request.md` only for session-state details |
| Implementing a milestone — decide build-inline vs dispatch a `spec-builder` | → **`references/build-workflow.md`** — the closed-brief chain, the inline-vs-dispatch gate, the build fleet. (Rationale: `build-workflow-doctrine.md`, only if a routing call is uncertain.) |

### Read verbs (v1.8.0+) — read-only, no skill flow

Always prefer these over walking the filesystem or hand-reading multiple PLAN/TASKS files.

| User says / context | Route to |
|---|---|
| `spectacular requests [--active\|--status\|--since\|--json]` | → CLI verb. Lists requests with frontmatter view. |
| `spectacular request <slug>` | → CLI verb. Cheap overview. `--brief [-mN]` compiles active implementation context; `--full` emits the ordered request-owned Markdown bundle. See [[request-workflow]]. |
| `spectacular spec <id> [--json]` | → CLI verb. Compact SPC state, intent, linked request, and safe next action; use `spec list` only to discover IDs. |
| `spectacular decisions [--tag\|--since\|--json]` | → CLI verb. Lists decisions. |
| `spectacular decision <slug>` | → CLI verb. Skim view of one decision. |
| `spectacular memories [--tag\|--since\|--json]` | → CLI verb. Lists memory entries. |
| `spectacular memory <slug>` | → CLI verb. Skim view of one memory. |
| `spectacular sessions [--status\|--since\|--json]` | → CLI verb. Lists sessions (read-only — distinct from `session start\|end` mutators). |
| `spectacular sessions show <slug>` | → CLI verb. Skim view of one session. |
| `spectacular show <doctype>` | → CLI verb. Dumps a canonical doc (prd/spec/principles/...). `--section <name>` filters to one H2. |
| `spectacular summary` | → CLI verb. One-page workspace overview (counts + active requests). Cheap cold-start. |
| `spectacular status --brief [--json]` | → CLI verb. Bounded orientation: blockers, request-health signals, fleet, and one safe next action. `--json` uses `spectacular.status.v2`; bare `status --json` remains the fleet-array contract. |
| `spectacular progress <slug>` | → CLI verb. Milestone tick rate parsed from TASKS.md. |
| `spectacular paths` | → CLI verb. JSON map of conventional paths. Use when locating files programmatically. |

**Universal flags:** `--status <s>`, `--since <Nd\|Nh\|Nw>`, `--limit N` (default 20), `--all`, `--json`. Detail verbs add `--full` to bypass skim mode.

**Cold-start pattern:** prefer `spectacular status --brief --json` → follow its named `request <slug>` (or `--brief` for active implementation). Use `summary` for a human-facing dashboard, not task selection.

### Doc-writing (generalized — works for any registered doc)

The conversational canonical grammar is verb-first: `/spectacular grill|refine|review <doc> [target]`. Document-first forms remain compatibility aliases. The verb defaults from the document's mode + state only when the user supplies the document without a verb. Each doc's rules file declares dispatch and gate checks.

| User says | Route to |
|---|---|
| `/spectacular grill <doc> [target]` | → resolve and print the target, then `references/grill.md` with `<doc-id>-rules.md` |
| `/spectacular refine <doc> [target]` | → resolve and print the target, then `references/refine.md` with `<doc-id>-rules.md` |
| `/spectacular review <doc> [target]` | → resolve and print the target, then `references/review.md` with `<doc-id>-rules.md` |
| `spectacular <doc>` (no verb) | → load `references/<doc-id>-rules.md`, resolve mode, dispatch |
| `spectacular <doc> grill` | → `references/grill.md` (with `<doc-id>-rules.md` context) |
| `spectacular <doc> refine` | → `references/refine.md` (with `<doc-id>-rules.md` context) |
| `spectacular <doc> review` | → `references/review.md` (with `<doc-id>-rules.md` context) |

**Registered docs:** the live registry is the set of `references/<doc-id>-rules.md` files; the catalog is `references/doc-index.md`. No hardcoded id list here — it drifts. (`spectacular prd …` is just this handler with `<doc> = prd`; bare `prd` → grill if empty, else review.)

### Where does this belong? — soft-DB routing

Deciding *which store* a piece of knowledge goes in (fact? decision? question? research? spike? fix? idea?) → **`references/soft-db-index.md`**, the canonical collection index. Load it whenever the routing isn't obvious.

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
| Scaffolding a new request (`spectacular new`) | → **[[plan-rules]] § 2-of-6 rule** (compact table; canonical: [[verify-authoring]]). Default: no VERIFY.md — `### Verification` group in TASKS.md or PLAN § Validation instead. |
| Grilling/refining a PLAN.md | → **[[plan-rules]] § 2-of-6 rule** — confirm result; ask user if VERIFY.md needed |
| Moving request `active → review` | → `lifecycle.md` § Verification artifact detection — pick artifact (VERIFY.md > TASKS Verification > PLAN Validation) |
| Moving request `review → verified` | → **`verify.md`** — the interactive validation walk (walk-only since b30), record to VERIFY-LOG, gate the transition. **Never skip.** |
| `spectacular verify <slug>` | → **`verify.md`** — the validation walk (skill-only; CLI redirects). |
| `spectacular sweep [<slug>]` / "audit the fleet" | → **`review-sweep.md`** — read-only request-auditor fan-out (review + ticked-active deep; planned batched overlap check). Never promotes — feeds the walk. |
| Automating a shipped scenario | → `verify-authoring.md` § Promoting checks to scripts — when to author `tests/verify/<slug>.test.sh`. |

Verification always runs against *some* artifact — "VERIFY.md is opt-in" means the *file*, never the act; the full doctrine lives in `verify-authoring.md` § "Verification always happens".

---

## State awareness

Load **only** what the task needs (principle 6 — progressive disclosure). Two authorities, no third list:

- **What to load per task type** — `.spectacular/AGENTS.md`'s context-loading table is authoritative; follow it over guessing or re-deriving a read list.
- **How to read state** — prefer the read verbs (§ Cold-start pattern above: `status --brief --json` → named `request` view) over walking the filesystem; the flow docs (`status.md`, `active-request.md`) own their own read steps.

Never read `archive/` during normal operation.

---

## Canonical rules (always apply)

- **Never overwrite canonical documents in place** — snapshot first (`PRD@v1.0.md`). See `references/versioning.md`.
- **Lifecycle state** lives in `PLAN.md` frontmatter (`status: planned | active | review | verified`). TASKS.md mirrors it for skim tooling; PLAN is authoritative — `doctor` repairs drift.
- **Entity lifecycles** are defined only by [[lifecycle-contract]]. Specs are historical execution context; code is authoritative, and only `approved` specs may seed requests.
- **Slugs** are kebab-case, skill-derived, user-overridable, uniqueness enforced.
- **Memory** (`spectacular remember this`) writes to `.spectacular/memories/` — git-committed, team-visible. Never to `.claude/` memory.
- Be proactive: surface stale state, propose lifecycle transitions, flag blocked requests.
- **Execution boundary:** when implementation reveals an unexpected requirement, tangent, or optimization, park it as an idea (or future target) instead of adding it to the active request's PLAN/TASKS. See [[wayfinding-sequencer]].
- **GitHub/Spectacular boundary:** GitHub owns capture, discussion, ownership, PRs, checks, merge, collaboration, and notifications. Spectacular owns reasoning, durable local context, decisions, plans, validation, and coordination. An Issue is a collaborative job card, not automatically a spec or request: use the smallest sufficient `direct | request | spec-first` path and link rather than mirror. See [[github-work-bridge]].
- **Know when to write to a collection, not just how** — the "When to act" trigger table in [[soft-db-index]]. Cheap/reversible writes on their natural trigger; permanent/team-visible writes (memory, decisions, archive) are proposed → human confirms → written, never autonomous.

### Task tracking — two layers

On-disk `requests/<slug>/TASKS.md` owns milestones (persistent, team-visible); harness `TaskCreate`/`TaskUpdate` owns ephemeral session micro-steps (finer-grained — never a one-for-one copy of TASKS.md lines). Full convention: `.spectacular/AGENTS.md` § Task tracking.

---

## Output format

Conversational briefing with a minimal embedded table. Never a raw dump. Identify the single highest-priority next action and ask what the user wants to do.

---

## References & templates index

No hand-list here — it drifts. Reference docs: `references/*.md`, cataloged in `references/doc-index.md`. Templates: the `templates/` tree; frontmatter stubs in `references/scaffold-reference.md`. Projects may override any template at `.spectacular/templates/<doc>/...` (same filenames, project-local wins).
