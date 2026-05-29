---
description: When to author tests/verify/<slug>.test.sh vs leave manual checklists in VERIFY.md.
when_to_use: Setting up executable verification for a request.
---

# VERIFY-as-tests — When and how to promote VERIFY scenarios to scripts

Loaded when a request's VERIFY.md contains scenarios that would benefit from automation, or when designing how to capture verification for a new request.

## When to author a `tests/verify/<slug>.sh` script

Two layers cover Spectacular's verification surface:

1. **`tests/cli/*.test.sh`** — feature-level test suites (init, doctor, migrate, pack, mutator, specs, conventions). Run on every commit via `tests/run.sh`. **This is where most automated verification lives.**
2. **`tests/verify/<slug>.sh`** — request-scoped scripts that exercise the *specific* end-to-end flow a VERIFY.md scenario describes, when that flow isn't fully covered by `tests/cli/`.

**Author a `tests/verify/<slug>.sh` when:**
- A VERIFY scenario describes a specific multi-verb workflow ("init → new → promote → archive → doctor passes") that no `tests/cli/` suite exercises end-to-end
- The scenario depends on fixture state that's too specific to belong in a generic test suite (e.g. a v0.4-shape workspace with flat SCHEMA-*.md contract docs from a real consumer project)
- The scenario validates a request's exit criteria after archive — i.e. the request shipped; would a regression bring it back?

**Don't author a `tests/verify/<slug>.sh` when:**
- The mechanical content is already covered by `tests/cli/<area>.test.sh` (the common case as of v0.7.x — 199+ asserts across 7 suites cover most VERIFY-able behavior)
- The scenario requires human judgment (interactive grill walkthroughs, UX QA, source-ingestion confidence calls). Leave these in VERIFY.md as `[ ]` checkboxes, tagged "manually verified" when human-walked.

## Convention

```
tests/verify/<slug>.sh
```

- Lowercase, kebab-case, matches the request slug
- Executable bash script with `set -euo pipefail`
- Sources nothing (self-contained); seeds its own temp workspace under `/tmp/`
- Exits 0 on pass, non-zero on fail
- Cleans up its temp workspace on success

## Template

```bash
#!/usr/bin/env bash
# tests/verify/<slug>.sh — verification scenarios from .spectacular/archive/<slug>/VERIFY.md
#
# Covers: <one-line scope — e.g. "v0.4 → v0.6 migration on Octopus-shape workspace">
# Originally authored: <YYYY-MM-DD>
# Shipped in: <version>
#
# Regression intent: if this script ever fails, a v0.x change has broken
# behavior that was explicitly verified for <slug> before archive.

set -euo pipefail

REPO_ROOT="$(cd "$(dirname "$0")/../.." && pwd)"
CLI="$REPO_ROOT/cli/spectacular"
DIR="/tmp/spectacular-verify-<slug>-$$"

cleanup() { rm -rf "$DIR"; }
trap cleanup EXIT

# Seed fixture (the specific state the VERIFY scenario assumes)
mkdir -p "$DIR" && cd "$DIR"
# ... seed steps ...

# Run the scenario
"$CLI" <verb> <args>

# Assert outcomes
[[ -f .spectacular/<expected-file> ]] || { echo "FAIL: <expected-file> not created"; exit 1; }
grep -qF "<expected-content>" .spectacular/<expected-file> || { echo "FAIL: content mismatch"; exit 1; }

echo "PASS: <slug> verify scenarios"
```

## How to wire into the test runner

`tests/run.sh` already discovers `tests/**/*.test.sh` recursively via `find`. **Scripts under `tests/verify/` need a `.test.sh` suffix to be picked up** — name them `tests/verify/<slug>.test.sh` (not just `<slug>.sh`).

Alternative (if you want to exclude verify scripts from the default run): name them `<slug>.sh` and add an explicit `tests/run-verify.sh` driver that finds + runs `tests/verify/*.sh`. v1: stick with `.test.sh` suffix — single discovery path.

## Backfill policy

**Don't backfill existing archived requests** unless a regression actually surfaces. As of v0.7.x, 7 `tests/cli/*.test.sh` suites cover most of what archived VERIFY.md files would script. Backfilling everything = duplication = drift risk.

The pattern is reserved for **new requests that ship behavior not already covered by an area-level test suite** — typically multi-verb workflows or workspace-shape-specific scenarios.

## Tagging in VERIFY.md

When marking a VERIFY scenario as verified, use one of:

- `[x] mechanically verified` — auto-checkable; either covered by `tests/cli/` OR by a `tests/verify/<slug>.test.sh` script
- `[x] manually verified` — walked by a human; cannot be scripted (interactive UX, judgment calls, multi-machine validation)

This convention lets future-you / future-agents grep VERIFY.md files to know which scenarios have an automated safety net vs. which depend on someone re-walking them.

## Related

- [[verification]] — the 2-of-6 rule that decides whether a request needs VERIFY.md at all
- [[lifecycle]] — review → verified transition; tests/verify/ scripts can be part of the gate
- [[doctor]] — substrate self-check; doctor catches many issues a VERIFY scenario would also catch
- `tests/run.sh` — discovers + runs all `tests/**/*.test.sh`
