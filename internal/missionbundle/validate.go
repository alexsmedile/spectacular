package missionbundle

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
	"time"

	"github.com/alexsmedile/spectacular/v2/internal/discovery"
	"github.com/alexsmedile/spectacular/v2/internal/domain"
	"github.com/alexsmedile/spectacular/v2/internal/workspace"
)

type validator struct {
	name string
	run  func(*discovery.Workspace, *Bundle) error
}

var registry = []validator{
	{"yaml-schema", validateSchema},
	{"uuidv7-identity", validateIdentity},
	{"reference-integrity", validateReferences},
	{"contract-binding", validateContract},
	{"baseline-binding", validateBaseline},
	{"activation-fingerprint", validateActivation},
	{"completion-claim-coverage", validateClaims},
	{"frozen-fallbacks", validateFallbacks},
	{"request-coverage", validateRequest},
	{"objective-dependency-dag", validateDAG},
	{"mission-order-integrity", validateMissionOrderIntegrity},
	{"mission-order-activation", validateMissionOrderActivation},
	{"run-state", validateRun},
	{"review-independence", validateReviews},
	{"authority-vocabulary", validateAuthority},
	{"mechanical-scope", validateScope},
	{"safe-file-layout", validateLayout},
	{"transition-atomicity", validateTransactions},
}

var (
	missionRefPattern = regexp.MustCompile(`^M[1-9][0-9]*$`)
	objectivePattern  = regexp.MustCompile(`^O[1-9][0-9]*$`)
	runPattern        = regexp.MustCompile(`^R[1-9][0-9]*$`)
	reviewPattern     = regexp.MustCompile(`^RV[1-9][0-9]*$`)
	shaPattern        = regexp.MustCompile(`^sha256:[0-9a-f]{64}$`)
	commitPattern     = regexp.MustCompile(`^[0-9a-f]{40}$`)
)

func Validate(ws *discovery.Workspace, b *Bundle) (Check, error) {
	check := Check{Ref: b.Ref, Path: b.Path, Schema: b.Validation.Schema, Valid: true, Legacy: b.Legacy}
	if b.Legacy {
		check.Schema = "legacy-v2"
		check.Checks = []string{"canonical-record", "uuidv7-identity", "read-only-legacy-decoder"}
		// Legacy records are exactly where the ref-spelling drift lives, so the
		// notice is reported before returning rather than skipped with the
		// mission.v2 validators.
		check.Notices = b.Notices()
		return check, nil
	}
	for _, item := range registry {
		if err := item.run(ws, b); err != nil {
			return Check{}, err
		}
		check.Checks = append(check.Checks, item.name)
	}
	if b.Activation != nil {
		check.Fingerprint = b.Activation.Fingerprint
	}
	check.Authority = b.AuthorityTable()
	check.Drift = b.Drift()
	check.Notices = b.Notices()
	return check, nil
}

func validateSchema(_ *discovery.Workspace, b *Bundle) error {
	if b.Validation.Schema != Schema {
		return invalid("validation.schema", "must be mission.v2")
	}
	if b.Validation.Mode != "planned" && b.Validation.Mode != "manual-bootstrap" && b.Validation.Mode != "cli" {
		return invalid("validation.mode", "must be planned, manual-bootstrap, or cli")
	}
	if b.Title == "" || b.Owner == "" || b.Outcome == "" || len(b.Completion) == 0 || len(b.Objectives) == 0 || b.RepairBudget < 0 || len(b.Stops) == 0 {
		return invalid("mission", "title, owner, outcome, completion, objectives, non-negative repair_budget, and stops are required")
	}
	if b.Review != "automatic" && b.Review != "clustered" && b.Review != "independent" {
		return invalid("review", "must be automatic, clustered, or independent")
	}
	if b.Status != "defined" && b.Status != "active" && b.Status != "completed" {
		return invalid("status", "compact Mission status must be defined, active, or completed")
	}
	if b.Status == "defined" && (b.Baseline != nil || b.Activation != nil || b.Run != nil || len(b.Runs) != 0) {
		return invalid("status", "defined Mission must not contain baseline, activation, or Run state")
	}
	if b.Status != "defined" && (b.Baseline == nil || b.Activation == nil || (b.Run == nil && len(b.Runs) == 0)) {
		return invalid("status", "active or completed Mission requires baseline, activation, and Run state")
	}
	return nil
}

