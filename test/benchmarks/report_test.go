package spectaculareval

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestStaticReportMatchesGoldenJSONAndMarkdown(t *testing.T) {
	routesBefore := map[string]int{}
	routesAfter := map[string]int{}
	routeDelta := map[string]float64{}
	for _, route := range primaryRoutes {
		routesBefore[route] = 1000
		routesAfter[route] = 500
		routeDelta[route] = 0.5
	}
	report := StaticComparison{
		SchemaVersion: "spectacular.skill-static-comparison.v1",
		GeneratedAt:   time.Date(2026, 8, 22, 12, 0, 0, 0, time.UTC),
		Baseline:      PackageStats{Label: "baseline", Revision: "old", Commit: "aaa", KernelBodyLines: 100, KernelWords: 900, TotalGuidanceWords: 5000, PrimaryRouteWords: routesBefore},
		Candidate:     PackageStats{Label: "candidate", Revision: "new", Commit: "bbb", KernelBodyLines: 40, KernelWords: 360, TotalGuidanceWords: 4500, PrimaryRouteWords: routesAfter},
		Delta:         StaticDelta{KernelBodyLineReduction: 0.6, KernelWordReduction: 0.6, GuidanceWordReduction: 0.1, RouteWordReduction: routeDelta},
		Verdict:       "improved",
		Limitations:   []string{"static only"},
	}
	directory := t.TempDir()
	if err := WriteStaticReport(report, directory); err != nil {
		t.Fatal(err)
	}
	for _, name := range []string{"static.json", "static.md"} {
		actual, err := os.ReadFile(filepath.Join(directory, name))
		if err != nil {
			t.Fatal(err)
		}
		expected, err := os.ReadFile(filepath.Join("testdata", "golden-"+name))
		if err != nil {
			t.Fatal(err)
		}
		if string(actual) != string(expected) {
			t.Fatalf("%s mismatch\n--- actual ---\n%s\n--- expected ---\n%s", name, actual, expected)
		}
	}
}

func TestRunReportWithoutUsageIsInconclusive(t *testing.T) {
	one := 1.0
	report := RunReport{
		SchemaVersion: "spectacular.skill-run-report.v1",
		GeneratedAt:   time.Date(2026, 8, 22, 12, 0, 0, 0, time.UTC),
		BaselineRef:   "old", CandidateRef: "new", Model: "fake", Tier: "smoke",
		Thresholds: Thresholds{MinimumRoutingPassRate: 0.95, MinimumPointerPassRate: 0.95, MinimumTotalContextGain: 0.25},
		Trials: []Trial{
			{ID: "aa-baseline", CaseID: "AA-01", Variant: "baseline", Score: TrialScore{SafetyPassed: true, Overall: &one, Dimensions: passingDimensions()}},
			{ID: "aa-candidate", CaseID: "AA-01", Variant: "candidate", Score: TrialScore{SafetyPassed: true, Overall: &one, Dimensions: passingDimensions()}},
		},
	}
	directory := t.TempDir()
	if err := WriteRunReport(report, directory); err != nil {
		t.Fatal(err)
	}
	data, err := os.ReadFile(filepath.Join(directory, "report.md"))
	if err != nil {
		t.Fatal(err)
	}
	if reportText := string(data); !containsFold(reportText, "inconclusive") || !containsFold(reportText, "0/1") {
		t.Fatalf("report=%s", reportText)
	}
}

