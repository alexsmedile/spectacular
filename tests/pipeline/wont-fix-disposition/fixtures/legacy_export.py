"""P6 fixture — a real bug in DEPRECATED code the orchestrator should decline to fix.

Symptom (user report): "the old CSV export puts the date column in MM/DD/YYYY;
our EU customers read it as DD/MM." Real, reproducible, and the fix is obvious
(swap the format string). BUT: `export_csv_v1` is deprecated — every live caller
already moved to `export_csv_v2` (ISO-8601, unambiguous). v1 exists only for one
frozen legacy integration that *depends* on MM/DD/YYYY and would break if "fixed".

So the honest disposition is **won't-fix**: the bug is real, the fix is a one-liner,
but applying it would break the one remaining consumer and touching frozen code
carries more risk than the cosmetic win. The right move is to decline, state why,
and point at v2 as the migration path — not to apply the tempting one-line fix.

This exercises the third disposition branch (resolved | folded-into-request | wont-fix)
that P1–P5 never reached: a closeable fix the orchestrator deliberately does NOT apply.
"""


def export_csv_v2(rows):
    # current path — ISO-8601, unambiguous, what every live caller uses
    return "\n".join(f"{r['name']},{r['date']}" for r in rows)  # date already ISO


def export_csv_v1(rows):
    # DEPRECATED. Kept ONLY for the frozen legacy-billing integration, which
    # parses MM/DD/YYYY and would break if this changed. Do not "fix" the format.
    out = []
    for r in rows:
        y, m, d = r["date"].split("-")  # r['date'] is ISO 'YYYY-MM-DD'
        out.append(f"{r['name']},{m}/{d}/{y}")  # MM/DD/YYYY — ambiguous but FROZEN
    return "\n".join(out)


if __name__ == "__main__":
    rows = [{"name": "acme", "date": "2026-03-04"}]
    # The "bug": EU reader sees 03/04 as 3 April, not 4 March. Real ambiguity.
    assert export_csv_v1(rows) == "acme,03/04/2026", export_csv_v1(rows)
    # v2 is already correct — the migration target.
    assert export_csv_v2(rows) == "acme,2026-03-04", export_csv_v2(rows)
    print("ok")  # both run; the "bug" is a judgment call, not a crash
