package selfupdate

// HostCommands are the shell commands that update a host's plugin. Each host
// differs, which is the reason detection is worth doing: Codex has no
// `plugin update`, and its marketplace upgrade moves the snapshot without
// installing anything, so the `add` is what actually lands the new version.
// Telling every user the same command would leave half of them stale.
func HostCommands(host Host) []string {
	switch host {
	case HostClaude:
		return []string{
			"claude plugin marketplace update spectacular",
			"claude plugin update spectacular@spectacular",
		}
	case HostCodex:
		return []string{
			"codex plugin marketplace upgrade spectacular",
			"codex plugin add spectacular@spectacular",
		}
	}
	return nil
}

// Assessment is what an audit found: the components, the newest release, and
// whether anything is behind.
type Assessment struct {
	Components []Component `json:"components"`
	Latest     string      `json:"latest,omitempty"`
	// Stale names the components older than Latest.
	Stale []string `json:"stale,omitempty"`
	// Notes carry problems that are not staleness: an unreachable release API,
	// a receipt that disagrees with the binary, a host CLI that is absent.
	Notes []string `json:"notes,omitempty"`
}

// Assess compares what is installed against the newest release. A missing
// latest version is not fatal: an offline machine still gets a report of what
// it has, which is most of what a diagnosis needs.
func Assess(components []Component, latest string) Assessment {
	assessment := Assessment{Components: components, Latest: latest}
	for _, component := range components {
		if component.Detail != "" {
			assessment.Notes = append(assessment.Notes, component.Name+": "+component.Detail)
		}
		if component.Installed && Outdated(component.Version, latest) {
			assessment.Stale = append(assessment.Stale, component.Name)
		}
	}
	return assessment
}

// BinaryStale reports whether the installed binary is behind the release.
func (a Assessment) BinaryStale() bool {
	for _, name := range a.Stale {
		if name == "binary" {
			return true
		}
	}
	return false
}

// StalePlugins names the plugin hosts that are behind, in report order.
func (a Assessment) StalePlugins() []Host {
	var hosts []Host
	for _, component := range a.Components {
		if component.Kind != "plugin" || !component.Installed {
			continue
		}
		if Outdated(component.Version, a.Latest) {
			hosts = append(hosts, Host(component.Name))
		}
	}
	return hosts
}
