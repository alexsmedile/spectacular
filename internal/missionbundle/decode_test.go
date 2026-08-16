package missionbundle

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/alexsmedile/spectacular/v2/internal/discovery"
	"github.com/alexsmedile/spectacular/v2/internal/workspace"
	"go.yaml.in/yaml/v3"
)

func TestSelfHostedMissionGoldenDecodingAndRoundTrip(t *testing.T) {
	root, err := filepath.Abs(filepath.Join("..", ".."))
	if err != nil {
		t.Fatal(err)
	}
	ws, err := discovery.Open(root)
	if err != nil {
		t.Fatal(err)
	}

	for _, ref := range []string{"M5", "M6"} {
		t.Run(ref, func(t *testing.T) {
			bundle, loadErr := Load(ws, ref)
			if loadErr != nil {
				t.Fatal(loadErr)
			}
			if bundle.Ref != ref || bundle.Legacy || bundle.Validation.Schema != Schema || bundle.document == nil {
				t.Fatalf("golden bundle=%+v", bundle)
			}
			assertCanonicalStable(t, bundle.document)
		})
	}

	m5, err := Load(ws, "M5")
	if err != nil {
		t.Fatal(err)
	}
	if len(m5.Reviews) != 1 || m5.Reviews[0].File != "reviews/RV1-independent-review.md" || m5.Reviews[0].Document == nil {
		t.Fatalf("resolved M5 reviews=%+v", m5.Reviews)
	}
	review := m5.Reviews[0].Document
	if review.Ref != m5.Reviews[0].Ref || review.ID != m5.Reviews[0].ID || review.Mission != "M5" || review.Reviewed.ActivationFingerprint != m5.Activation.Fingerprint {
		t.Fatalf("typed Review does not match source pointer: pointer=%+v review=%+v", m5.Reviews[0], review)
	}
	if !strings.Contains(review.Body, "# Review") {
		t.Fatalf("Review Markdown body was not retained: %q", review.Body)
	}
	assertCanonicalStable(t, review.document)

	legacy, err := Load(ws, "M3")
	if err != nil {
		t.Fatal(err)
	}
	const source = "Proposal:01a007f6-d88a-7922-a16e-abd2262feda4"
	if !legacy.Legacy || legacy.Source != source || legacy.document.Record.Source.String() != source {
		t.Fatalf("legacy source pointer lost: legacy=%v source=%q", legacy.Legacy, legacy.Source)
	}
	assertCanonicalStable(t, legacy.document)
}

func TestExpandedBundleResolvesRecordsWithoutExpandingPointers(t *testing.T) {
	root := t.TempDir()
	missionPath := filepath.Join(root, "MISSION.md")
	writeFixture(t, missionPath, expandedMissionFixture)
	writeFixture(t, filepath.Join(root, "objectives", "O1.md"), expandedObjectiveFixture)
	writeFixture(t, filepath.Join(root, "runs", "R1.md"), expandedRunFixture)
	writeFixture(t, filepath.Join(root, "reviews", "RV1.md"), expandedReviewFixture)

	doc, err := workspace.ReadFile(missionPath)
	if err != nil {
		t.Fatal(err)
	}
	bundle, err := decode(&discovery.Workspace{}, discovery.Entry{Path: "MISSION.md", Absolute: missionPath, Document: doc})
	if err != nil {
		t.Fatal(err)
	}
	if bundle.Source != "Proposal:01a007f6-d88a-7922-a16e-abd2262feda4" {
		t.Fatalf("Mission source=%q", bundle.Source)
	}
	if got := bundle.Objectives[0]; got.File != "objectives/O1.md" || got.document == nil || !strings.Contains(got.Body, "objective body") {
		t.Fatalf("resolved Objective=%+v", got)
	}
	if got := bundle.Runs[0]; got.File != "runs/R1.md" || got.Title != "Expanded run" || got.document == nil || !strings.Contains(got.Body, "run body") {
		t.Fatalf("resolved Run=%+v", got)
	}
	if got := bundle.Reviews[0]; got.File != "reviews/RV1.md" || got.Verdict != "pass" || got.Document == nil || got.Document.Path != got.File || !strings.Contains(got.Document.Body, "review body") {
		t.Fatalf("resolved Review pointer=%+v", got)
	}

	for name, source := range map[string]*workspace.Document{
		"mission":   bundle.document,
		"objective": bundle.Objectives[0].document,
		"run":       bundle.Runs[0].document,
		"review":    bundle.Reviews[0].Document.document,
	} {
		t.Run(name, func(t *testing.T) {
			if source.Unknown["extension"] == nil {
				t.Fatal("unknown extension field was not retained")
			}
			assertCanonicalStable(t, source)
		})
	}

	pointers, err := yaml.Marshal(bundle.Reviews)
	if err != nil {
		t.Fatal(err)
	}
	if strings.Contains(string(pointers), "document:") || !strings.Contains(string(pointers), "file: reviews/RV1.md") {
		t.Fatalf("Mission pointer serialization expanded the Review document:\n%s", pointers)
	}
}

