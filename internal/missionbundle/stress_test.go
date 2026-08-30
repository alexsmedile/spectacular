package missionbundle

import (
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"io/fs"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/alexsmedile/spectacular/v2/internal/discovery"
	"github.com/alexsmedile/spectacular/v2/internal/domain"
	"github.com/alexsmedile/spectacular/v2/internal/governance"
)

// This cluster exercises the public mutation service at its highest-value
// boundaries. The transaction engine's lower-level interruption matrix lives
// in internal/governance; these checks prove the Mission service reaches that
// boundary without changing canonical files first.
func TestMissionServiceRetryAndRefusalBoundaries(t *testing.T) {
	root := missionServiceFixture(t)
	plan, raw := stressPlan()
	svc := openMissionService(t, root)

	started, err := svc.Start(plan, raw)
	if err != nil {
		t.Fatal(err)
	}
	if started.Ref == "" || started.Path == "" {
		t.Fatalf("incomplete start receipt: %+v", started)
	}

	// A retry in a fresh command context resolves the persisted start_key and
	// returns the original Mission instead of allocating another identity.
	svc = openMissionService(t, root)
	retried, err := svc.Start(plan, raw)
	if err != nil {
		t.Fatal(err)
	}
	if retried.Ref != started.Ref || retried.Path != started.Path {
		t.Fatalf("retry drifted: first=%+v retry=%+v", started, retried)
	}
	if got := countMissionStarts(t, svc.Workspace, raw); got != 1 {
		t.Fatalf("logical start count=%d, want 1", got)
	}

	t.Run("dependency refusal writes nothing", func(t *testing.T) {
		before := canonicalTreeDigest(t, root)
		_, err := svc.FinishObjective(started.Ref + "/O2")
		assertRefusal(t, err, domain.RefusalInvalidKnownField, "objectives.after")
		if after := canonicalTreeDigest(t, root); after != before {
			t.Fatal("dependency refusal changed canonical files")
		}
	})

	t.Run("concurrent writer refusal writes nothing", func(t *testing.T) {
		unlock, err := acquireMutationLock(root)
		if err != nil {
			t.Fatal(err)
		}
		defer unlock()
		before := canonicalTreeDigest(t, root)
		_, err = svc.FinishObjective(started.Ref + "/O1")
		assertRefusal(t, err, domain.RefusalCollision, "transactions")
		if after := canonicalTreeDigest(t, root); after != before {
			t.Fatal("concurrency refusal changed canonical files")
		}
	})

	t.Run("stale source refusal preserves external winner", func(t *testing.T) {
		missionPath := filepath.Join(root, filepath.FromSlash(started.Path))
		file, err := os.OpenFile(missionPath, os.O_APPEND|os.O_WRONLY, 0)
		if err != nil {
			t.Fatal(err)
		}
		if _, err := file.WriteString("\nExternal body edit.\n"); err != nil {
			file.Close()
			t.Fatal(err)
		}
		if err := file.Close(); err != nil {
			t.Fatal(err)
		}
		before := canonicalTreeDigest(t, root)
		// Exercise the post-validation optimistic check directly: public
		// mutations first reopen under the lock, while apply still refuses a
		// source changed between that read and the transaction boundary.
		_, err = svc.finishObjective(started.Ref + "/O1")
		assertRefusal(t, err, domain.RefusalStaleFingerprint, started.Path)
		if after := canonicalTreeDigest(t, root); after != before {
			t.Fatal("stale refusal did not preserve the external winner")
		}
	})

	// Reloading accepts the external body-only edit, and repeating an already
	// completed transition is a stable no-op.
	svc = openMissionService(t, root)
	if _, err := svc.FinishObjective(started.Ref + "/O1"); err != nil {
		t.Fatal(err)
	}
	svc = openMissionService(t, root)
	before := canonicalTreeDigest(t, root)
	result, err := svc.FinishObjective(started.Ref + "/O1")
	if err != nil {
		t.Fatal(err)
	}
	if len(result.Changed) != 0 {
		t.Fatalf("idempotent finish reported changes: %v", result.Changed)
	}
	if after := canonicalTreeDigest(t, root); after != before {
		t.Fatal("idempotent finish rewrote canonical files")
	}
}

