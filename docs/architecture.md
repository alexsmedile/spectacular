# Architecture

Spectacular keeps five kinds of information in separate places. That makes it
clear what people read, what agents follow, and what the CLI can change.

| Surface | Audience | What it is |
|---|---|---|
| `.spectacular/` | governance | The records: Missions, Proposals, Decisions, Evidence |
| `.spectacular/raw/` | thinking | Unstructured sketchpad; gitignored, skip-listed, no frontmatter and no entity |
| `.spectacular/atlas/` | planning | Optional outcome and system maps; explanatory, mutable, and non-authoritative |
| `.spectacular/campaigns/` | planning | Optional durable roadmap maps; excluded from the record graph and CLI lifecycle |
| `skills/` | agents at runtime | Executable guidance the CLI and Skill load |
| `cmd/` + `internal/` | the machine | A typed CLI that validates and mutates records atomically |
| `docs/` | humans | What you are reading |

## Why files, not a database

Every record is a Markdown file stored in Git. There is no separate service or
hidden database. The small YAML section holds details the CLI checks; the prose
explains the decision to people and agents.

Because the records are files, you can review a change in a pull request and
recover the exact text that was approved at the time.

## The record types

Thirteen types, in two groups.

**Top-level** — they exist independently:

- **Proposal (`P<N>`)** — *"What could we do?"* Optional, mutable exploration. Carries no authority; the cheapest place to be wrong during brainstorming.
- **Decision (`D<N>`)** — *"Which path did we choose?"* Durable, immutable record of an owner choice and its rationale. Recorded atomically via `spectacular decide`.
- **Contract (`CC-<name>`)** — *"What does the system guarantee right now?"* Living, versioned specification of observable capabilities and invariants.
- **Mission (`M<N>`)** — *"What are we building and proving right now?"* Bounded, frozen agreement with execution authority, exact git baseline, and failable proof.

### The 4 Governance Primitives

| Primitive | Core Question | Lifecycle State | Authority Level | Location |
|---|---|---|---|---|
| **Proposal (`P<N>`)** | *"What could we do?"* | **Mutable & Exploratory** | **Zero Authority** (ideas, draft specs, open questions) | `.spectacular/proposals/` |
| **Decision (`D<N>`)** | *"Which path did we choose?"* | **Immutable & Permanent** | **Attributable Ruling** (owner choice, rationale, trade-offs) | `.spectacular/decisions/` |
| **Contract (`CC-<name>`)** | *"What does the system guarantee right now?"* | **Living Truth** (versioned) | **Governing Specification** (system invariants, behaviors, gaps) | `.spectacular/contracts/` |
| **Mission (`M<N>`)** | *"What are we building and proving right now?"* | **Frozen Execution Envelope** | **Execution Authority** (atomic claims, proof, reviews) | `.spectacular/missions/` |

```mermaid
flowchart LR
    P["💡 Proposal (P1)<br><i>Brainstorm & draft specs</i>"] --> D["⚖️ Decision (D1)<br><i>Owner locks key choices</i>"]
    D --> C["📜 Contract (CC-auth)<br><i>Incorporate living invariants</i>"]
    C --> M["🚀 Mission (M1)<br><i>Execute code & verify proof</i>"]
    M --> E["✅ Evidence & Archive<br><i>Mission complete, Proposal retired</i>"]
```

### First-Principles Synthesis: GTD, GSD, Superpowers, and Orca

Spectacular v2 synthesizes the proven strengths of the most effective task management and autonomous agent frameworks while replacing human discipline with **deterministic mechanical enforcement**:

