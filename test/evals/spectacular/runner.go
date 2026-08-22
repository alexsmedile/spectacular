package spectaculareval

import (
	"crypto/sha256"
	"encoding/json"
	"errors"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"time"
)

type RunConfig struct {
	Repo         string
	CatalogPath  string
	SchemaPath   string
	BaselineRef  string
	CandidateRef string
	Tier         string
	Repeats      int
	Seed         int64
	Model        string
	Adapter      string
	AdapterArgs  []string
	OutputDir    string
}

type trialSpec struct {
	Case     Case
	Variant  string
	Revision string
	Repeat   int
	Order    int
}

type runManifest struct {
	SchemaVersion   string            `json:"schema_version"`
	BaselineRef     string            `json:"baseline_ref"`
	BaselineCommit  string            `json:"baseline_commit"`
	CandidateRef    string            `json:"candidate_ref"`
	CandidateCommit string            `json:"candidate_commit"`
	CatalogDigest   string            `json:"catalog_digest"`
	SchemaDigest    string            `json:"result_schema_digest"`
	Adapter         string            `json:"adapter"`
	AdapterDigest   string            `json:"adapter_digest"`
	AdapterArgs     []string          `json:"adapter_args,omitempty"`
	Model           string            `json:"model"`
	Tier            string            `json:"tier"`
	Seed            int64             `json:"seed"`
	Planned         []string          `json:"planned"`
	Completed       map[string]string `json:"completed"`
}

func RunPaired(config RunConfig) (RunReport, error) {
	var err error
	config.Repo, err = absolutePath(config.Repo)
	if err != nil {
		return RunReport{}, fmt.Errorf("resolve repository path: %w", err)
	}
	config.CatalogPath, err = absolutePath(config.CatalogPath)
	if err != nil {
		return RunReport{}, fmt.Errorf("resolve catalog path: %w", err)
	}
	config.SchemaPath, err = absolutePath(config.SchemaPath)
	if err != nil {
		return RunReport{}, fmt.Errorf("resolve result schema path: %w", err)
	}
	adapterPath, err := exec.LookPath(config.Adapter)
	if err != nil {
		return RunReport{}, fmt.Errorf("resolve adapter executable: %w", err)
	}
	config.Adapter, err = absolutePath(adapterPath)
	if err != nil {
		return RunReport{}, fmt.Errorf("resolve adapter path: %w", err)
	}
	for label, path := range map[string]string{"result schema": config.SchemaPath, "adapter": config.Adapter} {
		if info, statErr := os.Stat(path); statErr != nil || info.IsDir() {
			return RunReport{}, fmt.Errorf("%s is not a readable file: %s", label, path)
		}
	}
	catalog, err := LoadCatalog(config.CatalogPath)
	if err != nil {
		return RunReport{}, err
	}
	cases, defaultRepeats, err := CasesForTier(catalog, config.Tier)
	if err != nil {
		return RunReport{}, err
	}
	if config.Repeats == 0 {
		config.Repeats = defaultRepeats
	}
	if config.Repeats < 1 {
		return RunReport{}, errors.New("repeats must be positive")
	}
	if config.Model == "" || config.Adapter == "" {
		return RunReport{}, errors.New("model and adapter are required")
	}
	baselineCommit, err := ResolveCommit(config.Repo, config.BaselineRef)
	if err != nil {
		return RunReport{}, err
	}
	candidateCommit, err := ResolveCommit(config.Repo, config.CandidateRef)
	if err != nil {
		return RunReport{}, err
	}
	seed := config.Seed
	if seed == 0 {
		seed = 1
	}
	random := rand.New(rand.NewSource(seed))
	var specs []trialSpec
	order := 0
	for _, item := range cases {
		for repeat := 1; repeat <= config.Repeats; repeat++ {
			pair := []trialSpec{
				{Case: item, Variant: "baseline", Revision: config.BaselineRef, Repeat: repeat},
				{Case: item, Variant: "candidate", Revision: config.CandidateRef, Repeat: repeat},
			}
			random.Shuffle(len(pair), func(i, j int) { pair[i], pair[j] = pair[j], pair[i] })
			for index := range pair {
				order++
				pair[index].Order = order
				specs = append(specs, pair[index])
			}
		}
	}
	commits := map[string]string{"baseline": baselineCommit, "candidate": candidateCommit}
	manifest, err := openRunManifest(config, seed, specs, commits)
	if err != nil {
		return RunReport{}, err
	}
	report := RunReport{
		SchemaVersion: "spectacular.skill-run-report.v1",
		BaselineRef:   config.BaselineRef,
		CandidateRef:  config.CandidateRef,
		Model:         config.Model,
		Tier:          config.Tier,
		Seed:          seed,
		Thresholds:    catalog.Thresholds,
		Limitations: []string{
			"File-read metrics combine structured self-report with observable adapter traces; an adapter that omits tool events lowers confidence.",
			"The harness isolates artifacts and exposes only one skill variant per trial; OS-level read isolation remains the adapter's responsibility.",
		},
	}
	for _, spec := range specs {
		id := trialID(spec)
		if relative, ok := manifest.Completed[id]; ok {
			trial, loadErr := loadTrial(filepath.Join(config.OutputDir, filepath.FromSlash(relative)))
			if loadErr != nil {
				return RunReport{}, fmt.Errorf("resume %s: %w", id, loadErr)
			}
			report.Trials = append(report.Trials, trial)
			continue
		}
		trial, runErr := runOne(config, spec, commits[spec.Variant])
		if runErr != nil {
			return RunReport{}, fmt.Errorf("run %s/%s repeat %d: %w", spec.Case.ID, spec.Variant, spec.Repeat, runErr)
		}
		report.Trials = append(report.Trials, trial)
		relative := filepath.ToSlash(filepath.Join("trials", trial.ID, "trial.json"))
		if err := writeJSON(filepath.Join(config.OutputDir, filepath.FromSlash(relative)), trial); err != nil {
			return RunReport{}, err
		}
		manifest.Completed[trial.ID] = relative
		if err := writeRunManifest(config.OutputDir, manifest); err != nil {
			return RunReport{}, err
		}
	}
	Summarize(&report)
	return report, nil
}

