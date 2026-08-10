package acceptance_test

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io/fs"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"sort"
	"strings"
	"testing"
)

var (
	repoRoot    string
	cliBinary   string
	smokeBinary string
	buildRoot   string
)

func TestMain(m *testing.M) {
	_, source, _, ok := runtime.Caller(0)
	if !ok {
		fmt.Fprintln(os.Stderr, "resolve acceptance source")
		os.Exit(1)
	}
	repoRoot = filepath.Dir(filepath.Dir(source))
	var err error
	buildRoot, err = os.MkdirTemp("", "spectacular-acceptance-")
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	defer os.RemoveAll(buildRoot)
	versionBytes, err := os.ReadFile(filepath.Join(repoRoot, "VERSION"))
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	version := strings.TrimSpace(string(versionBytes))
	cliBinary = filepath.Join(buildRoot, "spectacular")
	smokeBinary = filepath.Join(buildRoot, "release-smoke")
	ldflags := "-X github.com/alexsmedile/spectacular/v2/internal/buildinfo.Version=" + version + " -X github.com/alexsmedile/spectacular/v2/internal/buildinfo.Commit=acceptance"
	if output, buildErr := build("-trimpath", "-buildvcs=false", "-ldflags", ldflags, "-o", cliBinary, "./cmd/spectacular"); buildErr != nil {
		fmt.Fprintf(os.Stderr, "build CLI: %v\n%s", buildErr, output)
		os.Exit(1)
	}
	if output, buildErr := build("-trimpath", "-buildvcs=false", "-o", smokeBinary, "./cmd/release-smoke"); buildErr != nil {
		fmt.Fprintf(os.Stderr, "build smoke harness: %v\n%s", buildErr, output)
		os.Exit(1)
	}
	os.Exit(m.Run())
}

func build(args ...string) ([]byte, error) {
	cmd := exec.Command("go", append([]string{"build"}, args...)...)
	cmd.Dir = repoRoot
	cmd.Env = append(os.Environ(), "GOPROXY=off", "GOFLAGS=-mod=readonly")
	return cmd.CombinedOutput()
}

func TestHumanWorkspaceColdRecoveryUsesRealProcesses(t *testing.T) {
	workspace := t.TempDir()
	copyTree(t, filepath.Join(repoRoot, ".spectacular"), filepath.Join(workspace, ".spectacular"))
	before := snapshot(t, filepath.Join(workspace, ".spectacular"))

	project := command(t, workspace, 0, "anchor", "show", "project")
	for _, required := range []string{
		"Make governed agentic software work legible",
		"Current truth: PRODUCT",
		"Current truth: ARCHITECTURE",
		"Current truth: STACK",
		"Mission:       M1",
		"Owner gate:",
	} {
		if !strings.Contains(project.stdout, required) {
			t.Fatalf("project orientation omits %q:\n%s", required, project.stdout)
		}
	}
	mission := command(t, workspace, 0, "mission", "show", "M1")
	for _, required := range []string{"Objective: M1/O1", "Run: M1/R1", "Checkpoint: M1/R1/C1", "Gap: M1/G1-"} {
		if !strings.Contains(mission.stdout, required) {
			t.Fatalf("Mission card omits %q:\n%s", required, mission.stdout)
		}
	}

	commands := map[string]bool{}
	for _, seed := range [][]string{{"anchor", "show", "project", "--json"}, {"mission", "show", "M1", "--json"}} {
		result := command(t, workspace, 0, seed...)
		var value any
		if err := json.Unmarshal([]byte(result.stdout), &value); err != nil {
			t.Fatal(err)
		}
		collectShowCommands(value, commands)
	}
	if len(commands) < 8 {
		t.Fatalf("cold projections exposed only %d drill-down commands", len(commands))
	}
	ordered := make([]string, 0, len(commands))
	for line := range commands {
		ordered = append(ordered, line)
	}
	sort.Strings(ordered)
	for _, line := range ordered {
		fields := strings.Fields(line)
		if len(fields) < 3 || fields[0] != "spectacular" {
			t.Fatalf("invalid show command %q", line)
		}
		result := command(t, workspace, 0, append(fields[1:], "--json")...)
		var envelope struct {
			Schema string         `json:"schema_version"`
			Data   map[string]any `json:"data"`
		}
		if err := json.Unmarshal([]byte(result.stdout), &envelope); err != nil || envelope.Schema == "" || envelope.Data == nil {
			t.Fatalf("show command %q returned invalid envelope: %v\n%s", line, err, result.stdout)
		}
	}
	command(t, workspace, 0, "workspace", "validate", "project", "--json")
	if after := snapshot(t, filepath.Join(workspace, ".spectacular")); after != before {
		t.Fatal("cold recovery or pointer traversal mutated bytes, modes, or mtimes")
	}
}

func TestOfficialFixturesExerciseOnlyHumanLayout(t *testing.T) {
	for _, name := range []string{"scenario-a", "scenario-b-c"} {
		t.Run(name, func(t *testing.T) {
			root := filepath.Join(repoRoot, "testdata", name)
			if _, err := os.Stat(filepath.Join(root, ".spectacular", "records")); !os.IsNotExist(err) {
				t.Fatalf("fixture retains flat records directory: %v", err)
			}
			manifest, err := os.ReadFile(filepath.Join(root, ".spectacular", "workspace.yaml"))
			if err != nil {
				t.Fatal(err)
			}
			text := string(manifest)
			if !strings.Contains(text, "  - .\n") || !strings.Contains(text, "project_anchor: PROJECT.md") {
				t.Fatalf("fixture does not use current human layout:\n%s", text)
			}
			command(t, root, 0, "workspace", "validate", "project", "--json")
		})
	}
	command(t, filepath.Join(repoRoot, "testdata", "scenario-a"), 0, "mission", "show", "M1", "--json")
	command(t, filepath.Join(repoRoot, "testdata", "scenario-a"), 0, "checkpoint", "show", "M1/R1/C1", "--json")
}

