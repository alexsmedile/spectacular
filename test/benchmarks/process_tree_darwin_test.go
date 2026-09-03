//go:build darwin

package spectaculareval

import (
	"errors"
	"syscall"
)

func processGroupHasLiveMember(groupID int) (bool, error) {
	err := syscall.Kill(-groupID, 0)
	if errors.Is(err, syscall.ESRCH) {
		return false, nil
	}
	return err == nil, err
}