```text
┌─────────────────────────────────────────────────────────────────────────────┐
│ 🧠 THE FIRST-PRINCIPLES PIPELINE (GTD → GSD → ORCA → SPECTACULAR)           │
├───────────────────┬───────────────────────────┬─────────────────────────────┤
│ Framework Action  │ Spectacular Primitive     │ Operational Mechanism       │
├───────────────────┼───────────────────────────┼─────────────────────────────┤
│ 1. Inbox / Triage │ proposals/                │ Cheap, mutable exploration. │
│    & Capture      │                           │ Brainstorming with zero     │
│    (GTD)          │                           │ authority or commitment.    │
├───────────────────┼───────────────────────────┼─────────────────────────────┤
│ 2. Clarify &      │ contracts/                │ Living, versioned contracts │
│    Organize (GTD) │ decisions/                │ freeze agreed invariants &  │
│                   │                           │ irrevocable owner choices.  │
├───────────────────┼───────────────────────────┼─────────────────────────────┤
│ 3. Engage &       │ missions/                 │ Bounded execution envelope  │
│    Momentum (GSD) │ (Iterate to exit 0)       │ with explicit write scope.  │
│                   │                           │ Passing test suite = proof. │
├───────────────────┼───────────────────────────┼─────────────────────────────┤
│ 4. Subagent       │ charter --prompt          │ Single-turn naked worker    │
│    Isolation      │ (400–600 token budget)    │ envelope. Zero governance   │
│    (Orca)         │                           │ lore, zero file wandering.  │
├───────────────────┼───────────────────────────┼─────────────────────────────┤
│ 5. Modular        │ .agents/skills/           │ High-density micro-kernel   │
│    Tooling        │ (Progressive Disclosure)  │ with on-demand reference    │
│    (Superpowers)  │                           │ injection per domain.       │
├───────────────────┼───────────────────────────┼─────────────────────────────┤
│ 6. Mechanical     │ spectacular guard         │ OS-level sandbox watchdog,  │
│    Enforcement    │ mission check --verify    │ surgical quarantine, and    │
│    (Spectacular)  │ replay:                   │ crash replay verification.  │
└───────────────────┴───────────────────────────┴─────────────────────────────┘
```

**Mission-scoped** — they live inside a Mission bundle and have no meaning
outside it: Objective, Run, Checkpoint, Evidence, Assessment, Review, Handoff,
Gap, and Mission-local Decisions.

That split is why archiving a Mission archives its whole bundle: an Evidence
record separated from its Mission is not evidence of anything.

## The Lean 3-Layer Autopilot Model

Spectacular is designed for fast, token-efficient autonomous execution with minimal ceremony. All governance reduces to 3 layers:

1. **Layer 1: Living Truth & Decisions**: `PROJECT.md` (boundaries/non-goals) + `.spectacular/decisions/` (bulk-ideated architectural choices recorded with `spectacular decide`).
2. **Layer 2: Topological Flight Plan**: Multi-session roadmap in `.spectacular/campaigns/` (4–8 macro milestone blocks; unstarted blocks remain 4-line lightweight draft cards).
3. **Layer 3: Single-File Execution Envelopes**: Compact, self-contained Mission files (`.spectacular/missions/M<N>-<slug>/M<N>-<slug>.md`, $\le 500$ tokens) with inline objectives, deliverable checklists, and fail-fast stop triggers.

### Zero Sub-Record Sprawl Policy
Never create separate `checkpoints/`, `assessments/`, `runs/`, `handoffs/`, or multi-page manual evidence files for routine code tasks. The test suite passing (`exit 0`) and clean Git commit **is** the proof. Context flows across subagents and parallel sessions via lightweight prompts ($\le 300\text{--}500$ tokens) and thread links (`conversation://<id>`).

### Supervised Dispatch vs. Full Ownership Handoff
- **Supervised Dispatch (90% Default)**: In-session subagent delegation. The Orchestrator retains active Mission ownership, launches a worker subagent with a $\le 300$-token charter, and waits reactively for completion (`worker_done`). Zero governance files are created.
- **Full Ownership Handoff (10% Transfer)**: Permanent ownership transfer across distinct sessions, human engineers, or different AI harnesses. Formally recorded via `spectacular handoff record`.

### The Escalation & Decision Gate Protocol
When an autonomous worker subagent discovers an unrecorded architectural choice or boundary conflict, it halts immediately (Fail-Fast Stop) and sends an escalation to the Orchestrator. The Orchestrator records the ruling atomically via `spectacular decide` (`.spectacular/decisions/D<N>.md`) and resumes the worker with the locked decision ID.

### Channel Separation: Durable Git State vs. Ephemeral Channels
- **Git is for durable truth**: `PROJECT.md`, `decisions/`, `campaigns/`, and single-file `missions/`.
- **Host channels are for live coordination**: Ephemeral pings, ask/reply loops, and task dispatch stay inside host harness tools (`invoke_subagent`, `send_message`, `conversation://<id>`) with zero file pollution in Git.

