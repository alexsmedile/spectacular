# Handoff and Autopilot

Autopilot is explicit and non-default. Before dispatch, compile it with `mission autopilot --input
<json-file> --json`. Bind the exact Mission and owner Decision fingerprints, outcome/non-goals,
authoritative sources, delegated reversible decision domain, providers, allowed actions, the fixed
forbidden-effect ceiling, budgets, checks, expiry, stops, recovery, and return destination.

Create a Handoff only inside the same active Mission and as a subset of its scope, actions,
evidence claims, stops, budget, baseline, and authority. Include sender, actor, runtime-neutral
destination, optional non-authoritative host pointer, immutable inputs, expiry, and required return.
Validate immediately before dispatch and again after runtime replacement.

The receiver returns `succeeded | blocked | failed`, actor, final baseline/result, actions and
native-provider receipts, Evidence, remaining Gaps, budget use, recovery point, and exactly one next
action or owner gate. It never changes Mission lifecycle, claims evidence sufficiency, or gains
provider permission. Supersession creates a new linked Handoff; never edit the original.
