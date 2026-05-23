---
version: 1.0
updated: 2026-05-11
summary: "Technology stack and conventions for the Spectacular project"
---

# Stack

## Skill layer
- Claude Code slash command (`/spectacular`)
- SKILL.md lean orchestrator + `references/` subdocs
- Installed at `~/.claude/skills/spectacular/` (global) or `.claude/skills/spectacular/` (project-local)

## CLI layer
- Bootstrap tool (`spectacular init`)
- TBD — stack not yet decided

## Files
- Markdown for all human/agent-readable docs
- YAML frontmatter for machine-readable metadata
- `config.yaml` for project config

## Rules
- Skill reads frontmatter, not full file content, during briefings
- Small files over monolithic documents
- Never overwrite canonical documents — snapshot first
- State lives in frontmatter (PLAN.md for requests, specs/*/SPEC.md for per-capability state; top-level SPEC.md is the index)
