package spectaculareval

import (
	"archive/tar"
	"bufio"
	"bytes"
	"crypto/sha256"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
)

var primaryRoutes = []string{"orient", "prepare", "execute", "runtime", "close", "audit"}

func ResolveCommit(repo, revision string) (string, error) {
	command := exec.Command("git", "-C", repo, "rev-parse", "--verify", revision+"^{commit}")
	output, err := command.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("resolve immutable revision %q: %w: %s", revision, err, strings.TrimSpace(string(output)))
	}
	return strings.TrimSpace(string(output)), nil
}

func MaterializeSkill(repo, revision, destination string) (string, error) {
	commit, err := ResolveCommit(repo, revision)
	if err != nil {
		return "", err
	}
	command := exec.Command("git", "-C", repo, "archive", "--format=tar", commit, "skills/spectacular")
	output, err := command.Output()
	if err != nil {
		return "", fmt.Errorf("archive skill at %s: %w", commit, err)
	}
	reader := tar.NewReader(bytes.NewReader(output))
	for {
		header, nextErr := reader.Next()
		if nextErr == io.EOF {
			break
		}
		if nextErr != nil {
			return "", fmt.Errorf("read skill archive: %w", nextErr)
		}
		if header.Typeflag == tar.TypeXGlobalHeader || header.Typeflag == tar.TypeXHeader {
			continue
		}
		clean := filepath.Clean(header.Name)
		prefix := filepath.Join("skills", "spectacular")
		if header.Typeflag == tar.TypeDir && strings.HasPrefix(prefix, clean+string(filepath.Separator)) {
			continue
		}
		if clean != prefix && !strings.HasPrefix(clean, prefix+string(filepath.Separator)) {
			return "", fmt.Errorf("archive path escaped skill root: %s", header.Name)
		}
		relative, relErr := filepath.Rel(prefix, clean)
		if relErr != nil || relative == "." {
			continue
		}
		target := filepath.Join(destination, relative)
		if !strings.HasPrefix(filepath.Clean(target), filepath.Clean(destination)+string(filepath.Separator)) {
			return "", fmt.Errorf("archive target escaped destination: %s", target)
		}
		switch header.Typeflag {
		case tar.TypeDir:
			if err := os.MkdirAll(target, 0o755); err != nil {
				return "", err
			}
		case tar.TypeReg:
			if err := os.MkdirAll(filepath.Dir(target), 0o755); err != nil {
				return "", err
			}
			file, openErr := os.OpenFile(target, os.O_CREATE|os.O_EXCL|os.O_WRONLY, 0o644)
			if openErr != nil {
				return "", openErr
			}
			_, copyErr := io.Copy(file, reader)
			closeErr := file.Close()
			if copyErr != nil {
				return "", copyErr
			}
			if closeErr != nil {
				return "", closeErr
			}
		default:
			return "", fmt.Errorf("unsupported archive entry %s", header.Name)
		}
	}
	return commit, nil
}

func CopyTree(source, destination string) error {
	return filepath.WalkDir(source, func(path string, entry os.DirEntry, err error) error {
		if err != nil {
			return err
		}
		relative, err := filepath.Rel(source, path)
		if err != nil {
			return err
		}
		if relative == "." {
			return os.MkdirAll(destination, 0o755)
		}
		target := filepath.Join(destination, relative)
		if entry.IsDir() {
			return os.MkdirAll(target, 0o755)
		}
		if entry.Type()&os.ModeSymlink != 0 {
			return fmt.Errorf("symlink is not allowed in evaluation input: %s", path)
		}
		data, err := os.ReadFile(path)
		if err != nil {
			return err
		}
		return os.WriteFile(target, data, 0o644)
	})
}

