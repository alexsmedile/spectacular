package humanlayout

import (
	"strings"
	"testing"

	"github.com/alexsmedile/spectacular/v2/internal/discovery"
	"github.com/alexsmedile/spectacular/v2/internal/domain"
	"github.com/alexsmedile/spectacular/v2/internal/workspace"
)

func TestPlanBuildsReadableScopedMissionBundle(t *testing.T) {
	mission := document(t, domain.Mission, "019fe381-5d61-7223-b362-03a5f99a7b02", "Restore human operability")
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
		{mission, "M1", ".spectacular/missions/M1-restore-human-operability/MISSION.md"},
		{objective, "M1/O1", ".spectacular/missions/M1-restore-human-operability/objectives/O1-implement-readable-workspace.md"},
		{run, "M1/R1", ".spectacular/missions/M1-restore-human-operability/runs/R1-implement-layout/RUN.md"},
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
	root := string(indexes[".spectacular/index.md"])
	if !strings.Contains(root, "non-authoritative") || !strings.Contains(root, "`M1/R1/C1`") {
		t.Fatalf("index lacks projection boundary or scoped ref:\n%s", root)
	}
}

func TestIndexesClearAnEmptiedActiveCollection(t *testing.T) {
	mission := document(t, domain.Mission, "019fe381-5d61-7223-b362-03a5f99a7b02", "Archived Mission")
	workspace.SetString(mission, "human_ref", "M1")
	existing := []discovery.Entry{{
		Document: mission,
		Path:     ".spectacular/missions/M1-archived-mission/MISSION.md",
	}}
	archived := map[domain.ID]string{
		mission.Record.ID: ".spectacular/archive/missions/M1-archived-mission/MISSION.md",
	}
	indexes, err := Indexes(existing, []*workspace.Document{mission}, archived)
	if err != nil {
		t.Fatal(err)
	}
	active := string(indexes[".spectacular/missions/index.md"])
	if !strings.Contains(active, "non-authoritative") || strings.Contains(active, "`M1`") {
		t.Fatalf("active Mission index was not cleared:\n%s", active)
	}
	if archivedIndex := string(indexes[".spectacular/archive/missions/M1-archived-mission/index.md"]); !strings.Contains(archivedIndex, "`M1`") {
		t.Fatalf("archive index omits moved Mission:\n%s", archivedIndex)
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

func TestPlanKeepsExistingBundleDirectoriesStable(t *testing.T) {
	mission := document(t, domain.Mission, "019fe381-5d61-7223-b362-03a5f99a7b02", "Renamed Mission title")
	workspace.SetString(mission, "human_ref", "M1")
	existing := []discovery.Entry{{
		Document: mission,
		Path:     ".spectacular/missions/M1-original-bundle-name/MISSION.md",
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

func document(t *testing.T, noun domain.RecordType, raw, title string) *workspace.Document {
	t.Helper()
	id, err := domain.ParseID(raw)
	if err != nil {
		t.Fatal(err)
	}
	return &workspace.Document{Record: domain.Record{Type: noun, ID: id, Title: &title}}
}
