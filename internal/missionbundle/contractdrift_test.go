package missionbundle

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/alexsmedile/spectacular/v2/internal/discovery"
	"github.com/alexsmedile/spectacular/v2/internal/domain"
)

// A bound Contract's fingerprint exists to stop the Contract shifting under work
// in flight. Once a Mission is completed nothing is in flight, so the same change
// is a fact about the Contract rather than a defect in the Mission — and refusing
// it offered no legal correction, because a completed Mission is never rewritten
// to satisfy it. Every other status must keep refusing exactly as before.
func TestContractDriftRefusesLiveMissionsAndNoticesCompletedOnes(t *testing.T) {
	root, err := filepath.Abs(filepath.Join("..", ".."))
	if err != nil {
		t.Fatal(err)
	}
	ws, err := discovery.Open(root)
	if err != nil {
		t.Fatal(err)
	}
	source, err := Load(ws, "M9")
	if err != nil {
		t.Fatal(err)
	}

	for _, status := range []string{"defined", "active", "blocked", "completed"} {
		t.Run("matching/"+status, func(t *testing.T) {
			bundle := cloneBundle(t, source)
			bundle.Status = status
			if err := validateContract(ws, bundle); err != nil {
				t.Fatalf("status %s with a matching fingerprint must validate, got %v", status, err)
			}
			if bundle.contractDrift != "" {
				t.Fatalf("status %s recorded drift %q with a matching fingerprint", status, bundle.contractDrift)
			}
			for _, notice := range bundle.Notices() {
				if strings.HasPrefix(notice, "contract-drift:") {
					t.Fatalf("status %s reported %q with a matching fingerprint", status, notice)
				}
			}
		})

		t.Run("drifted/"+status, func(t *testing.T) {
			bundle := cloneBundle(t, source)
			bundle.Status = status
			bound := "sha256:" + strings.Repeat("a", 64)
			bundle.Contract.Fingerprint = bound

			err := validateContract(ws, bundle)
			if status != "completed" {
				if err == nil {
					t.Fatalf("status %s must refuse a drifted Contract", status)
				}
				refusal, ok := err.(*domain.Refusal)
				if !ok {
					t.Fatalf("status %s returned %T, want a typed refusal", status, err)
				}
				if refusal.Code != domain.RefusalStaleFingerprint || refusal.Field != "contract.fingerprint" {
					t.Fatalf("status %s refused %s on %q, want %s on contract.fingerprint",
						status, refusal.Code, refusal.Field, domain.RefusalStaleFingerprint)
				}
				if bundle.contractDrift != "" {
					t.Fatalf("status %s recorded drift while refusing", status)
				}
				return
			}

			if err != nil {
				t.Fatalf("a completed Mission must report drift, not refuse it; got %v", err)
			}
			if bundle.contractDrift == "" || bundle.contractDrift == bound {
				t.Fatalf("completed Mission recorded drift %q, want the Contract's actual digest", bundle.contractDrift)
			}
			notice := driftNotice(t, bundle)
			for _, want := range []string{bundle.Contract.Ref, bound, bundle.contractDrift} {
				if !strings.Contains(notice, want) {
					t.Fatalf("notice %q omits %q", notice, want)
				}
			}
			if !strings.Contains(notice, "after this Mission completed") {
				t.Fatalf("notice %q does not say when the Contract changed", notice)
			}
		})
	}
}

// The whole point of the gate is that a reader is told rather than blocked, so a
// completed Mission with a drifted Contract must still report valid through the
// full validator run — not merely survive validateContract in isolation.
func TestCompletedMissionWithDriftedContractStaysValid(t *testing.T) {
	root, err := filepath.Abs(filepath.Join("..", ".."))
	if err != nil {
		t.Fatal(err)
	}
	ws, err := discovery.Open(root)
	if err != nil {
		t.Fatal(err)
	}
	// Loaded rather than cloned: cloneBundle round-trips through JSON, which does
	// not reproduce the frozen envelope byte-for-byte, so a clone fails the
	// activation-fingerprint validator for reasons unrelated to Contract drift.
	bundle, err := Load(ws, "M9")
	if err != nil {
		t.Fatal(err)
	}
	if bundle.Status != "completed" {
		t.Fatalf("M9 status is %q, want completed; this test needs a completed Mission", bundle.Status)
	}
	bundle.Contract.Fingerprint = "sha256:" + strings.Repeat("b", 64)

	check, err := Validate(ws, bundle)
	if err != nil {
		t.Fatalf("Validate refused a completed Mission over Contract drift: %v", err)
	}
	if !check.Valid {
		t.Fatal("check reports invalid; drift must not fail validation")
	}
	var found bool
	for _, notice := range check.Notices {
		if strings.HasPrefix(notice, "contract-drift:") {
			found = true
		}
	}
	if !found {
		t.Fatalf("notices %v carry no contract-drift entry", check.Notices)
	}
}

// The failure this gate exists to fix, reproduced against a real amendment rather
// than a synthetic fingerprint. Rewriting one Gap's `blocked_on:` to `resolution:`
// on the Contract M7, M8, and M9 are bound to refused all three at once, with no
// legal correction available: the Gap was resolved in fact by M9 and could not be
// written down. The amendment now reports and every bound Mission stays valid.
func TestAmendingABoundContractNoLongerBreaksCompletedMissions(t *testing.T) {
	root, err := filepath.Abs(filepath.Join("..", ".."))
	if err != nil {
		t.Fatal(err)
	}
	contract := filepath.Join(root, ".spectacular", "contracts",
		"CC-projsurf-derived-state-and-dependency-shape.md")
	original, err := os.ReadFile(contract)
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() {
		if writeErr := os.WriteFile(contract, original, 0o644); writeErr != nil {
			t.Fatalf("restore the Contract: %v", writeErr)
		}
	})

	const openGap = "    blocked_on: A decision on whether removing them belongs with the Proposal schema work or with a separate cleanup."
	if !strings.Contains(string(original), openGap) {
		t.Skip("dead-v1-governance-code no longer reads blocked_on:; the amendment this test simulates has landed")
	}
	amended := strings.Replace(string(original), openGap,
		"    resolution: Closed by M9 as a separate cleanup, which was the decision this Gap asked for.", 1)
	if err := os.WriteFile(contract, []byte(amended), 0o644); err != nil {
		t.Fatal(err)
	}

	ws, err := discovery.Open(root)
	if err != nil {
		t.Fatal(err)
	}
	for _, ref := range []string{"M7", "M8", "M9"} {
		bundle, loadErr := Load(ws, ref)
		if loadErr != nil {
			t.Fatalf("%s: %v", ref, loadErr)
		}
		check, checkErr := Validate(ws, bundle)
		if checkErr != nil {
			t.Fatalf("%s refused after its bound Contract was amended: %v", ref, checkErr)
		}
		if !check.Valid {
			t.Fatalf("%s reports invalid after its bound Contract was amended", ref)
		}
		driftNotice(t, bundle)
	}
}

func driftNotice(t *testing.T, b *Bundle) string {
	t.Helper()
	for _, notice := range b.Notices() {
		if strings.HasPrefix(notice, "contract-drift:") {
			return notice
		}
	}
	t.Fatalf("notices %v carry no contract-drift entry", b.Notices())
	return ""
}
