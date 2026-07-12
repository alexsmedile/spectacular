---
description: Interactive validation walk — guides each VERIFY.md (or PLAN § Validation) check, records evidence, gates review→verified.
when_to_use: spectacular verify <slug>, or moving a request from review to verified.
---

# Verify — the interactive walk

Loaded when the user runs `spectacular verify <slug>` or asks to move a request `review → verified`. **Walk-only**: how checks run, evidence is captured, and the transition is gated — closing Principle 7. Authoring-time decisions — the 2-of-6 rule, fold patterns, promoting scenarios to `tests/verify/` scripts — live in [[verify-authoring]] (split v1.34, b30; quick 2-of-6 table in [[plan-rules]]).

> **Skill-only.** `verify` requires an LLM to read each check, ask for evidence, and judge pass/blocker. The CLI does not dispatch `verify` — at the terminal it prints a redirect to run inside Claude Code or Codex (same pattern as `grill`/`refine`).

## Core principle

**One check at a time. Each check verified by its own authority. No fabrication. Walk all, gate at end.**

The walk reads the request's verification artifact, iterates every check, **verifies each according to its kind**, records pass-or-blocker, and at the end reports the full picture and proposes the lifecycle transition. A check that can't be confirmed is a **blocker**, never a silent pass.

## What verification *is* — checks are typed

A "check" is one of **five kinds**, each with a different authority that decides pass/fail. This is the conceptual heart of the walk: it routes each check to the right judge. The five span a spine — **deterministic → judgment → human**:

