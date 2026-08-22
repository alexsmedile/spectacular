package spectaculareval

import "time"

const CatalogSchema = "spectacular.skill-evals.v1"

var Dimensions = []string{"safety", "task_success", "routing", "context", "interaction", "recovery"}

type Catalog struct {
	SchemaVersion string             `json:"schema_version"`
	Thresholds    Thresholds         `json:"thresholds"`
	Tiers         map[string]Tier    `json:"tiers"`
	Metrics       []MetricDefinition `json:"metrics"`
	Cases         []Case             `json:"cases"`
}

type Thresholds struct {
	MaximumSafetyFailures     int     `json:"maximum_safety_failures"`
	MaximumKernelBodyLines    int     `json:"maximum_kernel_body_lines"`
	MinimumTaskSuccessDelta   float64 `json:"minimum_task_success_delta"`
	MinimumTaskSuccessRate    float64 `json:"minimum_task_success_rate"`
	MinimumRoutingPassRate    float64 `json:"minimum_routing_pass_rate"`
	MinimumPointerPassRate    float64 `json:"minimum_pointer_pass_rate"`
	MinimumInitialContextGain float64 `json:"minimum_initial_context_reduction"`
	MinimumTotalContextGain   float64 `json:"minimum_total_context_reduction"`
	MinimumInteractionRate    float64 `json:"minimum_interaction_pass_rate"`
	MinimumRecoveryRate       float64 `json:"minimum_recovery_pass_rate"`
}

type Tier struct {
	Description string   `json:"description"`
	Repetitions int      `json:"repetitions"`
	Include     []string `json:"include"`
}

type MetricDefinition struct {
	Name        string `json:"name"`
	Source      string `json:"source"`
	Aggregation string `json:"aggregation"`
	Failure     string `json:"failure"`
}

type Case struct {
	ID          string             `json:"id"`
	Kind        string             `json:"kind"`
	Suite       string             `json:"suite,omitempty"`
	Tier        string             `json:"tier"`
	HeldOut     bool               `json:"held_out"`
	Fixture     string             `json:"fixture"`
	Prompt      string             `json:"prompt"`
	Intent      string             `json:"intent,omitempty"`
	Complexity  Complexity         `json:"complexity,omitempty"`
	Tags        []string           `json:"tags"`
	Environment map[string]string  `json:"environment,omitempty"`
	Expect      Expectation        `json:"expect"`
	Weights     map[string]float64 `json:"weights"`
}

type Complexity struct {
	Scope       int `json:"scope,omitempty"`
	Ambiguity   int `json:"ambiguity,omitempty"`
	Consequence int `json:"consequence,omitempty"`
	Continuity  int `json:"continuity,omitempty"`
}

type Expectation struct {
	Role                  string      `json:"role,omitempty"`
	ForbiddenRoles        []string    `json:"forbidden_roles,omitempty"`
	Phase                 string      `json:"phase,omitempty"`
	Status                string      `json:"status,omitempty"`
	ForbiddenStatuses     []string    `json:"forbidden_statuses,omitempty"`
	RequiredOutputTerms   []string    `json:"required_output_terms,omitempty"`
	RequiredCommands      []string    `json:"required_commands,omitempty"`
	ForbiddenAnyTerms     []string    `json:"forbidden_any_terms,omitempty"`
	RequiredTraceTerms    []string    `json:"required_trace_terms,omitempty"`
	ForbiddenTraceTerms   []string    `json:"forbidden_trace_terms,omitempty"`
	ExpectedReferences    []string    `json:"expected_references,omitempty"`
	ForbiddenReads        []string    `json:"forbidden_reads,omitempty"`
	AllowedChangedPaths   []string    `json:"allowed_changed_paths,omitempty"`
	ForbiddenChangedPaths []string    `json:"forbidden_changed_paths,omitempty"`
	MaximumOwnerQuestions *int        `json:"maximum_owner_questions,omitempty"`
	OwnerGateRequired     bool        `json:"owner_gate_required,omitempty"`
	ExactlyOnePrimaryRef  bool        `json:"exactly_one_primary_reference,omitempty"`
	RequireSingleReturn   bool        `json:"require_single_return,omitempty"`
	PostChecks            []PostCheck `json:"post_checks,omitempty"`
}

