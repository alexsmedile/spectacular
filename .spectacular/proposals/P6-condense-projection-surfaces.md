---
type: Proposal
id: 01a00a98-32b3-7612-b19a-b8ffa479505c
title: Condense projection surfaces
status: draft
human_ref: P6
created_by: Alex
created: "2026-08-16T12:42:21Z"
updated: "2026-08-16T12:42:21Z"
scope:
    - v2
target_contract: Contract:01a00a20-63dd-7670-97f1-9eb8e12adc3a
---

# Condense projection surfaces

Exploration for a possible Mission. Nothing here is frozen. The rendering rules are
firmer than the glyph choices, and the glyph choices are firmer than the exact
notation. Anything below may be dropped, split, or reversed at plan-freeze.

## Where this came from

A survey of high-density context representations for agents — compact serializations,
symbolic logic, decision tables, state-transition notation, graph syntax, telegraphic
prose. Most of that catalogue does not apply here and should be explicitly declined,
which is the first thing worth recording.

Spectacular is already lean where those techniques usually pay. `SKILL.md` is 110
lines. The six workflow references total 138. `MISSION.md` frontmatter is the primary
entry point and carries the whole control card. P3 and P4 already removed the obvious
slack from context loading and control ceremony. Re-serializing canonical files into a
denser format would cost readability and cold-recovery safety to save tokens that have
mostly already been saved.

What the survey did surface is a distinction the codebase already holds in one place
and does not yet exploit anywhere else.

## The distinction

`internal/projection` opens with: *"builds disposable, source-attributed recovery
views."* That doctrine is correct and under-used. It licenses a split the method does
not currently name:

| | Truth | Projection |
|---|---|---|
| Examples | `MISSION.md`, Contracts, Decisions | `mission show`, `mission check`, context bundles, diagrams |
| Optimized for | unambiguous meaning, human review, cold recovery | fitting a decision in one glance |
| Lifetime | durable, fingerprinted | disposable, regenerable |
| Verbosity | as long as meaning requires | as short as the decision allows |

This is orthogonal to the existing meaning-versus-mechanics split, and composes with
it. Meaning-versus-mechanics decides *who produces* a thing. Truth-versus-projection
decides *how dense it should be*. A projection can be aggressively condensed precisely
because it is disposable and carries a source pointer back to the record it summarizes.

The claim of this Proposal is narrow: condensation techniques apply to projections and
should be declined for truth.

## Problem

**`mission show` prints the record, not the state.** Current output for M6:

```text
M6 — Implement the compact Mission CLI
State: active
Outcome: Spectacular provides a small typed CLI that validates compact Mission
  bundles and performs only the repeated or atomic mechanics that are safer and
  cheaper than LLM-only execution.
Path: .spectacular/missions/M6-implement-compact-mission-cli/MISSION.md
Run: M6/R1 (awaiting-review)
Objective: M6/O1 — Build the shared typed Mission-bundle decoder, resolver,
  canonical writer, and schema registry. (implemented)
Objective: M6/O2 — Implement the minimal read, start, progress, expansion, review,
  and completion command surface. (implemented)
Objective: M6/O3 — Prove atomicity, safe refusals, representation equivalence,
  legacy readability, and compact distribution behavior. (implemented)
```

Nine lines that require the reader to reconstruct four facts the command already knows:
that every Objective is implemented, that the Run is therefore parked at review, that
review is the gate, and that the owner holds it. A resuming agent does that inference
on every cold start — the exact moment its context is scarcest and its risk of getting
it wrong is highest.

The data is all present in the bundle. The command declines to draw the conclusion.

**Objective order is stored but never drawn.** `after:` and `claims:` are per-Objective
YAML lists. Reading them means building a graph mentally and cross-checking two blocks
to answer "is any claim uncovered" or "what is startable now." On M6 the chain is
linear and this is merely tedious. On a branching DAG it is a real source of error, and
P5's dependency-split direction would make branching the normal case rather than the
exception.

**Authority is a list the agent must remember.** `authority.operator` and
`authority.requires_owner` are scanned linearly before every consequential action. The
failure mode is silent: an agent that does not re-read them proceeds on recollection.
The record is a decision table being stored as two flat arrays.

