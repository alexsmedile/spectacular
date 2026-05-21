---
version: 1.0
updated: 2026-05-11
summary: "Agent orchestration rules for the Spectacular project"
---

# Agents

## Context loading rules

| Task type | Load |
|---|---|
| Planning / design | PRD.md, DECISIONS.md, SPECTACULAR_PRD.md |
| Skill implementation | STACK.md, PLAN.md, TASKS.md, current/skill/ |
| CLI implementation | STACK.md, PLAN.md, TASKS.md, current/cli/ |
| Review / QA | VERIFY.md, current/<capability>, RISKS.md |

## Available skills

- `spectacular` — workspace management (self-referential: this project uses the skill it's building)

## Handoff conventions

- Summarize state in SESSION.md before handing off
- Load only the capability spec relevant to the current task, not all of current/
- The authoritative PRD for the product is `SPECTACULAR_PRD.md` at the repo root
