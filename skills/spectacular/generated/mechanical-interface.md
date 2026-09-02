# Mechanical interface

Generated from `internal/command.Registry`; do not edit by hand.

Release version: `2.15.0`

Release inspection: `spectacular --version [--json]` (`spectacular.build-info.v1`, `read-only`)

| Command | Arguments | Schema | Effect |
|---|---|---|---|
| `spectacular mission start` | `[<plan.md|->] [--data <json>] [--allow-main] [--create-branch] [--json]` | `spectacular.mission.start.v2` | `mutating` |
| `spectacular mission list` | `[--all] [--status <status>] [--json]` | `spectacular.mission.list.v2` | `read-only` |
| `spectacular mission show` | `<ref> [--json]` | `spectacular.mission.show.v2` | `read-only` |
| `spectacular mission check` | `<ref> [--verify] [--json]` | `spectacular.mission.check.v2` | `read-only` |
| `spectacular mission amend-scope` | `<ref> --add <paths> --by <owner> [--reason <text>] [--dry-run] [--json]` | `spectacular.mission.amend_scope.v2` | `mutating` |
| `spectacular mission close` | `<ref> [--by <owner>] [--json]` | `spectacular.mission.close.v2` | `mutating` |
| `spectacular objective show` | `<mission-ref>/<objective-ref> [--json]` | `spectacular.objective.show.v2` | `read-only` |
| `spectacular objective promote` | `<mission-ref>/<objective-ref> [--json]` | `spectacular.objective.promote.v2` | `mutating` |
| `spectacular objective finish` | `<mission-ref>/<objective-ref> [--json]` | `spectacular.objective.finish.v2` | `mutating` |
| `spectacular run show` | `<mission-ref>/<run-ref> [--json]` | `spectacular.run.show.v2` | `read-only` |
| `spectacular run start` | `<mission-ref>[/<objective-ref>] --title <title> [--json]` | `spectacular.run.start.v2` | `mutating` |
| `spectacular run transition` | `<target-ref> --to <state> --by <actor> --reason <text> [--next-action <action>] [--data <json>] [--json]` | `spectacular.run.transition.v2` | `mutating` |
| `spectacular review record` | `<mission-ref> [<review.md|->] [--data <json>] [--json]` | `spectacular.review.record.v2` | `mutating` |
| `spectacular handoff record` | `<mission-ref> [<handoff.md|->] [--by <sender>] [--data <json>] [--json]` | `spectacular.handoff.record.v2` | `mutating` |
| `spectacular evidence record` | `<mission-ref> [draft.md|-] [--from <test-output>] [--data <json>] [--json]` | `spectacular.evidence.record.v2` | `mutating` |
| `spectacular mission complete` | `<ref> [--by <owner>] [--json]` | `spectacular.mission.complete.v2` | `mutating` |
| `spectacular proposal check` | `<ref> [--json]` | `spectacular.proposal.check.v2` | `read-only` |
| `spectacular campaign check` | `<path> [--ascii] [--json]` | `spectacular.campaign.check.v2` | `read-only` |
| `spectacular contract amend` | `<contract-ref> --gap <gap-ref> --by <owner> [--resolution <text>] [--dry-run] [--json]` | `spectacular.contract.amend.v2` | `mutating` |
| `spectacular contract create` | `<ref> [--title <title>] [--json]` | `spectacular.contract.create.v2` | `mutating` |
| `spectacular charter` | `<mission-ref>/<objective-ref> [sources...] [--prompt] [--json]` | `spectacular.charter.show.v2` | `read-only` |
| `spectacular decide` | `[<decision.md|->] [--title <title>] [--disposition <accepted|rejected|deferred|superseded>] [--rationale <rationale>] [--actor <name>] [--supersedes <ref>] [--json]` | `spectacular.decision.record.v2` | `mutating` |
| `spectacular init` | `[<path>] [--name <name>] [--json]` | `spectacular.init.v2` | `mutating` |
| `spectacular guard` | `<mission-ref>/<objective-ref> [--watch] [--exec <command>] [--json] [-- <command...>]` | `spectacular.guard.v2` | `read-only` |
