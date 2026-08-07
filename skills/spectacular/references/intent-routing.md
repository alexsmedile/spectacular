---
description: Decide whether a user request needs Spectacular, a new specification, an existing request, or a direct PR-shaped change; capture confirmed intent before drafting a new SPC.
when_to_use: An ordinary change/build/plan request arrives in a Spectacular workspace, especially before a new request or SPC is drafted.
---

# Intent Routing

Do not let repository context choose the work for the user. First decide whether
Spectacular is needed at all; then, only if it is, decide whether the work needs
a new specification. A nearby request, an archive/closure prompt, or a
documentation-impact flag is evidence about the repository—not authorization to
change the user's requested outcome.

## What each layer is for

| Layer | Owns | Does not imply |
|---|---|---|
| Codebase review | Current code, tests, and named files | A workspace read, request, or spec |
| Spectacular read views | Durable workspace status, history, and request context | Permission to mutate a workspace artifact |
| SPC | A proposed durable behavior or contract that needs human approval before implementation | A request, implementation, or a session plan |
| Request `PLAN.md` | Cross-session goal, constraints, milestones, and validation | File-by-file execution steps |
| Request `TASKS.md` | Persistent milestone progress | A duplicate of the agent's current-session tasks |
| Codex/harness plan | Ephemeral in-chat sequencing for the current implementation session | A Spectacular request or any durable workspace write |
| In-chat agent execution | Applying the authorized change and checking it | A PR, commit, push, or a new SPC/request |

Use a Spectacular read view when workspace state answers the user's question.
Use normal repository inspection when code answers it. The presence of
`.spectacular/` never changes a direct code review into a planning ceremony.

## Route precedence

Classify the requested outcome before selecting a Spectacular verb or document
format. Apply these routes in order; stop at the first route that fits:

1. **Inspect or clarify** — read-only inspection; create nothing.
2. **Existing owner** — resume the named or clearly applicable live request,
   including its documentation-impact/spec-sync closure work.
3. **Bounded change** — direct PR-shaped code, documentation, or configuration
   work with a clear outcome and check.
4. **Unresolved meaning** — ask one focused routing question; create nothing.
5. **New durable boundary** — only then consider a new SPC and show its intent
   receipt before writing.

Words such as “spec,” “plan,” “documentation,” “diagram,” or a named file
describe a requested format or surface; they do not establish a new durable
boundary. A new SPC candidate must satisfy all of the following:

- the user wants to establish or change a durable behavior, contract,
  architecture, schema, or security boundary;
- the outcome is not already owned by a live request or its closure work; and
- the outcome cannot truthfully be delivered as one bounded direct change with
  a clear acceptance check.

An explicit natural-language request to draft a spec selects the drafting
format only after this classification. It never bypasses the confirmation
receipt. An explicit terminal `spectacular spec new` command remains the
user-directed mechanical exception described below.

## Route in order

1. **Restate the user's requested outcome.** Use their words where possible.
   Do not replace an outcome with a likely file list, a nearby request, or a
   maintenance opportunity.
2. **Choose the smallest route.**

   | If the request is… | Route |
   |---|---|
   | Asking what exists, what is active, or what a named request says | Read the codebase or the relevant Spectacular view. Do not create state. |
   | A bounded, reversible code/docs/configuration change with clear outcome, likely files, and check | Direct PR-shaped change. Keep the short scope and acceptance check in conversation; no Spectacular artifact is needed. |
   | Work already covered by a live request | Resume that request. Do not create a sibling SPC just because its docs need updating. |
   | A post-implementation documentation or system-spec sync | Use the owning request's docs-impact/spec-sync closure flow. Propose the delta; do not create an implementation SPC unless the user asks for a separate change. |
   | A new behavior/contract, material product or architecture choice, or multi-session/dependent implementation boundary | New SPC candidate. Show the intent receipt before writing. |
   | Missing a fact, feasibility result, business choice, or a future commitment | Route respectively to `RES`, `SPK`, `QUE`, or `IDEA`; do not use an SPC as a holding pen. |

