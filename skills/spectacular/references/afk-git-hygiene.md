---
description: Opt-in AFK Git isolation, branch naming, archive-first cleanup, and verified pull-request handoff.
when_to_use: Authorized away-from-keyboard discovery or execution needs a branch, playground cleanup, or PR handoff.
---

# AFK Git Hygiene

AFK Git behavior is disabled by default and every mutating command is dry-run first. Authorization has three layers: a durable goal-scoped AFK run, project policy (`afk_git.enabled: true`), and invocation confirmation (`--apply --yes`). None implies permission to merge, delete remote branches, approve breaking changes, access undeclared accounts/browser sessions/sensitive data, or decide product questions.

Explicit AFK continuation or a built-in `/goal` activates the durable run. It lasts until goal completion/cancellation or an irreversible/blocking gate. Approval resumes the same run; it never silently broadens the recorded goal or allowed actions.

## Branch classes

| Work | Spectacular class |
|---|---|
| Unconfirmed specification | `spec/draft-<name>` |
| Feasibility spike | `spike/prototype-<SPK-ID>` |
| Competing idea | `fork/idea-<IDEA-ID>` |
| Approved execution | `feat/v<version>-<name>` |

An optional host prefix is prepended, never replaced—for example `codex/spike/prototype-spk-001`. `afk preflight` requires a clean tree and refuses writes on configured primary branches.

## Lifecycle

1. `afk run start ... --apply --yes` records the exact goal, optional request, allowed actions, HITL gates, and start time under `.spectacular/afk/runs/`. The same run moves `active → gated → active` when approval is needed.
2. `afk status` inspects policy/repository state without mutation; `afk propose` prints a deterministic name.
3. `afk start ... --apply --yes` creates the branch only when a run is active, policy is enabled, and the tree is clean; provenance is appended to `.spectacular/afk/branches.md`.
4. Work and verification happen on the isolated branch. Spike code is evidence, not implementation.
5. `afk cleanup` requires outcome and evidence. Dry-run is default. Apply creates and verifies `refs/spectacular/archive/<branch>-<timestamp>`, records the tip SHA and restore command, then deletes the local branch. Remote deletion is refused.
6. `afk pr` requires a verified request, passing VERIFY-LOG outcome, approved source spec, fresh `--tests-passed`, policy permission, and explicit apply. Breaking changes require a separate approval flag. It opens `[Spectacular] Executed: <version> - <name>` and stops before merge.

Archive records contain a durable Git ref, tip SHA, and exact restore command. Recovery never depends on reflog retention; remote branches are untouched.
