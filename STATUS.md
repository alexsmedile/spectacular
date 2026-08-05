# STATUS — Issue backlog audit

**Updated:** 2026-08-05
**Repo:** `alexsmedile/spectacular`
**Audited against:** `main` @ `a5b4d5f` ("Merge pull request #21 from alexsmedile/codex/traffic-preflight")
**Open issues:** 19 (14 pre-existing + 5 filed by this audit)

---

## What this session did

Audited all open issues for **quality and currency** — are they written well enough for an independent session to act on, and do their factual claims still match the code? No implementation work was performed and no issue was closed.

- Verified every factual claim in all 14 pre-existing issues against the working tree
- Posted an audit comment to each of the 14
- Backfilled the missing `**One-liner:**` (per #26) into the 8 issues lacking it
- Filed 5 new issues for gaps that no existing issue named — **#28, #29, #30, #31, #32**
- Confirmed the `to-issues` agent **did fire reliably** — see below

### Issues filed by this audit

| # | Title | Origin |
|---|---|---|
| **28** | `summary --brief --json` schema collision between #11 and #12 | found by cross-reading #11 and #12 |
| **29** | AFK Git hygiene has no commit layer | found by verifying #5's own open question |
| **30** | Portfolio-level issue review workflow | user feedback — the recurring prompt |
| **31** | Graph engineering: nodes as agent execution substrate | user feedback — parent thesis for #20/#24 |
| **32** | PR bodies need a plain-language one-liner before the manifest | user feedback — corrects a mis-scope in #26 |

## to-issues agent: verified working

Four messages were sent to the `to-issues` agent. Three required filing; all three landed correctly on 2026-08-04 17:00–17:01, each with a `**One-liner:**` and correct cross-links.

| Message | Issue | Status |
|---|---|---|
| AFK end-goal gate on commits/PRs | **#25** | Filed correctly |
| ASCII node graph during orchestration | **#24** | Filed correctly |
| PR/issue one-liner requirement | **#26** | Filed correctly, self-applied |

**No missed filings. The suspected unreliability was not reproduced.**

---

## User items resolved 2026-08-05

Four enhancement requests, each checked against the code before answering.

### 1. Backlog review workflow → **filed #30**

The user re-sends this prompt to the orchestrator repeatedly:

> "any remaining issue open? how to act on the remaining ones? present, issue, status against current / assessment, reccomended action, why (optional). Impact order? Dedup? Grilling required? Separate branch? Let's review each, before sending an handoff prompt to independent session."

**Partially present.** `spectacular github triage <issue>` exists (`cli/spectacular:5344`, contract at `github-work-bridge.md:167-179`) and emits a good per-issue card — Meaning / Ready / Missing / Route / Why / Next.

**But it is single-issue and stateless.** No `--all`, no currency checking, no dedup, no impact ordering, no grill flag, no branch-collision detection, no handoff emission. Every cross-issue finding in this audit is invisible to per-issue triage. #30 is the direct generalization; it reuses the existing card as its per-issue unit.

### 2. AFK end-goal gate → **already filed (#25), rescope needed**

Already covered — see Headline finding 1. Most of the behaviour exists; the real gap is the missing commit layer, filed as **#29**.

### 3. Node graph → **filed #24, but under-scoped → parent filed as #31**

The user's framing is bigger than #24 as filed:

> "the new frontier of agentic workflows is to visualize nodes and give agents the ability to work on nodes. Humans have a more limited capacity to work on more then a few parallel nodes or branches. Agents and subagents can orchestrate that flawlessly."

#20 plans nodes, #24 draws them — **neither lets agents execute on them.** #31 files that as the parent thesis: node claiming, enforced dependency edges, parallel-eligibility derivation, per-node state, result reintegration. The ASCII view is reframed as the human supervision surface for a fan-out too wide to follow otherwise.

### 4. PR one-liner → **#26 was mis-scoped → filed #32**

**Correction to this audit's earlier reading.** #26 was interpreted as covering GitHub *Issue* bodies. The user clarified: *"pr = pull request on github. I had difficulties understanding what the pr, filed by spectacular, was about."*

The manifest at `github-work-bridge.md:242` requires no plain-language opening line; "purpose" is a change summary inside a structured body, and provenance fields lead. **#26 will not fix this.** Filed **#32** for the PR side. #26 remains valid for Issue bodies.

---

## Headline findings

### 1. #25 is largely already implemented — rescope before acting

`references/afk-git-hygiene.md` already provides goal-scoped AFK runs recording exact `goal` / `allowed_actions` / `hitl_gates` (`cli/spectacular:4869-4886`), three-layer authorization (durable run + project policy + `--apply --yes`), dry-run-first on every mutating command, and `afk pr` that opens **draft PRs only** and never merges.

The genuinely missing parts are narrower: no commit coverage at all, no scope-drift check at action time, no session-end reconciliation. **Acting on #25 as written would rebuild working behaviour.**

### 2. The AFK contract has a hole in the middle of its own lifecycle → filed #29

`grep -i "commit"` over `afk-git-hygiene.md` returns **zero matches**. The file governs branch creation, cleanup, and PR handoff — steps 1-3 and 5-6 of its own documented lifecycle. Step 4, "work and verification happen on the isolated branch," is where every commit happens and is entirely unspecified.

This also **answers #5's own open question** ("`afk-git-hygiene.md` may already cover part of this — needs review to confirm"): it does not.

### 3. #11 and #12 will collide on the same schema → filed #28

Both independently propose extending `summary --brief --json` with overlapping-but-unequal field lists. Whichever ships first defines a schema the other must break or fork. Neither issue can resolve this alone.

### 4. Stale figures in the SKILL.md pair

| Claim | Source | Actual |
|---|---|---|
| "230 lines" | #4 | **267 lines / 3,325 words** |
| "263 lines / 3,149 words" | #11 | **267 lines / 3,325 words** |
| "40+ reference files" | #4 | **72 files** |
| `references/glossary.md` | #4 acceptance criteria | **does not exist** — this is a create, not an edit |

Premise is unweakened — it is stronger. But two issues recorded two different counts weeks apart and both are now wrong, which is itself the argument for #11's P0 benchmark.

### 5. #10 is probably not a bug

The issue body investigates `buzz pr`; the maintainer's comment says the real context was the GitHub CLI via the Codex remote app. GitHub not permitting self-approval is intended platform behaviour. Recommend retitling to a self-merge **policy** decision and removing the `bug` label — it is currently the only bug-labelled issue, which overstates its urgency.

---

## Dependency map

```
GOAL / INTENT LANE
  #17 (validate goal)  ──►  #25 (scope gate)  ──►  #5 (commit split)
        │                         └──────────────►  #29 (commit layer)   [one branch]
        └──►  #32 (PR one-liner)                    [same manifest as #17 — one branch]

GRAPH LANE
  #31 (nodes as execution substrate — PARENT THESIS, grill with #20)
        ├──►  #20 (mission/node model, GRILL FIRST)
        └──►  #24 (render + state; needs `claimed` state if #31 proceeds)
                    ▲
  #18 (ASCII renderer, unblocked) ──┘

TOKEN / RETRIEVAL LANE
  #11 P0 (measure)  ──►  #28 (schema decision)  ──►  #11 P1 + #12 P1  ──►  #13
        └──►  #4 (trim + glossary, inside #11 P1)
                                                          └──►  #19 (grouping)

WORKFLOW LANE
  #30 (portfolio backlog review)  — builds on existing `github triage`, consumes #26

UNSEQUENCED
  #9 (evidence-blocked)   #10 (needs reframe)   #26 (needs target confirmation)
```

### Branch-collision groups

Issues that edit the same file/section and **must share a branch**:

| Group | Issues | Shared target |
|---|---|---|
| AFK commit layer | #25 → #5 → #29 | `references/afk-git-hygiene.md`, one absent section |
| PR manifest | #17 + #32 | `references/github-work-bridge.md` § Pull-request handoff |
| Brief schema | #11 P1 + #12 P1 (via #28) | `summary --brief --json` |
| SKILL.md | #4 inside #11 P1 | `skills/spectacular/SKILL.md` |

---

## Per-issue readiness

Legend: **Ready** = an independent session can act now · **Blocked** = needs a decision or evidence first · **Grill** = needs a design pass before dispatch

| # | Title | Verdict | What unblocks it |
|---|---|---|---|
| **18** | ASCII workflows/schemas | **Ready** | Nothing. Best first dispatch — build renderer with extensible per-node state so #24/#31 reuse it |
| **29** | AFK commit layer *(new)* | **Ready** | Nothing, but execute with #25 + #5 on one branch |
| **32** | PR one-liner *(new)* | **Ready\*** | \*Decide text source: new flag, request goal, or `--summary`. Share a branch with #17 |
| **17** | Intent/goal in schemas | **Ready\*** | \*Tighten "meaningful" → checkable (presence/length); `doctor` can't validate meaningfulness. Share a branch with #32 |
| **19** | Group by concept | **Ready\*** | \*Name one first target flow (recommend decision review); ideally fold into #12 |
| **30** | Backlog review workflow *(new)* | **Ready\*** | \*Decide command shape (`triage --all` vs `backlog review`) and whether it emits a prompt or a durable artifact |
| **26** | Issue one-liner | **Partial** | Body backfill ready (list below). Convention enforcement blocked: is the to-issues format defined in this repo or in the agent's own skill? |
| **5** | AFK commit scoping | **Blocked** | Merge with #25 + #29 onto one branch — three issues editing one absent section |
| **25** | AFK approval gate | **Blocked** | Rescope per finding 1 + decide: does scope-drift block, or warn and log? |
| **4** | Debloat SKILL.md | **Blocked** | Sequence after #11 P0; re-measure baseline; `glossary.md` must be created |
| **11** | Token footprint | **Blocked** | Two maintainer decisions: optimization target (median / worst-case / round-trips) and what % counts as success |
| **12** | Soft-DB projections | **Blocked** | #28 + shared P0 with #11 |
| **28** | Schema collision *(new)* | **Blocked** | Maintainer decision — command identity, field list, owning branch |
| **24** | ASCII mission graph | **Blocked** | Design work can start now; implementation waits on #20 |
| **10** | PR self-merge | **Blocked** | Which system (buzz vs GitHub CLI), exact command, exact error text |
| **9** | Spec/intent divergence | **Blocked** | Paste the original `unwire` prompt that produced SPC-002 — **most perishable evidence in the backlog** |
| **20** | Mission orchestrator | **Grill** | 3 architecture questions, one of which determines the whole data model. **Grill together with #31** |
| **31** | Nodes as execution substrate *(new)* | **Grill** | Framing issue, not a buildable brief. Its execution requirements constrain #20's data model — design them together or retrofit expensively |
| **13** | Semantic retrieval | **Park** | Correctly gated behind #11 + #12. Recommend labelling `blocked` and removing from the active queue |

---

## Recommended next actions

Ordered by leverage, not by issue number.

1. **Answer the 8 blocking questions** (below). Most of the backlog is blocked on maintainer decisions, not on work.
2. **Recover the #9 evidence now.** Session history decays; everything else in this backlog can wait, that cannot.
3. **Dispatch #18** as the first build. Unblocked, self-contained, de-risks #24 and #31.
4. **Merge #25 + #5 + #29 onto one branch**, order #25 → #5. Three issues, one absent section, guaranteed conflict otherwise.
5. **Grill #31 with #20**, not after. #31's execution requirements determine what #20's data model must support; designing the model first and retrofitting concurrency is the expensive order.
6. **Pair #17 + #32** on one branch — both edit the PR manifest, and both are candidates to derive text from PLAN `## Goal`.
7. **Consider #30 early despite its position.** It automates this entire audit. Every subsequent backlog review gets cheaper, and it builds on the existing `github triage` card rather than new machinery.
8. **Edit the 8 issue bodies** to add the one-liner — comments do not satisfy #26's stated purpose.

### The 8 blocking questions

| # | Question | Blocks |
|---|---|---|
| 1 | Optimization target — median context, worst-case, or tool round-trips? And what reduction is success? | #11, #4, #12, #13 |
| 2 | `summary --brief --json` vs `status --json` vs new `spectacular brief`; and the field list | #28, #11, #12 |
| 3 | Does an AFK scope-drift check block, or warn and log? | #25, #5, #29 |
| 4 | Is the to-issues agent's format defined in this repo, or in its own skill? | #26 |
| 5 | #10 — which system, which command, what error? | #10 |
| 6 | Mission/node = new doc type, linked requests, or generalized AFK run record? | #20, #24, #31 |
| 7 | Where does the PR one-liner text come from — new flag, request goal, or `--summary`? | #32, #17 |
| 8 | Backlog review command shape, and does it emit a prompt or a durable artifact? | #30 |

**Note on Q6:** the AFK run record is already a durable, goal-scoped, gate-aware object that survives sessions and moves `active → gated → active`. That is structurally close to what a node needs. Worth evaluating before inventing a new doc type.

**Note on Q7:** deriving the PR one-liner from PLAN `## Goal` risks reproducing the original problem — goal text is written for the request lifecycle, not for a PR reader. #17 recommends sourcing PR intent from that field, so Q7 is really one decision about what `## Goal` is for and who its audience is.

---

## One-liner backfill status (per #26)

Missing at audit start — 8 of 14. The split was chronological: everything filed 2026-08-04 08:41 or later had one.

| Had it | Missing → backfilled via audit comment |
|---|---|
| #18, #19, #20, #24, #25, #26 | #4, #5, #9, #10, #11, #12, #13, #17 |

**Still open:** these were added in comments. The issue **bodies** are unchanged, and a comment does not serve #26's stated purpose ("a reader unfamiliar with the thread can immediately understand"). Editing 8 bodies is a larger action than an audit should take unilaterally — flagged rather than done.

---

## Repo hygiene — checked, clean

- No large files tracked. Largest is `cli/spectacular` at 620KB (single-file Bash CLI, expected).
- `_archive/`, `_backups/`, `_archived/`, `_snapshots/` all gitignored. `agent-fleet.zip` (24K) correctly ignored.
- `.scrapekit/` and `.playwright-mcp/` already gitignored.
- Working tree clean at audit time.

---

## Scope note

This audit deliberately did **not**: close any issue, edit any issue body, change any label or title, or touch any code. Every recommendation above that would alter issue state is left as a recommendation. Retitles suggested for **#10** and **#25**; label change suggested for **#10** (drop `bug`) and **#13** (add `blocked`).
