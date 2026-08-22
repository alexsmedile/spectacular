package charter

import (
	"reflect"
	"testing"
)

func TestDeclaredSourceRefsPreservesContractualOrder(t *testing.T) {
	got := declaredSourceRefs(
		"Contract:one",
		[]string{"D1", "D2", "D1", "  "},
		[]string{"D3", "D2"},
		[]string{"D4", "D3"},
	)
	want := []string{"Contract:one", "D1", "D2", "D3", "D4"}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("declaredSourceRefs() = %#v, want %#v", got, want)
	}
}
