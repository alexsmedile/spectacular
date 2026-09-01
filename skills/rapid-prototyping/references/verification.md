# Proportional Verification

Use this when: a matrix round produces runnable files, rendered visuals, links, factual claims,
or a final integrated artifact. Skip this file for prose-only concept rounds.

## Choose checks from the artifact

Run the narrowest check that can falsify the claim being made.

| Artifact or claim | Minimum useful evidence |
|---|---|
| Code fragment | Parse/format plus the narrowest available lint, typecheck, or focused test |
| Responsive UI | Render at one narrow and one wide viewport; inspect clipping, hierarchy, and reflow |
| Interactive UI | Exercise the primary action, keyboard path, focus state, loading/error/empty states in scope |
| Accessibility claim | Run available semantic or automated accessibility checks and inspect labeled controls |
| Visual asset | Open the produced artifact and inspect dimensions, cropping, legibility, and required variants |
| Content or factual claim | Trace to provided sources; validate only links the option introduces |
| Final integration | Run the project's relevant build/test commands plus the visual or interaction checks above |

## Status vocabulary

- `Pass — <check>`: the named check ran successfully against this artifact.
- `Partial — <checks>`: some relevant checks ran; name the uncovered area.
- `Fail — <check>`: the check ran and found a defect.
- `Unavailable — <reason>`: a relevant tool or environment is missing.

Before using an artifact-backed status, discover the project's commands and available tools. Every
`Unavailable` status must name the commands or tools inspected, the attempted check or absence test,
and the precise missing prerequisite. Apply the same verifier set to comparable options.

In the matrix, keep evidence compact. At final delivery, report the command or inspection, its
result, and any uncovered risk. A passing command supports only the behavior it actually tests.
