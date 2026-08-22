// Package missionbundle owns the single typed representation of compact and
// expanded Spectacular Missions.
package missionbundle

import (
	"sort"

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

// Fallback records a seriously-considered rejected approach frozen at plan
// time. If repair exhausts mid-Run, the recorded fallbacks are surfaced to the
// owner alongside the failure that consumed the budget.
type Fallback struct {
	Approach        string `yaml:"approach" json:"approach"`
	RejectedBecause string `yaml:"rejected_because" json:"rejected_because"`
	InvalidatedIf   string `yaml:"invalidated_if" json:"invalidated_if"`
	Recommendation  bool   `yaml:"recommendation,omitempty" json:"recommendation,omitempty"`
}

// ResolvedGap declares one Gap on the bound Contract that this Mission closes,
// and the resolution text it will write. Both halves are frozen: the ref so a
// Mission cannot acquire the authority to amend a Contract after activation, and
// the text so the owner approves the exact wording at the activation gate rather
// than whatever completion happens to compose.
type ResolvedGap struct {
	Gap        string `yaml:"gap" json:"gap"`
	Resolution string `yaml:"resolution" json:"resolution"`
}

// Ask is one distinct thing the owner asked for, carrying what the plan decided
// to do about it. The disposition is authored, never inferred: a validator that
// guessed would produce false refusals on correct plans and false confidence on
// plausible-sounding ones.
type Ask struct {
	Ask string `yaml:"ask" json:"ask"`
	// Disposition is covered, deferred, or declined.
	Disposition string `yaml:"disposition" json:"disposition"`
	// Claims names the completion claims answering a covered ask. Every named
	// claim must exist; whether it genuinely answers the ask is review's job.
	Claims []string `yaml:"claims,omitempty" json:"claims,omitempty"`
	// Reason states why a deferred or declined ask is not being done.
	Reason string `yaml:"reason,omitempty" json:"reason,omitempty"`
}

// Request records what produced the plan, so a reader can tell what was asked
// for and what was dropped. Nothing else in the record holds the original ask:
// `outcome` is the agent's interpretation of it, and once frozen, no gate checks
// the interpretation against its source.
type Request struct {
	Source     string `yaml:"source" json:"source"`
	CapturedAt string `yaml:"captured_at" json:"captured_at"`
	Asks       []Ask  `yaml:"asks" json:"asks"`
}

// dispositions projects the part of a Request that the activation fingerprint
// freezes: the decisions, without the text that states them.
//
// The split is deliberate. Sharpening an ask mid-Mission must not invalidate
// activation, because punishing clarification is exactly the wrong incentive.
// Relabelling an ask from covered to deferred after activation is precisely the
// drift this record exists to catch, so it costs a re-activation.
func (r *Request) dispositions() []Ask {
	if r == nil {
		return nil
	}
	frozen := make([]Ask, 0, len(r.Asks))
	for _, ask := range r.Asks {
		claims := append([]string(nil), ask.Claims...)
		sort.Strings(claims)
		// Ask text and Reason are prose and stay outside the boundary.
		frozen = append(frozen, Ask{Disposition: ask.Disposition, Claims: claims})
	}
	return frozen
}

type TransitionHistory struct {
	At         string `yaml:"at" json:"at"`
	From       string `yaml:"from" json:"from"`
	To         string `yaml:"to" json:"to"`
	By         string `yaml:"by" json:"by"`
	Reason     string `yaml:"reason" json:"reason"`
	NextAction string `yaml:"next_action,omitempty" json:"next_action,omitempty"`
}

type Objective struct {
	Ref            string   `yaml:"ref" json:"ref"`
	ID             string   `yaml:"id" json:"id"`
	Source         string   `yaml:"-" json:"source,omitempty"`
	Sources        []string `yaml:"sources,omitempty" json:"sources,omitempty"`
	Outcome        string   `yaml:"outcome,omitempty" json:"outcome,omitempty"`
	Status         string   `yaml:"status,omitempty" json:"status,omitempty"`
	After          []string `yaml:"after,omitempty" json:"after,omitempty"`
	AfterInterface []string `yaml:"after_interface,omitempty" json:"after_interface,omitempty"`
	Claims         []string `yaml:"claims,omitempty" json:"claims,omitempty"`
	File           string   `yaml:"file,omitempty" json:"file,omitempty"`
	Body           string   `yaml:"-" json:"-"`
	Run            *Run     `yaml:"run,omitempty" json:"run,omitempty"`
	Runs           []Run    `yaml:"runs,omitempty" json:"runs,omitempty"`

	document *workspace.Document
}

type Run struct {
	Ref              string              `yaml:"ref" json:"ref"`
	ID               string              `yaml:"id" json:"id"`
	Source           string              `yaml:"-" json:"source,omitempty"`
	Title            string              `yaml:"title,omitempty" json:"title,omitempty"`
	Status           string              `yaml:"status" json:"status"`
	Operator         string              `yaml:"operator" json:"operator"`
	StartedAt        string              `yaml:"started_at" json:"started_at"`
	CurrentObjective string              `yaml:"current_objective,omitempty" json:"current_objective,omitempty"`
	Objective        string              `yaml:"objective,omitempty" json:"objective,omitempty"`
	Repairs          int                 `yaml:"repairs" json:"repairs"`
	File             string              `yaml:"file,omitempty" json:"file,omitempty"`
	History          []TransitionHistory `yaml:"history,omitempty" json:"history,omitempty"`
	Body             string              `yaml:"-" json:"-"`

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
	Request          *Request          `json:"request,omitempty"`
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
	ResolvesGaps     []ResolvedGap     `json:"resolves_gaps,omitempty"`
	Stops            []string          `json:"stops,omitempty"`
	Fallbacks        []Fallback        `json:"fallbacks,omitempty"`
	AfterMission     []string          `json:"after_mission,omitempty"`
	Reviews          []ReviewPointer   `json:"reviews,omitempty"`
	Handoffs         []HandoffPointer  `json:"handoffs,omitempty"`
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

	// contractDrift records that the bound Contract's content changed after this
	// Mission completed. A live Mission refuses on that change; a completed one
	// reports it, because its binding is a record of what the work was executed
	// against and re-hashing it asks a question with no remaining consequence.
	// Set by validateContract, read by Notices.
	contractDrift string

	// contractVersion is the bound Contract's declared version, reported rather
	// than enforced. Set by validateContractVersion, read by Notices.
	contractVersion int
}

type Check struct {
	Ref string `json:"ref"`
	// ContractVersion is the bound Contract's declared version. Reported, never
	// enforced: a Mission bound to an earlier version simply ran against it.
	ContractVersion int      `json:"contract_version,omitempty"`
	Path            string   `json:"path"`
	Schema          string   `json:"schema"`
	Valid           bool     `json:"valid"`
	Legacy          bool     `json:"legacy"`
	Checks          []string `json:"checks"`
	Fingerprint     string   `json:"fingerprint,omitempty"`
	// Authority is the decision table the authority-vocabulary validator
	// already resolves and previously discarded. It is derived on read.
	Authority []AuthorityAnswer `json:"authority,omitempty"`
	// Drift is the per-claim audit signal, most-flagged first. Derived on read.
	Drift []ClaimDrift `json:"drift,omitempty"`
	// Notices are reported, non-failing observations such as a record still
	// using the legacy `human_ref:` spelling. A notice never makes Valid false:
	// frozen records are not rewritten to finish a rename.
	Notices []string `json:"notices,omitempty"`
}

type Result struct {
	Operation   string   `json:"operation"`
	Ref         string   `json:"ref"`
	Path        string   `json:"path"`
	Fingerprint string   `json:"fingerprint,omitempty"`
	Changed     []string `json:"changed,omitempty"`
}