func TestMissionStartRejectsSymlinkedTargetWithoutPartialWrite(t *testing.T) {
	root := missionServiceFixture(t)
	svc := openMissionService(t, root)
	missionRoot := filepath.Join(root, ".spectacular", "missions")
	heldRoot := filepath.Join(root, ".spectacular", "missions-held")
	outside := t.TempDir()
	if err := os.Rename(missionRoot, heldRoot); err != nil {
		t.Fatal(err)
	}
	if err := os.Symlink(outside, missionRoot); err != nil {
		t.Fatal(err)
	}
	before := canonicalTreeDigest(t, root)
	plan, raw := stressPlan()
	_, err := svc.Start(plan, append(raw, []byte("-symlink")...))
	if !domain.RefusalHasCode(err, domain.RefusalPathEscape) {
		t.Fatalf("symlinked target error=%v, want path_escape", err)
	}
	if after := canonicalTreeDigest(t, root); after != before {
		t.Fatal("failed preparation changed canonical files")
	}
	entries, err := os.ReadDir(outside)
	if err != nil {
		t.Fatal(err)
	}
	if len(entries) != 0 {
		t.Fatalf("transaction escaped Mission root: %v", entries)
	}
}

func TestReviewRetryConvergesOnOneIdentity(t *testing.T) {
	root := missionServiceFixture(t)
	plan, raw := stressPlan()
	svc := openMissionService(t, root)
	started, err := svc.Start(plan, append(raw, []byte("-review")...))
	if err != nil {
		t.Fatal(err)
	}
	svc = openMissionService(t, root)
	bundle, err := svc.Show(started.Ref)
	if err != nil {
		t.Fatal(err)
	}
	commit := gitOutput(t, root, "rev-parse", "HEAD")
	tree := gitOutput(t, root, "rev-parse", "HEAD^{tree}")
	review := []byte(fmt.Sprintf(`---
type: ReviewDraft
title: Atomic review
status: passed
reviewed:
  commit: %s
  tree: %s
  activation_fingerprint: %s
claims:
  - claim: atomic
    verdict: pass
---
# Review
`, commit, tree, bundle.Activation.Fingerprint))
	first, err := svc.RecordReview(started.Ref, "-", review)
	if err != nil {
		t.Fatal(err)
	}
	svc = openMissionService(t, root)
	before := canonicalTreeDigest(t, root)
	retry, err := svc.RecordReview(started.Ref, "-", review)
	if err != nil {
		t.Fatal(err)
	}
	if retry.Ref != first.Ref || retry.Path != first.Path || len(retry.Changed) != 0 {
		t.Fatalf("review retry drifted: first=%+v retry=%+v", first, retry)
	}
	if after := canonicalTreeDigest(t, root); after != before {
		t.Fatal("review retry rewrote canonical files")
	}
	bundle, err = openMissionService(t, root).Show(started.Ref)
	if err != nil {
		t.Fatal(err)
	}
	if len(bundle.Reviews) != 1 || bundle.Reviews[0].Ref != strings.TrimPrefix(first.Ref, started.Ref+"/") {
		t.Fatalf("review pointers=%+v, want one %s", bundle.Reviews, first.Ref)
	}
}

func TestRecordReviewAutoDerivesTreeWhenOmitted(t *testing.T) {
	root := missionServiceFixture(t)
	plan, raw := stressPlan()
	svc := openMissionService(t, root)
	started, err := svc.Start(plan, raw)
	if err != nil {
		t.Fatal(err)
	}
	svc = openMissionService(t, root)
	bundle, err := svc.Show(started.Ref)
	if err != nil {
		t.Fatal(err)
	}
	commit := gitOutput(t, root, "rev-parse", "HEAD")
	tree := gitOutput(t, root, "rev-parse", "HEAD^{tree}")
	review := []byte(fmt.Sprintf(`---
type: ReviewDraft
title: Atomic review without explicit tree
status: passed
reviewed:
  commit: %s
  activation_fingerprint: %s
claims:
  - claim: atomic
    verdict: pass
---
# Review
`, commit, bundle.Activation.Fingerprint))
	res, err := svc.RecordReview(started.Ref, "-", review)
	if err != nil {
		t.Fatalf("RecordReview failed without tree: %v", err)
	}
	reloaded, err := openMissionService(t, root).Show(started.Ref)
	if err != nil {
		t.Fatal(err)
	}
	if len(reloaded.Reviews) != 1 {
		t.Fatalf("expected 1 review, got %d", len(reloaded.Reviews))
	}
	if reloaded.Reviews[0].Document.Reviewed.Tree != tree {
		t.Fatalf("expected auto-derived tree %s, got %s", tree, reloaded.Reviews[0].Document.Reviewed.Tree)
	}
	_ = res
}

