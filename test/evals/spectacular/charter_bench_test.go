package spectaculareval

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/alexsmedile/spectacular/v2/internal/campaign"
	"github.com/alexsmedile/spectacular/v2/internal/charter"
	"github.com/alexsmedile/spectacular/v2/internal/charter/tokenizer"
	"github.com/alexsmedile/spectacular/v2/internal/discovery"
)

func TestContextCharter_PairedContextEconomyAndDecisionFidelity(t *testing.T) {
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

	// 2. Measure Candidate: Compiled 3-layer Context Sandwich charter for M16/O1 with explicit steering sources
	sources := []string{
		"D12-isolation-and-context-compilation",
		"D21-context-sandwich-execution-gates",
	}
	compiledCharter, err := charter.Compile(ws, "M16", "O1", sources)
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

	if reductionPct < 40.0 {
		t.Fatalf("expected context reduction >= 40%%, got %.2f%%", reductionPct)
	}

	// 4. Decision Fidelity Invariant: Ensure all requested steering decisions exist in Layer 2
	decisionMap := make(map[string]bool)
	for _, d := range compiledCharter.Layer2.Decisions {
		decisionMap[d.Ref] = true
		if d.Title == "" {
			t.Errorf("decision %s has empty title in charter layer 2", d.Ref)
		}
	}
	for _, s := range sources {
		if !decisionMap[s] {
			t.Fatalf("decision fidelity failure: requested steering source %s missing from Layer 2", s)
		}
	}

	// 5. Perimeter Fidelity Invariants: Writes, stops, authority must be intact
	if len(compiledCharter.Layer1.Claims) == 0 {
		t.Fatal("safety regression: completion claims missing from charter")
	}
	if len(compiledCharter.Layer3.WritesPaths) == 0 {
		t.Fatal("safety regression: target writable paths missing from perimeter")
	}
	if len(compiledCharter.Layer3.AllowedActions) == 0 {
		t.Fatal("safety regression: operator authority missing from perimeter")
	}
	if len(compiledCharter.Layer3.RequiresOwner) == 0 {
		t.Fatal("safety regression: requires_owner missing from perimeter")
	}
	if len(compiledCharter.Layer3.Stops) == 0 {
		t.Fatal("safety regression: stops missing from perimeter")
	}
}

func TestContextCharter_MissingSourceRefusal(t *testing.T) {
	ws, err := discovery.Open(".")
	if err != nil {
		t.Fatalf("discovery.Open failed: %v", err)
	}

	// Request a non-existent decision ref
	_, err = charter.Compile(ws, "M16", "O1", []string{"D999-non-existent-decision"})
	if err == nil {
		t.Fatal("expected charter.Compile to refuse missing source, but got nil error")
	}
	if !strings.Contains(err.Error(), "could not be resolved") && !strings.Contains(err.Error(), "not found") {
		t.Fatalf("expected missing source refusal, got: %v", err)
	}
}

func TestProgressiveCampaign_SketchesPassValidation(t *testing.T) {
	ws, err := discovery.Open(".")
	if err != nil {
		t.Fatalf("discovery.Open failed: %v", err)
	}

	// Campaign with progressive horizon detailing: active block is detailed, downstream blocks are sketches
	campaignPath := ".spectacular/campaigns/context-sandwich-steering.md"
	if _, err := os.Stat(campaignPath); err == nil {
		res, err := campaign.Validate(ws, campaignPath)
		if err != nil {
			t.Fatalf("progressive campaign validation failed: %v", err)
		}
		if len(res.Blocks) == 0 {
			t.Fatal("expected campaign to have blocks, got 0")
		}
	}
}
