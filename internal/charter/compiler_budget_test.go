package charter

import (
	"strings"
	"testing"

	"github.com/alexsmedile/spectacular/v2/internal/discovery"
)

// The threshold helper has its own table test. This checks the integration
// boundary: Compile must turn an over-cap charter into a refusal, after its
// one allowed safe-compaction attempt, rather than returning unsafe output.
func TestCompileRefusesWhenConfiguredHardCapIsExceeded(t *testing.T) {
	ws, err := discovery.Open(".")
	if err != nil {
		t.Fatal(err)
	}
	ws.Config.TokenBudgets.Charter.Target = 1
	ws.Config.TokenBudgets.Charter.HardCap = 1

	_, err = Compile(ws, "M18", "O1", nil)
	if err == nil {
		t.Fatal("Compile succeeded despite exceeding the configured hard cap")
	}
	if !strings.Contains(err.Error(), "exceeds hard refusal ceiling") {
		t.Fatalf("hard-cap refusal = %v", err)
	}
}
