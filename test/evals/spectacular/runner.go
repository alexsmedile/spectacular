package spectaculareval

import (
	"context"
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
	Repo                      string
	CatalogPath               string
	SchemaPath                string
	BaselineRef               string
	BaselineMode              string
	CandidateRef              string
	CandidateMode             string
	Tier                      string
	Repeats                   int
	Seed                      int64
	Model                     string
	Adapter                   string
	AdapterArgs               []string
	SpectacularCLI            string
	OutputDir                 string
	AllowHeldOut              bool
	ReadIsolation             string
	MaxCalls                  int
	TrialTimeout              time.Duration
	RequireCertifiedTelemetry bool
}

type trialSpec struct {
	Case     Case
	Variant  string
	Revision string
	Repeat   int
	Order    int
}

type runManifest struct {
	SchemaVersion             string                    `json:"schema_version"`
	BaselineRef               string                    `json:"baseline_ref"`
	BaselineMode              string                    `json:"baseline_mode"`
	BaselineCommit            string                    `json:"baseline_commit"`
	CandidateRef              string                    `json:"candidate_ref"`
	CandidateMode             string                    `json:"candidate_mode"`
	CandidateCommit           string                    `json:"candidate_commit"`
	CatalogDigest             string                    `json:"catalog_digest"`
	HarnessDigest             string                    `json:"harness_inputs_digest"`
	SchemaDigest              string                    `json:"result_schema_digest"`
	Adapter                   string                    `json:"adapter"`
	AdapterDigest             string                    `json:"adapter_digest"`
	AdapterArgs               []string                  `json:"adapter_args,omitempty"`
	SpectacularCLI            string                    `json:"spectacular_cli"`
	SpectacularCLIDigest      string                    `json:"spectacular_cli_digest"`
	ReadIsolation             string                    `json:"read_isolation"`
	TrialTimeoutMS            int64                     `json:"trial_timeout_ms"`
	RequireCertifiedTelemetry bool                      `json:"require_certified_telemetry"`
	Model                     string                    `json:"model"`
	Tier                      string                    `json:"tier"`
	Seed                      int64                     `json:"seed"`
	Planned                   []string                  `json:"planned"`
	Completed                 map[string]completedTrial `json:"completed"`
}

type completedTrial struct {
	TrialPath      string `json:"trial_path"`
	ArtifactDigest string `json:"artifact_digest"`
}

var agentResultRequiredFields = []string{"role", "phase", "status", "summary", "next_action", "owner_gate", "owner_questions", "references_loaded", "files_read", "commands_run", "safety_notes"}
var allowedAgentRoles = []string{"none", "Orchestrator", "Runner", "Reviewer", "Autopilot"}
var allowedAgentPhases = []string{"", "orient", "prepare", "execute", "runtime", "close", "audit"}
var allowedAgentStatuses = []string{"done", "blocked", "owner-gate", "draft-only", "not-invoked"}
var allowedVariantModes = []string{"skill", "workspace-only", "native-direct"}

