---
doc-id: fixes
mode: index
location: .spectacular/fixes/
entries-dir: .spectacular/fixes/
scope: project-wide
template: templates/fixes/entry.md
snapshot-on-edit: false
summary: "Self-learning corpus of fully-resolved-and-verified bug fixes. Each F<N>.md records problem → intended → root cause → fix → success criteria → verified-by, plus a searchable signature so a future agent (or project) can recognise a recurring bug and reuse the fix. Not a symptom tracker."
status: active
---

# Fixes Rules

Soft-folder database of **verified bug fixes**. No top-level index file; folder listing is canonical. Each fix is one `F<N>.md` file.

**The rule that gates writing:** a fix entry is logged **only when the bug is fully resolved and verified** — not when it's spotted, not when a patch is drafted. The whole value is that every entry is *trustworthy and done*. In-progress diagnosis lives in [[audit-rules]], not here.

**Mode: `index`** — entry files at `entries-dir`, no regenerated index. Always a folder of files; no flat single-file mode (avoids `decisions/`-style dual-mode complexity).

**Mutator verbs (CLI, since v1.25.0):**
- `spectacular fix new "<title>"` — log `F<N>.md` (auto-numbered). Flags: `--severity`, `--bug/--cause/--fix <text>` (prefill slots), `--verified-by <text>`, `--from-audit <A<N>>` (validated — the audit must exist), `--dry-run`.
  - **The verified gate is a warning, not a block:** omitting `--verified-by` still creates the entry but sets `verified: null` and prints `⚠ no --verified-by — entry marked unverified`. Fill the **Verified by** slot before trusting it. (Soft gate chosen to allow scaffold-then-fill, matching `idea new`.)
- `spectacular fix list [--since <Nd>]` — tabular list; `unverified` shown in the VERIFIED column for null entries.

**Skill verbs (generic engine, when invoked):**
- `grill` → flesh out a thin entry (Bug / Root cause / Fix / Verified by).
- `refine` → sharpen vague slots.
- `review` → validate frontmatter; flag any entry missing `verified:`, **Verified by**, or a **Signature** (the retrieval field). An entry with `verified: null` is a draft, not a fix.

**Snapshot-on-edit: false** — a fix entry is written once, when settled; it isn't a versioned canonical doc.

**Entry frontmatter (required shape):**

```yaml
---
type: fix
opened: YYYY-MM-DD     # when the bug surfaced
verified: YYYY-MM-DD   # when the fix was proven — REQUIRED (a fix without this isn't done); null = unverified
severity: low | medium | high
from_audit: null       # A<N> if it graduated from an audit
signature: "..."       # symptoms/keywords for future retrieval (see below)
related: []
---
```

## The entry skeleton — problem → intended → fix → criteria

A fix entry follows the bug-fixing skeleton, so every entry answers the four questions that matter when you (or another agent) revisit it:

| Section | Question it answers |
|---|---|
| **Problem** | What was observed to be wrong? |
| **Intended behavior** | What *should* happen instead — the contract the bug violated? |
| **Root cause** | Why did it happen — the actual mechanism? |
| **Fix** | What was changed (definitive — this is the applied fix, not a proposal)? |
| **Success criteria** | What observable bar must the fix clear? (the *criterion*) |
| **Verified by** | The concrete check proving it clears the bar (the *evidence*) — a test path, a repro-now-passing, a described manual walk. "Tested it" is not acceptable. |
| **Signature** | Symptoms/keywords that let a future search recognise this bug again. |

**Criteria vs. Verified-by are different fields on purpose:** Success criteria is the *bar* ("exit 0 with no session open"); Verified by is the *evidence that clears it* ("`decide.test.sh` scenario_flat, green"). Don't collapse them.

## Signature — the self-learning field

`fixes/` is not just a record; it's a **retrieval corpus**. Before diagnosing a new bug, the skill greps existing fixes for a matching `signature:` (see [[fixes-loop]] / the skill's self-learning loop). A good signature names the *class* of the bug in searchable terms — the symptom string a future agent would type, plus the underlying pattern:

> `signature: "function returns exit 1 despite success; last statement is a trailing [[ ]] && short-circuit"`

Write it for the reader who hasn't seen this code — another agent, another project, future-you. This is what turns "we fixed it once" into "we know how to fix it." A per-project corpus today; the design leaves room for a shared `~/.spectacular/fixes/` pool or `spectacular fix export/import` across projects (not built yet).

## What this is **not**

- **Not an audit.** [[audit-rules]] is diagnosis in progress; `fixes/` is only the *verified outcome*. An audit graduates here via `from_audit: A<N>`.
- **Not a decision.** [[decisions-rules]] is *why we chose*; a fix is *what broke and how we resolved it*.
- **Not memory.** [[memory-rules]] holds durable standing facts/prefs; a fix is a dated event.
- **Not a changelog.** CHANGELOG.md is user-facing release notes; `fixes/` is the internal, root-cause-level record with its proof-of-fix.

**Related:** [[audit-rules]] (where a fix can originate), [[decisions-rules]], [[memory-rules]], [[new-request]], [[doc-index]].
