---
status: archived
priority: high
owner: alex
updated: 2026-08-02
build: b39
docs_impact: required
summary: "Make request creation, context retrieval, activation, verification, and command grammar coherent for humans and coding agents"
source_spec: SPC-002
source_spec_version: "1.0"
source_spec_digest: "sha256:2a54ef9e2ee2255804b01f1561ee9a594883f92137a2b9087cb501768f7cf43a"
activated_at: 2026-08-02
activated_by: user
activated_against: "325d7be53ccb43e5d2fe23a23e2a0bce21fb6bbd"
related:
  - PRD.md
  - ../../specs/SPC-002-request-workflow-interface.md
docs_impact_evidence: README, docs/commands.md, docs/workflow.md, docs/scaffold.md, skill references, PRD, architecture, system spec, and changelog updated
archived: 2026-08-02
---

# Plan — request-workflow-interface

<!--
  Canonical 7-slot PLAN template for a single request.
  Lives at: .spectacular/requests/<slug>/PLAN.md

  Rules:
  - PLAN is per-request. PRD is project-wide. Never put a PRD inside requests/.
  - This file's frontmatter `status:` is the single source of lifecycle state for the request.
  - The 7 required sections must appear IN ORDER, unnumbered:
      ## Goal, ## Constraints, ## Milestones, ## Tasks, ## Dependencies, ## Validation, ## Deliverables
    Extra sections (## Understanding, ## Decisions, request-specific) may appear
    BETWEEN them; doctor enforces the required set's presence + order, not a closed list.
  - All 7 required sections must be filled before this PLAN is considered usable.
  - Replace every <placeholder> with concrete content.
-->

## Goal

<!-- One sentence. What does this request change? -->
<!-- Compress the request's intent. Aligns with PRD's Vision or Goals — this is a slice, not a restatement. -->

Advance the PRD goals to reduce context rot and improve agent execution through structured retrieval by giving humans and coding agents one predictable request workflow that turns an approved specification into bounded implementation context, native session planning, verified execution, and clean archival without duplicating project history.

## Constraints

<!-- What's fixed before you start? Inherited from PRD/STACK/PRINCIPLES + request-specific limits. -->

- Spectacular owns durable request intent, lifecycle, milestones, cross-session state, and evidence; Codex/Claude built-in planning owns current-session micro-steps; subagents receive still narrower closed briefs.
- `PLAN.md`, `TASKS.md`, and `SESSION.md` remain the durable request sources. The implementation brief is generated on demand and is never stored as another drifting Markdown file.
- CLI entity operations use noun-first canonical grammar. Agentic document operations use verb-first grammar. Existing commands remain compatibility aliases during this request.
- `spectacular request <slug>` stays the overview; `--brief` returns authorized implementation context; `--full` returns the request-owned Markdown bundle; `--json` remains deterministic and machine-readable.
- The first slice supports milestone selection but does not add an arbitrary artifact/section query language. Agents can read a named file directly when the compiled brief is insufficient.
- `spectacular request new --from SPC-001` is mechanical scaffolding. `/spectacular spec act SPC-001` is the agentic authorization and execution flow.
- Only an approved specification may seed spec-driven implementation. Activation records the exact specification version and Git baseline without copying the specification into the request.
- `review → verified` remains evidence-gated through `/spectacular verify`; changing frontmatter alone is never equivalent to verification.
- Bash 3.2 compatibility, archive-first behavior, existing HITL gates, code/tests as implementation truth, and the execution scope boundary remain intact.

<!-- ## Understanding and ## Decisions below are OPTIONAL extra sections,
     allowed between Constraints and Milestones. -->

## Understanding

<!--
  OPTIONAL authoring slot, but REQUIRED before `planned → active` by the
  `understand-before-change` policy (@Implementation). Fill it here for a
  typical request; escalate to a dedicated requests/<slug>/UNDERSTANDING.md
  (same three subheads) for large ones — the policy is satisfied by EITHER.
  Not one of the 7 required authoring slots; it gates implementation, not planning.
-->

### How it works now

`spectacular new` creates PLAN and TASKS, while `spec act` currently does little more than call that scaffold and attach `source_spec`. During implementation the skill loads full request documents and manually reconstructs a five-slot milestone brief from TASKS, PLAN constraints, validation, and deliverables. `spectacular request <slug>` is only an index card; `--full` dumps PLAN alone. Workspace orientation is split across status, summary, next, requests, request, and progress. Lifecycle mutations and doc-writing verbs also mix noun-first, verb-first, CLI-only, and skill-only conventions.

### What changes

Make `request` the coherent execution namespace. Add generated overview, implementation brief, full-bundle, milestone, and JSON views; derive planned requests from approved specifications; record activation provenance; and make `/spectacular spec act` perform the gated agentic handoff into native Codex/Claude planning. Establish noun-first canonical CLI commands, verb-first conversational document commands, explicit transition behavior, compatibility aliases, and clearer advanced-command help without introducing needless selector or AFK subcommand layers.

### What stays the same

The five-state request lifecycle, the PLAN/TASKS bundle, Wayfinder discovery entities, specification evidence rules, archive closure gate, AFK authorization boundaries, and code/test authority remain unchanged. Existing command spellings keep working. Spectacular does not persist the native agent plan, infer product decisions from code, automatically merge, or load archived history during implementation.

## Decisions

<!--
  Design calls made inside this request. Format: chose X over Y — because Z.
  Rejected alternatives stay listed; deleting them re-litigates them later.
  Project-wide calls go to DECISIONS.md via `spectacular decide` instead
  (see decisions-rules.md routing table). Empty is fine — no decisions yet.
-->

- Chose a generated `request --brief` over a stored `BRIEF.md` because the prompt must reflect live request, session, Git, and code state without creating another synchronization obligation.
- Chose durable TASKS milestones plus native agent session plans plus narrow subagent briefs over one shared task granularity because each layer has a different lifetime and audience.
- Chose noun-first canonical CLI grammar with compatibility aliases over removing familiar top-level shortcuts because predictability can improve without breaking existing users.
- Chose verb-first agentic document commands over advertising two equal grammars because conversational actions read naturally while one canonical form prevents routing ambiguity.
- Chose `request new --from SPC-001` for mechanical scaffolding and `/spectacular spec act SPC-001` for agentic execution over letting one terminal command pretend to perform planning and authorization it cannot do.
- Chose overview, brief, full, milestone, and JSON views for the first slice over `--artifact` plus `--section` selectors because the compiled implementation prompt solves the real context problem without building a query language prematurely.
- Chose explicit activation provenance fields over copying the approved specification because an immutable version/commit baseline is reproducible and cheaper than duplicated content.
- Chose `/spectacular verify` as the normal owner of `review → verified` over unrestricted `advance --to verified` because verification evidence, not a status edit, earns the transition.
- Chose to keep `spectacular afk cleanup` under the existing AFK namespace and clarify that it cleans local branches over adding `afk branch cleanup` before another AFK cleanup type exists.
- Chose to fold documentation-impact assessment into verify/archive guidance while retaining its low-level compatibility command over making internal closure bookkeeping a prominent everyday command.

## Milestones

<!-- Ordered, demoable checkpoints. Outcomes, not tasks. -->
<!-- 3-7 milestones for a typical request. Each is something someone can see working. -->

- M1 — One documented request-workflow contract defines canonical grammar, view semantics, phase ownership, activation provenance, compatibility behavior, and transition gates.
- M2 — `spectacular request <slug>` provides deterministic overview, implementation brief, full-bundle, milestone, and JSON views without reading unrelated or archived context.
- M3 — An approved specification can generate a reviewed planned request, and `/spectacular spec act` safely turns that request into a baseline-anchored native-agent implementation session.
- M4 — Canonical noun-first request commands, verb-first agentic document commands, verification ownership, and compatibility aliases behave consistently without expanding the everyday command surface.
- M5 — Tests, doctor checks, skill routing, architecture, system spec, and user-facing docs teach and enforce the same lean workflow.

## Tasks

<!-- Pointer. The executable checklist lives in TASKS.md, grouped by milestone. -->

See `TASKS.md`.

## Dependencies

<!-- Other requests, skills, blocking decisions. Use [[request-slug]] notation. -->

- Executes approved `SPC-002`, which extends SPC-001's approved-spec execution handoff and execution-boundary requirements with the confirmed command and context contract.
- Builds on verified [[wayfinding-contract]], [[wayfinding-sequencer]], [[afk-git-hygiene]], and [[lifecycle-contract]].
- Must preserve the existing read-command, policy, verification, archive, undo, and Bash 3.2 substrates while rationalizing their public grammar.

## Validation

<!--
  How each milestone is verified. Per-milestone checks.
  Each check states its AUTHORITY: a run: command, an assertable property,
  a judgable artifact, or a human-observable behavior (see verify.md kinds).
  A check with no authority can't fail. Aspiration verbs (improve, enhance,
  optimize, handle gracefully) are not checks.
