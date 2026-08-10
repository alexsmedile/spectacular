package governance

// Mutation inputs are explicit, already-confirmed mechanical requests. The
// guided Skill owns authoring and judgment before these structures reach the
// CLI.
type Modification struct {
	From string `json:"from"`
	To   string `json:"to"`
}

type ContractCandidate struct {
	Purpose               string   `json:"purpose"`
	Outcome               string   `json:"outcome"`
	AppliesWhen           []string `json:"applies_when"`
	DoesNotApplyWhen      []string `json:"does_not_apply_when"`
	DoesNotProvide        []string `json:"does_not_provide"`
	RequiredBehavior      []string `json:"required_behavior"`
	OperatingCases        []string `json:"operating_cases"`
	PersistentInformation []string `json:"persistent_information"`
	ConformanceChecks     []string `json:"conformance_checks"`
	AuthorityFreshness    string   `json:"authority_freshness"`
	RelatedMaterial       []string `json:"related_material"`
}

type ProposalInput struct {
	ID                  string            `json:"id"`
	Title               string            `json:"title"`
	Actor               string            `json:"actor"`
	Status              string            `json:"status"`
	TargetContract      string            `json:"target_contract"`
	BaseVersion         string            `json:"base_version"`
	BaseFingerprint     string            `json:"base_fingerprint"`
	NewCapability       bool              `json:"new_capability"`
	Additions           []string          `json:"additions"`
	Modifications       []Modification    `json:"modifications"`
	Removals            []string          `json:"removals"`
	Candidate           ContractCandidate `json:"candidate"`
	Rationale           string            `json:"rationale"`
	Scope               []string          `json:"scope"`
	Gaps                []string          `json:"gaps"`
	Authorization       string            `json:"authorization"`
	IdempotencyKey      string            `json:"idempotency_key"`
	FreshnessValidUntil string            `json:"freshness_valid_until"`
}

type ObjectiveInput struct {
	ID            string   `json:"id"`
	Outcome       string   `json:"outcome"`
	Dependencies  []string `json:"dependencies"`
	ExpectedProof []string `json:"expected_proof"`
}

type MissionInput struct {
	ID                          string           `json:"id"`
	Title                       string           `json:"title"`
	Actor                       string           `json:"actor"`
	Proposal                    string           `json:"proposal"`
	Outcome                     string           `json:"outcome"`
	Objectives                  []ObjectiveInput `json:"objectives"`
	InitialRunID                string           `json:"initial_run_id"`
	DesignSufficiency           string           `json:"design_sufficiency"`
	SliceQuality                string           `json:"slice_quality"`
	Dependencies                []string         `json:"dependencies"`
	Gaps                        []string         `json:"gaps"`
	EvidenceClaims              []string         `json:"evidence_claims"`
	Scope                       []string         `json:"scope"`
	AllowedActions              []string         `json:"allowed_actions"`
	ForbiddenEffects            []string         `json:"forbidden_effects"`
	Baseline                    string           `json:"baseline"`
	BudgetUnits                 int              `json:"budget_units"`
	RepairBudget                int              `json:"repair_budget"`
	ExpiresAt                   string           `json:"expires_at"`
	Stops                       []string         `json:"stops"`
	RecoveryPoint               string           `json:"recovery_point"`
	ReturnDestination           string           `json:"return_destination"`
	Authorization               string           `json:"authorization"`
	ExpectedProposalFingerprint string           `json:"expected_proposal_fingerprint"`
	IdempotencyKey              string           `json:"idempotency_key"`
}

type HandoffInput struct {
	ID                         string   `json:"id"`
	Title                      string   `json:"title"`
	Mission                    string   `json:"mission"`
	Objective                  string   `json:"objective"`
	Run                        string   `json:"run"`
	Sender                     string   `json:"sender"`
	Actor                      string   `json:"actor"`
	Destination                string   `json:"destination"`
	HostPointer                string   `json:"host_pointer"`
	Scope                      []string `json:"scope"`
	Inputs                     []string `json:"inputs"`
	AllowedActions             []string `json:"allowed_actions"`
	ForbiddenEffects           []string `json:"forbidden_effects"`
	EvidenceClaims             []string `json:"evidence_claims"`
	BudgetUnits                int      `json:"budget_units"`
	ExpiresAt                  string   `json:"expires_at"`
	Stops                      []string `json:"stops"`
	RecoveryPoint              string   `json:"recovery_point"`
	ReturnDestination          string   `json:"return_destination"`
	Authorization              string   `json:"authorization"`
	ExpectedMissionFingerprint string   `json:"expected_mission_fingerprint"`
	IdempotencyKey             string   `json:"idempotency_key"`
	Supersedes                 string   `json:"supersedes"`
}

