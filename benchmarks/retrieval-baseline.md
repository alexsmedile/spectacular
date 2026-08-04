# Progressive-loading baseline

Measured on 2026-08-04 against `codex/retrieval-baseline`. This is a
measurement-only baseline for Issue #11: no router, CLI, cache, telemetry, or
workflow behavior changes are included.

Run it from the repository root:

```bash
bash tests/benchmarks/retrieval-baseline.sh
```

The harness creates a temporary, content-free workspace, exercises existing
read views, and prints byte counts rather than command output. Reference and
direct-read counts are declared in the adjacent fixture because a shell process
cannot observe model reference loading. A direct full-body read means the flow
must read the named target itself; CLI-compiled/skimmed views are not counted as
direct full-body reads.

## Baseline result

| Flow | References | CLI calls | Output bytes | Full-body reads | Repeated reads | Correct next action |
|---|---:|---:|---:|---:|---:|---|
| unknown workspace | 2 | 3 | 830 | 0 | 0 | Partial — selects active work, but does not rank it |
| named request inspection | 2 | 1 | 517 | 0 | 0 | Yes — current milestone and task |
| active implementation resume | 3 | 1 | 672 | 0 | 0 | Yes — selected milestone and execution boundary |
| simple bug fix | 2 | 1 | 162 | 0 | 0 | Yes — just-fix versus audit-first gate |
| document review | 3 | 1 | 73 | 1 | 0 | Yes — resolved target then review |
| verification | 2 | 1 | 517 | 1 | 0 | Yes — locate checks and begin gated walk |

## Specification and request document audit

| Flow | Refs | Calls | Bytes | Direct reads | Finding |
|---|---:|---:|---:|---:|---|
| Named specification | 2 | 1 | 218 | 1 | `spec list` identifies an SPC but has no compact named-spec view for state, intent, or safe next action. |
| Request inspection | 2 | 1 | 517 | 0 | `request <slug>` is already compact: frontmatter, progress, current task, and blockers. |
| Active request resume | 3 | 1 | 672 | 0 | `request <slug> --brief` already compiles only the approved implementation context; preserve this boundary. |
| Request-document review | 4 | 1 | 566 | 2 | `request <slug> --full` is appropriate only for review/debugging and correctly carries its two request-document reads. |

Implemented: the read-only `spec <id> [--json]` overview now supplies
frontmatter, lifecycle state, one-line intent, linked request, and a safe next
action in 303 bytes with no direct full-body read. Do not trim
`spec-lifecycle.md`, `request-workflow.md`,
`plan-rules.md`, or `tasks-rules.md` until that view exists and route coverage
proves the lifecycle/approval/verification gates remain discoverable.

Current ownership is already directionally clear:

- `spec-lifecycle.md` owns SPC state and approval/archive gates.
- `request-workflow.md` owns request namespace and compiled retrieval.
- `active-request.md` owns session-only exceptions and the implementation gate.
- `plan-rules.md` and `tasks-rules.md` own authoring/review structure; they are
  not startup context and should remain lazy.

## State-surface comparison

Assumption: the repository's root `AGENTS.md` is already injected by the
host, so its token cost is outside this comparison. `SKILL.md` remains the
required dispatcher whenever Spectacular is invoked. The candidate route is not
implemented; the `summary-candidate` fixture executes today's `summary --json`
as the smallest available state-view lower bound.

| State surface | Incremental references | CLI calls | Output bytes | Full-body reads | Next action |
|---|---:|---:|---:|---:|---|
| Current compact chain | 2 | 3 | 830 | 0 | Partial |
| `status --json` | 2 | 1 | 520 | 0 | Partial |
| `status --brief --json` | 2 | 1 | 731 | 0 | Yes |
| `summary --json` lower bound | 0 | 1 | 176 | 0 | No — not implemented by the current output |

This shows the decision seam, not a license to add a command or replace a
route. First decide whether `summary --json` can gain the missing facts or
whether `status --json` can become the one canonical state surface.

**Implemented recommendation:** `status --brief --json` is the compact,
action-oriented status surface. It preserves the legacy `status --json` fleet
array and adds blockers, local health signals, and a ranked next action in a
versioned envelope. It still needs the full status/doctor route when substrate
validation is required.

## Acceptance criteria for the selected state surface