func validateIdentity(_ *discovery.Workspace, b *Bundle) error {
	if _, err := domain.ParseID(b.ID); err != nil {
		return invalidCause("id", "must be canonical UUIDv7", err)
	}
	if !missionRefPattern.MatchString(b.Ref) {
		return invalid("ref", "Mission ref must match M<number>")
	}
	seenIDs := map[string]bool{b.ID: true}
	seenRefs := map[string]bool{}
	for _, objective := range b.Objectives {
		if _, err := domain.ParseID(objective.ID); err != nil {
			return invalidCause("objectives.id", "must be canonical UUIDv7", err)
		}
		if !objectivePattern.MatchString(objective.Ref) || seenRefs[objective.Ref] || seenIDs[objective.ID] {
			return invalid("objectives.ref", "Objective refs and identities must be unique O<number> values")
		}
		seenRefs[objective.Ref], seenIDs[objective.ID] = true, true
	}
	for _, run := range allRuns(b) {
		if _, err := domain.ParseID(run.ID); err != nil {
			return invalidCause("runs.id", "must be canonical UUIDv7", err)
		}
		if !runPattern.MatchString(run.Ref) || seenRefs[run.Ref] || seenIDs[run.ID] {
			return invalid("runs.ref", "Run refs and identities must be unique R<number> values")
		}
		seenRefs[run.Ref], seenIDs[run.ID] = true, true
	}
	for _, review := range b.Reviews {
		if _, err := domain.ParseID(review.ID); err != nil {
			return invalidCause("reviews.id", "must be canonical UUIDv7", err)
		}
		if !reviewPattern.MatchString(review.Ref) || seenRefs[review.Ref] || seenIDs[review.ID] {
			return invalid("reviews.ref", "review refs and identities must be unique RV<number> values")
		}
		seenRefs[review.Ref], seenIDs[review.ID] = true, true
	}
	return nil
}

func validateReferences(_ *discovery.Workspace, b *Bundle) error {
	for _, objective := range b.Objectives {
		if objective.Outcome == "" || objective.Status == "" || len(objective.Claims) == 0 {
			return invalid("objectives", "each Objective needs outcome, status, and claims")
		}
		if objective.Status != "pending" && objective.Status != "active" && objective.Status != "implemented" {
			return invalid("objectives.status", "must be pending, active, or implemented")
		}
	}
	return nil
}

func validateContract(ws *discovery.Workspace, b *Bundle) error {
	typed, err := domain.ParseReference(b.Contract.Ref)
	if err != nil || typed.Type != domain.Contract {
		return invalidCause("contract.ref", "must be an exact Contract:<UUIDv7> reference", err)
	}
	if !shaPattern.MatchString(b.Contract.Fingerprint) {
		return invalid("contract.fingerprint", "must be sha256:<64 lowercase hex>")
	}
	entry, err := ws.Lookup(b.Contract.Ref, domain.Contract)
	if err != nil {
		return err
	}
	data, err := os.ReadFile(entry.Absolute)
	if err != nil {
		return invalidCause("contract.ref", "cannot read bound Contract", err)
	}
	digest := sha256.Sum256(data)
	actual := "sha256:" + hex.EncodeToString(digest[:])
	if actual != b.Contract.Fingerprint {
		return domain.NewStateRefusal(domain.RefusalStaleFingerprint, "contract.fingerprint", "bound Contract content changed", b.Contract.Fingerprint, actual, "review the Contract delta and explicitly amend or restart the Mission", nil)
	}
	return nil
}

func validateBaseline(ws *discovery.Workspace, b *Bundle) error {
	if b.Status == "defined" {
		return nil
	}
	if !commitPattern.MatchString(b.Baseline.Commit) || strings.TrimSpace(b.Baseline.Branch) == "" {
		return invalid("baseline", "requires a full lowercase Git commit and branch")
	}
	cmd := exec.Command("git", "cat-file", "-e", b.Baseline.Commit+"^{commit}")
	cmd.Dir = ws.Root
	if err := cmd.Run(); err != nil {
		return invalidCause("baseline.commit", "commit is not available in this repository", err)
	}
	return nil
}

