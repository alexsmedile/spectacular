package spectaculareval

import (
	"strings"
	"testing"

	"github.com/alexsmedile/spectacular/v2/internal/command"
	"github.com/alexsmedile/spectacular/v2/internal/domain"
	"github.com/alexsmedile/spectacular/v2/internal/missionbundle"
)

func TestM21_ScopeHardeningAndFinalCampaignRegression(t *testing.T) {
	t.Run("claim: discriminating-guardrail-evaluation", func(t *testing.T) {
		// Valid relative paths pass
		valid := []string{"cmd/spectacular/main.go", "internal/charter/", "docs/architecture.md"}
		if err := missionbundle.ValidateWritePaths(valid); err != nil {
			t.Fatalf("expected valid paths to pass, got: %v", err)
		}

		// Traversal escape refuses
		if err := missionbundle.ValidateWritePaths([]string{"../secret.txt"}); err == nil || !domain.RefusalHasCode(err, domain.RefusalPathEscape) {
			t.Fatalf("expected RefusalPathEscape for parent traversal, got: %v", err)
		}

		// Glob wildcard refuses
		if err := missionbundle.ValidateWritePaths([]string{"internal/*.go"}); err == nil || !domain.RefusalHasCode(err, domain.RefusalInvalidScope) {
			t.Fatalf("expected RefusalInvalidScope for wildcard, got: %v", err)
		}
	})

	t.Run("claim: coherent-broad-perimeters", func(t *testing.T) {
		// Broad directory perimeters (e.g. 10 files inside a module) are legal without arbitrary numeric ceilings
		broadDirectory := []string{"internal/missionbundle/", "cmd/spectacular/"}
		if err := missionbundle.ValidateWritePaths(broadDirectory); err != nil {
			t.Fatalf("expected broad directory reservation to be legal, got: %v", err)
		}

		// Disjoint broad directories do not collide
		if missionbundle.PathsOverlap("internal/missionbundle/", "cmd/spectacular/") {
			t.Fatalf("expected disjoint broad directories to not overlap")
		}
	})

	t.Run("claim: final-campaign-regression", func(t *testing.T) {
		// 18 commands authorized
		if len(command.Registry) != 18 {
			t.Fatalf("expected 18 commands, got %d", len(command.Registry))
		}

		// 23 mandatory validators registered
		if len(missionbundle.ValidatorNames()) != 23 {
			t.Fatalf("expected 23 validators, got %d", len(missionbundle.ValidatorNames()))
		}

		// Verify key commands present
		requiredCommands := []string{
			"mission start", "mission show", "mission check",
			"objective show", "objective promote", "objective finish",
			"run show", "run start", "run transition",
			"review record", "handoff record", "evidence record",
			"mission complete", "proposal check", "campaign check",
			"contract amend", "charter", "decide",
		}
		for i, req := range requiredCommands {
			got := strings.Join(command.Registry[i].Words, " ")
			if got != req {
				t.Fatalf("command index %d: got %q, want %q", i, got, req)
			}
		}
	})
}
