---
version: 1.2
updated: 2026-07-05
summary: "Operating principles for Spectacular — and how the skill enforces each at runtime"
related:
  - PRD.md
  - ARCHITECTURE.md
  - AGENTS.md
---

# Spectacular — Operating Principles

These are the beliefs Spectacular is built on. Each principle is paired with a **How the skill enforces this** hook, because principles that aren't enforced at runtime are posters.

---

## 1. Context is the system

The project structure must optimize:
- retrieval quality
- context layering
- long-running continuity
- low context entropy

Agents should load only relevant information, progressively, by capability and task. Never the entire repository.

**How the skill enforces this:**
- Routing table in `SKILL.md` maps each command to exactly one reference doc
- `AGENTS.md` lists context-loading rules per task type — skill follows them
- Status briefings read frontmatter only, never full file bodies

---

## 2. Separate intent from truth

Strategic documents define why the product exists, its goals, philosophy, and constraints. Operational documents define current behavior, implementation truth, and active work. Temporary work should never pollute canonical truth.

**How the skill enforces this:**
- PRD, PRINCIPLES, ARCHITECTURE, ROADMAP live at `.spectacular/` root and change rarely
- Active work lives in `requests/<slug>/` and is archived on completion
- Skill refuses to write request-state into root canonical docs

---

## 3. Small files over giant documents

Prefer modular capability files, layered context, and local reasoning. Avoid giant PRDs, monolithic specs, and mega prompts.

**How the skill enforces this:**
- Anti-entropy check: warns when any canonical doc passes 500 lines
- PRD-review gate fails any PRD over the target shape (200 lines)
- Reference docs split aggressively — one concern per file

---

## 4. Humans and agents share the same workspace

The repository should remain readable by humans, structured for agents, supportive of retrieval systems, and supportive of automation.

**How the skill enforces this:**
- All canonical files are plain markdown — no proprietary formats
- Frontmatter is the machine signal layer; body text is human-readable
- Skill never produces output a human couldn't write by hand

---

## 5. Operational memory compounds

The system should preserve lessons, failures, architectural traps, recurring bugs, and implementation patterns. Agents should not repeatedly rediscover solved problems.

**How the skill enforces this:**
- `.spectacular/memory/` is git-committed and team-visible
- On archive, skill reviews the completed request and proposes memory entries
- `spectacular remember this` captures insights mid-session on confirmation
- Skill never writes to memory autonomously — humans confirm every entry

---

## 6. Progressive disclosure

Load context by need, not by precaution. Reference docs load on demand from `SKILL.md` routing — never all at once.

**How the skill enforces this:**
- `SKILL.md` is a lean orchestrator (under 100 lines); reference docs in `references/` are loaded only when their trigger fires
- Status briefings load frontmatter, not bodies
- `AGENTS.md` defines per-task-type loading rules; skill respects them
- The skill never reads `archive/` during normal operation

---

## 7. Three layers: intent → execution → validation

Every meaningful unit of work passes through three layers. Skipping any one is the failure mode.

- **Intent** — what we want and why (PRD project-wide; PLAN goal per-request, compressed)
- **Execution** — what we're doing about it (TASKS, code, agent runs)
- **Validation** — how we know it worked (VERIFY, review gates, lifecycle transitions)

**How the skill enforces this:**
- PRD review gate verifies intent before a project is "ready"
- PLAN.md goal carries compressed intent down to each request
- TASKS.md frontmatter `validates:` links task groups to milestones
- Lifecycle state cannot advance to `verified` without explicit validation
- The `prd review` and request-archive flows are validation checkpoints, not afterthoughts

---

## 8. Humans decide, agents propose

Irreversible or canonical-state-changing actions require human confirmation. Reversible local edits don't. Constant approval is friction; zero approval is unsafe — the line is "would undoing this be expensive?"

**How the skill enforces this:**
- Snapshot before overwrite on every canonical doc — automatic, not opt-in
- Lifecycle transitions (planned → active → review → verified → archived) are proposed, never auto-applied
- Archive, promote, memory writes — all confirmed
- Skill flags bulk operations (>5 files) before executing
- The skill's default tone is "I propose X — confirm?", not "I did X"

---

## 9. Feedback ≠ verification ≠ benchmark

Three distinct activities, often conflated, must stay separate:

- **Verification** — "did we ship what PLAN said?" Request-scoped, confirmatory, closed-ended. Lives in `VERIFY.md`. Terminates at `verified`.
- **Feedback loop** — "was that the right thing to ship?" System-scoped, exploratory, open-ended. Lives in `.spectacular/feedback/` or `requests/<slug>/feedback/`. Compounds over time. Never terminates.
- **Benchmark** — quantitative grading against a fixed task suite. Spectacular has no benchmarks and ships none. The word "evals" is intentionally avoided — it carries HumanEval/MMLU/accuracy-% baggage and pulls the wrong way.

