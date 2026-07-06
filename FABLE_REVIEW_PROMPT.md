# Fable Review Prompt — Spectacular Spec/Planning Quality Uplift

> Run this in a Fable session with a working directory of
> `/Users/alex/vault/data/skills_db/spectacular`. Copy everything below the line
> as the prompt. It is self-contained: it names the files to read, the real-world
> corpus to study, the metrics to move, and the exact deliverable shape.

---

You are reviewing **Spectacular**, an AI-native operational-workspace framework
for software projects. It ships as three coupled layers — a `.spectacular/`
directory convention, a Claude Code skill (`/spectacular`) that operates the
workspace, and a Bash CLI (`spectacular`) that is the deterministic mutator. The
skill is a *lean orchestrator*: `skills/spectacular/SKILL.md` reads triggers and
routes to one on-demand reference doc, which carries the actual instructions.

**Your one job:** make the skill produce *measurably better* output when it does
the three things it exists to do — **spec writing, planning, and strategy.**
Not: tidy the prose, not: add features, not: rename things. The bar is
"does this change make the next PLAN.md / SPEC.md / PRD.md the skill writes
sharper, more decision-dense, and less likely to drift or hand-wave."

## What "better" means here (the metrics you are moving)

Rank every proposed change against these. A change that doesn't move at least
one is out of scope — say so and drop it.

1. **Decision density** — specs/plans state *decisions with rationale*, not
   restatements of the request. Fewer "we will consider…" and more "X, because
   Y, ruling out Z."
2. **Falsifiability** — every plan carries checks that can actually fail
   (executable / assertable / observable), not vibes. The 2-of-6 verification
   rule (`references/verify.md`) is the existing lever — is it working?
3. **Diagnosis discipline** — for bug/investigation plans: hypotheses are
   ranked, each has a confirm/disprove check, and disproven ones are *kept with
   the evidence that killed them*. (The corpus has a gold-standard example —
   see below.)
4. **Strategic altitude** — PRD/ROADMAP/DECISIONS work reasons about *what &
   why & for whom* and sequencing, not just task lists. Does the skill push the
   human up to intent, or let them stay at the checkbox level?
5. **Progressive-disclosure fidelity** — the skill's own principle is small
   files, load only what's needed. Where does it over-load context (dulling the
   model) or under-load (missing the doc that would have made the output sharp)?
6. **Anti-drift** — canonical docs get snapshotted not overwritten; lifecycle
   state lives in one place; SPEC.md stays an index. Where does the guidance
   *let* the model produce drift?

## Read these first (the system under review)

Core skill + routing:
- `skills/spectacular/SKILL.md` — the orchestrator + full routing table
- `skills/spectacular/references/doc-index.md` — catalog of every doc type

The three target flows — read the reference docs that drive spec/plan/strategy:
- `references/grill.md`, `references/refine.md`, `references/review.md` — the
  generic engine for any doc
- `references/prd-rules.md`, `references/plan-rules.md`, `references/tasks-rules.md`,
  `references/spec-rules.md`, `references/roadmap-rules.md`, `references/decisions-rules.md`
- `references/new-request.md`, `references/active-request.md` — where plans are born
- `references/verify.md` — the falsifiability lever (2-of-6 rule)
- `references/bug-workflow.md` + `references/imagine.md` — the two planning axes
  (spec-driven diagnosis vs imagination-backed vision)

The CLI (skim, don't audit line-by-line — you're checking whether the mutators
*constrain the model toward good output*, not hunting Bash bugs):
- `cli/spectacular` — grep for the doc scaffolding paths (`doc_prd`, `doc_spec`,
  `doc_plan`…) and the frontmatter stubs they emit.

## Study real-world usage — this is the point, don't skip it

Sixteen real projects use `.spectacular/`. Study these to find *what the skill
actually produces in the wild* vs what the docs promise. Look for the gap.

