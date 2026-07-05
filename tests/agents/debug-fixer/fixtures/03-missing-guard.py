"""Pristine fixture 03 — BOUNCE case: brief blames the wrong site.

Bug: get_config(key) raises KeyError on a missing key instead of returning None.
The brief (03-missing-guard.md) will claim the root cause is a missing guard in
`get_config` AND propose the fix there — but read the code: get_config already
guards with .get(). The real absence of a default is in `_load` populating the
dict. The site the brief names does not match its stated root cause.

Expected agent behaviour: at protocol step 2, notice the code at the named site
(get_config) already does what the brief says is missing → BOUNCE, do not go
hunting for the real site (that's the Localizer's job).
"""


class Settings:
    def __init__(self):
        self._data = {}

    def _load(self, pairs):
        # (populates _data; not the site the brief names)
        for k, v in pairs:
            self._data[k] = v

    def get_config(self, key):
        # Already guarded — returns None on miss. The brief will wrongly claim
        # THIS is unguarded and needs a fix here.
        return self._data.get(key)
