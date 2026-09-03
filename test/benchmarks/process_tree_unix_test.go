//go:build darwin || linux

package spectaculareval

import (
	"context"
	"errors"
	"os/exec"
	"syscall"
	"testing"
	"time"
)

func TestTrialCancellationKillsAdapterProcessGroup(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 50*time.Millisecond)
	defer cancel()
	command := exec.CommandContext(ctx, "sh", "-c", "sleep 30 & wait")
	configureProcessTree(command)
	command.WaitDelay = 100 * time.Millisecond
	started := time.Now()
	err := command.Run()
	if err == nil || time.Since(started) > 2*time.Second {
		t.Fatalf("process-tree deadline failed: elapsed=%s err=%v", time.Since(started), err)
	}
	if command.Process == nil {
		t.Fatal("adapter process never started")
	}
	alive, checkErr := processGroupHasLiveMember(command.Process.Pid)
	if checkErr != nil {
		if errors.Is(checkErr, syscall.EPERM) {
			t.Skip("sandboxed environment prevents inspecting foreign process groups")
		}
		t.Fatalf("inspect adapter process group: %v", checkErr)
	}
	if alive {
		t.Fatal("adapter process group retained a live member after cancellation")
	}
}
