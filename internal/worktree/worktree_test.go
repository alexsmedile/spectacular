package worktree

import (
	"os"
	"os/exec"
	"path/filepath"
	"testing"
)

func createTestGitRepo(t *testing.T) string {
	t.Helper()
	dir := t.TempDir()

	cmd := exec.Command("git", "init")
	cmd.Dir = dir
	if err := cmd.Run(); err != nil {
		t.Fatalf("git init failed: %v", err)
	}

	cmd = exec.Command("git", "config", "user.name", "Test User")
	cmd.Dir = dir
	_ = cmd.Run()

	cmd = exec.Command("git", "config", "user.email", "test@example.com")
	cmd.Dir = dir
	_ = cmd.Run()

	readme := filepath.Join(dir, "README.md")
	if err := os.WriteFile(readme, []byte("# Test Repo\n"), 0644); err != nil {
		t.Fatalf("failed to write readme: %v", err)
	}

	cmd = exec.Command("git", "add", "README.md")
	cmd.Dir = dir
	if err := cmd.Run(); err != nil {
		t.Fatalf("git add failed: %v", err)
	}

	cmd = exec.Command("git", "commit", "-m", "Initial commit")
	cmd.Dir = dir
	if err := cmd.Run(); err != nil {
		t.Fatalf("git commit failed: %v", err)
	}

	return dir
}

func TestWorktree_ProvisionAndPrune(t *testing.T) {
	repoDir := createTestGitRepo(t)

	slug := "test-worker"
	branch := "feature/test-worker"

	path, err := Provision(repoDir, slug, branch, "")
	if err != nil {
		t.Fatalf("Provision failed: %v", err)
	}

	expectedPath := filepath.Join(repoDir, ".worktrees", slug)
	if path != expectedPath {
		t.Errorf("expected path %s, got %s", expectedPath, path)
	}

	// Verify path exists
	if _, err := os.Stat(path); os.IsNotExist(err) {
		t.Fatalf("expected worktree path %s to exist", path)
	}

	// Verify not dirty
	dirty, err := IsDirty(path)
	if err != nil {
		t.Fatalf("IsDirty failed: %v", err)
	}
	if dirty {
		t.Error("expected fresh worktree to be clean")
	}

	// Create dirty file
	dirtyFile := filepath.Join(path, "dirty.txt")
	if err := os.WriteFile(dirtyFile, []byte("dirty content"), 0644); err != nil {
		t.Fatalf("failed to write dirty file: %v", err)
	}

	dirty, err = IsDirty(path)
	if err != nil {
		t.Fatalf("IsDirty failed: %v", err)
	}
	if !dirty {
		t.Error("expected modified worktree to be dirty")
	}

	// Conservative prune should fail on dirty worktree
	if err := Prune(repoDir, slug, false); err == nil {
		t.Error("expected conservative Prune to fail on dirty worktree")
	}

	// Force prune should succeed
	if err := Prune(repoDir, slug, true); err != nil {
		t.Fatalf("force Prune failed: %v", err)
	}

	// Verify path no longer exists
	if _, err := os.Stat(path); !os.IsNotExist(err) {
		t.Errorf("expected worktree path %s to be removed", path)
	}
}
