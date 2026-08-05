# Soft-DB projection inventory

**Status:** implemented first slice for Issue `alexsmedile/spectacular#12`:
named entity details, literal `--full` escalation, queryable summary signals,
and omission regression coverage. Generic filters, caches, derived indexes,
and broad YAML parsing remain deferred. **Upstream baseline:** merged Issue #11 selected
`spectacular status --brief --json` as the workspace-orientation command. This
request extends that route only after a collection is selected; it does not
create or choose another cold-start command.

## Contract boundary

Frontmatter is the compact, queryable signal layer. The Markdown body is
durable evidence and is never discarded, summarized by a generic heading
slicer, or batch-read to emulate a query. A caller may read a directly named
file after it selects an entity, requests a named artifact, needs an
unprojected fact, or repairs/debugs the substrate.

Every high-use collection should have three predictable layers:

1. **List** — bounded rows used only to select a record or see whether an
   immediate action exists. It must show the record's identity, state, and the
   one discriminator necessary for that selection.
2. **Detail** — a named, entity-specific projection that supports one stated
   decision. It must name its literal `--full` escalation.
3. **Full evidence** — `spectacular <entity> <ref> --full`, emitting the
   source Markdown unchanged. Direct named-file access is an equivalent escape
   hatch once `<ref>` is known.

`request --brief` remains a special active-implementation compilation; it is
not a generic detail template.

Workspace orientation is already owned by #11:

```text
spectacular status --brief --json → request <slug> → request <slug> --brief (when active)
```

#12 begins at the selected collection/record boundary, not before it.

## Current command and view inventory

| Entity | Current list | Current named view / full | Smallest safe decision now | Gap or inconsistency |
|---|---|---|---|---|
| requests | `requests [--status/--active/--since/--limit/--all/--json]` | `request <slug>` overview; active `--brief`; `--full` bundle | Select a request by lifecycle, priority, update, and summary; overview identifies current milestone/task | Strongest family. List does not expose blockers/current milestone; overview has a generic skim plus derived state. |
| decisions | `decisions [--tag/--since/--limit/--all/--json]` | `decision <ref>` generic skim; `--full` | Select an ADR by date/tag/summary | Detail cannot safely answer the decision/rationale without body access; generic JSON returns raw frontmatter + headings. |
| memories | `memories [--tag/--since/--limit/--all/--json]` | `memory <ref>` generic skim; `--full` | Select a potentially relevant standing lesson | Detail does not expose the lesson, status/retraction, or applicability; list tag matching is substring-based. |
| questions | `question list [--status]` | none; direct file only | Identify an open human blocker by priority/question | List extracts `## Question` from body, has no limit/JSON, and no named detail or literal full route. |
| research | `research list [--status]` | none; direct file only | Identify an open/completed evidence task by goal | List extracts `## Goal` from body; no result/evidence/blocked-by or named detail/full. |
| spikes | `spike list [--status]` | none; direct file only | Identify feasibility work and whether execution needs approval | Same body extraction and missing detail; list hides result and `execution_requires_approval`. |
| ideas | `idea list [--status]` | none; `idea grill` is a workbench action, not read detail | Select a candidate by status, priority, and origin | No named read/full form; identity output is a filename-like slug rather than a consistent canonical reference. |
| audits | `audit list [--status]` | none; direct file only | Choose the investigation that needs diagnosis/closure | Row's ID is the filename (`A1-title`), not canonical `A1`; list omits problem/root cause/proposed disposition. |
| fixes | `fix list [--since]` | none; direct file only | Find a verified remediation candidate by title | No signature/root cause/remedy/proof detail; declared `--since` is currently parsed but not applied; CLI requires `--verified-by` while its rules doc still calls it a soft warning gate. |

All existing plural-list commands should retain their compatibility spelling.
The proposed singular forms add a read route only; they do not replace the
existing `new`, `list`, `resolve`, or `promote` subcommands.

## Proposed entity-specific details

The exact command after a list is deliberately uniform, but the output fields
are not. Each detail ends with `Full evidence: spectacular <entity> <ref>
--full`; no generic heading-to-field mapping is introduced.