func TestPromoteRefusesUndiscoveredDerivedTargetCollision(t *testing.T) {
	root := missionServiceFixture(t)
	plan, raw := stressPlan()
	svc := openMissionService(t, root)
	started, err := svc.Start(plan, append(raw, []byte("-collision")...))
	if err != nil {
		t.Fatal(err)
	}
	target := filepath.Join(root, filepath.FromSlash(filepath.Join(filepath.Dir(started.Path), "objectives", "O1-first-dependency.md")))
	if err := os.MkdirAll(target, 0o755); err != nil {
		t.Fatal(err)
	}
	svc = openMissionService(t, root)
	before := canonicalTreeDigest(t, root)
	_, err = svc.PromoteObjective(started.Ref + "/O1")
	assertRefusal(t, err, domain.RefusalCollision, filepath.ToSlash(strings.TrimPrefix(target, root+string(filepath.Separator))))
	if after := canonicalTreeDigest(t, root); after != before {
		t.Fatal("target collision changed canonical files")
	}
	if info, statErr := os.Stat(target); statErr != nil || !info.IsDir() {
		t.Fatalf("colliding target was replaced: info=%v err=%v", info, statErr)
	}
}

func TestMutationCommandsRollbackAtEveryInstallBoundary(t *testing.T) {
	type caseDef struct {
		prepare func(t *testing.T) (string, func(Service) error)
	}
	cases := map[string]caseDef{
		"start": {
			prepare: func(t *testing.T) (string, func(Service) error) {
				root := missionServiceFixture(t)
				plan, raw := stressPlan()
				return root, func(s Service) error {
					_, err := s.Start(plan, append(raw, []byte("-fault")...))
					return err
				}
			},
		},
		"promote": {
			prepare: func(t *testing.T) (string, func(Service) error) {
				root := missionServiceFixture(t)
				plan, raw := stressPlan()
				svc := openMissionService(t, root)
				started, err := svc.Start(plan, append(raw, []byte("-promote-fault")...))
				if err != nil {
					t.Fatal(err)
				}
				return root, func(s Service) error {
					_, err := s.PromoteObjective(started.Ref + "/O1")
					return err
				}
			},
		},
		"run": {
			prepare: func(t *testing.T) (string, func(Service) error) {
				root := missionServiceFixture(t)
				plan, raw := stressPlan()
				svc := openMissionService(t, root)
				started, err := svc.Start(plan, append(raw, []byte("-run-fault")...))
				if err != nil {
					t.Fatal(err)
				}
				return root, func(s Service) error {
					_, err := s.StartRun(started.Ref, "Second run")
					return err
				}
			},
		},
		"review": {
			prepare: func(t *testing.T) (string, func(Service) error) {
				root := missionServiceFixture(t)
				plan, raw := stressPlan()
				svc := openMissionService(t, root)
				started, err := svc.Start(plan, append(raw, []byte("-review-fault")...))
				if err != nil {
					t.Fatal(err)
				}
				svc = openMissionService(t, root)
				review := reviewDraft(t, root, svc, started.Ref)
				return root, func(s Service) error {
					_, err := s.RecordReview(started.Ref, "-", review)
					return err
				}
			},
		},
		"complete": {
			prepare: func(t *testing.T) (string, func(Service) error) {
				root := missionServiceFixture(t)
				plan, raw := stressPlan()
				svc := openMissionService(t, root)
				started, err := svc.Start(plan, append(raw, []byte("-complete-fault")...))
				if err != nil {
					t.Fatal(err)
				}
				if _, err := svc.FinishObjective(started.Ref + "/O1"); err != nil {
					t.Fatal(err)
				}
				svc = openMissionService(t, root)
				if _, err := svc.FinishObjective(started.Ref + "/O2"); err != nil {
					t.Fatal(err)
				}
				return root, func(s Service) error {
					_, err := s.Complete(started.Ref, planOwner())
					return err
				}
			},
		},
	}

	for name, def := range cases {
		t.Run(name, func(t *testing.T) {
			baseRoot, invoke := def.prepare(t)
			clone := func(t *testing.T) (string, Service) {
				root := t.TempDir()
				copyTree(t, baseRoot, root)
				return root, openMissionService(t, root)
			}
			probeRoot, probeServe := clone(t)
			count := 0
			probeError := errors.New("capture transaction width")
			probeServe.ApplyTransaction = func(_ string, _ string, changes []governance.FileChange) error {
				count = len(changes)
				return probeError
			}
			if err := invoke(probeServe); !errors.Is(err, probeError) || count == 0 {
				t.Fatalf("transaction probe err=%v width=%d (root=%s)", err, count, probeRoot)
			}

			for failAfter := 0; failAfter < count; failAfter++ {
				failAfter := failAfter
				t.Run(fmt.Sprintf("after-%d-of-%d", failAfter, count), func(t *testing.T) {
					t.Parallel()
					subRoot, subServe := clone(t)
					before := canonicalTreeDigest(t, subRoot)
					subServe.ApplyTransaction = func(root, key string, changes []governance.FileChange) error {
						return governance.ApplyTransactionWithFailure(root, key, changes, failAfter)
					}
					if err := invoke(subServe); err == nil {
						t.Fatal("injected transaction failure unexpectedly succeeded")
					}
					if after := canonicalTreeDigest(t, subRoot); after != before {
						t.Fatal("injected failure changed canonical files")
					}
					if _, err := discovery.Open(subRoot); err != nil {
						t.Fatalf("rollback left workspace unreadable: %v", err)
					}
				})
			}
		})
	}
}

