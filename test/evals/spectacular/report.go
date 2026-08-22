package spectaculareval

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"
)

func WriteStaticReport(report StaticComparison, directory string) error {
	if report.GeneratedAt.IsZero() {
		report.GeneratedAt = time.Now().UTC()
	}
	if err := os.MkdirAll(directory, 0o755); err != nil {
		return err
	}
	if err := writeJSON(filepath.Join(directory, "static.json"), report); err != nil {
		return err
	}
	return os.WriteFile(filepath.Join(directory, "static.md"), []byte(renderStatic(report)), 0o644)
}

func WriteRunReport(report RunReport, directory string) error {
	if report.GeneratedAt.IsZero() {
		report.GeneratedAt = time.Now().UTC()
	}
	Summarize(&report)
	if err := os.MkdirAll(directory, 0o755); err != nil {
		return err
	}
	if err := writeJSON(filepath.Join(directory, "report.json"), report); err != nil {
		return err
	}
	return os.WriteFile(filepath.Join(directory, "report.md"), []byte(renderRun(report)), 0o644)
}

func writeJSON(path string, value any) error {
	data, err := json.MarshalIndent(value, "", "  ")
	if err != nil {
		return err
	}
	data = append(data, '\n')
	return os.WriteFile(path, data, 0o644)
}

func renderStatic(report StaticComparison) string {
	var builder strings.Builder
	fmt.Fprintf(&builder, "# Spectacular static skill comparison\n\n")
	fmt.Fprintf(&builder, "Verdict: **%s**\n\n", report.Verdict)
	fmt.Fprintf(&builder, "| Measure | Baseline | Candidate | Reduction |\n|---|---:|---:|---:|\n")
	fmt.Fprintf(&builder, "| Kernel body lines | %d | %d | %.1f%% |\n", report.Baseline.KernelBodyLines, report.Candidate.KernelBodyLines, 100*report.Delta.KernelBodyLineReduction)
	fmt.Fprintf(&builder, "| Kernel words | %d | %d | %.1f%% |\n", report.Baseline.KernelWords, report.Candidate.KernelWords, 100*report.Delta.KernelWordReduction)
	fmt.Fprintf(&builder, "| Total guidance words | %d | %d | %.1f%% |\n", report.Baseline.TotalGuidanceWords, report.Candidate.TotalGuidanceWords, 100*report.Delta.GuidanceWordReduction)
	fmt.Fprintf(&builder, "\n## Primary route footprints\n\n| Route | Baseline words | Candidate words | Reduction |\n|---|---:|---:|---:|\n")
	for _, route := range primaryRoutes {
		fmt.Fprintf(&builder, "| `%s` | %d | %d | %.1f%% |\n", route, report.Baseline.PrimaryRouteWords[route], report.Candidate.PrimaryRouteWords[route], 100*report.Delta.RouteWordReduction[route])
	}
	fmt.Fprintf(&builder, "\n## Validation findings\n\n- Baseline: %d\n- Candidate: %d\n", len(report.Baseline.ValidationFindings), len(report.Candidate.ValidationFindings))
	if len(report.GateFailures) > 0 {
		fmt.Fprintf(&builder, "\n## Failed gates\n\n")
		for _, finding := range report.GateFailures {
			fmt.Fprintf(&builder, "- %s\n", finding)
		}
	}
	for _, finding := range report.Candidate.ValidationFindings {
		fmt.Fprintf(&builder, "  - %s\n", finding)
	}
	fmt.Fprintf(&builder, "\n## Limits\n\n")
	for _, limit := range report.Limitations {
		fmt.Fprintf(&builder, "- %s\n", limit)
	}
	return builder.String()
}

