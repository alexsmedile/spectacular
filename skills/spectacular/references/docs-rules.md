---
doc-id: docs-page
mode: stub
location: docs/<section>/<slug>.md
scope: per-request
template: templates/docs/page.md.tmpl
snapshot-on-edit: false
summary: "Single user-facing docs page (deprecated v1.2.0 — see pageworks)"
status: deprecated
covers-additional-doc-ids: [docs-manifest]
---

# Docs Rules — skill rules for `spectacular docs <verb>`

> **⚠ DEPRECATED in spectacular v1.2.0** — public-facing docs authoring is now owned by the [pageworks](https://github.com/alexsmedile/pageworks) skill. This reference will be removed in spectacular v2.0.0. The equivalent (pageworks-native) lives at `pageworks/skills/pageworks/references/authoring.md`.

Consumed by `grill.md` / `refine.md` / `review.md` when the active doc-id is `docs-page` or `docs-manifest`. Schema and folder contract live in [[docs-contract]].

> **Note on shared rules file:** this file declares dispatch for two doc-ids (`docs-page` as primary, `docs-manifest` via `covers-additional-doc-ids`). Special-case for the deprecated docs surface. New docs should follow the one-doc-one-rules-file pattern.

## Scope

This file drives **skill-side** verbs:

- `spectacular docs new <page>` — scaffold a page with frontmatter stub, prompt for section, update `docs.yaml`
- `spectacular docs new --section <name>` — declare a new section in `docs.yaml`, scaffold dir + placeholder page
- `spectacular docs review` — quality gate
- `spectacular docs status` — briefing scoped to docs/

The CLI handles `docs init` (mechanical scaffold) and `doctor docs` (validation) — see `docs-contract.md` § Validation rules and the CLI source.

## `docs new <page>` flow

1. **Resolve the page slug.** Convert `<page>` to kebab-case. Strip `.md` if user typed it. Reject if it collides with an existing page in any section.
2. **Pick the section.**
   - If `--section <name>` flag is set: use it. If the section doesn't exist in `docs.yaml`, error and suggest `docs new --section <name>` first.
   - If no flag: read `docs.yaml`, present the list of section ids + "create new section" as the last option. Ask the user to pick one. Never silently default.
3. **Ask for title and description.** Required frontmatter, can't be elided. One sentence each.
4. **Scaffold the file** at `docs/<section>/<page>.md` using `templates/docs/page.md.tmpl`. Fill `title`, `description`, `section`, `status: draft`, `updated: <today>`. Leave `order` and `since` empty (defaults handle them).
5. **Update `docs.yaml`** — append the page slug to that section's `pages:` list.
6. **Confirm the diff** with the user before any write. Show: new file path + the docs.yaml line being changed.

## `docs new --section <name>` flow

1. **Resolve section id.** kebab-case. Reject collision with existing sections.
2. **Ask for title** (display name). Required.
3. **Ask for order** — default: append (max existing order + 1). User can override.
4. **Scaffold the directory** `docs/<id>/` with a `.gitkeep`.
5. **Append section to `docs.yaml`** with `pages: []`.
6. **Confirm before write.** Show the docs.yaml delta.

If the user is creating a section *and* a page in the same flow, do both in order — section first, then page.

## `docs review` quality gate

Run all checks; report findings as a punch list. Pass = zero errors. Warnings don't block but should be addressed.

### Gate checks

| Severity | Check | Recovery |
|---|---|---|
| error | `docs.yaml` missing or unparseable | Run `spectacular docs init` (won't overwrite existing pages) |
| error | Any page declared in `docs.yaml` but missing on disk | Either create the page or remove the entry |
| error | Any page missing required frontmatter (`title`, `description`, `section`, `status`, `updated`) | Add the field; run `doctor docs --fix` for stubs |
| error | Page `section:` value doesn't match any section in `docs.yaml` | Fix typo or add the section |
| warning | Page file present but not declared in `docs.yaml` (orphan) | Add to docs.yaml or delete |
| warning | Page `updated:` is more than 14 days older than file mtime | Bump `updated:` or sync the content |
| warning | Section folder exists but has zero pages declared and zero files | Remove the section or add a page |
| info | Section has empty `pages:` but folder is empty too | Intentional empty section — fine |

### Output format

```text
docs review — found 2 errors, 1 warning

ERRORS
  ❌ docs/guides/team-billing.md — missing required frontmatter: description
  ❌ docs.yaml — page 'install' declared but docs/getting-started/install.md not found

WARNINGS
  ⚠️  docs/reference/cli.md — orphan (not in docs.yaml)

Suggested fixes:
  • Add `description:` to team-billing.md
  • Create install.md or remove from docs.yaml
  • Add cli.md to docs.yaml reference section, or delete it
```

## `docs status` briefing

Same shape as `/spectacular` no-arg briefing but scoped to `docs/`. Report:

- Site name + tagline (from `docs.yaml`)
- Page count by section: `Getting Started (3), Guides (1), Reference (2)`
- Draft pages: list slugs (so user sees what's unfinished)
- Stale pages: pages where `updated:` is >30 days behind file mtime
- One-line "next action" if anything obvious is open (e.g., "1 page missing frontmatter — run `spectacular docs review`")

Keep it short. Max 8 lines of body.

## Vibe → spec patterns (for `docs refine`, when needed)

`docs refine` is not in v1 scope, but if implemented later, here are the patterns:

| Vibe pattern | Rewrite to |
|---|---|
| "this page is about X" (first sentence) | "X is …" (lead with the thing, not the meta-description) |
| Long intro paragraph | One-sentence lead + immediate "## How to" subhead |
| "we" / "our" | "the CLI" / "Spectacular" (avoid first person plural in public docs) |
| TODO/FIXME inline | Strip and surface as a review-time warning |

## Anti-patterns

- **Don't auto-create pages without confirmation.** Always show the diff first.
- **Don't put per-page audience in frontmatter.** Folder is the audience boundary.
- **Don't deep-nest section folders.** Express nested grouping in `docs.yaml` `pages:` if absolutely needed; keep the filesystem flat.
- **Don't write to `docs/` from non-docs verbs.** Convention pack scaffold, doctor `--fix` for non-docs areas, archive flow — none of these should touch `docs/` files.
- **Don't conflate spec and doc content.** If a question is "how do I use this?" → docs. If it's "what's the contract?" → spec.
