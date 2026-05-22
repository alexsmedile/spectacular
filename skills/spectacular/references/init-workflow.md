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
| `.spectacular/<agents-file>` | Onboarding doc for agents; defaults to `AGENTS.md`, configurable via `--agents-file` |
| `.spectacular/requests/` | Request folders live here |
| `.spectacular/current/` | Capability specs (canonical truth) live here |

**Rationale:** PRD is the anchor every other doc references. config.yaml is how the skill discovers project state. AGENTS.md is the cold-start onboarding doc. requests/ and current/ are the working surfaces. Without these five, the workspace is unusable.

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
spectacular init --global                     # install skill to ~/.agents/ and ~/.claude/
spectacular init --update                     # re-download latest skill release
```

**Flag interactions:**
- `--kit` (default `blank`) sets the kit identity and triggers its always-docs.
- `--with <a,b,c>` adds explicit docs on top of kit always-docs (additive, deduplicated).
- `--minimal` overrides kit always-docs — only always-set is scaffolded, regardless of kit. Kit identity is still recorded in PRD frontmatter.
- Unknown kit or unknown doc-ID errors cleanly (non-zero exit, helpful message).

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
2. If `-i`: run interactive prompts (name, summary, agents-file, scope, kit, per-suggested-doc y/n)
3. Resolve doc-set: `always-set + (kit's always-docs unless --minimal) + --with entries`
4. Scaffold directories (`current/`, `requests/`)
5. Scaffold each resolved doc via `write_if_missing` (pre-flight rules apply per file)
6. Update `.gitignore` (append `.spectacular.local/` if absent)
7. Install skill (or skip if already installed per `skills.lock`)
8. Print summary

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

- [[doc-registry]] — registry the kit triggers-docs entries map into
- [[kits-contract]] — kit extension schema (adds-slots, modifies-slots, triggers-docs)
- [[verification]] — when VERIFY.md is justified for a CLI change like smart-init
- [[lifecycle]] — what happens to a workspace after init (request lifecycle, archival)
