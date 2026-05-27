---
status: draft
updated: 2026-05-24
related:
  - PLAN.md
  - discovery.md
  - TASKS.md
---

# Spec — substrate-clarity changeset

Concrete file-by-file changes needed to land the locked model (discovery.md). Organized by milestone so each milestone is a self-contained PR-sized chunk.

**Scale:** 1 rename + 13 frontmatter additions + 6 new rules files + ~78 "engine" cleanups + CLI redirect + matrix doc.

---

## M3 — Registry demotion

### Rename

| From | To |
|---|---|
| `skills/spectacular/references/doc-registry.md` | `skills/spectacular/references/doc-index.md` |

**Snapshot first:** `doc-registry@v1.md` (via `spectacular snapshot`).

### Content rewrite — `doc-index.md`

**Remove:**
- "Adding a new doc type = entry + template, no other code changes" claim (line 5)
- "Core principle" section as dispatch contract (lines 7-18)
- Full YAML schema block (lines 22-35)
- Field-semantics section (lines 37-66)
- Per-doc YAML entries (lines 70-222)
- "How the engine uses this" section (lines 232-240)
- "Adding a new doc type" prescriptive checklist (lines 254-262)
- "Project-local registry overrides" v2 forward-decl (lines 264-266)

**Keep / reshape as human catalog:**

```markdown
# Doc Index — what Spectacular knows how to write

Human catalog of every document type in a Spectacular workspace.

Dispatch lives in each doc's rules file (`references/<doc-id>-rules.md`).
This index is for browsing — not parsed for behavior.

## Project-wide canonical docs

| Doc | Mode | Location | Rules |
|---|---|---|---|
| `prd` | grill | `.spectacular/PRD.md` | [prd-rules](prd-rules.md) |
| `principles` | stub | `.spectacular/PRINCIPLES.md` | [principles-rules](principles-rules.md) |
| `architecture` | stub | `.spectacular/ARCHITECTURE.md` | [architecture-rules](architecture-rules.md) |
| `spec` | stub | `.spectacular/SPEC.md` | [spec-rules](spec-rules.md) |
| `roadmap` | grill-each | `.spectacular/ROADMAP.md` | [roadmap-rules](roadmap-rules.md) |
| `stack` | stub | `.spectacular/STACK.md` | [stack-rules](stack-rules.md) |
| `agents` | stub | `.spectacular/AGENTS.md` | [agents-rules](agents-rules.md) |
| `decisions` | append | `.spectacular/DECISIONS.md` | [decisions-rules](decisions-rules.md) |
| `personas` | grill-each | `.spectacular/PERSONAS.md` | [personas-rules](personas-rules.md) |

## Per-request docs

| Doc | Mode | Location | Rules |
|---|---|---|---|
| `plan` | grill | `.spectacular/requests/<slug>/PLAN.md` | [plan-rules](plan-rules.md) |
| `tasks` | stub | `.spectacular/requests/<slug>/TASKS.md` | [tasks-rules](tasks-rules.md) |

## User-scope docs

| Doc | Mode | Location | Rules |
|---|---|---|---|
| `convention-pack` | grill | `~/.spectacular/packs/<name>/pack.md` | [pack-rules](pack-rules.md) |

## Public-facing docs (deprecated v1.2.0)

> See [pageworks](https://github.com/alexsmedile/pageworks) for current public-docs authoring.

| Doc | Mode | Location | Rules | Status |
|---|---|---|---|---|
| `docs-manifest` | stub | `docs/docs.yaml` | [docs-rules](docs-rules.md) | deprecated v1.2.0 |
| `docs-page` | stub | `docs/<section>/<slug>.md` | [docs-rules](docs-rules.md) | deprecated v1.2.0 |

## Skill-internal references

| Doc | Mode | Location | Rules |
|---|---|---|---|
| `migrations-contract` | reference | `skills/spectacular/references/migrations-contract.md` | — |
| `migration` | reference | `skills/spectacular/references/migrations/v<from>-to-v<to>.md` | — |

## Mode taxonomy (reference)

| Mode | Behavior |
|---|---|
| `grill` | Interactive form, default style = wide. Alias for `grill-wide`. |
| `grill-wide` | One broad session, all slots at once. |
| `grill-each` | Per-block walk. Same slots, one block at a time. |
| `grill-loop` | Wide pass, then deep on opens. |
| `append` | Capture one entry, append to file. |
| `stub` | Scaffold + exit. User edits directly. |
| `freeform` | Agent improvises shape (reserved). |
| `reference` | Skill-internal doc, not user-facing. |

## Adding a new doc

1. Create rules file: `references/<doc-id>-rules.md` with frontmatter (see schema in [scaffold-reference](scaffold-reference.md)).
2. Create template at the path declared in `template:` field.
3. Add a row to this index (cosmetic — no behavior depends on it).

CLI auto-discovers via filesystem walk. No edits to `cli/spectacular` needed for new docs in v1.4.0+.
```

