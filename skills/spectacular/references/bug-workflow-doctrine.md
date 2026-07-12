---
doc-id: bug-workflow-doctrine
kind: reference
summary: "The rationale behind bug-workflow.md's rules — why each gate and disposition exists, the failure modes each prevents. Load only when a routing call feels uncertain or when editing the workflow itself; never needed for routine bug handling."
status: active
---

# Bug workflow — doctrine

The *why* behind every rule in [[bug-workflow]]. That doc is the runtime core (steps, gates,
tables, routing); this one holds the reasoning. Read the section matching the step you're unsure
about.

## Why the orchestrator bookends the fleet

The orchestrator is the only role that holds the whole bug, so triage, ceremony, fix-planning, and
every ledger write stay on the main thread. "Router" and "Validator" aren't separate agents — they
are you wearing the triage hat and the resolve hat; every main-thread box in the pipeline is the
orchestrator. The two real agents keep the honest-fallback invariant: the Investigator reports
`hypotheses-only` rather than a fake root cause; the Fixer bounces rather than freelance.
Pipeline: you triage + gate (Steps 0–1) → Investigator discovers (2b) → **you** plan (2b) → Fixer
applies ×N (1b) → **you** confirm + write the ledger (3).

## Step 0 — why prior-fixes-first

This is the payoff of the `fixes/` corpus: a signature written for a reader who hasn't seen the
code (another agent, another project) turns yesterday's investigation into today's grep hit. F2/F3
("trailing `&&` returns exit 1") is the archetype of a known *class* resolving a new symptom
instantly. The corpus is per-project today but designed to pool later (`~/.spectacular/fixes/`,
`fix export/import` — not built yet; see [[fixes-rules]]). When a lesson didn't stick and the same
bug recurs, logging a *new* fix with `related: [F<N>]` is itself signal — that's what F3 does.

## Step 1 — why ceremony scales with uncertainty

Opening an audit reflexively on a clear one-liner produces pure ceremony: you'd write
"Problem / Root cause / Fix" and immediately resolve it. The two-question gate (cause clear? site
single?) exists because the audit's value is *recording an investigation* — no investigation, no
audit. The other direction is worse: undocumented root-cause work on a cross-cutting bug loses the
trail that a future reader (or resuming orchestrator) needs. The infinite audit→plan→fix→audit
loop on a trivial bug is the anti-pattern this gate kills.

## Step 1b — the economics and the collision rule

**Why 3, not 2:** fan-out has a fixed cost (spawn + brief-write + collect). At 1–2 fixes that cost
roughly equals the saving; at 3+ concurrency clearly wins. Don't pay orchestration cost until it
pays back. Deliberately the same threshold as [[build-workflow]] B1 — identical economics.

**Why disjoint files:** two Fixers editing `users.py` in separate windows produce two diffs
against the *same* base; the second silently clobbers or conflicts. The fix is not coordination
machinery — it's **not having concurrent writers**. Serializing inline makes you the single serial
writer; there is no race to resolve. The re-verify between ordered fixes exists because fix A may
move fix B's line numbers or change what B needs to do — B planned against the pre-A file can be
wrong.

**Why no branches/worktrees:** branches solve *concurrent* writers; serialized-inline has none.
Per-Fixer worktrees + merge-back is real merge-conflict machinery re-solving a problem
sequential-inline already dissolves — YAGNI until same-file fan-out is measurably too slow, which
for a debug batch it won't be. The orchestrator verifies every fix, but there's nothing to
*merge* because it never forked.

**Why the closed-brief bar is yours, not the bounce's:** the bounce is the backstop for a *wrong*
brief (bad root cause, actually cross-cutting) — not the quality gate for a *vague* one. A vague
brief doesn't fail loudly; it comes back as a confident wrong fix, or wastes a spawn on a bounce
you could have prevented by closing the slot yourself.

## Step 1c — why the job opens exactly there, and why briefs persist

The folder + spine must exist *before* any agent spawns because slot assignment needs them, and
slot assignment is the orchestrator's because it's the only role that knows the folder and fan-out
count — that's what makes collisions impossible by construction. The brief is recorded in
`job.json` because the spawn prompt vanishes when the window closes: without the record, neither
you nor a resuming orchestrator can check the findings against what was actually asked (the
symmetry check), and a resumed job would be flying blind. Setup strictly before the spawn;
bookkeeping interleaves with returns. The whole step is skipped for just-fix bugs because the
trace exists for jobs that fan out — ceremony scales with uncertainty.

