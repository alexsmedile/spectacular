// Package buildinfo owns process-level release metadata. It deliberately has
// no product semantics and is populated by the deterministic release build.
package buildinfo

var (
	Version = "development"
	Commit  = "unknown"
)
