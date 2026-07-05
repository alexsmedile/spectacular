"""P5 capstone fixture — an OPEN bug that flows the whole pipeline.

Symptom (what the user reports): "the cart total is sometimes wrong after removing
an item." Cause + site are NOT given — the orchestrator must investigate.

Real cause (for the test author; the Investigator should discover this): remove_item
deletes from the _items list but never recomputes _total, and get_total returns the
stale cached _total. There are TWO sites that share the same root — add_item correctly
updates _total, remove_item does not, and there's a latent gap: clear() also doesn't.
So the honest fix is at the shared root (recompute-on-mutation), which the Investigator
should flag as blast-radius across mutators.

This is open (discover), has a shared root (planning judgment), and resolves to a
reusable lesson (cache-coherence on mutation) → exercises Investigator → plan → fix →
resolve → log, the full arc.
"""


class Cart:
    def __init__(self):
        self._items = []
        self._total = 0

    def add_item(self, price):
        self._items.append(price)
        self._total += price  # correctly maintained here

    def remove_item(self, price):
        self._items.remove(price)
        # BUG: _total not updated → get_total returns stale value

    def get_total(self):
        return self._total  # returns the stale cache


if __name__ == "__main__":
    c = Cart()
    c.add_item(10)
    c.add_item(5)
    c.remove_item(5)
    assert c.get_total() == 10, f"got {c.get_total()}"  # fails: still 15
    print("ok")
