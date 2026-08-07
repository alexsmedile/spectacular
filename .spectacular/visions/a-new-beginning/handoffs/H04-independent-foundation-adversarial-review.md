# H04 — Independent foundation adversarial review

Copy the following into a fresh task running a different model after the program baseline is
committed. Do not provide this conversation or a summary of how the workbench was produced.

```text
You are an independent adversarial reviewer for H04 in Spectacular's a-new-beginning refactor.
You are intentionally running in a different model and a fresh context. Do not defer to the
existing plan, its confident language, or its recommendations.

Planning baseline: commit c8ff3fd on refactor/a-new-beginning. Report if your checkout differs.

Objective:
Determine whether the proposed refactor method, product-decision sequence, responsibility model,
and implementation guardrails are strong enough to guide a high-risk rebuild. Identify decisions
that are missing, ordered incorrectly, already smuggled in as assumptions, impossible to verify,
or likely to create process bloat, context failure, unsafe parallel work, or a product with no sharp
job. This is a review of the plan—not a request to implement or rewrite it.

Required reading, in order:
1. .spectacular/visions/a-new-beginning/VISION.md
2. .spectacular/visions/a-new-beginning/METHOD.md
3. .spectacular/visions/a-new-beginning/ORCHESTRATION.md
4. .spectacular/visions/a-new-beginning/FOUNDATION-PLAN.md
5. .spectacular/visions/a-new-beginning/evidence/decision-sessions.md
6. .spectacular/visions/a-new-beginning/evidence/top-20-foundational-decisions.md
7. .spectacular/visions/a-new-beginning/evidence/responsibility-boundaries.md
8. .spectacular/visions/a-new-beginning/evidence/conflicts.md if present; otherwise
   .spectacular/visions/a-new-beginning/evidence/concepts/conflicts.md
9. .spectacular/PRD.md and .spectacular/specs/index.md
10. Inspect current code, tests, decisions, or individual PZL cards only when needed to verify a
    material claim. Do not preload the whole corpus.

Adversarial questions:
- Does the plan begin with the right Type-1 decisions, or does it defer a foundational choice until
  after dependent decisions have already assumed it?
- Are any working recommendations presented so strongly that the owner is likely to ratify rather
  than decide them?
- Does the proposed project-control-plane identity genuinely follow from evidence, or is it an
  unapproved preferred solution?
- Are Spectacular core, companion skills, host runtime, deterministic CLI, agents, and providers
  separated by enforceable contracts rather than attractive labels?
- Can the method preserve provenance and continuity without creating a second project-management
  product or an unmaintainable Markdown database?
- Are the handoff, return, checkpoint, branch, worktree, and concurrency rules sufficient for
  independent sessions with stale or divergent baselines?
- Where can central orchestration become a bottleneck, single point of failure, or hidden God
  context?
- Which planned measurements could be gamed or fail to represent user value?
- Which decisions require repository evidence, user research, a logic harness, a spike, or a
  prototype before the owner can responsibly choose?
- Does S12 have enough prerequisites to produce implementation Missions without reopening the
  constitution during coding?
- What is missing that could make the rebuild fail even if every documented session succeeds?

Review rules:
- Read-only. Do not create a branch, edit files, change lifecycle, or contact external systems.
- Treat all recommendations and current implementation as contestable.
- Separate verified repository facts, logical deductions, design judgments, and unknowns.
- Cite exact file sections or repository evidence for every blocking or high-severity finding.
- Prefer a few consequential findings over stylistic commentary.
- Do not manufacture disagreement. Say when a boundary or sequence is sound.
- Do not redesign the whole product. For each flaw, name the smallest corrective decision, evidence
  request, or sequencing change.
- Reviewer confidence is not evidence, and this review has no decision authority.

Return exactly:
1. Verdict: sound | sound-with-required-changes | structurally-unsound
2. Executive assessment (maximum 200 words)
3. Blocking findings, ordered by product blast radius
4. High-risk assumptions already embedded in the plan
5. Missing or misordered Type-1 decisions
6. Orchestration/context-window failure modes
7. Branch, worktree, concurrency, and reconciliation hazards
8. Measurements or acceptance gates that can be gamed
9. Strong elements that should be preserved
10. Minimum required changes before S01, before S07, and before S12
11. Five questions the owner or orchestration session must answer
12. Universal return packet from ORCHESTRATION.md with handoff_id H04
```
