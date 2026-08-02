---
status: archived
priority: medium
owner: alex
updated: 2026-08-02
build: b36
docs_impact: required
summary: "Reserve optional findings, fixes, bugs, security, and benchmark collections in init without activating unsupported workflows"
related:
  - PRD.md
docs_impact_evidence: Updated README.md, docs/commands.md, docs/scaffold.md, init workflow, architecture, system spec, and canonical collection contracts
archived: 2026-08-02
---

# Plan — advanced-engineering-collections

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

Let advanced projects opt into reserved findings, fixes, bugs, security, and benchmark collections during `spectacular init` without adding empty senior-engineering ceremony to lightweight workspaces or claiming unfinished workflows exist.

## Constraints

<!-- What's fixed before you start? Inherited from PRD/STACK/PRINCIPLES + request-specific limits. -->

- The default and kit-driven scaffold stay lightweight; every advanced collection is opt-in.
- Canonical folder names are `findings/`, `fixes/`, `bugs/`, `security/`, and `benchmarks/`; `security/` is a domain name, while `securities/` would mean financial instruments.
- Existing verified fixes use legacy `F<N>` IDs and remain readable. `FND` cannot claim `f1` until an explicit, previewed fix-ID migration removes that ambiguity.
- Reservation defines paths, prefixes, and intent only. It does not invent lifecycle states, mutators, doctor areas, or autonomous behavior.
- Init remains idempotent and Bash 3.2-compatible; existing collection contents are never overwritten.

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

`--with` accepts opt-in root docs and a few collection substrates, but the advanced engineering stores have no reserved init surface. The active fixes ledger writes `F1`, `F2`, … under `fixes/`, which conflicts with the proposed shorthand `f1` for `FND-001`.

### What changes

Allow `spectacular init --with findings,fixes,bugs,security,benchmarks`. Each selection creates an idempotent, Git-preservable collection directory. Document the future canonical prefixes `FND`, `FIX`, `BUG`, `SEC`, and `BMK`, with collision-safe aliases and explicit reserved status.

### What stays the same

The existing fix command, `F<N>` records, audit/debug pipeline, default scaffold, and all collection lifecycles remain unchanged. No migration runs and no record is allocated in a reserved collection.

## Decisions

<!--
  Design calls made inside this request. Format: chose X over Y — because Z.
  Rejected alternatives stay listed; deleting them re-litigates them later.
  Project-wide calls go to DECISIONS.md via `spectacular decide` instead
  (see decisions-rules.md routing table). Empty is fine — no decisions yet.
-->

- Chose `security/` over `securities/` because it names the engineering domain correctly.
- Chose explicit `--with` collection IDs over adding five folders to every coding workspace because advanced ceremony should be demand-driven.
- Chose reservation-only indexes over premature mutators and lifecycles because names and paths can stabilize before behavior is designed.
- Chose `fnd1` as the safe finding alias for now; the proposed `f1` remains reserved-but-unavailable while legacy fix `F<N>` records exist.
- Chose `bug1` as the safe bug alias; naked `b1` remains context-only because roadmap build IDs also use `b<N>`.

## Milestones

<!-- Ordered, demoable checkpoints. Outcomes, not tasks. -->
<!-- 3-7 milestones for a typical request. Each is something someone can see working. -->

- M1 — Init can scaffold any or all five advanced collection folders explicitly and idempotently.
- M2 — Canonical IDs, aliases, paths, and compatibility caveats are documented without changing existing fix behavior.
- M3 — Tests and user-facing scaffold documentation prove the default remains lean and advanced opt-in is discoverable.

## Tasks

<!-- Pointer. The executable checklist lives in TASKS.md, grouped by milestone. -->

See `TASKS.md`.

## Dependencies

<!-- Other requests, skills, blocking decisions. Use [[request-slug]] notation. -->

- Builds on verified [[lifecycle-contract]] and its unique-ID, migration, and no-false-claims rules.
- Existing `fixes/F<N>.md` compatibility blocks activation of the `f1` finding alias until a separately approved migration.

## Validation

<!--
  How each milestone is verified. Per-milestone checks.
  Each check states its AUTHORITY: a run: command, an assertable property,
  a judgable artifact, or a human-observable behavior (see verify.md kinds).
  A check with no authority can't fail. Aspiration verbs (improve, enhance,
  optimize, handle gracefully) are not checks.
-->

- M1 — run: `bash tests/cli/init.test.sh` proves explicit collection selection, all-five selection, idempotence, and unchanged bare init.
- M2 — assert: canonical ID and lifecycle references label new entities reserved and preserve legacy `F<N>` compatibility without ambiguous aliases.
- M3 — run: Bash syntax, version guard, focused init tests, and full suite pass; user-facing commands/scaffold docs list the opt-in surface.

## Deliverables

<!-- Artifacts that ship out of this request. Concrete files, docs, behaviors. -->

- Optional init scaffolders for `findings/`, `fixes/`, `bugs/`, `security/`, and `benchmarks/`.
- Reserved canonical ID/path contract for `FND`, `FIX`, `BUG`, `SEC`, and `BMK`.
- Regression tests and aligned architecture, system-spec, command, and scaffold documentation.
