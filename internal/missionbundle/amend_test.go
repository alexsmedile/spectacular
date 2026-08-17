package missionbundle

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/alexsmedile/spectacular/v2/internal/discovery"
	"github.com/alexsmedile/spectacular/v2/internal/domain"
	"github.com/alexsmedile/spectacular/v2/internal/governance"
	"go.yaml.in/yaml/v3"
)

// An amendment may reach the gaps: block and editorial frontmatter and nothing
// else. Without this guard "amendment" is "rewrite with extra steps", and the only
// thing separating the two is a reason string nothing validates.
func TestAmendmentRefusesFieldsThatStateWhatWasAgreed(t *testing.T) {
	const contract = `---
type: Contract
id: 01a00aae-8921-7b27-96a9-1a4c175e7dc6
purpose: original purpose
updated: "2026-08-16T00:00:00Z"
gaps:
  - ref: g1
    problem: p
    blocked_on: something
---
body
`
	tests := []struct {
		name    string
		mutate  func(string) string
		refused bool
	}{
		{"gap closed", func(s string) string {
			return strings.Replace(s, "    blocked_on: something", "    resolution: done", 1)
		}, false},
		{"editorial updated", func(s string) string {
			return strings.Replace(s, `updated: "2026-08-16T00:00:00Z"`, `updated: "2026-08-17T00:00:00Z"`, 1)
		}, false},
		{"purpose rewritten", func(s string) string {
			return strings.Replace(s, "purpose: original purpose", "purpose: something else", 1)
		}, true},
		{"field added", func(s string) string {
			return strings.Replace(s, "gaps:", "outcome: invented\ngaps:", 1)
		}, true},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			err := assertOnlyAmendableFieldsChanged(contract, test.mutate(contract))
			if test.refused && err == nil {
				t.Fatal("must refuse a change outside the amendable fields")
			}
			if !test.refused && err != nil {
				t.Fatalf("must permit the change, got %v", err)
			}
		})
	}
}

// The rewrite is textual so it leaves prose it did not change untouched. An
// amendment that reflows the whole Contract is not reviewable.
func TestGapRewriteTouchesOnlyTheNamedGap(t *testing.T) {
	const contract = `---
type: Contract
gaps:
  - ref: first
    problem: keep me exactly as I am
    blocked_on: an old reason
  - ref: second
    problem: also untouched
    blocked_on: >-
      a folded reason
      over two lines
---
body stays
`
	amended, err := rewriteGap(contract, "second", "closed for a stated reason")
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(amended, "problem: keep me exactly as I am") ||
		!strings.Contains(amended, "blocked_on: an old reason") {
		t.Fatal("the untargeted Gap was modified")
	}
	if strings.Contains(amended, "a folded reason") {
		t.Fatal("the targeted Gap's blocked_on block scalar was not removed")
	}
	if !strings.Contains(amended, "resolution: >-") || !strings.Contains(amended, "closed for a stated reason") {
		t.Fatalf("resolution was not written:\n%s", amended)
	}
	if !strings.HasSuffix(amended, "body stays\n") {
		t.Fatal("the Markdown body was modified")
	}

	// A Gap that is already closed has no blocked_on to rewrite, and a Gap that is
	// absent cannot be located; both are refusals rather than silent no-ops.
	if _, err := rewriteGap(contract, "missing", "x"); err == nil {
		t.Fatal("an absent Gap must refuse")
	}
}

// Long resolution text is emitted as a block scalar because a bare colon in a
// plain scalar silently breaks the document — a failure this project hit three
// times while authoring records by hand.
func TestResolutionIsEmittedAsAParseableBlockScalar(t *testing.T) {
	const contract = "---\ntype: Contract\ngaps:\n  - ref: g\n    blocked_on: x\n---\nbody\n"
	risky := "Closed by M9: the framing was too narrow, and a walk found more. " +
		strings.Repeat("Additional prose that pushes this past any sensible line width. ", 4)
	amended, err := rewriteGap(contract, "g", risky)
	if err != nil {
		t.Fatal(err)
	}
	gaps := struct {
		Gaps []contractGap `yaml:"gaps"`
	}{}
	if err := yamlUnmarshalFrontmatter(t, amended, &gaps); err != nil {
		t.Fatalf("the amended Contract does not parse: %v\n%s", err, amended)
	}
	if len(gaps.Gaps) != 1 {
		t.Fatalf("decoded %d Gaps", len(gaps.Gaps))
	}
	if !strings.Contains(gaps.Gaps[0].Resolution, "the framing was too narrow") {
		t.Fatalf("resolution round-tripped as %q", gaps.Gaps[0].Resolution)
	}
	if gaps.Gaps[0].BlockedOn != "" {
		t.Fatal("blocked_on survived the rewrite")
	}
}

