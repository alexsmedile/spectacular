package missionbundle

import (
	"path/filepath"
	"strings"
	"testing"

	"github.com/alexsmedile/spectacular/v2/internal/discovery"
	"github.com/alexsmedile/spectacular/v2/internal/domain"
	"github.com/alexsmedile/spectacular/v2/internal/workspace"
)

// A declaration without amend-contract in requires_owner claims an authority the
// record never granted, so it refuses before any Gap ref is examined.
func TestResolvedGapRequiresDeclaredOwnerAuthority(t *testing.T) {
	root, err := filepath.Abs(filepath.Join("..", ".."))
	if err != nil {
		t.Fatal(err)
	}
	ws, err := discovery.Open(root)
	if err != nil {
		t.Fatal(err)
	}
	source, err := Load(ws, "M9")
	if err != nil {
		t.Fatal(err)
	}
	bundle := cloneBundle(t, source)
	bundle.ResolvesGaps = []ResolvedGap{{Gap: "anything", Resolution: "text"}}
	err = validateResolvedGaps(ws, bundle)
	if err == nil {
		t.Fatal("a declaration without amend-contract authority must refuse")
	}
	refusal, ok := err.(*domain.Refusal)
	if !ok {
		t.Fatalf("returned %T, want a typed refusal", err)
	}
	if refusal.Field != "authority.requires_owner" {
		t.Fatalf("refused on %q, want authority.requires_owner", refusal.Field)
	}
}

// A Mission may only close Gaps the Contract it is bound to actually declares.
// The refs are held at plan-freeze rather than at completion: a Mission that
// discovers at the gate that it never had the authority to close a Gap has
// already done the work under a false premise.
func TestResolvedGapDeclarationIsHeldToTheBoundContract(t *testing.T) {
	root, err := filepath.Abs(filepath.Join("..", ".."))
	if err != nil {
		t.Fatal(err)
	}
	ws, err := discovery.Open(root)
	if err != nil {
		t.Fatal(err)
	}
	// M9 is bound to the Contract that carries the Gaps, so its own bound Contract
	// supplies both a real Gap ref and the absence of an invented one.
	source, err := Load(ws, "M9")
	if err != nil {
		t.Fatal(err)
	}
	gaps, err := ContractGaps(ws, source.Contract.Ref)
	if err != nil {
		t.Fatal(err)
	}
	if len(gaps) == 0 {
		t.Fatal("the Contract M9 is bound to declares no Gaps; this test needs at least one")
	}
	real := gaps[0].Ref

	tests := []struct {
		name    string
		declare []ResolvedGap
		field   string
	}{
		{"absent", nil, ""},
		{"empty", []ResolvedGap{}, ""},
		{"valid", []ResolvedGap{{Gap: real, Resolution: "closed by this Mission"}}, ""},
		{"dangling ref", []ResolvedGap{{Gap: "no-such-gap", Resolution: "closed"}}, "resolves_gaps.gap"},
		{"unnamed gap", []ResolvedGap{{Gap: "   ", Resolution: "closed"}}, "resolves_gaps.gap"},
		{"missing resolution", []ResolvedGap{{Gap: real}}, "resolves_gaps.resolution"},
		{"blank resolution", []ResolvedGap{{Gap: real, Resolution: "  "}}, "resolves_gaps.resolution"},
		{"declared twice", []ResolvedGap{
			{Gap: real, Resolution: "closed"},
			{Gap: real, Resolution: "closed again"},
		}, "resolves_gaps.gap"},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			bundle := cloneBundle(t, source)
			bundle.Authority.RequiresOwner = append(bundle.Authority.RequiresOwner, "amend-contract")
			bundle.ResolvesGaps = test.declare
			err := validateResolvedGaps(ws, bundle)
			if test.field == "" {
				if err != nil {
					t.Fatalf("declaration must validate, got %v", err)
				}
				return
			}
			if err == nil {
				t.Fatal("declaration must refuse")
			}
			refusal, ok := err.(*domain.Refusal)
			if !ok {
				t.Fatalf("returned %T, want a typed refusal", err)
			}
			if refusal.Field != test.field {
				t.Fatalf("refused on %q, want %q", refusal.Field, test.field)
			}
		})
	}
}

