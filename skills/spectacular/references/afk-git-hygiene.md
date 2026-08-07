---
description: Opt-in AFK Git isolation, branch naming, archive-first cleanup, and verified pull-request handoff.
when_to_use: Authorized away-from-keyboard discovery or execution needs a branch, playground cleanup, or PR handoff.
---

# AFK Git Hygiene

AFK Git behavior is disabled by default and every mutating command is dry-run first. Authorization has three layers: a durable goal-scoped AFK run, project policy (`afk_git.enabled: true`), and invocation confirmation (`--apply --yes`). None implies permission to merge, delete remote branches, approve breaking changes, access undeclared accounts/browser sessions/sensitive data, or decide product questions.

Explicit AFK continuation or a built-in `/goal` activates the durable run. It lasts until goal completion/cancellation or an irreversible/blocking gate. Approval resumes the same run; it never silently broadens the recorded goal or allowed actions.

`afk run start --authority-record` is an opt-in M1 record format. It records
the human, session, base, derived `afk/<run-id>/integration` branch, and goal
nodes; `afk run event` appends structured evidence through the CLI. The record
is append-only by CLI and doctor-validated, not tamper-proof Markdown. These
commands do not stage, commit, amend, push, open a PR, merge, reset, stash, or
otherwise mutate Git state.

`spectacular session end` may show a read-only commit review for any session,
including an AFK session. That inspection is not an AFK Git action and grants no
additional authority: it never stages, commits, amends, pushes, merges, resets,
or stashes, and it does not satisfy any AFK authorization layer.

## Branch classes

| Work | Spectacular class |
|---|---|
| Unconfirmed specification | `spec/draft-<name>` |
| Feasibility spike | `spike/prototype-<SPK-ID>` |
| Competing idea | `fork/idea-<IDEA-ID>` |
| Approved execution | `feat/v<version>-<name>` |

An optional host prefix is prepended, never replaced—for example `codex/spike/prototype-spk-001`. `afk preflight` requires a clean tree and refuses writes on configured primary branches.

## Lifecycle

1. `afk run start ... --apply --yes` records the exact goal, optional request, allowed actions, HITL gates, and start time under `.spectacular/afk/runs/`. With `--authority-record`, it also requires `--authorized-by`, `--session`, `--base`, and one or more `--goal-node`; `afk run event <type> ... --apply --yes` adds a structured event. The same run moves `active → gated → active` when approval is needed.
2. `afk status` inspects policy/repository state without mutation; `afk propose` prints a deterministic name.
3. `afk start ... --apply --yes` creates the branch only when a run is active, policy is enabled, and the tree is clean; provenance is appended to `.spectacular/afk/branches.md`.
4. Work and verification happen on the isolated branch. Spike code is evidence, not implementation.
5. `afk cleanup` requires outcome and evidence. Dry-run is default. For a merged branch, apply creates and verifies `refs/spectacular/archive/<branch>-<timestamp>`, reports the tip SHA and restore command, then deletes the local and matching `origin` branch under the same `--apply --yes` confirmation. A missing remote is already clean; a remote tip that differs from the verified merged tip is a blocker. `--keep-remote` is the explicit opt-out. Abandoned branches retain their remote branch. Cleanup leaves no untracked workspace receipt.
6. `afk pr` requires a verified request, passing VERIFY-LOG outcome, approved source spec, fresh `--tests-passed`, policy permission, and explicit apply. Breaking changes require a separate approval flag. It delegates to the shared GitHub work-bridge manifest and opens a **draft** `[Spectacular] Executed: <version> - <name>` PR; `spectacular github pr ready` owns the later current-head/check/confirmation gate. Every generated PR ends with `_Filed with Spectacular_` linked to the project. Pass repeatable `--summary "<change>"` and `--validation "<result>"` for precise bullets; otherwise the request goal supplies the summary. Never paste shell setup, environment variables, or opaque command transcripts into the PR body. The dry run prints the exact title and body before it opens the PR; neither command merges.

The durable Git recovery ref preserves the tip SHA and exact restore command printed by cleanup. Recovery never depends on reflog retention or an uncommitted workspace receipt.

These refs are cold recovery under [[artifact-retention]], not live branches or normal agent context. Cleanup may delete the local spike/fork branch after evidence is preserved and the ref verifies. Pruning old recovery refs is a separate, explicitly approved retention operation; routine AFK cleanup never does it.
