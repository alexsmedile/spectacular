# Explore, plan, and genesis

## One-Shot Genesis (Zero-to-One Kickoff)

When starting a project from scratch, or when receiving an initial PRD or prompt of intent:

1. **Scan for Intake Inputs**: Check for existing project kickoff documents in the project root or workspace (e.g. `./PRD.md`, `scratch/PRD.tmp.md`, intake notes, or output from `write-prd`). Treat any starter PRD as an ephemeral launchpad, not an eternal file. Distill its 8 foundational dimensions losslessly into:
   - **The Core Triad of Anchors**:
     * `.spectacular/PROJECT.md`: What & Why (Core scope, system boundaries, and strict non-goals).
     * `.spectacular/STACK.md`: What with (Languages, runtimes, database, libraries, baseline test command).
     * `.spectacular/ARCHITECTURE.md`: How (Directory layout, layers, component boundaries between DB, Server, API, and Domain).
   - **On-Demand Anchors (Earned only)**:
     * `.spectacular/VOCABULARY.md`: Defined terms and ubiquitous language, only when domain ontology or state machines are rich.
     * `.spectacular/SECURITY.md`: Non-standard compliance, isolation, or secret rules, only if project-specific.
     * `.spectacular/GUARDRAILS.md`: Custom AI operational rules, only upon explicit owner request.
     * `.spectacular/PRODUCT.md`: User personas, target market, and user journeys, only if separated from repository engineering.
   - **`M1-bootstrap` Mission Plan**: 2–4 executable claims with verifiable `pass_boundary` and `proof_requirement` directly derived from the PRD's measurable success criteria.

2. **Adopt Strong Defaults (Zero Grilling)**: Do not stall kickoff with multi-question interviews. Choose sane, production-grade defaults for toolchain and architecture; prompt only on irreversible semantic forks.

3. **Present One Genesis Preview**: Show the Core Triad summary and `M1-bootstrap` plan once in chat. On owner confirmation ("Yes" / "Proceed"), write the Core Anchors and activate `M1` with `spectacular mission start`.

For live templates and layouts, see [genesis-examples.md](genesis-examples.md).

## Campaign Planning (Mini-Roadmaps & Campaign Blocks)

A **Campaign** (synonyms: *Initiative*, *Milestone Arc*, *Flight Plan*, *Theme*) organizes 4–10 macro-concepts into a topological dependency sequence before freezing individual Missions. It is a high-level overview, not a heavy ceremonial document.

### What a Campaign Declares:
1. **Strategic Goal**: The overarching milestone outcome (e.g. *Launch Background Job Engine from zero to production*).
2. **Topological Map**: A dependency flow diagram (Mermaid primary, ASCII fallback) showing logical unblocking order.
3. **Exit Condition**: The observable milestone state certifying the campaign is achieved.

### What a Campaign Block Declares (The 4 Fields):
Each Campaign Block is an essential 4-field card (or diagram node), not a separate file:
- **Title / Theme**: Named macro-capability (e.g. `Block 2: Ingestion HTTP API`).
- **Capability Unlocked**: What the system will observably do when this block closes.
- **Prerequisites**: Upstream blocks that must be closed first.
- **Status & Mapping**: `PLANNED` | `IN PROGRESS -> M<N>` | `CLOSED (M<N>)`.

### Division of Responsibility:
- **Campaign / Block (Fluid Overview)**: Macro-capabilities, unblocking dependencies, fluid future blocks. Lives inline in chat/preview or as a 15-line `## Active Campaign` section in `.spectacular/PROJECT.md` (no extra files).
- **Mission (Frozen Execution)**: Atomic, frozen envelope with UUIDv7 identity, SHA-256 fingerprint, exact Git baseline, and strict verifiable claims (`pass_boundary` & `proof_requirement`).

### Block-to-Mission Mapping:
- One Campaign Block can resolve into one or more Missions upon execution (e.g. *Stripe Engine* $\to$ `M2-stripe-webhook` + `M3-subscription-lifecycle`).
- Multiple small adjacent blocks can be closed by a single cohesive Mission.
- Downstream blocks remain fluid until upstream proof and evidence are earned.

