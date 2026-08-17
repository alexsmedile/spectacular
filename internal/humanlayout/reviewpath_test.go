package humanlayout

import (
	"path/filepath"
	"testing"

	"github.com/alexsmedile/spectacular/v2/internal/discovery"
	"github.com/alexsmedile/spectacular/v2/internal/domain"
	"github.com/alexsmedile/spectacular/v2/internal/workspace"
)

// A Review's path is produced by the layout system rather than by a join written
// at the call site. Moving it there is only correct if it is behavior-preserving,
// so this asserts the layout rule reproduces the path every recorded Review
// already occupies. A rule that relocated recorded Reviews would be a migration,
// not a refactor.
func TestLayoutReproducesEveryRecordedReviewPath(t *testing.T) {
	root, err := filepath.Abs(filepath.Join("..", ".."))
	if err != nil {
		t.Fatal(err)
	}
	opened, err := discovery.Open(root)
	if err != nil {
		t.Fatal(err)
	}
	// M5's Review was hand-authored before `review record` existed. Its filename
	// does not derive from its title, so no rule reproduces it. It is excluded by
	// name, with its reason, rather than by a tolerance in the layout rule that
	// would silently accept any drift.
	legacy := map[string]bool{"M5": true}

	all := map[domain.ID]*workspace.Document{}
	for _, entry := range opened.Entries {
		all[entry.Document.Record.ID] = entry.Document
	}
	reviewed, skipped := 0, 0
	for _, entry := range opened.Entries {
		if entry.Document.Record.Type != domain.Review {
			continue
		}
		mission, _ := workspace.String(entry.Document, "mission", false)
		if legacy[mission] {
			skipped++
			continue
		}
		reviewed++
		planned, err := PlannedPath(opened.Entries, entry.Document)
		if err != nil {
			t.Fatalf("layout has no path for recorded Review %s: %v", HumanRef(entry.Document), err)
		}
		if actual := filepath.ToSlash(entry.Path); planned != actual {
			t.Fatalf("layout would move a recorded Review\n  ref:     %s\n  on disk: %s\n  layout:  %s", HumanRef(entry.Document), actual, planned)
		}
	}
	if reviewed == 0 {
		t.Fatal("no recorded Review found in the workspace; this test proves nothing")
	}
	if skipped != len(legacy) {
		t.Fatalf("the legacy exclusion names %d Mission(s) but matched %d; a named exception that matches nothing is stale", len(legacy), skipped)
	}
}

// A Review is Mission-scoped and lands in the Mission's reviews/ directory. This
// pins the shape independently of what happens to be recorded today.
func TestReviewPathIsMissionScoped(t *testing.T) {
	mission := &workspace.Document{Record: domain.Record{Type: domain.Mission, ID: mustID(t, "01a0102c-a360-71fe-a1be-8e1b010460b2"), Title: stringPointer("Some mission")}}
	workspace.SetString(mission, "ref", "M42")
	review := &workspace.Document{Record: domain.Record{Type: domain.Review, ID: mustID(t, "01a0103d-7fc8-74d4-9a8c-f3dae87cc687"), Title: stringPointer("Independent review of M42")}}
	workspace.SetString(review, "mission", "M42")
	workspace.SetString(review, "human_ref", "M42/RV1")
	all := map[domain.ID]*workspace.Document{mission.Record.ID: mission, review.Record.ID: review}
	planned, err := Path(review, all)
	if err != nil {
		t.Fatal(err)
	}
	expected := ".spectacular/missions/M42-some-mission/reviews/RV1-independent-review-of-m42.md"
	if planned != expected {
		t.Fatalf("review path\n  want %s\n  got  %s", expected, planned)
	}
}

func mustID(t *testing.T, raw string) domain.ID {
	t.Helper()
	id, err := domain.ParseID(raw)
	if err != nil {
		t.Fatal(err)
	}
	return id
}

func stringPointer(value string) *string { return &value }