**The lifecycle has no picture.** Proposal to Mission to Run to Evidence to review to
completion is described in prose spread across `SKILL.md` and five references.
`docs/diagrams/` shows the approach works — the matrix-proposal-loop cycle diagram
carries a process in one screen that would take a page of prose — but nothing
equivalent exists for the core lifecycle or for an individual Mission.

## Directions worth exploring

Four, listed cheapest-first and independent. Adopting one does not require any other.
All four are projection-only: no canonical file changes shape, no schema field is
added, no fingerprint input moves.

### A compact state line for `mission show`

Lead with a single line that states position in the lifecycle, then the next action.
Detail stays available behind existing pointers.

```text
M6 active · R1 awaiting-review · O1✓ O2✓ O3✓ · claims 4/4 · repairs 1/3 · gaps 0 · deps M5✓
NEXT: independent review — owner gate · baseline e0e2fdc · activation 2bb4d9ff
```

Every field is derived from data the bundle already carries, and two of them are
cheaper than they look. `claims 4/4` is coverage of frozen completion claims by
Objective `claims:` lists — already computed, since `completion-claim-coverage` is a
mandatory validator in the Contract. The result is computed on every check and then
discarded rather than shown. `repairs 1/3` reads the Run counter against
`repair_budget`, which P5 separately notes is written and never consumed.

The `NEXT:` line is the load-bearing part. Owner-gate computation already exists in
`internal/projection` as `OwnerGate{Code, Detail, Source}`, and the compact CLI path
does not surface it. This direction is largely wiring an existing computation to the
most-read output rather than building new logic.

Open: whether `--json` keeps today's full record and the state line is human-only, or
whether the JSON gains a parallel `state_line` field so an agent reading JSON gets the
same conclusion. Leaning toward the latter, since agents are the primary reader.

### Render the Objective graph by shape

The notation should be selected from the graph, not chosen by a flag. Two tiers cover
everything the current model can express.

**Tier 1 — linear chain.** Every node has in-degree and out-degree at most one. Render
inline in the state line, as above:

```text
O1✓ → O2✓ → O3▶
```

**Tier 2 — branching DAG.** Any fan-out or fan-in. Render as a small ASCII graph:

```text
O1✓ ─┬─ O2✓ ─┬─ O4▶
     └─ O3◐ ─┘
```

When width exceeds the terminal, fall back to level sets, which stay readable at any
branching factor:

```text
L0: O1✓
L1: O2✓ ∥ O3◐
L2: O4·  (after O2,O3 — blocked on O3)
```

Status glyphs, one character, no color dependency:

| Glyph | Meaning | Derived from |
|---|---|---|
| `✓` | implemented | `objective.status` |
| `◐` | in progress | `run.current_objective` |
| `▶` | ready — dependencies met, not started | `after:` satisfied |
| `·` | blocked — dependencies unmet | `after:` unsatisfied |
| `✗` | failed or gapped | Gap scoped to the Objective |

`▶` versus `·` is the distinction with the most operational value and the one no
current output makes: it answers "what could a second operator pick up right now"
without reading the graph by hand.

A `--graph` flag emits the same DAG as mermaid for human reading. Generated on demand,
never stored, so it cannot drift from the bundle.

Open: whether Tier 2 should ever appear in the state line or always break to its own
block. Leaning toward always breaking, since a wrapped ASCII graph is worse than no
graph.

**Dependency on P5's edge split.** The tier test above reads one edge set. If P5's
`after_interface:` direction is adopted, three things in this direction change and
should be settled together rather than sequentially:

- The shape test must run over the union of both edge sets. M6 is a strict `after:`
  chain today and selects Tier 1; under the split it becomes branching and selects
  Tier 2. Tier selection is not stable across that change.
- The glyph set needs to distinguish artifact-ready from interface-ready. Both satisfy
  "dependencies met," but only the second is startable at activation, and that
  difference is the entire value of the split. Five glyphs cannot express it — either a
  sixth is added or readiness is carried on the edge rather than the node.
- ASCII rendering needs two edge styles (`─` for `after:`, `╌` for interface edges),
  and `--graph` needs the mermaid equivalent.

