# Review and complete

## Assess each claim

Check every completion claim against its frozen pass boundary. For each one:

- does the Evidence meet the stated proof requirement
- is that Evidence attributable to a method, not just asserted
- is there contrary evidence
- is it fresh enough to still be true
- has the required review happened

When something fails, repair with the narrowest action justified by a new
hypothesis, new Evidence, or a materially smaller fix. Every repair consumes the
repair budget.

## Persist Evidence when it is earned

Write an Evidence file when the claim:

- enters assessment as a material claim
- depends on an external or provider observation
- is disputed
- must survive closure

Routine local progress does not need an Evidence file for every passing check.

## Create a review record only when required

Earned only when the Mission requires clustered or independent review.

```bash
spectacular review record <mission-ref> <review.md|-> --json
```

The record holds:

- reviewed tree and fingerprint
- reviewer identity, and their relation to the operator
- claim-by-claim verdicts
- findings and limitations
- time

**Independent review needs a reviewer who did not implement the reviewed scope.**
A fresh agent is not automatically an independent one. Freshness is not
independence.

## Complete in one flow

1. Run focused checks, then one full tree-bound gate — `bash test/verify.sh all`.
2. Verify every frozen claim, and every required Evidence or review:
   `spectacular mission check <ref> --json`.
3. Confirm the applicable product and specification edits are in the same worktree.
4. Present one owner gate, and stop. Do not proceed on your own reading of intent.
5. After attributable owner confirmation, close it:

   ```bash
   spectacular objective finish <ref>/<objective> --json
   spectacular mission complete <ref> --by <owner> --json
   ```

   `--by` is what makes the confirmation attributable — it names the person who
   accepted. Never pass a name the owner did not actually give you.

There is no separate Contract reconciliation command. Do not ask for one.

## Refuse completion

Refuse when any of these is true:

- the frozen criteria changed
- proof is stale, missing, or conflicting
- required independent review is incomplete
- a dependency or Gap still blocks
- the repair budget is exhausted
- the baseline drifted
- owner confirmation is missing

Completion and archival are different things. Archival proves nothing.
