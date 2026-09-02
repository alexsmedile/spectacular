# Process

How a piece of work moves through Spectacular.

![The Mission lifecycle](diagrams/lifecycle.svg)

## The Lean 3-Layer Autopilot Model

Spectacular provides an ultra-lean, token-efficient execution lifecycle with minimal ceremony:

1. **Layer 1: Living Ground Truth & Decisions**: `PROJECT.md` (boundaries/non-goals) + `.spectacular/decisions/` (bulk-ideated architectural choices locked with `spectacular decide`).
2. **Layer 2: Topological Flight Plan**: Multi-session roadmap in `.spectacular/campaigns/flight-plan.md` (4–8 macro milestone blocks).
3. **Layer 3: Single-File Execution Envelopes**: Compact, self-contained Mission files (`missions/M<N>.md`, $\le 500$ tokens) with inline deliverables checklist and fail-fast stop triggers.

### Mission Layout & Sub-Folder Selection Protocol
- **Tier 1: Single-File (90% Default)**: `M<N>.md` only. Routine features, bug fixes, refactors where tests passing (`exit 0`) is the proof. Zero sub-folders.
- **Tier 2: Hybrid Earned (~8%)**: `M<N>.md` + 1 earned sub-record (e.g. `evidence/` for live third-party API receipts, or `objectives/` for parallel worktrees).
- **Tier 3: Full Governed Bundle (~2%)**: High-stakes zero-downtime DB cutovers, auth/crypto, or payments requiring formal checkpoints (`checkpoints/`) or independent adversarial audit (`reviews/`).

### Dual-Lane Orchestration: Supervised Dispatch vs. Full Handoff
- **Supervised Dispatch (90% Default)**: In-session subagent delegation. The Orchestrator retains Mission ownership, compiles an immediate prompt via `spectacular charter <mission>/<objective> --prompt` (400–600 token sweet spot containing only objectives, authorized write paths, non-goals, interface definitions, and the test command), and dispatches to a subagent worker under `spectacular guard` without extra governance files.
- **Full Ownership Handoff (10% Transfer)**: Explicit ownership transfers across distinct sessions, human engineers, or different AI harnesses recorded via `spectacular handoff record`.

### The 3-Tier Question Escalator
- **Tier 1: Optimistic Consent**: Non-blocking 1-line statement of defaults for low-risk implementation choices (`"Proceeding with X unless you prefer Y"`).
- **Tier 2: Structured Batch Cards**: Numbered questions with lettered options (`1. Question ➔ A, B, C (Recommended)`), adaptive context depth, batch shorthand replies (`A, B, A`, `all defaults`), and open write-in support.
- **Tier 3: Trade-off Spectrum & Interactive Modals**: Framing competing design axes for unpredictable exploration, or leveraging interactive UI modals (`ask_question`).

> [!NOTE]
> **Standalone Interviewing**: The Question Escalator is not restricted to Mission activation. You can use Interview Mode during early ideation or architectural forks to align on choices, immediately lock the rulings via `spectacular decide`, and write code directly without drafting a Mission or Contract.

## Intake and probes

Ideas and small experiments do not need a Mission. Keep them in an issue,
`TODO.md`, or `scratch/`. A quick test is fine when it is time-boxed,
reversible, and does not touch external services, production data, or a release.

Promote the result when it becomes durable work: use a Proposal for a
consequential unresolved choice, or a compact Mission when the work has a
bounded outcome and needs authority, evidence, or review. Do not run an
unbounded intake queue on Autopilot.

## Campaigns: optional roadmap maps

A Campaign is one Markdown file in `.spectacular/campaigns/` for a genuinely
independent strategic arc. Aim for 4–10 roadmap blocks with their dependencies,
candidate or active Missions, and an exit condition; split an unwieldy map into
separate arcs. It is a planning map, not an automation queue or execution
authority. Campaigns are optional: use one only when several Missions need a
shared roadmap. Prefer a Mermaid dependency view when it makes order clearer.
`spectacular campaign check <path>` is read-only: it validates the optional
frontmatter map, resolves named Mission refs, detects cycles, and emits the
ordered Mermaid projection.

Campaign checking is an orchestration/planning action. Its `current` field is a
global map position, not a worker assignment; Mission workers normally follow
their Mission, Objective, and Run without loading a Campaign.

## Explore: the Proposal and Starter Inputs

