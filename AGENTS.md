# Contributor guide

This is Spectacular v2. The root Go module, `cmd/spectacular`, `skills/`,
`install/`, and `.spectacular/` are the only live product surface.

`CLAUDE.md` is a compatibility symlink to this file. `AGENTS.md` is
authoritative; edit it rather than the symlink.

Run `bash test/verify.sh all` before release changes or Mission completion.
Use tiered verification during development:
- `bash test/verify.sh quick`: static checks + unit tests in `cmd/`, `internal/`, `install/` (fastest inner loop).
- `bash test/verify.sh acceptance`: static checks + end-to-end acceptance fixtures.
- `bash test/verify.sh release`: 4-platform compilation, checksums, installer/rollback/recovery, and plugin manifests.
- `bash test/verify.sh all`: full race-detector test suite and release distribution gate.

Do not reintroduce v1 commands,
compatibility readers, migrations, generic record/search verbs, or a second
package root. Keep release version values aligned through `VERSION` and the
generated mechanical interface.

The public command surface is twelve commands. Adding one requires owner
authorization: argue the case in a Proposal, state the count before and after, and
let the owner decide. An agent never adds a command on its own reading of intent.
The number is reported, not defended — a thirteenth command the owner authorized is
correct, and a twelfth that nobody asked for is not. `proposal create` stays
forbidden.

A Contract is amended through `contract amend`, never by editing a bound Contract by
hand. An amendment may reach the `gaps:` block and editorial frontmatter only;
changing a field that states what was agreed is a `contract_version:` bump instead.
A Gap is never closed by deleting it — its entry survives with a stated resolution.

## `docs/` is human-facing product documentation

`docs/` holds the public documentation for Spectacular — the kind of material
that would be published to a `docs.<domain>` site. Its audience is a human
reading about the product, not an agent executing against it.

Write it as product documentation: concepts, guides, reference pages, and
diagrams that explain what Spectacular does and how to use it. Prose over
record structure.

Keep it distinct from the other surfaces:

| Surface | Audience | Purpose |
|---|---|---|
| `docs/` | humans reading the product docs | concepts, guides, reference, diagrams |
| `AGENTS.md` | coding agents | contributor rules and constraints for this repo |
| `skills/` | agents at runtime | executable guidance the CLI and Skill load |
| `.spectacular/` | governance | Missions, Proposals, Decisions, Evidence |

Rules:

- `docs/` is documentation only. It is never loaded as agent context and must
  not become a second home for Skill guidance or governance records.
- Governance records stay in `.spectacular/`. Do not narrate Mission state in
  `docs/`; link to the record instead.
- Nothing in `docs/` is authoritative for behavior. When docs and the generated
  mechanical interface disagree, the interface wins and the doc is stale.
- Keep it publishable: no local paths, no operator names, no scratch notes.
