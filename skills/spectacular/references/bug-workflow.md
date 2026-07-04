---
doc-id: bug-workflow
kind: reference
summary: "How the skill handles a bug: check prior fixes first, decide audit-first vs just-fix, then log a fix if it's reusable. Ties audit/ + fixes/ into a self-learning loop."
status: active
---

# Bug workflow — audit, fix, and the self-learning loop

Loaded when the user reports a bug, quirk, regression, or "why does X do Y". Governs how the skill routes between [[audit-rules]] (investigate), the request lifecycle (plan), and [[fixes-rules]] (log). The goal: **fix bugs at the right ceremony level — no infinite audit→plan→fix loop on a one-liner, no undocumented root-cause work on a cross-cutting one.**

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

## Step 2a — Just fix it

Make the change, verify it, move on. Then decide Step 3 (log or not).

## Step 2b — Audit, then route

`spectacular audit new "<title>" --problem "<symptom>" --intended "<what should happen>"`, fill Investigation/Root cause/Proposed fix/Success criteria as you learn them, then resolve via the disposition that fits (above).

---

## Step 3 — Log a fix? (build the corpus without spamming it)

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

**seen-it? (fixes/) → ceremony decision (audit vs just-fix) → resolve (fix · plan · won't-fix) → log if reusable (fixes/).**

Each verified, well-signed fix makes the next bug cheaper — that's the self-learning loop.

**Related:** [[audit-rules]], [[fixes-rules]], [[new-request]] (the fold/plan path), [[lifecycle]], [[decisions-rules]] (why-we-chose vs what-broke), [[doc-index]].