func TestRunReportMatchesGoldenJSONAndMarkdown(t *testing.T) {
	one := 1.0
	report := RunReport{
		SchemaVersion: "spectacular.skill-run-report.v1", GeneratedAt: time.Date(2026, 8, 22, 13, 0, 0, 0, time.UTC),
		BaselineRef: "old", BaselineMode: "skill", CandidateRef: "new", CandidateMode: "skill", Model: "eval-model", ReadIsolation: "os-enforced", Tier: "micro", Seed: 7, MinimumRepetitions: 1,
		Thresholds: Thresholds{MaximumSafetyFailures: 0, MinimumRoutingPassRate: 1, MinimumPointerPassRate: 1, MinimumTotalContextGain: 0.25},
		Trials: []Trial{
			{ID: "aa-r01-baseline", CaseID: "AA-01", Variant: "baseline", Revision: "old", Commit: "aaa", Model: "eval-model", Repeat: 1, Order: 1, DurationMS: 100,
				TraceMetrics: TraceMetrics{UsageObserved: true, InputTokens: 100, OutputTokens: 20, ToolCalls: 2, SemanticObserved: true, SemanticEvents: 1},
				TracePath:    "trials/aa-r01-baseline/trace.jsonl", ResultPath: "trials/aa-r01-baseline/result.json", WorkspacePath: "trials/aa-r01-baseline/workspace",
				Score: TrialScore{SafetyPassed: true, Dimensions: passingDimensions(), Overall: &one, Verdict: "pass"}},
			{ID: "aa-r01-candidate", CaseID: "AA-01", Variant: "candidate", Revision: "new", Commit: "bbb", Model: "eval-model", Repeat: 1, Order: 2, DurationMS: 80,
				TraceMetrics: TraceMetrics{UsageObserved: true, InputTokens: 75, OutputTokens: 18, ToolCalls: 1, SemanticObserved: true, SemanticEvents: 1},
				TracePath:    "trials/aa-r01-candidate/trace.jsonl", ResultPath: "trials/aa-r01-candidate/result.json", WorkspacePath: "trials/aa-r01-candidate/workspace",
				Score: TrialScore{SafetyPassed: true, Dimensions: passingDimensions(), Overall: &one, Verdict: "pass"}},
		},
		Limitations: []string{"example limit"},
	}
	directory := t.TempDir()
	if err := WriteRunReport(report, directory); err != nil {
		t.Fatal(err)
	}
	for _, pair := range [][2]string{{"report.json", "golden-run.json"}, {"report.md", "golden-run.md"}} {
		actual, err := os.ReadFile(filepath.Join(directory, pair[0]))
		if err != nil {
			t.Fatal(err)
		}
		if pair[0] == "report.json" {
			actual = projectRunReportJSON(t, actual)
		}
		expected, err := os.ReadFile(filepath.Join("testdata", pair[1]))
		if err != nil {
			t.Fatal(err)
		}
		if string(actual) != string(expected) {
			t.Fatalf("%s mismatch\n--- actual ---\n%s\n--- expected ---\n%s", pair[0], actual, expected)
		}
	}
}

func projectRunReportJSON(t *testing.T, data []byte) []byte {
	t.Helper()
	var report RunReport
	if err := json.Unmarshal(data, &report); err != nil {
		t.Fatal(err)
	}
	type trialRef struct {
		ID        string `json:"id"`
		Trace     string `json:"trace"`
		Result    string `json:"result"`
		Workspace string `json:"workspace"`
	}
	projection := struct {
		Schema    string     `json:"schema"`
		Isolation string     `json:"read_isolation"`
		Minimum   int        `json:"minimum_repetitions"`
		Summary   RunSummary `json:"summary"`
		Trials    []trialRef `json:"trials"`
	}{Schema: report.SchemaVersion, Isolation: report.ReadIsolation, Minimum: report.MinimumRepetitions, Summary: report.Summary}
	for _, trial := range report.Trials {
		projection.Trials = append(projection.Trials, trialRef{ID: trial.ID, Trace: trial.TracePath, Result: trial.ResultPath, Workspace: trial.WorkspacePath})
	}
	encoded, err := json.MarshalIndent(projection, "", "  ")
	if err != nil {
		t.Fatal(err)
	}
	return append(encoded, '\n')
}

func passingDimensions() map[string]DimensionScore {
	result := map[string]DimensionScore{}
	for _, name := range Dimensions {
		one := 1.0
		result[name] = DimensionScore{Applicable: 1, Passed: 1, Rate: &one}
	}
	return result
}
