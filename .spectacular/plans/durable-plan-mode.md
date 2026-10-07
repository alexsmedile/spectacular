---
type: Plan
title: Durable planning, workspace layout, and document profiles
version: '0.3'
created: '2026-10-07T18:42:01Z'
updated: "2026-10-07T19:11:33Z"
description: Planning route, container roles, and compact-default document profiles.
---

# Durable plan mode and workspace navigation

## Outcome

Use the host’s planning behavior through the Spectacular Skill and keep the
project copy in plans/. Make this workspace demonstrate small linked context.

## Approach and boundaries

- Add a Skill planning route; the public CLI stays at 26 commands.
- Honor native plan-mode restrictions. Report an unsaved draft honestly when
  the host prohibits writes; never claim a mode switch from prose.
- Reuse one plan through approach, implementation, progress, and outcome.
- Distinguish optional specs/audits from governed Contracts/Reviews/Handoffs.
- Update current anchors and Atlas guidance without rewriting historical
  Decisions, completed Missions, bound Contracts, or ignored raw captures.

## Implementation slices

1. Runtime: [plan route](../../skills/spectacular/references/plan.md) and
   [container guide](../../skills/spectacular/references/workspace.md).
2. Product docs: [active example](../../docs/human-workspace-contract.md) and
   [quickstart](../../docs/quickstart.md).
3. Self-hosting: [vocabulary](../VOCABULARY.md),
   [domain map](../atlas/domain-overview.md), and [entry points](../INDEX.md).

## Verification

Read back the saved files; inspect links and metadata with the optional context
checker. Lint and review the Skill route, including native-mode absence, forbidden
writes, repeated planning, and spec/Contract boundaries. Documentation changes
alone do not invoke this repository’s Go verification runner.

## Progress and remaining choices

The route, active-tree docs, and self-hosted context are implemented. Skill lint
reports zero errors; fresh review and issue verification pass with no blockers.
The optional checker inspected 26 maintained documents with zero findings;
Skill/product-doc file targets resolve. Temporary specs/audits fixtures passed
their agreement. The existing completed M22 still validates with its preserved
binding. No new governed artifact was manufactured to illustrate the tree.

The review also repaired existing resume, worker-timeout, fallback-path, and
external-review authorization instructions. No Go source or binary command was
changed, so the documentation-only route did not run the Go test runner.
Tag correction for v3.0.0-rc remains subject to the previously requested
publication approval.

## Document profile scale

The [folder agreement](../../skills/spectacular/knowledge-folders.yaml) defines
freeform, minimal, compact (soft default), extended, and separate governed profiles.
[Profile guidance](../../skills/spectacular/references/profiles.md) owns their
metadata, evolution, and diagnostic limits. No per-file profile label is required;
existing documents are not bulk-converted. This plan uses compact metadata.

Verification passed: Skill lint has zero errors, review/verification has no
blockers, all six existing diagnostic tests pass, and 26 maintained documents
have no findings. Minimal, compact, and extended temporary documents were accepted
without requiring optional fields; unknown domain metadata survived inspection.
File targets resolve. No Go behavior or public CLI command changed.
