---
type: Contract
id: 01a00aae-8921-7b27-96a9-1a4c175e7dc6
ref: CC-projsurf
title: Derived state and dependency shape
status: current
owner: Alex
created: "2026-08-16T13:06:45Z"
updated: "2026-08-16T13:06:45Z"
contract_version: "1"

purpose: Make the bundle state the system already holds legible at a glance, and let dependency declarations distinguish needing an artifact from needing only its frozen interface.
outcome: Projections state conclusions instead of records, audits aim at observed drift instead of memory, repair exhaustion returns a choice instead of a stop, and Objectives that need only a frozen interface start at activation.

supersedes_proposals:
  - Proposal:01a00a93-4757-7547-b64e-e91d2c291ce4
  - Proposal:01a00a98-32b3-7612-b19a-b8ffa479505c

applies_when:
  - A projection renders Mission, Objective, Run, claim, or authority state for a reader.
  - An audit target is selected without the owner naming one.
  - A Run exhausts its repair budget while frozen fallbacks exist.
  - An Objective declares a dependency on another Objective.
  - A Mission declares that it follows another Mission, or a reader needs the sequence of a multi-Mission job.
  - A Proposal is authored or checked.
does_not_apply_when:
  - A canonical record is being written, reviewed, or frozen. Condensation applies to projections and is declined for truth.
  - Concurrent live Runs are required; the Run model permits exactly one and that limit is out of scope here.
does_not_provide:
  - Owner judgment at any gate, automated matching of failures to fallbacks, concurrent-Run timelines, or generated lifecycle documentation.

required_behavior:
  - Lead compact projections with one state line stating lifecycle position, then a NEXT line stating the next action and who holds it.
  - Derive every state-line field from data the bundle already carries; add no canonical field to support rendering.
  - Surface the existing owner-gate computation on the compact command path rather than recomputing it per reader.
  - Report what changed since the previous session as a delta on the state line rather than re-deriving all preflight dimensions uniformly.
  - Compute per-claim drift as named flags derived from repairs consumed, evidence age and freshness, verdict state, and fingerprint age.
  - Default an unnamed audit target to the most-flagged claim, and show the flags that selected it.
  - Answer authority questions as a lookup against the declared authority block instead of requiring the reader to recall two arrays.
  - Render that lookup as a decision table from the `authority-vocabulary` validator already run by `mission check`, adding no command and no noun to the frozen surface.
  - Refuse an undeclared authority verb with code, field, problem, the declared vocabulary, and a safe correction.
  - Record each seriously-considered rejected approach at plan-freeze as approach, rejected_because, and invalidated_if.
  - Include fallbacks in the activation fingerprint from first release, so a fallback cannot be introduced mid-Run to escape a stop.
  - Hold invalidated_if to the same verbosity standard as pass_boundary and proof_requirement; it is frozen text that a gate reads, and ambiguity there is the expensive failure the freeze prevents.
  - Surface every recorded fallback at repair exhaustion, matched or not, alongside the failure that consumed the budget.
  - Permit ranking and recommending a fallback; never assert a match as fact, and never present a recommendation without the full set the owner can choose against.
  - Distinguish a dependency on a produced artifact from a dependency on a frozen interface, and treat the second as startable at activation.
  - Refuse an interface-only dependency whose target interface is not in fact frozen.
  - Select graph notation from the graph shape rather than from a flag.
  - Express readiness in a way that separates artifact-ready from interface-ready; a single ready state cannot carry the distinction the dependency split exists to create.
  - Carry Mission-to-Mission order as a list of Mission refs, validated for reference integrity and acyclicity.
  - Enforce exactly one rule on that order: refuse activating a Mission while any Mission it follows is not completed.
  - Leave existing prose `dependencies:` unchanged; it carries human conditions and is not the machine-checked edge.
  - Render multi-Mission sequence as an ASCII timeline drawn from Mission refs, status, creation, and completion time, generated on demand and never stored.
  - Render every projection identically whether an Objective is inline or promoted; the two decode paths must reach byte-identical output.
  - Resolve a Mission ref through one decoder that accepts `ref:` and legacy `human_ref:`, so Mission-order validation compares refs that are spelled one way regardless of how a record was authored.
  - Label a Mission in every view as `<ref> · <title>`; the ref is what the reader types, the title is what the reader recognises.
  - Provide four ASCII views: a Mission timeline, a compact Mission chain, an Objective graph, and Objective level sets.
  - Default the Objective view to the graph and fall back to level sets only when the graph would exceed the terminal width.
  - Render ASCII for straightforward shapes and escalate to mermaid only for graphs ASCII cannot render legibly, stating the reason for escalating.
  - Limit mermaid to two approved views: Objectives grouped by level with parallel work enclosed, and Mission edges distinguished by kind.
  - Treat every drawn example as guidance rather than a fixed template; adapt layout, grouping, and labelling to communicate a particular Mission, while holding the glyph vocabulary stable across all views.
  - Carry derived readiness into `--json` for Mission and Objective queries, so an agent reading JSON reaches the same conclusion a human reads from the drawing.
  - Define a compact Proposal schema matching hand-authored practice: type, id, ref, title, status, created_by, created, updated, scope, and target_contract.
  - Validate Proposals without providing a creation command; Proposals are authored as Markdown and checked, not generated through ceremony.
  - Require `ref:` on new records while still decoding legacy `human_ref:`, and report the legacy spelling as drift rather than refusing it.
  - Carry no candidate Contract body, freshness block, idempotency key, or authorization on a Proposal; a Proposal proposes, and the Contract it targets holds the frozen result.

