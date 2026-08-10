// Package discovery locates and reads an explicitly marked v2 workspace.
package discovery

import (
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"unicode/utf8"

	"github.com/alexsmedile/spectacular/v2/internal/domain"
	"github.com/alexsmedile/spectacular/v2/internal/workspace"
	"go.yaml.in/yaml/v3"
)

const SchemaVersion = "spectacular.workspace.v1"

type Manifest struct {
	SchemaVersion string   `yaml:"schema_version"`
	RecordRoots   []string `yaml:"record_roots"`
	ProjectAnchor string   `yaml:"project_anchor"`
	Guardrails    string   `yaml:"guardrails,omitempty"`
}

type Entry struct {
	Path        string
	Absolute    string
	Document    *workspace.Document
	Fingerprint string
}

type Workspace struct {
	Root        string
	MetadataDir string
	Manifest    Manifest
	Entries     []Entry
	byID        map[domain.ID]Entry
	byPath      map[string]Entry
	byHuman     map[string]Entry
}

func Open(start string) (*Workspace, error) {
	abs, err := filepath.Abs(start)
	if err != nil {
		return nil, refusal(domain.RefusalWorkspaceNotFound, start, "resolve starting directory", err)
	}
	info, err := os.Stat(abs)
	if err != nil {
		return nil, refusal(domain.RefusalWorkspaceNotFound, start, "inspect starting directory", err)
	}
	if !info.IsDir() {
		abs = filepath.Dir(abs)
	}
	for {
		meta := filepath.Join(abs, ".spectacular")
		metaInfo, metaErr := os.Lstat(meta)
		if metaErr == nil && metaInfo.Mode()&os.ModeSymlink != 0 {
			return nil, refusal(domain.RefusalPathEscape, meta, ".spectacular metadata directory must not be a symlink", nil)
		}
		if metaErr != nil && !os.IsNotExist(metaErr) {
			return nil, refusal(domain.RefusalInvalidManifest, meta, "inspect workspace metadata directory", metaErr)
		}
		marker := filepath.Join(meta, "workspace.yaml")
		if markerInfo, err := os.Lstat(marker); err == nil {
			if markerInfo.Mode()&os.ModeSymlink != 0 {
				return nil, refusal(domain.RefusalPathEscape, marker, "workspace marker must not be a symlink", nil)
			}
			if !markerInfo.Mode().IsRegular() {
				return nil, refusal(domain.RefusalInvalidManifest, marker, "workspace marker must be a regular file", nil)
			}
			return load(abs, marker)
		} else if !os.IsNotExist(err) {
			return nil, refusal(domain.RefusalInvalidManifest, marker, "inspect marker", err)
		}
		parent := filepath.Dir(abs)
		if parent == abs {
			break
		}
		abs = parent
	}
	return nil, refusal(domain.RefusalWorkspaceNotFound, start, "no explicit v2 .spectacular/workspace.yaml found", nil)
}

