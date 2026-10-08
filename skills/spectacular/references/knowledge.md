# Durable context in small linked objects

Use this when: maintaining or placing ordinary workspace knowledge, or starting
authorized work from a raw draft or plan without opting into a Mission.

## Work from the relevant context

Read the relevant anchors and nearby documents, follow useful links, then act.
Follow the task-sized reading route in [navigation.md](navigation.md): optional
INDEX → PROJECT → relevant Requirement → Spec → Plan, consulting Ontology, Atlas,
and Decisions when useful. Load only relevant linked context. Finish when the
requested outcome is verified and relevant durable knowledge is updated.
No CLI, manifest, governed lifecycle, or formal record is needed for this route.
Spectacular maintains durable context as navigable OKF objects and relationships.
Folder agreements live in [knowledge-folders.yaml](../knowledge-folders.yaml).

## Placement

| Home | Purpose |
|---|---|
| Root anchors | Current purpose, accepted architecture, vocabulary, constraints |
| raw/ | Unconstrained captures, sketches, scratchpad, quick-start material |
| atlas/ | Linked entities, concepts, relationships, lifecycles, and maps |
| decisions/ | Consequential choices and rationale |
| plans/ | Intended outcome, approach, progress, and results |
| requirements/ | Optional linked needs, constraints, and acceptance criteria |
| specs/ | Optional reusable behavior and interface specifications |
| audits/ | Optional revision-scoped inspection findings |
| contracts/ | Mechanically validated agreements and amendment/version pipeline |

Implementation plans belong in plans/, exploratory sketches in raw/. Product
documentation belongs in root docs/, managed by documentation tools such as
pageworks; Spectacular does not manage its structure or use it as a plan store.
Raw, sketch, and scratchpad name the same role across tools. Use the established
one rather than creating three copies; all are outside knowledge obligations.

ONTOLOGY.md is the preferred domain-model anchor name. Read an existing
VOCABULARY.md when it owns that model; do not create a duplicate authority or
rename historical bindings casually. A manual INDEX.md and generated index.json
are both valid navigation aids. Keep authored content out of generated indexes.

Update the existing owner of a topic before creating another file. Create a
folder or anchor only when content needs it. Requirements and design may live
in a plan; extract lasting knowledge only when it has an independent use.
For a distributed PRD, use [requirements.md](requirements.md). Optional drafting
aids live in [drafting-methods.md](drafting-methods.md); load them only when useful.
For planning requests use [plan.md](plan.md); for the full container and Mission
bundle layout use [workspace.md](workspace.md).
Mechanical Decision validation is recommended, not required. Descriptive names
need no governance flag. A numbered soft Decision may explicitly declare
governance: context to resolve naming ambiguity; governed identity/schema claims
still retain strict validation. Existing governed Decisions retain their rules.
Contracts always retain their mechanical validation and amendment/version process;
ordinary context never replaces an accepted Contract.

Keep a document stable as it matures. Move only when purpose changes, preserving
content and metadata and repairing incoming links. Never overwrite an occupied
destination or delete duplicates by inference. Report meaningful moves briefly.
Keep proposed ideas distinct from accepted constraints. Extract consequential
choices and accepted changes to project truth; explicitly supersede changed
Decisions rather than erasing their rationale. Ask about conflicting authority
or consequential ambiguity; choose routine placement autonomously.

## Small files and links

Aim for one retrievable question or outcome, usually 50–120 lines. This is a
judgment guide, not a limit. Split independent meaning, readers, or lifetime;
split by domain, code module, working phase, or delivered tranche when each part
stands alone; avoid fragments that need many other files merely to make sense. Keep a plan
short and link substantial reference material. No mandatory index or log.

Read relative Markdown links and Obsidian wikilinks, including aliases and
section targets. Markdown relative paths start at the source document; wikilink
paths start at .spectacular. Resolve explicit paths first; use a basename only
when unique. Report ambiguous or missing targets rather than inventing them.
Repair both link forms when moving a soft file. Use portable Markdown links in
public docs. These links navigate knowledge, not governed typed references;
no mechanical wikilink resolver is currently provided.

## Metadata profiles

Use [profiles.md](profiles.md) when choosing metadata or scaling a document.
Compact is the default for maintained soft context; minimal remains valid.
All soft profiles require type, document version, created, and updated. A short
description is recommended; title and extended retrieval fields remain optional.
The freeform raw/sketch/scratchpad role has no metadata obligations.

Governed schemas are independent of this scale. Existing records retain enforced
identity, metadata, mutation, freeze, and archival rules. Contracts keep their
mechanical validation and amendment/version pipeline. Ordinary edits do not
rewrite frozen Missions or amend bound Contracts by hand.

Reading or linking raw does not authorize publishing it. Keep ignored material
ignored unless retention is explicitly changed; governed proof preserves what it
relies on. Product checks follow the actual change and repository rules.

## Optional metadata and link diagnostics

Run `python3 scripts/check-knowledge.py <workspace>`, where workspace is the
.spectacular directory, from the skill directory when checking context quality.
Python 3.9+ and PyYAML are required. Exit 0 means no findings, 1 means advisory
findings, and 2 means unavailable inspection. The checker is read-only, skips raw
and governed records, checks folder metadata and file link targets, and reports
JSON. It does not certify section anchors, identities, or Contract bindings.

For an optional size estimate, run the resolved Skill's
`scripts/count-tokens.sh <file-path>` from the target project. It reports lines,
words, and an advisory heuristic estimate; it is not the official tokenizer,
never enforces a soft-file limit, and does not certify a Charter budget.
