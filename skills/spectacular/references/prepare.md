# Explore, plan, and genesis

Use this when: Orchestrator preparing Genesis, exploration, Proposal, or compact Mission preview.

## One-Shot Genesis (Zero-to-One Kickoff)

When starting a project from scratch, or when receiving an initial PRD or prompt of intent:

1. **Scan for Intake Inputs**: Check for existing project kickoff documents in the project root or workspace (e.g. `./PRD.md`, `scratch/PRD.tmp.md`, intake notes, or output from `write-prd`). Treat any starter PRD as an ephemeral launchpad, not an eternal file. Account for these eight foundational dimensions: (1) owner and user-observable outcome, (2) scope boundaries and non-goals, (3) inputs, outputs, and core workflow, (4) stack, runtime, and dependency constraints, (5) architecture and component boundaries, (6) data, state, security, and privacy constraints, (7) acceptance criteria and failable proof, and (8) operations, failure recovery, and delivery constraints. Distill them losslessly into:
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

3. **Prove coverage, then preview**: Before presenting anything for approval, make a source-to-destination coverage pass: map every one of the eight dimensions to a Core/earned Anchor or an `M1-bootstrap` claim, or mark it as an explicit open Gap. Only when all eight are accounted for, show the Core Triad summary and `M1-bootstrap` plan once. On owner confirmation ("Yes" / "Proceed"), write the Core Anchors and activate `M1` with `spectacular mission start`.

For live templates and layouts, see [genesis-examples.md](genesis-examples.md) or run `spectacular mission start --help` to emit the exact `MissionPlan` YAML frontmatter template directly to stdout.

## Campaign Planning (Mini-Roadmaps & Campaign Blocks)

A **Campaign** (synonyms: *Initiative*, *Milestone Arc*, *Flight Plan*, *Theme*) organizes macro-concepts into a topological dependency sequence before freezing individual Missions. Aim for 4–10 blocks: fewer may not need a Campaign, while substantially more usually deserves separate arcs. It is a high-level overview, not a heavy ceremonial document.

Campaign checks are an orchestration/planning operation. A Mission worker follows
the Mission, Objective, and Run it was assigned; it reads a Campaign only when
the Mission body explicitly cites it for non-binding strategic context.

### Useful Campaign content

Use judgment; these are lightweight guidance, not a schema or validation gate:
1. **Strategic Goal**: The overarching milestone outcome (e.g. *Launch Background Job Engine from zero to production*).
2. **Topological Map**: A dependency flow diagram (Mermaid primary, ASCII fallback) showing logical unblocking order.
3. **Exit Condition**: The observable milestone state certifying the campaign is achieved.

### Useful Campaign Block fields
Each Campaign Block is usually a compact card (or diagram node), not a separate file:
- **Title / Theme**: Named macro-capability (e.g. `Block 2: Ingestion HTTP API`).
- **Capability Unlocked**: What the system will observably do when this block closes.
- **Prerequisites**: Upstream blocks that must be closed first.
- **Status & Mapping**: `PLANNED` | `IN PROGRESS -> M<N>` | `CLOSED (M<N>)`.

### Division of Responsibility:
- **Campaign / Block (Fluid Overview)**: Macro-capabilities, unblocking dependencies, fluid future blocks. Lives in one optional Markdown file per independent arc under `.spectacular/campaigns/`; it never belongs in the stable `PROJECT.md` Anchor and grants no execution authority.
- **Mission (Frozen Execution)**: Atomic, frozen envelope with UUIDv7 identity, SHA-256 fingerprint, exact Git baseline, and strict verifiable claims (`pass_boundary` & `proof_requirement`).

### Block-to-Mission Mapping:
- One Campaign Block can resolve into one or more Missions upon execution (e.g. *Stripe Engine* $\to$ `M2-stripe-webhook` + `M3-subscription-lifecycle`).
- Multiple small adjacent blocks can be closed by a single cohesive Mission.
- Downstream blocks remain fluid until upstream proof and evidence are earned.

