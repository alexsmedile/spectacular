---
status: verified
priority: medium
owner: alex
updated: 2026-05-23
target_version: v1.1.0
archived_in: v1.1.0
summary: "Renderer adapters (MkDocs Material + Docusaurus) for the docs/ surface — shipped at v1.1.0 (M4 dogfood dropped, deferred to docs-writer chain)"
supersedes: PLAN@v1.0.md
related:
  - ../../archive/public-docs-foundation/PLAN.md
  - ../../specs/cli/SPEC.md
  - ../../../skills/spectacular/references/docs-contract.md
---

# Plan — Public Docs Advanced (narrowed)

> Earlier scope (4 renderers, 6 milestones, gated on external signals) preserved as `PLAN@v1.0.md`. This revision narrows to a dogfood activation against spectacular's own docs/.

## Goal

Ship the **MkDocs (Material)** and **Docusaurus** renderer adapters for the existing docs/ surface (shipped in v0.6.0 via `public-docs-foundation`), and dogfood by publishing spectacular's own docs/ to GitHub Pages.

## Why now

Spectacular's docs/ is the canonical example of the schema. Without a working render path, the schema is theoretical — every downstream user has to invent the adapter themselves. Shipping two structurally-different OSS adapters (YAML-driven MkDocs, JS-driven Docusaurus) also validates that the `docs.yaml` schema is genuinely renderer-agnostic.

**Activation override:** the PLAN@v1.0 gate required 2-of-3 external signals (Mintlify asks, 20+ page drift, community pack PR). None have fired. Activating anyway because the author wants this for the spectacular repo itself — internal dogfooding overrides the external-signal gate.

## Scope

**In scope**

- **M1 — Adapter reference doc** — `skills/spectacular/references/docs-renderer-adapters.md` documenting:
  - `docs.yaml` → `mkdocs.yml` mapping (nav, theme, plugins, hooks)
  - `docs.yaml` → `sidebars.js` + `docusaurus.config.js` mapping
  - Page frontmatter field translation table per renderer
  - Build-step boilerplate per renderer (commands, deps)
- **M2 — `spectacular docs export <renderer>` CLI** — writes config alongside docs/:
  - `mkdocs` → writes `docs/mkdocs.yml` at repo root (Material theme defaults)
  - `docusaurus` → writes `docs/docusaurus.config.js` + `docs/sidebars.js`
  - Both → write `.github/workflows/docs.yml` deploy action for GitHub Pages (idempotent: skip if exists, `--force` to overwrite)
  - `--out <path>` flag to override default location
- **M3 — `docs.yaml` schema extension** — add per-export-target hints:
  - `renderers: { mkdocs: {theme, plugins}, docusaurus: {preset, themeConfig} }`
  - Extension is additive; existing docs.yaml files remain valid
- **M5 — Tests + release** — smoke tests for both adapters, version bump to v1.1.0

**M4 — Dogfood: DROPPED** (architectural decision, 2026-05-23). Authoring spectacular's own docs/ is a job for a dedicated `docs-writer` agent backed by a new `docuworks` skill, not for the spectacular orchestrator. Adapters can be validated structurally without authoring real pages. Tracked in follow-on requests: `spectacular-bundle-restructure` → `docuworks-skill` + `docs-writer-agent` → `public-docs-dogfood`.

**Out of scope (deferred — own future requests if signal emerges)**

- Mintlify adapter (paid renderer)
- Fumadocs adapter (Next.js lock-in)
- Versioned docs (`docs publish` + `docs/versioned/`) — was M3 in PLAN@v1.0
- Spec → doc sync (`docs sync-from-spec`) — was M4 in PLAN@v1.0
- Convention-pack `docs-layout` rule category — was M5 in PLAN@v1.0; structurally blocked on `convention-pack-modules` v2
- Renderer round-trip (import existing config into docs.yaml)
- Hosted preview server (`docs serve`) — defer to user's `mkdocs serve` / `docusaurus start`
- Search index, i18n, MDX components, code-AST API reference

## Decisions

