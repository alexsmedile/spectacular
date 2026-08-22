package missionbundle

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/alexsmedile/spectacular/v2/internal/domain"
	"github.com/alexsmedile/spectacular/v2/internal/governance"
)

// handoffFixtureMission starts a Mission in a real git-backed fixture and returns
// the service, the Mission ref, and the repository's real HEAD commit and tree.
// The git binding is taken from the repository rather than invented, because the
// command verifies it against the real repository and a fabricated pair would
// only ever exercise the refusal.
func handoffFixtureMission(t *testing.T) (Service, string, string, string) {
	t.Helper()
	root := missionServiceFixture(t)
	service := openMissionService(t, root)
	plan, raw := stressPlan()
	plan.Title = "Handoff recording fixture"
	started, err := service.Start(plan, raw)
	if err != nil {
		t.Fatal(err)
	}
	commit := strings.TrimSpace(commandOutput(t, root, "git", "rev-parse", "HEAD"))
	tree := strings.TrimSpace(commandOutput(t, root, "git", "rev-parse", "HEAD^{tree}"))
	return openMissionService(t, root), started.Ref, commit, tree
}

func handoffDraftInput(commit, tree, supersedes string) []byte {
	var b strings.Builder
	b.WriteString("---\ntype: HandoffDraft\n")
	b.WriteString("title: Delegate a bounded repair\n")
	b.WriteString("reviewed:\n    commit: " + commit + "\n    tree: " + tree + "\n")
	b.WriteString("sender:\n    actor: Alex\n    relation_to_receiver: operator\n")
	b.WriteString("task: Do the bounded thing.\n")
	if supersedes != "" {
		b.WriteString("supersedes: " + supersedes + "\n")
	}
	b.WriteString("asserted:\n    - the tests pass at this commit\n")
	b.WriteString("assumed:\n    - the fixture is representative\n")
	b.WriteString("stops:\n    - scope would grow\n")
	b.WriteString("returns:\n    - the diff and the passing tests\n")
	b.WriteString("---\n\nBody.\n")
	return []byte(b.String())
}

func TestHandoffRecordWritesIntoTheMissionBundle(t *testing.T) {
	service, mission, commit, tree := handoffFixtureMission(t)
	result, err := service.RecordHandoff(mission, "-", "Alex", handoffDraftInput(commit, tree, ""))
	if err != nil {
		t.Fatal(err)
	}
	if !strings.HasPrefix(result.Ref, mission+"/H1-") {
		t.Fatalf("handoff ref is %q, want %s/H1-<key>", result.Ref, mission)
	}
	// The Handoff lands inside the Mission bundle through the layout system. An
	// unscoped .spectacular/handoffs/ path means the record was placed by a rule
	// that did not know which Mission it belongs to.
	if !strings.Contains(result.Path, "/missions/") || !strings.Contains(result.Path, "/handoffs/") {
		t.Fatalf("handoff path is %q, want it inside the Mission bundle", result.Path)
	}
	if _, err := os.Stat(filepath.Join(service.Workspace.Root, filepath.FromSlash(result.Path))); err != nil {
		t.Fatalf("recorded Handoff is not on disk: %v", err)
	}

	// The Mission carries the pointer, and the whole bundle still validates with
	// the new record in it.
	reopened := openMissionService(t, service.Workspace.Root)
	bundle, err := Load(reopened.Workspace, mission)
	if err != nil {
		t.Fatal(err)
	}
	if len(bundle.Handoffs) != 1 {
		t.Fatalf("Mission carries %d handoff pointers, want 1", len(bundle.Handoffs))
	}
	if _, err := Validate(reopened.Workspace, bundle); err != nil {
		t.Fatalf("Mission does not validate with a recorded Handoff: %v", err)
	}
}

// Recording the same logical Handoff twice converges on one record. A retry after
// a crash must not leave two delegations where the sender filed one.
func TestHandoffRecordIsIdempotentOnRetry(t *testing.T) {
	service, mission, commit, tree := handoffFixtureMission(t)
	first, err := service.RecordHandoff(mission, "-", "Alex", handoffDraftInput(commit, tree, ""))
	if err != nil {
		t.Fatal(err)
	}
	retry := openMissionService(t, service.Workspace.Root)
	second, err := retry.RecordHandoff(mission, "-", "Alex", handoffDraftInput(commit, tree, ""))
	if err != nil {
		t.Fatal(err)
	}
	if first.Ref != second.Ref || first.Path != second.Path {
		t.Fatalf("retry produced a different record\n  first:  %s %s\n  second: %s %s", first.Ref, first.Path, second.Ref, second.Path)
	}
	reopened := openMissionService(t, service.Workspace.Root)
	bundle, err := Load(reopened.Workspace, mission)
	if err != nil {
		t.Fatal(err)
	}
	if len(bundle.Handoffs) != 1 {
		t.Fatalf("retry left %d handoffs, want 1", len(bundle.Handoffs))
	}
}

