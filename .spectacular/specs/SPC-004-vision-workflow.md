---
id: SPC-004
type: specification
status: implemented
target_version: v1.38.0-execution
supersedes: ""
updated: 2026-08-04
summary: "Pre-request Vision workspaces turn grounded imagination, human-reviewed fragments, prototypes, and spike evidence into an explicitly approved direction before specification and execution planning"
related: []
version: 1.1
approved_at: 2026-08-04
approved_by: user
implemented_at: 2026-08-04
verified_against: b43
---

# SPC-004 — Pre-request Vision workspaces turn grounded imagination, human-reviewed fragments, prototypes, and spike evidence into an explicitly approved direction before specification and execution planning

## Intent

Spectacular must let a human and agent understand a prospective change, imagine
concrete strategies and experiences, react to reviewable fragments, and approve
a coherent north-star direction before a specification, request, or PLAN exists.
The durable result is an approved pre-request `VISION.md`; it is not executable
authorization. Approved Vision becomes the preferred source for a draft SPC,
while research and spikes remain separately-owned evidence and prototypes remain
human-reactable artifacts.

## Requirements

### R1 — Pre-request ownership

- New Vision workspaces live at `.spectacular/visions/<slug>/`, independently of
  requests and specifications.
- `spectacular imagine <slug>` scaffolds the workspace without requiring a
  request or PLAN.
- Existing `requests/<slug>/vision/` workspaces remain readable and diagnosable;
  the refactor does not move or rewrite them automatically.

### R2 — Explicit workflow and vocabulary

- The canonical loop is **Understand → Imagine → Probe → React → Confirm →
  Derive**.
- The human supplies aspiration, taste, and authority; the agent imagines
  concrete alternatives. `dream` remains natural language, not a second command.
- `prototype` is a showable artifact, not a standalone entity or required phase.
- `experiment` is not a hidden alias for feedback-loop. Natural-language
  experiments route by question: research for facts, spike for technical
  feasibility, Vision fragment/prototype for pre-spec human reaction, and
  feedback-loop for post-build learning.

### R3 — Vision spine and lifecycle

- `VISION.md` contains Intent, North star, Understanding, Experience signature,
  Strategies considered, Chosen direction, Boundaries, Evidence, Fragment
  manifest, and Approval.
- Vision status is `draft → proposed → approved`, with rejection allowed from
  draft/proposed. `approved` records `approved_by` and `approved_at`.
- Approval is an explicit human gate. Approved Vision may derive a draft
  specification but never authorizes production code or activates a request.

### R4 — Proportional fragments

- Fragments live under `fragments/` and support `strategy`, `story`, `flow`,
  `ui`, `arch`, and `prototype` kinds.
- No fragment kind is mandatory. The agent creates the cheapest set that makes
  the material uncertainty discussable.
- Every proposal fragment has a human reaction of `pending`, `approved`,
  `revise`, `rejected`, or `superseded`, plus an optional reaction note.
- Only approved fragments are load-bearing during specification derivation.

### R5 — Evidence boundary

- `evidence/` holds links or concise presentation artifacts, not copies of RES
  or SPK records.
- Research owns sourced facts. Spikes own hypothesis, experiment, evidence, and
  `supported | refuted | inconclusive` feasibility results.
- A spike result never approves a product direction. For Vision-linked spikes,
  the agent should present a showable conclusion for human reaction.

### R6 — Derivation and understanding boundary

- `spectacular vision derive <slug>` routes to the agentic derivation flow and
  requires an approved Vision.
- Derivation creates or updates a **draft SPC**, using only the approved Vision,
  approved fragments, and linked evidence. The SPC still requires its normal
  explicit approval.
- Vision Understanding covers current reality, users/needs, constraints, and
  material uncertainties. PLAN/UNDERSTANDING retains the later code-grounded
  implementation delta and is not automatically satisfied by Vision approval.

### R7 — CLI, doctor, and compatibility

- CLI supports `vision list`, `vision show`, `vision add`, `vision react`,
  `vision propose`, `vision approve`, and agentic `vision derive` routing.
- `doctor vision` validates both new and legacy locations, lifecycle vocabulary,
  required spine sections, fragment schema, manifest drift, and approval
  coherence. Mechanical fix remains limited to manifest regeneration.
- `spectacular paths` exposes the top-level visions collection.
- CLI behavior remains Bash 3.2-compatible.

### R8 — Routing composition

- Idea is cheap capture. When its destination remains unsettled, shaping routes
  to Imagine/Vision; direct promotion to SPC/request remains appropriate only
  when the execution destination is already accepted.
- The roadmap recommends `direction-validation` as the broad discovery phase;
  legacy `prototype` remains accepted for compatibility but is no longer the
  preferred phase name.
- Feedback-loop remains post-build/prototyping learning and does not own
  pre-spec approval.

### R9 — Verification and documentation

- Focused CLI tests cover pre-request scaffolding, fragment reactions, proposal
  and approval gates, derivation routing, doctor checks/fix, legacy discovery,
  experiment-alias removal, and Bash 3.2 syntax.
- Skill references, templates, help text, scaffold documentation, lifecycle,
  discovery routing, and canonical architecture/spec index are synchronized.
- Public `docs/` impact is assessed separately and handed to Pageworks when the
  user-facing rewrite is substantial.

## Evidence and decisions

- GitHub issue #7 asks for a per-request UX/experience-layer vibe check.
- The current `imagine` engine already provides per-fragment reaction, but it
  requires a request, mandates story/UI/architecture fragments, and derives
  directly to PLAN.
- SPC-002 established approved-spec-first request derivation, making the older
  Imagine-to-PLAN handoff temporally inconsistent.
- User decision (2026-08-04): proceed with the first-class pre-request Vision
  model, keep it opt-in, preserve Understanding, and implement through PR.

## Confirmation

Drafted from the explicitly accepted Vision workflow proposal. Approval must
confirm the requirements above without adding a mandatory gate to clearly
specified backend, maintenance, or direct work.

**Approved 2026-08-04 by user** — Explicitly approved the proposed pre-request Vision model and authorized autonomous implementation through a PR
