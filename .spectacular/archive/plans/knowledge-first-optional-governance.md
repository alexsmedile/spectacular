---
type: Plan
title: Knowledge-first workspace and optional governance
version: "0.6"
created: "2026-10-07T17:50:32Z"
updated: "2026-10-07T18:58:25Z"
description: Implement small linked knowledge documents and explicitly optional governance.
status: implemented
---

# Knowledge-first workspace and optional governance

Keep durable context as small linked OKF objects. Make raw material, Decisions, anchors, and Atlas the everyday center of
Spectacular. Start work directly from context or one small plan. Keep formal
execution governance available by explicit choice, with its guarantees intact.

The owner approved this direction on 2026-10-07 and requested implementation.
No Mission is activated by this plan. Public CLI changes still require the
explicit authorization specified in AGENTS.md.

## Read only what the job needs

- [Runtime folder roles, small documents, metadata, and links](../../../skills/spectacular/references/knowledge.md)
- [Mechanical blockers and proposed repairs](knowledge-governance-blockers.md)

The [durable planning update](durable-plan-mode.md) owns the subsequent Skill
route and self-hosted layout changes.

Runtime rules live in the skill; this plan links them instead of keeping a second
copy. The inventory holds bounded implementation findings. Progress stays here.

## Implementation order

1. **Inventory blockers — completed.** Runtime and discovery findings are in the
   linked inventory. Historical constraints and command boundaries are identified.
2. **Change runtime guidance — implemented.** Route ordinary work to relevant linked knowledge
   without CLI startup, Mission inference, or formal closure. Preserve explicit
   governed routing. Update only affected skill branches and review the result.
3. **Repair discovery — authorized and implemented.** Separate soft documents from
   governed records; allow plans and soft Decisions without poisoning unrelated
   commands. Malformed governed records remain errors. Catalog stays at 26 commands.
4. **Align operative rules and documentation — initial alignment implemented.** Update repository-owned AGENTS
   guidance, README, quickstart, workspace contract, and affected explanations.
   Preserve generated harness sections and historical Decisions. Resolve old
   rules through explicit supersession rather than rewriting accepted history.
5. **Validate first implementation — passed.** Exercise ordinary work without CLI or Mission;
   coexistence of soft/governed Decisions; malformed governed records; unrelated
   stale Missions; metadata preservation; link resolution and repair on moves.

## Implementation boundaries

- Retain existing governed record identities, formats, and historical bindings.
- No automatic conversion, archival, deletion, or publication of existing files.
- No v1 readers, migrations, second package root, or new public commands.
- Keep existing raw content ignored until a separate retention review. Permission
  to use a raw source is independent of permission to publish it.
- Review initializer behavior only after explicit authorization for that change.
- Keep product verification tied to actual changes, independent of optional
  knowledge checks. Markdown-only edits do not run `verify.sh`.
- After Go changes, run `quick`. Repair `preflight` before heavy tiers; use
  acceptance for changed boundaries and `all` at release/Mission completion.

## Completion criteria

Agents can complete ordinary work from raw or one plan without the binary,
manifest, fingerprints, or lifecycle gates. They maintain small linked documents
without duplicating truth. Explicit governed operations retain integrity checks.
Existing content survives, and proposed text is distinguishable from accepted
constraints. Report actual verification and remaining limitations here.

## Current state

- Runtime defaults to durable context and direct work. Contracts retain mandatory
  mechanical validation and amendment/version rules; Decision validation is
  recommended, not required. Raw/sketch/scratchpad have no content obligations.
- Folder agreements are centralized in the skill's knowledge-folders.yaml.
  The optional read-only checker validates maintained metadata and file link
  targets. It skips raw and governed records; it does not certify section anchors.
- Plain Anchor context and explicitly marked numbered soft Decisions no longer
  poison discovery. Identity/schema claims and strict collections remain governed.
- Manual INDEX.md is safely excluded from governed discovery. Transactions now
  generate scoped index.json inventories and never generate Markdown indexes.
  The root guide is manual; the case-only rename is recorded in Git.
- Maintained context metadata was normalized with dates recovered from Git history.
  ONTOLOGY.md is the preferred domain-anchor name; an existing ONTOLOGY.md
  remains supported without duplicating authority. Product docs remain owned by documentation tools.
- Context diagnostics: 23 documents, zero findings. Focused discovery/index tests,
  full race/acceptance/reproducibility/install/recovery release gate and six
  diagnostic tests passed for 3.0.0-rc. The official
  tokenizer cache is temporary. Skill lint: zero errors, five advisories.
- Review findings were repaired, including fallback context isolation and unsafe
  governed templates/splitting. A final issue-scoped review caught an unsupported
  Anchor-schema instruction; its repair passed an issue-scoped verification review.

Initializer behavior and historical governed records remain unchanged. Legacy
collection Markdown indexes remain historical projections and are no longer
regenerated. The public workflow docs were aligned for v3. Automatic section-anchor
validation remains separate work; no new public command was added (26 commands).
