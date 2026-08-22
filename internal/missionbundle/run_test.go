package missionbundle

import (
	"strings"
	"testing"
)

func TestValidateTransition(t *testing.T) {
	// Legal transitions
	legals := [][2]string{
		{"active", "paused"},
		{"active", "blocked"},
		{"active", "awaiting-review"},
		{"active", "completed"},
		{"active", "stopped"},
		{"paused", "active"},
		{"paused", "blocked"},
		{"paused", "stopped"},
		{"blocked", "active"},
		{"blocked", "stopped"},
		{"awaiting-review", "active"},
		{"awaiting-review", "completed"},
		{"awaiting-review", "stopped"},
	}

	for _, pair := range legals {
		if err := ValidateTransition(pair[0], pair[1]); err != nil {
			t.Errorf("expected transition from %s to %s to be legal, got: %v", pair[0], pair[1], err)
		}
	}

	// Illegal transitions
	illegals := [][2]string{
		{"completed", "active"},
		{"completed", "paused"},
		{"stopped", "active"},
		{"stopped", "completed"},
		{"paused", "completed"},
		{"blocked", "completed"},
	}

	for _, pair := range illegals {
		if err := ValidateTransition(pair[0], pair[1]); err == nil {
			t.Errorf("expected transition from %s to %s to be illegal, but succeeded", pair[0], pair[1])
		}
	}
}

func TestTransitionRun_Validation(t *testing.T) {
	s := Service{}
	_, err := s.TransitionRun("M18/R1", "paused", "", "reason", "")
	if err == nil || !strings.Contains(err.Error(), "actor identity is required") {
		t.Fatalf("expected missing actor error, got: %v", err)
	}

	_, err = s.TransitionRun("M18/R1", "paused", "Alex", "", "")
	if err == nil || !strings.Contains(err.Error(), "transition reason is required") {
		t.Fatalf("expected missing reason error, got: %v", err)
	}
}
