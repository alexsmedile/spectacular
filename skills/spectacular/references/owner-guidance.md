# Owner Interaction Guidance

Use this when: Orchestrator or operator approaches an owner gate, question, decision, or authority boundary.

## 1. When to Ask Owner Questions

Ask only when at least one of these is genuinely open:
- **Semantic fork**: Competing architectural or product directions with distinct trade-offs.
- **Mission boundary**: Scope expansion, modification of frozen completion criteria, or non-goal adjustments.
- **Authority gate**: Any action requiring `requires_owner` (e.g. `mission start`, `mission complete`, `contract amend`, `git push`, provider mutations).
- **Material risk / irreversible effect**: Destructive operations, data migrations, or external integrations.
- **Contract conflict**: Incompatibility with frozen capability specifications or upstream truth.

Do not run an interview when the Mission or project Anchors already answer the question. Carry accepted defaults forward.

### Assumption Calibration (Zero-Friction Progress)
- **Material / Semantic Assumptions**: If a choice affects user-visible behavior, public APIs, architectural boundaries, risk profile, or the completion claim, stop and present a four-part question or Gap.
- **Reversible Implementation Choices**: For low-level, reversible implementation choices strictly within authorized scope, choose the simplest path, state the assumption clearly in the handoff or Run body notes, and continue executing without interrupting the owner.

---

## 2. Four-Part Question Format

When presenting an open decision, use four structured parts:

```text
1. Plain Outcome: What happens from the user/product perspective.
2. Technical Basis: Why this choice exists (constraints, trade-offs, invariants).
3. Concrete Options: Specific choices formatted as "action -> consequence".
4. Recommended Default: The recommended path and the precise reason for it.
```

---

## 3. Ask for Authorization, Not for Labor

An owner-gated action needs the owner's **decision**, not the owner's **hands**. When reaching an owner gate, request authorization and specify the action you will execute upon approval:

> `mission complete` is owner-gated because completion freezes the claims and records attributable acceptance.
> Authorize me to run `spectacular mission complete M14 --by Alex`? (Y/N)

```text
┌─────────────────────────────────────────────────────────────┐
│ HOLD THE KEYBOARD PRINCIPLE                                │
│ • Ask for explicit confirmation, then run the command.     │
│ • Never hand back a command for the owner to copy/paste    │
│   when you can execute it upon explicit approval.           │
│ • Reserve owner-executed steps strictly for operations     │
│   where performing it IS the authority (credential entry,   │
│   interactive logins, legal signatures).                   │
└─────────────────────────────────────────────────────────────┘
```

When acting on an authorization, record who authorized and who performed. An operator acting on owner approval is not the owner acting.

---

## 4. Settle Execution Mode at Activation

When approaching the activation gate, route to [execute.md](execute.md) for the execution-mode table, checkpoint rules, and Run-body template. Settle owner involvement in the same activation exchange; do not ask piecemeal later.

---

## 5. Batch Gates; Ask Once

Collect every owner decision required for a phase and present them in a single exchange. Multiple piecemeal approvals create friction and delay.

- **Check prior decisions first**: Before asking, check whether the answer is already recorded in Anchors, Decisions, or earlier approvals.
- **Approvals carry forward**: An approval carries forward to the same class of action within the active phase.
- **Avoid duplicate gates**: Re-asking settled questions is noise, not governance.

---

## 6. State Boundaries Once; Refusal Format

Name an authority boundary the first time it becomes relevant, then act on it. Repeating constraints on every turn obscures where a decision is genuinely needed.

When you must decline an out-of-scope or unauthorized action, use this compact pattern:

```text
1. What cannot be done (1 sentence).
2. Nearest supported alternative (1 sentence).
3. Immediate unblocked next step (continue execution).
```

*Example:*
> "I cannot deploy directly to production without owner authorization. I can run the full staging verification suite and prepare the release artifact. Proceeding with staging tests."
