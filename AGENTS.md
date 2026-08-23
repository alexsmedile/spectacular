# Contributor guide

This is Spectacular v2. The root Go module, `cmd/spectacular`, `skills/`,
`install/`, and `.spectacular/` are the only live product surface.

`CLAUDE.md` is a compatibility symlink to this file. `AGENTS.md` is
authoritative; edit it rather than the symlink.

Run `bash test/verify.sh all` before release changes or Mission completion.
Use tiered verification during development:
- `bash test/verify.sh preflight`: Tier 0 static syntax/tree sanity + Tier 1 contract
  drift on the live Mission. Read-only, sub-second, emits a
  `spectacular.preflight-receipt.v1` JSON receipt on stdout. Run it before any heavy
  tier; if it fails, repair and do not run `acceptance`, `release`, or `all`.
  `PREFLIGHT_MISSION_REF=<ref>` pins the Mission checked;
  `PREFLIGHT_ALL_MISSIONS=1` sweeps every Mission instead of the live one.
- `bash test/verify.sh quick`: static checks + unit tests in `cmd/`, `internal/`, `install/` (fastest inner loop).
- `bash test/verify.sh acceptance`: static checks + end-to-end acceptance fixtures.
- `bash test/verify.sh release`: 4-platform compilation, checksums, installer/rollback/recovery, and plugin manifests.
- `bash test/verify.sh all`: full race-detector test suite and release distribution gate.

Never run `verify.sh` during orientation, conversational answers, status queries, or pure Markdown/documentation edits. `verify.sh` tests the Go codebase of this repository; use `quick` only after modifying Go source code, and `all` only at a final Mission completion or release gate.

Do not reintroduce v1 commands,
compatibility readers, migrations, generic record/search verbs, or a second
package root. Keep release version values aligned through `VERSION` and the
generated mechanical interface.

Adding or modifying public CLI commands requires explicit user/owner
authorization. An agent must never introduce or alter a command on its own
reading of intent. When proposing a new command, state the rationale, the current
and proposed command count, and wait for owner approval. `proposal create` stays
forbidden.

Every workspace document that names an entity declares `type:`. A `schema:` field is
a narrower claim: Spectacular governs this document and its frontmatter is under
mechanical check. Add `schema:` only when a command validates the document and
refuses on drift — a schema nobody enforces invites tooling to rely on a guarantee
that does not exist. An Atlas therefore carries `type:` alone. Mechanical checking
reaches the frontmatter; the body is not enforced and a body check may only warn.
See `D24-accepted`, which amends `D23-accepted`.

Never hand-write a frontmatter template into documentation or a test. A published
template is retrieved from `--schema` and is round tripped through the validator
that emitted it. A template that names a field the parser does not read produces a
document that validates while its meaning silently disappears.

`.spectacular/raw/` is gitignored and skip-listed. Nothing there is an entity and
nothing cites it. Never write a governed record, Proposal, or Decision into `raw/`:
it would never appear in review and never be committed.

A Contract is amended through `contract amend`, never by editing a bound Contract by
hand. An amendment may reach the `gaps:` block and editorial frontmatter only;
changing a field that states what was agreed is a `contract_version:` bump instead.
A Gap is never closed by deleting it — its entry survives with a stated resolution.

An amendment refuses while a bound Mission that did not declare the Gap is live. The
Mission that declared it is the exception and closes it while live; an owner
`--resolution` override is never exempt, because its wording was typed at a prompt
rather than approved at an activation gate.

A completed Mission's `contract.fingerprint` is a freeze point, not a stale pointer:
it records which agreement that Mission was executed against. Amendments re-point only
the live Mission. Never re-point, hand-edit, or otherwise "fix" a completed Mission's
binding — `mission check` reporting `contract-drift` on one is a notice, the Mission
stays `valid=true`, and `git log -S <fingerprint>` recovers the Contract text as it
was. See `D10-repoint`.

A Proposal that has shipped is retired, not left at `draft`. Nothing writes a Proposal's
status — `proposal check` validates the value it finds and no command advances it — so an
absorbed Proposal reads `draft` until an owner says otherwise. Retiring one means naming its
resolver in `resolved_by:`, setting `accepted`, and moving it to
`.spectacular/archive/proposals/` with `archive_authorization:` and
`archive_input_fingerprint:`, exactly as an archived Mission carries them. Write
`resolved_by:` before the move, never after: once the record leaves `proposals/`, that field
is the only thing tying it to the work that answered it. A Proposal is absorbed when the
question it asked was answered, not when most of it was — P5 shipped three of four
directions and stays live. Live `proposals/` holds open questions only. See
`D11-proposal-retirement`.

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