command_surface:
  # No new command and no new noun. CC-missioncli's ten-command surface is closed,
  # and M6 names a growing surface as a stop. Everything below is added output or a
  # flag on a command that already exists.
  - mission show <ref>            # gains state line and NEXT line
  - mission check <ref>           # gains per-claim drift flags and the authority decision table
  - mission show <ref> --graph    # ASCII Objective graph, level sets when too wide
  - mission show <ref> --timeline # multi-Mission ASCII sequence view

mandatory_validation:
  - proposal-schema-v2
  - ref-spelling-drift            # legacy human_ref reported, not refused
  - fallback-fingerprint-coverage
  - interface-dependency-frozen-target
  - objective-dependency-dag        # extended to two edge kinds
  - mission-order-integrity         # refs resolve, graph is acyclic
  - mission-order-activation        # refuse activation while a predecessor is incomplete

stress_properties:
  - A projection never writes to the canonical tree; no graph, timeline, or diagram is cached in the workspace.
  - Inline and promoted representations of the same Mission render identically across every projection, including graphs and timelines.
  - A state line and the record it summarizes never disagree; the line is derived, never stored.
  - Adding or altering a fallback after activation invalidates the activation fingerprint.
  - Repair exhaustion with recorded fallbacks returns an owner choice; repair exhaustion without them returns the existing stop.
  - An Objective depending only on frozen interfaces is startable at activation regardless of artifact progress.
  - Graph notation is a function of graph shape alone; the same bundle renders identically for every reader.
  - A Mission whose predecessor is incomplete cannot activate; a Mission with no declared order is unaffected.
  - The timeline renders from existing Mission fields with or without declared order, and never writes.