// --dry-run must describe the amendment and leave the tree byte-identical. The
// preview is the only thing the owner reads before authorizing, so it has to name
// every file the real run would write.
func TestDryRunDescribesTheAmendmentAndWritesNothing(t *testing.T) {
	root, err := filepath.Abs(filepath.Join("..", ".."))
	if err != nil {
		t.Fatal(err)
	}
	ws, err := discovery.Open(root)
	if err != nil {
		t.Fatal(err)
	}
	gap, contractRef := anOpenGap(t, ws)
	before := treeDigest(t, filepath.Join(root, ".spectacular"))

	service := Service{Workspace: ws}
	result, err := service.AmendContract(contractRef, gap, "Alex", "an owner-supplied resolution", true)
	if err != nil {
		t.Fatal(err)
	}
	if !result.DryRun {
		t.Fatal("result does not report itself as a dry run")
	}
	if result.From == result.To || result.To == "" {
		t.Fatalf("dry run reported from=%s to=%s", result.From, result.To)
	}
	if len(result.Changed) < 2 {
		t.Fatalf("dry run names %d files, want the Contract and its log", len(result.Changed))
	}
	if result.Log == "" || !strings.HasSuffix(result.Log, ".amendments.md") {
		t.Fatalf("log path is %q", result.Log)
	}
	if after := treeDigest(t, filepath.Join(root, ".spectacular")); after != before {
		t.Fatal("a dry run wrote to the workspace")
	}
}

// A live bound Mission means the Contract still constrains work in flight, which is
// exactly what the fingerprint exists to protect.
func TestAmendmentRefusesWhileABoundMissionIsLive(t *testing.T) {
	root, err := filepath.Abs(filepath.Join("..", ".."))
	if err != nil {
		t.Fatal(err)
	}
	ws, err := discovery.Open(root)
	if err != nil {
		t.Fatal(err)
	}
	// M11 is active and bound to CC-missioncli, so that Contract is protected.
	live, err := Load(ws, "M11")
	if err != nil {
		t.Fatal(err)
	}
	if live.Status != "active" {
		t.Skipf("M11 status is %q; this test needs a live Mission", live.Status)
	}
	service := Service{Workspace: ws}
	_, err = service.AmendContract(live.Contract.Ref, "any-gap", "Alex", "text", true)
	if err == nil {
		t.Fatal("amending a Contract with a live bound Mission must refuse")
	}
	if !strings.Contains(err.Error(), "live") {
		t.Fatalf("refusal does not name the live Mission: %v", err)
	}
}

// Without a declaring Mission and without an explicit owner override there is no
// approved wording to write, so the amendment refuses rather than inventing text.
func TestAmendmentRefusesWithNoApprovedWording(t *testing.T) {
	root, err := filepath.Abs(filepath.Join("..", ".."))
	if err != nil {
		t.Fatal(err)
	}
	ws, err := discovery.Open(root)
	if err != nil {
		t.Fatal(err)
	}
	gap, contractRef := anOpenGap(t, ws)
	service := Service{Workspace: ws}
	_, err = service.AmendContract(contractRef, gap, "Alex", "", true)
	if err == nil {
		t.Fatal("an amendment with no declared resolution and no override must refuse")
	}
	refusal, ok := err.(*domain.Refusal)
	if !ok {
		t.Fatalf("returned %T, want a typed refusal", err)
	}
	if refusal.Field != "gap" {
		t.Fatalf("refused on %q, want gap", refusal.Field)
	}
}

