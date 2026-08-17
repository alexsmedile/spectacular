# Explore and plan

## Explore

A Proposal is optional. Use one only when the exploration deserves a durable home
— for the problem, alternatives, open questions, research, or a draft
specification. Otherwise leave it in an issue or in the conversation.

A Proposal is mutable. It is neither authority nor current product truth.

Read the current Contracts and specifications before proposing observable
behavior. Once a direction is frozen, edit those files as ordinary Mission work.
There is no separate reconciliation lifecycle.

Check a Proposal that has a durable home:

```bash
spectacular proposal check <ref>
```

## Plan

Compare only approaches that genuinely differ and are outcome-sized. Weigh each on:

- observable result, and the proof it would need
- coherence with what exists
- dependencies and reversibility
- learning value
- integration path
- what happens if it is cancelled

Record one verdict: `sufficient | needs-evidence | needs-decision`.

Then grill only what is still unresolved — criteria, scope, dependencies, risks,
or blocking Gaps. Do not re-interview settled ground.

## Freeze a compact Mission preview

Frontmatter:

- title, owner, outcome, applicable Contract, Git baseline
- one completion claim per verifiable domain, each with a pass boundary and a
  proof requirement
- review level: `automatic | clustered | independent`, defaulted once when shared.
  Choose `independent` when any claim touches security, privacy, or rights;
  stored data or a migration; a shared or public interface; compatibility; more
  than one system boundary; an external provider; a destructive or
  hard-to-reverse effect; a production or observational claim; a material
  architecture change; a novel pattern; or evidence only the executor can see.
  Also choose it when the work is disputed. Otherwise `automatic` is honest, and
  `clustered` fits several small related claims. A reviewer who did not implement
  the scope is what makes it independent — see [close.md](close.md).
- Objectives, with dependencies and claim coverage
- initial Run and operator, authority, mechanical and semantic scope
- budgets, dependencies, Gaps, stops, recovery
- `resolves_gaps:` when the Mission closes a Gap on its bound Contract, as `gap`
  and `resolution` pairs. Both are frozen, so the owner approves the exact wording
  at activation and the Mission cannot gain amend authority afterwards. Completion
  refuses while a declared Gap is still open. Requires `amend-contract` in
  `requires_owner`.

Markdown body:

- origin and rationale
- the detailed execution plan
- conditional bootstrap and review notes

A claim is the part most often written too vaguely. It needs a boundary that can
fail, and a proof that names the test:

```yaml
completion:
    - claim: drift-flags
      pass_boundary: Each frozen completion claim carries named drift flags
          derived from repairs consumed, evidence age, verdict state, and
          fingerprint age.
      proof_requirement: Table-driven fixtures with known repair counts and
          evidence ages assert the exact flag set, the ranking, and the default
          selection including tie behavior.
```

`pass_boundary` states what must be observably true. `proof_requirement` states
what would demonstrate it. "Works correctly" is neither.

Present the preview **once**, in chat. Owner confirmation freezes the semantic
envelope.

The preview is a plan document, not yet a Mission. It carries no UUID, no
activation block, and no fingerprint — `mission start` generates those. For the
field-by-field shape of what it becomes, see
[mission-anatomy.md](mission-anatomy.md).

## Then activate

```bash
spectacular mission start plan.md --json   # or: ... start - --json  (stdin)
spectacular mission check <ref> --json     # confirm what was generated
```

It generates identities, bindings, activation, and the canonical path at
`.spectacular/missions/<slug>/MISSION.md` — atomically, from the approved plan.

Under a declared manual bootstrap, hand-author that file, generate valid
identities, and verify the structure directly — see [bootstrap.md](bootstrap.md).
