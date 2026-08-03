---
id: DEC-023
type: decision
status: verified
origin: user-derived
derived_from: "User migration interview confirmed on 2026-08-03"
supersedes: ""
evidence: "user direction"
tags: [migration, schema, git, branches, compatibility, request, lifecycle]
updated: 2026-08-03
---

# DEC-023 — Migrate workspace schema 2.0 to 3.0 through staged requests

**Context:**
The current repository already uses workspace_schema 2.0 while the roadmap ambiguously described a future v2 workspace contract. A breaking layout migration needs a distinct schema identity, compatibility window, preflight, and recovery boundary.

**Evidence:**
User direction recorded by explicit command.

**Alternatives considered:**
Recorded in the cited user context or decision discussion.

**Decision:**
Migrate workspace schema 2.0 to 3.0 through staged requests

**Consequences:**
Product and workspace-schema versions remain independent. A discovery-only readiness request precedes a separately approved migration request. Schema-3 fields soak additively under schema 2.0; dry-run and dual-shape validation ship before schema 3.0 becomes required. Migration fetches authorized remote state, uses a clean exact baseline and isolated branch, writes snapshots plus a manifest, and stops on divergence or conflict. Only clean fast-forwards and authorized rebases of isolated unpublished branches are mechanical; shared branches are never automatically rebased.
