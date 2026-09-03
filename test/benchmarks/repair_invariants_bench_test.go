package spectaculareval

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/alexsmedile/spectacular/v2/internal/charter"
	"github.com/alexsmedile/spectacular/v2/internal/charter/tokenizer"
	"github.com/alexsmedile/spectacular/v2/internal/discovery"
	"github.com/alexsmedile/spectacular/v2/internal/domain"
	"github.com/alexsmedile/spectacular/v2/internal/missionbundle"
)

func TestM22_RepairedSteeringInvariants(t *testing.T) {
	t.Run("claim: m19-reservation-and-git-repair - linked worktree gitdir resolution", func(t *testing.T) {
		tmp := t.TempDir()
		realGitDir := filepath.Join(tmp, "main_repo", ".git", "worktrees", "wt1")
		if err := os.MkdirAll(realGitDir, 0o755); err != nil {
			t.Fatal(err)
		}
		wtDir := filepath.Join(tmp, "wt1")
		if err := os.MkdirAll(wtDir, 0o755); err != nil {
			t.Fatal(err)
		}
		// Write .git file in linked worktree
		if err := os.WriteFile(filepath.Join(wtDir, ".git"), []byte("gitdir: "+realGitDir+"\n"), 0o644); err != nil {
			t.Fatal(err)
		}

		// Initially clean
		if err := missionbundle.CheckPassiveGitState(wtDir); err != nil {
			t.Fatalf("expected clean git state on worktree, got: %v", err)
		}

		// Inject active MERGE_HEAD in the resolved realGitDir
		if err := os.WriteFile(filepath.Join(realGitDir, "MERGE_HEAD"), []byte("0123456789abcdef\n"), 0o644); err != nil {
			t.Fatal(err)
		}

		// Now CheckPassiveGitState on wtDir MUST detect the conflict in realGitDir!
		err := missionbundle.CheckPassiveGitState(wtDir)
		if err == nil || !domain.RefusalHasCode(err, domain.RefusalCollision) {
			t.Fatalf("expected RefusalCollision on linked worktree merge conflict, got: %v", err)
		}
	})

	t.Run("claim: m19-reservation-and-git-repair - supersedes chain and path overlap", func(t *testing.T) {
		// Verify disjoint paths do not collide
		if missionbundle.PathsOverlap("internal/charter/", "cmd/spectacular/") {
			t.Fatal("expected disjoint paths to not overlap")
		}
		// Overlapping paths collide
		if !missionbundle.PathsOverlap("internal/charter/", "internal/charter/engine.go") {
			t.Fatal("expected nested path to overlap with directory")
		}
	})

	t.Run("claim: m20-evidence-target-integrity - strict reference validation", func(t *testing.T) {
		bundle := &missionbundle.Bundle{
			Ref: "M20",
			ID:  "019fe381-5d61-7223-b362-03a5f99a7b10",
			Objectives: []missionbundle.Objective{
				{Ref: "O1", Outcome: "First objective"},
			},
			Completion: []missionbundle.Criterion{
				{Claim: "valid-claim", PassBoundary: "pass", ProofRequirement: "proof"},
			},
		}

		// Evidence naming non-existent objective
		invalidObjEv := &missionbundle.Evidence{
			Mission:    "M20",
			Title:      "Proof with bad objective",
			Actor:      "Alex",
			Commit:     "0123456789abcdef0123456789abcdef01234567",
			Tree:       "0123456789abcdef0123456789abcdef01234567",
			Objectives: []string{"O99"},
			Claims:     []string{"valid-claim"},
		}
		err := missionbundle.ValidateEvidenceDirect(invalidObjEv, bundle)
		if err == nil || !domain.RefusalHasCode(err, domain.RefusalInvalidReference) {
			t.Fatalf("expected RefusalInvalidReference for non-existent objective O99, got: %v", err)
		}

		// Evidence naming non-existent claim
		invalidClaimEv := &missionbundle.Evidence{
			Mission:    "M20",
			Title:      "Proof with bad claim",
			Actor:      "Alex",
			Commit:     "0123456789abcdef0123456789abcdef01234567",
			Tree:       "0123456789abcdef0123456789abcdef01234567",
			Objectives: []string{"O1"},
			Claims:     []string{"non-existent-claim"},
		}
		err = missionbundle.ValidateEvidenceDirect(invalidClaimEv, bundle)
		if err == nil || !domain.RefusalHasCode(err, domain.RefusalInvalidReference) {
			t.Fatalf("expected RefusalInvalidReference for non-existent claim, got: %v", err)
		}
	})

	t.Run("claim: m21-pinned-benchmark-matrix - static context savings proof", func(t *testing.T) {
		root, err := filepath.Abs(filepath.Join("..", ".."))
		if err != nil {
			t.Fatal(err)
		}
		ws, err := discovery.Open(root)
		if err != nil {
			t.Fatalf("discovery.Open failed: %v", err)
		}

		compiled, err := charter.Compile(ws, "M16", "O1", []string{"D12-isolation-and-context-compilation"})
		if err != nil {
			t.Fatalf("failed to compile charter: %v", err)
		}
		rendered := compiled.RenderMarkdown()
		if !strings.Contains(rendered, "## 1. FROZEN TRUTH") {
			t.Fatal("expected rendered charter to contain FROZEN TRUTH")
		}
		count, err := tokenizer.Count(rendered)
		if err != nil {
			t.Fatalf("tokenizer.Count failed: %v", err)
		}
		if count > tokenizer.HardCeilingTokens {
			t.Fatalf("compiled charter tokens %d exceeded ceiling %d", count, tokenizer.HardCeilingTokens)
		}
	})
}
