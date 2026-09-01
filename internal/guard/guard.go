package guard

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/alexsmedile/spectacular/v2/internal/charter"
	"github.com/alexsmedile/spectacular/v2/internal/discovery"
	"github.com/alexsmedile/spectacular/v2/internal/domain"
)

const SchemaVersion = "spectacular.guard.v2"

// GuardResult represents the outcome of a guarded execution.
type GuardResult struct {
	SchemaVersion  string   `json:"schema_version"`
	MissionRef     string   `json:"mission_ref"`
	ObjectiveRef   string   `json:"objective_ref"`
	Mode           string   `json:"mode"`
	Status         string   `json:"status"` // "pass", "violation", "killed"
	ExitCode       int      `json:"exit_code"`
	AllowedPaths   []string `json:"allowed_paths"`
	PreservedPaths []string `json:"preserved_paths,omitempty"`
	EscapedPaths   []string `json:"escaped_paths,omitempty"`
	RolledBack     bool     `json:"rolled_back"`
	AutoHealed     bool     `json:"auto_healed"`
	FeedbackPrompt string   `json:"feedback_prompt,omitempty"`
	Output         string   `json:"output,omitempty"`
}

// Run executes a command under perimeter supervision (post-flight watchdog or real-time watcher).
func Run(ws *discovery.Workspace, targetRef string, watchMode bool, cmdArgs []string) (*GuardResult, error) {
	if ws == nil {
		return nil, errors.New("guard: workspace is nil")
	}
	if len(cmdArgs) == 0 {
		return nil, errors.New("guard: command arguments required after --")
	}

	parts := strings.Split(targetRef, "/")
	if len(parts) != 2 {
		return nil, domain.NewRefusal(domain.RefusalInvalidReference, targetRef, "expected <mission-ref>/<objective-ref> (e.g. M17/O1)", nil)
	}

	c, err := charter.Compile(ws, parts[0], parts[1], nil)
	if err != nil {
		return nil, fmt.Errorf("guard: compile charter: %w", err)
	}

	allowed := c.Layer3.WritesPaths
	testCmd := c.Layer3.VerificationCommand
	root := ws.Root

	preSnapshot, err := takeSnapshot(root)
	if err != nil {
		return nil, fmt.Errorf("guard: pre-flight snapshot: %w", err)
	}

	cmd := exec.Command(cmdArgs[0], cmdArgs[1:]...)
	cmd.Dir = root

	var combinedBuf bytes.Buffer
	cmd.Stdout = &combinedBuf
	cmd.Stderr = &combinedBuf

	mode := "watchdog"
	if watchMode {
		mode = "realtime-watcher"
	}

	res := &GuardResult{
		SchemaVersion: SchemaVersion,
		MissionRef:    parts[0],
		ObjectiveRef:  parts[1],
		Mode:          mode,
		AllowedPaths:  allowed,
	}

	if watchMode {
		// Real-time watcher mode
		if err := cmd.Start(); err != nil {
			return nil, fmt.Errorf("guard: start command: %w", err)
		}

		stopWatch := make(chan struct{})
		var watchMu sync.Mutex
		var violationPaths []string

		go func() {
			ticker := time.NewTicker(50 * time.Millisecond)
			defer ticker.Stop()
			for {
				select {
				case <-stopWatch:
					return
				case <-ticker.C:
					currSnapshot, _ := takeSnapshot(root)
					escaped, _ := diffPaths(preSnapshot, currSnapshot, allowed)
					if len(escaped) > 0 {
						watchMu.Lock()
						violationPaths = escaped
						watchMu.Unlock()
						_ = cmd.Process.Kill()
						return
					}
				}
			}
		}()

		waitErr := cmd.Wait()
		close(stopWatch)

		watchMu.Lock()
		escapedFound := violationPaths
		watchMu.Unlock()

		res.Output = combinedBuf.String()

		if len(escapedFound) > 0 {
			currSnapshot, _ := takeSnapshot(root)
			_, preserved := diffPaths(preSnapshot, currSnapshot, allowed)

			res.Status = "killed"
			res.EscapedPaths = escapedFound
			res.PreservedPaths = preserved
			res.RolledBack = true
			res.AutoHealed = true
			res.ExitCode = 2
			res.FeedbackPrompt = buildFeedbackPrompt(parts[0], parts[1], preserved, escapedFound, allowed, testCmd)
			surgicalQuarantine(root, preSnapshot, escapedFound)
			return res, nil
		}

		if waitErr != nil {
			if exitErr, ok := waitErr.(*exec.ExitError); ok {
				res.ExitCode = exitErr.ExitCode()
			} else {
				res.ExitCode = 1
			}
		} else {
			res.ExitCode = 0
		}
		res.Status = "pass"
		return res, nil
	}

	// Post-flight watchdog mode (default)
	waitErr := cmd.Run()
	res.Output = combinedBuf.String()
	if waitErr != nil {
		if exitErr, ok := waitErr.(*exec.ExitError); ok {
			res.ExitCode = exitErr.ExitCode()
		} else {
			res.ExitCode = 1
		}
	} else {
		res.ExitCode = 0
	}

	postSnapshot, err := takeSnapshot(root)
	if err != nil {
		return nil, fmt.Errorf("guard: post-flight snapshot: %w", err)
	}

	escaped, preserved := diffPaths(preSnapshot, postSnapshot, allowed)
	if len(escaped) > 0 {
		res.Status = "violation"
		res.EscapedPaths = escaped
		res.PreservedPaths = preserved
		res.RolledBack = true
		res.AutoHealed = true
		res.ExitCode = 2
		res.FeedbackPrompt = buildFeedbackPrompt(parts[0], parts[1], preserved, escaped, allowed, testCmd)
		surgicalQuarantine(root, preSnapshot, escaped)
		return res, nil
	}

	res.Status = "pass"
	res.PreservedPaths = preserved
	return res, nil
}