type HandoffReturnInput struct {
	ID                          string   `json:"id"`
	Title                       string   `json:"title"`
	Dispatch                    string   `json:"dispatch"`
	Status                      string   `json:"status"`
	Actor                       string   `json:"actor"`
	FinalBaseline               string   `json:"final_baseline"`
	Result                      string   `json:"result"`
	Actions                     []string `json:"actions"`
	ProviderReceipts            []string `json:"provider_receipts"`
	Evidence                    []string `json:"evidence"`
	RemainingGaps               []string `json:"remaining_gaps"`
	BudgetUsed                  int      `json:"budget_used"`
	RecoveryPoint               string   `json:"recovery_point"`
	NextAction                  string   `json:"next_action"`
	OwnerGate                   string   `json:"owner_gate"`
	ExpectedDispatchFingerprint string   `json:"expected_dispatch_fingerprint"`
	IdempotencyKey              string   `json:"idempotency_key"`
}

type EvidenceInput struct {
	ID                  string   `json:"id"`
	Title               string   `json:"title"`
	Mission             string   `json:"mission"`
	Objective           string   `json:"objective"`
	Checkpoint          string   `json:"checkpoint"`
	Claim               string   `json:"claim"`
	Classification      string   `json:"classification"`
	Scope               []string `json:"scope"`
	Method              string   `json:"method"`
	Actor               string   `json:"actor"`
	Target              string   `json:"target"`
	Environment         string   `json:"environment"`
	ObservedAt          string   `json:"observed_at"`
	FreshnessValidUntil string   `json:"freshness_valid_until"`
	Limitations         []string `json:"limitations"`
	ContraryEvidence    []string `json:"contrary_evidence"`
	RequiredChecks      []string `json:"required_checks"`
	CheckResults        []string `json:"check_results"`
	ReviewState         string   `json:"review_state"`
	ExecutorAuthored    bool     `json:"executor_authored"`
	Authorization       string   `json:"authorization"`
	IdempotencyKey      string   `json:"idempotency_key"`
}

type DecisionInput struct {
	ID                   string   `json:"id"`
	Title                string   `json:"title"`
	Actor                string   `json:"actor"`
	ActorRole            string   `json:"actor_role"`
	AuthorityBasis       string   `json:"authority_basis"`
	Question             string   `json:"question"`
	Scope                []string `json:"scope"`
	Disposition          string   `json:"disposition"`
	Rationale            string   `json:"rationale"`
	Alternatives         []string `json:"alternatives"`
	Targets              []string `json:"targets"`
	ExpectedFingerprints []string `json:"expected_fingerprints"`
	Operation            string   `json:"operation"`
	AuthorizedEffects    []string `json:"authorized_effects"`
	Conditions           []string `json:"conditions"`
	ExpiresAt            string   `json:"expires_at"`
	Evidence             []string `json:"evidence"`
	Supersedes           string   `json:"supersedes"`
	IdempotencyKey       string   `json:"idempotency_key"`
}

type RepairAttempt struct {
	AffectedClaims     []string `json:"affected_claims"`
	PreviousHypothesis string   `json:"previous_hypothesis"`
	NewHypothesis      string   `json:"new_hypothesis"`
	NewEvidence        []string `json:"new_evidence"`
	NarrowerAction     string   `json:"narrower_action"`
	Actor              string   `json:"actor"`
	BeforeEvidence     []string `json:"before_evidence"`
	AfterEvidence      []string `json:"after_evidence"`
	Checks             []string `json:"checks"`
	Result             string   `json:"result"`
	BudgetConsumed     int      `json:"budget_consumed"`
	RecoveryPoint      string   `json:"recovery_point"`
}

type AssessmentInput struct {
	ID               string          `json:"id"`
	Title            string          `json:"title"`
	Mission          string          `json:"mission"`
	Verdict          string          `json:"verdict"`
	Actor            string          `json:"actor"`
	Claims           []string        `json:"claims"`
	Evidence         []string        `json:"evidence"`
	BlockingFindings []string        `json:"blocking_findings"`
	Limitations      []string        `json:"limitations"`
	RepairAttempts   []RepairAttempt `json:"repair_attempts"`
	RecoveryPoint    string          `json:"recovery_point"`
	Authorization    string          `json:"authorization"`
	IdempotencyKey   string          `json:"idempotency_key"`
}

type TransitionInput struct {
	Mission             string
	To                  string
	Authorization       string
	ExpectedFingerprint string
	IdempotencyKey      string
	Disposition         string
	Assessment          string
	Reconciliation      string
	TerminalNextAction  string
	SatisfiedObjectives []string
}

type ReconcileInput struct {
	Contract            string `json:"contract"`
	Proposal            string `json:"proposal"`
	Authorization       string `json:"authorization"`
	ExpectedFingerprint string `json:"expected_fingerprint"`
	IdempotencyKey      string `json:"idempotency_key"`
}

type ReconcileSetInput struct {
	Items []ReconcileInput `json:"items"`
}

type ArchiveInput struct {
	Mission             string
	Authorization       string
	ExpectedFingerprint string
	IdempotencyKey      string
	TerminalPacket      string
}

type OperationResult struct {
	Operation        string   `json:"operation"`
	Ref              string   `json:"ref"`
	Path             string   `json:"path"`
	Fingerprint      string   `json:"fingerprint"`
	IdempotentReplay bool     `json:"idempotent_replay"`
	Sources          []string `json:"sources"`
	Receipt          string   `json:"receipt,omitempty"`
}