| Kind | Authority (what confirms it) | Nature | Example |
|---|---|---|---|
| **executable** | a command's **exit code** — run it, `0` = pass | deterministic — *external tool* reports | "tests pass" `run: npm test` |
| **assertable** | the **agent** checks a property of files/state | deterministic — *binary*, no opinion | "every spec has frontmatter" |
| **judgable** | the **LLM** reasons over named artifacts | fuzzy — *agent opinion*, fallible | "error messages read clearly" |
| **observable** | a **human** looks and confirms | opinion — *passive* (sees, doesn't act) | "the button is centered" |
| **manual** | a **human acts**, then confirms | opinion — *active* (must do something) | "run the rollback, confirm schema restored" |

> **The thin lines, made sharp:** *executable vs assertable* — executable trusts an **external command's** exit code; assertable is the **agent itself** checking a structural property (no subprocess). *judgable vs assertable* — judgable is **fuzzy/opinion**; assertable is **binary/mechanical**. *observable vs manual* — observable is **passive** (the human only looks); manual is **active** (the human must *do* something first). The walk must not collapse these five authorities.

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
- **`` `run: <cmd>` `` → `executable`** — a `run:` makes it executable regardless of other tags.
- **`{assert}` → `assertable`** — name the property + scope; the agent evaluates true/false.
- **`{judge}` → `judgable`** — name the artifact(s) to read in the check text.
- **`{manual}` → `manual`** — describe the *action* the human must perform, then confirm.

### Two accepted shapes — inline OR section-grouped

A check's kind can be set **per line** (inline, as above) or **per section** — a `##` heading tag applies to all its checks. **Prefer section-grouping when a phase is uniform** (kind read once, file scannable); **drop to inline where a group is mixed**.

```md
## Automated {run}
- [ ] npm test                  ← the line text IS the command
- [ ] npm run lint

## Manual QA {observable}
- [ ] button is centered

## Rollback {manual}
- [ ] run the rollback, confirm prior schema restored
```

- `## Title {run}` → every line under it is **executable**; the **line text is the command**.
- `## Title {observable|assert|judge|manual}` → every line is that kind (line text is the property/prompt).
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
- **Section wins, absolutely.** To mix kinds, use an *untagged* section with inline tags, or split into per-kind sections.
- **One kind per check.** Exactly one (first match). Don't combine inline tags.
- **Tags stay literal in the file.** Never stripped or moved — the file re-parses to the same kinds every walk.
- **Payload comes from the text:** executable ← the `run:` command (inline) or the whole line (under `{run}`); assertable/judgable ← the property/artifact named; observable/manual ← the line as the prompt.
- **Who sets it:** the skill proposes kinds when scaffolding VERIFY.md; untagged = observable, so a bare checklist still walks.

## Behavior

### 0. Pre-flight

1. **Resolve the request.** `<slug>` must exist at `.spectacular/requests/<slug>/`. If it's in `archive/`, refuse ("already archived; restore first").
2. **Check status.** The walk is meant for `status: review`.
   - If `review` → proceed.
   - If `active` → "`<slug>` is still `active`. Move to `review` first (`spectacular advance <slug> --to review`)?"
   - If `verified` already → "`<slug>` is already verified. Re-run the walk anyway? (records a fresh VERIFY-LOG entry)"
   - If `planned` → refuse: "nothing to verify yet — `<slug>` hasn't been built."
3. **Locate the verification artifact** (§ 1).

**Substrate check (auto-invoked on failure):** if the request's PLAN/VERIFY won't parse, auto-run `spectacular doctor workspace frontmatter` and surface findings before refusing. See [[doctor-substrate]].

### 1. Locate the checks (artifact resolution)

Resolve in this order:

1. **`requests/<slug>/VERIFY.md`** exists → walk its checkbox items (across all sections).
2. **No VERIFY.md** → fall back to **`PLAN.md § Validation`** items (and `TASKS.md § Verification` if present).
3. **Neither has checkable items** → tell the user verification has no defined checks; offer to (a) scaffold a VERIFY.md now if the 2-of-6 rule ([[plan-rules]]) warrants one, or (b) capture an ad-hoc check list for this walk.

Collect the full ordered list of checks before asking anything.

### 2. The walk loop — verify by kind

Parse each check's kind first, then route to its authority:

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
no human. true → PASS (state what was checked + the result). false → BLOCKER (what failed).
Deterministic: same inputs, same verdict. If the property is ambiguous to evaluate,
treat it as judgable instead — don't fake a binary on a fuzzy claim.
```

**judgable** (`{judge}`):
```
Read the named artifact(s). Reason over them; state the judgment + the evidence seen.
PASS if confident; if uncertain or the artifact is missing → BLOCKER
(never a confident pass on a guess). Surface the reasoning so the human can override.
```

**observable** (untagged or `{observable}`):
```
Ask: "How is this confirmed? (evidence / blocked / skip)"  — the human only looks.
• evidence given → record it as the human's confirmation → PASS.
  (The human is the authority — the agent records, it does not overrule.)
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

**Evidence stamp (manual + observable passes).** Before recording, ask once: "against what? (`<commit/build> · <app/extension identity>`)" — the stamp goes on the log row as `against:`. Default offer: current `git rev-parse --short HEAD`. Evidence without a stamp can't be told from stale evidence later ([[review-sweep]]).

After each: record the outcome (§ 3), advance.

**Exec safety (gating `run:` commands).** Default: **confirm each** executable check before running (show the command, `y/n/skip`). At walk start, offer **"run all executable checks this walk? (y)"** to batch-approve. Never run a command that wasn't shown to the user.

**Never fabricate a pass.** executable → only the exit code decides. assertable → the property is true or it's a blocker. judgable → uncertainty is a blocker, not an optimistic pass. observable → "looks done" without evidence is a blocker. manual → a confirmation that skipped the action is a blocker.

**Walk all checks** even after a blocker appears — the user gets the complete picture in one pass, not a premature exit.

### 3. Recording results (two artifacts)

The walk writes to **both**:

**a) VERIFY.md (or PLAN § Validation) — live state.**
- On PASS: tick the checkbox inline (`- [ ]` → `- [x]`).
- On BLOCKER / SKIP: leave the box unchecked. (The detail goes in the log, not as inline clutter.)

**b) VERIFY-LOG.md — the audit trail.**
- Location: `.spectacular/requests/<slug>/VERIFY-LOG.md`.
- Append one **timestamped walk entry** per run: every check + evidence + outcome + blocker reasons.
- Append-only, never overwritten. Scaffold from § VERIFY-LOG shape below if absent.

