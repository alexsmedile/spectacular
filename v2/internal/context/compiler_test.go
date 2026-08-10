package context

import (
	"io/fs"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/alexsmedile/spectacular/v2/internal/discovery"
)

func scenarioA(t *testing.T) *discovery.Workspace {
	t.Helper()
	root, err := filepath.Abs(filepath.Join("..", "..", "testdata", "scenario-a"))
	if err != nil {
		t.Fatal(err)
	}
	opened, err := discovery.Open(root)
	if err != nil {
		t.Fatal(err)
	}
	return opened
}

func TestCompileProjectAndMissionAreBoundedAndSourceBacked(t *testing.T) {
	now := func() time.Time { return time.Date(2026, 8, 10, 10, 0, 0, 0, time.UTC) }
	compiler := Compiler{Workspace: scenarioA(t), Now: now}
	project, err := compiler.Compile(Config{Scope: "project", Event: "@Orient"})
	if err != nil {
		t.Fatal(err)
	}
	if project.SchemaVersion != SchemaVersion || project.GenerationBasis == "" || len(project.Authoritative) != 2 || project.LoadedRecords >= project.AvailableRecords {
		t.Fatalf("unexpected project context %#v", project)
	}
	if project.Next.Kind != "continuation" || project.Next.Operation != "resume" || len(project.Omissions) != 1 {
		t.Fatalf("project continuation/absence not explicit: %#v", project)
	}
	mission, err := compiler.Compile(Config{Scope: "Mission:0198a1a0-0000-7000-8000-000000000002"})
	if err != nil {
		t.Fatal(err)
	}
	if mission.Scope == "project" || len(mission.Authoritative) < 4 || len(mission.Gaps) != 1 || mission.Next.Target == "" {
		t.Fatalf("unexpected Mission context %#v", mission)
	}
	again, err := compiler.Compile(Config{Scope: mission.Scope})
	if err != nil {
		t.Fatal(err)
	}
	if mission.GenerationBasis != again.GenerationBasis {
		t.Fatalf("generation basis changed: %s != %s", mission.GenerationBasis, again.GenerationBasis)
	}
}

func TestCompileRefusesInvalidScopeAndSelectorWithoutEvent(t *testing.T) {
	compiler := Compiler{Workspace: scenarioA(t)}
	if _, err := compiler.Compile(Config{Scope: "all"}); err == nil {
		t.Fatal("accepted generic context scope")
	}
	if _, err := compiler.Compile(Config{Scope: "project", Selector: "$git.commit"}); err == nil {
		t.Fatal("accepted selector without event")
	}
}

func TestCompileSelectsDeclaredGuardrailsWithoutMakingThemAuthority(t *testing.T) {
	source, err := filepath.Abs(filepath.Join("..", "..", "testdata", "scenario-a"))
	if err != nil {
		t.Fatal(err)
	}
	root := t.TempDir()
	if err := filepath.WalkDir(source, func(path string, entry fs.DirEntry, walkErr error) error {
		if walkErr != nil {
			return walkErr
		}
		relative, _ := filepath.Rel(source, path)
		if relative == "." {
			return nil
		}
		target := filepath.Join(root, relative)
		if entry.IsDir() {
			return os.Mkdir(target, 0o755)
		}
		data, readErr := os.ReadFile(path)
		if readErr != nil {
			return readErr
		}
		return os.WriteFile(target, data, 0o644)
	}); err != nil {
		t.Fatal(err)
	}
	guidance := "## @Run\nStop on baseline drift.\n## @Run $git.commit\nCommit coherent changes only.\n"
	if err := os.WriteFile(filepath.Join(root, ".spectacular", "GUARDRAILS.md"), []byte(guidance), 0o644); err != nil {
		t.Fatal(err)
	}
	opened, err := discovery.Open(root)
	if err != nil {
		t.Fatal(err)
	}
	opened.Manifest.Guardrails = "GUARDRAILS.md"
	bundle, err := (Compiler{Workspace: opened, Now: func() time.Time { return time.Date(2026, 8, 10, 10, 0, 0, 0, time.UTC) }}).Compile(Config{Scope: "project", Event: "@Run", Selector: "$git.commit"})
	if err != nil {
		t.Fatal(err)
	}
	if len(bundle.Guidance) != 2 || bundle.Guidance[1].Prose != "Commit coherent changes only." {
		t.Fatalf("Guardrails selection %#v", bundle.Guidance)
	}
	found := false
	for _, source := range bundle.ProjectionSources {
		if source.Role == "owner-guardrails" && source.Authority == "owner-guidance" {
			found = true
		}
	}
	if !found {
		t.Fatal("Guardrails source or non-authoritative role is missing")
	}
}
