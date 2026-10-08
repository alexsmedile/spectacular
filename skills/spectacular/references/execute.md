# Execute Micro-Kernel

Use this when: an orchestrator is activating or resuming a Mission, or a worker is executing an assigned code charter.

## 2. Start or resume

For a new owner-authorized activation, prepare the agreed branch/worktree and
run mission start once, then inspect and check the returned Mission. Finish
preparation when its named scope, binding, and authorized next work are known.

```bash
spectacular mission start plan.md --json
spectacular mission show <ref> --json
spectacular mission check <ref> --json
```

For resume, inspect the named Mission and its current Run/Objective first. Locate
and reuse its existing execution branch/worktree with native Git; preserve its
baseline, activation fingerprint, Contract binding, and current state. Resume only
when the check is valid, the Mission permits execution, the checkout matches its
execution context, and the next task is within existing authority. A completed or
stopped Mission is not new activation authority. If the checkout is missing or
drifted, report the exact mismatch and resolve it before changing governed work.
Do not run mission start again to recover an existing Mission.

## 3. Negative Constraints (DO NOT)
- **DO NOT** execute on `main` branch. Use the agreed execution branch; create one for a new activation, reuse it on resume.
- **DO NOT** output meta-planning narrative. Directly write code, run tests, and return results.
- **DO NOT** manage Spectacular files inside worker subagents (workers ignore `.spectacular/`).
- **DO NOT** touch files outside `allowed_changed_paths` defined in the charter.

## 4. Domain-specific implementation

For requested queue or batch processing, follow the target project's agreed
concurrency, retry, shutdown, and failure contracts. Choose its idioms from the
actual implementation rather than introducing generic worker/DLQ defaults.

## 5. Authority & Context Invariants
- **Progressive Context**: Drill down strictly: `Mission card -> current Objective -> exact sources`.
- **Authority Separation**: The `plan carries meaning`, while `tooling carries repeatability`.
- **Core-First Gate**: In shared engine setups, core independent verification must pass before downstream client integration.
- **Worktree Pre-warming**: Lead orchestrator must pre-warm native dependencies and env sandboxes in `.worktrees/<slug>` before worker dispatch.
- **Activation Boundary**: `A Decision is not activation authority` (only owner confirmation authorizes `mission start`).
- **Self-Hosting & Bootstrap**: When developing Spectacular, an `active Mission keeps the schema` frozen. Under declared `manual-bootstrap`, run `focused checks` directly.
