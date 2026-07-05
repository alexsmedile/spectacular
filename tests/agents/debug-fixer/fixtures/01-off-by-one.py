"""Pristine fixture 01 — off-by-one in a range boundary.

Bug: last_n_items should return the LAST n items, but the slice drops one.
This is a closed, single-site, mechanical fix.
"""


def last_n_items(items, n):
    # BUG: slice starts one too late — returns n-1 items when n < len
    return items[len(items) - n + 1:]


def _demo():
    assert last_n_items([1, 2, 3, 4, 5], 3) == [3, 4, 5], last_n_items([1, 2, 3, 4, 5], 3)
    assert last_n_items([1, 2, 3], 3) == [1, 2, 3]
    assert last_n_items([1, 2, 3], 1) == [3]
    print("ok")


if __name__ == "__main__":
    _demo()
