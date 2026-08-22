# Campaigns

Campaigns are optional, durable Markdown roadmap maps. They sequence several
candidate or active Missions around one strategic outcome, but they grant no
execution authority and are intentionally excluded from Spectacular's typed
record graph and CLI lifecycle.

Create one file per genuinely independent roadmap arc. Do not create a Campaign
for every project or use one as an unattended automation queue.

```md
# Campaign: Launch readiness

> Planning map only. It grants no execution authority.

## Strategic outcome

## Exit condition

## Blocks

| Block | Depends on | Candidate / active Mission | State |
| --- | --- | --- | --- |
| Release hardening | CI baseline | `M13` | complete |
| Hosted release | Release hardening | — | proposed |

## Decisions and non-goals
```

A Mission may link to a Campaign block in its Markdown origin or rationale. That
link is context only; it must not become a frozen Mission binding.