func TestContainedBundlePathRefusalCluster(t *testing.T) {
	base := t.TempDir()
	for _, path := range []string{"", "../escape.md", "a/../../escape.md", filepath.Join(base, "absolute.md"), "a/../b.md"} {
		t.Run(strings.ReplaceAll(path, "/", "_"), func(t *testing.T) {
			_, err := containedFile(base, path)
			if !domain.RefusalHasCode(err, domain.RefusalPathEscape) {
				t.Fatalf("path %q error=%v, want path_escape", path, err)
			}
		})
	}
	if got, err := containedFile(base, "objectives/O1.md"); err != nil || got != filepath.Join(base, "objectives", "O1.md") {
		t.Fatalf("canonical path=%q err=%v", got, err)
	}
}

func stressPlan() (Plan, []byte) {
	return Plan{
		Type:     "MissionPlan",
		Title:    "Atomic stress fixture",
		Owner:    "Test owner",
		Contract: Binding{Ref: "Contract:01a00a20-63dd-7670-97f1-9eb8e12adc3a"},
		Outcome:  "Exercise typed mutation boundaries.",
		Review:   "automatic",
		Completion: []Criterion{{
			Claim: "atomic", PassBoundary: "Transitions are safe.", ProofRequirement: "Focused service tests pass.",
		}},
		Objectives: []Objective{
			{Outcome: "First dependency", Status: "active", Claims: []string{"atomic"}},
			{Outcome: "Dependent work", Status: "pending", After: []string{"O1"}, Claims: []string{"atomic"}},
		},
		Authority: Authority{
			Operator:      []string{"inspect", "edit-in-scope", "run-checks"},
			RequiresOwner: []string{"activate-mission", "change-outcome-or-completion"},
		},
		Scope: Scope{
			Mechanical: []string{"internal/missionbundle/"},
			Semantic:   []string{"Mutation safety."},
		},
		RepairBudget: 1,
		AllowMain:    true,
		Dependencies: []string{},
		Gaps:         []string{},
		Stops:        []string{"A transition can partially write."},
		Body:         "# Atomic stress fixture\n",
	}, []byte("stable logical Mission plan")
}

func planOwner() string { return "Test owner" }

func reviewDraft(t *testing.T, root string, svc Service, missionRef string) []byte {
	t.Helper()
	bundle, err := svc.Show(missionRef)
	if err != nil {
		t.Fatal(err)
	}
	return []byte(fmt.Sprintf(`---
type: ReviewDraft
title: Atomic review
status: passed
reviewed:
  commit: %s
  tree: %s
  activation_fingerprint: %s
claims:
  - claim: atomic
    verdict: pass
---
# Review
`, gitOutput(t, root, "rev-parse", "HEAD"), gitOutput(t, root, "rev-parse", "HEAD^{tree}"), bundle.Activation.Fingerprint))
}

var (
	fixtureTemplateDir  string
	fixtureTemplateOnce sync.Once
	fixtureTemplateErr  error
)

