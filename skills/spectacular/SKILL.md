---
name: spectacular
description: >-
  Guide work only when the user invokes `$spectacular`, `/spectacular`, or in a `.spectacular/` workspace.
  Use for structured mission orchestration, bulk decisions (`spectacular decide`), single-file mission autopilot,
  supervised subagent dispatch, and multi-session campaigns. Triggers on "start mission", "spectacular decide",
  "flight plan", "autopilot", "supervised dispatch", "handoff", "mission check", or "complete mission".
  Do not invoke for generic planning, ungrounded chat, ordinary git operations, or tasks outside Spectacular.
metadata:
  version: "2.9.0"
---

# Spectacular

Run one bounded Mission at a time, from truth the owner already accepted.

## 1. Consolidated CLI Palette & Parameter Grammar

Commands output typed `.v2` JSON envelopes. Mutating commands execute as atomic transactions. `--by` and `--operator` auto-resolve from workspace config/git user if omitted:

```bash
# Core Lifecycle & Decisions
spectacular init [--name <project>]                               # Initialize fresh workspace
spectacular decide <file|-> [--json]                              # Record immutable decision D<N>
spectacular mission start <plan.md|-> [--json]                    # Activate single-file execution envelope
spectacular mission check <ref> [--json]                          # Verify frozen claims & proof (read-only)
spectacular mission complete <ref> [--by <owner>] [--json]        # Complete mission after owner gate

# Delegation, Autopilot & Audits
spectacular charter <mission-ref>[/<obj>] [--json]                # Compile context sandwich (≤1200 tokens)
spectacular handoff record <mission> <draft|-> [--by <actor>]     # Record cross-party handoff
spectacular review record <mission> <draft|-> [--json]            # Record independent review RV<N>
spectacular evidence record <mission> <draft|-> [--json]          # Record third-party proof E<N>
```

## 2. The Lean 3-Layer Autopilot Model

Spectacular is designed for fast, token-efficient autonomous execution with minimal ceremony. All governance reduces to 3 layers:

1. **Layer 1: Living Truth & Decisions**: `PROJECT.md` (boundaries/non-goals) + `.spectacular/decisions/` (bulk-ideated architectural choices recorded with `spectacular decide`).
2. **Layer 2: Topological Flight Plan**: Multi-session roadmap in `.spectacular/campaigns/` (4–8 macro milestone blocks; unstarted blocks remain 4-line lightweight draft cards).
3. **Layer 3: Single-File Execution Envelopes**: Compact, self-contained Mission files (`.spectacular/missions/M<N>-<slug>/M<N>-<slug>.md`, $\le 500$ tokens) with inline objectives, deliverable checklists, and fail-fast stop triggers.

### Mission Layout Judgment Protocol (3-Tier Decision Matrix)
Agents select the minimal sufficient layout tier before writing any mission:

| Tier | When to Select | Structure |
|---|---|---|
| **Tier 1: Single-File (90% Default)** | Routine features, bug fixes, refactors, local deterministic code where test suite passing (`exit 0`) is the proof. | `M<N>.md` **only** (zero sub-folders). |
| **Tier 2: Hybrid Earned (~8%)** | Needs live third-party API receipts (`evidence/`) or splitting work across parallel subagent worktrees (`objectives/`). | `M<N>.md` + only the 1 earned sub-record. |
| **Tier 3: Full Bundle (~2%)** | High-stakes zero-downtime DB cutovers, core auth/crypto, payments, or cross-org handoffs requiring formal checkpoints (`checkpoints/`) or independent adversarial audit (`reviews/`). | Full directory bundle with all needed governance records. |

**The Golden Rule**: Start Single-File. A sub-folder is earned only when a failable condition (external receipt, parallel worktree split, independent security review) explicitly demands it.

