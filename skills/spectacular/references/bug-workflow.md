---
doc-id: bug-workflow
kind: reference
summary: "How the skill handles a bug: check prior fixes first, decide audit-first vs just-fix, then log a fix if it's reusable. Ties audit/ + fixes/ into a self-learning loop."
status: active
---

# Bug workflow — audit, fix, and the self-learning loop

Loaded when the user reports a bug, quirk, regression, or "why does X do Y". Governs how the skill routes between [[audit-rules]] (investigate), the request lifecycle (plan), and [[fixes-rules]] (log). The goal: **fix bugs at the right ceremony level — no infinite audit→plan→fix loop on a one-liner, no undocumented root-cause work on a cross-cutting one.**

## The arc — read this first (orchestrator's whole job)

**You are the orchestrator: the only role that holds the whole bug and the only hand that mutates.** You triage, decide ceremony, open the job, write every brief, plan every fix, and own every ledger write. The agents are narrow delegable labor (read-only Investigator/Researcher, closed-brief Fixer); the two things you never delegate are **fix-planning** and **mutation**. Verbs, no overlap: **Investigator discovers → you plan → Fixer applies → you resolve.** You bookend the fleet.

The steps, in order (each only if earned — ceremony scales with uncertainty):

```
0   Seen it before?     grep fixes/ signatures ──── match → apply known remedy, done
1   Ceremony gate       root cause clear AND single site? → 2a just-fix; else fleet
1b  Fan-out gate        ≥3 independent closed disjoint-file fixes → fan out; else self-serve
1c  Open the job        mkdir debug/<slug>/ · write job.json · assign slots + briefs   (fleet only)
2a  Just fix            apply inline, verify                                            (no fleet)
2b  Investigate/plan    Investigator discovers → you plan the fix → route to 1b
3   Resolve + log       verify · [opt. code-reviewer / test-verifier] · outcome.json · summarize to audit/A<N> + fixes/F<N> if earned
```

**Resuming a job cold?** `job.json`'s `status` tells you where you are: `investigating` → you're mid-2b (Investigator out or its findings need planning) · `researching` → a Researcher is out · `planning` → findings are in, plan the fixes · `fixing` → Fixers are out (or their diffs need confirming) · `verifying` → run the Success criteria · `resolved` → write `outcome.json` + graduate to the library. Read the `timeline` array for what already happened; you're the single writer of the spine, so pick up where it left off.

> **@Debugging policy gate.** First, run `spectacular policy @Debugging` and follow every active policy returned. The four policies map one-to-one to the steps below (`check-prior-fixes` → Step 0, `ceremony-matches-uncertainty` → Step 1, `fix-root-not-symptom` → Step 2, `log-only-verified-reusable` → Step 3). All are `warn` — surface the finding, then proceed. See [[policy-injection]].

> **Traceability (fleet jobs only).** When a bug runs the fleet (Investigator / Researcher / fanned-out Fixers), the orchestrator opens a trace folder and owns its spine — **Step 1c** is where and when. In short: `.spectacular/debugs/<job-slug>/` with `job.json` (you write + update); each agent writes its own leaf artifact to a path *you* assign in its prompt (`investigation.json`, `research/research-NN.json`, `fixes/fix-NN.json`). This is the *raw pipeline* — kept as a trace, never pruned. On resolve you **summarize** it into the permanent library: `audits/A<N>` (the examination — if findings are worth keeping) and `fixes/F<N>` (the remedy — if reusable), cross-linked and both pointing back to the trace. **debugs/ = pipeline while it happens; audits/ + fixes/ = the two distilled summaries.** A **just-fix 1–2 line bug skips the folder entirely** — ceremony scales with uncertainty. Full schemas + the debug/audit/fixes distinction: [[debug-trace]].

---

## Step 0 — Have we seen this before? (the self-learning loop)

**Before diagnosing, search prior fixes.** This is the payoff of the `fixes/` corpus.

```
spectacular fix list          # scan titles + signatures
# then grep the signature field for the symptom in hand:
grep -rl "<symptom keywords>" .spectacular/fixes/
```