// A Gap ref that exists on some other Contract is still not one this Mission may
// close. The check is against the bound Contract, not the workspace.
func TestResolvedGapRefusesAGapOnADifferentContract(t *testing.T) {
	root, err := filepath.Abs(filepath.Join("..", ".."))
	if err != nil {
		t.Fatal(err)
	}
	ws, err := discovery.Open(root)
	if err != nil {
		t.Fatal(err)
	}
	// M10 is bound to CC-missioncli; the Gaps live on the Contract M9 is bound to.
	subject, err := Load(ws, "M10")
	if err != nil {
		t.Fatal(err)
	}
	elsewhere, err := Load(ws, "M9")
	if err != nil {
		t.Fatal(err)
	}
	if subject.Contract.Ref == elsewhere.Contract.Ref {
		t.Skip("M10 and M9 are bound to the same Contract; this test needs two")
	}
	gaps, err := ContractGaps(ws, elsewhere.Contract.Ref)
	if err != nil {
		t.Fatal(err)
	}
	if len(gaps) == 0 {
		t.Skip("the other Contract declares no Gaps to borrow")
	}

	bundle := cloneBundle(t, subject)
	bundle.Authority.RequiresOwner = append(bundle.Authority.RequiresOwner, "amend-contract")
	bundle.ResolvesGaps = []ResolvedGap{{Gap: gaps[0].Ref, Resolution: "closed"}}
	err = validateResolvedGaps(ws, bundle)
	if err == nil {
		t.Fatalf("a Gap on a different Contract must refuse; %q was accepted", gaps[0].Ref)
	}
	refusal, ok := err.(*domain.Refusal)
	if !ok {
		t.Fatalf("returned %T, want a typed refusal", err)
	}
	if refusal.Field != "resolves_gaps.gap" {
		t.Fatalf("refused on %q, want resolves_gaps.gap", refusal.Field)
	}
}

// The declaration is authority to amend a Contract, so it must sit inside the
// frozen envelope: a Mission that could add, retarget, or reword an entry after
// activation would be amending a Contract the owner never approved amending.
func TestResolvedGapsAreFrozenInTheActivationFingerprint(t *testing.T) {
	root, err := filepath.Abs(filepath.Join("..", ".."))
	if err != nil {
		t.Fatal(err)
	}
	ws, err := discovery.Open(root)
	if err != nil {
		t.Fatal(err)
	}
	source, err := Load(ws, "M9")
	if err != nil {
		t.Fatal(err)
	}

	// A Mission frozen before this field existed must still verify, so the field's
	// absence has to leave the envelope byte-identical. Asserted against the record
	// itself rather than a clone: cloneBundle round-trips through JSON and does not
	// reproduce the envelope exactly, so only the loaded Bundle can be compared to
	// the stored fingerprint.
	base, err := FrozenFingerprint(source)
	if err != nil {
		t.Fatal(err)
	}
	if source.ResolvesGaps != nil {
		t.Fatalf("M9 declares %d resolved Gaps; this test needs a Mission frozen without the field", len(source.ResolvesGaps))
	}
	if source.Activation == nil {
		t.Fatal("M9 carries no activation to compare against")
	}
	if base != source.Activation.Fingerprint {
		t.Fatalf("recomputed %s, want the stored %s", base, source.Activation.Fingerprint)
	}
	source.ResolvesGaps = []ResolvedGap{}
	if got, fpErr := FrozenFingerprint(source); fpErr != nil {
		t.Fatal(fpErr)
	} else if got != base {
		t.Fatalf("an empty declaration changed the fingerprint to %s, breaking every Mission frozen before the field existed", got)
	}
	source.ResolvesGaps = nil

	// The variants below are compared against each other rather than against the
	// stored fingerprint, so a clone's round-trip difference is common to all of
	// them and cannot hide a collision.
	cloneBase, err := FrozenFingerprint(cloneBundle(t, source))
	if err != nil {
		t.Fatal(err)
	}

	variants := map[string][]ResolvedGap{
		"one entry":     {{Gap: "a", Resolution: "x"}},
		"two entries":   {{Gap: "a", Resolution: "x"}, {Gap: "b", Resolution: "y"}},
		"reordered":     {{Gap: "b", Resolution: "y"}, {Gap: "a", Resolution: "x"}},
		"reworded":      {{Gap: "a", Resolution: "x."}},
		"retargeted":    {{Gap: "c", Resolution: "x"}},
		"resolution xy": {{Gap: "a", Resolution: "xy"}},
	}
	seen := map[string]string{cloneBase: "absent"}
	for name, declaration := range variants {
		bundle := cloneBundle(t, source)
		bundle.ResolvesGaps = declaration
		got, fpErr := FrozenFingerprint(bundle)
		if fpErr != nil {
			t.Fatal(fpErr)
		}
		if previous, collision := seen[got]; collision {
			t.Fatalf("%q and %q produce the same fingerprint %s", name, previous, got)
		}
		seen[got] = name
	}
}

