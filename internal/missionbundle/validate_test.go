package missionbundle

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"reflect"
	"strings"
	"testing"

	"github.com/alexsmedile/spectacular/v2/internal/discovery"
	"github.com/alexsmedile/spectacular/v2/internal/domain"
)

func TestSchemaRegistryOwnsEveryMandatoryValidation(t *testing.T) {
	want := []string{
		"activation-fingerprint", "authority-vocabulary", "baseline-binding", "completion-claim-coverage",
		"contract-binding", "mechanical-scope", "objective-dependency-dag", "reference-integrity",
		"review-independence", "run-state", "safe-file-layout", "transition-atomicity",
		"uuidv7-identity", "yaml-schema",
	}
	if got := ValidatorNames(); !reflect.DeepEqual(got, want) {
		t.Fatalf("validators=%v want=%v", got, want)
	}
}

func TestMandatoryValidatorsReturnTypedZeroMutationRefusals(t *testing.T) {
	root, err := filepath.Abs(filepath.Join("..", ".."))
	if err != nil {
		t.Fatal(err)
	}
	ws, err := discovery.Open(root)
	if err != nil {
		t.Fatal(err)
	}
	m6, err := Load(ws, "M6")
	if err != nil {
		t.Fatal(err)
	}
	m5, err := Load(ws, "M5")
	if err != nil {
		t.Fatal(err)
	}
	before := treeDigest(t, filepath.Join(root, ".spectacular"))
	tests := []struct {
		name   string
		base   *Bundle
		mutate func(*Bundle)
		check  func(*discovery.Workspace, *Bundle) error
		code   domain.RefusalCode
		field  string
	}{
		{"yaml-schema", m6, func(b *Bundle) { b.Validation.Schema = "invented" }, validateSchema, domain.RefusalInvalidKnownField, "validation.schema"},
		{"uuidv7-identity", m6, func(b *Bundle) { b.ID = "not-a-uuid" }, validateIdentity, domain.RefusalInvalidKnownField, "id"},
		{"reference-integrity", m6, func(b *Bundle) { b.Objectives[0].Outcome = "" }, validateReferences, domain.RefusalInvalidKnownField, "objectives"},
		{"contract-binding", m6, func(b *Bundle) { b.Contract.Fingerprint = "sha256:" + string(make([]byte, 64)) }, validateContract, domain.RefusalInvalidKnownField, "contract.fingerprint"},
		{"baseline-binding", m6, func(b *Bundle) { b.Baseline.Commit = "0000000000000000000000000000000000000000" }, validateBaseline, domain.RefusalInvalidKnownField, "baseline.commit"},
		{"activation-fingerprint", m6, func(b *Bundle) { b.Outcome += " drift" }, validateActivation, domain.RefusalStaleFingerprint, "activation.fingerprint"},
		{"completion-claim-coverage", m6, func(b *Bundle) { b.Objectives[0].Claims = []string{"unknown"} }, validateClaims, domain.RefusalInvalidKnownField, "objectives.claims"},
		{"objective-dependency-dag", m6, func(b *Bundle) { b.Objectives[0].After = []string{b.Objectives[0].Ref} }, validateDAG, domain.RefusalInvalidKnownField, "objectives.after"},
		{"run-state", m6, func(b *Bundle) { b.Run.Repairs = b.RepairBudget + 1 }, validateRun, domain.RefusalInvalidKnownField, "run"},
		{"review-independence", m5, func(b *Bundle) { b.Reviews[0].Verdict = "repair" }, validateReviews, domain.RefusalInvalidKnownField, "reviews"},
		{"authority-vocabulary", m6, func(b *Bundle) { b.Authority.Operator = append(b.Authority.Operator, "invent-authority") }, validateAuthority, domain.RefusalInvalidKnownField, "authority.operator"},
		{"mechanical-scope", m6, func(b *Bundle) { b.Scope.Mechanical = []string{"../escape"} }, validateScope, domain.RefusalInvalidKnownField, "scope.mechanical"},
		{"safe-file-layout", m6, func(b *Bundle) { b.Path = ".spectacular/MISSION.md" }, validateLayout, domain.RefusalInvalidKnownField, "path"},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			candidate := cloneBundle(t, test.base)
			test.mutate(candidate)
			err := test.check(ws, candidate)
			var refusal *domain.Refusal
			if !errors.As(err, &refusal) || refusal.Code != test.code || refusal.Field != test.field || refusal.Detail == "" || refusal.Recovery == "" {
				t.Fatalf("refusal=%+v want code=%s field=%s with problem and safe correction", refusal, test.code, test.field)
			}
		})
	}
	transactionRoot := t.TempDir()
	if err := os.MkdirAll(filepath.Join(transactionRoot, "transactions"), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(transactionRoot, "transactions", "pending.json"), []byte("{}"), 0o644); err != nil {
		t.Fatal(err)
	}
	err = validateTransactions(&discovery.Workspace{MetadataDir: transactionRoot}, m6)
	if !domain.RefusalHasCode(err, domain.RefusalTransactionRecovery) {
		t.Fatalf("transaction validator=%v", err)
	}
	if after := treeDigest(t, filepath.Join(root, ".spectacular")); after != before {
		t.Fatal("validation mutated the canonical tree")
	}
}

