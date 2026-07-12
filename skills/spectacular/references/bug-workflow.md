---
doc-id: bug-workflow
kind: reference
summary: "Runtime core for handling a bug: check prior fixes first, decide audit-first vs just-fix, run the fleet when earned, log a fix if reusable. Ties audit/ + fixes/ into the self-learning loop. Rationale lives in bug-workflow-doctrine.md; the build-direction mirror is build-workflow.md."
status: active
---

# Bug workflow — audit, fix, and the self-learning loop

Loaded when the user reports a bug, quirk, regression, or "why does X do Y". Goal: **fix bugs at
the right ceremony level** — no infinite audit→plan→fix loop on a one-liner, no undocumented
root-cause work on a cross-cutting one. **You are the orchestrator: the only role that holds the
whole bug and the only hand that mutates.** The agents are narrow delegable labor (read-only
Investigator/Researcher, closed-brief Fixer); the two things you never delegate are **fix-planning**
and **mutation**. Verbs, no overlap: **Investigator discovers → you plan → Fixer applies → you
resolve.**

> This is the runtime core: the steps, the gates, the tables, the routing. The *why* — rationale,
> failure modes, the fold/won't-fix reasoning — is in [[bug-workflow-doctrine]]. Load it only when
> a routing call feels uncertain or you're editing this workflow itself.

## The arc

```
0   Seen it before?     grep fixes/ signatures ──── match → apply known remedy, done
1   Ceremony gate       root cause clear AND single site? → 2a just-fix; else fleet
1b  Fan-out gate        ≥3 independent closed disjoint-file fixes → fan out; else self-serve
1c  Open the job        mkdir debugs/<slug>/ · write job.json · assign slots + briefs   (fleet only)
2a  Just fix            apply inline, verify                                            (no fleet)
2b  Investigate/plan    Investigator discovers → you plan the fix → route to 1b
3   Resolve + log       verify · [opt. code-reviewer / test-verifier] · outcome.json · summarize to audit/A<N> + fixes/F<N> if earned
```

**Resuming a job cold?** `job.json`'s `status` tells you where you are: `investigating` → mid-2b ·
`researching` → a Researcher is out · `planning` → findings in, plan the fixes · `fixing` → Fixers
out or diffs need confirming · `verifying` → run the Success criteria · `resolved` → write
`outcome.json` + graduate. Read the `timeline` for what happened; you're the single writer of the
trace, so pick up where it left off.

> **@Debugging policy gate.** First, run `spectacular policy @Debugging` and follow every active
> policy returned: `check-prior-fixes` → Step 0, `ceremony-matches-uncertainty` → Step 1,
> `fix-root-not-symptom` → Step 2, `log-only-verified-reusable` → Step 3. All are `warn`. See
> [[policy-injection]].

> **Traceability (fleet jobs only).** When a bug runs the fleet, open `.spectacular/debugs/<job-slug>/`
> with `job.json` (you write + update); as each agent returns, *you* persist its returned block as
> the leaf artifact at the slot you assigned (`investigation.json`, `research/research-NN.json`,
> `fixes/fix-NN.json`) — agents return blocks, they write no trace files. On resolve, **summarize**
> into the permanent library: `audits/A<N>` (if findings worth keeping) + `fixes/F<N>` (if
> reusable). A just-fix 1–2 line bug skips the folder entirely. Schemas: [[debug-trace]].

---

## Step 0 — Have we seen this before?

**Before diagnosing, search prior fixes:**

```
spectacular fix list          # scan titles + signatures
grep -rl "<symptom keywords>" .spectacular/fixes/
```

- **Match** → read that `F<N>.md`; its Root cause / Fix / Signature may resolve the bug immediately
  or name a known class. Apply the known remedy. If the recurrence itself is notable, log a new fix
  with `related: [F<N>]`.
- **No match** → Step 1; this is a candidate to *add* to the corpus.

## Step 1 — Audit-first, or just-fix? (the ceremony decision)

