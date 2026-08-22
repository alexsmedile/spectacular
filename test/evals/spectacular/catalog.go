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
		if _, ok := c.Tiers[item.Tier]; !ok {
			problems = append(problems, prefix+"unknown tier "+item.Tier)
		}
		if strings.TrimSpace(item.Prompt) == "" {
			problems = append(problems, prefix+"prompt is required")
		}
		fixture := filepath.Join(root, "fixtures", filepath.Clean(item.Fixture))
		if info, err := os.Stat(fixture); err != nil || !info.IsDir() {
			problems = append(problems, prefix+"fixture directory does not exist: "+item.Fixture)
		}
		if len(item.Weights) == 0 {
			problems = append(problems, prefix+"dimension weights are required")
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
	}
	if len(problems) > 0 {
		sort.Strings(problems)
		return fmt.Errorf("invalid benchmark catalog:\n- %s", strings.Join(problems, "\n- "))
	}
	return nil
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
