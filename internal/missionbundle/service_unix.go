//go:build !windows

package missionbundle

import (
	"os"
	"syscall"
)

// lockFile takes an exclusive, non-blocking advisory lock. LOCK_NB is what
// turns a busy workspace into an immediate refusal rather than a wait: a
// second mutation must refuse, never queue behind the first and then apply.
// The lock is released when the descriptor closes, so a killed process does
// not strand it.
func lockFile(file *os.File) error {
	return syscall.Flock(int(file.Fd()), syscall.LOCK_EX|syscall.LOCK_NB)
}

func unlockFile(file *os.File) error {
	return syscall.Flock(int(file.Fd()), syscall.LOCK_UN)
}
