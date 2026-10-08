# Internal module boundaries

The CLI remains one Go module and one executable. Internal packages own separate
capabilities; public CLI commands and mechanical schemas remain in the canonical
`command.Registry`.

| Package | Responsibility | Boundary |
|---|---|---|
| `domain` | Identity, record grammar, typed refusals | No workspace or process effects; UUID library explicitly allowed |
| `missionview` | Detached data for Mission readers | No imports, documents, services, or mutable storage references |
| `runpolicy` | Run transition rules | Pure decisions and typed refusals |
| `workspace`, `discovery` | Canonical documents and record discovery | Parsing and storage adapters |
| `missionbundle` | Governed use cases and atomic mutations | One decoder; projections exposed through `ReadView` |
| `charter` | Charter compilation and budget handling | Receives a `MissionLoader`; cannot import Mission mutation services |
| `guard` | Supervised process effects | Receives the loader from composition |
| `command` | Registry, argument handling, composition, output | Handlers grouped by Mission lifecycle, records, and execution context |

`missionbundle.Service` remains the use-case facade. Its implementations are
organized by capability in `service_*.go`; `git.go` owns Git integration and
`persistence.go` owns locking, canonical writes, and transaction application.
Completion and Gap resolution retain their coordinated transaction. The Run
transition decision delegates to `runpolicy` before changing any document.

`ReadView` projects through the existing decoder, including resolved expanded
records. It copies nested slices so consumers cannot mutate authoritative bundle
data. Charter and Guard receive the loader explicitly; production composition
supplies `missionbundle.ReadView`. Test consumers may supply detached views.

## Automatic enforcement

`scripts/architecture-boundaries.json` declares allowed internal imports for each
production package. `domain`, `missionview`, and `runpolicy` also have explicit
external-import allowlists. A new package or dependency requires an intentional
policy change. `internal/` alone restricts external consumers, not sibling imports.

Run `go run ./scripts/check-architecture` from the repository root. The checker
parses production imports for every platform, refuses unclassified packages and
forbidden edges, and detects cycles across the combined graph. Tests are excluded
from the production matrix so integration fixtures can compose real adapters.
The checker runs in the static verification gate; its negative tests run in the
quick gate. Build and runtime checks remain necessary alongside this check.
