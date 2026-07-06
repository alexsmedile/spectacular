# Fable Review — Spectacular Spec/Planning Quality Uplift

- **Date:** 2026-07-06
- **Branch:** `review/fable-spec-quality`
- **Reviewer:** Fable 5, per `FABLE_REVIEW_PROMPT.md`
- **Method:** three parallel passes — (A) corpus audit of 18 real artifacts across 5 live workspaces (harbor, octopus, wasabi, unwire, spectacular-self), (B) guidance-layer audit of the skill's reference docs + templates + CLI scaffolds, (C) comparative research on superpowers (obra), gstack (garrytan), and OpenSpec (Fission-AI).
- **Scope:** move six metrics — decision density · falsifiability · diagnosis discipline · strategic altitude · progressive disclosure · anti-drift. Everything else out of scope.

---

## 1. Corpus findings

Eighteen artifacts sampled (13 request PLANs/TASKS incl. 4 archived, 2 PRDs, 1 SPEC index, VERIFY files where present). Full scorecard in Appendix A. The distribution, per metric:

**Decision density — the framework's strongest metric.** 14/18 artifacts state decisions with rationale and rejected alternatives. Best-in-class: octopus `06-adapter-framework` (20 numbered locked decisions, each with a why and the alternative it killed), unwire `plugin-governance-model` ("the 'extract a shared `ir_core.py`' idea is **dropped**. Measured on real code, the genuinely-shared surface is ~30 trivial lines"). Every strong scorer cites a **grill session by date**; the weak ones (octopus PRD, `33-tui-visual-redesign`) never went through one. *Grill works — when it runs.*

