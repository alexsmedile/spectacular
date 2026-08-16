package workspace

// RefField is the current spelling of a record's human-readable reference.
const RefField = "ref"

// LegacyRefField is the superseded spelling. The rename from human_ref to ref
// stopped halfway: P1 through P6 and M2 through M4 use the legacy spelling,
// while M5, M6, and later records use the current one.
const LegacyRefField = "human_ref"

// Ref resolves a record's human-readable reference through one decoder, so
// every caller compares refs that are spelled one way regardless of how the
// record was authored. It returns the ref, whether the legacy spelling supplied
// it, and any decoding error.
//
// The legacy spelling is reported, not refused. Frozen records are not
// rewritten to finish a rename, so a reader must accept both while a checker
// reports the drift.
func Ref(document *Document) (string, bool, error) {
	current, err := String(document, RefField, false)
	if err != nil {
		return "", false, err
	}
	if current != "" {
		return current, false, nil
	}
	legacy, err := String(document, LegacyRefField, false)
	if err != nil {
		return "", false, err
	}
	return legacy, legacy != "", nil
}

// RefOrEmpty resolves a ref and discards the drift signal and any error. It
// suits presentation paths that must not fail on a malformed optional field;
// callers that report drift or refuse should use Ref.
func RefOrEmpty(document *Document) string {
	ref, _, err := Ref(document)
	if err != nil {
		return ""
	}
	return ref
}