type PostCheck struct {
	Command      []string `json:"command"`
	ExpectedExit int      `json:"expected_exit"`
}

type PostconditionResult struct {
	Command      []string `json:"command"`
	ExpectedExit int      `json:"expected_exit"`
	ActualExit   int      `json:"actual_exit"`
	Output       string   `json:"output,omitempty"`
	MutatedPaths []string `json:"mutated_paths,omitempty"`
	Passed       bool     `json:"passed"`
}

type AgentResult struct {
	Role             string   `json:"role"`
	Phase            string   `json:"phase"`
	Status           string   `json:"status"`
	Summary          string   `json:"summary"`
	NextAction       string   `json:"next_action"`
	OwnerGate        string   `json:"owner_gate"`
	OwnerQuestions   []string `json:"owner_questions"`
	ReferencesLoaded []string `json:"references_loaded"`
	FilesRead        []string `json:"files_read"`
	CommandsRun      []string `json:"commands_run"`
	SafetyNotes      []string `json:"safety_notes"`
}

type DimensionScore struct {
	Applicable int      `json:"applicable"`
	Passed     int      `json:"passed"`
	Rate       *float64 `json:"rate"`
	Findings   []string `json:"findings,omitempty"`
}

type TrialScore struct {
	SafetyPassed bool                      `json:"safety_passed"`
	HardFailures []string                  `json:"hard_failures,omitempty"`
	Dimensions   map[string]DimensionScore `json:"dimensions"`
	Overall      *float64                  `json:"overall"`
	Verdict      string                    `json:"verdict"`
}

type Trial struct {
	ID             string                `json:"id"`
	CaseID         string                `json:"case_id"`
	Suite          string                `json:"suite,omitempty"`
	Complexity     Complexity            `json:"complexity,omitempty"`
	Mode           string                `json:"mode,omitempty"`
	Tags           []string              `json:"tags,omitempty"`
	Variant        string                `json:"variant"`
	Revision       string                `json:"revision"`
	Commit         string                `json:"commit"`
	Model          string                `json:"model"`
	Repeat         int                   `json:"repeat"`
	Order          int                   `json:"order"`
	StartedAt      time.Time             `json:"started_at"`
	DurationMS     int64                 `json:"duration_ms"`
	ExitCode       int                   `json:"exit_code"`
	Result         AgentResult           `json:"result"`
	ChangedPaths   []string              `json:"changed_paths"`
	Postconditions []PostconditionResult `json:"postconditions,omitempty"`
	TraceMetrics   TraceMetrics          `json:"trace_metrics"`
	TracePath      string                `json:"trace_path"`
	ResultPath     string                `json:"result_path"`
	WorkspacePath  string                `json:"workspace_path"`
	Score          TrialScore            `json:"score"`
}

func (complexity Complexity) Total() int {
	return complexity.Scope + complexity.Ambiguity + complexity.Consequence + complexity.Continuity
}

type TraceMetrics struct {
	UsageObserved      bool     `json:"usage_observed"`
	InputTokens        int      `json:"input_tokens"`
	CachedInputTokens  int      `json:"cached_input_tokens"`
	OutputTokens       int      `json:"output_tokens"`
	ToolCalls          int      `json:"tool_calls"`
	Events             int      `json:"events"`
	SemanticEvents     int      `json:"semantic_events"`
	SemanticObserved   bool     `json:"semantic_observed"`
	ObservedFiles      []string `json:"observed_files,omitempty"`
	ObservedReferences []string `json:"observed_references,omitempty"`
	ObservedCommands   []string `json:"observed_commands,omitempty"`
}

type PackageStats struct {
	Label              string         `json:"label"`
	Revision           string         `json:"revision"`
	Commit             string         `json:"commit,omitempty"`
	KernelLines        int            `json:"kernel_lines"`
	KernelBodyLines    int            `json:"kernel_body_lines"`
	KernelWords        int            `json:"kernel_words"`
	KernelBytes        int            `json:"kernel_bytes"`
	ReferenceFiles     int            `json:"reference_files"`
	TotalGuidanceLines int            `json:"total_guidance_lines"`
	TotalGuidanceWords int            `json:"total_guidance_words"`
	PrimaryRouteWords  map[string]int `json:"primary_route_words"`
	ValidationFindings []string       `json:"validation_findings,omitempty"`
}

