// Package missionbundle owns the single typed representation of compact and
// expanded Spectacular Missions.
package missionbundle

import (
	"github.com/alexsmedile/spectacular/v2/internal/discovery"
	"github.com/alexsmedile/spectacular/v2/internal/workspace"
)

const Schema = "mission.v2"

type Binding struct {
	Ref         string `yaml:"ref" json:"ref"`
	Fingerprint string `yaml:"fingerprint" json:"fingerprint"`
}

type Baseline struct {
	Commit string `yaml:"commit" json:"commit"`
	Branch string `yaml:"branch" json:"branch"`
}

type Criterion struct {
	Claim            string `yaml:"claim" json:"claim"`
	PassBoundary     string `yaml:"pass_boundary" json:"pass_boundary"`
	ProofRequirement string `yaml:"proof_requirement" json:"proof_requirement"`
}

type Objective struct {
	Ref     string   `yaml:"ref" json:"ref"`
	ID      string   `yaml:"id" json:"id"`
	Source  string   `yaml:"-" json:"source,omitempty"`
	Outcome string   `yaml:"outcome,omitempty" json:"outcome,omitempty"`
	Status  string   `yaml:"status,omitempty" json:"status,omitempty"`
	After   []string `yaml:"after,omitempty" json:"after,omitempty"`
	Claims  []string `yaml:"claims,omitempty" json:"claims,omitempty"`
	File    string   `yaml:"file,omitempty" json:"file,omitempty"`
	Body    string   `yaml:"-" json:"-"`

	document *workspace.Document
}

type Run struct {
	Ref              string `yaml:"ref" json:"ref"`
	ID               string `yaml:"id" json:"id"`
	Source           string `yaml:"-" json:"source,omitempty"`
	Title            string `yaml:"title,omitempty" json:"title,omitempty"`
	Status           string `yaml:"status" json:"status"`
	Operator         string `yaml:"operator" json:"operator"`
	StartedAt        string `yaml:"started_at" json:"started_at"`
	CurrentObjective string `yaml:"current_objective" json:"current_objective"`
	Repairs          int    `yaml:"repairs" json:"repairs"`
	File             string `yaml:"file,omitempty" json:"file,omitempty"`
	Body             string `yaml:"-" json:"-"`

	document *workspace.Document
}

type Activation struct {
	By          string `yaml:"by" json:"by"`
	At          string `yaml:"at" json:"at"`
	Fingerprint string `yaml:"fingerprint" json:"fingerprint"`
}

type Validation struct {
	Schema string `yaml:"schema" json:"schema"`
	Mode   string `yaml:"mode" json:"mode"`
}

type Authority struct {
	Operator      []string `yaml:"operator" json:"operator"`
	RequiresOwner []string `yaml:"requires_owner" json:"requires_owner"`
}

type Scope struct {
	Mechanical []string `yaml:"mechanical" json:"mechanical"`
	Semantic   []string `yaml:"semantic" json:"semantic"`
}

type ReviewPointer struct {
	Ref      string  `yaml:"ref" json:"ref"`
	ID       string  `yaml:"id" json:"id"`
	File     string  `yaml:"file" json:"file"`
	Verdict  string  `yaml:"verdict" json:"verdict"`
	Document *Review `yaml:"-" json:"document,omitempty"`
}

type Reviewed struct {
	Commit                string `yaml:"commit" json:"commit"`
	Tree                  string `yaml:"tree" json:"tree"`
	ActivationFingerprint string `yaml:"activation_fingerprint" json:"activation_fingerprint"`
}

type ClaimVerdict struct {
	Claim   string `yaml:"claim" json:"claim"`
	Verdict string `yaml:"verdict" json:"verdict"`
}

// Review is the typed view of a Review record resolved from a Mission pointer.
// Its source Document stays attached so canonical rewrites preserve opaque YAML
// fields and the Markdown body without expanding the Mission's pointer.
type Review struct {
	ID          string         `json:"id"`
	Ref         string         `json:"ref"`
	Title       string         `json:"title"`
	Status      string         `json:"status"`
	Source      string         `json:"source,omitempty"`
	Mission     string         `json:"mission"`
	Created     string         `json:"created,omitempty"`
	Reviewed    Reviewed       `json:"reviewed"`
	Reviewer    Reviewer       `json:"reviewer"`
	Claims      []ClaimVerdict `json:"claims"`
	Findings    []string       `json:"findings,omitempty"`
	Limitations []string       `json:"limitations,omitempty"`
	Path        string         `json:"path"`
	Body        string         `json:"-"`

	document *workspace.Document
}

