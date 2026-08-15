// release-smoke drives the installed native binary through the accepted v2
// journey. Go is used only by this build-time proof harness; every product
// operation is a separate invocation of the installed artifact.
package main

import (
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io/fs"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"
	"time"
)

const (
	contractRef   = "Contract:0199b000-0000-7000-8000-000000000001"
	proposalRef   = "Proposal:0199c000-0000-7000-8000-000000000002"
	missionRef    = "Mission:0199c000-0000-7000-8000-000000000003"
	objectiveRef  = "Objective:0199c000-0000-7000-8000-000000000004"
	runRef        = "Run:0199c000-0000-7000-8000-000000000005"
	handoffRef    = "Handoff:0199c000-0000-7000-8000-000000000006"
	returnRef     = "Handoff:0199c000-0000-7000-8000-000000000007"
	evidenceRef   = "Evidence:0199c000-0000-7000-8000-000000000008"
	assessmentRef = "Assessment:0199c000-0000-7000-8000-000000000009"
)

type runner struct {
	binary    string
	workspace string
	inputs    string
	step      int
}

type envelope struct {
	Schema string         `json:"schema_version"`
	Data   map[string]any `json:"data"`
}

func main() {
	binary := flag.String("binary", "", "installed spectacular binary")
	fixture := flag.String("fixture", "testdata/scenario-b-c", "clean v2 fixture")
	workspace := flag.String("workspace", "", "empty disposable workspace destination")
	flag.Parse()
	if *binary == "" || *workspace == "" || flag.NArg() != 0 {
		fatal(errors.New("usage: release-smoke --binary <installed-binary> --workspace <empty-dir> [--fixture path]"))
	}
	absoluteBinary, err := filepath.Abs(*binary)
	if err != nil {
		fatal(err)
	}
	if err := copyFixture(*fixture, *workspace); err != nil {
		fatal(err)
	}
	inputRoot, err := os.MkdirTemp("", "spectacular-release-smoke-inputs-")
	if err != nil {
		fatal(err)
	}
	defer os.RemoveAll(inputRoot)
	r := &runner{binary: absoluteBinary, workspace: *workspace, inputs: inputRoot}
	if err := r.run(); err != nil {
		fatal(err)
	}
	fmt.Println("result=installed-binary-governed-closure-and-cold-resume-pass")
}

