// release-smoke drives an installed binary through the compact Mission loop.
package main

import (
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io"
	"io/fs"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"
)

const contractRef = "Contract:0199b000-0000-7000-8000-000000000001"

type envelope struct {
	Data map[string]any `json:"data"`
}

type runner struct {
	binary    string
	workspace string
	inputs    string
}

func main() {
	binary := flag.String("binary", "", "installed spectacular binary")
	fixture := flag.String("fixture", "testdata/scenario-b-c", "v2 fixture")
	workspace := flag.String("workspace", "", "empty disposable workspace")
	flag.Parse()
	if *binary == "" || *workspace == "" || flag.NArg() != 0 {
		fatal(errors.New("usage: release-smoke --binary <binary> --workspace <dir> [--fixture path]"))
	}
	abs, err := filepath.Abs(*binary)
	if err != nil {
		fatal(err)
	}
	if err := copyFixture(*fixture, *workspace); err != nil {
		fatal(err)
	}
	inputs, err := os.MkdirTemp("", "spectacular-release-smoke-")
	if err != nil {
		fatal(err)
	}
	defer os.RemoveAll(inputs)
	r := runner{binary: abs, workspace: *workspace, inputs: inputs}
	if err := r.initializeGit(); err != nil {
		fatal(err)
	}
	if err := r.run(); err != nil {
		fatal(err)
	}
	fmt.Println("result=installed-binary-compact-mission-pass")
}

func (r runner) run() error {
	version, err := exec.Command(r.binary, "--version").CombinedOutput()
	if err != nil || !regexp.MustCompile(`^spectacular 2\.`).Match(version) {
		return fmt.Errorf("version inspection: %w: %s", err, version)
	}
	plan := `---
type: MissionPlan
title: Installed release smoke
owner: release-owner
contract:
  ref: ` + contractRef + `
outcome: Prove the installed compact Mission lifecycle.
review: independent
completion:
  - claim: installed-loop
    pass_boundary: Installed commands complete and cold-check the Mission.
    proof_requirement: Separate processes return successful typed results.
objectives:
  - outcome: Exercise the installed loop.
    claims: [installed-loop]
authority:
  operator: [inspect, edit-in-scope, choose-reversible-implementation, run-checks, generate-derived-files, bounded-repair, commit-local]
  requires_owner: [activate-mission, change-outcome-or-completion, expand-scope, push, merge, release, irreversible-change, destructive-data, secret-change]
scope:
  mechanical: [.spectacular/]
  semantic: [Installed compact lifecycle proof.]
repair_budget: 1
dependencies: []
gaps: []
stops: [scope-drift]
---
# Mission

Exercise the installed binary.
`
	started, err := r.commandInput(plan, "mission", "start", "-", "--json")
	if err != nil {
		return err
	}
	missionRef := stringField(started, "ref")
	if _, err := r.command("objective", "promote", missionRef+"/O1", "--json"); err != nil {
		return err
	}
	if _, err := r.command("objective", "finish", missionRef+"/O1", "--json"); err != nil {
		return err
	}
	if _, err := r.command("run", "start", missionRef, "--title", "Cold completion run", "--json"); err != nil {
		return err
	}
	checked, err := r.command("mission", "check", missionRef, "--json")
	if err != nil {
		return err
	}
	commit, err := r.git("rev-parse", "HEAD")
	if err != nil {
		return err
	}
	tree, err := r.git("rev-parse", "HEAD^{tree}")
	if err != nil {
		return err
	}
	review := `---
type: ReviewDraft
title: Installed lifecycle review
status: passed
reviewed:
  commit: ` + commit + `
  tree: ` + tree + `
  activation_fingerprint: ` + stringField(checked, "fingerprint") + `
reviewer:
  actor: release-smoke-reviewer
  operator: release-smoke-operator
  relation_to_operator: independent
  implemented_reviewed_scope: false
  independence_basis: Separate installed-binary review step.
  evidence: [release-smoke:separate-process]
claims:
  - claim: installed-loop
    verdict: pass
findings: []
limitations: [Disposable local workspace.]
---
# Review

The installed command loop passes.
`
	reviewPath := filepath.Join(r.inputs, "review.md")
	if err := os.WriteFile(reviewPath, []byte(review), 0o644); err != nil {
		return err
	}
	if _, err := r.command("review", "record", missionRef, reviewPath, "--json"); err != nil {
		return err
	}
	if _, err := r.command("mission", "complete", missionRef, "--by", "release-owner", "--json"); err != nil {
		return err
	}
	final, err := r.command("mission", "check", missionRef, "--json")
	if err != nil || final.Data["valid"] != true {
		return fmt.Errorf("final Mission check failed: %w %#v", err, final.Data)
	}
	shown, err := r.command("mission", "show", missionRef, "--json")
	if err != nil || shown.Data["status"] != "completed" {
		return fmt.Errorf("cold Mission show failed: %w %#v", err, shown.Data)
	}
	return nil
}

func (r runner) command(args ...string) (envelope, error) {
	return r.commandInput("", args...)
}

func (r runner) commandInput(input string, args ...string) (envelope, error) {
	cmd := exec.Command(r.binary, args...)
	cmd.Dir = r.workspace
	if input != "" {
		cmd.Stdin = strings.NewReader(input)
	}
	output, err := cmd.CombinedOutput()
	if err != nil {
		return envelope{}, fmt.Errorf("spectacular %s: %w\n%s", strings.Join(args, " "), err, output)
	}
	var result envelope
	if err := json.Unmarshal(output, &result); err != nil {
		return envelope{}, fmt.Errorf("decode spectacular %s: %w\n%s", strings.Join(args, " "), err, output)
	}
	return result, nil
}

func (r runner) initializeGit() error {
	for _, args := range [][]string{{"init", "-q"}, {"config", "user.name", "Release Smoke"}, {"config", "user.email", "release-smoke@example.invalid"}, {"add", ".spectacular"}, {"commit", "-qm", "fixture"}} {
		if _, err := r.git(args...); err != nil {
			return err
		}
	}
	return nil
}

func (r runner) git(args ...string) (string, error) {
	cmd := exec.Command("git", args...)
	cmd.Dir = r.workspace
	output, err := cmd.CombinedOutput()
	return strings.TrimSpace(string(output)), err
}

func stringField(value envelope, name string) string {
	text, _ := value.Data[name].(string)
	return text
}

func copyFixture(source, destination string) error {
	return filepath.WalkDir(source, func(path string, entry fs.DirEntry, walkErr error) error {
		if walkErr != nil {
			return walkErr
		}
		relative, err := filepath.Rel(source, path)
		if err != nil {
			return err
		}
		target := filepath.Join(destination, relative)
		if entry.IsDir() {
			return os.MkdirAll(target, 0o755)
		}
		input, err := os.Open(path)
		if err != nil {
			return err
		}
		defer input.Close()
		output, err := os.OpenFile(target, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0o644)
		if err != nil {
			return err
		}
		if _, err := io.Copy(output, input); err != nil {
			output.Close()
			return err
		}
		return output.Close()
	})
}

func fatal(err error) {
	fmt.Fprintln(os.Stderr, err)
	os.Exit(1)
}
