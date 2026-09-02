# Prepare Micro-Kernel

## 1. Trigger Context
Orchestrator preparing Greenfield Genesis, Campaign roadmap, Architectural Decisions, or Mission drafting.

## 2. CLI Palette & Minimal Drafting
```bash
spectacular init [--name <project>]                               # Initialize workspace
spectacular decide --title "<Title>" --disposition "<Choice>" --rationale "<Why>" --json  # Direct Decision
spectacular mission start plan.md --json                          # Activate from minimal plan
spectacular mission check <ref> --json                            # Validate generated mission
```

## 3. Minimal Draft Grammar (Zero YAML Boilerplate)
Never hand-write timestamps, UUIDs, or hash fingerprints. Write minimal drafts; the CLI auto-populates metadata:

```yaml
# Minimal Mission Plan (plan.md)
ref: M1
slug: storage-engine
outcome: Implement in-memory priority queue with DLQ routing.
contract: CC-queue
completion:
  - claim: priority-scheduling
    pass_boundary: High priority jobs run before normal/low.
    proof_requirement: sh tests/check.sh exits 0.
objectives:
  - ref: O1
    outcome: Build queue worker pool and dlq routing.
    claims: [priority-scheduling]
```

## 4. Negative Constraints (DO NOT)
- **DO NOT** conduct multi-turn unprompted grilling. For Interview Mode, batch choices into Tier 2/3 cards (numbered questions, lettered options with recommended defaults) and confirm once.
- **DO NOT** create multi-folder mission directories for routine code (Tier 1 Single-File `M<N>.md` is default).
- **DO NOT** write Proposals for obvious tasks; draft directly to Mission or Decision.
- **DO NOT** echo YAML or generated mission contents into chat (Silent Mutation).

## 5. One-Shot Genesis Flow
1. Run `spectacular init`.
2. Author `.spectacular/PROJECT.md` (Boundaries & Constraints).
3. Record key technical choices via `spectacular decide --title ...`.
4. Run mutating launch `spectacular mission start plan.md --json` only after owner confirmation, then verify with read-only check `spectacular mission check <ref> --json`.
