package missionbundle

import (
	"bytes"
	"crypto/sha256"
	"encoding/binary"
	"encoding/hex"
	"encoding/json"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/alexsmedile/spectacular/v2/internal/discovery"
	"github.com/alexsmedile/spectacular/v2/internal/domain"
	"github.com/alexsmedile/spectacular/v2/internal/governance"
	"github.com/alexsmedile/spectacular/v2/internal/humanlayout"
	"github.com/alexsmedile/spectacular/v2/internal/workspace"
	"github.com/google/uuid"
	"go.yaml.in/yaml/v3"
)

type Service struct {
	Workspace        *discovery.Workspace
	Now              func() time.Time
	ApplyTransaction func(root, key string, changes []governance.FileChange) error
	TokenAnalyzer    func(ws *discovery.Workspace, bundle *Bundle) *TokenEfficiency
}

type Plan struct {
	Type         string      `yaml:"type"`
	Title        string      `yaml:"title"`
	Owner        string      `yaml:"owner"`
	Contract     Binding     `yaml:"contract"`
	Outcome      string      `yaml:"outcome"`
	Request      *Request    `yaml:"request,omitempty"`
	Review       string      `yaml:"review"`
	Completion   []Criterion `yaml:"completion"`
	Objectives   []Objective `yaml:"objectives"`
	Authority    Authority   `yaml:"authority"`
	Scope        Scope       `yaml:"scope"`
	RepairBudget int         `yaml:"repair_budget"`
	Dependencies []string    `yaml:"dependencies"`
	Gaps         []string    `yaml:"gaps"`
	Stops        []string    `yaml:"stops"`
	Fallbacks    []Fallback  `yaml:"fallbacks,omitempty"`
	AfterMission []string    `yaml:"after_mission,omitempty"`
	// ResolvesGaps names Gaps on the bound Contract that this Mission closes at
	// completion, with the resolution text it will write. Frozen at activation so
	// the authority to amend a Contract cannot be acquired afterwards.
	ResolvesGaps []ResolvedGap `yaml:"resolves_gaps,omitempty"`
	AllowMain    bool          `yaml:"allow_main,omitempty"`
	CreateBranch bool          `yaml:"create_branch,omitempty"`
	Body         string        `yaml:"-"`
}

type ReviewDraft struct {
	Type     string `yaml:"type"`
	Title    string `yaml:"title"`
	Status   string `yaml:"status"`
	Reviewed struct {
		Commit                string `yaml:"commit"`
		Tree                  string `yaml:"tree"`
		ActivationFingerprint string `yaml:"activation_fingerprint"`
	} `yaml:"reviewed"`
	Reviewer Reviewer `yaml:"reviewer"`
	Claims   []struct {
		Claim   string `yaml:"claim"`
		Verdict string `yaml:"verdict"`
	} `yaml:"claims"`
	Findings    []string `yaml:"findings"`
	Limitations []string `yaml:"limitations"`
	Body        string   `yaml:"-"`
}

func ReadPlan(path string, stdin []byte) (Plan, []byte, error) {
	data, err := readInput(path, stdin)
	if err != nil {
		return Plan{}, nil, err
	}
	trimmed := bytes.TrimSpace(data)
	if len(trimmed) > 0 && trimmed[0] == '{' {
		var plan Plan
		if err := json.Unmarshal(trimmed, &plan); err != nil {
			return Plan{}, nil, invalidCause("input", "decode Mission plan JSON", err)
		}
		if plan.Type == "" {
			plan.Type = "MissionPlan"
		}
		if plan.Type != "MissionPlan" {
			return Plan{}, nil, invalid("type", "Mission start input must declare type: MissionPlan")
		}
		return plan, data, nil
	}
	frontmatter, body, err := splitInput(data)
	if err != nil {
		return Plan{}, nil, err
	}
	var plan Plan
	if err := yaml.Unmarshal(frontmatter, &plan); err != nil {
		return Plan{}, nil, invalidCause("input", "decode Mission plan frontmatter", err)
	}
	if plan.Type != "MissionPlan" {
		return Plan{}, nil, invalid("type", "Mission start input must declare type: MissionPlan")
	}
	plan.Body = body
	return plan, data, nil
}

func (s Service) Start(plan Plan, raw []byte) (Result, error) {
	locked, unlock, err := s.beginMutation()
	if err != nil {
		return Result{}, err
	}
	defer unlock()
	return locked.start(plan, raw)
}