conformance_checks:
  - P5 and P6 validate against the compact Proposal schema unchanged; P3 and P4 validate with their legacy fields preserved and their `human_ref` reported as drift.
  - Golden bundles render a state line whose every field is traceable to a bundle field.
  - Table-driven authority lookups cover declared verbs, undeclared verbs, and owner-gated verbs, asserted against `mission check` output rather than a separate command.
  - A Mission with inline Objectives and the same Mission with those Objectives promoted produce byte-identical `show`, `--graph`, and state-line output.
  - Drift flags are asserted against fixtures with known repair, evidence age, verdict, and fingerprint state.
  - Fingerprint tests prove a mutated fallback invalidates activation and mutable progress does not.
  - Dependency tests cover artifact edges, interface edges, mixed graphs, and an interface edge onto an unfrozen target.
  - Mission-order tests cover a dangling ref, a cycle, activation ahead of an incomplete predecessor, and a Mission declaring no order.
  - Timeline fixtures cover sequential, overlapping, and in-progress Missions.
  - View fixtures assert the approved notation for all four views, including the width threshold that selects level sets over the graph.
  - JSON output for Mission and Objective queries is asserted to carry the same readiness conclusion the drawn views show.

gaps:
  - ref: lifecycle-diagram-ungenerated
    problem: The lifecycle diagram is hand-maintained because no explicit state model exists to generate it from; transitions are implicit in service.go control flow.
    blocked_on: Extraction of a declarative transition model.
  - ref: concurrent-run-timelines
    problem: Timelines across concurrently live Runs are not renderable; the Run model permits exactly one live Run.
    blocked_on: A Run-model change touching run start, fingerprints, atomicity, and review boundaries.
  - ref: dead-v1-governance-code
    problem: ProposalInput, CreateProposal, and the candidate_* machinery in internal/governance are unreachable; proposal create is in the forbidden-command test and no current command calls them.
    blocked_on: A decision on whether removing them belongs with the Proposal schema work or with a separate cleanup.
  - ref: mission-ref-frontmatter-drift
    problem: M2, M3, and M4 carry `human_ref:` while M5 and M6 carry `ref:`, so Mission-order refs would otherwise be compared across two spellings.
    resolution: Closed by this Contract. M7 normalizes decoding through one path that accepts both and reports the legacy spelling as drift; M8's mission-order-integrity resolves refs through that decoder. Completed Missions are not rewritten.
---
# Derived state and dependency shape

## What this Contract is for

Two explorations converged on one observation: the bundle already knows more than
it says. Claim coverage, owner gates, repair state, and dependency structure are
all present and all left for the reader to reconstruct. A resuming agent performs
that reconstruction at the moment its context is scarcest.

This Contract governs two kinds of work. The first makes existing state legible
without adding canonical fields. The second adds one genuinely new distinction to
the dependency model, because a single `after:` edge conflates needing a produced
artifact with needing only a frozen interface — and that conflation is what
serializes work that could run in parallel.

## Truth and projection

`internal/projection` already carries the doctrine: *"builds disposable,
source-attributed recovery views."*

| | Truth | Projection |
|---|---|---|
| Examples | `MISSION.md`, Contracts, Decisions | `mission show`, `mission check`, diagrams |
| Optimized for | unambiguous meaning, review, cold recovery | fitting a decision in one glance |
| Lifetime | durable, fingerprinted | disposable, regenerable |
| Verbosity | as long as meaning requires | as short as the decision allows |

Condensation applies to the right column and is declined for the left. This is
orthogonal to meaning-versus-mechanics: that split decides *who produces* a thing,
this one decides *how dense it may be*.

The rule has one consequence worth stating, because it cuts against the instinct
that shorter is better everywhere: `invalidated_if` is frozen text. It sits inside
the fingerprint and a gate reads it at repair exhaustion. It stays long.

## Mission decomposition

Two Missions. The boundary is whether any canonical field changes.

**M7 — projection and the Proposal schema.** The compact state line, per-claim drift
flags, the authority decision table, a provisional lifecycle diagram, the normalized ref
decoder, and the compact Proposal schema. No Mission field is added, no fingerprint
input moves. The projection work reads state the bundle already carries; the Proposal
schema records a shape already in use and adds validation for a record type the Mission
bundle does not contain. Neither touches dependency shape, so M7 can be frozen without
settling it.

