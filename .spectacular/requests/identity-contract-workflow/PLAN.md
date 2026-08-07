---
status: verified
priority: medium
owner: alex
updated: 2026-08-07
build: b46
docs_impact: none
summary: "Implement UUIDv7 identities, slug paths, and merged-spec-PR request gating"
related:
  - PRD.md
  - ../../specs/SPC-008-durable-identity-and-contract-pr-workflow.md
# Optional traffic evidence; omit until it is confirmed and complete.
# conflicts-with: [<request-slug>]
# traffic-boundaries: [<named, complete launch boundary>]
# release-constraints: [<shared release or migration constraint>]
docs_impact_reason: Workflow references and CLI help were updated; no public docs surface is owned by this repository.
---

# Plan — identity-contract-workflow

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

Implement SPC-008 so durable Spectacular records use branch-safe UUIDv7 identities, human slugs, and a merged-spec-PR approval gate without abandoning Markdown portability.

## Constraints

<!-- What's fixed before you start? Inherited from PRD/STACK/PRINCIPLES + request-specific limits. -->

- Preserve Markdown plus YAML frontmatter as the local, Git-reviewable source of durable workspace knowledge.
- Keep Bash 3.2 compatibility; introduce no database, server, or network requirement for ordinary identity allocation or retrieval.
- Keep the current `spec` and `request` command and path vocabulary in this request; the proposed terminology rename is deferred.
- Treat code and executable tests as current implementation truth. Capability contracts may become stale and must not claim otherwise.
- Make migration preview-first, reference-aware, and reversible through its explicit apply boundary; retain numeric aliases as read-only compatibility inputs.
- Keep forge mechanics in adapters. The core workflow must work locally and keep GitHub, GitLab, and Bitbucket origins as absolute URLs only.

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

The CLI derives durable `NNN` IDs by scanning the checked-out tree, so parallel
branches can allocate the same number. Specs use numbered filenames, and request
provenance duplicates source/spec/version/digest data. The current lifecycle
approves a spec through a mutable frontmatter status rather than its merged PR.

### What changes

Record creation, template validation, list/detail resolution, migration, and
lifecycle gates will adopt UUIDv7 IDs plus slugs. New specs will live at
`specs/<slug>.md`; shared records will carry `kind` and `scope`; and execution
will require the approved spec merge in its branch ancestry.

### What stays the same

Markdown remains canonical and human-editable, existing numeric aliases remain
readable, and the current `spec`/`request` terminology and compatibility commands
remain in place. No forge becomes the canonical spec source.

## Decisions

<!--
  Design calls made inside this request. Format: chose X over Y — because Z.
  Rejected alternatives stay listed; deleting them re-litigates them later.
  Project-wide calls go to DECISIONS.md via `spectacular decide` instead
  (see decisions-rules.md routing table). Empty is fine — no decisions yet.
-->

- Use UUIDv7 as the immutable durable identity and a slug as the human-facing locator — because branch-local counter allocation cannot provide global uniqueness.
- Treat the merged spec PR as approval — because Git history is the shared, reviewable confirmation rather than duplicate frontmatter provenance.
- Keep one shallow typed shared library with `scope` — because request folders should hold execution state, not parallel knowledge databases.

## Milestones

<!-- Ordered, demoable checkpoints. Outcomes, not tasks. -->
<!-- 3-7 milestones for a typical request. Each is something someone can see working. -->

- M1 — New durable records receive UUIDv7 identities and resolve through their human slugs without numeric allocation races.
- M2 — Templates and read views expose the minimal shared frontmatter spine, shared-library scope, and slug-only spec paths.
- M3 — A contract-derived request is accepted only after its spec PR is merged and its execution branch contains that merge.
- M4 — Existing numbered records migrate through a no-write preview and an explicit apply, preserving aliases and repairing duplicate SPC-005.
- M5 — The shipped contract, documentation, and regression suite prove compatibility across lifecycle, retrieval, and migration paths.

## Tasks

<!-- Pointer. The executable checklist lives in TASKS.md, grouped by milestone. -->

See `TASKS.md`.

## Dependencies

<!-- Other requests, skills, blocking decisions. Use [[request-slug]] notation. -->

- Approved contract: [[../../specs/SPC-008-durable-identity-and-contract-pr-workflow]].
- The current CLI templates and lifecycle validator are the migration baseline; their compatibility behavior must remain covered until the migration completes.
- No remote-provider integration is required for the core slice. Forge adapter work remains follow-up scope.

## Validation

<!--
  How each milestone is verified. Per-milestone checks.
  Each check states its AUTHORITY: a run: command, an assertable property,
  a judgable artifact, or a human-observable behavior (see verify.md kinds).
  A check with no authority can't fail. Aspiration verbs (improve, enhance,
  optimize, handle gracefully) are not checks.
-->

- M1 — run: focused identity/concurrency tests create records independently and prove distinct UUIDv7 IDs plus slug resolution.
- M2 — run: frontmatter/template/list-view tests prove `kind`, `scope`, and `slug` retrieval without body loading; an observable new spec path is `specs/<slug>.md`.
- M3 — run: lifecycle tests reject a draft/unmerged contract and an execution branch missing the approved merge, then accept the merged ancestry path.
- M4 — run: migration dry-run produces complete ID/path/reference mappings with no writes; apply leaves no duplicate durable IDs and resolves each legacy alias.
- M5 — run: `bash tests/run.sh`, `bash -n cli/spectacular cli/install.sh scripts/hooks/pre-commit`, `scripts/hooks/pre-commit --check`, and relevant `spectacular doctor` areas exit successfully.

## Deliverables

<!-- Artifacts that ship out of this request. Concrete files, docs, behaviors. -->

- UUIDv7/slug identity, resolver, and compatibility-alias substrate in the CLI.
- Updated templates, frontmatter validation, shared-library retrieval views, and spec path behavior.
- Merged-spec-PR approval and execution-branch ancestry gate, with forge scope configuration.
- Preview-first numbered-ID migration plus a mapping receipt and explicit duplicate-SPC-005 repair.
- Updated architecture/lifecycle/skill documentation and regression coverage.
