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
	uncertifiedAdapter := filepath.Join(t.TempDir(), "uncertified-adapter.sh")
	writeTestFile(t, uncertifiedAdapter, `#!/bin/sh
set -eu
printf '%s\n' '{"role":"Orchestrator","phase":"orient","status":"done","summary":"ok","next_action":"return","owner_gate":"","owner_questions":[],"references_loaded":["orient.md"],"files_read":[],"commands_run":[],"safety_notes":[]}' > "$SPECTACULAR_EVAL_RESULT"
printf '%s\n' '{"type":"spectacular.eval.usage","input_tokens":10,"output_tokens":1}' > "$SPECTACULAR_EVAL_TRACE"
`)
	if err := os.Chmod(uncertifiedAdapter, 0o755); err != nil {
		t.Fatal(err)
	}
	uncertified := config
	uncertified.Adapter = uncertifiedAdapter
	uncertified.OutputDir = t.TempDir()
	uncertified.RequireCertifiedTelemetry = true
	if _, err := RunPaired(uncertified); err == nil || !strings.Contains(err.Error(), "first-trial telemetry preflight") {
		t.Fatalf("expected first-trial telemetry refusal, got %v", err)
	}
	trialDirs, err := os.ReadDir(filepath.Join(uncertified.OutputDir, "trials"))
	if err != nil || len(trialDirs) != 1 {
		t.Fatalf("first-trial gate spent beyond one call: dirs=%v err=%v", trialDirs, err)
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

func TestFailedTrialArtifactsSurviveAdapterFailure(t *testing.T) {
	output := t.TempDir()
	temporary := t.TempDir()
	writeTestFile(t, filepath.Join(temporary, "trace.jsonl"), "adapter stderr\n")
	err := persistFailedTrial(output, "aa-r01-baseline", temporary, os.ErrInvalid)
	if err == nil || !strings.Contains(err.Error(), "artifacts preserved at") {
		t.Fatalf("error=%v", err)
	}
	matches, globErr := filepath.Glob(filepath.Join(output, "failed-trials", "aa-r01-baseline-*", "trace.jsonl"))
	if globErr != nil || len(matches) != 1 {
		t.Fatalf("matches=%v err=%v", matches, globErr)
	}
	if _, statErr := os.Stat(filepath.Join(filepath.Dir(matches[0]), "failure.txt")); statErr != nil {
		t.Fatal(statErr)
	}
}

func TestRunSafetyControlsRefuseExcessCallsAndUncertifiedTelemetry(t *testing.T) {
	if err := validateCallBudget(13, 12); err == nil || !strings.Contains(err.Error(), "13 model calls") {
		t.Fatalf("expected call-budget refusal, got %v", err)
	}
	if err := validateCallBudget(12, 12); err != nil {
		t.Fatal(err)
	}
	trial := Trial{TraceMetrics: TraceMetrics{UsageObserved: true}}
	if err := requireCertifiedTrial(trial); err == nil || !strings.Contains(err.Error(), "semantic host telemetry") {
		t.Fatalf("expected semantic telemetry refusal, got %v", err)
	}
	trial.TraceMetrics.SemanticObserved = true
	if err := requireCertifiedTrial(trial); err != nil {
		t.Fatal(err)
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

func TestControlModesExposeEquivalentInformationWithoutInstallingSkill(t *testing.T) {
	workspace := t.TempDir()
	writeTestFile(t, filepath.Join(workspace, "TASK.md"), "neutral task\n")
	writeTestFile(t, filepath.Join(workspace, ".spectacular", "PROJECT.md"), "canonical task\n")
	if err := prepareVariantWorkspace(RunConfig{}, workspace, "native-direct", ""); err != nil {
		t.Fatal(err)
	}
	if _, err := os.Stat(filepath.Join(workspace, ".spectacular")); !os.IsNotExist(err) {
		t.Fatalf("native control retained governance workspace: %v", err)
	}
	if _, err := os.Stat(filepath.Join(workspace, "TASK.md")); err != nil {
		t.Fatal(err)
	}

	workspace = t.TempDir()
	writeTestFile(t, filepath.Join(workspace, "TASK.md"), "neutral task\n")
	writeTestFile(t, filepath.Join(workspace, ".spectacular", "PROJECT.md"), "canonical task\n")
	if err := prepareVariantWorkspace(RunConfig{}, workspace, "workspace-only", ""); err != nil {
		t.Fatal(err)
	}
	if _, err := os.Stat(filepath.Join(workspace, "TASK.md")); !os.IsNotExist(err) {
		t.Fatalf("workspace-only control retained neutral projection: %v", err)
	}
	if _, err := os.Stat(filepath.Join(workspace, ".spectacular", "PROJECT.md")); err != nil {
		t.Fatal(err)
	}
	if oneOf("native-plan", allowedVariantModes...) {
		t.Fatal("prompt-only native-plan must not be advertised as an attributable control")
	}
}

func TestKnownArtifactOnlyAdaptersCannotClaimOSEnforcement(t *testing.T) {
	for _, name := range []string{"codex-adapter.sh", "claude-adapter.sh", "agy-adapter.sh", "opencode-adapter.sh"} {
		if err := validateAdapterIsolation(filepath.Join("scripts", name), "os-enforced"); err == nil {
			t.Fatalf("%s was relabeled os-enforced", name)
		}
	}
	if err := validateAdapterIsolation("external-container-adapter", "os-enforced"); err != nil {
		t.Fatal(err)
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
