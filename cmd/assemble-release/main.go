// assemble-release builds the complete deterministic native v2 distribution.
// It is a build-time tool only; end-user installation never invokes Go.
package main

import (
	"archive/tar"
	"bytes"
	"compress/gzip"
	"crypto/sha256"
	"debug/elf"
	"debug/macho"
	"encoding/hex"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io"
	"io/fs"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
	"sync"
	"time"

	"github.com/alexsmedile/spectacular/v2/internal/command"
)

const buildinfoPackage = "github.com/alexsmedile/spectacular/v2/internal/buildinfo"

type target struct {
	os   string
	arch string
}

type archiveEntry struct {
	name string
	mode int64
	data []byte
}

type releaseMetadata struct {
	SchemaVersion string `json:"schema_version"`
	Name          string `json:"name"`
	Version       string `json:"version"`
	Commit        string `json:"commit"`
	OS            string `json:"os"`
	Architecture  string `json:"architecture"`
}

var semver = regexp.MustCompile(`^[0-9]+\.[0-9]+\.[0-9]+(?:-[0-9A-Za-z.-]+)?$`)

func main() {
	output := flag.String("output", "", "empty or absent release output directory")
	versionFile := flag.String("version-file", "VERSION", "canonical release version file")
	commit := flag.String("commit", "unknown", "source commit embedded in release binaries")
	flag.Parse()
	if *output == "" || flag.NArg() != 0 {
		fatal(errors.New("usage: assemble-release --output <directory> [--version-file VERSION] [--commit SHA]"))
	}
	root, err := os.Getwd()
	if err != nil {
		fatal(err)
	}
	version, err := readVersion(filepath.Join(root, *versionFile))
	if err != nil {
		fatal(err)
	}
	if err := validateSources(root, version); err != nil {
		fatal(err)
	}
	if err := prepareOutput(*output); err != nil {
		fatal(err)
	}
	targets := []target{{"darwin", "amd64"}, {"darwin", "arm64"}, {"linux", "amd64"}, {"linux", "arm64"}}
	type buildResult struct {
		name     string
		data     []byte
		checksum string
		err      error
	}
	results := make([]buildResult, len(targets))
	var wg sync.WaitGroup
	for i, platform := range targets {
		wg.Add(1)
		go func(i int, platform target) {
			defer wg.Done()
			name := fmt.Sprintf("spectacular-v%s-%s-%s.tar.gz", version, platform.os, platform.arch)
			data, err := buildArchive(root, version, *commit, platform)
			if err != nil {
				results[i] = buildResult{err: fmt.Errorf("%s/%s: %w", platform.os, platform.arch, err)}
				return
			}
			digest := sha256.Sum256(data)
			results[i] = buildResult{
				name:     name,
				data:     data,
				checksum: hex.EncodeToString(digest[:]) + "  " + name,
			}
		}(i, platform)
	}
	wg.Wait()
	checksums := make([]string, 0, len(results))
	for _, res := range results {
		if res.err != nil {
			fatal(res.err)
		}
		if err := writeAtomic(filepath.Join(*output, res.name), res.data, 0o644); err != nil {
			fatal(err)
		}
		checksums = append(checksums, res.checksum)
	}
	sort.Strings(checksums)
	if err := writeAtomic(filepath.Join(*output, "SHA256SUMS"), []byte(strings.Join(checksums, "\n")+"\n"), 0o644); err != nil {
		fatal(err)
	}
	if err := writeAtomic(filepath.Join(*output, "VERSION"), []byte(version+"\n"), 0o644); err != nil {
		fatal(err)
	}
	fmt.Printf("release_version=%s\nartifacts=%d\noutput=%s\nresult=locally-assembled-unpublished\n", version, len(targets), *output)
}

func readVersion(path string) (string, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return "", fmt.Errorf("read version: %w", err)
	}
	version := strings.TrimSpace(string(data))
	if !semver.MatchString(version) {
		return "", fmt.Errorf("invalid release version %q", version)
	}
	return version, nil
}