When starting a project, use a short product brief or starter document. The
setup step turns it into the project files that describe the product, technical
stack, architecture, and first Mission.

### How Spectacular Disambiguates "Specs"

Because "spec" is an overloaded term, Spectacular routes specification work based on the maturity and intent of the request:

```mermaid
flowchart TD
    A["Request: 'Let's write/update specs for X'"] --> B{What stage is the spec in?}
    B -->|"1. Brainstorming / Exploring / Unresolved"| C["<b>Proposal</b><br><code>.spectacular/proposals/P&lt;N&gt;-&lt;slug&gt;.md</code><br><i>Mutable draft, open questions, alternatives</i>"]
    B -->|"2. Accepted Subsystem Behavior"| D["<b>Capability Contract</b><br><code>.spectacular/contracts/CC-&lt;name&gt;.md</code><br><i>Observable invariants, inputs/outputs, gaps</i>"]
    B -->|"3. High-Level System Scope"| E["<b>Project Anchors</b><br><code>.spectacular/PROJECT.md</code> & <code>ARCHITECTURE.md</code><br><i>System boundaries, directory layers, tech stack</i>"]
    B -->|"4. Ready to Execute & Prove"| F["<b>Mission</b><br><code>.spectacular/missions/M&lt;N&gt;/</code><br><i>Frozen verifiable claims & proof requirements</i>"]
```

### Iterating on Specs During Brainstorming

When you and an agent brainstorm in chat:
1. **Iterate Freely in a Proposal**: Draft the candidate specification in `.spectacular/proposals/P<N>-<slug>.md`. Proposals are completely mutable and carry no execution authority—you can rewrite, expand, or pivot the spec across dozens of chat turns without causing contract drift or validation failures.
2. **Lock Key Trade-offs with Decisions**: When an architectural or policy fork is settled, record it via `spectacular decide` (`D<N>`) to capture the immutable rationale.
3. **Promote to Living Truth**: When the specification is finalized, incorporate the observable invariants into a Capability Contract (`.spectacular/contracts/CC-<name>.md`) or Project Anchor (`PROJECT.md` / `ARCHITECTURE.md`).
4. **Retire the Proposal**: Mark the Proposal `status: accepted` with `resolved_by:` and move it to `.spectacular/archive/proposals/`.
5. **Execute via Mission**: Launch a frozen Mission (`spectacular mission start`) with verifiable claims and proof requirements.

A Proposal is optional. Write one when the approach is genuinely unclear and you
want to argue with it before committing.

It is mutable, carries no authority, and its status is an owner assertion rather
than a derived fact. Nothing about a Proposal grants permission to act — that is
the point. It is the cheapest place to be wrong.

When a Proposal's work ships, it is **retired**: it names the Mission that
absorbed it in `resolved_by:` and moves to `.spectacular/archive/proposals/`. A
Proposal is absorbed when the question it asked was answered, not when most of it
was. See `D11-proposal-retirement`.

## Prepare and freeze: the Mission

A Mission plan states an outcome, a stop condition, and the claims that must
hold. Activation freezes it and fingerprints the text.

### Upfront Alignment: Pattern Pass & Domain Vocabulary

Before code is written or delegated to subagents:
1. **The Architectural Pattern Pass (D29)**: The Orchestrator surveys standard library idioms, RFC specifications, and battle-tested open-source reference implementations (e.g. `xterm.js`, standard AST parsers) rather than generating bespoke abstractions from scratch. A 3-line Pattern Census is frozen into the Mission body.
2. **Domain Ontology & Banned Synonyms (`VOCABULARY.md`)**: Ubiquitous language is maintained under single-writer authority (Owner/Orchestrator only). Explicit Banned Synonyms prevent LLMs across fresh context windows from drifting between ambiguous verbs and states.

**What freezing means:** the Mission is judged against what it said at
activation. It is not edited later to match what actually happened. If the
agreement turns out to be wrong, that is a real event worth recording — amend it
through `contract amend`, which rewrites the Gap and re-points the live Mission
as one recoverable transaction.

**Branch before activating.** A Mission that runs on `main` has destroyed the
review and isolation boundary it depends on before it starts.

## Execute: Runs, Objectives, and Model Profiles

A Run is one bounded attempt at the frozen Mission.

