package humanlayout

import (
	"encoding/json"
	"strings"
	"testing"

	"github.com/alexsmedile/spectacular/v2/internal/discovery"
	"github.com/alexsmedile/spectacular/v2/internal/domain"
	"github.com/alexsmedile/spectacular/v2/internal/workspace"
)

func TestPlanBuildsReadableScopedMissionBundle(t *testing.T) {
	mission := document(t, domain.Mission, "019fe381-5d61-7223-b362-03a5f99a7b02", "Restore human operability")
	status := "active"
	mission.Record.Status = &status
	missionRef := "Mission:" + mission.Record.ID.String()
	objective := document(t, domain.Objective, "019fe381-5d61-7223-b362-03a5f99a7b03", "Implement readable workspace")
	workspace.SetString(objective, "mission", missionRef)
	run := document(t, domain.Run, "019fe381-5d61-7223-b362-03a5f99a7b05", "Implement layout")
	workspace.SetString(run, "mission", missionRef)
	checkpoint := document(t, domain.Checkpoint, "019fe381-5d61-7223-b362-03a5f99a7b06", "Layout approved")
	workspace.SetString(checkpoint, "run", "Run:"+run.Record.ID.String())
	evidence := document(t, domain.Evidence, "019fe381-5d61-7223-b362-03a5f99a7b07", "Filesystem proof")
	workspace.SetString(evidence, "mission", missionRef)

	docs := []*workspace.Document{evidence, checkpoint, run, objective, mission}
	paths, err := Plan(nil, docs)
	if err != nil {
		t.Fatal(err)
	}
	for _, expected := range []struct {
		doc  *workspace.Document
		ref  string
		path string
	}{
		{mission, "M1", ".spectacular/missions/M1-restore-human-operability/M1-restore-human-operability.md"},
		{objective, "M1/O1", ".spectacular/missions/M1-restore-human-operability/objectives/O1-implement-readable-workspace.md"},
		{run, "M1/R1", ".spectacular/missions/M1-restore-human-operability/runs/R1-implement-layout/R1-implement-layout.md"},
		{checkpoint, "M1/R1/C1", ".spectacular/missions/M1-restore-human-operability/runs/R1-implement-layout/checkpoints/C1-layout-approved.md"},
		{evidence, "M1/E1-t2lylz", ".spectacular/missions/M1-restore-human-operability/evidence/E1-t2lylz.md"},
	} {
		if got := HumanRef(expected.doc); got != expected.ref {
			t.Errorf("ref=%q want %q", got, expected.ref)
		}
		if got := paths[expected.doc.Record.ID]; got != expected.path {
			t.Errorf("path=%q want %q", got, expected.path)
		}
	}

	indexes, err := Indexes(nil, docs, paths)
	if err != nil {
		t.Fatal(err)
	}
	var catalog struct {
		Type    string     `json:"type"`
		Scope   string     `json:"scope"`
		Entries []indexRow `json:"entries"`
	}
	if err := json.Unmarshal(indexes[".spectacular/index.json"], &catalog); err != nil {
		t.Fatal(err)
	}
	if catalog.Type != "Index" || catalog.Scope != "governed-records" || len(catalog.Entries) != len(docs) {
		t.Fatalf("generated inventory claims the wrong scope: %#v", catalog)
	}
	for path := range indexes {
		if strings.HasSuffix(strings.ToLower(path), "index.md") {
			t.Fatalf("generated output can overwrite manual navigation: %s", path)
		}
	}

}

func TestIndexesClearAnEmptiedActiveCollection(t *testing.T) {
	mission := document(t, domain.Mission, "019fe381-5d61-7223-b362-03a5f99a7b02", "Archived Mission")
	workspace.SetString(mission, "human_ref", "M1")
	existing := []discovery.Entry{{
		Document: mission,
		Path:     ".spectacular/missions/M1-archived-mission/M1-archived-mission.md",
	}}
	archived := map[domain.ID]string{
		mission.Record.ID: ".spectacular/archive/missions/M1-archived-mission/M1-archived-mission.md",
	}
	indexes, err := Indexes(existing, []*workspace.Document{mission}, archived)
	if err != nil {
		t.Fatal(err)
	}
	active := string(indexes[".spectacular/missions/index.json"])
	if !strings.Contains(active, `"entries": []`) || strings.Contains(active, `"ref": "M1"`) {
		t.Fatalf("active Mission index was not cleared:\n%s", active)
	}
	if archivedIndex := string(indexes[".spectacular/archive/index.json"]); !strings.Contains(archivedIndex, `"ref": "M1"`) {
		t.Fatalf("archive collection index omits moved Mission:\n%s", archivedIndex)
	}
	if _, exists := indexes[".spectacular/archive/missions/M1-archived-mission/index.md"]; exists {
		t.Fatal("Mission-local index duplicates mission root")
	}
}

func TestShortKeyIsStableIdentityNotContent(t *testing.T) {
	id, err := domain.ParseID("019fe381-5d61-7223-b362-03a5f99a7b07")
	if err != nil {
		t.Fatal(err)
	}
	if got := ShortKey(id); got != "t2lylz" {
		t.Fatalf("short key=%q", got)
	}
}

func TestHumanRefPrefersCompactRef(t *testing.T) {
	doc := document(t, domain.Mission, "019fe381-5d61-7223-b362-03a5f99a7b02", "Compact Mission")
	workspace.SetString(doc, "human_ref", "M-old")
	workspace.SetString(doc, "ref", "M5")
	if got := HumanRef(doc); got != "M5" {
		t.Fatalf("HumanRef=%q want M5", got)
	}
}

func TestPlanKeepsExistingBundleDirectoriesStable(t *testing.T) {
	mission := document(t, domain.Mission, "019fe381-5d61-7223-b362-03a5f99a7b02", "Renamed Mission title")
	workspace.SetString(mission, "human_ref", "M1")
	existing := []discovery.Entry{{
		Document: mission,
		Path:     ".spectacular/missions/M1-original-bundle-name/M1-original-bundle-name.md",
	}}
	evidence := document(t, domain.Evidence, "019fe381-5d61-7223-b362-03a5f99a7b07", "Fresh proof")
	workspace.SetString(evidence, "mission", "Mission:"+mission.Record.ID.String())
	paths, err := Plan(existing, []*workspace.Document{evidence})
	if err != nil {
		t.Fatal(err)
	}
	if got := paths[evidence.Record.ID]; got != ".spectacular/missions/M1-original-bundle-name/evidence/E1-t2lylz.md" {
		t.Fatalf("new child escaped stable Mission bundle: %s", got)
	}
}

func TestPlanStandaloneReview(t *testing.T) {
	review := document(t, domain.Review, "019fe381-5d61-7223-b362-03a5f99a7b09", "Standalone PR Review")
	paths, err := Plan(nil, []*workspace.Document{review})
	if err != nil {
		t.Fatal(err)
	}
	if got := paths[review.Record.ID]; got != ".spectacular/reviews/RV1-standalone-pr-review.md" {
		t.Fatalf("standalone review placed incorrectly: %s", got)
	}
}

func document(t *testing.T, noun domain.RecordType, raw, title string) *workspace.Document {
	t.Helper()
	id, err := domain.ParseID(raw)
	if err != nil {
		t.Fatal(err)
	}
	return &workspace.Document{Record: domain.Record{Type: noun, ID: id, Title: &title}}
}
