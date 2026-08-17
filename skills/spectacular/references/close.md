# Review and complete

## Assess each claim

Check every completion claim against its frozen pass boundary using the FROST
framework (Frozen fit, Risk, Operability, System integrity, Truth of proof)
defined in [audit.md](audit.md). For each one:

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

The review input is a Markdown file with `ReviewDraft` frontmatter:

```yaml
---
type: ReviewDraft
title: Independent review of M<n>
status: passed
reviewed:
    commit: <40-char-git-commit>
    tree: <40-char-git-tree>
    activation_fingerprint: <sha256-mission-activation-fingerprint>
reviewer:
    actor: <identity-different-from-operator>
    operator: <operator-name>
    relation_to_operator: independent
    implemented_reviewed_scope: false
    independence_basis: <attributable statement of independence>
    evidence:
        - <attributable command or observation>
claims:
    - claim: <exact-claim-name-from-mission>
      verdict: pass
findings: []
limitations: []
---

# Review body
```

**Independent review needs a reviewer who did not implement the reviewed scope.**
A fresh agent is not automatically an independent one. Freshness is not
independence.

What earns it: the reviewer did not author the work, inspects primary evidence
rather than the executor's summary of it, and where the consequence is high uses
a different method, a qualified human, a specialist tool, or an independent
observation. The reviewer reports verdicts and findings; it never declares the
Mission complete.

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

There is no Contract reconciliation ritual at completion. Completion checks the
Gaps the Mission declared it would close and refuses while any is still open; it
never writes the Contract itself.

## Close a declared Gap

A Mission that declared `resolves_gaps:` closes each one when the work resolving
it lands, not at completion:

```bash
spectacular contract amend <contract-ref> --gap <gap-ref> --by <owner> --dry-run
spectacular contract amend <contract-ref> --gap <gap-ref> --by <owner>
```

`--dry-run` prints the resolution text, both fingerprints, and every Mission that
would be re-pointed, and writes nothing. The text is the wording frozen in the
Mission's declaration, so the owner approved it at activation.

The amendment reaches the `gaps:` block and editorial fields only. A semantic
change is a new `contract_version:`, not an amendment. It refuses while any bound
Mission is live, so amend between Missions or stop the live one first.

A Gap is never closed by deleting it.

## Refuse completion

Refuse when any of these is true:

- the frozen criteria changed
- proof is stale, missing, or conflicting
- required independent review is incomplete
- a dependency or Gap still blocks, or a declared resolved Gap is still open
- the repair budget is exhausted
- the baseline drifted
- owner confirmation is missing

Completion and archival are different things. Archival proves nothing.
