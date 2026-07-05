"""P4 fixture C — independent closed bug (wrong default)."""


def paginate(items, per_page=0):
    # bug: default 0 → ZeroDivisionError / empty page; should default to 10
    per_page = per_page or 0
    return [items[i:i + per_page] for i in range(0, len(items), per_page)] if per_page else []


if __name__ == "__main__":
    # with the fix (default 10), a 3-item list on default paging returns one page of 3
    assert paginate([1, 2, 3]) == [[1, 2, 3]], f"got {paginate([1,2,3])}"
    print("ok")
