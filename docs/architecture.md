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

## Delivery milestones

A milestone names who can complete which workflow, in which environment, and
what observation closes the outcome. Start with a thin integrated path and
resolve the uncertainties that could invalidate it. Prerequisites are justified
by the workflow or exposure they enable; an MVP label does not prescribe a
database, authentication system, event broker, or deployment platform.

Evaluate usability, integrity, distribution, operation, value validation, and
publication separately. A prototype can answer a design question; an integrated
workflow can establish technical usability; a representative user experiment
can test a value hypothesis. Alpha and beta have project-defined audiences,
limitations, and exit evidence.

Plans own delivery gates and proof links. Requirements own needs and acceptance.
Use MoSCoW, risk, and reversibility to bound the next useful increment. Companion
skills contribute architecture decisions, executable checks, data migration,
interface compatibility, and durable-delivery proof only when relevant.

## Implementation boundaries

The CLI uses a single command registry and application service facade. Mission
readers consume detached data through an injected loader. Pure Run transition
rules are separated from persistence effects, and coordinated writes retain
atomic transactions and recovery. Production dependency rules are checked
automatically, including platform-specific source files.
