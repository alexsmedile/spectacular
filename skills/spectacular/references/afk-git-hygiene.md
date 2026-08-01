---
description: Opt-in AFK Git isolation, branch naming, archive-first cleanup, and verified pull-request handoff.
when_to_use: Authorized away-from-keyboard discovery or execution needs a branch, playground cleanup, or PR handoff.
---

# AFK Git Hygiene

AFK Git behavior is disabled by default and every mutating command is dry-run first. Authorization has two independent layers: project policy (`afk_git.enabled: true`) and invocation confirmation (`--apply --yes`). Neither implies permission to merge, delete remote branches, approve breaking changes, or decide product questions.

## Branch classes

| Work | Spectacular class |
|---|---|
| Unconfirmed specification | `spec/draft-<name>` |
| Feasibility spike | `spike/prototype-<SPK-ID>` |
| Competing idea | `fork/idea-<IDEA-ID>` |
| Confirmed execution | `feat/v<version>-<name>` |

An optional host prefix is prepended, never replaced—for example `codex/spike/prototype-spk-001`. `afk preflight` requires a clean tree and refuses writes on configured primary branches.

## Lifecycle

1. `afk status` inspects policy/repository state without mutation.
2. `afk propose` prints a deterministic name.
3. `afk start ... --apply --yes` creates the branch only when policy is enabled and the tree is clean; provenance is appended to `.spectacular/afk/branches.md`.
4. Work and verification happen on the isolated branch. Spike code is evidence, not implementation.
5. `afk cleanup` requires outcome and evidence. Dry-run is default. Apply writes `.spectacular/archive/afk-branches/<timestamp>-<branch>.md` before deleting the local branch; remote deletion is refused.
6. `afk pr` requires a verified request, passing VERIFY-LOG outcome, current source spec, fresh `--tests-passed`, policy permission, and explicit apply. Breaking changes require a separate approval flag. It opens `[SpecTACular] Executed: <version> - <name>` and stops before merge.

Archive records state `recoverable_from: git reflog`. Local deletion is recoverable while reflog data remains; remote branches are untouched.
