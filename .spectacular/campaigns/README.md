# Campaigns

Campaigns are optional, durable Markdown roadmap maps. They sequence several
candidate or active Missions around one strategic outcome, but they grant no
execution authority and are intentionally excluded from Spectacular's typed
record graph and CLI lifecycle.

An Atlas in `../atlas/` explains the user journeys, capabilities, and system
connections behind a Campaign. A Campaign sequences the work; it does not
replace that explanation.

Create one file per genuinely independent roadmap arc. Do not create a Campaign
for every project or use one as an unattended automation queue.

```md
---
campaign_schema: spectacular.campaign.v1
title: Launch readiness
focus: Establish a safe release path.
atlas: ../atlas/release-confidence.md
current: B2
exit_condition: Releases are verified and repeatable.
blocks:
  - ref: B1
    title: Release hardening
    state: complete
    after: []
    advances: [publish-a-trustworthy-release]
    enables: [repeatable-release]
    missions: [M13]
  - ref: B2
    title: Hosted release
    state: active
    after: [B1]
    missions: []
---
# Campaign: Launch readiness

> Planning map only. It grants no execution authority.

## Decisions and non-goals
```

Every Campaign Block must state the value connection it advances. Cite one or
more journey steps with `advances:` when the work directly changes a user
workflow. For enabling work—security, reliability, migration, or developer
tooling—cite `enables:` and state its eventual user or business consequence in
the block's prose. Do not invent a customer journey merely to satisfy the map.

`advances:` and `enables:` are planning conventions today: `campaign check`
does not yet validate them. Keep their references readable and point readers to
the relevant Atlas file.

A Mission may link to a Campaign block in its Markdown origin or rationale. That
link is context only; it must not become a frozen Mission binding.

When a Campaign benefits from a mechanical order and Mermaid projection, keep
the compact frontmatter map shown above, then run:

```sh
spectacular campaign check .spectacular/campaigns/<campaign>.md
```

The command is read-only. It validates the map, current block, dependency order,
and any named Mission refs; detects cycles; and emits an ordered Mermaid
projection. A Campaign that does not need this projection remains ordinary
Markdown.

`current` is the Campaign's global map position, not an instruction to every
agent that reads it. Mission workers follow their assigned Mission, Objective,
and Run. An orchestrator may embed the generated Mermaid below the matching
markers; `campaign check` verifies it has not drifted:

````md
<!-- spectacular:campaign-mermaid:start -->
```mermaid
... output from campaign check ...
```
<!-- spectacular:campaign-mermaid:end -->
````
