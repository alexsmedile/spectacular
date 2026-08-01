---
id: DEC-012
type: decision
status: verified
origin: user-derived
derived_from: "User confirmations in request-workflow interview, 2026-08-02"
supersedes: ""
evidence: "user direction"
tags: [cli, request, context, prompt, agents]
updated: 2026-08-02
---

# DEC-012 — Generate implementation briefs on demand from durable request state instead of storing BRIEF.md

**Context:**
Implementation needs a compact self-retrieved starting prompt, but another persisted brief would drift from PLAN, TASKS, SESSION, specs, code, and Git.

**Evidence:**
User direction recorded by explicit command.

**Alternatives considered:**
Recorded in the cited user context or decision discussion.

**Decision:**
Generate implementation briefs on demand from durable request state instead of storing BRIEF.md.

**Consequences:**
request overview remains cheap; --brief compiles authorized execution context; --full returns the request bundle; native agents inspect current code and Git just in time. #cli #request #context #prompt #agents
