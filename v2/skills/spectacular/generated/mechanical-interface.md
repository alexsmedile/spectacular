# Mechanical interface

Generated from `internal/command.Registry`; do not edit by hand.

| Command | Arguments | Schema | Effect |
|---|---|---|---|
| `spectacular anchor show project` | `[--json]` | `spectacular.anchor.show.v1` | `read-only` |
| `spectacular mission list` | `[--json]` | `spectacular.mission.list.v1` | `read-only` |
| `spectacular mission show` | `<ref> [--json]` | `spectacular.mission.show.v1` | `read-only` |
| `spectacular gap list` | `--scope <ref> [--json]` | `spectacular.gap.list.v1` | `read-only` |
| `spectacular gap show` | `<ref> [--json]` | `spectacular.gap.show.v1` | `read-only` |
| `spectacular run show` | `<ref> [--json]` | `spectacular.run.show.v1` | `read-only` |
| `spectacular checkpoint show` | `<ref> [--json]` | `spectacular.checkpoint.show.v1` | `read-only` |
| `spectacular evidence show` | `<ref> [--json]` | `spectacular.evidence.show.v1` | `read-only` |
| `spectacular decision show` | `<ref> [--json]` | `spectacular.decision.show.v1` | `read-only` |
| `spectacular workspace validate` | `<scope> [--json]` | `spectacular.workspace.validate.v1` | `read-only` |
| `spectacular workspace context` | `<project|mission-ref> --event <@Event> [--selector <$domain.verb>] [--json]` | `spectacular.workspace.context.v1` | `read-only` |
| `spectacular proposal show` | `<ref> [--json]` | `spectacular.proposal.show.v1` | `read-only` |
| `spectacular proposal check-base` | `<ref> [--json]` | `spectacular.proposal.check-base.v1` | `read-only` |
| `spectacular proposal create` | `--input <json-file> [--json]` | `spectacular.proposal.create.v1` | `mutating` |
| `spectacular mission prepare` | `--input <json-file> [--json]` | `spectacular.mission.prepare.v1` | `read-only` |
| `spectacular mission create` | `--input <json-file> [--json]` | `spectacular.mission.create.v1` | `mutating` |
| `spectacular mission transition` | `<ref> --to <state> --authorization <decision-ref> --expected-fingerprint <sha> --idempotency-key <key> [--assessment <ref>] [--reconciliation <ref>] [--disposition <value>] [--terminal-next-action <text>] [--satisfied-objectives <ref,ref>] [--json]` | `spectacular.mission.transition.v1` | `mutating` |
| `spectacular mission autopilot` | `--input <json-file> [--json]` | `spectacular.mission.autopilot.v1` | `read-only` |
| `spectacular handoff show` | `<ref> [--json]` | `spectacular.handoff.show.v1` | `read-only` |
| `spectacular handoff validate` | `<ref> [--json]` | `spectacular.handoff.validate.v1` | `read-only` |
| `spectacular handoff create` | `--input <json-file> [--json]` | `spectacular.handoff.create.v1` | `mutating` |
| `spectacular handoff return` | `--input <json-file> [--json]` | `spectacular.handoff.return.v1` | `mutating` |
| `spectacular evidence create` | `--input <json-file> [--json]` | `spectacular.evidence.create.v1` | `mutating` |
| `spectacular decision create` | `--input <json-file> [--json]` | `spectacular.decision.create.v1` | `mutating` |
| `spectacular assessment record` | `--input <json-file> [--json]` | `spectacular.assessment.record.v1` | `mutating` |
| `spectacular contract show` | `<ref> [--json]` | `spectacular.contract.show.v1` | `read-only` |
| `spectacular contract reconcile` | `<ref> --proposal <ref> --authorization <decision-ref> --expected-fingerprint <sha|absent> --idempotency-key <key> [--json]` | `spectacular.contract.reconcile.v1` | `mutating` |
| `spectacular contract reconcile-set` | `--input <json-file> [--json]` | `spectacular.contract.reconcile-set.v1` | `mutating` |
| `spectacular mission archive` | `<ref> --authorization <decision-ref> --expected-fingerprint <sha> --idempotency-key <key> --terminal-packet <mission-ref> [--json]` | `spectacular.mission.archive.v1` | `mutating` |