func load(root, marker string) (*Workspace, error) {
	data, err := os.ReadFile(marker)
	if err != nil {
		return nil, refusal(domain.RefusalInvalidManifest, marker, "read marker", err)
	}
	if !utf8.Valid(data) {
		return nil, refusal(domain.RefusalInvalidManifest, marker, "marker must be valid UTF-8", nil)
	}
	var manifest Manifest
	if err := yaml.Unmarshal(data, &manifest); err != nil {
		return nil, refusal(domain.RefusalInvalidManifest, marker, "decode marker", err)
	}
	if manifest.SchemaVersion != SchemaVersion {
		return nil, refusal(domain.RefusalInvalidManifest, "schema_version", "expected "+SchemaVersion, nil)
	}
	if len(manifest.RecordRoots) == 0 || manifest.ProjectAnchor == "" {
		return nil, refusal(domain.RefusalInvalidManifest, marker, "record_roots and project_anchor are required", nil)
	}
	meta := filepath.Join(root, ".spectacular")
	metaReal, err := filepath.EvalSymlinks(meta)
	if err != nil {
		return nil, refusal(domain.RefusalPathEscape, meta, "resolve workspace metadata directory", err)
	}
	var paths []string
	seenRoots := map[string]bool{}
	for _, declared := range manifest.RecordRoots {
		if seenRoots[declared] {
			return nil, refusal(domain.RefusalInvalidManifest, "record_roots", "roots must be unique canonical relative paths", nil)
		}
		seenRoots[declared] = true
	}
	seenRoots = map[string]bool{}
	scanRoots := append([]string(nil), manifest.RecordRoots...)
	// Human-layout collections are canonical v2 roots. Discover any that
	// exist even when an older v2 fixture names only its seed root; mutations
	// never need to rewrite workspace.yaml merely to add an earned collection.
	if !containsRoot(manifest.RecordRoots, ".") {
		for _, standard := range []string{"contracts", "proposals", "missions", "evidence", "decisions", "gaps", "handoffs", "assessments", "archive/missions"} {
			if _, statErr := os.Stat(filepath.Join(meta, filepath.FromSlash(standard))); statErr == nil && !containsRoot(scanRoots, standard) {
				scanRoots = append(scanRoots, standard)
			}
		}
	}
	for _, declared := range scanRoots {
		if (declared != "." && !canonicalRelative(declared)) || seenRoots[declared] {
			return nil, refusal(domain.RefusalInvalidManifest, "record_roots", "roots must be unique canonical relative paths", nil)
		}
		seenRoots[declared] = true
		rootPath := filepath.Join(meta, filepath.FromSlash(declared))
		real, err := filepath.EvalSymlinks(rootPath)
		if err != nil {
			return nil, refusal(domain.RefusalInvalidManifest, declared, "declared root is missing or unreadable", err)
		}
		if !within(metaReal, real) {
			return nil, refusal(domain.RefusalPathEscape, declared, "declared root escapes .spectacular", nil)
		}
		err = filepath.WalkDir(rootPath, func(path string, d os.DirEntry, walkErr error) error {
			if walkErr != nil {
				return walkErr
			}
			if d.Type()&os.ModeSymlink != 0 {
				target, e := filepath.EvalSymlinks(path)
				if e != nil {
					return e
				}
				if !within(metaReal, target) {
					return refusal(domain.RefusalPathEscape, path, "symlink escapes .spectacular", nil)
				}
				if d.IsDir() {
					return filepath.SkipDir
				}
			}
			if d.IsDir() && (d.Name() == "history" || d.Name() == "transactions") {
				return filepath.SkipDir
			}
			// Generated indexes are committed navigation aids, never canonical
			// records. They must be rebuildable and must not enter authority.
			if !d.IsDir() && d.Name() != "index.md" && d.Name() != "GUARDRAILS.md" && strings.HasSuffix(d.Name(), ".md") {
				paths = append(paths, path)
			}
			return nil
		})
		if err != nil {
			var refused *domain.Refusal
			if errors.As(err, &refused) {
				return nil, err
			}
			return nil, refusal(domain.RefusalInvalidManifest, declared, "scan declared root", err)
		}
	}
	sort.Strings(paths)
	w := &Workspace{Root: root, MetadataDir: meta, Manifest: manifest, byID: map[domain.ID]Entry{}, byPath: map[string]Entry{}, byHuman: map[string]Entry{}}
	for _, absolute := range paths {
		real, err := filepath.EvalSymlinks(absolute)
		if err != nil || !within(metaReal, real) {
			return nil, refusal(domain.RefusalPathEscape, absolute, "record escapes .spectacular", err)
		}
		doc, err := workspace.ReadFile(absolute)
		if err != nil {
			return nil, err
		}
		fp, err := workspace.Fingerprint(doc)
		if err != nil {
			return nil, fmt.Errorf("fingerprint %q: %w", absolute, err)
		}
		rel, _ := filepath.Rel(root, absolute)
		public := filepath.ToSlash(rel)
		entry := Entry{Path: public, Absolute: absolute, Document: doc, Fingerprint: fp}
		if previous, ok := w.byPath[public]; ok {
			return nil, refusal(domain.RefusalDuplicatePath, public, "also claimed by "+previous.Document.Record.ID.String(), nil)
		}
		if previous, ok := w.byID[doc.Record.ID]; ok {
			return nil, refusal(domain.RefusalDuplicateID, doc.Record.ID.String(), "claimed by "+previous.Path+" and "+public, nil)
		}
		w.byPath[public] = entry
		w.byID[doc.Record.ID] = entry
		if human, humanErr := workspace.String(doc, "human_ref", false); humanErr != nil {
			return nil, humanErr
		} else if human != "" {
			if previous, exists := w.byHuman[human]; exists {
				return nil, refusal(domain.RefusalDuplicateID, human, "claimed by "+previous.Path+" and "+public, nil)
			}
			w.byHuman[human] = entry
		}
		w.Entries = append(w.Entries, entry)
	}
	anchorPath := ".spectacular/" + manifest.ProjectAnchor
	anchor, ok := w.byPath[anchorPath]
	if !ok {
		return nil, refusal(domain.RefusalRecordNotFound, anchorPath, "project_anchor does not resolve", nil)
	}
	if anchor.Document.Record.Type != domain.Anchor {
		return nil, refusal(domain.RefusalNounMismatch, anchorPath, "project_anchor must be Anchor", nil)
	}
	return w, nil
}

func containsRoot(roots []string, want string) bool {
	for _, root := range roots {
		if root == want {
			return true
		}
	}
	return false
}

func (w *Workspace) ProjectAnchor() Entry { return w.byPath[".spectacular/"+w.Manifest.ProjectAnchor] }