The direction of the coupling is worth stating plainly: the split makes this rendering
*more* valuable, not less. A linear chain is adequately served by inline arrows. Graph
rendering and the ready-versus-blocked distinction earn their cost precisely when
branching is the normal case, which is what the split produces.

### Answer authority questions instead of listing them

Add a read-only lookup against the existing `authority` block:

```text
$ spectacular authority check M6 push
DENIED · requires_owner · push · owner gate — return to owner

$ spectacular authority check M6 edit-in-scope
ALLOWED · operator · edit-in-scope
```

No schema change; the arrays are already there. This moves the check from "the agent
remembers the list" to "the agent asks," which is the same move the refusal format
already makes for validation. It also gives the two arrays a single rendering as a
decision table for human review, rather than two blocks that must be read together.

**This direction conflicts with the Contract as written and cannot be adopted in this
form.** `CC-missioncli` enumerates a closed ten-command surface and requires the CLI to
avoid added ceremony; M6's stops name a growing command surface as a stop condition.
`authority` would be an eleventh command and a new noun. Three ways out, in preference
order:

1. Fold the answer into existing output. `mission check` already validates
   `authority-vocabulary`; it could render the two arrays as a decision table without
   any new command. Cheapest, and stays inside the surface.
2. Add it as a flag on an existing command — `mission check M6 --authority push`.
   Cheap, but flags are how command surfaces grow sideways.
3. Amend the Contract to admit a read-only `authority` noun. Honest, but a Contract
   amendment for a convenience lookup is poor value.

Leaning toward the first. The underlying problem — authority read as a remembered list
rather than an answered question — is real regardless of which shape solves it.

Open: whether an unknown verb should refuse or default to `requires_owner`. Refusing is
more honest; defaulting to owner is safer. These conflict and the choice belongs to the
owner.

### One generated lifecycle diagram

Add a mermaid diagram of Proposal → Mission → Run → Evidence → review → completion to
`SKILL.md`, with the owner gates marked, and delete the routing prose it makes
redundant. `docs/diagrams/matrix-proposal-loop-cycle.md` is the format precedent.

The condition that makes this worth doing rather than decorative: it must be generated
from the command registry and state model, not hand-drawn. A hand-drawn diagram is one
more thing to keep in sync with `internal/command.Registry`, and drift in a diagram is
worse than no diagram because it is trusted at a glance. `internal/command/catalog.go`
already generates the mechanical interface table from the registry, so the precedent
for generated documentation exists in-tree.

**The diagram belongs in `SKILL.md`, outside `.spectacular/`.** The Contract requires
avoiding a `per-Mission index`, and M5's second repair deleted Mission-local `index.md`
generation on the principle that `MISSION.md` is already the bundle index and a second
generated navigation artifact beside a canonical one is duplication. That principle
reaches any generated artifact placed inside a Mission bundle, so this direction must
not produce one. A diagram in the Skill describes the method rather than an instance,
which is a different thing and not covered by that ban — but the distinction should be
stated at plan-freeze rather than assumed.

Open: whether the state model is currently explicit enough to generate from, or whether
transitions are implicit in `service.go` control flow and would need extracting first.
If the latter, this direction is more expensive than it looks and may not belong with
the other three.

## Constraints from `CC-missioncli`

Every direction here changes `show` or `check` output, so three Contract terms bind
this Proposal and should be treated as acceptance criteria rather than context.

**Inline and promoted must render identically.** A stress property requires that
"inline and promoted representations produce the same show, dependency, claim, and
completion results." Objective graph rendering is exactly the kind of change that
breaks this quietly — a promoted Objective resolves through a different path in the
decoder, and a renderer that reads the wrong one produces a different graph for the
same logical Mission. Any rendering added here needs a golden test asserting
byte-identical output across both representations of one Mission.

**`show` and `check` stay read-only.** Required behavior: "Keep read-only `show` and
`check` free of canonical writes." This rules out caching a rendered graph or diagram
into the workspace. Renderings are computed per invocation and thrown away, which is
consistent with treating them as projections and is the reason `--graph` emits to
stdout rather than to a file.

