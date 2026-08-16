package command

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"io/fs"
	"os"
	"os/exec"
	"path/filepath"
	"reflect"
	"strings"
	"testing"
	"time"

	"github.com/alexsmedile/spectacular/v2/internal/discovery"
	"github.com/alexsmedile/spectacular/v2/internal/domain"
	"github.com/alexsmedile/spectacular/v2/internal/governance"
	"github.com/alexsmedile/spectacular/v2/internal/humanlayout"
	"github.com/alexsmedile/spectacular/v2/internal/workspace"
)

const (
	missionID             = "0198a1a0-0000-7000-8000-000000000002"
	gapID                 = "0198a1a0-0000-7000-8000-000000000004"
	runID                 = "0198a1a0-0000-7000-8000-000000000005"
	checkpointID          = "0198a1a0-0000-7000-8000-000000000006"
	evidenceID            = "0198a1a0-0000-7000-8000-000000000007"
	decisionID            = "0198a1a0-0000-7000-8000-000000000008"
	fixtureProjectPath    = ".spectacular/PROJECT.md"
	fixtureProposalPath   = ".spectacular/proposals/P1-prove-cold-recovery.md"
	fixtureMissionPath    = ".spectacular/missions/M1-scenario-a-cold-recovery/MISSION.md"
	fixtureRunPath        = ".spectacular/missions/M1-scenario-a-cold-recovery/runs/R1-primary-recovery-run/RUN.md"
	fixtureCheckpointPath = ".spectacular/missions/M1-scenario-a-cold-recovery/runs/R1-primary-recovery-run/checkpoints/C1-implementation-ready.md"
	fixtureEvidencePath   = ".spectacular/missions/M1-scenario-a-cold-recovery/evidence/E1-6r35pl.md"
	fixtureDecisionPath   = ".spectacular/missions/M1-scenario-a-cold-recovery/decisions/D1-t77vmk.md"
	fixtureGapPath        = ".spectacular/missions/M1-scenario-a-cold-recovery/gaps/G1-uiaahv.md"
)

func fixture(t *testing.T) string {
	t.Helper()
	absolute, err := filepath.Abs(filepath.Join("..", "..", "testdata", "scenario-a"))
	if err != nil {
		t.Fatal(err)
	}
	return absolute
}
func fixedNow() time.Time { return time.Date(2026, 8, 10, 10, 0, 0, 0, time.UTC) }

func TestWorkspaceContextDefaultsToCompactPointersAndKeepsFullJSON(t *testing.T) {
	args := []string{"workspace", "context", "Mission:" + missionID, "--event", "@Orient"}
	var compact, stderr bytes.Buffer
	if exit := (Runner{Cwd: fixture(t), Stdout: &compact, Stderr: &stderr, Now: fixedNow}).Run(args); exit != 0 {
		t.Fatalf("compact context failed: %s", stderr.String())
	}
	for _, want := range []string{
		"Spectacular context — Mission:" + missionID,
		"Authority: mission-anchor Mission:" + missionID + " — " + fixtureMissionPath,
		"spectacular mission show Mission:" + missionID,
		"Gap:",
		"Full: rerun this command with --json.",
	} {
		if !strings.Contains(compact.String(), want) {
			t.Fatalf("compact context missing %q: %s", want, compact.String())
		}
	}
	if strings.Contains(compact.String(), `"schema_version"`) {
		t.Fatalf("default context is still a full JSON dump: %s", compact.String())
	}

	full := runJSON(t, fixture(t), append(args, "--json"))
	for _, want := range []string{`"schema_version":"spectacular.context.v1"`, `"path":"` + fixtureMissionPath + `"`, `"fingerprint":"`} {
		if !strings.Contains(full, want) {
			t.Fatalf("full context lost %q: %s", want, full)
		}
	}
}

func TestSelfHostedWorkspaceIsHumanOperable(t *testing.T) {
	root, err := filepath.Abs(filepath.Join("..", ".."))
	if err != nil {
		t.Fatal(err)
	}
	before := treeDigest(t, filepath.Join(root, ".spectacular"))
	text := runJSON(t, root, []string{"workspace", "context", "project", "--event", "@Orient", "--json"})
	for _, required := range []string{
		`"ref":"Contract:019fe381-5d61-7223-b362-03a5f99a7b10"`,
		`"path":".spectacular/PRODUCT.md"`,
		`"path":".spectacular/ARCHITECTURE.md"`,
		`"path":".spectacular/STACK.md"`,
		`"path":".spectacular/archive/missions/M1-human-operability/MISSION.md"`,
		`"kind":"continuation"`,
		`"operation":"publish v2.0.0-rc.2, then begin the v2.1.0 governed-autonomy Mission"`,
	} {
		if !strings.Contains(text, required) {
			t.Fatalf("self-hosted RC workspace missing %s: %s", required, text)
		}
	}
	if strings.Contains(text, "Scenario A") {
		t.Fatalf("self-hosted workspace still exposes fixture state: %s", text)
	}
	opened, err := discovery.Open(root)
	if err != nil {
		t.Fatal(err)
	}
	mission, err := opened.Lookup("Mission:019fe381-5d61-7223-b362-03a5f99a7b02", domain.Mission)
	if err != nil {
		t.Fatal(err)
	}
	if mission.Document.Record.Status == nil || *mission.Document.Record.Status != "resolved" {
		t.Fatalf("self-hosted Mission is not resolved: %#v", mission.Document.Record.Status)
	}
	archived, err := workspace.Bool(mission.Document, "archived", true)
	if err != nil || !archived {
		t.Fatalf("self-hosted Mission is not archived: archived=%v err=%v", archived, err)
	}
	byHuman, err := opened.Lookup("M1/R1/C1", domain.Checkpoint)
	if err != nil {
		t.Fatal(err)
	}
	if byHuman.Path != ".spectacular/archive/missions/M1-human-operability/runs/R1-implement-layout/checkpoints/C1-layout-in-progress.md" {
		t.Fatalf("human reference resolved %s", byHuman.Path)
	}
	gap, err := opened.Lookup("Gap:019fe381-5d61-7223-b362-03a5f99a7b04", domain.Gap)
	if err != nil {
		t.Fatal(err)
	}
	if gap.Document.Record.Title == nil || *gap.Document.Record.Title != "Owner disposition on the human-operable RC.2 candidate" {
		t.Fatalf("self-hosted Gap is not the review boundary: %#v", gap.Document.Record.Title)
	}
	if strings.Contains(text, `"kind":"owner-gate"`) {
		t.Fatalf("resolved self-hosted Mission still exposes an owner gate: %s", text)
	}
	if after := treeDigest(t, filepath.Join(root, ".spectacular")); after != before {
		t.Fatal("self-hosted workspace read mutated canonical state")
	}
}

