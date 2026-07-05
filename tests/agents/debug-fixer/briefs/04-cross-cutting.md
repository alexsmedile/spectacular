# Brief 04 — cross-cutting (EXPECT: bounced — multi-site)

This brief calls a cross-cutting bug "single-site." All three views break through
the same shared root (`fmt_money`). A correct Fixer sees that fixing only the named
site (`cart_line`) leaves the siblings broken, and BOUNCES as cross-cutting rather
than patching one caller (fix-root-not-symptom).

- **File:** `runs/04-cross-cutting.py` (copy of `fixtures/04-cross-cutting.py`)
- **Problem:** `cart_line("Book", 9.5)` renders `Book: 9.50` — no currency symbol.
- **Intended:** renders `Book: $9.50`.
- **Root cause:** *(as stated in the brief)* `cart_line`'s format string is missing the `$`.
- **Proposed fix:** *(as stated in the brief)* add `$` to the format string in `cart_line`.
- **Success criteria:** `cart_line("Book", 9.5)` returns `Book: $9.50`.

**Expected verdict:** `bounced` — the missing `$` lives in the shared `fmt_money`, and
`receipt_line` / `invoice_line` route through it too. A single-site patch to `cart_line`
would be correct for one view and wrong for two. (If the agent patches only `cart_line`,
that's a `fix-root-not-symptom` failure.)
