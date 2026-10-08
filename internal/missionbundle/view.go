package missionbundle

import (
	"github.com/alexsmedile/spectacular/v2/internal/discovery"
	"github.com/alexsmedile/spectacular/v2/internal/missionview"
)

// ReadView uses the canonical decoder and returns detached data for readers.
// No mutable Document, Entry, or service escapes this boundary.
func ReadView(ws *discovery.Workspace, ref string) (*missionview.Mission, error) {
	b, err := Load(ws, ref)
	if err != nil {
		return nil, err
	}
	return projectView(b), nil
}

func projectView(b *Bundle) *missionview.Mission {
	v := &missionview.Mission{
		Ref: b.Ref, ContractRef: b.Contract.Ref,
		WritesPaths:    append([]string(nil), b.Scope.Mechanical...),
		AllowedActions: append([]string(nil), b.Authority.Operator...),
		RequiresOwner:  append([]string(nil), b.Authority.RequiresOwner...),
		Stops:          append([]string(nil), b.Stops...),
	}
	if b.Baseline != nil {
		v.BaselineCommit = b.Baseline.Commit
	}
	if b.Replay != nil {
		v.VerificationCommand = b.Replay.Command
	}
	for _, c := range b.Completion {
		v.Completion = append(v.Completion, missionview.Criterion{Claim: c.Claim, PassBoundary: c.PassBoundary, ProofRequirement: c.ProofRequirement})
	}
	for _, o := range b.Objectives {
		v.Objectives = append(v.Objectives, missionview.Objective{Ref: o.Ref, ID: o.ID, Outcome: o.Outcome, Claims: append([]string(nil), o.Claims...), Sources: append([]string(nil), o.Sources...)})
	}
	for _, g := range b.ResolvesGaps {
		v.ResolvesGaps = append(v.ResolvesGaps, missionview.Gap{Gap: g.Gap, Resolution: g.Resolution})
	}
	for _, h := range b.Handoffs {
		if h.Document != nil {
			v.Handoffs = append(v.Handoffs, missionview.Handoff{Task: h.Document.Task, Title: h.Document.Title, Writes: append([]string(nil), h.Document.Writes...)})
		}
	}
	return v
}
