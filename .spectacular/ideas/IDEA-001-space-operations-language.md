---
id: IDEA-001
type: idea
status: parked
priority: medium
owner: alex
origin: Terminology exploration with Alex, 2026-08-04
updated: 2026-08-04
promoted_to: null
related: [../PRD.md, ../roadmaps/index.md, ../specs/index.md]
---

# IDEA-001 — space-operations-language

## Hypothesis

Give Spectacular a focused mission-control vocabulary for its distinctive operational experiences while retaining technical artifact names.

## Context

Spectacular should keep its product name. The space theme becomes a personal,
memorable layer for its distinctive operational experiences, rather than a
replacement vocabulary for every stored artifact or command.

The canonical technical layer remains explicit and searchable: `PRD`,
`ROADMAP`, `SPEC`, request, plan, task, milestone, verification, and archive.
This preserves interoperability, documentation clarity, and the meaning of
on-disk contracts.

**Language architecture:** first complete the technical layer: precise entity
names, CLI/API contracts, paths, schemas, and documentation must stand on their
own with no metaphor required. Then add the mystical/abstract storytelling
layer on top in headings, UI copy, onboarding, visualizations, and only
carefully chosen non-breaking aliases. The storytelling layer may enrich the
journey; it must never obscure what a command writes or what an artifact is.

The proposed mission-control layer:

| Technical surface | Spectacular language | Intended boundary |
|---|---|---|
| `status`, `summary`, `next` | **Navigator** | The orientation experience; not a second workspace entity. |
| `wayfind` | **Wayfinding** | Signature dependency-aware navigation feature. |
| unresolved/blocked graph state | **Fog** | What is not yet known or unblocked. |
| dependency-ready graph state | **Frontier** | The edge of known, actionable work. |
| request | **Mission** | A user-facing rendering of a request, not a replacement schema or CLI namespace. |
| `depends-on` / `blocks` | **Tether** | Graph/UI prose only; technical relation fields remain unchanged. |
| `afk run` | **Autopilot** | Human-facing name for bounded autonomous work with explicit gates. |
| `/spectacular act <SPC>` | **Launch** | Candidate alias only if it preserves every existing authorization gate. |
| `snapshot` | **Checkpoint** | User-facing phrase for creating a recoverable versioned copy. |
| `sweep` | **Recon Sweep** | Read-only evidence audit across the request fleet, using the agent fleet. |
| `audit` + `fix` | **Repair** | Branded bug-workflow arc; investigation and fix remain the technical records. |
| `doctor` | **Systems Check** | Product explanation for the existing integrity checker. |
| `verify` | **Telemetry Check** | Product explanation for evidence-based verification. |
| `archive` | **Mission Log** | Completion language; archive remains the technical operation. |

`Fleet` and `Sweep` are deliberately distinct: a **request fleet** is the
set of active/planned work; an **agent fleet** is the set of focused
subagents; a **Recon Sweep** is the audit that crosses the request fleet using
the agent fleet.

`SPC` remains a durable canonical ID prefix; prose and UI should normally say
**spec** or **specification**. The prior umbrella term “soft-DB” needs a
separate vocabulary review; the technically accurate candidate is
**frontmatter-indexed Markdown collections**.

### Mission: the user-facing identity of a request

**Mission** is the strongest candidate for a first-class branded concept. It
must mean exactly one thing: the user-facing identity of a durable
implementation **request**. It does not describe a spec, research record,
spike, AFK run, release, or the whole project.

Keep the technical contract unchanged:

- on disk: `requests/<slug>/`, `PLAN.md`, `TASKS.md`
- canonical CLI/API: `request` and `spectacular.request.v1`
- historical links, generated workspaces, and integrations: unchanged

The product layer can then say **Mission Briefing** (`request <slug> --brief`),
**Flight Plan** (`PLAN.md`), **Waypoints** (milestones), **Mission Fleet**
(active/planned requests), **Launch Mission** (`/spectacular act <SPC>`), and
**Log the Mission** (`archive`). A later `mission` CLI alias may be worthwhile
only if it is an exact non-breaking alias for `request`; a future `launch`
alias must run the same authorization gates as `act`.

## Open questions

- Should **Navigator** remain output/onboarding language or earn a dedicated
  command/interactive surface?
- Does **Checkpoint** apply only to `snapshot`, or can a milestone be called a
  checkpoint in UI copy without creating ambiguity?
- Is `launch` a worthwhile non-breaking alias for `/spectacular act`, or is it
  better as conversational/output language only?
- After the Mission language is proven in UI and docs, should `mission` become
  an exact CLI alias for `request`, including a plural `missions` alias for
  `requests`?
- Which of the proposed terms—if any—are strong enough to become branded
  feature headings in README and product documentation?
- When the collection vocabulary is reviewed, is “frontmatter-indexed Markdown
  collections” clear enough, or should individual collection names carry more
  of the meaning (decision ledger, session log, repair log)?

## Promoted to

—
