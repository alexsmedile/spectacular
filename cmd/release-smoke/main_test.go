package main

import "testing"

func TestValidV2VersionAcceptsStableAndPrerelease(t *testing.T) {
	for _, version := range []string{"spectacular 2.1.0", "spectacular 2.0.0-rc.2"} {
		if !validV2Version(version) {
			t.Fatalf("rejected %q", version)
		}
	}
	for _, version := range []string{"spectacular 1.9.0", "spectacular 2.bad.0", "spectacular 2.1"} {
		if validV2Version(version) {
			t.Fatalf("accepted %q", version)
		}
	}
}