func validateActivation(_ *discovery.Workspace, b *Bundle) error {
	if b.Status == "defined" {
		return nil
	}
	if b.Activation.By == "" || !timestamp(b.Activation.At) || !shaPattern.MatchString(b.Activation.Fingerprint) {
		return invalid("activation", "requires owner, RFC3339 time, and sha256 fingerprint")
	}
	want, err := FrozenFingerprint(b)
	if err != nil {
		return err
	}
	if want != b.Activation.Fingerprint {
		return domain.NewStateRefusal(domain.RefusalStaleFingerprint, "activation.fingerprint", "frozen Mission boundary changed", b.Activation.Fingerprint, want, "return the semantic change to the owner and record a new activation boundary", nil)
	}
	return nil
}

func FrozenFingerprint(b *Bundle) (string, error) {
	value := map[string]any{
		"outcome": b.Outcome, "review": b.Review, "completion": b.Completion,
		"authority": b.Authority, "scope": b.Scope, "repair_budget": b.RepairBudget,
		"dependencies": b.Dependencies, "gaps": b.Gaps, "stops": b.Stops,
	}
	if len(b.Fallbacks) > 0 {
		value["fallbacks"] = b.Fallbacks
	}
	if len(b.AfterMission) > 0 {
		value["after_mission"] = b.AfterMission
	}
	if b.Request != nil {
		value["request_dispositions"] = b.Request.dispositions()
	}
	data, err := json.Marshal(value)
	if err != nil {
		return "", err
	}
	digest := sha256.Sum256(data)
	return "sha256:" + hex.EncodeToString(digest[:]), nil
}

func validateMissionOrderIntegrity(ws *discovery.Workspace, b *Bundle) error {
	if len(b.AfterMission) == 0 {
		return nil
	}
	for _, predRef := range b.AfterMission {
		ref := strings.TrimSpace(predRef)
		if ref == "" {
			return invalid("after_mission", "empty Mission ref in after_mission")
		}
		if ref == b.Ref || ref == b.ID {
			return invalid("after_mission", "Mission cannot depend on itself: "+ref)
		}
		entry, err := ws.Lookup(ref, domain.Mission)
		if err != nil {
			return invalidCause("after_mission", "dangling Mission ref: "+ref, err)
		}
		if entry.Document.Record.ID.String() == b.ID {
			return invalid("after_mission", "Mission cannot depend on itself: "+ref)
		}
	}

	visiting, done := map[string]bool{}, map[string]bool{}
	var visit func(currentRef string) error
	visit = func(currentRef string) error {
		if visiting[currentRef] {
			return invalid("after_mission", "Mission order must be acyclic")
		}
		if done[currentRef] {
			return nil
		}
		visiting[currentRef] = true
		entry, err := ws.Lookup(currentRef, domain.Mission)
		if err == nil {
			targetBundle, decodeErr := decode(ws, entry)
			if decodeErr == nil && targetBundle != nil {
				for _, nextPred := range targetBundle.AfterMission {
					if err := visit(nextPred); err != nil {
						return err
					}
				}
			}
		}
		delete(visiting, currentRef)
		done[currentRef] = true
		return nil
	}
	visiting[b.Ref] = true
	for _, predRef := range b.AfterMission {
		if err := visit(predRef); err != nil {
			return err
		}
	}
	return nil
}

func validateMissionOrderActivation(ws *discovery.Workspace, b *Bundle) error {
	if len(b.AfterMission) == 0 || b.Status != "active" {
		return nil
	}
	for _, predRef := range b.AfterMission {
		ref := strings.TrimSpace(predRef)
		entry, err := ws.Lookup(ref, domain.Mission)
		if err != nil {
			return invalidCause("after_mission", "cannot load predecessor Mission: "+ref, err)
		}
		pred, err := decode(ws, entry)
		if err != nil {
			return invalidCause("after_mission", "cannot decode predecessor Mission: "+ref, err)
		}
		if pred.Status != "completed" {
			return invalid("after_mission", fmt.Sprintf("predecessor Mission %s is not completed (status: %s)", ref, pred.Status))
		}
	}
	return nil
}

func validateFallbacks(_ *discovery.Workspace, b *Bundle) error {
	if len(b.Fallbacks) == 0 {
		return nil
	}
	for _, fb := range b.Fallbacks {
		if strings.TrimSpace(fb.Approach) == "" {
			return invalid("fallbacks.approach", "fallback requires approach")
		}
		if strings.TrimSpace(fb.RejectedBecause) == "" {
			return invalid("fallbacks.rejected_because", "fallback requires rejected_because")
		}
		if strings.TrimSpace(fb.InvalidatedIf) == "" {
			return invalid("fallbacks.invalidated_if", "fallback requires invalidated_if")
		}
	}
	return nil
}