func TestSelfHostedIndexesAreDeterministicProjections(t *testing.T) {
	root, err := filepath.Abs(filepath.Join("..", ".."))
	if err != nil {
		t.Fatal(err)
	}
	opened, err := discovery.Open(root)
	if err != nil {
		t.Fatal(err)
	}
	generated, err := humanlayout.Indexes(opened.Entries, nil, nil)
	if err != nil {
		t.Fatal(err)
	}
	for path, expected := range generated {
		actual, readErr := os.ReadFile(filepath.Join(root, filepath.FromSlash(path)))
		if readErr != nil {
			t.Fatalf("read generated index %s: %v", path, readErr)
		}
		if !bytes.Equal(actual, expected) {
			t.Fatalf("generated index drifted: %s\n--- actual ---\n%s\n--- expected ---\n%s", path, actual, expected)
		}
	}
}

func TestSelfHostedHumanPointersExecuteWithoutMutation(t *testing.T) {
	root, err := filepath.Abs(filepath.Join("..", ".."))
	if err != nil {
		t.Fatal(err)
	}
	before := treeDigest(t, filepath.Join(root, ".spectacular"))
	commands := map[string]bool{}
	for _, seed := range [][]string{{"anchor", "show", "project", "--json"}, {"mission", "show", "M1", "--json"}} {
		var out bytes.Buffer
		exit := (Runner{Cwd: root, Stdout: &out, Stderr: &bytes.Buffer{}, Now: fixedNow}).Run(seed)
		if exit != 0 {
			t.Fatalf("seed %v exit=%d output=%s", seed, exit, out.String())
		}
		var payload any
		if err := json.Unmarshal(out.Bytes(), &payload); err != nil {
			t.Fatal(err)
		}
		collectShowCommands(payload, commands)
	}
	for commandLine := range commands {
		args := strings.Fields(strings.TrimPrefix(commandLine, "spectacular "))
		args = append(args, "--json")
		var out bytes.Buffer
		if exit := (Runner{Cwd: root, Stdout: &out, Stderr: &bytes.Buffer{}, Now: fixedNow}).Run(args); exit != 0 {
			t.Fatalf("emitted pointer failed: %s exit=%d output=%s", commandLine, exit, out.String())
		}
	}
	if after := treeDigest(t, filepath.Join(root, ".spectacular")); after != before {
		t.Fatal("human pointer traversal mutated canonical state")
	}
}

func TestPublicRegistryAndCommands(t *testing.T) {
	if len(Registry) != 34 {
		t.Fatalf("registry has %d commands, want 34", len(Registry))
	}
	seen := map[Operation]bool{}
	for _, spec := range Registry {
		if spec.Operation == 0 || seen[spec.Operation] {
			t.Fatalf("unbound or duplicate registry operation %d", spec.Operation)
		}
		seen[spec.Operation] = true
	}
	if len(seen) != int(opMissionAutopilot) {
		t.Fatalf("registry dispatch parity: %d bound operations, want %d", len(seen), opMissionAutopilot)
	}
	if VersionInspection.Command != "spectacular --version" || VersionInspection.Schema != "spectacular.build-info.v1" || VersionInspection.Effect != ReadOnly {
		t.Fatalf("release inspection metadata drifted: %#v", VersionInspection)
	}
	tests := []struct {
		args     []string
		schema   string
		contains string
	}{
		{[]string{"anchor", "show", "project", "--json"}, "spectacular.anchor.show.v1", `"authoritative":{"identity":{"noun":"anchor","ref":"project"`},
		{[]string{"mission", "list", "--json"}, "spectacular.mission.list.v1", `"items":[`},
		{[]string{"mission", "show", "Mission:" + missionID, "--json"}, "spectacular.mission.show.v1", `"path":".spectacular/proposals/P1-prove-cold-recovery.md"`},
		{[]string{"gap", "list", "--scope", missionID, "--json"}, "spectacular.gap.list.v1", `"noun":"Gap"`},
		{[]string{"gap", "show", "Gap:" + gapID, "--json"}, "spectacular.gap.show.v1", `"show_command":"spectacular mission show`},
		{[]string{"run", "show", fixtureRunPath, "--json"}, "spectacular.run.show.v1", `"show_command":"spectacular checkpoint show`},
		{[]string{"checkpoint", "show", checkpointID, "--json"}, "spectacular.checkpoint.show.v1", `"show_command":"spectacular evidence show`},
		{[]string{"evidence", "show", "Evidence:" + evidenceID, "--json"}, "spectacular.evidence.show.v1", `"show_command":"spectacular checkpoint show`},
		{[]string{"decision", "show", decisionID, "--json"}, "spectacular.decision.show.v1", `"show_command":"spectacular run show`},
		{[]string{"workspace", "validate", "project", "--json"}, "spectacular.workspace.validate.v1", `"valid":true`},
	}
	for _, test := range tests {
		t.Run(strings.Join(test.args, "_"), func(t *testing.T) {
			var out, errOut bytes.Buffer
			exit := Runner{Cwd: fixture(t), Stdout: &out, Stderr: &errOut, Now: fixedNow}.Run(test.args)
			if exit != 0 {
				t.Fatalf("exit=%d stderr=%s stdout=%s", exit, errOut.String(), out.String())
			}
			if !strings.Contains(out.String(), `"schema_version":"`+test.schema+`"`) || !strings.Contains(out.String(), test.contains) {
				t.Fatalf("unexpected JSON: %s", out.String())
			}
			if errOut.Len() != 0 {
				t.Fatalf("stderr=%s", errOut.String())
			}
		})
	}
}

