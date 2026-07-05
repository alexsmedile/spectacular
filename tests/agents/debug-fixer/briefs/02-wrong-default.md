# Brief 02 — wrong default (EXPECT: applied)

- **File:** `runs/02-wrong-default.js` (copy of `fixtures/02-wrong-default.js`)
- **Problem:** `paginate(items, 1)` returns an empty array — the first page is empty.
- **Intended:** default page size is 20; `paginate(items, 1)` returns items 0–19.
- **Root cause:** the `perPage` default is `0`, so `slice(0, 0)` returns nothing.
- **Proposed fix:** change the default parameter from `perPage = 0` to `perPage = 20`.
- **Success criteria:** `node runs/02-wrong-default.js` prints `ok` (the `_demo()` checks pass).