```
Is the root cause already clear?   AND   Does the fix touch a single site?
        │                                        │
    yes │ yes ───────────────────────────► JUST FIX (Step 2a). No audit.
        │                                        │
        └── no, on either ──────────────► AUDIT FIRST (Step 2b).
```

**Just-fix** when: reproducible, root cause understood, change local and non-breaking. (You may
still log a fix afterward — Step 3.)

**Audit-first** when *any* of: root cause unclear · symptom spans multiple callers/files (fix once
at the shared root) · user-reported but not reproduced · might not be a bug (the audit's
**Intended behavior** slot forces that call).

An audit exits via its **Disposition**: **one-line fix** (`audit resolve <A> --disposition "..."`)
· **folded into a request** (`--disposition "requests/<slug>"`, then scaffold the request) ·
**became fix F<N>** (`audit resolve <A> --into-fix`) · **won't-fix** (state why).

**Anti-pattern:** audit → plan → fix → audit again for a trivial bug. Two yeses → just the fix.

## Step 1b — Delegate or self-serve? (the fan-out decision)

Only a **closed fix** (root cause known, fix decided, single site) can go to a Fixer. Diagnosis and
site-hunting never delegate.

| Situation in hand | Orchestrator does it **itself** | **Fan out** to N Fixers |
|---|---|---|
| 1–2 closed fixes | ✓ — a spawn costs more than the edit | |
| **≥3 independent closed fixes** | | ✓ — concurrency wins |
| Fixes that share code / are order-dependent | ✓ — sequential, one context | |
| Fixes touching the **same file(s)** | ✓ — serialize inline, never parallel | |
| Cross-cutting (N symptoms, 1 root) | ✓ — it's *one* fix at the root | |
| Fix unclear / needs diagnosis | ✓ — not delegable | |

**The rule: fan out only when the fixes are independent, closed, and touch disjoint files — and
there are 3 or more.** Miss any → apply inline. (Fixed cost ≈ saving at 1–2; concurrency wins at
3+. Same gate as [[build-workflow]] B1.)

**Same-file fixes — serialize, don't parallelize (no branches/worktrees).** Three cases:
- **Independent hunks, same file** — apply one after another in one window; verify once at the end.
- **Ordered / blocking** — apply A → **re-verify** → *then* plan B against the new state → apply →
  verify.
- **Genuinely entangled** — it's a single fix; don't split what shares lines.

### The fan-out loop

1. **Write one closed brief per fix** — the five slots in each subagent's prompt: **Problem /
   Intended / Root cause / Proposed fix / Success criteria** ([[fixes-rules]]). Bar: each slot
   concrete enough that a Fixer who has never seen this bug applies it without a judgment call. A
   slot you can't fill → not closed → not a Fixer job (or it's audit-first). If the fix includes a
   test, add a one-line conventions note (the test file + harness to match). Don't lean on the
   bounce to catch a lazy brief — that gate is you, here.
2. **Spawn one `debug-fixer` per independent site, in parallel.** Each applies only its Proposed
   fix, verifies its Success criteria, returns `applied` or `bounced+reason`. Never writes the
   ledger.
3. **Collect all results before any ledger write**; persist each block to its `fixes/fix-NN.json`
   slot (Step 1c.4). Don't `fix new` as each returns — one coherent batch.
4. **Route by verdict** (well-formed first: `VERDICT` from the enum, a real `SITE`, and for
   `applied` a non-empty `CHANGED` + a `VERIFY` that ran; malformed = treat as bounce):
   - **applied** → confirm the change yourself — `git diff -- <the files in CHANGED>` (the block
     carries no diff; the edit lives in the tree) — then (if reusable) Step 3 `spectacular fix new`.
     Never log a fix you didn't confirm from the tree's diff + verify.
   - **bounced** → re-route to audit-first (Step 2b) — **never auto-retry the same brief.**

