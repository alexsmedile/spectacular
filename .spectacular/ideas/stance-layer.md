---
type: idea
status: parked
priority: medium
owner: alex
updated: 2026-07-10
origin: (captured from a scattered-notes triage session, 2026-07-10 — six improvement notes sorted; two developed here)
promoted_to: null
related:
  - ../PRINCIPLES.md
  - ../POLICY.md
  - ../decisions/index.md
  - ../specs/index.md
---

# Idea — the stance layer (architectural posture + MVP↔perfection knob)

> **Origin (2026-07-10):** six scattered improvement notes were triaged against the live
> workspace. Four turned out to be *already-solved* (they're the spine of the system) and are
> logged below as audit prompts, not new work. Two are genuine gaps worth building, and they
> share a theme — both are about **stance**: the thinking posture an agent takes (#5) and the
> strictness dial the project runs at (#6). This doc develops those two and parks the rest.

> **DECISIONS LOCKED (2026-07-11, via grill-me).** Every open decision below is now resolved.
> Summary — full rationale inline in each part:
> - **#5 architectural-stance** → a new `@Planning` **warn** policy that fires **only on a real
>   architectural fork** (crosses a module boundary / sets a precedent / two viable structures);
>   when it fires it **offers** `spectacular decide`, never forces it. Buildable — see Part 2.
> - **#6 grade knob** → **collapsed to a label only.** The severity-shifting *dial* was
>   **rejected as over-engineering** (config overrides already tune gates; a computed
>   `declared ± grade ± override` severity would *hurt* the mixed-model legibility we're fixing
>   in Part 4). Kept: a per-request `grade:` **label** — `prototype | mvp | standard |
>   production` — that `status` renders so prototype-green never reads as production-green.
>   No gate ever changes behavior on it. This dissolves Part 3's open-decisions #3/#4/#6 and Q9
>   (the "does prototype soften the blockers" safety fork) entirely — there's no softening.
> - **Part 4 legibility** → **L1 + L3 + L4 + L6 build now** (Wave 1); **L5** immediately after
>   (Wave 2); **L2 deferred** to its own scoped pass (it's the only item touching the
>   machine-read `check:` line). Marker/table syntax locked in Part 4.
> - **Sequencing** → the legibility patch **ships this session** (low-risk doc edits, fully
>   decided). The stance layer is **recorded + cut as a request**, not built this session.

---

## Part 1 — the six-note triage (what's already solved vs. what's a gap)

| # | Note | Verdict | Where it already lives |
|---|------|---------|------------------------|
| 1 | Manage durable instructions vs. mutable information | **Solved** — it's Principle 2 (*separate intent from truth*). Durable: PRD/PRINCIPLES/ARCHITECTURE/STACK. Mutable: `requests/`, specs, roadmaps, sessions. Practice layer: POLICY. | PRINCIPLES §2; `doctor` enforces "no request-state in root docs" |
| 2 | Eternal files as source of truth, consumed by scripts/skills | **Solved & actively sharpened** — the *generator-as-source-of-truth* pattern. `status` renders the fleet from `PLAN.md` frontmatter (the hand-cached CLAUDE.md table was **deleted** because it drifted — status-fleet-view, b23). `policy` reads POLICY.md; ROADMAP ledger is the single build→version source. | `spectacular status`, `spectacular policy`, ROADMAP ledger |
| 3 | Modular modules — `@`-tag a file, acts like a ref doc | **Mostly solved** — the whole context-loading model is this: `[[wikilinks]]`, the routing table, `references/<doc>-rules.md` on demand, `@<hook>` policy retrieval. **Gap:** it's all *skill-internal*; no first-class user/agent-facing "compose these N modules into a prompt" primitive. | SKILL.md routing table; POLICY `@<hook>` |
| 4 | How do I know if splitting docs + refs is a *good* approach? | **Genuine gap** — there's a *bias* (Principle 3) and one crude guard (500-line warning) but **no fitness function**. Over-splitting has a real cost (retrieval hops, orphaned refs) that nothing measures. | *parked separately — see note below* |
| 5 | Prompt like a senior dev → better architectural decisions | **Gap worth building** — the effect is *observed* but only thinly encoded. `understand-before-change` is the nearest hook; there's no policy that injects an **architectural-thinking stance** at prompt time. → **Part 2** | developed below |
| 6 | Rapid prototyping / MVP vs. perfection | **Gap worth building** — the prototyping half is well-built (feedback-loop mode, PRINCIPLES §9). Missing: an explicit **switchable strictness stance** so gate severity scales with the work's grade instead of being global. → **Part 3** | developed below |
| 7 | Gates | **Solved — strongest subsystem.** POLICY hooks, archive closure gate, 2-of-6 verification, `understand-before-change`, 18-check roadmap review. Note 6 is really "make gate *severity* context-dependent," which is the one open edge. | POLICY, closure gate, verify 2-of-6 |

**Not-developed-here but flagged:**
- **#4 (doc-split fitness function)** is the most interesting of the "unpicked" gaps — a metric + `doctor` check for whether a split *helped*: ref-hop count to answer a task, orphaned/dangling links, retrieval locality. Turns "small files" from a bias into a measurable. Cut a separate idea if pursued.
- **#3 (explicit `@`-module composition)** — a real seam but larger; would want its own request.

---

## Part 2 — Architectural-stance policy (#5)

### The observation
When the user prompts "act / decide like a senior dev would — think in terms of code and
architecture," agent output measurably improves: fewer quick-and-local hacks, more decisions
that account for the broader shape. The effect is real but currently depends on the user
remembering to say it every time. It should be a **runtime stance the workspace injects**, not
a hope.

### Why it's not already covered
- `understand-before-change` (@Implementation) gates *comprehension* — "did you read the code
  first" — not *architectural judgment*.
- `decisions/` (ADR log) is the **output home** for architectural decisions. What's missing is
  the **prompt-time stance** that *produces* good ADRs in the first place. Right now an agent
  can barrel into a request and never surface that a decision was even made.

### Sketch — a `@Planning` / `@Implementation` policy
Add a policy that injects an architectural posture and, crucially, ties it to the existing
`decisions/` substrate so the stance produces a durable artifact rather than evaporating.

```
## @Planning   (or @Implementation — decide which phase; likely @Planning)

### architectural-stance
- principle: (new or extend §… on judgment)
- severity: warn
Before proposing an approach, take the posture of a senior engineer responsible for the
whole system, not just this change. State the architectural decision being made and at least
one alternative you rejected and why. If the change sets a precedent, crosses a boundary, or
picks between two viable structures, that decision is worth an ADR — offer to `spectacular
decide "<text>"` it. Do not silently pick the locally-convenient option.
```

### Design questions — RESOLVED (2026-07-11, grill)
1. **Which hook?** → **`@Planning`, single hook.** Decision exists before code; sits with
   `scope-down` + `milestones-in-build-order` (the "shape the work before starting" family).
   Not @Implementation (too late), not both (dilutes).
2. **Block or warn?** → **warn.** A thinking posture, not a mechanically-checkable gate.
3. **Auto-ADR coupling.** → **offer, never require.** When the stance detects a real decision it
   offers `spectacular decide`; the human confirms (Principle 8). Require was rejected — it
   contradicts warn and can't be mechanically gated. (The #6 escalation that would have made
   production *require* the ADR is moot: #6 is now a label with no severity effect.)
4. **Avoid ceremony inflation.** → **resolved: conditional trigger.** The `check:` fires *only*
   when the change crosses a module boundary, sets a precedent future work follows, or picks
   between two viable structures — trivial edits pass silently. Same discipline as
   `ceremony-matches-uncertainty` @Debugging. Rejected the always-on "narrate alternatives"
   variant (it manufactures fake tradeoffs).

### Final policy shape (buildable)
```
## @Planning

### architectural-stance
- principle: (tag 11, or a new judgment principle — decide at build)
- severity: warn
- check: the plan names an architectural decision + one rejected alternative WHEN the change
  crosses a module boundary, sets a precedent future work will follow, or picks between two
  viable structures; trivial/local edits are exempt

Before fixing the approach, take the posture of a senior engineer responsible for the whole
system. Ask: does this change cross a module boundary, set a precedent others will copy, or
pick between two viable structures? If yes, name the decision and one alternative you rejected
and why, and offer to record it with `spectacular decide`. If no, proceed — not every edit is
an architectural decision, and manufacturing a fake alternative is its own noise. Surface the
decision; the human confirms the ADR.
```
*(Author it legibly from birth — this policy ships after the Part 4 legibility patch, so it
gets an L5 override clause and, being a warn, no L1 marker.)*

### Cost
Small. One POLICY block + wiring the "offer to `decide`" prompt (the `decide` verb already
exists). No CLI change strictly required.

---

## Part 3 — MVP ↔ perfection knob (#6)

### The observation
Gate severity today is **global**: a policy is `block` or `warn` for the whole project. But the
right strictness depends on the *grade of the work*. A throwaway prototype and a
production-critical change should not face the same gates. Prompt-time you already flex this by
hand ("this is just an MVP, don't over-engineer" vs "this is production, be rigorous"). That
flex should be a **first-class, switchable stance**, not tone-of-voice.

### Why it's not already covered
- The **prototyping *mode*** exists (feedback-loop, v1.6.0; PRINCIPLES §9 separates feedback
  from verification) — but that's a *workflow*, not a *strictness dial*.
- POLICY severity is fixed per-policy. There's no axis that says "run every gate one notch
  softer for this request because it's exploratory."

### RESOLVED (2026-07-11, grill) — the severity *dial* was rejected; keep only the label

The original sketch was a **dial**: a `grade` field that shifts every gate's severity one notch
(`prototype` softer, `production` harder). **We rejected the dial.** Two reasons:

1. **The problem it solved is already solved.** Gates are *already* tunable per-project through
   `config.yaml` overrides (`_policy_cfg_get`). If a repo is prototype-grade, soften the two or
   three gates that matter directly. A blanket dial is a second mechanism for the same job.
2. **It would *hurt* the Part 4 legibility work.** A dial makes a gate's real severity
   `declared ± grade ± config-override` — three inputs to reason about. We're spending Part 4
   making severity *more* legible to weak models; a computed severity pulls the opposite way.

**What we kept — the honesty half, a label only.** The one thing the dial genuinely bought that
config overrides don't: a request that runs at prototype rigor should *say so*, so nobody
mistakes prototype-green for production-green. That's a **label**, not a mechanism.

```yaml
# requests/<slug>/PLAN.md frontmatter
grade: mvp        # prototype | mvp | standard | production   (absent ⇒ standard)
```

Four-rung ladder (user-chosen — sharper than a binary, maps onto the system's own language):

| grade | meaning |
|---|---|
| `prototype` | throwaway / spike — learning only, not meant to ship (feedback-loop territory) |
| `mvp` | the smallest **shippable** slice — Principle 10's "finished block, scoped down, not half-done" |
| `standard` | normal production work — **today's default**; `absent` resolves here |
| `production` | production-critical — highest rigor expected |

**Behavior contract — this is the whole feature:**
- **`status` renders it** — fleet row, request card, and `--json` (`grade` field). A prototype
  passing green must *read* as prototype-green.
- **No gate ever changes behavior on it.** It is informational. Strictness stays where it already
  lives — the declared `severity:` + any `config.yaml` override. This is the deliberate cut.
- **`doctor lifecycle` warns on any value outside the four** (closed-enum check, exactly like the
  existing `status:` / `hold:` enum checks) so `grade: protoype` typos don't render as garbage.

**Dissolved by this decision:** Part 3 old open-decisions #3 (what a grade does to a policy),
#4 (closure-gate / verification-2-of-6 interaction), #6 (escalation re-arming gates), and the
whole Q9 "does prototype soften the ⛔ blockers" safety fork — all were about the severity
mechanism, which no longer exists. Old #1 (where it lives) → **PLAN frontmatter only**; a
`config.yaml` project default was dropped as unneeded for a pure label. Old #2 (how many rungs)
→ **four**, above.

### Why #5 and #6 no longer "belong together" mechanically
The original pitch was that one `grade` setting would flex *both* the thinking posture (#5) and
the enforcement posture (#6). With #6 reduced to a label that flexes nothing, that coupling is
gone — they're now two independent, small additions that happen to have been triaged together.
Ship them in one request for convenience, not because one depends on the other.

### Cost
Small. A frontmatter field + `status` render paths (fleet/card/json) + one closed-enum check in
`doctor lifecycle`. No policy-engine change at all (the rejected dial was the only thing that
would have touched `_policy_records`).

---

## Suggested next step

Cut **one narrow request** — "stance-layer" — shipping two independent, small additions:
1. **`architectural-stance`** — a `@Planning` warn policy, conditional trigger, offers
   `spectacular decide` (final shape in Part 2). Authored *after* the Part 4 legibility patch so
   it's born legible (gets an L5 override clause; warn ⇒ no L1 marker).
2. **`grade` label** — `prototype|mvp|standard|production` in PLAN frontmatter; `status` renders
   it; `doctor lifecycle` closed-enum-warns. No policy-engine change (final shape in Part 3).

The severity *dial* is explicitly out of scope (rejected). Keep #4 (doc-split fitness function)
and #3 (`@`-module composition) as separate parked ideas — they don't share this layer.

**This request is recorded but NOT built this session** — only the Part 4 legibility patch ships
now. See the [[stance-layer]] request PLAN once cut.

---

## Part 4 — Mixed-model-tier legibility of the rule layer (PRINCIPLES + POLICY)

> **Separate track from the stance layer (Parts 2–3).** This is not a new capability — it's a
> **legibility patch** to already-shipped material. The finding: PRINCIPLES.md + POLICY.md are
> written *at frontier level*. They assume the reader resolves cross-references, holds 11
> principles + 19 policies in working memory, and infers action from prose. A frontier model
> (Opus-class) does this; a **lower-intelligence model pattern-matches surface tokens and skims
> past the gate.** Every fix below is the same move — *make the implicit explicit at the point
> of use* — so a weaker model can't miss it while a frontier model reads the redundancy as
> harmless confirmation. No content softens; signal that's currently one hop away moves inline.

### The six findings (highest-leverage first)

| # | Finding | Who it bites | Fix | Cost |
|---|---------|--------------|-----|------|
| L1 | **Block-vs-warn is invisible to a skimmer.** The load-bearing fact ("blocks ONLY if it says `severity: block`") is in an HTML comment. Only 4/19 policies block; a weak model reads `- severity: warn` as one dashed line among four and can't tell a hard gate from a nudge. | Lower tier | Make severity unmissable at point of use — `### understand-before-change 🔒 BLOCK` vs `### scope-down · warn`, or a leading `> ⛔ THIS GATE REFUSES` on the 4 blockers. | Trivial, additive |
| L2 | **`check:` is written to be human-verified, not model-executed.** e.g. `check: milestones are ordered by dependency`. A frontier model turns it into a grep action; a weak model reads it as description. | Lower tier | Phrase checks as an imperative the model runs: "Grep `### M` headings; confirm no listed dependency appears in a later milestone." | Mechanical, per-policy |
| L3 | **Principle 10 vs 11 conflate.** Both are restraint; both cross-cite. Strong model keeps them separate (10 = *how much* / scope, 11 = *what order* / sequence); weak model collapses to "do less" and loses the ordering discipline that `build-order` + `milestones-in-build-order` depend on. | Lower tier | Two-line contrast box atop each: `#10: smallest slice (scope). #11: right order (sequence). Different failure modes.` | Cheap |
| L4 | **Excuse/Reality tables — your best weak-model tool — exist only @Debugging.** Those tables pre-empt the specific rationalization a model reaches for. The 4 debugging policies have them; `scope-down`, `build-order`, `understand-before-change` — the *most*-rationalized-away policies — have none. | Lower tier | Add a one-row Excuse/Reality + "Red flag — stop:" line to the 3–4 highest-value non-debugging warn policies. **Highest single leverage for weak agents in the set.** | Small |
| L5 | **Nothing tells a *strong* model when a warn is safe to override.** Docs are all "here's the rule." `commit-checkpoint` handles it well ("commit or explicitly defer, with a reason"); most don't. A frontier model either over-complies (treats warn as block) or invents its own bar. | Frontier | For warn policies, one clause naming the legitimate skip: "Override when X; record why." | Small, per-policy |
| L6 | **`Law:` lines assert but don't route.** `Law: no diagnosis before the fixes corpus is searched` — a weak model doesn't connect it to the mechanism (`spectacular fix list` + signature grep) one bullet up. | Lower tier | Inline the mechanism into the Law: "Law: no diagnosis before `spectacular fix list` + a signature grep has run." | Trivial |

### Recommended cut
Ship **L1 + L4 now** — highest-leverage, purely additive, testable, improves every future agent
run regardless of the stance layer. L2/L3/L6 are cheap follow-ons; L5 is the one frontier-facing
item and can ride along. This is a **separate patch from the stance layer** — it touches the same
two files but has no dependency on `grade` or `architectural-stance`.

### Open decisions — RESOLVED (2026-07-11, grill)
1. **Severity marker syntax** (L1) → **body line, not heading.** `_policy_records` parses the
   policy **id** from the whole `### …` heading, so any heading suffix would corrupt the id and
   break `_policy_lookup` + config overrides. Marker goes on a body line under the heading:
   `> ⛔ **BLOCKING** — refuses … until satisfied.` on the **4 blockers only**; warns unmarked.
   Wording flexes per-gate (transition / write / overwrite). `⛔` chosen over `🔒` (lock reads as
   "secure," ambiguous). The machine truth stays the `- severity: block` line; the marker is its
   legible echo.
2. **How many policies get Excuse/Reality** (L4) → **4:** `scope-down`, `build-order`,
   `earn-the-verification`, `understand-before-change`. Full Law + table + Red-flag shape (matches
   @Debugging). Reserved for high-temptation policies — not every warn.
3. **`check:` phrasing wholesale** (L2) → **deferred to its own scoped pass.** It's the only item
   touching the machine-read `check:` line, and "which checks are genuinely mechanizable vs
   judgment" is unresolved. Not in this patch. (New checks we author — e.g. architectural-stance —
   are written imperatively from birth; that's not the same as rewriting existing ones.)
4. **L1 redundant with a CLI `--blocks-only` view?** → **No — doc edit is the fix.** The finding
   is "a model *reading the file* can't see the gate." A CLI flag doesn't help a model reading
   POLICY.md directly (which the skill injects on hook entry). Rejected the CLI-only route.
5. **Sequencing vs the stance layer** → **legibility ships first, this session; stance layer is
   recorded + cut as a request, built later** — so its new policies are written legibly from
   birth. Within legibility: **Wave 1 = L1+L3+L4+L6**, **Wave 2 = L5**, **L2 deferred**.
