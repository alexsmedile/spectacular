---
status: planned
priority: medium
owner: alex
updated: 2026-05-23
summary: "Rename current/ → specs/, introduce .spectacular/SPEC.md as the always-on index; per-capability files become optional and only for complex work"
related:
  - ../../SPEC.md   # created by this request
  - ../public-docs-foundation/PLAN.md
  - ../../ARCHITECTURE.md
---

# Plan — Spec Rename

## Goal

Replace the `current/` folder convention with `specs/`, and add `.spectacular/SPEC.md` as a short, always-present index of how the system works *right now*. Per-capability spec files (`specs/<capability>/SPEC.md`) become **optional** — used only when a capability is complex enough to warrant its own surface, not by default.

## Why

Three accumulated problems:

1. **AI mis-routing** — "current" is a temporal word, not a content word. Agents (and humans) keep treating `current/` as recency state rather than system truth. TODO.md line 2 flags this directly: "ai has difficulties using current, maybe consider rename to specs/ and SPEC.md".
2. **Folder-first overhead** — today the convention assumes capabilities each get a folder. For small projects there are no capabilities yet, so `current/` sits empty (verified: the live `.spectacular/current/` in this repo is empty). The shape implies more structure than most projects need.
3. **Missing top-level surface** — there's no single "what is true about this system right now" doc. PRD says *what we want*; ARCHITECTURE says *how the workspace is shaped*; nothing says *what's actually built and how it actually behaves*. `SPEC.md` fills that gap as a 1-file index that can grow into a folder only when needed.

This also unblocks the spec-vs-doc clarification work (see `public-docs-foundation/`) — without a clean "spec" word, the docs/specs distinction stays muddy.

## Scope

**In scope**
- Folder rename: `.spectacular/current/` → `.spectacular/specs/`
- New always-on file: `.spectacular/SPEC.md` (index — bullet list of capabilities + 1-line summary each, or "no capabilities yet")
- `specs/<capability>/SPEC.md` becomes **optional** — created only when a capability needs more than a bullet in the index
- CLI: `spectacular init` scaffolds `SPEC.md` (always) + `specs/` (empty dir, .gitkeep)
- CLI: `spectacular doctor` migration — when `current/` exists with no `specs/`, propose mechanical rename via `--fix`
- One-cycle backwards compat: if `current/` exists alongside `specs/`, doctor flags as error (manual merge required)
- Doc-registry entry for `spec` doc type (template + slots + mode + location)
- Skill: doc verbs work uniformly — `spectacular spec`, `spectacular spec grill`, `spectacular spec review`
- Update all references: `ARCHITECTURE.md`, `AGENTS.md`, `skills/spectacular/references/*.md`, `docs/scaffold.md`, `README.md`, `CLAUDE.md`
- This repo's own `.spectacular/current/` → `.spectacular/specs/` migration (dogfood)

**Out of scope (deferred)**
- Auto-generated SPEC.md from request lifecycle (e.g., "verified request adds a line to SPEC.md") — flagged for v2
- Cross-linking specs ↔ requests as a graph
- Spec versioning beyond existing snapshot convention (`SPEC@v1.md` follows existing rules)
- Convention-pack `specs-layout` rule category

## Decisions

- **`SPEC.md` not `STATE.md` or `CURRENT.md`** — "spec" is the industry word; ties cleanly to "specification". Users immediately know what to put there.
- **Index-first, folder-optional** — small projects get 1 file; complex projects grow into `specs/<capability>/SPEC.md` per capability. No forced taxonomy.
- **No `current/` alias kept indefinitely** — one minor cycle of doctor migration support, then `current/` is treated as drift. Avoids two-name confusion long-term.
- **`SPEC.md` is canonical** — same snapshot rules as PRD/ARCHITECTURE (`SPEC@v1.md` etc.). Lives in always-set.
- **Per-capability files use bare `SPEC.md` filename** — not `<capability>-spec.md`. Folder name carries the capability identity; file name carries the type.

## Lifecycle impact

- Frontmatter unchanged — `spec` doc type uses the same registry-driven engine as PRD/PLAN/TASKS/etc.
- Context loading rules in `AGENTS.md` table: replace `current/<capability>` with `specs/<capability>/SPEC.md` (or just `SPEC.md` for index-only)
- Doctor adds `specs` area: validates SPEC.md exists, frontmatter parseable, capability subfolders (if any) each have a SPEC.md

## Validation

- Init in a fresh tmp dir → SPEC.md present, specs/ exists with .gitkeep, config.yaml unchanged
- Init with `--kit coding` → same plus STACK/ARCHITECTURE; SPEC.md still scaffolded
- Init in dir with existing `current/` → doctor flags migration available
- `spectacular doctor specs --fix` on legacy layout → renames current/ → specs/, preserves contents
- Pre-existing `specs/` + `current/` both present → doctor errors, refuses auto-fix
- Skill `/spectacular` briefing reads SPEC.md as part of system-truth context
- All references in repo updated; no remaining `current/` mentions outside archive/migration notes
- This repo's own dogfood migration committed
- Test suite extended: `tests/cli/specs.test.sh` covering 6 scenarios above

## Milestones

1. **M1 — Schema + scaffold** — SPEC.md template, doc-registry entry, init scaffolds it
2. **M2 — Doctor migration** — `specs` area + `--fix` handler for current → specs rename
3. **M3 — Reference updates** — all skill refs, ARCHITECTURE, AGENTS, README, docs/
4. **M4 — Dogfood** — migrate this repo's `current/` (empty, mechanical) + write the first real SPEC.md for spectacular itself
5. **M5 — Test + release** — tests/cli/specs.test.sh, version bump to 0.5.0, CHANGELOG, tag

## Risks

- **Other consumers** — if anyone has cloned this repo and uses `current/`, their workspace breaks. Mitigation: doctor migration is mechanical and reversible; CHANGELOG calls it out as breaking.
- **Naming collision** — "spec" is overloaded (test specs, OpenAPI spec, etc.). Mitigation: doc surfaces always say "system spec" or "capability spec" in prose.
- **Empty `specs/` confusion** — a fresh init has an empty folder. Mitigation: `.gitkeep` + a one-line `SPEC.md` template explaining the empty state.
