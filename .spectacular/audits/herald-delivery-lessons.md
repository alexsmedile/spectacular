---
type: Audit
version: "0.1"
created: 2026-10-08T16:45:30+02:00
updated: 2026-10-08T16:45:30+02:00
title: Product delivery lessons from Herald Mail
---
# Product delivery lessons from Herald Mail

Local review and instruction changes; no Mission activation or release claim.
Case evidence is the read-only Herald workspace supplied by the owner, inspected
on 2026-10-08. Its working-client plan, ROADMAP, delivery campaign, PROJECT,
PRODUCT, ARCHITECTURE, GUARDRAILS, workspace AGENTS and V1 PRD were read.
Herald's contributor instructions are in `.spectacular/AGENTS.md`, not its root.
The review preserves its deliberate uncommitted spikes and documents.

## Findings and attribution

- **Preserve:** Synthetic receive/recovery and permission checks explicitly do
  not prove provider or human acceptance. This is successful technical execution.
- **Execution and stale context:** The delivery campaign accumulates next steps
  after grants, credentials, receiving and adapter/file spikes. Its latest pointer
  repairs execution ownership, but earlier “now/next” sections remain easy to
  misread. Ending on each technical fragment and failing to reconcile owners is
  an execution problem; the Skill previously gave too little batch guidance.
- **Instruction gap:** Existing plans linked outcomes and checks without making
  usable milestones distinct from technical tasks, or requiring experiments to
  end in a decision. The current Herald plan/roadmap now separate pilot, usable
  headless client and terminal release. These names stay case-specific.
- **Already repaired:** Herald anchors now state ordinary work, optional OS
  isolation and no Docker prerequisite. Spectacular 3.1 already defaults to
  ordinary work. No new governance rule or Mission is needed.
- **Satellite mismatch:** Architecture/prototyping handoffs still said to invoke
  Mission governance. Owner guidance unconditionally recorded write-ins through
  `decide`. Both now retain ordinary context unless governance is selected.
- **Interaction gap:** Engineering option cards lacked audience calibration;
  architecture slices lacked a minimum path to the first real workflow; prototype
  levels did not explain how early UX research coexists with later integration.

## Applied owners

- [Delivery plan](../../skills/spectacular/references/plan.md): milestones,
  bounded batches/experiments, continuation and reconciliation; checkpoint authority.
- [Requirements](../../skills/spectacular/references/requirements.md) and
  [owner questions](../../skills/spectacular/references/owner-guidance.md): product
  choices at their gate, understandable engineering defaults, ordinary write-ins.
- [Spectacular route](../../skills/spectacular/SKILL.md): pointer from ordinary
  multi-step delivery to the existing plan reference; no new mandatory record.
- [Architecture](../../skills/system-architecture/SKILL.md) and
  [method](../../skills/system-architecture/references/architecture-method.md):
  minimum real workflow, bounded validation, actual constraints and opt-in handoff.
- [Prototyping](../../skills/rapid-prototyping/SKILL.md): experiment closure,
  early disposable UX evidence and service-backed integration; opt-in handoff.

## Scenario eval: manual instruction walkthrough

These are reviewed response traces, not measured host-agent executions or
benchmark pass rates. Re-run prompts against a host for behavioral regression.

| Case and input | Reviewed expected response | Result / instruction basis |
|---|---|---|
| Milestone/batch: “Continue the fixture-backed adapter work until a resumable receive workflow works; no live mail.” | Define integrated receive/restart/read result, keep real pilot open, finish mapping/handoff/recovery checks before return. | Pass in walkthrough: plan distinguishes milestone/task and evidence; no live authorization inferred. |
| Ordinary: `.spectacular/` exists, CLI unavailable, stale unrelated Mission; request local implementation. | Use existing plan and target checks; preserve frozen bindings; no activation, Docker or isolation prerequisite. | Pass in walkthrough: kernel ordinary route and architecture handoff; missing CLI blocks command-owned actions only. |
| Spike closure: two adapters pass synthetic body tests; streaming uses less RSS; production handoff untested. | Select a candidate conditionally, bound transactional handoff validation, integrate only after checks; no provider claim. If inconclusive name missing evidence and extension bound. | Pass in walkthrough: plan/method close experiments; no automatic new A/B/C exploration for a settled decision. |
| Reconciliation: plan owns execution, roadmap owns milestones, campaign owns receipts; each says a different “next”. | Reconcile implemented/verified/open in affected existing owners; campaign points to plan; one resume action; no new STATUS catalog. | Pass in walkthrough: plan batch-boundary rule; old receipts and frozen records survive. |
| Nontechnical owner: designer asks for usable mail; storage choice is reversible; real account/privacy not selected. | Explain embedded-storage default and proceed locally; ask account/privacy only before their effects. No queue/library questionnaire. | Pass in walkthrough: requirements and owner guidance; unresolved consequential forks still require a decision. |
| UX: headless service incomplete; user requests review-language mockups. | Produce disposable labels/flow with simulated state; defer production integration to functioning service; obtain human usability evidence. | Pass in walkthrough: prototyping early UX and Level 5 rule. |
| Boundary: “continue” follows skill edits; no publication permission. | Finish local validation/report; no commit, push, release or real mailbox access. | Pass in walkthrough: plan continuity and Git authority; no checkpoint permission inferred. |

## External handoffs

External skill sources were inspected through their installed symlinks; not edited.

- **project-status — Step 3 and Resume Protocol:** always creates root STATUS and
  resumes its first unfinished checkbox. In a workspace with existing execution
  and receipt owners this can duplicate or supersede accepted context. Proposed
  owner fix: reuse the established status/receipt owner, follow its linked current
  plan and authorization, create STATUS only when useful. Eval: existing roadmap,
  campaign and plan receive updates without a fourth catalog or blind checklist
  execution. Canonical source: `project-status/SKILL.md` in its own repository.
- **update-docs — Steps 1–2:** committed-since-tag scan can return NOTHING_TO_DO
  despite an explicit request to document an uncommitted batch. Proposed owner
  fix: inspect the authorized working diff first for that request, preserve dirty
  state, patch affected owners only. Eval: uncommitted feature gets accurate docs
  without commit/push. Canonical source: `git-stack/skills/update-docs/SKILL.md`.
- **Git integration:** git-ops already separates operations and verifies effects;
  no confirmed source defect from this case. Checkpoints remain separately
  authorized. No external Git skill edit proposed.

## Verification and limits

Diff whitespace and local links/frontmatter are checked separately from behavior.
Optional context tests use temporary PyYAML, without changing project dependencies.
No Go source changed; `verify.sh` is deliberately excluded by contributor rules.
No host benchmark or independent reviewer is claimed; manual walkthroughs cannot
prove agents will follow instructions under repeated continuation or compaction.
No Herald code, frozen records, generated interface, installed copy, commit,
push or tag was modified. Pre-existing staged work remains outside this patch.
