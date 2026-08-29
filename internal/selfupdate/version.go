package selfupdate

import (
	"strconv"
	"strings"
)

// CompareVersions orders two semantic versions, returning -1, 0, or 1. A
// leading "v" is tolerated because tags carry it and receipts do not.
//
// A pre-release sorts before the release it precedes, so 2.0.0-rc.2 is older
// than 2.0.0. Getting that backwards would make an update command report a
// release candidate as newer than the release that superseded it.
func CompareVersions(left, right string) int {
	leftCore, leftPre := splitVersion(left)
	rightCore, rightPre := splitVersion(right)

	for i := 0; i < 3; i++ {
		if leftCore[i] != rightCore[i] {
			if leftCore[i] < rightCore[i] {
				return -1
			}
			return 1
		}
	}
	switch {
	case leftPre == "" && rightPre == "":
		return 0
	case leftPre == "":
		return 1 // a release outranks its own pre-release
	case rightPre == "":
		return -1
	case leftPre < rightPre:
		return -1
	case leftPre > rightPre:
		return 1
	}
	return 0
}

// splitVersion separates the numeric core from any pre-release suffix. A
// component that does not parse counts as 0, so an unexpected directory name in
// a plugin cache sorts low instead of panicking a read-only report.
func splitVersion(version string) ([3]int, string) {
	version = strings.TrimSpace(version)
	version = strings.TrimPrefix(version, "v")
	core, pre, _ := strings.Cut(version, "-")
	var parts [3]int
	for i, field := range strings.SplitN(core, ".", 3) {
		if i > 2 {
			break
		}
		parts[i], _ = strconv.Atoi(field)
	}
	return parts, pre
}

// Outdated reports whether installed is older than latest. An empty installed
// version is not outdated: nothing is installed, which is a different state
// that the caller reports differently.
func Outdated(installed, latest string) bool {
	if installed == "" || latest == "" {
		return false
	}
	return CompareVersions(installed, latest) < 0
}