| Entity | Proposed named detail | Primary decision it safely supports | Required compact fields | Exact escalation |
|---|---|---|---|---|
| request | retain `request <slug>` overview; retain `request <slug> --brief` for active implementation | whether to resume, verify, or select the current task | lifecycle, source, current milestone/task, blockers, progress, request links | `spectacular request <slug> --full` |
| decision | `decision <ref>` | whether this ADR settles the current choice or must be reconsidered | ID/title, status, decided date, decision, context, consequences, related/tags | `spectacular decision <ref> --full` |
| memory | `memory <ref>` | whether the lesson applies now and is still active | ID/title, status/retraction marker, lesson, tags, created/updated, related | `spectacular memory <ref> --full` |
| question | `question <ref>` | whether work must stop for this human answer and what answer/options are needed | ID, status, priority, requires-user-input, question, context, options, blocked-by/target, resulting decision if resolved | `spectacular question <ref> --full` |
| research | `research <ref>` | whether the uncertainty is cleared and what evidence/result is available | ID, status, result, goal, outcome, evidence, blocked-by, related | `spectacular research <ref> --full` |
| spike | `spike <ref>` | whether a feasibility experiment is authorized/complete and what it established | ID, status, result, execution-requires-approval, goal, outcome, evidence, blocked-by | `spectacular spike <ref> --full` |
| idea | `idea <ref>` | whether the proposal merits further shaping, research, or explicit promotion | ID/title, status, priority, hypothesis, origin, open question, related/promoted-to | `spectacular idea <ref> --full` |
| audit | `audit <ref>` | whether diagnosis has enough cause/criteria to close, fix directly, or fold into a request | ID/title, status, severity, problem, intended behavior, root cause, proposed fix, disposition, related | `spectacular audit <ref> --full` |
| fix | `fix <ref>` | whether this verified remedy matches the present bug and may be reused | ID/title, verified-by/date, signature, severity, problem, root cause, fix, success criteria, from-audit/related | `spectacular fix <ref> --full` |

The list discriminator changes only where necessary: questions add
`requires_user_input` and blocking target; research/spikes add `result`;
spikes add approval state; audits add the canonical ID plus disposition; fixes
add signature. It does not require a generic `--where`, cache, or YAML parser.

## Proposed grammar and output rules

```text
spectacular <plural> [existing list options] [--limit N] [--json]
spectacular <singular> <id-or-alias> [--json]
spectacular <singular> <id-or-alias> --full
```

- Detail output is a stable, labelled projection whose fields are explicitly
  coded per entity. `--json` returns the same named fields under a
  versioned entity schema; it must not expose a raw frontmatter blob as the
  contract.
- Existing request flags stay exceptional and compatible:
  `request <slug> --brief [-mN]` is the active implementation view and
  `request <slug> --full` is its stable ordered bundle.
- New list options are limited to `--limit N` and `--json` where absent. Add
  semantic list filters only in separately approved work after the declared
  indexed-frontmatter subset exists.
- Sort explicitly by a documented stable key (canonical numeric ID for numbered
  entries; stable lifecycle/priority/slug ordering where the view is a queue),
  never by filesystem enumeration accidentally.
- Missing indexed data must render an explicit `unknown`/`not recorded` marker
  in human output and `null` in JSON. It must never silently substitute a body
  heading, except during an explicit legacy-compatibility path that is tested
  and labelled.

## Implementation map for the first slice

This is deliberately a narrow change set. It does not add a generic projection
engine: each entity handler extracts its own known slots, while shared code may
only handle flag parsing, ID resolution, deterministic row ordering, and JSON
escaping.

