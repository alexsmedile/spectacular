---
status: planned
updated: 2026-07-11
related:
  - PLAN.md
---

# Tasks — milestone-decomposition

<!-- Design pre-locked in PLAN.md + [[ideas/milestone-decomposition]]. Execution only.
     Note: M1 below is authored WITH nested-bullet sub-steps — it dogfoods the very
     checkpoint pattern this request ships. -->

## v1

### M1 — Step 1.5 size-and-decompose gate in build-workflow.md
- [ ] Add Step 1.5 (size-and-decompose gate) between Step 1 and Step 2 — lift the text from ideas/milestone-decomposition.md "The buildable shape"
  - [ ] sub: write the sizing test (one coherent pass vs multi-phase = several verify-points)
  - [ ] sub: write the decompose action (nested `- [ ]` sub-steps + mirror as harness TaskCreate)
  - [ ] sub: write the sequential build/dispatch rule (one closed sub-brief at a time, confirm each before the next; ordered ⇒ serialize, never parallel)
  - [ ] sub: write the single-phase skip path
- [ ] Update the arc diagram (the ``` block) to list Step 1.5 between Step 1 and Step 1b
- [ ] Update "The loop, in one line" to include the size-and-decompose beat
- [ ] → check: `grep -n "Step 1.5\|size-and-decompose\|sub-step" build-workflow.md` matches; diagram + loop updated

### M2 — decompose-large-milestone @Implementation warn policy
- [ ] Append `### decompose-large-milestone` under `## @Implementation` in POLICY.md (block text in the idea doc)
- [ ] principle: 10 · severity: warn · with an `**Override:**` clause · NO ⛔ marker (it's a warn)
- [ ] Bump POLICY.md frontmatter version
- [ ] → check: `spectacular policy @Implementation` lists it; `spectacular doctor policies` exits 0; block has Override, no ⛔

### M3 — mirror + convention notes
- [ ] `bug-workflow.md` Step 3: one line — a multi-phase fix decomposes the same way
- [ ] `AGENTS.md` § Task tracking: a note tying nested-bullet checkpoints to the two-layer model
- [ ] `tasks-rules.md`: a line blessing nested bullets as decomposition checkpoints (not only acceptance criteria)
- [ ] → check: grep confirms each of the three mentions

### M4 — docs synced
- [ ] `specs/index.md`: note the size-and-decompose gate in the Agent fleet / build-workflow area
- [ ] CHANGELOG entry; plugin version bump
- [ ] → check: grep specs/index.md + CHANGELOG; `spectacular doctor` exits 0

## v2 (deferred)

- [~] Live in-flight trace from a running agent — stays parked in [[ideas/builder-trace-and-fanout]] M3; only revisit if a single sub-step still runs too long to be visible