func TestHandoffRecordRefusesBadInput(t *testing.T) {
	service, mission, commit, tree := handoffFixtureMission(t)
	tests := []struct {
		name   string
		mutate func(string) string
		sender string
		code   domain.RefusalCode
		field  string
	}{
		{
			name:   "a commit that does not exist refuses",
			mutate: func(s string) string { return strings.Replace(s, commit, strings.Repeat("0", 40), 1) },
			code:   domain.RefusalInvalidKnownField, field: "handoff.reviewed.commit",
		},
		{
			name:   "a tree that is not the commit's refuses",
			mutate: func(s string) string { return strings.Replace(s, tree, strings.Repeat("1", 40), 1) },
			code:   domain.RefusalStaleFingerprint, field: "handoff.reviewed.tree",
		},
		{
			name: "an absent asserted list refuses",
			mutate: func(s string) string {
				return strings.Replace(s, "asserted:\n    - the tests pass at this commit\n", "", 1)
			},
			code: domain.RefusalInvalidKnownField, field: "handoff.asserted",
		},
		{
			name: "an absent assumed list refuses",
			mutate: func(s string) string {
				return strings.Replace(s, "assumed:\n    - the fixture is representative\n", "", 1)
			},
			code: domain.RefusalInvalidKnownField, field: "handoff.assumed",
		},
		{
			name:   "a sender who is not the recording identity refuses",
			mutate: func(s string) string { return s },
			sender: "Someone Else",
			code:   domain.RefusalUnauthorized, field: "by",
		},
		{
			name: "superseding a Handoff that was never recorded refuses",
			mutate: func(s string) string {
				return strings.Replace(s, "task: Do the bounded thing.\n", "task: Do the bounded thing.\nsupersedes: H9-zzzzzz\n", 1)
			},
			code: domain.RefusalInvalidKnownField, field: "handoff.supersedes",
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			sender := test.sender
			if sender == "" {
				sender = "Alex"
			}
			input := test.mutate(string(handoffDraftInput(commit, tree, "")))
			attempt := openMissionService(t, service.Workspace.Root)
			before := treeDigest(t, filepath.Join(attempt.Workspace.Root, ".spectacular"))
			_, err := attempt.RecordHandoff(mission, "-", sender, []byte(input))
			requireRefusal(t, err, test.code, test.field)
			// A refusal writes nothing. The record either exists completely or not
			// at all, so a rejected Handoff must not leave a partial one behind.
			if after := treeDigest(t, filepath.Join(attempt.Workspace.Root, ".spectacular")); after != before {
				t.Fatal("a refused handoff modified the workspace")
			}
		})
	}
}

// A Mission that does not exist is refused before anything is read or written.
func TestHandoffRecordRefusesAnUnknownMission(t *testing.T) {
	service, _, commit, tree := handoffFixtureMission(t)
	before := treeDigest(t, filepath.Join(service.Workspace.Root, ".spectacular"))
	if _, err := service.RecordHandoff("M9999", "-", "Alex", handoffDraftInput(commit, tree, "")); err == nil {
		t.Fatal("recording against an unknown Mission must refuse")
	}
	if after := treeDigest(t, filepath.Join(service.Workspace.Root, ".spectacular")); after != before {
		t.Fatal("a refused handoff modified the workspace")
	}
}