func (w *Workspace) Lookup(ref string, want domain.RecordType) (Entry, error) {
	var entry Entry
	var ok bool
	if strings.Contains(ref, ":") && !strings.HasPrefix(ref, ".spectacular/") {
		typed, err := domain.ParseReference(ref)
		if err != nil {
			return Entry{}, err
		}
		if typed.Type != want {
			return Entry{}, refusal(domain.RefusalNounMismatch, ref, "expected "+string(want), nil)
		}
		entry, ok = w.byID[typed.ID]
	} else if id, err := domain.ParseID(ref); err == nil {
		entry, ok = w.byID[id]
	} else if strings.HasPrefix(ref, ".spectacular/") && canonicalRelative(ref) {
		entry, ok = w.byPath[ref]
	} else if human, found := w.byHuman[ref]; found {
		entry, ok = human, true
	} else {
		return Entry{}, refusal(domain.RefusalInvalidReference, ref, "expected human reference, canonical UUIDv7, typed reference, or full .spectacular/ path", nil)
	}
	if !ok {
		return Entry{}, refusal(domain.RefusalRecordNotFound, ref, "record does not exist", nil)
	}
	if entry.Document.Record.Type != want {
		return Entry{}, refusal(domain.RefusalNounMismatch, ref, "resolved "+string(entry.Document.Record.Type)+", expected "+string(want), nil)
	}
	return entry, nil
}

func (w *Workspace) OfType(t domain.RecordType) []Entry {
	var out []Entry
	for _, entry := range w.Entries {
		if entry.Document.Record.Type == t {
			out = append(out, entry)
		}
	}
	return out
}

// ReadMetadataFile reads an optional manifest-declared support file without
// allowing it to escape the canonical .spectacular directory. Support files
// guide workflows; they are never indexed as canonical records.
func (w *Workspace) ReadMetadataFile(relative string) ([]byte, string, error) {
	if relative == "" || !canonicalRelative(relative) {
		return nil, "", refusal(domain.RefusalInvalidManifest, relative, "metadata file must be a canonical relative path", nil)
	}
	absolute := filepath.Join(w.MetadataDir, filepath.FromSlash(relative))
	metadataRoot, err := filepath.EvalSymlinks(w.MetadataDir)
	if err != nil {
		return nil, "", refusal(domain.RefusalPathEscape, relative, "resolve workspace metadata directory", err)
	}
	real, err := filepath.EvalSymlinks(absolute)
	if err != nil {
		return nil, "", refusal(domain.RefusalRecordNotFound, relative, "declared metadata file does not exist", err)
	}
	if !within(metadataRoot, real) {
		return nil, "", refusal(domain.RefusalPathEscape, relative, "declared metadata file escapes .spectacular", nil)
	}
	info, err := os.Stat(real)
	if err != nil || !info.Mode().IsRegular() {
		return nil, "", refusal(domain.RefusalInvalidManifest, relative, "declared metadata file must be regular", err)
	}
	data, err := os.ReadFile(real)
	if err != nil {
		return nil, "", refusal(domain.RefusalInvalidManifest, relative, "read declared metadata file", err)
	}
	if !utf8.Valid(data) {
		return nil, "", refusal(domain.RefusalInvalidManifest, relative, "declared metadata file must be valid UTF-8", nil)
	}
	sum := sha256.Sum256(data)
	return data, hex.EncodeToString(sum[:]), nil
}

// Source resolves a freshness basis to either an exact typed record or a
// canonical workspace-relative file. It never accepts a basename or follows
// a path outside the discovered workspace.
func (w *Workspace) Source(ref string) (Entry, string, error) {
	if strings.Contains(ref, ":") && !strings.HasPrefix(ref, ".spectacular/") {
		typed, err := domain.ParseReference(ref)
		if err != nil {
			return Entry{}, "", err
		}
		entry, err := w.Lookup(ref, typed.Type)
		return entry, entry.Fingerprint, err
	}
	if !strings.HasPrefix(ref, ".spectacular/") || !canonicalRelative(ref) {
		return Entry{}, "", refusal(domain.RefusalInvalidReference, "freshness_source", "expected exact typed reference or canonical .spectacular/ path", nil)
	}
	abs := filepath.Join(w.Root, filepath.FromSlash(ref))
	realRoot, err := filepath.EvalSymlinks(w.MetadataDir)
	if err != nil {
		return Entry{}, "", err
	}
	real, err := filepath.EvalSymlinks(abs)
	if err != nil {
		return Entry{}, "", refusal(domain.RefusalRecordNotFound, ref, "freshness source does not exist", err)
	}
	if !within(realRoot, real) {
		return Entry{}, "", refusal(domain.RefusalPathEscape, ref, "freshness source escapes .spectacular", nil)
	}
	data, err := os.ReadFile(real)
	if err != nil {
		return Entry{}, "", refusal(domain.RefusalRecordNotFound, ref, "read freshness source", err)
	}
	sum := sha256.Sum256(data)
	return Entry{Path: ref, Absolute: real}, hex.EncodeToString(sum[:]), nil
}

func canonicalRelative(p string) bool {
	return p != "" && !filepath.IsAbs(p) && filepath.ToSlash(filepath.Clean(p)) == p && p != "." && !strings.HasPrefix(p, "../") && !strings.Contains(p, "\\")
}
func within(root, path string) bool {
	rel, err := filepath.Rel(root, path)
	return err == nil && rel != ".." && !strings.HasPrefix(rel, ".."+string(filepath.Separator))
}
func refusal(code domain.RefusalCode, field, detail string, err error) error {
	return domain.NewRefusal(code, field, detail, err)
}
