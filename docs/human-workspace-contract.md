# Workspace layout

Create folders when content needs them. Spectacular keeps small typed Markdown
objects in containers; the minimal workspace can be a single project anchor.

## Example: a project actively working on releases

```text
project/
├── docs/                              # human-facing product docs
│   ├── guides/release.md
│   └── reference/artifacts.md
└── .spectacular/
    ├── PROJECT.md                     # purpose and current constraints
    ├── ONTOLOGY.md                  # domain objects and language
    ├── INDEX.md                       # optional manual navigation
    ├── index.json                     # generated governed inventory
    ├── raw/
    │   └── release-sketch.md           # freeform starting material
    ├── plans/
    │   ├── release-flow.md             # parent approach, progress, results
    │   └── linux-tranche.md            # independent linked delivery slice
    ├── requirements/
    │   ├── release-needs.md            # short linked overview; type Reference
    │   └── verified-download.md        # need and acceptance; type Requirement
    ├── specs/
    │   └── release-integrity.md        # reusable intended behavior; type Spec
    ├── scenarios/
    │   └── verify-download.md          # reusable journey; type Scenario
    ├── atlas/
    │   ├── release-artifact.md         # typed domain object
    │   └── release-lifecycle.md        # relationships and states
    ├── decisions/
    │   └── checksum-policy.md          # soft choice and rationale
    ├── audits/
    │   └── release-drift.md            # inspection scope/revision/findings
    ├── contracts/
    │   └── CC-r7u2p4-release-integrity.md     # mechanically bound agreement
    ├── proposals/
    │   └── P7-signing-options.md       # formally tracked open question
    ├── missions/
    │   └── M23-release-integrity/
    │       ├── M23-release-integrity.md  # explicitly selected frozen slice
    │       ├── reviews/RV1-release-review.md
    │       ├── evidence/E1-x4b8q2.md
    │       ├── handoffs/H1-k4p2a8.md
    │       ├── objectives/O2-linux-delivery.md  # only if promoted
    │       └── runs/R2-linux-attempt/R2-linux-attempt.md  # only if promoted
    ├── campaigns/release-readiness.md # optional sequencing across slices
    ├── retrospectives/release-lessons.md
    └── archive/                       # retired context and governed bundles
```

The names illustrate placement, not records to manufacture. The CLI assigns
exact governed filenames and reports their paths. Compact Missions keep inline
Objectives and Runs until supported promotion earns separate files. Mission-owned
reviews, evidence, handoffs, Gaps, and assessments normally live inside the bundle;
project-level collections are valid when their enforced format permits them.

Create only the containers whose content is useful. A single project anchor and
one plan can be enough. Raw aliases sketch/ and scratchpad/ have the same role;
retain the existing one rather than adding all three. The current initializer
still creates the broader governed layout.

## Choose a role

| File role | Question it answers | When to extract it |
|---|---|---|
| Plan | What will we do, and what happened? | A coherent work outcome needs continuity |
| Requirement | What is needed and how is acceptance observed? | A need has independent meaning, reuse, or delivery scope |
| Spec | What behavior or interface should exist? | Behavior or interfaces have independent readers or reuse |
| Scenario | Can a user reach a concrete final outcome? | A journey needs independent reuse and observable criteria |
| Atlas | What are the objects and how do they relate? | An explanation outlives a work plan |
| Decision | What did we choose and why? | Rationale or supersession needs retention |
| Audit | What did inspection find at this revision? | Findings have an independent scope and reader |
| Proposal | Which formally tracked question remains open? | Existing governed exploration is deliberately selected |
| Contract | What agreement is mechanically bound? | Accepted agreement needs validation and versioning |
| Mission | Which frozen slice is explicitly executing? | Owner selects governed work |
| Review / Handoff | What was formally evaluated or transferred? | The selected governed workflow needs that record |

A plan can contain requirements, inspection findings, and a continuity note.
Optional requirements/, specs/, scenarios/, and audits/ avoid overloading it only when those
parts stand alone. Existing Requirement files in specs/ remain valid; choose one
content owner and link it without duplication or mandatory migration. A short PRD
overview can link needs, domain objects, specs, decisions, and delivery plans.
See [linked requirements](../skills/spectacular/references/requirements.md).

Scenarios retain goals and expected outcomes across attempts. A campaign plan
selects cases by frequency, value, and risk; an Audit records software versions,
case revision, observed outputs, evidence, and limitations for each attempt.
Use passed, failed, blocked, or not verified verdicts; documented feasibility
alone does not establish a field-test pass. Keep reusable cases after a campaign
and retire obsolete cases to archive/scenarios/ with repaired links.

