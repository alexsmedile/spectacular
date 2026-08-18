package missionbundle

import (
	"errors"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/alexsmedile/spectacular/v2/internal/discovery"
	"github.com/alexsmedile/spectacular/v2/internal/domain"
)

// handoffFixture builds a complete, valid Handoff file body. Each test mutates
// one thing, so a refusal is attributable to the field the case names rather
// than to some other difference between two hand-written fixtures.
type handoffFixture struct {
	ref        string
	id         string
	mission    string
	title      string
	commit     string
	tree       string
	sender     string
	relation   string
	task       string
	asserted   *[]string
	assumed    *[]string
	stops      *[]string
	returns    *[]string
	supersedes string
}

func list(values ...string) *[]string { return &values }

func (f handoffFixture) write(t *testing.T, directory string) string {
	t.Helper()
	var b strings.Builder
	b.WriteString("---\ntype: Handoff\n")
	b.WriteString("id: " + f.id + "\n")
	b.WriteString("title: " + f.title + "\n")
	b.WriteString("human_ref: " + f.mission + "/" + f.ref + "\n")
	b.WriteString("mission: " + f.mission + "\n")
	b.WriteString("reviewed:\n    commit: " + f.commit + "\n    tree: " + f.tree + "\n")
	if f.sender != "" || f.relation != "" {
		b.WriteString("sender:\n    actor: " + f.sender + "\n    relation_to_receiver: " + f.relation + "\n")
	}
	if f.task != "" {
		b.WriteString("task: " + f.task + "\n")
	}
	if f.supersedes != "" {
		b.WriteString("supersedes: " + f.supersedes + "\n")
	}
	for _, field := range []struct {
		name   string
		values *[]string
	}{
		{"asserted", f.asserted}, {"assumed", f.assumed},
		{"stops", f.stops}, {"returns", f.returns},
	} {
		if field.values == nil {
			continue
		}
		if len(*field.values) == 0 {
			b.WriteString(field.name + ": []\n")
			continue
		}
		b.WriteString(field.name + ":\n")
		for _, value := range *field.values {
			b.WriteString("    - " + value + "\n")
		}
	}
	b.WriteString("---\n\nHandoff body.\n")

	path := filepath.Join(directory, f.ref+"-handoff.md")
	if err := os.WriteFile(path, []byte(b.String()), 0o644); err != nil {
		t.Fatal(err)
	}
	return path
}

func completeHandoff() handoffFixture {
	return handoffFixture{
		ref: "H1-aaaaaa", id: "01a0103d-7fc8-74d4-9a8c-f3dae87cc111", mission: "M12",
		title:  "Delegate the amendment repair",
		commit: strings.Repeat("a", 40), tree: strings.Repeat("b", 40),
		sender: "Alex", relation: "operator", task: "Repair the Gap rewrite",
		asserted: list("the adversarial fixture reproduces the defect"),
		assumed:  list("the amendment path is otherwise unchanged"),
		stops:    list("the rewrite would reflow an unrelated scalar"),
		returns:  list("the passing fixture and the diff"),
	}
}

// A Handoff is checked against the Mission that carries it, so the fixtures need
// a Mission bundle to point at them. This builds one in a temporary workspace.
func handoffBundle(t *testing.T, fixtures ...handoffFixture) (*discovery.Workspace, *Bundle) {
	t.Helper()
	root := t.TempDir()
	directory := filepath.Join(root, ".spectacular", "missions", "M12-handoffs", "handoffs")
	if err := os.MkdirAll(directory, 0o755); err != nil {
		t.Fatal(err)
	}
	bundle := &Bundle{
		Ref: "M12", ID: "01a010a6-01b0-7320-acc2-5c695bec2843",
		entry: discovery.Entry{Absolute: filepath.Join(root, ".spectacular", "missions", "M12-handoffs", "M12-handoffs.md")},
	}
	for _, fixture := range fixtures {
		fixture.write(t, directory)
		bundle.Handoffs = append(bundle.Handoffs, HandoffPointer{
			Ref: fixture.ref, ID: fixture.id,
			File: filepath.ToSlash(filepath.Join("handoffs", fixture.ref+"-handoff.md")),
		})
	}
	return &discovery.Workspace{Root: root}, bundle
}

