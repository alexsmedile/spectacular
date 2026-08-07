---
id: SPC-008
type: specification
status: draft
target_version: "tbd"
supersedes: ""
updated: 2026-08-07
summary: "Adopt UUIDv7 identities, slug-only spec paths, and merged spec-PR approval before request execution"
related: []
---

# SPC-008 — Adopt UUIDv7 identities, slug-only spec paths, and merged spec-PR approval before request execution

## Intent

Give every durable Spectacular record a branch-safe immutable identity while
keeping slugs as the normal human interface, and make a merged Markdown spec PR
the sole approval gate before a request may execute its capability contract.

## Requirements

- Generate a UUIDv7 `id:` for every new durable record: specification,
  request, decision, question, research, spike, idea, vision, memory, feedback,
  audit, fix, and session. Root anchors retain their fixed-path identity;
  request-local tasks and milestones retain their parent/request-local identity.
- Give each durable record a human-facing `slug:`. CLI conversation and ordinary
  projections resolve and display slugs; IDs remain available for exact links,
  automation, and repair.
- Replace numeric canonical IDs as durable references with UUIDv7 IDs. Existing
  `SPC-NNN`, `DEC-NNN`, and other numbered records migrate preview-first with
  an explicit mapping and remain resolvable as read-only compatibility aliases.
- Store new capability specs at `specs/<slug>.md`. A slug rename is a file rename
  only; durable references continue to use the immutable ID.
- Use the minimal common frontmatter spine: `id`, `slug`, `kind`, `scope`,
  `status`, `summary`, `related`, and `references`, plus only entity-specific
  lifecycle fields. Do not add digest, version-pin, commit-SHA, or
  provider-specific identity fields.
- Keep one shallow, typed shared soft library in the top-level collections.
  `scope: project` marks reusable or cross-cutting knowledge;
  `scope: request:<UUIDv7>` marks a durable record owned by one request.
  Request folders contain execution state and artifacts only, not duplicate
  local questions/research/spikes/decisions databases.
- Keep trivial request-only choices in `PLAN.md` rather than creating a record.
  Promote a record to `scope: project` only when another request or future
  session needs to retrieve it independently.
- Let a request derived from a capability contract carry only
  `contract: <UUIDv7>`. Remove redundant spec-source fields through the
  migration. Issue- or goal-originated work may record a single absolute URL in
  `origin:`.
- Treat a merged spec PR into the configured shared base branch as approval.
  A spec remains draft while it is worked on its spec branch. A request may
  activate only when its execution branch contains the merged spec PR.
- Keep Markdown specifications as the portable capability-contract layer.
  GitHub, GitLab, and Bitbucket work items remain linked origins/references;
  provider APIs and PR/MR mechanics belong to forge adapters, not frontmatter.
- Add a linked glossary/catalog that defines capability, capability contract,
  design and architectural work, execution work, identity, slug, and spec
  drift. Preserve `spec` and `request` as the current normative terms; record
  `capability`/`contract` and `mission` only as unconfirmed future vocabulary.

## Evidence and decisions

- User-confirmed direction, 2026-08-07:
  - UUIDv7 is the immutable ID format.
  - Specs use `specs/<slug>.md` paths.
  - Merging a spec PR is approval; a request begins only afterward on a branch
    containing that merge.
  - Forge work items stay collaborative origins, not canonical specs.
  - Existing numbered entities migrate; the broader terminology rename waits.
- [[glossary/index]] records the shared vocabulary and the distinction between
  capability contracts and execution requests.
- Current implementation evidence: `_next_canonical_id` allocates the next
  number from the checked-out tree and `doctor wayfinding` detects duplicates
  only after they coexist. The live duplicate `SPC-005` demonstrates the
  branch-collision failure this contract removes.

## Constraints

- Preserve Markdown and frontmatter as the cheap, local, Git-reviewable signal
  layer; do not introduce a database or a network requirement for identity
  allocation.
- Do not introduce a deep knowledge hierarchy below individual requests. The
  existing top-level typed collections remain the shared library; `requests/`
  remains the execution tree.
- Keep Bash 3.2 compatibility.
- Migration is preview-first, reference-aware, archive-safe, and reversible
  until its apply boundary is explicitly confirmed.
- Do not implement the future terminology rename (`spec` to `contract` or
  `capability`; `request` to `mission`) in this scope.
- Do not make GitHub, GitLab, or Bitbucket work items the canonical editable
  spec source in this scope.

## Milestones

- M1 — Vocabulary and target contract are published: glossary entries and the
  architecture/lifecycle documentation define the identity, slug, contract,
  scope, spec-PR, and request rules without renaming current commands.
- M2 — New records are branch-safe: CLI templates and validation create and
  accept UUIDv7 IDs plus slugs and scope, while normal views resolve by slug.
- M3 — Spec PRs gate execution: approval is derived from merge to the configured
  base; contract-derived request activation checks that the execution branch
  contains the approved merge.
- M4 — Existing records migrate safely: dry-run produces a complete mapping,
  apply rewrites references and paths, legacy numeric aliases remain readable,
  and the duplicate `SPC-005` is resolved explicitly.
- M5 — The contract is proven: focused concurrency, migration, link-resolution,
  ancestry, and compatibility tests pass; docs and templates describe shipped
  behavior.

## Tasks

Derive implementation tasks only after this spec is merged and approved.

## Dependencies

- No external provider integration is required for the core workflow.
- The migration must account for all existing numbered entity collections before
  the new identity validator becomes enforcing.

## Validation

- Two independently created records on separate branches receive distinct IDs
  and merge without identity conflict.
- `specs/<slug>.md` is created for each new spec; renaming its slug leaves all
  ID-based inbound references valid.
- A draft spec on a spec branch cannot seed or activate a contract-derived
  request; the merged spec can.
- A frontmatter-only list/filter retrieves records by `kind`, `scope`, `status`,
  and `summary` without loading their bodies; a selected record alone loads its
  full Markdown.
- A request-scoped research record is retrievable by its request scope, while a
  trivial PLAN decision creates no duplicate library record.
- An execution branch missing the approved spec merge is rejected with a clear
  remediation message.
- Migration dry-run makes no writes and reports every ID/path/reference mapping;
  applied migration leaves no duplicate durable IDs and resolves legacy aliases.
- GitHub, GitLab, and Bitbucket origin URLs survive round-trip parsing without
  provider-specific frontmatter fields.
- `bash tests/run.sh`, `scripts/hooks/pre-commit --check`, and relevant
  `spectacular doctor` areas pass after implementation.

## Deliverables

- UUIDv7/slug identity and resolution substrate in the CLI.
- Frontmatter-led shared-library list/find views with `scope` filtering.
- Updated record templates, lifecycle/architecture contracts, and glossary.
- Spec-branch/spec-PR approval and execution-branch ancestry gate.
- Preview-first numbered-ID migration with compatibility aliases and a mapping
  receipt.
- Regression coverage for concurrency, migration, lifecycle gating, and aliases.

## Confirmation

draft — not eligible for implementation until its spec pull request is merged
to the configured shared base branch.
