---
version: 1.0
updated: 2026-05-21
summary: "Spectacular roadmap — v1 status, v2 features, v3+ direction"
related:
  - PRD.md
  - ARCHITECTURE.md
---

# Spectacular — Roadmap

Time-ordered "what's next". Each item is a coarse target, not a commitment. Detail for in-flight work lives in `.spectacular/requests/<slug>/`.

---

# v1 (current)

Status: in progress.

- [x] `.spectacular/` directory convention finalized (see [ARCHITECTURE.md](ARCHITECTURE.md))
- [x] Skill (`/spectacular`) orchestrator + reference docs
- [x] PRD-craft flow (`prd / prd refine / prd review`) — see `requests/prd-craft/`
- [x] Canonical docs split (PRD / PRINCIPLES / ARCHITECTURE / ROADMAP)
- [ ] CLI bootstrap (`spectacular init`) — see `requests/cli-bootstrap/`
- [ ] First end-to-end dogfood on a fresh consumer project

---

# v2 — Workflows layer

Project-specific procedural sequences: release cycles, hotfix flows, migration procedures.

Each project handles these differently; the value is clear but the design is not yet finalized. Likely lives at `.spectacular/workflows/` with one file per workflow.

---

# v2 — Workspaces

Multiple teams or roles maintaining separate operational contexts within the same project.

**Naming convention:**
```
.spectacular/              ← default workspace (always present)
.spectacular.local/        ← personal overrides, always gitignored
.spectacular.<workspace>/  ← named team workspaces
```

Examples: `.spectacular.designteam/`, `.spectacular.devops/`, `.spectacular.builder/`.

**Rules:**
- Default workspace (`.spectacular/`) is always the base
- Named workspaces are fully committed and team-visible (same rules as default)
- `.spectacular.local/` is always gitignored regardless of workspace
- Skill reads the active workspace based on invocation context or explicit flag
- Named workspaces follow the same internal structure as the default
- Workspaces do not inherit from each other — each is independent

**Invocation (proposed):**
```
/spectacular                    ← operates on default workspace
/spectacular --workspace devops ← operates on .spectacular.devops/
spectacular status --workspace designteam
```

Workspace-switching UX is TBD.

---

# v2 — Nested workspaces

`.spectacular/` inside subdirectories of a monorepo, scoped to a specific app or package.

**Example:**
```
apps/builder/.spectacular/
apps/api/.spectacular/
packages/ui/.spectacular/
```

**Rules:**
- Nested workspaces are independent — they do not inherit from the repo-root `.spectacular/`
- The skill detects the nearest `.spectacular/` walking up from the current working directory
- Each nested workspace has its own `config.yaml`, root files, requests, memory
- Cross-workspace coordination is out of scope for v2

**Use case:** monorepos where separate teams own separate apps and want independent operational context without cross-contamination.

---

# v2 — Multi-agent orchestration

Subagent handoff conventions, parallel execution patterns, agent contracts.

**Likely scope:**
- Agent spec format (`role`, `responsibilities`, `capabilities`, `contract`)
- Handoff convention between sessions and between agents
- Spawn rules: when the main agent invokes a sub-agent

**Explicitly deferred until a real complex request exercises the need** — designing this before evidence produces the multi-agent research pipeline anti-pattern the PRD rejects.

---

# v2 — Hook-driven automation

- Auto-update `SESSION.md` on commit
- Auto-archive on merge to main
- Auto-propose lifecycle transitions when external signals fire (CI pass, PR merge, deploy)

Mechanism: Claude Code hooks + filesystem watchers. Out of scope for v1.

---

# v3+ — Context orchestration / Repository operating system

Long-term direction:

- **Smart retrieval across the full Spectacular structure** — semantic search over PRD/PLAN/TASKS/memory with citations back to source files
- **From solo builders to autonomous multi-agent engineering systems** — Spectacular as the substrate for coordinated agent teams operating on long-running products
- **Cross-project memory** — lessons that travel between projects without leaking project-specific detail

These are direction, not commitment.

---

# Related

- [PRD.md](PRD.md) — what Spectacular is and what v1 ships
- [ARCHITECTURE.md](ARCHITECTURE.md) — structures that v2+ items extend
- [PRINCIPLES.md](PRINCIPLES.md) — principles every future addition must respect
