# scripts/hooks/ — Git hooks

Local git hooks for this repo. Two guards fire on `git commit`:

- `pre-commit` — version consistency across plugin manifests, marketplace.json, and README badge. **Managed by [git-guard](https://github.com/alexsmedile/git-stack)** — do not hand-edit; it is regenerated.
- `pre-commit-wrapper` — runs project-local guards, then delegates to the git-guard `pre-commit`. Currently guards **SKILL.md `description` length** against Codex's 1024-char cap (uses [`../check-skill-desc.sh`](../check-skill-desc.sh); see [`.spectacular/decisions/D7.md`](../../.spectacular/decisions/D7.md)). This file is **not** git-guard-managed — extend it for any future commit-time check.

These are *not* plugin hooks. For Claude Code / Codex plugin runtime hooks loaded when spectacular is installed as a plugin, see [`/hooks/`](../../hooks/).

> **Same word, different runtimes.** Git hooks fire during `git commit` on this repo. Plugin hooks fire inside an installed Claude Code / Codex session. Don't conflate.

## Install

Point git at the `.active/` directory, which symlinks `pre-commit` → the wrapper (so both guards run, and git-guard's file stays untouched):

```bash
git config core.hooksPath scripts/hooks/.active
```

`.active/pre-commit` is a relative symlink to `../pre-commit-wrapper`, which in turn `exec`s `../pre-commit` (the git-guard version check). Editing the wrapper or re-running git-guard never breaks the chain.

> Legacy: `git config core.hooksPath scripts/hooks` runs the git-guard `pre-commit` **only** — skipping the SKILL.md guard. Use `.active/` to get both.

## Configuration

Edit `scripts/hooks/.git-guard.json` (if present) to configure git-guard's behavior. See the git-guard documentation for the full schema.
