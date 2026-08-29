package selfupdate

import "testing"

func TestCompareVersionsOrdersReleasesAndPreReleases(t *testing.T) {
	t.Parallel()

	tests := []struct {
		left, right string
		want        int
	}{
		{"2.7.2", "2.7.2", 0},
		{"2.7.1", "2.7.2", -1},
		{"2.7.2", "2.7.1", 1},
		{"2.6.0", "2.7.0", -1},
		{"1.34.0", "2.0.0", -1},
		// 34 > 8 numerically, but a string sort would put 1.34.0 before 1.8.3.
		// The plugin cache held exactly these two, so the wrong comparison would
		// report a stale version as the newest installed.
		{"1.34.0", "1.8.3", 1},
		// A pre-release precedes the release it leads to.
		{"2.0.0-rc.2", "2.0.0", -1},
		{"2.0.0", "2.0.0-rc.2", 1},
		{"2.0.0-rc.1", "2.0.0-rc.2", -1},
		// A leading v is tolerated: tags carry it, receipts do not.
		{"v2.7.2", "2.7.2", 0},
		{"v2.7.1", "2.7.2", -1},
	}
	for _, test := range tests {
		if got := CompareVersions(test.left, test.right); got != test.want {
			t.Errorf("CompareVersions(%q, %q) = %d, want %d", test.left, test.right, got, test.want)
		}
	}
}

func TestOutdatedTreatsAbsenceAsNotOutdated(t *testing.T) {
	t.Parallel()

	if !Outdated("2.7.1", "2.7.2") {
		t.Error("an older installed version is outdated")
	}
	if Outdated("2.7.2", "2.7.2") {
		t.Error("a current version is not outdated")
	}
	if Outdated("2.7.3", "2.7.2") {
		t.Error("a newer installed version is not outdated")
	}
	// Nothing installed is a different state from something old, and a caller
	// that conflated them would offer to update a component that is absent.
	if Outdated("", "2.7.2") {
		t.Error("an absent component is not outdated")
	}
	if Outdated("2.7.1", "") {
		t.Error("an unknown latest version cannot make anything outdated")
	}
}
