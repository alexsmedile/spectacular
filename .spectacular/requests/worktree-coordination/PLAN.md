---
id: 019fdd26-9c78-77cc-9d43-8d2c7ca2f54e
slug: worktree-coordination
kind: request
scope: project
status: verified
priority: medium
owner: alex
updated: 2026-08-07
build: b47
branch: feat/worktree-coordination
base: main
docs_impact: required
docs_impact_evidence: "docs/commands.md and skills/spectacular/SKILL.md updated with workspace coordination behavior"
summary: "Implement 019fdd18-96a0-73b5-81b5-38cda98d36f4: Add durable, evidence-backed coordination before Git mutations in concurrent worktrees"
related:
  - PRD.md
references: []
# Optional traffic evidence; omit until it is confirmed and complete.
# conflicts-with: [<request-slug>]
# traffic-boundaries: [<named, complete launch boundary>]
# release-constraints: [<shared release or migration constraint>]
contract: 019fdd18-96a0-73b5-81b5-38cda98d36f4
scaffolded_against: 8b1dee1f4060899636685c9bd03e7bccb1aef00a
activated_at: 2026-08-07
activated_by: alex
activated_against: 8b1dee1f4060899636685c9bd03e7bccb1aef00a
---

# Plan — worktree-coordination

<!--
  Canonical 7-slot PLAN template for a single request.
  Lives at: .spectacular/requests/<slug>/PLAN.md

  Rules:
  - PLAN is per-request. PRD is project-wide. Never put a PRD inside requests/.
  - This file's frontmatter `status:` is the single source of lifecycle state for the request.
  - The 7 required sections must appear IN ORDER, unnumbered:
      ## Goal, ## Constraints, ## Milestones, ## Tasks, ## Dependencies, ## Validation, ## Deliverables
    Extra sections (## Understanding, ## Decisions, request-specific) may appear
    BETWEEN them; doctor enforces the required set's presence + order, not a closed list.
  - All 7 required sections must be filled before this PLAN is considered usable.
  - Replace every <placeholder> with concrete content.
-->

## Goal

Give agents and maintainers deterministic, evidence-backed coordination before
Spectacular changes Git state in a concurrent or dirty worktree.
## Constraints

- Stay within the approved 019fdd18-96a0-73b5-81b5-38cda98d36f4 requirements; unresolved scope changes return to discovery.
- Preserve existing behavior outside this specification.
- Production code and invariant tests become implementation truth after verification.
- Preserve unrelated work; no reset, stash, discard, overwrite, or bulk staging
  is permitted as a recovery shortcut.
- Keep Bash 3.2 compatibility and support local-only repositories without a
  forge adapter or remote.
## Understanding

<!--
  OPTIONAL authoring slot, but REQUIRED before `planned → active` by the
  `understand-before-change` policy (@Implementation). Fill it here for a
  typical request; escalate to a dedicated requests/<slug>/UNDERSTANDING.md
  (same three subheads) for large ones — the policy is satisfied by EITHER.
  Not one of the 7 required authoring slots; it gates implementation, not planning.
-->

### How it works now

`afk` has isolated branch and cleanup helpers, but normal interactive commands
either require a clean checkout ad hoc or mutate without a shared path-level
ownership model. Requests know their contract but do not yet declare branch/base
intent for later sessions.

### What changes

The CLI gains `workspace preflight|plan|preserve|cleanup`, shared Git evidence
helpers, request branch/base frontmatter support, and common guards on
Spectacular-owned Git mutations. Skills/docs and regression tests describe the
provider-neutral contract.

### What stays the same

Git remains authoritative for ancestry. Forge APIs remain optional adapter
evidence. No per-file ownership database, mandatory network fetch, automatic
stash, remote branch deletion, or automatic merge is introduced.

## Decisions

<!--
  Design calls made inside this request. Format: chose X over Y — because Z.
  Rejected alternatives stay listed; deleting them re-litigates them later.
  Project-wide calls go to DECISIONS.md via `spectacular decide` instead
  (see decisions-rules.md routing table). Empty is fine — no decisions yet.
-->

- Compute path ownership from live evidence instead of persisting it — because
  durable file-by-file records become stale and mislead later sessions.
- Preserve unrelated paths with a named scoped branch/commit rather than stash
  — because commits are inspectable, shareable, and recoverable.
- Reuse AFK Git primitives where they satisfy the safety contract, but expose a
  separate interactive `workspace` namespace — because AFK authorization and
  ordinary concurrent-work coordination have different entry conditions.

## Milestones

- M1 — A deterministic local Git/request evidence model classifies every changed
  path without writes and understands declared request branch/base metadata.
- M2 — `workspace preflight` and `workspace plan` present the required report,
  blockers, exact prospective commands, and one safe next action.
- M3 — Scoped preservation creates a named branch and commit for selected paths
  without absorbing unrelated work.
- M4 — Cleanup inventory and apply path distinguish merged, open, declined, and
  unknown branches and enforce reachability/fresh-base deletion gates.
- M5 — Existing owned Git mutators run the preflight, docs/doctor cover the new
  metadata, and regression/compatibility checks pass.
## Tasks

<!-- Pointer. The executable checklist lives in TASKS.md, grouped by milestone. -->

See `TASKS.md`.

## Dependencies

- Approved source specification: 019fdd18-96a0-73b5-81b5-38cda98d36f4.
- Existing AFK Git-hygiene helpers and tests are the safety baseline.
- Linked decisions and questions remain in the source specification.
## Validation

- M1 — focused CLI scenarios prove request branch/base metadata and local Git
  evidence remain read-only and conservative.
- M2 — focused CLI scenarios prove zero writes in dirty/mixed worktrees and
  deterministic staged/unstaged/untracked path dispositions.
- M3 — preservation scenarios prove only declared paths are committed and no
  stash/reset/discard/overwrite occurs.
- M4 — cleanup scenarios prove stale-base, open-PR, declined, and no-provider
  states take the conservative disposition.
- M5 — `bash tests/run.sh`, `bash -n cli/spectacular cli/install.sh
  scripts/hooks/pre-commit`, `scripts/hooks/pre-commit --check`, and relevant
  `spectacular doctor` checks pass.
## Deliverables

- Provider-neutral workspace coordination CLI and common safety guard.
- Request frontmatter/doctor/skill documentation for branch and base evidence.
- Focused regression coverage plus verification evidence and docs-impact assessment.
