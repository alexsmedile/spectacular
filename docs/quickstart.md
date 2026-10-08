# Quickstart

Start from the context you have and do the work. Spectacular keeps that context
in small linked Markdown documents; a Mission is an explicitly chosen option.
Ordinary work needs neither the CLI nor a workspace manifest.

## Start with one useful document

Use an existing anchor or raw draft if it already explains the job. Otherwise,
write a short plan with the intended result, relevant constraints, an approach,
and how to check the result. Keep progress and results in that same plan.

Create `.spectacular/PROJECT.md` when you need durable project context. Additional
anchors and folders appear when content needs them; no starter collection of
empty directories is required for the ordinary workflow.

## Plan with Spectacular

Ask the Skill: `/spectacular plan improve the release flow`, or
`$spectacular plan improve the release flow` in hosts using that invocation.
It uses the host’s planning behavior and keeps the project copy at
`.spectacular/plans/release-flow.md`. This is a Skill route, not a binary command.

The host controls native plan mode, restrictions, and approvals. Spectacular
cannot toggle a mode through prose. When writes are prohibited, it returns an
unsaved draft and intended path; saving waits until the host permits it.
Repeated planning updates the same file. Implementation can later track progress,
checks, and outcomes there without creating a Mission.

## Choose a home by purpose

| Home | What belongs there |
|---|---|
| Root anchors | Current purpose, accepted architecture, vocabulary, constraints |
| `raw/` | Captures, research, unfinished drafts, quick-start material |
| `atlas/` | Domain maps and lasting explanations |
| `decisions/` | Consequential choices and why they were made |
| `plans/` | Intended work, approach, progress, and outcome |
| `requirements/` | Optional linked needs, constraints, and acceptance criteria |
| `specs/` | Optional reusable behavior and interfaces |
| `scenarios/` | Optional reusable user journeys, final outputs, and success criteria |
| `audits/` | Optional revision-scoped inspection findings |

See the [active-workspace example](human-workspace-contract.md) for Mission-owned
reviews, evidence, handoffs, and optional supporting folders.

A plan can include requirements and design. There is no required Proposal,
Contract, or Mission before implementation. Agents choose placement and evolve
ordinary documents within the authorized task; they ask when an accepted
constraint conflicts or a consequential action needs authorization.

Raw material remains provisional. Reading it does not authorize publishing it;
existing ignored captures stay ignored unless you choose a retention change.

## Keep documents small and connected

Prefer one question or coherent outcome per document, usually 50–120 lines.
Split by domain, code module, working phase, or delivered tranche when each
part has independent readers or reuse. Link it from
the plan rather than copying it. Avoid fragments that cannot stand on their own.

Use relative Markdown links or Obsidian wikilinks, with section targets and
aliases where useful. Explicit paths avoid ambiguous filenames. Agents follow
and maintain those links; the CLI does not currently supply a wikilink resolver.
Public documentation uses portable Markdown links.

Maintained knowledge carries foundational `type` metadata plus document `version`,
`created`, and `updated`. Compact is the default: add a short description when
useful, with title optional. Minimal remains valid; extended retrieval fields
are added only when needed. Raw is freeform. See the
[metadata scale](human-workspace-contract.md#metadata-scale). Preserve creation
time and unknown fields; omit unknown historical dates rather than inventing
one. Document versions track meaningful content changes; Git preserves edit
history. Missing optional metadata does not stop implementation.

## Do the work and verify it

An agent reads relevant anchors and linked context, implements the authorized
change, runs appropriate product checks, and records useful results. An unrelated
Mission or a missing Spectacular binary does not enroll or stop ordinary work.
Routine completion does not require a commit or a formal evidence package.

Soft Decisions use descriptive unnumbered filenames. Existing numbered Decisions
and other governed records retain their enforced identities and mutation rules.
Do not convert historical records merely to adopt the ordinary workflow.

## Choose governed execution when useful

Select a Mission explicitly when freezing an agreement and collecting formal
proof help the task. Governed operations need the compatible CLI and retain
activation, integrity, and completion rules. See [Process](process.md) for that
advanced workflow and [Installation](installation.md) for the CLI.

The current initializer creates a governed workspace and several collections.
It has not been changed to the minimal ordinary-work layout described above.
