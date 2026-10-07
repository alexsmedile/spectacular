# Workspace layout

Create folders when content needs them. Spectacular keeps small typed Markdown
objects in containers; the minimal workspace can be a single project anchor.

```text
.spectacular/
├── PROJECT.md
├── VOCABULARY.md            # optional domain vocabulary and relationships
├── INDEX.md                 # optional manual navigation
├── index.json               # generated governed-record navigation, when used
├── raw/                     # unconstrained captures and sketches
├── atlas/                   # lasting entities, concepts, maps, explanations
├── decisions/               # choices and rationale
├── plans/                   # intended work, progress, results
├── contracts/               # validated agreements
├── proposals/               # optional governed open questions
├── missions/                # explicitly selected governed work
├── evidence/                # governed proof
└── archive/                 # governed retirement under existing rules
```

This is an example, not a required scaffold. Governed workflows may also use
Gaps, handoffs, assessments, and reviews. The current initializer still creates
the broader governed layout. Raw aliases `sketch/` and `scratchpad/` have the same
freeform role; retain the existing one rather than adding all three.

## Metadata and folder agreements

Outside raw, maintained soft documents declare foundational `type`, document
`version`, immutable `created`, and meaningful-edit `updated`. Dates use ISO
timestamps with timezone; recover historical dates from reliable history and
leave unknown values unresolved rather than fabricating them. Preserve unknown
metadata. Types are extensible through the folder agreement.

The [folder agreement](../skills/spectacular/knowledge-folders.yaml) defines the
current types and metadata for anchors, Atlas, plans, Decisions, campaigns, and
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
