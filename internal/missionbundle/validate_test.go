package missionbundle

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
	"reflect"
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
	}{
		{"yaml-schema", m6, func(b *Bundle) { b.Validation.Schema = "invented" }, validateSchema},
		{"uuidv7-identity", m6, func(b *Bundle) { b.ID = "not-a-uuid" }, validateIdentity},
		{"reference-integrity", m6, func(b *Bundle) { b.Objectives[0].Outcome = "" }, validateReferences},
		{"contract-binding", m6, func(b *Bundle) { b.Contract.Fingerprint = "sha256:" + string(make([]byte, 64)) }, validateContract},
		{"baseline-binding", m6, func(b *Bundle) { b.Baseline.Commit = "0000000000000000000000000000000000000000" }, validateBaseline},
		{"activation-fingerprint", m6, func(b *Bundle) { b.Outcome += " drift" }, validateActivation},
		{"completion-claim-coverage", m6, func(b *Bundle) { b.Objectives[0].Claims = []string{"unknown"} }, validateClaims},
		{"objective-dependency-dag", m6, func(b *Bundle) { b.Objectives[0].After = []string{b.Objectives[0].Ref} }, validateDAG},
		{"run-state", m6, func(b *Bundle) { b.Run.Repairs = b.RepairBudget + 1 }, validateRun},
		{"review-independence", m5, func(b *Bundle) { b.Reviews[0].Verdict = "repair" }, validateReviews},
		{"authority-vocabulary", m6, func(b *Bundle) { b.Authority.Operator = append(b.Authority.Operator, "invent-authority") }, validateAuthority},
		{"mechanical-scope", m6, func(b *Bundle) { b.Scope.Mechanical = []string{"../escape"} }, validateScope},
		{"safe-file-layout", m6, func(b *Bundle) { b.Path = ".spectacular/MISSION.md" }, validateLayout},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			candidate := cloneBundle(t, test.base)
			test.mutate(candidate)
			err := test.check(ws, candidate)
			var refusal *domain.Refusal
			if !errors.As(err, &refusal) || refusal.Code == "" || refusal.Field == "" {
				t.Fatalf("error=%v want typed field refusal", err)
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
