// Package runtime compiles preparation and Autopilot inputs for a host coding
// runtime. It validates explicit authority; it never starts a runtime or
// performs a provider effect.
package runtime

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"sort"
	"strings"
	"time"

	"github.com/alexsmedile/spectacular/v2/internal/domain"
)

const (
	PreparationSchema = "spectacular.mission-preparation.v2"
	AutopilotSchema   = "spectacular.autopilot-charter.v1"
)

const (
	ReviewAutomatic   = "automatic"
	ReviewClustered   = "clustered"
	ReviewIndependent = "independent"
)

const (
	EnforcementHard        = "hard"
	EnforcementObserved    = "observed"
	EnforcementUnsupported = "unsupported"
)

type BoundSource struct {
	Ref         string `json:"ref"`
	Fingerprint string `json:"fingerprint"`
}

type CandidateSlice struct {
	Name                string   `json:"name"`
	Outcome             string   `json:"outcome"`
	Evidence            []string `json:"evidence"`
	Dependencies        []string `json:"dependencies"`
	CancellationState   string   `json:"cancellation_state"`
	Reversibility       string   `json:"reversibility"`
	StandaloneCoherence string   `json:"standalone_coherence"`
	IntegrationPath     string   `json:"integration_path"`
	LearningValue       string   `json:"learning_value"`
}

// CompletionCriterion is the frozen, minimal answer key for one Mission claim.
// Changing any field changes the preparation fingerprint and therefore needs a
// new owner-authorized Mission boundary.
type CompletionCriterion struct {
	Claim            string `json:"claim"`
	PassBoundary     string `json:"pass_boundary"`
	ProofRequirement string `json:"proof_requirement"`
	ReviewLevel      string `json:"review_level"`
}

// ReadinessIssue identifies exactly what preparation still needs to resolve.
// Guided workflows use these issues to grill adaptively instead of imposing a
// mandatory interview on already-sufficient work.
type ReadinessIssue struct {
	Code   string `json:"code"`
	Claim  string `json:"claim,omitempty"`
	Detail string `json:"detail"`
}

type PreparationInput struct {
	Proposal           BoundSource           `json:"proposal"`
	Baseline           string                `json:"baseline"`
	DirectionSources   []BoundSource         `json:"direction_sources"`
	Candidates         []CandidateSlice      `json:"candidates"`
	Selected           string                `json:"selected"`
	DesignSufficiency  string                `json:"design_sufficiency"`
	DesignRationale    string                `json:"design_rationale"`
	SliceQuality       string                `json:"slice_quality"`
	SliceRationale     string                `json:"slice_rationale"`
	BlockingGaps       []string              `json:"blocking_gaps"`
	CompletionCriteria []CompletionCriterion `json:"completion_criteria"`
	StopConditions     []string              `json:"stop_conditions"`
	EvidenceClaims     []string              `json:"evidence_claims"`
	FreshUntil         string                `json:"fresh_until"`
}

type PreparationReceipt struct {
	SchemaVersion      string                `json:"schema_version"`
	Fingerprint        string                `json:"fingerprint"`
	Proposal           BoundSource           `json:"proposal"`
	Baseline           string                `json:"baseline"`
	DirectionSources   []BoundSource         `json:"direction_sources"`
	Candidates         []CandidateSlice      `json:"candidates"`
	Selected           string                `json:"selected"`
	DesignSufficiency  string                `json:"design_sufficiency"`
	DesignRationale    string                `json:"design_rationale"`
	SliceQuality       string                `json:"slice_quality"`
	SliceRationale     string                `json:"slice_rationale"`
	BlockingGaps       []string              `json:"blocking_gaps"`
	CompletionCriteria []CompletionCriterion `json:"completion_criteria"`
	StopConditions     []string              `json:"stop_conditions"`
	EvidenceClaims     []string              `json:"evidence_claims"`
	FreshUntil         string                `json:"fresh_until"`
	Ready              bool                  `json:"ready"`
	UnmetRequirements  []ReadinessIssue      `json:"unmet_requirements"`
	Next               string                `json:"next"`
}

