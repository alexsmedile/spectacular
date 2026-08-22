package tokenizer

import "testing"

// These counts were generated with tiktoken 0.14.0's o200k_base encoding on
// 2026-08-23. They are deliberately varied: prose, code, Unicode, emoji, and
// URLs all expose counters that merely approximate words or bytes.
func TestCountMatchesPinnedO200KBaseOracle(t *testing.T) {
	tests := []struct {
		text string
		want int
	}{
		{"Hello world", 2},
		{"2 + 2 = 4", 7},
		{"antidisestablishmentarianism", 6},
		{"お誕生日おめでとう", 8},
		{"😀", 1},
		{"function add(a, b) { return a + b; }", 13},
		{"Caffè è già pronto.", 7},
		{"https://example.com/a?b=c&d=e", 11},
		{"foo_barBaz123", 4},
		{"🚀✨\n", 4},
	}
	for _, tt := range tests {
		got, err := Count(tt.text)
		if err != nil {
			t.Fatalf("Count(%q): %v", tt.text, err)
		}
		if got != tt.want {
			t.Errorf("Count(%q) = %d, want official o200k_base count %d", tt.text, got, tt.want)
		}
	}
}
