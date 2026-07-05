"""P4 fixture A — independent closed bug (off-by-one in a range)."""


def last_n(items, n):
    # bug: drops the final element; should be items[-n:]
    return items[-n - 1:-1]


if __name__ == "__main__":
    assert last_n([1, 2, 3, 4, 5], 2) == [4, 5], f"got {last_n([1,2,3,4,5],2)}"
    print("ok")