// TestPointerShowCommandsExecuteInFreshProcess enforces the public traversal
// promise: a projection may emit a show_command only when that exact command
// can be followed in a new process without changing the workspace.
func TestPointerShowCommandsExecuteInFreshProcess(t *testing.T) {
	if os.Getenv("SPECTACULAR_POINTER_HELPER") == "1" {
		args := os.Args
		separator := -1
		for i, arg := range args {
			if arg == "--" {
				separator = i
				break
			}
		}
		if separator < 0 {
			os.Exit(2)
		}
		os.Exit(Runner{Cwd: os.Getenv("SPECTACULAR_POINTER_WORKSPACE"), Stdout: os.Stdout, Stderr: os.Stderr, Now: fixedNow}.Run(args[separator+1:]))
	}
	root := fixture(t)
	before := treeDigest(t, root)
	seed := [][]string{
		{"anchor", "show", "project", "--json"}, {"mission", "list", "--json"},
		{"mission", "show", "Mission:" + missionID, "--json"}, {"gap", "list", "--scope", "Mission:" + missionID, "--json"},
		{"proposal", "show", "Proposal:0198a1a0-0000-7000-8000-000000000001", "--json"},
	}
	commands := map[string]bool{}
	for _, args := range seed {
		var out, errOut bytes.Buffer
		if exit := (Runner{Cwd: root, Stdout: &out, Stderr: &errOut, Now: fixedNow}).Run(args); exit != 0 {
			t.Fatalf("seed %v: %s", args, errOut.String())
		}
		var envelope map[string]any
		if err := json.Unmarshal(out.Bytes(), &envelope); err != nil {
			t.Fatal(err)
		}
		collectShowCommands(envelope, commands)
	}
	if len(commands) == 0 {
		t.Fatal("official fixture emitted no show commands")
	}
	for command := range commands {
		fields := strings.Fields(command)
		if len(fields) < 3 || fields[0] != "spectacular" {
			t.Fatalf("bad show command %q", command)
		}
		child := exec.Command(os.Args[0], "-test.run=^TestPointerShowCommandsExecuteInFreshProcess$", "--", fields[1], fields[2])
		if len(fields) > 3 {
			child.Args = append(child.Args, fields[3:]...)
		}
		child.Env = append(os.Environ(), "SPECTACULAR_POINTER_HELPER=1", "SPECTACULAR_POINTER_WORKSPACE="+root)
		if output, err := child.CombinedOutput(); err != nil {
			t.Fatalf("%q failed: %v: %s", command, err, output)
		}
	}
	if after := treeDigest(t, root); after != before {
		t.Fatalf("pointer traversal mutated workspace: before=%s after=%s", before, after)
	}
}

func collectShowCommands(value any, commands map[string]bool) {
	switch item := value.(type) {
	case map[string]any:
		if command, ok := item["show_command"].(string); ok {
			commands[command] = true
		}
		for _, child := range item {
			collectShowCommands(child, commands)
		}
	case []any:
		for _, child := range item {
			collectShowCommands(child, commands)
		}
	}
}

func treeDigest(t *testing.T, root string) string {
	t.Helper()
	hash := sha256.New()
	err := filepath.WalkDir(root, func(path string, entry fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if entry.IsDir() {
			return nil
		}
		data, err := os.ReadFile(path)
		if err != nil {
			return err
		}
		hash.Write([]byte(strings.TrimPrefix(path, root)))
		hash.Write(data)
		return nil
	})
	if err != nil {
		t.Fatal(err)
	}
	return hex.EncodeToString(hash.Sum(nil))
}

func TestGeneratedCatalogMatchesRegistryAndKeepsJudgmentOut(t *testing.T) {
	versionBytes, err := os.ReadFile(filepath.Join("..", "..", "VERSION"))
	if err != nil {
		t.Fatal(err)
	}
	version := strings.TrimSpace(string(versionBytes))
	var jsonOut, markdown bytes.Buffer
	if err := WriteCatalogJSONVersion(&jsonOut, version); err != nil {
		t.Fatal(err)
	}
	if err := WriteCatalogMarkdownVersion(&markdown, version); err != nil {
		t.Fatal(err)
	}
	for _, want := range []string{"spectacular --version", "spectacular workspace context", "spectacular mission prepare", "spectacular mission autopilot"} {
		if !strings.Contains(jsonOut.String(), want) || !strings.Contains(markdown.String(), want) {
			t.Fatalf("generated catalogs omit %s", want)
		}
	}
	for _, forbidden := range []string{"spectacular status", "spectacular inspect", "spectacular record", "spectacular decide", "spectacular resolve"} {
		if strings.Contains(jsonOut.String(), `"command": "`+forbidden+`"`) {
			t.Fatalf("generated catalog contains forbidden generic command %s", forbidden)
		}
	}
	for path, generated := range map[string]string{
		filepath.Join("..", "..", "skills", "spectacular", "generated", "mechanical-interface.json"): jsonOut.String(),
		filepath.Join("..", "..", "skills", "spectacular", "generated", "mechanical-interface.md"):   markdown.String(),
	} {
		checkedIn, err := os.ReadFile(path)
		if err != nil {
			t.Fatal(err)
		}
		if string(checkedIn) != generated {
			t.Fatalf("generated interface is stale: %s", path)
		}
	}
}

