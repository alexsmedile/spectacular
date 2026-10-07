# Architecture

Spectacular maintains durable context as navigable knowledge objects and their
relationships. The Markdown file is the unit of context; links connect the pieces
an agent needs for the current job. Folders group compatible objects without
requiring a workflow transition.

## Context and governance

| Surface | Purpose | Authority |
|---|---|---|
| Anchors | Current project truth, vocabulary, accepted constraints | Maintained context |
| Raw | Captures and unfinished drafts | Provisional |
| Atlas | Entities, concepts, maps, and lasting explanations | Maintained context |
| Decisions | Choices and rationale | Soft by default; governed validation optional |
| Plans | Outcomes, approaches, progress, and results | Maintained context |
| Specs | Reusable behavior and interface descriptions | Optional maintained context |
| Audits | Scope, inspected revision, findings, and limitations | Optional maintained context |
| Contracts | Agreements and their amendments | Mechanically governed |
| Missions and supporting records | Explicitly chosen execution and proof | Mechanically governed |

Agents use OKF object and entity principles: clear types, coherent objects,
explicit relationships, and navigable context. Read the relevant anchor and
follow links selectively. No full-workspace read or Mission is needed to start
ordinary work.

## Product parts

The Skill guides placement, interpretation, and work. Its planning route wraps
native host planning behavior and stores the project copy in plans/ when allowed.
It is useful alone. The Go
CLI validates governed frontmatter, computes fingerprints, and performs supported
atomic mutations. It retains 26 commands and the v2 module/schema identities in
product v3.

Discovery separates unclaimed soft context from governed records. A malformed
governed record remains an error; a soft plan or Atlas page does not cause an
unrelated command to refuse. Identity and schema claims preserve strict checking.
Contracts and other governed collections retain their enforced rules.

Generated `index.json` files describe governed records only. They are projections,
not context authority or a complete knowledge catalog. Authored `INDEX.md` files
provide optional navigation and are never overwritten by the generator.

## Boundaries

Root `docs/` is public product documentation managed independently by documentation
tools. `.spectacular/` holds durable context and optional governance. `skills/`
holds agent runtime guidance; `AGENTS.md` holds contributor rules. Plans and
scratch material stay out of product documentation.

Historical records remain valid under their original rules. Adopting soft context
does not rewrite completed Mission bindings, bypass Contract amendments, or
publish ignored raw drafts. See [Workspace layout](human-workspace-contract.md)
and [Process](process.md).