func (s Service) start(plan Plan, raw []byte) (Result, error) {
	if s.Workspace == nil {
		return Result{}, invalid("workspace", "workspace is required")
	}
	startDigest := sha256.Sum256(raw)
	startKey := "sha256:" + hex.EncodeToString(startDigest[:])
	for _, entry := range s.Workspace.OfType(domain.Mission) {
		if existing, _ := workspace.String(entry.Document, "start_key", false); existing == startKey {
			ref, _ := compactRef(entry.Document)
			return Result{Operation: "mission.start", Ref: ref, Path: entry.Path, Fingerprint: entry.Fingerprint}, nil
		}
	}
	if err := validatePlan(plan); err != nil {
		return Result{}, err
	}
	plan.Dependencies = presentStrings(plan.Dependencies)
	plan.Gaps = presentStrings(plan.Gaps)
	plan.Stops = presentStrings(plan.Stops)
	contract, err := resolveContract(s.Workspace, plan.Contract.Ref)
	if err != nil {
		return Result{}, err
	}
	commit, branch, baselineAt, err := gitBaseline(s.Workspace.Root)
	if err != nil {
		return Result{}, err
	}
	missionRef := nextMissionRef(s.Workspace)
	if (branch == "main" || branch == "master") && !plan.AllowMain {
		if plan.CreateBranch {
			featBranch := "feat/" + missionRef + "-" + humanlayout.Slug(plan.Title)
			cmd := exec.Command("git", "checkout", "-b", featBranch)
			cmd.Dir = s.Workspace.Root
			if out, err := cmd.CombinedOutput(); err != nil {
				return Result{}, domain.NewStateRefusal(domain.RefusalInvalidKnownField, "baseline.branch", "failed to create feature branch: "+string(out), "clean feature branch", branch, "check git status and create branch manually", err)
			}
			branch = featBranch
		} else {
			return Result{}, domain.NewStateRefusal(domain.RefusalInvalidKnownField, "baseline.branch", "activating directly on "+branch+" is prohibited without a dedicated feature branch", "feature branch (e.g. feat/"+missionRef+"-"+humanlayout.Slug(plan.Title)+")", branch, "checkout a feature branch ('git checkout -b <branch>') or supply --create-branch / --allow-main", nil)
		}
	}
	missionID, err := stableID(baselineAt, startKey+":mission")
	if err != nil {
		return Result{}, err
	}
	now := s.now()
	objectives := make([]Objective, len(plan.Objectives))
	for i, source := range plan.Objectives {
		id, idErr := stableID(baselineAt, startKey+":objective:"+strconv.Itoa(i+1))
		if idErr != nil {
			return Result{}, idErr
		}
		source.Ref = "O" + strconv.Itoa(i+1)
		source.ID = id.String()
		if source.Status == "" {
			source.Status = "pending"
		}
		objectives[i] = source
	}
	runID, err := stableID(baselineAt, startKey+":run:R1")
	if err != nil {
		return Result{}, err
	}
	doc := &workspace.Document{Record: domain.Record{
		Type: domain.Mission, ID: missionID, Title: stringPtr(plan.Title), Status: stringPtr("active"),
		Created: stringPtr(now), Updated: stringPtr(now),
	}, Unknown: map[string]*yaml.Node{}, Body: plan.Body}
	workspace.SetString(doc, "ref", missionRef)
	workspace.SetString(doc, "owner", plan.Owner)
	workspace.SetValue(doc, "contract", contract)
	workspace.SetValue(doc, "baseline", Baseline{Commit: commit, Branch: branch})
	workspace.SetString(doc, "outcome", plan.Outcome)
	if plan.Request != nil {
		workspace.SetValue(doc, "request", plan.Request)
	}
	workspace.SetString(doc, "review", plan.Review)
	workspace.SetValue(doc, "completion", plan.Completion)
	workspace.SetValue(doc, "objectives", objectives)
	run := Run{Ref: "R1", ID: runID.String(), Status: "active", Operator: plan.Owner, StartedAt: now, CurrentObjective: objectives[0].Ref}
	workspace.SetValue(doc, "run", run)
	workspace.SetValue(doc, "validation", Validation{Schema: Schema, Mode: "cli"})
	workspace.SetValue(doc, "authority", plan.Authority)
	workspace.SetValue(doc, "scope", plan.Scope)
	workspace.SetInt(doc, "repair_budget", plan.RepairBudget)
	workspace.SetStrings(doc, "dependencies", plan.Dependencies)
	workspace.SetStrings(doc, "gaps", plan.Gaps)
	workspace.SetStrings(doc, "stops", plan.Stops)
	if len(plan.Fallbacks) > 0 {
		workspace.SetValue(doc, "fallbacks", plan.Fallbacks)
	}
	if len(plan.AfterMission) > 0 {
		workspace.SetStrings(doc, "after_mission", plan.AfterMission)
	}
	if len(plan.ResolvesGaps) > 0 {
		workspace.SetValue(doc, "resolves_gaps", plan.ResolvesGaps)
	}
	workspace.SetString(doc, "start_key", startKey)
	temporary := &Bundle{Outcome: plan.Outcome, Request: plan.Request, Review: plan.Review, Completion: plan.Completion, Authority: plan.Authority, Scope: plan.Scope, RepairBudget: plan.RepairBudget, Dependencies: plan.Dependencies, Gaps: plan.Gaps, Stops: plan.Stops, Fallbacks: plan.Fallbacks, AfterMission: plan.AfterMission, ResolvesGaps: plan.ResolvesGaps}
	fingerprint, err := FrozenFingerprint(temporary)
	if err != nil {
		return Result{}, err
	}
	workspace.SetValue(doc, "activation", Activation{By: plan.Owner, At: now, Fingerprint: fingerprint})
	slug := humanlayout.Slug(plan.Title)
	path := filepath.ToSlash(filepath.Join(".spectacular", "missions", missionRef+"-"+slug, missionRef+"-"+slug+".md"))
	candidate := &Bundle{
		ID: missionID.String(), Ref: missionRef, Title: plan.Title, Status: "active", Owner: plan.Owner,
		Contract: contract, Baseline: &Baseline{Commit: commit, Branch: branch}, Outcome: plan.Outcome,
		Request: plan.Request, Review: plan.Review, Completion: plan.Completion, Objectives: objectives, Run: &run,
		Activation: &Activation{By: plan.Owner, At: now, Fingerprint: fingerprint}, Validation: Validation{Schema: Schema, Mode: "cli"},
		Authority: plan.Authority, Scope: plan.Scope, RepairBudget: plan.RepairBudget,
		Dependencies: plan.Dependencies, Gaps: plan.Gaps, Stops: plan.Stops, Fallbacks: plan.Fallbacks, AfterMission: plan.AfterMission,
		ResolvesGaps: plan.ResolvesGaps, Path: path,
		document: doc,
	}
	for _, check := range registry {
		if check.name == "safe-file-layout" {
			continue
		}
		if err := check.run(s.Workspace, candidate); err != nil {
			return Result{}, err
		}
	}
	return s.apply("mission.start:"+startKey, []*workspace.Document{doc}, map[domain.ID]string{missionID: path}, "mission.start", missionRef, path)
}

