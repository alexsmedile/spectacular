# Durable context in small linked objects

Use this when: maintaining or placing ordinary workspace knowledge, or starting
authorized work from a raw draft or plan without opting into a Mission.

## Work from the relevant context

Read the relevant anchors and nearby documents, follow useful links, then act.
An index may help navigation; load only the portion needed. Finish when the
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
| contracts/ | Mechanically validated agreements and amendment/version pipeline |

Implementation plans belong in plans/, exploratory sketches in raw/. Product
documentation belongs in root docs/, managed by documentation tools such as
pageworks; Spectacular does not manage its structure or use it as a plan store.
Raw, sketch, and scratchpad name the same role across tools. Use the established
one rather than creating three copies; all are outside knowledge obligations.

VOCABULARY.md is the preferred domain-model anchor name. Read an existing
ONTOLOGY.md when it owns that model; do not create a duplicate authority or
rename historical bindings casually. A manual INDEX.md and generated index.json
are both valid navigation aids. Keep authored content out of generated indexes.

Update the existing owner of a topic before creating another file. Create a
folder or anchor only when content needs it. Requirements and design may live
in a plan; extract lasting knowledge only when it has an independent use.
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

## Basic metadata

Use Markdown with YAML frontmatter. Type is foundational; types are extensible.
Outside raw, folder agreements require frontmatter with version, created, updated,
and a descriptive type. A short
description aids retrieval; title, tags, sources, and maturity status are optional.
Preserve unknown keys. Repair missing metadata when its values are known;
diagnostic findings do not activate Missions or block unrelated implementation.

Version is the document revision, independent of CLI releases and schemas.
Start at 0.1 and increment the minor number for meaningful content changes,
not typo fixes. Keep created immutable; update updated on meaningful edits.
Use RFC3339 timestamps with timezone. Recover old creation dates from reliable
history; omit unknown dates rather than fabricating them. Git preserves edits.

Raw has no metadata, naming, validation, or promotion obligation. Reading or linking provisional raw context
does not authorize publishing it. Keep existing ignored material ignored unless
retention is explicitly changed. Maintained knowledge can be normalized as useful,
without requiring promotion before work starts. Governed proof preserves what it
relies on rather than treating mutable unpublished text as a frozen guarantee.

Existing governed records keep their enforced identity, metadata, and mutation
rules. Ordinary edits do not rewrite frozen Missions, amend bound Contracts by
hand, or bypass archival protections. Product verification follows the actual
change and repository rules; optional knowledge checks are advisory.

## Optional metadata and link diagnostics

Run `python3 scripts/check-knowledge.py <workspace>`, where workspace is the
.spectacular directory, from the skill directory when checking context quality.
Python 3 and PyYAML are required. Exit 0 means no findings, 1 means advisory
findings, and 2 means unavailable inspection. The checker is read-only, skips raw
and governed records, checks folder metadata and file link targets, and reports
JSON. It does not certify section anchors, identities, or Contract bindings.
