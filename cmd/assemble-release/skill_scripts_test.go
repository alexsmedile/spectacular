package main

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

func TestDoctorRequiresExactCLICompatibility(t *testing.T) {
	script := filepath.Join("..", "..", "skills", "spectacular", "scripts", "doctor.sh")
	// The compatible version is read from VERSION rather than written here. A
	// hardcoded literal makes this test a release-blocking fixture: it passes until
	// the version is bumped, then fails inside the release workflow itself, long
	// after the bump looked clean. The release bumper does not edit Go test files.
	released, err := os.ReadFile(filepath.Join("..", "..", "VERSION"))
	if err != nil {
		t.Fatal(err)
	}
	current := strings.TrimSpace(string(released))
	for _, test := range []struct {
		name    string
		version string
		want    string
	}{
		{name: "compatible", version: current, want: "full — read, draft, and governed execution"},
		{name: "incompatible", version: "0.0.1", want: "reduced — read, explain, and draft only"},
	} {
		t.Run(test.name, func(t *testing.T) {
			bin := t.TempDir()
			fake := filepath.Join(bin, "spectacular")
			body := "#!/bin/sh\nprintf '%s\\n' '{\"commit\":\"test\",\"name\":\"spectacular\",\"schema_version\":\"spectacular.build-info.v1\",\"version\":\"" + test.version + "\"}'\n"
			if err := os.WriteFile(fake, []byte(body), 0o755); err != nil {
				t.Fatal(err)
			}
			command := exec.Command("sh", script)
			command.Env = append(os.Environ(), "PATH="+bin+string(os.PathListSeparator)+os.Getenv("PATH"))
			output, err := command.CombinedOutput()
			if err != nil {
				t.Fatalf("doctor: %v\n%s", err, output)
			}
			if !strings.Contains(string(output), test.want) {
				t.Fatalf("output does not contain %q:\n%s", test.want, output)
			}
		})
	}
}
