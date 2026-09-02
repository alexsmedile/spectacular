package worktree

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

// Info describes an active git worktree.
type Info struct {
	Slug   string `json:"slug"`
	Path   string `json:"path"`
	Branch string `json:"branch"`
	Clean  bool   `json:"clean"`
}

// Provision safely creates an isolated linked git worktree under .worktrees/<slug>.
func Provision(repoRoot string, slug string, branch string, baseCommit string) (string, error) {
	if slug == "" {
		return "", fmt.Errorf("worktree slug cannot be empty")
	}
	if branch == "" {
		return "", fmt.Errorf("worktree branch cannot be empty")
	}

	worktreesDir := filepath.Join(repoRoot, ".worktrees")
	if err := os.MkdirAll(worktreesDir, 0755); err != nil {
		return "", fmt.Errorf("failed to create .worktrees directory: %w", err)
	}

	targetPath := filepath.Join(worktreesDir, slug)
	if _, err := os.Stat(targetPath); err == nil {
		return "", fmt.Errorf("worktree already exists at path: %s", targetPath)
	}

	// Prepare git worktree add command
	args := []string{"worktree", "add"}
	if baseCommit != "" {
		args = append(args, "-b", branch, targetPath, baseCommit)
	} else {
		args = append(args, "-b", branch, targetPath)
	}

	cmd := exec.Command("git", args...)
	cmd.Dir = repoRoot
	var stderr bytes.Buffer
	cmd.Stderr = &stderr

	if err := cmd.Run(); err != nil {
		// If branch already exists, try attaching directly without -b
		if strings.Contains(stderr.String(), "already exists") {
			fallbackCmd := exec.Command("git", "worktree", "add", targetPath, branch)
			fallbackCmd.Dir = repoRoot
			fallbackCmd.Stderr = &stderr
			if fallbackErr := fallbackCmd.Run(); fallbackErr != nil {
				return "", fmt.Errorf("git worktree add failed: %s (%w)", strings.TrimSpace(stderr.String()), fallbackErr)
			}
			return targetPath, nil
		}
		return "", fmt.Errorf("git worktree add failed: %s (%w)", strings.TrimSpace(stderr.String()), err)
	}

	return targetPath, nil
}

// IsDirty checks whether the worktree at the given path has uncommitted changes.
func IsDirty(worktreePath string) (bool, error) {
	cmd := exec.Command("git", "status", "--porcelain")
	cmd.Dir = worktreePath
	var stdout bytes.Buffer
	cmd.Stdout = &stdout

	if err := cmd.Run(); err != nil {
		return false, fmt.Errorf("git status check failed: %w", err)
	}

	return strings.TrimSpace(stdout.String()) != "", nil
}

// Prune conservatively removes a worktree after confirming it is clean or if force is true.
func Prune(repoRoot string, slug string, force bool) error {
	targetPath := filepath.Join(repoRoot, ".worktrees", slug)
	if _, err := os.Stat(targetPath); os.IsNotExist(err) {
		return nil
	}

	if !force {
		dirty, err := IsDirty(targetPath)
		if err != nil {
			return fmt.Errorf("failed to check worktree status: %w", err)
		}
		if dirty {
			return fmt.Errorf("refusing to prune dirty worktree at %s (uncommitted changes present)", targetPath)
		}
	}

	args := []string{"worktree", "remove"}
	if force {
		args = append(args, "--force")
	}
	args = append(args, targetPath)

	cmd := exec.Command("git", args...)
	cmd.Dir = repoRoot
	var stderr bytes.Buffer
	cmd.Stderr = &stderr

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("git worktree remove failed: %s (%w)", strings.TrimSpace(stderr.String()), err)
	}

	return nil
}
