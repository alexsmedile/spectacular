---
name: spectacular
version: 2.6.0-rc1
description: Guide work in a canonical Spectacular v2 workspace through optional Proposal exploration, compact Mission preparation and activation, governed execution, earned Objective/Run expansion, Evidence and review, owner completion, audit, and cold recovery. Use for `/spectacular` jobs such as orient, explore, propose, plan, start, resume, handoff, review, complete, or audit; for compiling bounded runtime context or an Autopilot charter; and for safely continuing after session or runtime replacement.
---

# Spectacular

Run one bounded Mission at a time, from truth the owner already accepted.

## 1. Core Model & Mechanical Invariants

Markdown is canonical. Folders and explicit pointers provide navigable context. Agents recover from canonical records, not chat history. CLI owns deterministic validation, bindings, SHA-256 fingerprints, UUIDv7 identity, and atomic transitions. Missions are frozen execution envelopes; Objectives own outcome-sized claims; Runs are mutable attempts. Human docs (`docs/`) are not runtime governance authority.

| Surface | Responsibility |
|---|---|
| **Anchor** | Accepted truth: `PROJECT.md` (scope/boundaries), `STACK.md` (tools), `ARCHITECTURE.md` (layers). |
| **Campaign** | Optional roadmap sequence in `.spectacular/campaigns/` (planning context; no execution authority). |
| **Proposal** | Optional mutable exploration (`proposals/`); never current truth or execution authority. |
| **Mission / Contract** | Frozen execution envelope (`M<N>`) / modular capability specification (`CC-<module>`). |

### Mechanical Mode (3-State Model)
Invoke `spectacular --version` once at startup. If absent, unreadable, or incompatible → reduced mode:
- **CLI Usable**: Standard governed workflow and typed CLI validation.
- **CLI Absent**: Read/draft-only. Route to [reduced-mode.md](references/reduced-mode.md). Never emulate command-owned records, fabricate fingerprints, or claim atomic writes.
- **Declared `manual-bootstrap`**: Owner-approved exception only ([bootstrap.md](references/bootstrap.md)). Follow manual checklist; never cite CLI validation or atomicity. Delegated agents cannot declare bootstrap.

## 2. Role Resolution & Context Discipline

| Role | Entry Contract | Context Spine | Exit |
|---|---|---|---|
| **Orchestrator** | Top-level session | PROJECT Anchor → Phase ref → exact records/sources | Safe next action or owner gate |
| **Runner** | Handoff contract | Handoff → Objective → Run → named working inputs | Bounded result → Orchestrator |
| **Reviewer** | Review assignment | Review assignment → frozen claims → commit/tree → Evidence | Verdict & findings → Orchestrator |
| **Autopilot** | Charter | Charter → Objective/Run → permitted sources | Chartered return destination |

- **Role bootstrap**: Top-level owner session defaults to Orchestrator (reads `.spectacular/PROJECT.md` once). Delegated subagents without an entry contract stop and request one; they never self-promote.
- **Anti-escalation**: No entry contract or reference may grant authority above constitutional role ceilings.
- **Context boundary**: Runner reads only assigned inputs. Missing context produces one precise request to Orchestrator; no workspace scans. Campaigns are Orchestrator planning context, never worker selectors.

## 3. Preflight & Isolation

Evaluate branch and worktree isolation independently before mutation:
- **Branch** separates history; branch before activation (`git checkout -b <mission-slug>`).
- **Worktree** separates concurrent hands (`git worktree add`). Concurrent sessions require separate trees.
- Quick-patch directly on `main` is an explicit, non-default owner exception.

### Read-Only Preflight Contract
Check workspace (`PROJECT.md`), Git (branch & worktrees), bindings, identity, and blockers. Report:
1. **Plain outcome**: Current project direction, selected Mission, and lifecycle status.
2. **Technical evidence**: Git branch/worktree, commit SHA, Contract fingerprint, validation mode.
3. **Next action**: Exactly one safe next action, or one owner gate.

## 4. Primary Phase Router

Orchestrators and primary operators load exactly one primary phase reference:

| Phase | When Session Needs It | Primary Reference |
|---|---|---|
| `orient` | Ambiguous, cold-start, or uninitialized workspace state | [orient.md](references/orient.md) |
| `prepare` | Greenfield Genesis, exploration, Proposal, or Mission drafting | [prepare.md](references/prepare.md) |
| `execute` | Mission activation, execution mechanics, Git isolation | [execute.md](references/execute.md) |
| `runtime` | Packaging delegation, Handoff contracts, Autopilot charters | [runtime.md](references/runtime.md) |
| `close` | Claim assessment, Evidence, review, owner completion | [close.md](references/close.md) |
| `audit` | Retrospective proof or claim challenge using FROST | [audit.md](references/audit.md) |

Load a supporting reference only when the primary reference explicitly triggers it. When the phase changes, finish or stop the current phase before routing again.

## 5. Authority Constitution

- **Owner only**: Outcome, completion criteria, semantic scope, review independence, forbidden-effects.
- **Operator freely**: Reversible attempts, checks, and bounded repairs inside the Mission.
- **Return to owner**: Scope expansion, irreversible/provider effects, exhausted repairs, stops.
- **Proof separation**: Evidence, deterministic checks, independent review, owner acceptance, and completion are separate layers. A passing check proves only its specific observation.

## 6. Owner Maxims

- **Ask only when open**: Semantic forks, boundaries, authority, risks, irreversible effects, contract conflicts.
- **Four-part question formula**: (1) Plain outcome · (2) Technical basis · (3) Options (`action -> consequence`) · (4) Recommended default & why.
- **Authorization, not labor**: Request permission to act; hold the keyboard (see [owner-guidance.md](references/owner-guidance.md)).
- **Batch gates**: Check prior decisions first; approvals carry forward within the active phase; batch related approvals once.
- **State boundary once**: State constraints once, act on them, and use compact 3-part refusals.

## 7. Continuity & Precedence

- Return the state a cold session needs plus exactly one safe next action or owner gate.
- When Spectacular develops itself, an active Mission's schema and completion boundary remain frozen at activation; later changes apply only to later Missions.
- **Precedence**: Kernel owns invariants and authority. References own conditional procedure. Entry contracts select bounded context. Any conflict is documentation drift: stop and report.