func validateRequest(_ *discovery.Workspace, b *Bundle) error {
	if b.Request == nil {
		return nil
	}
	if strings.TrimSpace(b.Request.Source) == "" {
		return invalid("request.source", "request record requires a source")
	}
	if !timestamp(b.Request.CapturedAt) {
		return invalid("request.captured_at", "request record requires an RFC3339 captured_at timestamp")
	}
	if len(b.Request.Asks) == 0 {
		return invalid("request.asks", "request record requires at least one ask")
	}
	criteria := map[string]bool{}
	for _, criterion := range b.Completion {
		criteria[criterion.Claim] = true
	}
	for _, ask := range b.Request.Asks {
		if strings.TrimSpace(ask.Ask) == "" {
			return invalid("request.asks.ask", "request ask cannot be empty")
		}
		switch ask.Disposition {
		case "covered":
			if len(ask.Claims) == 0 {
				return invalid("request.asks.claims", "covered ask must declare at least one claim")
			}
			for _, claim := range ask.Claims {
				if !criteria[claim] {
					return invalid("request.asks.claims", "covered ask references an unknown completion claim")
				}
			}
		case "deferred", "declined":
			if strings.TrimSpace(ask.Reason) == "" {
				return invalid("request.asks.reason", "deferred or declined ask must declare a reason")
			}
		default:
			return invalid("request.asks.disposition", "ask disposition must be covered, deferred, or declined")
		}
	}
	return nil
}

func validateClaims(_ *discovery.Workspace, b *Bundle) error {
	criteria := map[string]bool{}
	for _, criterion := range b.Completion {
		if criterion.Claim == "" || criterion.PassBoundary == "" || criterion.ProofRequirement == "" || criteria[criterion.Claim] {
			return invalid("completion", "claims must be unique and include pass_boundary and proof_requirement")
		}
		criteria[criterion.Claim] = true
	}
	covered := map[string]bool{}
	for _, objective := range b.Objectives {
		for _, claim := range objective.Claims {
			if !criteria[claim] {
				return invalid("objectives.claims", "Objective references an unknown completion claim")
			}
			covered[claim] = true
		}
	}
	for claim := range criteria {
		if !covered[claim] {
			return invalid("completion", "every completion claim must be owned by an Objective")
		}
	}
	return nil
}

func validateDAG(_ *discovery.Workspace, b *Bundle) error {
	known := map[string]bool{}
	for _, objective := range b.Objectives {
		known[objective.Ref] = true
	}
	visiting, done := map[string]bool{}, map[string]bool{}
	deps := map[string][]string{}
	for _, objective := range b.Objectives {
		deps[objective.Ref] = append(append([]string(nil), objective.After...), objective.AfterInterface...)
		for _, dependency := range objective.After {
			if !known[dependency] || dependency == objective.Ref {
				return invalid("objectives.after", "dependencies must reference another Objective in the Mission: "+dependency)
			}
		}
		for _, dependency := range objective.AfterInterface {
			if dependency == objective.Ref {
				return invalid("objectives.after_interface", "interface dependency cannot reference self: "+dependency)
			}
			if !known[dependency] {
				return invalid("objectives.after_interface", "interface dependency references an unknown or unfrozen target: "+dependency)
			}
		}
	}
	var visit func(string) error
	visit = func(ref string) error {
		if visiting[ref] {
			return invalid("objectives.after", "Objective dependencies must be acyclic")
		}
		if done[ref] {
			return nil
		}
		visiting[ref] = true
		for _, dependency := range deps[ref] {
			if err := visit(dependency); err != nil {
				return err
			}
		}
		delete(visiting, ref)
		done[ref] = true
		return nil
	}
	for ref := range known {
		if err := visit(ref); err != nil {
			return err
		}
	}
	return nil
}

