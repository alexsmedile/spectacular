package charter

import (
	"errors"
	"fmt"
	"path/filepath"
	"strings"

	"github.com/alexsmedile/spectacular/v2/internal/charter/tokenizer"
	"github.com/alexsmedile/spectacular/v2/internal/discovery"
	"github.com/alexsmedile/spectacular/v2/internal/domain"
	"github.com/alexsmedile/spectacular/v2/internal/missionbundle"
	"github.com/alexsmedile/spectacular/v2/internal/workspace"
)

var (
	ErrObjectiveNotFound = errors.New("charter compiler: objective not found in mission")
	ErrRefusalOverCap    = errors.New("charter compiler: token count exceeds hard 1440 ceiling")
)

// Compile builds a read-only 3-layer Context Sandwich charter for the given Mission and Objective.
func Compile(ws *discovery.Workspace, missionRef string, objectiveRef string, extraSources []string) (*Charter, error) {
	if ws == nil {
		return nil, errors.New("charter compiler: workspace is nil")
	}

	bundle, err := missionbundle.Load(ws, missionRef)
	if err != nil {
		return nil, fmt.Errorf("charter compiler: load mission %s: %w", missionRef, err)
	}

	missionEntry, err := ws.Lookup(missionRef, domain.Mission)
	if err != nil {
		return nil, fmt.Errorf("charter compiler: lookup mission entry %s: %w", missionRef, err)
	}

	// Normalize objective ref (e.g., "M16/O1" -> "O1")
	objKey := objectiveRef
	if strings.Contains(objKey, "/") {
		parts := strings.Split(objKey, "/")
		objKey = parts[len(parts)-1]
	}

	var targetObj *missionbundle.Objective
	for i := range bundle.Objectives {
		if bundle.Objectives[i].Ref == objKey || bundle.Objectives[i].ID == objKey {
			targetObj = &bundle.Objectives[i]
			break
		}
	}
	if targetObj == nil {
		return nil, fmt.Errorf("%w: %s in mission %s", ErrObjectiveNotFound, objectiveRef, missionRef)
	}

	// 1. Build Layer 1: Frozen Truth
	baselineCommit := ""
	if bundle.Baseline != nil {
		baselineCommit = bundle.Baseline.Commit
	}

	layer1 := Layer1{
		ProjectAnchor: ws.Manifest.ProjectAnchor,
		MissionRef:    bundle.Ref,
		ObjectiveRef:  targetObj.Ref,
		Outcome:       targetObj.Outcome,
		ContractRef:   bundle.Contract.Ref,
		GitBaseline:   baselineCommit,
	}

	// Map completion claims for this objective
	claimSet := make(map[string]bool)
	for _, c := range targetObj.Claims {
		claimSet[c] = true
	}
	for _, comp := range bundle.Completion {
		if len(claimSet) == 0 || claimSet[comp.Claim] {
			layer1.Claims = append(layer1.Claims, ClaimItem{
				Claim:            comp.Claim,
				PassBoundary:     comp.PassBoundary,
				ProofRequirement: comp.ProofRequirement,
			})
		}
	}

	// 2. Resolve Sources in Strict Declaration Order: Bound Contract -> Mission sources -> Objective sources -> extraSources
	// Check Mission frontmatter for optional `sources:`.
	var missionSources []string
	if missionDoc := missionEntry.Document; missionDoc != nil {
		if sources, err := workspace.Strings(missionDoc, "sources", false); err == nil {
			missionSources = sources
		}
	}
	orderedSources := declaredSourceRefs(bundle.Contract.Ref, missionSources, targetObj.Sources, extraSources)

	// 3. Build Layer 2: Owner Steering & Layer 3: Perimeter
	var boundSources []BoundSource
	var decisions []DecisionItem
	var gaps []GapItem

	for _, sRef := range orderedSources {
		entry, fingerprint, found := resolveSource(ws, sRef)
		if !found || entry.Document == nil {
			return nil, domain.NewRefusal(domain.RefusalRecordNotFound, sRef, fmt.Sprintf("declared charter source %q could not be resolved", sRef), nil)
		}

		doc := entry.Document
		recordType := string(doc.Record.Type)
		title := ""
		if doc.Record.Title != nil {
			title = *doc.Record.Title
		}

		boundSources = append(boundSources, BoundSource{
			Ref:         sRef,
			Type:        recordType,
			Title:       title,
			Fingerprint: fingerprint,
		})

		if recordType == "Decision" || strings.HasPrefix(sRef, "D") || strings.Contains(entry.Path, "decisions") {
			disposition, _ := workspace.String(doc, "disposition", false)
			rationale, _ := workspace.String(doc, "rationale", false)
			decisions = append(decisions, DecisionItem{
				Ref:         sRef,
				Title:       title,
				Disposition: disposition,
				Rationale:   rationale,
			})
		}
	}

	// Add resolved gaps from mission if any
	for _, rg := range bundle.ResolvesGaps {
		gaps = append(gaps, GapItem{
			Ref:        rg.Gap,
			Resolution: rg.Resolution,
		})
	}

	layer2 := Layer2{
		Decisions: decisions,
		Gaps:      gaps,
	}

	// Build Layer 3
	layer3 := Layer3{
		WritesPaths:    bundle.Scope.Mechanical,
		AllowedActions: bundle.Authority.Operator,
		RequiresOwner:  bundle.Authority.RequiresOwner,
		Stops:          bundle.Stops,
	}

	c := &Charter{
		SchemaVersion: SchemaVersion,
		MissionRef:    bundle.Ref,
		ObjectiveRef:  targetObj.Ref,
		Sources:       boundSources,
		Layer1:        layer1,
		Layer2:        layer2,
		Layer3:        layer3,
	}

	// Render & Count Tokens
	rendered := c.RenderMarkdown()
	tokenCount, err := tokenizer.Count(rendered)
	if err != nil {
		return nil, fmt.Errorf("charter compiler: count tokens: %w", err)
	}

	targetCap := ws.Config.TokenBudgets.Charter.Target
	if targetCap == 0 {
		targetCap = tokenizer.MaxTargetTokens
	}
	hardCeiling := ws.Config.TokenBudgets.Charter.HardCap
	if hardCeiling == 0 {
		hardCeiling = tokenizer.HardCeilingTokens
	}

	// Safe Compaction if over target budget
	if tokenCount > targetCap {
		c.applySafeCompaction()
		rendered = c.RenderMarkdown()
		tokenCount, _ = tokenizer.Count(rendered)
	}

	c.TokenCount = tokenCount
	disp, _ := tokenizer.EvaluateDisposition(tokenCount)
	c.Disposition = disp

	if tokenCount > hardCeiling {
		return nil, domain.NewRefusal(domain.RefusalInvalidScope, "tokens", fmt.Sprintf("charter token count %d exceeds hard refusal ceiling %d even after compaction; objective split required", tokenCount, hardCeiling), nil)
	}

	return c, nil
}

