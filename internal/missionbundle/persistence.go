package missionbundle

import (
	"bytes"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/alexsmedile/spectacular/v2/internal/discovery"
	"github.com/alexsmedile/spectacular/v2/internal/domain"
	"github.com/alexsmedile/spectacular/v2/internal/governance"
	"github.com/alexsmedile/spectacular/v2/internal/humanlayout"
	"github.com/alexsmedile/spectacular/v2/internal/workspace"
)

func (s Service) apply(key string, docs []*workspace.Document, paths map[domain.ID]string, operation, ref, primaryPath string) (Result, error) {
	for _, doc := range docs {
		known := false
		for _, entry := range s.Workspace.Entries {
			if entry.Document.Record.ID != doc.Record.ID {
				continue
			}
			known = true
			current, readErr := workspace.ReadFile(entry.Absolute)
			if readErr != nil {
				return Result{}, readErr
			}
			fingerprint, fingerprintErr := workspace.Fingerprint(current)
			if fingerprintErr != nil {
				return Result{}, fingerprintErr
			}
			if fingerprint != entry.Fingerprint {
				return Result{}, domain.NewStateRefusal(domain.RefusalStaleFingerprint, entry.Path, "canonical source changed after command validation", entry.Fingerprint, fingerprint, "reload the Mission and retry the typed command", nil)
			}
			break
		}
		if !known {
			target := filepath.Clean(filepath.FromSlash(paths[doc.Record.ID]))
			if filepath.IsAbs(target) || target == ".." || strings.HasPrefix(target, ".."+string(filepath.Separator)) {
				return Result{}, domain.NewStateRefusal(domain.RefusalPathEscape, "path", "new record target escapes the workspace", "canonical workspace-relative path", paths[doc.Record.ID], "use the schema-derived path inside .spectacular", nil)
			}
			if _, statErr := os.Lstat(filepath.Join(s.Workspace.Root, target)); statErr == nil {
				return Result{}, domain.NewStateRefusal(domain.RefusalCollision, filepath.ToSlash(target), "new record target already exists", "unused schema-derived path", "occupied", "inspect the existing file and choose a distinct title or resolve the collision", nil)
			} else if !os.IsNotExist(statErr) {
				return Result{}, invalidCause(filepath.ToSlash(target), "inspect new record target", statErr)
			}
		}
	}
	changes := make([]governance.FileChange, 0, len(docs)+4)
	for _, doc := range docs {
		data, err := workspace.Canonical(doc)
		if err != nil {
			return Result{}, err
		}
		changes = append(changes, governance.FileChange{Path: paths[doc.Record.ID], Data: data, Mode: 0o644})
	}
	indexes, err := humanlayout.Indexes(s.Workspace.Entries, docs, paths)
	if err != nil {
		return Result{}, err
	}
	indexPaths := make([]string, 0, len(indexes))
	for path := range indexes {
		indexPaths = append(indexPaths, path)
	}
	sort.Strings(indexPaths)
	for _, path := range indexPaths {
		existingData, readErr := os.ReadFile(filepath.Join(s.Workspace.Root, filepath.FromSlash(path)))
		if readErr == nil && bytes.Equal(existingData, indexes[path]) {
			continue
		}
		changes = append(changes, governance.FileChange{Path: path, Data: indexes[path], Mode: 0o644})
	}
	apply := s.ApplyTransaction
	if apply == nil {
		apply = governance.ApplyTransaction
	}
	if err := apply(s.Workspace.Root, key, changes); err != nil {
		return Result{}, err
	}
	changed := make([]string, 0, len(changes))
	for _, change := range changes {
		changed = append(changed, filepath.ToSlash(change.Path))
	}
	return Result{Operation: operation, Ref: ref, Path: primaryPath, Changed: changed}, nil
}

func (s Service) beginMutation() (Service, func(), error) {
	if s.Workspace == nil {
		return Service{}, nil, invalid("workspace", "workspace is required")
	}
	unlock, err := acquireMutationLock(s.Workspace.Root)
	if err != nil {
		return Service{}, nil, err
	}
	fresh, err := discovery.Open(s.Workspace.Root)
	if err != nil {
		unlock()
		return Service{}, nil, err
	}
	s.Workspace = fresh
	return s, unlock, nil
}

// acquireMutationLock takes the workspace-wide exclusive lock that makes a
// concurrent mutation refuse instead of interleaving. The lock itself is
// load-bearing; only the kernel call that takes it is platform-specific, so
// lockFile and unlockFile are supplied by service_unix.go and
// service_windows.go.
func acquireMutationLock(root string) (func(), error) {
	directory := filepath.Join(root, ".spectacular", ".engine")
	if err := os.MkdirAll(directory, 0o755); err != nil {
		return nil, invalidCause("engine", "create mutation lock directory", err)
	}
	file, err := os.OpenFile(filepath.Join(directory, ".lock"), os.O_CREATE|os.O_RDWR, 0o600)
	if err != nil {
		return nil, invalidCause("engine", "open mutation lock", err)
	}
	if err := lockFile(file); err != nil {
		file.Close()
		return nil, domain.NewStateRefusal(domain.RefusalCollision, "transactions", "another Mission mutation is in progress", "exclusive mutation lock", "busy", "wait for the active command to finish, reload the Mission, and retry", err)
	}
	return func() {
		_ = unlockFile(file)
		_ = file.Close()
	}, nil
}
