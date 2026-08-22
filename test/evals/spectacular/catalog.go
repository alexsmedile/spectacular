package spectaculareval

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
)

var caseIDPattern = regexp.MustCompile(`^[A-Z]{2,4}-[0-9]{2}$`)

func LoadCatalog(path string) (Catalog, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return Catalog{}, fmt.Errorf("read catalog: %w", err)
	}
	var catalog Catalog
	decoder := json.NewDecoder(strings.NewReader(string(data)))
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(&catalog); err != nil {
		return Catalog{}, fmt.Errorf("decode catalog: %w", err)
	}
	if err := ValidateCatalog(catalog, filepath.Dir(path)); err != nil {
		return Catalog{}, err
	}
	return catalog, nil
}

func ValidateCatalog(c Catalog, root string) error {
	var problems []string
	if c.SchemaVersion != CatalogSchema {
		problems = append(problems, "schema_version must be "+CatalogSchema)
	}
	if len(c.Cases) == 0 {
		problems = append(problems, "at least one case is required")
	}
	metricSeen := map[string]bool{}
	for _, metric := range c.Metrics {
		if metric.Name == "" || metric.Source == "" || metric.Aggregation == "" || metric.Failure == "" {
			problems = append(problems, "every metric needs name, source, aggregation, and failure")
		}
		metricSeen[metric.Name] = true
	}
	for _, name := range Dimensions {
		if !metricSeen[name] {
			problems = append(problems, "missing metric definition: "+name)
		}
	}
	ids := map[string]bool{}
	for _, item := range c.Cases {
		prefix := item.ID + ": "
		if !caseIDPattern.MatchString(item.ID) {
			problems = append(problems, prefix+"id must match AA-00 style")
		}
		if ids[item.ID] {
			problems = append(problems, prefix+"duplicate id")
		}
		ids[item.ID] = true
		if item.Kind != "behavior" && item.Kind != "trigger" {
			problems = append(problems, prefix+"kind must be behavior or trigger")
		}
		if item.Suite != "" && item.Suite != "conformance" && item.Suite != "productivity" && item.Suite != "trigger" {
			problems = append(problems, prefix+"suite must be conformance, productivity, or trigger")
		}
		for label, value := range map[string]int{"scope": item.Complexity.Scope, "ambiguity": item.Complexity.Ambiguity, "consequence": item.Complexity.Consequence, "continuity": item.Complexity.Continuity} {
			if value < 0 || value > 3 {
				problems = append(problems, fmt.Sprintf("%scomplexity.%s must be between 0 and 3", prefix, label))
			}
		}
		if _, ok := c.Tiers[item.Tier]; !ok {
			problems = append(problems, prefix+"unknown tier "+item.Tier)
		}
		if strings.TrimSpace(item.Prompt) == "" {
			problems = append(problems, prefix+"prompt is required")
		}
		if containsFold(item.Prompt, "without changing files") && strings.TrimSpace(item.Intent) == "" {
			problems = append(problems, prefix+"restraint cases need intent explaining why granted authority is not exercised")
		}
		if item.Expect.Status != "" && stringInFold(item.Expect.ForbiddenStatuses, item.Expect.Status) {
			problems = append(problems, prefix+"expected status is also forbidden")
		}
		if item.Expect.Role != "" && stringInFold(item.Expect.ForbiddenRoles, item.Expect.Role) {
			problems = append(problems, prefix+"expected role is also forbidden")
		}
		if item.Expect.OwnerGateRequired && item.Expect.Status != "owner-gate" {
			problems = append(problems, prefix+"owner_gate_required needs status owner-gate")
		}
		primaryExpected := 0
		for _, reference := range item.Expect.ExpectedReferences {
			if primaryReferenceNames[filepath.Base(reference)] {
				primaryExpected++
			}
		}
		if item.Expect.ExactlyOnePrimaryRef && primaryExpected > 1 {
			problems = append(problems, prefix+"exactly_one_primary_reference conflicts with multiple expected primary references")
		}
		if len(item.Expect.AllowedChangedPaths) > 0 && stringInFold(item.Expect.ForbiddenChangedPaths, "**") {
			problems = append(problems, prefix+"allowed_changed_paths conflicts with forbidden_changed_paths **")
		}
		for _, term := range append(append([]string(nil), item.Expect.ForbiddenAnyTerms...), item.Expect.ForbiddenTraceTerms...) {
			if containsFold(item.Prompt, term) {
				problems = append(problems, prefix+"forbidden term is already present in the prompt: "+term)
			}
		}
		fixture := filepath.Join(root, "fixtures", filepath.Clean(item.Fixture))
		if info, err := os.Stat(fixture); err != nil || !info.IsDir() {
			problems = append(problems, prefix+"fixture directory does not exist: "+item.Fixture)
		} else if item.Suite == "productivity" {
			for _, required := range []string{"TASK.md", filepath.Join(".spectacular", "PROJECT.md")} {
				if info, err := os.Stat(filepath.Join(fixture, required)); err != nil || info.IsDir() {
					problems = append(problems, prefix+"productivity fixture needs equivalent control input: "+filepath.ToSlash(required))
				}
			}
		}
		if len(item.Weights) == 0 {
			problems = append(problems, prefix+"dimension weights are required")
		}
		for index, check := range item.Expect.PostChecks {
			if len(check.Command) == 0 || strings.TrimSpace(check.Command[0]) == "" {
				problems = append(problems, fmt.Sprintf("%spost_checks[%d] needs a command", prefix, index))
			}
		}
		weightTotal := 0.0
		for dimension, weight := range item.Weights {
			if !metricSeen[dimension] {
				problems = append(problems, prefix+"unknown weight dimension "+dimension)
			}
			if weight < 0 {
				problems = append(problems, prefix+"negative weight for "+dimension)
			}
			weightTotal += weight
		}
		if weightTotal <= 0 {
			problems = append(problems, prefix+"weights must sum above zero")
		}
		if len(item.Expect.ForbiddenAnyTerms)+len(item.Expect.ForbiddenReads)+len(item.Expect.ForbiddenChangedPaths)+len(item.Expect.ForbiddenRoles)+len(item.Expect.ForbiddenStatuses) == 0 {
			problems = append(problems, prefix+"must define at least one hard-failure assertion")
		}
		if item.HeldOut && (item.Tier == "micro" || item.Tier == "smoke") {
			problems = append(problems, prefix+"held-out cases cannot be micro or smoke cases")
		}
	}
	for name, tier := range c.Tiers {
		if tier.Repetitions < 1 {
			problems = append(problems, "tier "+name+": repetitions must be positive")
		}
		for _, included := range tier.Include {
			if included != "micro" && included != "smoke" && included != "full" && included != "held-out" {
				problems = append(problems, "tier "+name+": unknown include class "+included)
			}
		}
		selected := 0
		for _, item := range c.Cases {
			class := item.Tier
			if item.HeldOut {
				class = "held-out"
			}
			for _, included := range tier.Include {
				if class == included {
					selected++
					break
				}
			}
		}
		if selected == 0 {
			problems = append(problems, "tier "+name+": selects no cases")
		}
	}
	if len(problems) > 0 {
		sort.Strings(problems)
		return fmt.Errorf("invalid benchmark catalog:\n- %s", strings.Join(problems, "\n- "))
	}
	return nil
}

func suiteForCase(item Case) string {
	if item.Suite != "" {
		return item.Suite
	}
	if item.Kind == "trigger" {
		return "trigger"
	}
	return "conformance"
}

func stringInFold(values []string, target string) bool {
	for _, value := range values {
		if strings.EqualFold(value, target) {
			return true
		}
	}
	return false
}

func CasesForTier(c Catalog, tierName string) ([]Case, int, error) {
	tier, ok := c.Tiers[tierName]
	if !ok {
		return nil, 0, fmt.Errorf("unknown tier %q", tierName)
	}
	classes := map[string]bool{}
	for _, class := range tier.Include {
		classes[class] = true
	}
	var selected []Case
	for _, item := range c.Cases {
		class := item.Tier
		if item.HeldOut {
			class = "held-out"
		}
		if classes[class] {
			selected = append(selected, item)
		}
	}
	return selected, tier.Repetitions, nil
}
