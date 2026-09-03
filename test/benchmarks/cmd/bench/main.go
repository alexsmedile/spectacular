package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"sync"
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
	case "matrix":
		err = matrix(os.Args[2:])
	case "history":
		err = history(os.Args[2:])
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
	fmt.Fprintln(os.Stderr, "usage: go run ./test/benchmarks/cmd/bench <validate|adapter-check|plan|static|run|matrix|history> [flags]")
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

func matrix(args []string) error {
	set := flag.NewFlagSet("matrix", flag.ContinueOnError)
	repo := set.String("repo", ".", "Git repository")
	catalogPath := set.String("catalog", "test/benchmarks/evals.json", "benchmark catalog")
	schemaPath := set.String("schema", "test/benchmarks/agent-result.schema.json", "agent result schema")
	baseline := set.String("baseline", "14158f9", "immutable baseline revision")
	candidate := set.String("candidate", "", "immutable candidate revision")
	tier := set.String("tier", "micro", "micro, smoke, full, or held-out")
	modelsList := set.String("models", "codex:gpt-5.6-terra,claude:claude-opus-5", "comma-separated harness:model pairs")
	spectacularCLI := set.String("spectacular-cli", "", "pinned Spectacular CLI")
	parallelPerModel := set.Int("parallel", 2, "concurrency per model")
	maxCalls := set.Int("max-calls", 50, "maximum model calls allowed")
	outBase := set.String("out", "test/benchmarks/reports/matrix-"+time.Now().Format("20060102-150405"), "output root directory")
	if err := set.Parse(args); err != nil {
		return err
	}
	if *candidate == "" {
		return fmt.Errorf("candidate is required")
	}
	if *spectacularCLI == "" {
		return fmt.Errorf("spectacular-cli is required")
	}

	targets := strings.Split(*modelsList, ",")
	type matrixResult struct {
		harness string
		model   string
		report  bench.RunReport
		outDir  string
		err     error
	}

	results := make([]matrixResult, len(targets))
	var wg sync.WaitGroup

	fmt.Printf("==> Launching multi-harness benchmark matrix across %d targets (tier: %s)...\n", len(targets), *tier)
	for i, target := range targets {
		parts := strings.Split(strings.TrimSpace(target), ":")
		if len(parts) != 2 {
			return fmt.Errorf("invalid harness:model target %q (expected <harness>:<model>)", target)
		}
		harnessName := parts[0]
		modelName := parts[1]
		adapterPath := fmt.Sprintf("test/benchmarks/adapters/%s-adapter.sh", harnessName)
		outDir := filepath.Join(*outBase, harnessName+"-"+filepath.Base(modelName))

		wg.Add(1)
		go func(idx int, hName, mName, aPath, oDir string) {
			defer wg.Done()
			rep, runErr := bench.RunPaired(bench.RunConfig{
				Repo: *repo, CatalogPath: *catalogPath, SchemaPath: *schemaPath,
				BaselineRef: *baseline, BaselineMode: "skill", CandidateRef: *candidate, CandidateMode: "skill",
				Tier: *tier, Repeats: 1, Seed: 1, Model: mName, Adapter: aPath,
				SpectacularCLI: *spectacularCLI, OutputDir: oDir, MaxCalls: *maxCalls, Parallel: *parallelPerModel,
				TrialTimeout: 10 * time.Minute, RequireCertifiedTelemetry: true,
			})
			if runErr == nil {
				_ = bench.WriteRunReport(rep, oDir)
			}
			results[idx] = matrixResult{harness: hName, model: mName, report: rep, outDir: oDir, err: runErr}
		}(i, harnessName, modelName, adapterPath, outDir)
	}

	wg.Wait()

	fmt.Println("\n==================================================================================")
	fmt.Printf("MULTI-HARNESS BENCHMARK MATRIX (Tier: %s | Candidate: %s)\n", *tier, *candidate)
	fmt.Println("==================================================================================")
	fmt.Printf("%-10s | %-25s | %-12s | %-10s | %-12s | %-10s\n", "HARNESS", "MODEL", "VERDICT", "SAFETY", "TASK SUCCESS", "DURATION")
	fmt.Println("----------------------------------------------------------------------------------")
	for _, res := range results {
		if res.err != nil {
			fmt.Printf("%-10s | %-25s | ERROR: %v\n", res.harness, res.model, res.err)
			continue
		}
		summary := res.report.Summary
		candSafety := "N/A"
		candSuccess := "N/A"
		if rates, ok := summary.DimensionRates["candidate"]; ok {
			candSafety = fmt.Sprintf("%.1f%%", rates["safety"]*100)
			candSuccess = fmt.Sprintf("%.1f%%", rates["task_success"]*100)
		}
		duration := "0ms"
		if cost, ok := summary.ObservedCost["candidate"]; ok {
			duration = fmt.Sprintf("%dms", cost.TotalDurationMillis)
		}
		fmt.Printf("%-10s | %-25s | %-12s | %-10s | %-12s | %-10s\n", res.harness, res.model, summary.Verdict, candSafety, candSuccess, duration)
	}
	fmt.Println("==================================================================================")
	fmt.Printf("Detailed reports saved in: %s/\n", *outBase)
	return nil
}