func CompilePreparation(input PreparationInput, now time.Time) (PreparationReceipt, error) {
	if err := validateBoundSource(input.Proposal, domain.Proposal); err != nil {
		return PreparationReceipt{}, err
	}
	if input.Baseline == "" || len(input.DirectionSources) == 0 || len(input.Candidates) < 1 || len(input.Candidates) > 3 || input.Selected == "" || input.DesignRationale == "" || input.SliceRationale == "" || len(input.StopConditions) == 0 || len(input.EvidenceClaims) == 0 {
		return PreparationReceipt{}, invalid("preparation requires baseline, sources, one to three slices, rationales, completion, stops, and evidence")
	}
	for _, source := range input.DirectionSources {
		if err := validateAnyBoundSource(source); err != nil {
			return PreparationReceipt{}, err
		}
	}
	selected := false
	seen := map[string]bool{}
	for _, candidate := range input.Candidates {
		if candidate.Name == "" || seen[candidate.Name] || candidate.Outcome == "" || len(candidate.Evidence) == 0 || candidate.CancellationState == "" || candidate.Reversibility == "" || candidate.StandaloneCoherence == "" || candidate.IntegrationPath == "" || candidate.LearningValue == "" {
			return PreparationReceipt{}, invalid("each candidate slice needs a unique name, outcome, proof, cancellation state, reversibility, coherence, integration, and learning value")
		}
		seen[candidate.Name] = true
		selected = selected || candidate.Name == input.Selected
	}
	if !selected {
		return PreparationReceipt{}, invalid("selected slice does not exist")
	}
	if !oneOf(input.DesignSufficiency, "sufficient", "needs-evidence", "needs-decision") {
		return PreparationReceipt{}, invalid("design_sufficiency must be sufficient, needs-evidence, or needs-decision")
	}
	if !oneOf(input.SliceQuality, "coherent", "too-broad", "fragmented", "dependency-bound") {
		return PreparationReceipt{}, invalid("slice_quality must be coherent, too-broad, fragmented, or dependency-bound")
	}
	freshUntil, err := time.Parse(time.RFC3339, input.FreshUntil)
	if err != nil || !freshUntil.After(now) {
		return PreparationReceipt{}, invalid("fresh_until must be a future RFC3339 time")
	}
	criteria := append([]CompletionCriterion(nil), input.CompletionCriteria...)
	sort.Slice(criteria, func(i, j int) bool { return criteria[i].Claim < criteria[j].Claim })
	issues := readinessIssues(input, criteria)
	ready := len(issues) == 0
	next := "owner-activation"
	if !ready {
		next = "adaptive-grill"
	}
	receipt := PreparationReceipt{
		SchemaVersion: PreparationSchema, Proposal: input.Proposal, Baseline: input.Baseline,
		DirectionSources: input.DirectionSources, Candidates: input.Candidates, Selected: input.Selected,
		DesignSufficiency: input.DesignSufficiency, DesignRationale: input.DesignRationale,
		SliceQuality: input.SliceQuality, SliceRationale: input.SliceRationale, BlockingGaps: input.BlockingGaps,
		CompletionCriteria: criteria, StopConditions: input.StopConditions,
		EvidenceClaims: input.EvidenceClaims, FreshUntil: freshUntil.UTC().Format(time.RFC3339), Ready: ready,
		UnmetRequirements: issues, Next: next,
	}
	receipt.Fingerprint = fingerprint(receipt)
	return receipt, nil
}

func ValidatePreparationReceipt(receipt PreparationReceipt, now time.Time) error {
	if receipt.SchemaVersion != PreparationSchema || receipt.Fingerprint == "" {
		return invalid("preparation receipt schema and fingerprint are required")
	}
	stored := receipt.Fingerprint
	receipt.Fingerprint = ""
	if fingerprint(receipt) != stored {
		return invalid("preparation receipt fingerprint does not match its content")
	}
	if !receipt.Ready || receipt.DesignSufficiency != "sufficient" || receipt.SliceQuality != "coherent" || len(receipt.BlockingGaps) != 0 || len(receipt.UnmetRequirements) != 0 || receipt.Next != "owner-activation" {
		return invalid("preparation receipt does not permit owner activation")
	}
	if issues := readinessIssues(PreparationInput{DesignSufficiency: receipt.DesignSufficiency, SliceQuality: receipt.SliceQuality, BlockingGaps: receipt.BlockingGaps, EvidenceClaims: receipt.EvidenceClaims}, receipt.CompletionCriteria); len(issues) != 0 {
		return invalid("preparation receipt completion criteria are incomplete")
	}
	freshUntil, err := time.Parse(time.RFC3339, receipt.FreshUntil)
	if err != nil || !freshUntil.After(now) {
		return domain.NewRefusal(domain.RefusalExpiredAuthority, "preparation", "preparation receipt is stale", err)
	}
	return nil
}

