# Mechanical interface

Generated from `internal/command.Registry`; do not edit by hand.

Release version: `2.7.1`

Release inspection: `spectacular --version [--json]` (`spectacular.build-info.v1`, `read-only`)

| Command | Arguments | Schema | Effect |
|---|---|---|---|
| `spectacular mission start` | `<plan.md|-> [--json]` | `spectacular.mission.start.v2` | `mutating` |
| `spectacular mission show` | `<ref> [--json]` | `spectacular.mission.show.v2` | `read-only` |
| `spectacular mission check` | `<ref> [--json]` | `spectacular.mission.check.v2` | `read-only` |
| `spectacular objective show` | `<mission-ref>/<objective-ref> [--json]` | `spectacular.objective.show.v2` | `read-only` |
| `spectacular objective promote` | `<mission-ref>/<objective-ref> [--json]` | `spectacular.objective.promote.v2` | `mutating` |
| `spectacular objective finish` | `<mission-ref>/<objective-ref> [--json]` | `spectacular.objective.finish.v2` | `mutating` |
| `spectacular run show` | `<mission-ref>/<run-ref> [--json]` | `spectacular.run.show.v2` | `read-only` |
| `spectacular run start` | `<mission-ref>[/<objective-ref>] --title <title> [--json]` | `spectacular.run.start.v2` | `mutating` |
| `spectacular run transition` | `<target-ref> --to <state> --by <actor> --reason <text> [--next-action <action>] [--json]` | `spectacular.run.transition.v2` | `mutating` |
| `spectacular review record` | `<mission-ref> <review.md|-> [--json]` | `spectacular.review.record.v2` | `mutating` |
| `spectacular handoff record` | `<mission-ref> <handoff.md|-> --by <sender> [--json]` | `spectacular.handoff.record.v2` | `mutating` |
| `spectacular evidence record` | `<mission-ref> <evidence.md|-> [--json]` | `spectacular.evidence.record.v2` | `mutating` |
| `spectacular mission complete` | `<ref> --by <owner> [--json]` | `spectacular.mission.complete.v2` | `mutating` |
| `spectacular proposal check` | `<ref> [--json]` | `spectacular.proposal.check.v2` | `read-only` |
| `spectacular campaign check` | `<path> [--json]` | `spectacular.campaign.check.v2` | `read-only` |
| `spectacular contract amend` | `<contract-ref> --gap <gap-ref> --by <owner> [--resolution <text>] [--dry-run] [--json]` | `spectacular.contract.amend.v2` | `mutating` |
| `spectacular charter` | `<mission-ref>/<objective-ref> [sources...] [--json]` | `spectacular.charter.show.v2` | `read-only` |
| `spectacular decide` | `<decision.md|-> [--json]` | `spectacular.decision.record.v2` | `mutating` |
