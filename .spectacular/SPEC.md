---
version: 1.0
updated: 2026-05-23
summary: "What Spectacular actually is and does as of v0.5.0 — capabilities index"
related:
  - PRD.md
  - ARCHITECTURE.md
  - ROADMAP.md
---

# Spectacular — System Spec

## What this system is

Spectacular is an AI-native operational workspace for software projects. It ships in three layers: a `.spectacular/` directory convention, a `/spectacular` Claude Code / Codex skill that operates the workspace, and a `spectacular` CLI that bootstraps it. The skill reads frontmatter to build briefings without loading full file bodies, scaffolds requests, manages a five-state lifecycle, archives completed work, writes operational memory, and grills/refines/reviews any registered canonical doc through a single registry-driven engine.

## Capabilities

- **CLI bootstrap** — `spectacular init` scaffolds the always-set (PRD, SPEC, config, agents-file, requests/, specs/), installs the skill into `.agents/skills/spectacular/`, symlinks `.claude/skills/spectacular/`, and writes a version-pinned `skills.lock`. Idempotent — re-running fills empty stubs without overwriting content.
- **Smart-init kits** — five kits (`blank`, `coding`, `content`, `product`, `research`) declare which extra docs (`PRINCIPLES`, `ARCHITECTURE`, `ROADMAP`, `STACK`, `DECISIONS`) to scaffold via `triggers-docs`. Additive `--with <docs>` extends the kit; `--minimal` ignores kit defaults. See [[kits-contract]].
- **Doc-writing engine** — one generic `grill.md` + `refine.md` + `review.md` engine consumes [[doc-registry]] entries to handle any canonical doc verb (`spectacular <doc> grill|refine|review`). Per-doc behavior lives in `<doc>-overrides.md`. Registered v0.5.0: `prd`, `spec`, `plan`, `tasks`, `principles`, `architecture`, `roadmap`, `stack`, `agents`, `decisions`, `convention-pack`.
- **Lifecycle** — five states (`planned → active → review → verified → archived`) stored only in `PLAN.md` frontmatter. Skill proactively detects signals (all TASKS items checked → propose `review`; stale `active` requests → flag).
- **Verification 2-of-6 rule** — `requests/<slug>/VERIFY.md` only scaffolded when 2+ of 6 complexity criteria fire; otherwise verification folds into PLAN § Validation. See [[verification]].
- **Convention packs (v0.4.0+)** — opt-in repo-shape opinions declared via `config.yaml`'s `convention_pack:` block. Six rule categories (naming, taxonomy, root-files, gitignore, file-placement, project-types). Four scope locations (project-local → user → app-store → bundled). Three modes (suggest / scaffold / enforce). Bundled `minimal` pack + `alex-default` app-store reference. See [[packs-contract]].
- **Substrate doctor** — `spectacular doctor` runs read-only self-check across nine areas (`skill`, `workspace`, `frontmatter`, `snapshots`, `links`, `lifecycle`, `kits`, `conventions`, `specs`). `--fix` applies mechanical repairs (gitignore append, missing dirs, dangling symlinks, missing always-set stubs, pack-driven gitignore drift, legacy `current/` → `specs/` migration). See [[doctor]].
- **Two-layer task tracking** — harness `TaskCreate`/`TaskUpdate` for session-level micro-tracking (drives CLI live UI); on-disk `requests/<slug>/TASKS.md` for persistent milestone blocks. Anti-pattern: 1:1 duplication.
- **Memory** — `spectacular remember this` writes operational lessons to `.spectacular/memory/`. Git-committed, team-visible. Never to `.claude/` personal memory.
- **Versioning** — canonical docs (root layer + `SPEC.md` + `specs/<capability>/SPEC.md` + `config.yaml`) never overwritten in place. Snapshot as `<FILE>@vN.md` before edit; unversioned filename always points to current.
- **Distribution** — published as both Claude Code plugin marketplace (`alexsmedile/spectacular`) and Codex plugin marketplace. CLI installed via curl-one-liner (`cli/install.sh`).

## How to extend this file

- Add a bullet when a new capability ships (request → verified)
- Promote a bullet to `specs/<capability>/SPEC.md` when it grows past one line
- Snapshot before major rewrites: `spectacular snapshot .spectacular/SPEC.md`
