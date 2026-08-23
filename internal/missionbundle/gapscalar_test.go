package missionbundle

import (
	"strings"
	"testing"
)

// The Gap rewrite is textual, so it has to know which lines are keys and which
// are prose. A `blocked_on:` inside a block scalar body is text the sender wrote.
// M11's independent review predicted this failure; these fixtures are the ones it
// described.
//
// Each case states where the key is and where the decoy is, and asserts both that
// the right key moved and that the scalar body survived byte-for-byte. Asserting
// only the first would pass while the rewrite quietly ate someone's prose.
func TestGapRewriteKnowsItsScalars(t *testing.T) {
	tests := []struct {
		name        string
		contract    string
		wantRefusal bool
		// preserved is prose that must survive the rewrite unchanged.
		preserved []string
	}{
		{
			name: "the key at the entry's top level",
			contract: `---
type: Contract
gaps:
  - ref: g1
    problem: a plain problem
    blocked_on: an owner decision
---
body
`,
			preserved: []string{"problem: a plain problem"},
		},
		{
			name: "a decoy inside a scalar, with the real key after it",
			contract: `---
type: Contract
gaps:
  - ref: g1
    problem: >-
      A Gap whose problem body contains
      blocked_on: as literal text.
    blocked_on: an owner decision
---
body
`,
			preserved: []string{"A Gap whose problem body contains", "blocked_on: as literal text."},
		},
		{
			name: "a decoy inside a literal scalar, with the real key before it",
			contract: `---
type: Contract
gaps:
  - ref: g1
    blocked_on: an owner decision
    problem: |-
      blocked_on: this line is prose
      and so is this one
---
body
`,
			preserved: []string{"blocked_on: this line is prose", "and so is this one"},
		},
		{
			// A Gap recorded with only a problem statement is closable: the
			// resolution is appended to the entry. The decoy text must survive,
			// because a blocked_on inside a scalar body is prose, not a key.
			name: "a decoy and no real key appends the resolution",
			contract: `---
type: Contract
gaps:
  - ref: g1
    problem: >-
      This Gap mentions
      blocked_on: but carries no such key.
---
body
`,
			preserved: []string{"This Gap mentions", "blocked_on: but carries no such key."},
		},
		{
			// blocked_on may itself open a block scalar. Skipping every block-scalar
			// key before testing for blocked_on made the rewrite step over the very
			// key it was looking for.
			name: "the real key is itself a block scalar",
			contract: `---
type: Contract
gaps:
  - ref: g1
    problem: >-
      A Gap mentioning
      blocked_on: in prose.
    blocked_on: >-
      an owner decision
      spanning two lines
---
body
`,
			preserved: []string{"A Gap mentioning", "blocked_on: in prose."},
		},
		{
			name: "a scalar containing a blank line does not end the scalar early",
			contract: `---
type: Contract
gaps:
  - ref: g1
    problem: |-
      blocked_on: prose above a blank line

      blocked_on: prose below it
    blocked_on: an owner decision
---
body
`,
			preserved: []string{"blocked_on: prose above a blank line", "blocked_on: prose below it"},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			amended, err := rewriteGap(test.contract, "g1", "closed for a stated reason")
			if test.wantRefusal {
				if err == nil {
					t.Fatalf("a Gap with no blocked_on key must refuse, got:\n%s", amended)
				}
				return
			}
			if err != nil {
				t.Fatalf("refused: %v", err)
			}
			// The real key is gone and the resolution took its place.
			if strings.Contains(amended, "blocked_on: an owner decision") {
				t.Fatalf("the real blocked_on key survived:\n%s", amended)
			}
			if !strings.Contains(amended, "resolution: >-") || !strings.Contains(amended, "closed for a stated reason") {
				t.Fatalf("the resolution was not written:\n%s", amended)
			}
			// Every scalar body is byte-identical. A rewrite that reflowed prose it
			// did not change is not reviewable, which is why this is textual at all.
			for _, line := range test.preserved {
				if !strings.Contains(amended, line) {
					t.Fatalf("the rewrite damaged prose it did not own\n  lost: %q\n  got:\n%s", line, amended)
				}
			}
			// The amendment touches the Gap it names and nothing else.
			if !strings.HasSuffix(amended, "body\n") {
				t.Fatalf("the Markdown body was modified:\n%s", amended)
			}
			if err := assertOnlyAmendableFieldsChanged(test.contract, amended); err != nil {
				t.Fatalf("the rewrite reached outside the amendable fields: %v", err)
			}
		})
	}
}

