---
version: 2.0
updated: 2026-05-21
summary: "Onboarding doc for any agent or human landing inside .spectacular/ — how to operate, what to load, what not to touch"
related:
  - PRD.md
  - PRINCIPLES.md
  - ARCHITECTURE.md
---

# AGENTS.md — Working in `.spectacular/`

You just entered a Spectacular workspace. This doc tells you how to operate inside it without breaking the conventions.

---

## What this folder is

`.spectacular/` is an AI-native operational workspace. It separates:

- **Intent** (what we want) — `PRD.md`, `PRINCIPLES.md`
- **Structure** (how it's organized) — `ARCHITECTURE.md`, `ROADMAP.md`
- **Truth** (what currently exists) — `current/`
- **Work** (what's changing) — `requests/`
- **Memory** (what we learned) — `memory/`
- **History** (what's done) — `archive/`

Read [PRD.md](PRD.md) before planning anything. Read [ARCHITECTURE.md](ARCHITECTURE.md) if you need to understand where to put a file. Read [PRINCIPLES.md](PRINCIPLES.md) if you're unsure whether an action is allowed.

---

## How to operate

1. **Read frontmatter first**, file bodies second. Status, summary, version, related — these are the navigation layer.
2. **Load progressively.** Don't pre-load `specs/` or `requests/` wholesale; load only what the current task needs (see Context loading below). The top-level `SPEC.md` is cheap and always relevant.
3. **Snapshot before overwrite** on any canonical doc (root layer + `SPEC.md` + `specs/`). The snapshot name is `FILE@vN.md`. This is non-optional.
4. **Propose, don't act**, on irreversibles: archive, lifecycle promote, memory writes, bulk edits >5 files. Confirm with the human first.
5. **Never read `archive/`** during normal operation. It's write-only from your perspective.
6. **Write to `memory/` only on confirmation.** Never autonomously.

---

## Context loading by task

Load only what the task needs. Don't load the entire repository.

| Task type | Load |
|---|---|
| Planning / design | `PRD.md`, `PRINCIPLES.md`, `DECISIONS.md` |
| Refining intent / writing a PRD | `PRD.md`, skill references `prd-grill.md` / `prd-refine.md` / `prd-review.md` |
| Implementing a request | `STACK.md`, `requests/<slug>/PLAN.md`, `requests/<slug>/TASKS.md`, `SPEC.md`, relevant `specs/<capability>/SPEC.md` |
| Reviewing / QA | `requests/<slug>/VERIFY.md`, relevant `specs/<capability>/SPEC.md`, `requests/<slug>/RISKS.md` |
| Onboarding cold | `PRD.md`, `ARCHITECTURE.md`, this file |
| Structural questions | `ARCHITECTURE.md` only |
| Principles questions | `PRINCIPLES.md` only |
| Roadmap questions | `ROADMAP.md` only |

---

## Available skills

- **`spectacular`** — workspace management. Lean orchestrator that routes to reference docs. Trigger via `/spectacular` or commands like `spectacular new`, `spectacular archive`, `spectacular prd`, `spectacular remember this`, `spectacular snapshot`, `spectacular promote`, `spectacular status`.

When a new skill is added, update this list manually. The skill will warn if it detects a skill in `.spectacular/skills/` that isn't listed here.

---

## Creating requests

Use `spectacular new <description>` (or describe the work to `/spectacular` and let it derive the slug). The skill scaffolds `requests/<slug>/PLAN.md` + `TASKS.md` from the canonical template. Slugs are kebab-case; the skill confirms before creating.

Per-request folders contain `PLAN.md` + `TASKS.md` + (on demand) `SESSION.md`, `RISKS.md`, `VERIFY.md`, `specs/`, `artifacts/`. They **never** contain a `PRD.md` — product intent is project-wide and lives at the root.

---

## Task tracking — two layers, different scopes

Spectacular tracks work at **two distinct granularities** that must not be confused:

| Layer | Tool | Scope | Lifetime |
|---|---|---|---|
| **Milestone blocks** | On-disk `requests/<slug>/TASKS.md` | Whole request, persisted in git, team-visible | Days to weeks; survives sessions |
| **Session micro-tasks** | Harness `TaskCreate` / `TaskUpdate` (Claude Code built-in) | Current working session, ephemeral, agent-only | Minutes to hours; ends with session |

**Rule:** when starting a non-trivial session of work on any request, the agent creates harness tasks for the granular steps (file-by-file, gate-by-gate) and marks them `in_progress` / `completed` as work proceeds. On-disk `TASKS.md` continues to own the **milestone-level checklist** that the user reads to understand request status.

The two layers complement each other:
- Harness micro-tasks **silence the runtime's "task tools haven't been used" warning** while preserving the live progress signal in the CLI UI.
- On-disk `TASKS.md` items are **never duplicated** as harness tasks one-to-one. Harness tasks are *finer* — they decompose a single TASKS.md milestone into the concrete edits / commits / tests that complete it.
- When a session ends, harness tasks evaporate; their results live on in committed code + ticked on-disk TASKS items.

Anti-pattern: copying every line from `TASKS.md` into the harness one-for-one. That defeats the granularity split and creates duplicate maintenance.

---

## Don'ts

- **Don't touch `archive/`.** Write-only from your perspective.
- **Don't duplicate truth.** If a fact exists in `PRD.md`, don't restate it in a request's PLAN — link instead.
- **Don't overwrite canonical docs in place.** Always snapshot to `FILE@vN.md` first.
- **Don't write to `memory/` autonomously.** Propose; human confirms.
- **Don't create `requests/<slug>/PRD.md`.** It's an explicit anti-pattern.
- **Don't load `.spectacular/` wholesale.** Read frontmatter; load bodies on demand.
- **Don't trigger Claude Code's auto-memory** when writing to `.spectacular/memory/` — they're separate systems and double-capture is a bug.
- **Don't write into `docs/`.** That's [pageworks](https://github.com/alexsmedile/pageworks)' surface — spectacular's awareness ends at "does docs/ exist?" + "is a manifest present?". For schema, page authoring, renderer adapters, or maintenance: install pageworks and delegate. See `skills/spectacular/references/pageworks-handoff.md`.

---

## Skill boundary — spectacular vs pageworks (v1.2.0+)

Two skills, one workspace:

| Skill | Owns | Examples |
|---|---|---|
| **spectacular** | `.spectacular/` (internal workspace) | PRD, SPEC, specs/, plans, tasks, requests lifecycle, archive, doctor (workspace areas), memory |
| **pageworks** | `docs/` (public-facing) | docs.yaml schema, page authoring (Diátaxis), renderer export, doctor (docs validation), drift detection |

Spectacular's CLI keeps the `docs init`/`docs export`/etc. verbs working for backward compatibility (deprecated in v1.2.0, removed in v2.0.0). New work goes through pageworks.

After `spectacular archive <slug>` where the request touched SPEC.md, specs/, ARCHITECTURE.md, or PRD.md, spectacular prints a hint suggesting `pageworks audit` (suppress with `--no-docs-prompt` or per-project `docs.prompt_on_archive: false` in config.yaml).

---

## Handoff conventions

When ending a session or handing off to another agent:

- Update `requests/<slug>/SESSION.md` with current state, blockers, next actions
- Summarize decisions made; if any are architectural, propose a `DECISIONS.md` entry
- Surface any insights worth remembering via `spectacular remember this`

---

## Related

- [PRD.md](PRD.md) — what Spectacular is, what we're building
- [ARCHITECTURE.md](ARCHITECTURE.md) — the structure you're operating inside
- [PRINCIPLES.md](PRINCIPLES.md) — the rules behind the conventions
- [ROADMAP.md](ROADMAP.md) — what's coming next
