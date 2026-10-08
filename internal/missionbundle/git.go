package missionbundle

import (
	"os/exec"
	"strings"

	"github.com/alexsmedile/spectacular/v2/internal/domain"
)

func gitBaseline(root string) (string, string, string, error) {
	commitCommand := exec.Command("git", "rev-parse", "HEAD")
	commitCommand.Dir = root
	commit, err := commitCommand.Output()
	if err != nil {
		return "", "", "", invalidCause("baseline.commit", "read Git HEAD", err)
	}
	branchCommand := exec.Command("git", "branch", "--show-current")
	branchCommand.Dir = root
	branch, err := branchCommand.Output()
	if err != nil || strings.TrimSpace(string(branch)) == "" {
		return "", "", "", invalidCause("baseline.branch", "read current Git branch", err)
	}
	timeCommand := exec.Command("git", "show", "-s", "--format=%cI", strings.TrimSpace(string(commit)))
	timeCommand.Dir = root
	baselineAt, err := timeCommand.Output()
	if err != nil {
		return "", "", "", invalidCause("baseline.commit", "read Git commit time", err)
	}
	return strings.TrimSpace(string(commit)), strings.TrimSpace(string(branch)), strings.TrimSpace(string(baselineAt)), nil
}

func resolveCommitTree(root, commit string) (string, error) {
	cmd := exec.Command("git", "rev-parse", "--verify", commit+"^{tree}")
	cmd.Dir = root
	out, err := cmd.Output()
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(string(out)), nil
}

// verifyReviewedGit checks a record's git binding against the real repository.
// The field prefix names the record being verified, so a Handoff's refusal does
// not tell its sender to correct a review.
func verifyReviewedGit(root, commit, tree string, field string) error {
	if !commitPattern.MatchString(commit) || !commitPattern.MatchString(tree) {
		return invalid(field+".reviewed", "a bound record must carry canonical 40-character Git commit and tree IDs")
	}
	commitCommand := exec.Command("git", "rev-parse", "--verify", commit+"^{commit}")
	commitCommand.Dir = root
	resolvedCommit, err := commitCommand.Output()
	if err != nil || strings.TrimSpace(string(resolvedCommit)) != commit {
		return domain.NewStateRefusal(domain.RefusalInvalidKnownField, field+".reviewed.commit", "reviewed commit does not exist in this repository", "existing exact commit", commit, "review the committed tree, then record its exact commit and tree", err)
	}
	treeCommand := exec.Command("git", "rev-parse", "--verify", commit+"^{tree}")
	treeCommand.Dir = root
	resolvedTree, err := treeCommand.Output()
	actual := strings.TrimSpace(string(resolvedTree))
	if err != nil || actual != tree {
		return domain.NewStateRefusal(domain.RefusalStaleFingerprint, field+".reviewed.tree", "reviewed tree does not belong to the reviewed commit", actual, tree, "record the exact tree from git rev-parse <commit>^{tree}", err)
	}
	return nil
}