**The command surface is closed.** Ten commands are enumerated and M6 stops on a
growing surface. Three of the four directions here change the output of existing
commands and are unaffected; only `authority check` proposes a new noun, and it is
reshaped above to stay inside the surface.

## Explicitly declined

Recorded so a future Mission does not treat the source survey as a backlog.

**Re-serializing canonical files into TOON, compact YAML, or delimited rows.** The token
saving is real and the cost is worse. `MISSION.md` is read by humans during review, by
agents during cold recovery, and by `git diff` during every commit. All three degrade
under a denser encoding, and the third silently.

**Telegraphic or headline syntax in `completion`, `stops`, or Markdown rationale.**
These are the frozen semantic envelope. Ambiguity there is precisely the expensive
failure the freeze exists to prevent. `pass_boundary` and `proof_requirement` should
stay long enough to be unambiguous, and their current verbosity is correct rather than
slack.

**Symbolic logic notation for `authority` or `stops`.** `A ∧ ¬B ⇒ C` compresses well
and reviews badly. Owner review is the control that makes authority meaningful, and a
notation the owner must decode before disagreeing with weakens it.

**Hash-and-offload for tool output.** Standard advice for context-heavy agent loops, and
already handled here: the batched tree-bound gate plus reading detailed logs only on
failure is the same optimization arrived at by a different route.

**Delta-only state between sessions.** Genuinely attractive, and it is P5's preflight
trace direction rather than this one. If both are adopted, the trace supplies the delta
and this Proposal supplies its rendering.

## Gap — concurrent-Run timelines

The natural third tier of the graph rendering is a timeline across concurrently live
Runs with different operators:

```text
R1 codex   │ O1✓ ────── O2✓ ──── O4▶
R2 claude  │      O3◐ ─────────────┘
           └─ join at O4
```

This is not renderable today and is deliberately out of scope. The Run model permits
exactly one live Run — `run start` completes the previous one, and `Run` carries a
single `current_objective`. Multi-agent concurrency therefore requires a Run-model
change touching `run start`, fingerprints, atomicity, and review boundaries, which is
a different Mission with a much larger blast radius than four projection changes.

Recorded as a Gap so the notation question is not silently reopened later. If P5's
dependency split lands, latent parallelism in the DAG increases and the pressure for
this will follow; the two should be evaluated together at that point, not now.

## Open questions for the owner

- Is one Mission the right container for all four, or should the generated lifecycle
  diagram split off? The first three read the bundle and render; the fourth may require
  extracting an explicit state model first.
- Should the state line appear in `--json` as a derived field, or stay human-only?
- Do the first three land in the current worktree, given they add no canonical fields
  and change no frozen semantics, or wait for a clean activation?
- Should `authority check` refuse unknown verbs or default them to `requires_owner`?
- If P5's edge split is likely, should graph rendering wait for it rather than shipping
  a tier test that the split immediately invalidates?
- For the authority lookup: fold into `mission check` output, add a flag, or amend
  `CC-missioncli` to admit a read-only `authority` noun?

## Relationship to current work

M6 completed with independent review and `repairs: 1`, so it does not absorb any of this. Its frozen
schema and completion boundary stay as activated. P3 and P4 are accepted and reflected
in the Skill; this explores the projection surfaces neither reached.

### Interaction with P5

This Proposal is projection-only and P5 changes what is computed, so most of P5 lands
underneath this one rather than against it. Direction by direction:

| P5 direction | Effect here | Kind |
|---|---|---|
| Preflight trace with decay | None. The trace is session-scoped; this renders bundle state. | Independent |
| Aim audits at drift | Additive. The per-claim drift signal is a field the state line would want to carry. | Complementary |
| Keep rejected approaches as fallbacks | Additive. Adds an owner-gate path at repair exhaustion that the `NEXT:` line must render. | Small |
| Split `after:` into two edge types | Changes tier selection, the glyph set, and edge styling. | **Coupled** |

Only the fourth is coupled, and it is coupled in the useful direction — it increases
branching, which is what makes graph rendering worth building. The consequence for
sequencing is that the graph-rendering direction should not be frozen before the edge
split is decided, because its tier test and glyph set both depend on how many edge
kinds exist. The other three directions here are unaffected by every P5 direction and
can be frozen independently.
