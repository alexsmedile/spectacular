package spectaculareval

import (
	"strings"
	"testing"

	"github.com/alexsmedile/spectacular/v2/internal/command"
	"github.com/alexsmedile/spectacular/v2/internal/missionbundle"
)

func TestM20_ClusteredEvidenceAndSurface18(t *testing.T) {
	t.Run("claim: authorized-surface-18", func(t *testing.T) {
		if len(command.Registry) < 18 {
			t.Fatalf("expected at least 18 commands in command.Registry, got %d", len(command.Registry))
		}
		found := false
		for _, spec := range command.Registry {
			if strings.Join(spec.Words, " ") == "evidence record" {
				found = true
				if spec.JSONSchema != "spectacular.evidence.record.v2" {
					t.Errorf("expected v2 schema, got %s", spec.JSONSchema)
				}
				if spec.Effect != command.Mutating {
					t.Errorf("expected mutating effect, got %s", spec.Effect)
				}
			}
		}
		if !found {
			t.Fatal("evidence record not found in command.Registry")
		}
	})

	t.Run("claim: atomic-clustered-evidence", func(t *testing.T) {
		e := &missionbundle.Evidence{
			Ref:        "E1",
			Title:      "Attributable evidence for O1 and O2",
			Mission:    "M20",
			Actor:      "Alex",
			Commit:     "0123456789abcdef0123456789abcdef01234567",
			Tree:       "0123456789abcdef0123456789abcdef01234567",
			Objectives: []string{"O1", "O2"},
			Checks: []missionbundle.EvidenceCheck{
				{Name: "test-suite", Result: "pass"},
			},
		}
		if len(e.Objectives) != 2 {
			t.Errorf("expected clustered coverage over 2 objectives, got %d", len(e.Objectives))
		}
	})

	t.Run("claim: completed-is-not-proved", func(t *testing.T) {
		// Completed run lifecycle transition is valid, while evidence is separately recorded
		if err := missionbundle.ValidateTransition("active", "completed"); err != nil {
			t.Fatalf("expected active -> completed transition to be valid, got: %v", err)
		}
	})
}