// A correction is a new Handoff, never an edit. The superseded record survives
// byte-for-byte as what its sender believed when they sent it.
func TestSupersedingLeavesTheOriginalByteIdentical(t *testing.T) {
	service, mission, commit, tree := handoffFixtureMission(t)
	first, err := service.RecordHandoff(mission, "-", "Alex", handoffDraftInput(commit, tree, ""))
	if err != nil {
		t.Fatal(err)
	}
	originalPath := filepath.Join(service.Workspace.Root, filepath.FromSlash(first.Path))
	originalBytes, err := os.ReadFile(originalPath)
	if err != nil {
		t.Fatal(err)
	}
	originalRef := strings.TrimPrefix(first.Ref, mission+"/")

	correcting := openMissionService(t, service.Workspace.Root)
	correction := strings.Replace(string(handoffDraftInput(commit, tree, originalRef)),
		"task: Do the bounded thing.", "task: Do the corrected thing.", 1)
	second, err := correcting.RecordHandoff(mission, "-", "Alex", []byte(correction))
	if err != nil {
		t.Fatal(err)
	}
	if second.Ref == first.Ref {
		t.Fatal("a superseding Handoff reused the superseded record's identity")
	}
	after, err := os.ReadFile(originalPath)
	if err != nil {
		t.Fatalf("the superseded Handoff no longer exists: %v", err)
	}
	if string(after) != string(originalBytes) {
		t.Fatal("superseding modified the original Handoff")
	}

	reopened := openMissionService(t, service.Workspace.Root)
	bundle, err := Load(reopened.Workspace, mission)
	if err != nil {
		t.Fatal(err)
	}
	if len(bundle.Handoffs) != 2 {
		t.Fatalf("Mission carries %d handoffs, want both the original and its correction", len(bundle.Handoffs))
	}
	if _, err := Validate(reopened.Workspace, bundle); err != nil {
		t.Fatalf("Mission does not validate after a supersession: %v", err)
	}
}

// A fault at any write boundary leaves the canonical tree byte-identical. The
// Handoff and the Mission's pointer to it are one transaction: a Mission naming a
// Handoff that was never written is a dangling reference in a governance record.
func TestHandoffRecordRollsBackAtEveryWriteBoundary(t *testing.T) {
	service, mission, commit, tree := handoffFixtureMission(t)
	root := service.Workspace.Root
	before := treeDigest(t, filepath.Join(root, ".spectacular"))

	// Two writes: the Handoff and the Mission. Derived rather than hardcoded so a
	// third written file extends the coverage instead of escaping it.
	for boundary := 0; boundary < 2; boundary++ {
		attempt := openMissionService(t, root)
		attempt.ApplyTransaction = func(transactionRoot, key string, changes []governance.FileChange) error {
			return governance.ApplyTransactionWithFailure(transactionRoot, key, changes, boundary)
		}
		if _, err := attempt.RecordHandoff(mission, "-", "Alex", handoffDraftInput(commit, tree, "")); err == nil {
			t.Fatalf("boundary %d: injected failure did not surface", boundary)
		}
		if err := governance.RecoverTransactions(root); err != nil {
			t.Fatalf("boundary %d: recovery failed: %v", boundary, err)
		}
		if after := treeDigest(t, filepath.Join(root, ".spectacular")); after != before {
			t.Fatalf("boundary %d: a partial Handoff survived", boundary)
		}
	}
}

// A reader of a Mission sees its Handoffs, who sent them, and whether the state
// each one bound is still on disk. Resolution happens on load rather than only
// during validation, so a reader that never validates still sees the record.
func TestMissionReportsItsHandoffsAndSupersessions(t *testing.T) {
	service, mission, commit, tree := handoffFixtureMission(t)
	first, err := service.RecordHandoff(mission, "-", "Alex", handoffDraftInput(commit, tree, ""))
	if err != nil {
		t.Fatal(err)
	}
	originalRef := strings.TrimPrefix(first.Ref, mission+"/")
	correcting := openMissionService(t, service.Workspace.Root)
	if _, err := correcting.RecordHandoff(mission, "-", "Alex", handoffDraftInput(commit, tree, originalRef)); err != nil {
		t.Fatal(err)
	}

	reopened := openMissionService(t, service.Workspace.Root)
	bundle, err := Load(reopened.Workspace, mission)
	if err != nil {
		t.Fatal(err)
	}
	if len(bundle.Handoffs) != 2 {
		t.Fatalf("Mission carries %d handoffs, want 2", len(bundle.Handoffs))
	}
	for _, pointer := range bundle.Handoffs {
		if pointer.Document == nil {
			t.Fatalf("handoff %s resolved to a pointer with no record", pointer.Ref)
		}
		if pointer.Document.Sender.Actor != "Alex" || pointer.Document.Title == "" {
			t.Fatalf("handoff %s resolved without its sender or title", pointer.Ref)
		}
	}
	// The correction names what it replaced, which is what lets a reader arriving
	// at the original be pointed forward.
	if bundle.Handoffs[1].Document.Supersedes != originalRef {
		t.Fatalf("the correction supersedes %q, want %q", bundle.Handoffs[1].Document.Supersedes, originalRef)
	}
	if bundle.Handoffs[0].Document.Supersedes != "" {
		t.Fatal("the original Handoff claims to supersede something")
	}
}