func TestHandoffSchemaRefusesEachMissingRequirement(t *testing.T) {
	tests := []struct {
		name    string
		mutate  func(*handoffFixture)
		code    domain.RefusalCode
		field   string
		wantErr bool
	}{
		{name: "complete handoff is valid"},
		{
			name:   "absent asserted refuses",
			mutate: func(f *handoffFixture) { f.asserted = nil },
			code:   domain.RefusalMissingRequiredField, field: "asserted", wantErr: true,
		},
		{
			// An empty assumed list is a sender stating they assumed nothing. That
			// is a real thing to say and must not be confused with saying nothing.
			name:   "empty assumed is legal",
			mutate: func(f *handoffFixture) { f.assumed = list() },
		},
		{
			name:   "empty asserted is legal",
			mutate: func(f *handoffFixture) { f.asserted = list() },
		},
		{
			name:   "absent assumed refuses",
			mutate: func(f *handoffFixture) { f.assumed = nil },
			code:   domain.RefusalMissingRequiredField, field: "assumed", wantErr: true,
		},
		{
			name:   "absent stops refuses",
			mutate: func(f *handoffFixture) { f.stops = nil },
			code:   domain.RefusalMissingRequiredField, field: "stops", wantErr: true,
		},
		{
			name:   "empty stops refuses",
			mutate: func(f *handoffFixture) { f.stops = list() },
			code:   domain.RefusalInvalidKnownField, field: "handoff.stops", wantErr: true,
		},
		{
			name:   "absent returns refuses",
			mutate: func(f *handoffFixture) { f.returns = nil },
			code:   domain.RefusalMissingRequiredField, field: "returns", wantErr: true,
		},
		{
			name:   "absent task refuses",
			mutate: func(f *handoffFixture) { f.task = "" },
			code:   domain.RefusalInvalidKnownField, field: "handoff.task", wantErr: true,
		},
		{
			name:   "absent sender refuses",
			mutate: func(f *handoffFixture) { f.sender, f.relation = "", "" },
			code:   domain.RefusalMissingRequiredField, field: "sender", wantErr: true,
		},
		{
			name:   "sender without relation refuses",
			mutate: func(f *handoffFixture) { f.relation = "" },
			code:   domain.RefusalInvalidKnownField, field: "handoff.sender", wantErr: true,
		},
		{
			name:   "a commit that is not a full hash refuses",
			mutate: func(f *handoffFixture) { f.commit = "abc123" },
			code:   domain.RefusalInvalidKnownField, field: "handoff.reviewed", wantErr: true,
		},
		{
			name:   "a handoff naming a different Mission refuses",
			mutate: func(f *handoffFixture) { f.mission = "M11" },
			code:   domain.RefusalInvalidKnownField, field: "handoff.mission", wantErr: true,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			fixture := completeHandoff()
			if test.mutate != nil {
				test.mutate(&fixture)
			}
			ws, bundle := handoffBundle(t, fixture)
			err := validateHandoffs(ws, bundle)
			if !test.wantErr {
				if err != nil {
					t.Fatalf("expected a valid Handoff, got refusal: %v", err)
				}
				return
			}
			if err == nil {
				t.Fatal("expected a refusal, got none")
			}
			requireRefusal(t, err, test.code, test.field)
		})
	}
}

func TestHandoffSupersedesResolvesAgainstTheSameMission(t *testing.T) {
	original := completeHandoff()
	correction := completeHandoff()
	correction.ref, correction.id = "H2-bbbbbb", "01a0103d-7fc8-74d4-9a8c-f3dae87cc222"

	t.Run("superseding a recorded Handoff is valid", func(t *testing.T) {
		corrected := correction
		corrected.supersedes = original.ref
		ws, bundle := handoffBundle(t, original, corrected)
		if err := validateHandoffs(ws, bundle); err != nil {
			t.Fatalf("expected a valid supersession, got: %v", err)
		}
	})

	t.Run("superseding nothing refuses", func(t *testing.T) {
		corrected := correction
		corrected.supersedes = "H9-zzzzzz"
		ws, bundle := handoffBundle(t, original, corrected)
		requireRefusal(t, validateHandoffs(ws, bundle), domain.RefusalInvalidKnownField, "handoff.supersedes")
	})

	t.Run("superseding a Handoff on a different Mission refuses", func(t *testing.T) {
		// The superseded Handoff is not recorded on this Mission, which is how a
		// cross-Mission supersedes presents: the ref resolves nowhere here.
		corrected := correction
		corrected.supersedes = "H1-cccccc"
		ws, bundle := handoffBundle(t, original, corrected)
		requireRefusal(t, validateHandoffs(ws, bundle), domain.RefusalInvalidKnownField, "handoff.supersedes")
	})

	t.Run("superseding itself refuses", func(t *testing.T) {
		corrected := correction
		corrected.supersedes = corrected.ref
		ws, bundle := handoffBundle(t, original, corrected)
		requireRefusal(t, validateHandoffs(ws, bundle), domain.RefusalInvalidKnownField, "handoff.supersedes")
	})
}

// The asserted/assumed split is a claim the sender signs, not a fact the system
// can check. This asserts no validator reads either list's contents: a mechanism
// that scored them would be asserting something it cannot know, and the decision
// to report and never score it is one this Mission's stops forbid reversing.
func TestNoValidatorScoresAssertedOrAssumedContent(t *testing.T) {
	sentinel := "this text is never inspected by any validator"
	fixture := completeHandoff()
	fixture.asserted = list(sentinel)
	fixture.assumed = list(sentinel)
	ws, bundle := handoffBundle(t, fixture)
	if err := validateHandoffs(ws, bundle); err != nil {
		t.Fatalf("content of asserted/assumed changed the verdict: %v", err)
	}

	// The same Handoff with contradictory content is equally valid. If any
	// validator judged the claims, these two could not both pass.
	contradictory := completeHandoff()
	contradictory.asserted = list("the tests pass")
	contradictory.assumed = list("the tests do not pass")
	ws, bundle = handoffBundle(t, contradictory)
	if err := validateHandoffs(ws, bundle); err != nil {
		t.Fatalf("contradictory asserted/assumed content was scored: %v", err)
	}
}

func requireRefusal(t *testing.T, err error, code domain.RefusalCode, field string) {
	t.Helper()
	if err == nil {
		t.Fatal("expected a refusal, got none")
	}
	var refusal *domain.Refusal
	if !errors.As(err, &refusal) {
		t.Fatalf("expected a typed refusal, got %T: %v", err, err)
	}
	if refusal.Code != code || refusal.Field != field {
		t.Fatalf("refusal code=%s field=%s, want code=%s field=%s", refusal.Code, refusal.Field, code, field)
	}
}