func (s Service) AmendScope(ref string, addPaths []string, owner, reason string, dryRun bool) (Result, error) {
	if !dryRun {
		locked, unlock, err := s.beginMutation()
		if err != nil {
			return Result{}, err
		}
		defer unlock()
		return locked.amendScope(ref, addPaths, owner, reason, false)
	}
	fresh, err := discovery.Open(s.Workspace.Root)
	if err != nil {
		return Result{}, err
	}
	s.Workspace = fresh
	return s.amendScope(ref, addPaths, owner, reason, true)
}

func (s Service) amendScope(ref string, addPaths []string, owner, reason string, dryRun bool) (Result, error) {
	if strings.TrimSpace(owner) == "" {
		return Result{}, invalid("by", "an amendment requires the owner who authorizes it")
	}
	bundle, err := Load(s.Workspace, ref)
	if err != nil {
		return Result{}, err
	}
	if bundle.Legacy {
		return Result{}, invalid("mission", "legacy Mission is read-only")
	}
	if len(addPaths) == 0 {
		return Result{}, invalid("add", "at least one path must be specified to add to mechanical scope")
	}
	seen := map[string]bool{}
	for _, p := range bundle.Scope.Mechanical {
		seen[strings.TrimSuffix(p, "/")] = true
	}
	for _, p := range addPaths {
		clean := strings.TrimSuffix(p, "/")
		if clean == "" || filepath.IsAbs(clean) || filepath.ToSlash(filepath.Clean(clean)) != clean || strings.HasPrefix(clean, "../") || strings.Contains(clean, "\\") {
			return Result{}, invalid("scope.mechanical", "paths must be canonical workspace-relative paths: "+p)
		}
		if !seen[clean] {
			bundle.Scope.Mechanical = append(bundle.Scope.Mechanical, p)
			seen[clean] = true
		}
	}
	newFP, err := FrozenFingerprint(bundle)
	if err != nil {
		return Result{}, err
	}
	if dryRun {
		return Result{Operation: "mission.amend_scope", Ref: bundle.Ref, Path: bundle.Path, Fingerprint: newFP}, nil
	}
	if bundle.Activation != nil {
		bundle.Activation.Fingerprint = newFP
		workspace.SetValue(bundle.document, "activation", bundle.Activation)
	}
	workspace.SetValue(bundle.document, "scope", bundle.Scope)
	if strings.TrimSpace(reason) != "" {
		workspace.SetString(bundle.document, "scope_amendment_reason", reason)
	}
	bundle.document.Record.Updated = stringPtr(s.now())

	paths := map[domain.ID]string{bundle.document.Record.ID: bundle.Path}
	return s.apply("mission.amend_scope:"+bundle.ID+":"+newFP, []*workspace.Document{bundle.document}, paths, "mission.amend_scope", bundle.Ref, bundle.Path)
}

