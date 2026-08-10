package workspace

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/alexsmedile/spectacular/v2/internal/domain"
)

type replaceFile func(oldPath, newPath string) error

// WriteFile canonicalizes and atomically replaces path using a temporary file
// in the same directory. Before replacement, every failure removes only the
// temporary file and leaves an existing original untouched.
func WriteFile(path string, document *Document) error {
	return writeFileWithReplace(path, document, os.Rename)
}

func writeFileWithReplace(path string, document *Document, replace replaceFile) error {
	canonical, err := Canonical(document)
	if err != nil {
		return err
	}

	directory := filepath.Dir(path)
	base := filepath.Base(path)
	mode := os.FileMode(0o644)
	if info, statErr := os.Stat(path); statErr == nil {
		mode = info.Mode().Perm()
	} else if !os.IsNotExist(statErr) {
		return persistenceRefusal(path, "inspect original", statErr)
	}

	temporary, err := os.CreateTemp(directory, "."+base+".tmp-*")
	if err != nil {
		return persistenceRefusal(path, "create same-directory temporary file", err)
	}
	temporaryPath := temporary.Name()
	replaced := false
	defer func() {
		if !replaced {
			_ = os.Remove(temporaryPath)
		}
	}()

	if err := temporary.Chmod(mode); err != nil {
		_ = temporary.Close()
		return persistenceRefusal(path, "preserve file mode", err)
	}
	if _, err := temporary.Write(canonical); err != nil {
		_ = temporary.Close()
		return persistenceRefusal(path, "write temporary file", err)
	}
	if err := temporary.Sync(); err != nil {
		_ = temporary.Close()
		return persistenceRefusal(path, "sync temporary file", err)
	}
	if err := temporary.Close(); err != nil {
		return persistenceRefusal(path, "close temporary file", err)
	}
	if err := replace(temporaryPath, path); err != nil {
		return persistenceRefusal(path, "replace original", err)
	}
	replaced = true
	return nil
}

func persistenceRefusal(path, operation string, err error) error {
	return domain.NewRefusal(
		domain.RefusalPersistence,
		path,
		fmt.Sprintf("%s failed", operation),
		err,
	)
}