**Anti-patterns:** fanning a cross-cutting bug into N per-caller Fixers (it's one root fix);
Fixers on the same file (they collide); auto-retrying a bounced brief. When in doubt, self-serve.

## Step 1c — Open the job (only when the fleet runs)

The moment Step 1 says **audit-first** or Step 1b says **fan out**, before any agent spawns:

1. **Scaffold the trace folder** — `.spectacular/debugs/<job-slug>/`.
2. **Write `job.json`** — the spine. Stamp: `symptom_class`, `symptom`, `reporter`, `ceremony`,
   `status` (`investigating` or `fixing`); empty `artifacts` + `timeline`; `outcome: null`.
3. **Assign slots and hand out briefs** — Investigator → `investigation.json`; Fixers →
   `fixes/fix-01.json`, `fix-02.json`, …; Researcher → `research/research-NN.json`. Slots are
   *your* bookkeeping — agents get only their brief; they never write trace files. **Record each
   brief in `job.json`'s `brief` field** (the spawn prompt vanishes; the record enables the
   symmetry check in Step 2b and survives resume).

**→ agents spawn, run, return blocks.** After each returns:

4. **Persist + update the spine** — write the returned block to its slot as JSON (schemas:
   [[debug-trace]]), append a `timeline` entry, update `status` + `artifacts`. You're the single
   writer of the whole trace.

**Skip this step for a just-fix bug** — no folder, no json, no agents.

## Step 2a — Just fix it

Make the change, verify it, move on. Then decide Step 3.

## Step 2b — Investigate, plan, then route

**Delegate the investigation** when the bug needs real exploration and a fresh window beats
crowding the main thread; **investigate inline** when it's quick or you hold the context.

**The brief (going in):** **Symptom** (observed + expected) · **Where to look** (if known) · **Done
means** (root cause? repro? a specific question answered?). Scope the hunt like a task for a junior.

**The return (coming back):** the block is the Agent-tool result; you persist it to
`investigation.json` (Step 1c.4). Check it in order:
1. **Well-formed?** A `STATUS` from the enum; if `root-cause-found`, a non-empty Root cause + ≥1
   Suspected site + EVIDENCE. Malformed/empty → treat as `needs-more-context`: re-spawn sharper or
   investigate inline. Never plan from a broken block. A well-formed return also lists **ruled-out
   hypotheses with the evidence that killed each** — copy them into the `audit/A<N>` entry so no
   future walk re-opens a dead end.
2. **Against your brief?** Compare `STATUS` to the brief's Done-means (the brief is in `job.json`).
3. **Route by STATUS:**
   - **`root-cause-found`** → plan the fix(es).
   - **`hypotheses-only`** → route by `REASON`: `needs-research` → `debug-researcher` ·
     `needs-reproduction` → repro inline or ask the user (no Reproducer agent) · `needs-decision` →
     ask the human · `needs-more-context` → re-spawn with a sharper brief.
   The agent's STATUS ends *its turn*, not necessarily the investigation — `needs-more-context` and
   `needs-reproduction` re-open the hunt.

**3-strikes:** after 3 failed fix attempts (each fix reveals a new problem elsewhere), stop
spawning agents — lay out the fix history for the human and ask whether the surrounding design is
the actual bug. Never attempt a fourth fix without that conversation.

**`debug-researcher`** — when the bug smells external (library/platform quirk). Returns
`known-platform-bug` (apply documented workaround / pin-or-upgrade) · `genuinely-ours` (back to
Investigator, external ruled out) · `no-strong-match` (honest dead end).

**Optional second opinion (`codex-agent`, read-only)** — default **no**. Reach for it only when:
the fix would be cross-cutting/high blast-radius · you've already looped once and didn't converge ·
the root cause is plausible but unevidenced. Give it the same evidence; ask confirm/refute/
alternative — never to fix.

### Plan the fixes, then route

With findings in hand, **you plan** — turn the root cause into closed fixes (five slots each), then
apply Step 1b: ≥3 independent closed disjoint-file fixes → fan out; 1–2 → inline.

