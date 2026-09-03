---
name: spectacular
description: >-
  Guide work only when the user invokes `$spectacular`, `/spectacular`, or in a `.spectacular/` workspace.
  Use for structured mission orchestration, bulk decisions (`spectacular decide`), single-file mission autopilot,
  supervised subagent dispatch, and multi-session campaigns. Triggers on "start mission", "spectacular decide",
  "flight plan", "autopilot", "supervised dispatch", "handoff", "mission check", or "complete mission".
  Do not invoke for generic planning, ungrounded chat, ordinary git operations, or simple status/branch checks.
metadata:
  version: "2.17.0"
---

# Spectacular

Run one bounded Mission at a time, from truth the owner already accepted.

> **Fast Bailout**: If the query is a simple inspection (e.g. `git branch`, status check, diff, or questions without `$spectacular`), answer directly using native tools, report status: "done", and exit immediately with zero ceremony. Do not read `.spectacular/PROJECT.md` or load references.

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

## 2. Fast Autonomous Model & Foundational Anchors

Spectacular prioritizes execution over ceremony. Governance is managed strictly by the top-level **Orchestrator**; dispatched **Workers/Subagents do NOT manage Spectacular files** and execute purely against their code charter:

```
┌─────────────────────────────────────────────────────────────┐
│ 1. Living Truth (Orchestrator): PROJECT.md & decisions/     │
│ 2. Single-File Envelope: M<N>.md (frozen claims & checks)   │
│ 3. Direct Execution (Workers): Zero-preamble code & tests   │
└─────────────────────────────────────────────────────────────┘
```

### The 5 Foundational Anchors
1. **Boundaries & Non-Goals**: `PROJECT.md` (`boundaries:`, `constraints:`).
2. **Vocabulary & Ontology**: `VOCABULARY.md` (Canonical domain terms).
3. **Invariants & Failure Modes**: `GUARDRAILS.md` & `AGENTS.md` (Non-negotiable safety rules).
4. **Data Structures & Schemas**: Project-specific types/schemas in the codebase (cited in `contracts/`).
5. **State Machines & Lifecycles**: Non-governing visual Mermaid diagrams in `.spectacular/atlas/`.

### Graduated Governance Ladder (`governance:`)
Answers: *"What governs this slice of work?"*
- **`governance: inline` (Tier 0)**: Direct pair-programming in primary chat (`lead-checkout`). Zero files created.
- **`governance: board` (Tier 1)**: Gated dependency pipeline on non-governing `type: WorkBoard`. Lead tracks order and gates.
- **`governance: brief` (Tier 2)**: Temporary teammate in isolated `linked-worktree` with a plain-English Dispatch Brief.
- **`governance: mission` (Tier 3)**: Full immutable `M<N>.md` contract, Handoff, and compiled Charter for high-stakes milestones.

### Gated Waves & Operating Dial
- **Gated Waves (Sequential by Default)**: Stay sequential in the Lead session unless tasks have separate inputs, disjoint write scopes, and locked upstream interface contracts. Parallel side sessions are earned only after an interface gate passes.
- **"Returned ≠ Done"**: Side workers return code, diffs, and test receipts; only the Lead Orchestrator integrates branches and runs project-wide verification.
- **`mode: leverage` (Default)**: High autonomy; test suite passing (`exit 0`) + clean diff is primary proof.
- **`mode: control`**: High-precision mode for irreversible cutovers (auth, payments, DB migrations) requiring formal reviews and Evidence.
- **Minimal Drafts (Zero YAML Boilerplate)**: Prefer direct CLI flags (`spectacular decide --title ...`) or 3-line plans; the CLI auto-populates metadata.
- **Silent Mutation, One-Line Return**: State the 1-line outcome only (`"Recorded Decision D30. Next: run preflight"`). Never echo YAML frontmatter or full file bodies into chat.
- **Drop Collection Catalogs**: Never load catalogs or indexes into agent context. Query the CLI (`spectacular mission show <ref> --json`) directly.

### GitHub Native Integration Layer (`gh`)
Spectacular leverages native GitHub collaboration features via the `gh` CLI while keeping Git as the durable authority:
- **Intake (`gh issue view <id>`)**: When prompt or plan cites an issue (`#<id>`), fetch issue context to frame Mission outcome and acceptance claims.
- **PR Envelope (`gh pr create`)**: For consequential missions, branch `m<N>-<slug>` and link the Mission file in the PR description (`Closes #<id>`).
- **Review Mirroring (`gh pr review`)**: When an audit review is recorded (`.spectacular/reviews/RV<N>.md`), optionally submit the review body to GitHub's PR timeline.
- **Local / Airgapped Fallback**: Never require network access; operate 100% locally if `gh` is unauthenticated or the repo is offline.