func InspectPackage(label, revision, commit, root string) (PackageStats, error) {
	skillPath := filepath.Join(root, "SKILL.md")
	data, err := os.ReadFile(skillPath)
	if err != nil {
		return PackageStats{}, fmt.Errorf("read SKILL.md: %w", err)
	}
	stats := PackageStats{
		Label:             label,
		Revision:          revision,
		Commit:            commit,
		KernelLines:       lineCount(data),
		KernelBodyLines:   bodyLineCount(data),
		KernelWords:       len(strings.Fields(string(data))),
		KernelBytes:       len(data),
		PrimaryRouteWords: map[string]int{},
	}
	var referenceWords = map[string]int{}
	err = filepath.WalkDir(root, func(path string, entry os.DirEntry, walkErr error) error {
		if walkErr != nil {
			return walkErr
		}
		if entry.IsDir() {
			name := entry.Name()
			if name == "versions" || name == "evals" {
				return filepath.SkipDir
			}
			return nil
		}
		if entry.Name() == ".DS_Store" {
			stats.ValidationFindings = append(stats.ValidationFindings, "package contains .DS_Store: "+path)
			return nil
		}
		if filepath.Ext(path) != ".md" || strings.Contains(path, string(filepath.Separator)+"generated"+string(filepath.Separator)) {
			return nil
		}
		content, readErr := os.ReadFile(path)
		if readErr != nil {
			return readErr
		}
		stats.TotalGuidanceLines += lineCount(content)
		stats.TotalGuidanceWords += len(strings.Fields(string(content)))
		if strings.Contains(path, string(filepath.Separator)+"references"+string(filepath.Separator)) {
			stats.ReferenceFiles++
			base := strings.TrimSuffix(filepath.Base(path), ".md")
			referenceWords[base] = len(strings.Fields(string(content)))
			if !activationWithinFirstFiveLines(content) {
				stats.ValidationFindings = append(stats.ValidationFindings, "reference lacks early Use this when activation: "+filepath.Base(path))
			}
		}
		return nil
	})
	if err != nil {
		return PackageStats{}, err
	}
	for _, route := range primaryRoutes {
		if words, ok := referenceWords[route]; ok {
			stats.PrimaryRouteWords[route] = stats.KernelWords + words
		} else {
			stats.ValidationFindings = append(stats.ValidationFindings, "missing primary reference: "+route+".md")
		}
	}
	stats.ValidationFindings = append(stats.ValidationFindings, validateMarkdownLinks(root)...)
	stats.ValidationFindings = append(stats.ValidationFindings, validateReferenceReachability(root)...)
	sort.Strings(stats.ValidationFindings)
	return stats, nil
}

func ComparePackages(baseline, candidate PackageStats) StaticComparison {
	delta := StaticDelta{
		KernelBodyLineReduction: reduction(baseline.KernelBodyLines, candidate.KernelBodyLines),
		KernelWordReduction:     reduction(baseline.KernelWords, candidate.KernelWords),
		GuidanceWordReduction:   reduction(baseline.TotalGuidanceWords, candidate.TotalGuidanceWords),
		RouteWordReduction:      map[string]float64{},
	}
	for _, route := range primaryRoutes {
		delta.RouteWordReduction[route] = reduction(baseline.PrimaryRouteWords[route], candidate.PrimaryRouteWords[route])
	}
	verdict := "improved"
	if len(candidate.ValidationFindings) > len(baseline.ValidationFindings) || delta.KernelWordReduction < 0 {
		verdict = "regression"
	} else if len(candidate.ValidationFindings) > 0 {
		verdict = "improved-with-findings"
	}
	return StaticComparison{
		SchemaVersion: "spectacular.skill-static-comparison.v1",
		Baseline:      baseline,
		Candidate:     candidate,
		Delta:         delta,
		Verdict:       verdict,
		Limitations: []string{
			"Static word and line counts estimate available context; only model traces establish what was actually loaded.",
			"Static validation does not establish task success, authority compliance, or invocation accuracy.",
		},
	}
}

// ApplyStaticThresholds turns the measurement contract into deterministic
// static gates. It does not promote a static comparison to a behavioral pass.
func ApplyStaticThresholds(report *StaticComparison, thresholds Thresholds) {
	if thresholds.MaximumKernelBodyLines > 0 && report.Candidate.KernelBodyLines > thresholds.MaximumKernelBodyLines {
		report.GateFailures = append(report.GateFailures, fmt.Sprintf("candidate kernel body lines=%d, maximum=%d", report.Candidate.KernelBodyLines, thresholds.MaximumKernelBodyLines))
	}
	if report.Delta.KernelWordReduction < thresholds.MinimumInitialContextGain {
		report.GateFailures = append(report.GateFailures, fmt.Sprintf("kernel word reduction %.3f below %.3f", report.Delta.KernelWordReduction, thresholds.MinimumInitialContextGain))
	}
	if len(report.GateFailures) > 0 {
		report.Verdict = "regression"
	}
	sort.Strings(report.GateFailures)
}

func reduction(before, after int) float64 {
	if before == 0 {
		return 0
	}
	return float64(before-after) / float64(before)
}

