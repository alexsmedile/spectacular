package missionbundle

import (
	"path/filepath"
	"strings"
	"testing"

	"github.com/alexsmedile/spectacular/v2/internal/discovery"
	"github.com/alexsmedile/spectacular/v2/internal/domain"
	"github.com/alexsmedile/spectacular/v2/internal/workspace"
)

// Every name a Contract declares in mandatory_validation must correspond to
// something that runs. Before this check, CC-projsurf declared seven names and four
// resolved to nothing: one was registered under a different name, one was
// implemented inside another validator, one was a notice filed as a validator, and
// one was a validator built with no caller. A Contract could promise a validation
// that never executed and every Mission bound to it still reported valid.
func TestEveryDeclaredValidationResolves(t *testing.T) {
	root, err := filepath.Abs(filepath.Join("..", ".."))
	if err != nil {
		t.Fatal(err)
	}
	ws, err := discovery.Open(root)
	if err != nil {
		t.Fatal(err)
	}
	contracts := 0
	for _, entry := range ws.Entries {
		if entry.Document == nil || entry.Document.Record.Type != domain.Contract {
			continue
		}
		declared, err := workspace.Strings(entry.Document, "mandatory_validation", false)
		if err != nil || len(declared) == 0 {
			continue
		}
		contracts++
		for _, name := range declared {
			kind, ok := ResolveDeclaredValidation(name)
			if !ok {
				t.Errorf("%s declares %q, which resolves to no validator, notice, or Proposal check",
					filepath.Base(entry.Path), name)
				continue
			}
			t.Logf("%s: %s -> %s", filepath.Base(entry.Path), name, kind)
		}
	}
	if contracts == 0 {
		t.Fatal("no Contract in the workspace declares mandatory_validation")
	}
}

// A name matching nothing must not resolve, or the check above proves nothing.
func TestUnknownDeclaredValidationDoesNotResolve(t *testing.T) {
	for _, name := range []string{"", "   ", "invented-validator", "yaml-schema-typo"} {
		if kind, ok := ResolveDeclaredValidation(name); ok {
			t.Fatalf("%q resolved as %s", name, kind)
		}
	}
	// One of each kind must resolve, and to the right kind.
	for name, want := range map[string]string{
		"yaml-schema":        "validator",
		"ref-spelling-drift": "notice",
		"proposal-schema-v2": "proposal-validator",
	} {
		kind, ok := ResolveDeclaredValidation(name)
		if !ok || kind != want {
			t.Fatalf("%q resolved as (%q, %t), want %q", name, kind, ok, want)
		}
	}
}

// The two renamed validators keep their refusal identity. A rename that changed the
// code, field, or message would fix a declaration by breaking the contract the
// refusal itself makes with a reader.
func TestRenamedValidatorsKeepTheirRefusalIdentity(t *testing.T) {
	root, err := filepath.Abs(filepath.Join("..", ".."))
	if err != nil {
		t.Fatal(err)
	}
	ws, err := discovery.Open(root)
	if err != nil {
		t.Fatal(err)
	}
	source, err := Load(ws, "M9")
	if err != nil {
		t.Fatal(err)
	}

	tests := []struct {
		name   string
		mutate func(*Bundle)
		run    func(*discovery.Workspace, *Bundle) error
		field  string
	}{
		{"fallback-fingerprint-coverage", func(b *Bundle) {
			b.Fallbacks = []Fallback{{Approach: "", RejectedBecause: "r", InvalidatedIf: "i"}}
		}, validateFallbacks, "fallbacks.approach"},
		{"interface-dependency-frozen-target", func(b *Bundle) {
			b.Objectives = []Objective{{Ref: "O1", Outcome: "o", Status: "pending",
				Claims: []string{"c"}, AfterInterface: []string{"O9"}}}
		}, validateInterfaceDependencies, "objectives.after_interface"},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if _, ok := ResolveDeclaredValidation(test.name); !ok {
				t.Fatalf("%s is not registered", test.name)
			}
			bundle := cloneBundle(t, source)
			test.mutate(bundle)
			err := test.run(ws, bundle)
			if err == nil {
				t.Fatal("the mutation must refuse")
			}
			refusal, ok := err.(*domain.Refusal)
			if !ok {
				t.Fatalf("returned %T, want a typed refusal", err)
			}
			if refusal.Field != test.field {
				t.Fatalf("refused on %q, want %q", refusal.Field, test.field)
			}
			if !strings.Contains(strings.ToLower(string(refusal.Code)), "invalid") {
				t.Fatalf("refusal code is %q", refusal.Code)
			}
		})
	}
}
