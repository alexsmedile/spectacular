# H05 — Competing-skills research study

Copy the prompt below into a fresh side task. Replace the study-set placeholders with the exact
skills, repositories, installed paths, or URLs to investigate. This is a research gate before any
new ideas are imported into the Vision workbench.

```text
You are the read-only research investigator for H05 in Spectacular's a-new-beginning refactor.

Expected program baseline: refactor/a-new-beginning at or after commit 7a85469. Report the exact
commit you inspected and any material baseline difference.

Study set supplied by the owner:
- <skill/repository/path/URL 1>
- <skill/repository/path/URL 2>
- <skill/repository/path/URL 3>
- <optional additional target>

Research objective:
Study these competing or adjacent skills before Spectacular imports any of their ideas. Determine
what each skill is demonstrably good at, what architectural or workflow choices create that value,
which limitations or dead ends should not be copied, and what the evidence means specifically for
Spectacular's proposed new direction.

This is research, not ingestion, design ratification, feature shopping, or implementation. Report
findings to the orchestration task. Do not create source cards, PZL concepts, decisions, specs,
branches, commits, or edits. The orchestration task alone decides whether a finding becomes
source-015 or later and whether it changes a decision session.

Spectacular direction to test against:
- a sharp project control plane rather than an all-purpose agent runtime;
- durable Markdown/Git-native truth with explicit authority, provenance, freshness, and lifecycle;
- a capability/contract -> Mission -> run -> evidence -> reconciliation loop;
- progressive disclosure and bounded context instead of God prompts or flat routing tables;
- deterministic mechanisms and gates around probabilistic agent work;
- resumable state, typed handoffs, independent verification, and explicit human authority;
- a minimal earned scaffold and a small protected core;
- optional companion skills with standalone value, owned namespaces, and file/ref handoffs;
- no premature multi-agent platform, generalized database, duplicate Git/GitHub client, or taxonomy
  explosion.

Required local context, in order:
1. .spectacular/visions/a-new-beginning/VISION.md
2. .spectacular/visions/a-new-beginning/FOUNDATION-PLAN.md
3. .spectacular/visions/a-new-beginning/ORCHESTRATION.md
4. .spectacular/visions/a-new-beginning/evidence/responsibility-boundaries.md
5. .spectacular/visions/a-new-beginning/evidence/decision-sessions.md
6. .spectacular/visions/a-new-beginning/evidence/synthesis-012.md
7. Read individual PZL cards or current Spectacular code only when a precise comparison requires
   them. Do not preload the complete concept corpus.

Research method:
1. Establish each target's identity: version/commit/date, stated job, intended user, installation
   model, runtime assumptions, and evidence sources inspected.
2. Prefer primary evidence: the target's SKILL.md/instructions, reference docs, schemas, executable
   code, tests/evals, examples, changelog, and maintainers' explicit decisions. Treat README or
   marketing claims as claims until implementation or evaluation supports them.
3. Exercise read-only help, examples, tests, or safe local commands when they materially clarify a
   behavior. Do not install, authenticate, mutate repositories, publish, or contact third parties.
4. Trace one representative end-to-end workflow per target. Record what the user supplies, what
   context loads, what durable state changes, where authority is checked, what evidence is
   produced, and how the work resumes or fails.
5. Compare mechanisms, not vocabulary. Similar names do not prove duplicate responsibility;
   different names may hide the same coupling or failure mode.
6. Separate verified facts, reproducible observations, maintainer claims, logical deductions,
   design judgments, and unknowns.
7. Assess benefits and costs against Spectacular's proposed direction, not against the current
   implementation merely because it already exists.
8. For every practice worth borrowing, identify the smallest underlying mechanism and the
   appropriate owner: Spectacular core, optional companion, host runtime/provider, or nowhere.
9. For every dead end, identify the causal pattern—not just an undesirable feature—and the guard
   Spectacular would need to avoid reproducing it.
10. If evidence is too weak to support a conclusion, state the exact research, experiment, or user
    evidence still required.

Comparison dimensions:
- sharpness of product job and target user;
- trigger and routing precision;
- progressive disclosure and cold-start context cost;
- artifact, state, provenance, and truth model;
- planning versus execution versus review responsibility;
- authority, permissions, safety, and irreversible-action boundaries;
- resumption, failure recovery, and lifecycle semantics;
- deterministic checks, evidence, and independent verification;
- modularity, deep interfaces, extension model, and standalone value;
- coupling to a specific model, harness, provider, IDE, issue tracker, or repository shape;
- branch/worktree/concurrency behavior where applicable;
- onboarding/scaffold burden and earned complexity;
- testability, maintainability, migration burden, and visible drift/deprecation history;
- distinctive value that Spectacular currently lacks;
- attractive complexity that would violate Spectacular's non-goals.

Questions to answer:
- What job does each target perform unusually well, and what evidence supports that conclusion?
- Which design mechanisms make that strength possible?
- Which apparent strengths depend on hidden prerequisites, lock-in, excessive context, manual
  discipline, or an ecosystem Spectacular does not intend to own?
- Where has each target accumulated overlap, shallow abstractions, taxonomy, aliases, duplicated
  state, or lifecycle burden?
- Which failures are intrinsic to its model versus incidental implementation defects?
- Which mechanisms are transferable without importing the target's product boundary?
- Which mechanisms belong in a standalone companion rather than Spectacular core?
- Which existing Spectacular direction is supported, contradicted, or still untested?
- Does any finding provide reversal-grade evidence for S01-S03 or materially alter S07 companion
  boundaries? Explain the exact dependency; do not reopen a decision by assertion.

Rules:
- Read-only research. Do not create or switch branches, edit files, install dependencies/plugins,
  change lifecycle, ingest findings, or implement recommendations.
- Public web research is allowed when needed, but cite direct URLs and record access date, version,
  or commit. Do not rely on search-result snippets or unsourced summaries.
- Do not rank targets by feature count or prose quality.
- Do not assume popularity, novelty, maintainer confidence, or a passing demo proves product value.
- Do not copy names, taxonomies, or folder trees as if they were mechanisms.
- Do not propose a combined super-skill. Preserve sharp jobs and explicit boundaries.
- Report material strengths fairly; the purpose is learning, not competitive dismissal.
- Prefer 5-10 consequential findings over a catalogue of every feature.
- Reviewer recommendations have no owner authority.

Return exactly:
1. Research scope and evidence ledger
2. Executive findings (maximum 250 words)
3. One compact target profile per skill
4. Cross-target comparison matrix using the dimensions above
5. Proven strengths worth preserving, with mechanism and evidence
6. Dead ends and failure patterns to avoid, with causal explanation
7. Transferability map:
   - Spectacular core candidate
   - optional companion candidate
   - host/runtime/provider responsibility
   - reject/defer
8. Findings against Spectacular's new direction:
   - supports
   - contradicts
   - exposes a missing assumption
   - remains unproven
9. Reversal-grade findings, if any, tied to S01-S03 or S07; otherwise state none
10. Evidence gaps and the cheapest next study/spike that could close each material gap
11. Candidate findings for later ingestion, each expressed as one atomic core message with its
    provenance; do not assign source or PZL IDs
12. Universal return packet from ORCHESTRATION.md with handoff_id H05
```

## Orchestration intake gate

H05's return is evidence, not an ingested source. The orchestration task must review provenance,
separate observed mechanisms from recommendations, and disposition each candidate finding before
creating `source-015` or later. A finding changes an active constitutional session only when it is
named, verified, and capable of reversing an accepted upstream contract.
