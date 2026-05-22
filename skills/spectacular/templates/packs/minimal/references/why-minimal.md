# Why minimal?

Spectacular ships `minimal` as the default pack because **opinions about repo shape should be opt-in, not opt-out**. A project initialized today is rarely a project that knows what it'll become in three months. Forcing a `src/` / `tests/` / `docs/` skeleton on a one-script automation is friction; forcing a kebab-case naming rule on a Python package that's idiomatically snake_case is wrong.

So `minimal` enforces two things only:

1. **A `.gitignore` exists with safe defaults.** Almost every project benefits from ignoring `_archive/`, `_backups/`, `scratch/`, `.env.local`. Adding these is free; not adding them means a single careless `git add .` leaks a `.env.local` or commits a 500MB backup. Zero downside.

2. **A `README.md` exists with a Type/Stack/Run header.** The header is the difference between an agent or human spending 2 seconds vs 2 minutes figuring out what a repo is. It costs nothing to maintain (3 lines) and pays back every time someone clones or revisits the project.

Tool-generated hidden dirs (`.scrapekit/`, `.playwright-mcp/`, `.smart-env/`, `.obsidian/`) are explicitly in `never-auto-add`. The skill asks before touching them — per the global rule that auto-gitignoring tool state without consent is a behavior violation.

## When to install a heavier pack

You probably want more opinions when:

- You maintain a **mono-collection** of projects (e.g. `~/code/` with apps/libs/tools/) — install a pack that defines the top-level taxonomy and detection heuristic.
- You ship many projects of the **same type** (CLIs, libraries, skills, plugins) — install a pack with `project-types.<type>` scaffolds so `spectacular init --type cli` produces a consistent layout every time.
- You want **enforce-mode doctor checks** — install a pack and set `config.yaml`'s `convention_pack.mode: enforce` so doctor flags drift from your conventions during normal use.

Check `<repo>/packs/` (the Spectacular app store) for packs that match your workflow. `alex-default` is the canonical reference if you want to see what a fully-opinionated pack looks like.

## When not to install one

You're working on a one-off script, a quick spike, or a sandbox project. `minimal` is enough. Stronger opinions become tax, not benefit, on small surface area.
