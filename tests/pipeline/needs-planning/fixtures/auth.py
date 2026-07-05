"""P2 fixture — a bug whose fix is NOT mechanical (needs a request, not a Fixer).

The reported symptom is small ("expired sessions still authenticate"), but the root
cause is structural: there is NO expiry concept anywhere — sessions are a plain dict
with no timestamp, no TTL, no sweep. "Fixing" it means designing session expiry:
schema change (add issued_at), a TTL policy, a sweep or lazy-check, and touching every
create/read path. That's multi-step design work — a PLAN/TASKS request, not a closed
one-site fix.

The Investigator should find this and report it; the orchestrator should route it to
the request lifecycle, NOT fan out a Fixer. This scenario tests that routing fork.
"""

_sessions = {}


def create_session(user_id, token):
    # no issued_at, no ttl — the structural gap
    _sessions[token] = {"user": user_id}


def authenticate(token):
    # no expiry check possible — there's nothing to check against
    return _sessions.get(token, {}).get("user")


if __name__ == "__main__":
    create_session(1, "abc")
    # symptom: this "should" expire but there's no mechanism — always authenticates
    print("authenticated:", authenticate("abc"))
