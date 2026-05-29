---
description: Auto-invocation spec — when status/grill/onboarding/lifecycle silently run doctor first.
when_to_use: A skill flow needs to self-check substrate before proceeding.
---

# Doctor — skill-invoked substrate checks

Loaded when a skill flow (status, grill, onboarding, lifecycle transition) hits a substrate failure and needs to auto-invoke a scoped doctor area. Entry point: [[doctor]]; check tables: [[doctor-areas]].

## Principle

Doctor is not only a CLI subcommand. Skill flows that depend on parsing canonical artifacts (registry, kits, config, frontmatter) **must not silently proceed on parse failure** — they auto-run the relevant scoped doctor area and surface findings inline before continuing or aborting.

The skill does **not** run a full doctor sweep. It runs the relevant area(s) and surfaces results to the user with the suggestion to run `spectacular doctor --fix` (mechanical) and/or `/spectacular doctor --fix` (judgment).

## Invocation table

| Invoking flow | Doctor subset | Auto-trigger condition |
|---|---|---|
| `references/status.md` briefing | `workspace frontmatter kits` | If `config.yaml`, root doc frontmatter, or `doc-index.md` won't parse |
| `references/grill.md` pre-flight | `kits frontmatter` | If `doc-index.md`, the requested doc's `<doc>-rules.md`, or the active kit file won't parse |
| `references/onboarding.md` first-invocation | `workspace frontmatter` | First time the skill sees a workspace |
| `references/lifecycle.md` transition to `verified` | `lifecycle` (scoped to that request) | Always — verifies the verification artifact exists per [[verification]] |

## Behavior on substrate failure

1. **Run the scoped area** — invoke the CLI with the area args from the table above. Capture stdout + exit code.
2. **Surface inline** — present the findings before the briefing/refusal:
   ```
   Substrate issue detected while preparing <flow>:

   <doctor output>

   I'm pausing the <flow> until this is resolved.
   Run `spectacular doctor --fix` for mechanical repairs,
   or `/spectacular doctor --fix` to walk findings interactively.
   ```
3. **Decide: abort vs degrade** — if the failure blocks reading enough state to proceed (e.g. registry won't parse and the flow needs it), **abort**. If the flow can degrade (e.g. one root doc has a frontmatter warning but the briefing can still be built), **proceed with the inline note**.
4. **Never fabricate** — do not invent missing values, infer broken state, or proceed with partial data when the substrate is sick.

## CLI contract this depends on

`spectacular doctor <area1> <area2> <area3>` must run **all** named areas and return aggregate exit code. (Regression-tested in `tests/cli/doctor.test.sh:scenario_8b_multi_area` after 2026-05-23 bug fix.)

If a future skill flow needs a new auto-trigger, add a row to the table above and confirm the CLI supports that area combination.

## Anti-patterns

- **Silent fallback** — proceeding with empty state when the registry won't parse, then briefing the user as if everything were fine
- **Full sweep** — calling `spectacular doctor` (no args) from a skill flow. Always scope.
- **Bypassing on warnings** — substrate warnings (not just errors) still warrant a pause in flows that depend on the affected substrate
- **Auto-fixing without consent** — even if `--fix` could repair the finding mechanically, route to the user; don't `--fix` from inside another flow

## Related

- [[doctor]] — entry point
- [[doctor-areas]] — what each scoped check actually does
- [[doctor-repair]] — what happens when the user accepts the `/spectacular doctor --fix` suggestion