func FuzzMissionPlanInputNeverPanics(f *testing.F) {
	f.Add([]byte("---\ntype: MissionPlan\n---\n"))
	f.Add([]byte("not frontmatter"))
	f.Fuzz(func(t *testing.T, data []byte) {
		_, _, _ = ReadPlan("-", data)
	})
}

func TestMutationLockRefusesConcurrentWriter(t *testing.T) {
	root := t.TempDir()
	unlock, err := acquireMutationLock(root)
	if err != nil {
		t.Fatal(err)
	}
	defer unlock()
	second, err := acquireMutationLock(root)
	if second != nil {
		second()
	}
	if !domain.RefusalHasCode(err, domain.RefusalCollision) {
		t.Fatalf("concurrent lock error=%v", err)
	}
}

func TestReviewedGitBindingRejectsFabricatedCoordinates(t *testing.T) {
	root, err := filepath.Abs(filepath.Join("..", ".."))
	if err != nil {
		t.Fatal(err)
	}
	commit := strings.TrimSpace(commandOutput(t, root, "git", "rev-parse", "HEAD"))
	tree := strings.TrimSpace(commandOutput(t, root, "git", "rev-parse", "HEAD^{tree}"))
	if err := verifyReviewedGit(root, commit, tree); err != nil {
		t.Fatalf("exact commit/tree rejected: %v", err)
	}
	for _, test := range []struct {
		name, commit, tree, field string
		code                      domain.RefusalCode
	}{
		{"missing commit", strings.Repeat("0", 40), tree, "review.reviewed.commit", domain.RefusalInvalidKnownField},
		{"mismatched tree", commit, strings.Repeat("0", 40), "review.reviewed.tree", domain.RefusalStaleFingerprint},
	} {
		t.Run(test.name, func(t *testing.T) {
			err := verifyReviewedGit(root, test.commit, test.tree)
			var refusal *domain.Refusal
			if !errors.As(err, &refusal) || refusal.Code != test.code || refusal.Field != test.field || refusal.Detail == "" || refusal.Recovery == "" {
				t.Fatalf("refusal=%+v err=%v", refusal, err)
			}
		})
	}
}

func commandOutput(t *testing.T, directory, name string, args ...string) string {
	t.Helper()
	command := exec.Command(name, args...)
	command.Dir = directory
	output, err := command.Output()
	if err != nil {
		t.Fatal(err)
	}
	return string(output)
}

func cloneBundle(t *testing.T, source *Bundle) *Bundle {
	t.Helper()
	data, err := json.Marshal(source)
	if err != nil {
		t.Fatal(err)
	}
	var result Bundle
	if err := json.Unmarshal(data, &result); err != nil {
		t.Fatal(err)
	}
	result.entry = source.entry
	result.document = source.document
	return &result
}

func treeDigest(t *testing.T, root string) string {
	t.Helper()
	hash := sha256.New()
	if err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			return nil
		}
		rel, _ := filepath.Rel(root, path)
		data, readErr := os.ReadFile(path)
		if readErr != nil {
			return readErr
		}
		hash.Write([]byte(rel))
		hash.Write(data)
		return nil
	}); err != nil {
		t.Fatal(err)
	}
	return hex.EncodeToString(hash.Sum(nil))
}