func TestJSONIsDeterministicAndExactLookupFormsAgree(t *testing.T) {
	forms := []string{"M1", missionID, "Mission:" + missionID, fixtureMissionPath}
	var canonical map[string]any
	for _, ref := range forms {
		var out bytes.Buffer
		exit := Runner{Cwd: fixture(t), Stdout: &out, Stderr: &bytes.Buffer{}, Now: fixedNow}.Run([]string{"mission", "show", ref, "--json"})
		if exit != 0 {
			t.Fatalf("%s exit %d", ref, exit)
		}
		var got map[string]any
		if err := json.Unmarshal(out.Bytes(), &got); err != nil {
			t.Fatal(err)
		}
		if canonical == nil {
			canonical = got
		} else if !reflect.DeepEqual(canonical, got) {
			t.Fatalf("lookup form changed projection for %s", ref)
		}
	}
	var a, b bytes.Buffer
	r := Runner{Cwd: fixture(t), Stdout: &a, Stderr: &bytes.Buffer{}, Now: fixedNow}
	r.Run([]string{"anchor", "show", "project", "--json"})
	r.Stdout = &b
	r.Run([]string{"anchor", "show", "project", "--json"})
	if a.String() != b.String() {
		t.Fatalf("same inputs and clock produced different JSON")
	}
}

func TestUsageAndRefusalChannelsAndCodes(t *testing.T) {
	var out, errOut bytes.Buffer
	if got := (Runner{Cwd: fixture(t), Stdout: &out, Stderr: &errOut, Now: fixedNow}).Run([]string{"status", "--json"}); got != 2 {
		t.Fatalf("generic status exit=%d", got)
	}
	if !strings.Contains(out.String(), `"code":"usage"`) || errOut.Len() != 0 {
		t.Fatalf("JSON usage channels stdout=%q stderr=%q", out.String(), errOut.String())
	}
	if !strings.Contains(out.String(), `"invoked_command":"spectacular status --json"`) || !strings.Contains(out.String(), `"exit_status":2`) {
		t.Fatalf("usage refusal lacks invocation metadata: %s", out.String())
	}
	out.Reset()
	if got := (Runner{Cwd: fixture(t), Stdout: &out, Stderr: &bytes.Buffer{}, Now: fixedNow}).Run([]string{"mission", "list", "--json", "--json"}); got != 2 || !strings.Contains(out.String(), `"code":"usage"`) {
		t.Fatalf("duplicate --json exit=%d output=%s", got, out.String())
	}
	out.Reset()
	errOut.Reset()
	if got := (Runner{Cwd: fixture(t), Stdout: &out, Stderr: &errOut, Now: fixedNow}).Run([]string{"mission", "show", "Mission:" + gapID}); got != 3 {
		t.Fatalf("noun mismatch exit=%d", got)
	}
	if out.Len() != 0 || !strings.Contains(errOut.String(), "noun_mismatch") {
		t.Fatalf("human refusal channels stdout=%q stderr=%q", out.String(), errOut.String())
	}
}

func TestKnownCommandArgumentsAreValidatedBeforeWorkspaceDiscovery(t *testing.T) {
	nonWorkspace := t.TempDir()
	tests := [][]string{
		{"mission", "show", "--json"},
		{"anchor", "show", "project", "extra", "--json"},
		{"gap", "list", "Mission:" + missionID, "--json"},
	}
	for _, args := range tests {
		t.Run(strings.Join(args, "_"), func(t *testing.T) {
			var out, errOut bytes.Buffer
			exit := (Runner{Cwd: nonWorkspace, Stdout: &out, Stderr: &errOut, Now: fixedNow}).Run(args)
			if exit != 2 || !strings.Contains(out.String(), `"code":"usage"`) || strings.Contains(out.String(), `workspace_not_found`) {
				t.Fatalf("exit=%d stdout=%s stderr=%s", exit, out.String(), errOut.String())
			}
		})
	}
}

func TestOwnerGateAndStaleRefusal(t *testing.T) {
	root := copyFixture(t)
	decision := filepath.Join(root, filepath.FromSlash(fixtureDecisionPath))
	if err := os.Remove(decision); err != nil {
		t.Fatal(err)
	}
	var out bytes.Buffer
	if got := (Runner{Cwd: root, Stdout: &out, Stderr: &bytes.Buffer{}, Now: fixedNow}).Run([]string{"mission", "show", missionID, "--json"}); got != 0 {
		t.Fatalf("missing authority must be owner gate success, exit=%d output=%s", got, out.String())
	}
	if !strings.Contains(out.String(), `"code":"authorize_continuation"`) || strings.Contains(out.String(), `"continuation"`) {
		t.Fatalf("missing exact gate: %s", out.String())
	}
	root = copyFixture(t)
	evidence := filepath.Join(root, filepath.FromSlash(fixtureEvidencePath))
	data, err := os.ReadFile(evidence)
	if err != nil {
		t.Fatal(err)
	}
	data = bytes.Replace(data, []byte("freshness_valid_until: \"2026-12-31T23:59:59Z\""), []byte("freshness_valid_until: \"2026-08-10T09:00:00Z\""), 1)
	if err := os.WriteFile(evidence, data, 0o644); err != nil {
		t.Fatal(err)
	}
	out.Reset()
	if got := (Runner{Cwd: root, Stdout: &out, Stderr: &bytes.Buffer{}, Now: fixedNow}).Run([]string{"anchor", "show", "project", "--json"}); got != 3 {
		t.Fatalf("stale evidence exit=%d output=%s", got, out.String())
	}
	if !strings.Contains(out.String(), `"code":"stale_required_source"`) {
		t.Fatalf("stale refusal=%s", out.String())
	}
}

