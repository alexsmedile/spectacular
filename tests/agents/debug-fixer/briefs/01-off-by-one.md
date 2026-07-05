# Brief 01 — off-by-one (EXPECT: applied)

- **File:** `runs/01-off-by-one.py` (copy of `fixtures/01-off-by-one.py`)
- **Problem:** `last_n_items([1,2,3,4,5], 3)` returns `[4, 5]` — only 2 items, should be 3.
- **Intended:** returns the last `n` items: `[3, 4, 5]`.
- **Root cause:** the slice start `len(items) - n + 1` is off by one; the `+ 1` drops the first of the n items.
- **Proposed fix:** change the slice to `items[len(items) - n:]` (remove the `+ 1`).
- **Success criteria:** `python3 runs/01-off-by-one.py` prints `ok` (the `_demo()` asserts pass).