### 3b. Coherence pass (after the checks, before the gate)

The checks prove the milestones; coherence proves the *plan*. Skim the PLAN's `## Decisions` entries and confirm each shipped decision actually appears in the built artifact — a decision that silently didn't ship (chose X over Y, but the code does Y) is a finding. Grade it like a judgable check: advisory, reported in the walk summary, never a fabricated pass. Skip only when Decisions is empty.

### 4. The gate (end of walk)

After every check is walked, summarize:

```
Walk complete — <slug>
  ✓ <P> passed   ✗ <B> blocked   ⊘ <S> skipped   (of <N> checks)
```

**If all checks passed (B = 0, S = 0):**
- Default is **propose, human confirms**: "All <N> checks passed. Mark `<slug>` verified? (y/n)" — on `y` → `spectacular advance <slug> --to verified`. Never edit `status:` directly.
- **Configurable auto-flip:** `verify.auto_promote: true` in config (or `--auto`) skips the prompt and promotes immediately, reporting it.

**If any blocker or skip (B > 0 or S > 0):**
- Do **not** promote. Status stays `review`. Report the blocker list (with reasons) and point at the VERIFY-LOG entry.

### 5. Retrospective (end of a passing walk)

After a successful verify, ask once — skipping is fine:

> "Anything surprise you vs. the PRD/PLAN — a wrong assumption, a scope drift, a lesson? (skip to finish)"

If the user answers, route it to a memory entry via `spectacular remember` (never written autonomously — this just pre-fills it).

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
- ✓ [observe] <check> — <evidence the human gave> — against: <commit/build> · <identity>
- ✓ [manual] <check> — performed <action>, result: <result> — against: <commit/build> · <identity>
- ✗ [exec] <check> — BLOCKED: `<cmd>` exit 1 — <stderr tail>
- ⊘ [observe] <check> — skipped
- ⟳ [manual] <check> — pending-reverify: code moved past against-stamp (<who flagged> <date>)
**Outcome:** <verified | stayed in review — N blockers>
```

The `[kind]` tag records *which authority* confirmed each line — the audit trail shows not just "passed" but *how*. `[manual]`/`[observe]` rows additionally carry an **`against:` stamp** (commit/build + app/extension identity) so evidence is tied to what it proved. **`pending-reverify`** is a judgment flip — a sweep ([[review-sweep]]) or human appends a `⟳` row when the code moved past the stamp; it clears only by re-performing the check (a fresh stamped ✓ row). An old ✗ never means a current bug, and an old ✓ never means current proof — the stamp decides. `doctor lifecycle` warns on unstamped manual/observe rows and on `pending-reverify` rows in `review`/`verified` requests (mechanical presence checks only — no git heuristics).

## Lifecycle + archive tie-in

- The walk is the *intended* path to `verified`. `spectacular advance <slug> --to verified` still works directly (the CLI stays a dumb mutator), but the walk is what makes the checks real.
- **Archive warning:** `spectacular archive <slug>` should warn when a request reaches archive `verified` but with **no VERIFY-LOG.md** — flipped without ever being walked. Advisory, not blocking (see [[archive]]).

## CLI redirect

`spectacular verify <slug>` at the terminal prints:

> `verify` is a skill-side walk — it needs to read each check and judge your evidence. Run it inside Claude Code or Codex: `/spectacular` then `verify <slug>`.

## Related

- [[verify-authoring]] — 2-of-6 rule, fold patterns, VERIFY.md shape, promoting checks to scripts
- [[plan-rules]] § 2-of-6 rule — compact copy used at scaffold/grill time
- [[lifecycle]] — the `review → verified` transition the walk gates
- [[archive]] — verified precondition + the un-walked warning
- [[scaffold-reference]] — VERIFY.md stub
- [[principles]] — Principle 7 (intent → execution → validation); "agents propose, humans decide"
