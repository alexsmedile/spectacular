# Pageworks handoff — Vision workflow

Update Spectacular's user-facing documentation for the pre-request Vision
workflow implemented by SPC-004. Work from the implementation on branch
`codex/vision-workflow-refactor`; do not invent a second Dream, prototype, or
experiment lifecycle.

## Core message

- Human dreams; the agent `imagine`s; a durable `VISION.md` captures the agreed
  direction.
- The opt-in loop is **Understand → Imagine → Probe → React → Confirm → Derive**.
- Suggest it only when product, interaction, workflow, experience, or system
  shape uncertainty is material. Do not add a gate for clearly specified UI,
  backend, maintenance, migration, or direct work.
- New work lives at `.spectacular/visions/<slug>/` with `VISION.md`, flat typed
  `fragments/`, and `evidence/`.
- Fragment kinds are `strategy|story|flow|ui|arch|prototype`; none is mandatory.
- Each fragment receives `pending|approved|revise|rejected|superseded`; then the
  whole Vision receives explicit human approval with actor/date.
- An approved Vision derives a **draft SPC**, never PLAN or implementation.
  Approved SPC remains the source of request PLAN/TASKS.
- Prototype is a showable fragment, not an entity or required phase. Spike owns
  technical feasibility truth. Feedback-loop owns post-build learning.
- Legacy `requests/<slug>/vision/` remains readable/diagnosable; do not document
  an automatic migration.

## Pages to update

1. `docs/workflow.md`
   - Insert Vision between unsettled Idea/discovery and draft SPC:
     `idea/conversation → understand/imagine/probe → reacted fragments → approved Vision → draft SPC → approved SPC → request`.
   - Replace the current “prototype attached to request, vision, or feedback”
     row with the pre-spec Vision prototype vs post-build feedback distinction.
   - Add a short “when to suggest / when to skip” rule and one non-visual example.

2. `docs/commands.md`
   - Add the complete mechanical command surface:
     `imagine`, `vision list|show|add|react|propose|approve|reject|derive`.
   - Mark `vision derive` as an agentic redirect requiring approved Vision.
   - Explain that `experiment` is intentionally not a feedback alias.
   - Correct any table that presents `PRT` as a live prototype identity; it is
     reserved.

3. `docs/scaffold.md`
   - Remove `VISION.md` from the per-request file table.
   - Add top-level `visions/` to folder conventions, directory tree, and
     creation rules; show `VISION.md`, `fragments/`, and `evidence/`.
   - Replace request/vision ownership language for prototypes with Vision-owned
     pre-spec fragments and feedback-owned post-build artifacts.

4. `docs/visual-conventions.md`
   - Update the stored fragment path from
     `requests/<slug>/vision/ui/<name>.md` to
     `visions/<slug>/fragments/<name>.md`.
   - Explain that ASCII UI is one optional fragment type; strategies, flows,
     architecture, and prototypes may be non-visual or runnable.
   - Keep roadmap precision tier `vision` distinct from the Vision workflow.

5. `docs/troubleshooting.md`
   - Clarify that `imagine` scaffolding and Vision reaction/lifecycle mutations
     are mechanical CLI operations, while generative Imagine and `vision derive`
     require the skill.
   - Add a legacy-location note and `doctor vision --fix` boundary (manifest
     only; no approval/reaction repair).

6. `README.md`
   - Add `visions/` to the workspace tree and on-demand folder list.
   - Replace “attached prototype” discovery copy with the approved-Vision model.
   - Add one compact sentence connecting approved Vision → draft SPC → approved
     SPC → request.

7. Review `docs/integrations.md` for the Issue `spec-first` path. Mention Vision
   only when Issue direction is genuinely unsettled; do not make it universal.

## Source of truth

- `skills/spectacular/references/imagine.md`
- `skills/spectacular/references/vision-rules.md`
- `skills/spectacular/references/discovery-protocol.md`
- `skills/spectacular/references/spike-rules.md`
- `skills/spectacular/references/feedback-loop.md`
- `skills/spectacular/references/spec-lifecycle.md`
- `.spectacular/ARCHITECTURE.md`
- `.spectacular/specs/SPC-004-vision-workflow.md`

## Acceptance checks

- A new user can tell when to invoke Imagine and when to skip it.
- The docs show explicit human reaction on parts and approval of the whole.
- No page says Vision requires a request or derives PLAN.
- No page makes UI, a prototype, or a particular fragment kind mandatory.
- No page conflates a spike result with direction approval or pre-spec Vision
  with post-build feedback.
- Command examples match `spectacular vision --help` and the focused CLI tests.