**Falsifiability — bimodal, and the split follows the template, not the project.** Plans using the numbered 7-slot template carry real checks ("`scripts/validate-provider.sh` goes fully green", "`footprint` shows the WebContent pid gone — 743MB→0"). Plans written freeform skip validation entirely (unwire `plugin-governance-model` has **no validation section at all**; octopus 33's check is "screenshot shows the design language matches"). The template's §6 slot is load-bearing where present; nothing forces a check to be *assertable*, and nothing forces the template.

**Diagnosis discipline — excellent where exercised (n=2, both harbor, both post-bug-workflow-v1.25).** The gold standard keeps every disproven theory with its disproof; the discipline even leaks into feature plans (unwire `primitive-atoms` TASKS *correct the plan's own estimates* with evidence). No older-project bug plans exist to compare — the guidance is young and unproven at scale, but not visibly failing.

**Strategic altitude — PRDs hold altitude when short, sink when long.** Harbor's 105-line PRD is the model (non-goals with reasons, constraints encoding verified research). Octopus's 784-line PRD opens at altitude then absorbs full schema dumps and CLI verb tables — content that belongs in SPEC/specs/. Guidance gap: nothing ever says "your PRD is 8× the median and contains schema dumps."

**Progressive disclosure — good on requests, uneven on canonicals.** Request plans cluster at 100–250 lines with pointers outward. The dumps are canonical docs (octopus PRD) and design-heavy plans that embed drafts belonging in DECISIONS.md.

**Anti-drift — the corpus's clear weakest metric, and it is a lifecycle-TAIL failure.** Everything up to "the code works" is disciplined; everything after is best-effort:
- wasabi `m1-two-webviews` archived with TASKS `status: active`, 4 open boxes, and **every VERIFY.md checkbox unchecked** — while its sibling `sleep-policy-webviewpool`, same user same week, is a textbook close. The difference is enforcement, not skill.
- harbor `SPEC.md` frozen at *"Capabilities: nothing built yet"* against 49 requests, ~30+ shipped. The spec-sync contract was never exercised once on the corpus's most disciplined project.
- Non-vocab statuses proliferate exactly where reality outgrew the 4-state vocabulary: `fixed-pending-verify`, `in-progress`, `open`, `done`, `deferred-v2`.

**Verdict on guidance-gap vs user-gap:** decision density and diagnosis discipline are guidance *successes*. Plan-time falsifiability, PLAN→PRD traceability, and above all the **archive tail** are guidance gaps — the instructions are either advisory, contradictory, or point at structures that don't exist. The strongest observed behaviors (supersession blocks, evidence-ledger TASKS) were **invented by users** and have no framework home to amplify them.

---

## 2. What the guidance audit found (why the corpus looks like this)

The reference docs split cleanly into **sharpened blades** and **polite suggestions**:

- **Sharp:** `verify.md` ("Never fabricate a pass… only the exit code decides"; typed check authorities; append-only VERIFY-LOG) and `prd-rules.md` (vague-word lists, "≥1 number AND ≥1 verb AND ≥1 date" gate, before/after examples on every pattern). These two produce the corpus's strongest output. **prd-rules is the model the other rules files should be raised to.**
- **Toothless:** tasks-rules (all 5 checks syntactic — `- [ ] do the thing` passes), spec-rules (20 lines; index-ness lives only in a template comment), the autopilot path in new-request.md (skill drafts a whole PLAN with **no grill, no review gate** — the largest ungated content-production path), and the entire @Debugging policy layer (`warn` severity throughout).
- **Self-contradictory (trains the model against its own gates):**
  1. `target_version:` — plan-rules bans it ("Do not add it back"); roadmap-rules *depends* on it (Slot 6 autopopulation + the `full`-tier unlock trigger).
  2. `PLAN.md § Decisions` — decisions-rules routes request-scoped decisions there; **no template, slot, or check creates that section.** The destination is fictional. This is the direct cause of decision-drafts squatting in plan bodies (octopus 47 embeds four full decision drafts "to be assigned D-numbers on commit" — never committed).
  3. new-request.md embeds its own PLAN/TASKS templates that have **forked** from `templates/*/base.md` (embedded TASKS uses `## <Group>` headings that fail tasks-rules' own `### M<N>` check).
  4. review.md base check 3 requires `priority`,`owner` on PLANs; plan-rules' schema requires only `status`,`updated`,`summary`.
  5. Lifecycle "single home" (SKILL.md: state lives in PLAN frontmatter, "never duplicated") vs TASKS frontmatter carrying a mirrored `status:` plus sync-repair machinery.
  6. CLI fallback scaffold (cli/spectacular:~5804, when templates are missing) emits bare headings with **no placeholders** — the resulting PLAN can't fail the placeholder check because there's nothing to catch.

---

## 3. Comparative research — what the neighbors do better

| | superpowers (obra) | gstack (garrytan) | OpenSpec (Fission-AI) |
|---|---|---|---|
| **Core strength** | Behavioral enforcement under pressure | Interrogation quality + silent-failure hunting | Falsifiable specs + mechanical spec evolution |
| **Signature mechanism** | Iron Law + Excuse/Reality table + red-flag self-checks; plans written for "zero codebase context" readers with a placeholder ban | Evidence-before-questions ("cite `path:line` in your first question"); Error & Rescue Map ("zero silent failures"); scope-mode gate; default-on second-model review | SHALL-statement + GIVEN/WHEN/THEN scenario per requirement; ADDED/MODIFIED/REMOVED spec deltas merged at archive; `openspec validate` CI gate |
| **What it lacks** | Living spec state, lifecycle, decision log | A spec substrate at all; also the anti-pattern: ~800-line duplicated preambles per skill | Behavioral enforcement, debugging discipline, strategy layer |
| **Most relevant to Spectacular's gaps** | POLICY armor, brief-writing bar for spec-builder | Grill grounding, plan failure-mode section | **The archive-tail fix**: delta-based spec-sync + scenario-backed capability bullets |

Spectacular's architecture (lean orchestrator, progressive loading, living workspace, lifecycle) is *ahead* of all three. Its weaknesses are exactly where each neighbor is strong: OpenSpec has solved the spec-drift problem Spectacular's corpus exhibits; superpowers has solved instruction-compliance-under-pressure for the policies Spectacular marks `warn`; gstack has solved grill-grounding-in-evidence.

---

## 4. Ranked change list

Blast radius: **S** = surgical wording · **T** = template/schema (structural but contained) · **W** = workflow (touches a flow's shape).

| # | Change | Files | Metric(s) | Radius | Grounding |
|---|---|---|---|---|---|
| 1 | **Archive closure gate** — `spectacular archive` (skill flow + `doctor lifecycle`) blocks/loudly-warns unless: TASKS boxes all `[x]`/`[~]`-with-reason, VERIFY (if present) walked, and a **spec delta** declared (see #2). | `archive.md`, `lifecycle.md`, CLI `cmd_archive`, doctor | anti-drift, falsifiability | W | Corpus: m1-two-webviews archived w/ unchecked VERIFY; harbor SPEC "nothing built yet" @ 49 requests |
| 2 | **Delta-based spec-sync** — archive declares ADDED/MODIFIED/REMOVED capability bullets instead of free-form "proposed updates"; mechanically mergeable; `doctor specs` validates structurally. | `spec-sync.md`, `spec-rules.md`, `archive.md` | anti-drift | W | OpenSpec deltas; feeds the planned `spec-audit-mode` request |
| 3 | **Fix the six contradictions** (§2 above): retire `target_version:` from roadmap-rules; de-fork new-request.md's embedded templates (replace with pointers to `templates/*/base.md`); align review.md check 3 to per-doc schema; resolve the TASKS `status:` mirror (drop it or amend the "never duplicated" claim); add placeholders to the CLI fallback scaffold. | `roadmap-rules.md`, `new-request.md`, `review.md`, `tasks-rules.md`/`templates/tasks/base.md`, `cli/spectacular` | anti-drift, all | S | Guidance audit §4; forks train the model against the gates |
| 4 | **Give decisions a real home** — add `## Decisions` (chose X over Y — because Z) to `templates/plan/base.md` + plan-rules review check "each entry names an alternative". | `templates/plan/base.md`, `plan-rules.md` | decision density | T | Fictional `PLAN § Decisions` destination; octopus 47's squatting decision drafts |
| 5 | **Wire check-kind vocabulary into PLAN § Validation at authoring time** — Slot 6 prompt + gate: each check declares its authority (`run:`/`assert`/`judge`/`observable`/`manual`); flag aspiration-verb checks (`improve/enhance/optimize` pattern reused from roadmap-rules). | `plan-rules.md`, `templates/plan/base.md` | falsifiability | S | Corpus bimodality; vague checks only surface months later at the verify walk |
| 6 | **Sanction the supersession-block convention** — a named pattern for updating a live PLAN: `## SUPERSEDED <date> — <what changed>` prepended, original kept below, frontmatter `summary` stays ≤2 sentences. | `active-request.md`, `plan-rules.md` | diagnosis discipline, anti-drift | S | Corpus pattern #1 — the gold standard invented this; every file currently hand-rolls its own layering |
| 7 | **Autopilot passes through the gate** — when the skill drafts PLAN bodies itself, it must run `plan review` and show the punch list *before* asking for confirmation (mirrors imagine.md, which already enforces this). | `new-request.md` | all | S | Largest ungated content path; grill-fingerprint corpus finding |
| 8 | **Lifecycle vocabulary pressure-release** — either add one state (`blocked` or `paused`) or, cheaper, an explicit rule: non-vocab intent goes in a `note:` frontmatter field, `status:` stays vocab-only; `doctor lifecycle` flags non-vocab values with that remediation. | `lifecycle.md`, doctor | anti-drift | S | Corpus pattern #4: `fixed-pending-verify`, `in-progress`, `deferred-v2` encode real states |
| 9 | **Evidence-before-questions grill rule** — before any technical grill question on a code-touching doc, read ≥1 piece of codebase evidence and cite `path:line` in the first question; never ask what the code already answers. | `grill.md` | decision density, altitude | S | gstack `/spec` Phase 3; grounds grills in the repo instead of generic checklists |
| 10 | **Disproven-hypotheses artifact + 3-strikes escalation** — investigator returns must list hypotheses *ruled out with the evidence that killed them* (copied into `audit/A<N>`); after 3 failed fixes the orchestrator must stop and question the architecture with the human. | `bug-workflow.md`, `debug-trace.md` | diagnosis discipline | S | Gold standard's disproof ledger is currently voluntary; superpowers systematic-debugging |
| 11 | **Iron-Law armor for the 4 block policies** — each gets a one-line absolute law + a short Excuse/Reality table + red-flag self-check ("using 'should', 'probably', 'seems to'"). Consider raising the @Debugging hooks from `warn`. | `POLICY.md` templates, `policy-injection.md` | falsifiability, diagnosis | S | superpowers; the entire diagnosis layer is currently advisory |
| 12 | **PLAN Goal → PRD traceability gate** — review check: Goal names/links the PRD goal it serves; mini-refine trigger when Goal ≈ frontmatter `summary` ("that's the request restated — what does it *change*?"). | `plan-rules.md` | strategic altitude | S | Prompt-only today; octopus altitude sink |
| 13 | **TASKS acceptance stub + placeholder ban in briefs** — tasks template scaffolds one `→ check:` line per milestone group; spec-builder/debug-fixer briefs declare "TBD / appropriate error handling / similar to M<N>" as **brief failures**. | `templates/tasks/base.md`, `tasks-rules.md`, `build-workflow.md` | falsifiability | T | Corpus TASKS bimodality; superpowers writing-plans |
| 14 | **SPEC-index check** — spec-rules review override: any Capabilities bullet >2 lines or with sub-bullets → "promote to `specs/<cap>/` + leave a one-line pointer". New capability bullets prefer SHALL-strength observable phrasing; important ones earn a GIVEN/WHEN/THEN scenario in the capability spec. | `spec-rules.md`, `spec-sync.md` | anti-drift, falsifiability | S | Rule currently lives in a template comment; OpenSpec scenario discipline |
| 15 | **Verify: add the Coherence dimension + default-on retrospective** — the walk checks Completeness / Correctness / **Coherence** (do the PLAN's stated decisions actually appear in the code?); §5 retrospective changes from "optionally ask" to "ask once (skipping is fine)". | `verify.md` | falsifiability, anti-drift | S | OpenSpec verify; catches exactly the drift `spec-audit-mode` targets |

### Before/after for the top 5

**#1 — archive.md (new gate, inserted before the archive move):**
> *Now:* archive flow warns advisorily ("Advisory, not blocking") if the verification walk never ran.
> *After:* "**Closure gate.** Before moving the folder: (1) every TASKS box is `[x]` or `[~] <reason>` — an unexplained `[ ]` blocks; (2) if VERIFY.md exists, its walk ran (VERIFY-LOG has an entry) — an unwalked VERIFY blocks; (3) the spec delta (§ below) is declared, even if it is `NONE — <why this shipped nothing spec-visible>`. The user can override each block once, explicitly, and the override is recorded in the archive note."

**#2 — spec-sync.md:**
> *Now:* "propose SPEC.md + specs/ updates" (free-form).
> *After:* "The request's spec impact is declared as explicit deltas: `### ADDED` / `### MODIFIED` / `### REMOVED` capability bullets. MODIFIED quotes the current bullet and the replacement. Deltas merge mechanically into SPEC.md; `doctor specs` fails a delta that names a bullet which doesn't exist (REMOVED/MODIFIED) or already exists (ADDED)."

**#4 — templates/plan/base.md (after `## 2. Constraints`):**
```md
## Decisions
<!-- Design calls made inside this request. Format: chose X over Y — because Z.
     Rejected alternatives stay listed; deleting them re-litigates them later.
     Project-wide calls go to DECISIONS.md via `spectacular decide` instead. -->
- <DECISION — chose X over Y because Z>
```
Plus plan-rules gate check: "Each Decisions entry names an alternative (contains 'over', 'not', or 'instead of')."

**#5 — plan-rules.md Slot 6 prompt:**
> *Now:* "How each milestone is verified." (template placeholder: `- M1 — <VERIFICATION>`)
> *After:* "How each milestone is verified — each check states its **authority**: a `run:` command, an assertable property, a judgable artifact, or a human-observable behavior (see [[verify]] kinds). A check with no authority can't fail. Aspiration verbs (`improve`, `enhance`, `optimize`, `handle gracefully`) are not checks."

**#6 — active-request.md (new subsection):**
> "**Superseding a live plan.** When findings invalidate part of a PLAN: prepend a block `## SUPERSEDED <date> — <one line>` containing the corrected understanding, mark the affected sections `(superseded, kept for history)`, and keep them. Never delete a disproven section — the disproof trail is the plan's most valuable content (see the harbor `orphaned-local-files-shadow-drive` exemplar). Frontmatter `summary:` is rewritten to the *current* understanding, ≤2 sentences."

---

## 5. Explicit non-changes

| Considered | Rejected because |
|---|---|
| Folding reference docs into SKILL.md, or any consolidation of the lean-orchestrator | Breaks progressive disclosure on purpose; gstack's 800-line preambles are the cautionary tale |
| Adopting OpenSpec's full requirements format (SHALL + scenario for *every* bullet) | Ceremony mismatch — Spectacular's SPEC is an index by design; scenarios earn their place only on important capability bullets (#14 keeps the light version) |
| Default-on second-model (Codex) review for every plan | Already exists as a scoped trigger table for bugs (b21/d6091d2); making it default-on is a cost/latency decision for the user, not a guidance fix. Worth a POLICY opt-in at most |
| Making verify findings block archiving outright with no override | OpenSpec's own doctrine ("advisory… leaves the call to you") is right; #1 uses block-with-recorded-override instead |
| New lifecycle states beyond one pressure-release valve | The 4-state vocabulary is a feature; corpus drift is better fixed by the `note:` rule + doctor flag than by state proliferation |
| TDD enforcement (superpowers' Iron Law #1) | Host-project concern, not workspace-framework concern; STACK.md/POLICY.md is where a project opts into that |
| Rewriting octopus/older workspaces to comply | Out of scope — the fix is `doctor` + migrations catching it forward, not retroactive cleanup |
| CLI mutation-semantics changes | Load-bearing contract, explicitly out of scope per the prompt (the two CLI touches in #3 are stub content, not semantics) |

---

## 6. One highest-leverage change

**#1+#2, the archive closure gate with delta-based spec-sync.** The corpus is unambiguous: Spectacular's guidance already produces strong plans, strong diagnoses, and decision-dense docs — and then lets the ending rot. Work is archived with unchecked verification lists, statuses drift into invented vocabulary, and the flagship SPEC index of the most disciplined project in the corpus still says "nothing built yet" after ~30 shipped requests. One gate at the single convergence point every request already passes through (archive) fixes anti-drift, feeds the spec, and closes the falsifiability loop — without adding a gram of ceremony to the 95% of the lifecycle that already works.

---

## 7. A/B validation (added 2026-07-07, post-implementation)

After the W1+W2 changes shipped on this branch, the guidance was A/B-tested against `main`: identical synthetic scenarios (real-world inspired, traps planted), one blind agent per branch per scenario, each given ONLY its branch's reference docs and told to follow them exactly. Scenario files (re-runnable): `docs/reviews/ab-scenarios/`. Ten runs, 18 traps.

| # | Scenario (modeled on) | Trap | main | branch |
|---|---|---|---|---|
| A | Draft a PLAN from a fuzzy request (PDF export) | Surface the design fork (Chrome vs canvas) to the human | ⚠️ punted into M1 silently | ✅ `## Decisions` + [NEEDS CLARIFICATION] + asked directly |
| A | | Self-review before confirmation | ❌ not run | ✅ `plan review` punch list shown alongside draft |
| A | | Authority-tagged checks | ⚠️ decent, untagged | ✅ every check tagged (run/observable/assertable/judgable) |
| A | | Goal traced to PRD | ✅ G2 | ✅ G2, criterion quoted |
| A | | TASKS shape | ❌ used the forked embedded-template shape (fails tasks-rules check 5) | ✅ canonical `### M<N>` |
| B | Review a flawed PLAN (6 planted defects, csv-import) | Non-vocab `status: in-progress` | ✅ caught | ✅ caught + `note:` remediation |
| B | | Goal restates summary | ❌ noticed, declared "outside the gate" | ✅ gate failure (check 11), traced to PRD G1 |
| B | | 3 vibes-checks | ✅ "vague" | ✅ authority-less (check 7) |
| B | | Task-shaped milestones | ✅ caught (refine patterns) | ⚠️ missed — agent scoped strictly to gate checks |
| B | | Vague deliverable | ✅ | ✅ |
| B | | Frontmatter schema contradiction | ⚠️ silently resolved in the lucky direction | ✅ no ambiguity left |
| C | Debug orchestration: 3 failed fixes + dead ends in prose (harbor-style) | Stop before fix #4, escalate architecture | ❌ planned fix #4, decided tactical-vs-design itself | ✅ halted; fix-history table; "I won't attempt a fourth fix until you've made this call" |
| C | | Eliminated hypotheses preserved durably | ⚠️ improvised as ranked hypotheses w/ evidence_against (no schema home) | ✅ `ruled_out` array + copy-to-audit rule cited |
| D | "We were wrong — update the plan" (webhook-retries) | Disproven diagnosis kept | ❌ **deleted it** (one-line "revised" note) | ✅ SUPERSEDED block; all invalidated content kept + marked |
| D | | New design choice gets a home | ❌ buried in prose | ✅ Decisions entry with alternative + why |
| D | | Summary reflects current understanding | ✅ | ✅ |
| E | Verify walk: all checks pass but the shipped code contradicts a PLAN decision (fixed vs sliding window) | Drift caught | ⚠️ caught by **model initiative** (no prescribed step); walk log contains no trace of it | ✅ caught by the prescribed coherence pass; finding written into VERIFY-LOG |
| E | | Retrospective asked | ✅ (docs said "optionally") | ✅ per spec |

**Tally: branch 17/18 · main 9/18 by-the-book (~13/18 crediting model initiative).**

Reading: main often reaches good outcomes because a capable model reads everything and compensates — but those catches live in chat and depend on the model's mood; the branch turns them into prescribed steps with durable records (main's verify log was clean while its own chat flagged a real defect). The starkest deltas reproduced the corpus's real failures exactly: main deleted a disproven diagnosis (metric 3) and marched into fix #4 solo (3-strikes).

**Caveats:** n=1 per cell, same model on both sides, scenarios authored by the same reviewer. Directional evidence, not a controlled study.

**Residual finding → candidate follow-up:** the task-verb milestone check ("Implement X" is a task, not a milestone) lives only in refine patterns on both branches; the branch's strictly-scoped reviewer missed it in Scenario B. Consider promoting it to a plan review gate check (12).

## Appendix A — Corpus scorecard

Scores: **S** strong · **M** mixed · **W** weak · — n/a. Metrics: D=decision density, F=falsifiability, Dx=diagnosis discipline, A=strategic altitude, P=progressive disclosure, AD=anti-drift.

| Artifact | Type | Project | D | F | Dx | A | P | AD | Key evidence |
|---|---|---|---|---|---|---|---|---|---|
| orphaned-local-files-shadow-drive/PLAN | bug PLAN | harbor | S | S | S | — | S | S | Gold standard: disproven theories kept with per-theory evidence; CLOSE disposition with rationale |
| delete-dropped-after-inflight-edit/PLAN | bug PLAN | harbor | S | S | S | — | S | M | Root cause from `pending_ops` row dump; `status: fixed-pending-verify` non-vocab |
| sync-state-machine/PLAN+TASKS | feature PLAN | harbor | S | S | — | — | S | S | Re-scoped post-review with breadcrumbs; TASKS carry suite counts 277→288 |
| perf-idle-coordinator/PLAN+TASKS | perf PLAN | harbor | S | S | — | — | S | M | "EXPLAIN QUERY PLAN proof before/after; no speculative indexes"; `status: planned` w/ ticked boxes |
| settings-window-rework/PLAN+TASKS | UI PLAN | harbor | M | M | — | — | S | S | Success table w/ file:line; but `[manual] all settings actions reachable+working` |
| 06-adapter-framework/PLAN+TASKS | framework PLAN | octopus | S | M | — | — | M | S | 20 numbered locked decisions w/ rejected alternatives; tasks mostly bare ticks |
| 33-tui-visual-redesign/PLAN | design PLAN | octopus | M | W | — | — | M | W | "screenshot shows the design language matches"; `status: open`, no TASKS, stalled since 05-24 |
| 47-subtasks/PLAN | feature PLAN | octopus | S | M | — | — | M | W | Rich rejected-alternatives; `status: done` w/ all 12 deliverables unchecked |
| licensing/PLAN | feature PLAN | wasabi | S | S | — | — | M | M | "magic key… is a revenue hole"; `status: in-progress` non-vocab; status prose layered on original |
| session-restore/PLAN | small PLAN | wasabi | S | S | — | — | S | S | Diagnoses actual bug; every exit criterion observable |
| archive/m1-two-webviews/PLAN+TASKS+VERIFY | archived | wasabi | S | S | — | — | S | W | "≤1GB idle… ~1021MB MET but on the line"; archived w/ TASKS active + VERIFY fully unchecked |
| archive/sleep-policy-webviewpool/PLAN+TASKS | archived | wasabi | S | S | — | — | S | S | Textbook close: pid-gone evidence, `[~]` deferrals with reasons, `verified` matches ticks |
| plugin-governance-model/PLAN | strategy PLAN | unwire | S | W | — | S | M | S | Grill-driven drop with measurement; **zero validation section** |
| primitive-atoms/PLAN+TASKS | feature PLAN | unwire | S | S | — | — | S | S | M0 kill-switch gate; TASKS correct the plan's own estimates with evidence |
| builder-agent/PLAN+TASKS | meta PLAN | spectacular | S | M | — | — | S | S | Negative test in validation; some judgment-y checks |
| spec-audit-mode/PLAN | meta PLAN | spectacular | S | S | — | — | S | S | Each signal has a paired test (dirty→warn / clean→pass) |
| harbor PRD.md | PRD | harbor | S | S | — | S | S | S | Non-goals with reasons; constraints carry verified facts w/ dates |
| harbor SPEC.md | SPEC index | harbor | — | — | — | S | S | W | Pure index (good) but "nothing built yet" vs 49 requests (never fed) |
| octopus PRD.md | PRD | octopus | M | — | — | M | W | W | Strong vision §1, degrades into schema dumps; plain-text header, never grilled |

## Appendix B — Source material

- Corpus pass, guidance audit, and comparative research were run as three independent passes; raw comparative fetches under `.scrapekit/` (sp/, gs/, os/).
- Comparison projects: [obra/superpowers](https://github.com/obra/superpowers) · [garrytan/gstack](https://github.com/garrytan/gstack) · [Fission-AI/OpenSpec](https://github.com/Fission-AI/OpenSpec).