func TestProjectAuthorityProjectionAndDecisionAuthorizationAreStrict(t *testing.T) {
	var out bytes.Buffer
	if got := (Runner{Cwd: fixture(t), Stdout: &out, Stderr: &bytes.Buffer{}, Now: fixedNow}).Run([]string{"anchor", "show", "project", "--json"}); got != 0 {
		t.Fatalf("project exit=%d output=%s", got, out.String())
	}
	text := out.String()
	for _, required := range []string{`"authoritative"`, `"direction"`, `"boundaries"`, `"constraints"`, `"current_truth"`, `"projection"`, `"continuation"`, `"show_command":"spectacular anchor show project"`, `"basis":"explicit-validity-window+source-fingerprint"`, `"expected_fingerprint"`, `"actual_fingerprint"`} {
		if !strings.Contains(text, required) {
			t.Fatalf("project view missing %s: %s", required, text)
		}
	}
	for _, mutation := range []struct{ name, old, new, code string }{{"missing decided_by", "decided_by: owner\n", "", "missing_required_field"}, {"wrong operation", "operation: resume", "operation: advance", "conflicting_authority"}} {
		t.Run(mutation.name, func(t *testing.T) {
			root := copyFixture(t)
			path := filepath.Join(root, filepath.FromSlash(fixtureDecisionPath))
			data, err := os.ReadFile(path)
			if err != nil {
				t.Fatal(err)
			}
			data = bytes.Replace(data, []byte(mutation.old), []byte(mutation.new), 1)
			if err := os.WriteFile(path, data, 0o644); err != nil {
				t.Fatal(err)
			}
			var got bytes.Buffer
			exit := (Runner{Cwd: root, Stdout: &got, Stderr: &bytes.Buffer{}, Now: fixedNow}).Run([]string{"mission", "show", missionID, "--json"})
			if exit != 3 || !strings.Contains(got.String(), `"code":"`+mutation.code+`"`) {
				t.Fatalf("exit=%d refusal=%s", exit, got.String())
			}
		})
	}
}

func TestSupersededResumeDecisionKeepsHistoricalMissionFingerprint(t *testing.T) {
	root := copyFixture(t)
	opened, err := discovery.Open(root)
	if err != nil {
		t.Fatal(err)
	}
	decision, err := opened.Lookup("Decision:"+decisionID, domain.Decision)
	if err != nil {
		t.Fatal(err)
	}
	status := "superseded"
	decision.Document.Record.Status = &status
	writeCanonicalEntry(t, decision)

	opened, err = discovery.Open(root)
	if err != nil {
		t.Fatal(err)
	}
	mission, err := opened.Lookup("Mission:"+missionID, domain.Mission)
	if err != nil {
		t.Fatal(err)
	}
	workspace.SetString(mission.Document, "outcome", "A later terminal Mission revision preserves historical authority.")
	writeCanonicalEntry(t, mission)

	var out bytes.Buffer
	exit := (Runner{Cwd: root, Stdout: &out, Stderr: &bytes.Buffer{}, Now: fixedNow}).Run([]string{"decision", "show", decisionID, "--json"})
	if exit != 0 || !strings.Contains(out.String(), `"state":"superseded"`) {
		t.Fatalf("exit=%d output=%s", exit, out.String())
	}
}

func TestProposalCurrentTruthAndMissionSourceHaveHonestDrillDown(t *testing.T) {
	root := copyFixture(t)
	project := filepath.Join(root, filepath.FromSlash(fixtureProjectPath))
	replaceInFile(t, project, "current_truth:\n  - Mission:"+missionID, "current_truth:\n  - Mission:"+missionID+"\n  - Proposal:0198a1a0-0000-7000-8000-000000000001")

	for _, args := range [][]string{{"anchor", "show", "project", "--json"}, {"mission", "show", missionID, "--json"}} {
		var out bytes.Buffer
		exit := (Runner{Cwd: root, Stdout: &out, Stderr: &bytes.Buffer{}, Now: fixedNow}).Run(args)
		if exit != 0 {
			t.Fatalf("%v exit=%d output=%s", args, exit, out.String())
		}
		text := out.String()
		for _, required := range []string{`"noun":"proposal"`, `"ref":"Proposal:0198a1a0-0000-7000-8000-000000000001"`, `"human_ref":"P1"`, `"path":".spectacular/proposals/P1-prove-cold-recovery.md"`, `"fingerprint":"`, `"show_command":"spectacular proposal show P1"`} {
			if !strings.Contains(text, required) {
				t.Fatalf("Proposal drill-down missing %s: %s", required, text)
			}
		}
	}
}

func TestAllReadPathsLeaveWorkspaceBytesModesAndMtimesUntouched(t *testing.T) {
	root := copyFixture(t)
	before := snapshot(t, root)
	commands := [][]string{{"anchor", "show", "project", "--json"}, {"mission", "list"}, {"gap", "list", "--scope", missionID}, {"workspace", "validate", "project", "--json"}, {"mission", "show", "not-a-ref", "--json"}}
	for _, args := range commands {
		(Runner{Cwd: root, Stdout: &bytes.Buffer{}, Stderr: &bytes.Buffer{}, Now: fixedNow}).Run(args)
	}
	after := snapshot(t, root)
	if !reflect.DeepEqual(before, after) {
		t.Fatalf("read commands mutated workspace\nbefore=%#v\nafter=%#v", before, after)
	}
}

