# Contributor guide

## Default workspace workflow

Spectacular maintains durable context as small linked OKF objects. Ordinary work uses anchors for current project
truth, `raw/` for captures and quick-start drafts, `atlas/` for maps and lasting
explanations, `decisions/` for choices, and `plans/` for outcomes and approaches.
Agents judge placement, update an existing document before creating another, and
split content when it has independent meaning or lifetime. Aim for 50–120 lines
per soft document as guidance, not a mechanical limit. Follow relative Markdown
links or path-qualified wikilinks; repair incoming links on moves.

`type:` is foundational. Maintained soft documents use basic document version,
creation time, and update time; recover historical dates honestly and preserve
unknown metadata. No extra record identity or lifecycle is required. The runtime
details live in `skills/spectacular/references/knowledge.md`.
Metadata profiles live in `skills/spectacular/references/profiles.md`: compact is
the soft default, minimal remains valid, extended fields are optional, and raw is
freeform. These conventions require no profile label and do not select governance.

`/spectacular plan` is a Skill route over native host planning controls, with the
durable project copy in `.spectacular/plans/` when writes are allowed. Optional
`specs/` and `audits/` contain reusable behavior and inspection findings; create
them only for independently useful content. They do not replace governed
Contracts, Reviews, or Handoffs. Canonical Mission-owned records follow the path
reported by the CLI.
Folder-specific agreements live in `skills/spectacular/knowledge-folders.yaml`.
Outside raw/sketch/scratchpad, maintained documents require metadata appropriate
to those agreements. Raw has no metadata, naming, or promotion obligations.
Prefer VOCABULARY.md for new domain-model anchors; retain an existing ONTOLOGY.md
without duplicating its authority or rewriting historical bindings. Manual INDEX.md
and generated index.json are both valid; do not overwrite manual navigation
with generated output. Plans go in `.spectacular/plans/`, sketches in the existing
scratchpad-role folder, never in product `docs/`.

The directory does not activate governance. The user must explicitly select
governed work or a named Mission. CLI absence and unrelated Mission drift do not
block ordinary implementation. Existing governed records retain their mutation,
freeze, and archival rules. Soft Decisions use descriptive unnumbered filenames
without governed identity or schema claims; historical numbered Decisions remain
governed. Mechanical validation is recommended for Decisions, not obligatory.
Contracts retain their mechanical validation, amendment, and version pipeline.
This operative policy replaces D23/D24's restriction on consulting raw
material and requiring all Decisions to be governed, while retaining D24's
schema-honesty rule. Historical Decisions remain unchanged.

This is Spectacular v3; the Go module and mechanical record schemas retain v2 identities. The root Go module, `cmd/spectacular`, `skills/`,
`install/`, and `.spectacular/` are the only live product surface.

`CLAUDE.md` is a compatibility symlink to this file. `AGENTS.md` is
authoritative; edit it rather than the symlink.

Run `bash test/verify.sh all` before release changes or Mission completion.
Use tiered verification during development:
- `bash test/verify.sh preflight`: Tier 0 static syntax/tree sanity + Tier 1 contract
  drift on the live Mission. Read-only, emits a measured
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
See `D24-schema-field-mechanically-governs-frontmatter`, which amends `D23-workspace-entities-type-and-schema`.

Never hand-write a frontmatter template into documentation or a test. A published
template is retrieved from `--schema` and is round tripped through the validator
that emitted it. A template that names a field the parser does not read produces a
document that validates while its meaning silently disappears.

`.spectacular/raw/` is gitignored and outside the governed graph. Agents may read
and link provisional drafts there to start work. That permission does not publish
them or make them accepted truth. Governed proof preserves any material it relies
on. Never place a governed record in `raw/`: it would escape governed review.

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
| `.spectacular/` | durable context and optional governed execution | anchors, entities, maps, plans, Decisions, Contracts, Missions, proof |

Rules:

- `docs/` is documentation only. It is never loaded as agent context and must
  not become a second home for Skill guidance or governance records.
- Its structure is owned by documentation tools such as pageworks, not Spectacular.
- Governance records stay in `.spectacular/`. Do not narrate Mission state in
  `docs/`; link to the record instead.
- Nothing in `docs/` is authoritative for behavior. When docs and the generated
  mechanical interface disagree, the interface wins and the doc is stale.
- Keep it publishable: no local paths, no operator names, no scratch notes.
