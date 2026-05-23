---
status: verified
priority: medium
owner: alex
updated: 2026-05-23
shipped_in: v0.6.0
summary: "Public-facing docs/ as a first-class Spectacular surface — flat tree, single docs.yaml manifest, page frontmatter, doc verbs (init/new/review), doctor docs area. Renderer-agnostic."
related:
  - ../spec-rename/PLAN.md
  - ../public-docs-advanced/PLAN.md
---

# Plan — Public Docs Foundation

## Goal

Make `docs/` a first-class Spectacular surface — the same way `.spectacular/` is the workspace and `specs/` (post-rename) is system truth. Provide a **flat, opinionated, renderer-agnostic** structure for user-facing documentation: a single `docs.yaml` nav manifest, page-level frontmatter, and doc verbs in the skill (`spectacular docs init/new/review`).

Scope is foundation only: structure + verbs + validation. Renderer-specific adapters and spec→doc sync are deferred to `public-docs-advanced`.

## Why

Three audiences need three surfaces, but today Spectacular only formally owns two:

| Surface | For | Today |
|---|---|---|
| Strategy (PRD/PRINCIPLES/etc.) | Builders + agents planning | ✓ owned |
| Specs (`.spectacular/specs/`) | Builders + agents implementing | ✓ owned (post-rename) |
| **Docs (`docs/`)** | **End users + agents consuming the product** | **unowned — ad-hoc** |

`docs/` exists in this repo today as 5 flat markdown files with no nav, no frontmatter contract, no validation, no scaffold rules. As Spectacular grows, that breaks down: new pages get dropped wherever; the index lives in README; there's no way to validate the docs tree from CI; no way for the skill to know where a new page belongs.

The spec/doc distinction also matters: today `docs/configuration.md` mixes user-facing setup ("how to point at CLAUDE.md") with internal schema ("frontmatter fields"). Without a clear "doc" definition, that mixing continues. Making `docs/` first-class forces the question per-page: *who is this for?*

## Scope

**In scope**

- Convention: `docs/` lives at repo root (not `.spectacular/docs/`) — it's a public artifact
- Single manifest: `docs/docs.yaml` (nav + site metadata) — no per-folder `_section.yaml` files
- Page frontmatter contract: `title`, `description`, `section`, `order` (optional), `status`, `since` (optional), `updated`. **No `audience` field** — audience is folder-level (`docs/` = users + agents consuming the product; `specs/` = devs + coding agents building it).
- Flat preferred file tree:
  ```
  docs/
  ├── docs.yaml
  ├── index.md
  ├── getting-started/
  │   ├── install.md
  │   ├── quickstart.md
  │   └── concepts.md
  ├── guides/
  │   └── ...
  └── reference/
      └── ...
  ```
- No required subfolder beyond what `docs.yaml` declares; sections grow as needed
- **CLI surface** (`cli/spectacular`):
  - `spectacular docs init [--minimal]` — scaffold `docs/docs.yaml` + index.md + 3 default sections (or `--minimal` for just docs.yaml + index)
  - `spectacular doctor docs` — substrate validation (same checks as skill review verb)
- **Skill verbs** (driven by registry + `docs-overrides.md`):
  - `spectacular docs new <page>` — create a page; if `--section` omitted, prompt with section list from docs.yaml + "create new"; updates nav
  - `spectacular docs new --section <name>` — declare a new section in `docs.yaml`, scaffold dir + placeholder page
  - `spectacular docs review` — quality gate (frontmatter complete, nav ↔ fs consistent, no orphans, status field present)
  - `spectacular docs status` — briefing scoped to docs/ (page count by section, draft pages, stale `updated` dates)
- Doc-registry entries for `docs-page` and `docs-manifest` doc types
- Dogfood: migrate this repo's existing `docs/` (5 files) into the new shape — give each a frontmatter stub, populate `docs.yaml`, run review

**Out of scope (deferred to `public-docs-advanced`)**

- Renderer adapters (Mintlify / Docusaurus / Fumadocs / MkDocs export)
- `spectacular docs publish <version>` — snapshot to `docs/versioned/v<x.y.z>/`
- `spectacular docs sync-from-spec <spec>` — draft user-facing page from a spec
- Convention-pack `docs-layout` rule category
- Cross-page reference graph / broken-link detection beyond same-folder relative links
- Search index generation
- i18n / multi-locale

## Decisions

- **`docs.yaml`, no `_section.yaml`** — single manifest is simpler and easier to validate. Per-folder section files multiply maintenance touchpoints. User preference confirmed.
- **Renderer-agnostic in v1** — `docs.yaml` schema + page frontmatter are portable across Mintlify/Docusaurus/Fumadocs/MkDocs. Pick a renderer downstream; Spectacular does not render.
- **Flat tree, sections as folders** — section = folder = group in nav. Pages live one level deep. Subsubsections only via explicit nav nesting in `docs.yaml`. Most projects never need more.
- **`docs/` at repo root, not `.spectacular/docs/`** — public artifact; downstream renderers expect repo-root paths; matches Mintlify/Docusaurus/Fumadocs/Docus conventions.
- **No `audience` field** — the folder is the audience boundary. `docs/` = users + agents consuming the product (audience is universal there); `specs/` = devs + coding agents building it. Per-page audience would be ceremony. Saves one required field.
- **Three default sections on init** — `getting-started`, `guides`, `reference`. Industry-standard triad; users delete or extend.
- **CLI handles `init` + doctor; skill handles `new` + `review` + `status`** — mirrors how `pack` works. CLI scaffolds mechanically; skill drives interactive/judgment work.
- **`docs new` without `--section` prompts** — never silently places pages. Section list comes from docs.yaml; "create new section" is always an option.

