package command

import (
	"bytes"
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/alexsmedile/spectacular/v2/internal/discovery"
	"github.com/alexsmedile/spectacular/v2/internal/domain"
	"github.com/alexsmedile/spectacular/v2/internal/humanlayout"
	spectacularruntime "github.com/alexsmedile/spectacular/v2/internal/runtime"
	"github.com/alexsmedile/spectacular/v2/internal/workspace"
)

func TestScenarioSCleanV2DogfoodSurvivesRuntimeReplacement(t *testing.T) {
	root := copyFixture(t)
	seedScenarioSAuthority(t, root)
	before := snapshot(t, root)
	opened, err := discovery.Open(root)
	if err != nil {
		t.Fatal(err)
	}
	proposal, err := opened.Lookup("Proposal:"+missionSourceID, domain.Proposal)
	if err != nil {
		t.Fatal(err)
	}
	mission, err := opened.Lookup("Mission:"+missionID, domain.Mission)
	if err != nil {
		t.Fatal(err)
	}
	decision, err := opened.Lookup("Decision:"+autopilotDecisionID, domain.Decision)
	if err != nil {
		t.Fatal(err)
	}

	firstContext := runJSON(t, root, []string{"workspace", "context", "project", "--event", "@Orient", "--json"})
	if !strings.Contains(firstContext, `"schema_version":"spectacular.context.v1"`) || !strings.Contains(firstContext, `"kind":"continuation"`) {
		t.Fatalf("cold orientation lacks bounded continuation: %s", firstContext)
	}

	preparation := spectacularruntime.PreparationInput{
		Proposal: spectacularruntime.BoundSource{Ref: "Proposal:" + missionSourceID, Fingerprint: proposal.Fingerprint}, Baseline: "git:clean-v2-dogfood",
		DirectionSources: []spectacularruntime.BoundSource{{Ref: "Proposal:" + missionSourceID, Fingerprint: proposal.Fingerprint}},
		Candidates:       []spectacularruntime.CandidateSlice{{Name: "guided-vertical", Outcome: "compile safe runtime context", Evidence: []string{"cold-runtime comparison"}, CancellationState: "usable context compiler", Reversibility: "local fixture only", StandaloneCoherence: "complete guided slice", IntegrationPath: "registry to Skill", LearningValue: "runtime-neutral proof"}},
		Selected:         "guided-vertical", DesignSufficiency: "sufficient", DesignRationale: "accepted contracts and executable core fix the boundary", SliceQuality: "coherent", SliceRationale: "one integrated Skill/runtime seam",
		CompletionCriteria: []spectacularruntime.CompletionCriterion{{Claim: "claim:runtime-neutral-recovery", PassBoundary: "cold replacement emits the same source basis", ProofRequirement: "real-process scenario", ReviewLevel: spectacularruntime.ReviewIndependent}}, StopConditions: []string{"authority drift"}, EvidenceClaims: []string{"claim:runtime-neutral-recovery"}, FreshUntil: "2026-08-11T10:00:00Z",
	}
	preparationPath := writeJSONInput(t, preparation)
	prepared := runJSON(t, root, []string{"mission", "prepare", "--input", preparationPath, "--json"})
	if !strings.Contains(prepared, `"ready":true`) || !strings.Contains(prepared, `"next":"owner-activation"`) {
		t.Fatalf("preparation did not compile a ready owner gate: %s", prepared)
	}
	stalePreparation := preparation
	stalePreparation.Proposal.Fingerprint = strings.Repeat("0", 64)
	stalePath := writeJSONInput(t, stalePreparation)
	var staleOut, staleErr bytes.Buffer
	if exit := (Runner{Cwd: root, Stdout: &staleOut, Stderr: &staleErr, Now: fixedNow}).Run([]string{"mission", "prepare", "--input", stalePath, "--json"}); exit != 3 || !strings.Contains(staleOut.String(), `"code":"stale_fingerprint"`) {
		t.Fatalf("stale preparation source exit=%d stdout=%s stderr=%s", exit, staleOut.String(), staleErr.String())
	}

	autopilot := spectacularruntime.AutopilotInput{
		Enabled: true, Mission: spectacularruntime.BoundSource{Ref: "Mission:" + missionID, Fingerprint: mission.Fingerprint}, Authorization: spectacularruntime.BoundSource{Ref: "Decision:" + autopilotDecisionID, Fingerprint: decision.Fingerprint},
		Outcome: "reproduce cold v2 orientation", NonGoals: []string{"provider effects", "Mission resolution"}, AuthoritativeSources: []spectacularruntime.BoundSource{{Ref: "Proposal:" + missionSourceID, Fingerprint: proposal.Fingerprint}},
		DelegatedDecisionDomain: []string{"reversible local verification"}, AllowedProviders: []string{}, AllowedActions: []string{"inspect", "test"}, ForbiddenEffects: spectacularruntime.CanonicalForbiddenEffects(),
		BudgetUnits: 2, RepairBudget: 1, RequiredChecks: []string{"compare generation basis"}, ExpiresAt: "2026-08-11T10:00:00Z", StopConditions: []string{"baseline drift", "authority drift"}, RecoveryPoint: "git:clean-v2-dogfood", ReturnDestination: "Mission owner",
	}
	autopilotPath := writeJSONInput(t, autopilot)
	charter := runJSON(t, root, []string{"mission", "autopilot", "--input", autopilotPath, "--json"})
	if !strings.Contains(charter, `"runtime_neutral":true`) || !strings.Contains(charter, `"automatic_provider_effects":false`) || !strings.Contains(charter, `"Mission accountability remains with the sender"`) {
		t.Fatalf("unsafe Autopilot charter: %s", charter)
	}
	unsafe := autopilot
	unsafe.AllowedProviders = []string{"github"}
	unsafe.AllowedActions = []string{"merge", "push"}
	unsafe.BudgetUnits = 999
	unsafePath := writeJSONInput(t, unsafe)
	var unsafeOut, unsafeErr bytes.Buffer
	if exit := (Runner{Cwd: root, Stdout: &unsafeOut, Stderr: &unsafeErr, Now: fixedNow}).Run([]string{"mission", "autopilot", "--input", unsafePath, "--json"}); exit != 3 || !strings.Contains(unsafeOut.String(), `"code":"unauthorized"`) {
		t.Fatalf("out-of-envelope Autopilot exit=%d stdout=%s stderr=%s", exit, unsafeOut.String(), unsafeErr.String())
	}

	// A replacement runtime receives only cwd and canonical v2 state. The
	// fixed clock makes byte equality a strict cold-context proof.
	secondContext := runJSON(t, root, []string{"workspace", "context", "project", "--event", "@Orient", "--json"})
	if firstContext != secondContext {
		t.Fatalf("replacement runtime changed compiled context\nfirst=%s\nsecond=%s", firstContext, secondContext)
	}
	if after := snapshot(t, root); !sameSnapshot(before, after) {
		t.Fatal("read-only Skill/runtime compilation mutated canonical dogfood state")
	}
}