### The 3-Tier Question Escalator
- **Tier 1 (Optimistic Consent)**: Non-blocking 1-line default for low-risk implementation choices (`"Proceeding with X unless you prefer Y"`).
- **Tier 2 (Structured Batch Cards)**: Numbered questions with lettered options (`1. Question ➔ A, B, C (Recommended)`), calibrated context depth, batch shorthand replies (`A, B, A`, `all defaults`), and open write-in support.
- **Tier 3 (Trade-off Spectrum & Interactive Modals)**: Framing competing design axes for unpredictable exploration, or leveraging interactive UI modals (`ask_question`).

## Strategic Horizons: Anchors, Roadmaps, Campaigns, and Missions

Spectacular structures work across four distinct strategic altitudes:

1. **`PROJECT.md` (Anchor)**: System identity, mission, core boundaries, and strict non-goals.
2. **`ROADMAP.md` (Anchor)**: Macro-level product evolution and multi-quarter strategic horizons (6–18 months).
3. **Campaigns (`.spectacular/campaigns/`)**: Mid-term strategic arcs (2–6 weeks) with 4–10 sequenced blocks unblocking one milestone.
4. **Missions (`.spectacular/missions/`)**: Short-term atomic execution envelopes (1–3 days) with frozen verifiable claims and failable proof.

## Campaigns are plans, not records

A Campaign is an optional Markdown planning map in `.spectacular/campaigns/`.
It can sequence independent roadmap blocks and link candidate or active
Missions, but it grants no authority and is intentionally excluded from typed
CLI validation, identities, fingerprints, and lifecycle transitions. Campaigns
are mutable; a Mission's frozen envelope remains the execution authority.

A Mission may cite Campaign context in its Markdown body. Do not bind a Mission
to a Campaign in frozen frontmatter: roadmap edits must not create Mission drift.

## Atlases connect value to structure

An Atlas is an optional Markdown map for a coherent product-value slice. Its
outcome board shows the actor, journey steps, desired result, and success
signal. Its system board shows the capabilities, ownership boundaries,
dependencies, risks, and proof that make the result possible.

```text
desired future → journey step → capability → architecture boundary → proof
                              ↘ Campaign block → Mission
```

An Atlas makes this chain legible, but it is not a second source of execution
authority. Campaigns sequence potential work; Contracts state accepted behavior;
Missions freeze and authorize one implementation slice.

## The mechanical interface

The CLI command reference is generated into
[`mechanical-interface.md`](../skills/spectacular/generated/mechanical-interface.md).
The catalog is generated from the command registry, so it cannot drift from what
the binary does. When a document and the generated interface disagree, the
interface wins and the document is stale.

Commands are either `read-only` or `mutating`, and every mutating command is a
transaction: it either fully applies or fully rolls back. The test suite proves
this by injecting a fault at every write boundary and asserting the workspace is
unchanged.

Adding or changing a command requires owner approval.

## Give each agent only what it needs

Agents should not have to read the entire workspace to do one task. Spectacular
gives workers a short handoff with the approved work, allowed files, and named
sources. Larger details stay in linked files until they are needed.

## What the machine does and does not do

The CLI performs only the mechanics that are safer or cheaper done mechanically
than by an agent: validating a record graph, computing fingerprints, writing
atomically, resolving paths, and reporting drift.

It does not decide anything. It cannot complete a Mission, accept a Proposal,
grant authority, or judge whether work is good. Those are owner acts, and the
system's value comes from refusing to blur that line.

```text
        agent proposes and builds
                  │
                  ▼
    ┌─────────────────────────────┐
    │  CLI: validate, fingerprint │  mechanical — no judgment
    │  write atomically, report   │
    └─────────────────────────────┘
                  │
                  ▼
          owner decides at a gate     ← authority lives here
```

## Authority and gates

Authority is explicit and narrow. A Mission's frontmatter names what the
operator may do — inspect, edit in scope, run checks, generate derived files.
Anything outside that list requires an owner gate.

Three gates matter most:

1. **Activation** — freezes the agreement. Everything after is judged against it.
2. **Review** — proof is recorded as a record, with a verdict.
3. **Completion** — only the owner, and only when declared Gaps are closed.

A Gap is a stated limit, not a defect. It is never closed by deletion: its entry
survives with a written resolution, so the reason something was impossible stays
recoverable.

## Deterministic Replay & Crash Recovery Verification (`replay:`)

When autonomous agents implement stateful architectures (e.g. event sourcing, financial ledgers, materialized views, search indexes), they frequently create in-memory caches without persisting an underlying audit trail. If the cache is wiped or the process crashes, the system becomes corrupted.

