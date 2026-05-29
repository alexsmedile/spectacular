---
description: Interactive validation walk — guides each VERIFY.md (or PLAN § Validation) check, records evidence, gates review→verified.
when_to_use: spectacular verify <slug>, or moving a request from review to verified.
---

# Verify — interactive validation walk

Loaded when the user runs `spectacular verify <slug>` (or asks to move a request from `review` → `verified`).

This is the **execution half** of [[verification]]. That file decides *where* verification lives (the 2-of-6 rule: standalone VERIFY.md vs folded into PLAN § Validation / TASKS § Verification). This file *runs* the checks interactively, captures evidence, and gates the lifecycle transition — turning verification from a static checklist into an executed ritual (closes PRINCIPLES.md Principle 7).

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
   - If `active` → "`<slug>` is still `active`. Move to `review` first (`spectacular promote <slug> --to review`)?"
   - If `verified` already → "`<slug>` is already verified. Re-run the walk anyway? (records a fresh VERIFY-LOG entry)"
   - If `planned` → refuse: "nothing to verify yet — `<slug>` hasn't been built."
3. **Locate the verification artifact** (see § 1).

**Substrate check (auto-invoked on failure):** if the request's PLAN/VERIFY won't parse, auto-run `spectacular doctor workspace frontmatter` and surface findings before refusing. See [[doctor-substrate]].

### 1. Locate the checks (artifact resolution)

Per [[verification]], a request's checks live in one of two shapes. Resolve in this order:

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

### 4. The gate (end of walk)

After every check is walked, summarize:

```
Walk complete — <slug>
  ✓ <P> passed   ✗ <B> blocked   ⊘ <S> skipped   (of <N> checks)
```

**If all checks passed (B = 0, S = 0):**
- Default behavior is **propose, human confirms** (Spectacular principle: agents propose, humans decide):
  > "All <N> checks passed. Mark `<slug>` verified? (y/n)"
  - On `y` → `spectacular promote <slug> --to verified`. Never edit `status:` directly.
- **Configurable auto-flip:** if `verify.auto_promote: true` in `.spectacular/config.yaml` (or `--auto` flag), skip the prompt and promote immediately, reporting it. This setting is the seam where the future [[policy-engine]] severity model plugs in.

**If any blocker or skip (B > 0 or S > 0):**
- Do **not** promote. Status stays `review`.
- Report the blocker list (with reasons) and point at the VERIFY-LOG entry.
- The request returns to `review` with a clear punch list of what's left.

### 5. Retrospective (optional, end of a passing walk)

After a successful verify (before or after the promote), optionally ask:

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

## Lifecycle + archive tie-in

- The walk is the *intended* path to `verified`. `spectacular promote <slug> --to verified` still works directly (the CLI stays a dumb mutator), but the walk is what makes the checks real.
- **Archive warning:** `spectacular archive <slug>` should warn when a request reaches archive with `verified` status but **no VERIFY-LOG.md** — i.e. it was flipped verified without ever being walked. Advisory, not blocking (see [[archive]]).

## CLI redirect

`spectacular verify <slug>` at the terminal prints:

> `verify` is a skill-side walk — it needs to read each check and judge your evidence. Run it inside Claude Code or Codex: `/spectacular` then `verify <slug>`.

## Related

- [[verification]] — *where* checks live (2-of-6 rule); this file is *how* they're walked
- [[lifecycle]] — the `review → verified` transition the walk gates
- [[archive]] — verified precondition + the un-walked warning
- [[scaffold-reference]] — VERIFY.md stub
- [[policy-engine]] — future: the `auto_promote` seam generalizes into a configurable gate policy
- [[principles]] — Principle 7 (intent → execution → validation); "agents propose, humans decide"
