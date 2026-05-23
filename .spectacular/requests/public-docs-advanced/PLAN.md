---
status: planned
priority: low
owner: alex
updated: 2026-05-23
summary: "v2 docs surface — renderer adapters (Mintlify/Docusaurus/Fumadocs/MkDocs), docs versioning, spec→doc sync, convention-pack docs-layout category"
related:
  - ../public-docs-foundation/PLAN.md
  - ../spec-rename/PLAN.md
  - ../convention-pack-modules/PLAN.md
---

# Plan — Public Docs Advanced

## Goal

Build on `public-docs-foundation/` with three capabilities that need real-world usage before they're worth designing:

1. **Renderer adapters** — export `docs/` to Mintlify / Docusaurus / Fumadocs / MkDocs configurations
2. **Versioned docs** — snapshot docs/ at release time, navigable historical versions
3. **Spec ↔ doc sync** — draft user-facing pages from `.spectacular/specs/<x>/SPEC.md`; flag drift when specs change

Plus convention-pack integration: a `docs-layout` rule category so packs can declare required sections, renderer hints, frontmatter requirements.

## Why

Foundation ships a renderer-agnostic structure. That's the right v1, but it leaves three gaps that only appear with use:

- **Adapter friction** — users will adopt the schema then ask "how do I get this into Mintlify?" Without a documented adapter path, they fork or abandon.
- **Version drift** — once docs/ has 50 pages and 3 releases, "which version of the docs is this?" becomes a real question. Docusaurus and Mintlify both solve this with snapshot folders.
- **Spec/doc drift** — even with audience metadata, the same concept gets two pages (spec and doc) that diverge. A sync verb that drafts the user-facing version from the spec keeps them aligned by construction.

This request is **gated on real-world signals** — don't start until at least two of:

- 3+ users ask "how do I render this with Mintlify/Docusaurus"
- Docs/ in a downstream project hits 20+ pages and starts drifting from spec
- A community contribution arrives that would be cleaner as a pack `docs-layout` rule

Same pattern as `convention-pack-modules`: ship v1, learn from production, design v2 from data.

## Scope

**In scope (once activated)**

- Adapter docs in `skills/spectacular/references/docs-renderer-adapters.md`:
  - Mintlify (`mint.json` generation from `docs.yaml`)
  - Docusaurus (`sidebars.js` + `docusaurus.config.js` frontmatter mapping)
  - Fumadocs (`meta.json` per folder)
  - MkDocs (`mkdocs.yml` nav)
- CLI: `spectacular docs export <renderer> [--out path]` — write renderer config from `docs.yaml` + page frontmatter
- Versioning:
  - `spectacular docs publish <version>` — snapshot `docs/` → `docs/versioned/v<x.y.z>/`
  - `docs.yaml` extended with `versions: [...]` list
  - Doctor `docs` area: versioned snapshots are read-only, drift = error
- Spec sync:
  - `spectacular docs sync-from-spec <spec-path>` — read `.spectacular/specs/<x>/SPEC.md`, draft `docs/reference/<x>.md`
  - Sync-metadata in page frontmatter: `synced_from: .spectacular/specs/<x>/SPEC.md` + `synced_at: <date>`
  - Doctor `docs` area: when spec mtime > synced_at, warn "spec changed since last sync"
- Convention-pack `docs-layout` category:
  - Required sections (e.g., `[getting-started, reference]`)
  - Renderer hint (`renderer: mintlify | docusaurus | fumadocs | mkdocs | none`)
  - Required frontmatter fields (extend base contract)
  - `CHANGELOG.md` handling (`symlink | generate | external`)

**Out of scope (further deferred)**

- Hosted preview / server (`spectacular docs serve`)
- Search index generation
- i18n / multi-locale
- Auto-generated API reference from code AST
- MDX / custom component support beyond plain markdown
- Renderer round-trip (import existing Mintlify into our schema)

## Decisions (provisional — revisit at activation)

- **Export, not render** — Spectacular writes renderer configs alongside `docs/`; the renderer itself runs in the user's tooling. Avoids us owning a build pipeline.
- **`docs/versioned/v<x.y.z>/` snapshot model** — same as Docusaurus. Simpler than git-tag-driven reads; lets historical docs ship with the repo.
- **Spec→doc sync is human-curated** — skill drafts, user approves. No auto-overwrite. Sync metadata in frontmatter flags drift; resolution is always a conversation.
- **Adapter docs ship as reference, not as auto-running converters at first** — write the spec; let users invoke `docs export` explicitly. Auto-export on commit is a later add-on.

## Activation triggers

Don't start work until **2 of 3**:

1. 3+ users (issues, PRs, conversations) ask how to render with Mintlify or Docusaurus
2. Any project using docs/ hits 20+ pages and surfaces drift complaints
3. A community pack or PR ships that would be cleaner as `docs-layout` pack rule

## Validation (provisional)

- `docs export mintlify --out _mintlify/` produces working `mint.json` that Mintlify renders
- Same for docusaurus / fumadocs / mkdocs (smoke test each)
- `docs publish 0.5.0` snapshots cleanly; re-running is idempotent
- `docs sync-from-spec .spectacular/specs/auth/SPEC.md` produces a draft page; running again after spec edit warns about drift
- Pack with `docs-layout` enforces required sections via doctor

## Milestones (provisional)

1. **M1 — Renderer adapter spec doc** — `docs-renderer-adapters.md` with mapping tables for 4 renderers
2. **M2 — `docs export <renderer>`** — CLI + 4 adapter implementations
3. **M3 — Versioning** — `docs publish`, snapshot, docs.yaml extension, doctor handling
4. **M4 — Spec sync** — `docs sync-from-spec`, frontmatter metadata, drift detection
5. **M5 — Pack integration** — `docs-layout` rule category, doctor wiring
6. **M6 — Tests + release**

## Risks

- **Renderer surface churn** — Mintlify/Docusaurus/Fumadocs evolve; adapters need maintenance. Mitigation: ship adapters as reference docs + simple shell scripts users can copy, not as a built-in library we own forever.
- **Versioning storage cost** — snapshot per release inflates repo size. Mitigation: opt-in only; doctor flags but doesn't enforce.
- **Sync drift false positives** — spec edit that doesn't affect user-facing wording still flags drift. Mitigation: `synced_at` is updated by user via `docs sync-ack <path>` when they confirm no doc change needed.
- **Pack composition complexity** — `docs-layout` adds rule category #7. Mitigation: depends on `convention-pack-modules` v2 work to keep schema clean — sequence after that lands.