### Mechanical Mode (3-State Model)
Invoke `spectacular --version --json` once at startup and require `spectacular.build-info.v1` plus the exact release in `generated/mechanical-interface.json`:
- **CLI Usable**: Standard governed workflow and typed CLI validation.
- **CLI Absent**: Read/draft-only. Route to [reduced-mode.md](references/reduced-mode.md). Never emulate command-owned records or fabricate fingerprints.
- **Declared `manual-bootstrap`**: Owner-approved drafting exception only ([bootstrap.md](references/bootstrap.md)).

## 3. Role & Delegation Matrix

| Role | Responsibility | Context Spine | Output |
|---|---|---|---|
| **Lead (Orchestrator)** | Owns workspace truth, decisions (`decide`), & activation | `PROJECT.md` → Phase Ref | Next action or Owner Gate |
| **Worker / Side Session** | Executes code & tests in `linked-worktree`; **ignores governance** | Dispatch Brief ($\le 1200$ tok) | Return Receipt + Git diff |
| **Reviewer** | Inspects code against frozen claims (**Observe ≠ Act**) | Frozen claims → Diff | Structured verdict (`pass`/`fail`) |

- **Physical Workspaces**: `lead-checkout` (Lead/sequential), `linked-worktree` (`.worktrees/<slug>`), `sandbox` (disposable container/spike), `read-only` (reviewer).
- **Escalation Gate**: When a worker hits an architectural fork, it stops immediately. The Orchestrator records the choice via `spectacular decide` (`D<N>.md`) and resumes the worker.
- **Workers Never Edit Governance**: Subagents never create `checkpoints/`, `runs/`, or `missions/`.
- **Channel Separation**: Git is for durable truth (`PROJECT.md`, `decisions/`, `missions/`). Host channels are for ephemeral live coordination (`invoke_subagent`, `send_message`, `conversation://<id>`).
- **Token Discipline**: Worker prompt envelopes strictly bounded at $\le 1{,}200$ tokens (`o200k_base`).

## 4. Preflight & Verification Matrix

- **Branch Isolation**: Always `git checkout -b <slug>` before mission activation, or dispatch side workers to `.worktrees/<slug>`.
- **Verification Tiers**:
  - *Tier 1 (Quick)*: Executed by worker on every edit (`verify.sh quick` or domain test).
  - *Tier 0 (Preflight)*: Lint & syntax verification (`verify.sh preflight`).
  - *Tier 2/3 (Acceptance/Release)*: Executed at milestone completion / owner gate.


## 5. Primary Phase Router (Load $\le 1$ Reference)

| Phase | Trigger Context | Primary Reference |
|---|---|---|
| `orient` | Cold-start or ambiguous workspace | [orient.md](references/orient.md) |
| `prepare` | Greenfield ideation, Proposal, or Mission drafting | [prepare.md](references/prepare.md) |
| `execute` | Active Mission execution & concurrency invariants | [execute.md](references/execute.md) |
| `runtime` | Packaging subagent charters & handoffs | [runtime.md](references/runtime.md) |
| `close` | Completion claim check & Evidence | [close.md](references/close.md) |
| `audit` | Independent FROST claim challenge | [audit.md](references/audit.md) |

Load a supporting reference only when the primary reference explicitly triggers it. When the phase changes, finish or stop the current phase before routing again.

## 6. Authority & Execution Invariants
- **Authority**: Owner owns outcomes, boundaries, and acceptance. Operator freely attempts reversible checks and bounded repairs. `A Decision is not activation authority` (only owner confirms `mission start`).
- **Direct Greenfield Execution**: Skip meta-planning chat on direct builds. Write code and tests, run `tests/check.sh` / `verify.sh quick`, and report the terminal result.
- **Concurrency & Queues**: Bind concurrency to `--workers N`. Track retry attempts per item; route to `dlq.json` only after exceeding failure threshold ($\ge 3$).
- **Proof Separation**: Test passing (`exit 0`) proves deterministic mechanics. Independent reviews (`reviews/`) evaluate `Frozen fit` and `Truth of proof` without modifying code (Observe ≠ Act).

## 7. Owner Interaction & Continuity
- **Questions**: Ask only when open. Lead with the plain outcome and Technical basis; format options as action -> consequence (`1. Option A, B (Recommended default)`).
- **Self-Hosting**: When developing Spectacular, an active Mission keeps the schema frozen. Under declared `manual-bootstrap`, run focused checks directly.
- **Continuity**: Return cold-session state plus exactly one safe next action or owner gate. Kernel owns invariants; references own conditional procedures.