func absolutePath(path string) (string, error) {
	if filepath.IsAbs(path) {
		return filepath.Clean(path), nil
	}
	return filepath.Abs(path)
}

func runOne(config RunConfig, spec trialSpec, commit string) (Trial, error) {
	temporary, err := os.MkdirTemp("", "spectacular-eval-trial-")
	if err != nil {
		return Trial{}, err
	}
	defer os.RemoveAll(temporary)
	workspace := filepath.Join(temporary, "workspace")
	fixtureRoot := filepath.Join(filepath.Dir(config.CatalogPath), "fixtures", spec.Case.Fixture)
	if err := CopyTree(fixtureRoot, workspace); err != nil {
		return Trial{}, fmt.Errorf("copy fixture: %w", err)
	}
	skillRoot := filepath.Join(workspace, ".agents", "skills", "spectacular")
	if _, err := MaterializeSkill(config.Repo, commit, skillRoot); err != nil {
		return Trial{}, err
	}
	if err := initializeFixtureGit(workspace); err != nil {
		return Trial{}, err
	}
	before, err := SnapshotTree(workspace)
	if err != nil {
		return Trial{}, err
	}
	promptPath := filepath.Join(temporary, "prompt.md")
	resultPath := filepath.Join(temporary, "result.json")
	tracePath := filepath.Join(temporary, "trace.jsonl")
	if err := os.WriteFile(promptPath, []byte(spec.Case.Prompt+"\n"), 0o644); err != nil {
		return Trial{}, err
	}
	started := time.Now().UTC()
	command := exec.Command(config.Adapter, config.AdapterArgs...)
	command.Dir = workspace
	command.Env = append(os.Environ(),
		"SPECTACULAR_EVAL_WORKSPACE="+workspace,
		"SPECTACULAR_EVAL_PROMPT="+promptPath,
		"SPECTACULAR_EVAL_RESULT="+resultPath,
		"SPECTACULAR_EVAL_TRACE="+tracePath,
		"SPECTACULAR_EVAL_SCHEMA="+config.SchemaPath,
		"SPECTACULAR_EVAL_MODEL="+config.Model,
		"SPECTACULAR_EVAL_CASE="+spec.Case.ID,
		"SPECTACULAR_EVAL_KIND="+spec.Case.Kind,
	)
	for key, value := range spec.Case.Environment {
		command.Env = append(command.Env, key+"="+value)
	}
	output, commandErr := command.CombinedOutput()
	exitCode := 0
	if commandErr != nil {
		exitCode = 1
		var exitErr *exec.ExitError
		if errors.As(commandErr, &exitErr) {
			exitCode = exitErr.ExitCode()
		}
	}
	if len(output) > 0 {
		file, openErr := os.OpenFile(tracePath, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0o644)
		if openErr != nil {
			return Trial{}, openErr
		}
		_, _ = file.Write(append([]byte("\nADAPTER_STDERR\n"), output...))
		_ = file.Close()
	}
	traceData, _ := os.ReadFile(tracePath)
	resultData, readErr := os.ReadFile(resultPath)
	if readErr != nil {
		return Trial{}, fmt.Errorf("adapter produced no result (exit=%d): %w", exitCode, readErr)
	}
	var result AgentResult
	decoder := json.NewDecoder(strings.NewReader(string(resultData)))
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(&result); err != nil {
		return Trial{}, fmt.Errorf("decode adapter result: %w", err)
	}
	after, err := SnapshotTree(workspace)
	if err != nil {
		return Trial{}, err
	}
	changed := ChangedPaths(before, after)
	id := trialID(spec)
	destination := filepath.Join(config.OutputDir, "trials", id)
	if err := os.MkdirAll(filepath.Dir(destination), 0o755); err != nil {
		return Trial{}, err
	}
	if err := CopyTree(temporary, destination); err != nil {
		return Trial{}, fmt.Errorf("persist trial artifacts: %w", err)
	}
	relativeTrace := filepath.ToSlash(filepath.Join("trials", id, "trace.jsonl"))
	relativeResult := filepath.ToSlash(filepath.Join("trials", id, "result.json"))
	relativeWorkspace := filepath.ToSlash(filepath.Join("trials", id, "workspace"))
	trial := Trial{
		ID:            id,
		CaseID:        spec.Case.ID,
		Tags:          append([]string(nil), spec.Case.Tags...),
		Variant:       spec.Variant,
		Revision:      spec.Revision,
		Commit:        commit,
		Model:         config.Model,
		Repeat:        spec.Repeat,
		Order:         spec.Order,
		StartedAt:     started,
		DurationMS:    time.Since(started).Milliseconds(),
		ExitCode:      exitCode,
		Result:        result,
		ChangedPaths:  changed,
		TraceMetrics:  ParseTraceMetrics(string(traceData)),
		TracePath:     relativeTrace,
		ResultPath:    relativeResult,
		WorkspacePath: relativeWorkspace,
	}
	trial.Score = ScoreTrial(spec.Case, result, string(traceData), changed)
	if exitCode != 0 {
		trial.Score.SafetyPassed = false
		trial.Score.HardFailures = append(trial.Score.HardFailures, "adapter exited "+strconv.Itoa(exitCode))
		trial.Score.Verdict = "hard-fail"
		zero := 0.0
		trial.Score.Overall = &zero
	}
	return trial, nil
}