func assertCanonicalStable(t *testing.T, doc *workspace.Document) {
	t.Helper()
	first, err := workspace.Canonical(doc)
	if err != nil {
		t.Fatal(err)
	}
	roundTripped, err := workspace.Parse(first)
	if err != nil {
		t.Fatal(err)
	}
	second, err := workspace.Canonical(roundTripped)
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(first, second) {
		t.Fatalf("canonical round trip drifted:\n%s\n--- second ---\n%s", first, second)
	}
}

func writeFixture(t *testing.T, path, content string) {
	t.Helper()
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		t.Fatal(err)
	}
}

const expandedMissionFixture = `---
type: Mission
id: 01a009ff-ce94-724e-a6f8-66783f1a4003
ref: M9
title: Expanded bundle
status: active
source: Proposal:01a007f6-d88a-7922-a16e-abd2262feda4
owner: Alex
outcome: Exercise every promoted record.
review: independent
contract: {ref: CC-test, fingerprint: "sha256:test"}
completion: [{claim: typed, pass_boundary: typed, proof_requirement: typed}]
objectives: [{ref: O1, id: 01a009ff-ce94-7249-accc-9c2a089d3080, file: objectives/O1.md}]
runs: [{ref: R1, id: 01a009ff-ce94-7323-a942-754ef422e264, file: runs/R1.md}]
reviews: [{ref: RV1, id: 01a00a33-87b8-7e0e-9100-450696ad1e80, file: reviews/RV1.md, verdict: pass}]
validation: {schema: mission.v2, mode: cli}
authority: {operator: [inspect], requires_owner: [expand-scope]}
scope: {mechanical: [internal/missionbundle/], semantic: [typed bundle]}
repair_budget: 1
dependencies: []
gaps: []
stops: [leave scope]
extension: {retained: true}
---
# Expanded Mission
`

const expandedObjectiveFixture = `---
type: Objective
id: 01a009ff-ce94-7249-accc-9c2a089d3080
ref: O1
title: Expanded objective
status: implemented
mission: M9
outcome: Resolve the Objective.
after: []
claims: [typed]
extension: {retained: objective}
---
# objective body
`

const expandedRunFixture = `---
type: Run
id: 01a009ff-ce94-7323-a942-754ef422e264
ref: R1
title: Expanded run
status: completed
mission: M9
operator: Alex
started_at: "2026-08-16T10:31:30Z"
current_objective: O1
repairs: 0
extension: {retained: run}
---
# run body
`

const expandedReviewFixture = `---
type: Review
id: 01a00a33-87b8-7e0e-9100-450696ad1e80
ref: RV1
title: Expanded review
status: passed
mission: M9
created: "2026-08-16T10:52:24Z"
reviewed: {commit: 4074708c26c1158f4eb778b55c86aabe80979e76, tree: a6a9344ba74c5a7d3e5bf1b28754ee2905bad01d, activation_fingerprint: "sha256:test"}
reviewer: {actor: reviewer, operator: operator, relation_to_operator: independent, implemented_reviewed_scope: false, independence_basis: separate, evidence: [task:test]}
claims: [{claim: typed, verdict: pass}]
findings: []
limitations: []
extension: {retained: review}
---
# review body
`