func validateSources(root, version string) error {
	for _, manifest := range []string{".codex-plugin/plugin.json", ".claude-plugin/plugin.json"} {
		data, err := os.ReadFile(filepath.Join(root, manifest))
		if err != nil {
			return err
		}
		var value struct {
			Name    string `json:"name"`
			Version string `json:"version"`
		}
		if err := json.Unmarshal(data, &value); err != nil {
			return fmt.Errorf("parse %s: %w", manifest, err)
		}
		if value.Name != "spectacular" || value.Version != version {
			return fmt.Errorf("%s is not aligned to spectacular %s", manifest, version)
		}
	}
	checks := []struct {
		path   string
		render func(io.Writer) error
	}{
		{"skills/spectacular/generated/mechanical-interface.json", func(w io.Writer) error { return command.WriteCatalogJSONVersion(w, version) }},
		{"skills/spectacular/generated/mechanical-interface.md", func(w io.Writer) error { return command.WriteCatalogMarkdownVersion(w, version) }},
	}
	for _, check := range checks {
		var want bytes.Buffer
		if err := check.render(&want); err != nil {
			return err
		}
		got, err := os.ReadFile(filepath.Join(root, check.path))
		if err != nil {
			return err
		}
		if !bytes.Equal(got, want.Bytes()) {
			return fmt.Errorf("generated interface is stale: %s", check.path)
		}
	}
	skillRoot := filepath.Join(root, "skills", "spectacular")
	if err := validateSkill(skillRoot); err != nil {
		return err
	}
	skill, err := os.ReadFile(filepath.Join(skillRoot, "SKILL.md"))
	if err != nil {
		return err
	}
	if !bytes.Contains(skill, []byte("\nversion: "+version+"\n")) {
		return errors.New("canonical Skill version is not aligned")
	}
	return nil
}

func validateSkill(root string) error {
	banned := []string{"spectacular status", "spectacular new", "requests/", "universal doctor", "fallback reader"}
	return filepath.WalkDir(root, func(path string, entry fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if entry.IsDir() {
			return nil
		}
		if !entry.Type().IsRegular() {
			return fmt.Errorf("non-regular Skill entry: %s", path)
		}
		data, err := os.ReadFile(path)
		if err != nil {
			return err
		}
		for _, phrase := range banned {
			if bytes.Contains(bytes.ToLower(data), []byte(phrase)) {
				return fmt.Errorf("v1 runtime language %q in %s", phrase, path)
			}
		}
		return nil
	})
}

func prepareOutput(path string) error {
	info, err := os.Stat(path)
	if errors.Is(err, os.ErrNotExist) {
		return os.MkdirAll(path, 0o755)
	}
	if err != nil {
		return err
	}
	if !info.IsDir() {
		return errors.New("output is not a directory")
	}
	entries, err := os.ReadDir(path)
	if err != nil {
		return err
	}
	if len(entries) != 0 {
		return errors.New("output directory must be empty")
	}
	return nil
}

