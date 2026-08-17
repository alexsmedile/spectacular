package missionbundle

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/alexsmedile/spectacular/v2/internal/discovery"
	"github.com/alexsmedile/spectacular/v2/internal/domain"
)

func TestMissionOrderResolutionAcrossRepository(t *testing.T) {
	root, err := filepath.Abs(filepath.Join("..", ".."))
	if err != nil {
		t.Fatal(err)
	}
	ws, err := discovery.Open(root)
	if err != nil {
		t.Fatal(err)
	}

	t.Run("real repository missions M2 through M8 resolve as one order despite mixed ref and human_ref spellings", func(t *testing.T) {
		refs := []string{"M2", "M3", "M4", "M5", "M6", "M7", "M8"}
		bundles := make([]*Bundle, 0, len(refs))

		for _, ref := range refs {
			b, err := Load(ws, ref)
			if err != nil {
				t.Fatalf("failed to load %s: %v", ref, err)
			}
			bundles = append(bundles, b)
		}

		if len(bundles) != 7 {
			t.Fatalf("expected 7 bundles, got %d", len(bundles))
		}

		// Verify every ref resolves through both human and typed spelling where available
		for i, ref := range refs {
			entry, err := ws.Lookup(ref, domain.Mission)
			if err != nil {
				t.Fatalf("lookup by human ref %s failed: %v", ref, err)
			}
			typedRef := "Mission:" + entry.Document.Record.ID.String()
			typedEntry, err := ws.Lookup(typedRef, domain.Mission)
			if err != nil {
				t.Fatalf("lookup by typed ref %s failed: %v", typedRef, err)
			}
			if entry.Document.Record.ID != typedEntry.Document.Record.ID {
				t.Fatalf("mismatch between human and typed lookup for %s", ref)
			}
			if i > 0 {
				// Each mission follows the sequence
				prev := bundles[i-1]
				if prev.Status != "completed" && prev.Ref != "M8" {
					t.Fatalf("predecessor mission %s status=%s, want completed", prev.Ref, prev.Status)
				}
			}
		}
	})
}

func TestMissionOrderNegativeValidators(t *testing.T) {
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

	t.Run("refusal on dangling Mission ref", func(t *testing.T) {
		candidate := cloneBundle(t, base)
		candidate.AfterMission = []string{"M99-nonexistent"}
		err := validateMissionOrderIntegrity(ws, candidate)
		if err == nil {
			t.Fatal("expected refusal on dangling Mission ref, got nil")
		}
		var refusal *domain.Refusal
		if !errors.As(err, &refusal) {
			t.Fatalf("expected *domain.Refusal, got %v", err)
		}
		if refusal.Field != "after_mission" {
			t.Fatalf("refusal.Field=%s, want after_mission", refusal.Field)
		}
		if !strings.Contains(refusal.Detail, "M99-nonexistent") {
			t.Fatalf("refusal detail %q does not name missing ref", refusal.Detail)
		}
	})

	t.Run("refusal on self dependency", func(t *testing.T) {
		candidate := cloneBundle(t, base)
		candidate.AfterMission = []string{"M6"}
		err := validateMissionOrderIntegrity(ws, candidate)
		if err == nil {
			t.Fatal("expected refusal on self dependency, got nil")
		}
		var refusal *domain.Refusal
		if !errors.As(err, &refusal) {
			t.Fatalf("expected *domain.Refusal, got %v", err)
		}
		if refusal.Field != "after_mission" {
			t.Fatalf("refusal.Field=%s, want after_mission", refusal.Field)
		}
	})

	t.Run("refusal on Mission-level cycle", func(t *testing.T) {
		fixtureRoot := missionServiceFixture(t)
		plan1, _ := stressPlan()
		plan1.Title = "Cycle Mission One"
		svc := openMissionService(t, fixtureRoot)
		m1, err := svc.Start(plan1, []byte("plan 1"))
		if err != nil {
			t.Fatal(err)
		}

		plan2, _ := stressPlan()
		plan2.Title = "Cycle Mission Two"
		svc = openMissionService(t, fixtureRoot)
		m2, err := svc.Start(plan2, []byte("plan 2"))
		if err != nil {
			t.Fatal(err)
		}

		// Update m2's MISSION.md on disk to depend on m1
		m2Path := m2.Path
		if !filepath.IsAbs(m2Path) {
			m2Path = filepath.Join(fixtureRoot, m2Path)
		}
		m2Content, err := os.ReadFile(m2Path)
		if err != nil {
			t.Fatal(err)
		}
		updatedM2 := strings.Replace(string(m2Content), "\n---\n", fmt.Sprintf("\nafter_mission:\n  - %s\n---\n", m1.Ref), 1)
		if err := os.WriteFile(m2Path, []byte(updatedM2), 0644); err != nil {
			t.Fatal(err)
		}

		fixtureWs, err := discovery.Open(fixtureRoot)
		if err != nil {
			t.Fatal(err)
		}

		bundle1, err := Load(fixtureWs, m1.Ref)
		if err != nil {
			t.Fatal(err)
		}
		// Create cycle: m1 depends on m2, and m2 depends on m1
		bundle1.AfterMission = []string{m2.Ref}

		err = validateMissionOrderIntegrity(fixtureWs, bundle1)
		if err == nil {
			t.Fatal("expected refusal on Mission cycle, got nil")
		}
		var refusal *domain.Refusal
		if !errors.As(err, &refusal) {
			t.Fatalf("expected *domain.Refusal, got %v", err)
		}
		if refusal.Field != "after_mission" {
			t.Fatalf("refusal.Field=%s, want after_mission", refusal.Field)
		}
		if !strings.Contains(refusal.Detail, "acyclic") {
			t.Fatalf("refusal.Detail=%q, want mentioning acyclic", refusal.Detail)
		}
	})

	t.Run("refusal on activation ahead of incomplete predecessor", func(t *testing.T) {
		candidate := cloneBundle(t, base)
		candidate.Status = "active"
		candidate.AfterMission = []string{"M8"} // M8 is currently active, not completed!
		err := validateMissionOrderActivation(ws, candidate)
		if err == nil {
			t.Fatal("expected refusal on activation ahead of active/incomplete predecessor, got nil")
		}
		var refusal *domain.Refusal
		if !errors.As(err, &refusal) {
			t.Fatalf("expected *domain.Refusal, got %v", err)
		}
		if refusal.Field != "after_mission" {
			t.Fatalf("refusal.Field=%s, want after_mission", refusal.Field)
		}
		if !strings.Contains(refusal.Detail, "M8") || !strings.Contains(refusal.Detail, "not completed") {
			t.Fatalf("refusal detail %q does not state predecessor is not completed", refusal.Detail)
		}
	})

	t.Run("unactivated mission (defined) passes activation validator even with incomplete predecessor", func(t *testing.T) {
		candidate := cloneBundle(t, base)
		candidate.Status = "defined"
		candidate.AfterMission = []string{"M8"}
		if err := validateMissionOrderActivation(ws, candidate); err != nil {
			t.Fatalf("unactivated mission should pass activation check: %v", err)
		}
	})

	t.Run("activation with completed predecessor passes cleanly", func(t *testing.T) {
		candidate := cloneBundle(t, base)
		candidate.Status = "active"
		candidate.AfterMission = []string{"M5"} // M5 is completed
		if err := validateMissionOrderActivation(ws, candidate); err != nil {
			t.Fatalf("active mission with completed predecessor should pass: %v", err)
		}
	})
}
