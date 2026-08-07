---
id: 019fdd18-96a0-73b5-81b5-38cda98d36f4
slug: worktree-coordination
kind: spec
scope: project
status: draft
target_version: "tbd"
supersedes: ""
updated: 2026-08-07
summary: "Add durable, evidence-backed coordination before Git mutations in concurrent worktrees"
related: []
references: []
---

# Worktree coordination — durable, evidence-backed safety before Git mutation

## Intent

Give agents and maintainers one deterministic, read-only coordination view before
they create branches, commit, push, publish, merge, archive, or clean up Git
state. The safe path must preserve unrelated work, make ambiguous ownership
visible, and remain useful with local Git alone.

## Requirements

- Add `spectacular workspace preflight` as a read-only coordination report with
  exactly one proposed safe next action.
- Add `spectacular workspace plan` as a read-only, path-to-disposition plan.
- Before a branch, commit, push, PR/MR handoff, merge, archive, or cleanup
  mutation owned by Spectacular, run the same coordination preflight. Refuse
  the mutation when the plan exposes unresolved ownership or another blocker.
- Inspect, without mutation: current branch and upstream; configured shared
  base and divergence; local/remote base freshness; staged, unstaged, and
  untracked paths; local branch merge state; active/planned requests; and
  provider PR/MR facts when an adapter can supply them.
- Treat each changed path as unowned until evidence associates it with the
  current request, another request/spec, explicit user intent, or an existing
  branch/commit. Report its inferred owner, confidence, related request/spec/PR
  when known, and one of these dispositions:
  `belongs-to-current-request`, `belongs-to-known-work`,
  `needs-preservation-branch`, or `discardable-with-explicit-confirmation`.
- Keep Git ancestry as the core merge/preservation evidence. GitHub, GitLab,
  and Bitbucket PR/MR facts are optional adapter evidence; missing provider
  data is reported as `unknown`, never guessed.
- Add `spectacular workspace preserve <slug> --paths <...> --apply --yes`.
  Its preview must show the exact branch, scoped staged paths, commit message,
  and any push/PR action. Apply creates a clearly named branch, stages only the
  declared paths, creates a scoped preservation commit, pushes only when
  configured and authorized, and reports the resulting branch/commit/PR state.
- Add read-only `spectacular workspace cleanup`, which inventories merged,
  open, stale, declined, and unverifiable branches. Add
  `spectacular workspace cleanup <branch> --apply --yes` only for deletion
  after explicit verification and preview.
- Never use `git add -A` in mixed worktrees. Never silently reset, stash,
  discard, overwrite, or commit unrelated paths. Stash is not the normal
  preservation mechanism.
- Never merge directly to the configured shared base unless the user explicitly
  requests it. Never delete a branch until its remote merge/decline state is
  verified where available and its work remains accessible through a commit or
  PR/MR.
- A merged feature branch is deletable only after local shared base is
  fast-forwarded to its remote and the feature tip is reachable from it. An
  open PR/MR blocks deletion. A declined committed branch is merely eligible
  after explicit confirmation and durable preservation evidence.
- Extend request `PLAN.md` frontmatter with only the branch/provenance signal
  needed for ownership checks (optional `branch:` and `base:`); retain existing
  `contract:` and `origin:` semantics. Do not create a permanent per-file
  ownership registry, copy live PR state, or duplicate commit SHAs in Markdown.
- Persist only consequential evidence: scoped preservation commits and
  completed branch-deletion receipts. Preflight reports remain ephemeral.

## Evidence and decisions

- User-confirmed direction, 2026-08-07:
  - Safety must be automatic at dangerous moments but low-friction while
    ordinary coding proceeds.
  - Ambiguous ownership is surfaced rather than guessed; durable named branches
    plus scoped commits replace stash as the normal preservation path.
  - The core remains provider-neutral and Bash 3.2 compatible.
  - The requested identity-contract sequence is binding: draft this spec on
    `spec/worktree-coordination`, treat merge to configured shared base as
    approval, and create the implementation request branch only afterward.
- Existing evidence:
  - [[SPC-008-durable-identity-and-contract-pr-workflow]] establishes merged
    spec-PR approval and execution-branch ancestry as the implementation gate.
  - [[afk-git-hygiene]] already provides narrowly scoped AFK branch isolation
    and archive-first cleanup; this capability generalizes coordination without
    conflating normal interactive work with AFK authorization.

## Constraints

- Bash 3.2 is the runtime baseline.
- Markdown plus frontmatter remain the durable signal layer; no database,
  daemon, or mandatory network/provider dependency is introduced.
- Read-only commands must make no Git or workspace writes, including fetch,
  branch creation, staging, or status repair.
- Mutators are preview-first and require both `--apply` and `--yes`; a command
  must print its exact state-changing commands before execution.
- Existing request lifecycle, AFK Git hygiene, and forge adapter behavior remain
  compatible unless this contract explicitly strengthens their safety gate.

## Milestones

- M1 — Coordination evidence contract and minimal request branch metadata are
  defined, including provider-neutral unknown states and durable receipts.
- M2 — Read-only preflight and planning produce deterministic path dispositions,
  blockers, and a single safe next action without mutating Git state.
- M3 — Scoped preservation safely creates a dedicated branch and commit for
  declared paths while leaving other changes intact.
- M4 — Verified cleanup inventories candidates and deletes only explicitly
  confirmed, remotely/ancestrally preserved branches.
- M5 — Existing Spectacular Git-mutating flows invoke the common preflight, and
  the compatibility and regression suite proves no destructive escape path.

## Validation

- A dirty worktree containing one current-request file and one unrelated file
  yields separate read-only dispositions and no mutation.
- An untracked draft spec is preserved onto only its dedicated spec branch and
  scoped commit.
- A staged unrelated file makes a current-request commit refuse until scope is
  explicit.
- Cleanup refuses a merged branch while local shared base is stale, and refuses
  an open PR/MR branch.
- A declined but committed branch is reported as eligible only after explicit
  confirmation and accessible preservation evidence.
- With no remote, GitLab, or Bitbucket adapter data, core Git behavior works and
  provider facts are reported as `unknown`.
- Every destructive-path test proves no reset, stash, discard, or overwrite of
  unrelated work.
- `bash tests/run.sh`, Bash syntax checks, relevant doctor coverage, and AFK and
  request lifecycle regressions pass.

## Confirmation

draft — not eligible for implementation until its spec pull request is merged
to the configured shared base branch.