### Visual Presentation:
- **Mermaid Flowchart (Primary)**: Rendered in chat and markdown using `flowchart LR` or `flowchart TD` with block status (`CLOSED`, `IN PROGRESS`, `PLANNED`).
- **ASCII Diagram (Fallback)**: Compact plain-text layout for CLI and terminal logs.

## Explore

A Proposal is optional. Use one only when the exploration deserves a durable home
— for the problem, alternatives, open questions, research, or a draft
specification. Otherwise leave it in an issue or in the conversation.

A Proposal is mutable. It is neither authority nor current product truth.

Read the current Contracts and specifications before proposing observable
behavior. Once a direction is frozen, edit those files as ordinary Mission work.
There is no separate reconciliation lifecycle.

Check a Proposal that has a durable home:

```bash
spectacular proposal check <ref>
```

## Plan

Compare only approaches that genuinely differ and are outcome-sized. Weigh each on:

- observable result, and the proof it would need
- coherence with what exists
- dependencies and reversibility
- learning value
- integration path
- what happens if it is cancelled

Record one verdict: `sufficient | needs-evidence | needs-decision`.

Then grill only what is still unresolved — criteria, scope, dependencies, risks,
or blocking Gaps. Do not re-interview settled ground.

## Freeze a compact Mission preview

Frontmatter:

- title, owner, outcome, applicable Contract, Git baseline
- one completion claim per verifiable domain, each with a pass boundary and a
  proof requirement
- review level: `automatic | clustered | independent`, defaulted once when shared.
  Choose `independent` when any claim touches security, privacy, or rights;
  stored data or a migration; a shared or public interface; compatibility; more
  than one system boundary; an external provider; a destructive or
  hard-to-reverse effect; a production or observational claim; a material
  architecture change; a novel pattern; or evidence only the executor can see.
  Also choose it when the work is disputed. Otherwise `automatic` is honest, and
  `clustered` fits several small related claims. A reviewer who did not implement
  the scope is what makes it independent — see [close.md](close.md).
- Objectives, with dependencies and claim coverage
- initial Run and operator, authority, mechanical and semantic scope
- budgets, dependencies, Gaps, stops, recovery
- `resolves_gaps:` when the Mission closes a Gap on its bound Contract, as `gap`
  and `resolution` pairs. Both are frozen, so the owner approves the exact wording
  at activation and the Mission cannot gain amend authority afterwards. Completion
  refuses while a declared Gap is still open. Requires `amend-contract` in
  `requires_owner`.

Markdown body:

- origin and rationale
- the detailed execution plan
- conditional bootstrap and review notes

A claim is the part most often written too vaguely. It needs a boundary that can
fail, and a proof that names the test:

```yaml
completion:
    - claim: drift-flags
      pass_boundary: Each frozen completion claim carries named drift flags
          derived from repairs consumed, evidence age, verdict state, and
          fingerprint age.
      proof_requirement: Table-driven fixtures with known repair counts and
          evidence ages assert the exact flag set, the ranking, and the default
          selection including tie behavior.
```

`pass_boundary` states what must be observably true. `proof_requirement` states
what would demonstrate it. "Works correctly" is neither.

Present the preview **once**, in chat. Owner confirmation freezes the semantic
envelope.

The preview is a plan document, not yet a Mission. It carries no UUID, no
activation block, and no fingerprint — `mission start` generates those. For the
field-by-field shape of what it becomes, see
[mission-anatomy.md](mission-anatomy.md).

## Then activate

```bash
spectacular mission start plan.md --json   # or: ... start - --json  (stdin)
spectacular mission check <ref> --json     # confirm what was generated
```

It generates identities, bindings, activation, and the canonical path at
`.spectacular/missions/<slug>/<ref>-<slug>.md` — atomically, from the approved plan.

Under a declared manual bootstrap, hand-author that file, generate valid
identities, and verify the structure directly — see [bootstrap.md](bootstrap.md).