A feature can be `verified` while feedback reveals we built the wrong thing. Feedback can be glowing while VERIFY fails on one missed assertion. The axes are orthogonal — both must be exercised on user-visible changes.

**How the skill enforces this:**
- `feedback-loop` is a distinct mode; never auto-runs as part of `verify` or `review`
- Skill surfaces feedback-loop offers at three checkpoints only (milestone complete, request → review, archive flow) — never mid-flow
- Feedback files use the word "feedback" — never "eval", "benchmark", "score", or grading language
- `doctor feedback` is judgment-only (no `--fix`); feedback is never auto-resolved

---

## 10. Build the smallest verified slice, full scope in mind

The highest-impact version of a thing is rarely the most complete one. Default to the minimum slice that delivers the core value now — then verify it, learn, and extend. Build *less*, but build it as a finished block, not a stub: scoped down, not half-done. Hold the full scope in mind so today's slice stays future-proof; ship only the part you actually need today.

Over-engineering is the failure mode — speculative generality, abstractions for a second case that may never come, features nobody has asked for yet. An agent's instinct is to *build*; the discipline is to choose to build less, and to record what was deferred rather than building it now "while we're here."

**How the skill enforces this:**
- `@Planning` policy `scope-down` (warn) asks for the smallest high-impact slice before milestones are fixed
- PLAN milestones are demoable outcomes, not an exhaustive feature list (slot 3 rejects task-lists)
- Deferred scope goes to ROADMAP as explicit `v2+` — out-of-scope is recorded, not lost
- Non-goals are a first-class PRD/PLAN slot — saying no is part of the plan, not an omission

---

## 11. Earn each step — no rockets without the launchpad

Work has an order, and skipping a rung doesn't skip the work — it defers it to a worse moment. Don't reach for the moon before the rocket is built; don't start the rocket before the tooling that assembles it exists; don't run an integrity check on a thing that isn't built yet. A verification pass on a stub verifies nothing. A grand plan on an unproven foundation is a wish. Each step is a prerequisite for the next, not a parallel option — do them in the order that makes each one real.

This is the sequence complement to principle 10: #10 says build the *smallest* slice; this says build it in the *right order*. The failure mode is inverted ambition — pouring effort into the impressive far step (the moon shot, the integrity gate, the polish) while the near step it depends on is still missing, so the far step is hollow.

**How the skill enforces this:**
- Lifecycle is strictly ordered: `planned → active → review → verified` — a request cannot jump to `verified` without passing through the earlier states
- `@Implementation` policy `build-order` (warn) flags any step built on a stub, mock, or unbuilt prerequisite — build the lower layer first
- `@Implementation` policy `earn-the-verification` (warn) rejects a green check that exercises a placeholder instead of the real path
- `@Verification` policy `verification-present` (block) refuses `review → verified` while any check is unmet — `verify` drives the real flow, not the intention
- PLAN goal must exist before TASKS; TASKS before code — the doc order encodes the build order

---

# Agent Principles

Agents operating in a Spectacular workspace should:

- Load minimal viable context (per `AGENTS.md` rules)
- Prefer local capability reasoning over loading siblings
- Summarize state before any handoff (session, agent, or human)
- Avoid giant prompts — route to specific reference docs
- Preserve continuity through `.spectacular/memory/`, not chat scrollback
- Never read `archive/` during normal operation — write-only from the agent's perspective

---

# Anti-Entropy Rules

These are the maintenance rules that keep the workspace from rotting:

- Prefer spec files shorter than 500 lines
- Avoid duplicated truth — if a fact lives in two places, one is wrong
- Archive stale execution context aggressively
- Split capabilities — never let one capability spec swallow another
- Prefer explicit lifecycle states over implicit "I think this is done"
- Preserve operational memory — every archived request is a chance to learn
- Never overwrite canonical documents — snapshot first

---

# Closing — The most important insight

You are not designing a documentation system. You are designing a **temporal operational model**.

The names should express what's there. `SPEC.md` / `specs/` is what is true *now*; `requests/` is what is changing *next*; `memory/` is what was learned *before*.

- System truth → `load SPEC.md` (always) → drill into `specs/auth.md` only if needed
- Active work → `load requests/add-team-billing/`
- Past learning → `load memory/lessons.md`

Once you see Spectacular as a clock face rather than a filing cabinet, the rest of the structure explains itself.

---

# Related

- [PRD.md](PRD.md) — vision, deliverable, goals
- [ARCHITECTURE.md](ARCHITECTURE.md) — where these principles get implemented in file structure
- [AGENTS.md](AGENTS.md) — context-loading rules that operationalize principles 1 and 6
