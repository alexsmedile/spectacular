# Workspace containers and record ownership

Use this when: choosing among requirements, plans, specs, scenarios, proposals, Contracts, Atlas, audits,
reviews, handoffs, or explaining an active project’s tree.

## Choose by the question the file answers

| Home | Question | Mode |
|---|---|---|
| Root anchors | What is true about this project now? | Maintained context |
| raw/ | What did we capture or sketch? | Unconstrained |
| plans/ | What will we do, and what happened? | Maintained context |
| requirements/ | What is needed and how will acceptance be observed? | Optional maintained context |
| specs/ | What behavior or interface should exist? | Optional maintained context |
| scenarios/ | Can a user reach a concrete outcome, and how will we verify it? | Optional maintained context |
| atlas/ | What are the objects and how do they relate? | Maintained context |
| decisions/ | What did we choose and why? | Soft or explicitly governed |
| audits/ | What did inspection find at a named revision? | Optional maintained context |
| proposals/ | Which formally tracked question remains open? | Existing governed records |
| contracts/ | What agreement is mechanically bound? | Mechanically governed |
| missions/ | What explicitly authorized frozen slice is executing? | Mechanically governed |
| campaigns/ | How do several plans or Missions fit together? | Optional maintained context |
| retrospectives/ | What lessons deserve independent retention? | Optional maintained context |
| Root docs/ | How does a person use the shipped product? | External documentation surface |

Create optional requirements, specs, scenarios, or audits only when the content has independent use. A plan
can carry its own requirements and inspection findings. An audit note has scope,
revision inspected, findings, basis, and limitations; it does not certify governed
completion. A spec can state intended behavior without becoming an accepted Contract.
A reusable Scenario owns the journey and success criteria; a plan selects cases
and an Audit records observed outcomes. See [scenarios.md](scenarios.md).
Existing Proposal identity/lifecycle rules remain; `proposal create` is unavailable.
Ordinary unanswered questions can stay in raw or the plan without a Proposal.

## Governed records follow their owner

The CLI’s human-layout writer places Mission-owned reviews, evidence, handoffs,
Gaps, and assessments in the Mission bundle. Promoted Objectives and Runs appear
there only when their supported transitions create them; compact inline work does
not need extra folders. Project-level records can exist in their corresponding
root collections when their enforced format permits it. Use the path reported by
the command; do not infer a path from a diagram or hand-move a governed record.

A quick continuity note belongs in its plan. Extract a separate linked Reference
under plans only when it has an independent reader; a governed Handoff instead
uses the supported `handoff record` path. Audit observations and a governed Review
are different artifacts: record the latter only when that workflow is requested.

For naming, task-sized reading routes, configuration and inventory roles, and
retirement guidance, see [navigation.md](navigation.md).

## Example of an active project

These names illustrate placement, not templates or records to manufacture.

```text
project/
├── docs/                         # public docs, managed by documentation tools
│   └── guides/release.md
└── .spectacular/
    ├── PROJECT.md
    ├── ONTOLOGY.md
    ├── INDEX.md                  # optional manual routes
    ├── index.json                # generated governed inventory only
    ├── raw/release-sketch.md
    ├── plans/release-flow.md      # approach, progress, result; type Plan
    ├── plans/linux-tranche.md     # independent delivery slice; linked by parent
    ├── requirements/verified-download.md # need and acceptance; type Requirement
    ├── specs/release-integrity.md # behavior and interfaces; type Spec
    ├── scenarios/verify-download.md # user journey and final output; type Scenario
    ├── atlas/release-artifact.md  # entity and relationships; type Entity
    ├── decisions/checksum-policy.md  # soft Decision, no governed identity
    ├── audits/release-drift.md    # inspected revision and findings; type Audit
    ├── contracts/CC-r7u2p4-release-integrity.md
    ├── proposals/P7-signing-options.md
    └── missions/M23-release-integrity/
        ├── M23-release-integrity.md
        ├── reviews/RV1-release-review.md
        ├── evidence/E1-x4b8q2.md
        ├── handoffs/H1-k4p2a8.md
        └── objectives/O2-linux-delivery.md  # only if promoted
```

Archived governed bundles retain existing rules under archive/. Historical
collections and legacy indexes can remain. Never create this entire tree merely
to start a task. Requirements, plans, specs, scenarios, and audits use the basic folder metadata agreement;
governed files use their command-generated schema and validator. Raw is exempt.
Select metadata using [profiles.md](profiles.md), with compact as the soft default.
