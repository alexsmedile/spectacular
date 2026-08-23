# Human Workspace Contract

Spectacular's canonical Markdown workspace MUST be understandable from the
filesystem before any projection or JSON is generated. UUIDv7 remains the
durable identity and SHA-256 remains the exact-revision fingerprint; neither
is the primary navigation label.

## Project anchors

The project root is intentionally small:

```text
.spectacular/
├── PROJECT.md
├── PRODUCT.md
├── ARCHITECTURE.md
├── STACK.md
├── GUARDRAILS.md
├── atlas/
├── campaigns/
├── index.md
├── catalog.md
├── workspace.yaml
├── contracts/
├── proposals/
├── missions/
├── evidence/
├── decisions/
├── gaps/
└── archive/missions/
```

## Project anchors

`PROJECT.md` is the project Anchor. `PRODUCT.md`, `ARCHITECTURE.md`, and
`STACK.md` are authoritative project Anchors for their named questions.
`GUARDRAILS.md` is owner-authored guidance selected by the runtime. `index.md`
is a deterministic, committed, non-authoritative routing guide. `catalog.md` is
the corresponding complete, non-authoritative record inventory. Neither is
canonical or loaded as authority.

`campaigns/` holds optional, durable Markdown roadmap maps. Campaigns are
planning documents, not governed records: the CLI excludes them from the typed
record graph, and they grant no execution authority.

`atlas/` holds optional, mutable Markdown maps for coherent product-value
slices. An Atlas links user journeys and outcomes to capabilities and technical
boundaries. It explains why a Campaign block exists; it does not sequence work,
grant authority, or bind a Mission.

**Anchor naming rule**: Single-word uppercase filenames (`<NOUN>.md`) are reserved exclusively for Project Anchors and workspace landmark files (`README.md`, `AGENTS.md`, `HUMAN-WORKSPACE-CONTRACT.md`). All governed records carry their scoped prefix in their filename. The same shape is used by untracked local working files such as `TODO.md` and `FEEDBACKS.md`; a record that needs to rely on one quotes what it needs, because the file is not published alongside the record.

## Mission bundles

Each Mission is a cohesive directory whose name begins with its readable
Mission reference:

```text
missions/M4-human-operability/
├── M4-human-operability.md
├── objectives/O1-design-human-layout.md
├── runs/R1-implement-layout/
│   ├── R1-implement-layout.md
│   └── checkpoints/C1-layout-approved.md  # optional advanced/historical record
├── evidence/E1-grkfsd.md
├── decisions/D1-mpzktq.md
├── gaps/G1-nrcvhw.md
├── handoffs/H1-kxsdaf.md
└── assessments/A1-qptmzr.md
```

Mission, Objective, Run, and optional Checkpoint-record names combine a scoped
ordinal with a readable slug. Evidence, Decision, Gap, Handoff, and Assessment
names combine a scoped ordinal with a stable six-character key derived from the
UUID. The key is not a content hash. Slugs may change; the frontmatter `id` does
not.

Canonical human references are scoped: `M4`, `M4/O1`, `M4/R1`, `M4/R1/C1`,
and `M4/E1-grkfsd`. CLI lookup accepts them in addition to typed UUID refs.

Supporting records live inside their Mission by default. A record is promoted
to a project-level collection only when it genuinely applies across Missions.
Archival moves the complete Mission bundle to `archive/missions/` as one
recoverable transaction.

## Interface rule

Human output leads with readable reference, outcome, state, current Objective,
current Run, latest durable checkpoint when present, blocking Gaps, and exactly
one continuation or owner gate. Ordinary checkpoint notes live in the Run body.
UUIDs, fingerprints, paths, and generation basis remain available as source
detail and in `--json`; they do not dominate the default view.

No v1 path reader, migration, alias, or compatibility branch is part of this
contract. RC.2 is a clean correction of the v2 representation.