func readinessIssues(input PreparationInput, criteria []CompletionCriterion) []ReadinessIssue {
	issues := []ReadinessIssue{}
	if input.DesignSufficiency == "needs-evidence" {
		issues = append(issues, ReadinessIssue{Code: "needs-evidence", Detail: input.DesignRationale})
	}
	if input.DesignSufficiency == "needs-decision" {
		issues = append(issues, ReadinessIssue{Code: "needs-decision", Detail: input.DesignRationale})
	}
	if input.SliceQuality != "coherent" {
		issues = append(issues, ReadinessIssue{Code: "slice-not-coherent", Detail: input.SliceRationale})
	}
	for _, gap := range input.BlockingGaps {
		issues = append(issues, ReadinessIssue{Code: "blocking-gap", Detail: gap})
	}
	declared := map[string]bool{}
	for _, claim := range input.EvidenceClaims {
		if declared[claim] {
			issues = append(issues, ReadinessIssue{Code: "evidence-claim-duplicate", Claim: claim, Detail: "evidence claim is declared more than once"})
		}
		declared[claim] = true
	}
	seen := map[string]bool{}
	for _, criterion := range criteria {
		switch {
		case criterion.Claim == "":
			issues = append(issues, ReadinessIssue{Code: "criterion-missing-claim", Detail: "completion criterion has no claim"})
		case seen[criterion.Claim]:
			issues = append(issues, ReadinessIssue{Code: "criterion-duplicate-claim", Claim: criterion.Claim, Detail: "claim has more than one completion criterion"})
		case !declared[criterion.Claim]:
			issues = append(issues, ReadinessIssue{Code: "criterion-undeclared-claim", Claim: criterion.Claim, Detail: "criterion claim is not in evidence_claims"})
		}
		seen[criterion.Claim] = true
		if criterion.PassBoundary == "" {
			issues = append(issues, ReadinessIssue{Code: "criterion-missing-pass-boundary", Claim: criterion.Claim, Detail: "pass boundary is required"})
		}
		if criterion.ProofRequirement == "" {
			issues = append(issues, ReadinessIssue{Code: "criterion-missing-proof", Claim: criterion.Claim, Detail: "proof requirement is required"})
		}
		if !oneOf(criterion.ReviewLevel, ReviewAutomatic, ReviewClustered, ReviewIndependent) {
			issues = append(issues, ReadinessIssue{Code: "criterion-invalid-review-level", Claim: criterion.Claim, Detail: "review level must be automatic, clustered, or independent"})
		}
	}
	for _, claim := range input.EvidenceClaims {
		if !seen[claim] {
			issues = append(issues, ReadinessIssue{Code: "criterion-missing", Claim: claim, Detail: "claim needs a completion criterion"})
		}
	}
	sort.SliceStable(issues, func(i, j int) bool {
		left, right := issues[i].Code+"\x00"+issues[i].Claim+"\x00"+issues[i].Detail, issues[j].Code+"\x00"+issues[j].Claim+"\x00"+issues[j].Detail
		return left < right
	})
	return issues
}

var requiredForbidden = []string{
	"merge", "production-release", "production-configuration", "secret-change",
	"remote-deletion", "destructive-data", "security-privacy-rights-sensitive",
}

type AutopilotInput struct {
	Enabled                 bool            `json:"enabled"`
	Mission                 BoundSource     `json:"mission"`
	Authorization           BoundSource     `json:"authorization"`
	Outcome                 string          `json:"outcome"`
	NonGoals                []string        `json:"non_goals"`
	AuthoritativeSources    []BoundSource   `json:"authoritative_sources"`
	DelegatedDecisionDomain []string        `json:"delegated_decision_domain"`
	AllowedProviders        []string        `json:"allowed_providers"`
	AllowedActions          []string        `json:"allowed_actions"`
	ForbiddenEffects        []string        `json:"forbidden_effects"`
	BudgetUnits             int             `json:"budget_units"`
	RepairBudget            int             `json:"repair_budget"`
	ResourceLimits          []ResourceLimit `json:"resource_limits"`
	RequiredChecks          []string        `json:"required_checks"`
	ExpiresAt               string          `json:"expires_at"`
	StopConditions          []string        `json:"stop_conditions"`
	RecoveryPoint           string          `json:"recovery_point"`
	ReturnDestination       string          `json:"return_destination"`
}

