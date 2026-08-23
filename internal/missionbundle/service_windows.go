//go:build windows

package missionbundle

import (
	"os"

	"golang.org/x/sys/windows"
)

// lockRegion is the byte range LockFileEx operates on. The lock is advisory
// between Spectacular processes rather than a content guard, so locking a
// single byte is enough and works on an empty file.
const lockRegion = 1

// lockFile takes an exclusive, non-blocking lock on the workspace lock file.
// LOCKFILE_FAIL_IMMEDIATELY is the Windows counterpart of LOCK_NB and is what
// keeps the refusal semantics identical to the Unix path: a second mutation
// refuses at once instead of queueing behind the first and then applying.
// Windows releases the lock when the handle closes, including on abnormal
// termination, so a killed process does not strand it.
func lockFile(file *os.File) error {
	return windows.LockFileEx(
		windows.Handle(file.Fd()),
		windows.LOCKFILE_EXCLUSIVE_LOCK|windows.LOCKFILE_FAIL_IMMEDIATELY,
		0, lockRegion, 0,
		new(windows.Overlapped),
	)
}

func unlockFile(file *os.File) error {
	return windows.UnlockFileEx(
		windows.Handle(file.Fd()),
		0, lockRegion, 0,
		new(windows.Overlapped),
	)
}
