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
    ├── VOCABULARY.md                  # domain objects and language
    ├── INDEX.md                       # optional manual navigation
    ├── index.json                     # generated governed inventory
    ├── raw/
    │   └── release-sketch.md           # freeform starting material
    ├── plans/
    │   ├── release-flow.md             # parent approach, progress, results
    │   └── linux-tranche.md            # independent linked delivery slice
    ├── specs/
    │   └── release-integrity.md        # reusable intended behavior; type Spec
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
    └── archive/                       # authorized governed retirement
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
| Spec | What behavior or interface should exist? | Requirements have independent readers or reuse |
| Atlas | What are the objects and how do they relate? | An explanation outlives a work plan |
| Decision | What did we choose and why? | Rationale or supersession needs retention |
| Audit | What did inspection find at this revision? | Findings have an independent scope and reader |
| Proposal | Which formally tracked question remains open? | Existing governed exploration is deliberately selected |
| Contract | What agreement is mechanically bound? | Accepted agreement needs validation and versioning |
| Mission | Which frozen slice is explicitly executing? | Owner selects governed work |
| Review / Handoff | What was formally evaluated or transferred? | The selected governed workflow needs that record |

A plan can contain requirements, inspection findings, and a continuity note.
Optional specs/ and audits/ avoid overloading it only when those parts stand alone.
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
current types and metadata for anchors, Atlas, plans, specs, audits, Decisions, campaigns, and
retrospectives. `schema` means an enforced governed frontmatter claim; it is not
needed for soft context. Existing governed records retain their own metadata,
identity, and mutation rules. Raw has no metadata or naming obligations.

## Navigation and manageable files

Use one coherent question, object, or outcome per file. Prefer roughly 50–120
lines, splitting by independent meaning, lifetime, domain, module, phase, or
tranche when useful. Reuse links instead of copied context. Avoid mandatory
indexes, logs, and empty folders.

Relative Markdown links resolve from the source file. Path-qualified wikilinks
resolve from `.spectacular/`; use bare filenames only when unambiguous. Aliases
and section links help readers. Agents repair incoming links on moves.
`VOCABULARY.md` is preferred for new domain anchors; retain an existing
`ONTOLOGY.md` when it owns that meaning instead of duplicating authority.

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
Skill and require Python 3 with PyYAML. Findings guide cleanup; they do not enroll
work in a Mission. The checker does not certify section anchors or governed
bindings. Root `docs/` remains public product documentation, never a plan store.
