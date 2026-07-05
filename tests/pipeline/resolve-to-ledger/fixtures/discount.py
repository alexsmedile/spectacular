"""P1 fixture — a closed, reusable bug.

apply_discount is meant to cap the discount at 90% (floor price at 10% of list).
It caps at the wrong bound — uses max() where it needs min() — so a 95% discount
coupon charges *more* than list instead of flooring the price.

Root cause is obvious (max vs min), single site, reusable lesson (a classic
clamp-direction footgun) → this job should resolve AND graduate to fixes/F<N>.
"""


def apply_discount(price, pct):
    # bug: max() lets pct above the cap through; should be min(pct, 0.9)
    capped = max(pct, 0.9)
    return price * (1 - capped)


if __name__ == "__main__":
    # a 95% coupon should floor at 90% off → pay 10 on a 100 item
    # round() to dodge IEEE-754 noise (100*(1-0.9) == 9.999...); the bug under
    # test is clamp-direction, not float equality
    assert round(apply_discount(100, 0.95), 2) == 10.0, f"got {apply_discount(100, 0.95)}"
    assert round(apply_discount(100, 0.5), 2) == 50.0
    print("ok")
