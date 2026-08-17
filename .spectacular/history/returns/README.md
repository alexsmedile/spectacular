# Archived handoff returns

Pre-v2 handoff returns, kept as historical evidence. Nothing here is a v2 record
type, nothing here is authoritative, and nothing here is loaded as context.

They live under `history/`, which `internal/discovery/discovery.go` skips
entirely, so their v1 frontmatter does not have to satisfy any current schema.
Putting one of these files anywhere the scan reaches refuses every workspace
read with `invalid_type field type: unknown record noun` — the schema names
differ from this product's v2 despite the version number in them.

| File | Date | Origin | Disposition |
|---|---|---|---|
| `H16-sdlc-coherence-adversarial-review.md` | 2026-08-09 | Codex delegated review at `5b5a738` | Never accepted; `central_disposition: pending` |

## H16

An adversarial SDLC-coherence review of the v1 foundation. It audited four
scenarios and found no missing foundation, then captured an owner-approved
preparation loop and a two-role specialist slate.

Its return was never disposed of. The branch it audited was superseded by the v2
clean break, and the central orchestration its `next_action` names no longer
exists. It survived only as an untracked file in an abandoned worktree and was
found during worktree cleanup on 2026-08-17.

`P10` re-derives the two capabilities that still apply — a design-sufficiency
verdict and a slice-quality verdict — against the v2 surface, and deliberately
drops the rest: the eleven v1 contracts, the S10–S12 session plan, the
specialist roles as products, and the return schema. Read `P10` for the live
proposal; this file is kept only so the dropped material is recoverable if that
judgment turns out to be wrong.
