---
status: archived
priority: medium
owner: alex
updated: 2026-08-02
build: b38
docs_impact: required
summary: "Define live, stale-safe, temporary, and throwaway state with archive and synchronization rules"
related:
  - PRD.md
docs_impact_evidence: Updated README.md, docs/workflow.md, docs/scaffold.md, docs/commands.md, architecture, system spec, and lifecycle/retention references
archived: 2026-08-02
---

# Plan — artifact-retention-contract

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

Define which Spectacular and codebase artifacts must stay synchronized, may remain stale safely, are bounded temporary work, or may be garbage-collected after their durable learning is preserved.

## Constraints

<!-- What's fixed before you start? Inherited from PRD/STACK/PRINCIPLES + request-specific limits. -->

- Retention class is derived from entity, status, and path; no bulk frontmatter migration or daily freshness ritual.
- Production code and executable invariant/unit tests are the ultimate implementation truth after verified integration.
- Live planning surfaces must be trustworthy enough to sequence work: roadmap index, active requests, active-release specs before code generation, and unresolved questions.
- Historical Markdown is archived rather than purged; archive stays outside normal agent context.
- Disposable code may be deleted only after outcome/evidence and a recovery boundary are recorded.
- A resolved question produces a DEC only when the answer is a real choice between alternatives.

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

The roadmap already uses `roadmaps/index.md` as its canonical live entry and prunes shipped prose into `roadmaps/vX.Y.Z.md`. Open questions are surfaced by explicit session start, but resolved files remain in the active collection. Detailed specs can archive only after supersession/deprecation, while AFK spike cleanup already creates a hidden Git archive ref before deleting a local branch.

### What changes

Create one retention contract defining `live`, `stale-safe`, `temporary`, and `throwaway` classes plus synchronization, loading, and disposal rules. Make every status/onboarding session surface open user questions first. Move resolved questions into `archive/questions/` while keeping canonical dependencies satisfiable. Permit explicit archive-first closure of implemented, rejected, or abandoned detailed specs. Align user docs and lifecycle rules.

### What stays the same

No `ROADMAP_ARCHIVE.md`, retention frontmatter field, background synchronizer, automatic DEC creation, remote branch deletion, or hard deletion of durable Markdown. Decisions, findings, research outcomes, and indexes remain committed history. Code/test currency checks remain on-demand when historical specs or docs are reused.

## Decisions

<!--
  Design calls made inside this request. Format: chose X over Y — because Z.
  Rejected alternatives stay listed; deleting them re-litigates them later.
  Project-wide calls go to DECISIONS.md via `spectacular decide` instead
  (see decisions-rules.md routing table). Empty is fine — no decisions yet.
-->

- Chose four derived retention classes over per-file `retention:` metadata because the lifecycle already supplies the necessary signal and duplicated labels would drift.
- Chose `roadmaps/index.md` plus `roadmaps/vX.Y.Z.md` over `ROADMAP.md` plus `ROADMAP_ARCHIVE.md` because index mode already bounds prompt cost without parallel canonical entry points.
- Chose archive-on-question-resolution over leaving resolved blockers in `questions/` because that folder should represent only actionable fog.
- Chose optional decision linkage over automatic QUE→DEC conversion because factual answers and rejected assumptions are not architectural choices.
- Chose archive-first spec closure over purge because rejected reasoning remains valuable but should not pollute active context.
- Chose hidden verified Git archive refs before local spike-branch deletion because they provide recovery without appearing in normal branch or agent scans.

## Milestones

<!-- Ordered, demoable checkpoints. Outcomes, not tasks. -->
<!-- 3-7 milestones for a typical request. Each is something someone can see working. -->

- M1 — One canonical matrix assigns every major workspace/code artifact a freshness obligation, loading rule, and terminal disposition.
- M2 — Question and specification mutators enforce archive-first terminal behavior without breaking canonical dependency resolution.
- M3 — Session briefings, roadmap guidance, tests, and user-facing documentation agree on the same retention model.

## Tasks

<!-- Pointer. The executable checklist lives in TASKS.md, grouped by milestone. -->

See `TASKS.md`.

## Dependencies

<!-- Other requests, skills, blocking decisions. Use [[request-slug]] notation. -->

- Builds on [[lifecycle-contract]], [[discovery-evidence-protocol]], and [[afk-git-hygiene]].
- Preserves [[advanced-engineering-collections]] reservation-only status for `FND` while defining its future retention class.

## Validation

<!--
  How each milestone is verified. Per-milestone checks.
  Each check states its AUTHORITY: a run: command, an assertable property,
  a judgable artifact, or a human-observable behavior (see verify.md kinds).
  A check with no authority can't fail. Aspiration verbs (improve, enhance,
  optimize, handle gracefully) are not checks.
-->

- M1 — judge: each artifact has exactly one derived class and no class implies false continuous freshness.
- M2 — run: lifecycle/wayfinding tests prove resolved questions leave active context, remain valid satisfied dependencies, and implemented/draft/unconfirmed specs archive only through dry-run-first explicit apply.
- M3 — run: Bash syntax, version guard, focused tests, full suite, targeted doctor, and `git diff --check` pass.

## Deliverables

<!-- Artifacts that ship out of this request. Concrete files, docs, behaviors. -->

- Canonical `artifact-retention.md` contract and skill routing.
- Archive-first question/spec behavior with regression coverage.
- Aligned lifecycle, roadmap, discovery, architecture, system-spec, command, scaffold, workflow, and README guidance.