type CompletionRecord struct {
	By             string   `yaml:"by" json:"by"`
	At             string   `yaml:"at" json:"at"`
	Authorization  string   `yaml:"authorization" json:"authorization"`
	ReviewedCommit string   `yaml:"reviewed_commit" json:"reviewed_commit"`
	Review         string   `yaml:"review" json:"review"`
	Limitations    []string `yaml:"limitations,omitempty" json:"limitations,omitempty"`
}

type Reviewer struct {
	Actor                    string   `yaml:"actor" json:"actor"`
	Operator                 string   `yaml:"operator" json:"operator"`
	RelationToOperator       string   `yaml:"relation_to_operator" json:"relation_to_operator"`
	ImplementedReviewedScope bool     `yaml:"implemented_reviewed_scope" json:"implemented_reviewed_scope"`
	IndependenceBasis        string   `yaml:"independence_basis" json:"independence_basis"`
	Evidence                 []string `yaml:"evidence" json:"evidence"`
}

type Bundle struct {
	ID               string            `json:"id"`
	Ref              string            `json:"ref"`
	Title            string            `json:"title"`
	Status           string            `json:"status"`
	Source           string            `json:"source,omitempty"`
	Owner            string            `json:"owner,omitempty"`
	Created          string            `json:"created,omitempty"`
	Updated          string            `json:"updated,omitempty"`
	Contract         Binding           `json:"contract,omitempty"`
	Baseline         *Baseline         `json:"baseline,omitempty"`
	Outcome          string            `json:"outcome,omitempty"`
	Review           string            `json:"review,omitempty"`
	Completion       []Criterion       `json:"completion,omitempty"`
	Objectives       []Objective       `json:"objectives,omitempty"`
	Run              *Run              `json:"run,omitempty"`
	Runs             []Run             `json:"runs,omitempty"`
	Activation       *Activation       `json:"activation,omitempty"`
	Validation       Validation        `json:"validation,omitempty"`
	Authority        Authority         `json:"authority,omitempty"`
	Scope            Scope             `json:"scope,omitempty"`
	RepairBudget     int               `json:"repair_budget,omitempty"`
	Dependencies     []string          `json:"dependencies,omitempty"`
	Gaps             []string          `json:"gaps,omitempty"`
	Stops            []string          `json:"stops,omitempty"`
	Reviews          []ReviewPointer   `json:"reviews,omitempty"`
	CompletionRecord *CompletionRecord `json:"completion_record,omitempty"`
	Path             string            `json:"path"`
	Legacy           bool              `json:"legacy"`
	// State is derived on read so a JSON reader reaches the same conclusion a
	// human reads from the rendered output. It is never decoded from the file
	// and never written back: mutations set canonical fields individually, so
	// this cannot reach the canonical tree.
	State *State `yaml:"-" json:"state,omitempty"`
	Body  string `json:"-"`

	entry    discovery.Entry
	document *workspace.Document
}

type Check struct {
	Ref         string   `json:"ref"`
	Path        string   `json:"path"`
	Schema      string   `json:"schema"`
	Valid       bool     `json:"valid"`
	Legacy      bool     `json:"legacy"`
	Checks      []string `json:"checks"`
	Fingerprint string   `json:"fingerprint,omitempty"`
	// Authority is the decision table the authority-vocabulary validator
	// already resolves and previously discarded. It is derived on read.
	Authority []AuthorityAnswer `json:"authority,omitempty"`
	// Drift is the per-claim audit signal, most-flagged first. Derived on read.
	Drift []ClaimDrift `json:"drift,omitempty"`
}

type Result struct {
	Operation   string   `json:"operation"`
	Ref         string   `json:"ref"`
	Path        string   `json:"path"`
	Fingerprint string   `json:"fingerprint,omitempty"`
	Changed     []string `json:"changed,omitempty"`
}