// ResourceLimit states both the requested cap and what the selected host can
// truthfully do about it. Hard means measure and cancel; observed means measure
// only; unsupported promises neither.
type ResourceLimit struct {
	Resource          string `json:"resource"`
	Maximum           int    `json:"maximum"`
	Enforcement       string `json:"enforcement"`
	MeasureCapability string `json:"measure_capability,omitempty"`
	CancelCapability  string `json:"cancel_capability,omitempty"`
}

// AutopilotLimits is derived from the exact bound Mission and owner Decision.
// It is deliberately separate from AutopilotInput so requested authority can
// never validate itself.
type AutopilotLimits struct {
	Mission          BoundSource
	Authorization    BoundSource
	Outcome          string
	AllowedProviders []string
	AllowedActions   []string
	ForbiddenEffects []string
	BudgetUnits      int
	RepairBudget     int
	HostCapabilities []string
	ExpiresAt        string
	StopConditions   []string
}

type HandoffContract struct {
	Accountability   string   `json:"accountability"`
	Statuses         []string `json:"statuses"`
	RequiredReturn   []string `json:"required_return"`
	HostPointerRole  string   `json:"host_pointer_role"`
	ProviderBoundary string   `json:"provider_boundary"`
}

type AutopilotCharter struct {
	SchemaVersion            string          `json:"schema_version"`
	Fingerprint              string          `json:"fingerprint"`
	Enabled                  bool            `json:"enabled"`
	Mission                  BoundSource     `json:"mission"`
	Authorization            BoundSource     `json:"authorization"`
	Outcome                  string          `json:"outcome"`
	NonGoals                 []string        `json:"non_goals"`
	AuthoritativeSources     []BoundSource   `json:"authoritative_sources"`
	DelegatedDecisionDomain  []string        `json:"delegated_decision_domain"`
	AllowedProviders         []string        `json:"allowed_providers"`
	AllowedActions           []string        `json:"allowed_actions"`
	ForbiddenEffects         []string        `json:"forbidden_effects"`
	BudgetUnits              int             `json:"budget_units"`
	RepairBudget             int             `json:"repair_budget"`
	ResourceLimits           []ResourceLimit `json:"resource_limits"`
	HostCapabilities         []string        `json:"host_capabilities"`
	RequiredChecks           []string        `json:"required_checks"`
	ExpiresAt                string          `json:"expires_at"`
	StopConditions           []string        `json:"stop_conditions"`
	RecoveryPoint            string          `json:"recovery_point"`
	ReturnDestination        string          `json:"return_destination"`
	RuntimeNeutral           bool            `json:"runtime_neutral"`
	AutomaticProviderEffects bool            `json:"automatic_provider_effects"`
	Handoff                  HandoffContract `json:"handoff"`
}