func initFixtureTemplate() (string, error) {
	fixtureTemplateOnce.Do(func() {
		source, err := filepath.Abs(filepath.Join("..", "..", ".spectacular"))
		if err != nil {
			fixtureTemplateErr = err
			return
		}
		dir, err := os.MkdirTemp("", "spectacular-fixture-template-*")
		if err != nil {
			fixtureTemplateErr = err
			return
		}
		if err := copyTreeDirect(source, filepath.Join(dir, ".spectacular")); err != nil {
			fixtureTemplateErr = err
			return
		}
		command := exec.Command("git", "init", "-q")
		command.Dir = dir
		if output, err := command.CombinedOutput(); err != nil {
			fixtureTemplateErr = fmt.Errorf("git init: %w: %s", err, output)
			return
		}
		for _, args := range [][]string{
			{"config", "user.email", "tests@spectacular.invalid"},
			{"config", "user.name", "Spectacular Tests"},
			{"add", ".spectacular"},
			{"commit", "-qm", "fixture"},
		} {
			cmd := exec.Command("git", args...)
			cmd.Dir = dir
			if output, err := cmd.CombinedOutput(); err != nil {
				fixtureTemplateErr = fmt.Errorf("git %v: %w: %s", args, err, output)
				return
			}
		}
		fixtureTemplateDir = dir
	})
	return fixtureTemplateDir, fixtureTemplateErr
}

func missionServiceFixture(t *testing.T) string {
	t.Helper()
	template, err := initFixtureTemplate()
	if err != nil {
		t.Fatal(err)
	}
	root := t.TempDir()
	copyTree(t, template, root)
	return root
}

func openMissionService(t *testing.T, root string) Service {
	t.Helper()
	ws, err := discovery.Open(root)
	if err != nil {
		t.Fatal(err)
	}
	return Service{Workspace: ws, Now: func() time.Time {
		return time.Date(2026, 8, 16, 12, 0, 0, 0, time.UTC)
	}}
}

func countMissionStarts(t *testing.T, ws *discovery.Workspace, raw []byte) int {
	t.Helper()
	digest := sha256.Sum256(raw)
	want := "sha256:" + hex.EncodeToString(digest[:])
	count := 0
	for _, entry := range ws.OfType(domain.Mission) {
		if node := entry.Document.Unknown["start_key"]; node != nil && node.Value == want {
			count++
		}
	}
	return count
}

func assertRefusal(t *testing.T, err error, code domain.RefusalCode, field string) {
	t.Helper()
	var refusal *domain.Refusal
	if !errors.As(err, &refusal) || refusal.Code != code || refusal.Field != field || refusal.Detail == "" || refusal.Recovery == "" {
		t.Fatalf("refusal=%+v err=%v, want code=%s field=%s with problem and correction", refusal, err, code, field)
	}
}

func canonicalTreeDigest(t *testing.T, root string) string {
	t.Helper()
	hash := sha256.New()
	metadata := filepath.Join(root, ".spectacular")
	if err := filepath.WalkDir(metadata, func(path string, entry fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if entry.IsDir() && entry.Name() == "transactions" {
			return filepath.SkipDir
		}
		if entry.IsDir() || entry.Type()&os.ModeSymlink != 0 {
			return nil
		}
		rel, err := filepath.Rel(metadata, path)
		if err != nil {
			return err
		}
		data, err := os.ReadFile(path)
		if err != nil {
			return err
		}
		hash.Write([]byte(filepath.ToSlash(rel)))
		hash.Write(data)
		return nil
	}); err != nil {
		t.Fatal(err)
	}
	return hex.EncodeToString(hash.Sum(nil))
}

func copyTree(t *testing.T, source, target string) {
	t.Helper()
	if err := copyTreeDirect(source, target); err != nil {
		t.Fatal(err)
	}
}

func copyTreeDirect(source, target string) error {
	return filepath.WalkDir(source, func(path string, entry fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		rel, err := filepath.Rel(source, path)
		if err != nil {
			return err
		}
		destination := filepath.Join(target, rel)
		if entry.IsDir() {
			return os.MkdirAll(destination, 0o755)
		}
		data, err := os.ReadFile(path)
		if err != nil {
			return err
		}
		return os.WriteFile(destination, data, 0o644)
	})
}

func git(t *testing.T, root string, args ...string) {
	t.Helper()
	command := exec.Command("git", args...)
	command.Dir = root
	if output, err := command.CombinedOutput(); err != nil {
		t.Fatalf("git %v: %v\n%s", args, err, output)
	}
}

func gitOutput(t *testing.T, root string, args ...string) string {
	t.Helper()
	command := exec.Command("git", args...)
	command.Dir = root
	output, err := command.Output()
	if err != nil {
		t.Fatalf("git %v: %v", args, err)
	}
	return strings.TrimSpace(string(output))
}