| Surface | Introduce | Location |
|---|---|---|
| Indexed selection signal | A scalar `summary:` for new question, research, spike, idea, audit, and fix entries. Existing decision/memory/request summaries remain compatible. Legacy records may use a labelled compatibility fallback until a later migration/doctor scope is approved. | `skills/spectacular/templates/question/entry.md`, `idea/base.md`, `audit/entry.md`, `fixes/entry.md`; the inline discovery renderer and creation handlers in `cli/spectacular` |
| Named question/research/spike detail | `question <ref>`, `research <ref>`, and `spike <ref>`, each with `--json`/`--full`. Preserve their current `new|list|resolve` subcommands. | `cmd_question`, `cmd_discovery_record`, and their entity-specific helpers in `cli/spectacular` |
| Named idea/audit/fix detail | `idea <ref>`, `audit <ref>`, and `fix <ref>`, each with `--json`/`--full`. Preserve `idea promote` and audit/fix mutators. | `cmd_idea`, `cmd_audit`, `cmd_fix`, plus named entity helpers in `cli/spectacular` |
| Decision/memory replacement | Replace `_skim_file` as the default for these two commands with explicit decision and memory field renderers. Keep the generic helper only where its outline behavior is still explicitly wanted. | `cmd_decision` and `cmd_memory` in `cli/spectacular` |
| Selection rows | Read `summary:` and the small per-entity discriminator only; remove ordinary body-heading extraction from question/research/spike list paths. | `cmd_question_list`, `_discovery_list`, `cmd_idea_list`, `cmd_audit_list`, `cmd_fix_list` |
| Regression suite | One shared, end-to-end compact-view suite plus existing focused command tests for mutation compatibility. | New `tests/cli/read-projections.test.sh`; extend `decide.test.sh`, `idea.test.sh`, `audit-fix.test.sh`, `wayfinding-sequencer.test.sh`, and `request-workflow.test.sh` only where their existing fixtures own the behavior |
| Agent guidance | List → named detail → `--full`; direct named-file reads remain legal after selection and for repair/debugging. Do not change the #11 orientation route. | `skills/spectacular/SKILL.md`, `references/soft-db-index.md`, `docs/commands.md` |

The detail renderer may extract named, entity-owned body sections such as a
fix's `## Root cause`; this is a stable record contract, not a generic heading
slicer. The list renderer must rely on frontmatter so selection never needs to
load arbitrary bodies.

### Implementation checkpoint

The first slice now implements named `--json`/`--full` details for questions,
research, spikes, ideas, audits, fixes, decisions, and memories. It replaces
the decision/memory generic skim with explicit fields, keeps the established
request overview/brief intact, and adds `tests/cli/read-projections.test.sh`.
The re-run of `tests/benchmarks/retrieval-baseline.sh` retained the #11
`status --brief` orientation at 731 bytes with zero full-body reads; it offers
no evidence to add a compiled workspace briefing, generic filter, index, or
cache. Those remain future, separately scoped proposals.

## Decision-risk analysis

| Risk | Failure if omitted or generalized | Required mitigation |
|---|---|---|
| False progress | A request list lacks blockers/current milestone, so an agent starts the wrong work | Keep request overview/brief specialized; test blocker and active-task omission. |
| Re-litigated choice | A decision detail exposes headings but not choice/consequences | Project decision/context/consequences by name; retain full evidence. |
| Unsafe continuation | A question or spike hides human-input/approval state | Make `requires_user_input` and `execution_requires_approval` mandatory fields. |
| Evidence treated as conclusion | Research shows a goal but omits result/evidence | Project both `result` and evidence/outcome; `inconclusive` must remain visible. |
| Misapplied bug repair | A fix list/detail hides signature, root cause, or verification | Require all three in the fix detail; surface unverified/malformed records explicitly. |
| Parser overreach | A shared heading slicer or broad YAML parser produces a plausible but wrong field | Write entity extractors against known entry contracts; reject unsupported indexed values in later dedicated work. |
| Context regression | A collection view duplicates or contradicts the merged workspace orientation | Keep `status --brief --json` as the sole orientation view; collection detail begins only after entity selection. |
| Compatibility drift | New read grammar breaks aliases/mutator namespaces or unstable output | Retain all old spellings and assert text/JSON ordering, aliases, and error paths. |

## Test inventory

Tests belong in the existing CLI suites, with small purpose-built fixtures that
put a decision-critical value in each required slot. Do not use a body-heading
fixture as a substitute for indexed signal.

