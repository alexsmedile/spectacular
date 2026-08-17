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
	ws, gap, contractRef := amendableWorkspace(t)
	root := ws.Root
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
	// Whichever Mission is live is the one that protects its Contract. Naming a
	// specific Mission made this test skip as soon as that Mission completed,
	// which is precisely when it stops being watched.
	var contractRef string
	for _, entry := range ws.Entries {
		if entry.Document == nil || entry.Document.Record.Type != domain.Mission {
			continue
		}
		bundle, loadErr := Load(ws, entry.Document.Record.ID.String())
		if loadErr != nil || bundle.Status != "active" || bundle.Contract.Ref == "" {
			continue
		}
		contractRef = bundle.Contract.Ref
		break
	}
	if contractRef == "" {
		t.Skip("no live Mission is bound to a Contract in this workspace")
	}
	service := Service{Workspace: ws}
	_, err = service.AmendContract(contractRef, "any-gap", "Alex", "text", true)
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
	ws, gap, contractRef := amendableWorkspace(t)
	service := Service{Workspace: ws}
	_, err := service.AmendContract(contractRef, gap, "Alex", "", true)
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
	ws, gap, contractRef := amendableWorkspace(t)
	root := ws.Root
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

// anOpenGap finds a Gap still carrying blocked_on, so a test describes the
// workspace as it is rather than asserting a specific Gap that a later amendment
// would legitimately close. This is for tests that only need an open Gap to
// exist. A test that amends one needs amendableWorkspace instead, because the
// first open Gap in the real workspace may sit on a Contract with a live bound
// Mission, which the amendment path refuses by design.
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

// amendableWorkspace builds a workspace holding one Contract with one open Gap
// and no Mission bound to it.
//
// These tests previously opened the real repository and amended whichever open
// Gap they found first. That made them depend on which Missions happened to be
// live: amending a Contract a live Mission is bound to is refused by design, so
// activating any Mission on that Contract broke three tests of the amendment
// path for a reason that had nothing to do with the amendment path. Skipping
// instead of failing would have been worse — the suite would go green while the
// tests proved nothing.
//
// A fixture makes the precondition something the test establishes rather than
// something it hopes for, so these run unconditionally.
func amendableWorkspace(t *testing.T) (*discovery.Workspace, string, string) {
	t.Helper()
	root := t.TempDir()
	directory := filepath.Join(root, ".spectacular", "contracts")
	if err := os.MkdirAll(directory, 0o755); err != nil {
		t.Fatal(err)
	}
	const contract = `---
type: Contract
id: 01a00aae-8921-7b27-96a9-1a4c175e7dc6
human_ref: CC-fixture
title: Fixture contract
purpose: a contract that no live Mission is bound to
updated: "2026-08-16T00:00:00Z"
gaps:
    - ref: fixture-open-gap
      problem: the amendment path needs an open Gap that is safe to close
      blocked_on: an owner decision that has not been made
---

Fixture Contract body.
`
	if err := os.WriteFile(filepath.Join(directory, "CC-fixture-contract.md"), []byte(contract), 0o644); err != nil {
		t.Fatal(err)
	}
	const manifest = "schema_version: spectacular.workspace.v1\nrecord_roots:\n  - .\nproject_anchor: PROJECT.md\n"
	if err := os.WriteFile(filepath.Join(root, ".spectacular", "workspace.yaml"), []byte(manifest), 0o644); err != nil {
		t.Fatal(err)
	}
	const anchor = `---
type: Anchor
id: 019fe381-5d61-7223-b362-03a5f99a7b14
title: Fixture project
---

Fixture anchor.
`
	if err := os.WriteFile(filepath.Join(root, ".spectacular", "PROJECT.md"), []byte(anchor), 0o644); err != nil {
		t.Fatal(err)
	}
	ws, err := discovery.Open(root)
	if err != nil {
		t.Fatal(err)
	}
	contractRef := string(domain.Contract) + ":01a00aae-8921-7b27-96a9-1a4c175e7dc6"

	// The fixture is only meaningful if the amendment path agrees nothing is live
	// on it. Asking the production code keeps the fixture honest.
	if live, mission := (Service{Workspace: ws}).liveBoundMission(contractRef, ""); live {
		t.Fatalf("fixture workspace unexpectedly carries a live Mission %s", mission)
	}
	return ws, "fixture-open-gap", contractRef
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

// Completion enforces the declaration instead of executing it. A Mission that said
// it would close a Gap has not finished until the Gap is closed, and the refusal has
// to name the command that closes it — the failure the original stale_fingerprint
// refusal made was naming an amend path that did not exist.
func TestCompletionRefusesWhileADeclaredGapIsStillOpen(t *testing.T) {
	root, err := filepath.Abs(filepath.Join("..", ".."))
	if err != nil {
		t.Fatal(err)
	}
	ws, err := discovery.Open(root)
	if err != nil {
		t.Fatal(err)
	}
	gap, contractRef := anOpenGap(t, ws)
	bundle, err := Load(ws, "M9")
	if err != nil {
		t.Fatal(err)
	}
	service := Service{Workspace: ws}

	// Declaring nothing leaves completion exactly as it was.
	if err := service.assertDeclaredGapsClosed(bundle); err != nil {
		t.Fatalf("a Mission declaring no Gaps must be unaffected, got %v", err)
	}

	bundle.Contract.Ref = contractRef
	bundle.ResolvesGaps = []ResolvedGap{{Gap: gap, Resolution: "whatever this Mission promised"}}
	err = service.assertDeclaredGapsClosed(bundle)
	if err == nil {
		t.Fatal("completion must refuse while a declared Gap is still open")
	}
	if !strings.Contains(err.Error(), gap) {
		t.Fatalf("refusal does not name the Gap: %v", err)
	}
	// Recovery is where a refusal states the fix; Error() carries only code, field,
	// and problem, and the CLI prints the correction separately.
	refusal, ok := err.(*domain.Refusal)
	if !ok {
		t.Fatalf("returned %T, want a typed refusal", err)
	}
	if !strings.Contains(refusal.Recovery, "contract amend") || !strings.Contains(refusal.Recovery, gap) {
		t.Fatalf("correction does not name the command that closes it: %q", refusal.Recovery)
	}

	// A Gap that already carries a resolution satisfies the declaration.
	closed, closedContract := aClosedGap(t, ws)
	bundle.Contract.Ref = closedContract
	bundle.ResolvesGaps = []ResolvedGap{{Gap: closed, Resolution: "already written"}}
	if err := service.assertDeclaredGapsClosed(bundle); err != nil {
		t.Fatalf("a closed Gap must satisfy the declaration, got %v", err)
	}
}

func aClosedGap(t *testing.T, ws *discovery.Workspace) (string, string) {
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
			if gap.Resolution != "" {
				return gap.Ref, ref
			}
		}
	}
	t.Skip("no Contract in the workspace carries a closed Gap")
	return "", ""
}