// A Contract carrying several Gaps has exactly one of them rewritten, and the
// scalar bodies of the others are untouched even when they contain the decoy.
func TestGapRewriteLeavesOtherGapsAlone(t *testing.T) {
	const contract = `---
type: Contract
gaps:
  - ref: first
    problem: >-
      This Gap is not the target and mentions
      blocked_on: in its prose.
    blocked_on: the first reason
  - ref: second
    problem: also untouched
    blocked_on: the second reason
---
body
`
	amended, err := rewriteGap(contract, "second", "closed for a stated reason")
	if err != nil {
		t.Fatal(err)
	}
	for _, untouched := range []string{
		"This Gap is not the target and mentions",
		"blocked_on: in its prose.",
		"blocked_on: the first reason",
		"problem: also untouched",
	} {
		if !strings.Contains(amended, untouched) {
			t.Fatalf("the rewrite reached another Gap\n  lost: %q\n  got:\n%s", untouched, amended)
		}
	}
	if strings.Contains(amended, "blocked_on: the second reason") {
		t.Fatalf("the targeted Gap was not rewritten:\n%s", amended)
	}
}

// A Gap recorded with only a problem statement is closable. Before this, the
// amendment refused and the only way forward was to hand-edit a bound Contract
// to add a placeholder blocked_on, which is exactly what an amendment exists to
// prevent. The appended key must land inside the Gap it names, aligned with its
// siblings, leaving every neighbouring entry untouched.
func TestGapWithOnlyAProblemIsClosedByAppending(t *testing.T) {
	contract := `---
type: Contract
gaps:
  - ref: g1
    problem: the first problem
    blocked_on: an owner decision
  - ref: g2
    problem: >-
      the second problem, recorded
      with no blocked_on key
  - ref: g3
    problem: the third problem
---
body
`
	amended, err := rewriteGap(contract, "g2", "closed by shipping the fix")
	if err != nil {
		t.Fatalf("a Gap carrying only a problem must be closable: %v", err)
	}
	if !strings.Contains(amended, "    resolution: >-") {
		t.Fatalf("the resolution was not aligned with its sibling keys:\n%s", amended)
	}
	if !strings.Contains(amended, "closed by shipping the fix") {
		t.Fatalf("the resolution text was not written:\n%s", amended)
	}
	// The neighbours keep their own shape: g1's blocked_on is not consumed, and
	// g3 does not acquire the resolution that belongs to g2.
	if !strings.Contains(amended, "blocked_on: an owner decision") {
		t.Fatalf("the rewrite consumed a neighbouring Gap's blocked_on:\n%s", amended)
	}
	if strings.Count(amended, "resolution: >-") != 1 {
		t.Fatalf("exactly one resolution must be written:\n%s", amended)
	}
	// The appended key belongs to g2, so it must precede g3's entry.
	if strings.Index(amended, "closed by shipping the fix") > strings.Index(amended, "- ref: g3") {
		t.Fatalf("the resolution landed outside the Gap it names:\n%s", amended)
	}
	for _, line := range []string{"the second problem, recorded", "with no blocked_on key", "the third problem"} {
		if !strings.Contains(amended, line) {
			t.Fatalf("the rewrite damaged prose it did not own\n  lost: %q\n  got:\n%s", line, amended)
		}
	}
	if err := assertOnlyAmendableFieldsChanged(contract, amended); err != nil {
		t.Fatalf("the rewrite reached outside the amendable fields: %v", err)
	}
	if !strings.HasSuffix(amended, "body\n") {
		t.Fatalf("the Markdown body was modified:\n%s", amended)
	}
}
