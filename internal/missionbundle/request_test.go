package missionbundle

import (
	"errors"
	"path/filepath"
	"testing"

	"github.com/alexsmedile/spectacular/v2/internal/discovery"
	"github.com/alexsmedile/spectacular/v2/internal/domain"
)

func TestRequestFidelityValidationAndRefusals(t *testing.T) {
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

	t.Run("nil request validates cleanly", func(t *testing.T) {
		candidate := cloneBundle(t, base)
		candidate.Request = nil
		if err := validateRequest(ws, candidate); err != nil {
			t.Fatalf("expected nil error for nil request, got %v", err)
		}
	})

	t.Run("missing source refused", func(t *testing.T) {
		candidate := cloneBundle(t, base)
		candidate.Request = &Request{
			Source:     "",
			CapturedAt: "2026-08-16T20:59:04Z",
			Asks:       []Ask{{Ask: "do something", Disposition: "covered", Claims: []string{candidate.Completion[0].Claim}}},
		}
		err := validateRequest(ws, candidate)
		assertFieldRefusal(t, err, "request.source")
	})

	t.Run("invalid captured_at timestamp refused", func(t *testing.T) {
		candidate := cloneBundle(t, base)
		candidate.Request = &Request{
			Source:     "chat",
			CapturedAt: "not-a-timestamp",
			Asks:       []Ask{{Ask: "do something", Disposition: "covered", Claims: []string{candidate.Completion[0].Claim}}},
		}
		err := validateRequest(ws, candidate)
		assertFieldRefusal(t, err, "request.captured_at")
	})

	t.Run("empty asks list refused", func(t *testing.T) {
		candidate := cloneBundle(t, base)
		candidate.Request = &Request{
			Source:     "chat",
			CapturedAt: "2026-08-16T20:59:04Z",
			Asks:       []Ask{},
		}
		err := validateRequest(ws, candidate)
		assertFieldRefusal(t, err, "request.asks")
	})

	t.Run("empty ask text refused", func(t *testing.T) {
		candidate := cloneBundle(t, base)
		candidate.Request = &Request{
			Source:     "chat",
			CapturedAt: "2026-08-16T20:59:04Z",
			Asks:       []Ask{{Ask: "", Disposition: "covered", Claims: []string{candidate.Completion[0].Claim}}},
		}
		err := validateRequest(ws, candidate)
		assertFieldRefusal(t, err, "request.asks.ask")
	})

	t.Run("undispositioned ask refused", func(t *testing.T) {
		candidate := cloneBundle(t, base)
		candidate.Request = &Request{
			Source:     "chat",
			CapturedAt: "2026-08-16T20:59:04Z",
			Asks: []Ask{
				{Ask: "build something", Disposition: ""},
			},
		}
		err := validateRequest(ws, candidate)
		assertFieldRefusal(t, err, "request.asks.disposition")
	})

	t.Run("covered ask naming nonexistent claim refused", func(t *testing.T) {
		candidate := cloneBundle(t, base)
		candidate.Request = &Request{
			Source:     "chat",
			CapturedAt: "2026-08-16T20:59:04Z",
			Asks: []Ask{
				{Ask: "build something", Disposition: "covered", Claims: []string{"nonexistent-claim"}},
			},
		}
		err := validateRequest(ws, candidate)
		assertFieldRefusal(t, err, "request.asks.claims")
	})

	t.Run("covered ask with empty claims refused", func(t *testing.T) {
		candidate := cloneBundle(t, base)
		candidate.Request = &Request{
			Source:     "chat",
			CapturedAt: "2026-08-16T20:59:04Z",
			Asks: []Ask{
				{Ask: "build something", Disposition: "covered", Claims: nil},
			},
		}
		err := validateRequest(ws, candidate)
		assertFieldRefusal(t, err, "request.asks.claims")
	})

	t.Run("deferred ask without reason refused", func(t *testing.T) {
		candidate := cloneBundle(t, base)
		candidate.Request = &Request{
			Source:     "chat",
			CapturedAt: "2026-08-16T20:59:04Z",
			Asks: []Ask{
				{Ask: "defer something", Disposition: "deferred", Reason: ""},
			},
		}
		err := validateRequest(ws, candidate)
		assertFieldRefusal(t, err, "request.asks.reason")
	})

	t.Run("declined ask without reason refused", func(t *testing.T) {
		candidate := cloneBundle(t, base)
		candidate.Request = &Request{
			Source:     "chat",
			CapturedAt: "2026-08-16T20:59:04Z",
			Asks: []Ask{
				{Ask: "decline something", Disposition: "declined", Reason: "  "},
			},
		}
		err := validateRequest(ws, candidate)
		assertFieldRefusal(t, err, "request.asks.reason")
	})

	t.Run("validator never infers disposition from ask text matching claim", func(t *testing.T) {
		candidate := cloneBundle(t, base)
		candidate.Request = &Request{
			Source:     "chat",
			CapturedAt: "2026-08-16T20:59:04Z",
			Asks: []Ask{
				{Ask: candidate.Completion[0].Claim, Disposition: ""},
			},
		}
		err := validateRequest(ws, candidate)
		assertFieldRefusal(t, err, "request.asks.disposition")
	})
}