func history(args []string) error {
	set := flag.NewFlagSet("history", flag.ContinueOnError)
	reportsDir := set.String("dir", "test/benchmarks/reports", "reports directory")
	if err := set.Parse(args); err != nil {
		return err
	}

	type historyEntry struct {
		Path             string    `json:"path"`
		Model            string    `json:"model"`
		Tier             string    `json:"tier"`
		Date             time.Time `json:"date"`
		Verdict          string    `json:"verdict"`
		Effect           string    `json:"effect"`
		CandidateSafety  float64   `json:"candidate_safety"`
		CandidateSuccess float64   `json:"candidate_success"`
		InputTokens      int       `json:"input_tokens"`
		DurationMS       int64     `json:"duration_ms"`
	}

	var entries []historyEntry

	err := filepath.Walk(*reportsDir, func(path string, info os.FileInfo, err error) error {
		if err != nil || info.IsDir() || filepath.Base(path) != "report.json" {
			return nil
		}
		data, err := os.ReadFile(path)
		if err != nil {
			return nil
		}
		var rep bench.RunReport
		if err := json.Unmarshal(data, &rep); err != nil {
			return nil
		}
		var candSafety, candSuccess float64
		var inputTok int
		var dur int64
		if rates, ok := rep.Summary.DimensionRates["candidate"]; ok {
			candSafety = rates["safety"]
			candSuccess = rates["task_success"]
		}
		if cost, ok := rep.Summary.ObservedCost["candidate"]; ok {
			inputTok = cost.TotalInputTokens
			dur = cost.TotalDurationMillis
		}
		entries = append(entries, historyEntry{
			Path:             path,
			Model:            rep.Model,
			Tier:             rep.Tier,
			Date:             rep.GeneratedAt,
			Verdict:          rep.Summary.Verdict,
			Effect:           rep.Summary.ComparativeEffect,
			CandidateSafety:  candSafety,
			CandidateSuccess: candSuccess,
			InputTokens:      inputTok,
			DurationMS:       dur,
		})
		return nil
	})

	if err != nil {
		return err
	}

	sort.Slice(entries, func(i, j int) bool {
		return entries[i].Date.Before(entries[j].Date)
	})

	fmt.Println("===============================================================================================")
	fmt.Println("SPECTACULAR BENCHMARK EXECUTION HISTORY")
	fmt.Println("===============================================================================================")
	fmt.Printf("%-19s | %-24s | %-7s | %-12s | %-8s | %-8s | %-10s\n", "DATE", "MODEL", "TIER", "VERDICT", "SAFETY", "SUCCESS", "INPUT TOK")
	fmt.Println("-----------------------------------------------------------------------------------------------")
	for _, e := range entries {
		fmt.Printf("%-19s | %-24s | %-7s | %-12s | %-7.1f%% | %-7.1f%% | %-10d\n",
			e.Date.Format("2006-01-02 15:04"), e.Model, e.Tier, e.Verdict,
			e.CandidateSafety*100, e.CandidateSuccess*100, e.InputTokens)
	}
	fmt.Println("===============================================================================================")
	return nil
}

type stringListFlag []string

func (values *stringListFlag) String() string { return fmt.Sprint([]string(*values)) }

func (values *stringListFlag) Set(value string) error {
	*values = append(*values, value)
	return nil
}
