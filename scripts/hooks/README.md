# scripts/hooks/ — Git hooks

Local git hooks managed by the [git-guard](https://github.com/alexsmedile/git-stack) skill.

- `pre-commit` — version consistency check across plugin manifests, marketplace.json, and README badge. Fails the commit if any project-level version is drifting from the others.

These are *not* plugin hooks. For Claude Code / Codex plugin runtime hooks loaded when spectacular is installed as a plugin, see [`/hooks/`](../../hooks/).

> **Same word, different runtimes.** Git hooks fire during `git commit` on this repo. Plugin hooks fire inside an installed Claude Code / Codex session. Don't conflate.

## Install

To enable these hooks locally:

```bash
git config core.hooksPath scripts/hooks
```

Or, if git-guard is installed and configured, it will wire this up automatically.

## Configuration

Edit `scripts/hooks/.git-guard.json` (if present) to configure git-guard's behavior. See the git-guard documentation for the full schema.