func (r *runner) run() error {
	version, err := exec.Command(r.binary, "--version").CombinedOutput()
	if err != nil || !validV2Version(strings.TrimSpace(string(version))) {
		return fmt.Errorf("runtime version inspection failed: %w: %s", err, version)
	}
	if _, err := r.command("workspace", "context", "project", "--event", "@Orient", "--json"); err != nil {
		return err
	}
	contract, err := r.command("contract", "show", contractRef, "--json")
	if err != nil {
		return err
	}
	contractFP := stringField(contract.Data, "fingerprint")

	proposalDecision, err := r.decision("0199c000-0000-7000-8000-000000000010", "proposal.create", proposalRef, "absent", nil, "approve")
	if err != nil {
		return err
	}
	proposal, err := r.inputCommand(map[string]any{
		"id": strings.TrimPrefix(proposalRef, "Proposal:"), "title": "Ship release-ready v2", "actor": "owner", "status": "accepted",
		"target_contract": contractRef, "base_version": "1", "base_fingerprint": contractFP, "new_capability": false,
		"additions": []string{"Release only checksum-verified native v2 artifacts."}, "modifications": []any{}, "removals": []any{},
		"rationale": "Prove the installed governed release loop.", "scope": []string{"v2"}, "gaps": []any{},
		"authorization": proposalDecision, "idempotency_key": "scenario-r-proposal", "freshness_valid_until": r.future(),
	}, "proposal", "create", "--input", "@input", "--json")
	if err != nil {
		return err
	}
	proposalFP := stringField(proposal.Data, "fingerprint")

	prepared, err := r.inputCommand(map[string]any{
		"proposal": map[string]any{"ref": proposalRef, "fingerprint": proposalFP}, "baseline": contractFP,
		"direction_sources": []any{map[string]any{"ref": contractRef, "fingerprint": contractFP}},
		"candidates": []any{map[string]any{
			"name": "release-smoke", "outcome": "Prove installed governed closure and cold resume.", "evidence": []string{"claim:release-smoke"},
			"cancellation_state": "Release artifact remains locally inspectable.", "reversibility": "Disposable workspace only.",
			"standalone_coherence": "Complete governed v2 flow.", "integration_path": "guided orientation to mechanical closure", "learning_value": "native runtime proof",
		}},
		"selected": "release-smoke", "design_sufficiency": "sufficient", "design_rationale": "Accepted v2 contracts fix the loop.",
		"slice_quality": "coherent", "slice_rationale": "One installed-runtime acceptance slice.",
		"completion_criteria": []map[string]string{{"claim": "claim:release-smoke", "pass_boundary": "release smoke passes", "proof_requirement": "real-process smoke", "review_level": "automatic"}}, "stop_conditions": []string{"authority-drift"},
		"evidence_claims": []string{"claim:release-smoke"}, "fresh_until": r.future(),
	}, "mission", "prepare", "--input", "@input", "--json")
	if err != nil {
		return err
	}
	missionDecision, err := r.decision("0199c000-0000-7000-8000-000000000011", "mission.create", missionRef, "absent", nil, "activate")
	if err != nil {
		return err
	}
	mission, err := r.inputCommand(map[string]any{
		"id": strings.TrimPrefix(missionRef, "Mission:"), "title": "Installed release smoke", "actor": "owner", "proposal": proposalRef,
		"outcome":        "Prove installed governed closure and cold resume.",
		"objectives":     []any{map[string]any{"id": strings.TrimPrefix(objectiveRef, "Objective:"), "outcome": "Complete release smoke", "dependencies": []any{}, "expected_proof": []string{"claim:release-smoke"}}},
		"initial_run_id": strings.TrimPrefix(runRef, "Run:"), "design_sufficiency": "sufficient", "slice_quality": "coherent",
		"dependencies": []any{}, "gaps": []any{}, "evidence_claims": []string{"claim:release-smoke"}, "scope": []string{"v2"},
		"allowed_actions": []string{"inspect", "test"}, "forbidden_effects": []string{"provider-mutation"}, "baseline": contractFP,
		"budget_units": 2, "repair_budget": 1, "expires_at": r.future(), "stops": []string{"authority-drift"},
		"recovery_point": "local-release-artifact", "return_destination": "Mission owner", "authorization": missionDecision,
		"expected_proposal_fingerprint": proposalFP, "idempotency_key": "scenario-r-mission", "preparation": prepared.Data,
	}, "mission", "create", "--input", "@input", "--json")
	if err != nil {
		return err
	}
	missionFP := stringField(mission.Data, "fingerprint")

	activateDecision, err := r.decision("0199c000-0000-7000-8000-000000000012", "mission.transition.active", missionRef, missionFP, nil, "activate")
	if err != nil {
		return err
	}
	active, err := r.command("mission", "transition", missionRef, "--to", "active", "--authorization", activateDecision, "--expected-fingerprint", missionFP, "--idempotency-key", "scenario-r-active", "--json")
	if err != nil {
		return err
	}
	activeFP := stringField(active.Data, "fingerprint")

	handoffDecision, err := r.decision("0199c000-0000-7000-8000-000000000013", "handoff.create", handoffRef, "absent", nil, "dispatch")
	if err != nil {
		return err
	}
	handoff, err := r.inputCommand(map[string]any{
		"id": strings.TrimPrefix(handoffRef, "Handoff:"), "title": "Disposable installed-runtime check", "mission": missionRef,
		"objective": objectiveRef, "run": runRef, "sender": "owner", "actor": "executor", "destination": "replacement-runtime",
		"host_pointer": "local:disposable", "scope": []string{"v2"}, "inputs": []string{contractRef + "@" + contractFP},
		"allowed_actions": []string{"test"}, "forbidden_effects": []string{"provider-mutation"}, "evidence_claims": []string{"claim:release-smoke"},
		"budget_units": 1, "expires_at": r.future(), "stops": []string{"authority-drift"}, "recovery_point": "local-release-artifact",
		"return_destination": "Mission owner", "return_contract": []string{"result", "evidence", "recovery point", "one next action or owner gate"}, "authorization": handoffDecision, "expected_mission_fingerprint": activeFP,
		"idempotency_key": "scenario-r-handoff",
	}, "handoff", "create", "--input", "@input", "--json")
	if err != nil {
		return err
	}
	_, err = r.inputCommand(map[string]any{
		"id": strings.TrimPrefix(returnRef, "Handoff:"), "title": "Installed runtime returned", "dispatch": handoffRef, "status": "succeeded",
		"actor": "replacement-runtime", "final_baseline": contractFP, "result": "native v2 binary completed the bounded check",
		"actions": []string{"test"}, "provider_receipts": []any{}, "evidence": []string{evidenceRef}, "remaining_gaps": []any{},
		"budget_used": 1, "recovery_point": "local-release-artifact", "next_action": "record evidence", "owner_gate": "",
		"expected_dispatch_fingerprint": stringField(handoff.Data, "fingerprint"), "idempotency_key": "scenario-r-handoff-return",
	}, "handoff", "return", "--input", "@input", "--json")
	if err != nil {
		return err
	}

	evidenceDecision, err := r.decision("0199c000-0000-7000-8000-000000000014", "evidence.create", evidenceRef, "absent", nil, "record")
	if err != nil {
		return err
	}
	_, err = r.inputCommand(map[string]any{
		"id": strings.TrimPrefix(evidenceRef, "Evidence:"), "title": "Installed runtime proof", "mission": missionRef, "objective": objectiveRef,
		"claim": "claim:release-smoke", "classification": "direct", "scope": []string{"v2"}, "method": "installed native binary scenario",
		"actor": "executor", "target": proposalRef, "environment": "disposable local prefix", "observed_at": time.Now().UTC().Format(time.RFC3339),
		"freshness_valid_until": r.future(), "limitations": []string{"local unpublished artifacts"}, "contrary_evidence": []any{},
		"required_checks": []string{"installed-binary-flow"}, "check_results": []string{"pass:installed-binary-flow"},
		"review_state": "independent-accepted", "executor_authored": true, "authorization": evidenceDecision, "idempotency_key": "scenario-r-evidence",
	}, "evidence", "create", "--input", "@input", "--json")
	if err != nil {
		return err
	}

	awaitDecision, err := r.decision("0199c000-0000-7000-8000-000000000015", "mission.transition.awaiting-assessment", missionRef, activeFP, nil, "await")
	if err != nil {
		return err
	}
	awaiting, err := r.command("mission", "transition", missionRef, "--to", "awaiting-assessment", "--authorization", awaitDecision,
		"--expected-fingerprint", activeFP, "--idempotency-key", "scenario-r-await", "--json")
	if err != nil {
		return err
	}

	assessmentDecision, err := r.decision("0199c000-0000-7000-8000-000000000016", "assessment.record", assessmentRef, "absent", nil, "record")
	if err != nil {
		return err
	}
	_, err = r.inputCommand(map[string]any{
		"id": strings.TrimPrefix(assessmentRef, "Assessment:"), "title": "Installed release assessment", "mission": missionRef,
		"verdict": "ready-for-owner", "actor": "reviewer", "claims": []string{"claim:release-smoke"}, "evidence": []string{evidenceRef},
		"blocking_findings": []any{}, "limitations": []string{"local unpublished artifacts"}, "repair_attempts": []any{},
		"recovery_point": "local-release-artifact", "authorization": assessmentDecision, "idempotency_key": "scenario-r-assessment",
	}, "assessment", "record", "--input", "@input", "--json")
	if err != nil {
		return err
	}

	reconcileDecision, err := r.decision("0199c000-0000-7000-8000-000000000017", "contract.reconcile", contractRef, contractFP, []string{assessmentRef}, "accept")
	if err != nil {
		return err
	}
	reconciled, err := r.command("contract", "reconcile", contractRef, "--proposal", proposalRef, "--authorization", reconcileDecision,
		"--expected-fingerprint", contractFP, "--idempotency-key", "scenario-r-reconcile", "--json")
	if err != nil {
		return err
	}
	receipt := stringField(reconciled.Data, "receipt")
	resolvedAction := "inspect the cold-resume packet"
	resolveDecision, err := r.decision("0199c000-0000-7000-8000-000000000018", "mission.transition.resolved", missionRef,
		stringField(awaiting.Data, "fingerprint"), []string{assessmentRef, receipt}, "completed",
		"objective.satisfy:"+objectiveRef, "terminal-next-action:"+resolvedAction)
	if err != nil {
		return err
	}
	resolved, err := r.command("mission", "transition", missionRef, "--to", "resolved", "--authorization", resolveDecision,
		"--expected-fingerprint", stringField(awaiting.Data, "fingerprint"), "--idempotency-key", "scenario-r-resolve",
		"--assessment", assessmentRef, "--reconciliation", receipt, "--disposition", "completed",
		"--terminal-next-action", resolvedAction, "--satisfied-objectives", objectiveRef, "--json")
	if err != nil {
		return err
	}
	if _, err := r.command("workspace", "context", missionRef, "--event", "@Resume", "--json"); err != nil {
		return err
	}
	archiveDecision, err := r.decision("0199c000-0000-7000-8000-000000000019", "mission.archive", missionRef,
		stringField(resolved.Data, "fingerprint"), []string{assessmentRef, receipt}, "archive")
	if err != nil {
		return err
	}
	if _, err := r.command("mission", "archive", missionRef, "--authorization", archiveDecision,
		"--expected-fingerprint", stringField(resolved.Data, "fingerprint"), "--idempotency-key", "scenario-r-archive",
		"--terminal-packet", missionRef, "--json"); err != nil {
		return err
	}
	_, err = r.command("workspace", "context", "project", "--event", "@Orient", "--json")
	return err
}

