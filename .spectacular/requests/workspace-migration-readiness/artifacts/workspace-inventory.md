---
updated: 2026-08-03
against: 70745412435a846089ee5e3019427e6b2af4a3db
traffic: parallel
---

# Workspace inventory

## Baseline

- Local `main`, `origin/main`, and the request branch base all resolved to `70745412435a846089ee5e3019427e6b2af4a3db` after `git fetch origin`.
- `git rev-list --left-right --count origin/main...main` returned `0 0` before request edits.
- GitHub reported no open pull requests.
- Spectacular reported no active request before this audit; three unrelated requests were planned. This request was admitted as `parallel` because all outputs stay request-local.
- `workspace_schema: 2.0` matches `CURRENT_SCHEMA="2.0"`; the migration registry contains `0.4 → 0.5`, `0.5 → 0.6`, and `0.6 → 2.0`.
- Scaffolding build `b41` mechanically changed only `config.yaml:last_build` from 39 to 41 outside the request folder; b40 was already reserved by the GitHub-native lifecycle candidate.
- Full doctor reported zero errors. Before this request's expected `docs_impact: pending` warning, the only warnings were three legacy memory records lacking lifecycle status.

## Shared committed boundary

| Surface | Purpose / authority | Freshness and retention | Readiness disposition |
|---|---|---|---|
| Root anchors (`PRD`, `PRINCIPLES`, `ARCHITECTURE`, `POLICY`, `STACK`, `AGENTS`, `config`) | Project-wide intent, rules, structure, and configuration | Live; snapshot before canonical edits | Keep committed; schema contract must name required anchors |
| `specs/`, `roadmaps/`, `questions/` | Current capability context, runway, and user blockers | Live while active; archive/compact by lifecycle | Keep committed and progressively retrieved |
| `requests/` | Shared request intent, tasks, evidence, and handoff | Live until verified/archive | Keep committed; request slug is the collaboration namespace |
| `decisions/`, `memories/`, `sessions/`, `ideas/`, `fixes/` | Durable project knowledge | Collection-specific retention/indexing | Keep committed; never shadow from local state |
| `_snapshots/` | Canonical recovery history | Stale-safe; retention policy already reports prunable entries | Keep committed under current policy; pruning is separate work |
| `archive/` | Closed request/history store | Stale-safe; excluded from normal context | Keep committed; this audit listed paths only and read no bodies |
| `afk/` | Authorized autonomous-run provenance | Durable evidence | Keep committed and authority-bound |
| `debugs/` | Current debug-job contract | Temporary while open; archive after distillation | Keep as canonical plural path |
| `debug/` | Legacy singular trace path (2 tracked files) | Inconsistent with schema-2 migration contract | Judgment cleanup: classify trace before move/archive |
| `migrations.log` | Tracked root migration record | Not declared by the minimal-root architecture | Decide whether it is durable evidence and relocate or formally admit it |
| `.last-mutation` | Undo breadcrumb/session ephemera | Gitignored by policy but currently tracked | Untrack in a separately reviewed cleanup; never migrate as shared truth |
| `.DS_Store` | Operating-system metadata | Throwaway and gitignored, but currently tracked | Untrack in cleanup; no migration meaning |

Tracked-path counts at the baseline were: 217 archive files, 24 decision files, 22 snapshots, 10 roadmap files, 6 specs, 6 ideas, 6 active-request files before this request, 5 fixes, 4 memories, 2 legacy debug files, and 9 root files.

## Private local boundary

`.spectacular.local/` does not currently exist in this checkout. `.gitignore` covers the entire path. Filename-only checks found:

- no tracked `.spectacular.local/` path;
- no `.spectacular.local/` object path anywhere in reachable Git history;
- `.spectacular.local/ideas`, `.spectacular.local/github.yaml`, and `.spectacular.local/security` all resolve to the root ignore rule.

| Proposed local surface | Owner / creator | Sensitivity | Rule |
|---|---|---|---|
| `ideas/` | “Park this idea” | Private/incomplete thought | Create lazily; promote explicitly into committed `IDEA` |
| `github.yaml` | GitHub setup/activity | Account, host, fork, push-remote, expendable cache | May supplement machine operation; cannot replace canonical repository or authority |
| `security/` | Protected security workflow | Potentially confidential | Create lazily with restrictive permissions; never publish or migrate into shared storage |
| feature caches/preferences | Owning feature | Local operational state | Feature-owned, expendable unless explicitly declared otherwise |

There is no implemented general local-overrides loader in the CLI. Current code ensures the ignore rule and doctor warning/fix only. `init-workflow.md` presently claims that local state merges over shared state; D22 narrows that claim and the future contract must remove any implication that local data can override project truth.

## Leakage response

If a tracked local path is detected later:

1. Stop before reading or printing its contents.
2. Report pathnames and Git reachability only.
3. Ask the designated user/security authority to classify exposure.
4. Rotate exposed credentials before history repair when secrets may be involved.
5. Untracking, history rewriting, deletion, or disclosure each require explicit authorization and separate evidence.

## Baseline discrepancies

- The roadmap describes a future v1→v2 workspace migration, but schema 2.0 and its OKF migration are already shipped.
- The migration contract says schema 2.0 has no singular directories, yet tracked `debug/` remains.
- The architecture says the root contains anchors only, yet tracked `migrations.log`, `.last-mutation`, and `.DS_Store` remain.
- `status --against-latest` describes every unequal schema as “behind”; doctor correctly distinguishes a workspace newer than the CLI.
- No global evidence was found that ordinary mutators refuse an unknown newer schema. Schema 3 requires a fail-safe mutation guard before adoption.