func buildArchive(root, version, commit string, platform target) ([]byte, error) {
	temporary, err := os.MkdirTemp("", "spectacular-release-")
	if err != nil {
		return nil, err
	}
	defer os.RemoveAll(temporary)
	binary := filepath.Join(temporary, "spectacular")
	ldflags := fmt.Sprintf("-s -w -X %s.Version=%s -X %s.Commit=%s", buildinfoPackage, version, buildinfoPackage, commit)
	build := exec.Command("go", "build", "-trimpath", "-buildvcs=false", "-ldflags", ldflags, "-o", binary, "./cmd/spectacular")
	build.Dir = root
	build.Env = append(os.Environ(), "CGO_ENABLED=0", "GOOS="+platform.os, "GOARCH="+platform.arch, "GOTOOLCHAIN=local")
	if output, err := build.CombinedOutput(); err != nil {
		return nil, fmt.Errorf("go build: %w: %s", err, output)
	}
	binaryData, err := os.ReadFile(binary)
	if err != nil {
		return nil, err
	}
	if err := verifyBinaryPlatform(binary, platform); err != nil {
		return nil, err
	}
	metadata, err := json.MarshalIndent(releaseMetadata{
		SchemaVersion: "spectacular.release.v1", Name: "spectacular", Version: version,
		Commit: commit, OS: platform.os, Architecture: platform.arch,
	}, "", "  ")
	if err != nil {
		return nil, err
	}
	metadata = append(metadata, '\n')
	entries := []archiveEntry{
		{name: "spectacular/RELEASE.json", mode: 0o644, data: metadata},
		{name: "spectacular/VERSION", mode: 0o644, data: []byte(version + "\n")},
		{name: "spectacular/bin/spectacular", mode: 0o755, data: binaryData},
	}
	for _, source := range []struct{ from, to string }{
		{".codex-plugin/plugin.json", "spectacular/plugins/spectacular/.codex-plugin/plugin.json"},
		{".claude-plugin/plugin.json", "spectacular/plugins/spectacular/.claude-plugin/plugin.json"},
	} {
		data, err := os.ReadFile(filepath.Join(root, source.from))
		if err != nil {
			return nil, err
		}
		entries = append(entries, archiveEntry{name: source.to, mode: 0o644, data: data})
	}
	skillRoot := filepath.Join(root, "skills", "spectacular")
	err = filepath.WalkDir(skillRoot, func(path string, entry fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if entry.IsDir() {
			return nil
		}
		if !entry.Type().IsRegular() {
			return fmt.Errorf("non-regular Skill entry: %s", path)
		}
		relative, err := filepath.Rel(skillRoot, path)
		if err != nil {
			return err
		}
		data, err := os.ReadFile(path)
		if err != nil {
			return err
		}
		name := "spectacular/plugins/spectacular/skills/spectacular/" + filepath.ToSlash(relative)
		entries = append(entries, archiveEntry{name: name, mode: 0o644, data: data})
		return nil
	})
	if err != nil {
		return nil, err
	}
	sort.Slice(entries, func(i, j int) bool { return entries[i].name < entries[j].name })
	return encodeArchive(entries)
}

func verifyBinaryPlatform(path string, platform target) error {
	switch platform.os {
	case "darwin":
		file, err := macho.Open(path)
		if err != nil {
			return fmt.Errorf("parse Mach-O binary: %w", err)
		}
		defer file.Close()
		want := macho.CpuAmd64
		if platform.arch == "arm64" {
			want = macho.CpuArm64
		}
		if file.Cpu != want {
			return fmt.Errorf("Mach-O CPU is %s, want %s", file.Cpu, want)
		}
	case "linux":
		file, err := elf.Open(path)
		if err != nil {
			return fmt.Errorf("parse ELF binary: %w", err)
		}
		defer file.Close()
		want := elf.EM_X86_64
		if platform.arch == "arm64" {
			want = elf.EM_AARCH64
		}
		if file.Machine != want {
			return fmt.Errorf("ELF machine is %s, want %s", file.Machine, want)
		}
	default:
		return fmt.Errorf("unsupported target OS %q", platform.os)
	}
	return nil
}

func encodeArchive(entries []archiveEntry) ([]byte, error) {
	var output bytes.Buffer
	gzipWriter, err := gzip.NewWriterLevel(&output, gzip.BestCompression)
	if err != nil {
		return nil, err
	}
	gzipWriter.Header.ModTime = time.Unix(0, 0).UTC()
	gzipWriter.Header.OS = 255
	tarWriter := tar.NewWriter(gzipWriter)
	zero := time.Unix(0, 0).UTC()
	for _, entry := range entries {
		header := &tar.Header{
			Name: entry.name, Mode: entry.mode, Size: int64(len(entry.data)), Typeflag: tar.TypeReg,
			ModTime: zero, Uid: 0, Gid: 0, Uname: "root", Gname: "root", Format: tar.FormatUSTAR,
		}
		if err := tarWriter.WriteHeader(header); err != nil {
			return nil, err
		}
		if _, err := tarWriter.Write(entry.data); err != nil {
			return nil, err
		}
	}
	if err := tarWriter.Close(); err != nil {
		return nil, err
	}
	if err := gzipWriter.Close(); err != nil {
		return nil, err
	}
	return output.Bytes(), nil
}

func writeAtomic(path string, data []byte, mode fs.FileMode) error {
	temporary := path + ".tmp"
	if err := os.WriteFile(temporary, data, mode); err != nil {
		return err
	}
	if err := os.Chmod(temporary, mode); err != nil {
		return err
	}
	return os.Rename(temporary, path)
}

func fatal(err error) {
	fmt.Fprintln(os.Stderr, "release refused:", err)
	os.Exit(3)
}
