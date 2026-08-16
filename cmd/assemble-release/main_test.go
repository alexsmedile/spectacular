package main

import (
	"archive/tar"
	"bytes"
	"compress/gzip"
	"io"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

func TestCanonicalSkillDefinesLeanLaunchAndQuestionContract(t *testing.T) {
	root := filepath.Join("..", "..", "skills", "spectacular")
	parts := []string{"SKILL.md", filepath.Join("references", "prepare.md")}
	var content strings.Builder
	for _, part := range parts {
		data, err := os.ReadFile(filepath.Join(root, part))
		if err != nil {
			t.Fatal(err)
		}
		content.Write(data)
	}
	for _, required := range []string{
		"Read `.spectacular/PROJECT.md` first",
		"read-only launch preflight",
		"plain outcome; technical basis",
		"action -> consequence",
		"recommended default",
	} {
		if !strings.Contains(content.String(), required) {
			t.Fatalf("canonical Skill omits %q", required)
		}
	}
}

func TestCanonicalSkillDefinesLeanExecutionAndFROST(t *testing.T) {
	root := filepath.Join("..", "..", "skills", "spectacular")
	core, err := os.ReadFile(filepath.Join(root, "SKILL.md"))
	if err != nil {
		t.Fatal(err)
	}
	execute, err := os.ReadFile(filepath.Join(root, "references", "execute.md"))
	if err != nil {
		t.Fatal(err)
	}
	audit, err := os.ReadFile(filepath.Join(root, "references", "audit.md"))
	if err != nil {
		t.Fatal(err)
	}
	content := string(core) + string(execute) + string(audit)
	for _, required := range []string{
		"card -> claim packet -> exact sources -> full bundle",
		"disjoint claims + dependencies",
		"focused checks",
		"local green commit",
		"one push boundary",
		"Frozen fit",
		"Truth of proof",
		"active Mission remains governed by the preparation schema",
	} {
		if !strings.Contains(content, required) {
			t.Fatalf("canonical Skill omits %q", required)
		}
	}
	if strings.Contains(string(core), "Frozen fit**") {
		t.Fatal("core Skill embeds detailed FROST policy instead of routing to audit.md")
	}
	if strings.Contains(string(core), "--event <@Event> --json") {
		t.Fatal("core Skill still forces full JSON context")
	}
}

func TestArchiveMetadataIsCanonicalAndReproducible(t *testing.T) {
	entries := []archiveEntry{{name: "spectacular/VERSION", mode: 0o644, data: []byte("2.0.0\n")}}
	first, err := encodeArchive(entries)
	if err != nil {
		t.Fatal(err)
	}
	second, err := encodeArchive(entries)
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(first, second) {
		t.Fatal("identical archive inputs produced different bytes")
	}
	gzipReader, err := gzip.NewReader(bytes.NewReader(first))
	if err != nil {
		t.Fatal(err)
	}
	if !gzipReader.ModTime.IsZero() || gzipReader.OS != 255 || gzipReader.Name != "" || gzipReader.Comment != "" {
		t.Fatalf("non-canonical gzip header: %#v", gzipReader.Header)
	}
	tarReader := tar.NewReader(gzipReader)
	header, err := tarReader.Next()
	if err != nil {
		t.Fatal(err)
	}
	if header.Name != entries[0].name || header.Mode != entries[0].mode || header.Uid != 0 || header.Gid != 0 || header.Uname != "root" || header.Gname != "root" || !header.ModTime.Equal(time.Unix(0, 0).UTC()) {
		t.Fatalf("non-canonical tar header: %#v", header)
	}
	data, err := io.ReadAll(tarReader)
	if err != nil || !bytes.Equal(data, entries[0].data) {
		t.Fatalf("archive payload mismatch: %q err=%v", data, err)
	}
}