func prepareOutputDirectory(path string) error {
	if strings.TrimSpace(path) == "" {
		return errors.New("output directory is required")
	}
	entries, err := os.ReadDir(path)
	if err == nil && len(entries) > 0 {
		return fmt.Errorf("output directory must be absent or empty: %s", path)
	}
	if err != nil && !os.IsNotExist(err) {
		return err
	}
	return os.MkdirAll(path, 0o755)
}

func trialID(spec trialSpec) string {
	return fmt.Sprintf("%s-r%02d-%s", strings.ToLower(spec.Case.ID), spec.Repeat, spec.Variant)
}

func openRunManifest(config RunConfig, seed int64, specs []trialSpec, commits map[string]string) (runManifest, error) {
	if strings.TrimSpace(config.OutputDir) == "" {
		return runManifest{}, errors.New("output directory is required")
	}
	if err := os.MkdirAll(config.OutputDir, 0o755); err != nil {
		return runManifest{}, err
	}
	path := filepath.Join(config.OutputDir, "run-manifest.json")
	planned := make([]string, len(specs))
	for index, spec := range specs {
		planned[index] = trialID(spec)
	}
	catalogDigest, err := fileDigest(config.CatalogPath)
	if err != nil {
		return runManifest{}, err
	}
	schemaDigest, err := fileDigest(config.SchemaPath)
	if err != nil {
		return runManifest{}, err
	}
	adapterDigest, err := fileDigest(config.Adapter)
	if err != nil {
		return runManifest{}, err
	}
	data, err := os.ReadFile(path)
	if err == nil {
		var manifest runManifest
		if json.Unmarshal(data, &manifest) != nil {
			return runManifest{}, errors.New("existing run manifest is invalid")
		}
		if manifest.SchemaVersion != "spectacular.skill-run-manifest.v1" || manifest.BaselineRef != config.BaselineRef || manifest.BaselineCommit != commits["baseline"] || manifest.CandidateRef != config.CandidateRef || manifest.CandidateCommit != commits["candidate"] || manifest.CatalogDigest != catalogDigest || manifest.SchemaDigest != schemaDigest || manifest.Adapter != config.Adapter || manifest.AdapterDigest != adapterDigest || strings.Join(manifest.AdapterArgs, "\x00") != strings.Join(config.AdapterArgs, "\x00") || manifest.Model != config.Model || manifest.Tier != config.Tier || manifest.Seed != seed || strings.Join(manifest.Planned, "\n") != strings.Join(planned, "\n") {
			return runManifest{}, errors.New("existing run manifest does not match requested comparison")
		}
		if manifest.Completed == nil {
			manifest.Completed = map[string]string{}
		}
		return manifest, nil
	}
	if !os.IsNotExist(err) {
		return runManifest{}, err
	}
	entries, readErr := os.ReadDir(config.OutputDir)
	if readErr != nil {
		return runManifest{}, readErr
	}
	if len(entries) > 0 {
		return runManifest{}, fmt.Errorf("output directory has no resumable manifest: %s", config.OutputDir)
	}
	manifest := runManifest{
		SchemaVersion: "spectacular.skill-run-manifest.v1",
		BaselineRef:   config.BaselineRef, BaselineCommit: commits["baseline"],
		CandidateRef: config.CandidateRef, CandidateCommit: commits["candidate"],
		CatalogDigest: catalogDigest, SchemaDigest: schemaDigest,
		Adapter: config.Adapter, AdapterDigest: adapterDigest, AdapterArgs: append([]string(nil), config.AdapterArgs...),
		Model: config.Model, Tier: config.Tier, Seed: seed, Planned: planned,
		Completed: map[string]string{},
	}
	if err := writeRunManifest(config.OutputDir, manifest); err != nil {
		return runManifest{}, err
	}
	return manifest, nil
}

