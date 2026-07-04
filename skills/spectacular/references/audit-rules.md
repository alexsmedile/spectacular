---
doc-id: audit
mode: index
location: .spectacular/audit/
entries-dir: .spectacular/audit/
scope: project-wide
template: templates/audit/entry.md
snapshot-on-edit: false
summary: "Bug-investigation scratchpad — diagnose a symptom and find its root cause BEFORE deciding whether it needs a request. Soft-folder DB; each entry A<N>.md exits via a disposition."
status: active
---

# Audit Rules

Soft-folder database of **bug investigations** — the diagnostic step between "something's off" and "here's the plan". No top-level index file; folder listing is canonical. Each investigation is one `A<N>.md` file.

**Mode: `index`** — entry files at `entries-dir`, no regenerated index. Always a folder of files; there is no flat single-file mode (avoids the two-code-path complexity of `decisions/`).

**Mutator verbs (CLI, since v1.25.0):**
- `spectacular audit new "<title>"` — scaffold `A<N>.md` (auto-numbered), `status: open`. Flags: `--severity low|medium|high` (default medium), `--symptom <text>`, `--dry-run`.
- `spectacular audit list [--status open|resolved|folded|all]` — tabular list.
- `spectacular audit resolve <A<N>> --disposition "<text>"` — close: `status: resolved`, set disposition. `--into-fix` also scaffolds a `fixes/F<N>` entry (copies the audit's root cause, sets `from_audit: A<N>`, and owns the disposition string `became fix F<N>`); pass `--verified-by <text>` alongside to seed the fix's proof slot.

**Skill verbs (generic engine, when invoked):**
- `grill` → flesh out a thin investigation (Symptom / Investigation / Root cause).
- `refine` → sharpen vague slots.
- `review` → validate frontmatter across entries; flag stale `open` + missing `disposition` on `resolved`.

**Snapshot-on-edit: false** — audits are scratchpad, not versioned canonical docs.

**Entry frontmatter (required shape):**

```yaml
---
type: audit
status: open | resolved | folded
severity: low | medium | high
opened: YYYY-MM-DD
updated: YYYY-MM-DD
disposition: null   # set on close
related: []
---
```

## The entry skeleton — problem → intended → root cause → proposed fix → criteria

An audit walks the same bug-fixing skeleton as a fix, but the later slots are *unsettled* — it's diagnosis in progress:

| Section | Question it answers |
|---|---|
| **Problem** | What was observed to be wrong? |
| **Intended behavior** | What *should* happen instead — the contract the bug violates? |
| **Investigation** | What was checked, what was ruled out? |
| **Root cause** | The actual cause, once found (or "not yet found"). |
| **Proposed fix** | The *suggested* fix — still a proposal at audit stage, not yet applied or verified. |
| **Success criteria** | How we'll know it's fixed — the observable bar. |
| **Disposition** | The exit (see below). |

**Why intended-behavior + success-criteria live in the audit:** capturing them *before* fixing is what keeps the fix honest — you decide the bar before you're invested in a patch. When the audit graduates via `resolve --into-fix`, **every matching slot is copied forward** into the fix entry (Problem, Intended, Root cause, Proposed fix→Fix, Success criteria), so the investigation is never re-typed.

## Status lifecycle

```
open ──(root cause found, action decided)──► resolved
  │
  └──(rolled into a request instead)──────► folded
```

- **`open`** — under investigation. Default at creation.
- **`resolved`** — cause found + `disposition` set (one-line fix · won't-fix · became fix `F<N>`).
- **`folded`** — became (or joined) a `requests/<slug>` — the request now owns it.

## Disposition — the exit

Every closed audit carries a `disposition`. It is the single most important field: it says what *happened* to the investigation.

| Disposition | Meaning |
|---|---|
| `requests/<slug>` | Folded into a planned request — that request owns the fix. |
| one-line fix | Fixed on the spot; too small for a request. |
| won't-fix | Deliberately not fixing (state why). |
| became fix `F<N>` | Fixed + verified → logged as a [[fixes-rules]] entry. |

## What this is **not**

- **Not a fix.** A fix is *verified and done* — see [[fixes-rules]]. An audit is *diagnosis in progress*; it may never become a fix.
- **Not an idea.** [[idea-rules]] is feature exploration; audit is bug diagnosis.
- **Not a request.** Audits carry no PLAN/TASKS/lifecycle. A `folded` audit produces or joins a request.
- **Not a decision.** [[decisions-rules]] records *why we chose* between options; an audit records *what's wrong and why*.

**Related:** [[fixes-rules]] (the verified-fix log an audit can graduate into), [[idea-rules]] (feature vs bug), [[new-request]] (fold destination), [[decisions-rules]], [[doc-index]].
