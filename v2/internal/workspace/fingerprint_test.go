package workspace

import "testing"

const stableMission = "---\n" +
	"type: Mission\n" +
	"id: 018f2d8e-7b13-7aa1-9b34-acdeffedcba9\n" +
	"title: Implement M1\n" +
	"status: active\n" +
	"source: Proposal:018f2d8e-7b12-7cc3-8a45-123456789abc\n" +
	"priority: high\n" +
	"details:\n" +
	"  owner: central\n" +
	"  retries: 2\n" +
	"---\n" +
	"# Mission\n\nBuild the substrate.\n"

func TestFingerprintStableAcrossPropertyOrderQuotingAndLineEndings(t *testing.T) {
	t.Parallel()

	variant := "---\r\n" +
		"details:\r\n" +
		"  retries: 2\r\n" +
		"  owner: 'central'\r\n" +
		"priority: \"high\"\r\n" +
		"source: \"Proposal:018f2d8e-7b12-7cc3-8a45-123456789abc\"\r\n" +
		"status: \"active\"\r\n" +
		"title: 'Implement M1'\r\n" +
		"id: 018f2d8e-7b13-7aa1-9b34-acdeffedcba9\r\n" +
		"type: \"Mission\"\r\n" +
		"---\r\n" +
		"# Mission\r\n\r\nBuild the substrate.\r\n"

	baseFingerprint := fingerprintForText(t, stableMission)
	variantFingerprint := fingerprintForText(t, variant)
	if baseFingerprint != variantFingerprint {
		t.Fatalf("equivalent records have different fingerprints:\nbase:    %s\nvariant: %s", baseFingerprint, variantFingerprint)
	}
}

func TestFingerprintChangesWhenSemanticMeaningChanges(t *testing.T) {
	t.Parallel()

	base := fingerprintForText(t, stableMission)
	changes := map[string]string{
		"known field":   replaceOnce(t, stableMission, "title: Implement M1", "title: Implement M1 safely"),
		"unknown field": replaceOnce(t, stableMission, "priority: high", "priority: low"),
		"relationship": replaceOnce(t, stableMission,
			"Proposal:018f2d8e-7b12-7cc3-8a45-123456789abc",
			"Proposal:018f2d8e-7b18-7000-8000-000000000003"),
		"body": replaceOnce(t, stableMission, "Build the substrate.", "Build a different substrate."),
	}
	for name, changed := range changes {
		changed := changed
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			if got := fingerprintForText(t, changed); got == base {
				t.Fatalf("fingerprint did not change for %s", name)
			}
		})
	}
}

func TestFingerprintPreservesUnknownYAMLTagMeaning(t *testing.T) {
	t.Parallel()

	binary := "---\n" +
		"type: Proposal\n" +
		"id: 018f2d8e-7b12-7cc3-8a45-123456789abc\n" +
		"opaque: !!binary SGVsbG8=\n" +
		"---\n"
	plain := "---\n" +
		"type: Proposal\n" +
		"id: 018f2d8e-7b12-7cc3-8a45-123456789abc\n" +
		"opaque: Hello\n" +
		"---\n"
	binaryDocument, err := Parse([]byte(binary))
	if err != nil {
		t.Fatal(err)
	}
	binaryCanonical, err := Canonical(binaryDocument)
	if err != nil {
		t.Fatal(err)
	}
	binaryRoundTrip, err := Parse(binaryCanonical)
	if err != nil {
		t.Fatal(err)
	}
	if opaque := binaryRoundTrip.Unknown["opaque"]; opaque == nil || opaque.ShortTag() != "!!binary" {
		t.Fatalf("binary tag was not preserved: %#v\n%s", opaque, binaryCanonical)
	}
	if binaryFingerprint, plainFingerprint := fingerprintForText(t, binary), fingerprintForText(t, plain); binaryFingerprint == plainFingerprint {
		t.Fatalf("tag-distinct YAML values collapsed to fingerprint %s", binaryFingerprint)
	}
}

func TestFingerprintDistinguishesAliasGraphFromDuplicatedTree(t *testing.T) {
	t.Parallel()

	aliasGraph := "---\n" +
		"type: Proposal\n" +
		"id: 018f2d8e-7b12-7cc3-8a45-123456789abc\n" +
		"opaque:\n" +
		"  first: &shared\n" +
		"    value: retained\n" +
		"  second: *shared\n" +
		"---\n"
	duplicatedTree := "---\n" +
		"type: Proposal\n" +
		"id: 018f2d8e-7b12-7cc3-8a45-123456789abc\n" +
		"opaque:\n" +
		"  first:\n" +
		"    value: retained\n" +
		"  second:\n" +
		"    value: retained\n" +
		"---\n"
	aliasFingerprint := fingerprintForText(t, aliasGraph)
	duplicatedFingerprint := fingerprintForText(t, duplicatedTree)
	if aliasFingerprint == duplicatedFingerprint {
		t.Fatalf("alias graph and duplicated tree collapsed to fingerprint %s", aliasFingerprint)
	}
}

func fingerprintForText(t *testing.T, text string) string {
	t.Helper()
	document, err := Parse([]byte(text))
	if err != nil {
		t.Fatalf("Parse: %v", err)
	}
	fingerprint, err := Fingerprint(document)
	if err != nil {
		t.Fatalf("Fingerprint: %v", err)
	}
	return fingerprint
}

func replaceOnce(t *testing.T, input, old, replacement string) string {
	t.Helper()
	for index := 0; index+len(old) <= len(input); index++ {
		if input[index:index+len(old)] == old {
			return input[:index] + replacement + input[index+len(old):]
		}
	}
	t.Fatalf("%q not found", old)
	return ""
}