func TestAdversarialWorkspaceRefusalsAreDeterministic(t *testing.T) {
	tests := []struct {
		name   string
		mutate func(*testing.T, string)
		args   []string
		code   string
	}{
		{
			name: "duplicate identity is ambiguous",
			mutate: func(t *testing.T, root string) {
				data, err := os.ReadFile(filepath.Join(root, filepath.FromSlash(fixtureGapPath)))
				if err != nil {
					t.Fatal(err)
				}
				if err := os.WriteFile(filepath.Join(root, ".spectacular", "missions", "M1-scenario-a-cold-recovery", "gaps", "G2-duplicate.md"), data, 0o644); err != nil {
					t.Fatal(err)
				}
			},
			args: []string{"anchor", "show", "project", "--json"},
			code: "duplicate_id",
		},
		{
			name: "multiple decisions conflict",
			mutate: func(t *testing.T, root string) {
				path := filepath.Join(root, filepath.FromSlash(fixtureDecisionPath))
				data, err := os.ReadFile(path)
				if err != nil {
					t.Fatal(err)
				}
				data = bytes.Replace(data, []byte(decisionID), []byte("0198a1a0-0000-7000-8000-000000000009"), 1)
				data = bytes.Replace(data, []byte("human_ref: M1/D1-t77vmk"), []byte("human_ref: M1/D2-c3xuq5"), 1)
				if err := os.WriteFile(filepath.Join(root, ".spectacular", "missions", "M1-scenario-a-cold-recovery", "decisions", "D2-conflict.md"), data, 0o644); err != nil {
					t.Fatal(err)
				}
			},
			args: []string{"mission", "show", missionID, "--json"},
			code: "conflicting_authority",
		},
		{
			name: "missing freshness source",
			mutate: func(t *testing.T, root string) {
				replaceInFile(t, filepath.Join(root, filepath.FromSlash(fixtureEvidencePath)), ".spectacular/workspace.yaml", ".spectacular/missing.yaml")
			},
			args: []string{"evidence", "show", evidenceID, "--json"},
			code: "record_not_found",
		},
		{
			name: "malformed freshness fingerprint",
			mutate: func(t *testing.T, root string) {
				path := filepath.Join(root, filepath.FromSlash(fixtureEvidencePath))
				data, err := os.ReadFile(path)
				if err != nil {
					t.Fatal(err)
				}
				start := bytes.Index(data, []byte("freshness_source_fingerprint: \""))
				if start < 0 {
					t.Fatal("freshness_source_fingerprint not found")
				}
				start += len("freshness_source_fingerprint: \"")
				data = append(append(append([]byte{}, data[:start]...), []byte("bad")...), data[start+64:]...)
				if err := os.WriteFile(path, data, 0o644); err != nil {
					t.Fatal(err)
				}
			},
			args: []string{"evidence", "show", evidenceID, "--json"},
			code: "invalid_continuation_fingerprint",
		},
		{
			name: "changed source is stale",
			mutate: func(t *testing.T, root string) {
				path := filepath.Join(root, ".spectacular", "workspace.yaml")
				data, err := os.ReadFile(path)
				if err != nil {
					t.Fatal(err)
				}
				if err := os.WriteFile(path, append(data, []byte("# changed basis\n")...), 0o644); err != nil {
					t.Fatal(err)
				}
			},
			args: []string{"evidence", "show", evidenceID, "--json"},
			code: "stale_required_source",
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			root := copyFixture(t)
			test.mutate(t, root)
			var first, second bytes.Buffer
			for _, output := range []*bytes.Buffer{&first, &second} {
				exit := (Runner{Cwd: root, Stdout: output, Stderr: &bytes.Buffer{}, Now: fixedNow}).Run(test.args)
				if exit != 3 || !strings.Contains(output.String(), `"code":"`+test.code+`"`) {
					t.Fatalf("exit=%d refusal=%s", exit, output.String())
				}
			}
			if first.String() != second.String() {
				t.Fatalf("refusal changed across identical reads:\nfirst=%s\nsecond=%s", first.String(), second.String())
			}
		})
	}
}

func TestOwnerGatePathIsByteNonmutating(t *testing.T) {
	root := copyFixture(t)
	if err := os.Remove(filepath.Join(root, filepath.FromSlash(fixtureDecisionPath))); err != nil {
		t.Fatal(err)
	}
	before := snapshot(t, root)
	var out bytes.Buffer
	exit := (Runner{Cwd: root, Stdout: &out, Stderr: &bytes.Buffer{}, Now: fixedNow}).Run([]string{"mission", "show", missionID, "--json"})
	if exit != 0 || !strings.Contains(out.String(), `"owner_gate":{"code":"authorize_continuation"`) {
		t.Fatalf("exit=%d output=%s", exit, out.String())
	}
	if after := snapshot(t, root); !reflect.DeepEqual(before, after) {
		t.Fatalf("owner-gate read mutated workspace\nbefore=%#v\nafter=%#v", before, after)
	}
}

func TestAnchorJSONGoldenDigest(t *testing.T) {
	var out bytes.Buffer
	exit := (Runner{Cwd: fixture(t), Stdout: &out, Stderr: &bytes.Buffer{}, Now: fixedNow}).Run([]string{"anchor", "show", "project", "--json"})
	if exit != 0 {
		t.Fatalf("exit=%d output=%s", exit, out.String())
	}
	digest := sha256.Sum256(out.Bytes())
	got := hex.EncodeToString(digest[:])
	const want = "0415f253aa46babd1bc382eedccdeeffd002a4a43ee938798a434027103f9c54"
	if got != want {
		t.Fatalf("anchor JSON golden digest=%s want=%s\n%s", got, want, out.String())
	}
}

func TestGovernedCreateUsesStrictInputAndRegisteredReadback(t *testing.T) {
	root := copyTree(t, filepath.Join("..", "..", "testdata", "scenario-b-c"))
	input := map[string]any{
		"id":                    "0199b000-0000-7000-8000-000000000090",
		"title":                 "Authorize bounded operation",
		"actor":                 "Alex",
		"actor_role":            "owner",
		"authority_basis":       "accepted B+C contract",
		"question":              "May the bounded operation proceed?",
		"scope":                 []string{"v2"},
		"disposition":           "approve",
		"rationale":             "Owner explicitly authorized it.",
		"alternatives":          []string{"defer"},
		"targets":               []string{"Proposal:0199b000-0000-7000-8000-000000000091"},
		"expected_fingerprints": []string{"absent"},
		"operation":             "proposal.create",
		"authorized_effects":    []string{"create Proposal"},
		"conditions":            []string{"no provider effects"},
		"expires_at":            "2026-08-11T10:00:00Z",
		"evidence":              []string{},
		"supersedes":            "",
		"idempotency_key":       "cli-decision-create-1",
	}
	data, err := json.Marshal(input)
	if err != nil {
		t.Fatal(err)
	}
	path := filepath.Join(t.TempDir(), "decision.json")
	if err := os.WriteFile(path, data, 0o644); err != nil {
		t.Fatal(err)
	}
	var out bytes.Buffer
	runner := Runner{Cwd: root, Stdout: &out, Stderr: &bytes.Buffer{}, Now: fixedNow}
	if exit := runner.Run([]string{"decision", "create", "--input", path, "--json"}); exit != 0 || !strings.Contains(out.String(), `"operation":"decision.create"`) || !strings.Contains(out.String(), `"path":".spectacular/decisions/D1-`) {
		t.Fatalf("create exit=%d output=%s", exit, out.String())
	}
	indexBytes, err := os.ReadFile(filepath.Join(root, ".spectacular", "index.md"))
	if err != nil || !strings.Contains(string(indexBytes), "Authorize bounded operation") || !strings.Contains(string(indexBytes), "non-authoritative") {
		t.Fatalf("generated human index missing created Decision: %s err=%v", indexBytes, err)
	}
	out.Reset()
	if exit := runner.Run([]string{"decision", "show", "Decision:0199b000-0000-7000-8000-000000000090", "--json"}); exit != 0 || !strings.Contains(out.String(), `"noun":"Decision"`) {
		t.Fatalf("show exit=%d output=%s", exit, out.String())
	}
	input["unexpected_authority"] = true
	data, _ = json.Marshal(input)
	if err := os.WriteFile(path, data, 0o644); err != nil {
		t.Fatal(err)
	}
	out.Reset()
	if exit := runner.Run([]string{"decision", "create", "--input", path, "--json"}); exit != 3 || !strings.Contains(out.String(), `"code":"invalid_known_field"`) {
		t.Fatalf("unknown input exit=%d output=%s", exit, out.String())
	}
}

