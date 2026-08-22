package spectaculareval

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/alexsmedile/spectacular/v2/internal/charter"
	"github.com/alexsmedile/spectacular/v2/internal/charter/tokenizer"
	"github.com/alexsmedile/spectacular/v2/internal/discovery"
)

func TestContextCharter_PairedContextEconomyProof(t *testing.T) {
	ws, err := discovery.Open(".")
	if err != nil {
		t.Fatalf("discovery.Open failed: %v", err)
	}

	// 1. Measure Baseline: Whole-workspace raw governance ingestion (simulating full scan)
	var fullScanContent string
	err = filepath.Walk(ws.MetadataDir, func(path string, info os.FileInfo, err error) error {
		if err != nil || info.IsDir() {
			return nil
		}
		if filepath.Ext(path) == ".md" {
			data, _ := os.ReadFile(path)
			fullScanContent += string(data) + "\n"
		}
		return nil
	})
	if err != nil {
		t.Fatalf("walk metadata directory failed: %v", err)
	}

	baselineTokens, err := tokenizer.Count(fullScanContent)
	if err != nil {
		t.Fatalf("tokenizer count on baseline failed: %v", err)
	}
	if baselineTokens < 2000 {
		t.Fatalf("expected substantial baseline token count, got %d", baselineTokens)
	}

	// 2. Measure Candidate: Compiled 3-layer Context Sandwich charter for M16/O1
	compiledCharter, err := charter.Compile(ws, "M16", "O1", []string{
		"D12-isolation-and-context-compilation",
		"D21-context-sandwich-execution-gates",
	})
	if err != nil {
		t.Fatalf("charter.Compile failed: %v", err)
	}

	charterTokens := compiledCharter.TokenCount
	if charterTokens > tokenizer.MaxTargetTokens {
		t.Errorf("charter exceeds 1200 target: %d tokens", charterTokens)
	}

	// 3. Compute Ingestion Reduction
	reductionPct := (float64(baselineTokens-charterTokens) / float64(baselineTokens)) * 100.0
	t.Logf("Baseline Full-Scan Tokens: %d", baselineTokens)
	t.Logf("Compiled Charter Tokens:   %d", charterTokens)
	t.Logf("Context Ingestion Savings: %.2f%%", reductionPct)

	// Threshold gate: must be at least 40% reduction
	if reductionPct < 40.0 {
		t.Fatalf("expected context reduction >= 40%%, got %.2f%%", reductionPct)
	}

	// 4. Invariant Verification: Ensure 0 regression in authority, claims, and perimeters
	if len(compiledCharter.Layer1.Claims) == 0 {
		t.Fatal("safety regression: completion claims missing from charter")
	}
	if len(compiledCharter.Layer3.WritesPaths) == 0 {
		t.Fatal("safety regression: target writable paths missing from perimeter")
	}
	if len(compiledCharter.Layer3.AllowedActions) == 0 {
		t.Fatal("safety regression: operator authority missing from perimeter")
	}
}
