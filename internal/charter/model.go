// Package charter owns the read-only 3-layer Context Sandwich charter compiler.
package charter

import (
	"fmt"
	"strings"

	"github.com/alexsmedile/spectacular/v2/internal/charter/tokenizer"
)

const SchemaVersion = "spectacular.charter.v1"

// BoundSource tracks an explicitly resolved governance source in declaration order.
type BoundSource struct {
	Ref         string `json:"ref"`
	Type        string `json:"type"`
	Title       string `json:"title"`
	Fingerprint string `json:"fingerprint,omitempty"`
}

// Layer1 represents Frozen Truth: Project boundaries, active Mission/Objective, and completion criteria.
type Layer1 struct {
	ProjectAnchor string      `json:"project_anchor"`
	MissionRef    string      `json:"mission_ref"`
	ObjectiveRef  string      `json:"objective_ref"`
	Outcome       string      `json:"outcome"`
	Claims        []ClaimItem `json:"claims"`
	ContractRef   string      `json:"contract_ref"`
	GitBaseline   string      `json:"git_baseline"`
}

// ClaimItem represents a single atomic completion claim.
type ClaimItem struct {
	Claim            string `json:"claim"`
	PassBoundary     string `json:"pass_boundary"`
	ProofRequirement string `json:"proof_requirement"`
}

// Layer2 represents Owner Steering: Decision summaries, resolved Gaps, and non-goals.
type Layer2 struct {
	Decisions []DecisionItem `json:"decisions"`
	Gaps      []GapItem      `json:"gaps,omitempty"`
	NonGoals  []string       `json:"non_goals,omitempty"`
}

// DecisionItem captures an explicitly retrieved Decision.
type DecisionItem struct {
	Ref         string `json:"ref"`
	Title       string `json:"title"`
	Disposition string `json:"disposition"`
	Rationale   string `json:"rationale"`
}

// GapItem captures a resolved Gap.
type GapItem struct {
	Ref        string `json:"ref"`
	Problem    string `json:"problem"`
	Resolution string `json:"resolution"`
}

// Layer3 represents the Execution Perimeter: Target files, permissions, stops, and check commands.
type Layer3 struct {
	WritesPaths         []string `json:"writes_paths"`
	ReadsPaths          []string `json:"reads_paths"`
	AllowedActions      []string `json:"allowed_actions"`
	RequiresOwner       []string `json:"requires_owner"`
	Stops               []string `json:"stops"`
	VerificationCommand string   `json:"verification_command"`
}

// Charter represents the compiled, 3-layer Context Sandwich envelope.
type Charter struct {
	SchemaVersion string                `json:"schema_version"`
	MissionRef    string                `json:"mission_ref"`
	ObjectiveRef  string                `json:"objective_ref"`
	Sources       []BoundSource         `json:"sources"`
	Layer1        Layer1                `json:"layer1_frozen_truth"`
	Layer2        Layer2                `json:"layer2_owner_steering"`
	Layer3        Layer3                `json:"layer3_execution_perimeter"`
	TokenCount    int                   `json:"token_count"`
	Disposition   tokenizer.Disposition `json:"disposition"`
	Compacted     bool                  `json:"compacted"`
}

// RenderMarkdown formats the 3-layer charter into a clean, token-efficient prompt envelope.
func (c *Charter) RenderMarkdown() string {
	var b strings.Builder

	b.WriteString("# Context Charter: " + c.MissionRef + "/" + c.ObjectiveRef + "\n\n")

	// Layer 1: Frozen Truth
	b.WriteString("## 1. FROZEN TRUTH\n")
	b.WriteString("- **Mission**: " + c.MissionRef + "\n")
	b.WriteString("- **Objective**: " + c.ObjectiveRef + " — " + c.Layer1.Outcome + "\n")
	b.WriteString("- **Contract**: " + c.Layer1.ContractRef + "\n")
	b.WriteString("- **Baseline**: " + c.Layer1.GitBaseline + "\n\n")

	b.WriteString("### Completion Claims & Criteria:\n")
	for _, cl := range c.Layer1.Claims {
		b.WriteString(fmt.Sprintf("- **Claim**: `%s`\n  - **Pass Boundary**: %s\n  - **Proof**: %s\n", cl.Claim, cl.PassBoundary, cl.ProofRequirement))
	}
	b.WriteString("\n")

	// Layer 2: Owner Steering
	b.WriteString("## 2. OWNER STEERING\n")
	if len(c.Layer2.Decisions) > 0 {
		b.WriteString("### Relevant Decisions:\n")
		for _, d := range c.Layer2.Decisions {
			b.WriteString(fmt.Sprintf("- **%s** (`%s`): %s\n", d.Ref, d.Disposition, d.Title))
			if d.Rationale != "" {
				b.WriteString(fmt.Sprintf("  *Rationale*: %s\n", d.Rationale))
			}
		}
	}
	if len(c.Layer2.Gaps) > 0 {
		b.WriteString("### Resolved Gaps:\n")
		for _, g := range c.Layer2.Gaps {
			b.WriteString(fmt.Sprintf("- **%s**: %s (Resolution: %s)\n", g.Ref, g.Problem, g.Resolution))
		}
	}
	if len(c.Layer2.NonGoals) > 0 {
		b.WriteString("### Non-Goals:\n")
		for _, ng := range c.Layer2.NonGoals {
			b.WriteString("- " + ng + "\n")
		}
	}
	b.WriteString("\n")

	// Layer 3: Execution Perimeter
	b.WriteString("## 3. EXECUTION PERIMETER\n")
	b.WriteString("### Authorized Writes:\n")
	if len(c.Layer3.WritesPaths) == 0 {
		b.WriteString("- *(none declared)*\n")
	} else {
		for _, w := range c.Layer3.WritesPaths {
			b.WriteString("- `" + w + "`\n")
		}
	}

	if len(c.Layer3.ReadsPaths) > 0 {
		b.WriteString("### Allowed Reads:\n")
		for _, r := range c.Layer3.ReadsPaths {
			b.WriteString("- `" + r + "`\n")
		}
	}

	b.WriteString("### Authority & Stops:\n")
	b.WriteString("- **Operator Actions**: " + strings.Join(c.Layer3.AllowedActions, ", ") + "\n")
	b.WriteString("- **Requires Owner**: " + strings.Join(c.Layer3.RequiresOwner, ", ") + "\n")
	b.WriteString("- **Stops**: " + strings.Join(c.Layer3.Stops, ", ") + "\n")
	if c.Layer3.VerificationCommand != "" {
		b.WriteString("- **Verification Command**: `" + c.Layer3.VerificationCommand + "`\n")
	}

	return b.String()
}