func TestPublicReconcileSetAppliesTwoContractsAtomically(t *testing.T) {
	root := t.TempDir()
	meta := filepath.Join(root, ".spectacular")
	records := filepath.Join(meta, "records")
	if err := os.MkdirAll(records, 0o755); err != nil {
		t.Fatal(err)
	}
	manifest := []byte("schema_version: spectacular.workspace.v1\nrecord_roots:\n  - records\nproject_anchor: records/project.md\n")
	if err := os.WriteFile(filepath.Join(meta, "workspace.yaml"), manifest, 0o644); err != nil {
		t.Fatal(err)
	}
	markerHash := sha256.Sum256(manifest)
	markerFP := hex.EncodeToString(markerHash[:])
	anchor := testDocument(t, domain.Anchor, "0199c000-0000-7000-8000-000000000001", "Atomic set", "current")
	setTestFreshness(anchor, markerFP)
	workspace.SetString(anchor, "direction", "Apply an explicit atomic Contract set.")
	workspace.SetStrings(anchor, "boundaries", []string{"No provider effects."})
	workspace.SetStrings(anchor, "constraints", []string{"Owner Decision required."})
	workspace.SetStrings(anchor, "current_truth", []string{})
	writeTestDocument(t, filepath.Join(records, "project.md"), anchor)

	type item struct {
		contractID, proposalID, missionID, assessmentID, decisionID, addition string
	}
	items := []item{
		{"0199c000-0000-7000-8000-000000000010", "0199c000-0000-7000-8000-000000000011", "0199c000-0000-7000-8000-000000000012", "0199c000-0000-7000-8000-000000000013", "0199c000-0000-7000-8000-000000000014", "Atomic clause one."},
		{"0199c000-0000-7000-8000-000000000020", "0199c000-0000-7000-8000-000000000021", "0199c000-0000-7000-8000-000000000022", "0199c000-0000-7000-8000-000000000023", "0199c000-0000-7000-8000-000000000024", "Atomic clause two."},
	}
	var inputs []governance.ReconcileInput
	for _, item := range items {
		contractRef := "Contract:" + item.contractID
		proposalRef := "Proposal:" + item.proposalID
		missionRef := "Mission:" + item.missionID
		assessmentRef := "Assessment:" + item.assessmentID
		decisionRef := "Decision:" + item.decisionID
		contract := testDocument(t, domain.Contract, item.contractID, "Capability Contract", "current")
		setTestFreshness(contract, markerFP)
		workspace.SetString(contract, "contract_version", "1")
		workspace.SetString(contract, "purpose", "Atomic public reconciliation.")
		workspace.SetString(contract, "outcome", "Current truth remains complete.")
		workspace.SetStrings(contract, "applies_when", []string{"Governed delta."})
		workspace.SetStrings(contract, "does_not_apply_when", []string{"Read only."})
		workspace.SetStrings(contract, "does_not_provide", []string{"Provider authority."})
		workspace.SetStrings(contract, "required_behavior", []string{"Preserve truth."})
		workspace.SetStrings(contract, "operating_cases", []string{"Atomic set."})
		workspace.SetStrings(contract, "persistent_information", []string{"Provenance."})
		workspace.SetStrings(contract, "conformance_checks", []string{"Cold recovery."})
		workspace.SetString(contract, "authority_freshness", "Exact Decision.")
		workspace.SetStrings(contract, "related_material", []string{"B+C."})
		writeTestDocument(t, filepath.Join(records, "contract-"+item.contractID+".md"), contract)
		contractFP, err := workspace.Fingerprint(contract)
		if err != nil {
			t.Fatal(err)
		}
		proposal := testDocument(t, domain.Proposal, item.proposalID, "Accepted delta", "accepted")
		setTestFreshness(proposal, markerFP)
		workspace.SetString(proposal, "target_contract", contractRef)
		workspace.SetString(proposal, "base_version", "1")
		workspace.SetString(proposal, "base_fingerprint", contractFP)
		workspace.SetBool(proposal, "new_capability", false)
		workspace.SetStrings(proposal, "additions", []string{item.addition})
		workspace.SetStrings(proposal, "modification_from", []string{})
		workspace.SetStrings(proposal, "modification_to", []string{})
		workspace.SetStrings(proposal, "removals", []string{})
		workspace.SetStrings(proposal, "scope", []string{"v2"})
		writeTestDocument(t, filepath.Join(records, "proposal-"+item.proposalID+".md"), proposal)
		mission := testDocument(t, domain.Mission, item.missionID, "Assessed Mission", "awaiting-assessment")
		setTestFreshness(mission, markerFP)
		proposalTyped, _ := domain.ParseReference(proposalRef)
		mission.Record.Source = &proposalTyped
		workspace.SetStrings(mission, "evidence_claims", []string{})
		workspace.SetString(mission, "expires_at", "2026-08-11T10:00:00Z")
		writeTestDocument(t, filepath.Join(records, "mission-"+item.missionID+".md"), mission)
		assessment := testDocument(t, domain.Assessment, item.assessmentID, "Ready Assessment", "recorded")
		setTestFreshness(assessment, markerFP)
		workspace.SetString(assessment, "mission", missionRef)
		workspace.SetString(assessment, "verdict", "ready-for-owner")
		workspace.SetStrings(assessment, "claims", []string{})
		workspace.SetStrings(assessment, "evidence", []string{})
		writeTestDocument(t, filepath.Join(records, "assessment-"+item.assessmentID+".md"), assessment)
		decision := testDocument(t, domain.Decision, item.decisionID, "Reconcile", "recorded")
		setTestFreshness(decision, markerFP)
		workspace.SetString(decision, "actor_role", "owner")
		workspace.SetString(decision, "operation", "contract.reconcile")
		workspace.SetString(decision, "disposition", "accept")
		workspace.SetStrings(decision, "scope", []string{"v2"})
		workspace.SetStrings(decision, "authorized_effects", []string{"contract.reconcile"})
		workspace.SetStrings(decision, "conditions", []string{"no-provider-effects"})
		workspace.SetStrings(decision, "targets", []string{contractRef})
		workspace.SetStrings(decision, "expected_fingerprints", []string{contractFP})
		workspace.SetStrings(decision, "evidence", []string{assessmentRef})
		workspace.SetString(decision, "expires_at", "2026-08-11T10:00:00Z")
		writeTestDocument(t, filepath.Join(records, "decision-"+item.decisionID+".md"), decision)
		inputs = append(inputs, governance.ReconcileInput{Contract: contractRef, Proposal: proposalRef, Authorization: decisionRef, ExpectedFingerprint: contractFP, IdempotencyKey: "public-atomic-set"})
	}
	payload, err := json.Marshal(governance.ReconcileSetInput{Items: inputs})
	if err != nil {
		t.Fatal(err)
	}
	inputPath := filepath.Join(t.TempDir(), "set.json")
	if err := os.WriteFile(inputPath, payload, 0o644); err != nil {
		t.Fatal(err)
	}
	var out bytes.Buffer
	exit := (Runner{Cwd: root, Stdout: &out, Stderr: &bytes.Buffer{}, Now: fixedNow}).Run([]string{"contract", "reconcile-set", "--input", inputPath, "--json"})
	if exit != 0 || !strings.Contains(out.String(), `"schema_version":"spectacular.contract.reconcile-set.v1"`) || strings.Count(out.String(), `"operation":"contract.reconcile"`) != 2 {
		t.Fatalf("public reconcile-set exit=%d output=%s", exit, out.String())
	}
}

