package workspace

import (
	"bytes"
	"errors"
	"os"
	"path/filepath"
	"testing"

	"github.com/alexsmedile/spectacular/v2/internal/domain"
)

func TestWriteFileUsesCanonicalAtomicReplacement(t *testing.T) {
	t.Parallel()

	document, err := ReadFile(fixturePath("positive", "proposal.md"))
	if err != nil {
		t.Fatal(err)
	}
	want, err := Canonical(document)
	if err != nil {
		t.Fatal(err)
	}
	path := filepath.Join(t.TempDir(), "proposal.md")
	if err := WriteFile(path, document); err != nil {
		t.Fatalf("WriteFile: %v", err)
	}
	got, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(got, want) {
		t.Fatalf("written content is not canonical:\ngot:\n%s\nwant:\n%s", got, want)
	}
}

func TestFailedWritePreservesOriginalFile(t *testing.T) {
	t.Parallel()

	document, err := ReadFile(fixturePath("positive", "proposal.md"))
	if err != nil {
		t.Fatal(err)
	}
	directory := t.TempDir()
	path := filepath.Join(directory, "proposal.md")
	original := []byte("original bytes must survive\n")
	if err := os.WriteFile(path, original, 0o640); err != nil {
		t.Fatal(err)
	}

	replaceFailure := errors.New("injected replacement failure")
	err = writeFileWithReplace(path, document, func(_, _ string) error {
		return replaceFailure
	})
	if !domain.RefusalHasCode(err, domain.RefusalPersistence) || !errors.Is(err, replaceFailure) {
		t.Fatalf("writeFileWithReplace error = %v, want persistence refusal wrapping injected failure", err)
	}
	got, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(got, original) {
		t.Fatalf("failed write changed original:\ngot:  %q\nwant: %q", got, original)
	}
	entries, err := os.ReadDir(directory)
	if err != nil {
		t.Fatal(err)
	}
	if len(entries) != 1 || entries[0].Name() != "proposal.md" {
		t.Fatalf("failed write left temporary files: %#v", entries)
	}
}
