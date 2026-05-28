# Init Workflow — CLI Bootstrap

Triggered by: `spectacular init` CLI command, or first-time setup conversation.

---

## What init does

Sets up the `.spectacular/` scaffold on a new project. Run once per project. Safe to re-run — idempotent + non-destructive (never overwrites existing files).

As of v0.3.0, init scaffolds **only what the project needs** rather than all 7 root docs. The default is a 5-file **always-set**; extra docs come from the kit declaration or explicit flags.

---

## The always-set (5 files + 2 dirs)

Created on every init regardless of kit or flags:

| Path | Purpose |
|---|---|
| `.spectacular/PRD.md` | Product intent — the project's anchor doc |
| `.spectacular/config.yaml` | Machine-readable config (project name, naming rules, kit identity) |
| `.spectacular/SPEC.md` | System spec index — what's built, present tense |
| `.spectacular/<agents-file>` | Onboarding doc for agents; defaults to `AGENTS.md`, configurable via `--agents-file` |
| `.spectacular/requests/` | Request folders live here |
| `.spectacular/specs/` | Per-capability specs (optional; only when a capability outgrows a one-liner in SPEC.md) |

**Rationale:** PRD is the anchor every other doc references. SPEC.md is the always-on index of what's built. config.yaml is how the skill discovers project state. AGENTS.md is the cold-start onboarding doc. requests/ and specs/ are the working surfaces. Without these six, the workspace is unusable.

Everything else (`PRINCIPLES.md`, `ARCHITECTURE.md`, `ROADMAP.md`, `STACK.md`, `DECISIONS.md`) is opt-in via kit or `--with` flag. Empty stubs of unused docs create *stub fatigue* — the skill still reads them during briefings, diluting signal.

---

## Kit-driven scaffolding (v0.3.0+)

The user picks a **kit** that declares which extra docs the project needs. Kits are defined in `templates/prd/kits/<id>.md` (see [[kits-contract]]):

| Kit | Always-docs (auto-scaffolded) | Suggested-docs (interactive prompt) |
|---|---|---|
| `blank` | none | none |
| `coding` | STACK, ARCHITECTURE | PRINCIPLES, ROADMAP, DECISIONS |
| `content` | ROADMAP | PRINCIPLES, DECISIONS |
| `product` | ROADMAP | STACK, ARCHITECTURE, DECISIONS, PRINCIPLES |
| `research` | (none) | DECISIONS |

Selection:
- **Non-interactive default** — `blank` kit (no inference, no detection). Use `--kit <name>` to override.
- **Interactive (`-i`)** — menu prompts for kit, then asks y/n per suggested doc.

---

## CLI flags (v0.3.0)

```
spectacular init                              # always-set + blank kit
spectacular init -i                           # interactive: pick kit + per-suggested prompts
spectacular init --kit coding                 # always-set + coding's always-docs (STACK + ARCHITECTURE)
spectacular init --with principles,roadmap   # additive on top of kit defaults
spectacular init --kit coding --minimal       # always-set only; kit identity preserved, extras skipped
spectacular init --name my-app --agents-file CLAUDE.md
spectacular init --skill-scope global         # install skill to ~/.agents/ and ~/.claude/
spectacular init --skill-scope none           # scaffold only; install no skill
spectacular init --no-skill                   # alias for --skill-scope none
spectacular init --update                     # re-download latest skill release
```

**Flag interactions:**
- `--kit` (default `blank`) sets the kit identity and triggers its always-docs.
- `--with <a,b,c>` adds explicit docs on top of kit always-docs (additive, deduplicated).
- `--minimal` overrides kit always-docs — only always-set is scaffolded, regardless of kit. Kit identity is still recorded in PRD frontmatter.
- `--skill-scope <project|global|none>` (v1.8.3+) controls where (or whether) the skill is installed — see below. `--global` is a deprecated alias for `--skill-scope global`; `--no-skill` for `--skill-scope none`.
- Unknown kit, unknown doc-ID, or unknown skill-scope value errors cleanly (non-zero exit, helpful message).

---

## Skill scope + existing-install detection (v1.8.3+)

Init can install the skill in one of three scopes, controlled by `--skill-scope`:

| Scope | Installs to | Use when |
|---|---|---|
| `project` | `./.agents/skills/spectacular` + `./.claude/skills/spectacular` symlink | The skill should travel with this repo (committed `.agents/`) |
| `global` | `~/.agents/skills/spectacular` + `~/.claude/skills/spectacular` | One install shared across all your projects |
| `none` | nothing | The skill is already available some other way (plugin / global / upstream) |

**Before installing, init scans for an existing spectacular** in every location it could already be available:

1. **Current project** — `./.agents/...` or `./.claude/...`
2. **Up the worktree** — any parent directory between cwd and `$HOME`/`/` with `.agents/skills/spectacular` or `.claude/skills/spectacular`
3. **Global user scope** — `~/.agents/skills/spectacular`, `~/.claude/skills/spectacular`
4. **Plugin installs** — Claude Code (`~/.claude/plugins/cache/spectacular`), Codex (`~/.codex/plugins/cache/spectacular`), Gemini (`~/.gemini/extensions/spectacular`)

**Resolution when `--skill-scope` is *not* passed:**
- If any existing install is found → **warn (listing each location) and default to `none`** — scaffold proceeds, no redundant copy is installed. Force a copy with `--skill-scope project` (or `global`).
- If nothing is found → default to `project` (the historical behavior).

