"""Pristine fixture — OPEN bug: symptom given, cause+site unknown.

Symptom (what the "user" reports): `get_user(1)` sometimes returns the OLD name
after `update_name(1, "new")` was called. Intermittent-looking, but actually
deterministic once you find the cause. The report does NOT tell you where the bug
is or why — that's the Investigator's job to find.

An Investigator should: reproduce (run _demo, see the failure), localize (the cache
in get_user is never invalidated on update), diagnose (update_name writes the store
but not the cache), and emit a CLOSED brief: invalidate the cache entry in
update_name. Then a Fixer applies it.
"""


class Users:
    def __init__(self):
        self._store = {}
        self._cache = {}

    def create(self, uid, name):
        self._store[uid] = name

    def update_name(self, uid, name):
        # writes the store...
        self._store[uid] = name
        # ...but leaves self._cache[uid] holding the old name. (the bug — but the
        # Investigator must discover this; the brief handed in won't name it.)

    def get_user(self, uid):
        if uid in self._cache:
            return self._cache[uid]
        val = self._store.get(uid)
        self._cache[uid] = val
        return val


def _demo():
    u = Users()
    u.create(1, "alice")
    assert u.get_user(1) == "alice"      # caches "alice"
    u.update_name(1, "bob")
    assert u.get_user(1) == "bob", f"stale cache: got {u.get_user(1)!r}"
    print("ok")


if __name__ == "__main__":
    _demo()