func RunPaired(config RunConfig) (RunReport, error) {
	if strings.TrimSpace(config.OutputDir) == "" {
		return RunReport{}, errors.New("output directory is required")
	}
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
	config.OutputDir, err = absolutePath(config.OutputDir)
	if err != nil {
		return RunReport{}, fmt.Errorf("resolve output directory: %w", err)
	}
	adapterPath, err := exec.LookPath(config.Adapter)
	if err != nil {
		return RunReport{}, fmt.Errorf("resolve adapter executable: %w", err)
	}
	config.Adapter, err = absolutePath(adapterPath)
	if err != nil {
		return RunReport{}, fmt.Errorf("resolve adapter path: %w", err)
	}
	if strings.TrimSpace(config.SpectacularCLI) == "" {
		return RunReport{}, errors.New("pinned Spectacular CLI is required; pass --spectacular-cli")
	}
	config.SpectacularCLI, err = absolutePath(config.SpectacularCLI)
	if err != nil {
		return RunReport{}, fmt.Errorf("resolve pinned Spectacular CLI: %w", err)
	}
	for label, path := range map[string]string{"result schema": config.SchemaPath, "adapter": config.Adapter, "Spectacular CLI": config.SpectacularCLI} {
		if info, statErr := os.Stat(path); statErr != nil || info.IsDir() {
			return RunReport{}, fmt.Errorf("%s is not a readable file: %s", label, path)
		} else if label == "Spectacular CLI" && info.Mode()&0o111 == 0 {
			return RunReport{}, fmt.Errorf("%s is not executable: %s", label, path)
		}
	}
	catalog, err := LoadCatalog(config.CatalogPath)
	if err != nil {
		return RunReport{}, err
	}
	if config.ReadIsolation == "" {
		config.ReadIsolation = "artifact-only"
	}
	if config.BaselineMode == "" {
		config.BaselineMode = "skill"
	}
	if config.CandidateMode == "" {
		config.CandidateMode = "skill"
	}
	if !oneOf(config.BaselineMode, allowedVariantModes...) || !oneOf(config.CandidateMode, allowedVariantModes...) {
		return RunReport{}, fmt.Errorf("variant modes must be one of %v", allowedVariantModes)
	}
	if config.ReadIsolation != "artifact-only" && config.ReadIsolation != "os-enforced" {
		return RunReport{}, errors.New("read isolation must be artifact-only or os-enforced")
	}
	if err := validateAdapterIsolation(config.Adapter, config.ReadIsolation); err != nil {
		return RunReport{}, err
	}
	cases, defaultRepeats, err := CasesForTier(catalog, config.Tier)
	if err != nil {
		return RunReport{}, err
	}
	if config.Tier == "held-out" && !config.AllowHeldOut {
		return RunReport{}, errors.New("held-out tier requires explicit --allow-held-out; never use it while tuning")
	}
	if config.Tier == "held-out" && config.ReadIsolation != "os-enforced" {
		return RunReport{}, errors.New("held-out tier requires an externally OS-enforced read-isolation adapter")
	}
	if config.Repeats == 0 {
		config.Repeats = defaultRepeats
	}
	if config.Repeats < 1 {
		return RunReport{}, errors.New("repeats must be positive")
	}
	plannedCalls := len(cases) * config.Repeats * 2
	if err := validateCallBudget(plannedCalls, config.MaxCalls); err != nil {
		return RunReport{}, err
	}
	if config.TrialTimeout < 0 {
		return RunReport{}, errors.New("trial timeout cannot be negative")
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
		SchemaVersion:      "spectacular.skill-run-report.v1",
		BaselineRef:        config.BaselineRef,
		BaselineMode:       config.BaselineMode,
		CandidateRef:       config.CandidateRef,
		CandidateMode:      config.CandidateMode,
		Model:              config.Model,
		ReadIsolation:      config.ReadIsolation,
		Tier:               config.Tier,
		Seed:               seed,
		MinimumRepetitions: defaultRepeats,
		Thresholds:         catalog.Thresholds,
		Limitations: []string{
			"A conclusive verdict requires an adapter-authored spectacular.eval.observations event; model self-report remains visible but cannot establish semantic tool use.",
			"Token metrics use the maximum cumulative usage counters observed in trace events; adapters that emit per-turn rather than cumulative counters must normalize them before comparison.",
			"The harness isolates artifacts and exposes only one skill variant per trial; OS-level read isolation remains the adapter's responsibility.",
		},
	}
	for _, spec := range specs {
		id := trialID(spec)
		if completed, ok := manifest.Completed[id]; ok {
			trialPath, loadErr := containedPath(config.OutputDir, completed.TrialPath)
			if loadErr != nil {
				return RunReport{}, fmt.Errorf("resume %s: %w", id, loadErr)
			}
			trialDirectory := filepath.Dir(trialPath)
			digest, digestErr := directoryDigest(trialDirectory)
			if digestErr != nil || digest != completed.ArtifactDigest {
				return RunReport{}, fmt.Errorf("resume %s: artifact digest mismatch", id)
			}
			trial, loadErr := loadTrial(trialPath)
			if loadErr != nil {
				return RunReport{}, fmt.Errorf("resume %s: %w", id, loadErr)
			}
			mode := config.BaselineMode
			if spec.Variant == "candidate" {
				mode = config.CandidateMode
			}
			if loadErr = validateResumedTrial(trial, spec, commits[spec.Variant], config.Model, mode); loadErr != nil {
				return RunReport{}, fmt.Errorf("resume %s: %w", id, loadErr)
			}
			report.Trials = append(report.Trials, trial)
			if config.RequireCertifiedTelemetry && len(report.Trials) == 1 {
				if err := requireCertifiedTrial(trial); err != nil {
					return RunReport{}, fmt.Errorf("first-trial telemetry preflight: %w", err)
				}
			}
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
		artifactDigest, err := directoryDigest(filepath.Join(config.OutputDir, "trials", trial.ID))
		if err != nil {
			return RunReport{}, err
		}
		manifest.Completed[trial.ID] = completedTrial{TrialPath: relative, ArtifactDigest: artifactDigest}
		if err := writeRunManifest(config.OutputDir, manifest); err != nil {
			return RunReport{}, err
		}
		if config.RequireCertifiedTelemetry && len(report.Trials) == 1 {
			if err := requireCertifiedTrial(trial); err != nil {
				return RunReport{}, fmt.Errorf("first-trial telemetry preflight: %w; correct the adapter, then start a new output directory", err)
			}
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
	id := trialID(spec)
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
	mode := config.BaselineMode
	if spec.Variant == "candidate" {
		mode = config.CandidateMode
	}
	if err := prepareVariantWorkspace(config, workspace, mode, commit); err != nil {
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
	trialSchemaPath := filepath.Join(temporary, "result.schema.json")
	cliDirectory := filepath.Join(temporary, "bin")
	if err := stagePinnedCLI(config.SpectacularCLI, cliDirectory); err != nil {
		return Trial{}, err
	}
	prompt := variantPrompt(mode, spec.Case.Prompt)
	if err := os.WriteFile(promptPath, []byte(prompt+"\n"), 0o644); err != nil {
		return Trial{}, err
	}
	schemaData, err := os.ReadFile(config.SchemaPath)
	if err != nil {
		return Trial{}, err
	}
	if err := os.WriteFile(trialSchemaPath, schemaData, 0o644); err != nil {
		return Trial{}, err
	}
	started := time.Now().UTC()
	commandContext := context.Background()
	var cancel context.CancelFunc
	if config.TrialTimeout > 0 {
		ctx, stop := context.WithTimeout(commandContext, config.TrialTimeout)
		commandContext = ctx
		cancel = stop
	}
	command := exec.CommandContext(commandContext, config.Adapter, config.AdapterArgs...)
	configureProcessTree(command)
	if config.TrialTimeout > 0 {
		command.WaitDelay = time.Second
	}
	if cancel != nil {
		defer cancel()
	}
	command.Dir = workspace
	baseEnvironment := cleanEvalEnvironment(os.Environ(), config.Repo, config.OutputDir)
	baseEnvironment = withoutEnvironmentKey(baseEnvironment, "PATH")
	command.Env = append(baseEnvironment,
		"PATH="+cliDirectory+string(os.PathListSeparator)+os.Getenv("PATH"),
		"PWD="+workspace,
		"SPECTACULAR_EVAL_WORKSPACE="+workspace,
		"SPECTACULAR_EVAL_PROMPT="+promptPath,
		"SPECTACULAR_EVAL_RESULT="+resultPath,
		"SPECTACULAR_EVAL_TRACE="+tracePath,
		"SPECTACULAR_EVAL_SCHEMA="+trialSchemaPath,
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
	if errors.Is(commandContext.Err(), context.DeadlineExceeded) {
		cause := fmt.Errorf("adapter exceeded trial timeout %s", config.TrialTimeout)
		return Trial{}, persistFailedTrial(config.OutputDir, id, temporary, cause)
	}
	traceData, _ := os.ReadFile(tracePath)
	resultData, readErr := os.ReadFile(resultPath)
	if readErr != nil {
		cause := fmt.Errorf("adapter produced no result (exit=%d): %w", exitCode, readErr)
		return Trial{}, persistFailedTrial(config.OutputDir, id, temporary, cause)
	}
	result, err := decodeAgentResult(resultData)
	if err != nil {
		return Trial{}, persistFailedTrial(config.OutputDir, id, temporary, fmt.Errorf("decode adapter result: %w", err))
	}
	after, err := SnapshotTree(workspace)
	if err != nil {
		return Trial{}, err
	}
	changed := ChangedPaths(before, after)
	postconditions, err := runPostChecks(workspace, after, spec.Case.Expect.PostChecks)
	if err != nil {
		return Trial{}, err
	}
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
		ID:             id,
		CaseID:         spec.Case.ID,
		Suite:          suiteForCase(spec.Case),
		Complexity:     spec.Case.Complexity,
		Mode:           mode,
		Tags:           append([]string(nil), spec.Case.Tags...),
		Variant:        spec.Variant,
		Revision:       spec.Revision,
		Commit:         commit,
		Model:          config.Model,
		Repeat:         spec.Repeat,
		Order:          spec.Order,
		StartedAt:      started,
		DurationMS:     time.Since(started).Milliseconds(),
		ExitCode:       exitCode,
		Result:         result,
		ChangedPaths:   changed,
		Postconditions: postconditions,
		TraceMetrics:   ParseTraceMetrics(string(traceData)),
		TracePath:      relativeTrace,
		ResultPath:     relativeResult,
		WorkspacePath:  relativeWorkspace,
	}
	trial.Score = ScoreTrialWithPostconditions(spec.Case, result, string(traceData), changed, postconditions)
	if exitCode != 0 {
		trial.Score.SafetyPassed = false
		trial.Score.HardFailures = append(trial.Score.HardFailures, "adapter exited "+strconv.Itoa(exitCode))
		trial.Score.Verdict = "hard-fail"
		zero := 0.0
		trial.Score.Overall = &zero
	}
	return trial, nil
}

func requireCertifiedTrial(trial Trial) error {
	var findings []string
	if !trial.TraceMetrics.UsageObserved {
		findings = append(findings, "usage telemetry missing")
	}
	if !trial.TraceMetrics.SemanticObserved {
		findings = append(findings, "semantic host telemetry missing or self-reported only")
	}
	if len(findings) > 0 {
		return errors.New(strings.Join(findings, "; "))
	}
	return nil
}

func validateCallBudget(planned, maximum int) error {
	if maximum > 0 && planned > maximum {
		return fmt.Errorf("planned run requires %d model calls, exceeding --max-calls %d", planned, maximum)
	}
	return nil
}

func prepareVariantWorkspace(config RunConfig, workspace, mode, commit string) error {
	switch mode {
	case "skill":
		skillRoot := filepath.Join(workspace, ".agents", "skills", "spectacular")
		if _, err := MaterializeSkill(config.Repo, commit, skillRoot); err != nil {
			return err
		}
		_ = os.Remove(filepath.Join(workspace, "TASK.md"))
	case "workspace-only":
		_ = os.Remove(filepath.Join(workspace, "TASK.md"))
	case "native-direct":
		governanceRoot := filepath.Join(workspace, ".spectacular")
		if !strings.HasPrefix(governanceRoot, filepath.Clean(workspace)+string(os.PathSeparator)) {
			return errors.New("refuse unsafe native-control workspace transform")
		}
		if err := os.RemoveAll(governanceRoot); err != nil {
			return err
		}
	}
	return nil
}

func variantPrompt(mode, prompt string) string {
	switch mode {
	case "workspace-only":
		return "Use the canonical Markdown workspace records and folders directly. No Spectacular skill is installed; do not claim that it is.\n\n" + prompt
	case "native-direct":
		return "Use the host's normal direct execution behavior. Do not invoke or emulate Spectacular.\n\n" + prompt
	default:
		return prompt
	}
}

func validateAdapterIsolation(adapter, isolation string) error {
	if isolation != "os-enforced" {
		return nil
	}
	name := filepath.Base(adapter)
	for _, artifactOnly := range []string{"codex-adapter.sh", "claude-adapter.sh", "agy-adapter.sh", "opencode-adapter.sh"} {
		if name == artifactOnly {
			return fmt.Errorf("adapter %s declares artifact-only isolation and cannot be relabeled os-enforced", name)
		}
	}
	return nil
}

func persistFailedTrial(outputDir, id, temporary string, cause error) error {
	destination := filepath.Join(outputDir, "failed-trials", fmt.Sprintf("%s-%d", id, time.Now().UTC().UnixNano()))
	if err := os.MkdirAll(filepath.Dir(destination), 0o755); err != nil {
		return fmt.Errorf("%w; preserve failed trial: %v", cause, err)
	}
	if err := CopyTree(temporary, destination); err != nil {
		return fmt.Errorf("%w; preserve failed trial: %v", cause, err)
	}
	if err := os.WriteFile(filepath.Join(destination, "failure.txt"), []byte(cause.Error()+"\n"), 0o644); err != nil {
		return fmt.Errorf("%w; preserve failed trial: %v", cause, err)
	}
	return fmt.Errorf("%w; artifacts preserved at %s", cause, destination)
}

func cleanEvalEnvironment(environment []string, forbiddenRoots ...string) []string {
	cleaned := make([]string, 0, len(environment))
	for _, entry := range environment {
		name, value, _ := strings.Cut(entry, "=")
		if strings.HasPrefix(name, "SPECTACULAR_EVAL_") || name == "PWD" || name == "OLDPWD" || name == "INIT_CWD" {
			continue
		}
		exposesRoot := false
		for _, root := range forbiddenRoots {
			if root != "" && strings.Contains(value, root) {
				exposesRoot = true
				break
			}
		}
		if exposesRoot {
			continue
		}
		cleaned = append(cleaned, entry)
	}
	return cleaned
}

func withoutEnvironmentKey(environment []string, key string) []string {
	filtered := environment[:0]
	for _, entry := range environment {
		name, _, _ := strings.Cut(entry, "=")
		if name != key {
			filtered = append(filtered, entry)
		}
	}
	return filtered
}

func stagePinnedCLI(source, directory string) error {
	if err := os.MkdirAll(directory, 0o755); err != nil {
		return err
	}
	data, err := os.ReadFile(source)
	if err != nil {
		return fmt.Errorf("read pinned Spectacular CLI: %w", err)
	}
	if err := os.WriteFile(filepath.Join(directory, "spectacular"), data, 0o755); err != nil {
		return fmt.Errorf("stage pinned Spectacular CLI: %w", err)
	}
	return nil
}

func runPostChecks(workspace string, agentSnapshot map[string]string, checks []PostCheck) ([]PostconditionResult, error) {
	results := make([]PostconditionResult, 0, len(checks))
	for _, check := range checks {
		command := exec.Command(check.Command[0], check.Command[1:]...)
		command.Dir = workspace
		output, runErr := command.CombinedOutput()
		exitCode := 0
		if runErr != nil {
			exitCode = 1
			var exitErr *exec.ExitError
			if errors.As(runErr, &exitErr) {
				exitCode = exitErr.ExitCode()
			}
		}
		afterCheck, err := SnapshotTree(workspace)
		if err != nil {
			return nil, err
		}
		mutated := ChangedPaths(agentSnapshot, afterCheck)
		result := PostconditionResult{
			Command: append([]string(nil), check.Command...), ExpectedExit: check.ExpectedExit,
			ActualExit: exitCode, Output: truncate(string(output), 4096), MutatedPaths: mutated,
		}
		result.Passed = exitCode == check.ExpectedExit && len(mutated) == 0
		results = append(results, result)
	}
	return results, nil
}

func truncate(value string, maximum int) string {
	if len(value) <= maximum {
		return value
	}
	return value[:maximum] + "...[truncated]"
}

func decodeAgentResult(data []byte) (AgentResult, error) {
	var raw map[string]json.RawMessage
	if err := json.Unmarshal(data, &raw); err != nil {
		return AgentResult{}, err
	}
	for _, field := range agentResultRequiredFields {
		if _, ok := raw[field]; !ok {
			return AgentResult{}, fmt.Errorf("missing required field %q", field)
		}
	}
	var result AgentResult
	decoder := json.NewDecoder(strings.NewReader(string(data)))
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(&result); err != nil {
		return AgentResult{}, err
	}
	if !oneOf(result.Role, allowedAgentRoles...) {
		return AgentResult{}, fmt.Errorf("unknown role %q", result.Role)
	}
	if !oneOf(result.Phase, allowedAgentPhases...) {
		return AgentResult{}, fmt.Errorf("unknown phase %q", result.Phase)
	}
	if !oneOf(result.Status, allowedAgentStatuses...) {
		return AgentResult{}, fmt.Errorf("unknown status %q", result.Status)
	}
	return result, nil
}

func oneOf(value string, allowed ...string) bool {
	for _, candidate := range allowed {
		if value == candidate {
			return true
		}
	}
	return false
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
	cliDigest, err := fileDigest(config.SpectacularCLI)
	if err != nil {
		return runManifest{}, err
	}
	harnessDigest, err := benchmarkInputsDigest(filepath.Dir(config.CatalogPath), config.OutputDir)
	if err != nil {
		return runManifest{}, err
	}
	data, err := os.ReadFile(path)
	if err == nil {
		var manifest runManifest
		if json.Unmarshal(data, &manifest) != nil {
			return runManifest{}, errors.New("existing run manifest is invalid")
		}
		if manifest.SchemaVersion != "spectacular.skill-run-manifest.v1" || manifest.BaselineRef != config.BaselineRef || manifest.BaselineMode != config.BaselineMode || manifest.BaselineCommit != commits["baseline"] || manifest.CandidateRef != config.CandidateRef || manifest.CandidateMode != config.CandidateMode || manifest.CandidateCommit != commits["candidate"] || manifest.CatalogDigest != catalogDigest || manifest.HarnessDigest != harnessDigest || manifest.SchemaDigest != schemaDigest || manifest.Adapter != config.Adapter || manifest.AdapterDigest != adapterDigest || strings.Join(manifest.AdapterArgs, "\x00") != strings.Join(config.AdapterArgs, "\x00") || manifest.SpectacularCLI != config.SpectacularCLI || manifest.SpectacularCLIDigest != cliDigest || manifest.ReadIsolation != config.ReadIsolation || manifest.TrialTimeoutMS != config.TrialTimeout.Milliseconds() || manifest.RequireCertifiedTelemetry != config.RequireCertifiedTelemetry || manifest.Model != config.Model || manifest.Tier != config.Tier || manifest.Seed != seed || strings.Join(manifest.Planned, "\n") != strings.Join(planned, "\n") {
			return runManifest{}, errors.New("existing run manifest does not match requested comparison")
		}
		if manifest.Completed == nil {
			manifest.Completed = map[string]completedTrial{}
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
		BaselineRef:   config.BaselineRef, BaselineMode: config.BaselineMode, BaselineCommit: commits["baseline"],
		CandidateRef: config.CandidateRef, CandidateMode: config.CandidateMode, CandidateCommit: commits["candidate"],
		CatalogDigest: catalogDigest, HarnessDigest: harnessDigest, SchemaDigest: schemaDigest,
		Adapter: config.Adapter, AdapterDigest: adapterDigest, AdapterArgs: append([]string(nil), config.AdapterArgs...),
		SpectacularCLI: config.SpectacularCLI, SpectacularCLIDigest: cliDigest,
		ReadIsolation:  config.ReadIsolation,
		TrialTimeoutMS: config.TrialTimeout.Milliseconds(), RequireCertifiedTelemetry: config.RequireCertifiedTelemetry,
		Model: config.Model, Tier: config.Tier, Seed: seed, Planned: planned,
		Completed: map[string]completedTrial{},
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

func benchmarkInputsDigest(root, outputDir string) (string, error) {
	root = filepath.Clean(root)
	outputDir = filepath.Clean(outputDir)
	var paths []string
	err := filepath.WalkDir(root, func(path string, entry os.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if entry.IsDir() {
			if path == outputDir {
				return filepath.SkipDir
			}
			if path != root && (entry.Name() == "reports" || entry.Name() == "testdata") {
				return filepath.SkipDir
			}
			return nil
		}
		if entry.Name() == "README.md" || entry.Name() == ".DS_Store" {
			return nil
		}
		paths = append(paths, path)
		return nil
	})
	if err != nil {
		return "", err
	}
	sort.Strings(paths)
	hash := sha256.New()
	for _, path := range paths {
		relative, _ := filepath.Rel(root, path)
		data, err := os.ReadFile(path)
		if err != nil {
			return "", err
		}
		fmt.Fprintf(hash, "%s\x00%d\x00", filepath.ToSlash(relative), len(data))
		_, _ = hash.Write(data)
	}
	return fmt.Sprintf("sha256:%x", hash.Sum(nil)), nil
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

func containedPath(root, relative string) (string, error) {
	if filepath.IsAbs(relative) {
		return "", errors.New("artifact path must be relative")
	}
	root = filepath.Clean(root)
	path := filepath.Clean(filepath.Join(root, filepath.FromSlash(relative)))
	if path == root || !strings.HasPrefix(path, root+string(filepath.Separator)) {
		return "", errors.New("artifact path escapes output directory")
	}
	return path, nil
}

func validateResumedTrial(trial Trial, spec trialSpec, commit, model, mode string) error {
	if trial.ID != trialID(spec) || trial.CaseID != spec.Case.ID || trial.Suite != suiteForCase(spec.Case) || trial.Complexity != spec.Case.Complexity || trial.Mode != mode || trial.Variant != spec.Variant || trial.Revision != spec.Revision || trial.Commit != commit || trial.Model != model || trial.Repeat != spec.Repeat || trial.Order != spec.Order {
		return errors.New("trial identity does not match the run plan")
	}
	return nil
}

func directoryDigest(root string) (string, error) {
	var paths []string
	err := filepath.WalkDir(root, func(path string, entry os.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if path == root {
			return nil
		}
		if entry.Type()&os.ModeSymlink != 0 {
			return fmt.Errorf("symlink in trial artifacts: %s", path)
		}
		if !entry.IsDir() {
			paths = append(paths, path)
		}
		return nil
	})
	if err != nil {
		return "", err
	}
	sort.Strings(paths)
	hash := sha256.New()
	for _, path := range paths {
		relative, _ := filepath.Rel(root, path)
		data, err := os.ReadFile(path)
		if err != nil {
			return "", err
		}
		info, err := os.Stat(path)
		if err != nil {
			return "", err
		}
		fmt.Fprintf(hash, "%s\x00%o\x00%d\x00", filepath.ToSlash(relative), info.Mode().Perm(), len(data))
		_, _ = hash.Write(data)
	}
	return fmt.Sprintf("sha256:%x", hash.Sum(nil)), nil
}

func initializeFixtureGit(workspace string) error {
	commands := [][]string{
		{"init", "-q", "-b", "main"},
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
