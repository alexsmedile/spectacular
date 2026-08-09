---
type: owner-decision-evidence
decision: scenario-a-cold-recovery-defaults
version: 1.0
status: recorded
decided_by: owner
decided_at: 2026-08-10
source_thread: 019fe381-5d61-7223-b362-03a5f99a7b69
source_disposition: "A for all"
next_action: execute-scenario-a-autopilot
---

# Scenario A owner dispositions

The owner accepted the recommended option A for every consequential Scenario A decision in one
consolidated reply. This receipt records only those explicit dispositions; it does not accept the
implementation, resolve the Mission, alter the current Capability Contract, or authorize Scenario
B.

| Decision | Explicit disposition |
|---|---|
| Public slice | Implement exactly `anchor show project`, `mission list/show`, `gap list --scope`, `gap/run/checkpoint/evidence/decision show`, and scoped `workspace validate`, with JSON where applicable. |
| Anchor authority | Authoritative project-anchor material and Mission records own anchors; cards, Fog, and continuation are disposable projections. |
| Continuation | Emit a continuation only from a unique, fresh, conflict-free, explicitly authorized chain bound to the expected Mission fingerprint; otherwise emit the exact owner gate. Stored safety or next-action labels are not authority. |
| Human and JSON interfaces | Use versioned deterministic envelopes, source/fingerprint drill-down, explicit freshness and generation basis, stable refusals, and distinct usage/refusal exit codes. |
| Discovery and references | Discover the nearest explicit v2 marker, scan declared roots deterministically, support exact UUIDv7/typed/full-path references only, and expose noun-first drill-down for every consequential pointer. |
| Read guarantee | Reads create no workspace files, locks, caches, normalization, or temporary output; prove path, byte, mode, and mtime nonmutation. |
| Scenario proof | Use an adversarial cold fixture plus a real `v2/.spectacular/` self-hosted workspace with the accepted command, context, and local performance bounds. |
| Execution envelope | Use `codex/feat/v2-scenario-a-cold-recovery`, coherent local commits, no dependency unless demonstrably required, one Go command registry, two focused repair rounds, and fresh Investigator/Builder/Reviewer/Verifier roles. No remote or central-lifecycle effect. |

## Exact accepted thresholds

- Cold recovery uses at most 12 CLI invocations and 64 KiB aggregate JSON.
- Project orientation JSON is at most 24 KiB.
- No manual full-tree walk or authority error is permitted.
- Each command must remain below 500 ms p95 across 20 warm fixture runs on the execution host.
- Valid output, including an exact owner gate, exits 0; usage exits 2; deterministic input or
  workspace refusal exits 3.

## Authority boundary

This owner decision authorizes the bounded Scenario A implementation and its explicit Autopilot
charter. It does not authorize push, PR creation, merge, publish, release, provider mutation,
destructive cleanup, central acceptance, Mission resolution, current-contract reconciliation, or
Scenario B.
