package acceptance_test

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
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
		os.Exit(1)
	}
	repoRoot = filepath.Dir(filepath.Dir(source))
	var err error
	buildRoot, err = os.MkdirTemp("", "spectacular-acceptance-")
	if err != nil {
		os.Exit(1)
	}
	defer os.RemoveAll(buildRoot)
	version, _ := os.ReadFile(filepath.Join(repoRoot, "VERSION"))
	cliBinary = filepath.Join(buildRoot, "spectacular")
	smokeBinary = filepath.Join(buildRoot, "release-smoke")
	ldflags := "-X github.com/alexsmedile/spectacular/v2/internal/buildinfo.Version=" + strings.TrimSpace(string(version)) + " -X github.com/alexsmedile/spectacular/v2/internal/buildinfo.Commit=acceptance"
	if output, buildErr := build("-trimpath", "-buildvcs=false", "-ldflags", ldflags, "-o", cliBinary, "./cmd/spectacular"); buildErr != nil {
		fmt.Fprintln(os.Stderr, string(output))
		os.Exit(1)
	}
	if output, buildErr := build("-trimpath", "-buildvcs=false", "-o", smokeBinary, "./cmd/release-smoke"); buildErr != nil {
		fmt.Fprintln(os.Stderr, string(output))
		os.Exit(1)
	}
	os.Exit(m.Run())
}

func TestSelfHostedMissionsValidateWithInstalledBinary(t *testing.T) {
	for _, ref := range []string{"M5", "M6"} {
		result := command(t, repoRoot, 0, "mission", "check", ref, "--json")
		var payload struct {
			Schema string `json:"schema_version"`
			Data   struct {
				Valid  bool   `json:"valid"`
				Schema string `json:"schema"`
			} `json:"data"`
		}
		if err := json.Unmarshal([]byte(result), &payload); err != nil || !payload.Data.Valid || payload.Data.Schema != "mission.v2" {
			t.Fatalf("check %s: %v %s", ref, err, result)
		}
	}
	plain := command(t, repoRoot, 0, "mission", "show", "M6")
	if !strings.Contains(plain, "M6 — Implement the compact Mission CLI") || strings.Contains(plain, `"schema_version"`) {
		t.Fatalf("default output is not compact human text: %s", plain)
	}
}

func TestInstalledBinaryCompletesCompactLifecycle(t *testing.T) {
	workspace := t.TempDir()
	cmd := exec.Command(smokeBinary,
		"--binary", cliBinary,
		"--fixture", filepath.Join(repoRoot, "testdata", "scenario-b-c"),
		"--workspace", workspace,
	)
	cmd.Dir = repoRoot
	output, err := cmd.CombinedOutput()
	if err != nil || !strings.Contains(string(output), "installed-binary-compact-mission-pass") {
		t.Fatalf("release smoke: %v\n%s", err, output)
	}
	if matches, _ := filepath.Glob(filepath.Join(workspace, ".spectacular", "missions", "M1-*", "index.md")); len(matches) != 0 {
		t.Fatalf("compact Mission gained a local index: %v", matches)
	}
	command(t, workspace, 0, "mission", "check", "M1", "--json")
}

func TestLegacyMissionRemainsReadableWithoutOldCommands(t *testing.T) {
	fixture := filepath.Join(repoRoot, "testdata", "scenario-a")
	shown := command(t, fixture, 0, "mission", "show", "M1", "--json")
	if !strings.Contains(shown, `"legacy":true`) {
		t.Fatalf("legacy Mission not readable: %s", shown)
	}
	for _, old := range [][]string{{"proposal", "create"}, {"mission", "prepare"}, {"contract", "reconcile"}, {"workspace", "context"}} {
		command(t, fixture, 2, old...)
	}
}

func TestUsageAndRefusalAreTyped(t *testing.T) {
	usage := command(t, repoRoot, 2, "mission", "show", "--json")
	if !strings.Contains(usage, `"code":"usage"`) || !strings.Contains(usage, `"safe_correction"`) {
		t.Fatalf("usage=%s", usage)
	}
	refusal := command(t, repoRoot, 3, "mission", "check", "M999", "--json")
	if !strings.Contains(refusal, `"schema_version":"spectacular.refusal.v2"`) || !strings.Contains(refusal, `"mutation":"none"`) {
		t.Fatalf("refusal=%s", refusal)
	}
}

func build(args ...string) ([]byte, error) {
	cmd := exec.Command("go", append([]string{"build"}, args...)...)
	cmd.Dir = repoRoot
	cmd.Env = append(os.Environ(), "GOPROXY=off", "GOFLAGS=-mod=readonly")
	return cmd.CombinedOutput()
}

func command(t *testing.T, cwd string, expected int, args ...string) string {
	t.Helper()
	cmd := exec.Command(cliBinary, args...)
	cmd.Dir = cwd
	output, err := cmd.CombinedOutput()
	exit := 0
	if err != nil {
		var exitErr *exec.ExitError
		if !errors.As(err, &exitErr) {
			t.Fatalf("spectacular %s: %v", strings.Join(args, " "), err)
		}
		exit = exitErr.ExitCode()
	}
	if exit != expected {
		t.Fatalf("spectacular %s exit=%d want=%d\n%s", strings.Join(args, " "), exit, expected, output)
	}
	return string(output)
}