func TestRequestFingerprintPreservesTextEditsAndInvalidatesDispositionEdits(t *testing.T) {
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

	t.Run("nil request preserves historical fingerprint identically", func(t *testing.T) {
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
	bundle.Request = &Request{
		Source:     "chat, session opening",
		CapturedAt: "2026-08-16T20:59:04Z",
		Asks: []Ask{
			{
				Ask:         "implement the feature",
				Disposition: "covered",
				Claims:      []string{bundle.Completion[0].Claim},
			},
			{
				Ask:         "perform extra refactoring",
				Disposition: "deferred",
				Reason:      "outside scope",
			},
		},
	}
	initialFP, err := FrozenFingerprint(bundle)
	if err != nil {
		t.Fatal(err)
	}

	t.Run("editing ask prose text preserves fingerprint", func(t *testing.T) {
		edited := cloneBundle(t, bundle)
		edited.Request.Asks[0].Ask = "sharpened and elaborated ask prose description"
		edited.Request.Asks[1].Reason = "sharpened reason text with more explanation"
		edited.Request.Source = "updated source location"
		edited.Request.CapturedAt = "2026-08-17T00:00:00Z"
		fp, err := FrozenFingerprint(edited)
		if err != nil {
			t.Fatal(err)
		}
		if fp != initialFP {
			t.Fatalf("editing request text changed fingerprint: got=%s want=%s", fp, initialFP)
		}
	})

	t.Run("changing disposition from covered to deferred invalidates fingerprint", func(t *testing.T) {
		mutated := cloneBundle(t, bundle)
		mutated.Request.Asks[0].Disposition = "deferred"
		mutated.Request.Asks[0].Reason = "changed my mind"
		mutated.Request.Asks[0].Claims = nil
		fp, err := FrozenFingerprint(mutated)
		if err != nil {
			t.Fatal(err)
		}
		if fp == initialFP {
			t.Fatalf("changing disposition should invalidate fingerprint, got identical %s", fp)
		}
	})

	t.Run("changing covered claim invalidates fingerprint", func(t *testing.T) {
		mutated := cloneBundle(t, bundle)
		if len(bundle.Completion) > 1 {
			mutated.Request.Asks[0].Claims = []string{bundle.Completion[1].Claim}
			fp, err := FrozenFingerprint(mutated)
			if err != nil {
				t.Fatal(err)
			}
			if fp == initialFP {
				t.Fatalf("changing covered claim should invalidate fingerprint, got identical %s", fp)
			}
		}
	})
}

func assertFieldRefusal(t *testing.T, err error, field string) {
	t.Helper()
	var refusal *domain.Refusal
	if !errors.As(err, &refusal) {
		t.Fatalf("expected *domain.Refusal, got %v", err)
	}
	if refusal.Field != field {
		t.Fatalf("refusal.Field=%s, want=%s (code=%s detail=%s)", refusal.Field, field, refusal.Code, refusal.Detail)
	}
}

func TestCompletionRefusesWhenAskIsUndispositioned(t *testing.T) {
	root := missionServiceFixture(t)
	plan, raw := stressPlan()
	plan.Title = "Request completion test"
	plan.Request = &Request{
		Source:     "chat",
		CapturedAt: "2026-08-16T20:59:04Z",
		Asks: []Ask{
			{Ask: "Valid ask", Disposition: "covered", Claims: []string{plan.Completion[0].Claim}},
			{Ask: "Undispositioned ask", Disposition: ""},
		},
	}
	svc := openMissionService(t, root)
	_, err := svc.Start(plan, raw)
	if err == nil {
		t.Fatal("expected start/validation refusal on undispositioned ask, got nil")
	}
	assertFieldRefusal(t, err, "request.asks.disposition")
}
