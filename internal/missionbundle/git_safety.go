package missionbundle

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/alexsmedile/spectacular/v2/internal/domain"
)

// resolveGitDir locates the effective .git directory for a repository root,
// resolving linked worktree pointers when .git is a file.
func resolveGitDir(root string) (string, error) {
	gitPath := filepath.Join(root, ".git")
	info, err := os.Stat(gitPath)
	if err != nil {
		return "", err
	}
	if info.IsDir() {
		return gitPath, nil
	}
	// Linked worktree .git file containing "gitdir: <path>"
	data, err := os.ReadFile(gitPath)
	if err != nil {
		return "", err
	}
	content := strings.TrimSpace(string(data))
	if strings.HasPrefix(content, "gitdir:") {
		dir := strings.TrimSpace(strings.TrimPrefix(content, "gitdir:"))
		if !filepath.IsAbs(dir) {
			dir = filepath.Join(root, dir)
		}
		return filepath.Clean(dir), nil
	}
	return gitPath, nil
}

// CheckPassiveGitState verifies that the repository working tree is not in the middle
// of an unresolved merge, rebase, or cherry-pick conflict.
// It resolves linked worktrees and never mutates Git state or deletes worktrees.
func CheckPassiveGitState(root string) error {
	gitDir, err := resolveGitDir(root)
	if err != nil {
		// Not a git repo or inaccessible; do not block non-git tests
		return nil
	}

	// 1. Merge conflict in progress
	if _, err := os.Stat(filepath.Join(gitDir, "MERGE_HEAD")); err == nil {
		return domain.NewRefusal(domain.RefusalCollision, "git", "active merge conflict in progress; resolve git merge before continuing", nil)
	}

	// 2. Cherry-pick in progress
	if _, err := os.Stat(filepath.Join(gitDir, "CHERRY_PICK_HEAD")); err == nil {
		return domain.NewRefusal(domain.RefusalCollision, "git", "active cherry-pick in progress; resolve git cherry-pick before continuing", nil)
	}

	// 3. Rebase in progress (apply or merge)
	if _, err := os.Stat(filepath.Join(gitDir, "rebase-apply")); err == nil {
		return domain.NewRefusal(domain.RefusalCollision, "git", "active rebase in progress; resolve git rebase before continuing", nil)
	}
	if _, err := os.Stat(filepath.Join(gitDir, "rebase-merge")); err == nil {
		return domain.NewRefusal(domain.RefusalCollision, "git", "active rebase in progress; resolve git rebase before continuing", nil)
	}

	return nil
}
