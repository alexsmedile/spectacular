// Package missionview exposes the data needed by read-only Mission consumers.
// It has no persistence, serialization, or service dependencies.
package missionview

type Criterion struct{ Claim, PassBoundary, ProofRequirement string }
type Objective struct {
	Ref, ID, Outcome string
	Claims, Sources  []string
}
type Gap struct{ Gap, Resolution string }
type Handoff struct {
	Task, Title string
	Writes      []string
}
type Mission struct {
	Ref, BaselineCommit, ContractRef, VerificationCommand string
	Completion                                            []Criterion
	Objectives                                            []Objective
	ResolvesGaps                                          []Gap
	Handoffs                                              []Handoff
	WritesPaths, AllowedActions, RequiresOwner, Stops     []string
}
