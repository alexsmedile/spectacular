---
name: spectacular
version: 2.1.0
description: Guide work in a canonical Spectacular v2 workspace through optional Proposal exploration, compact Mission preparation and activation, governed execution, earned Objective/Run expansion, Evidence and review, owner completion, audit, and cold recovery. Use for `/spectacular` jobs such as orient, explore, propose, plan, start, resume, handoff, review, complete, or audit; for compiling bounded runtime context or an Autopilot charter; and for safely continuing after session or runtime replacement.
---

# Spectacular

Operate one bounded Mission from accepted project truth. Keep semantic judgment in the Skill and
host session, deterministic invariants in supported mechanical tooling, execution in the host
runtime, and provider effects in their native providers.

## Working model

| Surface | Responsibility |
|---|---|
| Proposal | Optional, mutable exploration in chat, an issue, or a Spectacular file |
| Mission | Frozen execution plan and primary entry point in `MISSION.md` |
| Contract/specification | Accepted product behavior; edit as ordinary Mission work |
| Decision | ADR-like durable choice and rationale, never routine lifecycle approval |
| Objective / Run | Inline while simple; promoted to a file only when independently useful |
| Evidence / review | Earned proof and assessment; neither silently changes success criteria |

## Start every workflow

1. Discover the explicit v2 workspace. Read `.spectacular/PROJECT.md` first, or root `PROJECT.md`
   only when the nested file is absent, then read the named Anchors and relevant Contracts.
2. Run one read-only launch preflight: workspace; Git branch, cleanliness, upstream/default-branch
   freshness and release state; selected Mission; Contract/baseline bindings; owner and activation
   fingerprint; validation mode; current Objective/Run; blocking dependencies, Gaps, and stops.
   Report a plain outcome, one technical evidence line, and one clean next action.
3. Enter one Mission through its `MISSION.md`. Read its compact frontmatter control card and body;
   follow promoted Objective, Run, Evidence, or review pointers only when the current work needs
   them. Do not preload project history, every record, or generated catalogs.
4. If the Mission declares supported mechanical validation, use the typed `show`/`check` command.
   If it declares `manual-bootstrap`, treat incompatible CLI mutation and validation as out of
   service: edit canonical Markdown directly and check YAML, UUIDv7 identity, refs, fingerprints,
   claim coverage, dependency DAG, Run state, authority, scope, and file layout without citing the
   old CLI as proof.
5. Load exactly one routed workflow reference below. Revalidate exact bindings, authority, budgets,
   and stops before consequential effects or resumed work.

## Divide meaning from mechanics

The plan supplies meaning: outcome, completion criteria, decomposition, semantic scope, authority,
dependencies, Gaps, stops, rationale, and prose. The LLM may draft or directly edit these canonical
files when that is the fastest clear path, especially during exploration or a declared bootstrap.

Mechanical tooling supplies repeatability: schema and vocabulary validation, UUIDv7/ref allocation,
fingerprints, baseline checks, dependency integrity, safe paths, atomic multi-file transitions,
concurrency/retry behavior, compact projections, and exact refusals. A Mission never chooses which
mandatory validators apply; the active schema registry owns them.

Use scripts or CLI when failure is expensive, the rule is exact and repeated, or a transition must
be atomic. Use LLM judgment when meaning is contextual, prose carries the value, several solutions
are valid, or encoding the work mechanically would cost more than checking the result.

## Mission structure

A simple active Mission begins with only `<mission>/MISSION.md`. Frontmatter holds identity, `ref`,
status, owner, exact Contract/baseline bindings, outcome, review level, frozen completion claims,
inline Objectives and current Run, activation fingerprint, validation mode, authority, mechanical
and semantic scope, budgets, dependencies, Gaps, and stops. Markdown holds origin, rationale,
detailed Objective plans, bootstrap conditions, examples, and review instructions.

Create `objectives/` when detail, delegation, concurrent ownership, or independent review earns a
separate Objective file. Create `runs/` when a distinct execution job, operator, baseline, or
recovery boundary earns another Run. Preserve UUID/ref identity and replace inline detail with one
pointer. `MISSION.md` remains the index; do not require a Mission-local `index.md`.

## Ask owner questions

Ask only when a semantic, Mission-boundary, authority, risk, irreversible-effect, or current
Contract fork remains. Use four short parts: plain outcome; technical basis; concrete
`action -> consequence` options; and the recommended default with its reason. Do not impose an
interview when the Mission is already sufficient. Make silent product assumptions explicit choices
or Gaps, and carry previously accepted defaults forward.

## Route the guided job

| Guided job | Load |
|---|---|
| `orient` | [orient.md](references/orient.md) |
| `explore`, `propose`, `plan` | [prepare.md](references/prepare.md) |
| `start`, `resume` | [execute.md](references/execute.md) |
| `handoff`, Autopilot | [runtime.md](references/runtime.md) |
| `review`, `complete` | [close.md](references/close.md) |
| `audit` | [audit.md](references/audit.md) |

## Authority and execution

- The Mission owner alone changes outcome, completion criteria, semantic scope, review independence,
  or the forbidden-effect ceiling. Record the owner and activation time plus a fingerprint over the
  frozen semantic envelope. Mutable Objective/Run progress is outside that fingerprint.
- The operator may choose reversible implementation attempts, tools, checks, and bounded repairs
  inside the Mission. Return scope expansion, irreversible/provider effects, exhausted repair, or a
  stop condition to the owner.
- Execute compactly: `Mission card -> current Objective -> exact sources -> work -> focused checks`.
  Batch cohesive work, then run one full tree-bound gate before independent review or completion.
- Fan out only outcome-sized, disjoint claim scopes with exact inputs, dependencies, authority,
  stops, and return contracts. Avoid recursive critic loops and tiny sessions.
- Keep Evidence, deterministic checks, independent review, owner acceptance, and completion
  distinct. A green check proves only its stated observation.

## Continuity

Return Mission, owner, current Objective/Run, Contract and Git baseline, validation mode, review and
Evidence state, remaining dependencies/Gaps, repair use, recovery point, and exactly one safe next
action or owner gate. When Spectacular develops itself, the active Mission retains the schema and
completion boundary frozen at activation; later product changes apply to later Missions.
