package test

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestFullVerificationGateRunsAcceptanceBeforeReleaseChecks(t *testing.T) {
	verify := readRepositoryFile(t, "test", "verify.sh")
	allMode := between(t, verify, "  all)\n", "  *)\n")

	acceptance := strings.Index(allMode, "acceptance_checks")
	release := strings.Index(allMode, "release_checks")
	if acceptance < 0 || release < 0 {
		t.Fatalf("all verification mode must run acceptance and release checks:\n%s", allMode)
	}
	if acceptance > release {
		t.Fatalf("all verification mode must run acceptance before release checks:\n%s", allMode)
	}
}

func TestReleaseWorkflowVerifiesBeforePublishing(t *testing.T) {
	workflow := readRepositoryFile(t, ".github", "workflows", "release.yml")
	fullGate := strings.Index(workflow, "run: bash test/verify.sh all")
	publish := strings.Index(workflow, "uses: softprops/action-gh-release")
	if fullGate < 0 || publish < 0 {
		t.Fatalf("release workflow must run the full gate and publish a release:\n%s", workflow)
	}
	if fullGate > publish {
		t.Fatal("release workflow publishes before the full verification gate")
	}
}

func TestVerificationWorkflowRunsAcceptanceOnSupportedPlatforms(t *testing.T) {
	workflow := readRepositoryFile(t, ".github", "workflows", "verify.yml")
	for _, runner := range []string{"ubuntu-latest", "macos-latest", "windows-latest"} {
		if !strings.Contains(workflow, runner) {
			t.Fatalf("verification workflow must run acceptance fixtures on %s", runner)
		}
	}
	if !strings.Contains(workflow, "run: go test -count=1 ./test/acceptance") {
		t.Fatal("verification workflow must run the acceptance fixtures in its platform matrix")
	}
}

func readRepositoryFile(t *testing.T, path ...string) string {
	t.Helper()
	data, err := os.ReadFile(filepath.Join(append([]string{".."}, path...)...))
	if err != nil {
		t.Fatal(err)
	}
	return string(data)
}

func between(t *testing.T, value, start, end string) string {
	t.Helper()
	from := strings.Index(value, start)
	if from < 0 {
		t.Fatalf("missing start marker %q", start)
	}
	to := strings.Index(value[from:], end)
	if to < 0 {
		t.Fatalf("missing end marker %q", end)
	}
	return value[from : from+to]
}
