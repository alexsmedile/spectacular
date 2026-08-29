// Package selfupdate reports and advances the versions of an installation: the
// native binary and every agent host that carries the Spectacular plugin.
//
// The package deliberately separates detection from action. Detection is
// read-only, works with no network, and never shells out to a host CLI, so
// `doctor` can report a machine that is offline or whose host tools are broken.
// Action builds on the same detection, so the two commands can never disagree
// about what is installed.
package selfupdate

import (
	"bufio"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

// Host names an agent runtime that can carry the Spectacular plugin.
type Host string

const (
	HostClaude Host = "claude"
	HostCodex  Host = "codex"
)

// Component is one updatable part of an installation. Absent components are
// reported rather than omitted: "not installed" and "unknown" are different
// answers, and a reader deciding whether to act needs to tell them apart.
type Component struct {
	Name string `json:"name"`
	// Kind is "binary" or "plugin".
	Kind string `json:"kind"`
	// Version is the installed version, empty when nothing is installed.
	Version string `json:"version,omitempty"`
	// Path is where the component was found, for a reader who wants to look.
	Path string `json:"path,omitempty"`
	// Installed distinguishes a component that is absent from one whose
	// version could not be read.
	Installed bool `json:"installed"`
	// Detail explains an absent or unreadable component.
	Detail string `json:"detail,omitempty"`
}

// Receipt is the install receipt written beside a binary installation. It
// records what was installed and how, so an update can reuse the same prefix,
// runtime, and platform rather than guessing them.
type Receipt struct {
	Version  string
	Runtime  string
	Platform string
	Prefix   string
}

// ReadReceipt parses the install receipt under prefix. A missing receipt is not
// an error: a binary can be present without one, and the caller decides whether
// that is a problem.
func ReadReceipt(prefix string) (Receipt, bool) {
	path := filepath.Join(prefix, "share", "spectacular", "install.receipt")
	file, err := os.Open(path)
	if err != nil {
		return Receipt{}, false
	}
	defer file.Close()

	receipt := Receipt{Prefix: prefix}
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		key, value, found := strings.Cut(scanner.Text(), "=")
		if !found {
			continue
		}
		switch strings.TrimSpace(key) {
		case "version":
			receipt.Version = strings.TrimSpace(value)
		case "runtime":
			receipt.Runtime = strings.TrimSpace(value)
		case "platform":
			receipt.Platform = strings.TrimSpace(value)
		}
	}
	if scanner.Err() != nil || receipt.Version == "" {
		return Receipt{}, false
	}
	return receipt, true
}

// DetectBinary reports the installed native binary. The receipt is preferred
// over running the binary: reading a file cannot hang, and a binary that no
// longer executes is exactly the state a doctor should still be able to report.
func DetectBinary(prefix string) Component {
	component := Component{Name: "binary", Kind: "binary"}
	binary := filepath.Join(prefix, "bin", "spectacular")
	if _, err := os.Stat(binary); err != nil {
		component.Detail = "no binary at " + binary
		return component
	}
	component.Path = binary
	component.Installed = true
	if receipt, ok := ReadReceipt(prefix); ok {
		component.Version = receipt.Version
		return component
	}
	component.Detail = "installed, but no readable install receipt"
	return component
}

// pluginRoots are the per-host cache directories a plugin install lands in.
// Each holds one directory per installed version.
func pluginRoots(home string) map[Host]string {
	return map[Host]string{
		HostClaude: filepath.Join(home, ".claude", "plugins", "cache", "spectacular", "spectacular"),
		HostCodex:  filepath.Join(home, ".codex", "plugins", "cache", "spectacular", "spectacular"),
	}
}

// DetectPlugin reports the plugin installed for one host. A host whose CLI is
// absent is reported as not installed rather than skipped, so a reader sees the
// whole picture and not just the parts that happened to work.
func DetectPlugin(home string, host Host) Component {
	component := Component{Name: string(host), Kind: "plugin"}
	root, ok := pluginRoots(home)[host]
	if !ok {
		component.Detail = "unknown host"
		return component
	}
	if _, err := exec.LookPath(string(host)); err != nil {
		component.Detail = string(host) + " CLI is not on PATH"
		return component
	}
	entries, err := os.ReadDir(root)
	if err != nil {
		component.Detail = "no plugin installed for " + string(host)
		return component
	}
	// A cache may hold several versions. The newest by semantic order is the
	// one a host resolves, so report that rather than whichever the filesystem
	// happened to list first.
	var versions []string
	for _, entry := range entries {
		if entry.IsDir() {
			versions = append(versions, entry.Name())
		}
	}
	if len(versions) == 0 {
		component.Detail = "no plugin installed for " + string(host)
		return component
	}
	newest := versions[0]
	for _, candidate := range versions[1:] {
		if CompareVersions(candidate, newest) > 0 {
			newest = candidate
		}
	}
	component.Version = newest
	component.Path = filepath.Join(root, newest)
	component.Installed = true
	if len(versions) > 1 {
		component.Detail = "cache holds " + itoa(len(versions)) + " versions; newest reported"
	}
	return component
}

// Detect reports every component of an installation.
func Detect(home, prefix string) []Component {
	components := []Component{DetectBinary(prefix)}
	for _, host := range []Host{HostClaude, HostCodex} {
		components = append(components, DetectPlugin(home, host))
	}
	return components
}

func itoa(value int) string {
	if value == 0 {
		return "0"
	}
	var digits []byte
	for value > 0 {
		digits = append([]byte{byte('0' + value%10)}, digits...)
		value /= 10
	}
	return string(digits)
}