Spectacular provides an optional `replay:` hook in `M<N>.md` frontmatter to enforce deterministic state reconstruction:

```yaml
replay:
  cache_paths:
    - "balances.json"
    - "balances.db"
  command: "python3 src/main.py reconcile"
```

### How Replay Verification Operates:
1. **Cache Eviction**: Spectacular deletes the declared `cache_paths` from disk.
2. **Replay Execution**: It invokes the declared recovery `command` against the raw event journal (e.g. `events.jsonl`).
3. **Equivalence & Test Assertion**: It verifies that reconstructed balances/indexes match reality and the test suite passes with `exit 0`.

This guarantees that derived data is completely ephemeral and that the raw journal remains the sole, immutable source of truth.

## Mechanical Perimeter Guard & Zero Wasted Work (`spectacular guard`)

When subagents are dispatched to implement an objective, they may occasionally produce unwanted side-effects outside their authorized perimeter (e.g., generating unneeded `.gitignore` files, polluting root configurations, or creating untracked scratch databases).

`spectacular guard <mission>/<objective> [--watch] [--json] -- <command...>` wraps any worker process in an OS-level sandbox watchdog:

```text
┌─────────────────────────────────────────────────────────────────────────────┐
│ 🛡️ SPECTACULAR GUARD PERIMETER ENGINE                                       │
├────────────────────────────────┬────────────────────────────────────────────┤
│ Mode                           │ Behavior & Enforcement                     │
├────────────────────────────────┼────────────────────────────────────────────┤
│ A. Post-Flight Watchdog        │ Default mode (0% runtime overhead). Runs   │
│    (spectacular guard M1/O1)   │ command, diffs snapshot on exit, auto-     │
│                                │ purges escaped files, exits 2 on escape.   │
├────────────────────────────────┼────────────────────────────────────────────┤
│ B. Real-Time Watcher           │ Spawns background filesystem listener. If  │
│    (spectacular guard --watch) │ the worker writes outside writes_paths,    │
│                                │ instantly kills rogue process (SIGKILL),   │
│                                │ purges illegal file, and emits refusal.    │
└────────────────────────────────┴────────────────────────────────────────────┘
```

### Surgical Quarantine (Zero Wasted Work)
When an escape occurs, `spectacular guard` does not wipe the worker's entire output. Instead, it performs a **surgical quarantine**:
1. **Purges Escaped Paths**: Deletes or restores *only* the files outside `writes_paths`.
2. **Preserves Valid Work**: Retains all changes and new code in authorized paths.
3. **Emits Session Continuation Prompt**: Outputs a tailored feedback turn (`FeedbackPrompt`) ready for headless session resumption (e.g. `claude -c`, `agy --continue`, `opencode --session`), allowing the worker to self-correct in seconds without starting over from scratch.

## The 4-Skill Starting Pack: Governance + Domain Execution

Spectacular ships with a curated suite of companion skills under `skills/` providing deep technical domain execution alongside mission governance:

| Skill | Role | Focus |
| :--- | :--- | :--- |
| **`spectacular`** | **Mission Governance Spine** | Orchestrates missions, contracts, objectives, flight plans, and subagent dispatch. |
| **`system-architecture`** | **System & Software Architecture** | C4 context/container models, bounded contexts, service boundaries, and ADRs. |
| **`data-modeling`** | **Data Modeling & Database Jobs** | Conceptual (Chen) $\to$ Logical (Crow's Foot) $\to$ Physical SQL DDL/ORMs, indexing, and zero-downtime migrations. |
| **`rapid-prototyping`** | **3-Option Tracer Prototyping** | 3 tracer fragments (A/B/C) across a 5-tier fidelity ladder (Atom $\to$ Integration) with decision ledgers. |

Each skill operates as a high-density, independent micro-kernel. When an operation crosses domain boundaries, skills declare clean **Expansion Handoffs** rather than duplicating instructions.

## Freeze points

A completed Mission's contract fingerprint is a freeze point, not a stale
pointer. It records which agreement that Mission was executed against.
Amendments re-point only the live Mission; a completed one keeps its binding
forever. `mission check` reporting drift on a completed Mission is a notice, not
an error — the Mission stays valid.

This is the same instinct as the archive: a finished record is kept as it was,
not updated to match the present.

![Spectacular architecture](diagrams/architecture.svg)

## See also

- [Quickstart](quickstart.md) — run one Mission end to end.
- [Process](process.md) — the lifecycle and its gates in detail.

