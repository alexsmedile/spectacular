package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"time"

	bench "github.com/alexsmedile/spectacular/v2/test/evals/spectacular"
)

func main() {
	if len(os.Args) < 2 {
		usage()
		os.Exit(2)
	}
	var err error
	switch os.Args[1] {
	case "validate":
		err = validate(os.Args[2:])
	case "static":
		err = static(os.Args[2:])
	case "run":
		err = run(os.Args[2:])
	default:
		usage()
		os.Exit(2)
	}
	if err != nil {
		fmt.Fprintln(os.Stderr, "spectacular eval:", err)
		os.Exit(1)
	}
}

func usage() {
	fmt.Fprintln(os.Stderr, "usage: go run ./test/evals/spectacular/cmd/bench <validate|static|run> [flags]")
}

func validate(args []string) error {
	set := flag.NewFlagSet("validate", flag.ContinueOnError)
	catalogPath := set.String("catalog", "test/evals/spectacular/evals.json", "benchmark catalog")
	if err := set.Parse(args); err != nil {
		return err
	}
	catalog, err := bench.LoadCatalog(*catalogPath)
	if err != nil {
		return err
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
	fmt.Printf("catalog valid: %d cases (%d behavior, %d trigger, %d held-out)\n", len(catalog.Cases), behavior, trigger, heldOut)
	return nil
}

func static(args []string) error {
	set := flag.NewFlagSet("static", flag.ContinueOnError)
	repo := set.String("repo", ".", "Git repository")
	baselineRef := set.String("baseline", "14158f9", "immutable baseline revision")
	candidateRef := set.String("candidate", "", "immutable candidate revision")
	candidateDir := set.String("candidate-dir", "", "candidate skill directory for provisional static inspection")
	output := set.String("out", "test/evals/spectacular/reports/static", "report directory")
	if err := set.Parse(args); err != nil {
		return err
	}
	if *candidateRef == "" && *candidateDir == "" {
		return fmt.Errorf("candidate or candidate-dir is required")
	}
	temporary, err := os.MkdirTemp("", "spectacular-static-")
	if err != nil {
		return err
	}
	defer os.RemoveAll(temporary)
	baselineRoot := filepath.Join(temporary, "baseline")
	baselineCommit, err := bench.MaterializeSkill(*repo, *baselineRef, baselineRoot)
	if err != nil {
		return err
	}
	baseline, err := bench.InspectPackage("baseline", *baselineRef, baselineCommit, baselineRoot)
	if err != nil {
		return err
	}
	var candidate bench.PackageStats
	if *candidateRef != "" {
		candidateRoot := filepath.Join(temporary, "candidate")
		candidateCommit, err := bench.MaterializeSkill(*repo, *candidateRef, candidateRoot)
		if err != nil {
			return err
		}
		candidate, err = bench.InspectPackage("candidate", *candidateRef, candidateCommit, candidateRoot)
		if err != nil {
			return err
		}
	} else {
		candidate, err = bench.InspectPackage("candidate", "working-tree", "", *candidateDir)
		if err != nil {
			return err
		}
	}
	report := bench.ComparePackages(baseline, candidate)
	report.GeneratedAt = time.Now().UTC()
	if *candidateDir != "" {
		report.Verdict = "provisional"
		report.Limitations = append(report.Limitations, "Candidate is a mutable directory; behavioral A/B trials require an immutable candidate revision.")
	}
	if err := bench.WriteStaticReport(report, *output); err != nil {
		return err
	}
	fmt.Printf("static comparison: %s (%s)\n", report.Verdict, filepath.Join(*output, "static.md"))
	return nil
}

func run(args []string) error {
	set := flag.NewFlagSet("run", flag.ContinueOnError)
	repo := set.String("repo", ".", "Git repository")
	catalogPath := set.String("catalog", "test/evals/spectacular/evals.json", "benchmark catalog")
	schemaPath := set.String("schema", "test/evals/spectacular/agent-result.schema.json", "agent result schema")
	baseline := set.String("baseline", "14158f9", "immutable baseline revision")
	candidate := set.String("candidate", "", "immutable candidate revision")
	tier := set.String("tier", "smoke", "smoke, full, or held-out")
	repeats := set.Int("repeats", 0, "override tier repetitions")
	seed := set.Int64("seed", 1, "pair randomization seed")
	model := set.String("model", "", "model identifier")
	adapter := set.String("adapter", "test/evals/spectacular/scripts/codex-adapter.sh", "adapter executable")
	output := set.String("out", "", "new output directory")
	if err := set.Parse(args); err != nil {
		return err
	}
	if *candidate == "" || *model == "" || *output == "" {
		return fmt.Errorf("candidate, model, and out are required")
	}
	report, err := bench.RunPaired(bench.RunConfig{
		Repo: *repo, CatalogPath: *catalogPath, SchemaPath: *schemaPath,
		BaselineRef: *baseline, CandidateRef: *candidate, Tier: *tier,
		Repeats: *repeats, Seed: *seed, Model: *model, Adapter: *adapter, OutputDir: *output,
	})
	if err != nil {
		return err
	}
	if err := bench.WriteRunReport(report, *output); err != nil {
		return err
	}
	fmt.Printf("paired benchmark: %s (%s)\n", report.Summary.Verdict, filepath.Join(*output, "report.md"))
	return nil
}
