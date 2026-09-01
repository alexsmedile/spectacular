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
	"syscall"
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
func Run(ws *discovery.Workspace, targetRef string, watchMode bool, execCmd string, cmdArgs []string) (*GuardResult, error) {
	if ws == nil {
		return nil, errors.New("guard: workspace is nil")
	}

	parts := strings.Split(targetRef, "/")
	if len(parts) != 2 {
		return nil, domain.NewRefusal(domain.RefusalInvalidReference, targetRef, "expected <mission-ref>/<objective-ref> (e.g. M17/O1)", nil)
	}

	c, err := charter.Compile(ws, parts[0], parts[1], nil)
	if err != nil {
		return nil, fmt.Errorf("guard: compile charter: %w", err)
	}

	if execCmd != "" && len(cmdArgs) > 0 {
		return nil, domain.NewRefusal(domain.RefusalInvalidReference, targetRef, "cannot combine --exec with command arguments after --", nil)
	}

	var cmd *exec.Cmd
	root := ws.Root

	if execCmd != "" {
		prompt := c.RenderPrompt()
		// Use shell execution to natively preserve quotes and argument boundaries
		cmd = exec.Command("sh", "-c", execCmd+` "$@"`, "spectacular-worker", prompt)
	} else if len(cmdArgs) > 0 {
		cmd = exec.Command(cmdArgs[0], cmdArgs[1:]...)
	} else {
		return nil, errors.New("guard: command arguments required after -- or via --exec")
	}

	cmd.Dir = root
	cmd.SysProcAttr = &syscall.SysProcAttr{Setpgid: true}

	allowed := c.Layer3.WritesPaths
	testCmd := c.Layer3.VerificationCommand

	preSnapshot, err := takeSnapshot(root, allowed)
	if err != nil {
		return nil, fmt.Errorf("guard: pre-flight snapshot: %w", err)
	}

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
			ticker := time.NewTicker(20 * time.Millisecond)
			defer ticker.Stop()
			for {
				select {
				case <-stopWatch:
					return
				case <-ticker.C:
					currSnapshot, _ := takeSnapshot(root, allowed)
					escaped, _ := diffPaths(root, preSnapshot, currSnapshot, allowed)
					if len(escaped) > 0 {
						watchMu.Lock()
						violationPaths = escaped
						watchMu.Unlock()
						killProcessGroup(cmd)
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

		// Always run post-flight diff to prevent short-lived race conditions
		postSnapshot, _ := takeSnapshot(root, allowed)
		finalEscaped, finalPreserved := diffPaths(root, preSnapshot, postSnapshot, allowed)
		if len(finalEscaped) > 0 {
			escapedFound = finalEscaped
		}

		if len(escapedFound) > 0 {
			res.Status = "violation"
			if len(violationPaths) > 0 {
				res.Status = "killed"
			}
			res.EscapedPaths = escapedFound
			res.PreservedPaths = finalPreserved
			res.RolledBack = true
			res.AutoHealed = true
			res.ExitCode = 2
			res.FeedbackPrompt = buildFeedbackPrompt(parts[0], parts[1], finalPreserved, escapedFound, allowed, testCmd)
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
		res.PreservedPaths = finalPreserved
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

	postSnapshot, err := takeSnapshot(root, allowed)
	if err != nil {
		return nil, fmt.Errorf("guard: post-flight snapshot: %w", err)
	}

	escaped, preserved := diffPaths(root, preSnapshot, postSnapshot, allowed)
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

func killProcessGroup(cmd *exec.Cmd) {
	if cmd == nil || cmd.Process == nil {
		return
	}
	pgid, err := syscall.Getpgid(cmd.Process.Pid)
	if err == nil {
		_ = syscall.Kill(-pgid, syscall.SIGKILL)
	} else {
		_ = cmd.Process.Kill()
	}
}

type fileMeta struct {
	modTime int64
	size    int64
	hash    string
	content []byte
}

type snapshot map[string]fileMeta

func takeSnapshot(root string, allowed []string) (snapshot, error) {
	snap := make(snapshot)
	err := filepath.WalkDir(root, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return nil
		}
		rel, relErr := filepath.Rel(root, path)
		if relErr != nil || rel == "." {
			return nil
		}
		// Skip non-project directories & directory symlinks
		if d.IsDir() {
			name := d.Name()
			if strings.HasPrefix(name, "_") || name == ".git" || name == "node_modules" || name == "vendor" || name == ".cache" || name == "__pycache__" || name == "target" || name == "dist" || name == ".next" || (strings.HasPrefix(rel, ".spectacular"+string(filepath.Separator)) && name == "raw") {
				return filepath.SkipDir
			}
			return nil
		}
		if d.Type()&os.ModeSymlink != 0 {
			info, sErr := os.Stat(path)
			if sErr == nil && info.IsDir() {
				return filepath.SkipDir
			}
		}

		info, infoErr := d.Info()
		if infoErr == nil {
			meta := fileMeta{
				modTime: info.ModTime().UnixNano(),
				size:    info.Size(),
			}
			relSlash := filepath.ToSlash(rel)
			// Compute hash & cache content for files outside allowed paths for cryptographic assurance
			if !matchesAny(relSlash, allowed) {
				if h, hErr := computeFileHash(path); hErr == nil {
					meta.hash = h
				}
				if info.Size() < 1024*1024 {
					if data, readErr := os.ReadFile(path); readErr == nil {
						meta.content = data
					}
				}
			}
			snap[relSlash] = meta
		}
		return nil
	})
	return snap, err
}

func computeFileHash(path string) (string, error) {
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

func diffPaths(root string, pre, post snapshot, allowed []string) (escaped []string, preserved []string) {
	if len(allowed) == 0 {
		return nil, nil
	}

	seen := make(map[string]bool)

	for path, postMeta := range post {
		preMeta, exists := pre[path]
		isAllowed := matchesAny(path, allowed)

		isModified := false
		if !exists {
			isModified = true
		} else if isAllowed {
			// Inside allowed: fast mtime+size check
			if preMeta.modTime != postMeta.modTime || preMeta.size != postMeta.size {
				isModified = true
			}
		} else {
			// Outside allowed: if metadata changed, confirm with cryptographic hash
			if preMeta.modTime != postMeta.modTime || preMeta.size != postMeta.size {
				full := filepath.Join(root, filepath.FromSlash(path))
				currentHash, _ := computeFileHash(full)
				if preMeta.hash == "" || preMeta.hash != currentHash {
					isModified = true
				}
			}
		}

		if isModified {
			seen[path] = true
			if isAllowed {
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

// surgicalQuarantine restores or deletes escaped paths while preserving valid work in authorized paths.
func surgicalQuarantine(root string, pre snapshot, escaped []string) {
	for _, path := range escaped {
		full := filepath.Join(root, filepath.FromSlash(path))
		meta, existed := pre[path]
		if !existed {
			// Rogue new file: remove it
			_ = os.Remove(full)
		} else if len(meta.content) > 0 {
			// Pre-existing file modified: restore exact pre-execution content
			_ = os.WriteFile(full, meta.content, 0o644)
		} else {
			// Fallback: restore via git checkout
			cmd := exec.Command("git", "checkout", "--", filepath.FromSlash(path))
			cmd.Dir = root
			_ = cmd.Run()
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