func TestAutopilotRefusesStaleCanonicalMissionAndDecisionAuthority(t *testing.T) {
	for _, test := range []struct {
		name string
		path string
	}{
		{name: "Mission freshness source", path: fixtureEvidencePath},
		{name: "Decision freshness source", path: ".spectacular/workspace.yaml"},
	} {
		t.Run(test.name, func(t *testing.T) {
			root := copyFixture(t)
			seedScenarioSAuthority(t, root)
			opened, err := discovery.Open(root)
			if err != nil {
				t.Fatal(err)
			}
			proposal, err := opened.Lookup("Proposal:"+missionSourceID, domain.Proposal)
			if err != nil {
				t.Fatal(err)
			}
			mission, err := opened.Lookup("Mission:"+missionID, domain.Mission)
			if err != nil {
				t.Fatal(err)
			}
			decision, err := opened.Lookup("Decision:"+autopilotDecisionID, domain.Decision)
			if err != nil {
				t.Fatal(err)
			}
			input := spectacularruntime.AutopilotInput{
				Enabled: true, Mission: spectacularruntime.BoundSource{Ref: "Mission:" + missionID, Fingerprint: mission.Fingerprint}, Authorization: spectacularruntime.BoundSource{Ref: "Decision:" + autopilotDecisionID, Fingerprint: decision.Fingerprint},
				Outcome: "reproduce cold v2 orientation", NonGoals: []string{"provider effects"}, AuthoritativeSources: []spectacularruntime.BoundSource{{Ref: "Proposal:" + missionSourceID, Fingerprint: proposal.Fingerprint}},
				DelegatedDecisionDomain: []string{"reversible local verification"}, AllowedActions: []string{"inspect", "test"}, ForbiddenEffects: spectacularruntime.CanonicalForbiddenEffects(),
				BudgetUnits: 2, RepairBudget: 1, RequiredChecks: []string{"compare generation basis"}, ExpiresAt: "2026-08-11T10:00:00Z", StopConditions: []string{"authority drift"}, RecoveryPoint: "git:clean-v2-dogfood", ReturnDestination: "Mission owner",
			}
			path := filepath.Join(root, filepath.FromSlash(test.path))
			data, err := os.ReadFile(path)
			if err != nil {
				t.Fatal(err)
			}
			if err := os.WriteFile(path, append(data, []byte("\n# freshness drift\n")...), 0o644); err != nil {
				t.Fatal(err)
			}
			inputPath := writeJSONInput(t, input)
			var stdout, stderr bytes.Buffer
			exit := (Runner{Cwd: root, Stdout: &stdout, Stderr: &stderr, Now: fixedNow}).Run([]string{"mission", "autopilot", "--input", inputPath, "--json"})
			if exit != 3 || !strings.Contains(stdout.String(), `"code":"insufficient_evidence"`) {
				t.Fatalf("stale canonical authority exit=%d stdout=%s stderr=%s", exit, stdout.String(), stderr.String())
			}
		})
	}
}

