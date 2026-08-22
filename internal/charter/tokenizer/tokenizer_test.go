package tokenizer

import (
	"strings"
	"testing"
)

func TestTokenizer_UTF8Validation(t *testing.T) {
	validText := "Hello, Spectacular v2 Context Sandwich!"
	count, err := Count(validText)
	if err != nil {
		t.Fatalf("unexpected error for valid UTF-8: %v", err)
	}
	if count <= 0 {
		t.Fatalf("expected positive count, got %d", count)
	}

	invalidBytes := string([]byte{0xff, 0xfe, 0xfd})
	_, err = Count(invalidBytes)
	if err != ErrInvalidUTF8 {
		t.Fatalf("expected ErrInvalidUTF8, got %v", err)
	}
}

func TestTokenizer_Digest(t *testing.T) {
	if TokenizerDataDigest == "" {
		t.Fatal("expected non-empty tokenizer data digest")
	}
	if len(TokenizerDataDigest) != 64 {
		t.Fatalf("expected 64 hex characters for SHA-256, got %d", len(TokenizerDataDigest))
	}
}

func TestTokenizer_ThresholdDispositions(t *testing.T) {
	tests := []struct {
		tokens      int
		disposition Disposition
	}{
		{0, DispositionPass},
		{500, DispositionPass},
		{1200, DispositionPass},
		{1201, DispositionWarn},
		{1350, DispositionWarn},
		{1400, DispositionWarn},
		{1401, DispositionSplitRecommended},
		{1420, DispositionSplitRecommended},
		{1440, DispositionSplitRecommended},
		{1441, DispositionRefusal},
		{2000, DispositionRefusal},
	}

	for _, tc := range tests {
		disp, _ := EvaluateDisposition(tc.tokens)
		if disp != tc.disposition {
			t.Errorf("at tokens %d: expected %s, got %s", tc.tokens, tc.disposition, disp)
		}
	}
}

func TestTokenizer_Repeatability(t *testing.T) {
	input := strings.Repeat("Layer 1: Frozen Truth\nLayer 2: Steering\nLayer 3: Execution Perimeter\n", 50)
	c1, err1 := Count(input)
	c2, err2 := Count(input)
	if err1 != nil || err2 != nil {
		t.Fatalf("unexpected error: %v, %v", err1, err2)
	}
	if c1 != c2 {
		t.Fatalf("expected identical counts, got %d and %d", c1, c2)
	}
}
