package spectaculareval

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/alexsmedile/spectacular/v2/internal/domain"
	"github.com/alexsmedile/spectacular/v2/internal/missionbundle"
)

func TestM19_DisjointWriteReservationsAndDependencyLocks(t *testing.T) {
	t.Run("claim: disjoint-write-reservations", func(t *testing.T) {
		validPaths := []string{"cmd/spectacular/main.go", "internal/charter/", "docs/architecture.md"}
		if err := missionbundle.ValidateWritePaths(validPaths); err != nil {
			t.Fatalf("expected valid paths to pass, got: %v", err)
		}

		if err := missionbundle.ValidateWritePaths([]string{"../escape.go"}); err == nil || !domain.RefusalHasCode(err, domain.RefusalPathEscape) {
			t.Fatalf("expected RefusalPathEscape for parent traversal, got: %v", err)
		}

		if err := missionbundle.ValidateWritePaths([]string{"cmd/*.go"}); err == nil || !domain.RefusalHasCode(err, domain.RefusalInvalidScope) {
			t.Fatalf("expected RefusalInvalidScope for wildcard, got: %v", err)
		}

		// Overlap tests
		if !missionbundle.PathsOverlap("internal/charter", "internal/charter/tokenizer") {
			t.Errorf("expected parent directory overlap to be true")
		}
		if !missionbundle.PathsOverlap("cmd/spectacular/main.go", "cmd/spectacular/main.go") {
			t.Errorf("expected exact file overlap to be true")
		}
		if missionbundle.PathsOverlap("internal/charter", "internal/discovery") {
			t.Errorf("expected disjoint directories to be false")
		}
	})

	t.Run("claim: passive-git-state-inspection", func(t *testing.T) {
		tempDir := t.TempDir()
		gitDir := filepath.Join(tempDir, ".git")
		if err := os.MkdirAll(gitDir, 0o755); err != nil {
			t.Fatal(err)
		}

		// Clean passes
		if err := missionbundle.CheckPassiveGitState(tempDir); err != nil {
			t.Fatalf("expected clean git state to pass, got: %v", err)
		}

		// MERGE_HEAD refuses
		mergeHead := filepath.Join(gitDir, "MERGE_HEAD")
		if err := os.WriteFile(mergeHead, []byte("commit1234\n"), 0o644); err != nil {
			t.Fatal(err)
		}
		err := missionbundle.CheckPassiveGitState(tempDir)
		if err == nil || !domain.RefusalHasCode(err, domain.RefusalCollision) {
			t.Fatalf("expected RefusalCollision for MERGE_HEAD, got: %v", err)
		}
		_ = os.Remove(mergeHead)
	})
}
