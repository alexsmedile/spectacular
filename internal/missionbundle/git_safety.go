package missionbundle

import (
	"os"
	"path/filepath"

	"github.com/alexsmedile/spectacular/v2/internal/domain"
)

// CheckPassiveGitState verifies that the repository working tree is not in the middle
// of an unresolved merge, rebase, or cherry-pick conflict.
// It never mutates the repository, creates branches, stashes changes, or deletes worktrees.
func CheckPassiveGitState(root string) error {
	gitDir := filepath.Join(root, ".git")
	info, err := os.Stat(gitDir)
	if err != nil {
		// Not a git repo or inaccessible; do not block non-git tests
		return nil
	}
	if !info.IsDir() {
		// Possibly a worktree pointer file .git -> resolve directory
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
