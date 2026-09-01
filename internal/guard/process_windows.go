//go:build windows

package guard

import (
	"os/exec"
)

func setProcessGroup(cmd *exec.Cmd) {
	// Windows process group handling
}

func killProcessGroup(cmd *exec.Cmd) {
	if cmd == nil || cmd.Process == nil {
		return
	}
	_ = cmd.Process.Kill()
}