// The log is append-only provenance: a second amendment adds an entry rather than
// replacing one, and the recorded `to` is the digest the amendment produced.
func TestAmendmentLogAppendsAndRecordsBothFingerprints(t *testing.T) {
	root := t.TempDir()
	path := "contracts/CC-x.amendments.md"
	if err := os.MkdirAll(filepath.Join(root, "contracts"), 0o755); err != nil {
		t.Fatal(err)
	}
	first, err := appendAmendmentLog(root, path, amendmentEntry{
		At: "t1", By: "Alex", Mission: "M9", Gap: "g1", Contract: "Contract:x",
		Source: "mission-declared", From: "sha256:aaa", To: "sha256:bbb",
	})
	if err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(root, filepath.FromSlash(path)), first, 0o644); err != nil {
		t.Fatal(err)
	}
	second, err := appendAmendmentLog(root, path, amendmentEntry{
		At: "t2", By: "Alex", Mission: "M10", Gap: "g2", Contract: "Contract:x",
		Source: "owner-supplied", From: "sha256:bbb", To: "sha256:ccc",
	})
	if err != nil {
		t.Fatal(err)
	}
	body := string(second)
	for _, want := range []string{"gap: g1", "gap: g2", "sha256:aaa", "sha256:ccc",
		"source: mission-declared", "source: owner-supplied"} {
		if !strings.Contains(body, want) {
			t.Fatalf("log lost %q after the second amendment:\n%s", want, body)
		}
	}
	if strings.Count(body, "# Amendments") != 1 {
		t.Fatal("the header was written twice")
	}
}

// Fault injection at every write boundary: an interrupted amendment must leave the
// Contract, the log, and every re-pointed Mission untouched, because a half-applied
// amendment is a Contract whose fingerprint no Mission agrees with.
func TestAmendmentRollsBackAtEveryWriteBoundary(t *testing.T) {
	root, err := filepath.Abs(filepath.Join("..", ".."))
	if err != nil {
		t.Fatal(err)
	}
	ws, err := discovery.Open(root)
	if err != nil {
		t.Fatal(err)
	}
	gap, contractRef := anOpenGap(t, ws)
	before := treeDigest(t, filepath.Join(root, ".spectacular"))

	// Derive the boundary count from the amendment itself rather than hardcoding it,
	// so adding a written file extends the coverage instead of silently escaping it.
	preview := Service{Workspace: ws}
	planned, err := preview.AmendContract(contractRef, gap, "Alex", "an owner-supplied resolution", true)
	if err != nil {
		t.Fatal(err)
	}
	for boundary := 0; boundary < len(planned.Changed); boundary++ {
		service := Service{Workspace: ws}
		service.ApplyTransaction = func(transactionRoot, key string, changes []governance.FileChange) error {
			return governance.ApplyTransactionWithFailure(transactionRoot, key, changes, boundary)
		}
		if _, err := service.AmendContract(contractRef, gap, "Alex", "an owner-supplied resolution", false); err == nil {
			t.Fatalf("boundary %d: injected failure did not surface", boundary)
		}
		if err := governance.RecoverTransactions(root); err != nil {
			t.Fatalf("boundary %d: recovery failed: %v", boundary, err)
		}
		if after := treeDigest(t, filepath.Join(root, ".spectacular")); after != before {
			t.Fatalf("boundary %d: workspace was left modified after recovery", boundary)
		}
	}
}

// anOpenGap finds a Gap still carrying blocked_on, so these tests describe the
// workspace as it is rather than asserting a specific Gap that a later amendment
// would legitimately close.
func anOpenGap(t *testing.T, ws *discovery.Workspace) (string, string) {
	t.Helper()
	for _, entry := range ws.Entries {
		if entry.Document == nil || entry.Document.Record.Type != domain.Contract {
			continue
		}
		ref := string(domain.Contract) + ":" + entry.Document.Record.ID.String()
		gaps, err := ContractGaps(ws, ref)
		if err != nil {
			continue
		}
		for _, gap := range gaps {
			if gap.BlockedOn != "" && gap.Resolution == "" {
				return gap.Ref, ref
			}
		}
	}
	t.Skip("no Contract in the workspace carries an open Gap")
	return "", ""
}

// yamlUnmarshalFrontmatter decodes a record's frontmatter so a test can prove the
// amended Contract still parses. Written here rather than reusing the workspace
// decoder because the fixtures are minimal Contracts, not full records.
func yamlUnmarshalFrontmatter(t *testing.T, record string, out any) error {
	t.Helper()
	_, frontmatter, _, err := splitRecord(record)
	if err != nil {
		return err
	}
	return yaml.Unmarshal([]byte(frontmatter), out)
}
