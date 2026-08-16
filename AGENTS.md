# Contributor guide

This is Spectacular v2. The root Go module, `cmd/spectacular`, `skills/`,
`install/`, and `.spectacular/` are the only live product surface.

`CLAUDE.md` is a compatibility symlink to this file. `AGENTS.md` is
authoritative; edit it rather than the symlink.

Run `bash test/verify.sh all` before release changes. Do not reintroduce v1 commands,
compatibility readers, migrations, generic record/search verbs, or a second
package root. Keep release version values aligned through `VERSION` and the
generated mechanical interface.

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
