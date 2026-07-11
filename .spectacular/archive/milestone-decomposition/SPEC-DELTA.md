### ADDED
- specs/index.md :: Agent fleet arc summary — the **size-and-decompose gate** (build-workflow Step 1.5 + `decompose-large-milestone` @Implementation policy): a multi-phase milestone is broken into sequential nested-`- [ ]` sub-step checkpoints built/dispatched one at a time, so a fat milestone never runs as one opaque block (visibility from sub-step boundaries, not a live trace on a running agent).

Note: the spec edit was made during this request's M4 (specs/index.md v1.10 → v1.11), so the
delta is already applied to current truth — this file records what shipped. No MODIFIED/REMOVED.
The change is doc-only (skill reference docs + POLICY.md + specs); no code or CLI surface changed.
