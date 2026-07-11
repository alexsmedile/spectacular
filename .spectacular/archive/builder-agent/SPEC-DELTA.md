NONE — the spec already reflects this work.

`specs/index.md` (the Agent fleet capability bullet) already names both deliverables:
**`spec-builder`** ("apply-only; a *closed* milestone brief → diff, bounces on planning")
and **`build-workflow.md`** as the build-direction routing arc. The spec was written
alongside the v1.29.0+/v1.30.0 fleet work, so this request's shipped surfaces
(`agents/spec-builder.md` + `skills/spectacular/references/build-workflow.md`) are the
mechanism the spec already describes — no ADDED/MODIFIED/REMOVED is owed.

Archived on M1+M2 (both verified). M3 (durable trace + CLI emit) was P11-gated and
never earned — deferred to [[ideas/builder-trace-and-fanout]], not built, so it adds
no spec surface. M2's fan-out walkthrough is verify-on-first-real-use (mechanism
shipped + parity-checked); also tracked in that idea.