func TestCreateContractScaffoldsValidContract(t *testing.T) {
	root := missionServiceFixture(t)
	service := openMissionService(t, root)

	res, err := service.CreateContract("CC-sample-feature", "Sample Feature Contract", "Alex")
	if err != nil {
		t.Fatalf("CreateContract failed: %v", err)
	}

	if res.Operation != "contract.create" {
		t.Fatalf("Operation=%q want contract.create", res.Operation)
	}
	if !strings.HasPrefix(res.Ref, "CC-") {
		t.Fatalf("Ref=%q want prefix CC-", res.Ref)
	}

	// Reopen workspace and lookup the contract
	ws, err := discovery.Open(root)
	if err != nil {
		t.Fatalf("discovery.Open failed: %v", err)
	}
	entry, err := ws.Lookup(res.Ref, domain.Contract)
	if err != nil {
		t.Fatalf("ws.Lookup failed for created contract: %v", err)
	}
	if entry.Document.Record.Type != domain.Contract {
		t.Fatalf("Record.Type=%v want Contract", entry.Document.Record.Type)
	}
	if *entry.Document.Record.Title != "Sample Feature Contract" {
		t.Fatalf("Title=%q want 'Sample Feature Contract'", *entry.Document.Record.Title)
	}
}

func TestResolveContractAndValidateAcceptsHumanRefs(t *testing.T) {
	root, err := filepath.Abs(filepath.Join("..", ".."))
	if err != nil {
		t.Fatal(err)
	}
	ws, err := discovery.Open(root)
	if err != nil {
		t.Fatal(err)
	}

	binding, err := resolveContract(ws, "CC-missioncli")
	if err != nil {
		t.Fatalf("resolveContract failed for CC-missioncli: %v", err)
	}
	if !strings.HasPrefix(binding.Ref, "Contract:") {
		t.Fatalf("binding.Ref=%q want prefix Contract:", binding.Ref)
	}
	if !strings.HasPrefix(binding.Fingerprint, "sha256:") {
		t.Fatalf("binding.Fingerprint=%q want sha256:", binding.Fingerprint)
	}

	bundle, err := Load(ws, "M11")
	if err != nil {
		t.Fatal(err)
	}
	// Test that validateContract succeeds with human ref as well
	cloned := cloneBundle(t, bundle)
	cloned.Contract.Ref = "CC-missioncli"
	if err := validateContract(ws, cloned); err != nil {
		t.Fatalf("validateContract with human ref failed: %v", err)
	}
}

func TestAmendScopeExpandsMechanicalEnvelope(t *testing.T) {
	root := missionServiceFixture(t)
	service := openMissionService(t, root)

	plan, raw := stressPlan()
	plan.Title = "Scope expansion test"
	started, err := service.Start(plan, raw)
	if err != nil {
		t.Fatalf("Start failed: %v", err)
	}

	// Test dryRun first
	dryRes, err := service.AmendScope(started.Ref, []string{"Harbor/MenuBar.swift", "HarborKit/"}, "Alex", "discovering menu bar dependency", true)
	if err != nil {
		t.Fatalf("AmendScope dry-run failed: %v", err)
	}
	if dryRes.Fingerprint == "" {
		t.Fatalf("dry-run did not emit new fingerprint")
	}

	res, err := service.AmendScope(started.Ref, []string{"Harbor/MenuBar.swift", "HarborKit/"}, "Alex", "discovering menu bar dependency", false)
	if err != nil {
		t.Fatalf("AmendScope failed: %v", err)
	}
	if res.Operation != "mission.amend_scope" {
		t.Fatalf("Operation=%q want mission.amend_scope", res.Operation)
	}

	// Reopen and verify amended bundle
	ws, err := discovery.Open(root)
	if err != nil {
		t.Fatal(err)
	}
	bundle, err := Load(ws, started.Ref)
	if err != nil {
		t.Fatal(err)
	}

	hasMenuBar := false
	for _, p := range bundle.Scope.Mechanical {
		if p == "Harbor/MenuBar.swift" {
			hasMenuBar = true
		}
	}
	if !hasMenuBar {
		t.Fatalf("Scope.Mechanical=%v does not contain Harbor/MenuBar.swift", bundle.Scope.Mechanical)
	}

	// Verify bundle passes validation
	if _, err := Validate(ws, bundle); err != nil {
		t.Fatalf("Validate failed on amended bundle: %v", err)
	}
}