func TestGovernedLoopPersistsAndArchivesHumanBundle(t *testing.T) {
	workspace := t.TempDir()
	cmd := exec.Command(smokeBinary,
		"--binary", cliBinary,
		"--fixture", filepath.Join(repoRoot, "testdata", "scenario-b-c"),
		"--workspace", workspace,
	)
	cmd.Dir = repoRoot
	if output, err := cmd.CombinedOutput(); err != nil {
		t.Fatalf("governed loop failed: %v\n%s", err, output)
	}
	archived := filepath.Join(workspace, ".spectacular", "archive", "missions", "M1-installed-release-smoke")
	if _, err := os.Stat(filepath.Join(archived, "MISSION.md")); err != nil {
		t.Fatalf("Mission bundle was not archived atomically: %v", err)
	}
	if _, err := os.Stat(filepath.Join(workspace, ".spectacular", "missions", "M1-installed-release-smoke")); !os.IsNotExist(err) {
		t.Fatalf("active Mission bundle survived archival: %v", err)
	}
	activeIndex, err := os.ReadFile(filepath.Join(workspace, ".spectacular", "missions", "index.md"))
	if err != nil {
		t.Fatal(err)
	}
	if strings.Contains(string(activeIndex), "`M1`") {
		t.Fatalf("active Mission index retained archived Mission:\n%s", activeIndex)
	}
	if err := filepath.WalkDir(filepath.Join(workspace, ".spectacular"), func(path string, entry fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if entry.IsDir() && entry.Name() == "records" {
			return fmt.Errorf("flat records directory created at %s", path)
		}
		return nil
	}); err != nil {
		t.Fatal(err)
	}
	command(t, workspace, 0, "workspace", "validate", "project", "--json")
	command(t, workspace, 0, "anchor", "show", "project")
}

func TestUsageAndRefusalAreZeroMutationProcessBoundaries(t *testing.T) {
	workspace := t.TempDir()
	before := snapshot(t, workspace)
	command(t, workspace, 2, "mission", "show")
	command(t, workspace, 3, "mission", "show", "M1", "--json")
	if after := snapshot(t, workspace); after != before {
		t.Fatal("usage or refusal created workspace state")
	}
}

type commandResult struct {
	stdout string
	stderr string
}

func command(t *testing.T, cwd string, expected int, args ...string) commandResult {
	t.Helper()
	cmd := exec.Command(cliBinary, args...)
	cmd.Dir = cwd
	var stdout, stderr strings.Builder
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	err := cmd.Run()
	exit := 0
	if err != nil {
		var exitErr *exec.ExitError
		if !errors.As(err, &exitErr) {
			t.Fatalf("spectacular %s: %v", strings.Join(args, " "), err)
		}
		exit = exitErr.ExitCode()
	}
	if exit != expected {
		t.Fatalf("spectacular %s exit=%d want=%d\nstdout=%s\nstderr=%s", strings.Join(args, " "), exit, expected, stdout.String(), stderr.String())
	}
	return commandResult{stdout: stdout.String(), stderr: stderr.String()}
}

func collectShowCommands(value any, commands map[string]bool) {
	switch item := value.(type) {
	case map[string]any:
		if line, ok := item["show_command"].(string); ok {
			commands[line] = true
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

func copyTree(t *testing.T, source, destination string) {
	t.Helper()
	err := filepath.WalkDir(source, func(path string, entry fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		relative, err := filepath.Rel(source, path)
		if err != nil {
			return err
		}
		target := filepath.Join(destination, relative)
		info, err := entry.Info()
		if err != nil {
			return err
		}
		if entry.IsDir() {
			return os.MkdirAll(target, info.Mode().Perm())
		}
		if !entry.Type().IsRegular() {
			return fmt.Errorf("non-regular fixture entry: %s", path)
		}
		data, err := os.ReadFile(path)
		if err != nil {
			return err
		}
		if err := os.WriteFile(target, data, info.Mode().Perm()); err != nil {
			return err
		}
		return os.Chtimes(target, info.ModTime(), info.ModTime())
	})
	if err != nil {
		t.Fatal(err)
	}
}

func snapshot(t *testing.T, root string) string {
	t.Helper()
	hash := sha256.New()
	err := filepath.WalkDir(root, func(path string, entry fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		relative, err := filepath.Rel(root, path)
		if err != nil {
			return err
		}
		info, err := entry.Info()
		if err != nil {
			return err
		}
		fmt.Fprintf(hash, "%s\x00%s\x00%d\x00", filepath.ToSlash(relative), info.Mode(), info.ModTime().UnixNano())
		if entry.Type().IsRegular() {
			data, err := os.ReadFile(path)
			if err != nil {
				return err
			}
			hash.Write(data)
		}
		return nil
	})
	if err != nil {
		t.Fatal(err)
	}
	return hex.EncodeToString(hash.Sum(nil))
}
