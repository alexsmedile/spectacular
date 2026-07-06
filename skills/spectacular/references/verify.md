---
description: Interactive validation walk — guides each VERIFY.md (or PLAN § Validation) check, records evidence, gates review→verified.
when_to_use: spectacular verify <slug>, or moving a request from review to verified.
---

# Verify — the complete verification reference

Loaded when the user runs `spectacular verify <slug>` (or asks to move a request from `review` → `verified`), decides whether a request needs a standalone VERIFY.md, or sets up executable verification.

This is the single home for verification (merged from `verify.md` + `verification.md` + `verify-tests.md` in v1.20.0). Three parts:

- **Part 1 — the interactive walk** (below): *how* checks are run, evidence captured, and the lifecycle transition gated. The executed ritual that closes PRINCIPLES.md Principle 7.
- **Part 2 — the 2-of-6 rule**: *where* verification lives — standalone VERIFY.md vs folded into PLAN § Validation / TASKS § Verification.
- **Part 3 — promoting checks to scripts**: when a scenario earns a permanent `tests/verify/<slug>.test.sh` regression net.

> **Skill-only.** `verify` requires an LLM to read each check, ask for evidence, and judge pass/blocker. The CLI (`cli/spectacular`) does not dispatch `verify` — when called at the terminal it prints a friendly redirect to run inside Claude Code or Codex (same pattern as `grill`/`refine`).

## Core principle

**One check at a time. Each check verified by its own authority. No fabrication. Walk all, gate at end.**

The walk reads the request's verification artifact, iterates every check, **verifies each according to its kind**, records pass-or-blocker, and at the end reports the full picture and proposes the lifecycle transition. A check that can't be confirmed is a **blocker**, never a silent pass.

## What verification *is* — checks are typed

Verification is **not one thing**. A "check" is one of **five kinds**, and each has a different authority that decides pass/fail. This is the conceptual heart of the walk: it routes each check to the right judge instead of pretending they're all the same. The five span a spine — **deterministic → judgment → human**:

