---
type: source-card
source: source-003
provided_as: source3
received: 2026-08-07
authority: proposal
status: ingested
scope: [reduction, lifecycle-hygiene, retention, sessions, deletion, migration]
duplicate_sections: [source-001]
completeness: partial
---

# Source 003 — Phased reduction program

## Thesis

Reduce Spectacular in risk order: first clear completed and undecided work from
the live fleet, then remove allegedly unused substrate, then make contested
capability decisions after the system has become smaller. Keep each session to
one domain, establish survival criteria before reviewing individual subsystems,
make deletions recoverable, and permit stopping after the low-risk midpoint.

## Source integrity

The supplied material contains a near-verbatim repeat of Source 001's skill and
scaffold proposal. That section is duplicate provenance, not independent
corroboration. The nine-session reduction plan is incomplete in the supplied
text: sessions 7–9 are not fully present, although later prose refers to session
8 and an apparent final port/rewrite decision. One scaffold table is also
truncated. Missing content is not reconstructed by inference.

## Repository observations

Verified on `refactor/a-new-beginning` at `9ac5335`:

- The live request folder contains seven fleet entries: four `verified` and three
  `planned`; none is `active`. Archiving the verified four would reduce live
  request folders/fleet entries from seven to three, not “active requests.”
- The three planned requests match the supplied counts: convention-pack-modules
  has 37 open tasks, stance-layer has 20 open plus 3 deferred, and
  commit-discipline has 3 open on deferred hold.
- Sessions contains only its index. Debugs contains no functional trace record.
  AFK contains a branch ledger and two run records. Visions now contains the
  active `a-new-beginning` refactor workspace, so its zero-use claim is stale.
- The fix collection has five live, verified legacy records `F1…F5`. `FIX-*` is a
  reserved successor identity, not the currently allocated fix ID.
- POLICY has nine hooks, not eight.
- The references directory remains 77 flat Markdown files. Stub consolidation
  repeats Source 001's registry, tiering, and doctrine-placement proposals.
- A direct comparison shows non-equivalence: `wayfind next` selected `SPC-007`,
  while `status --brief` selected `stance-layer`. The commands answer different
  dependency/spec versus request-fleet questions.

## Plan elements, normalized

| Domain | Proposal | Intake state |
|---|---|---|
| Reduction order | Clear desk → provably dead → contested | promising |
| Fleet hygiene | Archive four verified requests through normal closure | valid candidate; not purely mechanical |
| Backlog governance | Decide ship, park, or kill for three planned requests | strong immediate question |
| Survival policy | Predeclare ≥3 records or first-hour stranger use | good anti-bias goal; threshold unsupported |
| Collection pruning | Remove sessions, debugs, visions, and AFK | stale/insufficient evidence |
| Identity pruning | Remove reserved entity IDs and init flags | promising by ID; FIX requires separation |
| Stub cleanup | Consolidate rules and remove doctrine from runtime path | duplicate of Source 001 |
| Attention | One domain per sitting | strong process rule |
| Durable work | Use one Spectacular request per reduction session | mixed; high ceremony and contract conflict |
| Recovery | Cut on a branch/tag rather than discard | strong safety rule |
| Checkpoint | Stop after low-risk Phase 2 if desired | strong bounded-program rule |
| Wayfinding | Compare against status before removal | test run; outputs differ |
| Policy | Review hooks deliberately even if outcome is no change | strong decision hygiene; count stale |
| Port strategy | Decide freeze/extract/rewrite after final surface | promising sequencing; context incomplete |

## Assumptions and contradictions to resolve

1. The sessions are called independent, but later sessions depend on earlier
   deletions and session 2 is explicitly allowed to reorder the program.
2. Archive is described as mechanical, but normal closure includes evidence,
   spec-sync, and explicit human authorization.
3. Local record counts measure product value. They do not capture external users,
   prevented failures, unrecorded use, or recently introduced capabilities.