func validateRun(_ *discovery.Workspace, b *Bundle) error {
	if b.Status == "defined" {
		return nil
	}
	runs := allRuns(b)
	active := 0
	objectiveRefs := map[string]bool{}
	for _, objective := range b.Objectives {
		objectiveRefs[objective.Ref] = true
	}
	for _, run := range runs {
		if run.Status != "active" && run.Status != "awaiting-review" && run.Status != "completed" {
			return invalid("run.status", "must be active, awaiting-review, or completed")
		}
		if run.Status == "active" || run.Status == "awaiting-review" {
			active++
		}
		if run.Operator == "" || !timestamp(run.StartedAt) || !objectiveRefs[run.CurrentObjective] || run.Repairs < 0 || run.Repairs > b.RepairBudget {
			return invalid("run", "requires operator, RFC3339 start, current Objective, and repairs within budget")
		}
	}
	if b.Status == "active" && active != 1 {
		return invalid("run.status", "active Mission requires exactly one active or awaiting-review Run")
	}
	if b.Status == "completed" {
		for _, run := range runs {
			if run.Status != "completed" {
				return invalid("run.status", "completed Mission requires completed Runs")
			}
		}
		for _, objective := range b.Objectives {
			if objective.Status != "implemented" {
				return invalid("objectives.status", "completed Mission requires every Objective implemented")
			}
		}
		if b.CompletionRecord == nil || b.CompletionRecord.By != b.Owner || !timestamp(b.CompletionRecord.At) || b.CompletionRecord.Authorization == "" {
			return invalid("completion_record", "completed Mission requires attributable owner completion")
		}
	}
	return nil
}

func validateReviews(ws *discovery.Workspace, b *Bundle) error {
	if b.Status != "completed" || b.Review == "automatic" {
		return nil
	}
	if len(b.Reviews) == 0 {
		return invalid("reviews", "completed Mission requires its declared review")
	}
	base := filepath.Dir(b.entry.Absolute)
	passed := false
	for _, pointer := range b.Reviews {
		if pointer.Verdict != "pass" {
			continue
		}
		path, err := containedFile(base, pointer.File)
		if err != nil {
			return err
		}
		doc, err := workspace.ReadFile(path)
		if err != nil {
			return err
		}
		if doc.Record.Type != domain.Review || doc.Record.ID.String() != pointer.ID || value(doc.Record.Status) != "passed" {
			return invalid("reviews.file", "review pointer must resolve to the same passed Review identity")
		}
		if b.Review == "independent" {
			var reviewer Reviewer
			if err := workspace.DecodeValue(doc, "reviewer", &reviewer); err != nil {
				return err
			}
			if reviewer.Actor == "" || reviewer.Operator == "" || reviewer.Actor == reviewer.Operator || reviewer.RelationToOperator != "independent" || reviewer.ImplementedReviewedScope || reviewer.IndependenceBasis == "" || len(reviewer.Evidence) == 0 {
				return invalid("reviews.reviewer", "independent review requires distinct identities, non-implementation, basis, and attributable evidence")
			}
		}
		var reviewed struct {
			Commit                string `yaml:"commit"`
			Tree                  string `yaml:"tree"`
			ActivationFingerprint string `yaml:"activation_fingerprint"`
		}
		if err := workspace.DecodeValue(doc, "reviewed", &reviewed); err != nil {
			return err
		}
		if !commitPattern.MatchString(reviewed.Commit) || !commitPattern.MatchString(reviewed.Tree) || reviewed.ActivationFingerprint != b.Activation.Fingerprint {
			return invalid("reviews.reviewed", "review must bind an exact commit, tree, and current activation fingerprint")
		}
		if err := verifyReviewedGit(ws.Root, reviewed.Commit, reviewed.Tree); err != nil {
			return err
		}
		var claims []struct {
			Claim   string `yaml:"claim"`
			Verdict string `yaml:"verdict"`
		}
		if err := workspace.DecodeValue(doc, "claims", &claims); err != nil {
			return err
		}
		seen := map[string]bool{}
		for _, claim := range claims {
			if claim.Verdict != "pass" {
				return invalid("reviews.claims", "qualifying review claims must pass")
			}
			seen[claim.Claim] = true
		}
		for _, criterion := range b.Completion {
			if !seen[criterion.Claim] {
				return invalid("reviews.claims", "review must cover every frozen completion claim")
			}
		}
		passed = true
	}
	if !passed {
		return invalid("reviews", "no qualifying passed review was found")
	}
	return nil
}