### Zero Sub-Record Sprawl Policy
Never create separate `checkpoints/`, `assessments/`, `runs/`, `handoffs/`, or multi-page manual evidence files for routine code tasks. The test suite passing (`exit 0`) and clean Git commit **is** the proof. Context flows across subagents and parallel sessions via lightweight prompts ($\le 300\text{--}500$ tokens) and thread links (`conversation://<id>`).

| Surface | Responsibility |
|---|---|
| **Anchor / Decisions** | Accepted truth: `PROJECT.md` (scope/boundaries) + `decisions/D<N>.md` (durable choices). |
| **Campaign** | Multi-session roadmap sequence in `.spectacular/campaigns/` (planning context; no execution authority). |
| **Proposal** | Optional mutable exploration (`proposals/`); never current truth or execution authority. |
| **Mission / Contract** | Single-file frozen execution envelope (`M<N>`) / modular capability specification (`CC-<module>`). |

### Mechanical Mode (3-State Model)
Invoke `spectacular --version --json` once at startup and require `spectacular.build-info.v1` plus the exact release in `generated/mechanical-interface.json`. If absent, unreadable, or incompatible → reduced mode:
- **CLI Usable**: Standard governed workflow and typed CLI validation.
- **CLI Absent**: Read/draft-only. Route to [reduced-mode.md](references/reduced-mode.md). Never emulate command-owned records, fabricate fingerprints, or claim atomic writes.
- **Declared `manual-bootstrap`**: Owner-approved drafting exception only ([bootstrap.md](references/bootstrap.md)); it cannot create, activate, transition, or complete fingerprint-bound records. Delegated agents cannot declare bootstrap.

## 3. Role Resolution & Orchestration Discipline

| Role | Entry Contract | Context Spine | Exit |
|---|---|---|---|
| **Orchestrator** | Top-level session | PROJECT Anchor → Phase ref → exact records/sources | Safe next action or owner gate |
| **Runner** | Dispatch charter / Handoff | Charter/Handoff → Mission claim → assigned code paths | Bounded result → Orchestrator |
| **Reviewer** | Review assignment | Review assignment → frozen claims → commit/tree → primary proof | Verdict & findings → Orchestrator |
| **Autopilot** | Charter | Charter → Mission target → permitted sources | Chartered return destination |

### Orchestration Taxonomy: Supervised Dispatch vs. Full Handoff
- **Supervised Dispatch (90% Default)**: The Orchestrator retains active Mission ownership, dispatches a worker subagent with a $\le 300$-token charter, and waits reactively for completion (`worker_done`). The worker creates zero governance records; tests passing (`exit 0`) + Git commit is the proof.
- **Full Ownership Handoff (10% Transfer)**: Permanent ownership transfer across distinct sessions, human operators, or different AI harnesses (e.g. Claude $\to$ Codex $\to$ Antigravity). Formally recorded via `spectacular handoff record`.

### The Escalation & Decision Gate Protocol
- When a worker encounters an ambiguous interface, unrecorded architectural choice, or boundary conflict, it **must not guess or improvise**.
- It halts immediately (Fail-Fast Stop) and sends an **escalation** to the Orchestrator.
- The Orchestrator resolves the fork with `spectacular decide` (`.spectacular/decisions/D<N>.md`) and resumes the worker with the locked decision ID.

### Reviewer Role Hygiene (Observe ≠ Act)
- Reviewers inspect diffs, test logs, and primary evidence to evaluate FROST claims and return structured verdicts (`pass`/`fail` + findings).
- **Reviewers NEVER edit files, apply drive-by refactors, or fix observed defects.** Bounded repairs are returned to the Orchestrator to dispatch to a runner.

### Channel Separation: Durable Git State vs. Ephemeral Channels
- **Git is for durable truth**: `PROJECT.md`, `decisions/`, `campaigns/`, and single-file `missions/`.
- **Host channels are for live coordination**: Ephemeral pings, ask/reply loops, and task dispatch stay inside host harness tools (`invoke_subagent`, `send_message`, `conversation://<id>`) with zero file pollution in Git.

