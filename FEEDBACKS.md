# Feedback

## Matrix Proposal Loop

**Captured:** 2026-08-13  
**Status:** Review backlog; not yet implemented  
**Skill:** `skills/matrix-proposal-loop/`

### Overall assessment

The core workflow is strong: exactly three controlled options, compact comparison, one-key
steering, proportional verification, and implementation of only the selected lineage.

The main limitation is that the current workflow assumes mostly linear movement from small to large
fragments. Real design work often stays at the same scale for several sibling decisions, starts with
a large flow and zooms inward, or combines a small fragment with a mature interactive prototype.

### Real-world pressure tests

#### Existing mood board and dashboard concept

The visual language may already be settled while metric cards, filters, charts, and navigation remain
open. Selecting one element should not automatically grow into a larger composition when other
sibling elements still need alignment.

Implications:

- Classify mood-board decisions by authority.
- Permit several decisions at the same scale.
- Interpret the baseline relative to supplied references, not only the existing codebase.

#### Existing billing table redesign

The table may already include sorting, pagination, permissions, tests, empty states, and error states.
Design proposals should preserve behavior, use representative data and states, and remain isolated
from production files until selection.

Implications:

- Create proposals in temporary previews, stories, routes, or proposal directories.
- Verify behavior equivalence and representative states.
- Run targeted checks during exploration; run full checks after execution.

#### New onboarding flow

The first useful decision may be a Level 4 flow: linear, progressive, or adaptive. After choosing it,
the designer may need to zoom into its highest-risk screen instead of implementing the complete flow.

Implication: support top-down refinement as well as bottom-up growth.

#### Design-system component

A date-range picker is a small fragment but may already be a polished interactive prototype. This
shows that fragment size and implementation maturity are separate dimensions.

Suggested model:

```text
Fragment scale:    Signal → Element → Composition → Surface → Flow
Artifact maturity: Concept → Rendered → Interactive → Production
```

### Proposed improvements

#### Priority 1 — structural

1. **Rename “fidelity ladder” to “fragment scale.”**
   The current levels describe scope rather than fidelity. Track artifact maturity separately so
   verification matches what actually exists.

2. **Replace automatic growth with a readiness decision.**
   After `A`, `B`, `C`, or `M`, allow the next action to be:
   - refine another detail at the same scale;
   - grow into a larger fragment;
   - zoom into a smaller fragment;
   - execute the selected lineage.

   Advance only when the next fragment has enough locked inputs to be meaningfully designed.

3. **Add a reference contract for mood boards and existing concepts.**
   Classify supplied guidance as:
   - **Locked:** reproduce consistently;
   - **Preferred:** preserve unless evidence supports a change;
   - **Inspirational:** interpret rather than copy;
   - **Avoid:** explicitly excluded;
   - **Unresolved:** open to exploration.

4. **Rename option A from “Native” to “Baseline.”**
   Suggested option roles:
   - **A — Baseline:** closest to the supplied concept or current product;
   - **B — Focused:** optimizes the primary goal;
   - **C — Challenger:** tests one meaningful alternative assumption.

5. **Isolate proposal artifacts.**
   Keep A/B/C variants outside production files. Promote only the selected lineage after `X`.

#### Priority 2 — workflow depth

6. **Support top-down refinement.**
   After locking a flow, allow detailing its highest-risk surface before implementation.

7. **Support same-scale sibling exploration.**
   Resolve several elements or compositions without forcing the loop upward after every selection.

8. **Make recommendations fully evidence-based.**
   Remove the default bias toward B. Recommend A, B, or C according to the stated goals, constraints,
   and preflight evidence.

9. **Add representative content and state checks.**
   When applicable, test:
   - loading, empty, error, and partial-data states;
   - long content and localization;
   - roles and permissions;
   - keyboard and reduced-motion behavior;
   - mobile and narrow containers.

10. **Define explicit reopening syntax.**
    Possible commands:
    ```text
    R@1 Reopen the Level 1 card hierarchy
    R:color Reopen only the accent decision
    ```