var v2Version = regexp.MustCompile(`^spectacular 2\.[0-9]+\.[0-9]+(?:-[0-9A-Za-z.-]+)?$`)

func validV2Version(value string) bool {
	return v2Version.MatchString(value)
}

func (r *runner) future() string { return time.Now().UTC().Add(24 * time.Hour).Format(time.RFC3339) }

func (r *runner) decision(id, operation, target, expected string, evidence []string, disposition string, effects ...string) (string, error) {
	allEffects := append([]string{operation}, effects...)
	result, err := r.inputCommand(map[string]any{
		"id": id, "title": operation, "actor": "owner", "actor_role": "owner", "authority_basis": "approved Scenario R charter",
		"question": operation, "scope": []string{"v2"}, "disposition": disposition, "rationale": "Explicit local smoke authorization.",
		"alternatives": []any{}, "targets": []string{target}, "expected_fingerprints": []string{expected}, "operation": operation,
		"authorized_effects": allEffects, "conditions": []string{"no-provider-effects"}, "expires_at": r.future(), "evidence": evidence,
		"idempotency_key": "scenario-r-decision-" + id,
	}, "decision", "create", "--input", "@input", "--json")
	if err != nil {
		return "", err
	}
	return stringField(result.Data, "ref"), nil
}

func (r *runner) inputCommand(input map[string]any, args ...string) (envelope, error) {
	r.step++
	path := filepath.Join(r.inputs, fmt.Sprintf("%02d.json", r.step))
	data, err := json.MarshalIndent(input, "", "  ")
	if err != nil {
		return envelope{}, err
	}
	if err := os.WriteFile(path, append(data, '\n'), 0o600); err != nil {
		return envelope{}, err
	}
	for index := range args {
		if args[index] == "@input" {
			args[index] = path
		}
	}
	return r.command(args...)
}

