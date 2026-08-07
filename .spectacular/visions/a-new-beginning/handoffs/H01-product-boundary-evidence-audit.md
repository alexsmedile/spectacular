# H01 — Product-boundary evidence audit

Copy the following into a new side task after the program baseline is committed:

```text
You are the read-only evidence investigator for H01 in Spectacular's a-new-beginning refactor.

Objective:
Produce a compact factual packet for S01 Product Constitution. Determine what the current product
claims to be, what it demonstrably does, where those differ, which behaviors appear load-bearing,
and which current responsibilities contradict the inherited non-goals. Do not decide the future
product.

Required reading, in order:
1. .spectacular/visions/a-new-beginning/ORCHESTRATION.md
2. .spectacular/visions/a-new-beginning/VISION.md
3. .spectacular/PRD.md
4. .spectacular/specs/index.md
5. .spectacular/visions/a-new-beginning/evidence/synthesis-012.md
6. These concept cards only: PZL-012, PZL-013, PZL-019, PZL-020, PZL-038, PZL-062, PZL-064
7. Follow precise links into code/tests/decisions only when needed to verify a material claim.

Questions:
- What job and primary users does the current PRD claim?
- What responsibilities does the shipped capability index actually assign to Spectacular?
- Which three to seven behaviors have the strongest current evidence of unique value?
- Which capabilities contradict the current non-goals or one-time-bootstrap CLI claim?
- Which constraints are genuinely load-bearing versus inherited and reopenable?
- What evidence is missing before choosing among control plane, integrated executor, or
  knowledge/specification workspace?

Rules:
- Read-only: do not create a branch, edit files, change lifecycle, or contact external systems.
- Separate verified facts, document claims, inferences, and recommendations.
- Do not preload all 171 concept cards or archive history.
- Do not turn current implementation into permission to preserve it.
- If a repository claim cannot be verified cheaply, label it unknown.

Return exactly:
1. Executive finding (maximum 150 words)
2. Claimed product versus actual responsibility table
3. Load-bearing behavior candidates with evidence refs
4. Product-boundary contradictions
5. Inherited constraints requiring confirmation
6. Missing evidence that could change S01
7. One recommended framing question for the owner
8. Universal return packet from ORCHESTRATION.md with handoff_id H01
```