**Fork 1 — findings won't close into fixes** (the "fix" is design work: slots would read "design
X", a Fixer would have to invent decisions) → **fold into a request:**
- `spectacular new "<slug>"` — the open questions become the plan's decisions.
- `outcome.json` → `disposition: folded-into-request`, `request: <slug>`, `logged_fixes: []`; spine
  → `resolved`; keep the folder as trace. Optionally `audit new` first, then
  `audit resolve <A> --disposition "requests/<slug>"` — an audit is not required to fold.
- No `F<N>` is logged (nothing was fixed — it was planned).

**Fork 2 — the fix closes but shouldn't be applied** (would break a frozen consumer, touch
deprecated code with a live alternative, regress a deliberate trade-off, or cost more than the
symptom) → **won't-fix:**
- Leave the code untouched. `outcome.json` → `disposition: wont-fix` + a stated `reason` (why
  declining beats fixing, and the migration path if any), `logged_fixes: []`; spine → `resolved`.
- No `F<N>` is logged. Under just-fix ceremony (no folder), record the decline on an audit instead:
  `audit resolve <A> --disposition "won't-fix: <reason>"`.

Both forks close the job honestly without a fix landing.

**Common case — findings close AND the fix is worth applying** → continue: `audit new` if the
investigation earned a trail, fixes via Step 1b, ledger via Step 3. The ledger stays
single-threaded on you (`use-audit-fix-verbs`).

## Step 3 — Log a fix? (build the corpus without spamming it)

> **Multi-phase fix?** Occasionally a root cause spans several verify-point phases (change the
> shared function → update callers → add the regression test). Apply [[build-workflow]] B2:
> sequential sub-steps, confirm between each — never one opaque multi-phase brief.

> **Optional arms-length pass before calling it verified (judgment-gated, mirrors
> [[build-workflow]] D2):** [[code-reviewer]] over the fix diff when it's substantial or medium+
> blast radius (did the fix introduce a *new* problem?) · [[test-verifier]] when the Fixer
> self-reported the pass or blast radius is medium+ (independent pass/fail, or a regression test
> to the now-closed spec). On `fail`, the bug isn't resolved. Both optional; the ledger write
> stays yours.

After a bug is **resolved and verified**, log it when the fix carries **reusable knowledge**:
non-obvious root cause (footgun, platform quirk, ordering trap) · a bug class likely to recur ·
anything that took real investigation. **Don't** log typos, one-off content edits, renames — the
corpus is valuable because it's curated.

```
spectacular fix new "<title>" \
  --problem "..." --intended "..." --cause "..." --fix "..." \
  --criteria "..." --verified-by "<the check>" \
  --signature "<symptoms + pattern a future search would match>"
```

`--signature` is the most important flag — it's what makes the entry findable in Step 0. Omitting
`--verified-by` warns and sets `verified: null` (a draft). If the fix came from an audit, prefer
`spectacular audit resolve <A> --into-fix` (copies the five slots forward; you add Verified-by +
Signature).

---

## The loop, in one line

**seen-it? (fixes/) → ceremony decision (audit vs just-fix) → self-serve or fan out (≥3 independent
closed fixes → N× debug-fixer) → resolve (fix · fold · won't-fix) → log if reusable (fixes/).**

Each verified, well-signed fix makes the next bug cheaper — that's the self-learning loop.

## Policy backing — the `@Debugging` hook

| Step | Policy | Principle |
|---|---|---|
| 0 — seen it? | `check-prior-fixes` | 5 |
| 1 — ceremony decision | `ceremony-matches-uncertainty` | 11 |
| 2 — fix at the root | `fix-root-not-symptom` | 11 |
| 3 — log if reusable | `log-only-verified-reusable` | 5 |

All four are `warn` — debugging judgment stays with the human.

**Related:** [[bug-workflow-doctrine]] (the why), [[debug-trace]] (schemas), [[audit-rules]],
[[fixes-rules]], [[new-request]] (the fold path), [[build-workflow]] (the mirror), [[lifecycle]],
[[doc-index]].