## Schemas

### `docs/docs.yaml`

```yaml
site:
  name: Spectacular
  tagline: AI-native operational workspace
  base_url: https://spectacular.dev   # optional, for renderer adapters later

sections:
  - id: getting-started
    title: Getting Started
    order: 1
    pages: [install, quickstart, concepts]
  - id: guides
    title: Guides
    order: 2
    pages: [your-first-request, convention-packs, workflow]
  - id: reference
    title: Reference
    order: 3
    pages: [cli, configuration, frontmatter-schema, scaffold]

# optional top-level entries (no section):
extras:
  - changelog                          # symlinks to ../CHANGELOG.md
  - troubleshooting
```

### Page frontmatter

```yaml
---
title: Install
description: Get Spectacular running in two minutes.
section: getting-started
order: 1
status: stable                 # stable | draft | deprecated
since: 0.1.0
updated: 2026-05-23
---
```

| Field | Required | Notes |
|---|---|---|
| `title` | yes | Display title; defaults to first H1 if absent (warning) |
| `description` | yes | One sentence; used in nav previews and SEO |
| `section` | yes | Must match a section id in docs.yaml |
| `order` | no | Defaults to position in docs.yaml `pages` array |
| `status` | yes | `stable` / `draft` / `deprecated` |
| `since` | no | First version this page applied to |
| `updated` | yes | ISO date; doctor flags if older than file mtime |

## Lifecycle impact

- Adds `docs` area to doctor (warnings on drift, errors on broken manifest)
- Adds two doc types to registry: `docs-page` and `docs-manifest` (uses generic engine)
- Skill SKILL.md routing table extended with `docs` verbs
- Snapshot rule applies to `docs.yaml` (canonical) — page edits do not require snapshots

## Validation

- `spectacular docs init` in tmp dir → docs.yaml + index.md + 3 section dirs each with placeholder page, frontmatter complete
- `spectacular docs new install --section getting-started` → file scaffolded, docs.yaml updated, frontmatter stub
- `spectacular docs new --section api` → docs.yaml gets new section, dir scaffolded, no orphan
- `spectacular docs review` on dogfooded repo → passes; flags any page missing audience/status
- Doctor `docs` area: detects nav ↔ fs drift (page in docs.yaml missing from fs, or vice versa)
- Doctor `docs --fix`: mechanical fixes only — frontmatter stub injection for missing required fields, dedupe entries in docs.yaml
- This repo's `docs/` migration: each of the 5 files gets frontmatter + a section assignment, docs.yaml authored, review passes
- Test suite: `tests/cli/docs.test.sh` covering init, new, new --section, review pass, review fail with broken nav, doctor migration of legacy flat docs/

## Milestones

1. **M1 — Schema** — docs.yaml schema, page frontmatter contract, doc-registry entries, reference doc `docs-contract.md`
2. **M2 — CLI surface** — `spectacular docs init` (CLI subcommand, not just skill verb — like `pack`)
3. **M3 — Skill verbs** — `docs new`, `docs new --section`, `docs review` via engine + `docs-overrides.md`
4. **M4 — Doctor integration** — `docs` area, drift detection, mechanical fixes
5. **M5 — Dogfood** — migrate this repo's 5 docs files, author docs.yaml, run review until clean
6. **M6 — Tests + release** — tests/cli/docs.test.sh, version bump (depends on spec-rename ship order), CHANGELOG, tag

## Risks

- **Manifest staleness** — docs.yaml drifts from filesystem. Mitigation: doctor `docs` area is the primary defense; review verb in CI.
- **Frontmatter friction** — required fields slow page creation. Mitigation: `docs new` scaffolds the stub; only `title` and `description` need real values, the rest have sensible defaults.
- **Renderer lock-in fears** — users worry adopting our schema commits them to a renderer. Mitigation: docs.yaml is small + readable + portable; v2 ships adapter examples for top 4 renderers showing how cheap migration is.
- **Spec ↔ doc duplication** — content drifts between SPEC.md and docs/reference/. Mitigation: deferred to advanced request via `docs sync-from-spec`; for v1, audience metadata makes the boundary explicit, prose stays human-curated.
- **Folder depth creep** — users want sub-sub-sections. Mitigation: docs.yaml supports nested entries; filesystem stays flat. Re-evaluate if 3+ users complain.

## Dependencies

- **Soft dep on `spec-rename`** — clean spec/doc terminology makes audience field semantics clearer. Can ship in either order; prefer spec-rename first.
- No CLI dep — docs verbs are independent of pack and doctor surfaces.
