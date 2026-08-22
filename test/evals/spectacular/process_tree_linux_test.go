//go:build linux

package spectaculareval

import (
	"errors"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

func processGroupHasLiveMember(groupID int) (bool, error) {
	entries, err := os.ReadDir("/proc")
	if err != nil {
		return false, err
	}
	for _, entry := range entries {
		if _, err := strconv.Atoi(entry.Name()); err != nil {
			continue
		}
		data, err := os.ReadFile(filepath.Join("/proc", entry.Name(), "stat"))
		if errors.Is(err, os.ErrNotExist) || errors.Is(err, os.ErrPermission) {
			continue
		}
		if err != nil {
			return false, err
		}
		stat := string(data)
		closingParen := strings.LastIndexByte(stat, ')')
		if closingParen < 0 {
			continue
		}
		fields := strings.Fields(stat[closingParen+1:])
		if len(fields) < 3 {
			continue
		}
		processGroup, err := strconv.Atoi(fields[2])
		if err != nil || processGroup != groupID {
			continue
		}
		if fields[0] != "Z" && fields[0] != "X" {
			return true, nil
		}
	}
	return false, nil
}