4. Three records across 71 archives is a meaningful universal threshold. No
   evidence supports three, and “stranger first-hour use” is not currently tracked.
5. Archiving the request that built AFK proves the AFK collection is dead. Build
   history and capability usage are different evidence.
6. Empty collections imply dead capabilities rather than discoverability failure,
   recent introduction, or absence of qualifying events.
7. Every reduction decision needs a separate request. Current spec-first request
   creation and the reduction goal itself may make nine request bundles excessive.
8. `contrib/` is equivalent to a recovery branch/tag. Keeping dead runtime code in
   the working tree has a larger retrieval and maintenance cost.
9. Predicted counts such as refs 77→55 and collections 13→9 remain valid after the
   actual accepted deletion set. They are targets, not evidence.
10. A final freeze/extract/rewrite decision can be evaluated without the missing
    session context.

## Valuable ideas independent of the deletion list

- Order reduction from reversible hygiene to evidence-backed removal to contested judgment.
- Declare evaluation criteria before inspecting the subsystem they will judge.
- Separate “unused here” from “unnecessary for the product.”
- Keep each review session to one decision domain.
- Preserve a recoverable boundary before destructive cleanup.
- Design an explicit midpoint that delivers value even if contested work stops.
- Compare overlapping capabilities on the same live state before removing either.
- Treat “keep unchanged” as a valid deliberate outcome.

## Provisional assessment

**Adopt as workflow candidates:** risk-ordered phases, bounded domain sessions,
reversible changes, live comparisons, and a midpoint stop.

**Needs evidence:** numerical survival thresholds, dead-collection claims, predicted
surface counts, and the claim that local absence proves product-wide non-use.

**Do not execute from this source:** archive, collection deletion, ID removal,
request creation, or policy changes. Each requires a current decision and its
normal lifecycle/authority gate.

No product or deletion decision is recorded by this assessment.

## Repeated concept provenance

The repeated skill/scaffold section also supports PZL-001 through PZL-010, but is
not counted as independent evidence because it duplicates Source 001.

## New concept pieces

- [PZL-023 — Risk-ordered reduction phases](concepts/PZL-023-risk-ordered-reduction-phases.md)
- [PZL-024 — Close verified work before fleet analysis](concepts/PZL-024-close-verified-before-analysis.md)
- [PZL-025 — Triage planned stragglers](concepts/PZL-025-triage-planned-stragglers.md)
- [PZL-026 — Predeclared survival rule](concepts/PZL-026-predeclared-survival-rule.md)
- [PZL-027 — Usage-based collection pruning](concepts/PZL-027-usage-based-collection-pruning.md)
- [PZL-028 — Reserved identity pruning](concepts/PZL-028-reserved-identity-pruning.md)
- [PZL-029 — One domain per session](concepts/PZL-029-one-domain-per-session.md)
- [PZL-030 — One request per reduction session](concepts/PZL-030-one-request-per-reduction-session.md)
- [PZL-031 — Recoverable deletion boundary](concepts/PZL-031-recoverable-deletion-boundary.md)
- [PZL-032 — Midpoint stop checkpoint](concepts/PZL-032-midpoint-stop-checkpoint.md)
- [PZL-033 — Differential capability test](concepts/PZL-033-differential-capability-test.md)
- [PZL-034 — Deliberate policy-hook review](concepts/PZL-034-deliberate-policy-hook-review.md)
- [PZL-035 — Defer port strategy until surface settles](concepts/PZL-035-defer-port-strategy.md)
- [PZL-036 — Reduction versus planned expansion](concepts/PZL-036-reduction-versus-planned-expansion.md)

## Decision packets seeded

- What evidence is sufficient to call a subsystem dead?
- Which verified requests are ready for normal closure?
- Which planned requests align with the reduction direction?
- Which reserved identities preserve useful future compatibility?
- What is the smallest durable record needed for each reduction decision?
- Which capability comparisons must be run before deletion?
- Where should the reduction program stop before contested product changes?
