---
status: verified
updated: 2026-08-07
related:
  - PLAN.md
---

# Tasks — identity-contract-workflow

<!--
  Executable checklist for one request.
  Lives at: .spectacular/requests/<slug>/TASKS.md

  Rules:
  - Group tasks by milestone using `### M<N> — <name>` headings.
  - Flush-left checkboxes are the COUNTED units: `- [ ]` open, `- [x]` done,
    `- [~]` deferred (not-open-not-done; shown separately in progress).
  - Indented `  - [ ]` sub-bullets are allowed as a nested acceptance checklist
    under a task, but are NOT counted — progress counts top-level only, so
    x/total stays comparable across requests.
  - `status:` in frontmatter should match parent PLAN.md.
  - Tasks are owned by the user. Engine never adds/removes/reorders tasks.
-->

## v1

### M1 — UUIDv7 identity and resolution substrate
- [x] Define a Bash 3.2-safe UUIDv7 generator and validation helpers with deterministic test seams.
- [x] Add a shared durable-record resolver that accepts a UUIDv7 or slug and recognizes legacy numeric aliases read-only.
- [x] Replace counter-based allocation for each new durable record kind without changing root-anchor or request-local task/milestone identity rules.
- [x] → check: focused identity and parallel-allocation tests prove unique IDs and slug resolution.

### M2 — Frontmatter, paths, and library retrieval
- [x] Define and validate the minimal common frontmatter spine: `id`, `slug`, `kind`, `scope`, `status`, `summary`, `related`, and `references`.
- [x] Update durable-record templates and list/detail views to use the spine while retaining only necessary entity-specific fields.
- [x] Move new specification creation and discovery to `specs/<slug>.md`; make slug renames path-only changes with stable ID references.
- [x] Add `scope: project` and `scope: request:<UUIDv7>` retrieval for the shallow shared collections; preserve request folders for execution state and artifacts.
- [x] → check: template and frontmatter/list tests prove slug paths and scoped frontmatter-only retrieval.

### M3 — Merged-contract execution gate
- [x] Add recommended forge-scope configuration and a local Git ancestry check for the configured shared base branch.
- [x] Derive contract approval from its merged spec PR rather than a duplicated digest, revision, or provider-specific frontmatter field.
- [x] Simplify contract-derived request provenance to `contract: <UUIDv7>` and preserve a single absolute `origin:` only for issue- or goal-originated work.
- [x] Reject request creation or activation when the contract is unmerged or the execution branch lacks the approved merge; give clear remediation.
- [x] → check: lifecycle tests cover draft, merged, and missing-ancestry outcomes.

### M4 — Preview-first numbered-ID migration
- [x] Inventory every numbered durable entity and its inbound references, including the existing duplicate `SPC-005` collision.
- [x] Implement a no-write migration preview that reports IDs, aliases, paths, and rewritten references before apply.
- [x] Implement the explicit apply path with archive-safe rewrites, retained read-only aliases, and a mapping receipt.
- [x] Add migration recovery/idempotence coverage and verify that the duplicate SPC-005 becomes two unambiguous UUIDv7 records.
- [x] → check: dry-run makes no writes; apply produces no duplicate IDs and resolves every legacy alias.

### M5 — Contract closure and regression proof
- [x] Update architecture, lifecycle, CLI help, templates, and skill references to describe the shipped workflow without renaming `spec` or `request`.
- [x] Add regression coverage for UUIDv7 generation, aliases, migration, spec PR approval, branch ancestry, slug resolution, and scoped retrieval.
- [x] Run the full verification suite and record the migration and compatibility evidence in the request validation artifacts.
- [x] → check: required syntax, guard, test, and doctor commands pass.

## v2 (deferred)

- [~] Rename current `spec`/`request` terminology, commands, paths, and public JSON keys to `contract`/`mission` — deferred until the identity migration has shipped and compatibility impact is separately approved.
- [~] Add forge-provider adapters beyond the recommended core Git scope — deferred because Markdown contracts and absolute origin URLs are provider-neutral already.