### Visual presentation

Prefer a Mermaid flowchart when a Campaign has multiple dependencies; use
compact Markdown or ASCII only when it is clearer. The visual is a planning
projection, not authority. When the Campaign uses the documented frontmatter
map, run `spectacular campaign check <path>` to validate the current block and
order, then render the Mermaid projection.

## No-mission lane

Micro-tasks need no Mission when the owner approves: single-file edits, typo
fixes, doc passes, localized config tweaks — one domain, no semantic ambiguity,
no architectural risk. State the intended action, get the approval, apply the
reuse ladder anyway, do it directly, verify with the STACK.md baseline command.
Escalate to a Mission when scope grows past that, alters a shared interface, or
hits a non-trivial failure. The Git side is the quick-patch exception in
[execute.md](execute.md).

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

### Simplicity First & The 7-Rung Reuse Ladder

When designing the implementation for a frozen claim, prefer the smallest coherent solution. Climb the reuse ladder from top to bottom before writing custom logic:

1. **Does this need to exist?** $\to$ If no, skip (YAGNI).
2. **Already in this codebase?** $\to$ Reuse existing helpers and types; do not invent parallel abstractions.
3. **Standard library does it?** $\to$ Use the stdlib directly.
4. **Native platform / runtime feature?** $\to$ Use native capability.
5. **Installed dependency does it?** $\to$ Use existing packages; avoid adding dependencies.
6. **One clear line?** $\to$ Write the inline expression directly.
7. **Only then**: The minimal coherent custom implementation that satisfies the claim.

#### Non-Negotiable Preserve-List
Simplicity must never compromise integrity. Never simplify away:
- Strict input validation and sanitization.
- Attributable error handling, propagation, and structured diagnostics.
- Security boundaries, cryptographic parameters, auth checks, and safe defaults.
- Data integrity, transaction boundaries, and rollback protections.
- Accessibility guarantees.

After freezing, scope cuts return to the owner.

### Upfront Architectural Grilling vs. Progressive Horizon Detailing
- **Upfront Architectural Grilling**: Settle foundational architectural choices that span multiple blocks (e.g. B1 through B7) early at the Campaign/Decision level. Ask focused decision questions before freezing execution blocks.
- **Progressive Horizon Detailing**:
  > *"Fully detail the active/next mission; keep downstream ones as drafts / sketches. Small missions should stay direct and lean."*
  - Detail **only** the active or immediate next Mission block. Downstream blocks remain lightweight draft sketches without premature claim matrices.
  - Token budgets from `.spectacular/config.yaml` govern document sizing:
    - **Active Mission**: 400 – 900 tokens (upper limit: 1,200 tokens).
    - **Draft / Sketch Mission**: 100 – 300 tokens.
    - **Decision**: 150 – 400 tokens.
    - **Context Charter**: $\le 1{,}200$ tokens (hard cap: 1,440 tokens).

## Freeze a compact Mission preview

### Plan Style & Authoring Guidance (Lossless Compression)
1. **User Superpower at the Center**: Lead with what the developer or user gets (the concrete superpower and observable benefit).
2. **Lean & Direct for Small Missions**: Small, single-objective missions should stay direct, lean, and free of artificial ceremony. Avoid filler diagrams or forced multi-pillar structures on routine tasks.
3. **Hub-and-Spoke for Complex Milestones (Recommended)**: For large, multi-objective campaign blocks, structuring claims as distinct technical pillars (e.g. Compiler, Budget, Proof) clarifies the architectural nodes.
4. **Lossless Information Compression**:
   - **Process Chains**: Use arrows over verbose narrative (`B1 → B2 → B3`, `active → paused | blocked → completed`).
   - **Matrix Tuples**: Use compact key-value lists for scope and stops (`mechanical: [cmd/, internal/]`, `stops: [rewrite]`).
   - **Canonical Pointers**: Reference records by identifier (`Contract:019f...`, `D21`) rather than duplicating text.
   - **BLUF**: Bottom Line Up Front — strip conversational preamble and state verifiable outcomes directly.
