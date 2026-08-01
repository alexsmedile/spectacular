# Verify — wayfinding-sequencer

## Sequencing
- [x] Strict DAG, ranking, metaphor, boundary, and coherence fixtures pass. `run: bash tests/cli/wayfinding-sequencer.test.sh`
- [x] Existing canonical-ID and record behavior remains compatible. `run: bash tests/cli/wayfinding-contract.test.sh`
- [x] {assert} Status/order/next refuse dangling or cyclic graphs and never present a false frontier.
- [x] {assert} Priority precedes uncertainty; uncertainty deterministically prefers user questions, spikes, research, other questions, then specifications.

## Boundaries
- [x] {assert} Every metaphor route delegates to an existing gated verb; act-on-goal cannot bypass spec confirmation.
- [x] {assert} Cross-layer analysis is read-only and reports inferred edges and target inversions without editing source files or the roadmap.
- [x] Bash syntax and guarded versions remain valid. `run: bash -n cli/spectacular cli/install.sh scripts/hooks/pre-commit && scripts/hooks/pre-commit --check`
- [x] Complete regression suite passes. `run: bash tests/run.sh`