func CompileAutopilot(input AutopilotInput, limits AutopilotLimits, now time.Time) (AutopilotCharter, error) {
	if !input.Enabled {
		return AutopilotCharter{}, invalid("Autopilot must be explicitly enabled")
	}
	if err := validateBoundSource(input.Mission, domain.Mission); err != nil {
		return AutopilotCharter{}, err
	}
	if err := validateBoundSource(input.Authorization, domain.Decision); err != nil {
		return AutopilotCharter{}, err
	}
	if input.Outcome == "" || len(input.NonGoals) == 0 || len(input.AuthoritativeSources) == 0 || len(input.DelegatedDecisionDomain) == 0 || len(input.AllowedActions) == 0 || input.BudgetUnits < 1 || input.RepairBudget < 0 || len(input.RequiredChecks) == 0 || len(input.StopConditions) == 0 || input.RecoveryPoint == "" || input.ReturnDestination == "" {
		return AutopilotCharter{}, invalid("Autopilot requires outcome, non-goals, sources, decision domain, actions, budgets, checks, stops, recovery, and return")
	}
	if limits.Mission != input.Mission || limits.Authorization != input.Authorization || limits.Outcome == "" || limits.Outcome != input.Outcome {
		return AutopilotCharter{}, invalid("Autopilot request is not bound to the validated Mission outcome and authorization")
	}
	if !allContained(input.AllowedProviders, limits.AllowedProviders) || !allContained(input.AllowedActions, limits.AllowedActions) {
		return AutopilotCharter{}, invalid("Autopilot providers and actions exceed validated authority")
	}
	if !allContained(limits.ForbiddenEffects, input.ForbiddenEffects) || overlaps(input.AllowedActions, input.ForbiddenEffects) {
		return AutopilotCharter{}, invalid("Autopilot allowed actions conflict with, or omit, Mission forbidden effects")
	}
	if input.BudgetUnits > limits.BudgetUnits || input.RepairBudget > limits.RepairBudget {
		return AutopilotCharter{}, invalid("Autopilot budgets exceed the Mission envelope")
	}
	resourceLimits, hostCapabilities, err := validateResourceLimits(input.ResourceLimits, limits.HostCapabilities, input.RepairBudget)
	if err != nil {
		return AutopilotCharter{}, err
	}
	if !allContained(limits.StopConditions, input.StopConditions) {
		return AutopilotCharter{}, invalid("Autopilot stop conditions omit a Mission stop")
	}
	for _, source := range input.AuthoritativeSources {
		if err := validateAnyBoundSource(source); err != nil {
			return AutopilotCharter{}, err
		}
	}
	for _, forbidden := range requiredForbidden {
		if !contains(input.ForbiddenEffects, forbidden) {
			return AutopilotCharter{}, invalid("forbidden_effects must include " + forbidden)
		}
	}
	expires, err := time.Parse(time.RFC3339, input.ExpiresAt)
	if err != nil || !expires.After(now) {
		return AutopilotCharter{}, invalid("expires_at must be a future RFC3339 time")
	}
	limitExpiry, err := time.Parse(time.RFC3339, limits.ExpiresAt)
	if err != nil || expires.After(limitExpiry) {
		return AutopilotCharter{}, invalid("Autopilot expiry exceeds validated authority")
	}
	charter := AutopilotCharter{
		SchemaVersion: AutopilotSchema, Enabled: true, Mission: input.Mission,
		Authorization: input.Authorization, Outcome: input.Outcome, NonGoals: input.NonGoals,
		AuthoritativeSources: input.AuthoritativeSources, DelegatedDecisionDomain: input.DelegatedDecisionDomain,
		AllowedProviders: input.AllowedProviders, AllowedActions: input.AllowedActions,
		ForbiddenEffects: input.ForbiddenEffects, BudgetUnits: input.BudgetUnits, RepairBudget: input.RepairBudget,
		ResourceLimits: resourceLimits, HostCapabilities: hostCapabilities,
		RequiredChecks: input.RequiredChecks, ExpiresAt: expires.UTC().Format(time.RFC3339),
		StopConditions: input.StopConditions, RecoveryPoint: input.RecoveryPoint,
		ReturnDestination: input.ReturnDestination, RuntimeNeutral: true, AutomaticProviderEffects: false,
		Handoff: HandoffContract{
			Accountability:   "Mission accountability remains with the sender",
			Statuses:         []string{"succeeded", "blocked", "failed"},
			RequiredReturn:   []string{"actor", "final baseline", "result", "actions", "provider receipts", "Evidence", "remaining Gaps", "budget use", "recovery point", "exactly one next action or owner gate"},
			HostPointerRole:  "non-authoritative destination pointer",
			ProviderBoundary: "native providers retain permission and effect authority",
		},
	}
	charter.Fingerprint = fingerprint(charter)
	return charter, nil
}

func ValidateAutopilotCharter(charter AutopilotCharter, now time.Time) error {
	if charter.SchemaVersion != AutopilotSchema || charter.Fingerprint == "" || !charter.Enabled || !charter.RuntimeNeutral || charter.AutomaticProviderEffects {
		return invalid("Autopilot charter schema or safety flags are invalid")
	}
	stored := charter.Fingerprint
	charter.Fingerprint = ""
	if fingerprint(charter) != stored {
		return invalid("Autopilot charter fingerprint does not match its content")
	}
	expires, err := time.Parse(time.RFC3339, charter.ExpiresAt)
	if err != nil || !expires.After(now) {
		return invalid("Autopilot charter is stale")
	}
	for _, forbidden := range requiredForbidden {
		if !contains(charter.ForbiddenEffects, forbidden) {
			return invalid("Autopilot charter is below the initial effect ceiling")
		}
	}
	if _, _, err := validateResourceLimits(charter.ResourceLimits, charter.HostCapabilities, charter.RepairBudget); err != nil {
		return err
	}
	return nil
}

