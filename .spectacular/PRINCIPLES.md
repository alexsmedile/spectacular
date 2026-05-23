---
version: 1.0
updated: 2026-05-21
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

- System truth → `load SPEC.md` (always) → drill into `specs/auth/SPEC.md` only if needed
- Active work → `load requests/add-team-billing/`
- Past learning → `load memory/lessons.md`

Once you see Spectacular as a clock face rather than a filing cabinet, the rest of the structure explains itself.

---

# Related

- [PRD.md](PRD.md) — vision, deliverable, goals
- [ARCHITECTURE.md](ARCHITECTURE.md) — where these principles get implemented in file structure
- [AGENTS.md](AGENTS.md) — context-loading rules that operationalize principles 1 and 6
