# Pageworks Handoff — when spectacular delegates public-doc work

Loaded when: spectacular detects that an action may require public-doc updates and needs to surface the handoff (typically during `archive`, but also during `status` briefings or when the user asks about docs).

This file is the contract for **how** spectacular hands off to pageworks. Spectacular never invokes pageworks itself; it surfaces the signal and lets the user act.

## The boundary, in one line

**Spectacular owns `.spectacular/` (internal workspace). Pageworks owns `docs/` (public-facing).** Both can coexist; both are optional.

## When to suggest pageworks

Three triggers ranked by strength of signal:

| Trigger | Confidence | Signal text |
|---|---|---|
| `archive <slug>` where the request touched SPEC.md, specs/, or `docs/` references | high | "This request changed SPEC; public docs may need updating." |
| `status` briefing when `docs/` exists but pageworks is not installed | medium | "docs/ folder detected — install pageworks for public-doc work." |
| Any time user explicitly asks about docs/, writing a page, or rendering | high | "Pageworks owns docs/ end-to-end — install if you'd like the dedicated tooling." |

**Never** auto-invoke pageworks. The prompt is a signal; the action is the user's.

## The install hint (canonical phrasing)

When pageworks is not installed and the situation calls for it, use this exact phrasing so users learn to recognize the pattern:

```
Public docs work has moved to a dedicated skill: pageworks.
Install: https://github.com/alexsmedile/pageworks
```

If the user has expressed preference for a specific install tool earlier in the conversation (`apm`, `skizl`, manual symlink), append the appropriate command. Otherwise stop at the URL — installation choice is theirs.

## Archive-time prompt

When `spectacular archive <slug>` runs, spectacular inspects the archived request's tracked files. If any of these changed (per `git diff` against the previous archive snapshot, or by file mtime if no git):

- `.spectacular/SPEC.md`
- `.spectacular/specs/**`
- `.spectacular/ARCHITECTURE.md` (when present)
- `.spectacular/PRD.md` (when present)

…spectacular prints, **after** the archive is complete:

```
ℹ This request changed internal specs. Public docs/ may need updates.

  Suggested next steps:
    pageworks audit           # find pages whose synced_from: matches the changed spec
    pageworks review <page>   # work through each flagged page

  Or skip future prompts:
    spectacular archive <slug> --no-docs-prompt
    or set docs.prompt_on_archive: false in .spectacular/config.yaml
```

### Suppression mechanics

The prompt can be suppressed two ways:

- **Per-invocation**: `spectacular archive <slug> --no-docs-prompt`
- **Per-project**: set in `.spectacular/config.yaml`:
  ```yaml
  docs:
    prompt_on_archive: false   # default: true
  ```

When both are absent and `docs/` doesn't exist in the project, the prompt is silent regardless — there's nothing to update.

When `docs/` exists but pageworks isn't installed, the prompt includes the install hint above.

## Status-briefing reference

When `spectacular status` or `/spectacular` (no args) runs and the project has `docs/`:

- Mention "docs/ present" in the briefing.
- If pageworks installed: suggest running `pageworks status` for the docs briefing.
- If not installed: surface the install hint as a one-line follow-up.

Do not duplicate pageworks's briefing surface — spectacular reports presence, pageworks reports content.

## The reverse direction (pageworks → spectacular)

Pageworks's `maintenance.md` § "Common spec sources" knows that `.spectacular/specs/<x>/SPEC.md` is a typical source for `synced_from:` frontmatter. That's the only cross-reference pageworks makes back to spectacular; it never invokes spectacular's CLI.

If pageworks is installed but spectacular isn't, that's fine — pageworks's `synced_from:` works against any file path, not just `.spectacular/`.

## What this file is NOT

- Not a guide to writing docs (that's pageworks's `authoring.md` / `prose-patterns.md`).
- Not the full pageworks install/usage doc (that's pageworks's README).
- Not a runtime contract that spectacular and pageworks negotiate version compatibility — there's no API surface between them. Each runs on its own.

## Anti-patterns

- **Auto-installing pageworks** when the prompt fires — violates the global "never install without explicit instruction" rule. Surface, don't act.
- **Reading from `docs/`** during spectacular's verbs — spectacular's awareness is folder/manifest presence only. Anything past that is pageworks's job.
- **Recreating pageworks's behaviors** because the handoff "feels heavy" — if spectacular is doing schema validation, frontmatter checks, or renderer work on docs/, the boundary has been violated. Fix the call site.
- **Adding the install hint to non-docs flows** — spectacular shouldn't mention pageworks during workspace migrations, doctor checks unrelated to docs, pack management, etc. Hint only fires when docs/ is involved.