Skeletons provide optional starting outlines; anatomies explain artifact parts
and relationships. They add no lifecycle or validation gate. Draft directly in
the intended file or use raw; avoid a redundant skeleton beside its finished
artifact. Project-specific reusable methods can live in Atlas when useful. See
[drafting methods](../skills/spectacular/references/drafting-methods.md).
A quick handoff can remain in the plan; governed handoffs use their supported
command. Existing Proposals keep their rules; proposal create remains unavailable.

Root docs/ is managed independently by documentation tools. There is no default
.spectacular/docs/ copy. A spec is internal intended behavior; public documentation
explains shipped behavior. Specs do not replace the Contract amendment pipeline.

## Metadata and folder agreements

Outside raw, maintained soft documents declare foundational `type`, document
`version`, immutable `created`, and meaningful-edit `updated`. Dates use ISO
timestamps with timezone; recover historical dates from reliable history and
leave unknown values unresolved rather than fabricating them. Preserve unknown
metadata. Types are extensible through the folder agreement.

The [folder agreement](../skills/spectacular/knowledge-folders.yaml) defines the
current types and metadata for anchors, Atlas, requirements, plans, specs, audits, Decisions, campaigns, and
retrospectives. `schema` means an enforced governed frontmatter claim; it is not
needed for soft context. Existing governed records retain their own metadata,
identity, and mutation rules. Raw has no metadata or naming obligations.

## Metadata scale

Compact is the default for maintained context. These are document conventions,
not a progression toward mandatory governance.

| Profile | Recommended use | Metadata |
|---|---|---|
| Freeform | Raw, sketch, scratchpad | None required |
| Minimal | Small maintained context | type, version, created, updated |
| Compact (default) | Most requirements, plans, specs, Atlas pages, soft Decisions | Minimal plus recommended description; title when useful |
| Extended | Context needing extra retrieval or attribution | Compact plus relevant optional tags, sources, status, or domain fields |
| Governed | Explicitly governed records, including Contracts | Exact command-generated schema and validator |

Minimal remains valid. Optional fields stay optional, and raw can be detailed
without metadata. No profile label is required in the file. Existing documents
need no bulk conversion; preserve their metadata and choose additions by usefulness.
Extended does not mean a larger file, and governed is independent of the scale.
See the [profile guidance](../skills/spectacular/references/profiles.md) for versions,
dates, evolution, and diagnostic limits.

## Navigation and manageable files

Use one coherent question, object, or outcome per file. Prefer roughly 50–120
lines, splitting by independent meaning, lifetime, domain, module, phase, or
tranche when useful. Reuse links instead of copied context. Avoid mandatory
indexes, logs, and empty folders.

Relative Markdown links resolve from the source file. Path-qualified wikilinks
resolve from `.spectacular/`; use bare filenames only when unambiguous. Aliases
and section links help readers. Agents repair incoming links on moves.
`ONTOLOGY.md` is preferred for new domain anchors; retain an existing
`VOCABULARY.md` when it owns that meaning instead of duplicating authority.

Manual `INDEX.md` and generated `index.json` can coexist. Generated indexes cover
governed records only and never overwrite authored Markdown. Legacy navigation
files remain readable; there is no automatic migration.

## Where the boundaries break

Soft files carrying governed identities or schema claims are checked strictly.
Malformed governed records still refuse; ordinary documents must not impersonate
governed records. Never place governed records in ignored raw material.

Moving soft content requires link repair. Moving a governed record follows its
supported mutation and archival rules. Contract edits require the amendment or
version pipeline. Completed Mission bindings are preserved.

Optional metadata and file-link diagnostics use `check-knowledge.py` from the
Skill and require Python 3.9+ with PyYAML. Findings guide cleanup; they do not enroll
work in a Mission. The checker does not certify section anchors or governed
bindings. Root `docs/` remains public product documentation, never a plan store.

## Configuration and generated navigation

Ordinary Markdown work needs no CLI or configuration. The current governed CLI
requires workspace.yaml to locate records and the project Anchor. Optional
config.yaml supplies supported overrides; omit it when built-in defaults suffice.
YAML keeps those settings human-editable, consistent with record frontmatter.
JSON serves generated data instead.

Manual INDEX.md owns reading routes through current work, reusable knowledge, and
history. Generated index.json inventories governed records only. Optional collection
JSON files are filtered views of that same graph; none provides separate authority.
Do not generate Markdown indexes or a second catalog.json cache.

In this source repository, run `go run ./scripts/rebuild-workspace-index.go` to print
the governed root inventory, or append --write to refresh it. This requires Go and
the source checkout, adds no public CLI command, and preserves manual navigation.

Read by task: optional INDEX → PROJECT → relevant Requirement → Spec → Plan,
consulting Ontology, Atlas, Decisions, and guardrails as useful. Retire completed
work under archive/, preserving governed provenance and bindings; completed plans
and campaign inputs need not occupy live containers. Legacy imports can remain
byte-preserved under archive/raw/ with no current authority. See the
[naming and retirement guide](../skills/spectacular/references/navigation.md).