Missions and handoffs can name the kind of model work they need: careful
planning, fast implementation, or strict review. The host chooses the matching
model when it can.

Objectives are **earned, not planned**. Expand one when the work is real, rather
than enumerating a full tree upfront that will be wrong by the second Objective.
Dependencies come in two kinds:

- `after:` — this Objective needs the produced artifact, and is genuinely sequential.
- `after_interface:` — it needs only the contract shape, which exists the moment
  the plan freezes, so it can start at activation.

That split is what lets proof Objectives begin immediately instead of queueing
behind implementation.

`mission check` reports per-claim drift at any point: which claims are repaired,
how stale the evidence is, and where the next flag is.

### Checkpoints are planned Run-body gates

A checkpoint is an optional, planned point in a Run for reviewing progress,
running a check, collecting a decision, or choosing whether to resume. It is
not automatically a human-review gate and does not carry authority by itself.
Record the planned checkpoints in the Run body at activation, then update the
same section as the Run progresses.

Use the durable record that matches what happened at a checkpoint:

| If the checkpoint produces… | Record it as… |
| --- | --- |
| an owner choice or changed direction | a Decision |
| an observation or test result | Evidence |
| a verdict on the work | Review or Assessment |
| a transfer to another operator or runtime | Handoff |
| none of the above | a concise Run-body checkpoint note |

The Run-body note names the trigger, what was reviewed, the result, and the
next action. It supports a cold resume without creating a new record for every
ordinary progress check.

### Supervised Worker Dispatch & Prompt Budgeting (`spectacular charter --prompt`)

When delegating an Objective to a subagent worker, the Orchestrator avoids passing the entire workspace context. Instead, it extracts a lean **Worker Charter Prompt**:

```sh
spectacular charter M17/O1 --prompt
```

The prompt is structured into three self-contained sections designed for **Zero File Wandering**:
1. **Initial Code Grounding & Target Files**: Explicit authorized write paths and interface definitions, removing the need for blind `ls` or `find` exploratory turns.
2. **Test Expectation & Verification**: The exact test command (e.g. `sh tests/check.sh`) and pass criteria required to turn the build green.
3. **Worker Contract & Protocol**:
   - Strictly implement within authorized paths without touching `.spectacular/` governance files.
   - Iterate against the test runner until `exit 0`.
   - Single-Return Signals: `STATUS: DONE` (with changes summary) or `STATUS: BLOCKED` (with decision required).

#### Flexible Prompt Token Budgeting

| Tier | Task Complexity | Target Budget | Typical Payload |
| :--- | :--- | :--- | :--- |
| **Tier 1 (Micro)** | Targeted bugfix / single file | **250–400 tokens** | Target file, verification command, return signal |
| **Tier 2 (Standard)** | Feature / multi-function change | **400–600 tokens** *(Sweet Spot)* | Target files, interface signatures, test expectation, locked decisions |
| **Tier 3 (Genesis)** | Multi-file module / event journal | **600–900 tokens** *(Ceiling)* | Full mechanical perimeter, data models, crash replay assertion |

> **Splitting Rule**: If an Objective's charter prompt exceeds **900 tokens**, it is an architectural indicator that the Objective should be decomposed into smaller sequential sub-objectives (`O1`, `O2`).

### Atomic Verification Gate (`spectacular mission check --verify`)

During execution and before closing a Mission, `spectacular mission check <ref> --verify` executes a unified, atomic 4-point verification check in a single pass:
1. **Structural & Contract Drift**: Verifies Mission schema integrity and frozen contract bindings.
2. **Domain Test Suite**: Runs `tests/check.sh` (or `tier1_quick`).
3. **Replay Hook Execution**: If declared, evicts cache paths and verifies full state reconstruction.
4. **Git Cleanliness**: Asserts working tree cleanliness (`git status --porcelain`).

It returns a deterministic JSON receipt:
```json
{
  "ref": "M17",
  "valid": true,
  "checks": [
    "domain-verification-pass",
    "replay-reconstruction-pass",
    "git-working-tree-clean"
  ]
}
```

### Visual Campaign Inspection (`spectacular campaign check --ascii`)

To inspect a multi-mission Campaign dependency graph in the terminal without opening heavy visualization tools:

```sh
spectacular campaign check .spectacular/campaigns/genesis.md --ascii
```