func seedScenarioSAuthority(t *testing.T, root string) {
	t.Helper()
	opened, err := discovery.Open(root)
	if err != nil {
		t.Fatal(err)
	}
	mission, err := opened.Lookup("Mission:"+missionID, domain.Mission)
	if err != nil {
		t.Fatal(err)
	}
	workspace.SetString(mission.Document, "outcome", "reproduce cold v2 orientation")
	workspace.SetStrings(mission.Document, "scope", []string{"v2"})
	workspace.SetStrings(mission.Document, "allowed_actions", []string{"inspect", "test"})
	workspace.SetStrings(mission.Document, "forbidden_effects", spectacularruntime.CanonicalForbiddenEffects())
	workspace.SetStrings(mission.Document, "evidence_claims", []string{"claim:runtime-neutral-recovery"})
	workspace.SetInt(mission.Document, "budget_units", 4)
	workspace.SetInt(mission.Document, "repair_budget", 2)
	workspace.SetString(mission.Document, "expires_at", "2026-08-11T10:00:00Z")
	workspace.SetStrings(mission.Document, "stops", []string{"authority drift"})
	writeCanonicalEntry(t, mission)

	opened, err = discovery.Open(root)
	if err != nil {
		t.Fatal(err)
	}
	mission, err = opened.Lookup("Mission:"+missionID, domain.Mission)
	if err != nil {
		t.Fatal(err)
	}
	resumeDecision, err := opened.Lookup("Decision:"+decisionID, domain.Decision)
	if err != nil {
		t.Fatal(err)
	}
	workspace.SetString(resumeDecision.Document, "expected_mission_fingerprint", mission.Fingerprint)
	writeCanonicalEntry(t, resumeDecision)

	autopilotDecision := testDocument(t, domain.Decision, autopilotDecisionID, "Authorize bounded Autopilot", "recorded")
	workspace.SetString(autopilotDecision, "freshness_checked_at", "2026-08-10T10:00:00Z")
	workspace.SetString(autopilotDecision, "freshness_valid_until", "2026-08-11T10:00:00Z")
	workspace.SetString(autopilotDecision, "freshness_source", ".spectacular/workspace.yaml")
	_, markerFingerprint, err := opened.Source(".spectacular/workspace.yaml")
	if err != nil {
		t.Fatal(err)
	}
	workspace.SetString(autopilotDecision, "freshness_source_fingerprint", markerFingerprint)
	workspace.SetString(autopilotDecision, "actor_role", "owner")
	workspace.SetString(autopilotDecision, "operation", "mission.autopilot")
	workspace.SetStrings(autopilotDecision, "scope", []string{"v2"})
	workspace.SetString(autopilotDecision, "disposition", "approve")
	workspace.SetStrings(autopilotDecision, "targets", []string{"Mission:" + missionID})
	workspace.SetStrings(autopilotDecision, "expected_fingerprints", []string{mission.Fingerprint})
	workspace.SetStrings(autopilotDecision, "authorized_effects", []string{"mission.autopilot", "inspect", "test"})
	workspace.SetString(autopilotDecision, "expires_at", "2026-08-11T10:00:00Z")
	paths, err := humanlayout.Plan(opened.Entries, []*workspace.Document{autopilotDecision})
	if err != nil {
		t.Fatal(err)
	}
	writeTestDocument(t, filepath.Join(root, filepath.FromSlash(paths[autopilotDecision.Record.ID])), autopilotDecision)
}

func writeCanonicalEntry(t *testing.T, entry discovery.Entry) {
	t.Helper()
	data, err := workspace.Canonical(entry.Document)
	if err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(entry.Absolute, data, 0o644); err != nil {
		t.Fatal(err)
	}
}

const missionSourceID = "0198a1a0-0000-7000-8000-000000000001"
const autopilotDecisionID = "0198a1a0-0000-7000-8000-000000000009"

func runJSON(t *testing.T, root string, args []string) string {
	t.Helper()
	var stdout, stderr bytes.Buffer
	exit := (Runner{Cwd: root, Stdout: &stdout, Stderr: &stderr, Now: fixedNow}).Run(args)
	if exit != 0 {
		t.Fatalf("%v exit=%d stdout=%s stderr=%s", args, exit, stdout.String(), stderr.String())
	}
	return stdout.String()
}

func writeJSONInput(t *testing.T, value any) string {
	t.Helper()
	data, err := json.Marshal(value)
	if err != nil {
		t.Fatal(err)
	}
	path := filepath.Join(t.TempDir(), "input.json")
	if err := os.WriteFile(path, data, 0o644); err != nil {
		t.Fatal(err)
	}
	return path
}

func sameSnapshot(a, b map[string]fileState) bool {
	if len(a) != len(b) {
		return false
	}
	for path, left := range a {
		right, ok := b[path]
		if !ok || left.Mode != right.Mode || left.ModTime != right.ModTime || left.Data != right.Data {
			return false
		}
	}
	return true
}