- It returns, in one bounded structured response: project identity, active and
  review requests, open user blockers, selected current milestone, staleness or
  health warnings, and one ranked next action.
- Its selected request is inspectable through the existing `request <slug>` or
  `request <slug> --brief` view; it never expands request bodies by default.
- It preserves status substrate failure handling: malformed workspace state
  yields actionable doctor findings, never a fabricated brief.
- It does not replace lifecycle, approval, snapshot, or verification gates.
- It has no telemetry, cache, network call, or persisted read state.
- A regression test covers an open user blocker, a review-ready request, an
  active blocked request, and a malformed workspace before the default route
  changes.

## Surface contract

`SKILL.md` is loaded to choose the route and protect the workflow. The initial
command surface should stay small: state (`summary` or `status --json`), one
request (`request <slug>` / `--brief`), and recovery (`doctor`). Lifecycle
mutators stay discoverable through the router and CLI help; they do not belong
in every orientation payload.

The harness's output is the numeric source of truth: byte counts can change
when existing read views legitimately change. These numbers are observations,
not thresholds; a later change should explain a material delta rather than
preserve stale measurements mechanically.

## Routing findings

Unknown-workspace orientation has two incompatible prescriptions:

- `SKILL.md` calls its cold-start pattern `summary → requests --active → request <slug>`.
- `status.md` instead asks for `wayfind status --blockers-only`, config/root/index reads,
  `status --json`, memory listing, and `doctor specs`.

The benchmark uses the compact SKILL.md chain because this request measures
common retrieval. This is evidence of duplicated ownership, not permission to
remove the status safeguards; the status route owns blocker and substrate
checks that the compact chain currently does not expose.

## Duplicate-owner map

| Concern | Current owners | Baseline finding | Candidate canonical owner |
|---|---|---|---|
| SKILL router | `SKILL.md`, `doc-index.md` | Router contains trigger registry plus some workflow doctrine. | `SKILL.md` for applicability and route map; `doc-index.md` for catalog. |
| Intent routing | `SKILL.md` first-decision table, `intent-routing.md` | Same route classes appear twice; reference has the decision rules and examples. | `intent-routing.md`; router keeps one link/table row. |
| Status/orientation | `SKILL.md` cold-start pattern, `status.md`, CLI `summary`/`requests`/`status` | Two conflicting prescribed read sets. | `status.md` for no-arg safeguards and ranking; CLI for compiled facts. |
| Request workflow | `SKILL.md`, `request-workflow.md`, `active-request.md`, CLI `request` | CLI already supplies compact overview and brief; transition doctrine repeats. | `request-workflow.md`; `active-request.md` only session-state exceptions; CLI owns rendering. |
| Task layers | `SKILL.md`, `request-workflow.md`, `.spectacular/AGENTS.md` | The durable-versus-ephemeral split is explained in all three. | `.spectacular/AGENTS.md`; other files point to it. |
| Verification doctrine | `SKILL.md`, `verify.md`, `lifecycle.md`, `plan-rules.md` | The router's decision table is a useful index, while artifact selection and evidence gates belong elsewhere. | `verify.md` for walk/evidence; `lifecycle.md` for transition gate; `plan-rules.md` for 2-of-6 authoring. |

## Ranked lazy-loading candidates

1. **Task-layer explanation** — retain a one-line pointer in `SKILL.md` and
   `request-workflow.md`; load `.spectacular/AGENTS.md` only when managing
   session/milestone tracking. High duplicate reduction, no route loss.
2. **Intent-routing prose and examples** — keep the first-decision route map,
   defer the detailed receipt/examples to `intent-routing.md`. High-frequency
   core stays safe because the route classes remain visible.
3. **Status procedure detail** — preserve the no-arg route and its safety
   checks, but make `status.md` the only full procedure. Reconcile it with the
   compact CLI chain before moving anything.
4. **Verification-routing table** — keep only decision-point links in the
   router; retain 2-of-6 and evidence/transition safeguards in their existing
   owners. Do not collapse them into one document.
5. **Niche mode registries** (AFK, packs, imagine, feedback, migration) —
   candidates for an extended route index only after a route-coverage test
   proves aliases and safety gates remain discoverable.

Issue #4 remains overlap evidence only: its line-by-line trimming goal has not
started, and no content was removed or shortened here.