An explicit `--skill-scope` always wins over auto-detection — e.g. `--skill-scope project` installs a local copy even when a plugin install exists.

---

## Pre-flight non-overwrite

Init is **always idempotent + non-destructive**. Re-running on an initialized workspace is safe by design.

| State | Behavior | Stdout |
|---|---|---|
| File doesn't exist | Create with stub | `✓ created PRD.md` |
| File exists, empty (0 bytes or whitespace only) | Fill with stub | `✓ filled empty PRD.md` |
| File exists, has content | **Skip** — never overwrite | `⊘ PRD.md already present, leaving alone` |
| File exists, malformed (no frontmatter) | Skip; emit generic diagnostics | `⊘ PRD.md (issues detected — run diagnostics via spectacular doctor once available)` |
| Directory exists | No-op silently | (no output) |
| `.gitignore` entry missing | Append entry only | `✓ added .spectacular.local/ to .gitignore` |
| `.gitignore` entry present | No-op | (no output) |

**Decided exclusions:**
- **No `--force` flag.** Re-init can never overwrite content. To regenerate a stub: delete the file first, then re-init.
- **No schema migration.** Init creates; doctor diagnoses + opt-in repairs.
- **No project-type inference.** Bare init uses `blank` kit unconditionally.

**Adding a kit later is safe:** `spectacular init --kit coding` on an existing workspace adds only the kit's missing always-docs; never touches existing files. This is the canonical "upgrade from blank to coding" path.

---

## Init sequence (high-level)

1. Parse flags + validate (`--kit` known, `--with` doc IDs known)
2. Detect existing installs; if found and no `--skill-scope` given, warn + default scope to `none` (see "Skill scope" below). If `-i`: run interactive prompts (name, summary, agents-file, skill-scope, kit, per-suggested-doc y/n)
3. Resolve doc-set: `always-set + (kit's always-docs unless --minimal) + --with entries`
4. Scaffold directories (`specs/`, `requests/`)
5. Scaffold each resolved doc via `write_if_missing` (pre-flight rules apply per file)
6. Update `.gitignore` (append `.spectacular.local/` if absent)
7. **Pack consultation** — if `.spectacular/config.yaml` declares `convention_pack:` with `mode: scaffold` or `mode: enforce`, append the pack's `gitignore.always-add` entries (deduplicated). Pack source resolved via scope precedence (project-local → user → app-store → bundled). Always-set always wins on conflicts; pack never overwrites existing lines.
8. Install skill into the resolved scope (or skip if `--skill-scope none`, if an existing install was detected, or if already installed per `skills.lock`)
9. Print summary

## Convention packs (v0.4.0+)

Packs are opt-in repo-shape opinions. A pack declares rules across 6 categories (naming / taxonomy / root-files / gitignore / file-placement / project-types) — see [[packs-contract]] for the full schema.

A repo activates a pack via `config.yaml`:

```yaml
convention_pack:
  source: alex-default     # pack name resolved via scope precedence
  mode: suggest            # suggest | scaffold | enforce
```

**Three modes:**

| Mode | Init behavior | Doctor behavior |
|---|---|---|
| `suggest` | Pack is read but not applied. Skill may mention pack opinions during interactive work. | `conventions` check reports the pack is active but skips drift checks. |
| `scaffold` | Init appends pack's `gitignore.always-add` entries to `.gitignore` (idempotent). | `conventions` check flags missing gitignore entries as warnings. |
| `enforce` | Same as scaffold. | `conventions` check flags missing gitignore entries as errors (exit 2). |

**Pack precedence** (when same name exists in multiple scopes — first hit wins):
1. `<project>/.spectacular/packs/<name>/` (project-local)
2. `~/.spectacular/packs/<name>/` (user)
3. `<spectacular-repo>/packs/<name>/` (app-store, source-repo only)
4. `<skill>/templates/packs/<name>/` (bundled — currently just `minimal`)

**Why init does not auto-declare a pack:** on first init the user hasn't picked an opinion yet. `minimal` is the implicit baseline (it IS the bundled README + gitignore stub). To activate a heavier pack, user edits config.yaml after init and re-runs (idempotent) — or wires it via interactive init in a future version.

---

## .spectacular.local/

Personal override layer — never committed. Use for local dev overrides, personal config variations, sensitive local paths. The skill reads `.spectacular.local/` if present and merges with `.spectacular/`, with local taking precedence.

---

## Idempotency guarantees

- Re-running init is always safe.
- No file is ever overwritten.
- `.gitignore` entries are appended, never rewritten.
- Re-running with a different kit only adds missing docs; doesn't change existing content or kit identity in PRD frontmatter (existing PRD wins).

---

## Verification

Smart-init ships with a Bash test harness at `tests/cli/init.test.sh` covering all six core scenarios (bare init, `--kit coding`, `-i` interactive, idempotent re-run, `--with` flag, `--minimal` override). Run via `tests/run.sh cli`. Manual QA checklist lives in `requests/smart-init/VERIFY.md`.

## Related

- [[doc-index]] — registry the kit triggers-docs entries map into
- [[kits-contract]] — kit extension schema (adds-slots, modifies-slots, triggers-docs)
- [[verification]] — when VERIFY.md is justified for a CLI change like smart-init
- [[lifecycle]] — what happens to a workspace after init (request lifecycle, archival)