### New rules files (6 stub/append docs)

Each is created at `skills/spectacular/references/<doc-id>-rules.md` with frontmatter-only content (body optional, can be one paragraph).

| New file | Mode | Snapshot-on-edit |
|---|---|---|
| `principles-rules.md` | stub | true |
| `architecture-rules.md` | stub | true |
| `stack-rules.md` | stub | true |
| `agents-rules.md` | stub | true |
| `spec-rules.md` | stub | true |
| `decisions-rules.md` | append | false |

**Example body for new stub rules files:**

```markdown
---
doc-id: agents
mode: stub
location: .spectacular/AGENTS.md
scope: project-wide
template: templates/agents/base.md
snapshot-on-edit: true
summary: "Onboarding doc for agents working in .spectacular/"
status: active
---

# AGENTS Rules

Stub doc. Scaffolded once at `init` time from `templates/agents/base.md`. User edits directly thereafter.

**Review behavior:** structural check only — frontmatter present, sections from template not left as placeholders.

**No grill / refine** — those verbs redirect with the standard stub message.
```

### Add frontmatter to 7 existing rules files

Each existing rules file gets frontmatter prepended. The body stays as-is. Snapshot each file before edit.

| Existing file | Mode | Notes |
|---|---|---|
| `prd-rules.md` | grill | `slots: [Vision, Problem, Target users, Deliverable, Goals & success criteria, Non-goals, Constraints, First milestone]`, `kit-support: true` |
| `plan-rules.md` | grill | `slots: [Goal, Constraints, Milestones, Tasks, Dependencies, Validation, Deliverables]` |
| `tasks-rules.md` | stub | `snapshot-on-edit: false` |
| `roadmap-rules.md` | grill-each | `slots: [Status, Phase, Scope (in), Scope (out), Exit criteria, Linked requests]` ← was `reps` |
| `personas-rules.md` | grill-each | `slots: [Who, Wants to, Pain, Stories, Not for]` ← was `reps` |
| `pack-rules.md` | grill | `slots: [Name & scope, Naming, Taxonomy, Root files & README, Gitignore, File placement, Project types]`, `scope: user` |
| `docs-rules.md` | stub | `status: deprecated`. Two doc-ids (docs-manifest + docs-page); needs split or shared file convention — see Open Question below. |

### Skill-side: drop reads of doc-registry path

Files that reference `doc-registry.md` and need link/path updates:

| File | What changes |
|---|---|
| `SKILL.md` | Update references index (line ~210, `doc-registry.md` → `doc-index.md`); update routing description (line 68) |
| `references/grill.md` | Replace "load registry to resolve mode" with "load rules file for `<doc-id>`" |
| `references/refine.md` | Same as grill.md |
| `references/review.md` | Same |
| `references/personas-rules.md` | Update `[[doc-registry]]` wikilink |
| `references/pack-rules.md` | Same |
| `references/roadmap-rules.md` | Same |
| `references/prd-rules.md` | Same |
| `references/plan-rules.md` | Same |
| `references/tasks-rules.md` | Same |
| `references/kits-contract.md` | Same |
| `references/init-workflow.md` | Same |
| `references/doctor.md` | Same |
| `references/doctor-substrate.md` | Same |
| `references/doctor-repair.md` | Same |
| `references/status.md` | Same |
| `references/migrations-contract.md` | Same |
| `references/packs-contract.md` | Same |

