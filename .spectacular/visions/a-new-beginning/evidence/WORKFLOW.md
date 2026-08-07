---
type: refactor-intake-workflow
status: active-prototype
owner: alex
updated: 2026-08-07
---

# Refactor intake workflow

## Purpose

Turn heterogeneous plans and expert suggestions into small, filterable concept
pieces without letting source order, reputation, or document size decide the
product. This is a Vision-owned prototype. It may inform a reusable refactor
workflow skill later, but it is not yet a Spectacular product contract.

## Stores

| Store | Owns | Does not own |
|---|---|---|
| Source card | Provenance, authority, argument summary, source-wide assumptions | Final product choices |
| Concept piece | One reusable proposal, principle, warning, or measurement idea | The source's full narrative |
| Concept index | Filterable current view of pieces | Detailed reasoning |
| Contradiction matrix | Conflicts, dependencies, duplicate-authority risks, and drift | Resolution by implication |
| Vision fragment | A coherent direction presented for explicit human reaction | Raw intake |
| Decision record | A confirmed choice between real alternatives | Pending recommendations |
| Specification | The converged behavior or contract to implement | Rejected or unresolved pieces |

## Concept frontmatter

Every concept piece uses the same fields:

| Field | Meaning |
|---|---|
| `id` | Vision-local stable identifier; not a canonical Spectacular entity ID |
| `status` | `captured | compared | decided | promoted | parked` |
| `domain` | Primary comparison domain |
| `sources` | Source cards that independently contain or support the concept |
| `source_authority` | Authority of the current source evidence |
| `assessment` | Intake quality signal: `strong | promising | mixed | disputed | weak | unknown` |
| `evidence_status` | `unverified | partial | supported | refuted | not-needed` |
| `disposition` | Human choice: `pending | adopt | adapt | reject | defer | needs-evidence` |
| `depends_on` | Other pieces that must be settled first |
| `overlaps_with` | Related but not equivalent pieces |
| `conflicts_with` | Pieces that cannot both be accepted as written |
| `tags` | Secondary filter vocabulary |

The body owns the core message, value, assumptions, evidence, collisions,
trade-offs, recommendation, and decision state. Recommendation is advisory;
`disposition` stays `pending` until the owner chooses.

## Pass A — ingest each source

1. Register provenance, authority, date/context, and raw location.
2. Summarize the thesis and proposals without copying the whole source.
3. Verify cheap repository claims; mark expensive or external claims unverified.
4. Identify assumptions, contradictions, drift, duplication, and missing evidence.
5. Create or update one atomic concept piece per reusable idea.
6. Update source and concept indexes plus the contradiction matrix.

No specs, decisions, or implementation are produced during ordinary intake.

## Pass B — rolling comparison

After each source, link obvious duplicates, overlaps, dependencies, and conflicts.
After roughly three to five substantial sources—or earlier if a blocking conflict
appears—run a synthesis checkpoint:

- merge duplicate concepts by adding source provenance rather than copying files;
- split pieces that contain more than one independently decidable idea;
- update domain clusters and evidence strength;
- distinguish disagreement about facts from disagreement about values;
- identify the smallest set of decisions that resolves the most downstream fog.

## Pass C — decision packets

When intake is declared complete, or one domain is mature enough to decide,
produce decision packets containing:

1. the exact question;
2. viable options, including keep/remove where applicable;
3. supporting and contradicting pieces;
4. repository and external evidence;
5. user value, complexity, compatibility, migration, and maintenance trade-offs;
6. one clear recommendation;
7. the explicit human disposition.

Route research only when a fact could change the choice. Propose a spike only when
observed feasibility could change it. Prestige is provenance, not proof.

For this refactor, mature packets are grouped into the dependency-ordered program in
[`decision-sessions.md`](decision-sessions.md). Each session ends with an explicit contract,
downstream constraints, and human dispositions. The ranked architecture frontier is maintained in
[`top-20-foundational-decisions.md`](top-20-foundational-decisions.md). Neither artifact is a
substitute for recording the decisions in the underlying concept cards and eventual decision log.

## Pass D — converge and promote

Repack accepted pieces into the fewest coherent Vision fragments. Do not promote
one fragment per source or per concept. Resolve fragment reactions, approve the
whole Vision explicitly, then derive small specifications around real behavior or
contract boundaries. Approved specs seed ordered requests; requests seed code.

## Pass E — session retrospective

At the end of this refactor, review the workflow itself:

- Which fields were actually useful for retrieval and decisions?
- Which fields or files created busywork or duplicate truth?
- Did atomic pieces improve contradiction detection and recomposition?
- Where did IDs, indexes, or manual synchronization drift?
- Which steps require tooling, and which should remain agent judgment?
- Is the proven workflow general enough to become a standalone skill?

Only then design the reusable skill from observed behavior rather than this
first-pass schema.

## Quality gates

- One independently decidable idea per concept file.
- No concept without provenance and a core message.
- No factual claim presented as settled without evidence status.
- No `adopt`, `adapt`, or `reject` disposition inferred from an agent assessment.
- No new authority that duplicates an existing registry, index, or canonical doc.
- No accepted piece promoted until conflicts and dependencies are visible.
- No implementation task without an approved specification and validation path.
