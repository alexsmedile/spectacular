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
	tracePath := filepath.Join(output, "trials", report.Trials[0].ID, "trace.jsonl")
	originalTrace, err := os.ReadFile(tracePath)
	if err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(tracePath, append(originalTrace, []byte("tampered\n")...), 0o644); err != nil {
		t.Fatal(err)
	}
	if _, err := RunPaired(config); err == nil || !strings.Contains(err.Error(), "artifact digest mismatch") {
		t.Fatalf("expected artifact tampering refusal, got %v", err)
	}
	if err := os.WriteFile(tracePath, originalTrace, 0o644); err != nil {
		t.Fatal(err)
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

func TestDecodeAgentResultEnforcesAdapterContract(t *testing.T) {
	valid := `{"role":"Runner","phase":"","status":"done","summary":"ok","next_action":"return","owner_gate":"","owner_questions":[],"references_loaded":[],"files_read":[],"commands_run":[],"safety_notes":[]}`
	if _, err := decodeAgentResult([]byte(valid)); err != nil {
		t.Fatal(err)
	}
	for name, malformed := range map[string]string{
		"missing": `{"role":"Runner"}`,
		"role":    strings.Replace(valid, `"Runner"`, `"Superuser"`, 1),
		"status":  strings.Replace(valid, `"done"`, `"complete"`, 1),
		"unknown": strings.TrimSuffix(valid, "}") + `,"extra":true}`,
	} {
		t.Run(name, func(t *testing.T) {
			if _, err := decodeAgentResult([]byte(malformed)); err == nil {
				t.Fatal("expected malformed result refusal")
			}
		})
	}
}

func TestAgentResultSchemaMatchesRuntimeDecoder(t *testing.T) {
	data, err := os.ReadFile("agent-result.schema.json")
	if err != nil {
		t.Fatal(err)
	}
	var schema struct {
		Required   []string `json:"required"`
		Properties map[string]struct {
			Enum []string `json:"enum"`
		} `json:"properties"`
	}
	if err := json.Unmarshal(data, &schema); err != nil {
		t.Fatal(err)
	}
	for label, pair := range map[string][2][]string{
		"required": {schema.Required, agentResultRequiredFields},
		"roles":    {schema.Properties["role"].Enum, allowedAgentRoles},
		"phases":   {schema.Properties["phase"].Enum, allowedAgentPhases},
		"statuses": {schema.Properties["status"].Enum, allowedAgentStatuses},
	} {
		if strings.Join(pair[0], "\x00") != strings.Join(pair[1], "\x00") {
			t.Errorf("%s schema=%v runtime=%v", label, pair[0], pair[1])
		}
	}
}

func TestCleanEvalEnvironmentRemovesRepositoryPathLeaks(t *testing.T) {
	cleaned := cleanEvalEnvironment([]string{
		"PATH=/usr/bin", "PWD=/repo", "OLDPWD=/repo/parent", "INIT_CWD=/repo", "SPECTACULAR_EVAL_CASE=old",
		"SAFE=value", "LEAK=/repo/reports",
	}, "/repo")
	joined := strings.Join(cleaned, "\n")
	if joined != "PATH=/usr/bin\nSAFE=value" {
		t.Fatalf("cleaned=%q", joined)
	}
}

func TestPostChecksVerifyOutcomeAndRefuseMutatingChecks(t *testing.T) {
	workspace := t.TempDir()
	writeTestFile(t, filepath.Join(workspace, "value.txt"), "BAD\n")
	before, err := SnapshotTree(workspace)
	if err != nil {
		t.Fatal(err)
	}
	failed, err := runPostChecks(workspace, before, []PostCheck{{Command: []string{"grep", "-q", "OK", "value.txt"}, ExpectedExit: 0}})
	if err != nil || len(failed) != 1 || failed[0].Passed {
		t.Fatalf("failed=%+v err=%v", failed, err)
	}
	writeTestFile(t, filepath.Join(workspace, "value.txt"), "OK\n")
	before, _ = SnapshotTree(workspace)
	passed, err := runPostChecks(workspace, before, []PostCheck{{Command: []string{"grep", "-q", "OK", "value.txt"}, ExpectedExit: 0}})
	if err != nil || !passed[0].Passed {
		t.Fatalf("passed=%+v err=%v", passed, err)
	}
	mutated, err := runPostChecks(workspace, before, []PostCheck{{Command: []string{"sh", "-c", "printf MUTATED > value.txt"}, ExpectedExit: 0}})
	if err != nil || mutated[0].Passed || len(mutated[0].MutatedPaths) == 0 {
		t.Fatalf("mutated=%+v err=%v", mutated, err)
	}
}