M7 adds no command. CC-missioncli enumerates a closed ten-command surface and M6 names
a growing surface as a stop condition, so every addition here is output on a command
that already exists or a flag on one. The authority answer in particular costs nothing
new: `mission check` already runs the `authority-vocabulary` validator, which holds the
full operator and requires-owner vocabularies and currently discards them after
returning pass or fail. Rendering that table is a change to what is printed, not to what
is computed or to what can be invoked.

The ref decoder belongs here for the same reason it is cheap: M2, M3, and M4 carry
`human_ref:` while M5 and M6 carry `ref:`, and M8 cannot validate Mission order against
refs spelled two ways. One decoder that accepts both and reports the legacy spelling as
drift closes that without rewriting any completed Mission's frontmatter.

The multi-Mission timeline is projection and belongs to M7. It draws from Mission ref,
status, `created`, and `completion_record.at`, all of which exist today, and needs no
schema change to be useful — a timeline over prose dependencies simply omits the edges.

**M8 — schema.** Frozen fallbacks, the artifact-versus-interface edge split, Mission
order as typed refs, and the graph rendering that depends on the split. All add or
alter canonical fields. Keeping them together holds one freeze boundary around every
change to the record: `fallbacks` enters the activation fingerprint, `after_interface:`
enters the Objective model, `after_mission:` enters the Mission model, and the
rendering cannot be frozen before the edge kinds it draws are decided. `validate.go`
currently validates one dependency graph and must validate two Objective edge kinds
plus a Mission-level graph, including refusing an interface-only dependency onto an
unfrozen target and refusing activation ahead of an incomplete predecessor. That is
real work and a real source of new refusal modes, and it earns its own freeze.

Once `after_mission:` lands, the M7 timeline gains edges without changing: it renders
whatever order information exists.

The coupling runs in the useful direction. A linear chain barely needs a drawing;
branching is exactly when a drawing earns its keep. The split makes the rendering
more valuable, not less.

## The Proposal record

A Proposal proposes a change and drives a decision. That is what P5 and P6 do, and
`type: Proposal` is correct for them. What is stale is the schema, not the practice.

Two shapes are in the workspace. P3 and P4 carry thirty-five frontmatter fields:
eleven `candidate_*` fields duplicating a Contract body into the Proposal, four
freshness fields, plus `base_fingerprint`, `authorization`, and `idempotency_key`. P5
and P6 carry ten: `type`, `id`, `ref`, `title`, `status`, `created_by`, `created`,
`updated`, `scope`, `target_contract`.

The second shape is the live one. `proposal create` is not merely absent from the
compact CLI — it is in the forbidden list in `internal/command/command_test.go`, with
a test asserting it never returns. `ProposalInput` and its thirty-five fields are
unreachable v1 governance code that no current command calls. Nothing validates a
Proposal today, so the ten-field shape is what practice has settled on in the absence
of enforcement.

The schema follows practice. The compact Proposal schema is the ten fields above, and
it is validated but never generated: a Proposal is authored as Markdown and checked,
exactly as a compact Mission is. No candidate Contract body, because the Contract it
targets holds the frozen result. No freshness block, idempotency key, or authorization,
because those are the ceremony the compact CLI already declined.

`human_ref:` to `ref:` is a rename that stopped halfway. P1 through P6, M2 through M4,
and CC-v2prod use `human_ref:`; M5, M6, and CC-missioncli use `ref:`. New records
require `ref:`, the decoder continues to read `human_ref:`, and a check reports the
legacy spelling as drift rather than refusing it. Completed and frozen records are not
rewritten to finish a rename.

P5 and P6 are the reference examples of the shape.

## Multi-Mission sequence

A multi-Mission job has no view. `.spectacular/missions/index.md` is a flat generated
table of every record sorted by ref; it shows what exists and nothing about order,
dependency, or where the job has reached. No Mission file contains a diagram of any
kind. M2 through M6 form a chain of related work that is drawn nowhere.

It is also not drawable today, because the edges are prose. The four Missions that
declare dependencies use three different notations for the same relation, mixed with a
baseline requirement that is not a Mission edge at all:

```yaml
- "M5 completed with independent review and owner acceptance."   # M6
- "M3 implementation commit 87255b9"                             # M4
- "v2.0.0-rc.2 published baseline"                               # M2, M3, M4
[]                                                               # M5, which followed M4
```

`Dependencies` is `[]string`. It enters the activation fingerprint and is never
validated — no reference integrity, no acyclicity — while Objective `after:` receives
both. A dangling or misspelled Mission dependency is currently undetectable, and M5
declares no predecessor despite following M4.

The fix is deliberately thin. Order becomes a list of Mission refs and carries exactly
one machine-checked meaning: a Mission may not activate while a Mission it follows is
incomplete. No per-edge state, no condition objects, no second vocabulary. The existing
prose `dependencies:` is left exactly as it is, because the detail it carries —
*completed with independent review*, *at commit 87255b9* — is human meaning, and
turning it into schema would freeze judgment the owner should keep making.

Everything beyond that one rule is a view. Mission ref, status, `created`, and
`completion_record.at` are already present on every Mission, which is enough to draw a
timeline showing sequence, overlap, and current position. It renders as ASCII, is
generated on demand, and is never stored — so it cannot drift from the bundle.
Readability by an LLM is the design target, not decoration.

## Approved views

Four views, approved as drawn. Notation is part of the Contract because a view whose
glyphs drift between commands is worse than no view. Every field below is derived from
data the bundle already carries; none is stored.

A Mission is labelled `<ref> · <title>` throughout. The ref is what the reader types;
the title is what the reader recognises.

**Glyphs.** One character, no colour dependency, shared across all four views.

| Glyph | Meaning |
|---|---|
| `✓` | implemented |
| `◐` | in progress |
| `▶` | startable now — dependencies met, not started |
| `·` | blocked — dependencies unmet |
| `✗` | failed or gapped |

**A · Mission timeline.** Sequence, overlap, and current position across a
multi-Mission job.

```text
                       Aug15    Aug16
                       ────────────────────────
M5 · Compact Missions     ████████        completed
M6 · Mission CLI             ██████████   completed
M7 · Derived state                 ███▓   active · R1 · O2
M8 · Frozen schema                   ░░░  blocked · after M7
M9 · Timeline render                 ░░░  blocked · after M7
                       ────────────────────────
                       └ M8, M9 concurrent after M7

████ done   ▓ in progress   ░░░ not started
```

**B · Compact Mission chain.** The same order without a time axis, short enough to sit
in a state line.

```text
M5✓ → M6✓ → M7◐ ─┬─ M8·
                 └─ M9·
```

**C · Objective graph.** The default view inside one Mission.

```text
M7 · Render derived state

  O1✓ ─┬─ O2◐ ─┬─ O5·
       └─ O3▶ ─┘
  O4▶
```

**D · Objective level sets.** Used only when C would exceed the terminal width, since a
wrapped ASCII graph is worse than no graph. Level sets stay readable at any branching
factor.

```text
M7 · Render derived state

  L0  O1✓ derivation layer
  L1  O2◐ state line  ∥  O3▶ drift flags  ∥  O4▶ authority lookup
  L2  O5· timeline    (after O2, O3)

  ∥ = independent, can run in parallel
```

Each Objective carries a short name of two or three words alongside its ref, so a
reader is not forced to look up what `O3` refers to. One line per level; the name is a
shortening of the Objective outcome, not a second stored field.

The load-bearing distinction across all four views is `▶` versus `·`. It answers "what
could a second operator pick up right now" without reading the dependency lists by
hand, and no current output makes it. On the Mission drawn above, `O3` and `O4` are
both startable while today's `mission show` prints both as `planned`.

## Choosing ASCII or mermaid

ASCII is the default and carries the simple, straightforward case. Mermaid is for
graphs ASCII cannot render legibly, and two mermaid views are approved for that
purpose:

- **Grouped by level, parallel work enclosed.** A subgraph per level, so the band of
  independent Objectives is visually enclosed and labelled rather than inferred from
  spacing.
- **Mission edges by kind.** Solid for a dependency on the produced artifact, dotted
  for one needing only the frozen interface. Used when a Mission splits into related
  Missions that follow it.

Both are drawn in full at `docs/diagrams/mermaid-view-candidates.md`. The two remaining
candidates there — a colour-coded Objective graph and a native Gantt timeline — are
declined, because the ASCII Objective graph and ASCII timeline already read better for
their cases.

The renderer chooses. These signals indicate a graph has outgrown ASCII, and any one of
them is reason to consider mermaid rather than a rule that compels it:

- the graph fans out or fans in rather than forming a chain
- two edge kinds are present
- the graph is more than about two levels deep
- the ASCII rendering would exceed the terminal width

When the renderer escalates to mermaid it states why, so the reader can disagree.

**The drawings are guidance, not law.** Every view in this Contract is an approved
example of a notation, not a fixed template that every Mission must fit. Where a
particular Mission's shape is not served well by the example — an unusual dependency
pattern, a Run configuration the sample does not cover, a graph the notation renders
ambiguously — the renderer adapts the drawing to communicate that Mission clearly.
What must stay stable is the glyph vocabulary and the meaning of `▶` versus `·`,
because a reader who learns those in one view must be able to trust them in every
other. Layout, grouping, ordering, and labelling may vary to fit the case.

## Decisions taken

**Fallbacks are fingerprinted from first release.** The exploration asked whether an
advisory, unfingerprinted first version would be cheaper to trial. It would, and it
would leave open the exact hole the freeze exists to close: an operator inventing a
fallback mid-Run to escape a stop. Promoting the field into the fingerprint later is
a breaking schema move made twice. Written at plan-freeze or not at all.

**The system recommends; the owner decides.** At repair exhaustion the gate shows the
failure that consumed the budget and every recorded fallback, marks the one it
considers most likely, and asks a plain question. It never reports a match as fact.
A wrong assertion here would push an owner toward re-freezing on the wrong approach,
which is worse than the dead end it replaces. This also settles whether unmatched
fallbacks stay silent: they do not. The owner needs the options in order to decide.

**Drift is named flags, not a score.** A score defaults cleanly and cannot be argued
with. Flags let an owner disagree with the reason rather than the number, and owner
disagreement is the control that makes aimed audits safe to default.

**The authority answer is output, not a command.** An earlier draft of this Contract
introduced `authority check <mission-ref> <verb>` as an eleventh command and a new
noun. That is a direct conflict with CC-missioncli's closed ten-command surface and
with M6's stop on a growing surface, and the conflict is not worth paying: `mission
check` already computes the whole decision table inside `authority-vocabulary` and
throws it away. Folding the table into that output answers the same question, keeps the
surface frozen, and requires no Contract amendment. The underlying problem — authority
being read as two arrays to remember rather than a question to answer — is solved
either way; only the delivery changed.

**Inline and promoted must render identically, and it is tested.** CC-missioncli
already requires representation equivalence, but graph rendering is exactly where it
breaks quietly, because promoted Objectives resolve through a different decoder path
than inline ones. A byte-identical golden test across both representations is cheap
now and is the property most likely to rot silently as projection surfaces grow.

**Mission ref drift closes by decoding, not by rewriting.** M8 validates Mission order
against a list of refs, which is impossible while M2 through M4 say `human_ref:` and M5
and M6 say `ref:`. The fix is one decoder in M7 that accepts both spellings and reports
the legacy one as drift. Backfilling the three completed Missions was rejected: it
edits frozen frontmatter and changes fingerprints to finish a cosmetic rename, and
frozen records are not rewritten for tidiness.

**An undeclared authority verb is a refusal.** Defaulting an unknown verb to
`requires_owner` never permits what it shouldn't, but it answers a question the
record cannot answer and turns a typo into a confident wrong result. Refusal matches
how every other check in this system already behaves.

**The lifecycle diagram is hand-maintained, provisionally.** The exploration argued a
diagram must be generated from the command registry or not built at all, on the
grounds that a drifted diagram is worse than none because it is trusted at a glance.
That argument is sound and is overruled deliberately: no explicit state model exists
to generate from, transitions are implicit in `service.go`, and extracting one is
unscoped work that would dominate an otherwise cheap Mission. The diagram ships with
a visible hand-maintained marker and the debt is recorded as a Gap rather than left
to be rediscovered.

The exploration cited in-tree generated documentation as precedent for generating the
diagram. That precedent is `internal/command/catalog.go`, which generates the
mechanical interface *table* from the registry. No diagram in this repository is
generated; all three under `docs/diagrams/` are hand-written. Hand-maintaining this
one continues existing practice rather than departing from it.

The diagram follows the established convention: one file per diagram at
`docs/diagrams/<topic>-<aspect>.md`, a `# Title — Aspect` heading, a single fenced
`mermaid` `flowchart TD` block, and one short paragraph afterward stating the rule the
picture cannot draw. Nodes read `ID["N · Name<br/>detail · detail"]`, gates are diamond
nodes, and edge labels carry the branch condition. Owner gates are marked.