Renders an instant terminal DAG with live state indicators:
```text
Campaign DAG: Multi-Agent Platform Genesis
Focus: Autonomous subagent orchestration and verification gates.

[✓] B1: Architecture Core [complete]
    └── missions: M1, M2
[▶] B2: Replay Verification & Charter Dispatch [active] (after B1)
    └── missions: M10, M17
[ ] B3: Active Perimeter Jail [planned] (after B2)
    └── missions: M18
```

## Frontier Orchestration: Taming High Generation Velocity

Modern frontier models generate vast amounts of code and unit tests rapidly. Without structural containment, this high **Generation Velocity** leads to two major failure modes: **Blast Radius explosion** (touching or refactoring unrelated files) and **Context Amnesia** (losing architectural nuances and previous review feedback across fresh context windows).

Spectacular tames frontier models by pairing visual navigation with mechanical containment:

```mermaid
flowchart TD
    subgraph HUD ["1. Architectural HUD & Maps"]
        M["🗺️ Atlas Domain Maps (.spectacular/atlas/)<br><i>Non-governing Mermaid views of contexts & transitions</i>"]
        C["📊 Campaign DAG Projections (campaign check)<br><i>Milestone sequence & dependency flow</i>"]
    end

    subgraph Dispatch ["2. Supervised Dispatch"]
        Ch["📦 Charter Sandwich (spectacular charter --prompt)<br><i>Bounded context (≤400-600t): objectives & write paths</i>"]
        G["🛡️ Watchdog Jail (spectacular guard)<br><i>OS-level quarantine against blast radius sprawl</i>"]
    end

    subgraph Review ["3. Independent Review Sessions (Observe ≠ Act)"]
        F["🔍 4-Point FROST Review<br><i>Decision Compliance · Atlas Coverage · Blast Radius · Proof Validity</i>"]
        R["⚖️ Decision Rulings (spectacular decide)<br><i>Lock newly discovered architectural forks</i>"]
    end

    HUD --> Dispatch
    Dispatch --> Review
    Review -.->|Iterate & Repair| Dispatch
    Review --> R
```

### 1. The Architectural HUD & Maps
Instead of reading thousands of lines of raw diffs across complex repos, the Orchestrator maintains an **Architectural HUD**:
- **Visual Maps (`.spectacular/atlas/`)**: Non-governing visual Mermaid diagrams of bounded contexts, entity relationships, and state lifecycles. Governed by Decision D26, they navigate and explain without asserting schema claims or causing mission drift.
- **Topological Projections (`spectacular campaign check`)**: Terminal ASCII DAGs or Mermaid graphs showing milestone readiness and dependencies.

### 2. Controlling Dispatch & Containing Blast Radius
When the Orchestrator delegates implementation to a subagent:
- **The Context Sandwich (`spectacular charter M<N>/<obj> --prompt`)**: Workers receive only the lean charter prompt ($\le 400\text{--}600$ tokens: target objective, interface definitions, non-goals, and the test command). Workers are explicitly instructed to ignore `.spectacular/` and never create governance files.
- **Filesystem Watchdog (`spectacular guard`)**: The charter names explicit authorized write paths (`allowed_changed_paths`). `spectacular guard` acts as an OS-level watchdog, quarantining any rogue file edits outside authorized boundaries.
- **Fail-Fast Escalation Gate**: If the worker hits an unrecorded fork, it stops immediately (`STATUS: BLOCKED`) rather than hallucinating an ad-hoc pattern. The Orchestrator resolves the fork via `spectacular decide` (`D<N>.md`) and resumes the worker.

### 3. Review Sessions: Enforcing "Observe ≠ Act"
Frontier models suffer confirmation bias when reviewing their own fresh code. Spectacular enforces a strict separation between generation and review:
- **Observe ≠ Act**: Reviewers evaluate and report findings; they do not edit code.
- **The 4-Point Review Rubric**:
  1. **Decision Compliance**: Did the model strictly obey all locked `D<N>` rulings and stack constraints?
  2. **Atlas Coverage**: Did the code implement all declared state transitions in the visual Maps?
  3. **Blast Radius**: Did the changes remain strictly within authorized directories?
  4. **Proof Validity**: Are automated test receipts authentic, reproducible, and exiting `0`?

### 4. Where Reviews Live & The Traceability Dial
Spectacular calibrates ceremony based on risk:

- **Routine Fast Code (90% Default)**:
  For standard features, bug fixes, and clean refactors, **traceability is minimal**. Automated tests passing (`exit 0`) and clean Git commits are the primary proof. Reviews do not generate separate files—feedback is passed directly in session chat, PR descriptions, or commit notes.
- **High-Stakes Code (~10% Critical)**:
  For zero-downtime DB cutovers, cryptography, auth, or public API breaks, reviews are formally recorded using `spectacular review record M<N> review_draft.md --json`. This saves an immutable record at `.spectacular/missions/M<N>/reviews/RV<N>.md` capturing `reviewed.commit`, `activation_fingerprint`, and specific verdicts.
- **Are `reviews/` Standalone?**:
  In Spectacular's mechanical schema, formal `type: Review` (`RV<N>`) records are **mission-scoped** (`.spectacular/missions/M<N>/reviews/`) because they mechanically validate a Mission's frozen claims and activation fingerprint. In **A La Carte** mode (without a Mission), review findings are kept with zero ceremony in chat or PR notes, while any discovered architectural decisions are locked into `.spectacular/decisions/D<N>.md`, which is top-level and standalone.

## Prove: Evidence, Reviews, Handoffs

Keep proof in records, not chat messages.

A **Review** carries a verdict. An **Assessment** carries a judgment. **Evidence**
carries what was observed. None of them is a claim that something works — they
are the artifact a later reader uses to decide whether to trust it.

Independent review supports a **dual-path workflow**:
1. **In-Harness Subagents**: A clean-slate child agent is dispatched with the commit SHA and FROST criteria, writing `ReviewDraft` directly to `reviews/`.
2. **External Model / Human Handoff**: A copy-pasteable review prompt is generated in `handoffs/review-handoff-prompt.md` to run against external models (e.g. OpenAI o3, DeepSeek-R1) or peer reviewers, returning the verdict into `spectacular review record`.

A **Handoff** binds the exact commit and tree it was sent against, verified
against the repository, and splits what the sender knows into two lists:

- `asserted` — what the sender verified themselves
- `assumed` — what they are carrying on trust

Neither is scored. The split exists so the receiver knows exactly what to
re-verify before acting. A recorded Handoff is frozen; correcting one means
recording a new Handoff that carries `supersedes:`, and the original survives as
what its sender believed at the time.

## Close: completion is an owner act

```sh
spectacular mission complete M12 --by alex
```

An agent cannot complete its own Mission. Completion refuses while a declared Gap
is still open, so "we shipped it but this is still broken" cannot be recorded as
success.

### The "Mistake Tax" & Guardrail Feedback Loop (D29)
When a Mission completes having consumed repair budgets, hit unexpected regressions, or resolved non-trivial review findings, the failure root cause is codified into permanent rules:
- **Domain & System Invariants** $\to$ Appended directly to `.spectacular/GUARDRAILS.md`.
- **Harness & Tooling Invariants** $\to$ Appended directly to `AGENTS.md`.
- **Architectural Trade-offs** $\to$ Recorded as atomic Decisions (`.spectacular/decisions/`).

Ad-hoc `LEARNINGS.md` or `FIXES.md` sprawl is forbidden, ensuring every lesson is routed to its single permanent home.

A **Gap** is a stated limit rather than a defect. It is never closed by deleting
it — the entry survives with a written resolution, so the reason something was
impossible stays recoverable years later.

## After: freeze points and the archive

A completed Mission keeps its original contract fingerprint forever. Amendments
re-point only the live Mission. When `mission check` reports contract drift on a
completed Mission, that is a **notice**, not an error: the Mission stays valid,
and `git log -S <fingerprint>` recovers the Contract text as it was.

Completed Missions and absorbed Proposals move to `.spectacular/archive/`,
carrying the Decision that authorized the move and a fingerprint of what was
archived. Archived records are kept machine-readable — mainly for their numbering
and their reasoning — not for routine reading.

## The rule underneath all of it

The system does not decide anything. It validates, fingerprints, writes
atomically, and reports. Every judgment — is this good, is this done, may this
proceed — is an owner act at a named gate.

That is what makes an agent's work resumable: not a longer chat history, but a
record of what was agreed, what was attempted, what was proven, and who said yes.

## See also

- [Quickstart](quickstart.md) — run one Mission end to end.
- [Architecture](architecture.md) — the surfaces and record types.