func renderRun(report RunReport) string {
	var builder strings.Builder
	fmt.Fprintf(&builder, "# Spectacular paired behavior benchmark\n\n")
	fmt.Fprintf(&builder, "Verdict: **%s**<br>\nBaseline: `%s`<br>\nCandidate: `%s`<br>\nModel: `%s`<br>\nRead isolation: `%s`<br>\nTier: `%s`<br>\nMinimum repetitions: `%d`\n\n", report.Summary.Verdict, report.BaselineRef, report.CandidateRef, report.Model, report.ReadIsolation, report.Tier, report.MinimumRepetitions)
	fmt.Fprintf(&builder, "## Dimension rates\n\n| Dimension | Baseline | Candidate |\n|---|---:|---:|\n")
	for _, dimension := range Dimensions {
		fmt.Fprintf(&builder, "| %s | %.1f%% | %.1f%% |\n", dimension, 100*report.Summary.DimensionRates["baseline"][dimension], 100*report.Summary.DimensionRates["candidate"][dimension])
	}
	fmt.Fprintf(&builder, "\nSafety failures: baseline `%d`, candidate `%d`.\n", report.Summary.SafetyFailures["baseline"], report.Summary.SafetyFailures["candidate"])
	pairing := report.Summary.Pairing
	fmt.Fprintf(&builder, "\n## Paired outcomes\n\nPairs `%d`; candidate wins `%d`; candidate losses `%d`; both pass `%d`; both fail `%d`; discordant rate `%.1f%%`.", pairing.Pairs, pairing.CandidateWins, pairing.CandidateLosses, pairing.BothPass, pairing.BothFail, 100*pairing.DiscordantRate)
	if pairing.ExactSignPValue != nil {
		fmt.Fprintf(&builder, " Exact two-sided sign-test p-value `%.4f`.", *pairing.ExactSignPValue)
	}
	fmt.Fprintln(&builder)
	if len(pairing.UnstableCasePairs) > 0 {
		fmt.Fprintf(&builder, "\nUnstable case/variant outcomes: `%s`.\n", strings.Join(pairing.UnstableCasePairs, "`, `"))
	}
	fmt.Fprintf(&builder, "\n## Observed cost\n\n| Variant | Usage coverage | Median input tokens | Median cached tokens | Median output tokens | Median tool calls | Median duration |\n|---|---:|---:|---:|---:|---:|---:|\n")
	for _, variant := range []string{"baseline", "candidate"} {
		cost := report.Summary.ObservedCost[variant]
		fmt.Fprintf(&builder, "| %s | %d/%d | %.0f | %.0f | %.0f | %.1f | %.0fms |\n", variant, cost.TrialsWithUsage, cost.TotalTrials, cost.MedianInputTokens, cost.MedianCachedTokens, cost.MedianOutputTokens, cost.MedianToolCalls, cost.MedianDurationMillis)
	}
	if len(report.Summary.GateFailures) > 0 {
		fmt.Fprintf(&builder, "\n## Failed gates\n\n")
		for _, finding := range report.Summary.GateFailures {
			fmt.Fprintf(&builder, "- %s\n", finding)
		}
	}
	if len(report.Summary.PerCaseRegressions) > 0 {
		fmt.Fprintf(&builder, "\n## Regressions\n\n")
		for _, finding := range report.Summary.PerCaseRegressions {
			fmt.Fprintf(&builder, "- %s\n", finding)
		}
	}
	if len(report.Summary.InsufficientEvidence) > 0 {
		fmt.Fprintf(&builder, "\n## Insufficient evidence\n\n")
		for _, finding := range report.Summary.InsufficientEvidence {
			fmt.Fprintf(&builder, "- %s\n", finding)
		}
	}
	fmt.Fprintf(&builder, "\n## Trials\n\n| Trial | Case | Variant | Repeat | Verdict | Overall | Duration | Raw artifacts |\n|---|---|---|---:|---|---:|---:|---|\n")
	trials := append([]Trial(nil), report.Trials...)
	sort.Slice(trials, func(i, j int) bool { return trials[i].ID < trials[j].ID })
	for _, trial := range trials {
		overall := "n/a"
		if trial.Score.Overall != nil {
			overall = fmt.Sprintf("%.3f", *trial.Score.Overall)
		}
		artifacts := fmt.Sprintf("[result](%s) · [trace](%s) · [workspace](%s)", trial.ResultPath, trial.TracePath, trial.WorkspacePath)
		fmt.Fprintf(&builder, "| %s | %s | %s | %d | %s | %s | %dms | %s |\n", trial.ID, trial.CaseID, trial.Variant, trial.Repeat, trial.Score.Verdict, overall, trial.DurationMS, artifacts)
	}
	if len(report.Limitations) > 0 {
		fmt.Fprintf(&builder, "\n## Limits\n\n")
		for _, limit := range report.Limitations {
			fmt.Fprintf(&builder, "- %s\n", limit)
		}
	}
	return builder.String()
}
