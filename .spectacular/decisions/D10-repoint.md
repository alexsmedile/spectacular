---
type: Decision
id: 01a0120f-0563-7533-9061-2ef98f9c4780
title: Re-point only the live Mission; a completed Mission keeps the binding it agreed to
created_by: Alex
created: "2026-08-17T23:29:24Z"
updated: "2026-08-17T23:29:24Z"
actor: Alex
actor_role: owner
ref: D10-repoint
question: When an amendment changes a Contract's fingerprint, which bound Missions are re-pointed to the new one?
disposition: repoint-live-only
rationale: >-
    We rewrite the present to make the future work, not to correct a past that did not know about the
    future. A completed Mission is frozen in time: its binding is the historical fact of which Contract
    text it agreed to, and rewriting it writes today's answer over that fact. That destroys the record
    re-pointing was meant to protect, and it is the one thing needed to derive why work was done a
    particular way — probably because an older Contract said so. The stale binding was never a defect:
    `mission check` reports it as a contract-drift notice and the Mission stays `valid=true`, and
    `git log -S <old fingerprint>` recovers the exact Contract text in force. Nothing is lost and the
    freeze point becomes knowable. A live Mission is the opposite case and still re-points, because it
    is working against the Contract now and would otherwise report drift against an amendment it just
    made. The old behavior also touched four Mission files per amendment to achieve this.
alternatives:
    - re-point every bound Mission, as built, keeping every binding pointed at the current Contract text
    - re-point nothing, including the live Mission, accepting that a Mission amending its own Contract reports drift against itself for the rest of its life
authority_basis: >-
    Owner reviewed the re-pointing behavior after M12's amendments rewrote M6, M10, M11, and M12, judged
    the rewriting of completed Missions to be dead weight that discards the freeze point, and authorized
    narrowing re-pointing to live Missions.
authorized_effects:
    - amend.repoint-scope-narrowed
conditions:
    - completed-mission-records-left-byte-identical
    - contract-drift-remains-a-notice-not-a-refusal
expected_fingerprints:
    - sha256:95e351f2854fe11b4183c66b30e1899e3a804ab10884e488e51adaf671e20ec7
scope:
    - v2
targets:
    - Contract:01a00a20-63dd-7670-97f1-9eb8e12adc3a
evidence:
    - Mission:01a010a6-01b0-7320-acc2-5c695bec2843
supersedes: ""
---
# Decision

## What re-pointing is

Closing a Gap through `contract amend` rewrites the Contract, which changes its
fingerprint. Every Mission bound to that Contract carries the old fingerprint in
its `contract:` block. Re-pointing is the amendment's follow-up write that
replaces the old fingerprint with the new one on those Missions.

## What it used to do

Every bound Mission was re-pointed — live, completed, and archived alike. The
reasoning was that a binding should name the Contract as it currently reads, so a
reader is not chasing a fingerprint no file matches.

M12's two amendments made the cost visible: each one rewrote M6, M10, M11, and
M12, four files per amendment, three of them Missions that finished months of work
ago and were never going to be executed again.

## Why that is wrong

A completed Mission is a freeze in time. Its `contract.fingerprint` is not a
pointer that should stay fresh — it is the statement *this is the agreement I was
executed against*. That statement is the only thing that lets a later reader
derive why the work was shaped the way it was, and the answer is usually that an
older Contract said so.

Re-pointing overwrites exactly that. It replaces a true historical statement with
today's answer, and in doing so erases the record it claimed to be maintaining.

The premise that a stale binding is a problem does not hold. Reverting M6 and M11
to their pre-amendment fingerprints and running `mission check`:

```
M6  valid=true schema=mission.v2 checks=22 contract=v2
  notice: contract-drift: Contract:01a00a20… changed after this Mission completed;
          bound at sha256:80336b15…, now sha256:95e351f2…
M11 valid=true schema=mission.v2 checks=22 contract=v2
  notice: contract-drift: … bound at sha256:315f4ad5…, now sha256:95e351f2…
```

`valid=true`. The drift is a **notice, not a refusal** — the mechanism that
already exists for exactly this, and it states the freeze point rather than hiding
it. M6's binding at that point was `80336b15…`, from before M11's own amendments,
so the workspace had already been carrying drift on a completed Mission with no
consequence.

Nothing is lost either. The old Contract text is one command away:

```bash
git log -S sha256:80336b15… -- .spectacular/contracts/CC-<name>.md
```

| | Re-point everything | Re-point live only |
|---|---|---|
| What did M6 agree to? | overwritten | stated, and recoverable from git |
| Contract changed since? | hidden | reported as a drift notice |
| Files written per amendment | 4+ | 1 |

## Why the live Mission is different

A live Mission is not a historical record; it is working against the Contract
right now. Its binding has to track the current text, or a Mission that amends its
own Contract immediately reports drift against an amendment it just made and
carries that for the rest of its life. Re-pointing is correct there for the same
reason it is wrong for a completed Mission: the live Mission's binding is a
statement about the present.

## Consequence

`repointBoundMissions` skips any Mission whose status is not `active`. Completed
and archived Missions are left byte-identical; their drift surfaces as a notice on
`mission check`. The ambiguity refusal on multiple fingerprint occurrences is
unchanged and still applies to the live Mission.

This is a narrowing of an existing write, not a new capability, and needed no
command-surface change.
