package spectaculareval

import (
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

func passingDimensions() map[string]DimensionScore {
	result := map[string]DimensionScore{}
	for _, name := range Dimensions {
		one := 1.0
		result[name] = DimensionScore{Applicable: 1, Passed: 1, Rate: &one}
	}
	return result
}
