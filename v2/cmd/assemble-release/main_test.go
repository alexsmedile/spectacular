package main

import (
	"archive/tar"
	"bytes"
	"compress/gzip"
	"io"
	"testing"
	"time"
)

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