| Kind | Authority (what confirms it) | Nature | Walk behavior |
|---|---|---|---|
| **executable** | a command's **exit code** — run it, `0` = pass | deterministic — *external tool* reports | run the command |
| **assertable** | the **agent** checks a property of files/state | deterministic — *binary*, no opinion | agent evaluates true/false |
| **judgable** | the **LLM** reasons over named artifacts | fuzzy — *agent opinion*, fallible | agent judges + shows reasoning |
| **observable** | a **human** looks and confirms | opinion — *passive* (sees, doesn't act) | "confirm: <X>? (evidence)" |
| **manual** | a **human acts**, then confirms | opinion — *active* (must do something) | "perform: <action>, then confirm" |

Worked examples:
- "tests pass" → **executable** (`run: npm test`; exit code is ground truth — no opinion).
- "every spec has frontmatter" → **assertable** (agent greps/parses; true or false — no subprocess, no opinion).
- "error messages read clearly" → **judgable** (agent reads them and reasons; fallible — must show why).
- "the button is centered" → **observable** (a human looks; passive confirmation).
- "rollback restores the prior schema" → **manual** (a human must *run the rollback* first, then confirm).

The walk must not collapse these — running a test, grepping a property, an LLM's read, a human's glance, and a human performing an action are five distinct authorities.

> **The thin lines, made sharp:** *executable vs assertable* — executable trusts an **external command's** exit code; assertable is the **agent itself** checking a structural property (no subprocess). *judgable vs assertable* — judgable is **fuzzy/opinion** ("reads clearly"); assertable is **binary/mechanical** ("has frontmatter"). *observable vs manual* — observable is **passive** (the human only looks); manual is **active** (the human must *do* something before they can confirm).

### How a check declares its kind

Plain markdown, human-readable, git-diffable. The kind is an inline tag; executable checks carry the command:

```md
## Checks
- [ ] {observable} button is centered            ← human looks (default if untagged)
- [ ] tests pass `run: npm test`                 ← executable: walk runs it, exit 0 = pass
- [ ] {assert} every spec under specs/ has frontmatter   ← agent checks the property
- [ ] {judge} error messages read clearly        ← LLM reasons over the artifacts
- [ ] {manual} deploy to staging, confirm migration ran  ← human performs, then confirms
```

- **Untagged → `observable`** (the safe default — ask the human to look).
- **`` `run: <cmd>` `` → `executable`** — presence of a `run:` makes it executable regardless of other tags.
- **`{assert}` → `assertable`** — name the property + scope; the agent evaluates it mechanically (grep/parse/test), true/false.
- **`{judge}` → `judgable`** — the LLM evaluates; name the artifact(s) to read in the check text.
- **`{manual}` → `manual`** — describe the *action* the human must perform; the walk prompts them to do it, then confirm.

This makes verification's *nature* visible in the file — you can see at a glance which checks are deterministic (`run:` / `{assert}`) vs agent-judged (`{judge}`) vs human (passive `{observable}` / active `{manual}`).

### Two accepted shapes — inline OR section-grouped

A check's kind can be set **per line** (inline) or **per section** (a `##` heading tag applies to all its checks). Both are accepted.

**Prefer section-grouping when a phase is uniform.** Declaring the kind once on the heading keeps the file scannable (the phase's nature is obvious at the heading), keeps tokens tidy (the kind is read once per section, not re-parsed on every line), and groups checks into natural phases. **Reach for inline when kinds are mixed** within a group, or for a one-off check. Rule of thumb: group by phase first; drop to inline only where a group isn't uniform.

**Inline** — mixed kinds in one section, each line tagged:
```md
## Checks
- [ ] tests pass `run: npm test`
- [ ] {assert} every spec under specs/ has frontmatter
- [ ] {observable} button is centered
```

**Section-grouped** — a uniform group declares its kind once on the heading:
```md
## Automated {run}
- [ ] npm test                  ← the line text IS the command
- [ ] npm run lint

## Manual QA {observable}
- [ ] button is centered
- [ ] modal traps focus

## Rollback {manual}
- [ ] run the rollback, confirm prior schema restored
```

- `## Title {run}` → every line under it is **executable**; the **line text is the command** (no per-line `run:` needed; `## Automated {run}` + `- [ ] npm test` runs `npm test`).
- `## Title {observable|assert|judge|manual}` → every line under it is that kind (no payload — the line text is the property/prompt).
- Section-level `{run}` is the only way executable works section-wide; inline executable still uses `` `run: <cmd>` `` per line.

### Kind resolution (mechanics)

For each check, resolve kind in this order — **section is absolute**:

```
1. Enclosing ## section declares a kind ({run}/{assert}/{judge}/{observable}/{manual})?
   → use that for EVERY check in the section. Inline tags inside are IGNORED. (absolute)
2. Else parse the line's inline tag, first match wins:
   `run: <cmd>` → executable · {assert} → assertable · {judge} → judgable · {manual} → manual
3. Else → observable (default)
```

Rules that keep it predictable:
- **Section wins, absolutely.** To mix kinds, use an *untagged* section with inline tags, or split into per-kind sections. A section tag is not a default-you-can-override — it's the rule for that block.
- **One kind per check.** Exactly one (first match). Don't combine inline tags.
- **Tags stay literal in the file.** Never stripped or moved — VERIFY.md is the source; it re-parses to the same kinds every walk. What you see is what's parsed.
- **Payload comes from the text:** executable ← the `run:` command (inline) or the whole line (under `{run}` section); assertable/judgable ← the property/artifact named; observable/manual ← the line as the prompt.
- **Who sets it:** the skill proposes kinds when scaffolding VERIFY.md; you edit by hand (plain markdown); untagged = observable, so a bare checklist still walks (all human-confirmed).

## Behavior

### 0. Pre-flight

1. **Resolve the request.** `<slug>` must exist at `.spectacular/requests/<slug>/`. If it's in `archive/`, refuse ("already archived; restore first").
2. **Check status.** The walk is meant for `status: review`.
   - If `review` → proceed.
   - If `active` → "`<slug>` is still `active`. Move to `review` first (`spectacular advance <slug> --to review`)?"
   - If `verified` already → "`<slug>` is already verified. Re-run the walk anyway? (records a fresh VERIFY-LOG entry)"
   - If `planned` → refuse: "nothing to verify yet — `<slug>` hasn't been built."
3. **Locate the verification artifact** (see § 1).

**Substrate check (auto-invoked on failure):** if the request's PLAN/VERIFY won't parse, auto-run `spectacular doctor workspace frontmatter` and surface findings before refusing. See [[doctor-substrate]].

### 1. Locate the checks (artifact resolution)

Per Part 2 (the 2-of-6 rule) below, a request's checks live in one of two shapes. Resolve in this order:

1. **`requests/<slug>/VERIFY.md`** exists → walk its checkbox items (across all sections: Manual QA, Edge cases, Regression, Rollback validation).
2. **No VERIFY.md** → fall back to **`PLAN.md § Validation`** items (and `TASKS.md § Verification` if present).
3. **Neither has checkable items** → tell the user verification has no defined checks; offer to (a) scaffold a VERIFY.md now if the 2-of-6 rule warrants one, or (b) capture an ad-hoc check list for this walk.

Collect the full ordered list of checks before asking anything.

### 2. The walk loop — verify by kind

Parse each check's kind first (tag / `run:` / untagged-default-observable). Then, for each check in order, route to its authority:

**executable** (`run: <cmd>`):
```
Show the check + the exact command.
Gate (see Exec safety below): "run this? (y/n/skip)" — unless batch-allowed.
On run → exit 0 = PASS; non-zero = BLOCKER (capture stderr tail as the reason).
The exit code is ground truth. No human opinion, no LLM judgment.
```

**assertable** (`{assert}`):
```
The agent checks the named property mechanically (grep / parse / file test) — no subprocess,
no human. Evaluate true/false against the actual files/state.
true → PASS (state what was checked + the result).  false → BLOCKER (what failed).
Deterministic: same inputs, same verdict. If the property is ambiguous to evaluate,
treat it as judgable instead — don't fake a binary on a fuzzy claim.
```

**judgable** (`{judge}`):
```
Read the named artifact(s) the check points at.
Reason over them; state the judgment + the evidence the agent saw.
PASS if the agent is confident; if uncertain or the artifact is missing → BLOCKER
(never a confident pass on a guess). Surface the reasoning so the human can override.
```

**observable** (untagged or `{observable}`):
```
Ask: "How is this confirmed? (evidence / blocked / skip)"  — the human only looks.
• evidence given → record it as the human's confirmation → PASS.
  (The human is the authority here — the agent records, it does not overrule.)
• "blocked" → BLOCKER with reason.  • "skip" → SKIPPED.
```

**manual** (`{manual}`):
```
The human must PERFORM an action first, then confirm the result.
Prompt: "Perform: <action>. Done? (result / blocked / skip)"
Wait for the human to actually do it — don't accept a confirmation that skips the action.
• result reported → record action + result → PASS.
• "blocked" → BLOCKER (couldn't perform / it failed).  • "skip" → SKIPPED.
```

After each: record the outcome (§ 3), advance.

**Exec safety (gating `run:` commands).** Default: **confirm each** executable check before running (show the command, `y/n/skip`). At walk start, offer **"run all executable checks this walk? (y)"** to batch-approve — control by default, speed when the file is trusted. Mirrors the allowed-tools permission pattern. Never run a command that wasn't shown to the user.

**Never fabricate a pass.** executable → only the exit code decides. assertable → the property is true or it's a blocker (no fuzzy "close enough"). judgable → uncertainty is a blocker, not an optimistic pass. observable → "looks done" without evidence is a blocker. manual → a confirmation that skipped the action is a blocker.

**Walk all checks** even after a blocker appears — the user gets the complete picture in one pass, not a premature exit. (Mirrors the grill engine's "walk all, gate at end".)

### 3. Recording results (two artifacts)

The walk writes to **both**:

**a) VERIFY.md (or PLAN § Validation) — live state.**
- On PASS: tick the checkbox inline (`- [ ]` → `- [x]`).
- On BLOCKER / SKIP: leave the box unchecked. (The detail goes in the log, not as inline clutter.)
- This file always reflects current truth: checked = verified-this-pass.

**b) VERIFY-LOG.md — the audit trail.**
- Location: `.spectacular/requests/<slug>/VERIFY-LOG.md`.
- Append one **timestamped walk entry** per run, capturing every check + its evidence + outcome, and any blocker reasons.
- This is the durable history: who walked, when, what passed, what blocked. Never overwritten — append-only.
- Scaffold from the stub in § VERIFY-LOG shape below if it doesn't exist.

### 3b. Coherence pass (after the checks, before the gate)

The checks prove the milestones; coherence proves the *plan*. Skim the PLAN's `## Decisions` entries and confirm each shipped decision actually appears in the built artifact — a decision that silently didn't ship (chose X over Y, but the code does Y) is a finding. Grade it like a judgable check: advisory, reported in the walk summary, never a fabricated pass. One or two minutes; skip only when the Decisions section is empty.

### 4. The gate (end of walk)

After every check is walked, summarize:

```
Walk complete — <slug>
  ✓ <P> passed   ✗ <B> blocked   ⊘ <S> skipped   (of <N> checks)
```

**If all checks passed (B = 0, S = 0):**
- Default behavior is **propose, human confirms** (Spectacular principle: agents propose, humans decide):
  > "All <N> checks passed. Mark `<slug>` verified? (y/n)"
  - On `y` → `spectacular advance <slug> --to verified`. Never edit `status:` directly.
- **Configurable auto-flip:** if `verify.auto_promote: true` in `.spectacular/config.yaml` (or `--auto` flag), skip the prompt and promote immediately, reporting it. This setting is the seam where the future [[policy-engine]] severity model plugs in.

**If any blocker or skip (B > 0 or S > 0):**
- Do **not** promote. Status stays `review`.
- Report the blocker list (with reasons) and point at the VERIFY-LOG entry.
- The request returns to `review` with a clear punch list of what's left.

### 5. Retrospective (end of a passing walk)

After a successful verify (before or after the promote), ask once — skipping is fine:

> "Anything surprise you vs. the PRD/PLAN — a wrong assumption, a scope drift, a lesson? (skip to finish)"

If the user answers, route it to a `memory/` entry via `spectacular remember` (never written autonomously — the remember verb already requires the content; this just pre-fills it). This captures durable lessons at the moment they're freshest. Skipping is fine and common.

## VERIFY-LOG shape

Append-only. One entry per walk:

```md
---
updated: <today>
---

# Verify log — <slug>

## <YYYY-MM-DD HH:MM> — walk (<P> passed, <B> blocked, <S> skipped)

- ✓ [exec] <check> — `<cmd>` exit 0
- ✓ [assert] <check> — property holds: <what was checked>
- ✓ [judge] <check> — <agent reasoning + artifact seen>
- ✓ [observe] <check> — <evidence the human gave>
- ✓ [manual] <check> — performed <action>, result: <result>
- ✗ [exec] <check> — BLOCKED: `<cmd>` exit 1 — <stderr tail>
- ⊘ [observe] <check> — skipped
**Outcome:** <verified | stayed in review — N blockers>
```

The `[kind]` tag in each log line records *which authority* confirmed it — so the audit trail shows not just "passed" but *how* it was verified (a run command vs a human's eye vs agent judgment).

---

# Part 2 — Where verification lives (the 2-of-6 rule)

> Merged from the former `verification.md` (v1.20.0). This half decides *whether* a request needs a standalone `VERIFY.md` vs folding checks into PLAN/TASKS. Part 1 above is *how* the checks are walked once they exist.

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

When 2-of-6 triggers it, scaffold from the inline stub in `scaffold-reference.md § VERIFY.md`. Sections: Manual QA checklist · Edge cases to verify · Regression checklist · Rollback validation. Location: `.spectacular/requests/<slug>/VERIFY.md`.

## Anti-patterns

- **Skipping verification because VERIFY.md is "opt-in"** — opt-in refers to the *file*, not the *practice*.
- **VERIFY.md exists but is ignored** — if present, every unchecked item blocks `verified`.
- **VERIFY.md for every request** — empty noise; ignored over time.
- **VERIFY.md duplicating TASKS** — if it's just a copy of the implementation checklist, fold it back.
- **PLAN § Validation as a wishlist** — must be checkable, not aspirational.
- **No verification record at all** — `verified` without any artifact is unsupported by the lifecycle rule.

---

# Part 3 — Promoting checks to scripts

> Merged from the former `verify-tests.md` (v1.20.0). When a VERIFY scenario is worth a permanent regression net, promote it to a script.

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

---

## Lifecycle + archive tie-in

- The walk is the *intended* path to `verified`. `spectacular advance <slug> --to verified` still works directly (the CLI stays a dumb mutator), but the walk is what makes the checks real.
- **Archive warning:** `spectacular archive <slug>` should warn when a request reaches archive with `verified` status but **no VERIFY-LOG.md** — i.e. it was flipped verified without ever being walked. Advisory, not blocking (see [[archive]]).

## CLI redirect

`spectacular verify <slug>` at the terminal prints:

> `verify` is a skill-side walk — it needs to read each check and judge your evidence. Run it inside Claude Code or Codex: `/spectacular` then `verify <slug>`.

## Related

- [[lifecycle]] — the `review → verified` transition the walk gates
- [[archive]] — verified precondition + the un-walked warning
- [[scaffold-reference]] — VERIFY.md stub
- [[doctor]] — substrate self-check; catches many issues a VERIFY scenario would also catch
- [[principles]] — Principle 7 (intent → execution → validation); "agents propose, humans decide"
- `tests/run.sh` — discovers + runs all `tests/**/*.test.sh`
