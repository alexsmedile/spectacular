# Changelog

## [Unreleased]

## 2.3.0 — 2026-08-17

A signed Contract can now be amended, so a Contract that becomes wrong no longer
stays wrong, and every check a Contract declares either runs or fails loudly.

- Adds `contract amend <contract-ref> --gap <gap-ref> --by <owner>`, which closes one
  Gap by rewriting `blocked_on:` to `resolution:` and re-points the bound Missions, all
  as one recoverable transaction. `--dry-run` prints the resolution text, both
  fingerprints, and every Mission that would be re-pointed, and writes nothing. The
  amendment reaches the `gaps:` block and editorial fields only, and refuses while any
  bound Mission is live. A Gap is never closed by deleting it.
- Records each amendment in an append-only companion file beside the Contract, carrying
  time, owner, the Mission that declared the resolution, and the Contract fingerprint
  before and after. The log lives outside the Contract so an entry can record the
  fingerprint its own amendment produced.
- Stops refusing a completed Mission over a Contract that changed after it completed.
  The bound-Contract check previously gated on nothing, so amending a Contract refused
  every Mission ever bound to it — with no legal correction, because a completed Mission
  is never rewritten to satisfy a fingerprint. A completed Mission now reports the change
  as a notice; a live one still refuses.
- Adds `resolves_gaps:` to the Mission plan, naming the Gaps a Mission closes and the
  resolution text it will write. Both are frozen in the activation fingerprint, so the
  owner approves the exact wording at activation and a Mission cannot acquire the
  authority to amend a Contract afterwards. Completion refuses while a declared Gap is
  still open, and names the command that closes it.
- Adds `proposal check <ref>`, wiring the Proposal validator that existed with no caller.
  Proposals stay authored as Markdown wherever the author keeps them and are checked
  rather than generated; `proposal create` remains forbidden.
- Fixes four validation names that a Contract declared and nothing ran, so a Contract
  could promise a check that never executed. `frozen-fallbacks` is renamed to the
  declared `fallback-fingerprint-coverage`, the interface-dependency frozen-target check
  is split into its own registered validator, `ref-spelling-drift` is recognized as a
  notice rather than a validator, and `proposal-schema-v2` resolves to the newly wired
  Proposal check. Every declared name in every Contract now resolves, and one that
  matches nothing refuses.
- Reads `contract_version:`, which was carried by every Contract and used by nothing. It
  is validated as a positive integer and reported on `mission check`. A Mission bound to
  an earlier version is not refused or migrated: it ran against that version, and that is
  a true fact about it.
- Adds `amend-contract` to the owner authority vocabulary. A Mission declaring
  `resolves_gaps:` must list it under `requires_owner`, so the authority to change a
  signed agreement is declared rather than implied.
- Grows the command surface from ten to twelve, deliberately. Growth past twelve is a
  stop.

## 2.2.0 — 2026-08-17

Dead v1 surface removed, derived state reads recorded reviews, and Skill guidance
for branching, worktrees, and execution mode.

- Removes the unreachable v1 context-compiler chain — `internal/context`,
  `internal/projection`, and `internal/guardrails` — as one unit, and
  `internal/index`, the v1 predecessor of `discovery.Workspace.Lookup`. Prunes
  `internal/governance` to the transaction machinery that has live callers.
- Fixes the derived next action to consult the reviews a Mission carries. It
  previously fired on implemented Objectives alone, so a recorded review could not
  retire the instruction to record one. A review bound to a stale activation
  fingerprint keeps asking, since that is the drift the fingerprint exists to catch.
- Surfaces the underlying cause in human refusals. A YAML syntax error reported
  only `invalid_known_field`, sending readers hunting through field names while the
  parser's line number sat unprinted in the JSON envelope.
- Adds Skill guidance for choosing a branch and a worktree by what the job needs,
  including running a Mission session and a feedback session concurrently without
  either destroying the other's work.
- Adds an execution-mode question at activation — autopilot, checkpoints, or a
  named human-in-the-loop moment — so involvement is settled once instead of
  arriving one gate at a time.
- Accounts for every untracked working-tree path with a `.gitignore` rule that
  states its reason, and records `_snapshots/` as local recovery only.

## 2.1.1 — 2026-08-17

Test performance optimizations and contributor verification guidance.

- Optimizes `internal/missionbundle` test fixtures by caching template git repositories and parallelizing rollback subtests.
- Parallelizes 4-platform release archive compilation in `cmd/assemble-release`.
- Retains persistent Go test cache across verification runs and downstream installer/release test scripts.
- Codifies tiered verification guidance (`quick`, `acceptance`, `release`, `all`) in `AGENTS.md`.

## 2.1.0 — 2026-08-16

Governed-autonomy controls for preparation, delegation, review, and runtime limits.

- Adds adaptive preparation diagnostics and a frozen four-field completion criterion for every Mission claim.
- Adds criterion-driven automatic, clustered, and independent review without recursive critic loops.
- Validates Objective dependency DAGs and requires dependency-bound, disjoint Handoff claim scopes with explicit return contracts.
- Adds truthful hard, observed, and unsupported Autopilot caps for wall time, tokens, spend, parallel workers, and repair rounds.
- Fixes cold recovery for newly activated Missions before their first Checkpoint.

## 2.0.0-rc.2 — 2026-08-10

Human-operability correction for the v2 release candidate.

- Replaces flat UUID filenames with named project Anchors and cohesive Mission
  bundles while retaining UUIDv7 identity and SHA-256 revision fingerprints.
- Adds scoped human references for Missions, Objectives, Runs, Checkpoints,
  Evidence, Decisions, Gaps, Handoffs, and Assessments.
- Makes default CLI cards human-first and keeps exact machine data in `--json`.
- Commits deterministic, non-authoritative workspace and Mission indexes.
- Adds atomic whole-bundle Mission archival and a real self-hosted workspace.
- Replaces flat test fixtures with human-layout scenarios and adds a real-binary
  acceptance layer covering cold recovery, executable pointers, governed
  closure, archival, refusals, and zero-mutation reads.
- Fixes stale active indexes/directories after Mission archival, stable bundle
  placement after title changes, and invalid empty optional Evidence fields.

RC.1 is superseded because its machine-oriented workspace representation did
not satisfy the human-comprehension contract.

## 2.0.0-rc.1 — 2026-08-10

First externally consumable Spectacular v2 release candidate.

- Breaking: removes the v1 public product surface and compatibility paths; Git
  recovery pointers preserve the frozen history.
- Adds pointer-first retrieval and a governed Mission loop with explicit
  authority, Evidence, assessment, reconciliation, and closure boundaries.
- Ships native four-platform archives and a checksum-verifying installer.

Known limitations: v1 workspaces are unsupported; there is no migration or
compatibility reader; discovery is pointer-driven rather than broad search;
and release publication remains outside the local CLI.