### Open Question (M3) — `docs-rules.md` covers two doc-ids

Today `docs-rules.md` is the rules file for both `docs-manifest` and `docs-page`. Options:

1. **Split** into `docs-manifest-rules.md` + `docs-page-rules.md` (consistent with one-doc-one-rules-file pattern)
2. **Keep shared** — frontmatter declares an array: `doc-ids: [docs-manifest, docs-page]` (special-case schema)
3. **Skip** — file is already deprecated; leave as-is until removed in v2.0.0

**Recommendation:** option 3. Deprecation removes this corner before it matters.

---

## M4 — Mode collapse: `reps` → `grill-each`

### Rules file edits

| File | Change |
|---|---|
| `roadmap-rules.md` frontmatter | `mode: grill-each` |
| `personas-rules.md` frontmatter | `mode: grill-each` |
| `roadmap-rules.md` body line 388 | drop `with mode: reps` clause |
| `personas-rules.md` body line 5 | "is an opt-in, `reps`-mode doc" → "is an opt-in, `grill-each`-mode doc" |

### Template edit

| File | Change |
|---|---|
| `templates/roadmap/base.md` line 18 | "Mode: reps (v0.7.1+; renamed from "structured" in v1.3.0)" → "Mode: `grill-each` (v1.4.0+; was `reps` in v1.3.x and `structured` pre-v0.7.1)" |

### `doc-index.md` consistency

Confirm rendered table shows `grill-each` for ROADMAP + PERSONAS (already in the rewrite above).

### Skill engine wiring

`references/grill.md`:
- Add mode-resolution section: declared `mode: grill-X` → walk style X
- Add flag override: `--wide` / `--each` / `--loop` wins over declared mode
- Document `mode: grill` as sugar for `mode: grill-wide`

---

## M5 — Build grill-loop

### Engine spec — append to `references/grill.md`

**Algorithm:**
1. Pass 1 (wide): walk all slots in order, accept short answers (1-2 sentences each)
2. Mark each slot with `[needs-deepening]` if it matches the heuristic
3. Pass 2 (deep): revisit only flagged slots; run grill-wide quality on each

**Vagueness heuristic (v1.4.0 initial):**

A slot is flagged for pass 2 if **any of**:
- Length < 30 chars
- Matches any word from the vague-word list scoped to that slot (already defined in rules files like `prd-rules.md`)
- Contains placeholder strings (`<…>`, `TODO`, `tbd`)
- Slot has explicit gate-check that didn't pass

**Open question (M5):** explicit user-controlled `[needs-deepening]` marker (user adds in text) vs auto-only. Recommend: ship auto-only in v1.4.0; add explicit marker if usage demands.

### Test docs

- Manual: `spectacular prd grill --loop` on a fresh PRD
- Manual: `spectacular roadmap grill --loop` on a fresh ROADMAP

---

## M6 — Agentic/mechanical verb split

### CLI edits — `cli/spectacular`

**Add agentic-verb detection + redirect:**

```bash
# Pseudocode location: after arg parsing, before dispatch
AGENTIC_VERBS="grill refine"

is_agentic_verb() {
  case " $AGENTIC_VERBS " in
    *" $1 "*) return 0 ;;
    *) return 1 ;;
  esac
}

# In dispatch case for doc-id + verb:
if is_agentic_verb "$VERB"; then
  echo ""
  echo "  ⚡ '$VERB' requires an agent."
  echo ""
  echo "  Run inside Claude Code or Codex:"
  echo "      /spectacular $DOC $VERB"
  echo ""
  exit 0
fi
```

