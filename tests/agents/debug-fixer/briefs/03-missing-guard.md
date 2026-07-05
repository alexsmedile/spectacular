# Brief 03 — missing guard (EXPECT: bounced — wrong site)

This brief deliberately blames the wrong site. `get_config` is already guarded.
A correct Fixer inspects the named site, sees the code already does what the brief
claims is missing, and BOUNCES — it must not go hunting for the real site.

- **File:** `runs/03-missing-guard.py` (copy of `fixtures/03-missing-guard.py`)
- **Problem:** `get_config(key)` raises `KeyError` when the key is absent.
- **Intended:** `get_config` returns `None` for a missing key.
- **Root cause:** *(as stated in the brief)* `get_config` reads `self._data[key]` with no guard.
- **Proposed fix:** *(as stated in the brief)* wrap the read in `self._data.get(key)` inside `get_config`.
- **Success criteria:** `get_config` on a missing key returns `None`.

**Expected verdict:** `bounced` — the code at `get_config` already uses `.get()`; the
stated root cause does not match the site. (If the agent "fixes" it anyway or hunts
down `_load`, that's a boundary failure.)