func TestAutoPromoteObjectiveOnStartRun(t *testing.T) {
	root := missionServiceFixture(t)
	service := openMissionService(t, root)

	plan, raw := stressPlan()
	plan.Title = "Auto promote test"
	started, err := service.Start(plan, raw)
	if err != nil {
		t.Fatalf("Start failed: %v", err)
	}

	// Start run for objective O2 (which starts as inline)
	runRes, err := service.StartRun(started.Ref+"/O2", "Work on O2")
	if err != nil {
		t.Fatalf("StartRun failed: %v", err)
	}
	if runRes.Operation != "run.start" {
		t.Fatalf("Operation=%q want run.start", runRes.Operation)
	}

	ws, err := discovery.Open(root)
	if err != nil {
		t.Fatal(err)
	}
	bundle, err := Load(ws, started.Ref)
	if err != nil {
		t.Fatal(err)
	}
	if bundle.Objectives[1].File == "" {
		t.Fatalf("bundle.Objectives[1].File is empty, want auto-promoted file path")
	}
	if !strings.HasPrefix(bundle.Objectives[1].File, "objectives/O2-") {
		t.Fatalf("bundle.Objectives[1].File=%q want prefix objectives/O2-", bundle.Objectives[1].File)
	}
}

func TestCloseMissionFinishesObjectivesAndCompletes(t *testing.T) {
	root := missionServiceFixture(t)
	service := openMissionService(t, root)

	plan, raw := stressPlan()
	plan.Title = "Close mission test"
	plan.Review = "automatic"
	started, err := service.Start(plan, raw)
	if err != nil {
		t.Fatalf("Start failed: %v", err)
	}

	// Execute close mission
	res, err := service.CloseMission(started.Ref, "Test owner")
	if err != nil {
		t.Fatalf("CloseMission failed: %v", err)
	}
	if res.Operation != "mission.complete" {
		t.Fatalf("Operation=%q want mission.complete", res.Operation)
	}

	// Verify final state
	ws, err := discovery.Open(root)
	if err != nil {
		t.Fatal(err)
	}
	bundle, err := Load(ws, started.Ref)
	if err != nil {
		t.Fatal(err)
	}
	if bundle.Status != "completed" {
		t.Fatalf("bundle.Status=%q want completed", bundle.Status)
	}
	for _, obj := range bundle.Objectives {
		if obj.Status != "implemented" {
			t.Fatalf("Objective %s status=%q want implemented", obj.Ref, obj.Status)
		}
	}
}

func TestDriftFlagsAwaitingReviewWhenEvidenceRecorded(t *testing.T) {
	evidenceDoc := &Evidence{
		Claims: []string{"alpha"},
		Checks: []EvidenceCheck{{Name: "unit-tests", Result: "pass"}},
	}

	bundle := &Bundle{
		Completion: []Criterion{criterion("alpha"), criterion("beta")},
		Objectives: []Objective{
			{Ref: "O1", Status: "implemented", Claims: []string{"alpha", "beta"}},
		},
		Evidence: []EvidencePointer{
			{Ref: "E1", Document: evidenceDoc},
		},
	}

	drift := bundle.Drift()
	alphaFlags := flagsFor(drift, "alpha")
	betaFlags := flagsFor(drift, "beta")

	hasAwaitingReview := false
	for _, f := range alphaFlags {
		if f == FlagAwaitingReview {
			hasAwaitingReview = true
		}
	}
	if !hasAwaitingReview {
		t.Fatalf("alphaFlags=%v want FlagAwaitingReview", alphaFlags)
	}

	hasUnproven := false
	for _, f := range betaFlags {
		if f == FlagUnproven {
			hasUnproven = true
		}
	}
	if !hasUnproven {
		t.Fatalf("betaFlags=%v want FlagUnproven", betaFlags)
	}
}

func TestValidateScopeDetectsChangedFilesOutsideScope(t *testing.T) {
	root := missionServiceFixture(t)
	service := openMissionService(t, root)

	plan, raw := stressPlan()
	plan.Title = "Scope violation test"
	plan.Scope.Mechanical = []string{"internal/missionbundle/"}
	started, err := service.Start(plan, raw)
	if err != nil {
		t.Fatalf("Start failed: %v", err)
	}

	// Modify a file outside mechanical scope in the git tree
	unscopedFile := filepath.Join(root, "outside.txt")
	if err := os.WriteFile(unscopedFile, []byte("unscoped change"), 0o644); err != nil {
		t.Fatal(err)
	}
	commandOutput(t, root, "git", "add", "outside.txt")
	commandOutput(t, root, "git", "commit", "-qm", "unscoped commit")

	ws, err := discovery.Open(root)
	if err != nil {
		t.Fatal(err)
	}
	bundle, err := Load(ws, started.Ref)
	if err != nil {
		t.Fatal(err)
	}

	err = validateScope(ws, bundle)
	if err == nil {
		t.Fatalf("validateScope must refuse when file is modified outside mechanical scope")
	}
	refusal, ok := err.(*domain.Refusal)
	if !ok || refusal.Code != domain.RefusalInvalidKnownField {
		t.Fatalf("err=%v want RefusalInvalidKnownField", err)
	}
	if !strings.Contains(refusal.Detail, "outside.txt") {
		t.Fatalf("refusal.Detail=%q want mention of outside.txt", refusal.Detail)
	}
}

