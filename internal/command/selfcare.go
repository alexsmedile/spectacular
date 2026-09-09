package command

import (
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/alexsmedile/spectacular/v2/internal/selfupdate"
)

// Diagnosis is what `doctor` reports.
type Diagnosis struct {
	SchemaVersion string                 `json:"schema_version"`
	Installed     []selfupdate.Component `json:"installed"`
	Latest        string                 `json:"latest,omitempty"`
	Stale         []string               `json:"stale,omitempty"`
	Notes         []string               `json:"notes,omitempty"`
	Remedies      []string               `json:"remedies,omitempty"`
	Healthy       bool                   `json:"healthy"`
}

// UpdateResult is what `update` reports.
type UpdateResult struct {
	SchemaVersion string   `json:"schema_version"`
	Previous      string   `json:"previous_version,omitempty"`
	Installed     string   `json:"installed_version,omitempty"`
	Backup        string   `json:"backup_path,omitempty"`
	Applied       bool     `json:"applied"`
	Latest        string   `json:"latest,omitempty"`
	Stale         []string `json:"stale,omitempty"`
	Notes         []string `json:"notes,omitempty"`
	Remedies      []string `json:"remedies,omitempty"`
}

// installPrefix is the directory an installation was rooted at, derived from
// the running binary: <prefix>/bin/spectacular means <prefix> holds the
// receipt. It is used only to locate that receipt, never to write anything.
func installPrefix() string {
	target, err := selfupdate.TargetPath()
	if err != nil {
		return ""
	}
	return filepath.Dir(filepath.Dir(target))
}

// assess gathers what is installed and, when reachable, the newest release.
// A network failure is a note rather than a refusal: the local half of the
// report is the half a broken machine needs most, and it stays accurate offline.
func assess(ctx context.Context) (selfupdate.Assessment, bool) {
	home, _ := os.UserHomeDir()
	components := selfupdate.Detect(home, installPrefix())
	// A partial failure still carries a version: the release resolved but has
	// no archive for this platform. Reporting that version is what lets the
	// drift report stay useful while the download half is unavailable.
	release, err := selfupdate.LatestRelease(ctx, selfupdate.HTTPClient())
	assessment := selfupdate.Assess(components, release.Version)
	if err != nil {
		assessment.Notes = append(assessment.Notes, "release check: "+err.Error())
	}
	return assessment, err == nil && release.Archive != ""
}

// remedies lists the commands that would bring stale components current.
//
// The binary remedy is offered only when a downloadable archive actually
// exists. Three releases published no assets when the release workflow broke,
// and advising `update -y` against one of those sends the reader to a command
// that cannot succeed.
func remedies(assessment selfupdate.Assessment, binaryHandled, downloadable bool) []string {
	var lines []string
	if assessment.BinaryStale() && !binaryHandled && downloadable {
		lines = append(lines, "spectacular update -y   # replace the binary")
	}
	for _, host := range assessment.StalePlugins() {
		lines = append(lines, selfupdate.HostCommands(host)...)
	}
	return lines
}

func (r Runner) runDoctor(jsonMode bool, invoked string) int {
	assessment, downloadable := assess(context.Background())
	diagnosis := Diagnosis{
		SchemaVersion: "spectacular.doctor.v1",
		Installed:     assessment.Components,
		Latest:        assessment.Latest,
		Stale:         assessment.Stale,
		Notes:         assessment.Notes,
		Remedies:      remedies(assessment, false, downloadable),
		Healthy:       len(assessment.Stale) == 0 && len(assessment.Notes) == 0,
	}
	if jsonMode {
		if err := writeJSON(r.Stdout, diagnosis); err != nil {
			return r.refuse(true, invoked, err)
		}
		return 0
	}
	renderDiagnosis(r.Stdout, diagnosis)
	return 0
}