- **Role bootstrap**: Top-level owner session defaults to Orchestrator (reads `.spectacular/PROJECT.md` once). Delegated subagents without an entry contract stop and request one; they never self-promote.
- **Anti-escalation**: No entry contract or reference may grant authority above constitutional role ceilings.
- **Context boundary**: Runner reads only assigned inputs. Missing context produces one precise request to Orchestrator; no workspace scans. Campaigns are Orchestrator planning context, never worker selectors.
- **Context Sandwich & Token Discipline**: Worker agents receive a compiled prompt envelope (`spectacular charter`) strictly bounded at $\le 1{,}200$ tokens (`o200k_base`), leaving 99%+ of the model attention window free for codebase AST and test logs. Check token sizes using `bash skills/spectacular/scripts/count-tokens.sh <file|->`.

## 4. Preflight & Isolation

Evaluate branch and worktree isolation independently before mutation:
- **Branch** separates history; branch before activation (`git checkout -b <mission-slug>`).
- **Worktree** separates concurrent hands (`git worktree add`). Concurrent sessions require separate trees.
- Quick-patch directly on `main` is an explicit, non-default owner exception.

### Read-Only Preflight Contract
Check workspace (`PROJECT.md`), Git (branch & worktrees), bindings, identity, and blockers. Report:
1. **Plain outcome**: Current project direction, selected Mission, and lifecycle status.
2. **Technical evidence**: Git branch/worktree, commit SHA, Contract fingerprint, validation mode.
3. **Next action**: Exactly one safe next action, or one owner gate.

## 5. Primary Phase Router

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

## 6. Authority Constitution

- **Owner only**: Outcome, completion criteria, semantic scope, review independence, forbidden-effects.
- **Operator freely**: Reversible attempts, checks, and bounded repairs inside the Mission.
- **Smallest sensible change**: Inside authorized scope, implement only what the frozen claim needs. Do not add abstractions, configuration, refactors, or cleanup unless the claim requires them.
- **Return to owner**: Scope expansion, irreversible/provider effects, exhausted repairs, stops.
- **Proof separation**: Evidence, deterministic checks, independent review, owner acceptance, and completion are separate layers. A passing check proves only its specific observation.

## 7. Owner Maxims

- **Ask only when open**: Semantic forks, boundaries, authority, risks, irreversible effects, contract conflicts. Never re-ask settled decisions.
- **3-Tier Question Escalator**:
  - *Tier 1 (Optimistic Consent)*: State standard/reversible default and proceed non-blocking (`"Proceeding with X unless you prefer Y"`).
  - *Tier 2 (Batch Cards & Four-part formula)*: Lead with the plain outcome and Technical basis; format options as action -> consequence (`1. Question ➔ A, B, C (Recommended default)`). Accept shorthand (`A, B, A`, `all defaults`) and open write-ins.
  - *Tier 3 (Spectrum & Modals)*: Frame competing trade-off axes for open-ended design, or use interactive UI modals (`ask_question`).
- **Authorization, not labor**: Request permission to act; hold the keyboard (see [owner-guidance.md](references/owner-guidance.md)).
- **Report, don't widen (Observe ≠ Act)**: If you notice unrelated problems or defects mid-execution, report them to the owner or Orchestrator. Do not edit them or fold them into the current Mission.
- **Batch gates**: Check prior decisions first; approvals carry forward within the active phase; batch related approvals once.
- **State boundary once**: State constraints once, act on them, and use compact 3-part refusals.

## 8. Continuity & Precedence

- Return the state a cold session needs plus exactly one safe next action or owner gate.
- When Spectacular develops itself, an active Mission's schema and completion boundary remain frozen at activation; later changes apply only to later Missions.
- **Precedence**: Kernel owns invariants and authority. References own conditional procedure. Entry contracts select bounded context. Any conflict is documentation drift: stop and report.