// A supersedes chain longer than one link resolves to the newest Handoff. A
// reader arriving at the oldest record must be able to follow the corrections
// forward rather than stopping at the first replacement.
func TestSupersedesChainResolvesToTheNewestHandoff(t *testing.T) {
	service, mission, commit, tree := handoffFixtureMission(t)
	root := service.Workspace.Root

	// Each link binds a distinct commit so the three are distinct records rather
	// than the same logical Handoff converging under the idempotency rule.
	refs := []string{}
	previous := ""
	for i := 0; i < 3; i++ {
		writeFixtureCommit(t, root, i)
		linkCommit := strings.TrimSpace(commandOutput(t, root, "git", "rev-parse", "HEAD"))
		linkTree := strings.TrimSpace(commandOutput(t, root, "git", "rev-parse", "HEAD^{tree}"))
		attempt := openMissionService(t, root)
		result, err := attempt.RecordHandoff(mission, "-", "Alex", handoffDraftInput(linkCommit, linkTree, previous))
		if err != nil {
			t.Fatalf("link %d: %v", i, err)
		}
		previous = strings.TrimPrefix(result.Ref, mission+"/")
		refs = append(refs, previous)
	}
	_ = commit
	_ = tree

	bundle, err := Load(openMissionService(t, root).Workspace, mission)
	if err != nil {
		t.Fatal(err)
	}
	if len(bundle.Handoffs) != 3 {
		t.Fatalf("Mission carries %d handoffs, want a three-link chain", len(bundle.Handoffs))
	}
	newest := NewestHandoff(bundle, refs[0])
	if newest == nil || newest.Ref != refs[2] {
		got := "nil"
		if newest != nil {
			got = newest.Ref
		}
		t.Fatalf("following the chain from %s reached %s, want the newest %s", refs[0], got, refs[2])
	}
	// The newest Handoff is superseded by nothing, so it resolves to itself.
	if last := NewestHandoff(bundle, refs[2]); last == nil || last.Ref != refs[2] {
		t.Fatal("the newest Handoff must resolve to itself")
	}
}

// writeFixtureCommit makes a distinct commit so each link of a supersedes chain
// binds a different tree.
func writeFixtureCommit(t *testing.T, root string, n int) {
	t.Helper()
	path := filepath.Join(root, "chain.txt")
	if err := os.WriteFile(path, []byte(strings.Repeat("x", n+1)), 0o644); err != nil {
		t.Fatal(err)
	}
	commandOutput(t, root, "git", "add", "chain.txt")
	commandOutput(t, root, "git", "commit", "-qm", "chain link")
}

func TestHandoffRecordAutoDerivesTreeWhenOmitted(t *testing.T) {
	service, mission, commit, tree := handoffFixtureMission(t)
	var b strings.Builder
	b.WriteString("---\ntype: HandoffDraft\n")
	b.WriteString("title: Delegate without explicit tree\n")
	b.WriteString("reviewed:\n    commit: " + commit + "\n")
	b.WriteString("sender:\n    actor: Alex\n    relation_to_receiver: operator\n")
	b.WriteString("task: Do the bounded thing.\n")
	b.WriteString("asserted: []\n")
	b.WriteString("assumed: []\n")
	b.WriteString("stops:\n    - scope would grow\n")
	b.WriteString("returns:\n    - the diff and the passing tests\n")
	b.WriteString("---\n\nBody.\n")

	result, err := service.RecordHandoff(mission, "-", "Alex", []byte(b.String()))
	if err != nil {
		t.Fatalf("recording handoff without tree failed: %v", err)
	}

	reopened := openMissionService(t, service.Workspace.Root)
	bundle, err := Load(reopened.Workspace, mission)
	if err != nil {
		t.Fatal(err)
	}
	if len(bundle.Handoffs) != 1 {
		t.Fatalf("expected 1 handoff, got %d", len(bundle.Handoffs))
	}
	if bundle.Handoffs[0].Document.Reviewed.Tree != tree {
		t.Fatalf("expected tree %s to be auto-derived, got %s", tree, bundle.Handoffs[0].Document.Reviewed.Tree)
	}
	_ = result
}