func (r Runner) runUpdate(jsonMode bool, invoked string, apply bool) int {
	ctx := context.Background()
	assessment, downloadable := assess(ctx)
	result := UpdateResult{
		SchemaVersion: "spectacular.update.v1",
		Latest:        assessment.Latest,
		Stale:         assessment.Stale,
		Notes:         assessment.Notes,
		Previous:      installedVersion(assessment),
	}
	if !apply {
		result.Remedies = remedies(assessment, false, downloadable)
		return r.emitUpdate(jsonMode, invoked, result)
	}
	// -y covers the binary only. A plugin lives inside a host's own state tree,
	// and installing one behind that host's back leaves its records disagreeing
	// with what is on disk, so those stay a printed instruction.
	result.Remedies = remedies(assessment, true, downloadable)
	if !assessment.BinaryStale() {
		return r.emitUpdate(jsonMode, invoked, result)
	}
	target, err := selfupdate.TargetPath()
	if err != nil {
		return r.refuse(jsonMode, invoked, err)
	}
	if err := selfupdate.Writable(target); err != nil {
		return r.refuse(jsonMode, invoked, err)
	}
	release, err := selfupdate.LatestRelease(ctx, selfupdate.HTTPClient())
	if err != nil {
		return r.refuse(jsonMode, invoked, err)
	}
	binary, err := selfupdate.DownloadBinary(ctx, selfupdate.HTTPClient(), release)
	if err != nil {
		return r.refuse(jsonMode, invoked, err)
	}
	backup, err := selfupdate.ReplaceBinary(target, binary)
	if err != nil {
		return r.refuse(jsonMode, invoked, err)
	}
	result.Applied = true
	result.Installed = release.Version
	result.Backup = backup
	return r.emitUpdate(jsonMode, invoked, result)
}

func installedVersion(assessment selfupdate.Assessment) string {
	for _, component := range assessment.Components {
		if component.Kind == "binary" {
			return component.Version
		}
	}
	return ""
}

func (r Runner) emitUpdate(jsonMode bool, invoked string, result UpdateResult) int {
	if jsonMode {
		if err := writeJSON(r.Stdout, result); err != nil {
			return r.refuse(true, invoked, err)
		}
		return 0
	}
	renderUpdate(r.Stdout, result)
	return 0
}

func renderDiagnosis(w io.Writer, diagnosis Diagnosis) {
	fmt.Fprintln(w, "Installed:")
	for _, component := range diagnosis.Installed {
		renderComponent(w, component, diagnosis.Latest)
	}
	if diagnosis.Latest != "" {
		fmt.Fprintf(w, "\nLatest release: %s\n", diagnosis.Latest)
	}
	renderNotes(w, diagnosis.Notes)
	renderRemedies(w, diagnosis.Remedies)
	if diagnosis.Healthy {
		fmt.Fprintln(w, "\nNo problems found.")
	}
}

func renderComponent(w io.Writer, component selfupdate.Component, latest string) {
	status := "not installed"
	if component.Installed {
		status = component.Version
		if selfupdate.Outdated(component.Version, latest) {
			status += "  (outdated)"
		}
	}
	fmt.Fprintf(w, "  %-10s %s\n", component.Name, status)
}

func renderNotes(w io.Writer, notes []string) {
	if len(notes) == 0 {
		return
	}
	fmt.Fprintln(w, "\nNotes:")
	for _, note := range notes {
		fmt.Fprintf(w, "  %s\n", note)
	}
}

func renderRemedies(w io.Writer, lines []string) {
	if len(lines) == 0 {
		return
	}
	fmt.Fprintln(w, "\nTo update, run:")
	for _, line := range lines {
		fmt.Fprintf(w, "  %s\n", line)
	}
}

func renderUpdate(w io.Writer, result UpdateResult) {
	if result.Applied {
		fmt.Fprintf(w, "Updated the binary: %s -> %s\n", result.Previous, result.Installed)
		fmt.Fprintf(w, "Previous binary kept at %s\n", result.Backup)
	} else if result.Latest != "" && len(result.Stale) == 0 {
		fmt.Fprintf(w, "Everything is on %s.\n", result.Latest)
	} else if len(result.Stale) > 0 {
		fmt.Fprintf(w, "Behind release %s: %s\n", result.Latest, strings.Join(result.Stale, ", "))
	}
	renderNotes(w, result.Notes)
	renderRemedies(w, result.Remedies)
}
