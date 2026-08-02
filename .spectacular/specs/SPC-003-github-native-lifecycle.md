---
id: SPC-003
type: specification
status: draft
target_version: "tbd"
supersedes: ""
updated: 2026-08-02
summary: "GitHub-friendly collaboration projections for requests, bugs, ideas, pull requests, and verification evidence"
related:
  - ../decisions/DEC-018-keep-spectacular-as-the-primary-lifecycle.md
  - ../decisions/DEC-019-every-executed-spectacular-request-ends-through.md
  - ../decisions/DEC-020-keep-spectacular-questions-as-local-blockers.md
  - ../roadmaps/index.md
---

# SPC-003 — GitHub-friendly collaboration projections for requests, bugs, ideas, pull requests, and verification evidence

## Intent

Make Spectacular friendlier to collaborative projects by integrating with GitHub's native collaboration surfaces without turning GitHub into a competing source of lifecycle truth. Spectacular remains the primary driver; GitHub carries selected public/collaborative projections, review activity, and remote execution evidence.

Implementation is deliberately deferred. This draft preserves the current direction for a later grill session and is not eligible to seed a request.

## Requirements

### Confirmed constraints

- Spectacular remains authoritative for intent, dependencies, approved scope, decisions, and verification interpretation. GitHub provides linked collaboration and evidence surfaces rather than a duplicate canonical database.
- Every executed Spectacular request must end through a pull request before integration. Merge remains human-gated.
- Open `QUE` records remain local blockers by default; publishing internal ambiguity is never automatic.

### Working direction — not yet confirmed

- A bug may be projected to or linked with a GitHub Issue when collaboration benefits from it; not every bug must become an Issue.
- Imported GitHub Issues enter a triage boundary before becoming a BUG, request, idea, question, or no-action record. The durable location and command surface for this triage inbox remain undecided.
- Issue and PR comments may suggest or inform Spectacular work, but comments do not directly mutate lifecycle state, resolve a question, or create a decision.
- GitHub Discussions appear promising for open ideas and a public roadmap, while `.spectacular/ideas/` and `roadmaps/index.md` remain the private/canonical sources.
- Review feedback should eventually route by meaning—for example bug, finding, question, decision candidate, task, or parked idea—but the routing rules require a dedicated interview.

### Explicitly deferred to the grill

- When a request opens its PR: activation/draft, first push, review, or verification.
- Where imported Issue triage lives and how long unclassified imports persist.
- Which GitHub features beyond Issues, Pull Requests, comments, and Discussions earn first-class support: Actions/checks, Projects, Milestones, Releases, CODEOWNERS, environments, or security surfaces.
- Whether synchronization is pull-only, explicitly bidirectional, event-driven, or staged across versions.
- Failure, offline, authentication, permissions, fork, and multi-remote behavior.
- The exact route-by-meaning taxonomy and which routes require user confirmation.

## Evidence and decisions

- `DEC-018` — Spectacular is the primary lifecycle authority; GitHub is a collaboration/evidence projection.
- `DEC-019` — executed requests end through a pull request before integration.
- `DEC-020` — questions remain local blockers by default.
- User working direction recorded on 2026-08-02; all other choices above remain proposals for a future grill.

## Confirmation

Draft and intentionally deferred. Resume with `/spectacular grill spec SPC-003`; review and explicit approval are required before any implementation request may be created.