type snapshot map[string]string

func takeSnapshot(root string) (snapshot, error) {
	snap := make(snapshot)
	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return nil
		}
		rel, relErr := filepath.Rel(root, path)
		if relErr != nil || rel == "." {
			return nil
		}
		// Skip .git
		if rel == ".git" || strings.HasPrefix(rel, ".git"+string(filepath.Separator)) {
			if info.IsDir() {
				return filepath.SkipDir
			}
			return nil
		}

		if !info.IsDir() {
			hash, hashErr := fileHash(path)
			if hashErr == nil {
				snap[filepath.ToSlash(rel)] = hash
			}
		}
		return nil
	})
	return snap, err
}

func fileHash(path string) (string, error) {
	f, err := os.Open(path)
	if err != nil {
		return "", err
	}
	defer f.Close()
	h := sha256.New()
	if _, err := io.Copy(h, f); err != nil {
		return "", err
	}
	return hex.EncodeToString(h.Sum(nil)), nil
}

func diffPaths(pre, post snapshot, allowed []string) (escaped []string, preserved []string) {
	if len(allowed) == 0 {
		return nil, nil
	}

	seen := make(map[string]bool)

	for path, postHash := range post {
		preHash, exists := pre[path]
		if !exists || preHash != postHash {
			// Path was added or modified
			seen[path] = true
			if matchesAny(path, allowed) {
				preserved = append(preserved, path)
			} else {
				escaped = append(escaped, path)
			}
		}
	}

	// Check deletions
	for path := range pre {
		if _, exists := post[path]; !exists {
			if !seen[path] {
				if matchesAny(path, allowed) {
					preserved = append(preserved, path)
				} else {
					escaped = append(escaped, path)
				}
			}
		}
	}

	return escaped, preserved
}

func matchesAny(path string, patterns []string) bool {
	for _, pat := range patterns {
		normPat := filepath.ToSlash(pat)
		normPath := filepath.ToSlash(path)

		if strings.HasSuffix(normPat, "/**") {
			prefix := strings.TrimSuffix(normPat, "/**")
			if strings.HasPrefix(normPath, prefix+"/") || normPath == prefix {
				return true
			}
		} else if strings.HasSuffix(normPat, "/*") {
			prefix := strings.TrimSuffix(normPat, "/*")
			if filepath.Dir(normPath) == prefix {
				return true
			}
		} else if strings.HasPrefix(normPat, "*.") {
			ext := strings.TrimPrefix(normPat, "*")
			if strings.HasSuffix(normPath, ext) {
				return true
			}
		} else {
			if normPath == normPat || strings.HasPrefix(normPath, normPat+"/") {
				return true
			}
			matched, _ := filepath.Match(normPat, normPath)
			if matched {
				return true
			}
		}
	}
	return false
}

// surgicalQuarantine deletes ONLY escaped paths while preserving valid work in authorized paths.
func surgicalQuarantine(root string, pre snapshot, escaped []string) {
	for _, path := range escaped {
		full := filepath.Join(root, filepath.FromSlash(path))
		if _, existed := pre[path]; !existed {
			// Rogue new file: remove it
			_ = os.Remove(full)
		}
	}
}

func buildFeedbackPrompt(missionRef, objRef string, preserved, escaped []string, allowed []string, testCmd string) string {
	var b strings.Builder
	b.WriteString("[SPECTACULAR PERIMETER NOTICE]\n")
	if len(preserved) > 0 {
		b.WriteString(fmt.Sprintf("- Your valid changes in [%s] were preserved.\n", strings.Join(preserved, ", ")))
	}
	if len(escaped) > 0 {
		b.WriteString(fmt.Sprintf("- The following unauthorized write paths were quarantined and removed: [%s]\n", strings.Join(escaped, ", ")))
	}
	b.WriteString(fmt.Sprintf("- Authorized write paths for %s/%s are strictly: %s\n", missionRef, objRef, strings.Join(allowed, ", ")))
	if testCmd != "" {
		b.WriteString(fmt.Sprintf("- Please complete the objective within authorized paths, run `%s`, and finish with STATUS: DONE.\n", testCmd))
	} else {
		b.WriteString("- Please complete the objective within authorized paths and finish with STATUS: DONE.\n")
	}
	return b.String()
}
