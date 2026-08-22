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
	MinimumTaskSuccessDelta   float64 `json:"minimum_task_success_delta"`
	MinimumRoutingPassRate    float64 `json:"minimum_routing_pass_rate"`
	MinimumPointerPassRate    float64 `json:"minimum_pointer_pass_rate"`
	MinimumInitialContextGain float64 `json:"minimum_initial_context_reduction"`
	MinimumTotalContextGain   float64 `json:"minimum_total_context_reduction"`
}

type Tier struct {
	Description string `json:"description"`
	Repetitions int    `json:"repetitions"`
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
	Tier        string             `json:"tier"`
	HeldOut     bool               `json:"held_out"`
	Fixture     string             `json:"fixture"`
	Prompt      string             `json:"prompt"`
	Tags        []string           `json:"tags"`
	Environment map[string]string  `json:"environment,omitempty"`
	Expect      Expectation        `json:"expect"`
	Weights     map[string]float64 `json:"weights"`
}

type Expectation struct {
	Role                   string   `json:"role,omitempty"`
	Phase                  string   `json:"phase,omitempty"`
	Status                 string   `json:"status,omitempty"`
	RequiredOutputTerms    []string `json:"required_output_terms,omitempty"`
	ForbiddenAnyTerms      []string `json:"forbidden_any_terms,omitempty"`
	RequiredTraceTerms     []string `json:"required_trace_terms,omitempty"`
	ForbiddenTraceTerms    []string `json:"forbidden_trace_terms,omitempty"`
	ExpectedReferences     []string `json:"expected_references,omitempty"`
	ForbiddenReads         []string `json:"forbidden_reads,omitempty"`
	AllowedChangedPaths    []string `json:"allowed_changed_paths,omitempty"`
	ForbiddenChangedPaths  []string `json:"forbidden_changed_paths,omitempty"`
	MaximumOwnerQuestions  *int     `json:"maximum_owner_questions,omitempty"`
	ExactlyOnePrimaryRef   bool     `json:"exactly_one_primary_reference,omitempty"`
	RequireSingleReturn    bool     `json:"require_single_return,omitempty"`
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
	ID             string      `json:"id"`
	CaseID         string      `json:"case_id"`
	Variant        string      `json:"variant"`
	Revision       string      `json:"revision"`
	Commit         string      `json:"commit"`
	Model          string      `json:"model"`
	Repeat         int         `json:"repeat"`
	Order          int         `json:"order"`
	StartedAt      time.Time   `json:"started_at"`
	DurationMS     int64       `json:"duration_ms"`
	ExitCode       int         `json:"exit_code"`
	Result         AgentResult `json:"result"`
	ChangedPaths   []string    `json:"changed_paths"`
	TracePath      string      `json:"trace_path"`
	ResultPath     string      `json:"result_path"`
	WorkspacePath  string      `json:"workspace_path"`
	Score          TrialScore  `json:"score"`
}

type PackageStats struct {
	Label                 string         `json:"label"`
	Revision              string         `json:"revision"`
	Commit                string         `json:"commit,omitempty"`
	KernelLines           int            `json:"kernel_lines"`
	KernelBodyLines       int            `json:"kernel_body_lines"`
	KernelWords           int            `json:"kernel_words"`
	KernelBytes           int            `json:"kernel_bytes"`
	ReferenceFiles        int            `json:"reference_files"`
	TotalGuidanceLines    int            `json:"total_guidance_lines"`
	TotalGuidanceWords    int            `json:"total_guidance_words"`
	PrimaryRouteWords     map[string]int `json:"primary_route_words"`
	ValidationFindings    []string       `json:"validation_findings,omitempty"`
}

type StaticComparison struct {
	SchemaVersion string       `json:"schema_version"`
	GeneratedAt   time.Time    `json:"generated_at"`
	Baseline      PackageStats `json:"baseline"`
	Candidate     PackageStats `json:"candidate"`
	Delta         StaticDelta  `json:"delta"`
	Verdict       string       `json:"verdict"`
	Limitations   []string     `json:"limitations"`
}

type StaticDelta struct {
	KernelBodyLineReduction float64            `json:"kernel_body_line_reduction"`
	KernelWordReduction     float64            `json:"kernel_word_reduction"`
	GuidanceWordReduction   float64            `json:"guidance_word_reduction"`
	RouteWordReduction      map[string]float64 `json:"route_word_reduction"`
}

type RunReport struct {
	SchemaVersion string      `json:"schema_version"`
	GeneratedAt   time.Time   `json:"generated_at"`
	BaselineRef   string      `json:"baseline_ref"`
	CandidateRef  string      `json:"candidate_ref"`
	Model         string      `json:"model"`
	Tier          string      `json:"tier"`
	Seed          int64       `json:"seed"`
	Trials        []Trial     `json:"trials"`
	Summary       RunSummary  `json:"summary"`
	Limitations   []string    `json:"limitations,omitempty"`
}

type RunSummary struct {
	Verdict              string                        `json:"verdict"`
	SafetyFailures       map[string]int                `json:"safety_failures"`
	DimensionRates       map[string]map[string]float64 `json:"dimension_rates"`
	PerCaseRegressions   []string                      `json:"per_case_regressions,omitempty"`
	InsufficientEvidence []string                      `json:"insufficient_evidence,omitempty"`
}
