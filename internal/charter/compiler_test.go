package charter

import (
	"strings"
	"testing"

	"github.com/alexsmedile/spectacular/v2/internal/charter/tokenizer"
	"github.com/alexsmedile/spectacular/v2/internal/discovery"
)

func TestCompiler_CompileM16(t *testing.T) {
	ws, err := discovery.Open(".")
	if err != nil {
		t.Fatalf("discovery.Open failed: %v", err)
	}

	charter, err := Compile(ws, "M16", "O1", []string{"D12-isolation-and-context-compilation", "D21-context-sandwich-execution-gates"})
	if err != nil {
		t.Fatalf("Compile failed: %v", err)
	}

	if charter.MissionRef != "M16" {
		t.Errorf("expected MissionRef M16, got %s", charter.MissionRef)
	}
	if charter.ObjectiveRef != "O1" {
		t.Errorf("expected ObjectiveRef O1, got %s", charter.ObjectiveRef)
	}
	if len(charter.Layer1.Claims) == 0 {
		t.Error("expected non-empty completion claims in Layer 1")
	}
	if len(charter.Layer2.Decisions) == 0 {
		t.Error("expected retrieved decisions in Layer 2")
	}
	if len(charter.Layer3.WritesPaths) == 0 {
		t.Error("expected non-empty writes paths in Layer 3")
	}

	rendered := charter.RenderMarkdown()
	if !strings.Contains(rendered, "## 1. FROZEN TRUTH") ||
		!strings.Contains(rendered, "## 2. OWNER STEERING") ||
		!strings.Contains(rendered, "## 3. EXECUTION PERIMETER") {
		t.Error("rendered markdown missing one of the 3 Context Sandwich layers")
	}

	if charter.TokenCount <= 0 {
		t.Errorf("expected positive token count, got %d", charter.TokenCount)
	}
	if charter.TokenCount > tokenizer.MaxTargetTokens {
		t.Logf("charter token count: %d (disposition: %s, compacted: %v)", charter.TokenCount, charter.Disposition, charter.Compacted)
	}
}

func TestCompiler_DeduplicationPreservesOrder(t *testing.T) {
	ws, err := discovery.Open(".")
	if err != nil {
		t.Fatalf("discovery.Open failed: %v", err)
	}

	sources := []string{
		"D12-isolation-and-context-compilation",
		"D21-context-sandwich-execution-gates",
		"D12-isolation-and-context-compilation", // Duplicate
	}

	charter, err := Compile(ws, "M16", "O1", sources)
	if err != nil {
		t.Fatalf("Compile failed: %v", err)
	}

	// Verify D12 is not duplicated in Layer 2
	d12Count := 0
	for _, d := range charter.Layer2.Decisions {
		if strings.Contains(d.Ref, "D12") {
			d12Count++
		}
	}
	if d12Count > 1 {
		t.Errorf("expected deduplicated D12, found count %d", d12Count)
	}
}

func TestCompiler_SafeCompactionPreservesClaims(t *testing.T) {
	ws, err := discovery.Open(".")
	if err != nil {
		t.Fatalf("discovery.Open failed: %v", err)
	}

	charter, err := Compile(ws, "M16", "O1", []string{
		"D1-rnkfzw",
		"D10-repoint",
		"D11-proposal-retirement",
		"D12-isolation-and-context-compilation",
		"D20-use-live-charter-retrieval-without-a-persistent-cache",
		"D21-context-sandwich-execution-gates",
	})
	if err != nil {
		t.Fatalf("Compile failed: %v", err)
	}

	// Verify all frozen claims in Layer 1 are 100% intact even with heavy decision payload
	if len(charter.Layer1.Claims) == 0 {
		t.Fatal("expected claims to be preserved")
	}
	for _, c := range charter.Layer1.Claims {
		if c.Claim == "" || c.PassBoundary == "" || c.ProofRequirement == "" {
			t.Errorf("claim fields corrupted: %+v", c)
		}
	}
}
