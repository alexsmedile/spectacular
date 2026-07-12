---
description: Authoring-time verification decisions — the 2-of-6 rule (standalone VERIFY.md vs folded checks), fold patterns, VERIFY.md shape, and promoting scenarios to test scripts.
when_to_use: Scaffolding or grilling a request (does it need a VERIFY.md?), or deciding whether a verify scenario earns a permanent tests/verify/ script. NOT for running the walk — that's verify.md.
---

# Verify authoring — where checks live + when they become scripts

Loaded at **authoring time**: scaffolding a request, grilling a PLAN, or deciding a scenario's regression future. The **walk** — running the checks at `review → verified` — lives in [[verify]]; this doc never loads for it.

> History: split out of verify.md in v1.34 (b30). verify.md had merged three docs in v1.20.0 (`verify.md` + `verification.md` + `verify-tests.md`); this file carries the former Parts 2–3 (from `verification.md` / `verify-tests.md`), so the `review → verified` gate loads only walk content. The compact 2-of-6 table in [[plan-rules]] § 2-of-6 stays the quick-reference copy; **this doc is the canonical source**.

# Part 2 — Where verification lives (the 2-of-6 rule)

> **@Verification policy gate.** Before moving a request `review → verified`, run `spectacular policy @Verification` and follow every active policy. The default blocker is `verification-present`: every check in VERIFY.md (or PLAN § Validation) must be satisfied — or you stop. See [policy-injection.md](policy-injection.md).

## Verification always happens; the file is opt-in

To be precise about "opt-in" — it refers to **whether you scaffold a `VERIFY.md` file**, not whether you verify. Every request runs verification before reaching `verified`; the only question is *which artifact carries the checks*.

**Never skip verification because VERIFY.md is "optional."** If VERIFY.md exists, it is load-bearing — every unchecked item blocks the `verified` transition. If it doesn't exist, the checks live in PLAN § Validation or TASKS § Verification, and *those* block the transition. There is no path to `verified` without explicit verification of some artifact.

| Check type | Question it answers | When it runs | Natural home |
|---|---|---|---|
| Task completion | Did I do the work? | During `active` | `TASKS.md` checkboxes |
| Acceptance criteria | Did we build the right thing? | At `review` entry | `PLAN.md` § Validation |
| QA / risk verification | Did we build it correctly + safely? Will it break? | During `review` | `VERIFY.md` (when needed) |

The first two are universal. The third is **conditional** — it earns a standalone file only when risk and surface justify it.

## The 2-of-6 rule

Scaffold a standalone `VERIFY.md` when **at least two** are true:

1. **User-visible change** — anyone outside the team can observe behavior change.
2. **High reversibility cost** — migrations, schema changes, destructive operations.
3. **Multi-surface verification** — more than running tests: manual QA, browser checks, prod smoke.
4. **Risk surface non-trivial** — auth, billing, payments, PII, security, data integrity.
5. **External contract change** — public API, exported library shape, plugin interface, CLI flags.
6. **Rollback plan exists** — a non-trivial undo procedure that itself needs verification.

The common thread: checks that can't be expressed as `- [x]` next to a task because they're multi-step, time-sequenced, or require human judgment in a specific environment.

**Default to no file.** A new file must earn its keep. "No file" ≠ "no verification" — the checks still exist in PLAN or TASKS.

## When to fold into PLAN/TASKS instead

VERIFY.md is overkill — and the noise hurts — for: doc-only changes; internal refactors with test coverage; build/tooling changes; spec/template additions (dogfood test); config changes; anything single-step.

**Folded pattern A — PLAN § Validation** (preferred for small requests): per-milestone validation criteria in PLAN.md slot 6. Best when validation is descriptive ("this property holds").

```md
## 6. Validation
- M1 — `references/kits-contract.md` exists with full schema documented
- M3 — Dogfood test: coding kit produces 10-slot merged sequence with `kit: coding` frontmatter
```

**Folded pattern B — TASKS § Verification** (preferred when checks are step-by-step): a `### Verification` group at the end of TASKS.md, same checklist syntax. Closest to a "lite VERIFY.md" without the separate file.

## Standalone VERIFY.md shape