- **Two renderers, not four.** MkDocs Material covers the simplest case (single YAML); Docusaurus covers React/JS. Mintlify out (paid), Fumadocs out (Next.js coupling). Both deferred renderers are documented as "community-contributable" in M1's reference doc — explicit invitation, no built-in support promise.
- **Config alongside docs/, not in a sibling.** `docs/mkdocs.yml` and `docs/docusaurus.config.js` live with the content. Avoids a `docs/_export/` directory; renderer reads docs/ in place. User runs `mkdocs serve` / `docusaurus start` directly from the repo root.
- **`docs.yaml` is the source of truth.** No inference, no folder-walking heuristics. Adapter is a pure transformation. The `renderers:` block extension keeps cross-cutting renderer hints in one declared location.
- **GitHub Pages is the hosting target.** Adapter ships a working `.github/workflows/docs.yml`. Other targets (Vercel, Cloudflare, Netlify) documented in M1 reference but not pre-built — copy/paste from MkDocs/Docusaurus official docs.
- **Adapters export, they don't render.** Spectacular writes configs; the renderer's own toolchain builds the site. No build-pipeline ownership.
- **`--force` required to overwrite existing configs.** Bare `docs export` is safe by default. Idempotent re-run reports skipped files.

## Milestones

1. **M1 — Reference doc** ✓ shipped (commit pending)
   - `skills/spectacular/references/docs-renderer-adapters.md`
2. **M2 — `docs export` CLI** ✓ shipped (commit pending)
   - `spectacular docs export mkdocs|docusaurus [--out <path>] [--force] [--no-workflow]`
   - Idempotency, magic-comment pin, GitHub Pages workflow
3. **M3 — docs.yaml schema extension** ← in progress
   - Document `renderers:` block in `docs-contract.md` (referenced from adapter doc; needs full schema section)
   - Update `templates/docs/` template (if present) with commented `renderers:` example
   - Doctor `docs` area: validate `renderers:` if present (non-required; warn on unknown renderer keys)
4. ~~M4 — Dogfood~~ (dropped — see scope note above)
5. **M5 — Tests + release**
   - `tests/cli/docs-export.test.sh` — both adapters produce structurally-valid configs (YAML parse for mkdocs, JS syntax check for docusaurus)
   - Idempotency + --force + pin tests
   - Version bumps: cli/spectacular SPECTACULAR_VERSION, both plugin.json files, SKILL.md frontmatter
   - CHANGELOG entry under `[1.1.0]`
   - Tag v1.1.0, GitHub release, plugin marketplace bump

## Validation

- `spectacular docs export mkdocs` writes `mkdocs.yml` + `.github/workflows/docs.yml` from a docs/ tree (validated in sandbox 2026-05-23)
- `spectacular docs export docusaurus` writes `docusaurus.config.js` + `sidebars.js` (validated in sandbox 2026-05-23)
- Re-running reports skipped files; `--force` overwrites; `// spectacular: do-not-overwrite` pin respected even with `--force`
- Unknown renderer (mintlify/fumadocs) returns actionable error pointing at adapter reference
- Missing `docs.yaml` fails clean
- Empty sections dropped from nav (avoids invalid mkdocs YAML)
- Doctor `docs` area passes against extended docs.yaml

**Authoring validation deferred** — confirming the adapters produce render-quality output requires real pages, which is dogfood/docs-writer territory (separate request).

## Risks

- **Theme / plugin churn.** MkDocs Material releases monthly; Docusaurus v3 → v4 transitions happen. Mitigation: pin versions in the deploy workflow; document upgrade path in M1 reference. Adapter writes a tested-as-of-date version comment.
- **Docusaurus complexity creep.** `sidebars.js` is JS; users hit edge cases (auto-generated sidebars, category metadata). Mitigation: ship a minimal generator that covers docs.yaml's declared nav; document escape hatch (hand-edit sidebars.js, mark with `// spectacular: do-not-overwrite`).
- **`--force` overwrites user customization.** Once a user hand-edits the generated config, re-running with `--force` blows it away. Mitigation: clear warning in M1 reference; consider `// spectacular: do-not-overwrite` magic comment that even `--force` respects (defer to v1.2 if it becomes a real complaint).
- **GitHub Pages permissions friction.** `.github/workflows/docs.yml` needs `contents: write` + Pages enabled in repo settings. Mitigation: post-install message tells user the two clicks needed.

## Success criteria

- Both adapters produce structurally-valid configs from any docs.yaml that passes doctor
- Reference doc shipped at `skills/spectacular/references/docs-renderer-adapters.md`
- Schema extension (`renderers:` block) documented in `docs-contract.md`
- Doctor `docs` area validates `renderers:` block when present
- Smoke tests in `tests/cli/docs-export.test.sh` cover both adapters + idempotency + pin + --force
- v1.1.0 tagged + released

**Out of success criteria** — having a live hosted spectacular docs site (depends on docs-writer + dogfood request).