func testDocument(t *testing.T, noun domain.RecordType, rawID, title, status string) *workspace.Document {
	t.Helper()
	id, err := domain.ParseID(rawID)
	if err != nil {
		t.Fatal(err)
	}
	doc := &workspace.Document{Record: domain.Record{Type: noun, ID: id}, Body: "# " + string(noun) + "\n"}
	doc.Record.Title = testString(title)
	doc.Record.Status = testString(status)
	doc.Record.Created = testString("2026-08-10T10:00:00Z")
	doc.Record.Updated = testString("2026-08-10T10:00:00Z")
	return doc
}

func setTestFreshness(doc *workspace.Document, markerFP string) {
	workspace.SetString(doc, "freshness_checked_at", "2026-08-10T10:00:00Z")
	workspace.SetString(doc, "freshness_valid_until", "2026-08-11T10:00:00Z")
	workspace.SetString(doc, "freshness_source", ".spectacular/workspace.yaml")
	workspace.SetString(doc, "freshness_source_fingerprint", markerFP)
}

func writeTestDocument(t *testing.T, path string, doc *workspace.Document) {
	t.Helper()
	data, err := workspace.Canonical(doc)
	if err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(path, data, 0o644); err != nil {
		t.Fatal(err)
	}
}

func testString(value string) *string { return &value }

func replaceInFile(t *testing.T, path, old, replacement string) {
	t.Helper()
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	updated := bytes.Replace(data, []byte(old), []byte(replacement), 1)
	if bytes.Equal(data, updated) {
		t.Fatalf("%q not found in %s", old, path)
	}
	if err := os.WriteFile(path, updated, 0o644); err != nil {
		t.Fatal(err)
	}
}

type fileState struct {
	Mode    fs.FileMode
	ModTime time.Time
	Data    string
}

func snapshot(t *testing.T, root string) map[string]fileState {
	t.Helper()
	result := map[string]fileState{}
	err := filepath.WalkDir(root, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		rel, _ := filepath.Rel(root, path)
		info, e := d.Info()
		if e != nil {
			return e
		}
		state := fileState{Mode: info.Mode(), ModTime: info.ModTime()}
		if !d.IsDir() {
			data, e := os.ReadFile(path)
			if e != nil {
				return e
			}
			state.Data = string(data)
		}
		result[rel] = state
		return nil
	})
	if err != nil {
		t.Fatal(err)
	}
	return result
}
func copyFixture(t *testing.T) string {
	t.Helper()
	return copyTree(t, fixture(t))
}

func copyTree(t *testing.T, source string) string {
	t.Helper()
	target := t.TempDir()
	err := filepath.WalkDir(source, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		rel, _ := filepath.Rel(source, path)
		if rel == "." {
			return nil
		}
		dest := filepath.Join(target, rel)
		info, e := d.Info()
		if e != nil {
			return e
		}
		if d.IsDir() {
			return os.Mkdir(dest, info.Mode().Perm())
		}
		data, e := os.ReadFile(path)
		if e != nil {
			return e
		}
		if e = os.WriteFile(dest, data, info.Mode().Perm()); e != nil {
			return e
		}
		return os.Chtimes(dest, info.ModTime(), info.ModTime())
	})
	if err != nil {
		t.Fatal(err)
	}
	return target
}