// mission start decodes the plan through its own struct rather than through the
// Bundle decoder, so a field can validate everywhere else and still be dropped on
// the one path that freezes a Mission. This asserts the plan struct carries it: a
// silently discarded declaration would activate a Mission with no record of the
// authority it was granted, and the amendment at completion would find nothing.
func TestMissionPlanCarriesResolvedGaps(t *testing.T) {
	const plan = `---
type: MissionPlan
title: plan struct check
resolves_gaps:
    - gap: some-gap
      resolution: closed by this Mission
---
body
`
	decoded, _, err := ReadPlan("-", []byte(plan))
	if err != nil {
		t.Fatal(err)
	}
	if len(decoded.ResolvesGaps) != 1 {
		t.Fatalf("the plan struct decoded %d resolved Gaps, want 1; mission start would drop the declaration", len(decoded.ResolvesGaps))
	}
	if decoded.ResolvesGaps[0].Gap != "some-gap" || decoded.ResolvesGaps[0].Resolution == "" {
		t.Fatalf("decoded %+v", decoded.ResolvesGaps[0])
	}
}

// The declaration reaches the Bundle from the record rather than only from Go, so
// a hand-authored plan carrying it is decoded through the same optional-field path
// the decoder uses rather than silently dropped.
func TestResolvedGapsDecodeFromTheRecord(t *testing.T) {
	const record = `---
type: Mission
id: 01a00d6e-df08-75d3-a6c3-5975bb421630
title: decode check
resolves_gaps:
    - gap: some-gap
      resolution: closed by this Mission
    - gap: other-gap
      resolution: also closed
---
body
`
	doc, err := workspace.Parse([]byte(record))
	if err != nil {
		t.Fatal(err)
	}
	var bundle Bundle
	present, err := decodeOptional(doc, "resolves_gaps", &bundle.ResolvesGaps)
	if err != nil {
		t.Fatal(err)
	}
	if !present {
		t.Fatal("resolves_gaps was not decoded from the record")
	}
	if len(bundle.ResolvesGaps) != 2 {
		t.Fatalf("decoded %d entries, want 2", len(bundle.ResolvesGaps))
	}
	if bundle.ResolvesGaps[0].Gap != "some-gap" {
		t.Fatalf("gap=%q want some-gap", bundle.ResolvesGaps[0].Gap)
	}
	if !strings.Contains(bundle.ResolvesGaps[0].Resolution, "closed by this Mission") {
		t.Fatalf("resolution=%q", bundle.ResolvesGaps[0].Resolution)
	}

	// Absence must decode as absence rather than as an error, so a plan that
	// resolves no Gaps needs no field.
	silent, err := workspace.Parse([]byte("---\ntype: Mission\nid: 01a00d6e-df08-75d3-a6c3-5975bb421630\ntitle: none\n---\nbody\n"))
	if err != nil {
		t.Fatal(err)
	}
	var without Bundle
	present, err = decodeOptional(silent, "resolves_gaps", &without.ResolvesGaps)
	if err != nil {
		t.Fatal(err)
	}
	if present || len(without.ResolvesGaps) != 0 {
		t.Fatalf("absent field decoded as present=%t entries=%d", present, len(without.ResolvesGaps))
	}
}