func validateResourceLimits(limits []ResourceLimit, capabilities []string, repairBudget int) ([]ResourceLimit, []string, error) {
	canonicalLimits := append([]ResourceLimit(nil), limits...)
	sort.Slice(canonicalLimits, func(i, j int) bool { return canonicalLimits[i].Resource < canonicalLimits[j].Resource })
	canonicalCapabilities := sortedUniqueStrings(capabilities)
	seen := map[string]bool{}
	for _, limit := range canonicalLimits {
		if !oneOf(limit.Resource, "wall-time-seconds", "tokens", "spend-cents", "parallel-workers", "repair-rounds") {
			return nil, nil, invalid("resource limit must target wall-time-seconds, tokens, spend-cents, parallel-workers, or repair-rounds")
		}
		if seen[limit.Resource] {
			return nil, nil, invalid("resource limits may name each resource only once")
		}
		seen[limit.Resource] = true
		if limit.Maximum < 1 {
			return nil, nil, invalid("resource limit maximum must be positive")
		}
		if limit.Resource == "repair-rounds" && limit.Maximum > repairBudget {
			return nil, nil, invalid("repair-rounds limit exceeds the Autopilot repair budget")
		}
		switch limit.Enforcement {
		case EnforcementHard:
			if limit.MeasureCapability == "" || limit.CancelCapability == "" || !contains(canonicalCapabilities, limit.MeasureCapability) || !contains(canonicalCapabilities, limit.CancelCapability) {
				return nil, nil, invalid("hard resource limit requires declared host measurement and cancellation capabilities")
			}
		case EnforcementObserved:
			if limit.MeasureCapability == "" || !contains(canonicalCapabilities, limit.MeasureCapability) || limit.CancelCapability != "" {
				return nil, nil, invalid("observed resource limit requires declared measurement capability and no cancellation claim")
			}
		case EnforcementUnsupported:
			if limit.MeasureCapability != "" || limit.CancelCapability != "" {
				return nil, nil, invalid("unsupported resource limit cannot claim measurement or cancellation capabilities")
			}
		default:
			return nil, nil, invalid("resource limit enforcement must be hard, observed, or unsupported")
		}
	}
	return canonicalLimits, canonicalCapabilities, nil
}

func sortedUniqueStrings(values []string) []string {
	result := append([]string(nil), values...)
	sort.Strings(result)
	out := result[:0]
	for _, value := range result {
		if value != "" && (len(out) == 0 || out[len(out)-1] != value) {
			out = append(out, value)
		}
	}
	return out
}

func validateBoundSource(source BoundSource, want domain.RecordType) error {
	ref, err := domain.ParseReference(source.Ref)
	if err != nil || ref.Type != want {
		return invalid(fmt.Sprintf("source %q must be an exact %s reference", source.Ref, want))
	}
	return validateFingerprint(source.Fingerprint)
}

func validateAnyBoundSource(source BoundSource) error {
	if _, err := domain.ParseReference(source.Ref); err != nil {
		return invalid("authoritative source must use an exact typed reference")
	}
	return validateFingerprint(source.Fingerprint)
}

func validateFingerprint(value string) error {
	if len(value) != 64 || strings.ToLower(value) != value {
		return invalid("source fingerprint must be lowercase SHA-256 hex")
	}
	if _, err := hex.DecodeString(value); err != nil {
		return invalid("source fingerprint must be lowercase SHA-256 hex")
	}
	return nil
}

func invalid(detail string) error {
	return domain.NewRefusal(domain.RefusalInvalidKnownField, "runtime-input", detail, nil)
}

func oneOf(value string, allowed ...string) bool {
	for _, candidate := range allowed {
		if value == candidate {
			return true
		}
	}
	return false
}

func contains(values []string, want string) bool {
	for _, value := range values {
		if value == want {
			return true
		}
	}
	return false
}

func allContained(values, allowed []string) bool {
	for _, value := range values {
		if !contains(allowed, value) {
			return false
		}
	}
	return true
}

func overlaps(left, right []string) bool {
	for _, value := range left {
		if contains(right, value) {
			return true
		}
	}
	return false
}

func fingerprint(value any) string {
	data, _ := json.Marshal(value)
	sum := sha256.Sum256(data)
	return hex.EncodeToString(sum[:])
}

// CompletionFingerprint identifies only the frozen answer key, allowing
// execution records to change without obscuring a semantic criteria change.
func CompletionFingerprint(criteria []CompletionCriterion) string {
	canonical := append([]CompletionCriterion(nil), criteria...)
	sort.Slice(canonical, func(i, j int) bool { return canonical[i].Claim < canonical[j].Claim })
	return fingerprint(canonical)
}

// CanonicalForbiddenEffects returns a copy so guided workflows can render the
// fixed initial ceiling without duplicating it as editable prose.
func CanonicalForbiddenEffects() []string {
	out := append([]string{}, requiredForbidden...)
	sort.Strings(out)
	return out
}