// SupportedOperatorVerbs and SupportedOwnerVerbs are the closed vocabularies a
// Mission may declare. They are package-level so the authority table rendered to
// a reader and the vocabulary enforced by validation cannot drift apart.
var SupportedOperatorVerbs = []string{"inspect", "edit-in-scope", "choose-reversible-implementation", "run-checks", "generate-derived-files", "bounded-repair", "commit-local"}

var SupportedOwnerVerbs = []string{"activate-mission", "change-outcome-or-completion", "expand-scope", "push", "merge", "release", "irreversible-change", "destructive-data", "secret-change"}

func validateAuthority(_ *discovery.Workspace, b *Bundle) error {
	operator := allowedSet(SupportedOperatorVerbs...)
	owner := allowedSet(SupportedOwnerVerbs...)
	if len(b.Authority.Operator) == 0 || len(b.Authority.RequiresOwner) == 0 {
		return invalid("authority", "operator and requires_owner vocabularies are required")
	}
	for _, item := range b.Authority.Operator {
		if !operator[item] {
			return invalid("authority.operator", "contains an unsupported action")
		}
	}
	for _, item := range b.Authority.RequiresOwner {
		if !owner[item] {
			return invalid("authority.requires_owner", "contains an unsupported action")
		}
	}
	return nil
}

func validateScope(_ *discovery.Workspace, b *Bundle) error {
	if len(b.Scope.Mechanical) == 0 || len(b.Scope.Semantic) == 0 {
		return invalid("scope", "mechanical and semantic scope are required")
	}
	for _, path := range b.Scope.Mechanical {
		clean := strings.TrimSuffix(path, "/")
		if clean == "" || filepath.IsAbs(clean) || filepath.Clean(clean) != clean || strings.HasPrefix(clean, ".."+string(filepath.Separator)) || strings.Contains(clean, "\\") {
			return invalid("scope.mechanical", "paths must be canonical workspace-relative paths")
		}
	}
	for _, item := range b.Scope.Semantic {
		if strings.TrimSpace(item) == "" {
			return invalid("scope.semantic", "items must be non-empty")
		}
	}
	return nil
}

func validateLayout(ws *discovery.Workspace, b *Bundle) error {
	if filepath.Base(b.Path) != "MISSION.md" || !strings.HasPrefix(filepath.ToSlash(b.Path), ".spectacular/missions/") {
		return invalid("path", "Mission entry point must be .spectacular/missions/<bundle>/MISSION.md")
	}
	if _, err := os.Lstat(filepath.Join(ws.Root, filepath.FromSlash(b.Path))); err != nil {
		return invalidCause("path", "Mission entry point is unavailable", err)
	}
	return nil
}

func validateTransactions(ws *discovery.Workspace, _ *Bundle) error {
	entries, err := os.ReadDir(filepath.Join(ws.MetadataDir, "transactions"))
	if os.IsNotExist(err) {
		return nil
	}
	if err != nil {
		return invalidCause("transactions", "cannot inspect transaction recovery state", err)
	}
	for _, entry := range entries {
		if !entry.IsDir() && strings.HasSuffix(entry.Name(), ".json") {
			return domain.NewStateRefusal(domain.RefusalTransactionRecovery, "transactions", "an incomplete atomic transition requires recovery", "no pending journal", entry.Name(), "run a mutating command to invoke deterministic recovery, then re-check", nil)
		}
	}
	return nil
}

func allRuns(b *Bundle) []Run {
	if len(b.Runs) > 0 {
		return b.Runs
	}
	if b.Run != nil {
		return []Run{*b.Run}
	}
	return nil
}

func timestamp(raw string) bool {
	_, err := time.Parse(time.RFC3339, raw)
	return err == nil
}

func allowedSet(values ...string) map[string]bool {
	result := make(map[string]bool, len(values))
	for _, value := range values {
		result[value] = true
	}
	return result
}

func invalid(field, detail string) error {
	return domain.NewStateRefusal(domain.RefusalInvalidKnownField, field, detail, "schema-owned invariant", "invalid", "correct the named field and retry; no files were changed", nil)
}

func invalidCause(field, detail string, cause error) error {
	return domain.NewStateRefusal(domain.RefusalInvalidKnownField, field, detail, "schema-owned invariant", fmt.Sprint(cause), "correct the named field and retry; no files were changed", cause)
}

func ValidatorNames() []string {
	names := make([]string, 0, len(registry))
	for _, item := range registry {
		names = append(names, item.name)
	}
	sort.Strings(names)
	return names
}