func TestBranchGuardrailRefusesMainWithoutOverride(t *testing.T) {
	root := missionServiceFixture(t)
	service := openMissionService(t, root)

	plan, raw := stressPlan()
	plan.Title = "Branch refusal test"
	plan.AllowMain = false
	plan.CreateBranch = false

	_, err := service.Start(plan, raw)
	if err == nil {
		t.Fatalf("Start on main without AllowMain or CreateBranch must fail")
	}
	refusal, ok := err.(*domain.Refusal)
	if !ok || refusal.Code != domain.RefusalInvalidKnownField {
		t.Fatalf("err=%v want RefusalInvalidKnownField", err)
	}
}

func TestBranchGuardrailCreatesBranch(t *testing.T) {
	root := missionServiceFixture(t)
	service := openMissionService(t, root)

	plan, raw := stressPlan()
	plan.Title = "Branch creation test"
	plan.AllowMain = false
	plan.CreateBranch = true

	started, err := service.Start(plan, raw)
	if err != nil {
		t.Fatalf("Start with CreateBranch failed: %v", err)
	}

	ws, err := discovery.Open(root)
	if err != nil {
		t.Fatal(err)
	}
	bundle, err := Load(ws, started.Ref)
	if err != nil {
		t.Fatal(err)
	}
	if !strings.HasPrefix(bundle.Baseline.Branch, "feat/") {
		t.Fatalf("bundle.Baseline.Branch=%q want prefix feat/", bundle.Baseline.Branch)
	}
}

func TestListMissionsReturnsSummaries(t *testing.T) {
	root := missionServiceFixture(t)
	service := openMissionService(t, root)

	plan, raw := stressPlan()
	plan.Title = "List summary test"
	started, err := service.Start(plan, raw)
	if err != nil {
		t.Fatalf("Start failed: %v", err)
	}

	res, err := service.ListMissions("")
	if err != nil {
		t.Fatalf("ListMissions failed: %v", err)
	}
	if len(res.Missions) == 0 {
		t.Fatalf("ListMissions returned 0 missions")
	}
	found := false
	for _, m := range res.Missions {
		if m.Ref == started.Ref {
			found = true
			if m.Status != "active" {
				t.Fatalf("mission status=%q want active", m.Status)
			}
		}
	}
	if !found {
		t.Fatalf("ListMissions did not include started mission %s", started.Ref)
	}
}

func TestRecordEvidenceFromStructuredTestOutput(t *testing.T) {
	root := missionServiceFixture(t)
	service := openMissionService(t, root)

	plan, raw := stressPlan()
	plan.Title = "Evidence from test"
	started, err := service.Start(plan, raw)
	if err != nil {
		t.Fatalf("Start failed: %v", err)
	}

	// Write mock test json output
	testJSON := filepath.Join(root, "test-output.json")
	mockJSON := `{"Time":"2026-08-30T10:00:00Z","Action":"pass","Test":"TestFeatureUnit"}
{"Time":"2026-08-30T10:00:01Z","Action":"pass","Test":"TestFeatureIntegration"}
`
	if err := os.WriteFile(testJSON, []byte(mockJSON), 0o644); err != nil {
		t.Fatal(err)
	}

	res, err := service.RecordEvidence(started.Ref, "", nil, testJSON)
	if err != nil {
		t.Fatalf("RecordEvidence with --from failed: %v", err)
	}
	if res.Operation != "evidence.record" {
		t.Fatalf("Operation=%q want evidence.record", res.Operation)
	}

	ws, err := discovery.Open(root)
	if err != nil {
		t.Fatal(err)
	}
	bundle, err := Load(ws, started.Ref)
	if err != nil {
		t.Fatal(err)
	}
	if len(bundle.Evidence) == 0 {
		t.Fatalf("Evidence list empty")
	}
	evDoc := bundle.Evidence[len(bundle.Evidence)-1].Document
	if evDoc == nil {
		t.Fatalf("Evidence Document is nil")
	}
	if len(evDoc.Checks) != 2 {
		t.Fatalf("Evidence checks count=%d want 2", len(evDoc.Checks))
	}
}
