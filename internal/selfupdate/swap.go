package selfupdate

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
)

// ReplaceBinary installs data at path, keeping the displaced binary alongside
// it so a bad update can be undone.
//
// The swap is a rename, never a write over the target. A running program cannot
// be overwritten on Windows and is corrupted mid-flight if overwritten on Unix,
// but a rename is atomic on both: the path points at the old binary or the new
// one and never at a half-written file. The process that called this keeps
// executing from the displaced inode until it exits.
//
// The staged file is written into the target's own directory, because a rename
// across filesystems is not atomic and fails outright on most systems.
func ReplaceBinary(path string, data []byte) (backup string, err error) {
	if len(data) == 0 {
		return "", errors.New("refusing to install an empty binary")
	}
	directory := filepath.Dir(path)
	staged, err := os.CreateTemp(directory, ".spectacular-update-*")
	if err != nil {
		return "", fmt.Errorf("stage the new binary: %w", err)
	}
	stagedPath := staged.Name()
	defer func() {
		if err != nil {
			os.Remove(stagedPath)
		}
	}()
	if _, err = staged.Write(data); err != nil {
		staged.Close()
		return "", fmt.Errorf("stage the new binary: %w", err)
	}
	// Sync before the rename: a crash between the two would otherwise leave the
	// target pointing at a file whose contents never reached the disk.
	if err = staged.Sync(); err != nil {
		staged.Close()
		return "", fmt.Errorf("stage the new binary: %w", err)
	}
	if err = staged.Close(); err != nil {
		return "", fmt.Errorf("stage the new binary: %w", err)
	}
	if err = os.Chmod(stagedPath, 0o755); err != nil {
		return "", fmt.Errorf("stage the new binary: %w", err)
	}

	backup = path + ".old"
	os.Remove(backup)
	// Displace rather than delete. Until the second rename lands, the previous
	// binary is still on disk under a known name, so a failure here is
	// recoverable by hand rather than leaving no spectacular at all.
	if err = os.Rename(path, backup); err != nil && !errors.Is(err, os.ErrNotExist) {
		return "", fmt.Errorf("move the current binary aside: %w", err)
	}
	if err = os.Rename(stagedPath, path); err != nil {
		if restoreErr := os.Rename(backup, path); restoreErr != nil {
			return "", fmt.Errorf("install the new binary: %w (the previous binary is at %s and must be restored by hand)", err, backup)
		}
		return "", fmt.Errorf("install the new binary: %w (the previous binary was restored)", err)
	}
	return backup, nil
}

// TargetPath reports the binary an update would replace: the running executable
// itself, with symlinks resolved so an update rewrites the real file rather
// than replacing a link that points at it.
func TargetPath() (string, error) {
	executable, err := os.Executable()
	if err != nil {
		return "", fmt.Errorf("locate the running binary: %w", err)
	}
	resolved, err := filepath.EvalSymlinks(executable)
	if err != nil {
		return executable, nil
	}
	return resolved, nil
}

// Writable reports whether the binary can be replaced in place. An update needs
// to create and rename files in the binary's directory, so the directory's
// permissions decide, not the binary's.
func Writable(path string) error {
	directory := filepath.Dir(path)
	probe, err := os.CreateTemp(directory, ".spectacular-probe-*")
	if err != nil {
		return fmt.Errorf("%s is not writable: %w", directory, err)
	}
	name := probe.Name()
	probe.Close()
	os.Remove(name)
	return nil
}