11. **Persist a richer alignment snapshot.**
    A one-line trace can become lossy across long sessions or multiple agents. Consider:

    ```yaml
    scale: 2
    maturity: rendered
    baseline: moodboard-v3
    locked:
      - graphite palette
      - condensed headings
      - metric card B+C
    open:
      - dashboard composition
    rejected:
      - high-saturation alerts
    artifacts:
      - proposal-a
      - proposal-b
      - proposal-c
    ```

12. **Add a convergence rule.**
    When repeated variants are no longer materially different, reframe the open decision instead of
    generating cosmetic alternatives.

### Suggested evaluation coverage

Add scenarios for:

- a mood board plus an established baseline;
- several sibling decisions at the same fragment scale;
- a Level 4 flow followed by top-down detailing;
- an existing interface whose behavior must remain equivalent;
- continuation across sessions or agents;
- reopening one earlier decision;
- loading, error, permission, localization, and long-content states;
- evidence recommending A or C rather than routinely favoring B.

### Preserve these strengths

- Exactly three options per decision.
- A compact proposal matrix.
- One-letter steering as the primary interaction.
- Controlled variation across only one or two axes.
- Proportional, truthful verification.
- Implementation of only the selected lineage.

### Recommended first revision

Implement improvements 1–5 together. They resolve the most important structural gaps while keeping
the feedback menu fast and compact. Validate them with the mood-board, existing-table, and top-down
onboarding scenarios before adding the remaining enhancements.

## Chat output is too long and too jargon-heavy

**Captured:** 2026-08-16  
**Status:** Open; affects the Skill, not the CLI  
**Source:** Live session running the P5/P6 merge and the M7 plan

### What went wrong

Spectacular's canonical records are correctly dense. That density leaked into chat.
The agent narrated intermediate work in the vocabulary of the records themselves —
claims, pass boundaries, proof requirements, fingerprints, edge kinds — to a reader
who was deciding one thing.

Two concrete instances in this session:

- A question set the user could not act on: *"I don't really know what you are
  asking, I can't understand pros and cons, benefit, advantages, differences,
  options."*
- A ready-to-start Mission plan presented as a full frontmatter walkthrough — five
  claims with both boundaries each, eight Objectives, scope, stops — roughly a
  screen and a half, to answer one question: start it or change it?

The failure is not verbosity alone. It is that a **projection for a human in chat**
was rendered at the density of a **canonical record**. The Contract already draws
that exact line for CLI output and says nothing about the agent's own prose.

### Why it is a Spectacular problem, not only a model problem

The Skill tells the agent what to derive and what to freeze. It does not tell the
agent how much of that to say out loud, or in what register. So the agent defaults
to the register of the artifacts it is reading. Any capable model reading these
documents will drift the same way.

### Proposed rules for the Skill

1. **Plan preview is a summary, not a walkthrough.** Before activation, show:
   Mission title, claim count with one-line names, the Objective graph, and the
   stops. Not pass boundaries, not proof requirements. Those are in the file and
   are read on demand.
2. **A question is asked in the reader's terms.** Options state what changes for
   the user, not which field moves. Record vocabulary appears only when the user
   introduced it.
3. **Lead with the decision.** State what is being asked, then the minimum needed
   to answer. Reasoning and verification go after, or into the file.
4. **Verification is reported as a result, not a transcript.** "Checked X, it
   holds" beats reproducing the grep.
5. **The density rule the Contract applies to CLI output applies to agent prose.**
   Canonical records stay dense; anything addressed to a human in chat is a
   projection and condenses.

### Open

- Whether these belong in `SKILL.md` directly or in a referenced style file loaded
  at launch.
- Whether the Skill can carry a worked before/after example without inflating the
  launch read.

## Owner gates were handed back as homework

### What happened

Closing M7 required four owner-gated actions: re-activation, review mode change,
completion, and merge. The agent surfaced each one separately, in its own message,
each time re-explaining the same authority boundary and handing back a command for
the owner to type. The owner's response was: "I WANT TO FINISH THIS JOB, PLS HELP
ME", and then an explicit authorization.

