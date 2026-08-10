package command

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"io/fs"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"
	"time"
)

const (
	missionID    = "0198a1a0-0000-7000-8000-000000000002"
	gapID        = "0198a1a0-0000-7000-8000-000000000004"
	runID        = "0198a1a0-0000-7000-8000-000000000005"
	checkpointID = "0198a1a0-0000-7000-8000-000000000006"
	evidenceID   = "0198a1a0-0000-7000-8000-000000000007"
	decisionID   = "0198a1a0-0000-7000-8000-000000000008"
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

func TestPublicRegistryAndCommands(t *testing.T) {
	if len(Registry) != 25 {
		t.Fatalf("registry has %d commands, want 25", len(Registry))
	}
	seen := map[Operation]bool{}
	for _, spec := range Registry {
		if spec.Operation == 0 || seen[spec.Operation] {
			t.Fatalf("unbound or duplicate registry operation %d", spec.Operation)
		}
		seen[spec.Operation] = true
	}
	if len(seen) != int(opMissionArchive) {
		t.Fatalf("registry dispatch parity: %d bound operations, want %d", len(seen), opMissionArchive)
	}
	tests := []struct {
		args     []string
		schema   string
		contains string
	}{
		{[]string{"anchor", "show", "project", "--json"}, "spectacular.anchor.show.v1", `"authoritative":{"identity":{"noun":"anchor","ref":"project"`},
		{[]string{"mission", "list", "--json"}, "spectacular.mission.list.v1", `"items":[`},
		{[]string{"mission", "show", "Mission:" + missionID, "--json"}, "spectacular.mission.show.v1", `"path":".spectacular/records/proposal.md"`},
		{[]string{"gap", "list", "--scope", missionID, "--json"}, "spectacular.gap.list.v1", `"noun":"Gap"`},
		{[]string{"gap", "show", "Gap:" + gapID, "--json"}, "spectacular.gap.show.v1", `"show_command":"spectacular mission show`},
		{[]string{"run", "show", ".spectacular/records/run.md", "--json"}, "spectacular.run.show.v1", `"show_command":"spectacular checkpoint show`},
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

func TestJSONIsDeterministicAndExactLookupFormsAgree(t *testing.T) {
	forms := []string{missionID, "Mission:" + missionID, ".spectacular/records/mission.md"}
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
	decision := filepath.Join(root, ".spectacular", "records", "decision.md")
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
	evidence := filepath.Join(root, ".spectacular", "records", "evidence.md")
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
			path := filepath.Join(root, ".spectacular", "records", "decision.md")
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

func TestProposalCurrentTruthAndMissionSourceHaveHonestDrillDown(t *testing.T) {
	root := copyFixture(t)
	project := filepath.Join(root, ".spectacular", "records", "project.md")
	replaceInFile(t, project, "current_truth:\n  - Mission:"+missionID, "current_truth:\n  - Mission:"+missionID+"\n  - Proposal:0198a1a0-0000-7000-8000-000000000001")

	for _, args := range [][]string{{"anchor", "show", "project", "--json"}, {"mission", "show", missionID, "--json"}} {
		var out bytes.Buffer
		exit := (Runner{Cwd: root, Stdout: &out, Stderr: &bytes.Buffer{}, Now: fixedNow}).Run(args)
		if exit != 0 {
			t.Fatalf("%v exit=%d output=%s", args, exit, out.String())
		}
		text := out.String()
		for _, required := range []string{`"noun":"proposal"`, `"ref":"Proposal:0198a1a0-0000-7000-8000-000000000001"`, `"path":".spectacular/records/proposal.md"`, `"fingerprint":"`, `"show_command":"spectacular proposal show Proposal:`} {
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
				data, err := os.ReadFile(filepath.Join(root, ".spectacular", "records", "gap.md"))
				if err != nil {
					t.Fatal(err)
				}
				if err := os.WriteFile(filepath.Join(root, ".spectacular", "records", "duplicate-gap.md"), data, 0o644); err != nil {
					t.Fatal(err)
				}
			},
			args: []string{"anchor", "show", "project", "--json"},
			code: "duplicate_id",
		},
		{
			name: "multiple decisions conflict",
			mutate: func(t *testing.T, root string) {
				path := filepath.Join(root, ".spectacular", "records", "decision.md")
				data, err := os.ReadFile(path)
				if err != nil {
					t.Fatal(err)
				}
				data = bytes.Replace(data, []byte(decisionID), []byte("0198a1a0-0000-7000-8000-000000000009"), 1)
				if err := os.WriteFile(filepath.Join(root, ".spectacular", "records", "decision-2.md"), data, 0o644); err != nil {
					t.Fatal(err)
				}
			},
			args: []string{"mission", "show", missionID, "--json"},
			code: "conflicting_authority",
		},
		{
			name: "missing freshness source",
			mutate: func(t *testing.T, root string) {
				replaceInFile(t, filepath.Join(root, ".spectacular", "records", "evidence.md"), ".spectacular/workspace.yaml", ".spectacular/missing.yaml")
			},
			args: []string{"evidence", "show", evidenceID, "--json"},
			code: "record_not_found",
		},
		{
			name: "malformed freshness fingerprint",
			mutate: func(t *testing.T, root string) {
				path := filepath.Join(root, ".spectacular", "records", "evidence.md")
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
	if err := os.Remove(filepath.Join(root, ".spectacular", "records", "decision.md")); err != nil {
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
	const want = "9c2d6d213bd8e66c98643a0a0ea1009b93ecfa9dc1ef269445399284396da937"
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
	if exit := runner.Run([]string{"decision", "create", "--input", path, "--json"}); exit != 0 || !strings.Contains(out.String(), `"operation":"decision.create"`) {
		t.Fatalf("create exit=%d output=%s", exit, out.String())
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
