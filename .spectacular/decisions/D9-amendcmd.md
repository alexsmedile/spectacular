---
type: Decision
id: 01a0103a-ce5c-72bc-a1d9-fb5e2f3b27ad
title: Amend a Contract through its own command, not through Mission completion
created_by: Alex
created: "2026-08-17T14:58:04Z"
updated: "2026-08-17T14:58:04Z"
actor: Alex
actor_role: owner
ref: D9-amendcmd
question: Should amending a Contract ride mission complete, or become a twelfth command?
disposition: amend-frozen-boundary
rationale: >-
    Completion asserts that frozen claims were met; an amendment states that the agreement now says
    something different. Coupling them made `complete` perform two unrelated acts and produced three
    defects: the Gap is resolved when the work lands rather than at completion, so the Contract stays
    knowingly stale for the rest of the Mission; a Mission completes once, so a failed amendment has
    no re-entry; and Contract corrections that belong to no Mission — a typo, a stale `updated:`, the
    miscategorized notice M10/O4 must fix — stay unreachable. Contract is a primitive, as Proposal is.
alternatives:
    - keep the amendment at completion as originally frozen, accepting that Contract corrections outside a Mission stay impossible
    - stop M10 after O2 and open a separate Mission for the amendment command
authority_basis: Owner reviewed the coupling argument, rejected it, and authorized a twelfth command and the corresponding change to M10's frozen outcome, stop, and claims.
authorized_effects:
    - mission.amend-frozen-boundary
    - command.surface-growth
conditions:
    - no-review-bound-to-the-current-activation
    - completed-objectives-preserved
expected_fingerprints:
    - sha256:d6ca09488b06c53c09aae4e7bde6a49ef105339030f2fbc307ddf781805d406a
scope:
    - v2
targets:
    - Mission:01a00f98-0480-7ea2-9f3e-8e3a961aacc6
evidence:
    - Proposal:01a00f7b-e046-700f-9b13-ca4b04d03790
supersedes: ""
---
# Decision

M10 was activated with an outcome, a stop, and two claims stating that a resolved Gap
is closed "at completion in one owner gate" and that the command surface stops at
eleven. Building O3 surfaced that this is not the right shape.

The frozen boundary is therefore corrected rather than worked around: the surface goes
to twelve with `contract amend`, `resolves_gaps:` becomes a declaration of intent that
completion enforces rather than executes, and the amendment runs when the work that
resolves the Gap lands.

O1 and O2 are complete and their code is unaffected — the drift gate and the
`resolves_gaps:` declaration are both required by the corrected design. Re-activation
re-freezes them under the corrected boundary. No review is bound to the superseded
activation fingerprint, so nothing that attested to the old boundary is invalidated.
