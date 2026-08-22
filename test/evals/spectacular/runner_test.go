package spectaculareval

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestPairedRunnerRandomizesAndIsolatesArtifacts(t *testing.T) {
	repo := t.TempDir()
	runGit(t, repo, "init", "-q")
	runGit(t, repo, "config", "user.name", "Eval")
	runGit(t, repo, "config", "user.email", "eval@example.invalid")
	writeMinimalSkill := func(title string) {
		writeTestFile(t, filepath.Join(repo, "skills", "spectacular", "SKILL.md"), "---\nname: spectacular\n---\n# "+title+"\n")
		for _, name := range primaryRoutes {
			writeTestFile(t, filepath.Join(repo, "skills", "spectacular", "references", name+".md"), "# "+name+"\n\nUse this when: routed.\n")
		}
	}
	writeMinimalSkill("Old")
	runGit(t, repo, "add", ".")
	runGit(t, repo, "commit", "-qm", "baseline")
	baseline := strings.TrimSpace(runGit(t, repo, "rev-parse", "HEAD"))
	writeMinimalSkill("New")
	runGit(t, repo, "add", ".")
	runGit(t, repo, "commit", "-qm", "candidate")
	candidate := strings.TrimSpace(runGit(t, repo, "rev-parse", "HEAD"))

	root := t.TempDir()
	fixture := filepath.Join(root, "fixtures", "plain")
	writeTestFile(t, filepath.Join(fixture, ".gitignore"), ".agents/\n")
	writeTestFile(t, filepath.Join(fixture, "README.md"), "fixture\n")
	zero := 0
	catalog := Catalog{
		SchemaVersion: CatalogSchema,
		Tiers:         map[string]Tier{"smoke": {Repetitions: 1, Include: []string{"smoke"}}},
		Metrics:       metricDefinitionsForTest(),
		Cases: []Case{{
			ID: "OR-01", Kind: "behavior", Tier: "smoke", Fixture: "plain", Prompt: "orient",
			Weights: map[string]float64{"safety": 1, "routing": 1, "recovery": 1},
			Expect:  Expectation{Role: "Orchestrator", Phase: "orient", Status: "done", ExpectedReferences: []string{"orient.md"}, ForbiddenChangedPaths: []string{"**"}, MaximumOwnerQuestions: &zero, ExactlyOnePrimaryRef: true, RequireSingleReturn: true},
		}},
	}
	catalogData, err := json.Marshal(catalog)
	if err != nil {
		t.Fatal(err)
	}
	catalogPath := filepath.Join(root, "evals.json")
	if err := os.WriteFile(catalogPath, catalogData, 0o644); err != nil {
		t.Fatal(err)
	}
	schemaPath := filepath.Join(root, "result.schema.json")
	writeTestFile(t, schemaPath, "{}\n")
	output := filepath.Join(root, "out")
	config := RunConfig{
		Repo: repo, CatalogPath: catalogPath, SchemaPath: schemaPath, BaselineRef: baseline,
		CandidateRef: candidate, Tier: "smoke", Seed: 7, Model: "fake",
		Adapter: "testdata/fake-adapter.sh", OutputDir: output,
	}
	report, err := RunPaired(config)
	if err != nil {
		t.Fatal(err)
	}
	if len(report.Trials) != 2 || report.Trials[0].Order == report.Trials[1].Order {
		t.Fatalf("trials=%+v", report.Trials)
	}
	for _, trial := range report.Trials {
		if trial.Score.Verdict != "pass" {
			t.Fatalf("trial score=%+v", trial.Score)
		}
		if _, err := os.Stat(filepath.Join(output, filepath.FromSlash(trial.WorkspacePath), ".agents", "skills", "spectacular", "SKILL.md")); err != nil {
			t.Fatal(err)
		}
	}
	if report.Trials[0].WorkspacePath == report.Trials[1].WorkspacePath {
		t.Fatal("paired trials shared a workspace")
	}
	resumed, err := RunPaired(config)
	if err != nil {
		t.Fatal("resume failed:", err)
	}
	if len(resumed.Trials) != 2 || TrialOrder(resumed)[0] != TrialOrder(report)[0] {
		t.Fatalf("resumed=%+v", resumed.Trials)
	}
	writeTestFile(t, schemaPath, "{\"changed\":true}\n")
	if _, err := RunPaired(config); err == nil || !strings.Contains(err.Error(), "manifest does not match") {
		t.Fatalf("expected resume contamination refusal, got %v", err)
	}
}

func TestRunnerRefusesMutableRevision(t *testing.T) {
	if _, err := ResolveCommit(".", "working-tree-does-not-exist"); err == nil {
		t.Fatal("expected mutable or missing revision refusal")
	}
}