// resolveSource locates a source entry by exact ref, typed ref, path, or basename match in .spectacular/
func resolveSource(ws *discovery.Workspace, ref string) (discovery.Entry, string, bool) {
	// Try direct source resolution first
	if entry, fp, err := ws.Source(ref); err == nil {
		return entry, fp, true
	}

	// Try standard Lookup across kinds
	for _, kind := range []domain.RecordType{domain.Decision, domain.Contract, domain.Proposal, domain.Mission} {
		if entry, err := ws.Lookup(ref, kind); err == nil {
			return entry, entry.Fingerprint, true
		}
	}

	// Try entry match by path, basename, or ref prefix
	normalized := strings.TrimPrefix(ref, ".spectacular/")
	for _, entry := range ws.Entries {
		base := strings.TrimSuffix(filepath.Base(entry.Path), ".md")
		if base == ref || base == normalized || entry.Path == ref || entry.Path == ".spectacular/"+ref {
			return entry, entry.Fingerprint, true
		}
		// Match D12-isolation... from "D12" or "D12-isolation-and-context-compilation"
		if strings.HasPrefix(base, ref+"-") || strings.HasPrefix(ref, base+"-") {
			return entry, entry.Fingerprint, true
		}
		// Check ref inside document
		if docRef, _, _ := workspace.Ref(entry.Document); docRef != "" && docRef == ref {
			return entry, entry.Fingerprint, true
		}
	}

	return discovery.Entry{}, "", false
}

// applySafeCompaction shortens rationale and descriptive fields while strictly preserving claims, authority, and file perimeters.
func (c *Charter) applySafeCompaction() {
	c.Compacted = true
	for i := range c.Layer2.Decisions {
		r := c.Layer2.Decisions[i].Rationale
		if len(r) > 150 {
			// Compact rationale to first sentence
			if idx := strings.Index(r, ". "); idx > 0 && idx < 150 {
				c.Layer2.Decisions[i].Rationale = r[:idx+1]
			} else {
				c.Layer2.Decisions[i].Rationale = r[:147] + "..."
			}
		}
	}
}
