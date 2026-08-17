package missionbundle

import (
	"path/filepath"
	"strings"
	"testing"

	"github.com/alexsmedile/spectacular/v2/internal/discovery"
)

func TestFrozenFallbacksValidationAndFingerprinting(t *testing.T) {
	root, err := filepath.Abs(filepath.Join("..", ".."))
	if err != nil {
		t.Fatal(err)
	}
	ws, err := discovery.Open(root)
	if err != nil {
		t.Fatal(err)
	}
	base, err := Load(ws, "M6")
	if err != nil {
		t.Fatal(err)
	}

	t.Run("nil or empty fallbacks validates cleanly", func(t *testing.T) {
		candidate := cloneBundle(t, base)
		candidate.Fallbacks = nil
		if err := validateFallbacks(ws, candidate); err != nil {
			t.Fatalf("nil fallbacks error: %v", err)
		}
		candidate.Fallbacks = []Fallback{}
		if err := validateFallbacks(ws, candidate); err != nil {
			t.Fatalf("empty fallbacks error: %v", err)
		}
	})

	t.Run("missing approach refused", func(t *testing.T) {
		candidate := cloneBundle(t, base)
		candidate.Fallbacks = []Fallback{
			{Approach: "", RejectedBecause: "too slow", InvalidatedIf: "faster compiler"},
		}
		err := validateFallbacks(ws, candidate)
		assertFieldRefusal(t, err, "fallbacks.approach")
	})

	t.Run("missing rejected_because refused", func(t *testing.T) {
		candidate := cloneBundle(t, base)
		candidate.Fallbacks = []Fallback{
			{Approach: "two roots", RejectedBecause: "   ", InvalidatedIf: "single decoder fails"},
		}
		err := validateFallbacks(ws, candidate)
		assertFieldRefusal(t, err, "fallbacks.rejected_because")
	})

	t.Run("missing invalidated_if refused", func(t *testing.T) {
		candidate := cloneBundle(t, base)
		candidate.Fallbacks = []Fallback{
			{Approach: "two roots", RejectedBecause: "too complex", InvalidatedIf: ""},
		}
		err := validateFallbacks(ws, candidate)
		assertFieldRefusal(t, err, "fallbacks.invalidated_if")
	})
}

func TestFallbackFingerprintInvalidationVsRunProgress(t *testing.T) {
	root, err := filepath.Abs(filepath.Join("..", ".."))
	if err != nil {
		t.Fatal(err)
	}
	ws, err := discovery.Open(root)
	if err != nil {
		t.Fatal(err)
	}
	base, err := Load(ws, "M6")
	if err != nil {
		t.Fatal(err)
	}

	t.Run("nil fallbacks preserves historical fingerprints across M2-M7", func(t *testing.T) {
		for _, ref := range []string{"M2", "M3", "M4", "M5", "M6", "M7"} {
			m, err := Load(ws, ref)
			if err != nil {
				t.Fatalf("load %s: %v", ref, err)
			}
			if m.Activation == nil || m.Legacy {
				continue
			}
			fp, err := FrozenFingerprint(m)
			if err != nil {
				t.Fatalf("fingerprint %s: %v", ref, err)
			}
			if fp != m.Activation.Fingerprint {
				t.Fatalf("%s fingerprint mismatch: got=%s want=%s", ref, fp, m.Activation.Fingerprint)
			}
		}
	})

	bundle := cloneBundle(t, base)
	bundle.Fallbacks = []Fallback{
		{
			Approach:        "Keep two package roots",
			RejectedBecause: "Doubles the decode surface and migration cost",
			InvalidatedIf:   "A single decoder cannot preserve v2 readability without field loss",
		},
		{
			Approach:        "Use unstructured maps",
			RejectedBecause: "No schema type safety",
			InvalidatedIf:   "Type system is too restrictive",
		},
	}
	initialFP, err := FrozenFingerprint(bundle)
	if err != nil {
		t.Fatal(err)
	}

	t.Run("mutating fallback approach invalidates activation fingerprint", func(t *testing.T) {
		mutated := cloneBundle(t, bundle)
		mutated.Fallbacks[0].Approach = "Invent a third package root"
		fp, err := FrozenFingerprint(mutated)
		if err != nil {
			t.Fatal(err)
		}
		if fp == initialFP {
			t.Fatalf("mutating fallback approach must change fingerprint: got=%s", fp)
		}
	})

	t.Run("mutating fallback invalidated_if invalidates activation fingerprint", func(t *testing.T) {
		mutated := cloneBundle(t, bundle)
		mutated.Fallbacks[0].InvalidatedIf = "Never"
		fp, err := FrozenFingerprint(mutated)
		if err != nil {
			t.Fatal(err)
		}
		if fp == initialFP {
			t.Fatalf("mutating fallback condition must change fingerprint: got=%s", fp)
		}
	})

	t.Run("mutable Run progress does not invalidate activation fingerprint", func(t *testing.T) {
		mutated := cloneBundle(t, bundle)
		mutated.Run.Repairs += 2
		mutated.Run.CurrentObjective = "O3"
		mutated.Objectives[0].Status = "implemented"
		fp, err := FrozenFingerprint(mutated)
		if err != nil {
			t.Fatal(err)
		}
		if fp != initialFP {
			t.Fatalf("mutable run progress must preserve fingerprint: got=%s want=%s", fp, initialFP)
		}
	})
}

func TestRepairExhaustionSurfacesFullFallbackSetWithRecommendations(t *testing.T) {
	root, err := filepath.Abs(filepath.Join("..", ".."))
	if err != nil {
		t.Fatal(err)
	}
	ws, err := discovery.Open(root)
	if err != nil {
		t.Fatal(err)
	}
	base, err := Load(ws, "M6")
	if err != nil {
		t.Fatal(err)
	}

	bundle := cloneBundle(t, base)
	bundle.Status = "active"
	bundle.Run.Status = "active"
	bundle.RepairBudget = 2
	bundle.Run.Repairs = 2
	bundle.Fallbacks = []Fallback{
		{
			Approach:        "Keep two package roots",
			RejectedBecause: "Doubles the decode surface",
			InvalidatedIf:   "Single decoder cannot preserve readability",
			Recommendation:  true,
		},
		{
			Approach:        "Generate parallel schema adapters",
			RejectedBecause: "High code generation overhead",
			InvalidatedIf:   "Manual decoding becomes unmaintainable",
			Recommendation:  false,
		},
	}

	state := bundle.Derive()
	if !strings.Contains(state.Next, "repair budget is exhausted") {
		t.Fatalf("state.Next=%s want repair budget exhausted", state.Next)
	}
	if len(state.Fallbacks) != 2 {
		t.Fatalf("state.Fallbacks count=%d want 2", len(state.Fallbacks))
	}
	// Verify that state contains both the recommended and alternative fallback
	hasRec, hasAlt := false, false
	for _, fb := range state.Fallbacks {
		if fb.Recommendation && fb.Approach == "Keep two package roots" {
			hasRec = true
		}
		if !fb.Recommendation && fb.Approach == "Generate parallel schema adapters" {
			hasAlt = true
		}
	}
	if !hasRec || !hasAlt {
		t.Fatalf("state.Fallbacks must contain both recommended and alternative fallbacks: %+v", state.Fallbacks)
	}
}
