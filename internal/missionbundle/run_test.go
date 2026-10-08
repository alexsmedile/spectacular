package missionbundle

import (
	"strings"
	"testing"
)

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