| Coverage | Assertions required before shipping |
|---|---|
| Requests | List selects the correct lifecycle item; overview includes source, current milestone/task, blocker, and progress; active `--brief` remains unchanged; `--full` keeps ordered raw files. |
| Decisions | List selection; detail includes decision/context/consequences; tag/date fields survive JSON; absent consequences is explicit; full record is byte-for-byte body evidence. |
| Memories | Active and retracted records are distinguishable; detail exposes lesson/tags/related; legacy lesson fallback is explicit; full remains literal. |
| Questions | Open/deferred/resolved-history handling; priority + user-input + blocked-by fields; options/context appear in detail; no body scan is needed for the list row. |
| Research and spikes | `supported`, `refuted`, and `inconclusive` remain distinct; evidence/outcome survive detail; spike approval flag cannot be omitted; aliases resolve; completion evidence remains literal in full. |
| Ideas | Canonical ID and descriptive slug select the same record; status/priority/hypothesis/origin/open question/promoted target are projected; `idea new/list/promote` compatibility holds. |
| Audits and fixes | Canonical ID is separate from title slug; audit cause/proposed fix/disposition present; fix signature/cause/fix/verified-by present; unverified/malformed data is explicit; `fix --since` either works or is removed/documented in the implementation slice, never silently ignored. |
| Cross-cutting | Human rows and JSON have documented deterministic order; `--limit`, unknown alias, missing file, malformed required frontmatter, legacy record, and `--full` are covered; `bash -n cli/spectacular` passes under Bash 3.2-compatible syntax. |
| Baseline compatibility | `status --brief --json` remains the orientation route; collection tests compare list/detail/full reads against direct named-body reads without introducing another workspace briefing command. |

### First-suite fixture matrix

`tests/cli/read-projections.test.sh` should build one minimal workspace and
write deliberately different values into every field below. Assertions must
look for values, not merely headings, so deleting a projection field fails the
test.

| Fixture | Compact detail must expose | List must expose | Negative assertion |
|---|---|---|---|
| `QUE-001-release-owner` | `requires_user_input: true`, priority, question, options, `blocked_by` | ID, open/deferred state, priority, blocker signal, summary | A normal list must not print the long Context body. |
| `RES-001-provider-limit` | result, goal, outcome, evidence, dependency | ID, state, result, summary | `inconclusive` must not render as cleared/supported. |
| `SPK-001-parser-prototype` | result plus `execution_requires_approval: true`, outcome, evidence | ID, state, result, approval signal, summary | A list/detail cannot imply execution is authorized when the flag is true. |
| `IDEA-001-quick-import` | hypothesis, open question, status, priority, origin, promoted target | ID, status, priority, summary/origin | `idea promote` remains a mutator, not the default read path. |
| `A1-config-regression` | problem, root cause, proposed fix, disposition, severity | canonical `A1`, state, severity, disposition, summary | The displayed identifier must not be only the filename/title slug. |
| `F1-config-regression` | signature, root cause, applied fix, verified-by/date, criteria | canonical `F1`, verification state, signature, summary | A null/missing verification or signature must be explicit, never rendered as verified. |
| `DEC-001-storage-choice` | decision, context, consequences, tags/status | ID, date/status, summary/tags | The default detail must not collapse to a heading outline. |
| `M1-cli-convention` | lesson, active/retracted state, tags, related | ID, state, summary/tags | A retracted memory cannot render as an active instruction. |

For every fixture, assert all of the following:

```text
<entity> list        → bounded selection row only
<entity> <alias>     → entity-specific detail and literal full escalation hint
<entity> <alias> --json → stable, named field values; no raw-frontmatter contract
<entity> <alias> --full → original Markdown content, including a long evidence-only paragraph
```

Add error cases for an unknown alias, absent required frontmatter field, and an
unsupported flag. The error must name the entity and corrective command; it
must not silently fall back to a body scan.

## Staged implementation plan (begins after #11)

1. **Preserve the established boundary.** Keep #11's `status --brief --json`
   orientation route unchanged; begin #12 only once a collection/record needs
   inspection.
2. **Build entity-specific details.** Add named detail and literal full routes
   for the chosen missing collections, replacing decision/memory generic skims
   only when their compatibility tests pass. Do not add generic filters.
3. **Harden list projections.** Add only the list discriminators needed for
   safe selection plus deterministic sort/limit/JSON behavior.
4. **Prove omissions cannot regress.** Land the above fixtures and error/legacy
   tests; update CLI-first guidance only once its target commands exist.
5. **Measure before expanding.** Re-run the collection retrieval comparison. Consider
   compiled workspace briefing, constrained equality filters, frontmatter
   validation, derived indexes, or caching as separate proposals, not automatic
   follow-ons.
