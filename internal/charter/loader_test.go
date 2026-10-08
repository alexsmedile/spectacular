package charter

import (
	"errors"
	"testing"

	"github.com/alexsmedile/spectacular/v2/internal/discovery"
	"github.com/alexsmedile/spectacular/v2/internal/missionview"
)

func TestCompilerUsesReadModelPort(t *testing.T) {
	ws, err := discovery.Open(".")
	if err != nil {
		t.Fatal(err)
	}
	calls := 0
	load := func(got *discovery.Workspace, ref string) (*missionview.Mission, error) {
		calls++
		if got != ws || ref != "M16" {
			t.Fatalf("unexpected loader arguments: %p %s", got, ref)
		}
		return &missionview.Mission{Ref: ref, Objectives: []missionview.Objective{{Ref: "O1", Outcome: "isolated compiler"}}, WritesPaths: []string{"src/"}, AllowedActions: []string{"inspect"}}, nil
	}
	c, err := Compile(ws, "M16", "M16/O1", nil, load)
	if err != nil {
		t.Fatal(err)
	}
	if calls != 1 || c.Layer1.Outcome != "isolated compiler" || c.Layer3.WritesPaths[0] != "src/" {
		t.Fatalf("compiler did not use the supplied view: %+v", c)
	}
	sentinel := errors.New("loader failed")
	_, err = Compile(ws, "M16", "O1", nil, func(*discovery.Workspace, string) (*missionview.Mission, error) { return nil, sentinel })
	if !errors.Is(err, sentinel) {
		t.Fatalf("loader error lost: %v", err)
	}
	if _, err := Compile(ws, "M16", "O1", nil, nil); err == nil {
		t.Fatal("missing loader accepted")
	}
	if _, err := Compile(ws, "M16", "O1", nil, func(*discovery.Workspace, string) (*missionview.Mission, error) { return nil, nil }); err == nil {
		t.Fatal("nil view accepted")
	}
}