- **Match found** → read that `F<N>.md`. Its **Root cause**, **Fix**, and **Signature** may resolve the new bug immediately, or tell you it's a known class (e.g. F2/F3: "trailing `&& ` returns exit 1"). Apply the known remedy; you've just closed the learning loop. If the recurrence itself is notable (the lesson didn't stick), log a *new* fix that `related: [F<N>]` and says so — that's what F3 does.
- **No match** → proceed to Step 1. You're diagnosing something genuinely new; it's a candidate to *add* to the corpus.

The signature is written for a reader who hasn't seen the code — another agent, another project. A per-project search today; the corpus is designed to later pool into `~/.spectacular/fixes/` or move via `spectacular fix export/import` (not built yet — see [[fixes-rules]]).

---

## Step 1 — Audit-first, or just-fix? (the ceremony decision)

Ask two questions. Do **not** open an audit reflexively — most small bugs don't need one.

```
Is the root cause already clear?   AND   Does the fix touch a single site?
        │                                        │
    yes │ yes ───────────────────────────► JUST FIX (Step 2a). No audit.
        │                                        │
        └── no, on either ──────────────► AUDIT FIRST (Step 2b).
```

**Just-fix** when: reproducible, root cause understood, change is local and non-breaking. The audit would be pure ceremony — you'd write "Problem / Root cause / Fix" and immediately resolve it. Skip it. (You may still log a *fix* afterward if it's reusable — Step 3.)

**Audit-first** when *any* of:
- Root cause is unclear — you need to investigate before you can even name the fix.
- The symptom spans multiple callers / files / subsystems (a shared-function bug — fix once at the root, not per caller).
- User-reported but not yet reproduced — the audit records the investigation trail.
- It might not be a bug (could be intended behavior, or a design question) — the audit's **Intended behavior** slot forces that call.

Once an audit reaches a clear root cause + proposed fix, it exits via its **Disposition**:
- **one-line fix** → apply it, `audit resolve <A> --disposition "one-line fix: ..."`.
- **folded into a request** → the bug needs real planning (multi-step, needs a PLAN/TASKS). `audit resolve <A> --disposition "requests/<slug>"`, then scaffold the request. This is the audit→**plan**→fix path.
- **became fix F<N>** → resolved + verified now; graduate with `audit resolve <A> --into-fix` (copies all slots forward).
- **won't-fix** → deliberate; state why.

**Anti-pattern (the infinite loop you must avoid):** audit → plan → fix → audit again for a trivial bug. If the two questions both say "yes", there is no audit and no request — just the fix. Ceremony scales with uncertainty, not with every bug.

---

## Step 1b — Delegate or self-serve? (the fan-out decision)

Once Step 1 lands a bug (or a batch of bugs) in the **just-fix** quadrant, one more question: do *you* (the main thread) apply the fix, or do you **fan out** to `debug-fixer` subagents? The main thread is the Diagnostician + Localizer — diagnosis and site-hunting never delegate. Only a **closed fix** (root cause known, fix decided, single site) can go to a Fixer, and only when fanning out actually pays for itself.

### Decision table — self vs fan-out

| Situation in hand | Orchestrator does it **itself** | **Fan out** to N Fixers |
|---|---|---|
| 1–2 closed fixes | ✓ — a spawn costs more than the edit | |
| **≥3 independent closed fixes** | | ✓ — concurrency wins; each Fixer self-verifies |
| N fixes that share code / are order-dependent | ✓ — sequential, one context (see *Same-file fixes* below) | |
| N fixes touching the **same file(s)** | ✓ — serialize inline, never parallel (see below) | |
| Cross-cutting (N symptoms, 1 root) | ✓ — it's *one* fix at the root, not N | |
| Fix unclear / needs diagnosis | ✓ — that's the Diagnostician; not delegable | |

**The rule under the table:** fan out only when the fixes are **independent, closed, and touch disjoint files** — and there are **3 or more**. Miss any of those → the orchestrator applies them inline. Independence is the gate; "closed" is inherited from the Fixer contract; "disjoint files" prevents the parallel-edit collision the cross-cutting rule alone misses.

Why 3, not 2: fan-out has a fixed cost (spawn + brief-write + collect). At 1–2 fixes that cost roughly equals the saving; at 3+ concurrency clearly wins. Don't pay orchestration cost until it pays back.

