---
status: verified
priority: high
owner: alex
updated: 2026-08-04
build: b43
docs_impact: required
summary: "Implement SPC-004: Pre-request Vision workspaces turn grounded imagination, human-reviewed fragments, prototypes, and spike evidence into an explicitly approved direction before specification and execution planning"
related:
  - PRD.md
# Optional traffic evidence; omit until it is confirmed and complete.
# conflicts-with: [<request-slug>]
# traffic-boundaries: [<named, complete launch boundary>]
# release-constraints: [<shared release or migration constraint>]
source_spec: SPC-004
source_type: spec
source_ref: SPC-004
source_spec_version: 1.0
source_spec_digest: "sha256:38d28f566af0fa21105192ad109978bb8dc7fb0c93c2aac9614dc8ad3fd82310"
scaffolded_against: a5b4d5fbf3073c8e4f87de24255d5959d1e3f0ea
activated_at: 2026-08-04
activated_by: alex
activated_against: a5b4d5fbf3073c8e4f87de24255d5959d1e3f0ea
docs_impact_evidence: .spectacular/requests/vision-workflow/artifacts/pageworks-handoff.md — audited README.md and docs/{workflow,commands,scaffold,visual-conventions,troubleshooting,integrations}.md; Pageworks-owned rewrite deferred with exact acceptance checks
github_pr: "https://github.com/alexsmedile/spectacular/pull/22"
github_pr_opened_at: 2026-08-04
github_pr_head: 6431697f531b11a08d011aa6d3df4caef7b737f6
issue_resolution: on_merge
---

# Plan — vision-workflow

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

Advance the PRD goal of human-agent coherence by giving Spectacular an opt-in pre-request Vision workflow that turns grounded imagination and human-reviewed fragments into an explicitly approved direction from which an SPC can be derived.

## Constraints

- Stay within the approved SPC-004 requirements; unresolved scope changes return to discovery.
- Preserve legacy `requests/<slug>/vision/` workspaces without automatic movement or rewriting.
- Keep Vision optional; no new gate for clearly specified backend, maintenance, or direct work.
- Keep evidence truth (`supported | refuted | inconclusive`) separate from human direction approval.
- Bash 3.2 compatibility; no associative arrays or new runtime dependencies.
- Public `docs/` authoring belongs to Pageworks; this request assesses impact and produces a handoff when needed.

## Understanding

<!--
  OPTIONAL authoring slot, but REQUIRED before `planned → active` by the
  `understand-before-change` policy (@Implementation). Fill it here for a
  typical request; escalate to a dedicated requests/<slug>/UNDERSTANDING.md
  (same three subheads) for large ones — the policy is satisfied by EITHER.
  Not one of the 7 required authoring slots; it gates implementation, not planning.
-->

### How it works now

`spectacular imagine <slug>` requires an existing request, scaffolds `requests/<slug>/vision/`, mandates one story/UI/architecture fragment, records per-fragment `approved:` state, and derives directly to PLAN. Ideas can promote directly to requests. Research and spikes own evidence separately. `experiment` is a hidden alias for feedback-loop. PLAN Understanding is a later pre-implementation gate. The roadmap still presents `prototype` as a version phase.

### What changes

New Vision workspaces live under `.spectacular/visions/<slug>/` with a lifecycle-bearing `VISION.md`, proportional typed fragments, evidence links, explicit fragment reactions, whole-Vision approval, and an agentic approved-Vision-to-draft-SPC handoff. CLI, doctor, skill routing, templates, lifecycle/discovery contracts, and tests change together.

### What stays the same

SPCs remain the implementation authorization contract; PLAN/TASKS remain request execution artifacts; PLAN Understanding remains code-grounded; research and spikes retain their own evidence lifecycles; feedback-loop remains post-build learning; prototypes remain owned artifacts rather than canonical entities; legacy request visions remain readable.

## Decisions

<!--
  Design calls made inside this request. Format: chose X over Y — because Z.
  Rejected alternatives stay listed; deleting them re-litigates them later.
  Project-wide calls go to DECISIONS.md via `spectacular decide` instead
  (see decisions-rules.md routing table). Empty is fine — no decisions yet.
-->

- Chose pre-request `.spectacular/visions/<slug>/` over request-only Vision because human direction should precede SPC and PLAN derivation.
- Chose one `imagine` verb over separate Dream/Imagine modes because dreaming is human intent and imagining is the agent operation.
- Chose optional typed fragments over one-of-each scaffolding because the uncertainty should determine the artifact set.
- Chose linked RES/SPK evidence over copying evidence into Vision because truth and approval have different owners.
- Chose agentic Vision-to-SPC derivation over CLI text synthesis because derivation requires judgment and only approved fragments are load-bearing.
- Chose read-compatible legacy support over automatic migration because historical reaction records should not be rewritten without need.

## Milestones

- M1 — Pre-request Vision contract and templates define lifecycle, spine, fragments, and compatibility.
- M2 — CLI scaffolds, reads, reacts to, proposes, approves, and routes derivation for Vision workspaces.
- M3 — Discovery, Idea, Spike, Feedback, roadmap, specification, and request workflows compose around approved Vision.
- M4 — Doctor and focused regression tests prove lifecycle gates, manifest repair, legacy support, and Bash compatibility.
- M5 — Canonical project truth and documentation impact are synchronized; request verification evidence is complete.

## Tasks

<!-- Pointer. The executable checklist lives in TASKS.md, grouped by milestone. -->

See `TASKS.md`.

## Dependencies

- Approved source specification: SPC-004.
- Linked decisions and questions remain in the source specification.
## Validation

- M1 — assert: templates and reference contracts contain every required Vision section/state/kind and no mandatory fragment-kind rule remains.
- M2 — `run: bash tests/cli/vision.test.sh`; a fresh workspace can complete scaffold → fragments → reactions → propose → approve, while derive refuses non-approved Vision.
- M3 — assert: routing contains one authoritative role per Idea/Imagine/Vision/prototype/spike/feedback/spec/PLAN concept; `experiment` no longer dispatches to feedback-loop.
- M4 — `run: bash -n cli/spectacular` and `run: bash tests/run.sh`; doctor detects malformed/new/legacy Vision state and fixes only manifest drift.
- M5 — `run: scripts/hooks/pre-commit --check` and `run: ./cli/spectacular doctor`; canonical docs describe shipped behavior and public-doc impact has evidence or a Pageworks handoff.

## Deliverables

- Updated Vision/Imagine CLI and skill workflow.
- New pre-request Vision templates and compatibility-aware doctor behavior.
- Focused CLI regression coverage.
- Synchronized lifecycle, discovery, routing, scaffold, architecture, and system-spec contracts.
- Verification log and public-documentation impact/handoff.
