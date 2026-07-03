---
doc-id: fixes
mode: index
location: .spectacular/fixes/
entries-dir: .spectacular/fixes/
scope: project-wide
template: templates/fixes/entry.md
snapshot-on-edit: false
summary: "Log of fully-resolved-and-verified bug fixes. Not a symptom tracker — the settled record of what broke, why, what changed, and the check that proves it holds. Each entry F<N>.md."
status: active
---

# Fixes Rules

Soft-folder database of **verified bug fixes**. No top-level index file; folder listing is canonical. Each fix is one `F<N>.md` file.

**The rule that gates writing:** a fix entry is logged **only when the bug is fully resolved and verified** — not when it's spotted, not when a patch is drafted. The whole value is that every entry is *trustworthy and done*. In-progress diagnosis lives in [[audit-rules]], not here.

**Mode: `index`** — entry files at `entries-dir`, no regenerated index. Always a folder of files; no flat single-file mode (avoids `decisions/`-style dual-mode complexity).

**Scaffold-only in v1.25.0** — no CLI mutator verb yet. Hand-write entries from `templates/fixes/entry.md`; `<N>` is the next free `F<N>`. CLI verbs (`spectacular fix new|list|show`) are a deferred follow-up.

**Verbs (generic engine, when invoked):**
- `grill` → flesh out a thin entry (Bug / Root cause / Fix / Verified by).
- `refine` → sharpen vague slots.
- `review` → validate frontmatter; flag any entry missing `verified:` or **Verified by** (the load-bearing field).

**Snapshot-on-edit: false** — a fix entry is written once, when settled; it isn't a versioned canonical doc.

**Entry frontmatter (required shape):**

```yaml
---
type: fix
opened: YYYY-MM-DD     # when the bug surfaced
verified: YYYY-MM-DD   # when the fix was proven — REQUIRED (a fix without this isn't done)
severity: low | medium | high
from_audit: null       # A<N> if it graduated from an audit
related: []
---
```

**Required body sections:** Bug, Root cause, Fix, **Verified by**.

**Verified by** is the field that makes this collection worth keeping — it must name a *concrete* check: a test path (`tests/cli/decide.test.sh`), a repro that now passes, or a described manual walk. "Tested it" is not acceptable; name what would fail if the fix regressed.

## What this is **not**

- **Not an audit.** [[audit-rules]] is diagnosis in progress; `fixes/` is only the *verified outcome*. An audit graduates here via `from_audit: A<N>`.
- **Not a decision.** [[decisions-rules]] is *why we chose*; a fix is *what broke and how we resolved it*.
- **Not memory.** [[memory-rules]] holds durable standing facts/prefs; a fix is a dated event.
- **Not a changelog.** CHANGELOG.md is user-facing release notes; `fixes/` is the internal, root-cause-level record with its proof-of-fix.

**Related:** [[audit-rules]] (where a fix can originate), [[decisions-rules]], [[memory-rules]], [[new-request]], [[doc-index]].
