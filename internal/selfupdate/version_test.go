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

func TestOutdatedIgnoresNonVersions(t *testing.T) {
	t.Parallel()

	// A binary built from source reports "development". Parsing that as 0.0.0
	// would make every developer's own build look infinitely stale and prompt
	// them to replace it with a release.
	if Outdated("development", "2.18.0") {
		t.Error("a development build is not outdated")
	}
	if Outdated("2.17.0", "unknown") {
		t.Error("an unparseable latest version cannot make anything outdated")
	}
	// A real 0.0.x version is still a version and must compare normally.
	if !Outdated("0.0.1", "0.0.2") {
		t.Error("0.0.1 is older than 0.0.2")
	}
}
