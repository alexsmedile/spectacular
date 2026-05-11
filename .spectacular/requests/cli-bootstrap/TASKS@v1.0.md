---
updated: 2026-05-11
---

# Tasks — cli-bootstrap

## Decision

- [ ] Choose CLI stack (Node.js/npm vs shell script vs other)
- [ ] Define distribution method (npm package, homebrew tap, curl-install script, etc.)

## Core init command

- [ ] Scaffold `current/` and `requests/` directories
- [ ] Interactive prompt: project name + summary
- [ ] Write `config.yaml` from template
- [ ] Create stub root files with frontmatter: `PRD.md`, `STACK.md`, `DECISIONS.md`, `AGENTS.md`
- [ ] Update (or create) `.gitignore` with `.spectacular.local/`

## Skill installation

- [ ] Detect if `~/.claude/skills/spectacular/` exists (global install)
- [ ] If yes: create `.claude/skills/` and symlink into it
- [ ] If no: copy skill files into `.claude/skills/spectacular/` as fallback
- [ ] Report what was installed

## Idempotency

- [ ] Skip and report existing files (never overwrite)
- [ ] Skip existing `config.yaml`
- [ ] Safe to re-run as a "repair" for partial scaffolds

## CLI polish

- [ ] Clear success output listing what was created
- [ ] Exit with non-zero code on failure
- [ ] `--help` flag

## Testing

- [ ] Test on blank directory
- [ ] Test idempotency (re-run on existing scaffold)
- [ ] Test when global skill exists vs does not exist
- [ ] Test `.gitignore` creation and append behavior