When 2-of-6 triggers it, scaffold from the inline stub in `scaffold-reference.md § VERIFY.md`. Sections: Manual QA checklist · Edge cases to verify · Regression checklist · Rollback validation. Location: `.spectacular/requests/<slug>/VERIFY.md`. Check-kind tags (`{assert}` / `{judge}` / `{manual}` / `` `run:` `` / section-level `{run}`) are documented in [[verify]] — propose kinds when scaffolding so the walk routes each check to the right authority.

## Anti-patterns

- **Skipping verification because VERIFY.md is "opt-in"** — opt-in refers to the *file*, not the *practice*.
- **VERIFY.md exists but is ignored** — if present, every unchecked item blocks `verified`.
- **VERIFY.md for every request** — empty noise; ignored over time.
- **VERIFY.md duplicating TASKS** — if it's just a copy of the implementation checklist, fold it back.
- **PLAN § Validation as a wishlist** — must be checkable, not aspirational.
- **No verification record at all** — `verified` without any artifact is unsupported by the lifecycle rule.

# Part 3 — Promoting checks to scripts

When a VERIFY scenario is worth a permanent regression net, promote it to a script.

## When to author `tests/verify/<slug>.test.sh`

Two layers cover Spectacular's verification surface:

1. **`tests/cli/*.test.sh`** — feature-level suites (init, doctor, migrate, pack, mutator, specs, conventions), run on every commit via `tests/run.sh`. **Most automated verification lives here.**
2. **`tests/verify/<slug>.test.sh`** — request-scoped scripts exercising the *specific* end-to-end flow a VERIFY.md scenario describes, when `tests/cli/` doesn't cover it.

**Author one when:** a VERIFY scenario describes a multi-verb workflow (`init → new → advance → archive → doctor passes`) no `tests/cli/` suite exercises; or it depends on fixture state too specific for a generic suite; or it validates a request's exit criteria after archive (would a regression bring it back?).

**Don't when:** the mechanical content is already covered by `tests/cli/<area>.test.sh` (the common case); or the scenario requires human judgment (interactive grill walkthroughs, UX QA) — leave those in VERIFY.md as `[ ]`, tagged "manually verified".

## Convention + wiring

- `tests/verify/<slug>.test.sh` — lowercase kebab-case matching the slug; `set -euo pipefail`; self-contained (seeds its own `/tmp/` workspace); exits 0 on pass; cleans up on success.
- `tests/run.sh` discovers `tests/**/*.test.sh` recursively — the **`.test.sh` suffix is required** for pickup (not just `<slug>.sh`).

```bash
#!/usr/bin/env bash
# tests/verify/<slug>.test.sh — scenarios from .spectacular/archive/<slug>/VERIFY.md
# Regression intent: if this fails, a change broke behavior verified for <slug> before archive.
set -euo pipefail
REPO_ROOT="$(cd "$(dirname "$0")/../.." && pwd)"; CLI="$REPO_ROOT/cli/spectacular"
DIR="/tmp/spectacular-verify-<slug>-$$"; trap 'rm -rf "$DIR"' EXIT
mkdir -p "$DIR" && cd "$DIR"
# ... seed fixture, run scenario ...
"$CLI" <verb> <args>
[[ -f .spectacular/<expected-file> ]] || { echo "FAIL: not created"; exit 1; }
echo "PASS: <slug> verify scenarios"
```

## Backfill policy + tagging

**Don't backfill archived requests** unless a regression surfaces — `tests/cli/` already covers most of what archived VERIFY.md files would script; backfilling everything = drift risk. Reserve the pattern for **new requests shipping behavior not already covered** by an area-level suite.

When marking a VERIFY scenario verified, tag it: `[x] mechanically verified` (covered by `tests/cli/` or a `tests/verify/<slug>.test.sh`) or `[x] manually verified` (human-walked, can't be scripted). Lets future agents grep which scenarios have an automated safety net.

## Related

- [[verify]] — the interactive walk these checks feed (Part 1; loads at `review → verified`)
- [[plan-rules]] § 2-of-6 rule — the compact quick-reference copy used at scaffold/grill time
- [[scaffold-reference]] — VERIFY.md stub
- [[lifecycle]] — the transitions the artifacts gate
- `tests/run.sh` — discovers + runs all `tests/**/*.test.sh`