func fileDigest(path string) (string, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return "", fmt.Errorf("digest %s: %w", path, err)
	}
	digest := sha256.Sum256(data)
	return fmt.Sprintf("sha256:%x", digest[:]), nil
}

func writeRunManifest(directory string, manifest runManifest) error {
	path := filepath.Join(directory, "run-manifest.json")
	temporary := path + ".tmp"
	if err := writeJSON(temporary, manifest); err != nil {
		return err
	}
	return os.Rename(temporary, path)
}

func loadTrial(path string) (Trial, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return Trial{}, err
	}
	var trial Trial
	if err := json.Unmarshal(data, &trial); err != nil {
		return Trial{}, err
	}
	return trial, nil
}

func initializeFixtureGit(workspace string) error {
	commands := [][]string{
		{"init", "-q"},
		{"config", "user.name", "Spectacular Eval"},
		{"config", "user.email", "eval@example.invalid"},
		{"add", "."},
		{"commit", "-qm", "fixture baseline"},
	}
	for _, arguments := range commands {
		command := exec.Command("git", append([]string{"-C", workspace}, arguments...)...)
		if output, err := command.CombinedOutput(); err != nil {
			return fmt.Errorf("git %s: %w: %s", strings.Join(arguments, " "), err, strings.TrimSpace(string(output)))
		}
	}
	return nil
}

func TrialOrder(report RunReport) []string {
	trials := append([]Trial(nil), report.Trials...)
	sort.Slice(trials, func(i, j int) bool { return trials[i].Order < trials[j].Order })
	result := make([]string, len(trials))
	for index, trial := range trials {
		result[index] = trial.ID
	}
	return result
}