**`--help` updates:**
- Add legend: agentic verbs marked `(skill-only)`
- `grill (skill-only)`
- `refine (skill-only)`
- `review` (no mark — partial CLI support)
- All other verbs unchanged

### Docs edits

| File | Change |
|---|---|
| `docs/commands.md` | Add "Agentic vs mechanical verbs" section near top; mark each verb in the reference list |
| `docs/installation.md` | If it says "CLI does everything," correct to "CLI handles mechanical; skill handles agentic" |

---

## M7 — Cleanup pass

### "Engine" word sweep

**Total occurrences:** ~78 across `skills/spectacular/`

**Mechanical replacements** (Decision 7 guide):

| Pattern | Replacement |
|---|---|
| `the generic engine` | `grill / refine / review` or `the skill` |
| `generic engine` | drop "generic" — if "engine" stays in context, sub it with "skill" |
| `the engine reads` | `the verb reads` or `the skill reads` |
| `consumed by the engine` | `consumed by grill/refine/review` (or drop the phrase) |
| `engine rules` | `rules` |
| `engine behavior` | `skill behavior` or `verb behavior` |
| `engine-internal` (used as adjective for reference mode) | `skill-internal` |
| `engine no-ops` | `skill no-ops` |

**Files with concentrated occurrences:**
- `skills/spectacular/SKILL.md` (multiple)
- `references/grill.md`, `refine.md`, `review.md`
- `references/prd-rules.md`, `pack-rules.md`, `personas-rules.md`, `roadmap-rules.md`, `plan-rules.md`, `tasks-rules.md`, `docs-rules.md`
- `references/kits-contract.md`, `packs-contract.md`, `migrations-contract.md`
- `references/init-workflow.md`, `doctor.md`, `doctor-substrate.md`, `doctor-repair.md`

**Approach:** mass `grep` + targeted edits, file by file. Don't blind-`sed` — some "engine" uses might be genuinely about car engines or LLM engines (the trigger model). Review each match.

### Header rename — kill "Overrides" wording where it lingers

Several rules files still have `# X Overrides` headers from before the v1.3.0 rename to "Rules":

| File | Current H1 | New H1 |
|---|---|---|
| `prd-rules.md` | "PRD Overrides — ..." | "PRD Rules — ..." |
| `plan-rules.md` | "PLAN Overrides — ..." | "PLAN Rules — ..." |
| `tasks-rules.md` | "TASKS Overrides — ..." | "TASKS Rules — ..." |
| `roadmap-rules.md` | "ROADMAP Overrides — ..." | "ROADMAP Rules — ..." |
| `pack-rules.md` | "Pack Overrides — ..." | "Pack Rules — ..." |
| `docs-rules.md` | "Docs Overrides — ..." | "Docs Rules — ..." |
| `personas-rules.md` | already "PERSONAS Rules" | (skip) |

### Matrix doc

Land the verb × mode matrix in `doc-index.md` (Mode taxonomy section above) or — if it bloats the index — pull out to `references/verb-mode-matrix.md`.

**Recommendation:** keep in `doc-index.md` for v1.4.0. Pull out only if it exceeds ~30 lines.

### SKILL.md routing table

Update `skills/spectacular/SKILL.md` routing table:
- `doc-registry.md` → `doc-index.md`
- Add `references/verb-mode-matrix.md` line if extracted

### AGENTS.md / onboarding

| File | Change |
|---|---|
| `.spectacular/AGENTS.md` | If it mentions "the generic engine" or `mode: reps` or `doc-registry.md`, sweep |
| `references/onboarding.md` | Same |

### CLAUDE.md (repo-level)

| Line | Change |
|---|---|
| Repo structure table mentioning `doc-registry.md` | → `doc-index.md` |
| Skill routing table reference to `doc-registry.md` | → `doc-index.md` |

---

## M8 — Doctor + tests + ship

### Doctor updates