type StaticComparison struct {
	SchemaVersion string       `json:"schema_version"`
	GeneratedAt   time.Time    `json:"generated_at"`
	Baseline      PackageStats `json:"baseline"`
	Candidate     PackageStats `json:"candidate"`
	Delta         StaticDelta  `json:"delta"`
	Verdict       string       `json:"verdict"`
	GateFailures  []string     `json:"gate_failures,omitempty"`
	Limitations   []string     `json:"limitations"`
}

type StaticDelta struct {
	KernelBodyLineReduction float64            `json:"kernel_body_line_reduction"`
	KernelWordReduction     float64            `json:"kernel_word_reduction"`
	GuidanceWordReduction   float64            `json:"guidance_word_reduction"`
	RouteWordReduction      map[string]float64 `json:"route_word_reduction"`
}

type RunReport struct {
	SchemaVersion      string     `json:"schema_version"`
	GeneratedAt        time.Time  `json:"generated_at"`
	BaselineRef        string     `json:"baseline_ref"`
	BaselineMode       string     `json:"baseline_mode"`
	CandidateRef       string     `json:"candidate_ref"`
	CandidateMode      string     `json:"candidate_mode"`
	Model              string     `json:"model"`
	ReadIsolation      string     `json:"read_isolation"`
	Tier               string     `json:"tier"`
	Seed               int64      `json:"seed"`
	MinimumRepetitions int        `json:"minimum_repetitions"`
	Thresholds         Thresholds `json:"thresholds"`
	Trials             []Trial    `json:"trials"`
	Summary            RunSummary `json:"summary"`
	Limitations        []string   `json:"limitations,omitempty"`
}

type RunSummary struct {
	Verdict              string                        `json:"verdict"`
	MeasurementStatus    string                        `json:"measurement_status,omitempty"`
	ComparativeEffect    string                        `json:"comparative_effect,omitempty"`
	Readiness            string                        `json:"readiness,omitempty"`
	SafetyFailures       map[string]int                `json:"safety_failures"`
	DimensionRates       map[string]map[string]float64 `json:"dimension_rates"`
	ObservedCost         map[string]CostSummary        `json:"observed_cost"`
	Pairing              PairingSummary                `json:"pairing"`
	GateFailures         []string                      `json:"gate_failures,omitempty"`
	CostFindings         []string                      `json:"cost_findings,omitempty"`
	SharedFailures       []string                      `json:"shared_failures,omitempty"`
	PerCaseRegressions   []string                      `json:"per_case_regressions,omitempty"`
	InsufficientEvidence []string                      `json:"insufficient_evidence,omitempty"`
}

type PairingSummary struct {
	Pairs             int      `json:"pairs"`
	BothPass          int      `json:"both_pass"`
	BothFail          int      `json:"both_fail"`
	CandidateWins     int      `json:"candidate_wins"`
	CandidateLosses   int      `json:"candidate_losses"`
	DiscordantRate    float64  `json:"discordant_rate"`
	ExactSignPValue   *float64 `json:"exact_sign_p_value"`
	UnpairedTrialIDs  []string `json:"unpaired_trial_ids,omitempty"`
	UnstableCasePairs []string `json:"unstable_case_pairs,omitempty"`
}

type CostSummary struct {
	TrialsWithUsage      int     `json:"trials_with_usage"`
	TotalTrials          int     `json:"total_trials"`
	MedianInputTokens    float64 `json:"median_input_tokens"`
	MedianCachedTokens   float64 `json:"median_cached_input_tokens"`
	MedianOutputTokens   float64 `json:"median_output_tokens"`
	MedianToolCalls      float64 `json:"median_tool_calls"`
	MedianDurationMillis float64 `json:"median_duration_ms"`
	TotalInputTokens     int     `json:"total_input_tokens,omitempty"`
	TotalCachedTokens    int     `json:"total_cached_tokens,omitempty"`
	TotalOutputTokens    int     `json:"total_output_tokens,omitempty"`
	TotalToolCalls       int     `json:"total_tool_calls,omitempty"`
	TotalDurationMillis  int64   `json:"total_duration_ms,omitempty"`
	SuccessfulTrials     int     `json:"successful_trials,omitempty"`
	TokensPerSuccess     float64 `json:"input_tokens_per_success,omitempty"`
}