### Same-file fixes — serialize, don't parallelize (and don't branch)

Fan-out's disjoint-file rule exists because **parallel writers to one file collide** — two Fixers editing `users.py` in separate windows produce two diffs against the *same* base, and the second silently clobbers or conflicts with the first. The fix is not coordination machinery; it's **not having concurrent writers.** When ≥2 fixes touch the same file, the orchestrator applies them **inline, sequentially, in one context** — you *are* the single serial writer, so there's no race to resolve. Three cases:

- **Independent hunks, same file** — apply them one after another in the same window; each new edit sees the file the previous one already changed. Verify once at the end. No ordering to reason about beyond "don't lose an earlier hunk."
- **Ordered / blocking (fix B only makes sense after fix A)** — sequence explicitly: apply A → **re-verify** → *then* plan B against the new file state (A may have moved B's line numbers or changed what B needs to do) → apply B → verify. The re-verify between steps is the point; B planned against the pre-A file can be wrong.
- **Genuinely entangled (the two "fixes" are really one change to the same lines)** — it's a single fix, not two. Plan it as one closed edit; don't split what shares lines.

**Do we need git branches / worktrees for this? No.** Branches solve *concurrent* writers; serializing inline means there are none. Reaching for `isolation: worktree` per Fixer + a merge-back step would be real merge-conflict machinery to re-solve a problem sequential-inline already dissolves — YAGNI until same-file fan-out is measurably too slow, which for a debug batch it won't be. So: the orchestrator **verifies** every fix, but there's nothing to **merge** because it never forked. Parallelism is for disjoint files; the same file gets one careful serial hand.

### The fan-out loop

When the table says fan out:

1. **Write one closed brief per fix** — the five slots, in prose, in each subagent's prompt: **Problem / Intended / Root cause / Proposed fix / Success criteria** (same shape as the `fix new` flags; see [[fixes-rules]]). The bar: **each slot concrete enough that a Fixer who has never seen this bug applies it without a judgment call.** If any slot leaves the Fixer something to *decide* — which site, whether the root cause holds, what "done" means — it isn't closed; that decision is yours to make first, or it's audit-first (Step 2b). A slot you can't fill concretely means the fix isn't closed → not a Fixer job. Don't lean on the bounce to catch a lazy brief: the bounce is the backstop for a *wrong* brief, not the quality gate for a *vague* one — that gate is you, here.
2. **Spawn one `debug-fixer` per independent site, in parallel.** Each copies/edits only its own target, applies only its Proposed fix, verifies against its Success criteria, and returns `applied+diff` or `bounced+reason`. It never writes the ledger.
3. **Collect all results before any ledger write.** Each Fixer's returned block **is its Agent-tool result** (machine-read, not human prose); its `fixes/fix-NN.json` is the durable copy. Don't `fix new` as each returns — gather the full set first, so the corpus reflects one coherent batch.
4. **Route by verdict** (check the block is well-formed first — a `VERDICT` from the enum, a real `SITE`, and for `applied` a non-empty `DIFF` + a `VERIFY` that actually ran; a malformed or empty return is not a trustworthy `applied` — treat it as a bounce and re-diagnose):
   - **applied** → confirm the diff yourself, then (if reusable, Step 3) `spectacular fix new`. The ledger write stays single-threaded on the main thread — `use-audit-fix-verbs` holds. Never log a fix you didn't confirm from the returned diff + verify.
   - **bounced** → the brief was wrong (bad root cause, or actually cross-cutting/multi-site). **Re-route to audit-first (Step 2b) — never auto-retry the same brief.** A bounce is the boundary working; respect it.

**Fan-out anti-patterns:** fanning a cross-cutting bug into N per-caller Fixers (it's *one* root fix — `fix-root-not-symptom`); spawning Fixers for fixes that touch the same file (they collide); auto-retrying a bounced brief instead of re-diagnosing. When in doubt, self-serve — the table's default column is the safe one.

---

## Step 1c — Open the job (only when the fleet runs)

The moment Step 1 lands **audit-first**, or Step 1b says **fan out**, you're running the fleet — and *that's* the point the orchestrator does its setup. Step 1c **brackets the fan-out**: three setup acts run *before* any agent spawns, and one bookkeeping act runs *after each returns*. This is the orchestrator's first mutation and it's yours alone.

**Before the spawn** (the folder + spine must exist to hand out slot paths):

1. **Scaffold the trace folder** — `.spectacular/debugs/<job-slug>/` (slug = short kebab of the symptom).
2. **Write `job.json`** — the spine. Stamp at intake: `symptom_class` (`test_failure | runtime_error | wrong_behavior | build_error | performance | unknown`), `symptom`, `reporter` (`user | while-coding`), `ceremony` (`just-fix | audit-first` from Step 1), `status` (`investigating` or `fixing`). Empty `artifacts` index + `timeline`; `outcome: null`.
3. **Assign slots and hand out briefs** — you own slot assignment (you know the folder + fan-out count). Each agent's prompt carries **(a)** its brief and **(b)** the exact trace path it writes to: the Investigator gets `investigation.json`; each fanned-out Fixer gets `fixes/fix-01.json`, `fixes/fix-02.json`, …; a Researcher gets `research/research-NN.json`. No agent picks its own index — no collision. **Record the brief you hand out in `job.json`'s `brief` field** — the spawn prompt vanishes when the window closes, so persisting it is what lets you (or a resuming orchestrator) later check the findings against what you actually asked (the symmetry check in Step 2b).

**→ then the fan-out happens** — agents spawn (Step 1b loop / Step 2b), run, and each writes its own leaf artifact to the path you assigned.

**After each agent returns:**

4. **Update the spine** — append a `timeline` entry, update `status` + the `artifacts` index in `job.json`. You're the single writer of the spine; agents only write their own leaf artifact.

So the order is: **open (1–3) → fan out → update spine (4)**. Setup is strictly before the spawn; bookkeeping interleaves with returns.

**Skip this whole step for a just-fix 1–2 line bug** — no folder, no json, no agents (Step 2a). The trace exists for jobs that fan out; ceremony scales with uncertainty. What the brief itself contains is Step 1b (fix briefs) and Step 2b (investigation brief); this step is *where and when* the job opens. Full schemas: [[debug-trace]].

---

## Step 2a — Just fix it

Make the change, verify it, move on. Then decide Step 3 (log or not). This is the **self-serve** path from Step 1b — 1–2 fixes, or any fix that can't fan out.

## Step 2b — Investigate, plan, then route

Audit-first means the cause and/or site are unknown — **discovering the bug is genuinely harder than fixing it**, so that's the work to delegate. You (the orchestrator) can investigate inline, or hand the hunt to a `debug-investigator` subagent with its own context window. Three roles, three verbs, no overlap: **Investigator discovers → orchestrator plans → Fixer applies.**

**Delegate the investigation when:** the bug needs real exploration — trace the flow, reproduce, form and test hypotheses — and a fresh window will do it better than crowding the main thread. **Investigate inline when:** it's quick, or you already hold most of the context.

### Going in — write a well-crafted investigation brief

The Investigator gets a *scoped hunt*, not a loose "look at this." A good brief gives it:
- **Symptom** — what was observed, and the expected behaviour.
- **Where to look** — a starting point if you have one (file, subsystem, recent change).
- **Done means** — what closes the investigation: root cause found? repro achieved? a specific question answered?

A vague brief wastes the window. Scope the hunt the way you'd scope a task for a junior — enough to aim, not so much you've done the finding yourself.

### Coming back — findings, not a fix

**The channel:** the block the Investigator emits **is the Agent tool's return value** — you receive it as the tool result, not as a chat message. Its `investigation.json` is the *durable copy* of the same findings (survives a session loss; a resuming orchestrator reads it). Block = live hand-back; JSON = the record.

**First, check the return is well-formed.** Before routing, confirm the block parses: a `STATUS` from the enum, and — if `root-cause-found` — a non-empty Root cause + at least one Suspected site + EVIDENCE. A malformed or empty return (prose instead of the block, `root-cause-found` with no evidence, a STATUS off-enum) is **not** a plannable result — treat it like `needs-more-context` and re-spawn with a sharper brief, or investigate inline. Never plan a fix from a broken block.

A well-formed return also lists **hypotheses ruled out, each with the evidence that killed it** — an investigation that can't name what it eliminated hasn't tested anything. The orchestrator copies ruled-out hypotheses into the `audit/A<N>` entry so no future walk (a resuming orchestrator, a second Investigator, next month's you) re-opens a dead end. The disproof ledger is as valuable as the root cause.

**Then check it against your brief.** Compare the returned `STATUS` to the brief's **Done means**. Met → proceed to plan. Not met (you asked for a root cause, got `hypotheses-only`) → that's the loop signal, not a failure — decide plan-vs-loop by the REASON below. The brief you sent is recorded in `job.json`'s `brief` field, so this check survives resume too.

The Investigator reports **findings**, and it never proposes the edit — fix-planning is yours:
- **`STATUS: root-cause-found`** → root cause + suspected sites + blast radius + evidence. You now know *where* and *why*.
- **`STATUS: hypotheses-only`** → ranked hypotheses + a `REASON` (`needs-reproduction` / `needs-research` / `needs-decision` / `needs-more-context`). Route by reason: `needs-research` → a `debug-researcher` (is this a known platform/dependency bug?); `needs-reproduction` → **you handle it inline, or ask the user for a reliable repro** (there's no Reproducer agent — at this stage the orchestrator or user already holds the repro; a dedicated Reproducer is deferred to the Heisenbug case, see [[coding-agents]]); `needs-decision` → ask the human; `needs-more-context` → send it back with a sharper brief.

**The agent's `STATUS` ends *its turn*; it does not always end the *investigation*.** The investigation phase closes only when a STATUS you can plan a fix from comes back — `root-cause-found`, or a `hypotheses-only` whose REASON you can act on directly (`needs-decision` → you ask the human; `needs-research` → you route to a Researcher). The two looping reasons — `needs-more-context` and `needs-reproduction` — **re-open the hunt** (fresh Investigator with a sharper brief / get a reliable repro inline or from the user first), they don't close it. Don't mistake a returned block for a finished investigation and try to plan from hypotheses you can't act on — read the STATUS+REASON and decide: plan, or loop.

**3-strikes: after 3 failed fix attempts, question the architecture, not the hypothesis.** The trigger signature: each fix reveals a new problem in a different place. That is not a failed hypothesis — it's a wrong structure. Stop spawning Investigators/Fixers, lay out the fix history (what was tried, what each attempt broke) for the human, and ask whether the surrounding design is the actual bug. Never attempt a fourth fix without that conversation.

**The `debug-researcher`** is read-only web/forum/doc search with a debugging protocol (frame the symptom in others' vocabulary → diverse queries → judge relevancy hard → verdict). Use it when a bug smells *external* — a library or platform quirk others may have hit. It returns `VERDICT: known-platform-bug` (apply the documented workaround / pin-or-upgrade) · `genuinely-ours` (route back to an Investigator, external cause now ruled out) · `no-strong-match` (honest dead end). It prefers `scrapekit` for fetching, falls back to the harness `WebSearch`/`WebFetch`. Like the other agents: never writes the ledger.

**Optional: a second opinion from `codex-agent`.** Same read-only contract as the Investigator (Bash/Glob/Grep/Read/WebFetch — no Edit/Write), but a different model — a cross-check, not a replacement for the Investigator. Default is **no** — most `root-cause-found` findings are good enough to plan from as-is. Reach for it only when one of these is true:

| Trigger | Why it earns the spawn |
|---|---|
| The fix would be **cross-cutting or high blast-radius** (touches a shared function, many callers, or a hot path) | Planning wrong here is expensive — worth a second read before committing |
| You've **already looped once** (re-spawned the Investigator with a sharper brief and still got `hypotheses-only`, or `root-cause-found` that didn't hold up on the fix attempt) | Same-model re-spawn already failed to converge; a different reasoning model can break the stall |
| The root cause is **plausible but you're not confident** — no hard evidence, just the best of several hypotheses | A confirm/refute pass is cheap insurance before planning a fix on a guess |

**Don't** spawn it for a routine `root-cause-found` with solid evidence and a single obvious site — that's the common case, and a second opinion there is pure ceremony (mirrors the `ceremony-matches-uncertainty` policy: match effort to actual uncertainty, not to every bug).

Give it the same evidence the Investigator had (symptom, suspected sites, root cause so far) and ask it to confirm, refute, or propose an alternative — never to fix. Treat its answer like a second `hypotheses-only` opinion: it informs your plan, it doesn't replace your planning. It never writes the ledger.

### Plan the fixes, then fan out

With findings in hand, **you plan the fix(es)** — this is the orchestrator's job, the step the Investigator deliberately doesn't do. Turn the root cause into one or more *closed* fixes (each with the five slots: Problem / Intended / Root cause / Proposed fix / Success criteria). Then the Step 1b decision applies: **≥3 independent, closed, disjoint-file fixes → fan out `debug-fixer`s; 1–2 → fix inline.** Fan-out is usually faster *once the fixes are clear and well-scoped* — which they now are, because the investigation closed the unknowns.

**But first, a fork: can the findings even close into fixes?** Sometimes the investigation reveals the "fix" is *design work* — a schema change, a new policy, a concept the code doesn't have yet, touching every create/read path together (the Investigator will usually flag this as blast-radius + open questions, and size it "not a one-line fix"). You can tell because the five slots won't fill: **Proposed fix** would read "design X," and a Fixer handed that brief would have to *invent* the decisions (a TTL value, a clock source, an eviction policy). That fails the "concrete enough a Fixer applies without judgment" bar — so **do not fan out.** Instead, **fold it into a request** (the fleet→request bridge):

- `spectacular new "<slug>"` to scaffold the request (`PLAN.md`/`TASKS.md`) — the open questions become the plan's decisions.
- Write the job's `outcome.json` → `disposition: folded-into-request`, `request: <slug>`, `logged_fixes: []` (nothing was fixed — it was planned), flip the spine to `resolved`, keep the debug folder as trace.
- Optionally `spectacular audit new` first if the investigation earned an examination trail, then fold *that* (`audit resolve <A> --disposition "requests/<slug>"`). But an audit is **not required** to fold — the fleet path can route straight from `investigation.json` to a request. (This is the same fold as line 75's audit→request path, reached from the fleet instead of an open audit.)

This is a real disposition, not a failure: `folded-into-request` closes the debug job cleanly — the work continues in the request lifecycle, not the debug pipeline. **No `F<N>` is logged for a folded job** (Step 3 logs *verified fixes*; a fold has none yet).

**And a second fork: should the fix even be applied?** Sometimes the findings *do* close into a clean five-slot fix — and the right call is still **not to apply it.** The bug is real and the diff is obvious, but applying it would break a frozen consumer, touch deprecated code with a live alternative, regress a deliberate trade-off, or cost more than the symptom warrants (the Investigator often flags this as blast-radius on a path everything's migrated off). This is **won't-fix** (the third disposition) — a *decision*, not a failure to fix:

- **Do not apply the edit.** Weigh apply-vs-decline explicitly; if decline wins, leave the code untouched.
- Write the job's `outcome.json` → `disposition: wont-fix`, a **stated `reason`** (why declining beats fixing — and the migration/alternative path if there is one), `logged_fixes: []`, flip the spine to `resolved`, keep the debug folder as trace.
- **No `F<N>` is logged** — nothing was applied, so nothing graduates (Step 3 logs *verified fixes*; a decline has none). The `reason` is the durable record: a future reader learns *why it wasn't fixed*, which stops the same bug from being "rediscovered" and re-litigated.
- **Just-fix ceremony** (no debug/ folder was opened)? Record the decline on the audit instead: `spectacular audit new "..."` → `audit resolve <A> --disposition "won't-fix: <reason>"`. Same disposition, recorded where the examination lives. (This is line 77's audit `won't-fix`, reached inline instead of from a fleet job.)

`wont-fix` and `folded-into-request` are the two ways a debug job closes **without** a fix landing — one declines the work, the other relocates it to the request lifecycle. Both are honest resolutions; neither logs an `F<N>`.

**If the findings DO close AND the fix is worth applying** — the common case — continue:

Record it: `spectacular audit new` (from the findings) if the investigation earned an audit trail, and `spectacular fix new` per verified fix (Step 3). **The ledger stays single-threaded on the orchestrator** — neither the Investigator nor the Fixers write it (`use-audit-fix-verbs`).

**Pipeline (which role runs each step):** you triage + gate (Steps 0–1, the "Router") → Investigator discovers findings (agent, Step 2b) → **you** plan the fixes (Step 2b) → Fixer applies ×N (agent, Step 1b fan-out) → **you** confirm + write the ledger (Step 3, the "Validator"). "Router" and "Validator" are just *you* wearing the triage hat and the resolve hat — they aren't separate agents; every main-thread box is the orchestrator. The two real agents are delegable labor and keep the honest-fallback invariant: Investigator reports hypotheses-only rather than a fake root cause; Fixer bounces rather than freelance.

---

## Step 3 — Log a fix? (build the corpus without spamming it)

> **Before you call it verified — optional arms-length review + verify (judgment-gated).** Consider a
> fresh-window confirmation before treating the bug as resolved — same worth-it economics as fan-out;
> skip it for a trivial one-site fix. Two independent triggers (mirrors [[build-workflow]] Step 3):
> - **[[code-reviewer]]** (read-only) — when the fix diff is **substantial or medium+ blast radius**
>   (touches a shared module / schema / widely-used helper), over the fix diff — did the fix introduce a *new* problem
>   (a regression, an insecure path, dead code)? Returns severity-ranked findings; you triage and,
>   if needed, dispatch another `debug-fixer`.
> - **[[test-verifier]]** (apply-only, tests only) — when the `debug-fixer` **self-reported the pass**
>   or blast radius is medium+, independent pass/fail on the fix's check, or a regression test written
>   to the (now closed) reproduction spec. *The agent that wrote the fix shouldn't be the only one to
>   grade it.* On `fail`, the bug isn't resolved — plan the next fix.
>
> Both optional; the default (you confirm the Fixer's diff + verify yourself) stands for ordinary
> fixes. They inform resolution; the ledger write below stays yours.

After a bug is **resolved and verified**, decide whether it earns a `fixes/` entry. Log it when the fix carries **reusable knowledge** — a future agent or project would benefit:

- non-obvious root cause (a footgun, a platform quirk, an ordering trap)
- a bug *class* likely to recur (F2/F3's trailing-`&&` is the archetype)
- anything that took real investigation to understand

**Don't** log: typos, one-off content edits, a rename, anything whose "fix" teaches nothing. Not every bug is a fix entry — the corpus is valuable because it's curated, not exhaustive.

```
spectacular fix new "<title>" \
  --problem "..." --intended "..." --cause "..." --fix "..." \
  --criteria "..." --verified-by "<the check>" \
  --signature "<symptoms + pattern a future search would match>"
```

The `--signature` is the single most important flag — it's what makes the entry *findable* in Step 0 next time. Omitting `--verified-by` warns and sets `verified: null` (a draft, not a trustworthy fix).

If the fix came from an audit, prefer `spectacular audit resolve <A> --into-fix` — it copies Problem/Intended/Root cause/Proposed fix/Success criteria forward automatically; you only add Verified-by + Signature.

---

## The loop, in one line

**seen-it? (fixes/) → ceremony decision (audit vs just-fix) → self-serve or fan out (≥3 independent closed fixes → N× debug-fixer) → resolve (fix · plan · won't-fix) → log if reusable (fixes/).**

Each verified, well-signed fix makes the next bug cheaper — that's the self-learning loop.

---

## Policy backing — the `@Debugging` hook

Each step above is gated by a policy under the `@Debugging` hook in `POLICY.md` (retrieved when this doc loads). The steps are the *practice*; the policies are the *enforcement* — same order, same spine as principle 11 (earn each step):

| Step | Policy | Principle |
|---|---|---|
| 0 — seen it? | `check-prior-fixes` | 5 |
| 1 — ceremony decision | `ceremony-matches-uncertainty` | 11 |
| 2 — fix at the root | `fix-root-not-symptom` | 11 |
| 3 — log if reusable | `log-only-verified-reusable` | 5 |

All four are `warn` (surface + continue) — debugging judgment stays with the human; the policies make the right move the obvious one. See [[policies-contract]] for the hook set.

**Related:** [[audit-rules]], [[fixes-rules]], [[new-request]] (the fold/plan path), [[lifecycle]], [[decisions-rules]] (why-we-chose vs what-broke), [[doc-index]].
