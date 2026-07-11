---
status: verified
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
- [x] Add Step 1.5 (size-and-decompose gate) between Step 1 and Step 2 — lift the text from ideas/milestone-decomposition.md "The buildable shape"
  - [x] sub: write the sizing test (one coherent pass vs multi-phase = several verify-points)
  - [x] sub: write the decompose action (nested `- [ ]` sub-steps + mirror as harness TaskCreate)
  - [x] sub: write the sequential build/dispatch rule (one closed sub-brief at a time, confirm each before the next; ordered ⇒ serialize, never parallel)
  - [x] sub: write the single-phase skip path
- [x] Update the arc diagram (the fenced block) to list Step 1.5 between Step 1b and Step 2a
- [x] Update "The loop, in one line" to include the size-and-decompose beat
- [x] → check: `grep "Step 1.5|size-and-decompose|sub-step" build-workflow.md` → 16 matches; diagram + loop updated

### M2 — decompose-large-milestone @Implementation warn policy
- [x] Append `### decompose-large-milestone` under `## @Implementation` in POLICY.md (block text in the idea doc)
- [x] principle: 10 · severity: warn · with an `**Override:**` clause · NO ⛔ marker (it's a warn)
- [x] Bump POLICY.md frontmatter version (1.5 → 1.6)
- [x] → check: `spectacular policy @Implementation` lists it as `· warn`; `doctor policies` 21 blocks, 0 errors; Override present, no ⛔

### M3 — mirror + convention notes
- [x] `bug-workflow.md` Step 3: one line — a multi-phase fix decomposes the same way
- [x] `AGENTS.md` § Task tracking: a note tying nested-bullet checkpoints to the two-layer model
- [x] `tasks-rules.md`: a line blessing nested bullets as decomposition checkpoints (not only acceptance criteria)
- [x] → check: grep confirms each of the three mentions; doctor links 0 errors

### M4 — docs synced
- [x] `specs/index.md`: note the size-and-decompose gate in the arc summary (v1.10 → v1.11)
- [x] CHANGELOG `[Unreleased]` entry (plugin version bump deferred to release-time)
- [x] → check: grep specs/index.md + CHANGELOG match; `spectacular doctor` 0 errors 0 warnings

## v2 (deferred)

- [~] Live in-flight trace from a running agent — stays parked in [[ideas/builder-trace-and-fanout]] M3; only revisit if a single sub-step still runs too long to be visible
