package guardrails

import (
	"reflect"
	"testing"
)

func TestParseAndSelectReturnsVerbatimOwnerProse(t *testing.T) {
	doc, err := Parse([]byte("# Guardrails\n\n## @Run\nKeep the baseline exact.\n\n- Stop on drift.\n## @Run $git.commit\nCommit only coherent changes.\n## @Resolve\nOwner disposition is required.\n"))
	if err != nil {
		t.Fatal(err)
	}
	got, err := doc.Select("@Run", "$git.commit")
	if err != nil {
		t.Fatal(err)
	}
	want := []Section{{Event: "@Run", Prose: "Keep the baseline exact.\n\n- Stop on drift."}, {Event: "@Run", Selector: "$git.commit", Prose: "Commit only coherent changes."}}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("selected %#v, want %#v", got, want)
	}
}

func TestParseRefusesUnknownDuplicateAndInvalidSelectors(t *testing.T) {
	for _, input := range []string{
		"## @Launch\nno\n",
		"## @Run $bad\nno\n",
		"## @Run\none\n## @Run\ntwo\n",
		"## @Run\n\n",
	} {
		if _, err := Parse([]byte(input)); err == nil {
			t.Fatalf("accepted invalid Guardrails %q", input)
		}
	}
}
