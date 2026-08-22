package missionbundle

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/alexsmedile/spectacular/v2/internal/domain"
)

func TestCheckPassiveGitState(t *testing.T) {
	root := t.TempDir()
	gitDir := filepath.Join(root, ".git")
	if err := os.MkdirAll(gitDir, 0o755); err != nil {
		t.Fatal(err)
	}

	// 1. Clean git state passes
	if err := CheckPassiveGitState(root); err != nil {
		t.Fatalf("expected clean git state to pass, got: %v", err)
	}

	// 2. Active merge conflict refuses
	mergeHead := filepath.Join(gitDir, "MERGE_HEAD")
	if err := os.WriteFile(mergeHead, []byte("commit1234\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	err := CheckPassiveGitState(root)
	if err == nil || !domain.RefusalHasCode(err, domain.RefusalCollision) {
		t.Fatalf("expected RefusalCollision, got: %v", err)
	}
	_ = os.Remove(mergeHead)

	// 3. Active rebase-merge refuses
	rebaseMerge := filepath.Join(gitDir, "rebase-merge")
	if err := os.MkdirAll(rebaseMerge, 0o755); err != nil {
		t.Fatal(err)
	}
	err = CheckPassiveGitState(root)
	if err == nil || !domain.RefusalHasCode(err, domain.RefusalCollision) {
		t.Fatalf("expected RefusalCollision for rebase-merge, got: %v", err)
	}
	_ = os.RemoveAll(rebaseMerge)

	// 4. Active cherry-pick refuses
	cherryHead := filepath.Join(gitDir, "CHERRY_PICK_HEAD")
	if err := os.WriteFile(cherryHead, []byte("commit1234\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	err = CheckPassiveGitState(root)
	if err == nil || !domain.RefusalHasCode(err, domain.RefusalCollision) {
		t.Fatalf("expected RefusalCollision for CHERRY_PICK_HEAD, got: %v", err)
	}
}
