NONE — the spec already reflects this work.

`specs/index.md` (the Agent fleet capability bullet) already describes the two
workflow arcs as carrying "optional judgment-gated `code-reviewer` / `test-verifier`
steps before recording" and names `repo-explorer` as the build-side discover agent —
which is exactly what this request wired in (`repo-explorer` into `build-workflow.md`
Step 0a; `code-reviewer` + `test-verifier` as optional Step 3 gates in both
`build-workflow.md` and `bug-workflow.md`). The spec was written ahead of this
request during the v1.30.0 fleet work, so no ADDED/MODIFIED/REMOVED is owed.

The deliverables here are edits to skill *reference* docs (`build-workflow.md`,
`bug-workflow.md`, `SKILL.md`, `doc-index.md`) — the mechanism the spec already
promised, now actually present. Verified: all four PLAN §Validation greps pass;
`doctor links docs` clean.
