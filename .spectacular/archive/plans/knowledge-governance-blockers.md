---
type: Reference
title: Optional governance implementation blockers
version: "0.3"
created: "2026-10-07T17:50:32Z"
updated: "2026-10-07T18:10:43Z"
description: Inspected runtime, discovery, historical policy, and command boundaries.
status: draft
---

# Optional governance implementation blockers

Findings for [the implementation plan](knowledge-first-optional-governance.md).
Source inspection completed on 2026-10-07. Runtime, index generation,
and discovery repairs are implemented; findings below retain the inspected baseline.

## Runtime restrictions

| Source | Current behavior | Repair |
|---|---|---|
| skills/spectacular/SKILL.md, Mechanical Mode | Requires startup CLI/version agreement; absent CLI routes to read/draft-only | Require compatibility only for requested governed operations |
| references/orient.md | Forbids indexes; requires Mission/fingerprint status; routes missing setup to init | Read small routing aids when useful; orient on the task without enrolling it |
| references/prepare.md | Routes obvious work to Mission or Decision; forbids handwritten timestamps broadly | Default to direct work or one soft plan; command-owned identity rules stay local |
| references/execute.md | Mandatory Mission branch; generic worker/DLQ rules unrelated to many tasks | Restrict Mission rules to explicitly governed execution; remove unrelated mandates |
| references/reduced-mode.md | Frames missing CLI as read/draft fallback | Full ordinary execution and soft writing remain available |
| references/close.md | Routine proof requires clean commit; formal review rules selected by risk | Verify actual work; no implied commit or formal record requirement |

Paths above are relative to the repository, with references under
skills/spectacular/references/. The kernel already advertises graduated governance;
these stronger unconditional rules undermine its ordinary-work branch.

## Discovery and mechanical constraints

- internal/discovery/discovery.go walks declared record roots and reads all
  collected Markdown through workspace.ReadFile. It skip-lists raw, Atlas,
  campaigns, retrospectives, and machine/history directories, but not plans.
- This repository declares record_roots as the entire .spectacular directory.
  Our Plan files are therefore discovery candidates. Source inspection shows
  type Plan is rejected by the fixed record-type parser; a plain Decision without
  id would also be rejected. No failing CLI command was executed for this audit.
- internal/workspace/document.go requires a known type and UUID identity;
  internal/domain/record.go validates governed dates and identity. Keep these
  guarantees for governed parsing; add classification before it rather than
  weakening that parser for every record.
- Discovery currently loads all records before resolving an operation. A failure
  in an unrelated record can prevent reaching the requested record. Keep the
  initial repair focused on soft-file isolation; dependency-scoped validation
  needs separate review, especially duplicate identity detection.
- internal/command/init.go creates seven governed collection directories and a
  required workspace manifest. Changing init is a public behavior change.
- test/verify.sh preflight chooses a live Mission and invokes mission check.
  Product checks remain required under AGENTS; optional knowledge governance
  must not silently weaken release verification.
- No wikilink resolver was found in the inspected internal and runtime-script
  sources. Navigation support must not replace governed reference resolution.

## Proposed classification boundary

Implemented boundary: plans join the existing soft-directory skip list. Outside
those directories, the bound project anchor, canonical CC/numbered record names,
strict governed collections, and any id/ref/human_ref/schema/schema_version key
retain governed parsing. Malformed frontmatter also retains strict parsing.
Well-formed extensible types with no authority claim are soft; unnumbered
Decision files directly under decisions/ can be soft without id or other claims.
Soft Decisions never resolve as governed typed references. Regression fixtures
exercise coexistence and damaged identity/name/schema cases.

Do not publish a new governed frontmatter template by hand. Retrieve any template
through its validating interface. Soft metadata conventions are described in
[runtime knowledge guidance](../../../skills/spectacular/references/knowledge.md), without an enforcement claim.

## Historical policy and authorization

D23 and D24 make raw non-entity, uncitable material and place typed Decisions in
the governed system. Retain those records as history; revise operative rules and
explicitly supersede the conflicting portions. D24's schema-honesty principle
survives. D30 already supports inline work; simplify its default document policy.
D10 completed-Mission freeze and D11 governed Proposal retirement protections
remain. AGENTS and .gitignore repeat raw restrictions and need coordinated edits.

The generated command catalog lists 26 commands. Requested discovery behavior
change: existing commands tolerate soft Markdown while still rejecting malformed
governed records; catalog count remains 26. The owner explicitly authorized this
public behavior change on 2026-10-07, and the conservative repair is implemented.
No new command is proposed. Retention changes, automatic raw publication, and
initializer changes are outside that authorization request.

## Follow-up repairs

Manual INDEX.md is outside discovery. Generated transaction indexes now use JSON
with explicit governed-record scope; they cannot overwrite manual navigation.
Plain unbound Anchor context is allowed. Explicit governance: context disambiguates
numbered soft Decisions only when there are no governed claims and the collection
is not strict. Raw aliases are exempt. Folder metadata/link diagnostics are optional
and distinct from Contract validation. Existing creation dates were recovered from
Git; provisional raw material was neither normalized nor published.
