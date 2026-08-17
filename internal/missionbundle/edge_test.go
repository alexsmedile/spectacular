package missionbundle

import (
	"errors"
	"path/filepath"
	"strings"
	"testing"

	"github.com/alexsmedile/spectacular/v2/internal/discovery"
	"github.com/alexsmedile/spectacular/v2/internal/domain"
)

func TestInterfaceEdgeSplitValidation(t *testing.T) {
	root, err := filepath.Abs(filepath.Join("..", ".."))
	if err != nil {
		t.Fatal(err)
	}
	ws, err := discovery.Open(root)
	if err != nil {
		t.Fatal(err)
	}
	base, err := Load(ws, "M6")
	if err != nil {
		t.Fatal(err)
	}

	t.Run("table driven fixtures validate both edge kinds independently and combined", func(t *testing.T) {
		tests := []struct {
			name       string
			objectives []Objective
			valid      bool
			code       domain.RefusalCode
			field      string
			targetName string
		}{
			{
				name: "independent artifact dependency O2 after O1",
				objectives: []Objective{
					{Ref: "O1", Outcome: "O1", Status: "pending", Claims: []string{base.Completion[0].Claim}},
					{Ref: "O2", Outcome: "O2", Status: "pending", After: []string{"O1"}, Claims: []string{base.Completion[1].Claim}},
				},
				valid: true,
			},
			{
				name: "independent interface dependency O2 after_interface O1",
				objectives: []Objective{
					{Ref: "O1", Outcome: "O1", Status: "pending", Claims: []string{base.Completion[0].Claim}},
					{Ref: "O2", Outcome: "O2", Status: "pending", AfterInterface: []string{"O1"}, Claims: []string{base.Completion[1].Claim}},
				},
				valid: true,
			},
			{
				name: "combined mixed edges O3 after O2, O2 after_interface O1",
				objectives: []Objective{
					{Ref: "O1", Outcome: "O1", Status: "pending", Claims: []string{base.Completion[0].Claim}},
					{Ref: "O2", Outcome: "O2", Status: "pending", AfterInterface: []string{"O1"}, Claims: []string{base.Completion[1].Claim}},
					{Ref: "O3", Outcome: "O3", Status: "pending", After: []string{"O2"}, Claims: []string{base.Completion[2].Claim}},
				},
				valid: true,
			},
			{
				name: "cycle across interface edges only",
				objectives: []Objective{
					{Ref: "O1", Outcome: "O1", Status: "pending", AfterInterface: []string{"O2"}, Claims: []string{base.Completion[0].Claim}},
					{Ref: "O2", Outcome: "O2", Status: "pending", AfterInterface: []string{"O1"}, Claims: []string{base.Completion[1].Claim}},
				},
				valid: false,
				code:  domain.RefusalInvalidKnownField,
				field: "objectives.after",
			},
			{
				name: "cycle across mixed edge kinds O1 after O2 and O2 after_interface O1",
				objectives: []Objective{
					{Ref: "O1", Outcome: "O1", Status: "pending", After: []string{"O2"}, Claims: []string{base.Completion[0].Claim}},
					{Ref: "O2", Outcome: "O2", Status: "pending", AfterInterface: []string{"O1"}, Claims: []string{base.Completion[1].Claim}},
				},
				valid: false,
				code:  domain.RefusalInvalidKnownField,
				field: "objectives.after",
			},
			{
				name: "interface dependency on unfrozen target names the unfrozen target in refusal",
				objectives: []Objective{
					{Ref: "O1", Outcome: "O1", Status: "pending", AfterInterface: []string{"O_unfrozen_target_99"}, Claims: []string{base.Completion[0].Claim}},
				},
				valid:      false,
				code:       domain.RefusalInvalidKnownField,
				field:      "objectives.after_interface",
				targetName: "O_unfrozen_target_99",
			},
			{
				name: "self interface dependency refused",
				objectives: []Objective{
					{Ref: "O1", Outcome: "O1", Status: "pending", AfterInterface: []string{"O1"}, Claims: []string{base.Completion[0].Claim}},
				},
				valid:      false,
				code:       domain.RefusalInvalidKnownField,
				field:      "objectives.after_interface",
				targetName: "O1",
			},
		}

		for _, test := range tests {
			t.Run(test.name, func(t *testing.T) {
				candidate := cloneBundle(t, base)
				candidate.Objectives = test.objectives
				err := validateDAG(ws, candidate)
				if test.valid {
					if err != nil {
						t.Fatalf("unexpected error for valid graph: %v", err)
					}
					return
				}

				if err == nil {
					t.Fatalf("expected refusal for invalid graph, got nil")
				}
				var refusal *domain.Refusal
				if !errors.As(err, &refusal) {
					t.Fatalf("expected *domain.Refusal, got %v", err)
				}
				if refusal.Code != test.code {
					t.Fatalf("refusal.Code=%s, want=%s", refusal.Code, test.code)
				}
				if refusal.Field != test.field {
					t.Fatalf("refusal.Field=%s, want=%s", refusal.Field, test.field)
				}
				if test.targetName != "" && !strings.Contains(refusal.Detail, test.targetName) {
					t.Fatalf("refusal.Detail=%q does not name unfrozen target %q", refusal.Detail, test.targetName)
				}
			})
		}
	})
}