## Declined

**A cached preflight trace.** The exploration proposed persisting the last preflight's
observations with input hashes and a freshness horizon, to avoid re-deriving thirteen
dimensions every session. Declined on four grounds, recorded so it is not rediscovered
as an obvious win.

The cost removed is small — the preflight is an LLM reading `PROJECT.md` and Mission
frontmatter plus a few cheap git queries, the least expensive step in the loop — while
the mechanism is not: a persisted file, a hashing scheme, a freshness rule, and a new
trust path. It introduces a staleness failure mode where none exists today. The
genuinely volatile inputs must be re-checked every session regardless, which removes
most of the remaining value.

Its cited precedent is weaker than it appears. `.spectacular/.last-mutation` is
untracked, has no Git history, is referenced nowhere in the repository, and was last
written on 2026-08-07 under a slug matching no current Mission — before both M5 and M6
completed with full mutations. It is an abandoned v1-era mechanism that survived the v2
rewrite by being a dotfile. That is evidence the pattern was tried and dropped, not a
foundation to extend.

Session-scoped inheritance was also considered, on the reasoning that a Mission session
could scan once and its runners inherit rather than re-derive. It requires session
identity, which is not derivable today: the Go code has no session concept and reads no
environment variables anywhere. It would also optimize a fan-out that the Run model
currently forbids, since exactly one Run may be live.

The real insight underneath survives and is kept: the preflight reports all dimensions
uniformly whether or not anything moved. The fix is to report the delta, which the
state line already renders at no additional cost.

**Re-serializing canonical files into denser encodings.** `MISSION.md` is read by humans
during review, by agents during cold recovery, and by `git diff` at every commit. All
three degrade under a denser encoding and the third does so silently.

**Telegraphic syntax in frozen semantic fields.** `completion`, `stops`, `pass_boundary`,
`proof_requirement`, and `invalidated_if` are the frozen envelope. Their current
verbosity is correct rather than slack.

**Symbolic logic for authority or stops.** Compresses well, reviews badly. Owner review
is the control that makes authority meaningful, and a notation the owner must decode
before disagreeing with weakens it.

**Dropping the disjoint-claim requirement to widen fan-out.** Overlapping awareness is
not overlapping authority. Two runners owning one claim can return contradictory
evidence with no rule to arbitrate. The dependency split delivers the parallelism
without touching the guarantee.

**Autopilot deciding when to escalate.** Every direction here sits deliberately inside
the frozen envelope. Changes that blur who decides, or what counts as proof, are the
wrong trade regardless of the waste they remove.