Even after that authorization, the agent continued to narrate what it would not do
before doing what it was asked.

### The three distinct failures

1. **Refusal instead of authorization request.** The agent said "I cannot do this,
   here is the command" where the correct move was "this is owner-gated because X.
   Authorize me? (Y/N)". The owner has to be *asked*, not *assigned*.
2. **Gates surfaced one at a time.** Four decisions arrived across four exchanges.
   They were knowable together and should have been one.
3. **The same boundary restated at every step.** Naming an authority constraint
   once is diligence; naming it five times is stalling, and it buried the one
   place a real decision was needed.

### Why it is a Spectacular problem, not only a model problem

`SKILL.md` said *when* to ask an owner (semantic fork, boundary, risk,
irreversible effect) and *what four parts* a question has. It said nothing about
the register of a gate, whether to batch them, or that an authorization carries
forward. So the agent treated each gate as a fresh refusal event and re-derived the
whole justification each time.

The Skill also said "Return to the owner" without defining what returning means.
Read literally, handing back the remaining work satisfies it.

### Applied to the Skill

Three subsections added under "Ask owner questions", plus a clarification of
"Return to the owner":

- **Ask for authorization, not for labor** — an owner-gated action needs the
  owner's decision, not their hands. Offer to act on an explicit yes. Reserve
  owner-executed steps for cases where performing the act *is* the authority.
  Record authorizer and performer separately.
- **Batch gates; ask once** — collect a phase's decisions into one exchange, and
  check first whether an existing authorization, preference, or frozen default
  already answers it.
- **State a boundary once** — name it when it first becomes relevant, then act.
  Decline in one sentence, offer the nearest available action, continue with
  everything unblocked.
- **Return to the owner** now says returning means requesting a decision and
  proceeding on the answer, not transferring work, and that everything the gate
  does not block is finished before asking.

### Still open

- Whether "an approval carries forward to the same kind of action in the same
  phase" is safe as written. It should not let a single yes authorize an unbounded
  series of irreversible effects. The current wording is scoped to same-kind and
  same-phase, but the boundary deserves a worked example.
- `git-stack` and `git-ops` were named as related surfaces. Neither exists as an
  editable skill — `git-stack` is a plugin agent. If the same batching lesson
  applies to commit/push/merge flows, it needs a home there.

## `NEXT` asks for a review that already exists

**Found:** 2026-08-16, closing M7. **Where:** `internal/missionbundle/derive.go:193`.

After `review record` wrote RV1 and the Mission carried a passing review pointer,
`mission show M7` still printed:

```
NEXT: every Objective is implemented; record a review
```

`nextAction` selects that branch on `state.Done == len(b.Objectives)` alone. It
never reads `b.Reviews`, so a recorded review cannot retire the instruction. The
same line is emitted before and after the review exists.

### Why it matters more than a wrong string

M7 froze `state-line` on the boundary that every field is *derived from data the
bundle already carries*. The bundle carries `reviews` with verdicts. The NEXT line
ignores it, so the projection disagrees with the record it summarizes — which is
also, verbatim, an M7 stop condition:

> A state line, drift flag, or readiness conclusion disagrees with the record it
> summarizes.

The stop did not fire because nothing tests NEXT against a bundle that has both
implemented Objectives and a recorded review. The M7 golden fixtures cover the
pre-review state only.

### What the correct behavior probably is

Once every Objective is implemented and at least one review is bound to the
current activation fingerprint with a passing verdict, NEXT should name owner
completion, not another review. A review bound to a *stale* fingerprint should
keep asking, since that is exactly the drift the fingerprint exists to catch.

That distinction is a real design question, not a typo: it decides whether a
passing review survives a boundary amendment. It is written here rather than
patched silently.

### Not fixed here

M8 is active and froze its own scope. This is not in it, and amending an active
Mission's scope to absorb a found bug is the drift M8 exists to prevent. It waits
for the cleanup Mission — where it belongs alongside the unaddressed
"review the open gaps and dead weight" ask.
