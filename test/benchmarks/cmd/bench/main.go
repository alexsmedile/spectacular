package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"time"

	bench "github.com/alexsmedile/spectacular/v2/test/benchmarks"
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
	case "adapter-check":
		err = adapterCheck(os.Args[2:])
	case "plan":
		err = plan(os.Args[2:])
	case "static":
		err = static(os.Args[2:])
	case "run":
		err = run(os.Args[2:])
	default:
		usage()
		os.Exit(2)
	}
	if err != nil {
		fmt.Fprintln(os.Stderr, "spectacular bench:", err)
		os.Exit(1)
	}
}

func usage() {
	fmt.Fprintln(os.Stderr, "usage: go run ./test/benchmarks/cmd/bench <validate|adapter-check|plan|static|run> [flags]")
}

func adapterCheck(args []string) error {
	set := flag.NewFlagSet("adapter-check", flag.ContinueOnError)
	tracePath := set.String("trace", "", "captured raw or normalized adapter trace")
	allowZeroTools := set.Bool("allow-zero-tools", false, "permit a certified trace with no tool calls")
	if err := set.Parse(args); err != nil {
		return err
	}
	if *tracePath == "" {
		return fmt.Errorf("trace is required")
	}
	data, err := os.ReadFile(*tracePath)
	if err != nil {
		return err
	}
	result := bench.CertifyTrace(string(data), !*allowZeroTools)
	encoded, err := json.MarshalIndent(result, "", "  ")
	if err != nil {
		return err
	}
	fmt.Println(string(encoded))
	if !result.Valid {
		return fmt.Errorf("adapter trace is not certified")
	}
	return nil
}

func validate(args []string) error {
	set := flag.NewFlagSet("validate", flag.ContinueOnError)
	catalogPath := set.String("catalog", "test/benchmarks/evals.json", "benchmark catalog")
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

func plan(args []string) error {
	set := flag.NewFlagSet("plan", flag.ContinueOnError)
	catalogPath := set.String("catalog", "test/benchmarks/evals.json", "benchmark catalog")
	tier := set.String("tier", "micro", "micro, smoke, full, or held-out")
	repeats := set.Int("repeats", 0, "override tier repetitions")
	if err := set.Parse(args); err != nil {
		return err
	}
	catalog, err := bench.LoadCatalog(*catalogPath)
	if err != nil {
		return err
	}
	cases, defaultRepeats, err := bench.CasesForTier(catalog, *tier)
	if err != nil {
		return err
	}
	if *repeats == 0 {
		*repeats = defaultRepeats
	}
	if *repeats < 1 {
		return fmt.Errorf("repeats must be positive")
	}
	behavior, trigger, heldOut := 0, 0, 0
	ids := make([]string, 0, len(cases))
	for _, item := range cases {
		ids = append(ids, item.ID)
		if item.Kind == "behavior" {
			behavior++
		} else {
			trigger++
		}
		if item.HeldOut {
			heldOut++
		}
	}
	modelCalls := len(cases) * *repeats * 2
	fmt.Printf("tier %s: %d cases x %d repetitions x 2 variants = %d model calls\n", *tier, len(cases), *repeats, modelCalls)
	fmt.Printf("composition: %d behavior, %d trigger, %d held-out\n", behavior, trigger, heldOut)
	fmt.Printf("cases: %v\n", ids)
	return nil
}

func static(args []string) error {
	set := flag.NewFlagSet("static", flag.ContinueOnError)
	repo := set.String("repo", ".", "Git repository")
	catalogPath := set.String("catalog", "test/benchmarks/evals.json", "benchmark measurement contract")
	baselineRef := set.String("baseline", "14158f9", "immutable baseline revision")
	candidateRef := set.String("candidate", "", "immutable candidate revision")
	candidateDir := set.String("candidate-dir", "", "candidate skill directory for provisional static inspection")
	output := set.String("out", "test/benchmarks/reports/static", "report directory")
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
	catalog, err := bench.LoadCatalog(*catalogPath)
	if err != nil {
		return err
	}
	bench.ApplyStaticThresholds(&report, catalog.Thresholds)
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
	catalogPath := set.String("catalog", "test/benchmarks/evals.json", "benchmark catalog")
	schemaPath := set.String("schema", "test/benchmarks/agent-result.schema.json", "agent result schema")
	baseline := set.String("baseline", "14158f9", "immutable baseline revision")
	baselineMode := set.String("baseline-mode", "skill", "skill, workspace-only, or native-direct")
	candidate := set.String("candidate", "", "immutable candidate revision")
	candidateMode := set.String("candidate-mode", "skill", "skill, workspace-only, or native-direct")
	tier := set.String("tier", "smoke", "micro, smoke, full, or held-out")
	repeats := set.Int("repeats", 0, "override tier repetitions")
	seed := set.Int64("seed", 1, "pair randomization seed")
	model := set.String("model", "", "model identifier")
	adapter := set.String("adapter", "test/benchmarks/adapters/codex-adapter.sh", "adapter executable")
	spectacularCLI := set.String("spectacular-cli", "", "absolute path to the repository-built Spectacular CLI pinned for every trial")
	output := set.String("out", "", "new output directory")
	allowHeldOut := set.Bool("allow-held-out", false, "explicitly authorize a frozen held-out evidence run")
	readIsolation := set.String("read-isolation", "artifact-only", "artifact-only or os-enforced (external adapter attestation)")
	maxCalls := set.Int("max-calls", 12, "refuse runs requiring more model calls; set deliberately for larger tiers")
	parallel := set.Int("parallel", 1, "concurrency level for parallel trial execution")
	trialTimeout := set.Duration("trial-timeout", 10*time.Minute, "maximum duration for one model trial; zero disables")
	requireCertifiedTelemetry := set.Bool("require-certified-telemetry", true, "stop after the first trial unless usage and host semantic telemetry are observed")
	var adapterArgs stringListFlag
	set.Var(&adapterArgs, "adapter-arg", "argument forwarded to the adapter; repeatable")
	if err := set.Parse(args); err != nil {
		return err
	}
	if *candidate == "" || *model == "" || *output == "" {
		return fmt.Errorf("candidate, model, and out are required")
	}
	report, err := bench.RunPaired(bench.RunConfig{
		Repo: *repo, CatalogPath: *catalogPath, SchemaPath: *schemaPath,
		BaselineRef: *baseline, BaselineMode: *baselineMode, CandidateRef: *candidate, CandidateMode: *candidateMode, Tier: *tier,
		Repeats: *repeats, Seed: *seed, Model: *model, Adapter: *adapter, AdapterArgs: adapterArgs, SpectacularCLI: *spectacularCLI, OutputDir: *output, AllowHeldOut: *allowHeldOut, ReadIsolation: *readIsolation,
		MaxCalls: *maxCalls, Parallel: *parallel, TrialTimeout: *trialTimeout, RequireCertifiedTelemetry: *requireCertifiedTelemetry,
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

type stringListFlag []string

func (values *stringListFlag) String() string { return fmt.Sprint([]string(*values)) }

func (values *stringListFlag) Set(value string) error {
	*values = append(*values, value)
	return nil
}