func lineCount(data []byte) int {
	if len(data) == 0 {
		return 0
	}
	count := bytes.Count(data, []byte("\n"))
	if data[len(data)-1] != '\n' {
		count++
	}
	return count
}

func bodyLineCount(data []byte) int {
	scanner := bufio.NewScanner(bytes.NewReader(data))
	delimiters := 0
	lines := 0
	for scanner.Scan() {
		if scanner.Text() == "---" && delimiters < 2 {
			delimiters++
			continue
		}
		if delimiters >= 2 {
			lines++
		}
	}
	return lines
}

func activationWithinFirstFiveLines(data []byte) bool {
	scanner := bufio.NewScanner(bytes.NewReader(data))
	for line := 0; line < 5 && scanner.Scan(); line++ {
		if strings.HasPrefix(strings.TrimSpace(scanner.Text()), "Use this when:") {
			return true
		}
	}
	return false
}

var markdownLink = regexp.MustCompile(`\[[^\]]+\]\(([^)]+)\)`)

func validateMarkdownLinks(root string) []string {
	var findings []string
	_ = filepath.WalkDir(root, func(path string, entry os.DirEntry, err error) error {
		if err != nil || entry.IsDir() || filepath.Ext(path) != ".md" || strings.Contains(path, string(filepath.Separator)+"versions"+string(filepath.Separator)) {
			return nil
		}
		data, readErr := os.ReadFile(path)
		if readErr != nil {
			return nil
		}
		for _, match := range markdownLink.FindAllStringSubmatch(string(data), -1) {
			target := strings.SplitN(match[1], "#", 2)[0]
			if target == "" || strings.Contains(target, "://") {
				continue
			}
			resolved := filepath.Clean(filepath.Join(filepath.Dir(path), filepath.FromSlash(target)))
			if _, statErr := os.Stat(resolved); statErr != nil {
				findings = append(findings, fmt.Sprintf("broken link %s -> %s", path, match[1]))
			}
		}
		return nil
	})
	return findings
}

func validateReferenceReachability(root string) []string {
	root = filepath.Clean(root)
	queue := []string{filepath.Join(root, "SKILL.md")}
	seen := map[string]bool{}
	for len(queue) > 0 {
		path := queue[0]
		queue = queue[1:]
		if seen[path] {
			continue
		}
		seen[path] = true
		data, err := os.ReadFile(path)
		if err != nil {
			continue
		}
		for _, match := range markdownLink.FindAllStringSubmatch(string(data), -1) {
			target := strings.SplitN(match[1], "#", 2)[0]
			if target == "" || strings.Contains(target, "://") {
				continue
			}
			resolved := filepath.Clean(filepath.Join(filepath.Dir(path), filepath.FromSlash(target)))
			if filepath.Ext(resolved) != ".md" || !strings.HasPrefix(resolved, root+string(filepath.Separator)) {
				continue
			}
			if _, err := os.Stat(resolved); err == nil {
				queue = append(queue, resolved)
			}
		}
	}
	var findings []string
	references := filepath.Join(root, "references")
	_ = filepath.WalkDir(references, func(path string, entry os.DirEntry, err error) error {
		if err != nil || entry.IsDir() || filepath.Ext(path) != ".md" {
			return nil
		}
		if !seen[path] {
			relative, _ := filepath.Rel(root, path)
			findings = append(findings, "orphan reference unreachable from SKILL.md: "+filepath.ToSlash(relative))
		}
		return nil
	})
	return findings
}

func SnapshotTree(root string) (map[string]string, error) {
	result := map[string]string{}
	err := filepath.WalkDir(root, func(path string, entry os.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if entry.IsDir() {
			if entry.Name() == ".git" || entry.Name() == ".agents" {
				return filepath.SkipDir
			}
			return nil
		}
		data, readErr := os.ReadFile(path)
		if readErr != nil {
			return readErr
		}
		relative, relErr := filepath.Rel(root, path)
		if relErr != nil {
			return relErr
		}
		digest := sha256.Sum256(data)
		result[filepath.ToSlash(relative)] = fmt.Sprintf("%x", digest[:])
		return nil
	})
	return result, err
}

func ChangedPaths(before, after map[string]string) []string {
	seen := map[string]bool{}
	var paths []string
	for path, digest := range before {
		seen[path] = true
		if after[path] != digest {
			paths = append(paths, path)
		}
	}
	for path := range after {
		if !seen[path] {
			paths = append(paths, path)
		}
	}
	sort.Strings(paths)
	return paths
}
