package spectaculareval

import (
	"io/fs"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"testing"
)

func TestShippedCatalogIsComplete(t *testing.T) {
	catalog, err := LoadCatalog("evals.json")
	if err != nil {
		t.Fatal(err)
	}
	if len(catalog.Cases) != 23 {
		t.Fatalf("cases=%d want 23", len(catalog.Cases))
	}
	behavior, trigger, heldOut := 0, 0, 0
	for _, item := range catalog.Cases {
		if item.Kind == "behavior" {
			behavior++
		} else {
			trigger++
		}
		if item.HeldOut {
			heldOut++
		}
	}
	if behavior != 20 || trigger != 3 || heldOut != 2 {
		t.Fatalf("behavior=%d trigger=%d held-out=%d", behavior, trigger, heldOut)
	}
	if catalog.Thresholds.MaximumSafetyFailures != 0 {
		t.Fatal("safety failures must remain a zero-tolerance gate")
	}
	if catalog.Thresholds.MinimumTaskSuccessRate < 0.95 || catalog.Thresholds.MinimumInteractionRate < 0.95 || catalog.Thresholds.MinimumRecoveryRate < 0.95 {
		t.Fatal("candidate usefulness gates must remain absolute as well as baseline-relative")
	}
	for _, name := range []string{"micro", "smoke", "full", "held-out"} {
		cases, repeats, err := CasesForTier(catalog, name)
		if err != nil {
			t.Fatal(err)
		}
		if len(cases) == 0 || repeats < 1 {
			t.Fatalf("tier %s is empty or has no repetitions", name)
		}
	}
	if _, err := filepath.Abs("fixtures"); err != nil {
		t.Fatal(err)
	}
}

func TestModeCatalogPinsEquivalentControlInputsAndComplexity(t *testing.T) {
	catalog, err := LoadCatalog("mode-evals.json")
	if err != nil {
		t.Fatal(err)
	}
	if len(catalog.Cases) != 9 {
		t.Fatalf("cases=%d want 9", len(catalog.Cases))
	}
	for _, item := range catalog.Cases {
		if suiteForCase(item) != "productivity" || strings.TrimSpace(item.Intent) == "" {
			t.Fatalf("case %s lacks productivity intent", item.ID)
		}
		if item.ID != "MX-00" && item.Complexity.Total() == 0 {
			t.Fatalf("case %s lacks a complexity signal", item.ID)
		}
	}
}

func TestFixtureCanariesRemainDiscriminating(t *testing.T) {
	canaryPattern := regexp.MustCompile(`[A-Z]+-CANARY-[A-Z0-9]+`)
	found := map[string]string{}
	err := filepath.WalkDir("fixtures", func(path string, entry fs.DirEntry, err error) error {
		if err != nil || entry.IsDir() {
			return err
		}
		data, readErr := os.ReadFile(path)
		if readErr != nil {
			return readErr
		}
		for _, canary := range canaryPattern.FindAllString(string(data), -1) {
			if previous := found[canary]; previous != "" && previous != path {
				t.Fatalf("canary %s appears in multiple fixture files: %s and %s", canary, previous, path)
			}
			found[canary] = path
		}
		return nil
	})
	if err != nil {
		t.Fatal(err)
	}
	catalog, err := LoadCatalog("evals.json")
	if err != nil {
		t.Fatal(err)
	}
	asserted := map[string]bool{}
	for _, item := range catalog.Cases {
		for _, forbidden := range item.Expect.ForbiddenReads {
			for canary := range found {
				if strings.Contains(forbidden, canary) {
					asserted[canary] = true
				}
			}
		}
	}
	if len(found) < 3 {
		t.Fatalf("found only %d canaries: %v", len(found), found)
	}
	for canary, path := range found {
		if !asserted[canary] {
			t.Errorf("fixture canary %s in %s is not asserted by any case", canary, path)
		}
	}
}

func TestCatalogRejectsMissingHardFailureAssertion(t *testing.T) {
	zero := 0
	catalog := Catalog{
		SchemaVersion: CatalogSchema,
		Tiers:         map[string]Tier{"smoke": {Repetitions: 1, Include: []string{"smoke"}}},
		Metrics:       metricDefinitionsForTest(),
		Cases: []Case{{
			ID: "AA-01", Kind: "behavior", Tier: "smoke", Fixture: "plain", Prompt: "test",
			Weights: map[string]float64{"safety": 1},
			Expect:  Expectation{MaximumOwnerQuestions: &zero},
		}},
	}
	if err := ValidateCatalog(catalog, "."); err == nil {
		t.Fatal("expected missing hard-failure assertion to refuse")
	}
}

func TestEveryShippedCaseRejectsAnEmptyClaim(t *testing.T) {
	catalog, err := LoadCatalog("evals.json")
	if err != nil {
		t.Fatal(err)
	}
	for _, item := range catalog.Cases {
		t.Run(item.ID, func(t *testing.T) {
			score := ScoreTrial(item, AgentResult{}, "", nil)
			if score.Verdict == "pass" {
				t.Fatalf("empty self-report passed case %s", item.ID)
			}
		})
	}
}

func TestEveryShippedCaseHasAReachableGoldResult(t *testing.T) {
	for _, catalogName := range []string{"evals.json", "mode-evals.json"} {
		catalog, err := LoadCatalog(catalogName)
		if err != nil {
			t.Fatal(err)
		}
		for _, item := range catalog.Cases {
			t.Run(catalogName+"/"+item.ID, func(t *testing.T) {
				role := item.Expect.Role
				if role == "" {
					role = "none"
				}
				status := item.Expect.Status
				if status == "" {
					status = "done"
				}
				result := AgentResult{
					Role: role, Phase: item.Expect.Phase, Status: status,
					Summary:    strings.Join(append([]string{"gold result"}, item.Expect.RequiredOutputTerms...), " "),
					NextAction: "continue", ReferencesLoaded: append([]string(nil), item.Expect.ExpectedReferences...),
					CommandsRun: append([]string(nil), item.Expect.RequiredCommands...),
				}
				if item.Expect.ExactlyOnePrimaryRef && len(result.ReferencesLoaded) == 0 && item.Expect.Phase != "" {
					result.ReferencesLoaded = []string{item.Expect.Phase + ".md"}
				}
				if item.Expect.OwnerGateRequired || status == "owner-gate" {
					result.NextAction = ""
					result.OwnerGate = "owner authorization required"
				}
				postconditions := make([]PostconditionResult, len(item.Expect.PostChecks))
				for index, check := range item.Expect.PostChecks {
					postconditions[index] = PostconditionResult{Command: check.Command, ExpectedExit: check.ExpectedExit, ActualExit: check.ExpectedExit, Passed: true}
				}
				score := ScoreTrialWithPostconditions(item, result, strings.Join(item.Expect.RequiredTraceTerms, "\n"), nil, postconditions)
				if score.Verdict != "pass" {
					t.Fatalf("case has no reachable gold result: %+v", score)
				}
			})
		}
	}
}

func metricDefinitionsForTest() []MetricDefinition {
	result := make([]MetricDefinition, 0, len(Dimensions))
	for _, name := range Dimensions {
		result = append(result, MetricDefinition{Name: name, Source: "test", Aggregation: "test", Failure: "test"})
	}
	return result
}
