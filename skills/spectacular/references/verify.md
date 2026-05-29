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

Verification is **not one thing**. A "check" is one of three kinds, and each has a different authority that decides pass/fail. This is the conceptual heart of the walk: it routes each check to the right judge instead of pretending they're all the same.

| Kind | Authority (what confirms it) | Nature |
|---|---|---|
| **executable** | a command's **exit code** — run it, `0` = pass | deterministic / "mathematical" — no opinion |
| **observable** | the **human** provides evidence; the agent records it | opinionated — a human is the judge |
| **judgable** | the **LLM** reasons over named artifacts and decides | semi-opinionated — agent judgment, fallible |

So: "tests pass" is *executable* (run `npm test`, exit code is ground truth — there is no opinion). "The button is centered" is *observable* (a human looks). "Error messages read clearly" is *judgable* (the agent reads them and reasons). The walk must not collapse these — running a test is not the same authority as a human eyeballing a layout.

### How a check declares its kind

Plain markdown, human-readable, git-diffable. The kind is an inline tag; executable checks carry the command:

```md
## Manual QA
- [ ] {observable} button is centered          ← human evidence (default if untagged)
- [ ] tests pass `run: npm test`               ← executable: walk runs it, exit 0 = pass
- [ ] {judge} error messages read clearly      ← LLM reasons over the artifacts
```

- **Untagged → `observable`** (the safe default — ask the human).
- **`` `run: <cmd>` `` → `executable`** — presence of a `run:` makes it executable regardless of other tags.
- **`{judge}` → `judgable`** — the LLM evaluates; name the artifact(s) to read in the check text.

This makes verification's *nature* visible in the file — you can see at a glance which checks are deterministic (`run:`) vs opinionated (`{observable}`) vs agent-judged (`{judge}`).

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

**observable** (untagged or `{observable}`):
```
Ask: "How is this confirmed? (evidence / blocked / skip)"
• evidence given → record it as the human's confirmation → PASS.
  (The human is the authority here — the agent records, it does not overrule.)
• "blocked" → BLOCKER with reason.  • "skip" → SKIPPED.
```

**judgable** (`{judge}`):
```
Read the named artifact(s) the check points at.
Reason over them; state the judgment + the evidence the agent saw.
PASS if the agent is confident; if uncertain or the artifact is missing → BLOCKER
(never a confident pass on a guess). Surface the reasoning so the human can override.
```

After each: record the outcome (§ 3), advance.

**Exec safety (gating `run:` commands).** Default: **confirm each** executable check before running (show the command, `y/n/skip`). At walk start, offer **"run all executable checks this walk? (y)"** to batch-approve — control by default, speed when the file is trusted. Mirrors the allowed-tools permission pattern. Never run a command that wasn't shown to the user.

**Never fabricate a pass.** For observable, "looks done" without evidence is a blocker. For judgable, uncertainty is a blocker, not an optimistic pass. For executable, only the exit code decides.

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
- ✓ [observe] <check> — <evidence the human gave>
- ✓ [judge] <check> — <agent reasoning + artifact seen>
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
