package index

import (
	"path/filepath"
	"reflect"
	"strings"
	"testing"

	"github.com/alexsmedile/spectacular/v2/internal/domain"
	"github.com/alexsmedile/spectacular/v2/internal/workspace"
)

func TestExactIdentityAndFullPathLookup(t *testing.T) {
	t.Parallel()

	proposal := fixtureEntry(t, "records/proposals/semantic.md", "positive", "proposal.md")
	mission := fixtureEntry(t, "records/missions/m1.md", "positive", "mission.md")
	index, err := New([]Entry{mission, proposal})
	if err != nil {
		t.Fatalf("New with relationship target discovered later: %v", err)
	}

	byID, found := index.LookupID(proposal.Record.ID)
	if !found || byID.Path != proposal.Path {
		t.Fatalf("LookupID = (%#v, %t)", byID, found)
	}
	byPath, found := index.LookupPath(mission.Path)
	if !found || byPath.Record.ID != mission.Record.ID {
		t.Fatalf("LookupPath = (%#v, %t)", byPath, found)
	}
	if _, found := index.LookupPath(filepath.Base(proposal.Path)); found {
		t.Fatal("basename unexpectedly matched a full workspace-relative path")
	}
}

func TestRenameSafeIdentityLookupAndRelationship(t *testing.T) {
	t.Parallel()

	proposal := fixtureEntry(t, "records/proposals/original.md", "positive", "proposal.md")
	mission := fixtureEntry(t, "records/missions/m1.md", "positive", "mission.md")
	before, err := New([]Entry{proposal, mission})
	if err != nil {
		t.Fatal(err)
	}
	if _, found := before.LookupID(proposal.Record.ID); !found {
		t.Fatal("proposal identity absent before rename")
	}

	proposal.Path = "records/proposals/renamed.md"
	after, err := New([]Entry{mission, proposal})
	if err != nil {
		t.Fatalf("relationship broke after path rename: %v", err)
	}
	byID, found := after.LookupID(proposal.Record.ID)
	if !found || byID.Path != proposal.Path {
		t.Fatalf("renamed identity lookup = (%#v, %t)", byID, found)
	}
	if _, found := after.LookupPath("records/proposals/original.md"); found {
		t.Fatal("old path remained indexed")
	}
}

func TestDiscoveryOrderDoesNotAffectIndexResults(t *testing.T) {
	t.Parallel()

	proposal := fixtureEntry(t, "z/proposal.md", "positive", "proposal.md")
	mission := fixtureEntry(t, "a/mission.md", "positive", "mission.md")
	forward, err := New([]Entry{proposal, mission})
	if err != nil {
		t.Fatal(err)
	}
	reverse, err := New([]Entry{mission, proposal})
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(forward.Entries(), reverse.Entries()) {
		t.Fatalf("index results depend on discovery order:\nforward: %#v\nreverse: %#v", forward.Entries(), reverse.Entries())
	}
}

func TestDuplicateIdentityRefusalIsDeterministic(t *testing.T) {
	t.Parallel()

	first := fixtureEntry(t, "z/duplicate-b.md", "negative", "duplicate-b.md")
	second := fixtureEntry(t, "a/duplicate-a.md", "negative", "duplicate-a.md")
	_, forwardError := New([]Entry{first, second})
	_, reverseError := New([]Entry{second, first})
	if !domain.RefusalHasCode(forwardError, domain.RefusalDuplicateID) {
		t.Fatalf("forward error = %v, want duplicate_id", forwardError)
	}
	if !domain.RefusalHasCode(reverseError, domain.RefusalDuplicateID) {
		t.Fatalf("reverse error = %v, want duplicate_id", reverseError)
	}
	if forwardError.Error() != reverseError.Error() {
		t.Fatalf("duplicate refusal depends on discovery order:\nforward: %v\nreverse: %v", forwardError, reverseError)
	}
	if !strings.Contains(forwardError.Error(), "a/duplicate-a.md") ||
		!strings.Contains(forwardError.Error(), "z/duplicate-b.md") {
		t.Fatalf("duplicate refusal lacks both complete paths: %v", forwardError)
	}
}

func TestDuplicatePathRefusalIsDiscoveryOrderIndependent(t *testing.T) {
	t.Parallel()

	proposal := fixtureEntry(t, "same/record.md", "positive", "proposal.md")
	mission := fixtureEntry(t, "same/record.md", "negative", "wrong-type-target.md")
	_, forwardError := New([]Entry{proposal, mission})
	_, reverseError := New([]Entry{mission, proposal})
	if !domain.RefusalHasCode(forwardError, domain.RefusalDuplicatePath) {
		t.Fatalf("forward error = %v, want duplicate_path", forwardError)
	}
	if !domain.RefusalHasCode(reverseError, domain.RefusalDuplicatePath) {
		t.Fatalf("reverse error = %v, want duplicate_path", reverseError)
	}
	if forwardError.Error() != reverseError.Error() {
		t.Fatalf("duplicate-path refusal depends on discovery order:\nforward: %v\nreverse: %v", forwardError, reverseError)
	}
	for _, identity := range []string{proposal.Record.ID.String(), mission.Record.ID.String()} {
		if !strings.Contains(forwardError.Error(), identity) {
			t.Fatalf("duplicate-path refusal lacks identity %s: %v", identity, forwardError)
		}
	}
}

func TestBrokenTargetIsRefusedAfterDiscovery(t *testing.T) {
	t.Parallel()

	broken := fixtureEntry(t, "missions/broken.md", "negative", "broken-target.md")
	_, err := New([]Entry{broken})
	if !domain.RefusalHasCode(err, domain.RefusalTargetNotFound) {
		t.Fatalf("New error = %v, want target_not_found", err)
	}
}

func TestWrongTargetTypeIsRefusedAfterDiscovery(t *testing.T) {
	t.Parallel()

	source := fixtureEntry(t, "missions/source.md", "negative", "wrong-type-source.md")
	target := fixtureEntry(t, "missions/target.md", "negative", "wrong-type-target.md")
	_, err := New([]Entry{source, target})
	if !domain.RefusalHasCode(err, domain.RefusalTargetTypeMismatch) {
		t.Fatalf("New error = %v, want target_type_mismatch", err)
	}
}

func TestCanonicalWorkspaceRelativePathsAreRequired(t *testing.T) {
	t.Parallel()

	record := fixtureEntry(t, "valid/proposal.md", "positive", "proposal.md").Record
	for _, invalid := range []string{"", "/absolute.md", "../outside.md", "a/../record.md", `a\record.md`} {
		invalid := invalid
		t.Run(invalid, func(t *testing.T) {
			t.Parallel()
			_, err := New([]Entry{{Path: invalid, Record: record}})
			if !domain.RefusalHasCode(err, domain.RefusalInvalidWorkspacePath) {
				t.Fatalf("New(%q) error = %v, want invalid_workspace_path", invalid, err)
			}
		})
	}
}

func fixtureEntry(t *testing.T, workspacePath string, fixtureParts ...string) Entry {
	t.Helper()
	parts := append([]string{"..", "..", "testdata", "m1"}, fixtureParts...)
	document, err := workspace.ReadFile(filepath.Join(parts...))
	if err != nil {
		t.Fatalf("read fixture: %v", err)
	}
	return Entry{Path: workspacePath, Record: document.Record}
}