3. **Ask one routing question if two routes remain plausible.** For example:
   “Do you want this documentation sync folded into `<request>`, or are you
   asking for a separate, future-facing documentation capability?”

## Ask, propose, or proceed

| Condition | Agent behavior |
|---|---|
| The outcome, scope, and acceptance check are clear; the user says to implement | Proceed as a direct change. Make a concise in-chat Codex plan when the work has multiple steps. |
| The user explicitly asks for Spectacular, a spec, a request, or durable tracking | Select the matching workflow only after route classification. A natural-language spec request still needs a new durable boundary and the intent receipt. |
| “Implement this plan” points to a plan already visible in chat or named unambiguously | Execute it; do not manufacture a request. |
| “Implement this plan” has no recoverable plan, or could mean materially different outcomes | Ask which plan or restate the two plausible interpretations. Do not infer one from nearby repository context. |
| The task grows from a direct change into a contract/decision/dependency problem | Explain the escalation and propose the smallest Spectacular artifact; wait for confirmation before writing it. |

For the incident-shaped request, “implement the suggestions, draft a plan of
changes, no need for a formal request,” the correct route is **direct change +
in-chat Codex plan**. The suggestions in the immediately preceding conversation
are the scope; “no formal request” explicitly rules out a `PLAN.md`, `TASKS.md`,
or SPC. Inspect the affected skill files, plan the edits in chat, implement and
verify them. Ask only if “the suggestions” cannot be identified from the current
conversation.

## Intent receipt — required before a natural-language SPC draft

Show this compact receipt and wait for an explicit confirmation or correction
before invoking `spectacular spec new` or writing an SPC:

```text
I heard: <short, user-grounded restatement>
Outcome: <observable result the user wants>
Route: new SPC — <why a durable implementation boundary is needed>
Not doing: <closest plausible alternative>
Evidence: <only the repo records that support the interpretation>

Draft this SPC? (yes / correct it)
```

- `I heard` and `Outcome` are the authority. Evidence may refine constraints;
  it must not replace or widen them.
- `Not doing` must name the nearest tempting misroute when one exists—such as
  “not treating a docs-impact flag as a new docs feature.”
- Do not persist raw chat text by default. Once confirmed, the resulting SPC's
  Intent and Evidence sections record the agreed outcome and relevant durable
  evidence, not a transcript.
- An explicit terminal `spectacular spec new <slug> --summary <text>` is a
  user-directed mechanical write. It is not an agent license to infer a spec
  from ordinary conversation.

## Direct PR-shaped change

“PR-shaped” describes a bounded implementation unit, not permission to create
or push a pull request. It should normally have a one-line conversational
contract: outcome, named scope, and acceptance check. Use the normal repository
and Git workflow; create, push, or open a PR only when separately authorized.

Re-route into Spectacular if the work grows into a material contract decision,
cross-session plan, dependency chain, or unresolved user choice. State that
change in routing rather than silently upgrading the work.

## Examples

| User request | Correct route |
|---|---|
| “Fix the README link and run its link check.” | Direct PR-shaped change. |
| “Inspect the CMS editorial docs and tell me whether the README and diagrams disagree with the current contract.” | Read-only inspection; no SPC, request, or receipt. |
| “Update `README.md` and the architecture SVG to describe the already-approved CMS handoff, and run the docs check.” | Direct PR-shaped documentation change; no SPC or request. |
| “Close the current CMS request and make its docs truthful.” | Existing request's docs-impact/spec-sync flow. |
| “Define a durable CMS handoff contract that future adapters must implement.” | New SPC candidate; show receipt before drafting. |
| “Make the CMS docs right for future adapters.” | Ask whether this is an update to existing documentation or a new durable adapter contract; create nothing until answered. |
| “Should the CMS be Sanity or Strapi?” | `QUE` if the user must choose; `RES` if evidence is needed first. |

**Related:** [[new-request]], [[spec-lifecycle]], [[spec-sync]],
[[discovery-protocol]], [[question-rules]], and [[idea-rules]].
