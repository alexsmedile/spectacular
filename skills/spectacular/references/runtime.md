# Handoff and Autopilot

Autopilot is explicit and non-default. Before dispatch, compile it with `mission autopilot --input
<json-file> --json`. Bind the exact Mission and owner Decision fingerprints, outcome/non-goals,
authoritative sources, delegated reversible decision domain, providers, allowed actions, the fixed
forbidden-effect ceiling, budgets, checks, expiry, stops, recovery, and return destination.

Create a Handoff only inside the same active Mission and as a subset of its scope, actions,
evidence claims, stops, budget, baseline, and authority. Include sender, actor, runtime-neutral
destination, optional non-authoritative host pointer, immutable inputs, expiry, and required return.
Validate immediately before dispatch and again after runtime replacement.

Fan out only cohesive, outcome-sized work. Each Handoff names one Objective, that Objective's exact
dependencies, a disjoint claim scope, authoritative `ref@fingerprint` inputs, and an explicit return
contract. Do not create recursive critic loops or tiny sessions; finish working code and batch review
at the frozen `automatic | clustered | independent` level for each claim.

Optional wall-time, token, spend, parallel-worker, and repair-round caps must disclose
`hard | observed | unsupported` enforcement. `hard` is valid only when the independently validated
host envelope supplies both a measurement capability and a cancellation capability; `observed` declares measurement without a
kill switch; `unsupported` promises neither. Never describe a monitored estimate as a hard cap.

The receiver returns `succeeded | blocked | failed`, actor, final baseline/result, actions and
native-provider receipts, Evidence, remaining Gaps, budget use, recovery point, and exactly one next
action or owner gate. It never changes Mission lifecycle, claims evidence sufficiency, or gains
provider permission. Supersession creates a new linked Handoff; never edit the original.
