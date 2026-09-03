package spectaculareval

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

func TestInspectPackageMeasuresRoutedContextAndFindings(t *testing.T) {
	root := t.TempDir()
	links := ""
	for _, name := range primaryRoutes {
		links += "[" + name + "](references/" + name + ".md)\n"
	}
	writeTestFile(t, filepath.Join(root, "SKILL.md"), "---\nname: x\n---\n# X\n\n"+links)
	writeTestFile(t, filepath.Join(root, "references", "orient.md"), "# Orient\n\nUse this when: orienting.\n\nwords here\n")
	for _, name := range []string{"prepare", "execute", "runtime", "close", "audit"} {
		writeTestFile(t, filepath.Join(root, "references", name+".md"), "# Ref\n\nUse this when: routed.\n")
	}
	stats, err := InspectPackage("candidate", "working", "", root)
	if err != nil {
		t.Fatal(err)
	}
	if stats.KernelBodyLines == 0 || stats.PrimaryRouteWords["orient"] <= stats.KernelWords {
		t.Fatalf("stats=%+v", stats)
	}
	if len(stats.ValidationFindings) != 0 {
		t.Fatalf("findings=%v", stats.ValidationFindings)
	}
}

func TestInspectPackageFindsBrokenAndOrphanPointers(t *testing.T) {
	root := t.TempDir()
	writeTestFile(t, filepath.Join(root, "SKILL.md"), "---\nname: x\n---\n# X\n\n[Orient](references/orient.md)\n[Missing](references/missing.md)\n")
	writeTestFile(t, filepath.Join(root, "references", "orient.md"), "# Orient\n\nUse this when: orienting.\n")
	writeTestFile(t, filepath.Join(root, "references", "orphan.md"), "# Orphan\n\nUse this when: never.\n")
	stats, err := InspectPackage("candidate", "working", "", root)
	if err != nil {
		t.Fatal(err)
	}
	joined := strings.Join(stats.ValidationFindings, "\n")
	if !strings.Contains(joined, "broken link") || !strings.Contains(joined, "orphan reference") {
		t.Fatalf("findings=%v", stats.ValidationFindings)
	}
}

func TestApplyStaticThresholdsRejectsOversizedOrInsufficientKernelGain(t *testing.T) {
	report := ComparePackages(
		PackageStats{KernelBodyLines: 100, KernelWords: 1000, PrimaryRouteWords: map[string]int{}},
		PackageStats{KernelBodyLines: 91, KernelWords: 800, PrimaryRouteWords: map[string]int{}},
	)
	ApplyStaticThresholds(&report, Thresholds{MaximumKernelBodyLines: 90, MinimumInitialContextGain: 0.5})
	if report.Verdict != "regression" || len(report.GateFailures) != 2 {
		t.Fatalf("report=%+v", report)
	}
}

func TestMaterializeSkillUsesImmutableCommitAndCompletePackage(t *testing.T) {
	repo := t.TempDir()
	runGit(t, repo, "init", "-q")
	runGit(t, repo, "config", "user.name", "Eval")
	runGit(t, repo, "config", "user.email", "eval@example.invalid")
	writeTestFile(t, filepath.Join(repo, "skills", "spectacular", "SKILL.md"), "---\nname: spectacular\n---\n# Old\n")
	writeTestFile(t, filepath.Join(repo, "skills", "spectacular", "references", "orient.md"), "# Orient\n")
	writeTestFile(t, filepath.Join(repo, "skills", "spectacular", "scripts", "doctor.sh"), "#!/bin/sh\nexit 0\n")
	if err := os.Chmod(filepath.Join(repo, "skills", "spectacular", "scripts", "doctor.sh"), 0o755); err != nil {
		t.Fatal(err)
	}
	runGit(t, repo, "add", ".")
	runGit(t, repo, "commit", "-qm", "baseline")
	baseline := strings.TrimSpace(runGit(t, repo, "rev-parse", "HEAD"))
	writeTestFile(t, filepath.Join(repo, "skills", "spectacular", "SKILL.md"), "---\nname: spectacular\n---\n# New\n")
	runGit(t, repo, "add", ".")
	runGit(t, repo, "commit", "-qm", "candidate")
	destination := filepath.Join(t.TempDir(), "skill")
	commit, err := MaterializeSkill(repo, baseline, destination)
	if err != nil {
		t.Fatal(err)
	}
	if commit != baseline {
		t.Fatalf("commit=%s want %s", commit, baseline)
	}
	data, err := os.ReadFile(filepath.Join(destination, "SKILL.md"))
	if err != nil || !strings.Contains(string(data), "# Old") {
		t.Fatalf("data=%q err=%v", data, err)
	}
	if _, err := os.Stat(filepath.Join(destination, "references", "orient.md")); err != nil {
		t.Fatal("complete package reference missing:", err)
	}
	if info, err := os.Stat(filepath.Join(destination, "scripts", "doctor.sh")); err != nil || info.Mode().Perm() != 0o755 {
		t.Fatalf("executable mode was not preserved: mode=%v err=%v", info, err)
	}
}

func TestCopyTreePreservesExecutableModes(t *testing.T) {
	source := t.TempDir()
	destination := filepath.Join(t.TempDir(), "copy")
	path := filepath.Join(source, "check.sh")
	writeTestFile(t, path, "#!/bin/sh\n")
	if err := os.Chmod(path, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := CopyTree(source, destination); err != nil {
		t.Fatal(err)
	}
	info, err := os.Stat(filepath.Join(destination, "check.sh"))
	if err != nil || info.Mode().Perm() != 0o755 {
		t.Fatalf("mode=%v err=%v", info, err)
	}
}

func writeTestFile(t *testing.T, path, content string) {
	t.Helper()
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		t.Fatal(err)
	}
}

func runGit(t *testing.T, repo string, arguments ...string) string {
	t.Helper()
	command := exec.Command("git", append([]string{"-C", repo}, arguments...)...)
	output, err := command.CombinedOutput()
	if err != nil {
		t.Fatalf("git %v: %v: %s", arguments, err, output)
	}
	return string(output)
}