// decodeObjectiveGraph turns fuzz bytes into a small Objective graph. Each input
// byte is one edge: the high nibble selects the dependent Objective and the low
// nibble the dependency, both modulo the node count. Unmapped low nibbles become
// dangling references so the corpus reaches unknown-dependency refusals too.
func decodeObjectiveGraph(data []byte) []Objective {
	if len(data) == 0 {
		return nil
	}
	nodes := int(data[0]%8) + 1
	objectives := make([]Objective, nodes)
	for i := range objectives {
		objectives[i] = Objective{Ref: fmt.Sprintf("O%d", i+1)}
	}
	for _, edge := range data[1:] {
		from := int(edge>>4) % nodes
		target := int(edge & 0x0f)
		dependency := fmt.Sprintf("O%d", target+1)
		objectives[from].After = append(objectives[from].After, dependency)
	}
	return objectives
}

// acyclicByKahn is an independent oracle. validateDAG walks the graph with a
// recursive colour marking; this settles the same question by iteratively
// peeling zero-indegree nodes, so a bug in one is not reproduced by the other.
func acyclicByKahn(objectives []Objective) bool {
	known := map[string]bool{}
	for _, objective := range objectives {
		known[objective.Ref] = true
	}
	indegree := map[string]int{}
	dependents := map[string][]string{}
	for _, objective := range objectives {
		for _, dependency := range objective.After {
			if !known[dependency] {
				continue
			}
			indegree[objective.Ref]++
			dependents[dependency] = append(dependents[dependency], objective.Ref)
		}
	}
	queue := []string{}
	for ref := range known {
		if indegree[ref] == 0 {
			queue = append(queue, ref)
		}
	}
	settled := 0
	for len(queue) > 0 {
		ref := queue[0]
		queue = queue[1:]
		settled++
		for _, dependent := range dependents[ref] {
			indegree[dependent]--
			if indegree[dependent] == 0 {
				queue = append(queue, dependent)
			}
		}
	}
	return settled == len(known)
}

// FuzzObjectiveDependencyGraph searches generated dependency graphs for any case
// where validateDAG disagrees with an independent topological sort, refuses with
// the wrong code or field, or panics. Hand-written cases cover only a self-cycle
// and one unresolved dependency; generated graphs reach multi-node cycles,
// diamonds, and repeated edges that nobody thought to write down.
func FuzzObjectiveDependencyGraph(f *testing.F) {
	f.Add([]byte{0})                   // single Objective, no edges
	f.Add([]byte{2, 0x10})             // O2 depends on O1: acyclic
	f.Add([]byte{0, 0x00})             // O1 depends on itself: self-cycle
	f.Add([]byte{2, 0x10, 0x01})       // O1 <-> O2: two-node cycle
	f.Add([]byte{3, 0x10, 0x21, 0x02}) // O1 -> O2 -> O3 -> O1: three-node cycle
	f.Add([]byte{3, 0x10, 0x20, 0x21}) // diamond
	f.Add([]byte{1, 0x0f})             // dangling dependency on a missing Objective

	f.Fuzz(func(t *testing.T, data []byte) {
		objectives := decodeObjectiveGraph(data)
		if len(objectives) == 0 {
			return
		}
		bundle := &Bundle{Objectives: objectives}

		err := validateDAG(nil, bundle)

		known := map[string]bool{}
		for _, objective := range objectives {
			known[objective.Ref] = true
		}
		selfOrUnknown := false
		for _, objective := range objectives {
			for _, dependency := range objective.After {
				if !known[dependency] || dependency == objective.Ref {
					selfOrUnknown = true
				}
			}
		}
		acyclic := acyclicByKahn(objectives)

		switch {
		case selfOrUnknown || !acyclic:
			if err == nil {
				t.Fatalf("accepted an invalid dependency graph: %#v", objectives)
			}
			if !domain.RefusalHasCode(err, domain.RefusalInvalidKnownField) {
				t.Fatalf("refusal code=%v for %#v", err, objectives)
			}
			var refusal *domain.Refusal
			if !errors.As(err, &refusal) {
				t.Fatalf("refusal was not typed: %v", err)
			}
			if refusal.Field != "objectives.after" {
				t.Fatalf("refusal field=%q, want objectives.after", refusal.Field)
			}
			if refusal.Detail == "" || refusal.Recovery == "" {
				t.Fatalf("refusal lacked problem or correction: %+v", refusal)
			}
		default:
			if err != nil {
				t.Fatalf("refused a valid acyclic graph %#v: %v", objectives, err)
			}
		}
	})
}
