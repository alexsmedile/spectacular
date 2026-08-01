---
description: Canonical freshness, loading, archival, and garbage-collection contract for live, stale-safe, temporary, and throwaway artifacts.
when_to_use: Deciding what must stay synchronized, what may remain stale, what leaves active context, or what may be deleted.
---

# Artifact Retention Contract

Retention is derived from an artifact's entity, lifecycle status, and path. Do not add a `retention:` field to every file: duplicated labels would become another freshness problem.

## Four classes

| Class | Promise | Normal loading | Terminal rule |
|---|---|---|---|
| **Live** | Trustworthy enough to plan or execute current work | Briefings and task-scoped reads | Re-evaluate at its named lifecycle checkpoints |
| **Stale-safe** | Durable history; accurate about what was known/done then, not necessarily now | Index first, body only on explicit historical lookup | Keep or compact behind an index; check against code before reuse |
| **Temporary** | Bounded working context; allowed to evolve while its owner is active | Only with its owning request/session/discovery flow | Promote durable learning or archive the record when the owner closes |
| **Throwaway** | The file/code itself has no lasting authority | Never part of normal retrieval | Preserve outcome/evidence and recovery boundary, then delete |

“Stale-safe” is not neglect. It is an explicit statement that historical material may be useful without competing with current implementation truth.

## Live state

| Artifact | Freshness obligation | Checkpoint |
|---|---|---|
| Production code + executable unit/invariant tests | Ultimate implementation truth after verified integration | Every implementation/verification run |
| `roadmaps/index.md` | Current build ledger, active/planned release direction, and shipped-file index agree | Major roadmap/release change, heavy request closure, session end |
| `specs/index.md` | Cheap present-tense capability summary agrees with verified code at its last review | Request archive/spec delta and explicit currency audit |
| Open/deferred `questions/` | Contains only unresolved fog; open `requires_user_input` blockers are accurate | Start of every human-agent session and before sequencing |
| Active request `PLAN.md`/`TASKS.md`/`SESSION.md` | Current goal, scope, progress, blockers, and next action agree | Meaningful work boundary and session handoff |
| Draft/unconfirmed/approved SPC for the active release | Developer/user decisions stay synchronized until code generation starts | Interview/decision change and approval/action gate |
| PRD/principles/architecture/policy | Current strategic and operating intent, not implementation detail | Major product/architecture decision |

Live does not mean continuously polled. It means a documented event owns re-evaluation. Indexes summarize; they do not duplicate full bodies.

At the beginning of every human-agent session, run `spectacular wayfind status --blockers-only` before lower-priority briefing content. Surface highest priority, then downstream impact. Deferred questions are not immediate blockers.

## Temporary working context

- An approved spec becomes a temporary execution contract once code generation begins. It may diverge as implementation teaches us; record behavior-changing departures in the request/spec delta rather than pretending continuous synchronization.
- Active request artifacts, prototype mocks, SESSION notes, and unconfirmed AFK drafts are temporary while their owner is active.
- Temporary does not mean disposable without review. On closure, promote the durable outcome to code/tests, a DEC, evidence record, memory, capability index, or archived historical record.

## Stale-safe history

| Artifact | Why it stays | Compaction |
|---|---|---|
| Implemented/rejected/abandoned detailed SPC | Records the contract and evidence used to reach code or reject a path | Move to `archive/specs/`; retain ID, prior status, dates, reason, and verification reference |
| Resolved QUE | Records the ambiguity, options, answer, actor, and provenance | Move to `archive/questions/`; create/link DEC only for a genuine choice |
| DEC and future FND | Durable rationale/takeaway | Keep index summary live; after scale threshold, bodies may move to cold archive through a previewed compaction while IDs remain indexed |
| Completed RES/SPK evidence | Explains what was tested and learned | Keep as evidence; index/compact when scale requires, never treat as current vendor truth without re-checking |
| Shipped roadmap prose | Historical release intent/outcome | `roadmaps/vX.Y.Z.md`; `roadmaps/index.md` remains the only live entry point |
| Archived requests/sessions/fixes/memories | Operational history | Outside normal reads; retrieve explicitly by ID/signature |

There is no root `ROADMAP.md` or `ROADMAP_ARCHIVE.md` in the current contract. `roadmaps/index.md` owns the live ledger and index; `spectacular roadmap migrate` moves shipped prose into `roadmaps/vX.Y.Z.md` while keeping recent history bounded.

Decision bodies remain individually addressable and their frontmatter stays strong (`id`, status, origin/provenance, evidence, tags, dates, supersession). Only the live index compacts as it grows:

- newest 50 decisions: one unchanged line per decision;
- preceding 50: blocks of 10 (`D51–60`, etc.);
- older decisions: blocks of 50 (`D1–50`, etc.).

Thus D60 introduces `D1–10`; D70 adds `D11–20`; at D150 the index shows `D1–50` as one old block, D51–100 in ten-entry blocks, and D101–150 individually. Block summaries retain ID ranges and topic tags so agents can query frontmatter and open only the relevant bodies. Compaction never merges, renumbers, or deletes the underlying decision records.

Before using an archived spec, decision, finding, research result, or human-facing doc to change code, check it against current code/tests and current vendor documentation when external behavior matters. Historical records are evidence, not present-tense claims.

## Throwaway and garbage collection

Throwaway candidates include spike/fork branches, prototype UI/CLI sandbox code, scratch fixtures, generated caches, and abandoned local experiment output. Delete only after all applicable preconditions hold:

1. Record the question/hypothesis and outcome (`supported|refuted|inconclusive` where applicable).
2. Preserve the minimum evidence needed to justify the outcome.
3. Link any resulting QUE, DEC, SPC, request, or reusable finding.
4. Confirm no production code, active dependency, secret-recovery need, or unique user data lives only there.
5. For AFK Git branches, create and verify the hidden `refs/spectacular/archive/*` recovery ref, record tip SHA + restore command, then delete the local branch. Remote deletion stays HITL.

Hidden archive refs are cold recovery, not active context: normal branch listings, status briefings, and agent reads do not load them. A future explicit cleanup may prune old recovery refs only under a separately approved retention policy.

## Terminal transitions

- `QUE open|deferred → archived (archived_from: resolved)` after answer provenance is recorded. A factual answer does not manufacture a DEC.
- `SPC draft|unconfirmed|approved → archived` when rejected/abandoned, with reason.
- `SPC approved → implemented → archived` after verified integration/merge, retaining `implemented_at` and `verified_against`.
- `SPC implemented → superseded|deprecated → archived` remains valid when replacement/removal history is useful before compaction.
- Durable Markdown is archived, never hard-deleted by routine lifecycle commands.

Archive placement and `archived_from` preserve the terminal fact while keeping active collections pristine. Canonical ID lookup may resolve archived dependencies, but normal wayfinding/status iteration never loads archive bodies.

**Related:** [[lifecycle-contract]], [[discovery-protocol]], [[question-rules]], [[spec-lifecycle]], [[roadmap-rules]], [[archive]], and [[afk-git-hygiene]].
