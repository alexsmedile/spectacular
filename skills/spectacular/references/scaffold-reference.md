# Scaffold Reference — File Templates

Canonical templates for all Spectacular files. Used by the skill when creating new files.

---

## Root layer

### PRD.md
```md
---
version: 1.0
updated: <today>
summary: "<one-line product/business intent>"
---

# Product Requirements Document

## Vision
<What this product is and why it exists>

## Goals
- 

## Non-goals
- 

## Target users
- 
```

### STACK.md
```md
---
version: 1.0
updated: <today>
summary: "Technology stack and engineering rules"
---

# Stack

## Frontend
- 

## Backend
- 

## Infrastructure
- 

## Rules
- 
```

### DECISIONS.md
```md
---
version: 1.0
updated: <today>
summary: "Architectural and product decisions log"
---

# Decisions

## <YYYY-MM-DD>

**Decision:** 
**Why:** 
**Tradeoffs:** 
```

### AGENTS.md
```md
---
version: 1.0
updated: <today>
summary: "Agent orchestration rules for this project"
---

# Agents

## Context loading rules

| Task type | Load |
|---|---|
| Planning | PRD.md, DECISIONS.md |
| Implementation | STACK.md, PLAN.md, TASKS.md, current/<capability> |
| Review | VERIFY.md, current/<capability>, RISKS.md |

## Available skills

- spectacular — workspace management

## Handoff conventions

- Summarize state in SESSION.md before handing off to another agent
```

### config.yaml
```yaml
project:
  name: <project-name>
  summary: "<one-liner>"

naming:
  requests: kebab-case
  prefix: ""

required_files:
  requests:
    - PLAN.md
    - TASKS.md

agents:
  default_context:
    - PRD.md
    - STACK.md
    - DECISIONS.md

skills:
  symlink_on_init: []
```

---

## Requests layer

### PLAN.md
```md
---
status: planned
priority: medium
owner: 
updated: <today>
summary: "<one-line description>"
related:
  - current/<capability>
---

# <Request title>

## Goal
<What are we trying to achieve?>

## Why
<Why now? What problem does this solve?>

## Scope
- 

## Out of scope
- 

## Approach
<High-level implementation approach>

## Success criteria
- 
```

### TASKS.md
```md
---
updated: <today>
---

# Tasks — <slug>

## <Group name>

- [ ] 
- [ ] 
```

### SESSION.md
```md
---
updated: <today>
---

# Session — <slug>

## Current state
<What's been done so far>

## Active task
<What's being worked on right now>

## Blockers
<Anything blocking progress — or "none">

## Next actions
- 
```

### RISKS.md
```md
---
updated: <today>
---

# Risks — <slug>

## <Risk title>

**Likelihood:** low | medium | high
**Impact:** low | medium | high

**Description:** 

**Mitigation:** 
```

### VERIFY.md
```md
---
updated: <today>
---

# Verify — <slug>

> VERIFY answers "did we build it correctly and safely?" (PLAN answers "did we build the right thing?")

## Manual QA checklist

- [ ] 
- [ ] 

## Edge cases to verify

- [ ] 

## Regression checklist

- [ ] 

## Rollback validation

- [ ] 
```

---

## Current layer

### current/<capability>.md
```md
---
status: draft
updated: <today>
summary: "<one-line description of what this capability does>"
---

# <Capability name>

## Purpose
<What this capability does>

## Requirements
- 

## Scenarios
- 

## Security considerations
- 

## Performance expectations
- 
```