## Step 2b — why discovery delegates but planning never does

Audit-first means the cause and/or site are unknown — **discovering the bug is genuinely harder
than fixing it**, so discovery is the work worth a fresh window. But the Investigator deliberately
stops at findings: it maps the solution space (approaches + trade-offs) without prescribing the
edit, because choosing the approach requires holding the whole bug, its blast radius, and the
project's constraints — which only the orchestrator does. Three roles, three verbs, no overlap.

**Why the well-formed check:** an LLM return can be prose instead of the block, `root-cause-found`
with no evidence, a STATUS off-enum. Planning a fix from a broken block launders a guess into a
diagnosis. Same reason the ruled-out list matters: an investigation that can't name what it
eliminated hasn't tested anything, and the disproof ledger is as valuable as the root cause — it
stops a future walk (a resuming orchestrator, a second Investigator, next month's you) from
re-opening a dead end.

**Why STATUS ends the turn, not the investigation:** the phase closes only when a *plannable*
status arrives — `root-cause-found`, or a `hypotheses-only` whose REASON you can act on directly
(`needs-decision` → ask; `needs-research` → Researcher). The two looping reasons
(`needs-more-context`, `needs-reproduction`) re-open the hunt. Planning from hypotheses you can't
act on is how wrong fixes get confident.

**Why 3-strikes questions the architecture:** when each fix reveals a new problem in a different
place, that's not a failed hypothesis — it's a wrong structure. A fourth attempt without the human
conversation just relocates the symptom again. (There's no Reproducer agent for
`needs-reproduction` because at that stage the orchestrator or user already holds the repro; a
dedicated Reproducer is deferred to the Heisenbug case — see [[coding-agents]].)

**Why the codex second opinion defaults to no:** most `root-cause-found` findings with solid
evidence and a single obvious site are good enough to plan from — a second opinion there is pure
ceremony (`ceremony-matches-uncertainty` again). It earns the spawn only when planning wrong is
expensive (cross-cutting / high blast radius), when same-model re-spawning already failed to
converge (a different reasoning model can break the stall), or when the best hypothesis has no
hard evidence. Treat its answer like a second `hypotheses-only` opinion — it informs your plan, it
doesn't replace your planning.

## The two no-fix dispositions — why they're resolutions, not failures

**`folded-into-request`:** sometimes the investigation reveals the "fix" is design work — a schema
change, a new concept, decisions a Fixer would have to *invent* (a TTL value, a clock source, an
eviction policy). You can tell because the five slots won't fill: Proposed fix would read "design
X", which fails the applies-without-judgment bar. Fanning out anyway delegates design to an agent
that can't hold it. The fold closes the debug job cleanly and relocates the work to the request
lifecycle, where open questions become the plan's decisions. No `F<N>` — Step 3 logs *verified
fixes*, and a fold has none yet.

**`wont-fix`:** sometimes the findings *do* close into a clean five-slot fix and the right call is
still not to apply it — it would break a frozen consumer, touch deprecated code with a live
alternative, regress a deliberate trade-off, or cost more than the symptom warrants. This is a
*decision*, not a failure to fix, and the stated `reason` is the durable record: a future reader
learns *why it wasn't fixed*, which stops the same bug from being rediscovered and re-litigated.
No `F<N>` — nothing was applied, so nothing graduates.

Both dispositions matter precisely because they're the honest exits: without them, every opened
job pressures toward *some* landed diff, which is how frozen consumers get broken and design work
gets smuggled in as "fixes".

## Step 3 — why the corpus is curated, not exhaustive

The corpus's value is its hit rate in Step 0: every typo and rename logged is noise that lowers
it. Log only what teaches — non-obvious causes, recurring classes, real investigations. The
signature is the most important field because it's the *search key*: written in the vocabulary a
future symptom would grep for, not in this codebase's internal names. `verified: null` drafts
exist so an unverified remedy can't masquerade as a trusted one.

**Why the arms-length pass mirrors build D2:** the agent that wrote the fix shouldn't be the only
one to grade it, and a fix diff can introduce a *new* problem the symptom check never exercises.
Same worth-it economics as fan-out — skip for a trivial one-site fix; reach for it on blast
radius or self-reported passes.

**Related:** [[bug-workflow]] (the runtime core this explains), [[build-workflow-doctrine]] (the
build-direction mirror's rationale, incl. the relation table), [[debug-trace]], [[audit-rules]],
[[fixes-rules]], [[policies-contract]].