func (r *runner) command(args ...string) (envelope, error) {
	cmd := exec.Command(r.binary, args...)
	cmd.Dir = r.workspace
	output, err := cmd.CombinedOutput()
	if err != nil {
		return envelope{}, fmt.Errorf("spectacular %s: %w\n%s", strings.Join(args, " "), err, output)
	}
	var parsed envelope
	if err := json.Unmarshal(output, &parsed); err != nil {
		return envelope{}, fmt.Errorf("spectacular %s returned invalid JSON: %w\n%s", strings.Join(args, " "), err, output)
	}
	if parsed.Schema == "" || parsed.Data == nil {
		return envelope{}, fmt.Errorf("spectacular %s returned an incomplete envelope", strings.Join(args, " "))
	}
	fmt.Printf("command=spectacular %s\nschema=%s\n", strings.Join(args, " "), parsed.Schema)
	return parsed, nil
}

func stringField(value map[string]any, path ...string) string {
	var current any = value
	for _, key := range path {
		object, ok := current.(map[string]any)
		if !ok {
			fatal(fmt.Errorf("field %s is not an object", strings.Join(path, ".")))
		}
		current = object[key]
	}
	text, ok := current.(string)
	if !ok || text == "" {
		fatal(fmt.Errorf("field %s is not a non-empty string", strings.Join(path, ".")))
	}
	return text
}

func copyFixture(source, destination string) error {
	entries, err := os.ReadDir(destination)
	if err != nil {
		return err
	}
	if len(entries) != 0 {
		return errors.New("smoke workspace destination must be empty")
	}
	return filepath.WalkDir(source, func(path string, entry fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		relative, err := filepath.Rel(source, path)
		if err != nil || relative == "." {
			return err
		}
		target := filepath.Join(destination, relative)
		if entry.IsDir() {
			return os.MkdirAll(target, 0o755)
		}
		if !entry.Type().IsRegular() {
			return fmt.Errorf("fixture contains non-regular entry: %s", path)
		}
		data, err := os.ReadFile(path)
		if err != nil {
			return err
		}
		return os.WriteFile(target, data, 0o644)
	})
}

func fatal(err error) {
	fmt.Fprintln(os.Stderr, "smoke refused:", err)
	os.Exit(3)
}