**Heavy users (mine for patterns — good and bad):**
- `/Users/alex/code/apps/harbor/.spectacular/requests/` — 49 requests, a
  shipping macOS file-provider app. Read 4–5 PLAN.md files here.
- `/Users/alex/vault/data/skills_db/octopus/.spectacular/requests/` — 40 requests.
- `/Users/alex/code/apps/wasabi/.spectacular/` — 7 active + 6 **archived**
  (full lifecycle: planned→verified→archived).
- `/Users/alex/vault/data/skills_db/unwire/.spectacular/` — 19 reqs + 3 archived.
- `/Users/alex/vault/data/skills_db/spectacular/.spectacular/requests/` — the
  framework managing *itself*; the most disciplined corpus.

**Gold-standard artifact — read this in full, it is the bar:**
`/Users/alex/code/apps/harbor/.spectacular/requests/orphaned-local-files-shadow-drive/PLAN.md`
A bug plan that got its diagnosis *corrected and then sharpened across three
dates*, with a full block of disproven store-side causes kept as history and a
cheap evidence-giving fix direction. This is what metric #3 (diagnosis
discipline) looks like when it works. **Find where the skill's guidance would
NOT reliably produce this** — that gap is your highest-value finding.

**Method for the corpus pass:** sample ~15 PLAN.md/SPEC.md/PRD.md across the
heavy users. For each, ask: is this decision-dense or a restated request? Does
it carry falsifiable checks? For bug plans, is the diagnosis ranked+disproven or
a single guess? Tabulate what you see. The *distribution* of quality tells you
where the guidance is load-bearing and where it's ornamental.

## Deliverable

Write a single markdown report to
`/Users/alex/vault/data/skills_db/spectacular/docs/reviews/fable-spec-quality-review.md`.
Structure:

1. **Corpus findings** (½ page) — the quality distribution you observed across
   the ~15 sampled artifacts, per metric. Concrete: name files, quote lines.
   Where does real output already hit the bar? Where does it fall short, and is
   the shortfall a *guidance gap* (fixable in the skill) or a *user gap*?

2. **Ranked change list** — the improvements, most-impactful first. Each entry:
   - **What** — the specific edit (which file, which section, what changes).
   - **Which metric(s)** it moves and *why you believe it moves them* — tie it
     to a corpus observation, not a hunch.
   - **Blast radius** — is this a surgical wording change or a structural shift?
     Prefer surgical. Spectacular is at v1.26 and self-hosting; a change that
     forces re-writing every reference doc is almost never worth it.
   - **Concrete before/after** for the top 5 changes — show the actual
     prose/rule as it is now and as you'd write it.

3. **Explicit non-changes** — things you considered and rejected, with the
   one-line reason. This is as valuable as the change list; it stops the next
   reviewer re-litigating.

4. **One highest-leverage change** — if only one edit ships, which, and the
   single sentence for why.

## Constraints on your recommendations

- **Bias to deletion and sharpening over addition.** The skill is already large.
  A rule that makes an existing instruction *bite harder* beats a new rule.
- **Respect progressive disclosure.** Don't propose folding reference docs into
  SKILL.md — that breaks the load-only-what's-needed architecture on purpose.
- **Every recommendation must be falsifiable itself** — "this would improve
  clarity" is not actionable; "plans currently skip the disprove-check because
  `plan-rules.md` frames hypotheses as optional — make them required with the
  harbor example inlined" is.
- **Don't touch the CLI's mutation semantics.** You may recommend better
  frontmatter *stubs* or scaffolded prompts, but the CLI-is-the-mutator contract
  is load-bearing and out of scope.
- If you find the guidance is already good and the shortfall is user discipline,
  **say that** — a report that honestly finds "the docs are fine, the lever that's
  missing is X" is more useful than a padded change list.

Do the corpus pass before you write a single recommendation. The recommendations
are only as good as the gap you actually found in real output.
