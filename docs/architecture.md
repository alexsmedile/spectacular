# Architecture

Spectacular is five surfaces that do not overlap. Most of its design follows
from keeping them separate.

| Surface | Audience | What it is |
|---|---|---|
| `.spectacular/` | governance | The records: Missions, Proposals, Decisions, Evidence |
| `.spectacular/campaigns/` | planning | Optional durable roadmap maps; excluded from the record graph and CLI lifecycle |
| `skills/` | agents at runtime | Executable guidance the CLI and Skill load |
| `cmd/` + `internal/` | the machine | A typed CLI that validates and mutates records atomically |
| `docs/` | humans | What you are reading |

## Why records, not a database

Every governed record is git-versioned Markdown with typed YAML frontmatter.
There is no daemon and no opaque store. The frontmatter carries what a machine
must check — identities (UUIDv7), fingerprints (SHA-256), claim boundaries,
bindings — and the prose carries what a human and an agent must understand.

This is a deliberate trade. A database would be faster to query and impossible
to read in a diff. Because the records are files, `git log -S <fingerprint>`
recovers the exact text an agreement had when it was signed, and a reviewer can
see a governance change in a pull request like any other change.

## The record types

Thirteen types, in two groups.

**Top-level** — they exist independently:

- **Contract** — an accepted specification. Amended through `contract amend`,
  never by hand.
- **Mission** — a bounded, frozen agreement about one piece of work.
- **Proposal** — optional, mutable exploration. Carries no authority.
- **Decision** — a durable record of a choice and its reasoning.

**Mission-scoped** — they live inside a Mission bundle and have no meaning
outside it: Objective, Run, Checkpoint, Evidence, Assessment, Review, Handoff,
Gap, and Mission-local Decisions.

That split is why archiving a Mission archives its whole bundle: an Evidence
record separated from its Mission is not evidence of anything.

## Campaigns are plans, not records

A Campaign is an optional Markdown planning map in `.spectacular/campaigns/`.
It can sequence independent roadmap blocks and link candidate or active
Missions, but it grants no authority and is intentionally excluded from typed
CLI validation, identities, fingerprints, and lifecycle transitions. Campaigns
are mutable; a Mission's frozen envelope remains the execution authority.

A Mission may cite Campaign context in its Markdown body. Do not bind a Mission
to a Campaign in frozen frontmatter: roadmap edits must not create Mission drift.

## The mechanical interface

The CLI is thirteen commands, generated into a catalog at
[`mechanical-interface.md`](../skills/spectacular/generated/mechanical-interface.md).
The catalog is generated from the command registry, so it cannot drift from what
the binary does. When a document and the generated interface disagree, the
interface wins and the document is stale.

Commands are either `read-only` or `mutating`, and every mutating command is a
transaction: it either fully applies or fully rolls back. The test suite proves
this by injecting a fault at every write boundary and asserting the workspace is
unchanged.

Adding a command requires owner authorization. The count is reported, not
defended — a fourteenth command the owner authorized is correct, and a
thirteenth nobody asked for is not.

## What the machine does and does not do

The CLI performs only the mechanics that are safer or cheaper done mechanically
than by an agent: validating a record graph, computing fingerprints, writing
atomically, resolving paths, and reporting drift.

It does not decide anything. It cannot complete a Mission, accept a Proposal,
grant authority, or judge whether work is good. Those are owner acts, and the
system's value comes from refusing to blur that line.

```text
        agent proposes and builds
                  │
                  ▼
    ┌─────────────────────────────┐
    │  CLI: validate, fingerprint │  mechanical — no judgment
    │  write atomically, report   │
    └─────────────────────────────┘
                  │
                  ▼
          owner decides at a gate     ← authority lives here
```

## Authority and gates

Authority is explicit and narrow. A Mission's frontmatter names what the
operator may do — inspect, edit in scope, run checks, generate derived files.
Anything outside that list requires an owner gate.

Three gates matter most:

1. **Activation** — freezes the agreement. Everything after is judged against it.
2. **Review** — proof is recorded as a record, with a verdict.
3. **Completion** — only the owner, and only when declared Gaps are closed.

A Gap is a stated limit, not a defect. It is never closed by deletion: its entry
survives with a written resolution, so the reason something was impossible stays
recoverable.

## Freeze points

A completed Mission's contract fingerprint is a freeze point, not a stale
pointer. It records which agreement that Mission was executed against.
Amendments re-point only the live Mission; a completed one keeps its binding
forever. `mission check` reporting drift on a completed Mission is a notice, not
an error — the Mission stays valid.

This is the same instinct as the archive: a finished record is kept as it was,
not updated to match the present.

![Spectacular architecture](diagrams/architecture.svg)

## See also

- [Quickstart](quickstart.md) — run one Mission end to end.
- [Process](process.md) — the lifecycle and its gates in detail.
