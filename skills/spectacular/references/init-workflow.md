# Init Workflow — CLI Bootstrap

Triggered by: `spectacular init` CLI command, or first-time setup conversation.

---

## What init does

Sets up the `.spectacular/` scaffold on a new project. Run once per project.

---

## Init sequence

### Step 1 — Scaffold directory structure

Create only the minimal required directories:

```
.spectacular/
├── current/
└── requests/
```

Remaining directories are created **on demand** when first needed:

| Directory | Created when |
|---|---|
| `ideas/` | First idea file is saved |
| `skills/` | First project skill is added |
| `memory/` | First `spectacular remember this` is confirmed |
| `archive/` | First request is archived |
| `archive/ideas/` | First idea is promoted to a request |

### Step 2 — Write config.yaml

Prompt for:
- Project name
- One-line summary

Write `.spectacular/config.yaml` using the template from `scaffold-reference.md`.

### Step 3 — Create stub root files

Create with frontmatter stubs (user fills in content). Templates live in `scaffold-reference.md`.

**Required (always created):**
- `.spectacular/PRD.md` — product intent
- `.spectacular/STACK.md` — host project tech choices
- `.spectacular/DECISIONS.md` — ADR-style log
- `.spectacular/AGENTS.md` — onboarding doc for agents in this workspace

**Optional but recommended (offer during init):**
- `.spectacular/PRINCIPLES.md` — operating principles + enforcement hooks
- `.spectacular/ARCHITECTURE.md` — `.spectacular/` structure reference
- `.spectacular/ROADMAP.md` — versioned future work

In `-i` (interactive) mode, ask once: *"Create the full canonical doc set (PRD + PRINCIPLES + ARCHITECTURE + ROADMAP), or PRD only?"* Default for non-interactive: PRD only — the others can be added by the user later or via a follow-up scaffold command.

### Step 4 — Install skill

Options:
- **Project-local**: install to `.claude/skills/spectacular/` (symlink from global or copy)
- **Global**: install to `~/.claude/skills/spectacular/`

Default: project-local symlink to global install if global exists, otherwise copy.

### Step 5 — Update .gitignore

Add to `.gitignore`:
```
.spectacular.local/
```

`.spectacular/` itself is **fully committed** — all files including SESSION.md.

---

## Post-init state

After init, the project has:
- `current/` and `requests/` directories
- `config.yaml` with project name/summary
- Required stub root files (`PRD.md`, `STACK.md`, `DECISIONS.md`, `AGENTS.md`) with frontmatter
- Optional canonical docs (`PRINCIPLES.md`, `ARCHITECTURE.md`, `ROADMAP.md`) if user opted in
- Skill installed and accessible
- `.spectacular.local/` gitignored
- Other directories (`ideas/`, `memory/`, `skills/`, `archive/`) created on demand

---

## .spectacular.local/

Personal override layer — never committed.

Use for:
- Local dev overrides
- Personal config variations
- Sensitive local paths

The skill reads `.spectacular.local/` if present and merges with `.spectacular/` config, with local taking precedence.

---

## Idempotency

Init should be safe to re-run:
- Never overwrite existing files
- Never overwrite existing config.yaml
- Skip steps where files already exist, report what was skipped
- Can be used to "repair" a partial scaffold
