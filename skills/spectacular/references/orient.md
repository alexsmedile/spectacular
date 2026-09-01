# Orient Micro-Kernel

## 1. Trigger Context
Orchestrator resolving ambiguous, cold-start, or uninitialized workspace state.

## 2. CLI Palette & Inspection
```bash
spectacular mission show <ref> --json    # Current execution state
spectacular mission check <ref> --json   # Read-only validation check
spectacular init [--name <project>]      # Initialize greenfield workspace
```

## 3. Negative Constraints (DO NOT)
- **DO NOT** preload git history, old closed missions, or whole directory listings.
- **DO NOT** read `.spectacular/catalog.md` or `index.md` (use CLI `--json` instead).
- **DO NOT** combine multiple active missions into one session; pick exactly one.
- **DO NOT** invent missing anchors; if uninitialized, run `spectacular init` or route to [prepare.md](prepare.md).

## 4. One-Line Preflight Report Format
Emit strictly a 3-line status and proceed without echoing files:
1. **Direction**: Current project outcome and active Mission reference (`M<N>`).
2. **Technical Proof**: Git branch, commit SHA, Contract fingerprint, validation mode.
3. **Next Action**: Exactly one safe next action or one owner gate.
