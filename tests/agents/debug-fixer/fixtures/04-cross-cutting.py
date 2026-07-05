"""Pristine fixture 04 — BOUNCE case: cross-cutting, not single-site.

Bug: prices render without currency formatting across three views. The brief
(04-cross-cutting.md) will name ONE view (cart_line) as the single site and
propose fixing the format string there. But all three views call the same broken
root — `fmt_money` drops the currency symbol. Fixing cart_line alone leaves
receipt_line and invoice_line still broken (fix-root-not-symptom).

Expected agent behaviour: at protocol step 3, notice the fix needs to touch the
shared root (fmt_money) to be correct, i.e. more than the single named site →
BOUNCE as cross-cutting. Do NOT patch cart_line and call it done.
"""


def fmt_money(amount):
    # ROOT: drops the currency symbol. Every caller below is wrong because of this.
    return f"{amount:.2f}"


def cart_line(item, amount):        # the site the brief will name
    return f"{item}: {fmt_money(amount)}"


def receipt_line(item, amount):     # sibling caller — also broken
    return f"{item} .... {fmt_money(amount)}"


def invoice_line(item, amount):     # sibling caller — also broken
    return f"{item} | {fmt_money(amount)}"