-->

- M1 — assert: one canonical reference defines every command as read-only, mechanical mutation, or agentic workflow; aliases and non-goals are explicit and no two canonical commands claim the same responsibility.
- M2 — run: focused request-view tests prove overview, brief authorization, full-bundle ordering, `-m 2`/`-m2`/`--milestone M2` normalization, JSON stability, missing-artifact handling, and archive exclusion.
- M3 — run: focused spec-handoff tests prove approved-only `--from`, derived PLAN/TASKS content, duplicate refusal, flat activation baseline fields, blocker enforcement, and skill-only `spec act` routing into native planning.
- M4 — run: command-routing tests prove noun-first canonical forms and old aliases reach the same mutators, document verbs resolve and display their target, direct verified transitions cannot bypass evidence, and AFK/docs-impact behavior is unchanged.
- M5 — run: `bash -n cli/spectacular cli/install.sh scripts/hooks/pre-commit`, `scripts/hooks/pre-commit --check`, focused tests, full `bash tests/run.sh`, and scoped doctor checks all pass; observable: README and command/workflow docs present the same concise happy path.

## Deliverables

<!-- Artifacts that ship out of this request. Concrete files, docs, behaviors. -->

- Approved specification revision or successor defining the request execution-context contract.
- Bash 3.2-compatible request views, spec-derived scaffolding, activation-baseline helpers, canonical namespace routing, compatibility aliases, and verification guards.
- Agentic routing for verb-first document operations and `/spectacular spec act` native-plan handoff.
- Focused regression tests for request reads, milestone normalization, spec handoff, lifecycle gates, aliases, and machine-readable output.
- Updated skill references, architecture/system spec, CLI help, workflow/command/scaffold documentation, and request closure evidence.
