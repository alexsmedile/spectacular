---
type: synthesis-checkpoint
checkpoint: 002
sources: [source-001, source-002, source-003, source-004, source-005, source-006]
concepts: 79
human_dispositions: 0
updated: 2026-08-07
---

# Synthesis checkpoint 002

## What changed since checkpoint 001

Source 005 supplied code-backed CLI defects and a v2 surface proposal. Source 006
supplied a coherent replacement product loop centered on autonomous Missions. The
database now contains enough material to expose a foundational product fork, but
not enough human disposition to choose it.

## Strong convergence

Across otherwise different proposals, the following outcomes repeatedly survive:

1. **A small trustworthy context baseline** should orient a cold agent without
   loading the entire workspace.
2. **Capability-level behavioral contracts** should express intended truth while
   code/tests remain implementation truth.
3. **Work should be approved as outcomes, boundaries, and evidence**, not predicted
   file edits.
4. **Execution context must be bounded and reconstructable**, with explicit Git,
   permission, retry, and stop semantics.
5. **Evidence must map to promised behavior**, not merely list commands that ran.
6. **Closure must reconcile current contracts and preserve history.**
7. **Cold-agent continuation is a primary product test**, not incidental session hygiene.
8. **Spectacular should not build a model host, queue, daemon, or general agent platform.**
9. **New structure must be earned by repeated need**, not scaffolded or registered speculatively.

These are protected-outcome candidates. They do not imply agreement on filenames,
commands, lifecycle states, or who executes work.

## The three coherent worlds now visible

### World A — Lifecycle control layer

- Preserve the current PRD boundary: the skill judges and orchestrates; the CLI is
  deterministic and mechanical.
- Simplify requests, specs, status, evidence, skill routing, scaffold, and CLI.
- Host coding agents perform implementation under Spectacular's durable conventions.
- Autonomy remains a bounded workflow/recipe, not a product-owned runtime.

**Strength:** closest to current contracts and easiest to reduce incrementally.

**Risk:** the end-to-end experience may remain fragmented across many artifacts and commands.

### World B — Guarded mission runner

- Replace the main product model with PROJECT/SYSTEM, capabilities, Missions,
  RUN/EVIDENCE, and archive.
- `mission plan|run|close` owns planning, autonomous execution, proof, and reconciliation.
- One lifecycle and one coding agent replace many specialized records and workflows.

**Strength:** clearest single product loop and strongest user-facing story.

**Risk:** reverses current CLI, orchestration, mutation, approval, and workspace contracts;
it is a product rewrite, not a refactor by subtraction.

### World C — Durable mission contract, host-owned execution

- Adopt capability-centered contracts, a compact Mission delta, clause-mapped
  evidence, closure reconciliation, and cold resume.
- Keep agentic planning/execution inside the selected host runtime or skill.
- Spectacular persists authority and state but does not launch or schedule agents.
- Retain or simplify current artifacts only where they implement unique behavior.

**Strength:** captures Source 006's coherent loop while respecting the “no agent
platform” and mechanical CLI boundaries.

**Risk:** may become an awkward compromise unless command ownership and the Mission
artifact are designed as one consistent system.

World C is an inferred comparison option, not a recommendation or accepted merge.

## Independent stabilization lane

PZL-047 and PZL-048 are correctness/safety defects, not v2 preferences:

- canonical `kind` records are read through legacy `type` fields in Wayfinder;
- cleanup can delete a remote branch without separate remote-deletion consent.

They can be fixed before the product-world decision without adopting Source 005's
surface plan or Source 006's Mission model.

## Decisions that must not be bundled

1. **Product responsibility:** durable control layer, guarded runner, or host-owned execution.
2. **Truth model:** accepted behavior versus code/tests versus observed reality versus history.
3. **Canonical work unit:** current SPC + request, one Mission, or a Mission projection over them.
4. **Approval model:** which gates are distinct authorities and which are duplicate ceremony.
5. **Lifecycle granularity:** which types/states change behavior, safety, evidence, or retrieval.
6. **Artifact authority:** replace, project, or retain current Anchor/spec/request/evidence files.
7. **Command authority:** mechanical CLI versus agentic skill/runtime and the public v2 grammar.
8. **Migration policy:** compatibility, recovery, and treatment of already-shipped capabilities.

## Current recommendation about process

Continue source intake without implementation or product dispositions. When intake
closes, begin with one decision packet: **Does Spectacular own autonomous execution?**
That answer constrains vocabulary, workspace shape, CLI grammar, AFK/worktree scope,
approval gates, lifecycle, and implementation architecture. Deciding names or file
trees first would force this choice implicitly.
