package missionbundle

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"github.com/alexsmedile/spectacular/v2/internal/discovery"
	"github.com/alexsmedile/spectacular/v2/internal/domain"
	"github.com/alexsmedile/spectacular/v2/internal/workspace"
)

func (s Service) Show(ref string) (*Bundle, error) {
	bundle, err := Load(s.Workspace, ref)
	if err != nil {
		return nil, err
	}
	state := bundle.Derive()
	bundle.State = &state
	return bundle, nil
}

func (s Service) Check(ref string) (Check, error) {
	bundle, err := Load(s.Workspace, ref)
	if err != nil {
		return Check{}, err
	}
	return Validate(s.Workspace, bundle)
}

func (s Service) CheckWithVerify(ref string) (Check, error) {
	bundle, err := Load(s.Workspace, ref)
	if err != nil {
		return Check{}, err
	}
	check, err := Validate(s.Workspace, bundle)
	if err != nil {
		return check, err
	}

	root := ""
	if s.Workspace != nil {
		root = s.Workspace.Root
	}

	// 1. Run domain verification if configured or script exists
	testCmd := ""
	if s.Workspace != nil && s.Workspace.Config.Verification.Tier1Quick != "" {
		testCmd = s.Workspace.Config.Verification.Tier1Quick
	} else {
		checkScriptPath := filepath.Join(root, "tests", "check.sh")
		if stat, statErr := os.Stat(checkScriptPath); statErr == nil && !stat.IsDir() {
			testCmd = "sh tests/check.sh"
		}
	}

	if testCmd != "" {
		ctxVerify, cancelVerify := context.WithTimeout(context.Background(), 30*time.Second)
		cmd := exec.CommandContext(ctxVerify, "sh", "-c", testCmd)
		cmd.Dir = root
		out, runErr := cmd.CombinedOutput()
		cancelVerify()
		if ctxVerify.Err() == context.DeadlineExceeded {
			check.Valid = false
			check.Notices = append(check.Notices, "domain verification timed out after 30s")
			return check, nil
		}
		if runErr != nil {
			check.Valid = false
			check.Notices = append(check.Notices, fmt.Sprintf("domain verification failed: %s", strings.TrimSpace(string(out))))
			return check, nil
		}
		check.Checks = append(check.Checks, "domain-verification-pass")
	}

	// 2. Replay check if declared
	if bundle.Replay != nil && bundle.Replay.Command != "" {
		for _, p := range bundle.Replay.CachePaths {
			full := filepath.Join(root, p)
			_ = os.RemoveAll(full)
		}
		ctxReplay, cancelReplay := context.WithTimeout(context.Background(), 30*time.Second)
		replayCmd := exec.CommandContext(ctxReplay, "sh", "-c", bundle.Replay.Command)
		replayCmd.Dir = root
		out, runErr := replayCmd.CombinedOutput()
		cancelReplay()
		if ctxReplay.Err() == context.DeadlineExceeded {
			check.Valid = false
			check.Notices = append(check.Notices, "replay command timed out after 30s")
			return check, nil
		}
		if runErr != nil {
			check.Valid = false
			check.Notices = append(check.Notices, fmt.Sprintf("replay command failed: %s", strings.TrimSpace(string(out))))
			return check, nil
		}

		if testCmd != "" {
			ctxPost, cancelPost := context.WithTimeout(context.Background(), 30*time.Second)
			postCmd := exec.CommandContext(ctxPost, "sh", "-c", testCmd)
			postCmd.Dir = root
			postOut, postErr := postCmd.CombinedOutput()
			cancelPost()
			if ctxPost.Err() == context.DeadlineExceeded {
				check.Valid = false
				check.Notices = append(check.Notices, "post-replay domain verification timed out after 30s")
				return check, nil
			}
			if postErr != nil {
				check.Valid = false
				check.Notices = append(check.Notices, fmt.Sprintf("post-replay verification failed: %s", strings.TrimSpace(string(postOut))))
				return check, nil
			}
		}
		check.Checks = append(check.Checks, "replay-reconstruction-pass")
	}

	// 3. Git working tree cleanliness check
	ctxGit, cancelGit := context.WithTimeout(context.Background(), 3*time.Second)
	gitStatus := exec.CommandContext(ctxGit, "git", "status", "--porcelain")
	gitStatus.Dir = root
	out, errGit := gitStatus.Output()
	cancelGit()
	if ctxGit.Err() == context.DeadlineExceeded {
		check.Notices = append(check.Notices, "git status timed out after 3s")
	} else if errGit == nil {
		if len(bytes.TrimSpace(out)) == 0 {
			check.Checks = append(check.Checks, "git-working-tree-clean")
		} else {
			check.Notices = append(check.Notices, "git working tree contains untracked or uncommitted changes")
		}
	}

	// 4. Token efficiency derivation
	if s.TokenAnalyzer != nil {
		check.TokenEfficiency = s.TokenAnalyzer(s.Workspace, bundle)
	}

	return check, nil
}

func (s Service) ListMissions(statusFilter string, all bool) (MissionListResult, error) {
	if s.Workspace != nil {
		if fresh, err := discovery.Open(s.Workspace.Root); err == nil {
			s.Workspace = fresh
		}
	}
	var results []MissionSummary
	seen := map[string]bool{}
	for _, entry := range s.Workspace.Entries {
		if entry.Document.Record.Type != domain.Mission {
			continue
		}
		human, _, err := workspace.Ref(entry.Document)
		if err != nil || human == "" || seen[human] {
			continue
		}
		seen[human] = true
		bundle, err := Load(s.Workspace, human)
		if err != nil {
			continue
		}
		if statusFilter != "" && statusFilter != "all" {
			if bundle.Status != statusFilter {
				continue
			}
		} else if !all && statusFilter != "all" {
			// Active-first: exclude completed, resolved, superseded, cancelled
			switch bundle.Status {
			case "completed", "resolved", "superseded", "cancelled":
				continue
			}
		}
		derived := bundle.Derive()
		title := bundle.Title
		if title == "" && entry.Document.Record.Title != nil {
			title = *entry.Document.Record.Title
		}
		results = append(results, MissionSummary{
			Ref:    bundle.Ref,
			Title:  title,
			Status: bundle.Status,
			Holder: derived.Holder,
			Next:   derived.Next,
			Path:   bundle.Path,
		})
	}
	sort.Slice(results, func(i, j int) bool {
		return results[i].Ref < results[j].Ref
	})
	return MissionListResult{
		SchemaVersion: "spectacular.mission.list.v2",
		Missions:      results,
	}, nil
}

type MissionSummary struct {
	Ref    string `json:"ref"`
	Title  string `json:"title"`
	Status string `json:"status"`
	Holder string `json:"holder"`
	Next   string `json:"next"`
	Path   string `json:"path"`
}

type MissionListResult struct {
	SchemaVersion string           `json:"schema_version"`
	Missions      []MissionSummary `json:"missions"`
}