5. **Atomic Single-Invariant Claims**: One claim = **one** observable invariant + **one** verifiable proof check.
6. **Key Deliverables in Body**: Include a direct action checklist of target files and verification commands.

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
- Objectives, with dependencies, claim coverage, and optional `sources:` references
- initial Run and operator, authority, mechanical and semantic scope
- budgets, dependencies, Gaps, stops, recovery
- `resolves_gaps:` when the Mission closes a Gap on its bound Contract, as `gap`
  and `resolution` pairs. Both are frozen, so the owner approves the exact wording
  at activation and the Mission cannot gain amend authority afterwards. Completion
  refuses while a declared Gap is still open. Requires `amend-contract` in
  `requires_owner`.

Markdown body:

- `## Purpose & Scope`: Concise 2-3 sentence overview.
- `## Key Deliverables & Actions`: Direct file-by-file action checklist.
- origin and rationale, if non-obvious.

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

**Adapting proof to the project**: A proof requirement does not demand an elaborate test framework. If the repository lacks automated test suites (e.g. documentation, static sites, CLI scripts, prototypes), name the most direct, objective verification method available:
- **Build / Run checks**: `npm run build` compiles clean; `python script.py --test-flag` exits 0.
- **Reproducible scratch fixtures**: A small verification script in `scratch/` testing exact input/output pairs.
- **Structural / Schema checks**: `spectacular mission check <ref>` validates cleanly; `markdownlint` passes; internal links resolve.
- The invariant is simply that the proof is **failable, objective, and attributable**—never mere self-assertion.

## Reference Attachments Pattern (Progressive Disclosure)

When a Mission requires extensive technical details (such as a 200-line JSON schema, SQL DDL migrations, long API fragments, or complex sample payloads), **never inline them directly into the core Mission file**.

Inlining large payloads bloats the active Mission past the 900-token threshold and pollutes the agent's attention on unrelated objectives.

**The Solution: Reference Files**:
1. Place large technical specifications in `references/` (e.g. `.spectacular/missions/M.../references/schema.json` or workspace `references/`).
2. Link the reference in the Mission's `sources:` list or the specific Objective's `sources:` list:
   ```yaml
   objectives:
     - outcome: Implement user schema validation.
       claims: [user-schema-validation]
       sources: [references/user-schema.json]
   ```
3. The core Mission stays lean (400–900 tokens), and the Context Charter compiler progressively resolves the exact reference only for the Objective that needs it.

## Contract Updates: Amend vs. Version Bump

Updating a Contract follows two explicit paths:

| Update Kind | When to Use | Action |
|---|---|---|
| **Path A: Amend** | Closing a declared Gap or updating editorial frontmatter | Run `spectacular contract amend <contract-ref> --gap <gap-ref> --by <owner>`. Does not bump `contract_version`. |
| **Path B: Version Bump** | Changing behavioral invariants, command surfaces, or schema | In an authorized Mission, edit `.spectacular/contracts/<file>.md`, increment `contract_version: "N" -> "N+1"`, and check with `spectacular mission check`. Past missions stay frozen to their old version SHA. |

Present the preview **once**, in chat. Owner confirmation freezes the semantic
envelope.

The preview is a plan document, not yet a Mission. It carries no UUID, no
activation block, and no fingerprint — `mission start` generates those. For the
field-by-field shape of what it becomes, see
[mission-anatomy.md](mission-anatomy.md).

## Then activate

After owner confirmation, run the mutating launch once, then verify the generated Mission with the read-only check:

```bash
spectacular mission start plan.md --json   # or: ... start - --json  (stdin)
spectacular mission check <ref> --json     # confirm what was generated
```

It generates identities, bindings, activation, and the canonical path at
`.spectacular/missions/<slug>/<ref>-<slug>.md` — atomically, from the approved plan.

Under declared `manual-bootstrap`, draft the future shape outside the governed
lifecycle; do not generate identities or activate it by hand — see
[bootstrap.md](bootstrap.md).
