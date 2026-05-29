---
updated: 2026-05-30
---

# Verify — verify-walk

> VERIFY answers "did we build it correctly and safely?" (PLAN answers "did we build the right thing?")
> Checks are TYPED — each verified by its own authority. Walked by `spectacular verify verify-walk`.
> Dogfood note: verify-walk is verified by its own mechanism — these checks exercise the five kinds.

## Automated {run}
- [x] bash -n cli/spectacular
- [x] bash cli/spectacular verify verify-walk 2>&1 | grep -q "interactive skill flow"

## Properties {assert}
- [x] references/verify.md documents all five kinds (executable, assertable, judgable, observable, manual)
- [x] SKILL.md routes review→verified to verify.md and lists `verify <slug>`
- [x] scaffold-reference.md contains both a VERIFY.md typed stub and a VERIFY-LOG.md stub
- [x] cmd_archive warns when a verified request has no VERIFY-LOG.md

## Reads clearly {judge}
- [x] verify.md's five-kind model + two-shape syntax read coherently with no overlap
- [x] docs/commands.md verify section explains the kinds to a new user

## Manual QA {observable}
- [x] the CLI redirect message names all five kinds and points to /spectacular verify <slug>

## Actions {manual}
- [x] this VERIFY.md was itself walked end-to-end (the dogfood), producing a VERIFY-LOG entry