func validatePlan(plan Plan) error {
	if plan.Title == "" || plan.Owner == "" || plan.Outcome == "" || plan.Contract.Ref == "" || len(plan.Completion) == 0 || len(plan.Objectives) == 0 || len(plan.Stops) == 0 {
		return invalid("input", "plan requires title, owner, Contract ref, outcome, completion, Objectives, and stops")
	}
	if plan.Review != "automatic" && plan.Review != "clustered" && plan.Review != "independent" {
		return invalid("review", "must be automatic, clustered, or independent")
	}
	for _, objective := range plan.Objectives {
		if objective.Outcome == "" || len(objective.Claims) == 0 {
			return invalid("objectives", "each plan Objective needs outcome and claims")
		}
	}
	return nil
}

func nextMissionRef(ws *discovery.Workspace) string {
	max := 0
	for _, entry := range ws.OfType(domain.Mission) {
		ref := workspace.RefOrEmpty(entry.Document)
		if missionRefPattern.MatchString(ref) {
			value, _ := strconv.Atoi(strings.TrimPrefix(ref, "M"))
			if value > max {
				max = value
			}
		}
	}
	return "M" + strconv.Itoa(max+1)
}

func stableID(at, key string) (domain.ID, error) {
	instant, err := time.Parse(time.RFC3339, at)
	if err != nil {
		return "", invalidCause("id", "derive retry-stable UUIDv7 timestamp", err)
	}
	digest := sha256.Sum256([]byte(key))
	var raw [16]byte
	milliseconds := uint64(instant.UnixMilli())
	var encoded [8]byte
	binary.BigEndian.PutUint64(encoded[:], milliseconds)
	copy(raw[:6], encoded[2:])
	copy(raw[6:], digest[:10])
	raw[6] = (raw[6] & 0x0f) | 0x70
	raw[8] = (raw[8] & 0x3f) | 0x80
	return domain.ParseID(uuid.UUID(raw).String())
}

// missionRecordPath resolves a Mission-scoped record's canonical location through
// the layout system rather than through a join written at the call site, so a
// record created by any path lands where the layout rule says it lands. It
// returns the workspace-relative path and the bundle-relative pointer the
// Mission stores.
func (s Service) missionRecordPath(bundle *Bundle, doc *workspace.Document, ref string) (string, string, error) {
	// The layout system reads the scoped ref from the field workspace.Ref reads,
	// so the Mission scope is written there for the resolution and the record's
	// own ref is restored afterward. A record stores its leaf ref and names its
	// Mission in mission:; the scoped spelling is how the layout rule is asked
	// where that record belongs, not something the file carries.
	unscoped := workspace.RefOrEmpty(doc)
	workspace.SetString(doc, workspace.RefField, bundle.Ref+"/"+ref)
	path, err := humanlayout.PlannedPath(s.Workspace.Entries, doc)
	workspace.SetString(doc, workspace.RefField, unscoped)
	if err != nil {
		return "", "", err
	}
	path = filepath.ToSlash(path)
	bundleDirectory := filepath.ToSlash(filepath.Dir(bundle.Path))
	relative, err := filepath.Rel(bundleDirectory, path)
	if err != nil {
		return "", "", invalidCause("record", "resolve record path relative to its Mission", err)
	}
	return path, filepath.ToSlash(relative), nil
}

func (s Service) now() string {
	now := s.Now
	if now == nil {
		now = time.Now
	}
	return now().UTC().Truncate(time.Second).Format(time.RFC3339)
}

func presentStrings(values []string) []string {
	if values == nil {
		return []string{}
	}
	return values
}

func readInput(path string, stdin []byte) ([]byte, error) {
	if path == "-" || (path == "" && len(stdin) > 0) {
		if len(stdin) == 0 {
			return nil, invalid("input", "stdin is empty")
		}
		return stdin, nil
	}
	if path == "" {
		return nil, invalid("input", "input path is empty")
	}
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, domain.NewRefusal(domain.RefusalRecordNotFound, "input", "read Markdown input", err)
	}
	return data, nil
}

func splitInput(data []byte) ([]byte, string, error) {
	text := strings.ReplaceAll(string(data), "\r\n", "\n")
	if !strings.HasPrefix(text, "---\n") {
		return nil, "", invalid("input", "Markdown input must start with YAML frontmatter")
	}
	rest := strings.TrimPrefix(text, "---\n")
	index := strings.Index(rest, "\n---\n")
	if index < 0 {
		return nil, "", invalid("input", "Markdown input needs a closing --- delimiter")
	}
	return []byte(rest[:index] + "\n"), strings.TrimLeft(rest[index+5:], "\n"), nil
}

func stringPtr(value string) *string { return &value }