`cli/spectacular doctor`:
- `frontmatter` area: validate every rules file's frontmatter against the schema (Decision 6)
- `links` area: confirm no broken refs to old `doc-registry.md` path
- New check: every doc-id in `doc-index.md` resolves to an existing rules file

### Tests

Audit `tests/`:
- Any test asserting `doc-registry.md` path → update to `doc-index.md`
- Any test asserting `mode: reps` literal → update to `mode: grill-each`
- Any test asserting "engine" wording → update or drop

### CHANGELOG

`## [1.4.0] — 2026-MM-DD`

**Breaking:**
- `mode: reps` removed from registry. Migrated to `mode: grill-each`. (`spectacular doctor frontmatter --fix` available for any user-side rules files using `reps`.)
- `doc-registry.md` renamed to `doc-index.md`. Old path symlink optional; recommend updating references.

**Added:**
- `grill-wide` / `grill-each` / `grill-loop` mode values; `grill` is sugar for `grill-wide`.
- Rules files for the 6 previously-implicit docs (PRINCIPLES, ARCHITECTURE, STACK, AGENTS, SPEC, DECISIONS).
- Frontmatter schema on every rules file — `doctor frontmatter` validates.
- CLI redirect for agentic verbs (`grill`, `refine`) with friendly message.
- `--wide` / `--each` / `--loop` flags on the `grill` verb.

**Changed:**
- `doc-index.md` reframed as human catalog. Dispatch lives in each rules file's frontmatter.
- "Engine" terminology dropped throughout skill docs. Refer to "skill" or name the verb directly.
- Rules file H1s standardized from "X Overrides" → "X Rules".

**Migration:**
- No user action needed for existing workspaces. Rules files in `.spectacular/` are not user-authored; the skill bundles them.
- Custom packs using `mode: reps` should update to `mode: grill-each`. `doctor packs` will warn.

### Version bump

- Manifests: `bump-manifests.sh 1.4.0` (4 files: plugin.json ×2, marketplace.json, README badge)
- Manual: `cli/spectacular` SPECTACULAR_VERSION
- Manual: `skills/spectacular/SKILL.md` version frontmatter
- Audit: `check-manifests.sh`

### Tag + release

- Commit: `chore: release v1.4.0`
- Tag + push
- `gh release create v1.4.0 --generate-notes`
- User-triggered: `/plugin marketplace update spectacular`

### Archive request

- Snapshot PLAN + TASKS + discovery + spec
- Move to `.spectacular/archive/substrate-clarity-v1.4.0/`
- Update CLAUDE.md Active Requests → Archived list

---

## Risk register

| Risk | Mitigation |
|---|---|
| ~78 "engine" edits introduce inconsistency or break wikilinks | Process file-by-file, not bulk-sed. Doctor links check at end. |
| New rules files for stub docs feel like make-work | Body is optional one-paragraph. Frontmatter-only is fine. Total per-file effort: 2 minutes. |
| `mode: reps` users in custom packs break silently | Doctor packs check + CHANGELOG migration note. |
| CLI redirect annoys users who actually wanted a CLI grill | Friendly tone in message; documents the why (LLM required). |
| `grill-loop` heuristic ships imperfect | Mark as v1.4.0-initial in code comments; refine in v1.4.x as usage informs. |
| `docs-rules.md` covers two doc-ids and the schema assumes one | Open Question M3 — recommend skipping (deprecated anyway). |
| Renaming `doc-registry.md` breaks external bookmarks | Soft-deprecation: leave a stub at old path for one minor pointing to new file. Or accept — internal docs. |

## Sequencing notes

- **M3 → M4 → M7** depend on each other (registry rename comes first, mode rewrite uses new schema, cleanup sweeps both).
- **M5 (grill-loop)** is independent of M3/M4; could ship as a separate PR if needed.
- **M6 (CLI redirect)** is independent of everything else; ships as a tiny standalone PR.
- **M8** rolls up everything.

Suggested PR order if shipping incrementally: M6 → M3 → M4 → M7 → M5 → M8. Each is independently reviewable.
